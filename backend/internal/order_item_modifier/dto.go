package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderItemModifiersResponse represents a order_item_modifiers response
type OrderItemModifiersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	OrderItemId uuid.UUID `json:"order_item_id"`
	
	ModifierId uuid.UUID `json:"modifier_id"`
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id"`
	
	ModifierName string `json:"modifier_name"`
	
	Quantity *int64 `json:"quantity"`
	
	PriceAdjustment *float64 `json:"price_adjustment"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateOrderItemModifiersRequest represents a request to create a order_item_modifiers
type CreateOrderItemModifiersRequest struct {
	
	OrderItemId uuid.UUID `json:"order_item_id" validate:"required"`
	
	ModifierId uuid.UUID `json:"modifier_id" validate:"required"`
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id"`
	
	ModifierName string `json:"modifier_name" validate:"required"`
	
	Quantity *int64 `json:"quantity"`
	
	PriceAdjustment *float64 `json:"price_adjustment"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateOrderItemModifiersRequest) Validate() error {
	
	if r.OrderItemId == uuid.Nil {
		return fmt.Errorf("order_item_id is required")
	}
	
	if r.ModifierId == uuid.Nil {
		return fmt.Errorf("modifier_id is required")
	}
	
	if r.ModifierName == "" {
		return fmt.Errorf("modifier_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateOrderItemModifiersRequest represents a request to update a order_item_modifiers
type UpdateOrderItemModifiersRequest struct {
	
	OrderItemId *uuid.UUID `json:"order_item_id,omitempty" validate:"omitempty,required"`
	
	ModifierId *uuid.UUID `json:"modifier_id,omitempty" validate:"omitempty,required"`
	
	ModifierGroupId *uuid.UUID `json:"modifier_group_id,omitempty"`
	
	ModifierName *string `json:"modifier_name,omitempty" validate:"omitempty,required"`
	
	Quantity *int64 `json:"quantity,omitempty"`
	
	PriceAdjustment *float64 `json:"price_adjustment,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateOrderItemModifiersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.OrderItemId != nil {
		hasUpdate = true
	}
	
	if r.ModifierId != nil {
		hasUpdate = true
	}
	
	if r.ModifierGroupId != nil {
		hasUpdate = true
	}
	
	if r.ModifierName != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.PriceAdjustment != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
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

// OrderItemModifiersListResponse represents a paginated list of order_item_modifiers records
type OrderItemModifiersListResponse struct {
	Items      []*OrderItemModifiersResponse `json:"items"`
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
