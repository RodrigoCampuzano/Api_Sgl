package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type UserRepositoryPostgres struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Role, user.Status).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepositoryPostgres) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL`
	err := r.db.Get(&user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryPostgres) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, role = $3, status = $4,
		    failed_login_attempts = $5, last_login = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, user.Email, user.PasswordHash, user.Role, user.Status,
		user.FailedLoginAttempts, user.LastLogin, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryPostgres) Delete(id uuid.UUID) error {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT * FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&users, query, limit, offset)
	return users, err
}

func (r *UserRepositoryPostgres) Count(filters map[string]interface{}) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	err := r.db.Get(&count, query)
	return count, err
}
