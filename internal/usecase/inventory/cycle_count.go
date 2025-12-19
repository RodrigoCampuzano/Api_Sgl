package inventory

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// PerformCycleCountUseCase implementa HU-15: Conteo cíclico con selección aleatoria
type PerformCycleCountUseCase struct {
	cycleCountRepo domain.CycleCountRepository
	inventoryRepo  domain.InventoryRepository
	movementRepo   domain.InventoryMovementRepository
	productRepo    domain.ProductRepository
	auditRepo      domain.AuditRepository
}

func NewPerformCycleCountUseCase(
	cycleCountRepo domain.CycleCountRepository,
	inventoryRepo domain.InventoryRepository,
	movementRepo domain.InventoryMovementRepository,
	productRepo domain.ProductRepository,
	auditRepo domain.AuditRepository,
) *PerformCycleCountUseCase {
	return &PerformCycleCountUseCase{
		cycleCountRepo: cycleCountRepo,
		inventoryRepo:  inventoryRepo,
		movementRepo:   movementRepo,
		productRepo:    productRepo,
		auditRepo:      auditRepo,
	}
}

// GenerateDailyCycleCounts genera 5 ubicaciones aleatorias para conteo diario (HU-15)
func (uc *PerformCycleCountUseCase) GenerateDailyCycleCounts() ([]*domain.CycleCount, error) {
	// Obtener productos activos (priorizando alto valor - simulado)
	products, err := uc.productRepo.List(nil, 100, 0)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, errors.New("no hay productos disponibles")
	}

	// Seleccionar 5 productos aleatorios
	rand.Seed(time.Now().UnixNano())
	selectedCount := 5
	if len(products) < 5 {
		selectedCount = len(products)
	}

	selectedProducts := make([]*domain.Product, selectedCount)
	perm := rand.Perm(len(products))
	for i := 0; i < selectedCount; i++ {
		selectedProducts[i] = products[perm[i]]
	}

	// Crear conteos cíclicos
	var cycleCounts []*domain.CycleCount
	scheduledDate := time.Now()

	for _, product := range selectedProducts {
		// Obtener stock actual
		stock, _ := uc.inventoryRepo.GetStockByProduct(product.ID)

		count := &domain.CycleCount{
			ScheduledDate:    scheduledDate,
			Location:         "Random", // En producción, seleccionar ubicación específica
			ProductID:        product.ID,
			ExpectedQuantity: &stock,
			Status:           "PENDIENTE",
		}

		if err := uc.cycleCountRepo.Create(count); err != nil {
			continue
		}

		cycleCounts = append(cycleCounts, count)
	}

	return cycleCounts, nil
}

type PerformCountInput struct {
	CountID         uuid.UUID `json:"count_id"`
	CountedQuantity int       `json:"counted_quantity"`
	UserID          uuid.UUID `json:"-"`
}

// PerformCount registra el conteo y ajusta inventario si hay varianza
func (uc *PerformCycleCountUseCase) PerformCount(input PerformCountInput) error {
	// 1. Obtener el conteo cíclico
	count, err := uc.cycleCountRepo.FindByID(input.CountID)
	if err != nil {
		return errors.New("conteo cíclico no encontrado")
	}

	if count.Status != "PENDIENTE" {
		return errors.New("el conteo ya fue realizado")
	}

	// 2. Registrar el conteo
	now := time.Now()
	count.CountedQuantity = &input.CountedQuantity
	count.CountedBy = &input.UserID
	count.CountedAt = &now
	count.Status = "COMPLETADO"

	if err := uc.cycleCountRepo.Update(count); err != nil {
		return err
	}

	// 3. Si hay varianza, crear ajuste de inventario
	if count.Variance != nil && *count.Variance != 0 {
		// Obtener inventarios del producto
		inventories, _ := uc.inventoryRepo.FindByProduct(count.ProductID)

		if len(inventories) > 0 {
			// Ajustar el primer inventario disponible (simplificado)
			inventory := inventories[0]
			previousQty := inventory.Quantity
			inventory.Quantity = input.CountedQuantity

			movement := &domain.InventoryMovement{
				InventoryID:      inventory.ID,
				MovementType:     domain.MovementAjuste,
				Quantity:         *count.Variance,
				PreviousQuantity: previousQty,
				NewQuantity:      inventory.Quantity,
				ReferenceID:      &count.ID,
				ReferenceType:    "CYCLE_COUNT",
				Reason:           "Ajuste por conteo cíclico",
				PerformedBy:      input.UserID,
			}

			_ = uc.movementRepo.Create(movement)
			_ = uc.inventoryRepo.Update(inventory)
		}
	}

	// 4. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "CYCLE_COUNT",
		EntityType: "CYCLE_COUNT",
		EntityID:   &count.ID,
		NewValues: map[string]interface{}{
			"expected": count.ExpectedQuantity,
			"counted":  input.CountedQuantity,
			"variance": count.Variance,
		},
	})

	return nil
}
