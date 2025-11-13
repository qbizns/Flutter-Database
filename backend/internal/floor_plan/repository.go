package floor_plan

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

// Repository handles database operations for FloorPlans
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FloorPlans repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FloorPlans represents a floor_plans entity
type FloorPlans struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	FloorName string `json:"floor_name" db:"floor_name"`
	FloorLevel *int64 `json:"floor_level" db:"floor_level"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	LayoutConfig json.RawMessage `json:"layout_config" db:"layout_config"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new floor_plans record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FloorPlans) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "floor_plans", duration, nil)
	}()

	query := `
		INSERT INTO floor_plans (
			, organization_id
			, location_id
			, floor_name
			, floor_level
			, display_order
			, layout_config
			, is_active
			, is_default
			, description
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
			, $15
			, $16
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorName,
		entity.FloorLevel,
		entity.DisplayOrder,
		entity.LayoutConfig,
		entity.IsActive,
		entity.IsDefault,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create floor_plans", zap.Error(err))
		return fmt.Errorf("failed to create floor_plans: %w", err)
	}

	r.logger.Info("created floor_plans",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a floor_plans by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FloorPlans, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "floor_plans", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_name
			, floor_level
			, display_order
			, layout_config
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM floor_plans
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FloorPlans
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.FloorName,
		&entity.FloorLevel,
		&entity.DisplayOrder,
		&entity.LayoutConfig,
		&entity.IsActive,
		&entity.IsDefault,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("floor_plans not found")
	}

	if err != nil {
		r.logger.Error("failed to get floor_plans", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get floor_plans: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of floor_plans records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FloorPlans, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "floor_plans", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM floor_plans
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count floor_plans records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_name
			, floor_level
			, display_order
			, layout_config
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM floor_plans
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list floor_plans", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list floor_plans: %w", err)
	}
	defer rows.Close()

	var entities []*FloorPlans
	for rows.Next() {
		var entity FloorPlans
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorName,
			&entity.FloorLevel,
			&entity.DisplayOrder,
			&entity.LayoutConfig,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan floor_plans: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating floor_plans rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing floor_plans record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FloorPlans) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "floor_plans", duration, nil)
	}()

	query := `
		UPDATE floor_plans
		SET
			, organization_id = $2
			, location_id = $3
			, floor_name = $4
			, floor_level = $5
			, display_order = $6
			, layout_config = $7
			, is_active = $8
			, is_default = $9
			, description = $10
			, notes = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, created_by = $16
			, updated_by = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorName,
		entity.FloorLevel,
		entity.DisplayOrder,
		entity.LayoutConfig,
		entity.IsActive,
		entity.IsDefault,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update floor_plans", zap.Error(err))
		return fmt.Errorf("failed to update floor_plans: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("floor_plans not found or already deleted")
	}

	r.logger.Info("updated floor_plans",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a floor_plans record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "floor_plans", duration, nil)
	}()

	query := `
		UPDATE floor_plans
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete floor_plans", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete floor_plans: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("floor_plans not found or already deleted")
	}

	r.logger.Info("deleted floor_plans", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves floor_plans records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*FloorPlans, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "floor_plans", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM floor_plans
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count floor_plans records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_name
			, floor_level
			, display_order
			, layout_config
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM floor_plans
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list floor_plans by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list floor_plans: %w", err)
	}
	defer rows.Close()

	var entities []*FloorPlans
	for rows.Next() {
		var entity FloorPlans
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorName,
			&entity.FloorLevel,
			&entity.DisplayOrder,
			&entity.LayoutConfig,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan floor_plans: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

