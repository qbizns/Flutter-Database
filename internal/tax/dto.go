package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TaxesResponse represents a taxes response
type TaxesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	TaxGroupId *uuid.UUID `json:"tax_group_id"`
	
	TaxCode string `json:"tax_code"`
	
	TaxName string `json:"tax_name"`
	
	TaxRate float64 `json:"tax_rate"`
	
	TaxScope string `json:"tax_scope"`
	
	IsPriceInclusive *bool `json:"is_price_inclusive"`
	
	TaxAccountId uuid.UUID `json:"tax_account_id"`
	
	TaxRefundAccountId *uuid.UUID `json:"tax_refund_account_id"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateTaxesRequest represents a request to create a taxes
type CreateTaxesRequest struct {
	
	TaxGroupId *uuid.UUID `json:"tax_group_id"`
	
	TaxCode string `json:"tax_code" validate:"required"`
	
	TaxName string `json:"tax_name" validate:"required"`
	
	TaxRate float64 `json:"tax_rate" validate:"required"`
	
	TaxScope string `json:"tax_scope" validate:"required"`
	
	IsPriceInclusive *bool `json:"is_price_inclusive"`
	
	TaxAccountId uuid.UUID `json:"tax_account_id" validate:"required"`
	
	TaxRefundAccountId *uuid.UUID `json:"tax_refund_account_id"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateTaxesRequest) Validate() error {
	
	if r.TaxCode == "" {
		return fmt.Errorf("tax_code is required")
	}
	
	if r.TaxName == "" {
		return fmt.Errorf("tax_name is required")
	}
	
	if r.TaxRate == nil {
		return fmt.Errorf("tax_rate is required")
	}
	
	if r.TaxScope == "" {
		return fmt.Errorf("tax_scope is required")
	}
	
	if r.TaxAccountId == uuid.Nil {
		return fmt.Errorf("tax_account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTaxesRequest represents a request to update a taxes
type UpdateTaxesRequest struct {
	
	TaxGroupId *uuid.UUID `json:"tax_group_id,omitempty"`
	
	TaxCode *string `json:"tax_code,omitempty" validate:"omitempty,required"`
	
	TaxName *string `json:"tax_name,omitempty" validate:"omitempty,required"`
	
	TaxRate *float64 `json:"tax_rate,omitempty" validate:"omitempty,required"`
	
	TaxScope *string `json:"tax_scope,omitempty" validate:"omitempty,required"`
	
	IsPriceInclusive *bool `json:"is_price_inclusive,omitempty"`
	
	TaxAccountId *uuid.UUID `json:"tax_account_id,omitempty" validate:"omitempty,required"`
	
	TaxRefundAccountId *uuid.UUID `json:"tax_refund_account_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTaxesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TaxGroupId != nil {
		hasUpdate = true
	}
	
	if r.TaxCode != nil {
		hasUpdate = true
	}
	
	if r.TaxName != nil {
		hasUpdate = true
	}
	
	if r.TaxRate != nil {
		hasUpdate = true
	}
	
	if r.TaxScope != nil {
		hasUpdate = true
	}
	
	if r.IsPriceInclusive != nil {
		hasUpdate = true
	}
	
	if r.TaxAccountId != nil {
		hasUpdate = true
	}
	
	if r.TaxRefundAccountId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// TaxesListResponse represents a paginated list of taxes records
type TaxesListResponse struct {
	Items      []*TaxesResponse `json:"items"`
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
