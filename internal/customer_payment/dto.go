package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerPaymentsResponse represents a customer_payments response
type CustomerPaymentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PaymentNumber string `json:"payment_number"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	PaymentDate time.Time `json:"payment_date"`
	
	PaymentMethod string `json:"payment_method"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PaymentAmount float64 `json:"payment_amount"`
	
	DepositAccountId *uuid.UUID `json:"deposit_account_id"`
	
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

// CreateCustomerPaymentsRequest represents a request to create a customer_payments
type CreateCustomerPaymentsRequest struct {
	
	PaymentNumber string `json:"payment_number" validate:"required"`
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	PaymentDate time.Time `json:"payment_date" validate:"required"`
	
	PaymentMethod string `json:"payment_method" validate:"required"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	PaymentAmount float64 `json:"payment_amount" validate:"required"`
	
	DepositAccountId *uuid.UUID `json:"deposit_account_id"`
	
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
func (r *CreateCustomerPaymentsRequest) Validate() error {
	
	if r.PaymentNumber == "" {
		return fmt.Errorf("payment_number is required")
	}
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
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

// UpdateCustomerPaymentsRequest represents a request to update a customer_payments
type UpdateCustomerPaymentsRequest struct {
	
	PaymentNumber *string `json:"payment_number,omitempty" validate:"omitempty,required"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	PaymentDate *time.Time `json:"payment_date,omitempty" validate:"omitempty,required"`
	
	PaymentMethod *string `json:"payment_method,omitempty" validate:"omitempty,required"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	PaymentAmount *float64 `json:"payment_amount,omitempty" validate:"omitempty,required"`
	
	DepositAccountId *uuid.UUID `json:"deposit_account_id,omitempty"`
	
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
func (r *UpdateCustomerPaymentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PaymentNumber != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
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
	
	if r.DepositAccountId != nil {
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

// CustomerPaymentsListResponse represents a paginated list of customer_payments records
type CustomerPaymentsListResponse struct {
	Items      []*CustomerPaymentsResponse `json:"items"`
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
