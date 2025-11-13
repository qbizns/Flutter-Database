package sale_item

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SaleItemsResponse represents a sale_items response
type SaleItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	SaleId uuid.UUID `json:"sale_id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductName string `json:"product_name"`
	
	ProductSku *string `json:"product_sku"`
	
	Quantity float64 `json:"quantity"`
	
	Unit *string `json:"unit"`
	
	UnitPrice float64 `json:"unit_price"`
	
	CostPrice *float64 `json:"cost_price"`
	
	Subtotal float64 `json:"subtotal"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	Total float64 `json:"total"`
	
	DiscountType *string `json:"discount_type"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	Notes *string `json:"notes"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
}

// CreateSaleItemsRequest represents a request to create a sale_items
type CreateSaleItemsRequest struct {
	
	SaleId uuid.UUID `json:"sale_id" validate:"required"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductName string `json:"product_name" validate:"required"`
	
	ProductSku *string `json:"product_sku"`
	
	Quantity float64 `json:"quantity" validate:"required"`
	
	Unit *string `json:"unit"`
	
	UnitPrice float64 `json:"unit_price" validate:"required"`
	
	CostPrice *float64 `json:"cost_price"`
	
	Subtotal float64 `json:"subtotal" validate:"required"`
	
	TaxRate *float64 `json:"tax_rate"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	Total float64 `json:"total" validate:"required"`
	
	DiscountType *string `json:"discount_type"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: CustomFields json.RawMessage `json:"custom_fields"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateSaleItemsRequest) Validate() error {
	
	if r.SaleId == uuid.Nil {
		return fmt.Errorf("sale_id is required")
	}
	
	if r.ProductName == "" {
		return fmt.Errorf("product_name is required")
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
	
	if r.Total == nil {
		return fmt.Errorf("total is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSaleItemsRequest represents a request to update a sale_items
type UpdateSaleItemsRequest struct {
	
	SaleId *uuid.UUID `json:"sale_id,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty"`
	
	ProductName *string `json:"product_name,omitempty" validate:"omitempty,required"`
	
	ProductSku *string `json:"product_sku,omitempty"`
	
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,required"`
	
	Unit *string `json:"unit,omitempty"`
	
	UnitPrice *float64 `json:"unit_price,omitempty" validate:"omitempty,required"`
	
	CostPrice *float64 `json:"cost_price,omitempty"`
	
	Subtotal *float64 `json:"subtotal,omitempty" validate:"omitempty,required"`
	
	TaxRate *float64 `json:"tax_rate,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	Total *float64 `json:"total,omitempty" validate:"omitempty,required"`
	
	DiscountType *string `json:"discount_type,omitempty"`
	
	DiscountValue *float64 `json:"discount_value,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CustomFields *json.RawMessage `json:"custom_fields,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSaleItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductName != nil {
		hasUpdate = true
	}
	
	if r.ProductSku != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.Unit != nil {
		hasUpdate = true
	}
	
	if r.UnitPrice != nil {
		hasUpdate = true
	}
	
	if r.CostPrice != nil {
		hasUpdate = true
	}
	
	if r.Subtotal != nil {
		hasUpdate = true
	}
	
	if r.TaxRate != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.Total != nil {
		hasUpdate = true
	}
	
	if r.DiscountType != nil {
		hasUpdate = true
	}
	
	if r.DiscountValue != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CustomFields != nil {
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

// SaleItemsListResponse represents a paginated list of sale_items records
type SaleItemsListResponse struct {
	Items      []*SaleItemsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	// Duplicate removed: Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
