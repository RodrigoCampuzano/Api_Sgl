package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/usecase/inventory"
)

type InventoryHandler struct {
	getStockUC       *inventory.GetStockUseCase
	getFEFOLotsUC    *inventory.GetFEFOLotsUseCase
	registerDamageUC *inventory.RegisterDamageUseCase
	cycleCountUC     *inventory.PerformCycleCountUseCase
}

func NewInventoryHandler(
	getStockUC *inventory.GetStockUseCase,
	getFEFOLotsUC *inventory.GetFEFOLotsUseCase,
	registerDamageUC *inventory.RegisterDamageUseCase,
	cycleCountUC *inventory.PerformCycleCountUseCase,
) *InventoryHandler {
	return &InventoryHandler{
		getStockUC:       getStockUC,
		getFEFOLotsUC:    getFEFOLotsUC,
		registerDamageUC: registerDamageUC,
		cycleCountUC:     cycleCountUC,
	}
}

// GetStock godoc
// @Summary      Monitor de stock (HU-05)
// @Description  Obtiene inventario en tiempo real con alertas de punto de reorden
// @Tags         inventory
// @Produce      json
// @Param        brand     query     string  false  "Filtrar por marca"
// @Param        category  query     string  false  "Filtrar por categoría"
// @Success      200       {array}   inventory.StockItem
// @Security     Bearer
// @Router       /api/v1/inventory/stock [get]
func (h *InventoryHandler) GetStock(c *gin.Context) {
	var brand *domain.Brand
	var category *string

	if brandStr := c.Query("brand"); brandStr != "" {
		b := domain.Brand(brandStr)
		brand = &b
	}

	if cat := c.Query("category"); cat != "" {
		category = &cat
	}

	stock, err := h.getStockUC.Execute(brand, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetFEFOLots godoc
// @Summary      Lotes FEFO (HU-06)
// @Description  Obtiene lotes ordenados por First Expired First Out
// @Tags         inventory
// @Produce      json
// @Param        product_id  path      string  true  "Product ID"
// @Success      200         {array}   inventory.FEFOLot
// @Security     Bearer
// @Router       /api/v1/inventory/fefo/{product_id} [get]
func (h *InventoryHandler) GetFEFOLots(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de producto inválido"})
		return
	}

	lots, err := h.getFEFOLotsUC.Execute(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lots)
}

// RegisterDamage godoc
// @Summary      Registrar merma (HU-13)
// @Description  Registra daño/merma con foto de evidencia obligatoria
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        damage  body      inventory.RegisterDamageInput  true  "Datos de la merma"
// @Success      200     {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/inventory/damages [post]
func (h *InventoryHandler) RegisterDamage(c *gin.Context) {
	var input inventory.RegisterDamageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	output, err := h.registerDamageUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// GenerateCycleCounts godoc
// @Summary      Generar conteos cíclicos (HU-15)
// @Description  Genera 5 ubicaciones aleatorias para conteo diario
// @Tags         inventory
// @Produce      json
// @Success      200  {array}   domain.CycleCount
// @Security     Bearer
// @Router       /api/v1/inventory/cycle-counts/generate [post]
func (h *InventoryHandler) GenerateCycleCounts(c *gin.Context) {
	counts, err := h.cycleCountUC.GenerateDailyCycleCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// PerformCycleCount godoc
// @Summary      Realizar conteo cíclico (HU-15)
// @Description  Registra el conteo físico y ajusta inventario si hay varianza
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        count  body      inventory.PerformCountInput  true  "Datos del conteo"
// @Success      200    {object}  map[string]string
// @Security     Bearer
// @Router       /api/v1/inventory/cycle-counts/perform [post]
func (h *InventoryHandler) PerformCycleCount(c *gin.Context) {
	var input inventory.PerformCountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))
	input.UserID = userID

	if err := h.cycleCountUC.PerformCount(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conteo registrado exitosamente"})
}
