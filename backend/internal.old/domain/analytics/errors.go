package analytics

import "errors"

var (
	// Analytic Plan errors
	ErrAnalyticPlanCodeRequired  = errors.New("analytic plan code is required")
	ErrAnalyticPlanNameRequired  = errors.New("analytic plan name is required")
	ErrAnalyticPlanNotFound      = errors.New("analytic plan not found")
	ErrDuplicateAnalyticPlanCode = errors.New("analytic plan code already exists")

	// Analytic Account errors
	ErrAnalyticAccountCodeRequired  = errors.New("analytic account code is required")
	ErrAnalyticAccountNameRequired  = errors.New("analytic account name is required")
	ErrAnalyticAccountNotFound      = errors.New("analytic account not found")
	ErrDuplicateAnalyticAccountCode = errors.New("analytic account code already exists")
	ErrCircularAnalyticReference    = errors.New("circular reference detected in analytic account hierarchy")

	// Deferred Revenue errors
	ErrDeferredRevenueNotFound      = errors.New("deferred revenue contract not found")
	ErrInvalidDeferralDates         = errors.New("end date must be after start date")
	ErrDeferralAlreadyCompleted     = errors.New("deferral contract is already completed")
	ErrDeferralAmountExceeded       = errors.New("total recognized amount exceeds contract amount")
	ErrInvalidDeferralAmount        = errors.New("deferral amount must be positive")

	// Deferred Expense errors
	ErrDeferredExpenseNotFound      = errors.New("deferred expense contract not found")

	// Deferral Schedule errors
	ErrDeferralScheduleNotFound     = errors.New("deferral schedule entry not found")
	ErrDeferralScheduleAlreadyPosted = errors.New("deferral schedule already posted")
	ErrInvalidScheduleLineNumber    = errors.New("invalid schedule line number")

	// Budget errors
	ErrBudgetNotFound              = errors.New("budget not found")
	ErrBudgetCodeRequired          = errors.New("budget code is required")
	ErrBudgetNameRequired          = errors.New("budget name is required")
	ErrDuplicateBudgetCode         = errors.New("budget code already exists")
	ErrBudgetAlreadyApproved       = errors.New("budget is already approved")
	ErrBudgetAlreadyActive         = errors.New("budget is already active")
	ErrBudgetCannotBeModified      = errors.New("budget cannot be modified in current status")

	// Budget Line errors
	ErrBudgetLineNotFound          = errors.New("budget line not found")
	ErrBudgetLineNeedsDimension    = errors.New("budget line must have account or analytic account")

	// Localization errors
	ErrLocalizationPackageNotFound = errors.New("localization package not found")

	// General errors
	ErrOrganizationIDRequired      = errors.New("organization_id is required")
	ErrInvalidRecognitionMethod    = errors.New("invalid recognition method")
	ErrInvalidBudgetType           = errors.New("invalid budget type")
)
