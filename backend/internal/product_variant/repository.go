package product_variant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for ProductVariants
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ProductVariants repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ProductVariants represents a product_variants entity
type ProductVariants struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	VariantName string `json:"variant_name" db:"variant_name"`
	Sku *string `json:"sku" db:"sku"`
	Barcode *string `json:"barcode" db:"barcode"`
	Attributes json.RawMessage `json:"attributes" db:"attributes"`
	CostPrice *float64 `json:"cost_price" db:"cost_price"`
	SellingPrice *float64 `json:"selling_price" db:"selling_price"`
	CompareAtPrice *float64 `json:"compare_at_price" db:"compare_at_price"`
	CurrentStock *float64 `json:"current_stock" db:"current_stock"`
	ReorderLevel *float64 `json:"reorder_level" db:"reorder_level"`
	ReorderQuantity *float64 `json:"reorder_quantity" db:"reorder_quantity"`
	Weight *float64 `json:"weight" db:"weight"`
	WeightUnit *string `json:"weight_unit" db:"weight_unit"`
	Dimensions json.RawMessage `json:"dimensions" db:"dimensions"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	ImageUrl *string `json:"image_url" db:"image_url"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	(costPrice *string `json:"(cost_price" db:"(cost_price"`
	(sellingPrice *string `json:"(selling_price" db:"(selling_price"`
	(compareAtPrice *string `json:"(compare_at_price" db:"(compare_at_price"`
}

// Create inserts a new product_variants record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ProductVariants) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "product_variants", duration, nil)
	}()

	query := `
		INSERT INTO product_variants (
			, organization_id
			, product_id
			, variant_name
			, sku
			, barcode
			, attributes
			, cost_price
			, selling_price
			, compare_at_price
			, current_stock
			, reorder_level
			, reorder_quantity
			, weight
			, weight_unit
			, dimensions
			, is_active
			, is_default
			, sort_order
			, image_url
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, (cost_price
			, (selling_price
			, (compare_at_price
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.VariantName,
		entity.Sku,
		entity.Barcode,
		entity.Attributes,
		entity.CostPrice,
		entity.SellingPrice,
		entity.CompareAtPrice,
		entity.CurrentStock,
		entity.ReorderLevel,
		entity.ReorderQuantity,
		entity.Weight,
		entity.WeightUnit,
		entity.Dimensions,
		entity.IsActive,
		entity.IsDefault,
		entity.SortOrder,
		entity.ImageUrl,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(costPrice,
		entity.(sellingPrice,
		entity.(compareAtPrice,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create product_variants", zap.Error(err))
		return fmt.Errorf("failed to create product_variants: %w", err)
	}

	r.logger.Info("created product_variants",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a product_variants by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ProductVariants, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_variants", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, variant_name
			, sku
			, barcode
			, attributes
			, cost_price
			, selling_price
			, compare_at_price
			, current_stock
			, reorder_level
			, reorder_quantity
			, weight
			, weight_unit
			, dimensions
			, is_active
			, is_default
			, sort_order
			, image_url
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (cost_price
			, (selling_price
			, (compare_at_price
		FROM product_variants
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ProductVariants
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.VariantName,
		&entity.Sku,
		&entity.Barcode,
		&entity.Attributes,
		&entity.CostPrice,
		&entity.SellingPrice,
		&entity.CompareAtPrice,
		&entity.CurrentStock,
		&entity.ReorderLevel,
		&entity.ReorderQuantity,
		&entity.Weight,
		&entity.WeightUnit,
		&entity.Dimensions,
		&entity.IsActive,
		&entity.IsDefault,
		&entity.SortOrder,
		&entity.ImageUrl,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.(costPrice,
		&entity.(sellingPrice,
		&entity.(compareAtPrice,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("product_variants not found")
	}

	if err != nil {
		r.logger.Error("failed to get product_variants", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get product_variants: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of product_variants records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ProductVariants, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_variants", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_variants
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_variants records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, variant_name
			, sku
			, barcode
			, attributes
			, cost_price
			, selling_price
			, compare_at_price
			, current_stock
			, reorder_level
			, reorder_quantity
			, weight
			, weight_unit
			, dimensions
			, is_active
			, is_default
			, sort_order
			, image_url
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (cost_price
			, (selling_price
			, (compare_at_price
		FROM product_variants
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_variants", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_variants: %w", err)
	}
	defer rows.Close()

	var entities []*ProductVariants
	for rows.Next() {
		var entity ProductVariants
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.VariantName,
			&entity.Sku,
			&entity.Barcode,
			&entity.Attributes,
			&entity.CostPrice,
			&entity.SellingPrice,
			&entity.CompareAtPrice,
			&entity.CurrentStock,
			&entity.ReorderLevel,
			&entity.ReorderQuantity,
			&entity.Weight,
			&entity.WeightUnit,
			&entity.Dimensions,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.SortOrder,
			&entity.ImageUrl,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(costPrice,
			&entity.(sellingPrice,
			&entity.(compareAtPrice,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_variants: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating product_variants rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing product_variants record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ProductVariants) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_variants", duration, nil)
	}()

	query := `
		UPDATE product_variants
		SET
			, organization_id = $2
			, product_id = $3
			, variant_name = $4
			, sku = $5
			, barcode = $6
			, attributes = $7
			, cost_price = $8
			, selling_price = $9
			, compare_at_price = $10
			, current_stock = $11
			, reorder_level = $12
			, reorder_quantity = $13
			, weight = $14
			, weight_unit = $15
			, dimensions = $16
			, is_active = $17
			, is_default = $18
			, sort_order = $19
			, image_url = $20
			, notes = $21
			, metadata = $22
			, updated_at = $24
			, deleted_at = $25
			, created_by = $26
			, updated_by = $27
			, (cost_price = $28
			, (selling_price = $29
			, (compare_at_price = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.VariantName,
		entity.Sku,
		entity.Barcode,
		entity.Attributes,
		entity.CostPrice,
		entity.SellingPrice,
		entity.CompareAtPrice,
		entity.CurrentStock,
		entity.ReorderLevel,
		entity.ReorderQuantity,
		entity.Weight,
		entity.WeightUnit,
		entity.Dimensions,
		entity.IsActive,
		entity.IsDefault,
		entity.SortOrder,
		entity.ImageUrl,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(costPrice,
		entity.(sellingPrice,
		entity.(compareAtPrice,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update product_variants", zap.Error(err))
		return fmt.Errorf("failed to update product_variants: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_variants not found or already deleted")
	}

	r.logger.Info("updated product_variants",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a product_variants record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_variants", duration, nil)
	}()

	query := `
		UPDATE product_variants
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete product_variants", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete product_variants: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_variants not found or already deleted")
	}

	r.logger.Info("deleted product_variants", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves product_variants records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ProductVariants, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_variants", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_variants
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_variants records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, variant_name
			, sku
			, barcode
			, attributes
			, cost_price
			, selling_price
			, compare_at_price
			, current_stock
			, reorder_level
			, reorder_quantity
			, weight
			, weight_unit
			, dimensions
			, is_active
			, is_default
			, sort_order
			, image_url
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (cost_price
			, (selling_price
			, (compare_at_price
		FROM product_variants
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_variants by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_variants: %w", err)
	}
	defer rows.Close()

	var entities []*ProductVariants
	for rows.Next() {
		var entity ProductVariants
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.VariantName,
			&entity.Sku,
			&entity.Barcode,
			&entity.Attributes,
			&entity.CostPrice,
			&entity.SellingPrice,
			&entity.CompareAtPrice,
			&entity.CurrentStock,
			&entity.ReorderLevel,
			&entity.ReorderQuantity,
			&entity.Weight,
			&entity.WeightUnit,
			&entity.Dimensions,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.SortOrder,
			&entity.ImageUrl,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(costPrice,
			&entity.(sellingPrice,
			&entity.(compareAtPrice,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_variants: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

