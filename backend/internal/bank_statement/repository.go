package bank_statement

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

// Repository handles database operations for BankStatements
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankStatements repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankStatements represents a bank_statements entity
type BankStatements struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BankAccountId uuid.UUID `json:"bank_account_id" db:"bank_account_id"`
	StatementNumber *string `json:"statement_number" db:"statement_number"`
	StatementDate time.Time `json:"statement_date" db:"statement_date"`
	PeriodStartDate time.Time `json:"period_start_date" db:"period_start_date"`
	PeriodEndDate time.Time `json:"period_end_date" db:"period_end_date"`
	OpeningBalance float64 `json:"opening_balance" db:"opening_balance"`
	ClosingBalance float64 `json:"closing_balance" db:"closing_balance"`
	ImportSource *string `json:"import_source" db:"import_source"`
	ImportSource *string `json:"import_source" db:"import_source"`
	ImportFileName *string `json:"import_file_name" db:"import_file_name"`
	Status *string `json:"status" db:"status"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new bank_statements record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankStatements) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_statements", duration, nil)
	}()

	query := `
		INSERT INTO bank_statements (
			, organization_id
			, bank_account_id
			, statement_number
			, statement_date
			, period_start_date
			, period_end_date
			, opening_balance
			, closing_balance
			, import_source
			, import_source
			, import_file_name
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
			, $14
			, $15
			, $16
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BankAccountId,
		entity.StatementNumber,
		entity.StatementDate,
		entity.PeriodStartDate,
		entity.PeriodEndDate,
		entity.OpeningBalance,
		entity.ClosingBalance,
		entity.ImportSource,
		entity.ImportSource,
		entity.ImportFileName,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create bank_statements", zap.Error(err))
		return fmt.Errorf("failed to create bank_statements: %w", err)
	}

	r.logger.Info("created bank_statements",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a bank_statements by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankStatements, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statements", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_number
			, statement_date
			, period_start_date
			, period_end_date
			, opening_balance
			, closing_balance
			, import_source
			, import_source
			, import_file_name
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM bank_statements
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BankStatements
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BankAccountId,
		&entity.StatementNumber,
		&entity.StatementDate,
		&entity.PeriodStartDate,
		&entity.PeriodEndDate,
		&entity.OpeningBalance,
		&entity.ClosingBalance,
		&entity.ImportSource,
		&entity.ImportSource,
		&entity.ImportFileName,
		&entity.Status,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_statements not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_statements", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_statements: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_statements records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankStatements, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statements", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_statements
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_statements records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_number
			, statement_date
			, period_start_date
			, period_end_date
			, opening_balance
			, closing_balance
			, import_source
			, import_source
			, import_file_name
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM bank_statements
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_statements", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_statements: %w", err)
	}
	defer rows.Close()

	var entities []*BankStatements
	for rows.Next() {
		var entity BankStatements
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankAccountId,
			&entity.StatementNumber,
			&entity.StatementDate,
			&entity.PeriodStartDate,
			&entity.PeriodEndDate,
			&entity.OpeningBalance,
			&entity.ClosingBalance,
			&entity.ImportSource,
			&entity.ImportSource,
			&entity.ImportFileName,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_statements: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_statements rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_statements record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankStatements) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statements", duration, nil)
	}()

	query := `
		UPDATE bank_statements
		SET
			, organization_id = $2
			, bank_account_id = $3
			, statement_number = $4
			, statement_date = $5
			, period_start_date = $6
			, period_end_date = $7
			, opening_balance = $8
			, closing_balance = $9
			, import_source = $10
			, import_source = $11
			, import_file_name = $12
			, status = $13
			, notes = $14
			, created_by = $15
			, updated_by = $16
			, updated_at = $18
			, deleted_at = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BankAccountId,
		entity.StatementNumber,
		entity.StatementDate,
		entity.PeriodStartDate,
		entity.PeriodEndDate,
		entity.OpeningBalance,
		entity.ClosingBalance,
		entity.ImportSource,
		entity.ImportSource,
		entity.ImportFileName,
		entity.Status,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_statements", zap.Error(err))
		return fmt.Errorf("failed to update bank_statements: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statements not found or already deleted")
	}

	r.logger.Info("updated bank_statements",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a bank_statements record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statements", duration, nil)
	}()

	query := `
		UPDATE bank_statements
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_statements", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_statements: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statements not found or already deleted")
	}

	r.logger.Info("deleted bank_statements", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves bank_statements records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BankStatements, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statements", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_statements
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_statements records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bank_account_id
			, statement_number
			, statement_date
			, period_start_date
			, period_end_date
			, opening_balance
			, closing_balance
			, import_source
			, import_source
			, import_file_name
			, status
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM bank_statements
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_statements by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_statements: %w", err)
	}
	defer rows.Close()

	var entities []*BankStatements
	for rows.Next() {
		var entity BankStatements
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BankAccountId,
			&entity.StatementNumber,
			&entity.StatementDate,
			&entity.PeriodStartDate,
			&entity.PeriodEndDate,
			&entity.OpeningBalance,
			&entity.ClosingBalance,
			&entity.ImportSource,
			&entity.ImportSource,
			&entity.ImportFileName,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_statements: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

