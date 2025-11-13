package sale_return_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SaleReturnItemsResponse represents a sale_return_items response
type SaleReturnItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SaleReturnId uuid.UUID `json:"sale_return_id"`
	
	OriginalSaleItemId *uuid.UUID `json:"original_sale_item_id"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	Quantity float64 `json:"quantity"`
	
	UnitPrice float64 `json:"unit_price"`
	
	Subtotal float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	TotalAmount float64 `json:"total_amount"`
	
	ReturnReasonId *uuid.UUID `json:"return_reason_id"`
	
	ReturnReasonNotes *string `json:"return_reason_notes"`
	
	ItemCondition *string `json:"item_condition"`
	
	ItemCondition *string `json:"item_condition"`
	
	IsRestockable *bool `json:"is_restockable"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateSaleReturnItemsRequest represents a request to create a sale_return_items
type CreateSaleReturnItemsRequest struct {
	
	SaleReturnId uuid.UUID `json:"sale_return_id" validate:"required"`
	
	OriginalSaleItemId *uuid.UUID `json:"original_sale_item_id"`
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	UnitPrice float64 `json:"unit_price" validate:"required"`
	
	Subtotal float64 `json:"subtotal" validate:"required"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	TotalAmount float64 `json:"total_amount" validate:"required"`
	
	ReturnReasonId *uuid.UUID `json:"return_reason_id"`
	
	ReturnReasonNotes *string `json:"return_reason_notes"`
	
	ItemCondition *string `json:"item_condition"`
	
	ItemCondition *string `json:"item_condition"`
	
	IsRestockable *bool `json:"is_restockable"`
	
}

// Validate validates the create request
func (r *CreateSaleReturnItemsRequest) Validate() error {
	
	if r.SaleReturnId == uuid.Nil {
		return fmt.Errorf("sale_return_id is required")
	}
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.Quantity == nil {
		return fmt.Errorf("quantity is required")
	}
	
	if r.UnitPrice == nil {
		return fmt.Errorf("unit_price is required")
	}
	
	if r.Subtotal == nil {
		return fmt.Errorf("subtotal is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSaleReturnItemsRequest represents a request to update a sale_return_items
type UpdateSaleReturnItemsRequest struct {
	
	SaleReturnId *uuid.UUID `json:"sale_return_id,omitempty" validate:"omitempty,required"`
	
	OriginalSaleItemId *uuid.UUID `json:"original_sale_item_id,omitempty"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	UnitPrice *float64 `json:"unit_price,omitempty" validate:"omitempty,required"`
	
	Subtotal *float64 `json:"subtotal,omitempty" validate:"omitempty,required"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,required"`
	
	ReturnReasonId *uuid.UUID `json:"return_reason_id,omitempty"`
	
	ReturnReasonNotes *string `json:"return_reason_notes,omitempty"`
	
	ItemCondition *string `json:"item_condition,omitempty"`
	
	ItemCondition *string `json:"item_condition,omitempty"`
	
	IsRestockable *bool `json:"is_restockable,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSaleReturnItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SaleReturnId != nil {
		hasUpdate = true
	}
	
	if r.OriginalSaleItemId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.UnitPrice != nil {
		hasUpdate = true
	}
	
	if r.Subtotal != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.ReturnReasonId != nil {
		hasUpdate = true
	}
	
	if r.ReturnReasonNotes != nil {
		hasUpdate = true
	}
	
	if r.ItemCondition != nil {
		hasUpdate = true
	}
	
	if r.ItemCondition != nil {
		hasUpdate = true
	}
	
	if r.IsRestockable != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// SaleReturnItemsListResponse represents a paginated list of sale_return_items records
type SaleReturnItemsListResponse struct {
	Items      []*SaleReturnItemsResponse `json:"items"`
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
