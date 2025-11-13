package bank_statement_reconciliation

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

// Repository handles database operations for BankStatementReconciliations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankStatementReconciliations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankStatementReconciliations represents a bank_statement_reconciliations entity
type BankStatementReconciliations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BankStatementLineId uuid.UUID `json:"bank_statement_line_id" db:"bank_statement_line_id"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	PaymentId *uuid.UUID `json:"payment_id" db:"payment_id"`
	MatchedAmount float64 `json:"matched_amount" db:"matched_amount"`
	MatchedBy *uuid.UUID `json:"matched_by" db:"matched_by"`
	MatchedAt *time.Time `json:"matched_at" db:"matched_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new bank_statement_reconciliations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankStatementReconciliations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_statement_reconciliations", duration, nil)
	}()

	query := `
		INSERT INTO bank_statement_reconciliations (
			, organization_id
			, bank_statement_line_id
			, journal_entry_id
			, payment_id
			, matched_amount
			, matched_by
			, matched_at
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
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BankStatementLineId,
		entity.JournalEntryId,
		entity.PaymentId,
		entity.MatchedAmount,
		entity.MatchedBy,
		entity.MatchedAt,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create bank_statement_reconciliations", zap.Error(err))
		return fmt.Errorf("failed to create bank_statement_reconciliations: %w", err)
	}

	r.logger.Info("created bank_statement_reconciliations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a bank_statement_reconciliations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankStatementReconciliations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statement_reconciliations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, bank_statement_line_id
			, journal_entry_id
			, payment_id
			, matched_amount
			, matched_by
			, matched_at
			, deleted_at
		FROM bank_statement_reconciliations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BankStatementReconciliations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BankStatementLineId,
		&entity.JournalEntryId,
		&entity.PaymentId,
		&entity.MatchedAmount,
		&entity.MatchedBy,
		&entity.MatchedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_statement_reconciliations not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_statement_reconciliations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_statement_reconciliations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_statement_reconciliations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankStatementReconciliations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statement_reconciliations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_statement_reconciliations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_statement_reconciliations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_statement_line_id
			, journal_entry_id
			, payment_id
			, matched_amount
			, matched_by
			, matched_at
			, deleted_at
		FROM bank_statement_reconciliations
		WHERE deleted_at IS NULL
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_statement_reconciliations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_statement_reconciliations: %w", err)
	}
	defer rows.Close()

	var entities []*BankStatementReconciliations
	for rows.Next() {
		var entity BankStatementReconciliations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankStatementLineId,
			&entity.JournalEntryId,
			&entity.PaymentId,
			&entity.MatchedAmount,
			&entity.MatchedBy,
			&entity.MatchedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_statement_reconciliations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_statement_reconciliations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_statement_reconciliations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankStatementReconciliations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statement_reconciliations", duration, nil)
	}()

	query := `
		UPDATE bank_statement_reconciliations
		SET
			, organization_id = $2
			, bank_statement_line_id = $3
			, journal_entry_id = $4
			, payment_id = $5
			, matched_amount = $6
			, matched_by = $7
			, matched_at = $8
			, deleted_at = $9
			
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BankStatementLineId,
		entity.JournalEntryId,
		entity.PaymentId,
		entity.MatchedAmount,
		entity.MatchedBy,
		entity.MatchedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_statement_reconciliations", zap.Error(err))
		return fmt.Errorf("failed to update bank_statement_reconciliations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statement_reconciliations not found or already deleted")
	}

	r.logger.Info("updated bank_statement_reconciliations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a bank_statement_reconciliations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statement_reconciliations", duration, nil)
	}()

	query := `
		UPDATE bank_statement_reconciliations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_statement_reconciliations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_statement_reconciliations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statement_reconciliations not found or already deleted")
	}

	r.logger.Info("deleted bank_statement_reconciliations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves bank_statement_reconciliations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BankStatementReconciliations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statement_reconciliations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_statement_reconciliations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_statement_reconciliations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_statement_line_id
			, journal_entry_id
			, payment_id
			, matched_amount
			, matched_by
			, matched_at
			, deleted_at
		FROM bank_statement_reconciliations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_statement_reconciliations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_statement_reconciliations: %w", err)
	}
	defer rows.Close()

	var entities []*BankStatementReconciliations
	for rows.Next() {
		var entity BankStatementReconciliations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankStatementLineId,
			&entity.JournalEntryId,
			&entity.PaymentId,
			&entity.MatchedAmount,
			&entity.MatchedBy,
			&entity.MatchedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_statement_reconciliations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

