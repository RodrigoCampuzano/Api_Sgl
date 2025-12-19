package fleet

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// RegisterMaintenanceUseCase implementa HU-16: Control de mantenimiento
type RegisterMaintenanceUseCase struct {
	maintenanceRepo domain.VehicleMaintenanceRepository
	vehicleRepo     domain.VehicleRepository
	auditRepo       domain.AuditRepository
}

func NewRegisterMaintenanceUseCase(
	maintenanceRepo domain.VehicleMaintenanceRepository,
	vehicleRepo domain.VehicleRepository,
	auditRepo domain.AuditRepository,
) *RegisterMaintenanceUseCase {
	return &RegisterMaintenanceUseCase{
		maintenanceRepo: maintenanceRepo,
		vehicleRepo:     vehicleRepo,
		auditRepo:       auditRepo,
	}
}

type RegisterMaintenanceInput struct {
	VehicleID       uuid.UUID `json:"vehicle_id"`
	MaintenanceType string    `json:"maintenance_type"` // "PREVENTIVO", "CORRECTIVO", "NEUMATICOS"
	Description     string    `json:"description"`
	Cost            float64   `json:"cost"`
	PerformedBy     string    `json:"performed_by"`
	UserID          uuid.UUID `json:"-"`
}

func (uc *RegisterMaintenanceUseCase) Execute(input RegisterMaintenanceInput) error {
	// 1. Verificar vehículo
	vehicle, err := uc.vehicleRepo.FindByID(input.VehicleID)
	if err != nil {
		return errors.New("vehículo no encontrado")
	}

	// 2. Crear registro de mantenimiento
	maintenance := &domain.VehicleMaintenance{
		VehicleID:       input.VehicleID,
		MaintenanceType: input.MaintenanceType,
		Description:     input.Description,
		Cost:            input.Cost,
		StartDate:       time.Now(),
		PerformedBy:     input.PerformedBy,
	}

	if err := uc.maintenanceRepo.Create(maintenance); err != nil {
		return err
	}

	// 3. Actualizar estado del vehículo
	vehicle.Status = domain.VehicleEnTaller
	now := time.Now()
	vehicle.LastMaintenanceDate = &now

	// Programar próximo mantenimiento (ejemplo: 3 meses después)
	nextMaintenance := now.AddDate(0, 3, 0)
	vehicle.NextMaintenanceDate = &nextMaintenance

	_ = uc.vehicleRepo.Update(vehicle)

	// 4. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "REGISTER_MAINTENANCE",
		EntityType: "VEHICLE",
		EntityID:   &vehicle.ID,
		NewValues: map[string]interface{}{
			"maintenance_type":      input.MaintenanceType,
			"cost":                  input.Cost,
			"next_maintenance_date": nextMaintenance,
		},
	})

	return nil
}

// PerformPreDepartureCheckUseCase implementa HU-17: Check-list pre-salida
type PerformPreDepartureCheckUseCase struct {
	checklistRepo domain.PreDepartureChecklistRepository
	routeRepo     domain.RouteRepository
	vehicleRepo   domain.VehicleRepository
	auditRepo     domain.AuditRepository
}

func NewPerformPreDepartureCheckUseCase(
	checklistRepo domain.PreDepartureChecklistRepository,
	routeRepo domain.RouteRepository,
	vehicleRepo domain.VehicleRepository,
	auditRepo domain.AuditRepository,
) *PerformPreDepartureCheckUseCase {
	return &PerformPreDepartureCheckUseCase{
		checklistRepo: checklistRepo,
		routeRepo:     routeRepo,
		vehicleRepo:   vehicleRepo,
		auditRepo:     auditRepo,
	}
}

type PerformCheckInput struct {
	RouteID        uuid.UUID `json:"route_id"`
	TireCondition  string    `json:"tire_condition"` // "BUENO", "REGULAR", "MALO"
	FuelLevel      int       `json:"fuel_level"`     // 0-100%
	OilLevel       string    `json:"oil_level"`      // "OK", "BAJO"
	LightsOk       bool      `json:"lights_ok"`
	DamagePhotoURL string    `json:"damage_photo_url,omitempty"` // Foto de daños previos
	Notes          string    `json:"notes,omitempty"`
	UserID         uuid.UUID `json:"-"`
}

func (uc *PerformPreDepartureCheckUseCase) Execute(input PerformCheckInput) error {
	// 1. Verificar ruta
	route, err := uc.routeRepo.FindByID(input.RouteID)
	if err != nil {
		return errors.New("ruta no encontrada")
	}

	// 2. Validaciones de seguridad (HU-17)
	if input.FuelLevel < 25 {
		return errors.New("nivel de combustible insuficiente (<25%). Cargar combustible antes de partir")
	}

	if input.TireCondition == "MALO" {
		return errors.New("condición de llantas MALA. Cambiar llantas antes de partir")
	}

	if !input.LightsOk {
		return errors.New("luces no funcionan correctamente. Reparar antes de partir")
	}

	// 3. Crear check-list
	checklist := &domain.PreDepartureChecklist{
		RouteID:        input.RouteID,
		DriverID:       route.DriverID,
		TireCondition:  input.TireCondition,
		FuelLevel:      input.FuelLevel,
		OilLevel:       input.OilLevel,
		LightsOk:       input.LightsOk,
		DamagePhotoURL: input.DamagePhotoURL,
		Notes:          input.Notes,
		CheckedAt:      time.Now(),
	}

	if err := uc.checklistRepo.Create(checklist); err != nil {
		return err
	}

	// 4. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "PRE_DEPARTURE_CHECK",
		EntityType: "ROUTE",
		EntityID:   &route.ID,
		NewValues: map[string]interface{}{
			"tire_condition": input.TireCondition,
			"fuel_level":     input.FuelLevel,
			"lights_ok":      input.LightsOk,
		},
	})

	return nil
}
