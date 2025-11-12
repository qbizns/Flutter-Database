package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SaleReturnsResponse represents a sale_returns response
type SaleReturnsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ReturnNumber string `json:"return_number"`
	
	OriginalSaleId *uuid.UUID `json:"original_sale_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	ReturnDate time.Time `json:"return_date"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	RefundAmount *float64 `json:"refund_amount"`
	
	RestockingFee *float64 `json:"restocking_fee"`
	
	RefundMethod *string `json:"refund_method"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateSaleReturnsRequest represents a request to create a sale_returns
type CreateSaleReturnsRequest struct {
	
	ReturnNumber string `json:"return_number" validate:"required"`
	
	OriginalSaleId *uuid.UUID `json:"original_sale_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	ReturnDate time.Time `json:"return_date" validate:"required"`
	
	TotalAmount *float64 `json:"total_amount"`
	
	RefundAmount *float64 `json:"refund_amount"`
	
	RestockingFee *float64 `json:"restocking_fee"`
	
	RefundMethod *string `json:"refund_method"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateSaleReturnsRequest) Validate() error {
	
	if r.ReturnNumber == "" {
		return fmt.Errorf("return_number is required")
	}
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.ReturnDate == nil {
		return fmt.Errorf("return_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSaleReturnsRequest represents a request to update a sale_returns
type UpdateSaleReturnsRequest struct {
	
	ReturnNumber *string `json:"return_number,omitempty" validate:"omitempty,required"`
	
	OriginalSaleId *uuid.UUID `json:"original_sale_id,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	ReturnDate *time.Time `json:"return_date,omitempty" validate:"omitempty,required"`
	
	TotalAmount *float64 `json:"total_amount,omitempty"`
	
	RefundAmount *float64 `json:"refund_amount,omitempty"`
	
	RestockingFee *float64 `json:"restocking_fee,omitempty"`
	
	RefundMethod *string `json:"refund_method,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSaleReturnsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ReturnNumber != nil {
		hasUpdate = true
	}
	
	if r.OriginalSaleId != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.ReturnDate != nil {
		hasUpdate = true
	}
	
	if r.TotalAmount != nil {
		hasUpdate = true
	}
	
	if r.RefundAmount != nil {
		hasUpdate = true
	}
	
	if r.RestockingFee != nil {
		hasUpdate = true
	}
	
	if r.RefundMethod != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// SaleReturnsListResponse represents a paginated list of sale_returns records
type SaleReturnsListResponse struct {
	Items      []*SaleReturnsResponse `json:"items"`
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
