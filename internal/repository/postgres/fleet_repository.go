package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type VehicleRepositoryPostgres struct {
	db *sqlx.DB
}

func NewVehicleRepository(db *sqlx.DB) domain.VehicleRepository {
	return &VehicleRepositoryPostgres{db: db}
}

func (r *VehicleRepositoryPostgres) Create(vehicle *domain.Vehicle) error {
	query := `
		INSERT INTO vehicles (plate_number, vehicle_type, brand, model, year, capacity_kg, capacity_m3, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, vehicle.PlateNumber, vehicle.VehicleType, vehicle.Brand,
		vehicle.Model, vehicle.Year, vehicle.CapacityKg, vehicle.CapacityM3, vehicle.Status).
		Scan(&vehicle.ID, &vehicle.CreatedAt, &vehicle.UpdatedAt)
}

func (r *VehicleRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	query := `SELECT * FROM vehicles WHERE id = $1`
	err := r.db.Get(&vehicle, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &vehicle, nil
}

func (r *VehicleRepositoryPostgres) Update(vehicle *domain.Vehicle) error {
	query := `
		UPDATE vehicles
		SET status = $1, last_maintenance_date = $2, next_maintenance_date = $3, 
		    is_active = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	result, err := r.db.Exec(query, vehicle.Status, vehicle.LastMaintenanceDate,
		vehicle.NextMaintenanceDate, vehicle.IsActive, vehicle.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *VehicleRepositoryPostgres) ListAvailable(filters map[string]interface{}) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	query := `
		SELECT * FROM vehicles 
		WHERE status = 'DISPONIBLE' AND is_active = true
		ORDER BY vehicle_type
	`
	err := r.db.Select(&vehicles, query)
	return vehicles, err
}

func (r *VehicleRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	query := `SELECT * FROM vehicles ORDER BY plate_number LIMIT $1 OFFSET $2`
	err := r.db.Select(&vehicles, query, limit, offset)
	return vehicles, err
}

// DriverRepositoryPostgres implementa el repositorio de choferes
type DriverRepositoryPostgres struct {
	db *sqlx.DB
}

func NewDriverRepository(db *sqlx.DB) domain.DriverRepository {
	return &DriverRepositoryPostgres{db: db}
}

func (r *DriverRepositoryPostgres) Create(driver *domain.Driver) error {
	query := `
		INSERT INTO drivers (user_id, license_number, license_expiry, phone, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, driver.UserID, driver.LicenseNumber, driver.LicenseExpiry,
		driver.Phone, driver.Status).Scan(&driver.ID, &driver.CreatedAt, &driver.UpdatedAt)
}

func (r *DriverRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Driver, error) {
	var driver domain.Driver
	query := `SELECT * FROM drivers WHERE id = $1`
	err := r.db.Get(&driver, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepositoryPostgres) FindByUserID(userID uuid.UUID) (*domain.Driver, error) {
	var driver domain.Driver
	query := `SELECT * FROM drivers WHERE user_id = $1`
	err := r.db.Get(&driver, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *DriverRepositoryPostgres) Update(driver *domain.Driver) error {
	query := `
		UPDATE drivers
		SET phone = $1, status = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	result, err := r.db.Exec(query, driver.Phone, driver.Status, driver.IsActive, driver.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *DriverRepositoryPostgres) ListAvailable() ([]*domain.Driver, error) {
	var drivers []*domain.Driver
	query := `
		SELECT * FROM drivers 
		WHERE status = 'DISPONIBLE' AND is_active = true
		ORDER BY id
	`
	err := r.db.Select(&drivers, query)
	return drivers, err
}

func (r *DriverRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Driver, error) {
	var drivers []*domain.Driver
	query := `SELECT * FROM drivers ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.Select(&drivers, query, limit, offset)
	return drivers, err
}

// RouteRepositoryPostgres implementa el repositorio de rutas
type RouteRepositoryPostgres struct {
	db *sqlx.DB
}

func NewRouteRepository(db *sqlx.DB) domain.RouteRepository {
	return &RouteRepositoryPostgres{db: db}
}

func (r *RouteRepositoryPostgres) Create(route *domain.Route) error {
	query := `
		INSERT INTO routes (route_number, order_id, vehicle_id, driver_id, route_type, 
		                   departure_date, estimated_arrival, status, assigned_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, route.RouteNumber, route.OrderID, route.VehicleID,
		route.DriverID, route.RouteType, route.DepartureDate, route.EstimatedArrival,
		route.Status, route.AssignedBy).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
}

func (r *RouteRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Route, error) {
	var route domain.Route
	query := `SELECT * FROM routes WHERE id = $1`
	err := r.db.Get(&route, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &route, nil
}

func (r *RouteRepositoryPostgres) Update(route *domain.Route) error {
	query := `
		UPDATE routes
		SET departure_date = $1, actual_arrival = $2, status = $3, 
		    invoice_pdf_url = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	result, err := r.db.Exec(query, route.DepartureDate, route.ActualArrival,
		route.Status, route.InvoicePDFURL, route.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *RouteRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Route, error) {
	var routes []*domain.Route
	query := `SELECT * FROM routes ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&routes, query, limit, offset)
	return routes, err
}

// VehicleMaintenanceRepositoryPostgres implementa el repositorio de mantenimiento
type VehicleMaintenanceRepositoryPostgres struct {
	db *sqlx.DB
}

func NewVehicleMaintenanceRepository(db *sqlx.DB) domain.VehicleMaintenanceRepository {
	return &VehicleMaintenanceRepositoryPostgres{db: db}
}

func (r *VehicleMaintenanceRepositoryPostgres) Create(maintenance *domain.VehicleMaintenance) error {
	query := `
		INSERT INTO vehicle_maintenance (vehicle_id, maintenance_type, description, cost, start_date, performed_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, maintenance.VehicleID, maintenance.MaintenanceType,
		maintenance.Description, maintenance.Cost, maintenance.StartDate, maintenance.PerformedBy).
		Scan(&maintenance.ID, &maintenance.CreatedAt)
}

func (r *VehicleMaintenanceRepositoryPostgres) FindByVehicleID(vehicleID uuid.UUID, limit, offset int) ([]*domain.VehicleMaintenance, error) {
	var maintenances []*domain.VehicleMaintenance
	query := `
		SELECT * FROM vehicle_maintenance
		WHERE vehicle_id = $1
		ORDER BY start_date DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.Select(&maintenances, query, vehicleID, limit, offset)
	return maintenances, err
}

// PreDepartureChecklistRepositoryPostgres implementa el repositorio de check-list
type PreDepartureChecklistRepositoryPostgres struct {
	db *sqlx.DB
}

func NewPreDepartureChecklistRepository(db *sqlx.DB) domain.PreDepartureChecklistRepository {
	return &PreDepartureChecklistRepositoryPostgres{db: db}
}

func (r *PreDepartureChecklistRepositoryPostgres) Create(checklist *domain.PreDepartureChecklist) error {
	query := `
		INSERT INTO pre_departure_checklist (route_id, driver_id, tire_condition, fuel_level, 
		                                    oil_level, lights_ok, damage_photo_url, notes, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return r.db.QueryRow(query, checklist.RouteID, checklist.DriverID, checklist.TireCondition,
		checklist.FuelLevel, checklist.OilLevel, checklist.LightsOk, checklist.DamagePhotoURL,
		checklist.Notes, checklist.CheckedAt).Scan(&checklist.ID)
}

func (r *PreDepartureChecklistRepositoryPostgres) FindByRouteID(routeID uuid.UUID) (*domain.PreDepartureChecklist, error) {
	var checklist domain.PreDepartureChecklist
	query := `SELECT * FROM pre_departure_checklist WHERE route_id = $1`
	err := r.db.Get(&checklist, query, routeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &checklist, nil
}
