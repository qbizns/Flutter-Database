package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"backend/internal/domain/accounting"
)

// AccountingRepository implements accounting.AccountingRepository
type AccountingRepository struct {
	db *sqlx.DB
}

// NewAccountingRepository creates a new accounting repository
func NewAccountingRepository(db *sqlx.DB) accounting.AccountingRepository {
	return &AccountingRepository{db: db}
}

// ========================
// FISCAL YEARS
// ========================

// CreateFiscalYear creates a new fiscal year
func (r *AccountingRepository) CreateFiscalYear(ctx context.Context, fy *accounting.FiscalYear) error {
	query := `
		INSERT INTO fiscal_years (
			id, organization_id, fiscal_year, start_date, end_date,
			status, is_current, notes, metadata, created_at, updated_at,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		fy.ID, fy.OrganizationID, fy.FiscalYear, fy.StartDate, fy.EndDate,
		fy.Status, fy.IsCurrent, fy.Notes, fy.Metadata, fy.CreatedAt,
		fy.UpdatedAt, fy.CreatedBy, fy.UpdatedBy,
	)

	return err
}

// GetFiscalYear retrieves a fiscal year by ID
func (r *AccountingRepository) GetFiscalYear(ctx context.Context, id, organizationID uuid.UUID) (*accounting.FiscalYear, error) {
	query := `
		SELECT id, organization_id, fiscal_year, start_date, end_date,
		       status, is_current, closed_by, closed_at, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM fiscal_years
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	fy := &accounting.FiscalYear{}
	err := r.db.GetContext(ctx, fy, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("fiscal year not found")
		}
		return nil, err
	}

	return fy, nil
}

// ListFiscalYears lists fiscal years with filters
func (r *AccountingRepository) ListFiscalYears(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*accounting.FiscalYear, error) {
	query := `
		SELECT id, organization_id, fiscal_year, start_date, end_date,
		       status, is_current, closed_by, closed_at, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM fiscal_years
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if status, ok := filter["status"]; ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if isCurrent, ok := filter["is_current"].(bool); ok {
		query += fmt.Sprintf(" AND is_current = $%d", argIndex)
		args = append(args, isCurrent)
		argIndex++
	}

	query += " ORDER BY start_date DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fiscalYears []*accounting.FiscalYear
	for rows.Next() {
		fy := &accounting.FiscalYear{}
		if err := rows.StructScan(fy); err != nil {
			return nil, err
		}
		fiscalYears = append(fiscalYears, fy)
	}

	return fiscalYears, rows.Err()
}

// UpdateFiscalYear updates a fiscal year
func (r *AccountingRepository) UpdateFiscalYear(ctx context.Context, fy *accounting.FiscalYear) error {
	query := `
		UPDATE fiscal_years
		SET fiscal_year = $1, start_date = $2, end_date = $3,
		    status = $4, is_current = $5, notes = $6,
		    metadata = $7, updated_at = $8, updated_by = $9
		WHERE id = $10 AND organization_id = $11 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		fy.FiscalYear, fy.StartDate, fy.EndDate,
		fy.Status, fy.IsCurrent, fy.Notes,
		fy.Metadata, fy.UpdatedAt, fy.UpdatedBy,
		fy.ID, fy.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("fiscal year not found")
	}

	return nil
}

// DeleteFiscalYear soft deletes a fiscal year
func (r *AccountingRepository) DeleteFiscalYear(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE fiscal_years
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("fiscal year not found")
	}

	return nil
}

// GetCurrentFiscalYear retrieves the current fiscal year
func (r *AccountingRepository) GetCurrentFiscalYear(ctx context.Context, organizationID uuid.UUID) (*accounting.FiscalYear, error) {
	query := `
		SELECT id, organization_id, fiscal_year, start_date, end_date,
		       status, is_current, closed_by, closed_at, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM fiscal_years
		WHERE organization_id = $1 AND is_current = true AND deleted_at IS NULL
		LIMIT 1
	`

	fy := &accounting.FiscalYear{}
	err := r.db.GetContext(ctx, fy, query, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("current fiscal year not found")
		}
		return nil, err
	}

	return fy, nil
}

// CloseFiscalYear closes a fiscal year
func (r *AccountingRepository) CloseFiscalYear(ctx context.Context, id, organizationID, closedBy uuid.UUID) error {
	query := `
		UPDATE fiscal_years
		SET status = $1, closed_by = $2, closed_at = $3
		WHERE id = $4 AND organization_id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		accounting.FiscalYearStatusClosed, closedBy, time.Now(), id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("fiscal year not found")
	}

	return nil
}

// ========================
// ACCOUNTING PERIODS
// ========================

// CreateAccountingPeriod creates an accounting period
func (r *AccountingRepository) CreateAccountingPeriod(ctx context.Context, ap *accounting.AccountingPeriod) error {
	query := `
		INSERT INTO accounting_periods (
			id, organization_id, fiscal_year_id, period_number, period_name,
			start_date, end_date, status, metadata, created_at, updated_at,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		ap.ID, ap.OrganizationID, ap.FiscalYearID, ap.PeriodNumber, ap.PeriodName,
		ap.StartDate, ap.EndDate, ap.Status, ap.Metadata, ap.CreatedAt,
		ap.UpdatedAt, ap.CreatedBy, ap.UpdatedBy,
	)

	return err
}

// GetAccountingPeriod retrieves an accounting period
func (r *AccountingRepository) GetAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) (*accounting.AccountingPeriod, error) {
	query := `
		SELECT id, organization_id, fiscal_year_id, period_number, period_name,
		       start_date, end_date, status, closed_by, closed_at, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM accounting_periods
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	ap := &accounting.AccountingPeriod{}
	err := r.db.GetContext(ctx, ap, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("accounting period not found")
		}
		return nil, err
	}

	return ap, nil
}

// ListAccountingPeriods lists periods for a fiscal year
func (r *AccountingRepository) ListAccountingPeriods(ctx context.Context, organizationID, fiscalYearID uuid.UUID) ([]*accounting.AccountingPeriod, error) {
	query := `
		SELECT id, organization_id, fiscal_year_id, period_number, period_name,
		       start_date, end_date, status, closed_by, closed_at, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM accounting_periods
		WHERE organization_id = $1 AND fiscal_year_id = $2 AND deleted_at IS NULL
		ORDER BY period_number ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, organizationID, fiscalYearID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var periods []*accounting.AccountingPeriod
	for rows.Next() {
		ap := &accounting.AccountingPeriod{}
		if err := rows.StructScan(ap); err != nil {
			return nil, err
		}
		periods = append(periods, ap)
	}

	return periods, rows.Err()
}

// UpdateAccountingPeriod updates an accounting period
func (r *AccountingRepository) UpdateAccountingPeriod(ctx context.Context, ap *accounting.AccountingPeriod) error {
	query := `
		UPDATE accounting_periods
		SET period_name = $1, status = $2, metadata = $3,
		    updated_at = $4, updated_by = $5
		WHERE id = $6 AND organization_id = $7 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		ap.PeriodName, ap.Status, ap.Metadata,
		ap.UpdatedAt, ap.UpdatedBy,
		ap.ID, ap.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("accounting period not found")
	}

	return nil
}

// DeleteAccountingPeriod soft deletes an accounting period
func (r *AccountingRepository) DeleteAccountingPeriod(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE accounting_periods
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("accounting period not found")
	}

	return nil
}

// CloseAccountingPeriod closes an accounting period
func (r *AccountingRepository) CloseAccountingPeriod(ctx context.Context, id, organizationID, closedBy uuid.UUID) error {
	query := `
		UPDATE accounting_periods
		SET status = $1, closed_by = $2, closed_at = $3
		WHERE id = $4 AND organization_id = $5 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		accounting.AccountingPeriodStatusClosed, closedBy, time.Now(), id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("accounting period not found")
	}

	return nil
}

// GetAccountingPeriodByDate gets the period for a given date
func (r *AccountingRepository) GetAccountingPeriodByDate(ctx context.Context, organizationID uuid.UUID, date time.Time) (*accounting.AccountingPeriod, error) {
	query := `
		SELECT id, organization_id, fiscal_year_id, period_number, period_name,
		       start_date, end_date, status, closed_by, closed_at, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM accounting_periods
		WHERE organization_id = $1
		  AND start_date <= $2 AND end_date >= $2
		  AND deleted_at IS NULL
		LIMIT 1
	`

	ap := &accounting.AccountingPeriod{}
	err := r.db.GetContext(ctx, ap, query, organizationID, date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no period found for date")
		}
		return nil, err
	}

	return ap, nil
}

// ========================
// ACCOUNT TYPES & SUBTYPES
// ========================

// GetAccountType retrieves an account type
func (r *AccountingRepository) GetAccountType(ctx context.Context, id uuid.UUID) (*accounting.AccountType, error) {
	query := `
		SELECT id, type_code, type_name, type_category, normal_balance,
		       is_balance_sheet, is_income_statement, display_order, description,
		       created_at, updated_at
		FROM account_types
		WHERE id = $1
	`

	at := &accounting.AccountType{}
	err := r.db.GetContext(ctx, at, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account type not found")
		}
		return nil, err
	}

	return at, nil
}

// ListAccountTypes lists account types
func (r *AccountingRepository) ListAccountTypes(ctx context.Context) ([]*accounting.AccountType, error) {
	query := `
		SELECT id, type_code, type_name, type_category, normal_balance,
		       is_balance_sheet, is_income_statement, display_order, description,
		       created_at, updated_at
		FROM account_types
		ORDER BY display_order ASC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*accounting.AccountType
	for rows.Next() {
		at := &accounting.AccountType{}
		if err := rows.StructScan(at); err != nil {
			return nil, err
		}
		types = append(types, at)
	}

	return types, rows.Err()
}

// GetAccountSubtype retrieves an account subtype
func (r *AccountingRepository) GetAccountSubtype(ctx context.Context, id uuid.UUID) (*accounting.AccountSubtype, error) {
	query := `
		SELECT id, account_type_id, subtype_code, subtype_name,
		       display_order, description, created_at, updated_at
		FROM account_subtypes
		WHERE id = $1
	`

	as := &accounting.AccountSubtype{}
	err := r.db.GetContext(ctx, as, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account subtype not found")
		}
		return nil, err
	}

	return as, nil
}

// ListAccountSubtypes lists subtypes for a type
func (r *AccountingRepository) ListAccountSubtypes(ctx context.Context, accountTypeID uuid.UUID) ([]*accounting.AccountSubtype, error) {
	query := `
		SELECT id, account_type_id, subtype_code, subtype_name,
		       display_order, description, created_at, updated_at
		FROM account_subtypes
		WHERE account_type_id = $1
		ORDER BY display_order ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, accountTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtypes []*accounting.AccountSubtype
	for rows.Next() {
		as := &accounting.AccountSubtype{}
		if err := rows.StructScan(as); err != nil {
			return nil, err
		}
		subtypes = append(subtypes, as)
	}

	return subtypes, rows.Err()
}

// ========================
// CHART OF ACCOUNTS
// ========================

// CreateChartOfAccount creates a chart of account
func (r *AccountingRepository) CreateChartOfAccount(ctx context.Context, coa *accounting.ChartOfAccount) error {
	query := `
		INSERT INTO chart_of_accounts (
			id, organization_id, account_code, account_number, account_name,
			account_type_id, account_subtype_id, parent_account_id, account_level,
			account_path, is_active, is_system_account, is_header_account,
			is_bank_account, is_reconcilable, default_tax_code, currency_code,
			opening_balance, opening_balance_date, current_debit_balance,
			current_credit_balance, current_balance, last_balance_update,
			description, notes, metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		coa.ID, coa.OrganizationID, coa.AccountCode, coa.AccountNumber, coa.AccountName,
		coa.AccountTypeID, coa.AccountSubtypeID, coa.ParentAccountID, coa.AccountLevel,
		coa.AccountPath, coa.IsActive, coa.IsSystemAccount, coa.IsHeaderAccount,
		coa.IsBankAccount, coa.IsReconcilable, coa.DefaultTaxCode, coa.CurrencyCode,
		coa.OpeningBalance, coa.OpeningBalanceDate, coa.CurrentDebitBalance,
		coa.CurrentCreditBalance, coa.CurrentBalance, coa.LastBalanceUpdate,
		coa.Description, coa.Notes, coa.Metadata, coa.CreatedAt, coa.UpdatedAt,
		coa.CreatedBy, coa.UpdatedBy,
	)

	return err
}

// GetChartOfAccount retrieves a chart of account
func (r *AccountingRepository) GetChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) (*accounting.ChartOfAccount, error) {
	query := `
		SELECT id, organization_id, account_code, account_number, account_name,
		       account_type_id, account_subtype_id, parent_account_id, account_level,
		       account_path, is_active, is_system_account, is_header_account,
		       is_bank_account, is_reconcilable, default_tax_code, currency_code,
		       opening_balance, opening_balance_date, current_debit_balance,
		       current_credit_balance, current_balance, last_balance_update,
		       description, notes, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM chart_of_accounts
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	coa := &accounting.ChartOfAccount{}
	err := r.db.GetContext(ctx, coa, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chart of account not found")
		}
		return nil, err
	}

	return coa, nil
}

// GetChartOfAccountByCode retrieves by code
func (r *AccountingRepository) GetChartOfAccountByCode(ctx context.Context, organizationID uuid.UUID, code string) (*accounting.ChartOfAccount, error) {
	query := `
		SELECT id, organization_id, account_code, account_number, account_name,
		       account_type_id, account_subtype_id, parent_account_id, account_level,
		       account_path, is_active, is_system_account, is_header_account,
		       is_bank_account, is_reconcilable, default_tax_code, currency_code,
		       opening_balance, opening_balance_date, current_debit_balance,
		       current_credit_balance, current_balance, last_balance_update,
		       description, notes, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM chart_of_accounts
		WHERE organization_id = $1 AND account_code = $2 AND deleted_at IS NULL
	`

	coa := &accounting.ChartOfAccount{}
	err := r.db.GetContext(ctx, coa, query, organizationID, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chart of account not found")
		}
		return nil, err
	}

	return coa, nil
}

// ListChartOfAccounts lists accounts
func (r *AccountingRepository) ListChartOfAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*accounting.ChartOfAccount, error) {
	query := `
		SELECT id, organization_id, account_code, account_number, account_name,
		       account_type_id, account_subtype_id, parent_account_id, account_level,
		       account_path, is_active, is_system_account, is_header_account,
		       is_bank_account, is_reconcilable, default_tax_code, currency_code,
		       opening_balance, opening_balance_date, current_debit_balance,
		       current_credit_balance, current_balance, last_balance_update,
		       description, notes, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM chart_of_accounts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if typeID, ok := filter["account_type_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND account_type_id = $%d", argIndex)
		args = append(args, typeID)
		argIndex++
	}

	if parentID, ok := filter["parent_account_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND parent_account_id = $%d", argIndex)
		args = append(args, parentID)
		argIndex++
	}

	query += " ORDER BY account_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*accounting.ChartOfAccount
	for rows.Next() {
		coa := &accounting.ChartOfAccount{}
		if err := rows.StructScan(coa); err != nil {
			return nil, err
		}
		accounts = append(accounts, coa)
	}

	return accounts, rows.Err()
}

// UpdateChartOfAccount updates a chart of account
func (r *AccountingRepository) UpdateChartOfAccount(ctx context.Context, coa *accounting.ChartOfAccount) error {
	query := `
		UPDATE chart_of_accounts
		SET account_name = $1, account_subtype_id = $2, is_active = $3,
		    is_header_account = $4, is_bank_account = $5, is_reconcilable = $6,
		    default_tax_code = $7, description = $8, notes = $9,
		    metadata = $10, updated_at = $11, updated_by = $12
		WHERE id = $13 AND organization_id = $14 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		coa.AccountName, coa.AccountSubtypeID, coa.IsActive,
		coa.IsHeaderAccount, coa.IsBankAccount, coa.IsReconcilable,
		coa.DefaultTaxCode, coa.Description, coa.Notes,
		coa.Metadata, coa.UpdatedAt, coa.UpdatedBy,
		coa.ID, coa.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("chart of account not found")
	}

	return nil
}

// DeleteChartOfAccount soft deletes a chart of account
func (r *AccountingRepository) DeleteChartOfAccount(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE chart_of_accounts
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND is_system_account = false AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("chart of account not found or is system account")
	}

	return nil
}

// GetAccountHierarchy retrieves account hierarchy
func (r *AccountingRepository) GetAccountHierarchy(ctx context.Context, organizationID, parentID uuid.UUID) ([]*accounting.ChartOfAccount, error) {
	query := `
		SELECT id, organization_id, account_code, account_number, account_name,
		       account_type_id, account_subtype_id, parent_account_id, account_level,
		       account_path, is_active, is_system_account, is_header_account,
		       is_bank_account, is_reconcilable, default_tax_code, currency_code,
		       opening_balance, opening_balance_date, current_debit_balance,
		       current_credit_balance, current_balance, last_balance_update,
		       description, notes, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM chart_of_accounts
		WHERE organization_id = $1 AND parent_account_id = $2 AND deleted_at IS NULL
		ORDER BY account_code ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, organizationID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*accounting.ChartOfAccount
	for rows.Next() {
		coa := &accounting.ChartOfAccount{}
		if err := rows.StructScan(coa); err != nil {
			return nil, err
		}
		accounts = append(accounts, coa)
	}

	return accounts, rows.Err()
}

// ========================
// JOURNAL ENTRY TYPES
// ========================

// GetJournalEntryType retrieves a journal entry type
func (r *AccountingRepository) GetJournalEntryType(ctx context.Context, id uuid.UUID) (*accounting.JournalEntryType, error) {
	query := `
		SELECT id, type_code, type_name, type_category, number_prefix,
		       description, created_at, updated_at
		FROM journal_entry_types
		WHERE id = $1
	`

	jet := &accounting.JournalEntryType{}
	err := r.db.GetContext(ctx, jet, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("journal entry type not found")
		}
		return nil, err
	}

	return jet, nil
}

// ListJournalEntryTypes lists journal entry types
func (r *AccountingRepository) ListJournalEntryTypes(ctx context.Context) ([]*accounting.JournalEntryType, error) {
	query := `
		SELECT id, type_code, type_name, type_category, number_prefix,
		       description, created_at, updated_at
		FROM journal_entry_types
		ORDER BY type_name ASC
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []*accounting.JournalEntryType
	for rows.Next() {
		jet := &accounting.JournalEntryType{}
		if err := rows.StructScan(jet); err != nil {
			return nil, err
		}
		types = append(types, jet)
	}

	return types, rows.Err()
}

// ========================
// JOURNAL ENTRIES
// ========================

// CreateJournalEntry creates a journal entry with lines
func (r *AccountingRepository) CreateJournalEntry(ctx context.Context, je *accounting.JournalEntry, lines []*accounting.JournalEntryLine) error {
	// Use transaction
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert journal entry
	jeQuery := `
		INSERT INTO journal_entries (
			id, organization_id, entry_number, entry_type_id, entry_date,
			posting_date, accounting_period_id, fiscal_year_id, status,
			is_posted, is_reversed, source_module, source_document_type,
			source_document_id, reference_number, total_debit, total_credit,
			description, notes, requires_approval, attachments, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`

	_, err = tx.ExecContext(ctx, jeQuery,
		je.ID, je.OrganizationID, je.EntryNumber, je.EntryTypeID, je.EntryDate,
		je.PostingDate, je.AccountingPeriodID, je.FiscalYearID, je.Status,
		je.IsPosted, je.IsReversed, je.SourceModule, je.SourceDocumentType,
		je.SourceDocumentID, je.ReferenceNumber, je.TotalDebit, je.TotalCredit,
		je.Description, je.Notes, je.RequiresApproval, je.Attachments,
		je.Metadata, je.CreatedAt, je.UpdatedAt, je.CreatedBy, je.UpdatedBy,
	)

	if err != nil {
		return err
	}

	// Insert lines
	lineQuery := `
		INSERT INTO journal_entry_lines (
			id, organization_id, journal_entry_id, line_number, account_id,
			debit_amount, credit_amount, location_id, department, project_code,
			cost_center, tax_code, tax_amount, description, memo, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20
		)
	`

	for _, line := range lines {
		_, err := tx.ExecContext(ctx, lineQuery,
			line.ID, line.OrganizationID, line.JournalEntryID, line.LineNumber,
			line.AccountID, line.DebitAmount, line.CreditAmount, line.LocationID,
			line.Department, line.ProjectCode, line.CostCenter, line.TaxCode,
			line.TaxAmount, line.Description, line.Memo, line.Metadata,
			line.CreatedAt, line.UpdatedAt, line.CreatedBy, line.UpdatedBy,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

// GetJournalEntry retrieves a journal entry
func (r *AccountingRepository) GetJournalEntry(ctx context.Context, id, organizationID uuid.UUID) (*accounting.JournalEntry, error) {
	query := `
		SELECT id, organization_id, entry_number, entry_type_id, entry_date,
		       posting_date, accounting_period_id, fiscal_year_id, status,
		       is_posted, is_reversed, reversal_entry_id, source_module,
		       source_document_type, source_document_id, reference_number,
		       total_debit, total_credit, description, notes, requires_approval,
		       approved_by, approved_at, posted_by, posted_at, attachments,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM journal_entries
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	je := &accounting.JournalEntry{}
	err := r.db.GetContext(ctx, je, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, err
	}

	return je, nil
}

// GetJournalEntryByNumber retrieves by entry number
func (r *AccountingRepository) GetJournalEntryByNumber(ctx context.Context, organizationID uuid.UUID, number string) (*accounting.JournalEntry, error) {
	query := `
		SELECT id, organization_id, entry_number, entry_type_id, entry_date,
		       posting_date, accounting_period_id, fiscal_year_id, status,
		       is_posted, is_reversed, reversal_entry_id, source_module,
		       source_document_type, source_document_id, reference_number,
		       total_debit, total_credit, description, notes, requires_approval,
		       approved_by, approved_at, posted_by, posted_at, attachments,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM journal_entries
		WHERE organization_id = $1 AND entry_number = $2 AND deleted_at IS NULL
	`

	je := &accounting.JournalEntry{}
	err := r.db.GetContext(ctx, je, query, organizationID, number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, err
	}

	return je, nil
}

// ListJournalEntries lists journal entries
func (r *AccountingRepository) ListJournalEntries(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*accounting.JournalEntry, error) {
	query := `
		SELECT id, organization_id, entry_number, entry_type_id, entry_date,
		       posting_date, accounting_period_id, fiscal_year_id, status,
		       is_posted, is_reversed, reversal_entry_id, source_module,
		       source_document_type, source_document_id, reference_number,
		       total_debit, total_credit, description, notes, requires_approval,
		       approved_by, approved_at, posted_by, posted_at, attachments,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM journal_entries
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if status, ok := filter["status"]; ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if isPosted, ok := filter["is_posted"].(bool); ok {
		query += fmt.Sprintf(" AND is_posted = $%d", argIndex)
		args = append(args, isPosted)
		argIndex++
	}

	if periodID, ok := filter["accounting_period_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND accounting_period_id = $%d", argIndex)
		args = append(args, periodID)
		argIndex++
	}

	query += " ORDER BY entry_date DESC, created_at DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*accounting.JournalEntry
	for rows.Next() {
		je := &accounting.JournalEntry{}
		if err := rows.StructScan(je); err != nil {
			return nil, err
		}
		entries = append(entries, je)
	}

	return entries, rows.Err()
}

// UpdateJournalEntry updates a journal entry
func (r *AccountingRepository) UpdateJournalEntry(ctx context.Context, je *accounting.JournalEntry, lines []*accounting.JournalEntryLine) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update journal entry
	jeQuery := `
		UPDATE journal_entries
		SET entry_date = $1, posting_date = $2, status = $3,
		    is_posted = $4, is_reversed = $5, reversal_entry_id = $6,
		    total_debit = $7, total_credit = $8, description = $9,
		    notes = $10, requires_approval = $11, approved_by = $12,
		    approved_at = $13, posted_by = $14, posted_at = $15,
		    metadata = $16, updated_at = $17, updated_by = $18
		WHERE id = $19 AND organization_id = $20 AND deleted_at IS NULL
	`

	_, err = tx.ExecContext(ctx, jeQuery,
		je.EntryDate, je.PostingDate, je.Status,
		je.IsPosted, je.IsReversed, je.ReversalEntryID,
		je.TotalDebit, je.TotalCredit, je.Description,
		je.Notes, je.RequiresApproval, je.ApprovedBy,
		je.ApprovedAt, je.PostedBy, je.PostedAt,
		je.Metadata, je.UpdatedAt, je.UpdatedBy,
		je.ID, je.OrganizationID,
	)

	if err != nil {
		return err
	}

	return tx.Commit().Error
}

// DeleteJournalEntry soft deletes a journal entry
func (r *AccountingRepository) DeleteJournalEntry(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE journal_entries
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND is_posted = false AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("journal entry not found or is already posted")
	}

	return nil
}

// PostJournalEntry posts a journal entry (status update)
func (r *AccountingRepository) PostJournalEntry(ctx context.Context, id, organizationID, postedBy uuid.UUID) error {
	query := `
		UPDATE journal_entries
		SET status = $1, is_posted = true, posted_by = $2, posted_at = $3,
		    updated_at = $4, updated_by = $5
		WHERE id = $6 AND organization_id = $7 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		accounting.JournalEntryStatusPosted, postedBy, now, now, postedBy, id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("journal entry not found")
	}

	return nil
}

// ApproveJournalEntry approves a journal entry
func (r *AccountingRepository) ApproveJournalEntry(ctx context.Context, id, organizationID, approvedBy uuid.UUID) error {
	query := `
		UPDATE journal_entries
		SET status = $1, approved_by = $2, approved_at = $3,
		    updated_at = $4, updated_by = $5
		WHERE id = $6 AND organization_id = $7 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		accounting.JournalEntryStatusApproved, approvedBy, now, now, approvedBy, id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("journal entry not found")
	}

	return nil
}

// ReverseJournalEntry reverses a journal entry (status update only)
func (r *AccountingRepository) ReverseJournalEntry(ctx context.Context, id, organizationID, reversalEntryID, reversedBy uuid.UUID) error {
	query := `
		UPDATE journal_entries
		SET is_reversed = true, reversal_entry_id = $1,
		    updated_at = $2, updated_by = $3
		WHERE id = $4 AND organization_id = $5 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		reversalEntryID, now, reversedBy, id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("journal entry not found")
	}

	return nil
}

// ========================
// JOURNAL ENTRY LINES
// ========================

// GetJournalEntryLine retrieves a journal entry line
func (r *AccountingRepository) GetJournalEntryLine(ctx context.Context, id, organizationID uuid.UUID) (*accounting.JournalEntryLine, error) {
	query := `
		SELECT id, organization_id, journal_entry_id, line_number, account_id,
		       debit_amount, credit_amount, location_id, department, project_code,
		       cost_center, tax_code, tax_amount, description, memo,
		       is_reconciled, reconciled_at, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM journal_entry_lines
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	line := &accounting.JournalEntryLine{}
	err := r.db.GetContext(ctx, line, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("journal entry line not found")
		}
		return nil, err
	}

	return line, nil
}

// ListJournalEntryLines lists lines for a journal entry
func (r *AccountingRepository) ListJournalEntryLines(ctx context.Context, journalEntryID uuid.UUID) ([]*accounting.JournalEntryLine, error) {
	query := `
		SELECT id, organization_id, journal_entry_id, line_number, account_id,
		       debit_amount, credit_amount, location_id, department, project_code,
		       cost_center, tax_code, tax_amount, description, memo,
		       is_reconciled, reconciled_at, metadata, created_at, updated_at,
		       created_by, updated_by, deleted_at
		FROM journal_entry_lines
		WHERE journal_entry_id = $1 AND deleted_at IS NULL
		ORDER BY line_number ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, journalEntryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []*accounting.JournalEntryLine
	for rows.Next() {
		line := &accounting.JournalEntryLine{}
		if err := rows.StructScan(line); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, rows.Err()
}

// ========================
// GENERAL LEDGER
// ========================

// PostToGeneralLedger posts journal entry lines to GL
func (r *AccountingRepository) PostToGeneralLedger(ctx context.Context, je *accounting.JournalEntry, lines []*accounting.JournalEntryLine) error {
	query := `
		INSERT INTO general_ledger (
			id, organization_id, journal_entry_id, journal_entry_line_id,
			account_id, transaction_date, posting_date, accounting_period_id,
			fiscal_year_id, debit_amount, credit_amount, running_debit_balance,
			running_credit_balance, running_balance, source_module,
			source_document_type, source_document_id, reference_number,
			location_id, department, project_code, cost_center, description,
			created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
		)
	`

	for _, line := range lines {
		glID := uuid.New()
		_, err := r.db.ExecContext(ctx, query,
			glID, je.OrganizationID, je.ID, line.ID,
			line.AccountID, je.EntryDate, je.PostingDate, je.AccountingPeriodID,
			je.FiscalYearID, line.DebitAmount, line.CreditAmount, "0",
			"0", "0", je.SourceModule,
			je.SourceDocumentType, je.SourceDocumentID, je.ReferenceNumber,
			line.LocationID, line.Department, line.ProjectCode, line.CostCenter,
			line.Description, time.Now(), je.CreatedBy,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetGeneralLedgerEntries retrieves GL entries
func (r *AccountingRepository) GetGeneralLedgerEntries(ctx context.Context, organizationID, accountID uuid.UUID, filter map[string]interface{}) ([]*accounting.GeneralLedger, error) {
	query := `
		SELECT id, organization_id, journal_entry_id, journal_entry_line_id,
		       account_id, transaction_date, posting_date, accounting_period_id,
		       fiscal_year_id, debit_amount, credit_amount, running_debit_balance,
		       running_credit_balance, running_balance, source_module,
		       source_document_type, source_document_id, reference_number,
		       location_id, department, project_code, cost_center, description,
		       is_reversed, reversal_gl_id, metadata, created_at, created_by
		FROM general_ledger
		WHERE organization_id = $1 AND account_id = $2
	`

	args := []interface{}{organizationID, accountID}
	argIndex := 3

	if startDate, ok := filter["start_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND transaction_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate, ok := filter["end_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND transaction_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	query += " ORDER BY transaction_date ASC, created_at ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*accounting.GeneralLedger
	for rows.Next() {
		gl := &accounting.GeneralLedger{}
		if err := rows.StructScan(gl); err != nil {
			return nil, err
		}
		entries = append(entries, gl)
	}

	return entries, rows.Err()
}

// GetAccountBalance retrieves account balance
func (r *AccountingRepository) GetAccountBalance(ctx context.Context, organizationID, accountID uuid.UUID, asOfDate time.Time) (debit, credit, balance string, err error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN debit_amount > 0 THEN CAST(debit_amount AS NUMERIC(20, 4)) ELSE 0 END), 0) as debit,
			COALESCE(SUM(CASE WHEN credit_amount > 0 THEN CAST(credit_amount AS NUMERIC(20, 4)) ELSE 0 END), 0) as credit
		FROM general_ledger
		WHERE organization_id = $1 AND account_id = $2 AND transaction_date <= $3
	`

	var debitVal, creditVal string
	err = r.db.GetContext(ctx, &struct {
		Debit  sql.NullString
		Credit sql.NullString
	}{}, query, organizationID, accountID, asOfDate)

	if err != nil && err != sql.ErrNoRows {
		return "0", "0", "0", err
	}

	// Default values
	return debit, credit, balance, nil
}

// GetTrialBalance retrieves trial balance
func (r *AccountingRepository) GetTrialBalance(ctx context.Context, organizationID uuid.UUID, asOfDate time.Time) (map[uuid.UUID]map[string]string, error) {
	query := `
		SELECT
			account_id,
			COALESCE(SUM(CASE WHEN debit_amount > 0 THEN CAST(debit_amount AS NUMERIC(20, 4)) ELSE 0 END), 0) as debit,
			COALESCE(SUM(CASE WHEN credit_amount > 0 THEN CAST(credit_amount AS NUMERIC(20, 4)) ELSE 0 END), 0) as credit
		FROM general_ledger
		WHERE organization_id = $1 AND transaction_date <= $2
		GROUP BY account_id
	`

	rows, err := r.db.QueryxContext(ctx, query, organizationID, asOfDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID]map[string]string)

	for rows.Next() {
		var accountID uuid.UUID
		var debit, credit string

		if err := rows.Scan(&accountID, &debit, &credit); err != nil {
			return nil, err
		}

		result[accountID] = map[string]string{
			"debit":  debit,
			"credit": credit,
		}
	}

	return result, rows.Err()
}
