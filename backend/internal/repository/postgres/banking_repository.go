package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/your-org/pos-backend/internal/domain/banking"
)

// BankingRepository implements banking.BankingRepository
type BankingRepository struct {
	db *sqlx.DB
}

// NewBankingRepository creates a new banking repository
func NewBankingRepository(db *sqlx.DB) banking.BankingRepository {
	return &BankingRepository{db: db}
}

// ========================
// BANK ACCOUNTS
// ========================

// CreateBankAccount creates a new bank account
func (r *BankingRepository) CreateBankAccount(ctx context.Context, account *banking.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (
			id, organization_id, chart_account_id, bank_name, account_number,
			account_type, routing_number, swift_code, currency_code,
			current_balance, statement_balance, last_statement_date,
			is_active, online_banking_enabled, last_sync_date,
			notes, metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		account.ID, account.OrganizationID, account.ChartAccountID, account.BankName, account.AccountNumber,
		account.AccountType, account.RoutingNumber, account.SwiftCode, account.CurrencyCode,
		account.CurrentBalance, account.StatementBalance, account.LastStatementDate,
		account.IsActive, account.OnlineBankingEnabled, account.LastSyncDate,
		account.Notes, account.Metadata, account.CreatedAt, account.UpdatedAt,
		account.CreatedBy, account.UpdatedBy,
	)

	return err
}

// GetBankAccount retrieves a bank account by ID
func (r *BankingRepository) GetBankAccount(ctx context.Context, id, organizationID uuid.UUID) (*banking.BankAccount, error) {
	query := `
		SELECT id, organization_id, chart_account_id, bank_name, account_number,
		       account_type, routing_number, swift_code, currency_code,
		       current_balance, statement_balance, last_statement_date,
		       is_active, online_banking_enabled, last_sync_date,
		       notes, metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM bank_accounts
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	account := &banking.BankAccount{}
	err := r.db.GetContext(ctx, account, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrBankAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

// GetBankAccountByNumber retrieves a bank account by account number
func (r *BankingRepository) GetBankAccountByNumber(ctx context.Context, organizationID uuid.UUID, accountNumber string) (*banking.BankAccount, error) {
	query := `
		SELECT id, organization_id, chart_account_id, bank_name, account_number,
		       account_type, routing_number, swift_code, currency_code,
		       current_balance, statement_balance, last_statement_date,
		       is_active, online_banking_enabled, last_sync_date,
		       notes, metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM bank_accounts
		WHERE organization_id = $1 AND account_number = $2 AND deleted_at IS NULL
	`

	account := &banking.BankAccount{}
	err := r.db.GetContext(ctx, account, query, organizationID, accountNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrBankAccountNotFound
		}
		return nil, err
	}

	return account, nil
}

// ListBankAccounts lists bank accounts with filters
func (r *BankingRepository) ListBankAccounts(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*banking.BankAccount, error) {
	query := `
		SELECT id, organization_id, chart_account_id, bank_name, account_number,
		       account_type, routing_number, swift_code, currency_code,
		       current_balance, statement_balance, last_statement_date,
		       is_active, online_banking_enabled, last_sync_date,
		       notes, metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM bank_accounts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if accountType, ok := filter["account_type"].(string); ok {
		query += fmt.Sprintf(" AND account_type = $%d", argIndex)
		args = append(args, accountType)
		argIndex++
	}

	if currencyCode, ok := filter["currency_code"].(string); ok {
		query += fmt.Sprintf(" AND currency_code = $%d", argIndex)
		args = append(args, currencyCode)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*banking.BankAccount
	for rows.Next() {
		account := &banking.BankAccount{}
		if err := rows.StructScan(account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, rows.Err()
}

// UpdateBankAccount updates a bank account
func (r *BankingRepository) UpdateBankAccount(ctx context.Context, account *banking.BankAccount) error {
	query := `
		UPDATE bank_accounts
		SET bank_name = $1, account_type = $2, routing_number = $3, swift_code = $4,
		    is_active = $5, online_banking_enabled = $6, last_sync_date = $7,
		    notes = $8, metadata = $9, updated_at = $10, updated_by = $11
		WHERE id = $12 AND organization_id = $13 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		account.BankName, account.AccountType, account.RoutingNumber, account.SwiftCode,
		account.IsActive, account.OnlineBankingEnabled, account.LastSyncDate,
		account.Notes, account.Metadata, account.UpdatedAt, account.UpdatedBy,
		account.ID, account.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrBankAccountNotFound
	}

	return nil
}

// DeleteBankAccount soft deletes a bank account
func (r *BankingRepository) DeleteBankAccount(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE bank_accounts
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
		return banking.ErrBankAccountNotFound
	}

	return nil
}

// ========================
// BANK RECONCILIATIONS
// ========================

// CreateBankReconciliation creates a new bank reconciliation
func (r *BankingRepository) CreateBankReconciliation(ctx context.Context, reconciliation *banking.BankReconciliation) error {
	query := `
		INSERT INTO bank_reconciliations (
			id, organization_id, bank_account_id, statement_date, statement_balance,
			reconciliation_date, book_balance, cleared_balance, difference,
			status, is_reconciled, accounting_period_id, notes, metadata,
			created_at, updated_at, reconciled_by, reconciled_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		reconciliation.ID, reconciliation.OrganizationID, reconciliation.BankAccountID,
		reconciliation.StatementDate, reconciliation.StatementBalance,
		reconciliation.ReconciliationDate, reconciliation.BookBalance, reconciliation.ClearedBalance,
		reconciliation.Difference, reconciliation.Status, reconciliation.IsReconciled,
		reconciliation.AccountingPeriodID, reconciliation.Notes, reconciliation.Metadata,
		reconciliation.CreatedAt, reconciliation.UpdatedAt, reconciliation.ReconciledBy,
		reconciliation.ReconciledAt, reconciliation.CreatedBy, reconciliation.UpdatedBy,
	)

	return err
}

// GetBankReconciliation retrieves a bank reconciliation by ID
func (r *BankingRepository) GetBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*banking.BankReconciliation, error) {
	query := `
		SELECT id, organization_id, bank_account_id, statement_date, statement_balance,
		       reconciliation_date, book_balance, cleared_balance, difference,
		       status, is_reconciled, accounting_period_id, notes, metadata,
		       created_at, updated_at, reconciled_by, reconciled_at, created_by, updated_by, deleted_at
		FROM bank_reconciliations
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	reconciliation := &banking.BankReconciliation{}
	err := r.db.GetContext(ctx, reconciliation, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrReconciliationNotFound
		}
		return nil, err
	}

	return reconciliation, nil
}

// ListBankReconciliations lists bank reconciliations with filters
func (r *BankingRepository) ListBankReconciliations(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*banking.BankReconciliation, error) {
	query := `
		SELECT id, organization_id, bank_account_id, statement_date, statement_balance,
		       reconciliation_date, book_balance, cleared_balance, difference,
		       status, is_reconciled, accounting_period_id, notes, metadata,
		       created_at, updated_at, reconciled_by, reconciled_at, created_by, updated_by, deleted_at
		FROM bank_reconciliations
		WHERE organization_id = $1 AND bank_account_id = $2 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID, bankAccountID}
	argIndex := 3

	if status, ok := filter["status"].(string); ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if isReconciled, ok := filter["is_reconciled"].(bool); ok {
		query += fmt.Sprintf(" AND is_reconciled = $%d", argIndex)
		args = append(args, isReconciled)
		argIndex++
	}

	query += " ORDER BY statement_date DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reconciliations []*banking.BankReconciliation
	for rows.Next() {
		reconciliation := &banking.BankReconciliation{}
		if err := rows.StructScan(reconciliation); err != nil {
			return nil, err
		}
		reconciliations = append(reconciliations, reconciliation)
	}

	return reconciliations, rows.Err()
}

// UpdateBankReconciliation updates a bank reconciliation
func (r *BankingRepository) UpdateBankReconciliation(ctx context.Context, reconciliation *banking.BankReconciliation) error {
	query := `
		UPDATE bank_reconciliations
		SET statement_balance = $1, reconciliation_date = $2, book_balance = $3,
		    cleared_balance = $4, difference = $5, status = $6, is_reconciled = $7,
		    notes = $8, metadata = $9, updated_at = $10, updated_by = $11
		WHERE id = $12 AND organization_id = $13 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		reconciliation.StatementBalance, reconciliation.ReconciliationDate, reconciliation.BookBalance,
		reconciliation.ClearedBalance, reconciliation.Difference, reconciliation.Status,
		reconciliation.IsReconciled, reconciliation.Notes, reconciliation.Metadata,
		reconciliation.UpdatedAt, reconciliation.UpdatedBy,
		reconciliation.ID, reconciliation.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrReconciliationNotFound
	}

	return nil
}

// DeleteBankReconciliation soft deletes a bank reconciliation
func (r *BankingRepository) DeleteBankReconciliation(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE bank_reconciliations
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
		return banking.ErrReconciliationNotFound
	}

	return nil
}

// MarkReconciliationComplete marks a reconciliation as complete
func (r *BankingRepository) MarkReconciliationComplete(ctx context.Context, id, organizationID, reconciledBy uuid.UUID) error {
	query := `
		UPDATE bank_reconciliations
		SET status = $1, is_reconciled = $2, reconciled_by = $3, reconciled_at = $4, updated_at = $5
		WHERE id = $6 AND organization_id = $7 AND deleted_at IS NULL
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		banking.BankReconciliationStatusReconciled, true, reconciledBy, now, now,
		id, organizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrReconciliationNotFound
	}

	return nil
}

// ========================
// BANK RECONCILIATION ITEMS
// ========================

// CreateBankReconciliationItem creates a new reconciliation item
func (r *BankingRepository) CreateBankReconciliationItem(ctx context.Context, item *banking.BankReconciliationItem) error {
	query := `
		INSERT INTO bank_reconciliation_items (
			id, organization_id, bank_reconciliation_id, general_ledger_id, journal_entry_line_id,
			is_cleared, cleared_date, created_at, cleared_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.OrganizationID, item.BankReconciliationID, item.GeneralLedgerID,
		item.JournalEntryLineID, item.IsCleared, item.ClearedDate, item.CreatedAt, item.ClearedBy,
	)

	return err
}

// GetBankReconciliationItem retrieves a reconciliation item by ID
func (r *BankingRepository) GetBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) (*banking.BankReconciliationItem, error) {
	query := `
		SELECT id, organization_id, bank_reconciliation_id, general_ledger_id, journal_entry_line_id,
		       is_cleared, cleared_date, created_at, cleared_by
		FROM bank_reconciliation_items
		WHERE id = $1 AND organization_id = $2
	`

	item := &banking.BankReconciliationItem{}
	err := r.db.GetContext(ctx, item, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("reconciliation item not found")
		}
		return nil, err
	}

	return item, nil
}

// ListBankReconciliationItems lists reconciliation items for a reconciliation
func (r *BankingRepository) ListBankReconciliationItems(ctx context.Context, reconciliationID uuid.UUID) ([]*banking.BankReconciliationItem, error) {
	query := `
		SELECT id, organization_id, bank_reconciliation_id, general_ledger_id, journal_entry_line_id,
		       is_cleared, cleared_date, created_at, cleared_by
		FROM bank_reconciliation_items
		WHERE bank_reconciliation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, reconciliationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*banking.BankReconciliationItem
	for rows.Next() {
		item := &banking.BankReconciliationItem{}
		if err := rows.StructScan(item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// UpdateBankReconciliationItem updates a reconciliation item
func (r *BankingRepository) UpdateBankReconciliationItem(ctx context.Context, item *banking.BankReconciliationItem) error {
	query := `
		UPDATE bank_reconciliation_items
		SET is_cleared = $1, cleared_date = $2, cleared_by = $3
		WHERE id = $4 AND organization_id = $5
	`

	result, err := r.db.ExecContext(ctx, query,
		item.IsCleared, item.ClearedDate, item.ClearedBy,
		item.ID, item.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("reconciliation item not found")
	}

	return nil
}

// DeleteBankReconciliationItem soft deletes a reconciliation item
func (r *BankingRepository) DeleteBankReconciliationItem(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		DELETE FROM bank_reconciliation_items
		WHERE id = $1 AND organization_id = $2
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
		return fmt.Errorf("reconciliation item not found")
	}

	return nil
}

// MarkReconciliationItemCleared marks an item as cleared
func (r *BankingRepository) MarkReconciliationItemCleared(ctx context.Context, id, organizationID, clearedBy uuid.UUID) error {
	query := `
		UPDATE bank_reconciliation_items
		SET is_cleared = true, cleared_date = $1, cleared_by = $2
		WHERE id = $3 AND organization_id = $4
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, now, clearedBy, id, organizationID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("reconciliation item not found")
	}

	return nil
}

// ========================
// BANK STATEMENTS
// ========================

// CreateBankStatement creates a new bank statement
func (r *BankingRepository) CreateBankStatement(ctx context.Context, statement *banking.BankStatement) error {
	query := `
		INSERT INTO bank_statements (
			id, organization_id, bank_account_id, statement_number, statement_date,
			period_start_date, period_end_date, opening_balance, closing_balance,
			import_source, import_file_name, status, notes,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		statement.ID, statement.OrganizationID, statement.BankAccountID,
		statement.StatementNumber, statement.StatementDate,
		statement.PeriodStartDate, statement.PeriodEndDate,
		statement.OpeningBalance, statement.ClosingBalance,
		statement.ImportSource, statement.ImportFileName, statement.Status,
		statement.Notes, statement.CreatedAt, statement.UpdatedAt,
		statement.CreatedBy, statement.UpdatedBy,
	)

	return err
}

// GetBankStatement retrieves a bank statement by ID
func (r *BankingRepository) GetBankStatement(ctx context.Context, id, organizationID uuid.UUID) (*banking.BankStatement, error) {
	query := `
		SELECT id, organization_id, bank_account_id, statement_number, statement_date,
		       period_start_date, period_end_date, opening_balance, closing_balance,
		       import_source, import_file_name, status, notes,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM bank_statements
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	statement := &banking.BankStatement{}
	err := r.db.GetContext(ctx, statement, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrBankStatementNotFound
		}
		return nil, err
	}

	return statement, nil
}

// ListBankStatements lists bank statements with filters
func (r *BankingRepository) ListBankStatements(ctx context.Context, organizationID, bankAccountID uuid.UUID, filter map[string]interface{}) ([]*banking.BankStatement, error) {
	query := `
		SELECT id, organization_id, bank_account_id, statement_number, statement_date,
		       period_start_date, period_end_date, opening_balance, closing_balance,
		       import_source, import_file_name, status, notes,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM bank_statements
		WHERE organization_id = $1 AND bank_account_id = $2 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID, bankAccountID}
	argIndex := 3

	if status, ok := filter["status"].(string); ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if importSource, ok := filter["import_source"].(string); ok {
		query += fmt.Sprintf(" AND import_source = $%d", argIndex)
		args = append(args, importSource)
		argIndex++
	}

	query += " ORDER BY statement_date DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []*banking.BankStatement
	for rows.Next() {
		statement := &banking.BankStatement{}
		if err := rows.StructScan(statement); err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	return statements, rows.Err()
}

// UpdateBankStatement updates a bank statement
func (r *BankingRepository) UpdateBankStatement(ctx context.Context, statement *banking.BankStatement) error {
	query := `
		UPDATE bank_statements
		SET statement_number = $1, opening_balance = $2, closing_balance = $3,
		    status = $4, notes = $5, updated_at = $6, updated_by = $7
		WHERE id = $8 AND organization_id = $9 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		statement.StatementNumber, statement.OpeningBalance, statement.ClosingBalance,
		statement.Status, statement.Notes, statement.UpdatedAt, statement.UpdatedBy,
		statement.ID, statement.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrBankStatementNotFound
	}

	return nil
}

// DeleteBankStatement soft deletes a bank statement
func (r *BankingRepository) DeleteBankStatement(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE bank_statements
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
		return banking.ErrBankStatementNotFound
	}

	return nil
}

// ========================
// BANK STATEMENT LINES
// ========================

// CreateBankStatementLine creates a new statement line
func (r *BankingRepository) CreateBankStatementLine(ctx context.Context, line *banking.BankStatementLine) error {
	query := `
		INSERT INTO bank_statement_lines (
			id, bank_statement_id, line_number, transaction_date, value_date,
			amount, currency_code, description, reference,
			counterparty_name, counterparty_account, bank_reference, check_number,
			status, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		line.ID, line.BankStatementID, line.LineNumber, line.TransactionDate, line.ValueDate,
		line.Amount, line.CurrencyCode, line.Description, line.Reference,
		line.CounterpartyName, line.CounterpartyAccount, line.BankReference, line.CheckNumber,
		line.Status, line.Notes, line.CreatedAt, line.UpdatedAt,
	)

	return err
}

// GetBankStatementLine retrieves a statement line by ID
func (r *BankingRepository) GetBankStatementLine(ctx context.Context, id uuid.UUID) (*banking.BankStatementLine, error) {
	query := `
		SELECT id, bank_statement_id, line_number, transaction_date, value_date,
		       amount, currency_code, description, reference,
		       counterparty_name, counterparty_account, bank_reference, check_number,
		       status, notes, created_at, updated_at, deleted_at
		FROM bank_statement_lines
		WHERE id = $1 AND deleted_at IS NULL
	`

	line := &banking.BankStatementLine{}
	err := r.db.GetContext(ctx, line, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrBankStatementLineNotFound
		}
		return nil, err
	}

	return line, nil
}

// ListBankStatementLines lists statement lines with filters
func (r *BankingRepository) ListBankStatementLines(ctx context.Context, statementID uuid.UUID, filter map[string]interface{}) ([]*banking.BankStatementLine, error) {
	query := `
		SELECT id, bank_statement_id, line_number, transaction_date, value_date,
		       amount, currency_code, description, reference,
		       counterparty_name, counterparty_account, bank_reference, check_number,
		       status, notes, created_at, updated_at, deleted_at
		FROM bank_statement_lines
		WHERE bank_statement_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{statementID}
	argIndex := 2

	if status, ok := filter["status"].(string); ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	query += " ORDER BY line_number ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []*banking.BankStatementLine
	for rows.Next() {
		line := &banking.BankStatementLine{}
		if err := rows.StructScan(line); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	return lines, rows.Err()
}

// UpdateBankStatementLine updates a statement line
func (r *BankingRepository) UpdateBankStatementLine(ctx context.Context, line *banking.BankStatementLine) error {
	query := `
		UPDATE bank_statement_lines
		SET amount = $1, description = $2, reference = $3,
		    counterparty_name = $4, counterparty_account = $5,
		    status = $6, notes = $7, updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		line.Amount, line.Description, line.Reference,
		line.CounterpartyName, line.CounterpartyAccount,
		line.Status, line.Notes, line.UpdatedAt,
		line.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrBankStatementLineNotFound
	}

	return nil
}

// DeleteBankStatementLine soft deletes a statement line
func (r *BankingRepository) DeleteBankStatementLine(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE bank_statement_lines
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrBankStatementLineNotFound
	}

	return nil
}

// ========================
// BANK STATEMENT RECONCILIATIONS
// ========================

// CreateBankStatementReconciliation creates a new statement reconciliation
func (r *BankingRepository) CreateBankStatementReconciliation(ctx context.Context, reconciliation *banking.BankStatementReconciliation) error {
	query := `
		INSERT INTO bank_statement_reconciliations (
			id, organization_id, bank_statement_line_id, journal_entry_id, payment_id,
			matched_amount, matched_by, matched_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		reconciliation.ID, reconciliation.OrganizationID, reconciliation.BankStatementLineID,
		reconciliation.JournalEntryID, reconciliation.PaymentID,
		reconciliation.MatchedAmount, reconciliation.MatchedBy, reconciliation.MatchedAt,
	)

	return err
}

// GetBankStatementReconciliation retrieves a statement reconciliation by ID
func (r *BankingRepository) GetBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) (*banking.BankStatementReconciliation, error) {
	query := `
		SELECT id, organization_id, bank_statement_line_id, journal_entry_id, payment_id,
		       matched_amount, matched_by, matched_at, deleted_at
		FROM bank_statement_reconciliations
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	reconciliation := &banking.BankStatementReconciliation{}
	err := r.db.GetContext(ctx, reconciliation, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("statement reconciliation not found")
		}
		return nil, err
	}

	return reconciliation, nil
}

// ListBankStatementReconciliations lists statement reconciliations for a line
func (r *BankingRepository) ListBankStatementReconciliations(ctx context.Context, statementLineID uuid.UUID) ([]*banking.BankStatementReconciliation, error) {
	query := `
		SELECT id, organization_id, bank_statement_line_id, journal_entry_id, payment_id,
		       matched_amount, matched_by, matched_at, deleted_at
		FROM bank_statement_reconciliations
		WHERE bank_statement_line_id = $1 AND deleted_at IS NULL
		ORDER BY matched_at DESC
	`

	rows, err := r.db.QueryxContext(ctx, query, statementLineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reconciliations []*banking.BankStatementReconciliation
	for rows.Next() {
		reconciliation := &banking.BankStatementReconciliation{}
		if err := rows.StructScan(reconciliation); err != nil {
			return nil, err
		}
		reconciliations = append(reconciliations, reconciliation)
	}

	return reconciliations, rows.Err()
}

// DeleteBankStatementReconciliation soft deletes a statement reconciliation
func (r *BankingRepository) DeleteBankStatementReconciliation(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE bank_statement_reconciliations
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
		return fmt.Errorf("statement reconciliation not found")
	}

	return nil
}

// ========================
// RECONCILIATION RULE MODELS
// ========================

// CreateReconciliationRuleModel creates a new reconciliation rule
func (r *BankingRepository) CreateReconciliationRuleModel(ctx context.Context, rule *banking.ReconciliationRuleModel) error {
	query := `
		INSERT INTO reconciliation_rule_models (
			id, organization_id, rule_name, rule_code, sequence,
			amount_min, amount_max, description_pattern, counterparty_pattern, reference_pattern,
			journal_id, account_id, analytic_account_id, tax_id,
			is_active, auto_apply, notes,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		rule.ID, rule.OrganizationID, rule.RuleName, rule.RuleCode, rule.Sequence,
		rule.AmountMin, rule.AmountMax, rule.DescriptionPattern, rule.CounterpartyPattern,
		rule.ReferencePattern, rule.JournalID, rule.AccountID, rule.AnalyticAccountID,
		rule.TaxID, rule.IsActive, rule.AutoApply, rule.Notes,
		rule.CreatedAt, rule.UpdatedAt, rule.CreatedBy, rule.UpdatedBy,
	)

	return err
}

// GetReconciliationRuleModel retrieves a rule by ID
func (r *BankingRepository) GetReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) (*banking.ReconciliationRuleModel, error) {
	query := `
		SELECT id, organization_id, rule_name, rule_code, sequence,
		       amount_min, amount_max, description_pattern, counterparty_pattern, reference_pattern,
		       journal_id, account_id, analytic_account_id, tax_id,
		       is_active, auto_apply, notes,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM reconciliation_rule_models
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rule := &banking.ReconciliationRuleModel{}
	err := r.db.GetContext(ctx, rule, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, banking.ErrReconciliationRuleNotFound
		}
		return nil, err
	}

	return rule, nil
}

// ListReconciliationRuleModels lists rules with filters
func (r *BankingRepository) ListReconciliationRuleModels(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*banking.ReconciliationRuleModel, error) {
	query := `
		SELECT id, organization_id, rule_name, rule_code, sequence,
		       amount_min, amount_max, description_pattern, counterparty_pattern, reference_pattern,
		       journal_id, account_id, analytic_account_id, tax_id,
		       is_active, auto_apply, notes,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM reconciliation_rule_models
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	if autoApply, ok := filter["auto_apply"].(bool); ok {
		query += fmt.Sprintf(" AND auto_apply = $%d", argIndex)
		args = append(args, autoApply)
		argIndex++
	}

	query += " ORDER BY sequence ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*banking.ReconciliationRuleModel
	for rows.Next() {
		rule := &banking.ReconciliationRuleModel{}
		if err := rows.StructScan(rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

// UpdateReconciliationRuleModel updates a rule
func (r *BankingRepository) UpdateReconciliationRuleModel(ctx context.Context, rule *banking.ReconciliationRuleModel) error {
	query := `
		UPDATE reconciliation_rule_models
		SET rule_name = $1, rule_code = $2, sequence = $3,
		    amount_min = $4, amount_max = $5, description_pattern = $6,
		    counterparty_pattern = $7, reference_pattern = $8,
		    journal_id = $9, account_id = $10, analytic_account_id = $11, tax_id = $12,
		    is_active = $13, auto_apply = $14, notes = $15,
		    updated_at = $16, updated_by = $17
		WHERE id = $18 AND organization_id = $19 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		rule.RuleName, rule.RuleCode, rule.Sequence,
		rule.AmountMin, rule.AmountMax, rule.DescriptionPattern,
		rule.CounterpartyPattern, rule.ReferencePattern,
		rule.JournalID, rule.AccountID, rule.AnalyticAccountID, rule.TaxID,
		rule.IsActive, rule.AutoApply, rule.Notes,
		rule.UpdatedAt, rule.UpdatedBy,
		rule.ID, rule.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return banking.ErrReconciliationRuleNotFound
	}

	return nil
}

// DeleteReconciliationRuleModel soft deletes a rule
func (r *BankingRepository) DeleteReconciliationRuleModel(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE reconciliation_rule_models
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
		return banking.ErrReconciliationRuleNotFound
	}

	return nil
}
