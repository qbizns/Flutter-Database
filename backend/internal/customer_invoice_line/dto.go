package customer_invoice_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerInvoiceLinesResponse represents a customer_invoice_lines response
type CustomerInvoiceLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerInvoiceId uuid.UUID `json:"customer_invoice_id"`
	
	LineNumber int64 `json:"line_number"`
	
	RevenueAccountId uuid.UUID `json:"revenue_account_id"`
	
	Description string `json:"description"`
	
	Quantity *float64 `json:"quantity"`
	
	UnitPrice float64 `json:"unit_price"`
	
	Amount float64 `json:"amount"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	TaxCode *string `json:"tax_code"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateCustomerInvoiceLinesRequest represents a request to create a customer_invoice_lines
type CreateCustomerInvoiceLinesRequest struct {
	
	CustomerInvoiceId uuid.UUID `json:"customer_invoice_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	RevenueAccountId uuid.UUID `json:"revenue_account_id" validate:"required"`
	
	Description string `json:"description" validate:"required"`
	
	Quantity *float64 `json:"quantity"`
	
	UnitPrice float64 `json:"unit_price" validate:"required"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	TaxCode *string `json:"tax_code"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateCustomerInvoiceLinesRequest) Validate() error {
	
	if r.CustomerInvoiceId == uuid.Nil {
		return fmt.Errorf("customer_invoice_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.RevenueAccountId == uuid.Nil {
		return fmt.Errorf("revenue_account_id is required")
	}
	
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	
	if r.UnitPrice == nil {
		return fmt.Errorf("unit_price is required")
	}
	
	if r.Amount == nil {
		return fmt.Errorf("amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerInvoiceLinesRequest represents a request to update a customer_invoice_lines
type UpdateCustomerInvoiceLinesRequest struct {
	
	CustomerInvoiceId *uuid.UUID `json:"customer_invoice_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	RevenueAccountId *uuid.UUID `json:"revenue_account_id,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty" validate:"omitempty,required"`
	
	Quantity *float64 `json:"quantity,omitempty"`
	
	UnitPrice *float64 `json:"unit_price,omitempty" validate:"omitempty,required"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	Department *string `json:"department,omitempty"`
	
	ProjectCode *string `json:"project_code,omitempty"`
	
	TaxCode *string `json:"tax_code,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomerInvoiceLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerInvoiceId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.RevenueAccountId != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Quantity != nil {
		hasUpdate = true
	}
	
	if r.UnitPrice != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.Department != nil {
		hasUpdate = true
	}
	
	if r.ProjectCode != nil {
		hasUpdate = true
	}
	
	if r.TaxCode != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CustomerInvoiceLinesListResponse represents a paginated list of customer_invoice_lines records
type CustomerInvoiceLinesListResponse struct {
	Items      []*CustomerInvoiceLinesResponse `json:"items"`
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
