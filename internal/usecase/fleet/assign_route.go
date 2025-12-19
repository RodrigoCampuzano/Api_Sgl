package fleet

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

// AssignRouteUseCase implementa HU-10: Asignación inteligente de rutas
type AssignRouteUseCase struct {
	routeRepo   domain.RouteRepository
	vehicleRepo domain.VehicleRepository
	driverRepo  domain.DriverRepository
	orderRepo   domain.OrderRepository
	auditRepo   domain.AuditRepository
}

func NewAssignRouteUseCase(
	routeRepo domain.RouteRepository,
	vehicleRepo domain.VehicleRepository,
	driverRepo domain.DriverRepository,
	orderRepo domain.OrderRepository,
	auditRepo domain.AuditRepository,
) *AssignRouteUseCase {
	return &AssignRouteUseCase{
		routeRepo:   routeRepo,
		vehicleRepo: vehicleRepo,
		driverRepo:  driverRepo,
		orderRepo:   orderRepo,
		auditRepo:   auditRepo,
	}
}

type AssignRouteInput struct {
	OrderID          uuid.UUID        `json:"order_id"`
	VehicleID        *uuid.UUID       `json:"vehicle_id,omitempty"` // Opcional
	DriverID         *uuid.UUID       `json:"driver_id,omitempty"`  // Opcional
	RouteType        domain.RouteType `json:"route_type"`
	DepartureDate    time.Time        `json:"departure_date"`
	EstimatedArrival time.Time        `json:"estimated_arrival"`
	UserID           uuid.UUID        `json:"-"`
}

type AssignRouteOutput struct {
	Route           *domain.Route   `json:"route"`
	AssignedVehicle *domain.Vehicle `json:"assigned_vehicle"`
	AssignedDriver  *domain.Driver  `json:"assigned_driver"`
	AutoAssigned    bool            `json:"auto_assigned"`
}

func (uc *AssignRouteUseCase) Execute(input AssignRouteInput) (*AssignRouteOutput, error) {
	// 1. Verificar que el pedido existe
	order, err := uc.orderRepo.FindByID(input.OrderID)
	if err != nil {
		return nil, errors.New("pedido no encontrado")
	}

	autoAssigned := false

	// 2. HU-10: Asignación inteligente de vehículo si no se especificó
	var vehicleID uuid.UUID
	if input.VehicleID == nil {
		// Buscar vehículo disponible que cumpla con el tipo sugerido
		availableVehicles, _ := uc.vehicleRepo.ListAvailable(nil)

		var selectedVehicle *domain.Vehicle
		for _, v := range availableVehicles {
			if order.SuggestedVehicle != nil && v.VehicleType == *order.SuggestedVehicle {
				selectedVehicle = v
				break
			}
		}

		if selectedVehicle == nil && len(availableVehicles) > 0 {
			selectedVehicle = availableVehicles[0] // Tomar el primero disponible
		}

		if selectedVehicle == nil {
			return nil, domain.ErrVehicleNotAvailable
		}

		vehicleID = selectedVehicle.ID
		autoAssigned = true
	} else {
		vehicleID = *input.VehicleID
	}

	// Verificar disponibilidad y actualizar estado
	vehicle, err := uc.vehicleRepo.FindByID(vehicleID)
	if err != nil {
		return nil, errors.New("vehículo no encontrado")
	}

	if !vehicle.IsAvailableForRoute() {
		return nil, domain.ErrVehicleNotAvailable
	}

	// 3. Asignación inteligente de chofer si no se especificó
	var driverID uuid.UUID
	if input.DriverID == nil {
		availableDrivers, _ := uc.driverRepo.ListAvailable()
		if len(availableDrivers) == 0 {
			return nil, domain.ErrDriverNotAvailable
		}
		driverID = availableDrivers[0].ID
		autoAssigned = true
	} else {
		driverID = *input.DriverID
	}

	driver, err := uc.driverRepo.FindByID(driverID)
	if err != nil {
		return nil, errors.New("chofer no encontrado")
	}

	if !driver.IsAvailableForRoute() {
		return nil, domain.ErrDriverNotAvailable
	}

	// 4. Crear la ruta
	routeNumber := fmt.Sprintf("RTA-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	route := &domain.Route{
		RouteNumber:      routeNumber,
		OrderID:          input.OrderID,
		VehicleID:        vehicleID,
		DriverID:         driverID,
		RouteType:        input.RouteType,
		DepartureDate:    &input.DepartureDate,
		EstimatedArrival: &input.EstimatedArrival,
		Status:           domain.OrderConfirmado,
		AssignedBy:       input.UserID,
	}

	if err := uc.routeRepo.Create(route); err != nil {
		return nil, err
	}

	// 5. Actualizar estados
	vehicle.Status = domain.VehicleEnRuta
	_ = uc.vehicleRepo.Update(vehicle)

	driver.Status = domain.DriverEnRuta
	_ = uc.driverRepo.Update(driver)

	order.Status = domain.OrderEnRuta
	_ = uc.orderRepo.Update(order)

	// 6. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &input.UserID,
		Action:     "ASSIGN_ROUTE",
		EntityType: "ROUTE",
		EntityID:   &route.ID,
		NewValues: map[string]interface{}{
			"route_number":  routeNumber,
			"vehicle_id":    vehicleID,
			"driver_id":     driverID,
			"auto_assigned": autoAssigned,
		},
	})

	return &AssignRouteOutput{
		Route:           route,
		AssignedVehicle: vehicle,
		AssignedDriver:  driver,
		AutoAssigned:    autoAssigned,
	}, nil
}

// GenerateInvoiceUseCase implementa HU-11: Generación de remisión PDF
type GenerateInvoiceUseCase struct {
	routeRepo     domain.RouteRepository
	orderRepo     domain.OrderRepository
	orderLineRepo domain.OrderLineRepository
	customerRepo  domain.CustomerRepository
	auditRepo     domain.AuditRepository
}

func NewGenerateInvoiceUseCase(
	routeRepo domain.RouteRepository,
	orderRepo domain.OrderRepository,
	orderLineRepo domain.OrderLineRepository,
	customerRepo domain.CustomerRepository,
	auditRepo domain.AuditRepository,
) *GenerateInvoiceUseCase {
	return &GenerateInvoiceUseCase{
		routeRepo:     routeRepo,
		orderRepo:     orderRepo,
		orderLineRepo: orderLineRepo,
		customerRepo:  customerRepo,
		auditRepo:     auditRepo,
	}
}

func (uc *GenerateInvoiceUseCase) Execute(routeID uuid.UUID, userID uuid.UUID) (string, error) {
	// 1. Obtener ruta
	route, err := uc.routeRepo.FindByID(routeID)
	if err != nil {
		return "", errors.New("ruta no encontrada")
	}

	// 2. Obtener pedido y cliente
	order, _ := uc.orderRepo.FindByID(route.OrderID)
	customer, _ := uc.customerRepo.FindByID(order.CustomerID)
	lines, _ := uc.orderLineRepo.FindByOrderID(order.ID)

	// HU-11: En producción, aquí se generaría el PDF real
	// Por ahora, simulamos la URL
	pdfURL := fmt.Sprintf("/invoices/invoice_%s_%s.pdf", route.RouteNumber, time.Now().Format("20060102"))

	// Actualizar ruta con la URL del PDF
	route.InvoicePDFURL = pdfURL
	_ = uc.routeRepo.Update(route)

	// Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &userID,
		Action:     "GENERATE_INVOICE",
		EntityType: "ROUTE",
		EntityID:   &route.ID,
		NewValues: map[string]interface{}{
			"invoice_url":   pdfURL,
			"customer_name": customer.Name,
			"total_lines":   len(lines),
		},
	})

	return pdfURL, nil
}
