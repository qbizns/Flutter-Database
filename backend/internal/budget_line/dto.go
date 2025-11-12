package budget_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BudgetLinesResponse represents a budget_lines response
type BudgetLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	BudgetId uuid.UUID `json:"budget_id"`
	
	AccountId *uuid.UUID `json:"account_id"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	PeriodStartDate *time.Time `json:"period_start_date"`
	
	PeriodEndDate *time.Time `json:"period_end_date"`
	
	PlannedAmount float64 `json:"planned_amount"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	AccountId string `json:"account_id"`
	
}

// CreateBudgetLinesRequest represents a request to create a budget_lines
type CreateBudgetLinesRequest struct {
	
	BudgetId uuid.UUID `json:"budget_id" validate:"required"`
	
	AccountId *uuid.UUID `json:"account_id"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	PeriodStartDate *time.Time `json:"period_start_date"`
	
	PeriodEndDate *time.Time `json:"period_end_date"`
	
	PlannedAmount float64 `json:"planned_amount" validate:"required"`
	
	Notes *string `json:"notes"`
	
	AccountId string `json:"account_id" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateBudgetLinesRequest) Validate() error {
	
	if r.BudgetId == uuid.Nil {
		return fmt.Errorf("budget_id is required")
	}
	
	if r.PlannedAmount == nil {
		return fmt.Errorf("planned_amount is required")
	}
	
	if r.AccountId == "" {
		return fmt.Errorf("account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBudgetLinesRequest represents a request to update a budget_lines
type UpdateBudgetLinesRequest struct {
	
	BudgetId *uuid.UUID `json:"budget_id,omitempty" validate:"omitempty,required"`
	
	AccountId *uuid.UUID `json:"account_id,omitempty"`
	
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id,omitempty"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	PeriodStartDate *time.Time `json:"period_start_date,omitempty"`
	
	PeriodEndDate *time.Time `json:"period_end_date,omitempty"`
	
	PlannedAmount *float64 `json:"planned_amount,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	AccountId *string `json:"account_id,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateBudgetLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BudgetId != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	
	if r.AnalyticAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
		hasUpdate = true
	}
	
	if r.PeriodStartDate != nil {
		hasUpdate = true
	}
	
	if r.PeriodEndDate != nil {
		hasUpdate = true
	}
	
	if r.PlannedAmount != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BudgetLinesListResponse represents a paginated list of budget_lines records
type BudgetLinesListResponse struct {
	Items      []*BudgetLinesResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
