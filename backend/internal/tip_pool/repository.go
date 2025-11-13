package tip_pool

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

// Repository handles database operations for TipPools
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TipPools repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TipPools represents a tip_pools entity
type TipPools struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	PoolName string `json:"pool_name" db:"pool_name"`
	PoolType *string `json:"pool_type" db:"pool_type"`
	Description *string `json:"description" db:"description"`
	DistributionMethod *string `json:"distribution_method" db:"distribution_method"`
	DistributionConfig json.RawMessage `json:"distribution_config" db:"distribution_config"`
	EligiblePositions *string `json:"eligible_positions" db:"eligible_positions"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'daily', *string `json:"'daily'," db:"'daily',"`
	'equal', *string `json:"'equal'," db:"'equal',"`
}

// Create inserts a new tip_pools record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TipPools) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "tip_pools", duration, nil)
	}()

	query := `
		INSERT INTO tip_pools (
			, organization_id
			, location_id
			, pool_name
			, pool_type
			, description
			, distribution_method
			, distribution_config
			, eligible_positions
			, is_active
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'daily',
			, 'equal',
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
			, $14
			, $15
			, $16
			, $17
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.PoolName,
		entity.PoolType,
		entity.Description,
		entity.DistributionMethod,
		entity.DistributionConfig,
		entity.EligiblePositions,
		entity.IsActive,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'daily',,
		entity.'equal',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create tip_pools", zap.Error(err))
		return fmt.Errorf("failed to create tip_pools: %w", err)
	}

	r.logger.Info("created tip_pools",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a tip_pools by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TipPools, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_pools", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, pool_name
			, pool_type
			, description
			, distribution_method
			, distribution_config
			, eligible_positions
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'daily',
			, 'equal',
		FROM tip_pools
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TipPools
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.PoolName,
		&entity.PoolType,
		&entity.Description,
		&entity.DistributionMethod,
		&entity.DistributionConfig,
		&entity.EligiblePositions,
		&entity.IsActive,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'daily',,
		&entity.'equal',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tip_pools not found")
	}

	if err != nil {
		r.logger.Error("failed to get tip_pools", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get tip_pools: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of tip_pools records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TipPools, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_pools", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tip_pools
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tip_pools records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, pool_name
			, pool_type
			, description
			, distribution_method
			, distribution_config
			, eligible_positions
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'daily',
			, 'equal',
		FROM tip_pools
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tip_pools", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tip_pools: %w", err)
	}
	defer rows.Close()

	var entities []*TipPools
	for rows.Next() {
		var entity TipPools
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.PoolName,
			&entity.PoolType,
			&entity.Description,
			&entity.DistributionMethod,
			&entity.DistributionConfig,
			&entity.EligiblePositions,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'daily',,
			&entity.'equal',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tip_pools: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating tip_pools rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing tip_pools record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TipPools) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tip_pools", duration, nil)
	}()

	query := `
		UPDATE tip_pools
		SET
			, organization_id = $2
			, location_id = $3
			, pool_name = $4
			, pool_type = $5
			, description = $6
			, distribution_method = $7
			, distribution_config = $8
			, eligible_positions = $9
			, is_active = $10
			, metadata = $11
			, updated_at = $13
			, created_by = $14
			, updated_by = $15
			, deleted_at = $16
			, 'daily', = $17
			, 'equal', = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.PoolName,
		entity.PoolType,
		entity.Description,
		entity.DistributionMethod,
		entity.DistributionConfig,
		entity.EligiblePositions,
		entity.IsActive,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'daily',,
		entity.'equal',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update tip_pools", zap.Error(err))
		return fmt.Errorf("failed to update tip_pools: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tip_pools not found or already deleted")
	}

	r.logger.Info("updated tip_pools",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a tip_pools record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tip_pools", duration, nil)
	}()

	query := `
		UPDATE tip_pools
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete tip_pools", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete tip_pools: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tip_pools not found or already deleted")
	}

	r.logger.Info("deleted tip_pools", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves tip_pools records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TipPools, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_pools", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tip_pools
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tip_pools records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, pool_name
			, pool_type
			, description
			, distribution_method
			, distribution_config
			, eligible_positions
			, is_active
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'daily',
			, 'equal',
		FROM tip_pools
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tip_pools by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tip_pools: %w", err)
	}
	defer rows.Close()

	var entities []*TipPools
	for rows.Next() {
		var entity TipPools
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.PoolName,
			&entity.PoolType,
			&entity.Description,
			&entity.DistributionMethod,
			&entity.DistributionConfig,
			&entity.EligiblePositions,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'daily',,
			&entity.'equal',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tip_pools: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

