package chart_of_account

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

// Repository handles database operations for ChartOfAccounts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ChartOfAccounts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ChartOfAccounts represents a chart_of_accounts entity
type ChartOfAccounts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	AccountCode string `json:"account_code" db:"account_code"`
	AccountNumber string `json:"account_number" db:"account_number"`
	AccountName string `json:"account_name" db:"account_name"`
	AccountTypeId uuid.UUID `json:"account_type_id" db:"account_type_id"`
	AccountSubtypeId *uuid.UUID `json:"account_subtype_id" db:"account_subtype_id"`
	ParentAccountId *uuid.UUID `json:"parent_account_id" db:"parent_account_id"`
	AccountLevel *int64 `json:"account_level" db:"account_level"`
	AccountPath *string `json:"account_path" db:"account_path"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsSystemAccount *bool `json:"is_system_account" db:"is_system_account"`
	IsHeaderAccount *bool `json:"is_header_account" db:"is_header_account"`
	IsBankAccount *bool `json:"is_bank_account" db:"is_bank_account"`
	IsReconcilable *bool `json:"is_reconcilable" db:"is_reconcilable"`
	DefaultTaxCode *string `json:"default_tax_code" db:"default_tax_code"`
	CurrencyCode *string `json:"currency_code" db:"currency_code"`
	OpeningBalance *float64 `json:"opening_balance" db:"opening_balance"`
	OpeningBalanceDate *time.Time `json:"opening_balance_date" db:"opening_balance_date"`
	CurrentDebitBalance *float64 `json:"current_debit_balance" db:"current_debit_balance"`
	CurrentCreditBalance *float64 `json:"current_credit_balance" db:"current_credit_balance"`
	CurrentBalance *float64 `json:"current_balance" db:"current_balance"`
	LastBalanceUpdate *time.Time `json:"last_balance_update" db:"last_balance_update"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new chart_of_accounts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ChartOfAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "chart_of_accounts", duration, nil)
	}()

	query := `
		INSERT INTO chart_of_accounts (
			, organization_id
			, account_code
			, account_number
			, account_name
			, account_type_id
			, account_subtype_id
			, parent_account_id
			, account_level
			, account_path
			, is_active
			, is_system_account
			, is_header_account
			, is_bank_account
			, is_reconcilable
			, default_tax_code
			, currency_code
			, opening_balance
			, opening_balance_date
			, current_debit_balance
			, current_credit_balance
			, current_balance
			, last_balance_update
			, description
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
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $29
			, $30
			, $31
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.AccountCode,
		entity.AccountNumber,
		entity.AccountName,
		entity.AccountTypeId,
		entity.AccountSubtypeId,
		entity.ParentAccountId,
		entity.AccountLevel,
		entity.AccountPath,
		entity.IsActive,
		entity.IsSystemAccount,
		entity.IsHeaderAccount,
		entity.IsBankAccount,
		entity.IsReconcilable,
		entity.DefaultTaxCode,
		entity.CurrencyCode,
		entity.OpeningBalance,
		entity.OpeningBalanceDate,
		entity.CurrentDebitBalance,
		entity.CurrentCreditBalance,
		entity.CurrentBalance,
		entity.LastBalanceUpdate,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create chart_of_accounts", zap.Error(err))
		return fmt.Errorf("failed to create chart_of_accounts: %w", err)
	}

	r.logger.Info("created chart_of_accounts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a chart_of_accounts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ChartOfAccounts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "chart_of_accounts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, account_code
			, account_number
			, account_name
			, account_type_id
			, account_subtype_id
			, parent_account_id
			, account_level
			, account_path
			, is_active
			, is_system_account
			, is_header_account
			, is_bank_account
			, is_reconcilable
			, default_tax_code
			, currency_code
			, opening_balance
			, opening_balance_date
			, current_debit_balance
			, current_credit_balance
			, current_balance
			, last_balance_update
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM chart_of_accounts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ChartOfAccounts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.AccountCode,
		&entity.AccountNumber,
		&entity.AccountName,
		&entity.AccountTypeId,
		&entity.AccountSubtypeId,
		&entity.ParentAccountId,
		&entity.AccountLevel,
		&entity.AccountPath,
		&entity.IsActive,
		&entity.IsSystemAccount,
		&entity.IsHeaderAccount,
		&entity.IsBankAccount,
		&entity.IsReconcilable,
		&entity.DefaultTaxCode,
		&entity.CurrencyCode,
		&entity.OpeningBalance,
		&entity.OpeningBalanceDate,
		&entity.CurrentDebitBalance,
		&entity.CurrentCreditBalance,
		&entity.CurrentBalance,
		&entity.LastBalanceUpdate,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("chart_of_accounts not found")
	}

	if err != nil {
		r.logger.Error("failed to get chart_of_accounts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get chart_of_accounts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of chart_of_accounts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ChartOfAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "chart_of_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM chart_of_accounts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chart_of_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, account_code
			, account_number
			, account_name
			, account_type_id
			, account_subtype_id
			, parent_account_id
			, account_level
			, account_path
			, is_active
			, is_system_account
			, is_header_account
			, is_bank_account
			, is_reconcilable
			, default_tax_code
			, currency_code
			, opening_balance
			, opening_balance_date
			, current_debit_balance
			, current_credit_balance
			, current_balance
			, last_balance_update
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM chart_of_accounts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list chart_of_accounts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list chart_of_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*ChartOfAccounts
	for rows.Next() {
		var entity ChartOfAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AccountCode,
			&entity.AccountNumber,
			&entity.AccountName,
			&entity.AccountTypeId,
			&entity.AccountSubtypeId,
			&entity.ParentAccountId,
			&entity.AccountLevel,
			&entity.AccountPath,
			&entity.IsActive,
			&entity.IsSystemAccount,
			&entity.IsHeaderAccount,
			&entity.IsBankAccount,
			&entity.IsReconcilable,
			&entity.DefaultTaxCode,
			&entity.CurrencyCode,
			&entity.OpeningBalance,
			&entity.OpeningBalanceDate,
			&entity.CurrentDebitBalance,
			&entity.CurrentCreditBalance,
			&entity.CurrentBalance,
			&entity.LastBalanceUpdate,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan chart_of_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating chart_of_accounts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing chart_of_accounts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ChartOfAccounts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "chart_of_accounts", duration, nil)
	}()

	query := `
		UPDATE chart_of_accounts
		SET
			, organization_id = $2
			, account_code = $3
			, account_number = $4
			, account_name = $5
			, account_type_id = $6
			, account_subtype_id = $7
			, parent_account_id = $8
			, account_level = $9
			, account_path = $10
			, is_active = $11
			, is_system_account = $12
			, is_header_account = $13
			, is_bank_account = $14
			, is_reconcilable = $15
			, default_tax_code = $16
			, currency_code = $17
			, opening_balance = $18
			, opening_balance_date = $19
			, current_debit_balance = $20
			, current_credit_balance = $21
			, current_balance = $22
			, last_balance_update = $23
			, description = $24
			, notes = $25
			, metadata = $26
			, updated_at = $28
			, created_by = $29
			, updated_by = $30
			, deleted_at = $31
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $32
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.AccountCode,
		entity.AccountNumber,
		entity.AccountName,
		entity.AccountTypeId,
		entity.AccountSubtypeId,
		entity.ParentAccountId,
		entity.AccountLevel,
		entity.AccountPath,
		entity.IsActive,
		entity.IsSystemAccount,
		entity.IsHeaderAccount,
		entity.IsBankAccount,
		entity.IsReconcilable,
		entity.DefaultTaxCode,
		entity.CurrencyCode,
		entity.OpeningBalance,
		entity.OpeningBalanceDate,
		entity.CurrentDebitBalance,
		entity.CurrentCreditBalance,
		entity.CurrentBalance,
		entity.LastBalanceUpdate,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update chart_of_accounts", zap.Error(err))
		return fmt.Errorf("failed to update chart_of_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("chart_of_accounts not found or already deleted")
	}

	r.logger.Info("updated chart_of_accounts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a chart_of_accounts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "chart_of_accounts", duration, nil)
	}()

	query := `
		UPDATE chart_of_accounts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete chart_of_accounts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete chart_of_accounts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("chart_of_accounts not found or already deleted")
	}

	r.logger.Info("deleted chart_of_accounts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves chart_of_accounts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ChartOfAccounts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "chart_of_accounts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM chart_of_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count chart_of_accounts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, account_code
			, account_number
			, account_name
			, account_type_id
			, account_subtype_id
			, parent_account_id
			, account_level
			, account_path
			, is_active
			, is_system_account
			, is_header_account
			, is_bank_account
			, is_reconcilable
			, default_tax_code
			, currency_code
			, opening_balance
			, opening_balance_date
			, current_debit_balance
			, current_credit_balance
			, current_balance
			, last_balance_update
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM chart_of_accounts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list chart_of_accounts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list chart_of_accounts: %w", err)
	}
	defer rows.Close()

	var entities []*ChartOfAccounts
	for rows.Next() {
		var entity ChartOfAccounts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AccountCode,
			&entity.AccountNumber,
			&entity.AccountName,
			&entity.AccountTypeId,
			&entity.AccountSubtypeId,
			&entity.ParentAccountId,
			&entity.AccountLevel,
			&entity.AccountPath,
			&entity.IsActive,
			&entity.IsSystemAccount,
			&entity.IsHeaderAccount,
			&entity.IsBankAccount,
			&entity.IsReconcilable,
			&entity.DefaultTaxCode,
			&entity.CurrencyCode,
			&entity.OpeningBalance,
			&entity.OpeningBalanceDate,
			&entity.CurrentDebitBalance,
			&entity.CurrentCreditBalance,
			&entity.CurrentBalance,
			&entity.LastBalanceUpdate,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan chart_of_accounts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

