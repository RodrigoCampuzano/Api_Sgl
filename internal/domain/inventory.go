package domain

import (
	"time"

	"github.com/google/uuid"
)

// StockStatus representa el estado del stock
type StockStatus string

const (
	StockDisponible StockStatus = "DISPONIBLE"
	StockReservado  StockStatus = "RESERVADO"
	StockBloqueado  StockStatus = "BLOQUEADO"
	StockCuarentena StockStatus = "CUARENTENA"
	StockCaducado   StockStatus = "CADUCADO"
)

// MovementType representa el tipo de movimiento de inventario
type MovementType string

const (
	MovementEntrada       MovementType = "ENTRADA"
	MovementSalida        MovementType = "SALIDA"
	MovementAjuste        MovementType = "AJUSTE"
	MovementMerma         MovementType = "MERMA"
	MovementDevolucion    MovementType = "DEVOLUCION"
	MovementTransferencia MovementType = "TRANSFERENCIA"
)

// Inventory representa el inventario de un producto
type Inventory struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	ProductID         uuid.UUID   `json:"product_id" db:"product_id"`
	LotNumber         string      `json:"lot_number" db:"lot_number"`
	ExpirationDate    *time.Time  `json:"expiration_date,omitempty" db:"expiration_date"`
	Quantity          int         `json:"quantity" db:"quantity"`
	Status            StockStatus `json:"status" db:"status"`
	WarehouseLocation string      `json:"warehouse_location,omitempty" db:"warehouse_location"`
	LastMovementAt    *time.Time  `json:"last_movement_at,omitempty" db:"last_movement_at"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

// IsExpired verifica si el producto está caducado
func (i *Inventory) IsExpired() bool {
	if i.ExpirationDate == nil {
		return false
	}
	return time.Now().After(*i.ExpirationDate)
}

// DaysUntilExpiration retorna cuántos días faltan para la caducidad
func (i *Inventory) DaysUntilExpiration() *int {
	if i.ExpirationDate == nil {
		return nil
	}
	days := int(time.Until(*i.ExpirationDate).Hours() / 24)
	return &days
}

// InventoryMovement representa un movimiento de inventario
type InventoryMovement struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	InventoryID      uuid.UUID    `json:"inventory_id" db:"inventory_id"`
	MovementType     MovementType `json:"movement_type" db:"movement_type"`
	Quantity         int          `json:"quantity" db:"quantity"`
	PreviousQuantity int          `json:"previous_quantity" db:"previous_quantity"`
	NewQuantity      int          `json:"new_quantity" db:"new_quantity"`
	ReferenceID      *uuid.UUID   `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType    string       `json:"reference_type,omitempty" db:"reference_type"`
	Reason           string       `json:"reason,omitempty" db:"reason"`
	EvidencePhotoURL string       `json:"evidence_photo_url,omitempty" db:"evidence_photo_url"`
	PerformedBy      uuid.UUID    `json:"performed_by" db:"performed_by"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
}

// CycleCount representa un conteo cíclico
type CycleCount struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ScheduledDate    time.Time  `json:"scheduled_date" db:"scheduled_date"`
	Location         string     `json:"location,omitempty" db:"location"`
	ProductID        uuid.UUID  `json:"product_id" db:"product_id"`
	ExpectedQuantity *int       `json:"expected_quantity,omitempty" db:"expected_quantity"`
	CountedQuantity  *int       `json:"counted_quantity,omitempty" db:"counted_quantity"`
	Variance         *int       `json:"variance,omitempty" db:"variance"`
	CountedBy        *uuid.UUID `json:"counted_by,omitempty" db:"counted_by"`
	CountedAt        *time.Time `json:"counted_at,omitempty" db:"counted_at"`
	Status           string     `json:"status" db:"status"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// InventoryRepository define los métodos para inventario
type InventoryRepository interface {
	Create(inventory *Inventory) error
	FindByID(id uuid.UUID) (*Inventory, error)
	FindByProduct(productID uuid.UUID) ([]*Inventory, error)
	FindByProductFEFO(productID uuid.UUID) ([]*Inventory, error) // First Expired First Out
	Update(inventory *Inventory) error
	ListAvailable(filters map[string]interface{}, limit, offset int) ([]*Inventory, error)
	GetStockByProduct(productID uuid.UUID) (int, error)
}

// InventoryMovementRepository define los métodos para movimientos
type InventoryMovementRepository interface {
	Create(movement *InventoryMovement) error
	FindByInventoryID(inventoryID uuid.UUID, limit, offset int) ([]*InventoryMovement, error)
	FindByDateRange(from, to time.Time, limit, offset int) ([]*InventoryMovement, error)
}

// CycleCountRepository define los métodos para conteo cíclico
type CycleCountRepository interface {
	Create(count *CycleCount) error
	FindByID(id uuid.UUID) (*CycleCount, error)
	Update(count *CycleCount) error
	ListPending(limit, offset int) ([]*CycleCount, error)
}
