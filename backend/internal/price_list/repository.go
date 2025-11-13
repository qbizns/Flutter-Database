package price_list

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

// Repository handles database operations for PriceLists
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PriceLists repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PriceLists represents a price_lists entity
type PriceLists struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PriceListCode string `json:"price_list_code" db:"price_list_code"`
	PriceListName string `json:"price_list_name" db:"price_list_name"`
	PriceListType *string `json:"price_list_type" db:"price_list_type"`
	PriceListType *string `json:"price_list_type" db:"price_list_type"`
	EffectiveFrom *time.Time `json:"effective_from" db:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to" db:"effective_to"`
	BasePriceAdjustmentType *string `json:"base_price_adjustment_type" db:"base_price_adjustment_type"`
	BasePriceAdjustmentValue *float64 `json:"base_price_adjustment_value" db:"base_price_adjustment_value"`
	Priority *int64 `json:"priority" db:"priority"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new price_lists record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PriceLists) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "price_lists", duration, nil)
	}()

	query := `
		INSERT INTO price_lists (
			, organization_id
			, price_list_code
			, price_list_name
			, price_list_type
			, price_list_type
			, effective_from
			, effective_to
			, base_price_adjustment_type
			, base_price_adjustment_value
			, priority
			, is_active
			, description
			, created_by
			, updated_by
			, deleted_at
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
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PriceListCode,
		entity.PriceListName,
		entity.PriceListType,
		entity.PriceListType,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.BasePriceAdjustmentType,
		entity.BasePriceAdjustmentValue,
		entity.Priority,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create price_lists", zap.Error(err))
		return fmt.Errorf("failed to create price_lists: %w", err)
	}

	r.logger.Info("created price_lists",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a price_lists by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PriceLists, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "price_lists", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, price_list_code
			, price_list_name
			, price_list_type
			, price_list_type
			, effective_from
			, effective_to
			, base_price_adjustment_type
			, base_price_adjustment_value
			, priority
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM price_lists
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PriceLists
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PriceListCode,
		&entity.PriceListName,
		&entity.PriceListType,
		&entity.PriceListType,
		&entity.EffectiveFrom,
		&entity.EffectiveTo,
		&entity.BasePriceAdjustmentType,
		&entity.BasePriceAdjustmentValue,
		&entity.Priority,
		&entity.IsActive,
		&entity.Description,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("price_lists not found")
	}

	if err != nil {
		r.logger.Error("failed to get price_lists", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get price_lists: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of price_lists records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PriceLists, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "price_lists", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM price_lists
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count price_lists records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, price_list_code
			, price_list_name
			, price_list_type
			, price_list_type
			, effective_from
			, effective_to
			, base_price_adjustment_type
			, base_price_adjustment_value
			, priority
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM price_lists
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list price_lists", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list price_lists: %w", err)
	}
	defer rows.Close()

	var entities []*PriceLists
	for rows.Next() {
		var entity PriceLists
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PriceListCode,
			&entity.PriceListName,
			&entity.PriceListType,
			&entity.PriceListType,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.BasePriceAdjustmentType,
			&entity.BasePriceAdjustmentValue,
			&entity.Priority,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan price_lists: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating price_lists rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing price_lists record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PriceLists) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "price_lists", duration, nil)
	}()

	query := `
		UPDATE price_lists
		SET
			, organization_id = $2
			, price_list_code = $3
			, price_list_name = $4
			, price_list_type = $5
			, price_list_type = $6
			, effective_from = $7
			, effective_to = $8
			, base_price_adjustment_type = $9
			, base_price_adjustment_value = $10
			, priority = $11
			, is_active = $12
			, description = $13
			, created_by = $14
			, updated_by = $15
			, updated_at = $17
			, deleted_at = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PriceListCode,
		entity.PriceListName,
		entity.PriceListType,
		entity.PriceListType,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.BasePriceAdjustmentType,
		entity.BasePriceAdjustmentValue,
		entity.Priority,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update price_lists", zap.Error(err))
		return fmt.Errorf("failed to update price_lists: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("price_lists not found or already deleted")
	}

	r.logger.Info("updated price_lists",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a price_lists record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "price_lists", duration, nil)
	}()

	query := `
		UPDATE price_lists
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete price_lists", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete price_lists: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("price_lists not found or already deleted")
	}

	r.logger.Info("deleted price_lists", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves price_lists records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PriceLists, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "price_lists", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM price_lists
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count price_lists records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, price_list_code
			, price_list_name
			, price_list_type
			, price_list_type
			, effective_from
			, effective_to
			, base_price_adjustment_type
			, base_price_adjustment_value
			, priority
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM price_lists
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list price_lists by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list price_lists: %w", err)
	}
	defer rows.Close()

	var entities []*PriceLists
	for rows.Next() {
		var entity PriceLists
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PriceListCode,
			&entity.PriceListName,
			&entity.PriceListType,
			&entity.PriceListType,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.BasePriceAdjustmentType,
			&entity.BasePriceAdjustmentValue,
			&entity.Priority,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan price_lists: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

