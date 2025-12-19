package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sgl-disasur/api/internal/usecase/auth"
)

type AuthHandler struct {
	loginUseCase    *auth.LoginUseCase
	registerUseCase *auth.RegisterUserUseCase
}

func NewAuthHandler(loginUC *auth.LoginUseCase, registerUC *auth.RegisterUserUseCase) *AuthHandler {
	return &AuthHandler{
		loginUseCase:    loginUC,
		registerUseCase: registerUC,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Login de usuario
// @Description  Autentica un usuario y retorna un JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Credenciales de login"
// @Success      200         {object}  auth.LoginOutput
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	input := auth.LoginInput{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	result, err := h.loginUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Logout godoc
// @Summary      Logout de usuario
// @Description  Cierra la sesión del usuario eliminando el token
// @Tags         auth
// @Security     Bearer
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada exitosamente"})
}

// Register godoc
// @Summary      Registrar usuario
// @Description  Crea un nuevo usuario en el sistema (público para crear admin inicial)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      auth.RegisterUserInput  true  "Datos del usuario"
// @Success      201   {object}  auth.RegisterUserOutput
// @Failure      400   {object}  map[string]string
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input auth.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	input.IPAddress = c.ClientIP()

	result, err := h.registerUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
