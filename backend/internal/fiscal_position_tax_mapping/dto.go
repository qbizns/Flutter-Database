package fiscal_position_tax_mapping

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FiscalPositionTaxMappingsResponse represents a fiscal_position_tax_mappings response
type FiscalPositionTaxMappingsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	FiscalPositionId uuid.UUID `json:"fiscal_position_id"`
	
	SourceTaxId uuid.UUID `json:"source_tax_id"`
	
	DestinationTaxId *uuid.UUID `json:"destination_tax_id"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateFiscalPositionTaxMappingsRequest represents a request to create a fiscal_position_tax_mappings
type CreateFiscalPositionTaxMappingsRequest struct {
	
	FiscalPositionId uuid.UUID `json:"fiscal_position_id" validate:"required"`
	
	SourceTaxId uuid.UUID `json:"source_tax_id" validate:"required"`
	
	DestinationTaxId *uuid.UUID `json:"destination_tax_id"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateFiscalPositionTaxMappingsRequest) Validate() error {
	
	if r.FiscalPositionId == uuid.Nil {
		return fmt.Errorf("fiscal_position_id is required")
	}
	
	if r.SourceTaxId == uuid.Nil {
		return fmt.Errorf("source_tax_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFiscalPositionTaxMappingsRequest represents a request to update a fiscal_position_tax_mappings
type UpdateFiscalPositionTaxMappingsRequest struct {
	
	FiscalPositionId *uuid.UUID `json:"fiscal_position_id,omitempty" validate:"omitempty,required"`
	
	SourceTaxId *uuid.UUID `json:"source_tax_id,omitempty" validate:"omitempty,required"`
	
	DestinationTaxId *uuid.UUID `json:"destination_tax_id,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFiscalPositionTaxMappingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FiscalPositionId != nil {
		hasUpdate = true
	}
	
	if r.SourceTaxId != nil {
		hasUpdate = true
	}
	
	if r.DestinationTaxId != nil {
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

// FiscalPositionTaxMappingsListResponse represents a paginated list of fiscal_position_tax_mappings records
type FiscalPositionTaxMappingsListResponse struct {
	Items      []*FiscalPositionTaxMappingsResponse `json:"items"`
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
