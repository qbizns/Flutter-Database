package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReturnReasonsResponse represents a return_reasons response
type ReturnReasonsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	ReasonCode string `json:"reason_code"`
	
	ReasonName string `json:"reason_name"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	AffectsInventory *bool `json:"affects_inventory"`
	
	IsRestockable *bool `json:"is_restockable"`
	
	IsActive *bool `json:"is_active"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateReturnReasonsRequest represents a request to create a return_reasons
type CreateReturnReasonsRequest struct {
	
	ReasonCode string `json:"reason_code" validate:"required"`
	
	ReasonName string `json:"reason_name" validate:"required"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	AffectsInventory *bool `json:"affects_inventory"`
	
	IsRestockable *bool `json:"is_restockable"`
	
	IsActive *bool `json:"is_active"`
	
	DisplayOrder *int64 `json:"display_order"`
	
}

// Validate validates the create request
func (r *CreateReturnReasonsRequest) Validate() error {
	
	if r.ReasonCode == "" {
		return fmt.Errorf("reason_code is required")
	}
	
	if r.ReasonName == "" {
		return fmt.Errorf("reason_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateReturnReasonsRequest represents a request to update a return_reasons
type UpdateReturnReasonsRequest struct {
	
	ReasonCode *string `json:"reason_code,omitempty" validate:"omitempty,required"`
	
	ReasonName *string `json:"reason_name,omitempty" validate:"omitempty,required"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	AffectsInventory *bool `json:"affects_inventory,omitempty"`
	
	IsRestockable *bool `json:"is_restockable,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateReturnReasonsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ReasonCode != nil {
		hasUpdate = true
	}
	
	if r.ReasonName != nil {
		hasUpdate = true
	}
	
	if r.RequiresApproval != nil {
		hasUpdate = true
	}
	
	if r.AffectsInventory != nil {
		hasUpdate = true
	}
	
	if r.IsRestockable != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ReturnReasonsListResponse represents a paginated list of return_reasons records
type ReturnReasonsListResponse struct {
	Items      []*ReturnReasonsResponse `json:"items"`
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
