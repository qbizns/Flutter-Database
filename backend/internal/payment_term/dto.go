package payment_term

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentTermsResponse represents a payment_terms response
type PaymentTermsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	TermCode string `json:"term_code"`
	
	TermName string `json:"term_name"`
	
	Note *string `json:"note"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePaymentTermsRequest represents a request to create a payment_terms
type CreatePaymentTermsRequest struct {
	
	TermCode string `json:"term_code" validate:"required"`
	
	TermName string `json:"term_name" validate:"required"`
	
	Note *string `json:"note"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePaymentTermsRequest) Validate() error {
	
	if r.TermCode == "" {
		return fmt.Errorf("term_code is required")
	}
	
	if r.TermName == "" {
		return fmt.Errorf("term_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePaymentTermsRequest represents a request to update a payment_terms
type UpdatePaymentTermsRequest struct {
	
	TermCode *string `json:"term_code,omitempty" validate:"omitempty,required"`
	
	TermName *string `json:"term_name,omitempty" validate:"omitempty,required"`
	
	Note *string `json:"note,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePaymentTermsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TermCode != nil {
		hasUpdate = true
	}
	
	if r.TermName != nil {
		hasUpdate = true
	}
	
	if r.Note != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// PaymentTermsListResponse represents a paginated list of payment_terms records
type PaymentTermsListResponse struct {
	Items      []*PaymentTermsResponse `json:"items"`
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
