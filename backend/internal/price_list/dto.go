package price_list

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PriceListsResponse represents a price_lists response
type PriceListsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PriceListCode string `json:"price_list_code"`
	
	PriceListName string `json:"price_list_name"`
	
	PriceListType *string `json:"price_list_type"`
	
	PriceListType *string `json:"price_list_type"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	BasePriceAdjustmentType *string `json:"base_price_adjustment_type"`
	
	BasePriceAdjustmentValue *float64 `json:"base_price_adjustment_value"`
	
	Priority *int64 `json:"priority"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePriceListsRequest represents a request to create a price_lists
type CreatePriceListsRequest struct {
	
	PriceListCode string `json:"price_list_code" validate:"required"`
	
	PriceListName string `json:"price_list_name" validate:"required"`
	
	PriceListType *string `json:"price_list_type"`
	
	PriceListType *string `json:"price_list_type"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	BasePriceAdjustmentType *string `json:"base_price_adjustment_type"`
	
	BasePriceAdjustmentValue *float64 `json:"base_price_adjustment_value"`
	
	Priority *int64 `json:"priority"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePriceListsRequest) Validate() error {
	
	if r.PriceListCode == "" {
		return fmt.Errorf("price_list_code is required")
	}
	
	if r.PriceListName == "" {
		return fmt.Errorf("price_list_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePriceListsRequest represents a request to update a price_lists
type UpdatePriceListsRequest struct {
	
	PriceListCode *string `json:"price_list_code,omitempty" validate:"omitempty,required"`
	
	PriceListName *string `json:"price_list_name,omitempty" validate:"omitempty,required"`
	
	PriceListType *string `json:"price_list_type,omitempty"`
	
	PriceListType *string `json:"price_list_type,omitempty"`
	
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	
	BasePriceAdjustmentType *string `json:"base_price_adjustment_type,omitempty"`
	
	BasePriceAdjustmentValue *float64 `json:"base_price_adjustment_value,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePriceListsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PriceListCode != nil {
		hasUpdate = true
	}
	
	if r.PriceListName != nil {
		hasUpdate = true
	}
	
	if r.PriceListType != nil {
		hasUpdate = true
	}
	
	if r.PriceListType != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	
	if r.EffectiveTo != nil {
		hasUpdate = true
	}
	
	if r.BasePriceAdjustmentType != nil {
		hasUpdate = true
	}
	
	if r.BasePriceAdjustmentValue != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
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

// PriceListsListResponse represents a paginated list of price_lists records
type PriceListsListResponse struct {
	Items      []*PriceListsResponse `json:"items"`
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
