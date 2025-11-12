package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SalesResponse represents a sales response
type SalesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SaleNumber string `json:"sale_number"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	TransactionType string `json:"transaction_type"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	CashierId *uuid.UUID `json:"cashier_id"`
	
	Subtotal float64 `json:"subtotal"`
	
	TaxAmount float64 `json:"tax_amount"`
	
	DiscountAmount float64 `json:"discount_amount"`
	
	TotalAmount float64 `json:"total_amount"`
	
	PaidAmount float64 `json:"paid_amount"`
	
	ChangeAmount float64 `json:"change_amount"`
	
	OutstandingAmount float64 `json:"outstanding_amount"`
	
	PaymentStatus string `json:"payment_status"`
	
	DiscountType *string `json:"discount_type"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	DiscountReason *string `json:"discount_reason"`
	
	TransactionDate time.Time `json:"transaction_date"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Notes *string `json:"notes"`
	
	InternalNotes *string `json:"internal_notes"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateSalesRequest represents a request to create a sales
type CreateSalesRequest struct {
	
	SaleNumber string `json:"sale_number" validate:"required"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	TransactionType string `json:"transaction_type" validate:"required"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	CashierId *uuid.UUID `json:"cashier_id"`
	
	Subtotal float64 `json:"subtotal" validate:"required"`
	
	TaxAmount float64 `json:"tax_amount" validate:"required"`
	
	DiscountAmount float64 `json:"discount_amount" validate:"required"`
	
	TotalAmount float64 `json:"total_amount" validate:"required"`
	
	PaidAmount float64 `json:"paid_amount" validate:"required"`
	
	ChangeAmount float64 `json:"change_amount" validate:"required"`
	
	OutstandingAmount float64 `json:"outstanding_amount" validate:"required"`
	
	PaymentStatus string `json:"payment_status" validate:"required"`
	
	DiscountType *string `json:"discount_type"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	DiscountReason *string `json:"discount_reason"`
	
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	Notes *string `json:"notes"`
	
	InternalNotes *string `json:"internal_notes"`
	
	CustomFields json.RawMessage `json:"custom_fields"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateSalesRequest) Validate() error {
	
	if r.SaleNumber == "" {
		return fmt.Errorf("sale_number is required")
	}
	
	if r.TransactionType == "" {
		return fmt.Errorf("transaction_type is required")
	}
	
	if r.Subtotal == nil {
		return fmt.Errorf("subtotal is required")
	}
	
	if r.TaxAmount == nil {
		return fmt.Errorf("tax_amount is required")
	}
	
	if r.DiscountAmount == nil {
		return fmt.Errorf("discount_amount is required")
	}
	
	if r.TotalAmount == nil {
		return fmt.Errorf("total_amount is required")
	}
	
	if r.PaidAmount == nil {
		return fmt.Errorf("paid_amount is required")
	}
	
	if r.ChangeAmount == nil {
		return fmt.Errorf("change_amount is required")
	}
	
	if r.OutstandingAmount == nil {
		return fmt.Errorf("outstanding_amount is required")
	}
	
	if r.PaymentStatus == "" {
		return fmt.Errorf("payment_status is required")
	}
	
	if r.TransactionDate == nil {
		return fmt.Errorf("transaction_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSalesRequest represents a request to update a sales
type UpdateSalesRequest struct {
	
	SaleNumber *string `json:"sale_number,omitempty" validate:"omitempty,required"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,required"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	CashierId *uuid.UUID `json:"cashier_id,omitempty"`
	
	Subtotal *float64 `json:"subtotal,omitempty" validate:"omitempty,required"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty" validate:"omitempty,required"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty" validate:"omitempty,required"`
	
	TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,required"`
	
	PaidAmount *float64 `json:"paid_amount,omitempty" validate:"omitempty,required"`
	
	ChangeAmount *float64 `json:"change_amount,omitempty" validate:"omitempty,required"`
	
	OutstandingAmount *float64 `json:"outstanding_amount,omitempty" validate:"omitempty,required"`
	
	PaymentStatus *string `json:"payment_status,omitempty" validate:"omitempty,required"`
	
	DiscountType *string `json:"discount_type,omitempty"`
	
	DiscountValue *float64 `json:"discount_value,omitempty"`
	
	DiscountReason *string `json:"discount_reason,omitempty"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty" validate:"omitempty,required"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	InternalNotes *string `json:"internal_notes,omitempty"`
	
	CustomFields *json.RawMessage `json:"custom_fields,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSalesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SaleNumber != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.CashierId != nil {
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
	
	if r.PaidAmount != nil {
		hasUpdate = true
	}
	
	if r.ChangeAmount != nil {
		hasUpdate = true
	}
	
	if r.OutstandingAmount != nil {
		hasUpdate = true
	}
	
	if r.PaymentStatus != nil {
		hasUpdate = true
	}
	
	if r.DiscountType != nil {
		hasUpdate = true
	}
	
	if r.DiscountValue != nil {
		hasUpdate = true
	}
	
	if r.DiscountReason != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.InternalNotes != nil {
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

// SalesListResponse represents a paginated list of sales records
type SalesListResponse struct {
	Items      []*SalesResponse `json:"items"`
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
