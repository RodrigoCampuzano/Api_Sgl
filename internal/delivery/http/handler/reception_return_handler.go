package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProcessReturn implementa HU-14: Registro de devoluciones
// ProcessReturn godoc
// @Summary      Procesar devolución (HU-14)
// @Description  Registra devolución clasificándola como APTA o DESECHO
// @Tags         reception
// @Accept       json
// @Produce      json
// @Param        return  body      map[string]interface{}  true  "Datos de devolución"
// @Success      200     {object}  map[string]interface{}
// @Security     Bearer
// @Router       /api/v1/reception/returns [post]
func (h *ReceptionHandler) ProcessReturn(c *gin.Context) {
	var input struct {
		ProductID uuid.UUID `json:"product_id" binding:"required"`
		Quantity  int       `json:"quantity" binding:"required,min=1"`
		Condition string    `json:"condition" binding:"required"` // APTA o DESECHO
		Reason    string    `json:"reason" binding:"required"`
		PhotoURL  string    `json:"photo_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Validar condición
	if input.Condition != "APTA" && input.Condition != "DESECHO" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Condición debe ser APTA o DESECHO"})
		return
	}

	// HU-14: Lógica de devoluciones
	message := ""
	location := ""

	if input.Condition == "APTA" {
		location = "CUARENTENA"
		message = fmt.Sprintf("Devolución procesada: %d unidades ingresadas a CUARENTENA", input.Quantity)
	} else {
		location = "DESECHO"
		message = fmt.Sprintf("Devolución procesada: %d unidades marcadas como DESECHO", input.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   message,
		"condition": input.Condition,
		"location":  location,
		"quantity":  input.Quantity,
	})
}
