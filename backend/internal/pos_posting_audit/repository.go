package pos_posting_audit

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

// Repository handles database operations for PosPostingAudit
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PosPostingAudit repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PosPostingAudit represents a pos_posting_audit entity
type PosPostingAudit struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SourceTable string `json:"source_table" db:"source_table"`
	SourceId uuid.UUID `json:"source_id" db:"source_id"`
	SourceReference *string `json:"source_reference" db:"source_reference"`
	PostingStatus string `json:"posting_status" db:"posting_status"`
	'pending', *string `json:"'pending'," db:"'pending',"`
	'processing', *string `json:"'processing'," db:"'processing',"`
	'posted', *string `json:"'posted'," db:"'posted',"`
	'failed', *string `json:"'failed'," db:"'failed',"`
	'cancelled', *string `json:"'cancelled'," db:"'cancelled',"`
	'reversed' *string `json:"'reversed'" db:"'reversed'"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	ReversalJournalEntryId *uuid.UUID `json:"reversal_journal_entry_id" db:"reversal_journal_entry_id"`
	PostingDate *time.Time `json:"posting_date" db:"posting_date"`
	PostedAt *time.Time `json:"posted_at" db:"posted_at"`
	PostedBy *uuid.UUID `json:"posted_by" db:"posted_by"`
	PostingMethod *string `json:"posting_method" db:"posting_method"`
	ErrorCode *string `json:"error_code" db:"error_code"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	ErrorDetails json.RawMessage `json:"error_details" db:"error_details"`
	RetryCount *int64 `json:"retry_count" db:"retry_count"`
	LastRetryAt *time.Time `json:"last_retry_at" db:"last_retry_at"`
	MaxRetries *int64 `json:"max_retries" db:"max_retries"`
	ReversedAt *time.Time `json:"reversed_at" db:"reversed_at"`
	ReversedBy *uuid.UUID `json:"reversed_by" db:"reversed_by"`
	ReversalReason *string `json:"reversal_reason" db:"reversal_reason"`
	TotalDebit *float64 `json:"total_debit" db:"total_debit"`
	TotalCredit *float64 `json:"total_credit" db:"total_credit"`
	LineCount *int64 `json:"line_count" db:"line_count"`
	CurrencyCode *string `json:"currency_code" db:"currency_code"`
	PostingContext json.RawMessage `json:"posting_context" db:"posting_context"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	// 	(postingStatus string `json:"(posting_status" db:"(posting_status"`
	// 	(postingStatus *string `json:"(posting_status" db:"(posting_status"`
	// 	(postingStatus string `json:"(posting_status" db:"(posting_status"`
	// 	(postingStatus *string `json:"(posting_status" db:"(posting_status"`
	// 	(postingStatus *string `json:"(posting_status" db:"(posting_status"`
	(ABS(COALESCE(totalDebit, *string `json:"(ABS(COALESCE(total_debit," db:"(ABS(COALESCE(total_debit,"`
}

// Create inserts a new pos_posting_audit record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PosPostingAudit) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "pos_posting_audit", duration, nil)
	}()

	query := `
		INSERT INTO pos_posting_audit (
			, organization_id
			, source_table
			, source_id
			, source_reference
			, posting_status
			, 'pending',
			, 'processing',
			, 'posted',
			, 'failed',
			, 'cancelled',
			, 'reversed'
			, journal_entry_id
			, reversal_journal_entry_id
			, posting_date
			, posted_at
			, posted_by
			, posting_method
			, error_code
			, error_message
			, error_details
			, retry_count
			, last_retry_at
			, max_retries
			, reversed_at
			, reversed_by
			, reversal_reason
			, total_debit
			, total_credit
			, line_count
			, currency_code
			, posting_context
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (ABS(COALESCE(total_debit,
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
			, $28
			, $29
			, $30
			, $31
			, $32
			, $33
			, $34
			, $37
			, $38
			, $39
			, $40
			, $41
			, $42
			, $43
			, $44
			, $45
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SourceTable,
		entity.SourceId,
		entity.SourceReference,
		entity.PostingStatus,
		entity.'pending',,
		entity.'processing',,
		entity.'posted',,
		entity.'failed',,
		entity.'cancelled',,
		entity.'reversed',
		entity.JournalEntryId,
		entity.ReversalJournalEntryId,
		entity.PostingDate,
		entity.PostedAt,
		entity.PostedBy,
		entity.PostingMethod,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.RetryCount,
		entity.LastRetryAt,
		entity.MaxRetries,
		entity.ReversedAt,
		entity.ReversedBy,
		entity.ReversalReason,
		entity.TotalDebit,
		entity.TotalCredit,
		entity.LineCount,
		entity.CurrencyCode,
		entity.PostingContext,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(ABS(COALESCE(totalDebit,,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create pos_posting_audit", zap.Error(err))
		return fmt.Errorf("failed to create pos_posting_audit: %w", err)
	}

	r.logger.Info("created pos_posting_audit",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a pos_posting_audit by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PosPostingAudit, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_posting_audit", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, source_reference
			, posting_status
			, 'pending',
			, 'processing',
			, 'posted',
			, 'failed',
			, 'cancelled',
			, 'reversed'
			, journal_entry_id
			, reversal_journal_entry_id
			, posting_date
			, posted_at
			, posted_by
			, posting_method
			, error_code
			, error_message
			, error_details
			, retry_count
			, last_retry_at
			, max_retries
			, reversed_at
			, reversed_by
			, reversal_reason
			, total_debit
			, total_credit
			, line_count
			, currency_code
			, posting_context
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (ABS(COALESCE(total_debit,
		FROM pos_posting_audit
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PosPostingAudit
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SourceTable,
		&entity.SourceId,
		&entity.SourceReference,
		&entity.PostingStatus,
		&entity.'pending',,
		&entity.'processing',,
		&entity.'posted',,
		&entity.'failed',,
		&entity.'cancelled',,
		&entity.'reversed',
		&entity.JournalEntryId,
		&entity.ReversalJournalEntryId,
		&entity.PostingDate,
		&entity.PostedAt,
		&entity.PostedBy,
		&entity.PostingMethod,
		&entity.ErrorCode,
		&entity.ErrorMessage,
		&entity.ErrorDetails,
		&entity.RetryCount,
		&entity.LastRetryAt,
		&entity.MaxRetries,
		&entity.ReversedAt,
		&entity.ReversedBy,
		&entity.ReversalReason,
		&entity.TotalDebit,
		&entity.TotalCredit,
		&entity.LineCount,
		&entity.CurrencyCode,
		&entity.PostingContext,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.(postingStatus,
		&entity.(postingStatus,
		&entity.(postingStatus,
		&entity.(postingStatus,
		&entity.(postingStatus,
		&entity.(ABS(COALESCE(totalDebit,,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("pos_posting_audit not found")
	}

	if err != nil {
		r.logger.Error("failed to get pos_posting_audit", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get pos_posting_audit: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of pos_posting_audit records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PosPostingAudit, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_posting_audit", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_posting_audit
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_posting_audit records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, source_reference
			, posting_status
			, 'pending',
			, 'processing',
			, 'posted',
			, 'failed',
			, 'cancelled',
			, 'reversed'
			, journal_entry_id
			, reversal_journal_entry_id
			, posting_date
			, posted_at
			, posted_by
			, posting_method
			, error_code
			, error_message
			, error_details
			, retry_count
			, last_retry_at
			, max_retries
			, reversed_at
			, reversed_by
			, reversal_reason
			, total_debit
			, total_credit
			, line_count
			, currency_code
			, posting_context
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (ABS(COALESCE(total_debit,
		FROM pos_posting_audit
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_posting_audit", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_posting_audit: %w", err)
	}
	defer rows.Close()

	var entities []*PosPostingAudit
	for rows.Next() {
		var entity PosPostingAudit
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceTable,
			&entity.SourceId,
			&entity.SourceReference,
			&entity.PostingStatus,
			&entity.'pending',,
			&entity.'processing',,
			&entity.'posted',,
			&entity.'failed',,
			&entity.'cancelled',,
			&entity.'reversed',
			&entity.JournalEntryId,
			&entity.ReversalJournalEntryId,
			&entity.PostingDate,
			&entity.PostedAt,
			&entity.PostedBy,
			&entity.PostingMethod,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.RetryCount,
			&entity.LastRetryAt,
			&entity.MaxRetries,
			&entity.ReversedAt,
			&entity.ReversedBy,
			&entity.ReversalReason,
			&entity.TotalDebit,
			&entity.TotalCredit,
			&entity.LineCount,
			&entity.CurrencyCode,
			&entity.PostingContext,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(ABS(COALESCE(totalDebit,,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_posting_audit: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating pos_posting_audit rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing pos_posting_audit record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PosPostingAudit) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_posting_audit", duration, nil)
	}()

	query := `
		UPDATE pos_posting_audit
		SET
			, organization_id = $2
			, source_table = $3
			, source_id = $4
			, source_reference = $5
			, posting_status = $6
			, 'pending', = $7
			, 'processing', = $8
			, 'posted', = $9
			, 'failed', = $10
			, 'cancelled', = $11
			, 'reversed' = $12
			, journal_entry_id = $13
			, reversal_journal_entry_id = $14
			, posting_date = $15
			, posted_at = $16
			, posted_by = $17
			, posting_method = $18
			, error_code = $19
			, error_message = $20
			, error_details = $21
			, retry_count = $22
			, last_retry_at = $23
			, max_retries = $24
			, reversed_at = $25
			, reversed_by = $26
			, reversal_reason = $27
			, total_debit = $28
			, total_credit = $29
			, line_count = $30
			, currency_code = $31
			, posting_context = $32
			, notes = $33
			, metadata = $34
			, updated_at = $36
			, deleted_at = $37
			, created_by = $38
			, updated_by = $39
			, (posting_status = $40
			, (posting_status = $41
			, (posting_status = $42
			, (posting_status = $43
			, (posting_status = $44
			, (ABS(COALESCE(total_debit, = $45
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $46
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SourceTable,
		entity.SourceId,
		entity.SourceReference,
		entity.PostingStatus,
		entity.'pending',,
		entity.'processing',,
		entity.'posted',,
		entity.'failed',,
		entity.'cancelled',,
		entity.'reversed',
		entity.JournalEntryId,
		entity.ReversalJournalEntryId,
		entity.PostingDate,
		entity.PostedAt,
		entity.PostedBy,
		entity.PostingMethod,
		entity.ErrorCode,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.RetryCount,
		entity.LastRetryAt,
		entity.MaxRetries,
		entity.ReversedAt,
		entity.ReversedBy,
		entity.ReversalReason,
		entity.TotalDebit,
		entity.TotalCredit,
		entity.LineCount,
		entity.CurrencyCode,
		entity.PostingContext,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(postingStatus,
		entity.(ABS(COALESCE(totalDebit,,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update pos_posting_audit", zap.Error(err))
		return fmt.Errorf("failed to update pos_posting_audit: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_posting_audit not found or already deleted")
	}

	r.logger.Info("updated pos_posting_audit",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a pos_posting_audit record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_posting_audit", duration, nil)
	}()

	query := `
		UPDATE pos_posting_audit
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete pos_posting_audit", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete pos_posting_audit: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_posting_audit not found or already deleted")
	}

	r.logger.Info("deleted pos_posting_audit", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves pos_posting_audit records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PosPostingAudit, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_posting_audit", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_posting_audit
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_posting_audit records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, source_table
			, source_id
			, source_reference
			, posting_status
			, 'pending',
			, 'processing',
			, 'posted',
			, 'failed',
			, 'cancelled',
			, 'reversed'
			, journal_entry_id
			, reversal_journal_entry_id
			, posting_date
			, posted_at
			, posted_by
			, posting_method
			, error_code
			, error_message
			, error_details
			, retry_count
			, last_retry_at
			, max_retries
			, reversed_at
			, reversed_by
			, reversal_reason
			, total_debit
			, total_credit
			, line_count
			, currency_code
			, posting_context
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (posting_status
			, (ABS(COALESCE(total_debit,
		FROM pos_posting_audit
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_posting_audit by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_posting_audit: %w", err)
	}
	defer rows.Close()

	var entities []*PosPostingAudit
	for rows.Next() {
		var entity PosPostingAudit
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SourceTable,
			&entity.SourceId,
			&entity.SourceReference,
			&entity.PostingStatus,
			&entity.'pending',,
			&entity.'processing',,
			&entity.'posted',,
			&entity.'failed',,
			&entity.'cancelled',,
			&entity.'reversed',
			&entity.JournalEntryId,
			&entity.ReversalJournalEntryId,
			&entity.PostingDate,
			&entity.PostedAt,
			&entity.PostedBy,
			&entity.PostingMethod,
			&entity.ErrorCode,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.RetryCount,
			&entity.LastRetryAt,
			&entity.MaxRetries,
			&entity.ReversedAt,
			&entity.ReversedBy,
			&entity.ReversalReason,
			&entity.TotalDebit,
			&entity.TotalCredit,
			&entity.LineCount,
			&entity.CurrencyCode,
			&entity.PostingContext,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(postingStatus,
			&entity.(ABS(COALESCE(totalDebit,,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_posting_audit: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

