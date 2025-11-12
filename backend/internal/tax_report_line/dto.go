package tax_report_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TaxReportLinesResponse represents a tax_report_lines response
type TaxReportLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	TaxReportDefinitionId uuid.UUID `json:"tax_report_definition_id"`
	
	LineCode string `json:"line_code"`
	
	LineName string `json:"line_name"`
	
	Sequence *int64 `json:"sequence"`
	
	ParentLineId *uuid.UUID `json:"parent_line_id"`
	
	FormulaType *string `json:"formula_type"`
	
	Formula *string `json:"formula"`
	
	TaxGroupIds *uuid.UUID `json:"tax_group_ids"`
	
	AccountIds *uuid.UUID `json:"account_ids"`
	
	TaxIds *uuid.UUID `json:"tax_ids"`
	
	IsSubtotal *bool `json:"is_subtotal"`
	
	IsTotal *bool `json:"is_total"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateTaxReportLinesRequest represents a request to create a tax_report_lines
type CreateTaxReportLinesRequest struct {
	
	TaxReportDefinitionId uuid.UUID `json:"tax_report_definition_id" validate:"required"`
	
	LineCode string `json:"line_code" validate:"required"`
	
	LineName string `json:"line_name" validate:"required"`
	
	Sequence *int64 `json:"sequence"`
	
	ParentLineId *uuid.UUID `json:"parent_line_id"`
	
	FormulaType *string `json:"formula_type"`
	
	Formula *string `json:"formula"`
	
	TaxGroupIds *uuid.UUID `json:"tax_group_ids"`
	
	AccountIds *uuid.UUID `json:"account_ids"`
	
	TaxIds *uuid.UUID `json:"tax_ids"`
	
	IsSubtotal *bool `json:"is_subtotal"`
	
	IsTotal *bool `json:"is_total"`
	
	Notes *string `json:"notes"`
	
}

// Validate validates the create request
func (r *CreateTaxReportLinesRequest) Validate() error {
	
	if r.TaxReportDefinitionId == uuid.Nil {
		return fmt.Errorf("tax_report_definition_id is required")
	}
	
	if r.LineCode == "" {
		return fmt.Errorf("line_code is required")
	}
	
	if r.LineName == "" {
		return fmt.Errorf("line_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTaxReportLinesRequest represents a request to update a tax_report_lines
type UpdateTaxReportLinesRequest struct {
	
	TaxReportDefinitionId *uuid.UUID `json:"tax_report_definition_id,omitempty" validate:"omitempty,required"`
	
	LineCode *string `json:"line_code,omitempty" validate:"omitempty,required"`
	
	LineName *string `json:"line_name,omitempty" validate:"omitempty,required"`
	
	Sequence *int64 `json:"sequence,omitempty"`
	
	ParentLineId *uuid.UUID `json:"parent_line_id,omitempty"`
	
	FormulaType *string `json:"formula_type,omitempty"`
	
	Formula *string `json:"formula,omitempty"`
	
	TaxGroupIds *uuid.UUID `json:"tax_group_ids,omitempty"`
	
	AccountIds *uuid.UUID `json:"account_ids,omitempty"`
	
	TaxIds *uuid.UUID `json:"tax_ids,omitempty"`
	
	IsSubtotal *bool `json:"is_subtotal,omitempty"`
	
	IsTotal *bool `json:"is_total,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTaxReportLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TaxReportDefinitionId != nil {
		hasUpdate = true
	}
	
	if r.LineCode != nil {
		hasUpdate = true
	}
	
	if r.LineName != nil {
		hasUpdate = true
	}
	
	if r.Sequence != nil {
		hasUpdate = true
	}
	
	if r.ParentLineId != nil {
		hasUpdate = true
	}
	
	if r.FormulaType != nil {
		hasUpdate = true
	}
	
	if r.Formula != nil {
		hasUpdate = true
	}
	
	if r.TaxGroupIds != nil {
		hasUpdate = true
	}
	
	if r.AccountIds != nil {
		hasUpdate = true
	}
	
	if r.TaxIds != nil {
		hasUpdate = true
	}
	
	if r.IsSubtotal != nil {
		hasUpdate = true
	}
	
	if r.IsTotal != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// TaxReportLinesListResponse represents a paginated list of tax_report_lines records
type TaxReportLinesListResponse struct {
	Items      []*TaxReportLinesResponse `json:"items"`
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
