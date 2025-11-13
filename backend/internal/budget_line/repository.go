package budget_line

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

// Repository handles database operations for BudgetLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BudgetLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BudgetLines represents a budget_lines entity
type BudgetLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	BudgetId uuid.UUID `json:"budget_id" db:"budget_id"`
	AccountId *uuid.UUID `json:"account_id" db:"account_id"`
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id" db:"analytic_account_id"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	PeriodStartDate *time.Time `json:"period_start_date" db:"period_start_date"`
	PeriodEndDate *time.Time `json:"period_end_date" db:"period_end_date"`
	PlannedAmount float64 `json:"planned_amount" db:"planned_amount"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	AccountId string `json:"account_id" db:"account_id"`
}

// Create inserts a new budget_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BudgetLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "budget_lines", duration, nil)
	}()

	query := `
		INSERT INTO budget_lines (
			, budget_id
			, account_id
			, analytic_account_id
			, accounting_period_id
			, period_start_date
			, period_end_date
			, planned_amount
			, notes
			, deleted_at
			, account_id
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $12
			, $13
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.BudgetId,
		entity.AccountId,
		entity.AnalyticAccountId,
		entity.AccountingPeriodId,
		entity.PeriodStartDate,
		entity.PeriodEndDate,
		entity.PlannedAmount,
		entity.Notes,
		entity.DeletedAt,
		entity.AccountId,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create budget_lines", zap.Error(err))
		return fmt.Errorf("failed to create budget_lines: %w", err)
	}

	r.logger.Info("created budget_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a budget_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BudgetLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "budget_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, budget_id
			, account_id
			, analytic_account_id
			, accounting_period_id
			, period_start_date
			, period_end_date
			, planned_amount
			, notes
			, created_at
			, updated_at
			, deleted_at
			, account_id
		FROM budget_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BudgetLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.BudgetId,
		&entity.AccountId,
		&entity.AnalyticAccountId,
		&entity.AccountingPeriodId,
		&entity.PeriodStartDate,
		&entity.PeriodEndDate,
		&entity.PlannedAmount,
		&entity.Notes,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.AccountId,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get budget_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get budget_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of budget_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BudgetLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "budget_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM budget_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count budget_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, budget_id
			, account_id
			, analytic_account_id
			, accounting_period_id
			, period_start_date
			, period_end_date
			, planned_amount
			, notes
			, created_at
			, updated_at
			, deleted_at
			, account_id
		FROM budget_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list budget_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list budget_lines: %w", err)
	}
	defer rows.Close()

	var entities []*BudgetLines
	for rows.Next() {
		var entity BudgetLines
		err := rows.Scan(
			&entity.Id,
			&entity.BudgetId,
			&entity.AccountId,
			&entity.AnalyticAccountId,
			&entity.AccountingPeriodId,
			&entity.PeriodStartDate,
			&entity.PeriodEndDate,
			&entity.PlannedAmount,
			&entity.Notes,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.AccountId,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan budget_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating budget_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing budget_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BudgetLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "budget_lines", duration, nil)
	}()

	query := `
		UPDATE budget_lines
		SET
			, budget_id = $2
			, account_id = $3
			, analytic_account_id = $4
			, accounting_period_id = $5
			, period_start_date = $6
			, period_end_date = $7
			, planned_amount = $8
			, notes = $9
			, updated_at = $11
			, deleted_at = $12
			, account_id = $13
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.BudgetId,
		entity.AccountId,
		entity.AnalyticAccountId,
		entity.AccountingPeriodId,
		entity.PeriodStartDate,
		entity.PeriodEndDate,
		entity.PlannedAmount,
		entity.Notes,
		time.Now(),
		entity.DeletedAt,
		entity.AccountId,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update budget_lines", zap.Error(err))
		return fmt.Errorf("failed to update budget_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("budget_lines not found or already deleted")
	}

	r.logger.Info("updated budget_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a budget_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "budget_lines", duration, nil)
	}()

	query := `
		UPDATE budget_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete budget_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete budget_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("budget_lines not found or already deleted")
	}

	r.logger.Info("deleted budget_lines", zap.String("id", id.String()))
	return nil
}



