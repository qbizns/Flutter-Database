package modifier

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

// Repository handles database operations for Modifiers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Modifiers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Modifiers represents a modifiers entity
type Modifiers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ModifierGroupId *uuid.UUID `json:"modifier_group_id" db:"modifier_group_id"`
	ModifierName string `json:"modifier_name" db:"modifier_name"`
	ModifierCode *string `json:"modifier_code" db:"modifier_code"`
	DisplayName *string `json:"display_name" db:"display_name"`
	PriceAdjustment *float64 `json:"price_adjustment" db:"price_adjustment"`
	PriceType *string `json:"price_type" db:"price_type"`
	IsAvailable *bool `json:"is_available" db:"is_available"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	TrackInventory *bool `json:"track_inventory" db:"track_inventory"`
	CurrentStock *float64 `json:"current_stock" db:"current_stock"`
	LowStockThreshold *float64 `json:"low_stock_threshold" db:"low_stock_threshold"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	ImageUrl *string `json:"image_url" db:"image_url"`
	Description *string `json:"description" db:"description"`
	AllergenInfo *string `json:"allergen_info" db:"allergen_info"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new modifiers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Modifiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "modifiers", duration, nil)
	}()

	query := `
		INSERT INTO modifiers (
			, organization_id
			, modifier_group_id
			, modifier_name
			, modifier_code
			, display_name
			, price_adjustment
			, price_type
			, is_available
			, is_default
			, track_inventory
			, current_stock
			, low_stock_threshold
			, display_order
			, image_url
			, description
			, allergen_info
			, notes
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
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ModifierGroupId,
		entity.ModifierName,
		entity.ModifierCode,
		entity.DisplayName,
		entity.PriceAdjustment,
		entity.PriceType,
		entity.IsAvailable,
		entity.IsDefault,
		entity.TrackInventory,
		entity.CurrentStock,
		entity.LowStockThreshold,
		entity.DisplayOrder,
		entity.ImageUrl,
		entity.Description,
		entity.AllergenInfo,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create modifiers", zap.Error(err))
		return fmt.Errorf("failed to create modifiers: %w", err)
	}

	r.logger.Info("created modifiers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a modifiers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Modifiers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifiers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, modifier_group_id
			, modifier_name
			, modifier_code
			, display_name
			, price_adjustment
			, price_type
			, is_available
			, is_default
			, track_inventory
			, current_stock
			, low_stock_threshold
			, display_order
			, image_url
			, description
			, allergen_info
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM modifiers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Modifiers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ModifierGroupId,
		&entity.ModifierName,
		&entity.ModifierCode,
		&entity.DisplayName,
		&entity.PriceAdjustment,
		&entity.PriceType,
		&entity.IsAvailable,
		&entity.IsDefault,
		&entity.TrackInventory,
		&entity.CurrentStock,
		&entity.LowStockThreshold,
		&entity.DisplayOrder,
		&entity.ImageUrl,
		&entity.Description,
		&entity.AllergenInfo,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("modifiers not found")
	}

	if err != nil {
		r.logger.Error("failed to get modifiers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get modifiers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of modifiers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Modifiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM modifiers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count modifiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, modifier_group_id
			, modifier_name
			, modifier_code
			, display_name
			, price_adjustment
			, price_type
			, is_available
			, is_default
			, track_inventory
			, current_stock
			, low_stock_threshold
			, display_order
			, image_url
			, description
			, allergen_info
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM modifiers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list modifiers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list modifiers: %w", err)
	}
	defer rows.Close()

	var entities []*Modifiers
	for rows.Next() {
		var entity Modifiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ModifierGroupId,
			&entity.ModifierName,
			&entity.ModifierCode,
			&entity.DisplayName,
			&entity.PriceAdjustment,
			&entity.PriceType,
			&entity.IsAvailable,
			&entity.IsDefault,
			&entity.TrackInventory,
			&entity.CurrentStock,
			&entity.LowStockThreshold,
			&entity.DisplayOrder,
			&entity.ImageUrl,
			&entity.Description,
			&entity.AllergenInfo,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan modifiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating modifiers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing modifiers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Modifiers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "modifiers", duration, nil)
	}()

	query := `
		UPDATE modifiers
		SET
			, organization_id = $2
			, modifier_group_id = $3
			, modifier_name = $4
			, modifier_code = $5
			, display_name = $6
			, price_adjustment = $7
			, price_type = $8
			, is_available = $9
			, is_default = $10
			, track_inventory = $11
			, current_stock = $12
			, low_stock_threshold = $13
			, display_order = $14
			, image_url = $15
			, description = $16
			, allergen_info = $17
			, notes = $18
			, metadata = $19
			, updated_at = $21
			, deleted_at = $22
			, created_by = $23
			, updated_by = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ModifierGroupId,
		entity.ModifierName,
		entity.ModifierCode,
		entity.DisplayName,
		entity.PriceAdjustment,
		entity.PriceType,
		entity.IsAvailable,
		entity.IsDefault,
		entity.TrackInventory,
		entity.CurrentStock,
		entity.LowStockThreshold,
		entity.DisplayOrder,
		entity.ImageUrl,
		entity.Description,
		entity.AllergenInfo,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update modifiers", zap.Error(err))
		return fmt.Errorf("failed to update modifiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("modifiers not found or already deleted")
	}

	r.logger.Info("updated modifiers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a modifiers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "modifiers", duration, nil)
	}()

	query := `
		UPDATE modifiers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete modifiers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete modifiers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("modifiers not found or already deleted")
	}

	r.logger.Info("deleted modifiers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves modifiers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Modifiers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "modifiers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM modifiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count modifiers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, modifier_group_id
			, modifier_name
			, modifier_code
			, display_name
			, price_adjustment
			, price_type
			, is_available
			, is_default
			, track_inventory
			, current_stock
			, low_stock_threshold
			, display_order
			, image_url
			, description
			, allergen_info
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM modifiers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list modifiers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list modifiers: %w", err)
	}
	defer rows.Close()

	var entities []*Modifiers
	for rows.Next() {
		var entity Modifiers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ModifierGroupId,
			&entity.ModifierName,
			&entity.ModifierCode,
			&entity.DisplayName,
			&entity.PriceAdjustment,
			&entity.PriceType,
			&entity.IsAvailable,
			&entity.IsDefault,
			&entity.TrackInventory,
			&entity.CurrentStock,
			&entity.LowStockThreshold,
			&entity.DisplayOrder,
			&entity.ImageUrl,
			&entity.Description,
			&entity.AllergenInfo,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan modifiers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

