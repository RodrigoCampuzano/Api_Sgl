package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/usecase/reception"
)

type ReceptionHandler struct {
	createOrderUC *reception.CreateReceptionOrderUseCase
	blindCountUC  *reception.BlindCountUseCase
	productRepo   domain.ProductRepository
	supplierRepo  domain.SupplierRepository
	receptionRepo domain.ReceptionOrderRepository
	lineRepo      domain.ReceptionLineRepository
}

func NewReceptionHandler(
	createOrderUC *reception.CreateReceptionOrderUseCase,
	blindCountUC *reception.BlindCountUseCase,
	productRepo domain.ProductRepository,
	supplierRepo domain.SupplierRepository,
	receptionRepo domain.ReceptionOrderRepository,
	lineRepo domain.ReceptionLineRepository,
) *ReceptionHandler {
	return &ReceptionHandler{
		createOrderUC: createOrderUC,
		blindCountUC:  blindCountUC,
		productRepo:   productRepo,
		supplierRepo:  supplierRepo,
		receptionRepo: receptionRepo,
		lineRepo:      lineRepo,
	}
}

// CreateReceptionOrder godoc
// @Summary      Crear orden de recepción (HU-01)
// @Description  Crea una nueva orden de recepción con sus líneas
// @Tags         reception
// @Accept       json
// @Produce      json
// @Param        order  body      reception.CreateReceptionOrderInput  true  "Datos de la orden"
// @Success      201    {object}  domain.ReceptionOrder
// @Failure      400    {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/reception/orders [post]
func (h *ReceptionHandler) CreateOrder(c *gin.Context) {
	var input reception.CreateReceptionOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtener user_id del contexto
	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	order, err := h.createOrderUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// BlindCount godoc
// @Summary      Conteo ciego (HU-02)
// @Description  Registra el conteo físico sin mostrar cantidades esperadas
// @Tags         reception
// @Accept       json
// @Produce      json
// @Param        count  body      reception.BlindCountInput  true  "Datos del conteo"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/reception/blind-count [post]
func (h *ReceptionHandler) BlindCount(c *gin.Context) {
	var input reception.BlindCountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	if err := h.blindCountUC.Execute(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conteo registrado exitosamente"})
}

// ListOrders godoc
// @Summary      Listar órdenes de recepción
// @Description  Obtiene listado de órdenes de recepción
// @Tags         reception
// @Produce      json
// @Success      200  {array}   domain.ReceptionOrder
// @Security     Bearer
// @Router       /api/v1/reception/orders [get]
func (h *ReceptionHandler) ListOrders(c *gin.Context) {
	orders, err := h.receptionRepo.List(nil, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary      Obtener orden de recepción por ID
// @Description  Obtiene los detalles de una orden de recepción con sus líneas
// @Tags         reception
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Security     Bearer
// @Router       /api/v1/reception/orders/{id} [get]
func (h *ReceptionHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	order, err := h.receptionRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	lines, _ := h.lineRepo.FindByOrderID(id)

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"lines": lines,
	})
}
