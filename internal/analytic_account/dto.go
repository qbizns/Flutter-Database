package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AnalyticAccountsResponse represents a analytic_accounts response
type AnalyticAccountsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	AnalyticPlanId *uuid.UUID `json:"analytic_plan_id"`
	
	AccountCode string `json:"account_code"`
	
	AccountName string `json:"account_name"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id"`
	
	AccountLevel *int64 `json:"account_level"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateAnalyticAccountsRequest represents a request to create a analytic_accounts
type CreateAnalyticAccountsRequest struct {
	
	AnalyticPlanId *uuid.UUID `json:"analytic_plan_id"`
	
	AccountCode string `json:"account_code" validate:"required"`
	
	AccountName string `json:"account_name" validate:"required"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id"`
	
	AccountLevel *int64 `json:"account_level"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateAnalyticAccountsRequest) Validate() error {
	
	if r.AccountCode == "" {
		return fmt.Errorf("account_code is required")
	}
	
	if r.AccountName == "" {
		return fmt.Errorf("account_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAnalyticAccountsRequest represents a request to update a analytic_accounts
type UpdateAnalyticAccountsRequest struct {
	
	AnalyticPlanId *uuid.UUID `json:"analytic_plan_id,omitempty"`
	
	AccountCode *string `json:"account_code,omitempty" validate:"omitempty,required"`
	
	AccountName *string `json:"account_name,omitempty" validate:"omitempty,required"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id,omitempty"`
	
	AccountLevel *int64 `json:"account_level,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAnalyticAccountsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.AnalyticPlanId != nil {
		hasUpdate = true
	}
	
	if r.AccountCode != nil {
		hasUpdate = true
	}
	
	if r.AccountName != nil {
		hasUpdate = true
	}
	
	if r.ParentAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccountLevel != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// AnalyticAccountsListResponse represents a paginated list of analytic_accounts records
type AnalyticAccountsListResponse struct {
	Items      []*AnalyticAccountsResponse `json:"items"`
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
