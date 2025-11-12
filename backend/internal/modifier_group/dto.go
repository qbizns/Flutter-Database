package modifier_group

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ModifierGroupsResponse represents a modifier_groups response
type ModifierGroupsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	GroupName string `json:"group_name"`
	
	GroupCode *string `json:"group_code"`
	
	DisplayName *string `json:"display_name"`
	
	SelectionType *string `json:"selection_type"`
	
	MinSelections *int64 `json:"min_selections"`
	
	MaxSelections *int64 `json:"max_selections"`
	
	ExactSelections *int64 `json:"exact_selections"`
	
	IsRequired *bool `json:"is_required"`
	
	AffectsPrice *bool `json:"affects_price"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	MinSelections *string `json:"min_selections"`
	
	(maxSelections *string `json:"(max_selections"`
	
	(exactSelections *string `json:"(exact_selections"`
	
}

// CreateModifierGroupsRequest represents a request to create a modifier_groups
type CreateModifierGroupsRequest struct {
	
	GroupName string `json:"group_name" validate:"required"`
	
	GroupCode *string `json:"group_code"`
	
	DisplayName *string `json:"display_name"`
	
	SelectionType *string `json:"selection_type"`
	
	MinSelections *int64 `json:"min_selections"`
	
	MaxSelections *int64 `json:"max_selections"`
	
	ExactSelections *int64 `json:"exact_selections"`
	
	IsRequired *bool `json:"is_required"`
	
	AffectsPrice *bool `json:"affects_price"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	MinSelections *string `json:"min_selections"`
	
	(maxSelections *string `json:"(max_selections"`
	
	(exactSelections *string `json:"(exact_selections"`
	
}

// Validate validates the create request
func (r *CreateModifierGroupsRequest) Validate() error {
	
	if r.GroupName == "" {
		return fmt.Errorf("group_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateModifierGroupsRequest represents a request to update a modifier_groups
type UpdateModifierGroupsRequest struct {
	
	GroupName *string `json:"group_name,omitempty" validate:"omitempty,required"`
	
	GroupCode *string `json:"group_code,omitempty"`
	
	DisplayName *string `json:"display_name,omitempty"`
	
	SelectionType *string `json:"selection_type,omitempty"`
	
	MinSelections *int64 `json:"min_selections,omitempty"`
	
	MaxSelections *int64 `json:"max_selections,omitempty"`
	
	ExactSelections *int64 `json:"exact_selections,omitempty"`
	
	IsRequired *bool `json:"is_required,omitempty"`
	
	AffectsPrice *bool `json:"affects_price,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	MinSelections *string `json:"min_selections,omitempty"`
	
	(maxSelections *string `json:"(max_selections,omitempty"`
	
	(exactSelections *string `json:"(exact_selections,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateModifierGroupsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.GroupName != nil {
		hasUpdate = true
	}
	
	if r.GroupCode != nil {
		hasUpdate = true
	}
	
	if r.DisplayName != nil {
		hasUpdate = true
	}
	
	if r.SelectionType != nil {
		hasUpdate = true
	}
	
	if r.MinSelections != nil {
		hasUpdate = true
	}
	
	if r.MaxSelections != nil {
		hasUpdate = true
	}
	
	if r.ExactSelections != nil {
		hasUpdate = true
	}
	
	if r.IsRequired != nil {
		hasUpdate = true
	}
	
	if r.AffectsPrice != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	
	if r.MinSelections != nil {
		hasUpdate = true
	}
	
	if r.(maxSelections != nil {
		hasUpdate = true
	}
	
	if r.(exactSelections != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ModifierGroupsListResponse represents a paginated list of modifier_groups records
type ModifierGroupsListResponse struct {
	Items      []*ModifierGroupsResponse `json:"items"`
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
