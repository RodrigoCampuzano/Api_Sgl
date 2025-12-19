package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus representa el estado de un pedido
type OrderStatus string

const (
	OrderBorrador      OrderStatus = "BORRADOR"
	OrderConfirmado    OrderStatus = "CONFIRMADO"
	OrderEnPreparacion OrderStatus = "EN_PREPARACION"
	OrderListo         OrderStatus = "LISTO"
	OrderEnRuta        OrderStatus = "EN_RUTA"
	OrderEntregado     OrderStatus = "ENTREGADO"
	OrderCancelado     OrderStatus = "CANCELADO"
)

// VehicleType representa el tipo de vehículo
type VehicleType string

const (
	VehicleVan       VehicleType = "VAN"
	VehicleCamioneta VehicleType = "CAMIONETA"
	VehicleCamion35  VehicleType = "CAMION_3_5"
	VehicleTorton    VehicleType = "TORTON"
)

// Customer representa un cliente
type Customer struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	RFC         string    `json:"rfc,omitempty" db:"rfc"`
	Address     string    `json:"address,omitempty" db:"address"`
	City        string    `json:"city,omitempty" db:"city"`
	State       string    `json:"state,omitempty" db:"state"`
	PostalCode  string    `json:"postal_code,omitempty" db:"postal_code"`
	Phone       string    `json:"phone,omitempty" db:"phone"`
	Email       string    `json:"email,omitempty" db:"email"`
	CreditLimit float64   `json:"credit_limit" db:"credit_limit"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Order representa un pedido
type Order struct {
	ID               uuid.UUID    `json:"id" db:"id"`
	OrderNumber      string       `json:"order_number" db:"order_number"`
	CustomerID       uuid.UUID    `json:"customer_id" db:"customer_id"`
	Status           OrderStatus  `json:"status" db:"status"`
	TotalWeightKg    float64      `json:"total_weight_kg" db:"total_weight_kg"`
	TotalVolumeM3    float64      `json:"total_volume_m3" db:"total_volume_m3"`
	TotalCost        float64      `json:"total_cost" db:"total_cost"`
	SuggestedVehicle *VehicleType `json:"suggested_vehicle,omitempty" db:"suggested_vehicle"`
	HasFragileItems  bool         `json:"has_fragile_items" db:"has_fragile_items"`
	HasHeavyItems    bool         `json:"has_heavy_items" db:"has_heavy_items"`
	LoadingAlert     string       `json:"loading_alert,omitempty" db:"loading_alert"`
	CreatedBy        uuid.UUID    `json:"created_by" db:"created_by"`
	CreatedAt        time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt        *time.Time   `json:"-" db:"deleted_at"`
}

// SuggestVehicle sugiere el tipo de vehículo basado en volumen (HU-08)
func (o *Order) SuggestVehicle() VehicleType {
	if o.TotalVolumeM3 < 10 {
		return VehicleVan
	} else if o.TotalVolumeM3 < 20 {
		return VehicleCamion35
	}
	return VehicleTorton
}

// GenerateLoadingAlert genera alerta de estiba si hay productos frágiles y pesados (HU-09)
func (o *Order) GenerateLoadingAlert() string {
	if o.HasFragileItems && o.HasHeavyItems {
		return "ADVERTENCIA: Cuidado al estibar. Evitar colocar productos pesados sobre frágiles."
	}
	if o.HasFragileItems {
		return "PRECAUCIÓN: El pedido contiene productos frágiles. Manejar con cuidado."
	}
	return ""
}

// OrderLine representa una línea de pedido
type OrderLine struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	OrderID     uuid.UUID  `json:"order_id" db:"order_id"`
	ProductID   uuid.UUID  `json:"product_id" db:"product_id"`
	InventoryID *uuid.UUID `json:"inventory_id,omitempty" db:"inventory_id"`
	Quantity    int        `json:"quantity" db:"quantity"`
	UnitPrice   float64    `json:"unit_price" db:"unit_price"`
	Subtotal    float64    `json:"subtotal" db:"subtotal"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
}

// OrderRepository define los métodos para pedidos
type OrderRepository interface {
	Create(order *Order) error
	FindByID(id uuid.UUID) (*Order, error)
	FindByOrderNumber(orderNumber string) (*Order, error)
	Update(order *Order) error
	Delete(id uuid.UUID) error
	List(filters map[string]interface{}, limit, offset int) ([]*Order, error)
	FindStuckOrders(hours int, limit, offset int) ([]*Order, error) // HU-24
}

// OrderLineRepository define los métodos para líneas de pedido
type OrderLineRepository interface {
	Create(line *OrderLine) error
	CreateBatch(lines []*OrderLine) error
	FindByOrderID(orderID uuid.UUID) ([]*OrderLine, error)
}

// CustomerRepository define los métodos para clientes
type CustomerRepository interface {
	Create(customer *Customer) error
	FindByID(id uuid.UUID) (*Customer, error)
	Update(customer *Customer) error
	List(filters map[string]interface{}, limit, offset int) ([]*Customer, error)
}
