package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sgl-disasur/api/internal/domain"
)

type ProductRepositoryPostgres struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &ProductRepositoryPostgres{db: db}
}

func (r *ProductRepositoryPostgres) Create(product *domain.Product) error {
	query := `
		INSERT INTO products (sku, name, brand, category, barcode, weight_kg, length_cm, width_cm, height_cm, is_fragile, unit_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, product.SKU, product.Name, product.Brand, product.Category, product.Barcode,
		product.WeightKg, product.LengthCm, product.WidthCm, product.HeightCm, product.IsFragile, product.UnitPrice).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *ProductRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT * FROM products WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryPostgres) FindBySKU(sku string) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT * FROM products WHERE sku = $1 AND deleted_at IS NULL`
	err := r.db.Get(&product, query, sku)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryPostgres) FindByBarcode(barcode string) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT * FROM products WHERE barcode = $1 AND deleted_at IS NULL`
	err := r.db.Get(&product, query, barcode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryPostgres) Update(product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, brand = $2, category = $3, barcode = $4, weight_kg = $5,
		    length_cm = $6, width_cm = $7, height_cm = $8, is_fragile = $9, unit_price = $10,
		    is_active = $11, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12 AND deleted_at IS NULL
	`
	result, err := r.db.Exec(query, product.Name, product.Brand, product.Category, product.Barcode,
		product.WeightKg, product.LengthCm, product.WidthCm, product.HeightCm, product.IsFragile,
		product.UnitPrice, product.IsActive, product.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProductRepositoryPostgres) Delete(id uuid.UUID) error {
	query := `UPDATE products SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProductRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	query := `SELECT * FROM products WHERE deleted_at IS NULL`

	// Agregar filtros
	if brand, ok := filters["brand"]; ok {
		query += fmt.Sprintf(" AND brand = '%s'", brand)
	}
	if category, ok := filters["category"]; ok {
		query += fmt.Sprintf(" AND category = '%s'", category)
	}

	query += " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	err := r.db.Select(&products, query, limit, offset)
	return products, err
}

func (r *ProductRepositoryPostgres) Count(filters map[string]interface{}) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`
	err := r.db.Get(&count, query)
	return count, err
}

// SupplierRepositoryPostgres implementa el repositorio de proveedores
type SupplierRepositoryPostgres struct {
	db *sqlx.DB
}

func NewSupplierRepository(db *sqlx.DB) domain.SupplierRepository {
	return &SupplierRepositoryPostgres{db: db}
}

func (r *SupplierRepositoryPostgres) Create(supplier *domain.Supplier) error {
	query := `
		INSERT INTO suppliers (name, brand, rfc, contact_name, phone, email)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, supplier.Name, supplier.Brand, supplier.RFC, supplier.ContactName,
		supplier.Phone, supplier.Email).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)
}

func (r *SupplierRepositoryPostgres) FindByID(id uuid.UUID) (*domain.Supplier, error) {
	var supplier domain.Supplier
	query := `SELECT * FROM suppliers WHERE id = $1`
	err := r.db.Get(&supplier, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepositoryPostgres) Update(supplier *domain.Supplier) error {
	query := `
		UPDATE suppliers
		SET name = $1, brand = $2, rfc = $3, contact_name = $4, phone = $5, email = $6,
		    is_active = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`
	result, err := r.db.Exec(query, supplier.Name, supplier.Brand, supplier.RFC, supplier.ContactName,
		supplier.Phone, supplier.Email, supplier.IsActive, supplier.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *SupplierRepositoryPostgres) List(filters map[string]interface{}, limit, offset int) ([]*domain.Supplier, error) {
	var suppliers []*domain.Supplier
	query := `SELECT * FROM suppliers WHERE is_active = true ORDER BY name LIMIT $1 OFFSET $2`
	err := r.db.Select(&suppliers, query, limit, offset)
	return suppliers, err
}
