package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// VendorPaymentsResponse represents a vendor_payments response
type VendorPaymentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PaymentNumber string `json:"payment_number"`
	
	SupplierId uuid.UUID `json:"supplier_id"`
	
	PaymentDate time.Time `json:"payment_date"`
	
	PaymentMethod string `json:"payment_method"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PaymentAmount float64 `json:"payment_amount"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	Memo *string `json:"memo"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateVendorPaymentsRequest represents a request to create a vendor_payments
type CreateVendorPaymentsRequest struct {
	
	PaymentNumber string `json:"payment_number" validate:"required"`
	
	SupplierId uuid.UUID `json:"supplier_id" validate:"required"`
	
	PaymentDate time.Time `json:"payment_date" validate:"required"`
	
	PaymentMethod string `json:"payment_method" validate:"required"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PaymentAmount float64 `json:"payment_amount" validate:"required"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	IsPosted *bool `json:"is_posted"`
	
	Memo *string `json:"memo"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateVendorPaymentsRequest) Validate() error {
	
	if r.PaymentNumber == "" {
		return fmt.Errorf("payment_number is required")
	}
	
	if r.SupplierId == uuid.Nil {
		return fmt.Errorf("supplier_id is required")
	}
	
	if r.PaymentDate == nil {
		return fmt.Errorf("payment_date is required")
	}
	
	if r.PaymentMethod == "" {
		return fmt.Errorf("payment_method is required")
	}
	
	if r.PaymentAmount == nil {
		return fmt.Errorf("payment_amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateVendorPaymentsRequest represents a request to update a vendor_payments
type UpdateVendorPaymentsRequest struct {
	
	PaymentNumber *string `json:"payment_number,omitempty" validate:"omitempty,required"`
	
	SupplierId *uuid.UUID `json:"supplier_id,omitempty" validate:"omitempty,required"`
	
	PaymentDate *time.Time `json:"payment_date,omitempty" validate:"omitempty,required"`
	
	PaymentMethod *string `json:"payment_method,omitempty" validate:"omitempty,required"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	PaymentAmount *float64 `json:"payment_amount,omitempty" validate:"omitempty,required"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id,omitempty"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	IsPosted *bool `json:"is_posted,omitempty"`
	
	Memo *string `json:"memo,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateVendorPaymentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PaymentNumber != nil {
		hasUpdate = true
	}
	
	if r.SupplierId != nil {
		hasUpdate = true
	}
	
	if r.PaymentDate != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.PaymentAmount != nil {
		hasUpdate = true
	}
	
	if r.BankAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.IsPosted != nil {
		hasUpdate = true
	}
	
	if r.Memo != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// VendorPaymentsListResponse represents a paginated list of vendor_payments records
type VendorPaymentsListResponse struct {
	Items      []*VendorPaymentsResponse `json:"items"`
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
