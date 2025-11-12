package journal_entry

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

// Repository handles database operations for JournalEntries
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new JournalEntries repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// JournalEntries represents a journal_entries entity
type JournalEntries struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	EntryNumber string `json:"entry_number" db:"entry_number"`
	EntryTypeId uuid.UUID `json:"entry_type_id" db:"entry_type_id"`
	EntryDate time.Time `json:"entry_date" db:"entry_date"`
	PostingDate time.Time `json:"posting_date" db:"posting_date"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	FiscalYearId *uuid.UUID `json:"fiscal_year_id" db:"fiscal_year_id"`
	Status *string `json:"status" db:"status"`
	IsPosted *bool `json:"is_posted" db:"is_posted"`
	IsReversed *bool `json:"is_reversed" db:"is_reversed"`
	ReversalEntryId *uuid.UUID `json:"reversal_entry_id" db:"reversal_entry_id"`
	SourceModule *string `json:"source_module" db:"source_module"`
	SourceDocumentType *string `json:"source_document_type" db:"source_document_type"`
	SourceDocumentId *uuid.UUID `json:"source_document_id" db:"source_document_id"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	TotalDebit *float64 `json:"total_debit" db:"total_debit"`
	TotalCredit *float64 `json:"total_credit" db:"total_credit"`
	Description string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	PostedBy *uuid.UUID `json:"posted_by" db:"posted_by"`
	PostedAt *time.Time `json:"posted_at" db:"posted_at"`
	Attachments json.RawMessage `json:"attachments" db:"attachments"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	(isPosted *string `json:"(is_posted" db:"(is_posted"`
	(ABS(totalDebit *string `json:"(ABS(total_debit" db:"(ABS(total_debit"`
}

// Create inserts a new journal_entries record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *JournalEntries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "journal_entries", duration, nil)
	}()

	query := `
		INSERT INTO journal_entries (
			, organization_id
			, entry_number
			, entry_type_id
			, entry_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, status
			, is_posted
			, is_reversed
			, reversal_entry_id
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, total_debit
			, total_credit
			, description
			, notes
			, requires_approval
			, approved_by
			, approved_at
			, posted_by
			, posted_at
			, attachments
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, (is_posted
			, (ABS(total_debit
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
			, $27
			, $30
			, $31
			, $32
			, $33
			, $34
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.EntryNumber,
		entity.EntryTypeId,
		entity.EntryDate,
		entity.PostingDate,
		entity.AccountingPeriodId,
		entity.FiscalYearId,
		entity.Status,
		entity.IsPosted,
		entity.IsReversed,
		entity.ReversalEntryId,
		entity.SourceModule,
		entity.SourceDocumentType,
		entity.SourceDocumentId,
		entity.ReferenceNumber,
		entity.TotalDebit,
		entity.TotalCredit,
		entity.Description,
		entity.Notes,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.PostedBy,
		entity.PostedAt,
		entity.Attachments,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.(isPosted,
		entity.(ABS(totalDebit,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create journal_entries", zap.Error(err))
		return fmt.Errorf("failed to create journal_entries: %w", err)
	}

	r.logger.Info("created journal_entries",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a journal_entries by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*JournalEntries, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entries", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, entry_number
			, entry_type_id
			, entry_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, status
			, is_posted
			, is_reversed
			, reversal_entry_id
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, total_debit
			, total_credit
			, description
			, notes
			, requires_approval
			, approved_by
			, approved_at
			, posted_by
			, posted_at
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (is_posted
			, (ABS(total_debit
		FROM journal_entries
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity JournalEntries
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.EntryNumber,
		&entity.EntryTypeId,
		&entity.EntryDate,
		&entity.PostingDate,
		&entity.AccountingPeriodId,
		&entity.FiscalYearId,
		&entity.Status,
		&entity.IsPosted,
		&entity.IsReversed,
		&entity.ReversalEntryId,
		&entity.SourceModule,
		&entity.SourceDocumentType,
		&entity.SourceDocumentId,
		&entity.ReferenceNumber,
		&entity.TotalDebit,
		&entity.TotalCredit,
		&entity.Description,
		&entity.Notes,
		&entity.RequiresApproval,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.PostedBy,
		&entity.PostedAt,
		&entity.Attachments,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.(isPosted,
		&entity.(ABS(totalDebit,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("journal_entries not found")
	}

	if err != nil {
		r.logger.Error("failed to get journal_entries", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get journal_entries: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of journal_entries records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*JournalEntries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journal_entries
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal_entries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, entry_number
			, entry_type_id
			, entry_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, status
			, is_posted
			, is_reversed
			, reversal_entry_id
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, total_debit
			, total_credit
			, description
			, notes
			, requires_approval
			, approved_by
			, approved_at
			, posted_by
			, posted_at
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (is_posted
			, (ABS(total_debit
		FROM journal_entries
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journal_entries", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journal_entries: %w", err)
	}
	defer rows.Close()

	var entities []*JournalEntries
	for rows.Next() {
		var entity JournalEntries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.EntryNumber,
			&entity.EntryTypeId,
			&entity.EntryDate,
			&entity.PostingDate,
			&entity.AccountingPeriodId,
			&entity.FiscalYearId,
			&entity.Status,
			&entity.IsPosted,
			&entity.IsReversed,
			&entity.ReversalEntryId,
			&entity.SourceModule,
			&entity.SourceDocumentType,
			&entity.SourceDocumentId,
			&entity.ReferenceNumber,
			&entity.TotalDebit,
			&entity.TotalCredit,
			&entity.Description,
			&entity.Notes,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.PostedBy,
			&entity.PostedAt,
			&entity.Attachments,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.(isPosted,
			&entity.(ABS(totalDebit,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal_entries: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating journal_entries rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing journal_entries record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *JournalEntries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journal_entries", duration, nil)
	}()

	query := `
		UPDATE journal_entries
		SET
			, organization_id = $2
			, entry_number = $3
			, entry_type_id = $4
			, entry_date = $5
			, posting_date = $6
			, accounting_period_id = $7
			, fiscal_year_id = $8
			, status = $9
			, is_posted = $10
			, is_reversed = $11
			, reversal_entry_id = $12
			, source_module = $13
			, source_document_type = $14
			, source_document_id = $15
			, reference_number = $16
			, total_debit = $17
			, total_credit = $18
			, description = $19
			, notes = $20
			, requires_approval = $21
			, approved_by = $22
			, approved_at = $23
			, posted_by = $24
			, posted_at = $25
			, attachments = $26
			, metadata = $27
			, updated_at = $29
			, created_by = $30
			, updated_by = $31
			, deleted_at = $32
			, (is_posted = $33
			, (ABS(total_debit = $34
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $35
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.EntryNumber,
		entity.EntryTypeId,
		entity.EntryDate,
		entity.PostingDate,
		entity.AccountingPeriodId,
		entity.FiscalYearId,
		entity.Status,
		entity.IsPosted,
		entity.IsReversed,
		entity.ReversalEntryId,
		entity.SourceModule,
		entity.SourceDocumentType,
		entity.SourceDocumentId,
		entity.ReferenceNumber,
		entity.TotalDebit,
		entity.TotalCredit,
		entity.Description,
		entity.Notes,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.PostedBy,
		entity.PostedAt,
		entity.Attachments,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.(isPosted,
		entity.(ABS(totalDebit,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update journal_entries", zap.Error(err))
		return fmt.Errorf("failed to update journal_entries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entries not found or already deleted")
	}

	r.logger.Info("updated journal_entries",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a journal_entries record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journal_entries", duration, nil)
	}()

	query := `
		UPDATE journal_entries
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete journal_entries", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete journal_entries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entries not found or already deleted")
	}

	r.logger.Info("deleted journal_entries", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves journal_entries records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*JournalEntries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journal_entries
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal_entries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, entry_number
			, entry_type_id
			, entry_date
			, posting_date
			, accounting_period_id
			, fiscal_year_id
			, status
			, is_posted
			, is_reversed
			, reversal_entry_id
			, source_module
			, source_document_type
			, source_document_id
			, reference_number
			, total_debit
			, total_credit
			, description
			, notes
			, requires_approval
			, approved_by
			, approved_at
			, posted_by
			, posted_at
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, (is_posted
			, (ABS(total_debit
		FROM journal_entries
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journal_entries by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journal_entries: %w", err)
	}
	defer rows.Close()

	var entities []*JournalEntries
	for rows.Next() {
		var entity JournalEntries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.EntryNumber,
			&entity.EntryTypeId,
			&entity.EntryDate,
			&entity.PostingDate,
			&entity.AccountingPeriodId,
			&entity.FiscalYearId,
			&entity.Status,
			&entity.IsPosted,
			&entity.IsReversed,
			&entity.ReversalEntryId,
			&entity.SourceModule,
			&entity.SourceDocumentType,
			&entity.SourceDocumentId,
			&entity.ReferenceNumber,
			&entity.TotalDebit,
			&entity.TotalCredit,
			&entity.Description,
			&entity.Notes,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.PostedBy,
			&entity.PostedAt,
			&entity.Attachments,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.(isPosted,
			&entity.(ABS(totalDebit,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal_entries: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

