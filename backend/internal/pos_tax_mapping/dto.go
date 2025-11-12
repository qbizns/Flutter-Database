package pos_tax_mapping

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PosTaxMappingsResponse represents a pos_tax_mappings response
type PosTaxMappingsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PosTaxCode *string `json:"pos_tax_code"`
	
	TaxCategoryCode *string `json:"tax_category_code"`
	
	PosTaxRate *float64 `json:"pos_tax_rate"`
	
	AccountingTaxId *uuid.UUID `json:"accounting_tax_id"`
	
	DefaultTaxAccountId *uuid.UUID `json:"default_tax_account_id"`
	
	DefaultTaxExpenseAccountId *uuid.UUID `json:"default_tax_expense_account_id"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	IsInclusive *bool `json:"is_inclusive"`
	
	AppliesToSales *bool `json:"applies_to_sales"`
	
	AppliesToPurchases *bool `json:"applies_to_purchases"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	EffectiveFrom *string `json:"effective_from"`
	
}

// CreatePosTaxMappingsRequest represents a request to create a pos_tax_mappings
type CreatePosTaxMappingsRequest struct {
	
	PosTaxCode *string `json:"pos_tax_code"`
	
	TaxCategoryCode *string `json:"tax_category_code"`
	
	PosTaxRate *float64 `json:"pos_tax_rate"`
	
	AccountingTaxId *uuid.UUID `json:"accounting_tax_id"`
	
	DefaultTaxAccountId *uuid.UUID `json:"default_tax_account_id"`
	
	DefaultTaxExpenseAccountId *uuid.UUID `json:"default_tax_expense_account_id"`
	
	IsDefault *bool `json:"is_default"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	IsInclusive *bool `json:"is_inclusive"`
	
	AppliesToSales *bool `json:"applies_to_sales"`
	
	AppliesToPurchases *bool `json:"applies_to_purchases"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	EffectiveFrom *string `json:"effective_from"`
	
}

// Validate validates the create request
func (r *CreatePosTaxMappingsRequest) Validate() error {
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePosTaxMappingsRequest represents a request to update a pos_tax_mappings
type UpdatePosTaxMappingsRequest struct {
	
	PosTaxCode *string `json:"pos_tax_code,omitempty"`
	
	TaxCategoryCode *string `json:"tax_category_code,omitempty"`
	
	PosTaxRate *float64 `json:"pos_tax_rate,omitempty"`
	
	AccountingTaxId *uuid.UUID `json:"accounting_tax_id,omitempty"`
	
	DefaultTaxAccountId *uuid.UUID `json:"default_tax_account_id,omitempty"`
	
	DefaultTaxExpenseAccountId *uuid.UUID `json:"default_tax_expense_account_id,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	IsInclusive *bool `json:"is_inclusive,omitempty"`
	
	AppliesToSales *bool `json:"applies_to_sales,omitempty"`
	
	AppliesToPurchases *bool `json:"applies_to_purchases,omitempty"`
	
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	EffectiveFrom *string `json:"effective_from,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePosTaxMappingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PosTaxCode != nil {
		hasUpdate = true
	}
	
	if r.TaxCategoryCode != nil {
		hasUpdate = true
	}
	
	if r.PosTaxRate != nil {
		hasUpdate = true
	}
	
	if r.AccountingTaxId != nil {
		hasUpdate = true
	}
	
	if r.DefaultTaxAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultTaxExpenseAccountId != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.IsInclusive != nil {
		hasUpdate = true
	}
	
	if r.AppliesToSales != nil {
		hasUpdate = true
	}
	
	if r.AppliesToPurchases != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	
	if r.EffectiveTo != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PosTaxMappingsListResponse represents a paginated list of pos_tax_mappings records
type PosTaxMappingsListResponse struct {
	Items      []*PosTaxMappingsResponse `json:"items"`
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
