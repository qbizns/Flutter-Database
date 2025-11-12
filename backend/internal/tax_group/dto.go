package tax_group

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TaxGroupsResponse represents a tax_groups response
type TaxGroupsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	GroupCode string `json:"group_code"`
	
	GroupName string `json:"group_name"`
	
	Sequence *int64 `json:"sequence"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateTaxGroupsRequest represents a request to create a tax_groups
type CreateTaxGroupsRequest struct {
	
	GroupCode string `json:"group_code" validate:"required"`
	
	GroupName string `json:"group_name" validate:"required"`
	
	Sequence *int64 `json:"sequence"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateTaxGroupsRequest) Validate() error {
	
	if r.GroupCode == "" {
		return fmt.Errorf("group_code is required")
	}
	
	if r.GroupName == "" {
		return fmt.Errorf("group_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTaxGroupsRequest represents a request to update a tax_groups
type UpdateTaxGroupsRequest struct {
	
	GroupCode *string `json:"group_code,omitempty" validate:"omitempty,required"`
	
	GroupName *string `json:"group_name,omitempty" validate:"omitempty,required"`
	
	Sequence *int64 `json:"sequence,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTaxGroupsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.GroupCode != nil {
		hasUpdate = true
	}
	
	if r.GroupName != nil {
		hasUpdate = true
	}
	
	if r.Sequence != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// TaxGroupsListResponse represents a paginated list of tax_groups records
type TaxGroupsListResponse struct {
	Items      []*TaxGroupsResponse `json:"items"`
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
