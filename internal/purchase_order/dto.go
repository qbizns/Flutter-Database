package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PurchaseOrdersResponse represents a purchase_orders response
type PurchaseOrdersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PoNumber string `json:"po_number"`
	
	SupplierId uuid.UUID `json:"supplier_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	OrderDate time.Time `json:"order_date"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	ActualDeliveryDate *time.Time `json:"actual_delivery_date"`
	
	SubtotalAmount *float64 `json:"subtotal_amount"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	ShippingAmount *float64 `json:"shipping_amount"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	PaymentDueDate *time.Time `json:"payment_due_date"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePurchaseOrdersRequest represents a request to create a purchase_orders
type CreatePurchaseOrdersRequest struct {
	
	PoNumber string `json:"po_number" validate:"required"`
	
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	OrderDate time.Time `json:"order_date" validate:"required"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date"`
	
	ActualDeliveryDate *time.Time `json:"actual_delivery_date"`
	
	SubtotalAmount *float64 `json:"subtotal_amount"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	ShippingAmount *float64 `json:"shipping_amount"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	PaymentDueDate *time.Time `json:"payment_due_date"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePurchaseOrdersRequest) Validate() error {
	
	if r.PoNumber == "" {
		return fmt.Errorf("po_number is required")
	}
	
	if r.SupplierId == uuid.Nil {
		return fmt.Errorf("supplier_id is required")
	}
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.OrderDate == nil {
		return fmt.Errorf("order_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePurchaseOrdersRequest represents a request to update a purchase_orders
type UpdatePurchaseOrdersRequest struct {
	
	PoNumber *string `json:"po_number,omitempty" validate:"omitempty,required"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty" validate:"omitempty,required"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	OrderDate *time.Time `json:"order_date,omitempty" validate:"omitempty,required"`
	
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date,omitempty"`
	
	ActualDeliveryDate *time.Time `json:"actual_delivery_date,omitempty"`
	
	SubtotalAmount *float64 `json:"subtotal_amount,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	ShippingAmount *float64 `json:"shipping_amount,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty"`
	
	PaymentTerms *string `json:"payment_terms,omitempty"`
	
	PaymentDueDate *time.Time `json:"payment_due_date,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePurchaseOrdersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PoNumber != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.OrderDate != nil {
		hasUpdate = true
	}
	
	if r.ExpectedDeliveryDate != nil {
		hasUpdate = true
	}
	
	if r.ActualDeliveryDate != nil {
		hasUpdate = true
	}
	
	if r.SubtotalAmount != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.ShippingAmount != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.PaymentTerms != nil {
		hasUpdate = true
	}
	
	if r.PaymentDueDate != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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

// PurchaseOrdersListResponse represents a paginated list of purchase_orders records
type PurchaseOrdersListResponse struct {
	Items      []*PurchaseOrdersResponse `json:"items"`
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
