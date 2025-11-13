package accounting_period

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

// Repository handles database operations for AccountingPeriods
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AccountingPeriods repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AccountingPeriods represents a accounting_periods entity
type AccountingPeriods struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	FiscalYearId uuid.UUID `json:"fiscal_year_id" db:"fiscal_year_id"`
	PeriodNumber int64 `json:"period_number" db:"period_number"`
	PeriodName string `json:"period_name" db:"period_name"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate time.Time `json:"end_date" db:"end_date"`
	Status *string `json:"status" db:"status"`
	ClosedBy *uuid.UUID `json:"closed_by" db:"closed_by"`
	ClosedAt *time.Time `json:"closed_at" db:"closed_at"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new accounting_periods record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AccountingPeriods) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "accounting_periods", duration, nil)
	}()

	query := `
		INSERT INTO accounting_periods (
			, organization_id
			, fiscal_year_id
			, period_number
			, period_name
			, start_date
			, end_date
			, status
			, closed_by
			, closed_at
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
		entity.FiscalYearId,
		entity.PeriodNumber,
		entity.PeriodName,
		entity.StartDate,
		entity.EndDate,
		entity.Status,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create accounting_periods", zap.Error(err))
		return fmt.Errorf("failed to create accounting_periods: %w", err)
	}

	r.logger.Info("created accounting_periods",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a accounting_periods by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AccountingPeriods, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "accounting_periods", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, fiscal_year_id
			, period_number
			, period_name
			, start_date
			, end_date
			, status
			, closed_by
			, closed_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM accounting_periods
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity AccountingPeriods
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.FiscalYearId,
		&entity.PeriodNumber,
		&entity.PeriodName,
		&entity.StartDate,
		&entity.EndDate,
		&entity.Status,
		&entity.ClosedBy,
		&entity.ClosedAt,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("accounting_periods not found")
	}

	if err != nil {
		r.logger.Error("failed to get accounting_periods", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get accounting_periods: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of accounting_periods records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AccountingPeriods, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "accounting_periods", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM accounting_periods
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count accounting_periods records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fiscal_year_id
			, period_number
			, period_name
			, start_date
			, end_date
			, status
			, closed_by
			, closed_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM accounting_periods
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list accounting_periods", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list accounting_periods: %w", err)
	}
	defer rows.Close()

	var entities []*AccountingPeriods
	for rows.Next() {
		var entity AccountingPeriods
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FiscalYearId,
			&entity.PeriodNumber,
			&entity.PeriodName,
			&entity.StartDate,
			&entity.EndDate,
			&entity.Status,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan accounting_periods: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating accounting_periods rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing accounting_periods record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AccountingPeriods) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "accounting_periods", duration, nil)
	}()

	query := `
		UPDATE accounting_periods
		SET
			, organization_id = $2
			, fiscal_year_id = $3
			, period_number = $4
			, period_name = $5
			, start_date = $6
			, end_date = $7
			, status = $8
			, closed_by = $9
			, closed_at = $10
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
		entity.FiscalYearId,
		entity.PeriodNumber,
		entity.PeriodName,
		entity.StartDate,
		entity.EndDate,
		entity.Status,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update accounting_periods", zap.Error(err))
		return fmt.Errorf("failed to update accounting_periods: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("accounting_periods not found or already deleted")
	}

	r.logger.Info("updated accounting_periods",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a accounting_periods record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "accounting_periods", duration, nil)
	}()

	query := `
		UPDATE accounting_periods
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete accounting_periods", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete accounting_periods: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("accounting_periods not found or already deleted")
	}

	r.logger.Info("deleted accounting_periods", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves accounting_periods records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*AccountingPeriods, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "accounting_periods", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM accounting_periods
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count accounting_periods records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fiscal_year_id
			, period_number
			, period_name
			, start_date
			, end_date
			, status
			, closed_by
			, closed_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM accounting_periods
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list accounting_periods by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list accounting_periods: %w", err)
	}
	defer rows.Close()

	var entities []*AccountingPeriods
	for rows.Next() {
		var entity AccountingPeriods
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FiscalYearId,
			&entity.PeriodNumber,
			&entity.PeriodName,
			&entity.StartDate,
			&entity.EndDate,
			&entity.Status,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan accounting_periods: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

