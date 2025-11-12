package product_component

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductComponentsResponse represents a product_components response
type ProductComponentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ParentProductId uuid.UUID `json:"parent_product_id"`
	
	ComponentProductId *uuid.UUID `json:"component_product_id"`
	
	ComponentVariantId *uuid.UUID `json:"component_variant_id"`
	
	Quantity float64 `json:"quantity"`
	
	InheritPrice *bool `json:"inherit_price"`
	
	PriceOverride *float64 `json:"price_override"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsOptional *bool `json:"is_optional"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	ComponentProductId string `json:"component_product_id"`
	
}

// CreateProductComponentsRequest represents a request to create a product_components
type CreateProductComponentsRequest struct {
	
	ParentProductId uuid.UUID `json:"parent_product_id" validate:"required"`
	
	ComponentProductId *uuid.UUID `json:"component_product_id"`
	
	ComponentVariantId *uuid.UUID `json:"component_variant_id"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	InheritPrice *bool `json:"inherit_price"`
	
	PriceOverride *float64 `json:"price_override"`
	
	DisplayOrder *int64 `json:"display_order"`
	
	IsOptional *bool `json:"is_optional"`
	
	ComponentProductId string `json:"component_product_id" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateProductComponentsRequest) Validate() error {
	
	if r.ParentProductId == uuid.Nil {
		return fmt.Errorf("parent_product_id is required")
	}
	
	if r.Quantity == nil {
		return fmt.Errorf("quantity is required")
	}
	
	if r.ComponentProductId == "" {
		return fmt.Errorf("component_product_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductComponentsRequest represents a request to update a product_components
type UpdateProductComponentsRequest struct {
	
	ParentProductId *uuid.UUID `json:"parent_product_id,omitempty" validate:"omitempty,required"`
	
	ComponentProductId *uuid.UUID `json:"component_product_id,omitempty"`
	
	ComponentVariantId *uuid.UUID `json:"component_variant_id,omitempty"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	InheritPrice *bool `json:"inherit_price,omitempty"`
	
	PriceOverride *float64 `json:"price_override,omitempty"`
	
	DisplayOrder *int64 `json:"display_order,omitempty"`
	
	IsOptional *bool `json:"is_optional,omitempty"`
	
	ComponentProductId *string `json:"component_product_id,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateProductComponentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ParentProductId != nil {
		hasUpdate = true
	}
	
	if r.ComponentProductId != nil {
		hasUpdate = true
	}
	
	if r.ComponentVariantId != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.InheritPrice != nil {
		hasUpdate = true
	}
	
	if r.PriceOverride != nil {
		hasUpdate = true
	}
	
	if r.DisplayOrder != nil {
		hasUpdate = true
	}
	
	if r.IsOptional != nil {
		hasUpdate = true
	}
	
	if r.ComponentProductId != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ProductComponentsListResponse represents a paginated list of product_components records
type ProductComponentsListResponse struct {
	Items      []*ProductComponentsResponse `json:"items"`
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
