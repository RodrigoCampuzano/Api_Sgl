package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type OrderRepositoryPostgres struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) domain.OrderRepository {
	return &OrderRepositoryPostgres{db: db}
}

func (r *OrderRepositoryPostgres) Create(order *domain.Order) error {
	query := `
		INSERT INTO orders (order_number, customer_id, status, total_weight_kg, total_volume_m3, 
		                   total_cost, suggested_vehicle, has_fragile_items, has_heavy_items, 
		                   loading_alert, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, order.OrderNumber, order.CustomerID, order.Status,
		order.TotalWeightKg, order.TotalVolumeM3, order.TotalCost, order.SuggestedVehicle,
		order.HasFragileItems, order.HasHeavyItems, order.LoadingAlert, order.CreatedBy).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *OrderRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	query := `SELECT * FROM orders WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryPostgres) FindByOrderNumber(orderNumber string) (*domain.Order, error) {
	var order domain.Order
	query := `SELECT * FROM orders WHERE order_number = $1 AND deleted_at IS NULL`
	err := r.db.Get(&order, query, orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepositoryPostgres) Update(order *domain.Order) error {
	query := `
		UPDATE orders
		SET status = $1, total_weight_kg = $2, total_volume_m3 = $3, total_cost = $4,
		    suggested_vehicle = $5, has_fragile_items = $6, has_heavy_items = $7,
		    loading_alert = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, order.Status, order.TotalWeightKg, order.TotalVolumeM3,
		order.TotalCost, order.SuggestedVehicle, order.HasFragileItems, order.HasHeavyItems,
		order.LoadingAlert, order.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *OrderRepositoryPostgres) Delete(id uuid.UUID) error {
	query := `UPDATE orders SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *OrderRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	query := `SELECT * FROM orders WHERE deleted_at IS NULL`

	if status, ok := filters["status"]; ok {
		query += fmt.Sprintf(" AND status = '%s'", status)
	}

	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	err := r.db.Select(&orders, query, limit, offset)
	return orders, err
}

// FindStuckOrders implementa HU-24: Pedidos atorados > X horas
func (r *OrderRepositoryPostgres) FindStuckOrders(hours int, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	query := `
		SELECT * FROM orders
		WHERE status IN ('EN_PREPARACION', 'CONFIRMADO')
		  AND created_at < NOW() - INTERVAL '%d hours'
		  AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`
	query = fmt.Sprintf(query, hours)
	err := r.db.Select(&orders, query, limit, offset)
	return orders, err
}

// OrderLineRepositoryPostgres implementa el repositorio de líneas de pedido
type OrderLineRepositoryPostgres struct {
	db *sqlx.DB
}

func NewOrderLineRepository(db *sqlx.DB) domain.OrderLineRepository {
	return &OrderLineRepositoryPostgres{db: db}
}

func (r *OrderLineRepositoryPostgres) Create(line *domain.OrderLine) error {
	query := `
		INSERT INTO order_lines (order_id, product_id, inventory_id, quantity, unit_price, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, line.OrderID, line.ProductID, line.InventoryID,
		line.Quantity, line.UnitPrice, line.Subtotal).Scan(&line.ID, &line.CreatedAt)
}

func (r *OrderLineRepositoryPostgres) CreateBatch(lines []*domain.OrderLine) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO order_lines (order_id, product_id, inventory_id, quantity, unit_price, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	for _, line := range lines {
		err := tx.QueryRow(query, line.OrderID, line.ProductID, line.InventoryID,
			line.Quantity, line.UnitPrice, line.Subtotal).Scan(&line.ID, &line.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderLineRepositoryPostgres) FindByOrderID(orderID uuid.UUID) ([]*domain.OrderLine, error) {
	var lines []*domain.OrderLine
	query := `SELECT * FROM order_lines WHERE order_id = $1 ORDER BY created_at`
	err := r.db.Select(&lines, query, orderID)
	return lines, err
}

// CustomerRepositoryPostgres implementa el repositorio de clientes
type CustomerRepositoryPostgres struct {
	db *sqlx.DB
}

func NewCustomerRepository(db *sqlx.DB) domain.CustomerRepository {
	return &CustomerRepositoryPostgres{db: db}
}

func (r *CustomerRepositoryPostgres) Create(customer *domain.Customer) error {
	query := `
		INSERT INTO customers (name, rfc, address, city, state, postal_code, phone, email, credit_limit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, customer.Name, customer.RFC, customer.Address, customer.City,
		customer.State, customer.PostalCode, customer.Phone, customer.Email, customer.CreditLimit).
		Scan(&customer.ID, &customer.CreatedAt, &customer.UpdatedAt)
}

func (r *CustomerRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	query := `SELECT * FROM customers WHERE id = $1`
	err := r.db.Get(&customer, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepositoryPostgres) Update(customer *domain.Customer) error {
	query := `
		UPDATE customers
		SET name = $1, address = $2, phone = $3, email = $4, credit_limit = $5,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`
	result, err := r.db.Exec(query, customer.Name, customer.Address, customer.Phone,
		customer.Email, customer.CreditLimit, customer.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CustomerRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	query := `SELECT * FROM customers ORDER BY name LIMIT $1 OFFSET $2`
	err := r.db.Select(&customers, query, limit, offset)
	return customers, err
}
