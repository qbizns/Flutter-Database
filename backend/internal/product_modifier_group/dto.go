package product_modifier_group

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductModifierGroupsResponse represents a product_modifier_groups response
type ProductModifierGroupsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ModifierGroupId uuid.UUID `json:"modifier_group_id"`
	
	IsRequired *bool `json:"is_required"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	OverrideMinSelections *int64 `json:"override_min_selections"`
	
	OverrideMaxSelections *int64 `json:"override_max_selections"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateProductModifierGroupsRequest represents a request to create a product_modifier_groups
type CreateProductModifierGroupsRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ModifierGroupId uuid.UUID `json:"modifier_group_id" validate:"required"`
	
	IsRequired *bool `json:"is_required"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsActive *bool `json:"is_active"`
	
	OverrideMinSelections *int64 `json:"override_min_selections"`
	
	OverrideMaxSelections *int64 `json:"override_max_selections"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateProductModifierGroupsRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.ModifierGroupId == uuid.Nil {
		return fmt.Errorf("modifier_group_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductModifierGroupsRequest represents a request to update a product_modifier_groups
type UpdateProductModifierGroupsRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id,omitempty" validate:"omitempty,required"`
	
	IsRequired *bool `json:"is_required,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	OverrideMinSelections *int64 `json:"override_min_selections,omitempty"`
	
	OverrideMaxSelections *int64 `json:"override_max_selections,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateProductModifierGroupsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ModifierGroupId != nil {
		hasUpdate = true
	}
	
	if r.IsRequired != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.OverrideMinSelections != nil {
		hasUpdate = true
	}
	
	if r.OverrideMaxSelections != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ProductModifierGroupsListResponse represents a paginated list of product_modifier_groups records
type ProductModifierGroupsListResponse struct {
	Items      []*ProductModifierGroupsResponse `json:"items"`
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
