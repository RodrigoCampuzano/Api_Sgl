package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole define los roles de usuario en el sistema
type UserRole string

const (
	RoleAdminTI         UserRole = "ADMIN_TI"
	RoleGerente         UserRole = "GERENTE"
	RoleJefeAlmacen     UserRole = "JEFE_ALMACEN"
	RoleAuxiliar        UserRole = "AUXILIAR"
	RoleSupervisor      UserRole = "SUPERVISOR"
	RoleRecepcionista   UserRole = "RECEPCIONISTA"
	RoleVendedor        UserRole = "VENDEDOR"
	RoleJefeTrafico     UserRole = "JEFE_TRAFICO"
	RoleChofer          UserRole = "CHOFER"
	RoleMontacarguista  UserRole = "MONTACARGUISTA"
	RoleCargador        UserRole = "CARGADOR"
	RolePlanificador    UserRole = "PLANIFICADOR"
	RoleFlota           UserRole = "FLOTA"
	RoleAuditor         UserRole = "AUDITOR"
	RoleServicioCliente UserRole = "SERVICIO_CLIENTE"
)

// UserStatus define el estado de un usuario
type UserStatus string

const (
	UserStatusActivo    UserStatus = "ACTIVO"
	UserStatusBloqueado UserStatus = "BLOQUEADO"
	UserStatusInactivo  UserStatus = "INACTIVO"
)

// User representa un usuario del sistema
type User struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	Username            string     `json:"username" db:"username"`
	Email               string     `json:"email" db:"email"`
	PasswordHash        string     `json:"-" db:"password_hash"`
	Role                UserRole   `json:"role" db:"role"`
	Status              UserStatus `json:"status" db:"status"`
	FailedLoginAttempts int        `json:"failed_login_attempts" db:"failed_login_attempts"`
	LastLogin           *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt           *time.Time `json:"-" db:"deleted_at"`
}

// IsActive verifica si el usuario está activo
func (u *User) IsActive() bool {
	return u.Status == UserStatusActivo && u.DeletedAt == nil
}

// IsLocked verifica si el usuario está bloqueado
func (u *User) IsLocked() bool {
	return u.Status == UserStatusBloqueado || u.FailedLoginAttempts >= 3
}

// IncrementFailedAttempts incrementa el contador de intentos fallidos
// y bloquea la cuenta si alcanza 3 intentos
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= 3 {
		u.Status = UserStatusBloqueado
	}
}

// ResetFailedAttempts resetea el contador de intentos fallidos
// y desbloquea la cuenta si estaba bloqueada por intentos
func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
	if u.Status == UserStatusBloqueado {
		u.Status = UserStatusActivo
	}
}

// Session representa una sesión activa de usuario
type Session struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// IsExpired verifica si la sesión ha expirado
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// AuditLog representa un registro de auditoría
type AuditLog struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	UserID     *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Action     string                 `json:"action" db:"action"`
	EntityType string                 `json:"entity_type,omitempty" db:"entity_type"`
	EntityID   *uuid.UUID             `json:"entity_id,omitempty" db:"entity_id"`
	OldValues  map[string]interface{} `json:"old_values,omitempty" db:"old_values"`
	NewValues  map[string]interface{} `json:"new_values,omitempty" db:"new_values"`
	IPAddress  string                 `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent  string                 `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}

// UserRepository define los métodos de repositorio para usuarios
type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(filters map[string]interface{}, limit, offset int) ([]*User, error)
	Count(filters map[string]interface{}) (int, error)
}

// SessionRepository define los métodos de repositorio para sesiones
type SessionRepository interface {
	Create(session *Session) error
	FindByToken(token string) (*Session, error)
	Delete(token string) error
	DeleteExpired() error
	DeleteByUserID(userID uuid.UUID) error
}

// AuditRepository define los métodos de repositorio para auditoría
type AuditRepository interface {
	Log(log AuditLog) error
	FindByUserID(userID uuid.UUID, limit, offset int) ([]*AuditLog, error)
	FindByEntity(entityType string, entityID uuid.UUID) ([]*AuditLog, error)
	FindByDateRange(from, to time.Time, limit, offset int) ([]*AuditLog, error)
}
