package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sgl-disasur/api/internal/delivery/http/handler"
	"github.com/sgl-disasur/api/internal/delivery/http/middleware"
)

type RouterConfig struct {
	AuthHandler      *handler.AuthHandler
	ProductHandler   *handler.ProductHandler
	ReceptionHandler *handler.ReceptionHandler
	InventoryHandler *handler.InventoryHandler
	OrderHandler     *handler.OrderHandler
	FleetHandler     *handler.FleetHandler
	FileHandler      *handler.FileHandler
	SecretKey        string
}

func SetupRouter(config *RouterConfig) *gin.Engine {
	r := gin.Default()

	// Middleware global
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "SGL-DISASUR API",
			"version": "1.0.0",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Rutas públicas - Autenticación
		auth := v1.Group("/auth")
		{
			auth.POST("/login", config.AuthHandler.Login)
			auth.POST("/register", config.AuthHandler.Register) // Registro público
		}

		// Rutas protegidas
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(config.SecretKey))
		{
			// Logout (requiere autenticación)
			auth.POST("/logout", config.AuthHandler.Logout)

			// === MÓDULO 0: USUARIOS ===
			users := protected.Group("/users")
			users.Use(middleware.RequireRole("ADMIN_TI", "GERENTE"))
			{
				users.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Lista de usuarios - TODO"})
				})
			}

			// === MÓDULO 1: PRODUCTOS (HU-04) ===
			products := protected.Group("/products")
			{
				products.GET("", config.ProductHandler.List)
				products.GET("/:id", config.ProductHandler.GetByID)
				products.POST("",
					middleware.RequireRole("ADMIN_TI", "JEFE_ALMACEN"),
					config.ProductHandler.Create)
			}

			// === MÓDULO 1: RECEPCIÓN ===
			reception := protected.Group("/reception")
			{
				// HU-01: Alta de órdenes
				reception.POST("/orders",
					middleware.RequireRole("JEFE_ALMACEN", "SUPERVISOR"),
					config.ReceptionHandler.CreateOrder)

				reception.GET("/orders", config.ReceptionHandler.ListOrders)
				reception.GET("/orders/:id", config.ReceptionHandler.GetOrder)

				// HU-02/03: Conteo ciego y validación
				reception.POST("/blind-count",
					middleware.RequireRole("AUXILIAR", "RECEPCIONISTA"),
					config.ReceptionHandler.BlindCount)

				// HU-14: Devoluciones - ADMIN_TI AGREGADO
				reception.POST("/returns",
					middleware.RequireRole("RECEPCIONISTA", "JEFE_ALMACEN", "ADMIN_TI"),
					config.ReceptionHandler.ProcessReturn)
			}

			// === MÓDULO 2: INVENTARIO ===
			inventory := protected.Group("/inventory")
			{
				// HU-05: Monitor de stock
				inventory.GET("/stock", config.InventoryHandler.GetStock)

				// HU-06: FEFO
				inventory.GET("/fefo/:product_id", config.InventoryHandler.GetFEFOLots)

				// HU-13: Registro de mermas
				inventory.POST("/damages",
					middleware.RequireRole("JEFE_ALMACEN", "SUPERVISOR"),
					config.InventoryHandler.RegisterDamage)

				// HU-15: Conteo cíclico
				inventory.POST("/cycle-counts/generate",
					middleware.RequireRole("JEFE_ALMACEN"),
					config.InventoryHandler.GenerateCycleCounts)

				inventory.POST("/cycle-counts/perform",
					middleware.RequireRole("AUXILIAR", "MONTACARGUISTA"),
					config.InventoryHandler.PerformCycleCount)
			}

			// === MÓDULO 3: PEDIDOS ===
			orders := protected.Group("/orders")
			{
				// HU-07/08/09/18: Crear pedido
				orders.POST("",
					middleware.RequireRole("VENDEDOR", "JEFE_TRAFICO"),
					config.OrderHandler.CreateOrder)

				orders.GET("", config.OrderHandler.ListOrders)
				orders.GET("/:id", config.OrderHandler.GetOrder)

				// HU-24: Pedidos atorados
				orders.GET("/stuck", config.OrderHandler.GetStuckOrders)
			}

			// Clientes
			customers := protected.Group("/customers")
			{
				customers.GET("", config.OrderHandler.ListCustomers)
				customers.POST("",
					middleware.RequireRole("ADMIN_TI", "JEFE_TRAFICO", "VENDEDOR"),
					config.OrderHandler.CreateCustomer)
			}

			// === MÓDULO 4: FLOTA ===
			fleet := protected.Group("/fleet")
			{
				// Vehículos
				fleet.GET("/vehicles", config.FleetHandler.ListVehicles)

				// HU-16: Control de mantenimiento
				fleet.POST("/vehicles/maintenance",
					middleware.RequireRole("FLOTA", "JEFE_TRAFICO"),
					config.FleetHandler.RegisterMaintenance)

				// Choferes
				fleet.GET("/drivers", config.FleetHandler.ListDrivers)

				// Rutas
				fleet.GET("/routes", config.FleetHandler.ListRoutes)

				// HU-10: Asignar ruta
				fleet.POST("/routes",
					middleware.RequireRole("JEFE_TRAFICO"),
					config.FleetHandler.AssignRoute)

				// HU-11: Generar remisión
				fleet.POST("/routes/:route_id/invoice",
					middleware.RequireRole("JEFE_TRAFICO", "VENDEDOR"),
					config.FleetHandler.GenerateInvoice)

				// HU-17: Check-list pre-salida
				fleet.POST("/routes/pre-departure-check",
					middleware.RequireRole("CHOFER"),
					config.FleetHandler.PerformPreDepartureCheck)
			}

			// === MÓDULO 6: REPORT ES ===
			reports := protected.Group("/reports")
			reports.Use(middleware.RequireRole("GERENTE", "ADMIN_TI"))
			{
				// HU-12: Dashboard
				reports.GET("/dashboard", func(c *gin.Context) {
					c.JSON(200, gin.H{
						"message": "Dashboard - En desarrollo",
						"kpis": map[string]interface{}{
							"ventas_hoy":      0,
							"camiones_fuera":  0,
							"pedidos_activos": 0,
						},
					})
				})

				// HU-23: Reporte de rotación
				reports.GET("/rotation", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Reporte de rotación - En desarrollo"})
				})

				// HU-24: Pedidos atorados
				reports.GET("/stuck-orders", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Pedidos atorados - En desarrollo"})
				})
			}

			// === FILES (UPLOAD) ===
			files := protected.Group("/files")
			{
				files.POST("/upload", config.FileHandler.UploadFile)
			}
		}
	}

	// Servir archivos subidos (público)
	r.GET("/uploads/*filepath", config.FileHandler.ServeFile)

	return r
}
