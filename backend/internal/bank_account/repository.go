package bank_account

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

// Repository handles database operations for BankAccounts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BankAccounts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BankAccounts represents a bank_accounts entity
type BankAccounts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ChartAccountId uuid.UUID `json:"chart_account_id" db:"chart_account_id"`
	BankName string `json:"bank_name" db:"bank_name"`
	AccountNumber string `json:"account_number" db:"account_number"`
	AccountType *string `json:"account_type" db:"account_type"`
	RoutingNumber *string `json:"routing_number" db:"routing_number"`
	SwiftCode *string `json:"swift_code" db:"swift_code"`
	CurrencyCode *string `json:"currency_code" db:"currency_code"`
	CurrentBalance *float64 `json:"current_balance" db:"current_balance"`
	StatementBalance *float64 `json:"statement_balance" db:"statement_balance"`
	LastStatementDate *time.Time `json:"last_statement_date" db:"last_statement_date"`
	IsActive *bool `json:"is_active" db:"is_active"`
	OnlineBankingEnabled *bool `json:"online_banking_enabled" db:"online_banking_enabled"`
	LastSyncDate *time.Time `json:"last_sync_date" db:"last_sync_date"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new bank_accounts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BankAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "bank_accounts", duration, nil)
	}()

	query := `
		INSERT INTO bank_accounts (
			, organization_id
			, chart_account_id
			, bank_name
			, account_number
			, account_type
			, routing_number
			, swift_code
			, currency_code
			, current_balance
			, statement_balance
			, last_statement_date
			, is_active
			, online_banking_enabled
			, last_sync_date
			, notes
			, metadata
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
			, $17
			, $20
			, $21
			, $22
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ChartAccountId,
		entity.BankName,
		entity.AccountNumber,
		entity.AccountType,
		entity.RoutingNumber,
		entity.SwiftCode,
		entity.CurrencyCode,
		entity.CurrentBalance,
		entity.StatementBalance,
		entity.LastStatementDate,
		entity.IsActive,
		entity.OnlineBankingEnabled,
		entity.LastSyncDate,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create bank_accounts", zap.Error(err))
		return fmt.Errorf("failed to create bank_accounts: %w", err)
	}

	r.logger.Info("created bank_accounts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a bank_accounts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BankAccounts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_accounts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, chart_account_id
			, bank_name
			, account_number
			, account_type
			, routing_number
			, swift_code
			, currency_code
			, current_balance
			, statement_balance
			, last_statement_date
			, is_active
			, online_banking_enabled
			, last_sync_date
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_accounts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity BankAccounts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ChartAccountId,
		&entity.BankName,
		&entity.AccountNumber,
		&entity.AccountType,
		&entity.RoutingNumber,
		&entity.SwiftCode,
		&entity.CurrencyCode,
		&entity.CurrentBalance,
		&entity.StatementBalance,
		&entity.LastStatementDate,
		&entity.IsActive,
		&entity.OnlineBankingEnabled,
		&entity.LastSyncDate,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("bank_accounts not found")
	}

	if err != nil {
		r.logger.Error("failed to get bank_accounts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get bank_accounts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of bank_accounts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BankAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_accounts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, chart_account_id
			, bank_name
			, account_number
			, account_type
			, routing_number
			, swift_code
			, currency_code
			, current_balance
			, statement_balance
			, last_statement_date
			, is_active
			, online_banking_enabled
			, last_sync_date
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*BankAccounts
	for rows.Next() {
		var entity BankAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ChartAccountId,
			&entity.BankName,
			&entity.AccountNumber,
			&entity.AccountType,
			&entity.RoutingNumber,
			&entity.SwiftCode,
			&entity.CurrencyCode,
			&entity.CurrentBalance,
			&entity.StatementBalance,
			&entity.LastStatementDate,
			&entity.IsActive,
			&entity.OnlineBankingEnabled,
			&entity.LastSyncDate,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating bank_accounts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing bank_accounts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BankAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_accounts", duration, nil)
	}()

	query := `
		UPDATE bank_accounts
		SET
			, organization_id = $2
			, chart_account_id = $3
			, bank_name = $4
			, account_number = $5
			, account_type = $6
			, routing_number = $7
			, swift_code = $8
			, currency_code = $9
			, current_balance = $10
			, statement_balance = $11
			, last_statement_date = $12
			, is_active = $13
			, online_banking_enabled = $14
			, last_sync_date = $15
			, notes = $16
			, metadata = $17
			, updated_at = $19
			, created_by = $20
			, updated_by = $21
			, deleted_at = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ChartAccountId,
		entity.BankName,
		entity.AccountNumber,
		entity.AccountType,
		entity.RoutingNumber,
		entity.SwiftCode,
		entity.CurrencyCode,
		entity.CurrentBalance,
		entity.StatementBalance,
		entity.LastStatementDate,
		entity.IsActive,
		entity.OnlineBankingEnabled,
		entity.LastSyncDate,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update bank_accounts", zap.Error(err))
		return fmt.Errorf("failed to update bank_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_accounts not found or already deleted")
	}

	r.logger.Info("updated bank_accounts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a bank_accounts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "bank_accounts", duration, nil)
	}()

	query := `
		UPDATE bank_accounts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete bank_accounts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete bank_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("bank_accounts not found or already deleted")
	}

	r.logger.Info("deleted bank_accounts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves bank_accounts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BankAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "bank_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM bank_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count bank_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, chart_account_id
			, bank_name
			, account_number
			, account_type
			, routing_number
			, swift_code
			, currency_code
			, current_balance
			, statement_balance
			, last_statement_date
			, is_active
			, online_banking_enabled
			, last_sync_date
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM bank_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list bank_accounts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list bank_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*BankAccounts
	for rows.Next() {
		var entity BankAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ChartAccountId,
			&entity.BankName,
			&entity.AccountNumber,
			&entity.AccountType,
			&entity.RoutingNumber,
			&entity.SwiftCode,
			&entity.CurrencyCode,
			&entity.CurrentBalance,
			&entity.StatementBalance,
			&entity.LastStatementDate,
			&entity.IsActive,
			&entity.OnlineBankingEnabled,
			&entity.LastSyncDate,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan bank_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

