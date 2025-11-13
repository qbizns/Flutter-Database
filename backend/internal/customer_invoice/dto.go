package customer_invoice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerInvoicesResponse represents a customer_invoices response
type CustomerInvoicesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	InvoiceNumber string `json:"invoice_number"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	InvoiceDate time.Time `json:"invoice_date"`
	
	DueDate time.Time `json:"due_date"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	Subtotal *float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	TotalAmount float64 `json:"total_amount"`
	
	PaidAmount *float64 `json:"paid_amount"`
	
	BalanceDue *float64 `json:"balance_due"`
	
	Status *string `json:"status"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Memo *string `json:"memo"`
	
	Attachments json.RawMessage `json:"attachments"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCustomerInvoicesRequest represents a request to create a customer_invoices
type CreateCustomerInvoicesRequest struct {
	
	InvoiceNumber string `json:"invoice_number" validate:"required"`
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	InvoiceDate time.Time `json:"invoice_date" validate:"required"`
	
	DueDate time.Time `json:"due_date" validate:"required"`
	
	PaymentTerms *string `json:"payment_terms"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	Subtotal *float64 `json:"subtotal"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	TotalAmount float64 `json:"total_amount" validate:"required"`
	
	PaidAmount *float64 `json:"paid_amount"`
	
	BalanceDue *float64 `json:"balance_due"`
	
	// 	Status *string `json:"status"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Memo *string `json:"memo"`
	
	// Duplicate removed: Attachments json.RawMessage `json:"attachments"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateCustomerInvoicesRequest) Validate() error {
	
	if r.InvoiceNumber == "" {
		return fmt.Errorf("invoice_number is required")
	}
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	
	if r.InvoiceDate.IsZero() {
		return fmt.Errorf("invoice_date is required")
	}
	
	if r.DueDate.IsZero() {
		return fmt.Errorf("due_date is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerInvoicesRequest represents a request to update a customer_invoices
type UpdateCustomerInvoicesRequest struct {
	
	InvoiceNumber *string `json:"invoice_number,omitempty" validate:"omitempty,required"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	InvoiceDate *time.Time `json:"invoice_date,omitempty" validate:"omitempty,required"`
	
	DueDate *time.Time `json:"due_date,omitempty" validate:"omitempty,required"`
	
	PaymentTerms *string `json:"payment_terms,omitempty"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	Subtotal *float64 `json:"subtotal,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	TotalAmount *float64 `json:"total_amount,omitempty" validate:"omitempty,required"`
	
	PaidAmount *float64 `json:"paid_amount,omitempty"`
	
	BalanceDue *float64 `json:"balance_due,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	IsPosted *bool `json:"is_posted,omitempty"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Memo *string `json:"memo,omitempty"`
	
	Attachments *json.RawMessage `json:"attachments,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomerInvoicesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.InvoiceNumber != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.InvoiceDate != nil {
		hasUpdate = true
	}
	
	if r.DueDate != nil {
		hasUpdate = true
	}
	
	if r.PaymentTerms != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
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
	
	if r.BalanceDue != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.IsPosted != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Memo != nil {
		hasUpdate = true
	}
	
	if r.Attachments != nil {
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

// CustomerInvoicesListResponse represents a paginated list of customer_invoices records
type CustomerInvoicesListResponse struct {
	Items      []*CustomerInvoicesResponse `json:"items"`
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
