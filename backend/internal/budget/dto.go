package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BudgetsResponse represents a budgets response
type BudgetsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BudgetCode string `json:"budget_code"`
	
	BudgetName string `json:"budget_name"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate time.Time `json:"end_date"`
	
	BudgetType *string `json:"budget_type"`
	
	BudgetType *string `json:"budget_type"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBudgetsRequest represents a request to create a budgets
type CreateBudgetsRequest struct {
	
	BudgetCode string `json:"budget_code" validate:"required"`
	
	BudgetName string `json:"budget_name" validate:"required"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate time.Time `json:"end_date" validate:"required"`
	
	BudgetType *string `json:"budget_type"`
	
	BudgetType *string `json:"budget_type"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateBudgetsRequest) Validate() error {
	
	if r.BudgetCode == "" {
		return fmt.Errorf("budget_code is required")
	}
	
	if r.BudgetName == "" {
		return fmt.Errorf("budget_name is required")
	}
	
	if r.StartDate == nil {
		return fmt.Errorf("start_date is required")
	}
	
	if r.EndDate == nil {
		return fmt.Errorf("end_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBudgetsRequest represents a request to update a budgets
type UpdateBudgetsRequest struct {
	
	BudgetCode *string `json:"budget_code,omitempty" validate:"omitempty,required"`
	
	BudgetName *string `json:"budget_name,omitempty" validate:"omitempty,required"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id,omitempty"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty" validate:"omitempty,required"`
	
	BudgetType *string `json:"budget_type,omitempty"`
	
	BudgetType *string `json:"budget_type,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBudgetsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BudgetCode != nil {
		hasUpdate = true
	}
	
	if r.BudgetName != nil {
		hasUpdate = true
	}
	
	if r.FiscalYearId != nil {
		hasUpdate = true
	}
	
	if r.StartDate != nil {
		hasUpdate = true
	}
	
	if r.EndDate != nil {
		hasUpdate = true
	}
	
	if r.BudgetType != nil {
		hasUpdate = true
	}
	
	if r.BudgetType != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BudgetsListResponse represents a paginated list of budgets records
type BudgetsListResponse struct {
	Items      []*BudgetsResponse `json:"items"`
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
