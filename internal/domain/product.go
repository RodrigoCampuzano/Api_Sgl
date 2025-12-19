package domain

import (
	"time"

	"github.com/google/uuid"
)

// Brand representa las marcas disponibles
type Brand string

const (
	BrandCostena   Brand = "COSTENA"
	BrandJumex     Brand = "JUMEX"
	BrandPronto    Brand = "PRONTO"
	BrandLaCostena Brand = "LA_COSTENA"
	BrandOtros     Brand = "OTROS"
)

// Product representa un producto en el catálogo
type Product struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	SKU       string     `json:"sku" db:"sku"`
	Name      string     `json:"name" db:"name"`
	Brand     Brand      `json:"brand" db:"brand"`
	Category  string     `json:"category" db:"category"`
	Barcode   string     `json:"barcode,omitempty" db:"barcode"`
	WeightKg  float64    `json:"weight_kg" db:"weight_kg"`
	LengthCm  float64    `json:"length_cm" db:"length_cm"`
	WidthCm   float64    `json:"width_cm" db:"width_cm"`
	HeightCm  float64    `json:"height_cm" db:"height_cm"`
	IsFragile bool       `json:"is_fragile" db:"is_fragile"`
	UnitPrice float64    `json:"unit_price" db:"unit_price"`
	IsActive  bool       `json:"is_active" db:"is_active"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// CalculateVolume calcula el volumen en m³
func (p *Product) CalculateVolume() float64 {
	return (p.LengthCm * p.WidthCm * p.HeightCm) / 1000000.0
}

// Supplier representa un proveedor
type Supplier struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Brand       Brand     `json:"brand" db:"brand"`
	RFC         string    `json:"rfc,omitempty" db:"rfc"`
	ContactName string    `json:"contact_name,omitempty" db:"contact_name"`
	Phone       string    `json:"phone,omitempty" db:"phone"`
	Email       string    `json:"email,omitempty" db:"email"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductRepository define los métodos de repositorio para productos
type ProductRepository interface {
	Create(product *Product) error
	FindByID(id uuid.UUID) (*Product, error)
	FindBySKU(sku string) (*Product, error)
	FindByBarcode(barcode string) (*Product, error)
	Update(product *Product) error
	Delete(id uuid.UUID) error
	List(filters map[string]interface{}, limit, offset int) ([]*Product, error)
	Count(filters map[string]interface{}) (int, error)
}

// SupplierRepository define los métodos de repositorio para proveedores
type SupplierRepository interface {
	Create(supplier *Supplier) error
	FindByID(id uuid.UUID) (*Supplier, error)
	Update(supplier *Supplier) error
	List(filters map[string]interface{}, limit, offset int) ([]*Supplier, error)
}
