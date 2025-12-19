package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type SessionRepositoryPostgres struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) domain.SessionRepository {
	return &SessionRepositoryPostgres{db: db}
}

func (r *SessionRepositoryPostgres) Create(session *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, ip_address, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, session.UserID, session.Token, session.IPAddress,
		session.UserAgent, session.ExpiresAt).Scan(&session.ID, &session.CreatedAt)
}

func (r *SessionRepositoryPostgres) FindByToken(token string) (*domain.Session, error) {
	var session domain.Session
	query := `SELECT * FROM sessions WHERE token = $1`
	err := r.db.Get(&session, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSessionExpired
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepositoryPostgres) Delete(token string) error {
	query := `DELETE FROM sessions WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

func (r *SessionRepositoryPostgres) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query)
	return err
}

func (r *SessionRepositoryPostgres) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

// AuditRepositoryPostgres implementa el repositorio de auditoría
type AuditRepositoryPostgres struct {
	db *sqlx.DB
}

func NewAuditRepository(db *sqlx.DB) domain.AuditRepository {
	return &AuditRepositoryPostgres{db: db}
}

func (r *AuditRepositoryPostgres) Log(log domain.AuditLog) error {
	oldValuesJSON, _ := json.Marshal(log.OldValues)
	newValuesJSON, _ := json.Marshal(log.NewValues)

	query := `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, log.UserID, log.Action, log.EntityType, log.EntityID,
		oldValuesJSON, newValuesJSON, log.IPAddress, log.UserAgent)
	return err
}

func (r *AuditRepositoryPostgres) FindByUserID(userID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `SELECT * FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.Select(&logs, query, userID, limit, offset)
	return logs, err
}

func (r *AuditRepositoryPostgres) FindByEntity(entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `SELECT * FROM audit_logs WHERE entity_type = $1 AND entity_id = $2 ORDER BY created_at DESC`
	err := r.db.Select(&logs, query, entityType, entityID)
	return logs, err
}

func (r *AuditRepositoryPostgres) FindByDateRange(from, to time.Time, limit, offset int) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	query := `SELECT * FROM audit_logs WHERE created_at BETWEEN $1 AND $2 ORDER BY created_at DESC LIMIT $3 OFFSET $4`
	err := r.db.Select(&logs, query, from, to, limit, offset)
	return logs, err
}
