package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/usecase/orders"
)

type OrderHandler struct {
	createOrderUC *orders.CreateOrderUseCase
	orderRepo     domain.OrderRepository
	orderLineRepo domain.OrderLineRepository
	customerRepo  domain.CustomerRepository
}

func NewOrderHandler(
	createOrderUC *orders.CreateOrderUseCase,
	orderRepo domain.OrderRepository,
	orderLineRepo domain.OrderLineRepository,
	customerRepo domain.CustomerRepository,
) *OrderHandler {
	return &OrderHandler{
		createOrderUC: createOrderUC,
		orderRepo:     orderRepo,
		orderLineRepo: orderLineRepo,
		customerRepo:  customerRepo,
	}
}

// CreateOrder godoc
// @Summary      Crear pedido (HU-07, 08, 09, 18)
// @Description  Crea pedido multi-marca con sugerencia de vehículo y alertas
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body      orders.CreateOrderInput  true  "Datos del pedido"
// @Success      201    {object}  orders.CreateOrderOutput
// @Security     Bearer
// @Router       /api/v1/orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input orders.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	result, err := h.createOrderUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// ListOrders godoc
// @Summary      Listar pedidos
// @Description  Obtiene listado de pedidos con filtros
// @Tags         orders
// @Produce      json
// @Param        status  query     string  false  "Filtrar por estado"
// @Success      200     {array}   domain.Order
// @Security     Bearer
// @Router       /api/v1/orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	orders, err := h.orderRepo.List(filters, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary      Obtener pedido por ID
// @Description  Obtiene detalles de un pedido con sus líneas
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Security     Bearer
// @Router       /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := h.orderRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}

	lines, _ := h.orderLineRepo.FindByOrderID(id)

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"lines": lines,
	})
}

// CreateCustomer godoc
// @Summary      Crear cliente
// @Description  Crea un nuevo cliente en el sistema
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        customer  body      domain.Customer  true  "Datos del cliente"
// @Success      201       {object}  domain.Customer
// @Security     Bearer
// @Router       /api/v1/customers [post]
func (h *OrderHandler) CreateCustomer(c *gin.Context) {
	var customer domain.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.customerRepo.Create(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// ListCustomers godoc
// @Summary      Listar clientes
// @Description  Obtiene catálogo de clientes activos
// @Tags         customers
// @Produce      json
// @Success      200  {array}   domain.Customer
// @Security     Bearer
// @Router       /api/v1/customers [get]
func (h *OrderHandler) ListCustomers(c *gin.Context) {
	customers, err := h.customerRepo.List(make(map[string]interface{}), 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

// GetStuckOrders godoc
// @Summary      Pedidos atorados (HU-24)
// @Description  Obtiene pedidos que llevan más de X horas sin avanzar
// @Tags         orders
// @Produce      json
// @Param        hours  query     int  false  "Horas mínimas (default: 4)"
// @Success      200    {array}   domain.Order
// @Security     Bearer
// @Router       /api/v1/orders/stuck [get]
func (h *OrderHandler) GetStuckOrders(c *gin.Context) {
	hours := 4 // Default: 4 horas
	if hoursParam := c.Query("hours"); hoursParam != "" {
		fmt.Sscanf(hoursParam, "%d", &hours)
	}

	orders, err := h.orderRepo.FindStuckOrders(hours, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
