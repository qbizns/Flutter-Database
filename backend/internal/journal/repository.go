package journal

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

// Repository handles database operations for Journals
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Journals repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Journals represents a journals entity
type Journals struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	JournalCode string `json:"journal_code" db:"journal_code"`
	JournalName string `json:"journal_name" db:"journal_name"`
	JournalType string `json:"journal_type" db:"journal_type"`
	'sale', *string `json:"'sale'," db:"'sale',"`
	BankAccountId *uuid.UUID `json:"bank_account_id" db:"bank_account_id"`
	DefaultDebitAccountId *uuid.UUID `json:"default_debit_account_id" db:"default_debit_account_id"`
	DefaultCreditAccountId *uuid.UUID `json:"default_credit_account_id" db:"default_credit_account_id"`
	SequencePrefix *string `json:"sequence_prefix" db:"sequence_prefix"`
	SequenceNumber *int64 `json:"sequence_number" db:"sequence_number"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	(journalType string `json:"(journal_type" db:"(journal_type"`
}

// Create inserts a new journals record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Journals) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "journals", duration, nil)
	}()

	query := `
		INSERT INTO journals (
			, organization_id
			, journal_code
			, journal_name
			, journal_type
			, 'sale',
			, bank_account_id
			, default_debit_account_id
			, default_credit_account_id
			, sequence_prefix
			, sequence_number
			, is_active
			, notes
			, created_by
			, updated_by
			, deleted_at
			, (journal_type
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
			, $18
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.JournalCode,
		entity.JournalName,
		entity.JournalType,
		entity.'sale',,
		entity.BankAccountId,
		entity.DefaultDebitAccountId,
		entity.DefaultCreditAccountId,
		entity.SequencePrefix,
		entity.SequenceNumber,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.(journalType,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create journals", zap.Error(err))
		return fmt.Errorf("failed to create journals: %w", err)
	}

	r.logger.Info("created journals",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a journals by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Journals, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journals", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, journal_code
			, journal_name
			, journal_type
			, 'sale',
			, bank_account_id
			, default_debit_account_id
			, default_credit_account_id
			, sequence_prefix
			, sequence_number
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
			, (journal_type
		FROM journals
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Journals
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.JournalCode,
		&entity.JournalName,
		&entity.JournalType,
		&entity.'sale',,
		&entity.BankAccountId,
		&entity.DefaultDebitAccountId,
		&entity.DefaultCreditAccountId,
		&entity.SequencePrefix,
		&entity.SequenceNumber,
		&entity.IsActive,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.(journalType,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("journals not found")
	}

	if err != nil {
		r.logger.Error("failed to get journals", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get journals: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of journals records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Journals, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journals", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journals
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journals records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_code
			, journal_name
			, journal_type
			, 'sale',
			, bank_account_id
			, default_debit_account_id
			, default_credit_account_id
			, sequence_prefix
			, sequence_number
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
			, (journal_type
		FROM journals
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journals", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journals: %w", err)
	}
	defer rows.Close()

	var entities []*Journals
	for rows.Next() {
		var entity Journals
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalCode,
			&entity.JournalName,
			&entity.JournalType,
			&entity.'sale',,
			&entity.BankAccountId,
			&entity.DefaultDebitAccountId,
			&entity.DefaultCreditAccountId,
			&entity.SequencePrefix,
			&entity.SequenceNumber,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.(journalType,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journals: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating journals rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing journals record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Journals) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journals", duration, nil)
	}()

	query := `
		UPDATE journals
		SET
			, organization_id = $2
			, journal_code = $3
			, journal_name = $4
			, journal_type = $5
			, 'sale', = $6
			, bank_account_id = $7
			, default_debit_account_id = $8
			, default_credit_account_id = $9
			, sequence_prefix = $10
			, sequence_number = $11
			, is_active = $12
			, notes = $13
			, created_by = $14
			, updated_by = $15
			, updated_at = $17
			, deleted_at = $18
			, (journal_type = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.JournalCode,
		entity.JournalName,
		entity.JournalType,
		entity.'sale',,
		entity.BankAccountId,
		entity.DefaultDebitAccountId,
		entity.DefaultCreditAccountId,
		entity.SequencePrefix,
		entity.SequenceNumber,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		time.Now(),
		entity.DeletedAt,
		entity.(journalType,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update journals", zap.Error(err))
		return fmt.Errorf("failed to update journals: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journals not found or already deleted")
	}

	r.logger.Info("updated journals",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a journals record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journals", duration, nil)
	}()

	query := `
		UPDATE journals
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete journals", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete journals: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journals not found or already deleted")
	}

	r.logger.Info("deleted journals", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves journals records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Journals, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journals", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journals
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journals records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, journal_code
			, journal_name
			, journal_type
			, 'sale',
			, bank_account_id
			, default_debit_account_id
			, default_credit_account_id
			, sequence_prefix
			, sequence_number
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
			, (journal_type
		FROM journals
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journals by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journals: %w", err)
	}
	defer rows.Close()

	var entities []*Journals
	for rows.Next() {
		var entity Journals
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JournalCode,
			&entity.JournalName,
			&entity.JournalType,
			&entity.'sale',,
			&entity.BankAccountId,
			&entity.DefaultDebitAccountId,
			&entity.DefaultCreditAccountId,
			&entity.SequencePrefix,
			&entity.SequenceNumber,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.(journalType,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journals: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

