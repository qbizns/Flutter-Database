package uom_conversion

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UomConversionsResponse represents a uom_conversions response
type UomConversionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	FromUomId uuid.UUID `json:"from_uom_id"`
	
	ToUomId uuid.UUID `json:"to_uom_id"`
	
	ConversionFactor float64 `json:"conversion_factor"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateUomConversionsRequest represents a request to create a uom_conversions
type CreateUomConversionsRequest struct {
	
	FromUomId uuid.UUID `json:"from_uom_id" validate:"required"`
	
	ToUomId uuid.UUID `json:"to_uom_id" validate:"required"`
	
	ConversionFactor float64 `json:"conversion_factor" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateUomConversionsRequest) Validate() error {
	
	if r.FromUomId == uuid.Nil {
		return fmt.Errorf("from_uom_id is required")
	}
	
	if r.ToUomId == uuid.Nil {
		return fmt.Errorf("to_uom_id is required")
	}
	
	if r.ConversionFactor == nil {
		return fmt.Errorf("conversion_factor is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUomConversionsRequest represents a request to update a uom_conversions
type UpdateUomConversionsRequest struct {
	
	FromUomId *uuid.UUID `json:"from_uom_id,omitempty" validate:"omitempty,required"`
	
	ToUomId *uuid.UUID `json:"to_uom_id,omitempty" validate:"omitempty,required"`
	
	ConversionFactor *float64 `json:"conversion_factor,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateUomConversionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FromUomId != nil {
		hasUpdate = true
	}
	
	if r.ToUomId != nil {
		hasUpdate = true
	}
	
	if r.ConversionFactor != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UomConversionsListResponse represents a paginated list of uom_conversions records
type UomConversionsListResponse struct {
	Items      []*UomConversionsResponse `json:"items"`
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
