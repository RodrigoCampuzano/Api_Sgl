package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sgl-disasur/api/internal/domain"
)

type ProductHandler struct {
	productRepo domain.ProductRepository
}

func NewProductHandler(productRepo domain.ProductRepository) *ProductHandler {
	return &ProductHandler{
		productRepo: productRepo,
	}
}

// ListProducts godoc
// @Summary      Listar productos (HU-04)
// @Description  Obtiene el catálogo de productos con filtros
// @Tags         products
// @Produce      json
// @Param        brand     query     string  false  "Filtrar por marca"
// @Param        category  query     string  false  "Filtrar por categoría"
// @Success      200       {array}   domain.Product
// @Security     Bearer
// @Router       /api/v1/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	filters := make(map[string]interface{})

	if brand := c.Query("brand"); brand != "" {
		filters["brand"] = brand
	}
	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	products, err := h.productRepo.List(filters, 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary      Obtener producto por ID o SKU
// @Description  Obtiene los detalles de un producto
// @Tags         products
// @Produce      json
// @Param        id   path      string  true  "Product ID o SKU"
// @Success      200  {object}  domain.Product
// @Security     Bearer
// @Router       /api/v1/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")

	// Intentar buscar por UUID
	if id, err := uuid.Parse(idStr); err == nil {
		product, err := h.productRepo.FindByID(id)
		if err == nil {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	// Si no es UUID, buscar por SKU
	product, err := h.productRepo.FindBySKU(idStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary      Crear producto
// @Description  Crea un nuevo producto en el catálogo
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      domain.Product  true  "Datos del producto"
// @Success      201      {object}  domain.Product
// @Security     Bearer
// @Router       /api/v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.productRepo.Create(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}
