package budget

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

// Repository handles database operations for Budgets
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Budgets repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Budgets represents a budgets entity
type Budgets struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BudgetCode string `json:"budget_code" db:"budget_code"`
	BudgetName string `json:"budget_name" db:"budget_name"`
	FiscalYearId *uuid.UUID `json:"fiscal_year_id" db:"fiscal_year_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate time.Time `json:"end_date" db:"end_date"`
	BudgetType *string `json:"budget_type" db:"budget_type"`
	BudgetType *string `json:"budget_type" db:"budget_type"`
	Status *string `json:"status" db:"status"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new budgets record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Budgets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "budgets", duration, nil)
	}()

	query := `
		INSERT INTO budgets (
			, organization_id
			, budget_code
			, budget_name
			, fiscal_year_id
			, start_date
			, end_date
			, budget_type
			, budget_type
			, status
			, notes
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
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BudgetCode,
		entity.BudgetName,
		entity.FiscalYearId,
		entity.StartDate,
		entity.EndDate,
		entity.BudgetType,
		entity.BudgetType,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create budgets", zap.Error(err))
		return fmt.Errorf("failed to create budgets: %w", err)
	}

	r.logger.Info("created budgets",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a budgets by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Budgets, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "budgets", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, budget_code
			, budget_name
			, fiscal_year_id
			, start_date
			, end_date
			, budget_type
			, budget_type
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM budgets
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Budgets
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BudgetCode,
		&entity.BudgetName,
		&entity.FiscalYearId,
		&entity.StartDate,
		&entity.EndDate,
		&entity.BudgetType,
		&entity.BudgetType,
		&entity.Status,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budgets not found")
	}

	if err != nil {
		r.logger.Error("failed to get budgets", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of budgets records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Budgets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "budgets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM budgets
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budgets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, budget_code
			, budget_name
			, fiscal_year_id
			, start_date
			, end_date
			, budget_type
			, budget_type
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM budgets
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list budgets", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list budgets: %w", err)
	}
	defer rows.Close()

	var entities []*Budgets
	for rows.Next() {
		var entity Budgets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BudgetCode,
			&entity.BudgetName,
			&entity.FiscalYearId,
			&entity.StartDate,
			&entity.EndDate,
			&entity.BudgetType,
			&entity.BudgetType,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budgets: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating budgets rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing budgets record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Budgets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "budgets", duration, nil)
	}()

	query := `
		UPDATE budgets
		SET
			, organization_id = $2
			, budget_code = $3
			, budget_name = $4
			, fiscal_year_id = $5
			, start_date = $6
			, end_date = $7
			, budget_type = $8
			, budget_type = $9
			, status = $10
			, notes = $11
			, created_by = $12
			, updated_by = $13
			, updated_at = $15
			, deleted_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BudgetCode,
		entity.BudgetName,
		entity.FiscalYearId,
		entity.StartDate,
		entity.EndDate,
		entity.BudgetType,
		entity.BudgetType,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update budgets", zap.Error(err))
		return fmt.Errorf("failed to update budgets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("budgets not found or already deleted")
	}

	r.logger.Info("updated budgets",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a budgets record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "budgets", duration, nil)
	}()

	query := `
		UPDATE budgets
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete budgets", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete budgets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("budgets not found or already deleted")
	}

	r.logger.Info("deleted budgets", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves budgets records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Budgets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "budgets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM budgets
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budgets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, budget_code
			, budget_name
			, fiscal_year_id
			, start_date
			, end_date
			, budget_type
			, budget_type
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM budgets
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list budgets by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list budgets: %w", err)
	}
	defer rows.Close()

	var entities []*Budgets
	for rows.Next() {
		var entity Budgets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BudgetCode,
			&entity.BudgetName,
			&entity.FiscalYearId,
			&entity.StartDate,
			&entity.EndDate,
			&entity.BudgetType,
			&entity.BudgetType,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budgets: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

