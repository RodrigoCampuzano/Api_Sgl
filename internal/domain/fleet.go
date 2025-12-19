package domain

import (
	"time"

	"github.com/google/uuid"
)

// VehicleStatus representa el estado de un vehículo
type VehicleStatus string

const (
	VehicleDisponible    VehicleStatus = "DISPONIBLE"
	VehicleEnRuta        VehicleStatus = "EN_RUTA"
	VehicleEnTaller      VehicleStatus = "EN_TALLER"
	VehicleFueraServicio VehicleStatus = "FUERA_SERVICIO"
)

// DriverStatus representa el estado de un chofer
type DriverStatus string

const (
	DriverDisponible DriverStatus = "DISPONIBLE"
	DriverEnRuta     DriverStatus = "EN_RUTA"
	DriverDescanso   DriverStatus = "DESCANSO"
	DriverLicencia   DriverStatus = "LICENCIA"
)

// RouteType representa el tipo de ruta
type RouteType string

const (
	RouteLocal   RouteType = "LOCAL"
	RouteForanea RouteType = "FORANEA"
)

// Vehicle representa un vehículo de la flota
type Vehicle struct {
	ID                  uuid.UUID     `json:"id" db:"id"`
	PlateNumber         string        `json:"plate_number" db:"plate_number"`
	VehicleType         VehicleType   `json:"vehicle_type" db:"vehicle_type"`
	Brand               string        `json:"brand,omitempty" db:"brand"`
	Model               string        `json:"model,omitempty" db:"model"`
	Year                int           `json:"year,omitempty" db:"year"`
	CapacityKg          float64       `json:"capacity_kg" db:"capacity_kg"`
	CapacityM3          float64       `json:"capacity_m3" db:"capacity_m3"`
	Status              VehicleStatus `json:"status" db:"status"`
	LastMaintenanceDate *time.Time    `json:"last_maintenance_date,omitempty" db:"last_maintenance_date"`
	NextMaintenanceDate *time.Time    `json:"next_maintenance_date,omitempty" db:"next_maintenance_date"`
	IsActive            bool          `json:"is_active" db:"is_active"`
	CreatedAt           time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at" db:"updated_at"`
}

// IsAvailableForRoute verifica si el vehículo está disponible
func (v *Vehicle) IsAvailableForRoute() bool {
	return v.Status == VehicleDisponible && v.IsActive
}

// Driver representa un chofer
type Driver struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	UserID        uuid.UUID    `json:"user_id" db:"user_id"`
	LicenseNumber string       `json:"license_number" db:"license_number"`
	LicenseExpiry time.Time    `json:"license_expiry" db:"license_expiry"`
	Phone         string       `json:"phone,omitempty" db:"phone"`
	Status        DriverStatus `json:"status" db:"status"`
	IsActive      bool         `json:"is_active" db:"is_active"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// IsAvailableForRoute verifica si el chofer está disponible
func (d *Driver) IsAvailableForRoute() bool {
	return d.Status == DriverDisponible && d.IsActive && time.Now().Before(d.LicenseExpiry)
}

// Route representa una ruta/viaje
type Route struct {
	ID               uuid.UUID   `json:"id" db:"id"`
	RouteNumber      string      `json:"route_number" db:"route_number"`
	OrderID          uuid.UUID   `json:"order_id" db:"order_id"`
	VehicleID        uuid.UUID   `json:"vehicle_id" db:"vehicle_id"`
	DriverID         uuid.UUID   `json:"driver_id" db:"driver_id"`
	RouteType        RouteType   `json:"route_type" db:"route_type"`
	DepartureDate    *time.Time  `json:"departure_date,omitempty" db:"departure_date"`
	EstimatedArrival *time.Time  `json:"estimated_arrival,omitempty" db:"estimated_arrival"`
	ActualArrival    *time.Time  `json:"actual_arrival,omitempty" db:"actual_arrival"`
	Status           OrderStatus `json:"status" db:"status"`
	InvoicePDFURL    string      `json:"invoice_pdf_url,omitempty" db:"invoice_pdf_url"`
	AssignedBy       uuid.UUID   `json:"assigned_by" db:"assigned_by"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// VehicleMaintenance representa un registro de mantenimiento
type VehicleMaintenance struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	VehicleID       uuid.UUID  `json:"vehicle_id" db:"vehicle_id"`
	MaintenanceType string     `json:"maintenance_type" db:"maintenance_type"`
	Description     string     `json:"description,omitempty" db:"description"`
	Cost            float64    `json:"cost,omitempty" db:"cost"`
	StartDate       time.Time  `json:"start_date" db:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty" db:"end_date"`
	PerformedBy     string     `json:"performed_by,omitempty" db:"performed_by"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// PreDepartureChecklist representa el check-list pre-salida
type PreDepartureChecklist struct {
	ID             uuid.UUID `json:"id" db:"id"`
	RouteID        uuid.UUID `json:"route_id" db:"route_id"`
	DriverID       uuid.UUID `json:"driver_id" db:"driver_id"`
	TireCondition  string    `json:"tire_condition,omitempty" db:"tire_condition"`
	FuelLevel      int       `json:"fuel_level" db:"fuel_level"`
	OilLevel       string    `json:"oil_level,omitempty" db:"oil_level"`
	LightsOk       bool      `json:"lights_ok" db:"lights_ok"`
	DamagePhotoURL string    `json:"damage_photo_url,omitempty" db:"damage_photo_url"`
	Notes          string    `json:"notes,omitempty" db:"notes"`
	CheckedAt      time.Time `json:"checked_at" db:"checked_at"`
}

// VehicleRepository define los métodos para vehículos
type VehicleRepository interface {
	Create(vehicle *Vehicle) error
	FindByID(id uuid.UUID) (*Vehicle, error)
	Update(vehicle *Vehicle) error
	ListAvailable(filters map[string]interface{}) ([]*Vehicle, error)
	List(filters map[string]interface{}, limit, offset int) ([]*Vehicle, error)
}

// DriverRepository define los métodos para choferes
type DriverRepository interface {
	Create(driver *Driver) error
	FindByID(id uuid.UUID) (*Driver, error)
	FindByUserID(userID uuid.UUID) (*Driver, error)
	Update(driver *Driver) error
	ListAvailable() ([]*Driver, error)
	List(filters map[string]interface{}, limit, offset int) ([]*Driver, error)
}

// RouteRepository define los métodos para rutas
type RouteRepository interface {
	Create(route *Route) error
	FindByID(id uuid.UUID) (*Route, error)
	Update(route *Route) error
	List(filters map[string]interface{}, limit, offset int) ([]*Route, error)
}

// VehicleMaintenanceRepository define los métodos para mantenimiento
type VehicleMaintenanceRepository interface {
	Create(maintenance *VehicleMaintenance) error
	FindByVehicleID(vehicleID uuid.UUID, limit, offset int) ([]*VehicleMaintenance, error)
}

// PreDepartureChecklistRepository define los métodos para check-list
type PreDepartureChecklistRepository interface {
	Create(checklist *PreDepartureChecklist) error
	FindByRouteID(routeID uuid.UUID) (*PreDepartureChecklist, error)
}
