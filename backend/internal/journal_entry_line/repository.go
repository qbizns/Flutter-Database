package journal_entry_line

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

// Repository handles database operations for JournalEntryLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new JournalEntryLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// JournalEntryLines represents a journal_entry_lines entity
type JournalEntryLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	JournalEntryId uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	AccountId uuid.UUID `json:"account_id" db:"account_id"`
	DebitAmount *float64 `json:"debit_amount" db:"debit_amount"`
	CreditAmount *float64 `json:"credit_amount" db:"credit_amount"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	Department *string `json:"department" db:"department"`
	ProjectCode *string `json:"project_code" db:"project_code"`
	CostCenter *string `json:"cost_center" db:"cost_center"`
	TaxCode *string `json:"tax_code" db:"tax_code"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	Description *string `json:"description" db:"description"`
	Memo *string `json:"memo" db:"memo"`
	IsReconciled *bool `json:"is_reconciled" db:"is_reconciled"`
	ReconciledAt *time.Time `json:"reconciled_at" db:"reconciled_at"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	(debitAmount *string `json:"(debit_amount" db:"(debit_amount"`
	(creditAmount *string `json:"(credit_amount" db:"(credit_amount"`
	(debitAmount *string `json:"(debit_amount" db:"(debit_amount"`
}

// Create inserts a new journal_entry_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *JournalEntryLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "journal_entry_lines", duration, nil)
	}()

	query := `
		INSERT INTO journal_entry_lines (
			, organization_id
			, journal_entry_id
			, line_number
			, account_id
			, debit_amount
			, credit_amount
			, location_id
			, department
			, project_code
			, cost_center
			, tax_code
			, tax_amount
			, description
			, memo
			, is_reconciled
			, reconciled_at
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, (debit_amount
			, (credit_amount
			, (debit_amount
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
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.JournalEntryId,
		entity.LineNumber,
		entity.AccountId,
		entity.DebitAmount,
		entity.CreditAmount,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.CostCenter,
		entity.TaxCode,
		entity.TaxAmount,
		entity.Description,
		entity.Memo,
		entity.IsReconciled,
		entity.ReconciledAt,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.(debitAmount,
		entity.(creditAmount,
		entity.(debitAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create journal_entry_lines", zap.Error(err))
		return fmt.Errorf("failed to create journal_entry_lines: %w", err)
	}

	r.logger.Info("created journal_entry_lines",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a journal_entry_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*JournalEntryLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entry_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, line_number
			, account_id
			, debit_amount
			, credit_amount
			, location_id
			, department
			, project_code
			, cost_center
			, tax_code
			, tax_amount
			, description
			, memo
			, is_reconciled
			, reconciled_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (debit_amount
			, (credit_amount
			, (debit_amount
		FROM journal_entry_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity JournalEntryLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.JournalEntryId,
		&entity.LineNumber,
		&entity.AccountId,
		&entity.DebitAmount,
		&entity.CreditAmount,
		&entity.LocationId,
		&entity.Department,
		&entity.ProjectCode,
		&entity.CostCenter,
		&entity.TaxCode,
		&entity.TaxAmount,
		&entity.Description,
		&entity.Memo,
		&entity.IsReconciled,
		&entity.ReconciledAt,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.(debitAmount,
		&entity.(creditAmount,
		&entity.(debitAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("journal_entry_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get journal_entry_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get journal_entry_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of journal_entry_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*JournalEntryLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entry_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journal_entry_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal_entry_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, line_number
			, account_id
			, debit_amount
			, credit_amount
			, location_id
			, department
			, project_code
			, cost_center
			, tax_code
			, tax_amount
			, description
			, memo
			, is_reconciled
			, reconciled_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (debit_amount
			, (credit_amount
			, (debit_amount
		FROM journal_entry_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journal_entry_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journal_entry_lines: %w", err)
	}
	defer rows.Close()

	var entities []*JournalEntryLines
	for rows.Next() {
		var entity JournalEntryLines
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalEntryId,
			&entity.LineNumber,
			&entity.AccountId,
			&entity.DebitAmount,
			&entity.CreditAmount,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.CostCenter,
			&entity.TaxCode,
			&entity.TaxAmount,
			&entity.Description,
			&entity.Memo,
			&entity.IsReconciled,
			&entity.ReconciledAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.(debitAmount,
			&entity.(creditAmount,
			&entity.(debitAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal_entry_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating journal_entry_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing journal_entry_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *JournalEntryLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journal_entry_lines", duration, nil)
	}()

	query := `
		UPDATE journal_entry_lines
		SET
			, organization_id = $2
			, journal_entry_id = $3
			, line_number = $4
			, account_id = $5
			, debit_amount = $6
			, credit_amount = $7
			, location_id = $8
			, department = $9
			, project_code = $10
			, cost_center = $11
			, tax_code = $12
			, tax_amount = $13
			, description = $14
			, memo = $15
			, is_reconciled = $16
			, reconciled_at = $17
			, metadata = $18
			, updated_at = $20
			, created_by = $21
			, updated_by = $22
			, deleted_at = $23
			, (debit_amount = $24
			, (credit_amount = $25
			, (debit_amount = $26
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.JournalEntryId,
		entity.LineNumber,
		entity.AccountId,
		entity.DebitAmount,
		entity.CreditAmount,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.CostCenter,
		entity.TaxCode,
		entity.TaxAmount,
		entity.Description,
		entity.Memo,
		entity.IsReconciled,
		entity.ReconciledAt,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.(debitAmount,
		entity.(creditAmount,
		entity.(debitAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update journal_entry_lines", zap.Error(err))
		return fmt.Errorf("failed to update journal_entry_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entry_lines not found or already deleted")
	}

	r.logger.Info("updated journal_entry_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a journal_entry_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journal_entry_lines", duration, nil)
	}()

	query := `
		UPDATE journal_entry_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete journal_entry_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete journal_entry_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entry_lines not found or already deleted")
	}

	r.logger.Info("deleted journal_entry_lines", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves journal_entry_lines records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*JournalEntryLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entry_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journal_entry_lines
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal_entry_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_entry_id
			, line_number
			, account_id
			, debit_amount
			, credit_amount
			, location_id
			, department
			, project_code
			, cost_center
			, tax_code
			, tax_amount
			, description
			, memo
			, is_reconciled
			, reconciled_at
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (debit_amount
			, (credit_amount
			, (debit_amount
		FROM journal_entry_lines
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journal_entry_lines by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journal_entry_lines: %w", err)
	}
	defer rows.Close()

	var entities []*JournalEntryLines
	for rows.Next() {
		var entity JournalEntryLines
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalEntryId,
			&entity.LineNumber,
			&entity.AccountId,
			&entity.DebitAmount,
			&entity.CreditAmount,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.CostCenter,
			&entity.TaxCode,
			&entity.TaxAmount,
			&entity.Description,
			&entity.Memo,
			&entity.IsReconciled,
			&entity.ReconciledAt,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.(debitAmount,
			&entity.(creditAmount,
			&entity.(debitAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal_entry_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

