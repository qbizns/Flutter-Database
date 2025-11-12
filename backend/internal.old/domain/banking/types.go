package banking

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ========================
// BANK ACCOUNTS
// ========================

type BankAccountType string

const (
	BankAccountTypeChecking      BankAccountType = "checking"
	BankAccountTypeSavings       BankAccountType = "savings"
	BankAccountTypeCreditCard    BankAccountType = "credit_card"
	BankAccountTypeLineOfCredit  BankAccountType = "line_of_credit"
)

func (t BankAccountType) String() string {
	return string(t)
}

type BankAccount struct {
	ID                    uuid.UUID        `db:"id" json:"id"`
	OrganizationID        uuid.UUID        `db:"organization_id" json:"organization_id"`
	ChartAccountID        uuid.UUID        `db:"chart_account_id" json:"chart_account_id"`
	BankName              string           `db:"bank_name" json:"bank_name"`
	AccountNumber         string           `db:"account_number" json:"account_number"`
	AccountType           BankAccountType  `db:"account_type" json:"account_type"`
	RoutingNumber         *string          `db:"routing_number" json:"routing_number,omitempty"`
	SwiftCode             *string          `db:"swift_code" json:"swift_code,omitempty"`
	CurrencyCode          string           `db:"currency_code" json:"currency_code"`
	CurrentBalance        string           `db:"current_balance" json:"current_balance"`
	StatementBalance      string           `db:"statement_balance" json:"statement_balance"`
	LastStatementDate     *time.Time       `db:"last_statement_date" json:"last_statement_date,omitempty"`
	IsActive              bool             `db:"is_active" json:"is_active"`
	OnlineBankingEnabled  bool             `db:"online_banking_enabled" json:"online_banking_enabled"`
	LastSyncDate          *time.Time       `db:"last_sync_date" json:"last_sync_date,omitempty"`
	Notes                 *string          `db:"notes" json:"notes,omitempty"`
	Metadata              json.RawMessage  `db:"metadata" json:"metadata,omitempty"`
	CreatedAt             time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time        `db:"updated_at" json:"updated_at"`
	CreatedBy             *uuid.UUID       `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy             *uuid.UUID       `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt             *time.Time       `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateBankAccountRequest struct {
	ChartAccountID       uuid.UUID              `json:"chart_account_id" validate:"required"`
	BankName             string                 `json:"bank_name" validate:"required,min=1,max=255"`
	AccountNumber        string                 `json:"account_number" validate:"required,min=1,max=100"`
	AccountType          BankAccountType        `json:"account_type" validate:"required"`
	RoutingNumber        *string                `json:"routing_number" validate:"omitempty,max=50"`
	SwiftCode            *string                `json:"swift_code" validate:"omitempty,max=50"`
	CurrencyCode         string                 `json:"currency_code" validate:"required,len=3"`
	IsActive             bool                   `json:"is_active"`
	OnlineBankingEnabled bool                   `json:"online_banking_enabled"`
	Notes                *string                `json:"notes"`
	Metadata             map[string]interface{} `json:"metadata"`
}

type UpdateBankAccountRequest struct {
	BankName             *string                `json:"bank_name" validate:"omitempty,min=1,max=255"`
	AccountType          *BankAccountType       `json:"account_type"`
	RoutingNumber        *string                `json:"routing_number" validate:"omitempty,max=50"`
	SwiftCode            *string                `json:"swift_code" validate:"omitempty,max=50"`
	IsActive             *bool                  `json:"is_active"`
	OnlineBankingEnabled *bool                  `json:"online_banking_enabled"`
	Notes                *string                `json:"notes"`
	Metadata             map[string]interface{} `json:"metadata"`
}

// ========================
// BANK RECONCILIATIONS
// ========================

type BankReconciliationStatus string

const (
	BankReconciliationStatusInProgress BankReconciliationStatus = "in_progress"
	BankReconciliationStatusReconciled BankReconciliationStatus = "reconciled"
	BankReconciliationStatusLocked     BankReconciliationStatus = "locked"
)

func (s BankReconciliationStatus) String() string {
	return string(s)
}

type BankReconciliation struct {
	ID                   uuid.UUID                  `db:"id" json:"id"`
	OrganizationID       uuid.UUID                  `db:"organization_id" json:"organization_id"`
	BankAccountID        uuid.UUID                  `db:"bank_account_id" json:"bank_account_id"`
	StatementDate        time.Time                  `db:"statement_date" json:"statement_date"`
	StatementBalance     string                     `db:"statement_balance" json:"statement_balance"`
	ReconciliationDate   *time.Time                 `db:"reconciliation_date" json:"reconciliation_date,omitempty"`
	BookBalance          string                     `db:"book_balance" json:"book_balance"`
	ClearedBalance       string                     `db:"cleared_balance" json:"cleared_balance"`
	Difference           string                     `db:"difference" json:"difference"`
	Status               BankReconciliationStatus   `db:"status" json:"status"`
	IsReconciled         bool                       `db:"is_reconciled" json:"is_reconciled"`
	AccountingPeriodID   *uuid.UUID                 `db:"accounting_period_id" json:"accounting_period_id,omitempty"`
	Notes                *string                    `db:"notes" json:"notes,omitempty"`
	Metadata             json.RawMessage            `db:"metadata" json:"metadata,omitempty"`
	CreatedAt            time.Time                  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time                  `db:"updated_at" json:"updated_at"`
	ReconciledBy         *uuid.UUID                 `db:"reconciled_by" json:"reconciled_by,omitempty"`
	ReconciledAt         *time.Time                 `db:"reconciled_at" json:"reconciled_at,omitempty"`
	CreatedBy            *uuid.UUID                 `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID                 `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt            *time.Time                 `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateBankReconciliationRequest struct {
	BankAccountID      uuid.UUID              `json:"bank_account_id" validate:"required"`
	StatementDate      time.Time              `json:"statement_date" validate:"required"`
	StatementBalance   string                 `json:"statement_balance" validate:"required"`
	AccountingPeriodID *uuid.UUID             `json:"accounting_period_id"`
	Notes              *string                `json:"notes"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type UpdateBankReconciliationRequest struct {
	StatementBalance string                 `json:"statement_balance"`
	Status           *BankReconciliationStatus `json:"status"`
	Notes            *string                `json:"notes"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// ========================
// BANK RECONCILIATION ITEMS
// ========================

type BankReconciliationItem struct {
	ID                    uuid.UUID   `db:"id" json:"id"`
	OrganizationID        uuid.UUID   `db:"organization_id" json:"organization_id"`
	BankReconciliationID  *uuid.UUID  `db:"bank_reconciliation_id" json:"bank_reconciliation_id,omitempty"`
	GeneralLedgerID       uuid.UUID   `db:"general_ledger_id" json:"general_ledger_id"`
	JournalEntryLineID    uuid.UUID   `db:"journal_entry_line_id" json:"journal_entry_line_id"`
	IsCleared             bool        `db:"is_cleared" json:"is_cleared"`
	ClearedDate           *time.Time  `db:"cleared_date" json:"cleared_date,omitempty"`
	CreatedAt             time.Time   `db:"created_at" json:"created_at"`
	ClearedBy             *uuid.UUID  `db:"cleared_by" json:"cleared_by,omitempty"`
}

type CreateBankReconciliationItemRequest struct {
	BankReconciliationID *uuid.UUID `json:"bank_reconciliation_id"`
	GeneralLedgerID      uuid.UUID  `json:"general_ledger_id" validate:"required"`
	JournalEntryLineID   uuid.UUID  `json:"journal_entry_line_id" validate:"required"`
}

// ========================
// BANK STATEMENTS
// ========================

type BankStatementStatus string

const (
	BankStatementStatusDraft       BankStatementStatus = "draft"
	BankStatementStatusReconciling BankStatementStatus = "reconciling"
	BankStatementStatusReconciled  BankStatementStatus = "reconciled"
	BankStatementStatusClosed      BankStatementStatus = "closed"
)

func (s BankStatementStatus) String() string {
	return string(s)
}

type BankStatementImportSource string

const (
	BankStatementImportSourceManual    BankStatementImportSource = "manual"
	BankStatementImportSourceFileImport BankStatementImportSource = "file_import"
	BankStatementImportSourceAPI       BankStatementImportSource = "api"
	BankStatementImportSourceBankFeed  BankStatementImportSource = "bank_feed"
)

func (s BankStatementImportSource) String() string {
	return string(s)
}

type BankStatement struct {
	ID               uuid.UUID                  `db:"id" json:"id"`
	OrganizationID   uuid.UUID                  `db:"organization_id" json:"organization_id"`
	BankAccountID    uuid.UUID                  `db:"bank_account_id" json:"bank_account_id"`
	StatementNumber  *string                    `db:"statement_number" json:"statement_number,omitempty"`
	StatementDate    time.Time                  `db:"statement_date" json:"statement_date"`
	PeriodStartDate  time.Time                  `db:"period_start_date" json:"period_start_date"`
	PeriodEndDate    time.Time                  `db:"period_end_date" json:"period_end_date"`
	OpeningBalance   string                     `db:"opening_balance" json:"opening_balance"`
	ClosingBalance   string                     `db:"closing_balance" json:"closing_balance"`
	ImportSource     BankStatementImportSource  `db:"import_source" json:"import_source"`
	ImportFileName   *string                    `db:"import_file_name" json:"import_file_name,omitempty"`
	Status           BankStatementStatus        `db:"status" json:"status"`
	Notes            *string                    `db:"notes" json:"notes,omitempty"`
	CreatedAt        time.Time                  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time                  `db:"updated_at" json:"updated_at"`
	CreatedBy        *uuid.UUID                 `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy        *uuid.UUID                 `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt        *time.Time                 `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateBankStatementRequest struct {
	BankAccountID   uuid.UUID                 `json:"bank_account_id" validate:"required"`
	StatementNumber *string                   `json:"statement_number"`
	StatementDate   time.Time                 `json:"statement_date" validate:"required"`
	PeriodStartDate time.Time                 `json:"period_start_date" validate:"required"`
	PeriodEndDate   time.Time                 `json:"period_end_date" validate:"required"`
	OpeningBalance  string                    `json:"opening_balance" validate:"required"`
	ClosingBalance  string                    `json:"closing_balance" validate:"required"`
	ImportSource    BankStatementImportSource `json:"import_source"`
	ImportFileName  *string                   `json:"import_file_name"`
	Notes           *string                   `json:"notes"`
}

type UpdateBankStatementRequest struct {
	StatementNumber *string                 `json:"statement_number"`
	OpeningBalance  *string                 `json:"opening_balance"`
	ClosingBalance  *string                 `json:"closing_balance"`
	Status          *BankStatementStatus    `json:"status"`
	Notes           *string                 `json:"notes"`
}

// ========================
// BANK STATEMENT LINES
// ========================

type BankStatementLineStatus string

const (
	BankStatementLineStatusUnmatched    BankStatementLineStatus = "unmatched"
	BankStatementLineStatusMatched      BankStatementLineStatus = "matched"
	BankStatementLineStatusPartialMatch BankStatementLineStatus = "partial_match"
	BankStatementLineStatusIgnored      BankStatementLineStatus = "ignored"
)

func (s BankStatementLineStatus) String() string {
	return string(s)
}

type BankStatementLine struct {
	ID                  uuid.UUID                   `db:"id" json:"id"`
	BankStatementID     uuid.UUID                   `db:"bank_statement_id" json:"bank_statement_id"`
	LineNumber          int                         `db:"line_number" json:"line_number"`
	TransactionDate     time.Time                   `db:"transaction_date" json:"transaction_date"`
	ValueDate           *time.Time                  `db:"value_date" json:"value_date,omitempty"`
	Amount              string                      `db:"amount" json:"amount"`
	CurrencyCode        string                      `db:"currency_code" json:"currency_code"`
	Description         *string                     `db:"description" json:"description,omitempty"`
	Reference           *string                     `db:"reference" json:"reference,omitempty"`
	CounterpartyName    *string                     `db:"counterparty_name" json:"counterparty_name,omitempty"`
	CounterpartyAccount *string                     `db:"counterparty_account" json:"counterparty_account,omitempty"`
	BankReference       *string                     `db:"bank_reference" json:"bank_reference,omitempty"`
	CheckNumber         *string                     `db:"check_number" json:"check_number,omitempty"`
	Status              BankStatementLineStatus     `db:"status" json:"status"`
	Notes               *string                     `db:"notes" json:"notes,omitempty"`
	CreatedAt           time.Time                   `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time                   `db:"updated_at" json:"updated_at"`
	DeletedAt           *time.Time                  `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateBankStatementLineRequest struct {
	LineNumber          int     `json:"line_number" validate:"required,min=1"`
	TransactionDate     time.Time `json:"transaction_date" validate:"required"`
	ValueDate           *time.Time `json:"value_date"`
	Amount              string    `json:"amount" validate:"required"`
	CurrencyCode        string    `json:"currency_code" validate:"required,len=3"`
	Description         *string   `json:"description"`
	Reference           *string   `json:"reference"`
	CounterpartyName    *string   `json:"counterparty_name"`
	CounterpartyAccount *string   `json:"counterparty_account"`
	BankReference       *string   `json:"bank_reference"`
	CheckNumber         *string   `json:"check_number"`
	Notes               *string   `json:"notes"`
}

type UpdateBankStatementLineRequest struct {
	Amount              *string                      `json:"amount"`
	Description         *string                      `json:"description"`
	Reference           *string                      `json:"reference"`
	CounterpartyName    *string                      `json:"counterparty_name"`
	CounterpartyAccount *string                      `json:"counterparty_account"`
	Status              *BankStatementLineStatus     `json:"status"`
	Notes               *string                      `json:"notes"`
}

// ========================
// BANK STATEMENT RECONCILIATIONS
// ========================

type BankStatementReconciliation struct {
	ID                   uuid.UUID   `db:"id" json:"id"`
	OrganizationID       uuid.UUID   `db:"organization_id" json:"organization_id"`
	BankStatementLineID  uuid.UUID   `db:"bank_statement_line_id" json:"bank_statement_line_id"`
	JournalEntryID       *uuid.UUID  `db:"journal_entry_id" json:"journal_entry_id,omitempty"`
	PaymentID            *uuid.UUID  `db:"payment_id" json:"payment_id,omitempty"`
	MatchedAmount        string      `db:"matched_amount" json:"matched_amount"`
	MatchedBy            *uuid.UUID  `db:"matched_by" json:"matched_by,omitempty"`
	MatchedAt            time.Time   `db:"matched_at" json:"matched_at"`
	DeletedAt            *time.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateBankStatementReconciliationRequest struct {
	BankStatementLineID uuid.UUID  `json:"bank_statement_line_id" validate:"required"`
	JournalEntryID      *uuid.UUID `json:"journal_entry_id"`
	PaymentID           *uuid.UUID `json:"payment_id"`
	MatchedAmount       string     `json:"matched_amount" validate:"required"`
}

// ========================
// RECONCILIATION RULE MODELS
// ========================

type ReconciliationRuleModel struct {
	ID                  uuid.UUID   `db:"id" json:"id"`
	OrganizationID      uuid.UUID   `db:"organization_id" json:"organization_id"`
	RuleName            string      `db:"rule_name" json:"rule_name"`
	RuleCode            *string     `db:"rule_code" json:"rule_code,omitempty"`
	Sequence            int         `db:"sequence" json:"sequence"`
	AmountMin           *string     `db:"amount_min" json:"amount_min,omitempty"`
	AmountMax           *string     `db:"amount_max" json:"amount_max,omitempty"`
	DescriptionPattern  *string     `db:"description_pattern" json:"description_pattern,omitempty"`
	CounterpartyPattern *string     `db:"counterparty_pattern" json:"counterparty_pattern,omitempty"`
	ReferencePattern    *string     `db:"reference_pattern" json:"reference_pattern,omitempty"`
	JournalID           *uuid.UUID  `db:"journal_id" json:"journal_id,omitempty"`
	AccountID           *uuid.UUID  `db:"account_id" json:"account_id,omitempty"`
	AnalyticAccountID   *uuid.UUID  `db:"analytic_account_id" json:"analytic_account_id,omitempty"`
	TaxID               *uuid.UUID  `db:"tax_id" json:"tax_id,omitempty"`
	IsActive            bool        `db:"is_active" json:"is_active"`
	AutoApply           bool        `db:"auto_apply" json:"auto_apply"`
	Notes               *string     `db:"notes" json:"notes,omitempty"`
	CreatedAt           time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time   `db:"updated_at" json:"updated_at"`
	CreatedBy           *uuid.UUID  `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID  `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt           *time.Time  `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateReconciliationRuleModelRequest struct {
	RuleName            string                 `json:"rule_name" validate:"required,min=1,max=255"`
	RuleCode            *string                `json:"rule_code" validate:"omitempty,max=50"`
	Sequence            int                    `json:"sequence"`
	AmountMin           *string                `json:"amount_min"`
	AmountMax           *string                `json:"amount_max"`
	DescriptionPattern  *string                `json:"description_pattern"`
	CounterpartyPattern *string                `json:"counterparty_pattern"`
	ReferencePattern    *string                `json:"reference_pattern"`
	JournalID           *uuid.UUID             `json:"journal_id"`
	AccountID           *uuid.UUID             `json:"account_id"`
	AnalyticAccountID   *uuid.UUID             `json:"analytic_account_id"`
	TaxID               *uuid.UUID             `json:"tax_id"`
	IsActive            bool                   `json:"is_active"`
	AutoApply           bool                   `json:"auto_apply"`
	Notes               *string                `json:"notes"`
}

type UpdateReconciliationRuleModelRequest struct {
	RuleName            *string                `json:"rule_name"`
	RuleCode            *string                `json:"rule_code"`
	Sequence            *int                   `json:"sequence"`
	AmountMin           *string                `json:"amount_min"`
	AmountMax           *string                `json:"amount_max"`
	DescriptionPattern  *string                `json:"description_pattern"`
	CounterpartyPattern *string                `json:"counterparty_pattern"`
	ReferencePattern    *string                `json:"reference_pattern"`
	JournalID           *uuid.UUID             `json:"journal_id"`
	AccountID           *uuid.UUID             `json:"account_id"`
	AnalyticAccountID   *uuid.UUID             `json:"analytic_account_id"`
	TaxID               *uuid.UUID             `json:"tax_id"`
	IsActive            *bool                  `json:"is_active"`
	AutoApply           *bool                  `json:"auto_apply"`
	Notes               *string                `json:"notes"`
}

// ========================
// VALIDATION ERRORS
// ========================

var (
	ErrInvalidAccountType           = errors.New("invalid bank account type")
	ErrInvalidCurrencyCode          = errors.New("invalid currency code")
	ErrInvalidAmount                = errors.New("invalid amount value")
	ErrDuplicateAccountNumber       = errors.New("account number already exists for this organization")
	ErrBankAccountNotFound          = errors.New("bank account not found")
	ErrBankAccountInactive          = errors.New("bank account is inactive")
	ErrReconciliationAlreadyExists  = errors.New("reconciliation already exists for this statement")
	ErrReconciliationNotFound       = errors.New("bank reconciliation not found")
	ErrReconciliationLocked         = errors.New("reconciliation is locked and cannot be modified")
	ErrBankStatementNotFound        = errors.New("bank statement not found")
	ErrBankStatementLineNotFound    = errors.New("bank statement line not found")
	ErrStatementLineAlreadyMatched  = errors.New("statement line is already matched")
	ErrInvalidMatchAmount           = errors.New("matched amount must be positive")
	ErrReconciliationRuleNotFound   = errors.New("reconciliation rule model not found")
	ErrInvalidDateRange             = errors.New("period end date must be after period start date")
	ErrBalanceMismatch              = errors.New("statement balance does not match calculated balance")
	ErrUnmatchedTransactions        = errors.New("unmatched transactions remain")
)

// Value implements the database driver interface for custom scanning
func (t BankAccountType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan implements the database driver interface for custom scanning
func (t *BankAccountType) Scan(value interface{}) error {
	if value == nil {
		*t = ""
		return nil
	}
	val, ok := value.(string)
	if !ok {
		return ErrInvalidAccountType
	}
	*t = BankAccountType(val)
	return nil
}
