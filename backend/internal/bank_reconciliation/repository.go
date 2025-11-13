package bank_reconciliation

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

// Repository handles database operations for BankReconciliations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankReconciliations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankReconciliations represents a bank_reconciliations entity
type BankReconciliations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BankAccountId uuid.UUID `json:"bank_account_id" db:"bank_account_id"`
	StatementDate time.Time `json:"statement_date" db:"statement_date"`
	StatementBalance float64 `json:"statement_balance" db:"statement_balance"`
	ReconciliationDate *time.Time `json:"reconciliation_date" db:"reconciliation_date"`
	BookBalance *float64 `json:"book_balance" db:"book_balance"`
	ClearedBalance *float64 `json:"cleared_balance" db:"cleared_balance"`
	Difference *float64 `json:"difference" db:"difference"`
	Status *string `json:"status" db:"status"`
	IsReconciled *bool `json:"is_reconciled" db:"is_reconciled"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	ReconciledBy *uuid.UUID `json:"reconciled_by" db:"reconciled_by"`
	ReconciledAt *time.Time `json:"reconciled_at" db:"reconciled_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new bank_reconciliations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankReconciliations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_reconciliations", duration, nil)
	}()

	query := `
		INSERT INTO bank_reconciliations (
			, organization_id
			, bank_account_id
			, statement_date
			, statement_balance
			, reconciliation_date
			, book_balance
			, cleared_balance
			, difference
			, status
			, is_reconciled
			, accounting_period_id
			, notes
			, metadata
			, reconciled_by
			, reconciled_at
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
			, $17
			, $18
			, $19
			, $20
			, $21
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BankAccountId,
		entity.StatementDate,
		entity.StatementBalance,
		entity.ReconciliationDate,
		entity.BookBalance,
		entity.ClearedBalance,
		entity.Difference,
		entity.Status,
		entity.IsReconciled,
		entity.AccountingPeriodId,
		entity.Notes,
		entity.Metadata,
		entity.ReconciledBy,
		entity.ReconciledAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create bank_reconciliations", zap.Error(err))
		return fmt.Errorf("failed to create bank_reconciliations: %w", err)
	}

	r.logger.Info("created bank_reconciliations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a bank_reconciliations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankReconciliations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_date
			, statement_balance
			, reconciliation_date
			, book_balance
			, cleared_balance
			, difference
			, is_reconciled
			, accounting_period_id
			, notes
			, metadata
			, created_at
			, updated_at
			, reconciled_by
			, reconciled_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_reconciliations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BankReconciliations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BankAccountId,
		&entity.StatementDate,
		&entity.StatementBalance,
		&entity.ReconciliationDate,
		&entity.BookBalance,
		&entity.ClearedBalance,
		&entity.Difference,
		&entity.Status,
		&entity.IsReconciled,
		&entity.AccountingPeriodId,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.ReconciledBy,
		&entity.ReconciledAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_reconciliations not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_reconciliations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_reconciliations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_reconciliations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankReconciliations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_reconciliations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_reconciliations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_date
			, statement_balance
			, reconciliation_date
			, book_balance
			, cleared_balance
			, difference
			, is_reconciled
			, accounting_period_id
			, notes
			, metadata
			, created_at
			, updated_at
			, reconciled_by
			, reconciled_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_reconciliations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_reconciliations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_reconciliations: %w", err)
	}
	defer rows.Close()

	var entities []*BankReconciliations
	for rows.Next() {
		var entity BankReconciliations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankAccountId,
			&entity.StatementDate,
			&entity.StatementBalance,
			&entity.ReconciliationDate,
			&entity.BookBalance,
			&entity.ClearedBalance,
			&entity.Difference,
			&entity.Status,
			&entity.IsReconciled,
			&entity.AccountingPeriodId,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.ReconciledBy,
			&entity.ReconciledAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_reconciliations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_reconciliations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_reconciliations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankReconciliations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_reconciliations", duration, nil)
	}()

	query := `
		UPDATE bank_reconciliations
		SET
			, organization_id = $2
			, bank_account_id = $3
			, statement_date = $4
			, statement_balance = $5
			, reconciliation_date = $6
			, book_balance = $7
			, cleared_balance = $8
			, difference = $9
			, status = $10
			, is_reconciled = $11
			, accounting_period_id = $12
			, notes = $13
			, metadata = $14
			, updated_at = $16
			, reconciled_by = $17
			, reconciled_at = $18
			, created_by = $19
			, updated_by = $20
			, deleted_at = $21
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $22
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BankAccountId,
		entity.StatementDate,
		entity.StatementBalance,
		entity.ReconciliationDate,
		entity.BookBalance,
		entity.ClearedBalance,
		entity.Difference,
		entity.Status,
		entity.IsReconciled,
		entity.AccountingPeriodId,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.ReconciledBy,
		entity.ReconciledAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_reconciliations", zap.Error(err))
		return fmt.Errorf("failed to update bank_reconciliations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_reconciliations not found or already deleted")
	}

	r.logger.Info("updated bank_reconciliations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a bank_reconciliations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_reconciliations", duration, nil)
	}()

	query := `
		UPDATE bank_reconciliations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_reconciliations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_reconciliations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_reconciliations not found or already deleted")
	}

	r.logger.Info("deleted bank_reconciliations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves bank_reconciliations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BankReconciliations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_reconciliations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_reconciliations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_reconciliations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_date
			, statement_balance
			, reconciliation_date
			, book_balance
			, cleared_balance
			, difference
			, status
			, is_reconciled
			, accounting_period_id
			, notes
			, metadata
			, created_at
			, updated_at
			, reconciled_by
			, reconciled_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_reconciliations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_reconciliations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_reconciliations: %w", err)
	}
	defer rows.Close()

	var entities []*BankReconciliations
	for rows.Next() {
		var entity BankReconciliations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankAccountId,
			&entity.StatementDate,
			&entity.StatementBalance,
			&entity.ReconciliationDate,
			&entity.BookBalance,
			&entity.ClearedBalance,
			&entity.Difference,
			&entity.Status,
			&entity.IsReconciled,
			&entity.AccountingPeriodId,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.ReconciledBy,
			&entity.ReconciledAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_reconciliations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

