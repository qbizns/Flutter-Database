package product

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProductsResponse represents a products response
type ProductsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	Sku *string `json:"sku"`
	
	Barcode *string `json:"barcode"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	CostPrice *float64 `json:"cost_price"`
	
	SellingPrice float64 `json:"selling_price"`
	
	CompareAtPrice *float64 `json:"compare_at_price"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	IsTaxInclusive *bool `json:"is_tax_inclusive"`
	
	TrackInventory *bool `json:"track_inventory"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold"`
	
	Unit *string `json:"unit"`
	
	IsService *bool `json:"is_service"`
	
	IsComposite *bool `json:"is_composite"`
	
	HasVariants *bool `json:"has_variants"`
	
	ImageUrl *string `json:"image_url"`
	
	Images json.RawMessage `json:"images"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsActive *bool `json:"is_active"`
	
	IsFeatured *bool `json:"is_featured"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateProductsRequest represents a request to create a products
type CreateProductsRequest struct {
	
	Sku *string `json:"sku"`
	
	Barcode *string `json:"barcode"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	CategoryId *uuid.UUID `json:"category_id"`
	
	CostPrice *float64 `json:"cost_price"`
	
	SellingPrice float64 `json:"selling_price" validate:"required"`
	
	CompareAtPrice *float64 `json:"compare_at_price"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	IsTaxInclusive *bool `json:"is_tax_inclusive"`
	
	TrackInventory *bool `json:"track_inventory"`
	
	CurrentStock *float64 `json:"current_stock"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold"`
	
	Unit *string `json:"unit"`
	
	IsService *bool `json:"is_service"`
	
	IsComposite *bool `json:"is_composite"`
	
	HasVariants *bool `json:"has_variants"`
	
	ImageUrl *string `json:"image_url" validate:"url"`
	
	// Duplicate removed: Images json.RawMessage `json:"images"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsActive *bool `json:"is_active"`
	
	IsFeatured *bool `json:"is_featured"`
	
	// Duplicate removed: CustomFields json.RawMessage `json:"custom_fields"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateProductsRequest) Validate() error {
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.SellingPrice == nil {
		return fmt.Errorf("selling_price is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateProductsRequest represents a request to update a products
type UpdateProductsRequest struct {
	
	Sku *string `json:"sku,omitempty"`
	
	Barcode *string `json:"barcode,omitempty"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	CategoryId *uuid.UUID `json:"category_id,omitempty"`
	
	CostPrice *float64 `json:"cost_price,omitempty"`
	
	SellingPrice *float64 `json:"selling_price,omitempty" validate:"omitempty,required"`
	
	CompareAtPrice *float64 `json:"compare_at_price,omitempty"`
	
	TaxRate *float64 `json:"tax_rate,omitempty"`
	
	IsTaxInclusive *bool `json:"is_tax_inclusive,omitempty"`
	
	TrackInventory *bool `json:"track_inventory,omitempty"`
	
	CurrentStock *float64 `json:"current_stock,omitempty"`
	
	LowStockThreshold *float64 `json:"low_stock_threshold,omitempty"`
	
	Unit *string `json:"unit,omitempty"`
	
	IsService *bool `json:"is_service,omitempty"`
	
	IsComposite *bool `json:"is_composite,omitempty"`
	
	HasVariants *bool `json:"has_variants,omitempty"`
	
	ImageUrl *string `json:"image_url,omitempty" validate:"omitempty,url"`
	
	Images *json.RawMessage `json:"images,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsFeatured *bool `json:"is_featured,omitempty"`
	
	CustomFields *json.RawMessage `json:"custom_fields,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateProductsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Sku != nil {
		hasUpdate = true
	}
	
	if r.Barcode != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.CategoryId != nil {
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
	
	if r.TaxRate != nil {
		hasUpdate = true
	}
	
	if r.IsTaxInclusive != nil {
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
	
	if r.Unit != nil {
		hasUpdate = true
	}
	
	if r.IsService != nil {
		hasUpdate = true
	}
	
	if r.IsComposite != nil {
		hasUpdate = true
	}
	
	if r.HasVariants != nil {
		hasUpdate = true
	}
	
	if r.ImageUrl != nil {
		hasUpdate = true
	}
	
	if r.Images != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsFeatured != nil {
		hasUpdate = true
	}
	
	if r.CustomFields != nil {
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

// ProductsListResponse represents a paginated list of products records
type ProductsListResponse struct {
	Items      []*ProductsResponse `json:"items"`
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
