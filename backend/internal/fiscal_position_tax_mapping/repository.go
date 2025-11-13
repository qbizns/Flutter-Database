package fiscal_position_tax_mapping

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

// Repository handles database operations for FiscalPositionTaxMappings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FiscalPositionTaxMappings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FiscalPositionTaxMappings represents a fiscal_position_tax_mappings entity
type FiscalPositionTaxMappings struct {
	Id *uuid.UUID `json:"id" db:"id"`
	FiscalPositionId uuid.UUID `json:"fiscal_position_id" db:"fiscal_position_id"`
	SourceTaxId uuid.UUID `json:"source_tax_id" db:"source_tax_id"`
	DestinationTaxId *uuid.UUID `json:"destination_tax_id" db:"destination_tax_id"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new fiscal_position_tax_mappings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FiscalPositionTaxMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "fiscal_position_tax_mappings", duration, nil)
	}()

	query := `
		INSERT INTO fiscal_position_tax_mappings (
			, fiscal_position_id
			, source_tax_id
			, destination_tax_id
			, created_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $7
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.FiscalPositionId,
		entity.SourceTaxId,
		entity.DestinationTaxId,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create fiscal_position_tax_mappings", zap.Error(err))
		return fmt.Errorf("failed to create fiscal_position_tax_mappings: %w", err)
	}

	r.logger.Info("created fiscal_position_tax_mappings",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a fiscal_position_tax_mappings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FiscalPositionTaxMappings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_position_tax_mappings", duration, nil)
	}()

	query := `
		SELECT
			id
			, fiscal_position_id
			, source_tax_id
			, destination_tax_id
			, created_by
			, created_at
			, deleted_at
		FROM fiscal_position_tax_mappings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FiscalPositionTaxMappings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.FiscalPositionId,
		&entity.SourceTaxId,
		&entity.DestinationTaxId,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fiscal_position_tax_mappings not found")
	}

	if err != nil {
		r.logger.Error("failed to get fiscal_position_tax_mappings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get fiscal_position_tax_mappings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of fiscal_position_tax_mappings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FiscalPositionTaxMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_position_tax_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fiscal_position_tax_mappings
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fiscal_position_tax_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, fiscal_position_id
			, source_tax_id
			, destination_tax_id
			, created_by
			, created_at
			, deleted_at
		FROM fiscal_position_tax_mappings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fiscal_position_tax_mappings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fiscal_position_tax_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*FiscalPositionTaxMappings
	for rows.Next() {
		var entity FiscalPositionTaxMappings
		err := rows.Scan(
			&entity.Id,
			&entity.FiscalPositionId,
			&entity.SourceTaxId,
			&entity.DestinationTaxId,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fiscal_position_tax_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating fiscal_position_tax_mappings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing fiscal_position_tax_mappings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FiscalPositionTaxMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_position_tax_mappings", duration, nil)
	}()

	query := `
		UPDATE fiscal_position_tax_mappings
		SET
			, fiscal_position_id = $2
			, source_tax_id = $3
			, destination_tax_id = $4
			, created_by = $5
			, deleted_at = $7
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.FiscalPositionId,
		entity.SourceTaxId,
		entity.DestinationTaxId,
		entity.CreatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update fiscal_position_tax_mappings", zap.Error(err))
		return fmt.Errorf("failed to update fiscal_position_tax_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_position_tax_mappings not found or already deleted")
	}

	r.logger.Info("updated fiscal_position_tax_mappings",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a fiscal_position_tax_mappings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_position_tax_mappings", duration, nil)
	}()

	query := `
		UPDATE fiscal_position_tax_mappings
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete fiscal_position_tax_mappings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete fiscal_position_tax_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_position_tax_mappings not found or already deleted")
	}

	r.logger.Info("deleted fiscal_position_tax_mappings", zap.String("id", id.String()))
	return nil
}



