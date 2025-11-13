package general_ledger

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

// Repository handles database operations for GeneralLedger
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new GeneralLedger repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// GeneralLedger represents a general_ledger entity
type GeneralLedger struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	JournalEntryId uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id" db:"journal_entry_line_id"`
	AccountId uuid.UUID `json:"account_id" db:"account_id"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	PostingDate time.Time `json:"posting_date" db:"posting_date"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	FiscalYearId *uuid.UUID `json:"fiscal_year_id" db:"fiscal_year_id"`
	DebitAmount *float64 `json:"debit_amount" db:"debit_amount"`
	CreditAmount *float64 `json:"credit_amount" db:"credit_amount"`
	RunningDebitBalance *float64 `json:"running_debit_balance" db:"running_debit_balance"`
	RunningCreditBalance *float64 `json:"running_credit_balance" db:"running_credit_balance"`
	RunningBalance *float64 `json:"running_balance" db:"running_balance"`
	SourceModule *string `json:"source_module" db:"source_module"`
	SourceDocumentType *string `json:"source_document_type" db:"source_document_type"`
	SourceDocumentId *uuid.UUID `json:"source_document_id" db:"source_document_id"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	Department *string `json:"department" db:"department"`
	ProjectCode *string `json:"project_code" db:"project_code"`
	CostCenter *string `json:"cost_center" db:"cost_center"`
	Description *string `json:"description" db:"description"`
	IsReversed *bool `json:"is_reversed" db:"is_reversed"`
	ReversalGlId *uuid.UUID `json:"reversal_gl_id" db:"reversal_gl_id"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	(debitAmount *string `json:"(debit_amount" db:"(debit_amount"`
	(creditAmount *string `json:"(credit_amount" db:"(credit_amount"`
}

// Create inserts a new general_ledger record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *GeneralLedger) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "general_ledger", duration, nil)
	}()

	query := `
		INSERT INTO general_ledger (
			, organization_id
			, journal_entry_id
			, journal_entry_line_id
			, account_id
			, transaction_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, debit_amount
			, credit_amount
			, running_debit_balance
			, running_credit_balance
			, running_balance
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, location_id
			, department
			, project_code
			, cost_center
			, description
			, is_reversed
			, reversal_gl_id
			, metadata
			, created_by
			, (debit_amount
			, (credit_amount
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
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $28
			, $29
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.JournalEntryId,
		entity.JournalEntryLineId,
		entity.AccountId,
		entity.TransactionDate,
		entity.PostingDate,
		entity.AccountingPeriodId,
		entity.FiscalYearId,
		entity.DebitAmount,
		entity.CreditAmount,
		entity.RunningDebitBalance,
		entity.RunningCreditBalance,
		entity.RunningBalance,
		entity.SourceModule,
		entity.SourceDocumentType,
		entity.SourceDocumentId,
		entity.ReferenceNumber,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.CostCenter,
		entity.Description,
		entity.IsReversed,
		entity.ReversalGlId,
		entity.Metadata,
		entity.CreatedBy,
		entity.(debitAmount,
		entity.(creditAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create general_ledger", zap.Error(err))
		return fmt.Errorf("failed to create general_ledger: %w", err)
	}

	r.logger.Info("created general_ledger",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a general_ledger by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*GeneralLedger, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "general_ledger", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, journal_entry_line_id
			, account_id
			, transaction_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, debit_amount
			, credit_amount
			, running_debit_balance
			, running_credit_balance
			, running_balance
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, location_id
			, department
			, project_code
			, cost_center
			, description
			, is_reversed
			, reversal_gl_id
			, metadata
			, created_at
			, created_by
			, (debit_amount
			, (credit_amount
		FROM general_ledger
		WHERE id = $1
		
	`

	var entity GeneralLedger
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.JournalEntryId,
		&entity.JournalEntryLineId,
		&entity.AccountId,
		&entity.TransactionDate,
		&entity.PostingDate,
		&entity.AccountingPeriodId,
		&entity.FiscalYearId,
		&entity.DebitAmount,
		&entity.CreditAmount,
		&entity.RunningDebitBalance,
		&entity.RunningCreditBalance,
		&entity.RunningBalance,
		&entity.SourceModule,
		&entity.SourceDocumentType,
		&entity.SourceDocumentId,
		&entity.ReferenceNumber,
		&entity.LocationId,
		&entity.Department,
		&entity.ProjectCode,
		&entity.CostCenter,
		&entity.Description,
		&entity.IsReversed,
		&entity.ReversalGlId,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.CreatedBy,
		&entity.(debitAmount,
		&entity.(creditAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("general_ledger not found")
	}

	if err != nil {
		r.logger.Error("failed to get general_ledger", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get general_ledger: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of general_ledger records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*GeneralLedger, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "general_ledger", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM general_ledger
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count general_ledger records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, journal_entry_line_id
			, account_id
			, transaction_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, debit_amount
			, credit_amount
			, running_debit_balance
			, running_credit_balance
			, running_balance
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, location_id
			, department
			, project_code
			, cost_center
			, description
			, is_reversed
			, reversal_gl_id
			, metadata
			, created_at
			, created_by
			, (debit_amount
			, (credit_amount
		FROM general_ledger
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list general_ledger", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list general_ledger: %w", err)
	}
	defer rows.Close()

	var entities []*GeneralLedger
	for rows.Next() {
		var entity GeneralLedger
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalEntryId,
			&entity.JournalEntryLineId,
			&entity.AccountId,
			&entity.TransactionDate,
			&entity.PostingDate,
			&entity.AccountingPeriodId,
			&entity.FiscalYearId,
			&entity.DebitAmount,
			&entity.CreditAmount,
			&entity.RunningDebitBalance,
			&entity.RunningCreditBalance,
			&entity.RunningBalance,
			&entity.SourceModule,
			&entity.SourceDocumentType,
			&entity.SourceDocumentId,
			&entity.ReferenceNumber,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.CostCenter,
			&entity.Description,
			&entity.IsReversed,
			&entity.ReversalGlId,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.(debitAmount,
			&entity.(creditAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan general_ledger: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating general_ledger rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing general_ledger record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *GeneralLedger) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "general_ledger", duration, nil)
	}()

	query := `
		UPDATE general_ledger
		SET
			, organization_id = $2
			, journal_entry_id = $3
			, journal_entry_line_id = $4
			, account_id = $5
			, transaction_date = $6
			, posting_date = $7
			, accounting_period_id = $8
			, fiscal_year_id = $9
			, debit_amount = $10
			, credit_amount = $11
			, running_debit_balance = $12
			, running_credit_balance = $13
			, running_balance = $14
			, source_module = $15
			, source_document_type = $16
			, source_document_id = $17
			, reference_number = $18
			, location_id = $19
			, department = $20
			, project_code = $21
			, cost_center = $22
			, description = $23
			, is_reversed = $24
			, reversal_gl_id = $25
			, metadata = $26
			, created_by = $28
			, (debit_amount = $29
			, (credit_amount = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.JournalEntryId,
		entity.JournalEntryLineId,
		entity.AccountId,
		entity.TransactionDate,
		entity.PostingDate,
		entity.AccountingPeriodId,
		entity.FiscalYearId,
		entity.DebitAmount,
		entity.CreditAmount,
		entity.RunningDebitBalance,
		entity.RunningCreditBalance,
		entity.RunningBalance,
		entity.SourceModule,
		entity.SourceDocumentType,
		entity.SourceDocumentId,
		entity.ReferenceNumber,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.CostCenter,
		entity.Description,
		entity.IsReversed,
		entity.ReversalGlId,
		entity.Metadata,
		entity.CreatedBy,
		entity.(debitAmount,
		entity.(creditAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update general_ledger", zap.Error(err))
		return fmt.Errorf("failed to update general_ledger: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("general_ledger not found or already deleted")
	}

	r.logger.Info("updated general_ledger",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a general_ledger record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "general_ledger", duration, nil)
	}()

	query := `DELETE FROM general_ledger WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete general_ledger", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete general_ledger: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("general_ledger not found")
	}

	r.logger.Info("deleted general_ledger", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves general_ledger records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*GeneralLedger, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "general_ledger", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM general_ledger
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count general_ledger records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, journal_entry_line_id
			, account_id
			, transaction_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, debit_amount
			, credit_amount
			, running_debit_balance
			, running_credit_balance
			, running_balance
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, location_id
			, department
			, project_code
			, cost_center
			, description
			, is_reversed
			, reversal_gl_id
			, metadata
			, created_at
			, created_by
			, (debit_amount
			, (credit_amount
		FROM general_ledger
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list general_ledger by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list general_ledger: %w", err)
	}
	defer rows.Close()

	var entities []*GeneralLedger
	for rows.Next() {
		var entity GeneralLedger
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalEntryId,
			&entity.JournalEntryLineId,
			&entity.AccountId,
			&entity.TransactionDate,
			&entity.PostingDate,
			&entity.AccountingPeriodId,
			&entity.FiscalYearId,
			&entity.DebitAmount,
			&entity.CreditAmount,
			&entity.RunningDebitBalance,
			&entity.RunningCreditBalance,
			&entity.RunningBalance,
			&entity.SourceModule,
			&entity.SourceDocumentType,
			&entity.SourceDocumentId,
			&entity.ReferenceNumber,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.CostCenter,
			&entity.Description,
			&entity.IsReversed,
			&entity.ReversalGlId,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.(debitAmount,
			&entity.(creditAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan general_ledger: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

