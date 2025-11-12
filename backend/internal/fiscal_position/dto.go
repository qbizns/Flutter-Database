package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FiscalPositionsResponse represents a fiscal_positions response
type FiscalPositionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PositionCode string `json:"position_code"`
	
	PositionName string `json:"position_name"`
	
	AutoApply *bool `json:"auto_apply"`
	
	CountryId *string `json:"country_id"`
	
	StateProvince *string `json:"state_province"`
	
	ZipPostalCodeRange *string `json:"zip_postal_code_range"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateFiscalPositionsRequest represents a request to create a fiscal_positions
type CreateFiscalPositionsRequest struct {
	
	PositionCode string `json:"position_code" validate:"required"`
	
	PositionName string `json:"position_name" validate:"required"`
	
	AutoApply *bool `json:"auto_apply"`
	
	CountryId *string `json:"country_id"`
	
	StateProvince *string `json:"state_province"`
	
	ZipPostalCodeRange *string `json:"zip_postal_code_range"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateFiscalPositionsRequest) Validate() error {
	
	if r.PositionCode == "" {
		return fmt.Errorf("position_code is required")
	}
	
	if r.PositionName == "" {
		return fmt.Errorf("position_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFiscalPositionsRequest represents a request to update a fiscal_positions
type UpdateFiscalPositionsRequest struct {
	
	PositionCode *string `json:"position_code,omitempty" validate:"omitempty,required"`
	
	PositionName *string `json:"position_name,omitempty" validate:"omitempty,required"`
	
	AutoApply *bool `json:"auto_apply,omitempty"`
	
	CountryId *string `json:"country_id,omitempty"`
	
	StateProvince *string `json:"state_province,omitempty"`
	
	ZipPostalCodeRange *string `json:"zip_postal_code_range,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFiscalPositionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PositionCode != nil {
		hasUpdate = true
	}
	
	if r.PositionName != nil {
		hasUpdate = true
	}
	
	if r.AutoApply != nil {
		hasUpdate = true
	}
	
	if r.CountryId != nil {
		hasUpdate = true
	}
	
	if r.StateProvince != nil {
		hasUpdate = true
	}
	
	if r.ZipPostalCodeRange != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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

// FiscalPositionsListResponse represents a paginated list of fiscal_positions records
type FiscalPositionsListResponse struct {
	Items      []*FiscalPositionsResponse `json:"items"`
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
