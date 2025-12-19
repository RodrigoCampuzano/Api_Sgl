package domain

import (
	"time"

	"github.com/google/uuid"
)

// ReceptionStatus representa el estado de una orden de recepción
type ReceptionStatus string

const (
	ReceptionPendiente     ReceptionStatus = "PENDIENTE"
	ReceptionEnConteo      ReceptionStatus = "EN_CONTEO"
	ReceptionValidada      ReceptionStatus = "VALIDADA"
	ReceptionConIncidencia ReceptionStatus = "CON_INCIDENCIA"
	ReceptionCompletada    ReceptionStatus = "COMPLETADA"
)

// ProductCondition representa la condición de un producto
type ProductCondition string

const (
	ConditionApto       ProductCondition = "APTO"
	ConditionDesecho    ProductCondition = "DESECHO"
	ConditionCuarentena ProductCondition = "CUARENTENA"
)

// DiscrepancyStatus representa el estado de una discrepancia
type DiscrepancyStatus string

const (
	DiscrepancyDetectada  DiscrepancyStatus = "DETECTADA"
	DiscrepancyEnRevision DiscrepancyStatus = "EN_REVISION"
	DiscrepancyResuelta   DiscrepancyStatus = "RESUELTA"
	DiscrepancyAceptada   DiscrepancyStatus = "ACEPTADA"
)

// ReceptionOrder representa una orden de recepción
type ReceptionOrder struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	OrderNumber    string          `json:"order_number" db:"order_number"`
	SupplierID     uuid.UUID       `json:"supplier_id" db:"supplier_id"`
	Brand          Brand           `json:"brand" db:"brand"`
	InvoiceNumber  string          `json:"invoice_number,omitempty" db:"invoice_number"`
	InvoiceFileURL string          `json:"invoice_file_url,omitempty" db:"invoice_file_url"`
	Status         ReceptionStatus `json:"status" db:"status"`
	ReceivedBy     *uuid.UUID      `json:"received_by,omitempty" db:"received_by"`
	ReceivedAt     *time.Time      `json:"received_at,omitempty" db:"received_at"`
	ValidatedBy    *uuid.UUID      `json:"validated_by,omitempty" db:"validated_by"`
	ValidatedAt    *time.Time      `json:"validated_at,omitempty" db:"validated_at"`
	Notes          string          `json:"notes,omitempty" db:"notes"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// ReceptionLine representa una línea de una orden de recepción
type ReceptionLine struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	ReceptionOrderID uuid.UUID        `json:"reception_order_id" db:"reception_order_id"`
	ProductID        uuid.UUID        `json:"product_id" db:"product_id"`
	ExpectedQuantity int              `json:"expected_quantity" db:"expected_quantity"`
	CountedQuantity  *int             `json:"counted_quantity,omitempty" db:"counted_quantity"`
	Discrepancy      *int             `json:"discrepancy,omitempty" db:"discrepancy"`
	LotNumber        string           `json:"lot_number,omitempty" db:"lot_number"`
	ExpirationDate   *time.Time       `json:"expiration_date,omitempty" db:"expiration_date"`
	Condition        ProductCondition `json:"condition" db:"condition"`
	CountedBy        *uuid.UUID       `json:"counted_by,omitempty" db:"counted_by"`
	CountedAt        *time.Time       `json:"counted_at,omitempty" db:"counted_at"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// HasDiscrepancy verifica si hay discrepancia en la línea
func (rl *ReceptionLine) HasDiscrepancy() bool {
	if rl.CountedQuantity == nil {
		return false
	}
	return *rl.CountedQuantity != rl.ExpectedQuantity
}

// ReceptionDiscrepancy representa una discrepancia en la recepción
type ReceptionDiscrepancy struct {
	ID              uuid.UUID         `json:"id" db:"id"`
	ReceptionLineID uuid.UUID         `json:"reception_line_id" db:"reception_line_id"`
	ExpectedQty     int               `json:"expected_qty" db:"expected_qty"`
	CountedQty      int               `json:"counted_qty" db:"counted_qty"`
	Difference      int               `json:"difference" db:"difference"`
	Status          DiscrepancyStatus `json:"status" db:"status"`
	ResolutionNotes string            `json:"resolution_notes,omitempty" db:"resolution_notes"`
	ResolvedBy      *uuid.UUID        `json:"resolved_by,omitempty" db:"resolved_by"`
	ResolvedAt      *time.Time        `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// ReceptionOrderRepository define los métodos para órdenes de recepción
type ReceptionOrderRepository interface {
	Create(order *ReceptionOrder) error
	FindByID(id uuid.UUID) (*ReceptionOrder, error)
	FindByOrderNumber(orderNumber string) (*ReceptionOrder, error)
	Update(order *ReceptionOrder) error
	List(filters map[string]interface{}, limit, offset int) ([]*ReceptionOrder, error)
}

// ReceptionLineRepository define los métodos para líneas de recepción
type ReceptionLineRepository interface {
	Create(line *ReceptionLine) error
	CreateBatch(lines []*ReceptionLine) error
	FindByID(id uuid.UUID) (*ReceptionLine, error)
	FindByOrderID(orderID uuid.UUID) ([]*ReceptionLine, error)
	Update(line *ReceptionLine) error
}

// ReceptionDiscrepancyRepository define los métodos para discrepancias
type ReceptionDiscrepancyRepository interface {
	Create(discrepancy *ReceptionDiscrepancy) error
	FindByLineID(lineID uuid.UUID) (*ReceptionDiscrepancy, error)
	Update(discrepancy *ReceptionDiscrepancy) error
	ListPending(limit, offset int) ([]*ReceptionDiscrepancy, error)
}
