package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TaxReportDefinitionsResponse represents a tax_report_definitions response
type TaxReportDefinitionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	LocalizationPackageId *uuid.UUID `json:"localization_package_id"`
	
	ReportCode string `json:"report_code"`
	
	ReportName string `json:"report_name"`
	
	Jurisdiction *string `json:"jurisdiction"`
	
	Authority *string `json:"authority"`
	
	ReportFrequency *string `json:"report_frequency"`
	
	Version *string `json:"version"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateTaxReportDefinitionsRequest represents a request to create a tax_report_definitions
type CreateTaxReportDefinitionsRequest struct {
	
	LocalizationPackageId *uuid.UUID `json:"localization_package_id"`
	
	ReportCode string `json:"report_code" validate:"required"`
	
	ReportName string `json:"report_name" validate:"required"`
	
	Jurisdiction *string `json:"jurisdiction"`
	
	Authority *string `json:"authority"`
	
	ReportFrequency *string `json:"report_frequency"`
	
	Version *string `json:"version"`
	
	EffectiveFrom *time.Time `json:"effective_from"`
	
	EffectiveTo *time.Time `json:"effective_to"`
	
	IsActive *bool `json:"is_active"`
	
	Description *string `json:"description"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateTaxReportDefinitionsRequest) Validate() error {
	
	if r.ReportCode == "" {
		return fmt.Errorf("report_code is required")
	}
	
	if r.ReportName == "" {
		return fmt.Errorf("report_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTaxReportDefinitionsRequest represents a request to update a tax_report_definitions
type UpdateTaxReportDefinitionsRequest struct {
	
	LocalizationPackageId *uuid.UUID `json:"localization_package_id,omitempty"`
	
	ReportCode *string `json:"report_code,omitempty" validate:"omitempty,required"`
	
	ReportName *string `json:"report_name,omitempty" validate:"omitempty,required"`
	
	Jurisdiction *string `json:"jurisdiction,omitempty"`
	
	Authority *string `json:"authority,omitempty"`
	
	ReportFrequency *string `json:"report_frequency,omitempty"`
	
	Version *string `json:"version,omitempty"`
	
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	
	EffectiveTo *time.Time `json:"effective_to,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTaxReportDefinitionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocalizationPackageId != nil {
		hasUpdate = true
	}
	
	if r.ReportCode != nil {
		hasUpdate = true
	}
	
	if r.ReportName != nil {
		hasUpdate = true
	}
	
	if r.Jurisdiction != nil {
		hasUpdate = true
	}
	
	if r.Authority != nil {
		hasUpdate = true
	}
	
	if r.ReportFrequency != nil {
		hasUpdate = true
	}
	
	if r.Version != nil {
		hasUpdate = true
	}
	
	if r.EffectiveFrom != nil {
		hasUpdate = true
	}
	
	if r.EffectiveTo != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// TaxReportDefinitionsListResponse represents a paginated list of tax_report_definitions records
type TaxReportDefinitionsListResponse struct {
	Items      []*TaxReportDefinitionsResponse `json:"items"`
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
