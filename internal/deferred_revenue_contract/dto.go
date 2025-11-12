package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeferredRevenueContractsResponse represents a deferred_revenue_contracts response
type DeferredRevenueContractsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerInvoiceId *uuid.UUID `json:"customer_invoice_id"`
	
	InvoiceLineId *uuid.UUID `json:"invoice_line_id"`
	
	ContractName *string `json:"contract_name"`
	
	TotalDeferredAmount float64 `json:"total_deferred_amount"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate time.Time `json:"end_date"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	DeferredAccountId uuid.UUID `json:"deferred_account_id"`
	
	RevenueAccountId uuid.UUID `json:"revenue_account_id"`
	
	Status *string `json:"status"`
	
	RecognizedAmount *float64 `json:"recognized_amount"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateDeferredRevenueContractsRequest represents a request to create a deferred_revenue_contracts
type CreateDeferredRevenueContractsRequest struct {
	
	CustomerInvoiceId *uuid.UUID `json:"customer_invoice_id"`
	
	InvoiceLineId *uuid.UUID `json:"invoice_line_id"`
	
	ContractName *string `json:"contract_name"`
	
	TotalDeferredAmount float64 `json:"total_deferred_amount" validate:"required"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate time.Time `json:"end_date" validate:"required"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	RecognitionMethod *string `json:"recognition_method"`
	
	DeferredAccountId uuid.UUID `json:"deferred_account_id" validate:"required"`
	
	RevenueAccountId uuid.UUID `json:"revenue_account_id" validate:"required"`
	
	Status *string `json:"status"`
	
	RecognizedAmount *float64 `json:"recognized_amount"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateDeferredRevenueContractsRequest) Validate() error {
	
	if r.TotalDeferredAmount == nil {
		return fmt.Errorf("total_deferred_amount is required")
	}
	
	if r.StartDate == nil {
		return fmt.Errorf("start_date is required")
	}
	
	if r.EndDate == nil {
		return fmt.Errorf("end_date is required")
	}
	
	if r.DeferredAccountId == uuid.Nil {
		return fmt.Errorf("deferred_account_id is required")
	}
	
	if r.RevenueAccountId == uuid.Nil {
		return fmt.Errorf("revenue_account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeferredRevenueContractsRequest represents a request to update a deferred_revenue_contracts
type UpdateDeferredRevenueContractsRequest struct {
	
	CustomerInvoiceId *uuid.UUID `json:"customer_invoice_id,omitempty"`
	
	InvoiceLineId *uuid.UUID `json:"invoice_line_id,omitempty"`
	
	ContractName *string `json:"contract_name,omitempty"`
	
	TotalDeferredAmount *float64 `json:"total_deferred_amount,omitempty" validate:"omitempty,required"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty" validate:"omitempty,required"`
	
	RecognitionMethod *string `json:"recognition_method,omitempty"`
	
	RecognitionMethod *string `json:"recognition_method,omitempty"`
	
	DeferredAccountId *uuid.UUID `json:"deferred_account_id,omitempty" validate:"omitempty,required"`
	
	RevenueAccountId *uuid.UUID `json:"revenue_account_id,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	RecognizedAmount *float64 `json:"recognized_amount,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeferredRevenueContractsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerInvoiceId != nil {
		hasUpdate = true
	}
	
	if r.InvoiceLineId != nil {
		hasUpdate = true
	}
	
	if r.ContractName != nil {
		hasUpdate = true
	}
	
	if r.TotalDeferredAmount != nil {
		hasUpdate = true
	}
	
	if r.StartDate != nil {
		hasUpdate = true
	}
	
	if r.EndDate != nil {
		hasUpdate = true
	}
	
	if r.RecognitionMethod != nil {
		hasUpdate = true
	}
	
	if r.RecognitionMethod != nil {
		hasUpdate = true
	}
	
	if r.DeferredAccountId != nil {
		hasUpdate = true
	}
	
	if r.RevenueAccountId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RecognizedAmount != nil {
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

// DeferredRevenueContractsListResponse represents a paginated list of deferred_revenue_contracts records
type DeferredRevenueContractsListResponse struct {
	Items      []*DeferredRevenueContractsResponse `json:"items"`
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
