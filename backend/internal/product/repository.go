package product

import (
	"encoding/json"
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

// Repository handles database operations for Products
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Products repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Products represents a products entity
type Products struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	Sku *string `json:"sku" db:"sku"`
	Barcode *string `json:"barcode" db:"barcode"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	CategoryId *uuid.UUID `json:"category_id" db:"category_id"`
	CostPrice *float64 `json:"cost_price" db:"cost_price"`
	SellingPrice float64 `json:"selling_price" db:"selling_price"`
	CompareAtPrice *float64 `json:"compare_at_price" db:"compare_at_price"`
	TaxRate *float64 `json:"tax_rate" db:"tax_rate"`
	IsTaxInclusive *bool `json:"is_tax_inclusive" db:"is_tax_inclusive"`
	TrackInventory *bool `json:"track_inventory" db:"track_inventory"`
	CurrentStock *float64 `json:"current_stock" db:"current_stock"`
	LowStockThreshold *float64 `json:"low_stock_threshold" db:"low_stock_threshold"`
	Unit *string `json:"unit" db:"unit"`
	IsService *bool `json:"is_service" db:"is_service"`
	IsComposite *bool `json:"is_composite" db:"is_composite"`
	HasVariants *bool `json:"has_variants" db:"has_variants"`
	ImageUrl *string `json:"image_url" db:"image_url"`
	Images json.RawMessage `json:"images" db:"images"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsFeatured *bool `json:"is_featured" db:"is_featured"`
	CustomFields json.RawMessage `json:"custom_fields" db:"custom_fields"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new products record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Products) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "products", duration, nil)
	}()

	query := `
		INSERT INTO products (
			, organization_id
			, sku
			, barcode
			, name
			, description
			, category_id
			, cost_price
			, selling_price
			, compare_at_price
			, tax_rate
			, is_tax_inclusive
			, track_inventory
			, current_stock
			, low_stock_threshold
			, unit
			, is_service
			, is_composite
			, has_variants
			, image_url
			, images
			, sort_order
			, is_active
			, is_featured
			, custom_fields
			, metadata
			, deleted_at
			, created_by
			, updated_by
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
			, $23
			, $24
			, $25
			, $26
			, $29
			, $30
			, $31
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Sku,
		entity.Barcode,
		entity.Name,
		entity.Description,
		entity.CategoryId,
		entity.CostPrice,
		entity.SellingPrice,
		entity.CompareAtPrice,
		entity.TaxRate,
		entity.IsTaxInclusive,
		entity.TrackInventory,
		entity.CurrentStock,
		entity.LowStockThreshold,
		entity.Unit,
		entity.IsService,
		entity.IsComposite,
		entity.HasVariants,
		entity.ImageUrl,
		entity.Images,
		entity.SortOrder,
		entity.IsActive,
		entity.IsFeatured,
		entity.CustomFields,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create products", zap.Error(err))
		return fmt.Errorf("failed to create products: %w", err)
	}

	r.logger.Info("created products",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a products by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Products, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "products", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, sku
			, barcode
			, name
			, description
			, category_id
			, cost_price
			, selling_price
			, compare_at_price
			, tax_rate
			, is_tax_inclusive
			, track_inventory
			, current_stock
			, low_stock_threshold
			, unit
			, is_service
			, is_composite
			, has_variants
			, image_url
			, images
			, sort_order
			, is_active
			, is_featured
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM products
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Products
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Sku,
		&entity.Barcode,
		&entity.Name,
		&entity.Description,
		&entity.CategoryId,
		&entity.CostPrice,
		&entity.SellingPrice,
		&entity.CompareAtPrice,
		&entity.TaxRate,
		&entity.IsTaxInclusive,
		&entity.TrackInventory,
		&entity.CurrentStock,
		&entity.LowStockThreshold,
		&entity.Unit,
		&entity.IsService,
		&entity.IsComposite,
		&entity.HasVariants,
		&entity.ImageUrl,
		&entity.Images,
		&entity.SortOrder,
		&entity.IsActive,
		&entity.IsFeatured,
		&entity.CustomFields,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("products not found")
	}

	if err != nil {
		r.logger.Error("failed to get products", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of products records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Products, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "products", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sku
			, barcode
			, name
			, description
			, category_id
			, cost_price
			, selling_price
			, compare_at_price
			, tax_rate
			, is_tax_inclusive
			, track_inventory
			, current_stock
			, low_stock_threshold
			, unit
			, is_service
			, is_composite
			, has_variants
			, image_url
			, images
			, sort_order
			, is_active
			, is_featured
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM products
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list products", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var entities []*Products
	for rows.Next() {
		var entity Products
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Sku,
			&entity.Barcode,
			&entity.Name,
			&entity.Description,
			&entity.CategoryId,
			&entity.CostPrice,
			&entity.SellingPrice,
			&entity.CompareAtPrice,
			&entity.TaxRate,
			&entity.IsTaxInclusive,
			&entity.TrackInventory,
			&entity.CurrentStock,
			&entity.LowStockThreshold,
			&entity.Unit,
			&entity.IsService,
			&entity.IsComposite,
			&entity.HasVariants,
			&entity.ImageUrl,
			&entity.Images,
			&entity.SortOrder,
			&entity.IsActive,
			&entity.IsFeatured,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan products: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating products rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing products record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Products) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "products", duration, nil)
	}()

	query := `
		UPDATE products
		SET
			, organization_id = $2
			, sku = $3
			, barcode = $4
			, name = $5
			, description = $6
			, category_id = $7
			, cost_price = $8
			, selling_price = $9
			, compare_at_price = $10
			, tax_rate = $11
			, is_tax_inclusive = $12
			, track_inventory = $13
			, current_stock = $14
			, low_stock_threshold = $15
			, unit = $16
			, is_service = $17
			, is_composite = $18
			, has_variants = $19
			, image_url = $20
			, images = $21
			, sort_order = $22
			, is_active = $23
			, is_featured = $24
			, custom_fields = $25
			, metadata = $26
			, updated_at = $28
			, deleted_at = $29
			, created_by = $30
			, updated_by = $31
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $32
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Sku,
		entity.Barcode,
		entity.Name,
		entity.Description,
		entity.CategoryId,
		entity.CostPrice,
		entity.SellingPrice,
		entity.CompareAtPrice,
		entity.TaxRate,
		entity.IsTaxInclusive,
		entity.TrackInventory,
		entity.CurrentStock,
		entity.LowStockThreshold,
		entity.Unit,
		entity.IsService,
		entity.IsComposite,
		entity.HasVariants,
		entity.ImageUrl,
		entity.Images,
		entity.SortOrder,
		entity.IsActive,
		entity.IsFeatured,
		entity.CustomFields,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update products", zap.Error(err))
		return fmt.Errorf("failed to update products: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("products not found or already deleted")
	}

	r.logger.Info("updated products",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a products record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "products", duration, nil)
	}()

	query := `
		UPDATE products
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete products", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete products: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("products not found or already deleted")
	}

	r.logger.Info("deleted products", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves products records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Products, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "products", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sku
			, barcode
			, name
			, description
			, category_id
			, cost_price
			, selling_price
			, compare_at_price
			, tax_rate
			, is_tax_inclusive
			, track_inventory
			, current_stock
			, low_stock_threshold
			, unit
			, is_service
			, is_composite
			, has_variants
			, image_url
			, images
			, sort_order
			, is_active
			, is_featured
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM products
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list products by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var entities []*Products
	for rows.Next() {
		var entity Products
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Sku,
			&entity.Barcode,
			&entity.Name,
			&entity.Description,
			&entity.CategoryId,
			&entity.CostPrice,
			&entity.SellingPrice,
			&entity.CompareAtPrice,
			&entity.TaxRate,
			&entity.IsTaxInclusive,
			&entity.TrackInventory,
			&entity.CurrentStock,
			&entity.LowStockThreshold,
			&entity.Unit,
			&entity.IsService,
			&entity.IsComposite,
			&entity.HasVariants,
			&entity.ImageUrl,
			&entity.Images,
			&entity.SortOrder,
			&entity.IsActive,
			&entity.IsFeatured,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan products: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

