package product_component

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

// Repository handles database operations for ProductComponents
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ProductComponents repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ProductComponents represents a product_components entity
type ProductComponents struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ParentProductId uuid.UUID `json:"parent_product_id" db:"parent_product_id"`
	ComponentProductId *uuid.UUID `json:"component_product_id" db:"component_product_id"`
	ComponentVariantId *uuid.UUID `json:"component_variant_id" db:"component_variant_id"`
	Quantity float64 `json:"quantity" db:"quantity"`
	InheritPrice *bool `json:"inherit_price" db:"inherit_price"`
	PriceOverride *float64 `json:"price_override" db:"price_override"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	IsOptional *bool `json:"is_optional" db:"is_optional"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	ComponentProductId string `json:"component_product_id" db:"component_product_id"`
}

// Create inserts a new product_components record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ProductComponents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "product_components", duration, nil)
	}()

	query := `
		INSERT INTO product_components (
			, organization_id
			, parent_product_id
			, component_product_id
			, component_variant_id
			, quantity
			, inherit_price
			, price_override
			, display_order
			, is_optional
			, deleted_at
			, component_product_id
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
			, $13
			, $14
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ParentProductId,
		entity.ComponentProductId,
		entity.ComponentVariantId,
		entity.Quantity,
		entity.InheritPrice,
		entity.PriceOverride,
		entity.DisplayOrder,
		entity.IsOptional,
		entity.DeletedAt,
		entity.ComponentProductId,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create product_components", zap.Error(err))
		return fmt.Errorf("failed to create product_components: %w", err)
	}

	r.logger.Info("created product_components",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a product_components by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ProductComponents, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_components", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, parent_product_id
			, component_product_id
			, component_variant_id
			, quantity
			, inherit_price
			, price_override
			, display_order
			, is_optional
			, created_at
			, updated_at
			, deleted_at
			, component_product_id
		FROM product_components
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ProductComponents
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ParentProductId,
		&entity.ComponentProductId,
		&entity.ComponentVariantId,
		&entity.Quantity,
		&entity.InheritPrice,
		&entity.PriceOverride,
		&entity.DisplayOrder,
		&entity.IsOptional,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.ComponentProductId,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("product_components not found")
	}

	if err != nil {
		r.logger.Error("failed to get product_components", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get product_components: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of product_components records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ProductComponents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_components", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_components
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_components records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, parent_product_id
			, component_product_id
			, component_variant_id
			, quantity
			, inherit_price
			, price_override
			, display_order
			, is_optional
			, created_at
			, updated_at
			, deleted_at
			, component_product_id
		FROM product_components
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_components", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_components: %w", err)
	}
	defer rows.Close()

	var entities []*ProductComponents
	for rows.Next() {
		var entity ProductComponents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ParentProductId,
			&entity.ComponentProductId,
			&entity.ComponentVariantId,
			&entity.Quantity,
			&entity.InheritPrice,
			&entity.PriceOverride,
			&entity.DisplayOrder,
			&entity.IsOptional,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.ComponentProductId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_components: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating product_components rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing product_components record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ProductComponents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_components", duration, nil)
	}()

	query := `
		UPDATE product_components
		SET
			, organization_id = $2
			, parent_product_id = $3
			, component_product_id = $4
			, component_variant_id = $5
			, quantity = $6
			, inherit_price = $7
			, price_override = $8
			, display_order = $9
			, is_optional = $10
			, updated_at = $12
			, deleted_at = $13
			, component_product_id = $14
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $15
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ParentProductId,
		entity.ComponentProductId,
		entity.ComponentVariantId,
		entity.Quantity,
		entity.InheritPrice,
		entity.PriceOverride,
		entity.DisplayOrder,
		entity.IsOptional,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.ComponentProductId,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update product_components", zap.Error(err))
		return fmt.Errorf("failed to update product_components: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_components not found or already deleted")
	}

	r.logger.Info("updated product_components",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a product_components record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_components", duration, nil)
	}()

	query := `
		UPDATE product_components
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete product_components", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete product_components: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_components not found or already deleted")
	}

	r.logger.Info("deleted product_components", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves product_components records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ProductComponents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_components", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_components
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_components records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, parent_product_id
			, component_product_id
			, component_variant_id
			, quantity
			, inherit_price
			, price_override
			, display_order
			, is_optional
			, created_at
			, updated_at
			, deleted_at
			, component_product_id
		FROM product_components
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_components by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_components: %w", err)
	}
	defer rows.Close()

	var entities []*ProductComponents
	for rows.Next() {
		var entity ProductComponents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ParentProductId,
			&entity.ComponentProductId,
			&entity.ComponentVariantId,
			&entity.Quantity,
			&entity.InheritPrice,
			&entity.PriceOverride,
			&entity.DisplayOrder,
			&entity.IsOptional,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.ComponentProductId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_components: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

