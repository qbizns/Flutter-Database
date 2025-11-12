package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UnitsOfMeasureResponse represents a units_of_measure response
type UnitsOfMeasureResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	UomCode string `json:"uom_code"`
	
	UomName string `json:"uom_name"`
	
	UomType string `json:"uom_type"`
	
	IsBaseUnit *bool `json:"is_base_unit"`
	
	IsActive *bool `json:"is_active"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateUnitsOfMeasureRequest represents a request to create a units_of_measure
type CreateUnitsOfMeasureRequest struct {
	
	UomCode string `json:"uom_code" validate:"required"`
	
	UomName string `json:"uom_name" validate:"required"`
	
	UomType string `json:"uom_type" validate:"required"`
	
	IsBaseUnit *bool `json:"is_base_unit"`
	
	IsActive *bool `json:"is_active"`
	
}

// Validate validates the create request
func (r *CreateUnitsOfMeasureRequest) Validate() error {
	
	if r.UomCode == "" {
		return fmt.Errorf("uom_code is required")
	}
	
	if r.UomName == "" {
		return fmt.Errorf("uom_name is required")
	}
	
	if r.UomType == "" {
		return fmt.Errorf("uom_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUnitsOfMeasureRequest represents a request to update a units_of_measure
type UpdateUnitsOfMeasureRequest struct {
	
	UomCode *string `json:"uom_code,omitempty" validate:"omitempty,required"`
	
	UomName *string `json:"uom_name,omitempty" validate:"omitempty,required"`
	
	UomType *string `json:"uom_type,omitempty" validate:"omitempty,required"`
	
	IsBaseUnit *bool `json:"is_base_unit,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateUnitsOfMeasureRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UomCode != nil {
		hasUpdate = true
	}
	
	if r.UomName != nil {
		hasUpdate = true
	}
	
	if r.UomType != nil {
		hasUpdate = true
	}
	
	if r.IsBaseUnit != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UnitsOfMeasureListResponse represents a paginated list of units_of_measure records
type UnitsOfMeasureListResponse struct {
	Items      []*UnitsOfMeasureResponse `json:"items"`
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
