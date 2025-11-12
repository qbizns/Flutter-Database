package analytics

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RecognitionMethod defines how deferred items are recognized
type RecognitionMethod string

const (
	RecognitionMethodStraightLine RecognitionMethod = "straight_line"
	RecognitionMethodCustom       RecognitionMethod = "custom"
	RecognitionMethodMilestone    RecognitionMethod = "milestone"
	RecognitionMethodUsageBased   RecognitionMethod = "usage_based"
)

// DeferralStatus represents the state of a deferral contract
type DeferralStatus string

const (
	DeferralStatusDraft     DeferralStatus = "draft"
	DeferralStatusActive    DeferralStatus = "active"
	DeferralStatusCompleted DeferralStatus = "completed"
	DeferralStatusCanceled  DeferralStatus = "canceled"
)

// ScheduleStatus represents the state of a schedule entry
type ScheduleStatus string

const (
	ScheduleStatusPending  ScheduleStatus = "pending"
	ScheduleStatusPosted   ScheduleStatus = "posted"
	ScheduleStatusCanceled ScheduleStatus = "canceled"
)

// BudgetType represents the type of budget
type BudgetType string

const (
	BudgetTypeOperating    BudgetType = "operating"
	BudgetTypeCapital      BudgetType = "capital"
	BudgetTypeCashFlow     BudgetType = "cash_flow"
	BudgetTypeProject      BudgetType = "project"
	BudgetTypeDepartmental BudgetType = "departmental"
)

// BudgetStatus represents the state of a budget
type BudgetStatus string

const (
	BudgetStatusDraft    BudgetStatus = "draft"
	BudgetStatusApproved BudgetStatus = "approved"
	BudgetStatusActive   BudgetStatus = "active"
	BudgetStatusClosed   BudgetStatus = "closed"
)

// AnalyticPlan represents a dimension for cost accounting
type AnalyticPlan struct {
	ID          uuid.UUID
	OrganizationID uuid.UUID
	PlanCode    string
	PlanName    string
	IsActive    bool
	Description *string
	CreatedBy   *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// CreateAnalyticPlanRequest contains input for creating an analytic plan
type CreateAnalyticPlanRequest struct {
	PlanCode    string  `json:"plan_code" validate:"required,max=20"`
	PlanName    string  `json:"plan_name" validate:"required,max=255"`
	Description *string `json:"description"`
}

// UpdateAnalyticPlanRequest contains input for updating an analytic plan
type UpdateAnalyticPlanRequest struct {
	PlanName    *string `json:"plan_name"`
	IsActive    *bool   `json:"is_active"`
	Description *string `json:"description"`
}

// AnalyticAccount represents a cost center or department
type AnalyticAccount struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	AnalyticPlanID *uuid.UUID
	AccountCode    string
	AccountName    string
	ParentAccountID *uuid.UUID
	AccountLevel   int
	IsActive       bool
	Description    *string
	CreatedBy      *uuid.UUID
	UpdatedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// CreateAnalyticAccountRequest contains input for creating an analytic account
type CreateAnalyticAccountRequest struct {
	AnalyticPlanID  *uuid.UUID `json:"analytic_plan_id"`
	AccountCode     string     `json:"account_code" validate:"required,max=50"`
	AccountName     string     `json:"account_name" validate:"required,max=255"`
	ParentAccountID *uuid.UUID `json:"parent_account_id"`
	Description     *string    `json:"description"`
}

// UpdateAnalyticAccountRequest contains input for updating an analytic account
type UpdateAnalyticAccountRequest struct {
	AccountName     *string    `json:"account_name"`
	ParentAccountID *uuid.UUID `json:"parent_account_id"`
	IsActive        *bool      `json:"is_active"`
	Description     *string    `json:"description"`
}

// DeferredRevenueContract represents revenue received upfront but recognized over time
type DeferredRevenueContract struct {
	ID                  uuid.UUID
	OrganizationID      uuid.UUID
	CustomerInvoiceID   *uuid.UUID
	InvoiceLineID       *uuid.UUID
	ContractName        string
	TotalDeferredAmount float64
	StartDate           time.Time
	EndDate             time.Time
	RecognitionMethod   RecognitionMethod
	DeferredAccountID   uuid.UUID
	RevenueAccountID    uuid.UUID
	Status              DeferralStatus
	RecognizedAmount    float64
	Notes               *string
	CreatedBy           *uuid.UUID
	UpdatedBy           *uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

// CreateDeferredRevenueRequest contains input for creating a deferred revenue contract
type CreateDeferredRevenueRequest struct {
	CustomerInvoiceID   *uuid.UUID         `json:"customer_invoice_id"`
	InvoiceLineID       *uuid.UUID         `json:"invoice_line_id"`
	ContractName        string             `json:"contract_name" validate:"required,max=255"`
	TotalDeferredAmount float64            `json:"total_deferred_amount" validate:"required,gt=0"`
	StartDate           time.Time          `json:"start_date" validate:"required"`
	EndDate             time.Time          `json:"end_date" validate:"required"`
	RecognitionMethod   RecognitionMethod  `json:"recognition_method" validate:"required"`
	DeferredAccountID   uuid.UUID          `json:"deferred_account_id" validate:"required"`
	RevenueAccountID    uuid.UUID          `json:"revenue_account_id" validate:"required"`
	Notes               *string            `json:"notes"`
}

// UpdateDeferredRevenueRequest contains input for updating a deferred revenue contract
type UpdateDeferredRevenueRequest struct {
	ContractName        *string            `json:"contract_name"`
	EndDate             *time.Time         `json:"end_date"`
	RecognitionMethod   *RecognitionMethod `json:"recognition_method"`
	Status              *DeferralStatus    `json:"status"`
	Notes               *string            `json:"notes"`
}

// DeferredRevenueSchedule tracks revenue recognition over time
type DeferredRevenueSchedule struct {
	ID                 uuid.UUID
	ContractID         uuid.UUID
	LineNumber         int
	RecognitionDate    time.Time
	RecognitionAmount  float64
	Status             ScheduleStatus
	JournalEntryID     *uuid.UUID
	CreatedAt          time.Time
	PostedAt           *time.Time
	DeletedAt          *time.Time
}

// CreateDeferredRevenueScheduleRequest contains input for creating a schedule entry
type CreateDeferredRevenueScheduleRequest struct {
	ContractID         uuid.UUID `json:"contract_id" validate:"required"`
	LineNumber         int       `json:"line_number" validate:"required,gt=0"`
	RecognitionDate    time.Time `json:"recognition_date" validate:"required"`
	RecognitionAmount  float64   `json:"recognition_amount" validate:"required,gt=0"`
}

// UpdateDeferredRevenueScheduleRequest contains input for updating a schedule entry
type UpdateDeferredRevenueScheduleRequest struct {
	RecognitionAmount *float64 `json:"recognition_amount"`
	Status            *ScheduleStatus `json:"status"`
}

// DeferredExpenseContract represents expenses paid upfront but recognized over time
type DeferredExpenseContract struct {
	ID                  uuid.UUID
	OrganizationID      uuid.UUID
	VendorBillID        *uuid.UUID
	BillLineID          *uuid.UUID
	ContractName        string
	TotalDeferredAmount float64
	StartDate           time.Time
	EndDate             time.Time
	RecognitionMethod   RecognitionMethod
	DeferredAccountID   uuid.UUID
	ExpenseAccountID    uuid.UUID
	Status              DeferralStatus
	RecognizedAmount    float64
	Notes               *string
	CreatedBy           *uuid.UUID
	UpdatedBy           *uuid.UUID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

// CreateDeferredExpenseRequest contains input for creating a deferred expense contract
type CreateDeferredExpenseRequest struct {
	VendorBillID        *uuid.UUID         `json:"vendor_bill_id"`
	BillLineID          *uuid.UUID         `json:"bill_line_id"`
	ContractName        string             `json:"contract_name" validate:"required,max=255"`
	TotalDeferredAmount float64            `json:"total_deferred_amount" validate:"required,gt=0"`
	StartDate           time.Time          `json:"start_date" validate:"required"`
	EndDate             time.Time          `json:"end_date" validate:"required"`
	RecognitionMethod   RecognitionMethod  `json:"recognition_method" validate:"required"`
	DeferredAccountID   uuid.UUID          `json:"deferred_account_id" validate:"required"`
	ExpenseAccountID    uuid.UUID          `json:"expense_account_id" validate:"required"`
	Notes               *string            `json:"notes"`
}

// UpdateDeferredExpenseRequest contains input for updating a deferred expense contract
type UpdateDeferredExpenseRequest struct {
	ContractName        *string            `json:"contract_name"`
	EndDate             *time.Time         `json:"end_date"`
	RecognitionMethod   *RecognitionMethod `json:"recognition_method"`
	Status              *DeferralStatus    `json:"status"`
	Notes               *string            `json:"notes"`
}

// DeferredExpenseSchedule tracks expense recognition over time
type DeferredExpenseSchedule struct {
	ID                 uuid.UUID
	ContractID         uuid.UUID
	LineNumber         int
	RecognitionDate    time.Time
	RecognitionAmount  float64
	Status             ScheduleStatus
	JournalEntryID     *uuid.UUID
	CreatedAt          time.Time
	PostedAt           *time.Time
	DeletedAt          *time.Time
}

// CreateDeferredExpenseScheduleRequest contains input for creating a schedule entry
type CreateDeferredExpenseScheduleRequest struct {
	ContractID         uuid.UUID `json:"contract_id" validate:"required"`
	LineNumber         int       `json:"line_number" validate:"required,gt=0"`
	RecognitionDate    time.Time `json:"recognition_date" validate:"required"`
	RecognitionAmount  float64   `json:"recognition_amount" validate:"required,gt=0"`
}

// UpdateDeferredExpenseScheduleRequest contains input for updating a schedule entry
type UpdateDeferredExpenseScheduleRequest struct {
	RecognitionAmount *float64 `json:"recognition_amount"`
	Status            *ScheduleStatus `json:"status"`
}

// Budget represents a budget scenario for planning
type Budget struct {
	ID          uuid.UUID
	OrganizationID uuid.UUID
	BudgetCode  string
	BudgetName  string
	FiscalYearID *uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
	BudgetType  BudgetType
	Status      BudgetStatus
	Notes       *string
	CreatedBy   *uuid.UUID
	UpdatedBy   *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// CreateBudgetRequest contains input for creating a budget
type CreateBudgetRequest struct {
	BudgetCode   string     `json:"budget_code" validate:"required,max=50"`
	BudgetName   string     `json:"budget_name" validate:"required,max=255"`
	FiscalYearID *uuid.UUID `json:"fiscal_year_id"`
	StartDate    time.Time  `json:"start_date" validate:"required"`
	EndDate      time.Time  `json:"end_date" validate:"required"`
	BudgetType   BudgetType `json:"budget_type" validate:"required"`
	Notes        *string    `json:"notes"`
}

// UpdateBudgetRequest contains input for updating a budget
type UpdateBudgetRequest struct {
	BudgetName *string     `json:"budget_name"`
	Status     *BudgetStatus `json:"status"`
	Notes      *string     `json:"notes"`
}

// BudgetLine represents a budget allocation for an account or dimension
type BudgetLine struct {
	ID                  uuid.UUID
	BudgetID            uuid.UUID
	AccountID           *uuid.UUID
	AnalyticAccountID   *uuid.UUID
	AccountingPeriodID  *uuid.UUID
	PeriodStartDate     *time.Time
	PeriodEndDate       *time.Time
	PlannedAmount       float64
	Notes               *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

// CreateBudgetLineRequest contains input for creating a budget line
type CreateBudgetLineRequest struct {
	BudgetID            uuid.UUID  `json:"budget_id" validate:"required"`
	AccountID           *uuid.UUID `json:"account_id"`
	AnalyticAccountID   *uuid.UUID `json:"analytic_account_id"`
	AccountingPeriodID  *uuid.UUID `json:"accounting_period_id"`
	PeriodStartDate     *time.Time `json:"period_start_date"`
	PeriodEndDate       *time.Time `json:"period_end_date"`
	PlannedAmount       float64    `json:"planned_amount" validate:"required"`
	Notes               *string    `json:"notes"`
}

// UpdateBudgetLineRequest contains input for updating a budget line
type UpdateBudgetLineRequest struct {
	PlannedAmount *float64 `json:"planned_amount"`
	Notes         *string  `json:"notes"`
}

// LocalizationPackage contains country-specific accounting rules
type LocalizationPackage struct {
	ID          uuid.UUID
	PackageCode string
	PackageName string
	CountryCode *string
	Region      *string
	Description *string
	Version     *string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// AnalyticQuery filters for fetching analytic data
type AnalyticQuery struct {
	OrganizationID    uuid.UUID
	AnalyticPlanID    *uuid.UUID
	IsActive          *bool
	ParentAccountID   *uuid.UUID
	Limit             int
	Offset            int
}

// DeferralQuery filters for fetching deferral data
type DeferralQuery struct {
	OrganizationID   uuid.UUID
	Status           *DeferralStatus
	DateFrom         *time.Time
	DateTo           *time.Time
	Limit            int
	Offset           int
}

// BudgetQuery filters for fetching budget data
type BudgetQuery struct {
	OrganizationID  uuid.UUID
	FiscalYearID    *uuid.UUID
	BudgetType      *BudgetType
	Status          *BudgetStatus
	Limit           int
	Offset          int
}

// BudgetComparison contains actual vs budgeted amounts
type BudgetComparison struct {
	BudgetID       uuid.UUID
	AccountID      *uuid.UUID
	AnalyticAccountID *uuid.UUID
	PlannedAmount  float64
	ActualAmount   float64
	Variance       float64
	VariancePercent float64
}

// DeferralScheduleItem represents a single recognition entry
type DeferralScheduleItem struct {
	ID               uuid.UUID
	ContractID       uuid.UUID
	LineNumber       int
	RecognitionDate  time.Time
	RecognitionAmount float64
	Status           ScheduleStatus
	IsPosted         bool
}

// Scan implements sql.Scanner interface
func (rm *RecognitionMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*rm = RecognitionMethod(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (rm RecognitionMethod) Value() (driver.Value, error) {
	return string(rm), nil
}

// Scan implements sql.Scanner interface
func (ds *DeferralStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*ds = DeferralStatus(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (ds DeferralStatus) Value() (driver.Value, error) {
	return string(ds), nil
}

// Scan implements sql.Scanner interface
func (ss *ScheduleStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*ss = ScheduleStatus(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (ss ScheduleStatus) Value() (driver.Value, error) {
	return string(ss), nil
}

// Scan implements sql.Scanner interface
func (bt *BudgetType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*bt = BudgetType(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (bt BudgetType) Value() (driver.Value, error) {
	return string(bt), nil
}

// Scan implements sql.Scanner interface
func (bs *BudgetStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	*bs = BudgetStatus(value.(string))
	return nil
}

// Value implements driver.Valuer interface
func (bs BudgetStatus) Value() (driver.Value, error) {
	return string(bs), nil
}
