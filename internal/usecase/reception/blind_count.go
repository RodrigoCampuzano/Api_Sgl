package reception

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// BlindCountUseCase HU-02: Conteo ciego (sin mostrar cantidad esperada)
type BlindCountUseCase struct {
	receptionOrderRepo domain.ReceptionOrderRepository
	receptionLineRepo  domain.ReceptionLineRepository
	discrepancyRepo    domain.ReceptionDiscrepancyRepository
	auditRepo          domain.AuditRepository
}

func NewBlindCountUseCase(
	receptionOrderRepo domain.ReceptionOrderRepository,
	receptionLineRepo domain.ReceptionLineRepository,
	discrepancyRepo domain.ReceptionDiscrepancyRepository,
	auditRepo domain.AuditRepository,
) *BlindCountUseCase {
	return &BlindCountUseCase{
		receptionOrderRepo: receptionOrderRepo,
		receptionLineRepo:  receptionLineRepo,
		discrepancyRepo:    discrepancyRepo,
		auditRepo:          auditRepo,
	}
}

type BlindCountLineInput struct {
	LineID          uuid.UUID               `json:"line_id"`
	CountedQuantity int                     `json:"counted_quantity"`
	Condition       domain.ProductCondition `json:"condition,omitempty"`
}

type BlindCountInput struct {
	ReceptionOrderID uuid.UUID             `json:"reception_order_id"`
	Lines            []BlindCountLineInput `json:"lines"`
	UserID           uuid.UUID             `json:"-"` // Del contexto (AUXILIAR o RECEPCIONISTA)
}

func (uc *BlindCountUseCase) Execute(input BlindCountInput) error {
	// 1. Verificar que la orden existe y está pendiente
	order, err := uc.receptionOrderRepo.FindByID(input.ReceptionOrderID)
	if err != nil {
		return errors.New("orden de recepción no encontrada")
	}

	if order.Status != domain.ReceptionPendiente {
		return errors.New("la orden ya fue contada o validada")
	}

	// 2. Actualizar las líneas con el conteo
	now := time.Now()
	hasDiscrepancies := false

	for _, countInput := range input.Lines {
		line, err := uc.receptionLineRepo.FindByID(countInput.LineID)
		if err != nil {
			return err
		}

		// Verificar que pertenece a esta orden
		if line.ReceptionOrderID != input.ReceptionOrderID {
			return errors.New("línea no pertenece a esta orden")
		}

		// Actualizar con el conteo
		line.CountedQuantity = &countInput.CountedQuantity
		line.CountedBy = &input.UserID
		line.CountedAt = &now
		if countInput.Condition != "" {
			line.Condition = countInput.Condition
		}

		if err := uc.receptionLineRepo.Update(line); err != nil {
			return err
		}

		// HU-03: Detectar discrepancias automáticamente
		if line.HasDiscrepancy() {
			hasDiscrepancies = true
			discrepancy := &domain.ReceptionDiscrepancy{
				ReceptionLineID: line.ID,
				ExpectedQty:     line.ExpectedQuantity,
				CountedQty:      *line.CountedQuantity,
				Difference:      *line.CountedQuantity - line.ExpectedQuantity,
				Status:          domain.DiscrepancyDetectada,
			}
			_ = uc.discrepancyRepo.Create(discrepancy)
		}
	}

	// 3. Actualizar estado de la orden
	if hasDiscrepancies {
		order.Status = domain.ReceptionConIncidencia
	} else {
		order.Status = domain.ReceptionEnConteo
	}
	order.ReceivedBy = &input.UserID
	receivedAt := time.Now()
	order.ReceivedAt = &receivedAt

	if err := uc.receptionOrderRepo.Update(order); err != nil {
		return err
	}

	// 4. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "BLIND_COUNT",
		EntityType: "RECEPTION_ORDER",
		EntityID:   &order.ID,
		NewValues: map[string]interface{}{
			"lines_counted":     len(input.Lines),
			"has_discrepancies": hasDiscrepancies,
		},
	})

	return nil
}
