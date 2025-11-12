package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductVariantsResponse represents a product_variants response
type ProductVariantsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	VariantName string `json:"variant_name"`
	
	Sku *string `json:"sku"`
	
	Barcode *string `json:"barcode"`
	
	Attributes json.RawMessage `json:"attributes"`
	
	CostPrice *float64 `json:"cost_price"`
	
	SellingPrice *float64 `json:"selling_price"`
	
	CompareAtPrice *float64 `json:"compare_at_price"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	ReorderLevel *float64 `json:"reorder_level"`
	
	ReorderQuantity *float64 `json:"reorder_quantity"`
	
	Weight *float64 `json:"weight"`
	
	WeightUnit *string `json:"weight_unit"`
	
	Dimensions json.RawMessage `json:"dimensions"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	SortOrder *int64 `json:"sort_order"`
	
	ImageUrl *string `json:"image_url"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(costPrice *string `json:"(cost_price"`
	
	(sellingPrice *string `json:"(selling_price"`
	
	(compareAtPrice *string `json:"(compare_at_price"`
	
}

// CreateProductVariantsRequest represents a request to create a product_variants
type CreateProductVariantsRequest struct {
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	VariantName string `json:"variant_name" validate:"required"`
	
	Sku *string `json:"sku"`
	
	Barcode *string `json:"barcode"`
	
	Attributes json.RawMessage `json:"attributes"`
	
	CostPrice *float64 `json:"cost_price"`
	
	SellingPrice *float64 `json:"selling_price"`
	
	CompareAtPrice *float64 `json:"compare_at_price"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	ReorderLevel *float64 `json:"reorder_level"`
	
	ReorderQuantity *float64 `json:"reorder_quantity"`
	
	Weight *float64 `json:"weight"`
	
	WeightUnit *string `json:"weight_unit"`
	
	Dimensions json.RawMessage `json:"dimensions"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	SortOrder *int64 `json:"sort_order"`
	
	ImageUrl *string `json:"image_url" validate:"url"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(costPrice *string `json:"(cost_price"`
	
	(sellingPrice *string `json:"(selling_price"`
	
	(compareAtPrice *string `json:"(compare_at_price"`
	
}

// Validate validates the create request
func (r *CreateProductVariantsRequest) Validate() error {
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.VariantName == "" {
		return fmt.Errorf("variant_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductVariantsRequest represents a request to update a product_variants
type UpdateProductVariantsRequest struct {
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	VariantName *string `json:"variant_name,omitempty" validate:"omitempty,required"`
	
	Sku *string `json:"sku,omitempty"`
	
	Barcode *string `json:"barcode,omitempty"`
	
	Attributes *json.RawMessage `json:"attributes,omitempty"`
	
	CostPrice *float64 `json:"cost_price,omitempty"`
	
	SellingPrice *float64 `json:"selling_price,omitempty"`
	
	CompareAtPrice *float64 `json:"compare_at_price,omitempty"`
	
	CurrentStock *float64 `json:"current_stock,omitempty"`
	
	ReorderLevel *float64 `json:"reorder_level,omitempty"`
	
	ReorderQuantity *float64 `json:"reorder_quantity,omitempty"`
	
	Weight *float64 `json:"weight,omitempty"`
	
	WeightUnit *string `json:"weight_unit,omitempty"`
	
	Dimensions *json.RawMessage `json:"dimensions,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	ImageUrl *string `json:"image_url,omitempty" validate:"omitempty,url"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(costPrice *string `json:"(cost_price,omitempty"`
	
	(sellingPrice *string `json:"(selling_price,omitempty"`
	
	(compareAtPrice *string `json:"(compare_at_price,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateProductVariantsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.VariantName != nil {
		hasUpdate = true
	}
	
	if r.Sku != nil {
		hasUpdate = true
	}
	
	if r.Barcode != nil {
		hasUpdate = true
	}
	
	if r.Attributes != nil {
		hasUpdate = true
	}
	
	if r.CostPrice != nil {
		hasUpdate = true
	}
	
	if r.SellingPrice != nil {
		hasUpdate = true
	}
	
	if r.CompareAtPrice != nil {
		hasUpdate = true
	}
	
	if r.CurrentStock != nil {
		hasUpdate = true
	}
	
	if r.ReorderLevel != nil {
		hasUpdate = true
	}
	
	if r.ReorderQuantity != nil {
		hasUpdate = true
	}
	
	if r.Weight != nil {
		hasUpdate = true
	}
	
	if r.WeightUnit != nil {
		hasUpdate = true
	}
	
	if r.Dimensions != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.ImageUrl != nil {
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
	
	if r.(costPrice != nil {
		hasUpdate = true
	}
	
	if r.(sellingPrice != nil {
		hasUpdate = true
	}
	
	if r.(compareAtPrice != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ProductVariantsListResponse represents a paginated list of product_variants records
type ProductVariantsListResponse struct {
	Items      []*ProductVariantsResponse `json:"items"`
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
