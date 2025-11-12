package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AnalyticPlansResponse represents a analytic_plans response
type AnalyticPlansResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PlanCode string `json:"plan_code"`
	
	PlanName string `json:"plan_name"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateAnalyticPlansRequest represents a request to create a analytic_plans
type CreateAnalyticPlansRequest struct {
	
	PlanCode string `json:"plan_code" validate:"required"`
	
	PlanName string `json:"plan_name" validate:"required"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateAnalyticPlansRequest) Validate() error {
	
	if r.PlanCode == "" {
		return fmt.Errorf("plan_code is required")
	}
	
	if r.PlanName == "" {
		return fmt.Errorf("plan_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAnalyticPlansRequest represents a request to update a analytic_plans
type UpdateAnalyticPlansRequest struct {
	
	PlanCode *string `json:"plan_code,omitempty" validate:"omitempty,required"`
	
	PlanName *string `json:"plan_name,omitempty" validate:"omitempty,required"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAnalyticPlansRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PlanCode != nil {
		hasUpdate = true
	}
	
	if r.PlanName != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// AnalyticPlansListResponse represents a paginated list of analytic_plans records
type AnalyticPlansListResponse struct {
	Items      []*AnalyticPlansResponse `json:"items"`
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
