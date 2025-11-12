package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ModifiersResponse represents a modifiers response
type ModifiersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id"`
	
	ModifierName string `json:"modifier_name"`
	
	ModifierCode *string `json:"modifier_code"`
	
	DisplayName *string `json:"display_name"`
	
	PriceAdjustment *float64 `json:"price_adjustment"`
	
	PriceType *string `json:"price_type"`
	
	IsAvailable *bool `json:"is_available"`
	
	IsDefault *bool `json:"is_default"`
	
	TrackInventory *bool `json:"track_inventory"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ImageUrl *string `json:"image_url"`
	
	Description *string `json:"description"`
	
	AllergenInfo *string `json:"allergen_info"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateModifiersRequest represents a request to create a modifiers
type CreateModifiersRequest struct {
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id"`
	
	ModifierName string `json:"modifier_name" validate:"required"`
	
	ModifierCode *string `json:"modifier_code"`
	
	DisplayName *string `json:"display_name"`
	
	PriceAdjustment *float64 `json:"price_adjustment"`
	
	PriceType *string `json:"price_type"`
	
	IsAvailable *bool `json:"is_available"`
	
	IsDefault *bool `json:"is_default"`
	
	TrackInventory *bool `json:"track_inventory"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	ImageUrl *string `json:"image_url" validate:"url"`
	
	Description *string `json:"description"`
	
	AllergenInfo *string `json:"allergen_info"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateModifiersRequest) Validate() error {
	
	if r.ModifierName == "" {
		return fmt.Errorf("modifier_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateModifiersRequest represents a request to update a modifiers
type UpdateModifiersRequest struct {
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id,omitempty"`
	
	ModifierName *string `json:"modifier_name,omitempty" validate:"omitempty,required"`
	
	ModifierCode *string `json:"modifier_code,omitempty"`
	
	DisplayName *string `json:"display_name,omitempty"`
	
	PriceAdjustment *float64 `json:"price_adjustment,omitempty"`
	
	PriceType *string `json:"price_type,omitempty"`
	
	IsAvailable *bool `json:"is_available,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	TrackInventory *bool `json:"track_inventory,omitempty"`
	
	CurrentStock *float64 `json:"current_stock,omitempty"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	ImageUrl *string `json:"image_url,omitempty" validate:"omitempty,url"`
	
	Description *string `json:"description,omitempty"`
	
	AllergenInfo *string `json:"allergen_info,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateModifiersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ModifierGroupId != nil {
		hasUpdate = true
	}
	
	if r.ModifierName != nil {
		hasUpdate = true
	}
	
	if r.ModifierCode != nil {
		hasUpdate = true
	}
	
	if r.DisplayName != nil {
		hasUpdate = true
	}
	
	if r.PriceAdjustment != nil {
		hasUpdate = true
	}
	
	if r.PriceType != nil {
		hasUpdate = true
	}
	
	if r.IsAvailable != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.TrackInventory != nil {
		hasUpdate = true
	}
	
	if r.CurrentStock != nil {
		hasUpdate = true
	}
	
	if r.LowStockThreshold != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.ImageUrl != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.AllergenInfo != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ModifiersListResponse represents a paginated list of modifiers records
type ModifiersListResponse struct {
	Items      []*ModifiersResponse `json:"items"`
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
