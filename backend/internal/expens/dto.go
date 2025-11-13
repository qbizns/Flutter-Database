package expens

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ExpensesResponse represents a expenses response
type ExpensesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	ExpenseNumber string `json:"expense_number"`
	
	ExpenseDate time.Time `json:"expense_date"`
	
	Category string `json:"category"`
	
	Subcategory *string `json:"subcategory"`
	
	PayeeName string `json:"payee_name"`
	
	PaymentMethod *string `json:"payment_method"`
	
	Amount float64 `json:"amount"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	TotalAmount float64 `json:"total_amount"`
	
	Currency *string `json:"currency"`
	
	Status *string `json:"status"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	ReceiptUrl *string `json:"receipt_url"`
	
	AttachmentUrls json.RawMessage `json:"attachment_urls"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Amount *string `json:"amount"`
	
	TaxAmount *string `json:"tax_amount"`
	
	TotalAmount *string `json:"total_amount"`
	
}

// CreateExpensesRequest represents a request to create a expenses
type CreateExpensesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	ExpenseNumber string `json:"expense_number" validate:"required"`
	
	ExpenseDate time.Time `json:"expense_date" validate:"required"`
	
	Category string `json:"category" validate:"required"`
	
	Subcategory *string `json:"subcategory"`
	
	PayeeName string `json:"payee_name" validate:"required"`
	
	PaymentMethod *string `json:"payment_method"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	TotalAmount float64 `json:"total_amount" validate:"required"`
	
	Currency *string `json:"currency"`
	
	// 	Status *string `json:"status"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id"`
	
	ReceiptUrl *string `json:"receipt_url" validate:"url"`
	
	// Duplicate removed: AttachmentUrls json.RawMessage `json:"attachment_urls"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Amount *string `json:"amount"`
	
	TaxAmount *string `json:"tax_amount"`
	
	TotalAmount *string `json:"total_amount"`
	
}

// Validate validates the create request
func (r *CreateExpensesRequest) Validate() error {
	
	if r.ExpenseNumber == "" {
		return fmt.Errorf("expense_number is required")
	}
	
	if r.ExpenseDate.IsZero() {
		return fmt.Errorf("expense_date is required")
	}
	
	if r.Category == "" {
		return fmt.Errorf("category is required")
	}
	
	if r.PayeeName == "" {
		return fmt.Errorf("payee_name is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateExpensesRequest represents a request to update a expenses
type UpdateExpensesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	ExpenseNumber *string `json:"expense_number,omitempty" validate:"omitempty,required"`
	
	ExpenseDate *time.Time `json:"expense_date,omitempty" validate:"omitempty,required"`
	
	Category *string `json:"category,omitempty" validate:"omitempty,required"`
	
	Subcategory *string `json:"subcategory,omitempty"`
	
	PayeeName *string `json:"payee_name,omitempty" validate:"omitempty,required"`
	
	PaymentMethod *string `json:"payment_method,omitempty"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,required"`
	
	Currency *string `json:"currency,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id,omitempty"`
	
	ReceiptUrl *string `json:"receipt_url,omitempty" validate:"omitempty,url"`
	
	AttachmentUrls *json.RawMessage `json:"attachment_urls,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	Amount *string `json:"amount,omitempty"`
	
	TaxAmount *string `json:"tax_amount,omitempty"`
	
	TotalAmount *string `json:"total_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateExpensesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.ExpenseNumber != nil {
		hasUpdate = true
	}
	
	if r.ExpenseDate != nil {
		hasUpdate = true
	}
	
	if r.Category != nil {
		hasUpdate = true
	}
	
	if r.Subcategory != nil {
		hasUpdate = true
	}
	
	if r.PayeeName != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.Currency != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.PurchaseOrderId != nil {
		hasUpdate = true
	}
	
	if r.ReceiptUrl != nil {
		hasUpdate = true
	}
	
	if r.AttachmentUrls != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ExpensesListResponse represents a paginated list of expenses records
type ExpensesListResponse struct {
	Items      []*ExpensesResponse `json:"items"`
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
