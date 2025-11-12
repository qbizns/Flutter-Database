package units_of_measure

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

// Repository handles database operations for UnitsOfMeasure
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new UnitsOfMeasure repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// UnitsOfMeasure represents a units_of_measure entity
type UnitsOfMeasure struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	UomCode string `json:"uom_code" db:"uom_code"`
	UomName string `json:"uom_name" db:"uom_name"`
	UomType string `json:"uom_type" db:"uom_type"`
	IsBaseUnit *bool `json:"is_base_unit" db:"is_base_unit"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new units_of_measure record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *UnitsOfMeasure) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "units_of_measure", duration, nil)
	}()

	query := `
		INSERT INTO units_of_measure (
			, organization_id
			, uom_code
			, uom_name
			, uom_type
			, is_base_unit
			, is_active
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.UomCode,
		entity.UomName,
		entity.UomType,
		entity.IsBaseUnit,
		entity.IsActive,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create units_of_measure", zap.Error(err))
		return fmt.Errorf("failed to create units_of_measure: %w", err)
	}

	r.logger.Info("created units_of_measure",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a units_of_measure by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*UnitsOfMeasure, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "units_of_measure", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, uom_code
			, uom_name
			, uom_type
			, is_base_unit
			, is_active
			, created_at
			, deleted_at
		FROM units_of_measure
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity UnitsOfMeasure
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.UomCode,
		&entity.UomName,
		&entity.UomType,
		&entity.IsBaseUnit,
		&entity.IsActive,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("units_of_measure not found")
	}

	if err != nil {
		r.logger.Error("failed to get units_of_measure", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get units_of_measure: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of units_of_measure records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*UnitsOfMeasure, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "units_of_measure", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM units_of_measure
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count units_of_measure records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, uom_code
			, uom_name
			, uom_type
			, is_base_unit
			, is_active
			, created_at
			, deleted_at
		FROM units_of_measure
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list units_of_measure", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list units_of_measure: %w", err)
	}
	defer rows.Close()

	var entities []*UnitsOfMeasure
	for rows.Next() {
		var entity UnitsOfMeasure
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UomCode,
			&entity.UomName,
			&entity.UomType,
			&entity.IsBaseUnit,
			&entity.IsActive,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan units_of_measure: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating units_of_measure rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing units_of_measure record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *UnitsOfMeasure) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "units_of_measure", duration, nil)
	}()

	query := `
		UPDATE units_of_measure
		SET
			, organization_id = $2
			, uom_code = $3
			, uom_name = $4
			, uom_type = $5
			, is_base_unit = $6
			, is_active = $7
			, deleted_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.UomCode,
		entity.UomName,
		entity.UomType,
		entity.IsBaseUnit,
		entity.IsActive,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update units_of_measure", zap.Error(err))
		return fmt.Errorf("failed to update units_of_measure: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("units_of_measure not found or already deleted")
	}

	r.logger.Info("updated units_of_measure",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a units_of_measure record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "units_of_measure", duration, nil)
	}()

	query := `
		UPDATE units_of_measure
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete units_of_measure", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete units_of_measure: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("units_of_measure not found or already deleted")
	}

	r.logger.Info("deleted units_of_measure", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves units_of_measure records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*UnitsOfMeasure, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "units_of_measure", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM units_of_measure
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count units_of_measure records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, uom_code
			, uom_name
			, uom_type
			, is_base_unit
			, is_active
			, created_at
			, deleted_at
		FROM units_of_measure
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list units_of_measure by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list units_of_measure: %w", err)
	}
	defer rows.Close()

	var entities []*UnitsOfMeasure
	for rows.Next() {
		var entity UnitsOfMeasure
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UomCode,
			&entity.UomName,
			&entity.UomType,
			&entity.IsBaseUnit,
			&entity.IsActive,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan units_of_measure: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

