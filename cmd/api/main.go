package main

import (
	"fmt"
	"log"

	"github.com/sgl-disasur/api/internal/delivery/http"
	"github.com/sgl-disasur/api/internal/delivery/http/handler"
	"github.com/sgl-disasur/api/internal/infrastructure/config"
	"github.com/sgl-disasur/api/internal/infrastructure/database"
	"github.com/sgl-disasur/api/internal/infrastructure/logger"
	"github.com/sgl-disasur/api/internal/repository/postgres"
	"github.com/sgl-disasur/api/internal/usecase/auth"
	"github.com/sgl-disasur/api/internal/usecase/fleet"
	"github.com/sgl-disasur/api/internal/usecase/inventory"
	"github.com/sgl-disasur/api/internal/usecase/orders"
	"github.com/sgl-disasur/api/internal/usecase/reception"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/sgl-disasur/api/docs" // Swagger docs
)

//	@title			SGL-DISASUR API
//	@version		1.0
//	@description	Sistema de Gestión Logística Multi-Marca
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.email	support@sgl-disasur.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// 1. Cargar configuración
	cfg := config.Load()

	// 2. Inicializar logger
	logger.Init(cfg.LogLevel)
	defer logger.Sync()

	logger.Log.Info("Starting SGL-DISASUR API...")
	logger.Log.Infof("Environment: %s", cfg.GinMode)

	// 3. Conectar a la base de datos
	db, err := database.NewDatabase(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Log.Info("Database connection established")

	// 4. Inicializar repositorios
	userRepo := postgres.NewUserRepository(db.DB)
	sessionRepo := postgres.NewSessionRepository(db.DB)
	auditRepo := postgres.NewAuditRepository(db.DB)
	productRepo := postgres.NewProductRepository(db.DB)
	supplierRepo := postgres.NewSupplierRepository(db.DB)
	receptionOrderRepo := postgres.NewReceptionOrderRepository(db.DB)
	receptionLineRepo := postgres.NewReceptionLineRepository(db.DB)
	receptionDiscrepancyRepo := postgres.NewReceptionDiscrepancyRepository(db.DB)
	inventoryRepo := postgres.NewInventoryRepository(db.DB)
	inventoryMovementRepo := postgres.NewInventoryMovementRepository(db.DB)
	cycleCountRepo := postgres.NewCycleCountRepository(db.DB)
	orderRepo := postgres.NewOrderRepository(db.DB)
	orderLineRepo := postgres.NewOrderLineRepository(db.DB)
	customerRepo := postgres.NewCustomerRepository(db.DB)
	vehicleRepo := postgres.NewVehicleRepository(db.DB)
	driverRepo := postgres.NewDriverRepository(db.DB)
	routeRepo := postgres.NewRouteRepository(db.DB)
	maintenanceRepo := postgres.NewVehicleMaintenanceRepository(db.DB)
	checklistRepo := postgres.NewPreDepartureChecklistRepository(db.DB)

	// 5. Inicializar casos de uso
	// Auth
	loginUseCase := auth.NewLoginUseCase(
		userRepo,
		sessionRepo,
		auditRepo,
		cfg.JWTSecretKey,
		cfg.JWTExpirationHours,
	)
	registerUserUseCase := auth.NewRegisterUserUseCase(userRepo, auditRepo)

	// Reception
	createReceptionOrderUC := reception.NewCreateReceptionOrderUseCase(
		receptionOrderRepo,
		receptionLineRepo,
		supplierRepo,
		productRepo,
		auditRepo,
	)

	blindCountUC := reception.NewBlindCountUseCase(
		receptionOrderRepo,
		receptionLineRepo,
		receptionDiscrepancyRepo,
		auditRepo,
	)

	// Inventory
	getStockUC := inventory.NewGetStockUseCase(inventoryRepo, productRepo)
	getFEFOLotsUC := inventory.NewGetFEFOLotsUseCase(inventoryRepo)
	registerDamageUC := inventory.NewRegisterDamageUseCase(inventoryRepo, inventoryMovementRepo, auditRepo)
	cycleCountUC := inventory.NewPerformCycleCountUseCase(
		cycleCountRepo,
		inventoryRepo,
		inventoryMovementRepo,
		productRepo,
		auditRepo,
	)

	// Orders
	createOrderUC := orders.NewCreateOrderUseCase(
		orderRepo,
		orderLineRepo,
		customerRepo,
		productRepo,
		inventoryRepo,
		auditRepo,
	)

	// Fleet
	assignRouteUC := fleet.NewAssignRouteUseCase(routeRepo, vehicleRepo, driverRepo, orderRepo, auditRepo)
	generateInvoiceUC := fleet.NewGenerateInvoiceUseCase(routeRepo, orderRepo, orderLineRepo, customerRepo, auditRepo)
	registerMaintenanceUC := fleet.NewRegisterMaintenanceUseCase(maintenanceRepo, vehicleRepo, auditRepo)
	preDepartureCheckUC := fleet.NewPerformPreDepartureCheckUseCase(checklistRepo, routeRepo, vehicleRepo, auditRepo)

	// 6. Inicializar handlers
	authHandler := handler.NewAuthHandler(loginUseCase, registerUserUseCase)
	productHandler := handler.NewProductHandler(productRepo)
	receptionHandler := handler.NewReceptionHandler(
		createReceptionOrderUC,
		blindCountUC,
		productRepo,
		supplierRepo,
		receptionOrderRepo,
		receptionLineRepo,
	)
	inventoryHandler := handler.NewInventoryHandler(
		getStockUC,
		getFEFOLotsUC,
		registerDamageUC,
		cycleCountUC,
	)
	orderHandler := handler.NewOrderHandler(
		createOrderUC,
		orderRepo,
		orderLineRepo,
		customerRepo,
	)
	fleetHandler := handler.NewFleetHandler(
		assignRouteUC,
		generateInvoiceUC,
		registerMaintenanceUC,
		preDepartureCheckUC,
		vehicleRepo,
		driverRepo,
		routeRepo,
		maintenanceRepo,
	)

	// File upload handler
	uploadPath := cfg.StoragePath
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	fileHandler := handler.NewFileHandler(uploadPath, 10) // 10MB max

	// 7. Configurar router
	routerConfig := &http.RouterConfig{
		AuthHandler:      authHandler,
		ProductHandler:   productHandler,
		ReceptionHandler: receptionHandler,
		InventoryHandler: inventoryHandler,
		OrderHandler:     orderHandler,
		FleetHandler:     fleetHandler,
		FileHandler:      fileHandler,
		SecretKey:        cfg.JWTSecretKey,
	}
	router := http.SetupRouter(routerConfig)

	// Swagger Documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 8. Iniciar servidor
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Log.Infof("Server starting on %s", addr)
	logger.Log.Info("API Documentation available at /swagger/index.html")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
