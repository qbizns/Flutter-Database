package invoice_payment_schedule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InvoicePaymentSchedulesResponse represents a invoice_payment_schedules response
type InvoicePaymentSchedulesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SourceType string `json:"source_type"`
	
	SourceId uuid.UUID `json:"source_id"`
	
	LineNumber int64 `json:"line_number"`
	
	DueDate time.Time `json:"due_date"`
	
	AmountDue float64 `json:"amount_due"`
	
	AmountPaid *float64 `json:"amount_paid"`
	
	Status *string `json:"status"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateInvoicePaymentSchedulesRequest represents a request to create a invoice_payment_schedules
type CreateInvoicePaymentSchedulesRequest struct {
	
	SourceType string `json:"source_type" validate:"required"`
	
	SourceId uuid.UUID `json:"source_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	DueDate time.Time `json:"due_date" validate:"required"`
	
	AmountDue float64 `json:"amount_due" validate:"required"`
	
	AmountPaid *float64 `json:"amount_paid"`
	
	Status *string `json:"status"`
	
}

// Validate validates the create request
func (r *CreateInvoicePaymentSchedulesRequest) Validate() error {
	
	if r.SourceType == "" {
		return fmt.Errorf("source_type is required")
	}
	
	if r.SourceId == uuid.Nil {
		return fmt.Errorf("source_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.DueDate == nil {
		return fmt.Errorf("due_date is required")
	}
	
	if r.AmountDue == nil {
		return fmt.Errorf("amount_due is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateInvoicePaymentSchedulesRequest represents a request to update a invoice_payment_schedules
type UpdateInvoicePaymentSchedulesRequest struct {
	
	SourceType *string `json:"source_type,omitempty" validate:"omitempty,required"`
	
	SourceId *uuid.UUID `json:"source_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	DueDate *time.Time `json:"due_date,omitempty" validate:"omitempty,required"`
	
	AmountDue *float64 `json:"amount_due,omitempty" validate:"omitempty,required"`
	
	AmountPaid *float64 `json:"amount_paid,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateInvoicePaymentSchedulesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SourceType != nil {
		hasUpdate = true
	}
	
	if r.SourceId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.DueDate != nil {
		hasUpdate = true
	}
	
	if r.AmountDue != nil {
		hasUpdate = true
	}
	
	if r.AmountPaid != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// InvoicePaymentSchedulesListResponse represents a paginated list of invoice_payment_schedules records
type InvoicePaymentSchedulesListResponse struct {
	Items      []*InvoicePaymentSchedulesResponse `json:"items"`
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
