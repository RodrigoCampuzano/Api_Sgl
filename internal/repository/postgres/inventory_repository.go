package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type InventoryRepositoryPostgres struct {
	db *sqlx.DB
}

func NewInventoryRepository(db *sqlx.DB) domain.InventoryRepository {
	return &InventoryRepositoryPostgres{db: db}
}

func (r *InventoryRepositoryPostgres) Create(inventory *domain.Inventory) error {
	query := `
		INSERT INTO inventory (product_id, lot_number, expiration_date, quantity, status, warehouse_location)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, inventory.ProductID, inventory.LotNumber, inventory.ExpirationDate,
		inventory.Quantity, inventory.Status, inventory.WarehouseLocation).
		Scan(&inventory.ID, &inventory.CreatedAt, &inventory.UpdatedAt)
}

func (r *InventoryRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Inventory, error) {
	var inventory domain.Inventory
	query := `SELECT * FROM inventory WHERE id = $1`
	err := r.db.Get(&inventory, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &inventory, nil
}

func (r *InventoryRepositoryPostgres) FindByProduct(productID uuid.UUID) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	query := `SELECT * FROM inventory WHERE product_id = $1 ORDER BY expiration_date`
	err := r.db.Select(&inventories, query, productID)
	return inventories, err
}

// FindByProductFEFO implementa HU-06: First Expired First Out
func (r *InventoryRepositoryPostgres) FindByProductFEFO(productID uuid.UUID) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	query := `
		SELECT * FROM inventory 
		WHERE product_id = $1 
		  AND status = 'DISPONIBLE' 
		  AND quantity > 0
		ORDER BY expiration_date ASC NULLS LAST, created_at ASC
	`
	err := r.db.Select(&inventories, query, productID)
	return inventories, err
}

func (r *InventoryRepositoryPostgres) Update(inventory *domain.Inventory) error {
	now := time.Now()
	inventory.LastMovementAt = &now

	query := `
		UPDATE inventory
		SET quantity = $1, status = $2, last_movement_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`
	result, err := r.db.Exec(query, inventory.Quantity, inventory.Status, inventory.LastMovementAt, inventory.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ListAvailable implementa HU-05: Monitor de stock
func (r *InventoryRepositoryPostgres) ListAvailable(filters map[string]interface{}, limit, offset int) ([]*domain.Inventory, error) {
	var inventories []*domain.Inventory
	query := `
		SELECT i.*, p.name as product_name, p.brand, p.category
		FROM inventory i
		JOIN products p ON i.product_id = p.id
		WHERE i.status IN ('DISPONIBLE', 'RESERVADO')
		  AND i.quantity > 0
		ORDER BY i.created_at DESC
		LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&inventories, query, limit, offset)
	return inventories, err
}

func (r *InventoryRepositoryPostgres) GetStockByProduct(productID uuid.UUID) (int, error) {
	var totalStock int
	query := `
		SELECT COALESCE(SUM(quantity), 0)
		FROM inventory
		WHERE product_id = $1 AND status = 'DISPONIBLE'
	`
	err := r.db.Get(&totalStock, query, productID)
	return totalStock, err
}

// InventoryMovementRepositoryPostgres implementa el repositorio de movimientos
type InventoryMovementRepositoryPostgres struct {
	db *sqlx.DB
}

func NewInventoryMovementRepository(db *sqlx.DB) domain.InventoryMovementRepository {
	return &InventoryMovementRepositoryPostgres{db: db}
}

func (r *InventoryMovementRepositoryPostgres) Create(movement *domain.InventoryMovement) error {
	query := `
		INSERT INTO inventory_movements 
		(inventory_id, movement_type, quantity, previous_quantity, new_quantity, reference_id, 
		 reference_type, reason, evidence_photo_url, performed_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, movement.InventoryID, movement.MovementType, movement.Quantity,
		movement.PreviousQuantity, movement.NewQuantity, movement.ReferenceID, movement.ReferenceType,
		movement.Reason, movement.EvidencePhotoURL, movement.PerformedBy).
		Scan(&movement.ID, &movement.CreatedAt)
}

func (r *InventoryMovementRepositoryPostgres) FindByInventoryID(inventoryID uuid.UUID, limit, offset int) ([]*domain.InventoryMovement, error) {
	var movements []*domain.InventoryMovement
	query := `
		SELECT * FROM inventory_movements
		WHERE inventory_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.Select(&movements, query, inventoryID, limit, offset)
	return movements, err
}

func (r *InventoryMovementRepositoryPostgres) FindByDateRange(from, to time.Time, limit, offset int) ([]*domain.InventoryMovement, error) {
	var movements []*domain.InventoryMovement
	query := `
		SELECT * FROM inventory_movements
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	err := r.db.Select(&movements, query, from, to, limit, offset)
	return movements, err
}

// CycleCountRepositoryPostgres implementa el repositorio de conteo cíclico
type CycleCountRepositoryPostgres struct {
	db *sqlx.DB
}

func NewCycleCountRepository(db *sqlx.DB) domain.CycleCountRepository {
	return &CycleCountRepositoryPostgres{db: db}
}

func (r *CycleCountRepositoryPostgres) Create(count *domain.CycleCount) error {
	query := `
		INSERT INTO cycle_counts (scheduled_date, location, product_id, expected_quantity, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, count.ScheduledDate, count.Location, count.ProductID,
		count.ExpectedQuantity, count.Status).Scan(&count.ID, &count.CreatedAt)
}

func (r *CycleCountRepositoryPostgres) FindByID(id uuid.UUID) (*domain.CycleCount, error) {
	var count domain.CycleCount
	query := `SELECT * FROM cycle_counts WHERE id = $1`
	err := r.db.Get(&count, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &count, nil
}

func (r *CycleCountRepositoryPostgres) Update(count *domain.CycleCount) error {
	query := `
		UPDATE cycle_counts
		SET counted_quantity = $1, counted_by = $2, counted_at = $3, status = $4
		WHERE id = $5
	`
	result, err := r.db.Exec(query, count.CountedQuantity, count.CountedBy, count.CountedAt,
		count.Status, count.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CycleCountRepositoryPostgres) ListPending(limit, offset int) ([]*domain.CycleCount, error) {
	var counts []*domain.CycleCount
	query := `
		SELECT * FROM cycle_counts
		WHERE status = 'PENDIENTE'
		ORDER BY scheduled_date
		LIMIT $1 OFFSET $2
	`
	err := r.db.Select(&counts, query, limit, offset)
	return counts, err
}
