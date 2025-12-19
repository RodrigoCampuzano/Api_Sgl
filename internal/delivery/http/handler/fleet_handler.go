package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/usecase/fleet"
)

type FleetHandler struct {
	assignRouteUC         *fleet.AssignRouteUseCase
	generateInvoiceUC     *fleet.GenerateInvoiceUseCase
	registerMaintenanceUC *fleet.RegisterMaintenanceUseCase
	preDepartureCheckUC   *fleet.PerformPreDepartureCheckUseCase
	vehicleRepo           domain.VehicleRepository
	driverRepo            domain.DriverRepository
	routeRepo             domain.RouteRepository
	maintenanceRepo       domain.VehicleMaintenanceRepository
}

func NewFleetHandler(
	assignRouteUC *fleet.AssignRouteUseCase,
	generateInvoiceUC *fleet.GenerateInvoiceUseCase,
	registerMaintenanceUC *fleet.RegisterMaintenanceUseCase,
	preDepartureCheckUC *fleet.PerformPreDepartureCheckUseCase,
	vehicleRepo domain.VehicleRepository,
	driverRepo domain.DriverRepository,
	routeRepo domain.RouteRepository,
	maintenanceRepo domain.VehicleMaintenanceRepository,
) *FleetHandler {
	return &FleetHandler{
		assignRouteUC:         assignRouteUC,
		generateInvoiceUC:     generateInvoiceUC,
		registerMaintenanceUC: registerMaintenanceUC,
		preDepartureCheckUC:   preDepartureCheckUC,
		vehicleRepo:           vehicleRepo,
		driverRepo:            driverRepo,
		routeRepo:             routeRepo,
		maintenanceRepo:       maintenanceRepo,
	}
}

// AssignRoute godoc
// @Summary      Asignar ruta (HU-10)
// @Description  Asigna vehículo y chofer a un pedido de forma inteligente
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Param        route  body      fleet.AssignRouteInput  true  "Datos de la ruta"
// @Success      201    {object}  fleet.AssignRouteOutput
// @Security     Bearer
// @Router       /api/v1/fleet/routes [post]
func (h *FleetHandler) AssignRoute(c *gin.Context) {
	var input fleet.AssignRouteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	result, err := h.assignRouteUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GenerateInvoice godoc
// @Summary      Generar remisión (HU-11)
// @Description  Genera PDF de remisión para una ruta
// @Tags         fleet
// @Produce      json
// @Param        route_id  path      string  true  "Route ID"
// @Success      200       {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/fleet/routes/{route_id}/invoice [post]
func (h *FleetHandler) GenerateInvoice(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := uuid.Parse(routeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ruta inválido"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	pdfURL, err := h.generateInvoiceUC.Execute(routeID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Remisión generada exitosamente",
		"invoice_url": pdfURL,
	})
}

// RegisterMaintenance godoc
// @Summary      Registrar mantenimiento (HU-16)
// @Description  Registra mantenimiento de vehículo
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Param        maintenance  body      fleet.RegisterMaintenanceInput  true  "Datos del mantenimiento"
// @Success      200          {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/fleet/vehicles/maintenance [post]
func (h *FleetHandler) RegisterMaintenance(c *gin.Context) {
	var input fleet.RegisterMaintenanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	if err := h.registerMaintenanceUC.Execute(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mantenimiento registrado exitosamente"})
}

// PerformPreDepartureCheck godoc
// @Summary      Check-list pre-salida (HU-17)
// @Description  Registra check-list de seguridad antes de partir
// @Tags         fleet
// @Accept       json
// @Produce      json
// @Param        check  body      fleet.PerformCheckInput  true  "Datos del check-list"
// @Success      200    {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/fleet/routes/pre-departure-check [post]
func (h *FleetHandler) PerformPreDepartureCheck(c *gin.Context) {
	var input fleet.PerformCheckInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	if err := h.preDepartureCheckUC.Execute(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Check-list completado. Vehículo listo para partir"})
}

// ListVehicles godoc
// @Summary      Listar vehículos
// @Description  Obtiene listado de vehículos de la flota
// @Tags         fleet
// @Produce      json
// @Success      200  {array}   domain.Vehicle
// @Security     Bearer
// @Router       /api/v1/fleet/vehicles [get]
func (h *FleetHandler) ListVehicles(c *gin.Context) {
	vehicles, err := h.vehicleRepo.List(nil, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vehicles)
}

// ListDrivers godoc
// @Summary      Listar choferes
// @Description  Obtiene listado de choferes
// @Tags         fleet
// @Produce      json
// @Success      200  {array}   domain.Driver
// @Security     Bearer
// @Router       /api/v1/fleet/drivers [get]
func (h *FleetHandler) ListDrivers(c *gin.Context) {
	drivers, err := h.driverRepo.List(nil, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}

// ListRoutes godoc
// @Summary      Listar rutas
// @Description  Obtiene listado de rutas/viajes
// @Tags         fleet
// @Produce      json
// @Success      200  {array}   domain.Route
// @Security     Bearer
// @Router       /api/v1/fleet/routes [get]
func (h *FleetHandler) ListRoutes(c *gin.Context) {
	routes, err := h.routeRepo.List(nil, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}
