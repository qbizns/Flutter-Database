package product_modifier_group

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

// Repository handles database operations for ProductModifierGroups
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ProductModifierGroups repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ProductModifierGroups represents a product_modifier_groups entity
type ProductModifierGroups struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ModifierGroupId uuid.UUID `json:"modifier_group_id" db:"modifier_group_id"`
	IsRequired *bool `json:"is_required" db:"is_required"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	IsActive *bool `json:"is_active" db:"is_active"`
	OverrideMinSelections *int64 `json:"override_min_selections" db:"override_min_selections"`
	OverrideMaxSelections *int64 `json:"override_max_selections" db:"override_max_selections"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new product_modifier_groups record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ProductModifierGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "product_modifier_groups", duration, nil)
	}()

	query := `
		INSERT INTO product_modifier_groups (
			, organization_id
			, product_id
			, modifier_group_id
			, is_required
			, display_order
			, is_active
			, override_min_selections
			, override_max_selections
			, notes
			, metadata
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ModifierGroupId,
		entity.IsRequired,
		entity.DisplayOrder,
		entity.IsActive,
		entity.OverrideMinSelections,
		entity.OverrideMaxSelections,
		entity.Notes,
		entity.Metadata,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create product_modifier_groups", zap.Error(err))
		return fmt.Errorf("failed to create product_modifier_groups: %w", err)
	}

	r.logger.Info("created product_modifier_groups",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a product_modifier_groups by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ProductModifierGroups, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_modifier_groups", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, modifier_group_id
			, is_required
			, display_order
			, is_active
			, override_min_selections
			, override_max_selections
			, notes
			, metadata
			, created_at
			, updated_at
		FROM product_modifier_groups
		WHERE id = $1
		
	`

	var entity ProductModifierGroups
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.ModifierGroupId,
		&entity.IsRequired,
		&entity.DisplayOrder,
		&entity.IsActive,
		&entity.OverrideMinSelections,
		&entity.OverrideMaxSelections,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("product_modifier_groups not found")
	}

	if err != nil {
		r.logger.Error("failed to get product_modifier_groups", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get product_modifier_groups: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of product_modifier_groups records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ProductModifierGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_modifier_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_modifier_groups
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_modifier_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, modifier_group_id
			, is_required
			, display_order
			, is_active
			, override_min_selections
			, override_max_selections
			, notes
			, metadata
			, created_at
			, updated_at
		FROM product_modifier_groups
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_modifier_groups", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_modifier_groups: %w", err)
	}
	defer rows.Close()

	var entities []*ProductModifierGroups
	for rows.Next() {
		var entity ProductModifierGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ModifierGroupId,
			&entity.IsRequired,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.OverrideMinSelections,
			&entity.OverrideMaxSelections,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_modifier_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating product_modifier_groups rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing product_modifier_groups record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ProductModifierGroups) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "product_modifier_groups", duration, nil)
	}()

	query := `
		UPDATE product_modifier_groups
		SET
			, organization_id = $2
			, product_id = $3
			, modifier_group_id = $4
			, is_required = $5
			, display_order = $6
			, is_active = $7
			, override_min_selections = $8
			, override_max_selections = $9
			, notes = $10
			, metadata = $11
			, updated_at = $13
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.ModifierGroupId,
		entity.IsRequired,
		entity.DisplayOrder,
		entity.IsActive,
		entity.OverrideMinSelections,
		entity.OverrideMaxSelections,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update product_modifier_groups", zap.Error(err))
		return fmt.Errorf("failed to update product_modifier_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_modifier_groups not found or already deleted")
	}

	r.logger.Info("updated product_modifier_groups",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a product_modifier_groups record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "product_modifier_groups", duration, nil)
	}()

	query := `DELETE FROM product_modifier_groups WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete product_modifier_groups", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete product_modifier_groups: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("product_modifier_groups not found")
	}

	r.logger.Info("deleted product_modifier_groups", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves product_modifier_groups records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ProductModifierGroups, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "product_modifier_groups", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM product_modifier_groups
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count product_modifier_groups records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, modifier_group_id
			, is_required
			, display_order
			, is_active
			, override_min_selections
			, override_max_selections
			, notes
			, metadata
			, created_at
			, updated_at
		FROM product_modifier_groups
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list product_modifier_groups by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list product_modifier_groups: %w", err)
	}
	defer rows.Close()

	var entities []*ProductModifierGroups
	for rows.Next() {
		var entity ProductModifierGroups
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.ModifierGroupId,
			&entity.IsRequired,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.OverrideMinSelections,
			&entity.OverrideMaxSelections,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product_modifier_groups: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

