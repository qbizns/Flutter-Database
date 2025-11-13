package bank_statement_line

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

// Repository handles database operations for BankStatementLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankStatementLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankStatementLines represents a bank_statement_lines entity
type BankStatementLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	BankStatementId uuid.UUID `json:"bank_statement_id" db:"bank_statement_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	ValueDate *time.Time `json:"value_date" db:"value_date"`
	Amount float64 `json:"amount" db:"amount"`
	CurrencyCode *string `json:"currency_code" db:"currency_code"`
	Description *string `json:"description" db:"description"`
	Reference *string `json:"reference" db:"reference"`
	CounterpartyName *string `json:"counterparty_name" db:"counterparty_name"`
	CounterpartyAccount *string `json:"counterparty_account" db:"counterparty_account"`
	BankReference *string `json:"bank_reference" db:"bank_reference"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new bank_statement_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankStatementLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_statement_lines", duration, nil)
	}()

	query := `
		INSERT INTO bank_statement_lines (
			, bank_statement_id
			, line_number
			, transaction_date
			, value_date
			, amount
			, currency_code
			, description
			, reference
			, counterparty_name
			, counterparty_account
			, bank_reference
			, status
			, notes
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
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.BankStatementId,
		entity.LineNumber,
		entity.TransactionDate,
		entity.ValueDate,
		entity.Amount,
		entity.CurrencyCode,
		entity.Description,
		entity.Reference,
		entity.CounterpartyName,
		entity.CounterpartyAccount,
		entity.BankReference,
		entity.Status,
		entity.Status,
		entity.Notes,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create bank_statement_lines", zap.Error(err))
		return fmt.Errorf("failed to create bank_statement_lines: %w", err)
	}

	r.logger.Info("created bank_statement_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a bank_statement_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankStatementLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statement_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, bank_statement_id
			, line_number
			, transaction_date
			, value_date
			, amount
			, currency_code
			, description
			, reference
			, counterparty_name
			, counterparty_account
			, bank_reference
			, notes
			, created_at
			, updated_at
			, deleted_at
		FROM bank_statement_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BankStatementLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.BankStatementId,
		&entity.LineNumber,
		&entity.TransactionDate,
		&entity.ValueDate,
		&entity.Amount,
		&entity.CurrencyCode,
		&entity.Description,
		&entity.Reference,
		&entity.CounterpartyName,
		&entity.CounterpartyAccount,
		&entity.BankReference,
		&entity.Status,
		&entity.Status,
		&entity.Notes,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_statement_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_statement_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_statement_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_statement_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankStatementLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_statement_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_statement_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_statement_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, bank_statement_id
			, line_number
			, transaction_date
			, value_date
			, amount
			, currency_code
			, description
			, reference
			, counterparty_name
			, counterparty_account
			, bank_reference
			, notes
			, created_at
			, updated_at
			, deleted_at
		FROM bank_statement_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_statement_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_statement_lines: %w", err)
	}
	defer rows.Close()

	var entities []*BankStatementLines
	for rows.Next() {
		var entity BankStatementLines
		err := rows.Scan(
			&entity.Id,
			&entity.BankStatementId,
			&entity.LineNumber,
			&entity.TransactionDate,
			&entity.ValueDate,
			&entity.Amount,
			&entity.CurrencyCode,
			&entity.Description,
			&entity.Reference,
			&entity.CounterpartyName,
			&entity.CounterpartyAccount,
			&entity.BankReference,
			&entity.Status,
			&entity.Status,
			&entity.Notes,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_statement_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_statement_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_statement_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankStatementLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statement_lines", duration, nil)
	}()

	query := `
		UPDATE bank_statement_lines
		SET
			, bank_statement_id = $2
			, line_number = $3
			, transaction_date = $4
			, value_date = $5
			, amount = $6
			, currency_code = $7
			, description = $8
			, reference = $9
			, counterparty_name = $10
			, counterparty_account = $11
			, bank_reference = $12
			, status = $13
			, notes = $15
			, updated_at = $17
			, deleted_at = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.BankStatementId,
		entity.LineNumber,
		entity.TransactionDate,
		entity.ValueDate,
		entity.Amount,
		entity.CurrencyCode,
		entity.Description,
		entity.Reference,
		entity.CounterpartyName,
		entity.CounterpartyAccount,
		entity.BankReference,
		entity.Status,
		entity.Status,
		entity.Notes,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_statement_lines", zap.Error(err))
		return fmt.Errorf("failed to update bank_statement_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statement_lines not found or already deleted")
	}

	r.logger.Info("updated bank_statement_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a bank_statement_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_statement_lines", duration, nil)
	}()

	query := `
		UPDATE bank_statement_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_statement_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_statement_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_statement_lines not found or already deleted")
	}

	r.logger.Info("deleted bank_statement_lines", zap.String("id", id.String()))
	return nil
}



