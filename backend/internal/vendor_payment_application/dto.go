package vendor_payment_application

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// VendorPaymentApplicationsResponse represents a vendor_payment_applications response
type VendorPaymentApplicationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	VendorPaymentId uuid.UUID `json:"vendor_payment_id"`
	
	VendorBillId uuid.UUID `json:"vendor_bill_id"`
	
	AppliedAmount float64 `json:"applied_amount"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateVendorPaymentApplicationsRequest represents a request to create a vendor_payment_applications
type CreateVendorPaymentApplicationsRequest struct {
	
	VendorPaymentId uuid.UUID `json:"vendor_payment_id" validate:"required"`
	
	VendorBillId uuid.UUID `json:"vendor_bill_id" validate:"required"`
	
	AppliedAmount float64 `json:"applied_amount" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateVendorPaymentApplicationsRequest) Validate() error {
	
	if r.VendorPaymentId == uuid.Nil {
		return fmt.Errorf("vendor_payment_id is required")
	}
	
	if r.VendorBillId == uuid.Nil {
		return fmt.Errorf("vendor_bill_id is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateVendorPaymentApplicationsRequest represents a request to update a vendor_payment_applications
type UpdateVendorPaymentApplicationsRequest struct {
	
	VendorPaymentId *uuid.UUID `json:"vendor_payment_id,omitempty" validate:"omitempty,required"`
	
	VendorBillId *uuid.UUID `json:"vendor_bill_id,omitempty" validate:"omitempty,required"`
	
	AppliedAmount *float64 `json:"applied_amount,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateVendorPaymentApplicationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.VendorPaymentId != nil {
		hasUpdate = true
	}
	
	if r.VendorBillId != nil {
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

// VendorPaymentApplicationsListResponse represents a paginated list of vendor_payment_applications records
type VendorPaymentApplicationsListResponse struct {
	Items      []*VendorPaymentApplicationsResponse `json:"items"`
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
