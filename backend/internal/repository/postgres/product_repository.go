package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/products"
)

// ProductRepository implements products.Repository
type ProductRepository struct {
	db *DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// List retrieves products with filters
func (r *ProductRepository) List(ctx context.Context, orgID uuid.UUID, filters products.ProductFilters) ([]products.Product, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	// Build query
	query := `
		SELECT
			id, organization_id, name, sku, barcode, description,
			category_id, unit_price, cost, tax_rate, unit,
			current_stock, min_stock_level, max_stock_level,
			is_active, is_track_inventory, allow_negative_stock,
			image_url, product_type, accounting_account_id, tax_code_id,
			created_at, updated_at, created_by, updated_by
		FROM products
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	// Add ordering
	query += " ORDER BY name ASC"

	// Add pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var productList []products.Product
	for rows.Next() {
		var p products.Product
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Name, &p.SKU, &p.Barcode, &p.Description,
			&p.CategoryID, &p.UnitPrice, &p.Cost, &p.TaxRate, &p.Unit,
			&p.CurrentStock, &p.MinStockLevel, &p.MaxStockLevel,
			&p.IsActive, &p.IsTrackInventory, &p.AllowNegativeStock,
			&p.ImageURL, &p.ProductType, &p.AccountingAccountID, &p.TaxCodeID,
			&p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		productList = append(productList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productList, nil
}

// Count counts products matching filters
func (r *ProductRepository) Count(ctx context.Context, orgID uuid.UUID, filters products.ProductFilters) (int64, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM products WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	// Apply filters (same as List)
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR sku ILIKE $%d OR barcode ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND category_id = $%d", argCount)
		args = append(args, *filters.CategoryID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *products.Product) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, product.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO products (
			id, organization_id, name, sku, barcode, description,
			category_id, unit_price, cost, tax_rate, unit,
			current_stock, min_stock_level, max_stock_level,
			is_active, is_track_inventory, allow_negative_stock,
			image_url, product_type, accounting_account_id, tax_code_id,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21,
			$22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		product.ID, product.OrganizationID, product.Name, product.SKU, product.Barcode, product.Description,
		product.CategoryID, product.UnitPrice, product.Cost, product.TaxRate, product.Unit,
		product.CurrentStock, product.MinStockLevel, product.MaxStockLevel,
		product.IsActive, product.IsTrackInventory, product.AllowNegativeStock,
		product.ImageURL, product.ProductType, product.AccountingAccountID, product.TaxCodeID,
		product.CreatedAt, product.UpdatedAt, product.CreatedBy,
	)

	return err
}

// Get retrieves a product by ID
func (r *ProductRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*products.Product, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, name, sku, barcode, description,
			category_id, unit_price, cost, tax_rate, unit,
			current_stock, min_stock_level, max_stock_level,
			is_active, is_track_inventory, allow_negative_stock,
			image_url, product_type, accounting_account_id, tax_code_id,
			created_at, updated_at, created_by, updated_by
		FROM products
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var p products.Product
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&p.ID, &p.OrganizationID, &p.Name, &p.SKU, &p.Barcode, &p.Description,
		&p.CategoryID, &p.UnitPrice, &p.Cost, &p.TaxRate, &p.Unit,
		&p.CurrentStock, &p.MinStockLevel, &p.MaxStockLevel,
		&p.IsActive, &p.IsTrackInventory, &p.AllowNegativeStock,
		&p.ImageURL, &p.ProductType, &p.AccountingAccountID, &p.TaxCodeID,
		&p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(ctx context.Context, orgID uuid.UUID, sku string) (*products.Product, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, name, sku, barcode, description,
			category_id, unit_price, cost, tax_rate, unit,
			current_stock, min_stock_level, max_stock_level,
			is_active, is_track_inventory, allow_negative_stock,
			image_url, product_type, accounting_account_id, tax_code_id,
			created_at, updated_at, created_by, updated_by
		FROM products
		WHERE organization_id = $1 AND LOWER(sku) = LOWER($2) AND deleted_at IS NULL
	`

	var p products.Product
	err := r.db.Pool.QueryRow(ctx, query, orgID, sku).Scan(
		&p.ID, &p.OrganizationID, &p.Name, &p.SKU, &p.Barcode, &p.Description,
		&p.CategoryID, &p.UnitPrice, &p.Cost, &p.TaxRate, &p.Unit,
		&p.CurrentStock, &p.MinStockLevel, &p.MaxStockLevel,
		&p.IsActive, &p.IsTrackInventory, &p.AllowNegativeStock,
		&p.ImageURL, &p.ProductType, &p.AccountingAccountID, &p.TaxCodeID,
		&p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(ctx context.Context, product *products.Product) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, product.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE products SET
			name = $3, sku = $4, barcode = $5, description = $6,
			category_id = $7, unit_price = $8, cost = $9, tax_rate = $10, unit = $11,
			min_stock_level = $12, max_stock_level = $13,
			is_active = $14, is_track_inventory = $15, allow_negative_stock = $16,
			image_url = $17, product_type = $18, accounting_account_id = $19, tax_code_id = $20,
			updated_at = $21, updated_by = $22
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		product.OrganizationID, product.ID,
		product.Name, product.SKU, product.Barcode, product.Description,
		product.CategoryID, product.UnitPrice, product.Cost, product.TaxRate, product.Unit,
		product.MinStockLevel, product.MaxStockLevel,
		product.IsActive, product.IsTrackInventory, product.AllowNegativeStock,
		product.ImageURL, product.ProductType, product.AccountingAccountID, product.TaxCodeID,
		product.UpdatedAt, product.UpdatedBy,
	)

	return err
}

// Delete soft-deletes a product
func (r *ProductRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE products
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// UpdateStock updates product stock quantity
func (r *ProductRepository) UpdateStock(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, quantity float64) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE products
		SET
			current_stock = current_stock + $3,
			updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, productID, quantity, time.Now())
	return err
}
