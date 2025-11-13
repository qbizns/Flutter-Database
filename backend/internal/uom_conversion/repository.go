package uom_conversion

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

// Repository handles database operations for UomConversions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new UomConversions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// UomConversions represents a uom_conversions entity
type UomConversions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	FromUomId uuid.UUID `json:"from_uom_id" db:"from_uom_id"`
	ToUomId uuid.UUID `json:"to_uom_id" db:"to_uom_id"`
	ConversionFactor float64 `json:"conversion_factor" db:"conversion_factor"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new uom_conversions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *UomConversions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "uom_conversions", duration, nil)
	}()

	query := `
		INSERT INTO uom_conversions (
			, from_uom_id
			, to_uom_id
			, conversion_factor
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $6
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.FromUomId,
		entity.ToUomId,
		entity.ConversionFactor,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create uom_conversions", zap.Error(err))
		return fmt.Errorf("failed to create uom_conversions: %w", err)
	}

	r.logger.Info("created uom_conversions",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a uom_conversions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*UomConversions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "uom_conversions", duration, nil)
	}()

	query := `
		SELECT
			id
			, from_uom_id
			, to_uom_id
			, conversion_factor
			, created_at
			, deleted_at
		FROM uom_conversions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity UomConversions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.FromUomId,
		&entity.ToUomId,
		&entity.ConversionFactor,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("uom_conversions not found")
	}

	if err != nil {
		r.logger.Error("failed to get uom_conversions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get uom_conversions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of uom_conversions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*UomConversions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "uom_conversions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM uom_conversions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count uom_conversions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, from_uom_id
			, to_uom_id
			, conversion_factor
			, created_at
			, deleted_at
		FROM uom_conversions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list uom_conversions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list uom_conversions: %w", err)
	}
	defer rows.Close()

	var entities []*UomConversions
	for rows.Next() {
		var entity UomConversions
		err := rows.Scan(
			&entity.Id,
			&entity.FromUomId,
			&entity.ToUomId,
			&entity.ConversionFactor,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan uom_conversions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating uom_conversions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing uom_conversions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *UomConversions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "uom_conversions", duration, nil)
	}()

	query := `
		UPDATE uom_conversions
		SET
			, from_uom_id = $2
			, to_uom_id = $3
			, conversion_factor = $4
			, deleted_at = $6
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.FromUomId,
		entity.ToUomId,
		entity.ConversionFactor,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update uom_conversions", zap.Error(err))
		return fmt.Errorf("failed to update uom_conversions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("uom_conversions not found or already deleted")
	}

	r.logger.Info("updated uom_conversions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a uom_conversions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "uom_conversions", duration, nil)
	}()

	query := `
		UPDATE uom_conversions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete uom_conversions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete uom_conversions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("uom_conversions not found or already deleted")
	}

	r.logger.Info("deleted uom_conversions", zap.String("id", id.String()))
	return nil
}



