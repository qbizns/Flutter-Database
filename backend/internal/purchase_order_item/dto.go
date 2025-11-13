package purchase_order_item

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PurchaseOrderItemsResponse represents a purchase_order_items response
type PurchaseOrderItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PurchaseOrderId uuid.UUID `json:"purchase_order_id"`
	
	LineNumber int64 `json:"line_number"`
	
	ProductId uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	QuantityOrdered float64 `json:"quantity_ordered"`
	
	QuantityReceived *float64 `json:"quantity_received"`
	
	UnitCost float64 `json:"unit_cost"`
	
	Subtotal float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	TotalAmount float64 `json:"total_amount"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePurchaseOrderItemsRequest represents a request to create a purchase_order_items
type CreatePurchaseOrderItemsRequest struct {
	
	PurchaseOrderId uuid.UUID `json:"purchase_order_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	ProductId uuid.UUID `json:"product_id" validate:"required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	QuantityOrdered float64 `json:"quantity_ordered" validate:"required"`
	
	QuantityReceived *float64 `json:"quantity_received"`
	
	UnitCost float64 `json:"unit_cost" validate:"required"`
	
	Subtotal float64 `json:"subtotal" validate:"required"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	TotalAmount float64 `json:"total_amount" validate:"required"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	Notes *string `json:"notes"`
	
}

// Validate validates the create request
func (r *CreatePurchaseOrderItemsRequest) Validate() error {
	
	if r.PurchaseOrderId == uuid.Nil {
		return fmt.Errorf("purchase_order_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.ProductId == uuid.Nil {
		return fmt.Errorf("product_id is required")
	}
	
	if r.QuantityOrdered == nil {
		return fmt.Errorf("quantity_ordered is required")
	}
	
	if r.UnitCost == nil {
		return fmt.Errorf("unit_cost is required")
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

// UpdatePurchaseOrderItemsRequest represents a request to update a purchase_order_items
type UpdatePurchaseOrderItemsRequest struct {
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty" validate:"omitempty,required"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	QuantityOrdered *float64 `json:"quantity_ordered,omitempty" validate:"omitempty,required"`
	
	QuantityReceived *float64 `json:"quantity_received,omitempty"`
	
	UnitCost *float64 `json:"unit_cost,omitempty" validate:"omitempty,required"`
	
	Subtotal *float64 `json:"subtotal,omitempty" validate:"omitempty,required"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,required"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePurchaseOrderItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PurchaseOrderId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.QuantityOrdered != nil {
		hasUpdate = true
	}
	
	if r.QuantityReceived != nil {
		hasUpdate = true
	}
	
	if r.UnitCost != nil {
		hasUpdate = true
	}
	
	if r.Subtotal != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.ExpectedDeliveryDate != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PurchaseOrderItemsListResponse represents a paginated list of purchase_order_items records
type PurchaseOrderItemsListResponse struct {
	Items      []*PurchaseOrderItemsResponse `json:"items"`
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
