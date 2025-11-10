package accounting

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ========================
// FISCAL YEARS
// ========================

type FiscalYearStatus string

const (
	FiscalYearStatusOpen   FiscalYearStatus = "open"
	FiscalYearStatusClosed FiscalYearStatus = "closed"
	FiscalYearStatusLocked FiscalYearStatus = "locked"
)

func (s FiscalYearStatus) String() string {
	return string(s)
}

type FiscalYear struct {
	ID            uuid.UUID       `db:"id" json:"id"`
	OrganizationID uuid.UUID      `db:"organization_id" json:"organization_id"`
	FiscalYear    string          `db:"fiscal_year" json:"fiscal_year"`
	StartDate     time.Time       `db:"start_date" json:"start_date"`
	EndDate       time.Time       `db:"end_date" json:"end_date"`
	Status        FiscalYearStatus `db:"status" json:"status"`
	IsCurrent     bool            `db:"is_current" json:"is_current"`
	ClosedBy      *uuid.UUID      `db:"closed_by" json:"closed_by,omitempty"`
	ClosedAt      *time.Time      `db:"closed_at" json:"closed_at,omitempty"`
	Notes         *string         `db:"notes" json:"notes,omitempty"`
	Metadata      json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
	CreatedBy     *uuid.UUID      `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy     *uuid.UUID      `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt     *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateFiscalYearRequest struct {
	FiscalYear string    `json:"fiscal_year" validate:"required,min=4,max=10"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
	IsCurrent  bool      `json:"is_current"`
	Notes      *string   `json:"notes"`
	Metadata   map[string]interface{} `json:"metadata"`
}

type UpdateFiscalYearRequest struct {
	FiscalYear *string `json:"fiscal_year" validate:"omitempty,min=4,max=10"`
	StartDate  *time.Time `json:"start_date"`
	EndDate    *time.Time `json:"end_date"`
	IsCurrent  *bool   `json:"is_current"`
	Notes      *string `json:"notes"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ========================
// ACCOUNTING PERIODS
// ========================

type AccountingPeriodStatus string

const (
	AccountingPeriodStatusOpen   AccountingPeriodStatus = "open"
	AccountingPeriodStatusClosed AccountingPeriodStatus = "closed"
	AccountingPeriodStatusLocked AccountingPeriodStatus = "locked"
)

func (s AccountingPeriodStatus) String() string {
	return string(s)
}

type AccountingPeriod struct {
	ID             uuid.UUID              `db:"id" json:"id"`
	OrganizationID uuid.UUID              `db:"organization_id" json:"organization_id"`
	FiscalYearID   uuid.UUID              `db:"fiscal_year_id" json:"fiscal_year_id"`
	PeriodNumber   int                    `db:"period_number" json:"period_number"`
	PeriodName     string                 `db:"period_name" json:"period_name"`
	StartDate      time.Time              `db:"start_date" json:"start_date"`
	EndDate        time.Time              `db:"end_date" json:"end_date"`
	Status         AccountingPeriodStatus `db:"status" json:"status"`
	ClosedBy       *uuid.UUID             `db:"closed_by" json:"closed_by,omitempty"`
	ClosedAt       *time.Time             `db:"closed_at" json:"closed_at,omitempty"`
	Metadata       json.RawMessage        `db:"metadata" json:"metadata,omitempty"`
	CreatedAt      time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at" json:"updated_at"`
	CreatedBy      *uuid.UUID             `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID             `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt      *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateAccountingPeriodRequest struct {
	FiscalYearID uuid.UUID              `json:"fiscal_year_id" validate:"required"`
	PeriodNumber int                    `json:"period_number" validate:"required,min=1,max=12"`
	PeriodName   string                 `json:"period_name" validate:"required"`
	StartDate    time.Time              `json:"start_date" validate:"required"`
	EndDate      time.Time              `json:"end_date" validate:"required"`
	Metadata     map[string]interface{} `json:"metadata"`
}

type UpdateAccountingPeriodRequest struct {
	PeriodName *string                `json:"period_name"`
	Status     *AccountingPeriodStatus `json:"status"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ========================
// ACCOUNT TYPES & SUBTYPES
// ========================

type AccountTypeCategory string

const (
	AccountTypeCategoryBalanceSheet     AccountTypeCategory = "balance_sheet"
	AccountTypeCategoryIncomeStatement  AccountTypeCategory = "income_statement"
)

type NormalBalance string

const (
	NormalBalanceDebit  NormalBalance = "debit"
	NormalBalanceCredit NormalBalance = "credit"
)

type AccountType struct {
	ID                 uuid.UUID            `db:"id" json:"id"`
	TypeCode           string               `db:"type_code" json:"type_code"`
	TypeName           string               `db:"type_name" json:"type_name"`
	TypeCategory       AccountTypeCategory  `db:"type_category" json:"type_category"`
	NormalBalance      NormalBalance        `db:"normal_balance" json:"normal_balance"`
	IsBalanceSheet     bool                 `db:"is_balance_sheet" json:"is_balance_sheet"`
	IsIncomeStatement  bool                 `db:"is_income_statement" json:"is_income_statement"`
	DisplayOrder       int                  `db:"display_order" json:"display_order"`
	Description        *string              `db:"description" json:"description,omitempty"`
	CreatedAt          time.Time            `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time            `db:"updated_at" json:"updated_at"`
}

type AccountSubtype struct {
	ID             uuid.UUID `db:"id" json:"id"`
	AccountTypeID  uuid.UUID `db:"account_type_id" json:"account_type_id"`
	SubtypeCode    string    `db:"subtype_code" json:"subtype_code"`
	SubtypeName    string    `db:"subtype_name" json:"subtype_name"`
	DisplayOrder   int       `db:"display_order" json:"display_order"`
	Description    *string   `db:"description" json:"description,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// ========================
// CHART OF ACCOUNTS
// ========================

type ChartOfAccount struct {
	ID                  uuid.UUID       `db:"id" json:"id"`
	OrganizationID      uuid.UUID       `db:"organization_id" json:"organization_id"`
	AccountCode         string          `db:"account_code" json:"account_code"`
	AccountNumber       string          `db:"account_number" json:"account_number"`
	AccountName         string          `db:"account_name" json:"account_name"`
	AccountTypeID       uuid.UUID       `db:"account_type_id" json:"account_type_id"`
	AccountSubtypeID    *uuid.UUID      `db:"account_subtype_id" json:"account_subtype_id,omitempty"`
	ParentAccountID     *uuid.UUID      `db:"parent_account_id" json:"parent_account_id,omitempty"`
	AccountLevel        int             `db:"account_level" json:"account_level"`
	AccountPath         *string         `db:"account_path" json:"account_path,omitempty"`
	IsActive            bool            `db:"is_active" json:"is_active"`
	IsSystemAccount     bool            `db:"is_system_account" json:"is_system_account"`
	IsHeaderAccount     bool            `db:"is_header_account" json:"is_header_account"`
	IsBankAccount       bool            `db:"is_bank_account" json:"is_bank_account"`
	IsReconcilable      bool            `db:"is_reconcilable" json:"is_reconcilable"`
	DefaultTaxCode      *string         `db:"default_tax_code" json:"default_tax_code,omitempty"`
	CurrencyCode        string          `db:"currency_code" json:"currency_code"`
	OpeningBalance      string          `db:"opening_balance" json:"opening_balance"`
	OpeningBalanceDate  *time.Time      `db:"opening_balance_date" json:"opening_balance_date,omitempty"`
	CurrentDebitBalance string          `db:"current_debit_balance" json:"current_debit_balance"`
	CurrentCreditBalance string         `db:"current_credit_balance" json:"current_credit_balance"`
	CurrentBalance      string          `db:"current_balance" json:"current_balance"`
	LastBalanceUpdate   *time.Time      `db:"last_balance_update" json:"last_balance_update,omitempty"`
	Description         *string         `db:"description" json:"description,omitempty"`
	Notes               *string         `db:"notes" json:"notes,omitempty"`
	Metadata            json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt           time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time       `db:"updated_at" json:"updated_at"`
	CreatedBy           *uuid.UUID      `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID      `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt           *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateChartOfAccountRequest struct {
	AccountCode       string                 `json:"account_code" validate:"required,min=1,max=50"`
	AccountNumber     string                 `json:"account_number" validate:"required,min=1,max=50"`
	AccountName       string                 `json:"account_name" validate:"required,min=1,max=255"`
	AccountTypeID     uuid.UUID              `json:"account_type_id" validate:"required"`
	AccountSubtypeID  *uuid.UUID             `json:"account_subtype_id"`
	ParentAccountID   *uuid.UUID             `json:"parent_account_id"`
	IsActive          bool                   `json:"is_active"`
	IsSystemAccount   bool                   `json:"is_system_account"`
	IsHeaderAccount   bool                   `json:"is_header_account"`
	IsBankAccount     bool                   `json:"is_bank_account"`
	IsReconcilable    bool                   `json:"is_reconcilable"`
	DefaultTaxCode    *string                `json:"default_tax_code"`
	CurrencyCode      string                 `json:"currency_code"`
	OpeningBalance    *string                `json:"opening_balance"`
	OpeningBalanceDate *time.Time            `json:"opening_balance_date"`
	Description       *string                `json:"description"`
	Notes             *string                `json:"notes"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type UpdateChartOfAccountRequest struct {
	AccountName       *string                `json:"account_name"`
	AccountSubtypeID  *uuid.UUID             `json:"account_subtype_id"`
	IsActive          *bool                  `json:"is_active"`
	IsHeaderAccount   *bool                  `json:"is_header_account"`
	IsBankAccount     *bool                  `json:"is_bank_account"`
	IsReconcilable    *bool                  `json:"is_reconcilable"`
	DefaultTaxCode    *string                `json:"default_tax_code"`
	Description       *string                `json:"description"`
	Notes             *string                `json:"notes"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// ========================
// JOURNAL ENTRY TYPES
// ========================

type JournalEntryTypeCategory string

const (
	JournalEntryTypeCategoryStandard    JournalEntryTypeCategory = "standard"
	JournalEntryTypeCategoryAdjusting   JournalEntryTypeCategory = "adjusting"
	JournalEntryTypeCategoryClosing     JournalEntryTypeCategory = "closing"
	JournalEntryTypeCategoryReversing   JournalEntryTypeCategory = "reversing"
)

type JournalEntryType struct {
	ID          uuid.UUID                   `db:"id" json:"id"`
	TypeCode    string                      `db:"type_code" json:"type_code"`
	TypeName    string                      `db:"type_name" json:"type_name"`
	TypeCategory *JournalEntryTypeCategory  `db:"type_category" json:"type_category,omitempty"`
	NumberPrefix *string                    `db:"number_prefix" json:"number_prefix,omitempty"`
	Description *string                     `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time                   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time                   `db:"updated_at" json:"updated_at"`
}

// ========================
// JOURNAL ENTRIES
// ========================

type JournalEntryStatus string

const (
	JournalEntryStatusDraft    JournalEntryStatus = "draft"
	JournalEntryStatusPosted   JournalEntryStatus = "posted"
	JournalEntryStatusApproved JournalEntryStatus = "approved"
	JournalEntryStatusReversed JournalEntryStatus = "reversed"
	JournalEntryStatusVoided   JournalEntryStatus = "voided"
)

func (s JournalEntryStatus) String() string {
	return string(s)
}

type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

type JournalEntry struct {
	ID                   uuid.UUID               `db:"id" json:"id"`
	OrganizationID       uuid.UUID               `db:"organization_id" json:"organization_id"`
	EntryNumber          string                  `db:"entry_number" json:"entry_number"`
	EntryTypeID          uuid.UUID               `db:"entry_type_id" json:"entry_type_id"`
	EntryDate            time.Time               `db:"entry_date" json:"entry_date"`
	PostingDate          time.Time               `db:"posting_date" json:"posting_date"`
	AccountingPeriodID   *uuid.UUID              `db:"accounting_period_id" json:"accounting_period_id,omitempty"`
	FiscalYearID         *uuid.UUID              `db:"fiscal_year_id" json:"fiscal_year_id,omitempty"`
	Status               JournalEntryStatus      `db:"status" json:"status"`
	IsPosted             bool                    `db:"is_posted" json:"is_posted"`
	IsReversed           bool                    `db:"is_reversed" json:"is_reversed"`
	ReversalEntryID      *uuid.UUID              `db:"reversal_entry_id" json:"reversal_entry_id,omitempty"`
	SourceModule         *string                 `db:"source_module" json:"source_module,omitempty"`
	SourceDocumentType   *string                 `db:"source_document_type" json:"source_document_type,omitempty"`
	SourceDocumentID     *uuid.UUID              `db:"source_document_id" json:"source_document_id,omitempty"`
	ReferenceNumber      *string                 `db:"reference_number" json:"reference_number,omitempty"`
	TotalDebit           string                  `db:"total_debit" json:"total_debit"`
	TotalCredit          string                  `db:"total_credit" json:"total_credit"`
	Description          string                  `db:"description" json:"description"`
	Notes                *string                 `db:"notes" json:"notes,omitempty"`
	RequiresApproval     bool                    `db:"requires_approval" json:"requires_approval"`
	ApprovedBy           *uuid.UUID              `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt           *time.Time              `db:"approved_at" json:"approved_at,omitempty"`
	PostedBy             *uuid.UUID              `db:"posted_by" json:"posted_by,omitempty"`
	PostedAt             *time.Time              `db:"posted_at" json:"posted_at,omitempty"`
	Attachments          json.RawMessage         `db:"attachments" json:"attachments,omitempty"`
	Metadata             json.RawMessage         `db:"metadata" json:"metadata,omitempty"`
	CreatedAt            time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time               `db:"updated_at" json:"updated_at"`
	CreatedBy            *uuid.UUID              `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID              `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt            *time.Time              `db:"deleted_at" json:"deleted_at,omitempty"`
}

type JournalEntryLine struct {
	ID                 uuid.UUID       `db:"id" json:"id"`
	OrganizationID     uuid.UUID       `db:"organization_id" json:"organization_id"`
	JournalEntryID     uuid.UUID       `db:"journal_entry_id" json:"journal_entry_id"`
	LineNumber         int             `db:"line_number" json:"line_number"`
	AccountID          uuid.UUID       `db:"account_id" json:"account_id"`
	DebitAmount        string          `db:"debit_amount" json:"debit_amount"`
	CreditAmount       string          `db:"credit_amount" json:"credit_amount"`
	LocationID         *uuid.UUID      `db:"location_id" json:"location_id,omitempty"`
	Department         *string         `db:"department" json:"department,omitempty"`
	ProjectCode        *string         `db:"project_code" json:"project_code,omitempty"`
	CostCenter         *string         `db:"cost_center" json:"cost_center,omitempty"`
	TaxCode            *string         `db:"tax_code" json:"tax_code,omitempty"`
	TaxAmount          string          `db:"tax_amount" json:"tax_amount"`
	Description        *string         `db:"description" json:"description,omitempty"`
	Memo               *string         `db:"memo" json:"memo,omitempty"`
	IsReconciled       bool            `db:"is_reconciled" json:"is_reconciled"`
	ReconciledAt       *time.Time      `db:"reconciled_at" json:"reconciled_at,omitempty"`
	Metadata           json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt          time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updated_at"`
	CreatedBy          *uuid.UUID      `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID      `db:"updated_by" json:"updated_by,omitempty"`
	DeletedAt          *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateJournalEntryLineRequest struct {
	AccountID     uuid.UUID              `json:"account_id" validate:"required"`
	LineNumber    int                    `json:"line_number" validate:"required,min=1"`
	DebitAmount   *string                `json:"debit_amount"`
	CreditAmount  *string                `json:"credit_amount"`
	LocationID    *uuid.UUID             `json:"location_id"`
	Department    *string                `json:"department"`
	ProjectCode   *string                `json:"project_code"`
	CostCenter    *string                `json:"cost_center"`
	TaxCode       *string                `json:"tax_code"`
	TaxAmount     *string                `json:"tax_amount"`
	Description   *string                `json:"description"`
	Memo          *string                `json:"memo"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type CreateJournalEntryRequest struct {
	EntryTypeID          uuid.UUID                       `json:"entry_type_id" validate:"required"`
	EntryDate            time.Time                       `json:"entry_date" validate:"required"`
	PostingDate          time.Time                       `json:"posting_date" validate:"required"`
	AccountingPeriodID   *uuid.UUID                      `json:"accounting_period_id"`
	FiscalYearID         *uuid.UUID                      `json:"fiscal_year_id"`
	SourceModule         *string                         `json:"source_module"`
	SourceDocumentType   *string                         `json:"source_document_type"`
	SourceDocumentID     *uuid.UUID                      `json:"source_document_id"`
	ReferenceNumber      *string                         `json:"reference_number"`
	Description          string                          `json:"description" validate:"required"`
	Notes                *string                         `json:"notes"`
	RequiresApproval     bool                            `json:"requires_approval"`
	Lines                []CreateJournalEntryLineRequest `json:"lines" validate:"required,min=2"`
	Attachments          []Attachment                    `json:"attachments"`
	Metadata             map[string]interface{}          `json:"metadata"`
}

type UpdateJournalEntryRequest struct {
	EntryDate       *time.Time                      `json:"entry_date"`
	PostingDate     *time.Time                      `json:"posting_date"`
	Description     *string                         `json:"description"`
	Notes           *string                         `json:"notes"`
	RequiresApproval *bool                          `json:"requires_approval"`
	Lines           []CreateJournalEntryLineRequest `json:"lines"`
	Metadata        map[string]interface{}          `json:"metadata"`
}

// ========================
// GENERAL LEDGER
// ========================

type GeneralLedger struct {
	ID                 uuid.UUID       `db:"id" json:"id"`
	OrganizationID     uuid.UUID       `db:"organization_id" json:"organization_id"`
	JournalEntryID     uuid.UUID       `db:"journal_entry_id" json:"journal_entry_id"`
	JournalEntryLineID uuid.UUID       `db:"journal_entry_line_id" json:"journal_entry_line_id"`
	AccountID          uuid.UUID       `db:"account_id" json:"account_id"`
	TransactionDate    time.Time       `db:"transaction_date" json:"transaction_date"`
	PostingDate        time.Time       `db:"posting_date" json:"posting_date"`
	AccountingPeriodID *uuid.UUID      `db:"accounting_period_id" json:"accounting_period_id,omitempty"`
	FiscalYearID       *uuid.UUID      `db:"fiscal_year_id" json:"fiscal_year_id,omitempty"`
	DebitAmount        string          `db:"debit_amount" json:"debit_amount"`
	CreditAmount       string          `db:"credit_amount" json:"credit_amount"`
	RunningDebitBalance string          `db:"running_debit_balance" json:"running_debit_balance"`
	RunningCreditBalance string         `db:"running_credit_balance" json:"running_credit_balance"`
	RunningBalance     string          `db:"running_balance" json:"running_balance"`
	SourceModule       *string         `db:"source_module" json:"source_module,omitempty"`
	SourceDocumentType *string         `db:"source_document_type" json:"source_document_type,omitempty"`
	SourceDocumentID   *uuid.UUID      `db:"source_document_id" json:"source_document_id,omitempty"`
	ReferenceNumber    *string         `db:"reference_number" json:"reference_number,omitempty"`
	LocationID         *uuid.UUID      `db:"location_id" json:"location_id,omitempty"`
	Department         *string         `db:"department" json:"department,omitempty"`
	ProjectCode        *string         `db:"project_code" json:"project_code,omitempty"`
	CostCenter         *string         `db:"cost_center" json:"cost_center,omitempty"`
	Description        *string         `db:"description" json:"description,omitempty"`
	IsReversed         bool            `db:"is_reversed" json:"is_reversed"`
	ReversalGLID       *uuid.UUID      `db:"reversal_gl_id" json:"reversal_gl_id,omitempty"`
	Metadata           json.RawMessage `db:"metadata" json:"metadata,omitempty"`
	CreatedAt          time.Time       `db:"created_at" json:"created_at"`
	CreatedBy          *uuid.UUID      `db:"created_by" json:"created_by,omitempty"`
}

// ========================
// VALIDATION ERRORS
// ========================

var (
	ErrInvalidDateRange          = errors.New("end date must be after start date")
	ErrJournalEntryNotBalanced   = errors.New("journal entry debit and credit amounts must be equal")
	ErrMinimumTwoLines           = errors.New("journal entry must have at least 2 lines")
	ErrDebitCreditBoth           = errors.New("line must have either debit or credit amount, not both")
	ErrHeaderAccountPosting      = errors.New("cannot post to header accounts")
	ErrInactiveAccount           = errors.New("cannot post to inactive accounts")
	ErrPeriodClosed              = errors.New("cannot post to closed periods")
	ErrFiscalYearClosed          = errors.New("cannot post to closed fiscal years")
	ErrInvalidAccountLevel       = errors.New("invalid account level")
	ErrCyclicAccountHierarchy    = errors.New("cannot create cyclic account hierarchy")
	ErrSystemAccountModification = errors.New("cannot modify system accounts")
)
