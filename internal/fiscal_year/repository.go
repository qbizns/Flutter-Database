package fiscal_year

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

// Repository handles database operations for FiscalYears
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FiscalYears repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FiscalYears represents a fiscal_years entity
type FiscalYears struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	FiscalYear string `json:"fiscal_year" db:"fiscal_year"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate time.Time `json:"end_date" db:"end_date"`
	Status *string `json:"status" db:"status"`
	IsCurrent *bool `json:"is_current" db:"is_current"`
	ClosedBy *uuid.UUID `json:"closed_by" db:"closed_by"`
	ClosedAt *time.Time `json:"closed_at" db:"closed_at"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new fiscal_years record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FiscalYears) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "fiscal_years", duration, nil)
	}()

	query := `
		INSERT INTO fiscal_years (
			, organization_id
			, fiscal_year
			, start_date
			, end_date
			, status
			, is_current
			, closed_by
			, closed_at
			, notes
			, metadata
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
			, $14
			, $15
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.FiscalYear,
		entity.StartDate,
		entity.EndDate,
		entity.Status,
		entity.IsCurrent,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create fiscal_years", zap.Error(err))
		return fmt.Errorf("failed to create fiscal_years: %w", err)
	}

	r.logger.Info("created fiscal_years",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a fiscal_years by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FiscalYears, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_years", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, fiscal_year
			, start_date
			, end_date
			, status
			, is_current
			, closed_by
			, closed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fiscal_years
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FiscalYears
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.FiscalYear,
		&entity.StartDate,
		&entity.EndDate,
		&entity.Status,
		&entity.IsCurrent,
		&entity.ClosedBy,
		&entity.ClosedAt,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fiscal_years not found")
	}

	if err != nil {
		r.logger.Error("failed to get fiscal_years", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get fiscal_years: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of fiscal_years records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FiscalYears, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_years", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fiscal_years
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fiscal_years records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fiscal_year
			, start_date
			, end_date
			, status
			, is_current
			, closed_by
			, closed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fiscal_years
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fiscal_years", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fiscal_years: %w", err)
	}
	defer rows.Close()

	var entities []*FiscalYears
	for rows.Next() {
		var entity FiscalYears
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FiscalYear,
			&entity.StartDate,
			&entity.EndDate,
			&entity.Status,
			&entity.IsCurrent,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fiscal_years: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating fiscal_years rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing fiscal_years record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FiscalYears) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_years", duration, nil)
	}()

	query := `
		UPDATE fiscal_years
		SET
			, organization_id = $2
			, fiscal_year = $3
			, start_date = $4
			, end_date = $5
			, status = $6
			, is_current = $7
			, closed_by = $8
			, closed_at = $9
			, notes = $10
			, metadata = $11
			, updated_at = $13
			, created_by = $14
			, updated_by = $15
			, deleted_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.FiscalYear,
		entity.StartDate,
		entity.EndDate,
		entity.Status,
		entity.IsCurrent,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update fiscal_years", zap.Error(err))
		return fmt.Errorf("failed to update fiscal_years: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_years not found or already deleted")
	}

	r.logger.Info("updated fiscal_years",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a fiscal_years record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_years", duration, nil)
	}()

	query := `
		UPDATE fiscal_years
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete fiscal_years", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete fiscal_years: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_years not found or already deleted")
	}

	r.logger.Info("deleted fiscal_years", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves fiscal_years records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*FiscalYears, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_years", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fiscal_years
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fiscal_years records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fiscal_year
			, start_date
			, end_date
			, status
			, is_current
			, closed_by
			, closed_at
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fiscal_years
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fiscal_years by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fiscal_years: %w", err)
	}
	defer rows.Close()

	var entities []*FiscalYears
	for rows.Next() {
		var entity FiscalYears
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FiscalYear,
			&entity.StartDate,
			&entity.EndDate,
			&entity.Status,
			&entity.IsCurrent,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fiscal_years: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

