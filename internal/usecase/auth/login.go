package auth

import (
	"time"

	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/infrastructure/security"
)

type LoginUseCase struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
	auditRepo   domain.AuditRepository
	secretKey   string
	jwtExpHours int
}

func NewLoginUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	auditRepo domain.AuditRepository,
	secretKey string,
	jwtExpHours int,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		auditRepo:   auditRepo,
		secretKey:   secretKey,
		jwtExpHours: jwtExpHours,
	}
}

type LoginInput struct {
	Username  string
	Password  string
	IPAddress string
	UserAgent string
}

type LoginOutput struct {
	Token     string       `json:"token"`
	User      *domain.User `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

func (uc *LoginUseCase) Execute(input LoginInput) (*LoginOutput, error) {
	// 1. Buscar usuario
	user, err := uc.userRepo.FindByUsername(input.Username)
	if err != nil {
		_ = uc.auditRepo.Log(domain.AuditLog{
			Action:    "LOGIN_FAILED",
			IPAddress: input.IPAddress,
			UserAgent: input.UserAgent,
		})
		return nil, domain.ErrInvalidCredentials
	}

	// 2. Verificar si está bloqueado (HU-00: 3 intentos)
	if user.IsLocked() {
		_ = uc.auditRepo.Log(domain.AuditLog{
			UserID:    &user.ID,
			Action:    "LOGIN_BLOCKED",
			IPAddress: input.IPAddress,
		})
		return nil, domain.ErrAccountLocked
	}

	// 3. Verificar contraseña
	if !security.CheckPasswordHash(input.Password, user.PasswordHash) {
		// Incrementar intentos fallidos
		user.IncrementFailedAttempts()
		_ = uc.userRepo.Update(user)
		_ = uc.auditRepo.Log(domain.AuditLog{
			UserID:    &user.ID,
			Action:    "LOGIN_FAILED",
			IPAddress: input.IPAddress,
		})
		return nil, domain.ErrInvalidCredentials
	}

	// 4. Reset intentos y actualizar last_login
	user.ResetFailedAttempts()
	now := time.Now()
	user.LastLogin = &now
	_ = uc.userRepo.Update(user)

	// 5. Generar JWT token
	token, err := security.GenerateToken(user.ID, user.Username, string(user.Role), uc.secretKey, uc.jwtExpHours)
	if err != nil {
		return nil, err
	}

	// 6. Crear sesión
	expiresAt := time.Now().Add(time.Duration(uc.jwtExpHours) * time.Hour)
	session := &domain.Session{
		UserID:    user.ID,
		Token:     token,
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
		ExpiresAt: expiresAt,
	}
	_ = uc.sessionRepo.Create(session)

	// 7. Auditar login exitoso
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:    &user.ID,
		Action:    "LOGIN_SUCCESS",
		IPAddress: input.IPAddress,
		UserAgent: input.UserAgent,
	})

	return &LoginOutput{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}
