package auth

import (
	"errors"

	"github.com/sgl-disasur/api/internal/domain"
	"github.com/sgl-disasur/api/internal/infrastructure/security"
)

type RegisterUserUseCase struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditRepository
}

func NewRegisterUserUseCase(
	userRepo domain.UserRepository,
	auditRepo domain.AuditRepository,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:  userRepo,
		auditRepo: auditRepo,
	}
}

type RegisterUserInput struct {
	Username  string          `json:"username" binding:"required"`
	Email     string          `json:"email" binding:"required"`
	Password  string          `json:"password" binding:"required"`
	Role      domain.UserRole `json:"role" binding:"required"`
	IPAddress string          `json:"-"`
}

type RegisterUserOutput struct {
	User *domain.User `json:"user"`
}

func (uc *RegisterUserUseCase) Execute(input RegisterUserInput) (*RegisterUserOutput, error) {
	// 1. Verificar que el username no existe
	existingUser, _ := uc.userRepo.FindByUsername(input.Username)
	if existingUser != nil {
		return nil, domain.ErrUsernameExists
	}

	// 2. Verificar que el email no existe
	existingEmail, _ := uc.userRepo.FindByEmail(input.Email)
	if existingEmail != nil {
		return nil, domain.ErrEmailExists
	}

	// 3. Hashear la contraseña
	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("error al hashear contraseña")
	}

	// 4. Crear el usuario
	user := &domain.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         input.Role,
		Status:       domain.UserStatusActivo,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 5. Auditar
	_ = uc.auditRepo.Log(domain.AuditLog{
		UserID:     &user.ID,
		Action:     "USER_REGISTERED",
		EntityType: "USER",
		EntityID:   &user.ID,
		IPAddress:  input.IPAddress,
		NewValues: map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})

	return &RegisterUserOutput{
		User: user,
	}, nil
}
