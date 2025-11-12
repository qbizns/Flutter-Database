package customer_payment_application

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerPaymentApplicationsResponse represents a customer_payment_applications response
type CustomerPaymentApplicationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerPaymentId uuid.UUID `json:"customer_payment_id"`
	
	CustomerInvoiceId uuid.UUID `json:"customer_invoice_id"`
	
	AppliedAmount float64 `json:"applied_amount"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCustomerPaymentApplicationsRequest represents a request to create a customer_payment_applications
type CreateCustomerPaymentApplicationsRequest struct {
	
	CustomerPaymentId uuid.UUID `json:"customer_payment_id" validate:"required"`
	
	CustomerInvoiceId uuid.UUID `json:"customer_invoice_id" validate:"required"`
	
	AppliedAmount float64 `json:"applied_amount" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateCustomerPaymentApplicationsRequest) Validate() error {
	
	if r.CustomerPaymentId == uuid.Nil {
		return fmt.Errorf("customer_payment_id is required")
	}
	
	if r.CustomerInvoiceId == uuid.Nil {
		return fmt.Errorf("customer_invoice_id is required")
	}
	
	if r.AppliedAmount == nil {
		return fmt.Errorf("applied_amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerPaymentApplicationsRequest represents a request to update a customer_payment_applications
type UpdateCustomerPaymentApplicationsRequest struct {
	
	CustomerPaymentId *uuid.UUID `json:"customer_payment_id,omitempty" validate:"omitempty,required"`
	
	CustomerInvoiceId *uuid.UUID `json:"customer_invoice_id,omitempty" validate:"omitempty,required"`
	
	AppliedAmount *float64 `json:"applied_amount,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateCustomerPaymentApplicationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerPaymentId != nil {
		hasUpdate = true
	}
	
	if r.CustomerInvoiceId != nil {
		hasUpdate = true
	}
	
	if r.AppliedAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CustomerPaymentApplicationsListResponse represents a paginated list of customer_payment_applications records
type CustomerPaymentApplicationsListResponse struct {
	Items      []*CustomerPaymentApplicationsResponse `json:"items"`
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
