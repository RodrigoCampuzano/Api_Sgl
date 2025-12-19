package domain

import "errors"

var (
	// Errores generales
	ErrNotFound       = errors.New("recurso no encontrado")
	ErrAlreadyExists  = errors.New("el recurso ya existe")
	ErrInvalidInput   = errors.New("datos de entrada inválidos")
	ErrUnauthorized   = errors.New("no autorizado")
	ErrForbidden      = errors.New("acceso denegado")
	ErrInternalServer = errors.New("error interno del servidor")

	// Errores de autenticación
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrAccountLocked      = errors.New("cuenta bloqueada por intentos fallidos")
	ErrSessionExpired     = errors.New("sesión expirada")
	ErrInvalidToken       = errors.New("token inválido")

	// Errores de usuarios
	ErrUserNotFound   = errors.New("usuario no encontrado")
	ErrUsernameExists = errors.New("nombre de usuario ya existe")
	ErrEmailExists    = errors.New("email ya existe")

	// Errores de inventario
	ErrInsufficientStock = errors.New("stock insuficiente")
	ErrExpiredProduct    = errors.New("producto caducado")
	ErrInvalidLot        = errors.New("lote inválido")

	// Errores de pedidos
	ErrOrderNotFound         = errors.New("pedido no encontrado")
	ErrOrderAlreadyProcessed = errors.New("pedido ya procesado")
	ErrInvalidOrderStatus    = errors.New("estado de pedido inválido")

	// Errores de flota
	ErrVehicleNotAvailable = errors.New("vehículo no disponible")
	ErrDriverNotAvailable  = errors.New("chofer no disponible")
)
