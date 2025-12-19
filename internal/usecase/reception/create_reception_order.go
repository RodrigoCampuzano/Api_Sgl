package reception

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// CreateReceptionOrderUseCase HU-01: Alta de órdenes de recepción
type CreateReceptionOrderUseCase struct {
	receptionOrderRepo domain.ReceptionOrderRepository
	receptionLineRepo  domain.ReceptionLineRepository
	supplierRepo       domain.SupplierRepository
	productRepo        domain.ProductRepository
	auditRepo          domain.AuditRepository
}

func NewCreateReceptionOrderUseCase(
	receptionOrderRepo domain.ReceptionOrderRepository,
	receptionLineRepo domain.ReceptionLineRepository,
	supplierRepo domain.SupplierRepository,
	productRepo domain.ProductRepository,
	auditRepo domain.AuditRepository,
) *CreateReceptionOrderUseCase {
	return &CreateReceptionOrderUseCase{
		receptionOrderRepo: receptionOrderRepo,
		receptionLineRepo:  receptionLineRepo,
		supplierRepo:       supplierRepo,
		productRepo:        productRepo,
		auditRepo:          auditRepo,
	}
}

type ReceptionLineInput struct {
	ProductID        uuid.UUID  `json:"product_id"`
	ExpectedQuantity int        `json:"expected_quantity"`
	LotNumber        string     `json:"lot_number"`
	ExpirationDate   *time.Time `json:"expiration_date,omitempty"`
}

type CreateReceptionOrderInput struct {
	SupplierID     uuid.UUID            `json:"supplier_id"`
	InvoiceNumber  string               `json:"invoice_number"`
	InvoiceFileURL string               `json:"invoice_file_url,omitempty"`
	Notes          string               `json:"notes,omitempty"`
	Lines          []ReceptionLineInput `json:"lines"`
	UserID         uuid.UUID            `json:"-"` // Del contexto
}

func (uc *CreateReceptionOrderUseCase) Execute(input CreateReceptionOrderInput) (*domain.ReceptionOrder, error) {
	// 1. Verificar que el proveedor existe
	supplier, err := uc.supplierRepo.FindByID(input.SupplierID)
	if err != nil {
		return nil, errors.New("proveedor no encontrado")
	}

	// 2. Validar que haya líneas
	if len(input.Lines) == 0 {
		return nil, errors.New("la orden debe tener al menos una línea")
	}

	// 3. Generar número de orden único
	orderNumber := fmt.Sprintf("REC-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	// 4. Crear la orden de recepción
	order := &domain.ReceptionOrder{
		OrderNumber:    orderNumber,
		SupplierID:     input.SupplierID,
		Brand:          supplier.Brand,
		InvoiceNumber:  input.InvoiceNumber,
		InvoiceFileURL: input.InvoiceFileURL,
		Status:         domain.ReceptionPendiente,
		Notes:          input.Notes,
	}

	if err := uc.receptionOrderRepo.Create(order); err != nil {
		return nil, err
	}

	// 5. Crear las líneas de la orden
	var lines []*domain.ReceptionLine
	for _, lineInput := range input.Lines {
		// Verificar que el producto existe
		_, err := uc.productRepo.FindByID(lineInput.ProductID)
		if err != nil {
			return nil, fmt.Errorf("producto %s no encontrado", lineInput.ProductID)
		}

		line := &domain.ReceptionLine{
			ReceptionOrderID: order.ID,
			ProductID:        lineInput.ProductID,
			ExpectedQuantity: lineInput.ExpectedQuantity,
			LotNumber:        lineInput.LotNumber,
			ExpirationDate:   lineInput.ExpirationDate,
			Condition:        domain.ConditionApto,
		}
		lines = append(lines, line)
	}

	if err := uc.receptionLineRepo.CreateBatch(lines); err != nil {
		return nil, err
	}

	// 6. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "CREATE_RECEPTION_ORDER",
		EntityType: "RECEPTION_ORDER",
		EntityID:   &order.ID,
		NewValues: map[string]interface{}{
			"order_number": orderNumber,
			"supplier_id":  input.SupplierID,
			"lines_count":  len(lines),
		},
	})

	return order, nil
}
