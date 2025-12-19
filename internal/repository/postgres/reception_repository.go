package postgres

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type ReceptionOrderRepositoryPostgres struct {
	db *sqlx.DB
}

func NewReceptionOrderRepository(db *sqlx.DB) domain.ReceptionOrderRepository {
	return &ReceptionOrderRepositoryPostgres{db: db}
}

func (r *ReceptionOrderRepositoryPostgres) Create(order *domain.ReceptionOrder) error {
	query := `
		INSERT INTO reception_orders (order_number, supplier_id, brand, invoice_number, invoice_file_url, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, order.OrderNumber, order.SupplierID, order.Brand, order.InvoiceNumber,
		order.InvoiceFileURL, order.Status, order.Notes).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *ReceptionOrderRepositoryPostgres) FindByID(id uuid.UUID) (*domain.ReceptionOrder, error) {
	var order domain.ReceptionOrder
	query := `SELECT * FROM reception_orders WHERE id = $1`
	err := r.db.Get(&order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *ReceptionOrderRepositoryPostgres) FindByOrderNumber(orderNumber string) (*domain.ReceptionOrder, error) {
	var order domain.ReceptionOrder
	query := `SELECT * FROM reception_orders WHERE order_number = $1`
	err := r.db.Get(&order, query, orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *ReceptionOrderRepositoryPostgres) Update(order *domain.ReceptionOrder) error {
	query := `
		UPDATE reception_orders
		SET status = $1, received_by = $2, received_at = $3, validated_by = $4,
		    validated_at = $5, notes = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`
	result, err := r.db.Exec(query, order.Status, order.ReceivedBy, order.ReceivedAt,
		order.ValidatedBy, order.ValidatedAt, order.Notes, order.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ReceptionOrderRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.ReceptionOrder, error) {
	var orders []*domain.ReceptionOrder
	query := `SELECT * FROM reception_orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.Select(&orders, query, limit, offset)
	return orders, err
}

// ReceptionLineRepositoryPostgres implementa el repositorio de líneas de recepción
type ReceptionLineRepositoryPostgres struct {
	db *sqlx.DB
}

func NewReceptionLineRepository(db *sqlx.DB) domain.ReceptionLineRepository {
	return &ReceptionLineRepositoryPostgres{db: db}
}

func (r *ReceptionLineRepositoryPostgres) Create(line *domain.ReceptionLine) error {
	query := `
		INSERT INTO reception_lines (reception_order_id, product_id, expected_quantity, lot_number, expiration_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, line.ReceptionOrderID, line.ProductID, line.ExpectedQuantity,
		line.LotNumber, line.ExpirationDate).Scan(&line.ID, &line.CreatedAt, &line.UpdatedAt)
}

func (r *ReceptionLineRepositoryPostgres) CreateBatch(lines []*domain.ReceptionLine) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO reception_lines (reception_order_id, product_id, expected_quantity, lot_number, expiration_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	for _, line := range lines {
		err := tx.QueryRow(query, line.ReceptionOrderID, line.ProductID, line.ExpectedQuantity,
			line.LotNumber, line.ExpirationDate).Scan(&line.ID, &line.CreatedAt, &line.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReceptionLineRepositoryPostgres) FindByID(id uuid.UUID) (*domain.ReceptionLine, error) {
	var line domain.ReceptionLine
	query := `SELECT * FROM reception_lines WHERE id = $1`
	err := r.db.Get(&line, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &line, nil
}

func (r *ReceptionLineRepositoryPostgres) FindByOrderID(orderID uuid.UUID) ([]*domain.ReceptionLine, error) {
	var lines []*domain.ReceptionLine
	query := `SELECT * FROM reception_lines WHERE reception_order_id = $1 ORDER BY created_at`
	err := r.db.Select(&lines, query, orderID)
	return lines, err
}

func (r *ReceptionLineRepositoryPostgres) Update(line *domain.ReceptionLine) error {
	query := `
		UPDATE reception_lines
		SET counted_quantity = $1, condition = $2, counted_by = $3, counted_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	result, err := r.db.Exec(query, line.CountedQuantity, line.Condition, line.CountedBy,
		line.CountedAt, line.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ReceptionDiscrepancyRepositoryPostgres implementa el repositorio de discrepancias
type ReceptionDiscrepancyRepositoryPostgres struct {
	db *sqlx.DB
}

func NewReceptionDiscrepancyRepository(db *sqlx.DB) domain.ReceptionDiscrepancyRepository {
	return &ReceptionDiscrepancyRepositoryPostgres{db: db}
}

func (r *ReceptionDiscrepancyRepositoryPostgres) Create(discrepancy *domain.ReceptionDiscrepancy) error {
	query := `
		INSERT INTO reception_discrepancies (reception_line_id, expected_qty, counted_qty, difference, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, discrepancy.ReceptionLineID, discrepancy.ExpectedQty,
		discrepancy.CountedQty, discrepancy.Difference, discrepancy.Status).
		Scan(&discrepancy.ID, &discrepancy.CreatedAt, &discrepancy.UpdatedAt)
}

func (r *ReceptionDiscrepancyRepositoryPostgres) FindByLineID(lineID uuid.UUID) (*domain.ReceptionDiscrepancy, error) {
	var discrepancy domain.ReceptionDiscrepancy
	query := `SELECT * FROM reception_discrepancies WHERE reception_line_id = $1`
	err := r.db.Get(&discrepancy, query, lineID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &discrepancy, nil
}

func (r *ReceptionDiscrepancyRepositoryPostgres) Update(discrepancy *domain.ReceptionDiscrepancy) error {
	query := `
		UPDATE reception_discrepancies
		SET status = $1, resolution_notes = $2, resolved_by = $3, resolved_at = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	result, err := r.db.Exec(query, discrepancy.Status, discrepancy.ResolutionNotes,
		discrepancy.ResolvedBy, discrepancy.ResolvedAt, discrepancy.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ReceptionDiscrepancyRepositoryPostgres) ListPending(limit, offset int) ([]*domain.ReceptionDiscrepancy, error) {
	var discrepancies []*domain.ReceptionDiscrepancy
	query := `
		SELECT * FROM reception_discrepancies 
		WHERE status IN ('DETECTADA', 'EN_REVISION') 
		ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&discrepancies, query, limit, offset)
	return discrepancies, err
}
