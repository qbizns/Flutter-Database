package price_list_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PriceListItemsResponse represents a price_list_items response
type PriceListItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PriceListId uuid.UUID `json:"price_list_id"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	OverridePrice *float64 `json:"override_price"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	MarkupPercentage *float64 `json:"markup_percentage"`
	
	MinPrice *float64 `json:"min_price"`
	
	MaxPrice *float64 `json:"max_price"`
	
	MinQuantity *float64 `json:"min_quantity"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	ProductId string `json:"product_id"`
	
}

// CreatePriceListItemsRequest represents a request to create a price_list_items
type CreatePriceListItemsRequest struct {
	
	PriceListId uuid.UUID `json:"price_list_id" validate:"required"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	OverridePrice *float64 `json:"override_price"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	MarkupPercentage *float64 `json:"markup_percentage"`
	
	MinPrice *float64 `json:"min_price"`
	
	MaxPrice *float64 `json:"max_price"`
	
	MinQuantity *float64 `json:"min_quantity"`
	
	ProductId string `json:"product_id" validate:"required"`
	
}

// Validate validates the create request
func (r *CreatePriceListItemsRequest) Validate() error {
	
	if r.PriceListId == uuid.Nil {
		return fmt.Errorf("price_list_id is required")
	}
	
	if r.ProductId == "" {
		return fmt.Errorf("product_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePriceListItemsRequest represents a request to update a price_list_items
type UpdatePriceListItemsRequest struct {
	
	PriceListId *uuid.UUID `json:"price_list_id,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	CategoryId *uuid.UUID `json:"category_id,omitempty"`
	
	OverridePrice *float64 `json:"override_price,omitempty"`
	
	DiscountPercentage *float64 `json:"discount_percentage,omitempty"`
	
	MarkupPercentage *float64 `json:"markup_percentage,omitempty"`
	
	MinPrice *float64 `json:"min_price,omitempty"`
	
	MaxPrice *float64 `json:"max_price,omitempty"`
	
	MinQuantity *float64 `json:"min_quantity,omitempty"`
	
	ProductId *string `json:"product_id,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdatePriceListItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PriceListId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.CategoryId != nil {
		hasUpdate = true
	}
	
	if r.OverridePrice != nil {
		hasUpdate = true
	}
	
	if r.DiscountPercentage != nil {
		hasUpdate = true
	}
	
	if r.MarkupPercentage != nil {
		hasUpdate = true
	}
	
	if r.MinPrice != nil {
		hasUpdate = true
	}
	
	if r.MaxPrice != nil {
		hasUpdate = true
	}
	
	if r.MinQuantity != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PriceListItemsListResponse represents a paginated list of price_list_items records
type PriceListItemsListResponse struct {
	Items      []*PriceListItemsResponse `json:"items"`
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
