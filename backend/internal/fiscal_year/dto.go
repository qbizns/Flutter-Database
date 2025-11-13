package fiscal_year

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// FiscalYearsResponse represents a fiscal_years response
type FiscalYearsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	FiscalYear string `json:"fiscal_year"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate time.Time `json:"end_date"`
	
	Status *string `json:"status"`
	
	IsCurrent *bool `json:"is_current"`
	
	ClosedBy *uuid.UUID `json:"closed_by"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateFiscalYearsRequest represents a request to create a fiscal_years
type CreateFiscalYearsRequest struct {
	
	FiscalYear string `json:"fiscal_year" validate:"required"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate time.Time `json:"end_date" validate:"required"`
	
	// 	Status *string `json:"status"`
	
	IsCurrent *bool `json:"is_current"`
	
	ClosedBy *uuid.UUID `json:"closed_by"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateFiscalYearsRequest) Validate() error {
	
	if r.FiscalYear == "" {
		return fmt.Errorf("fiscal_year is required")
	}
	
	if r.StartDate.IsZero() {
		return fmt.Errorf("start_date is required")
	}
	
	if r.EndDate.IsZero() {
		return fmt.Errorf("end_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateFiscalYearsRequest represents a request to update a fiscal_years
type UpdateFiscalYearsRequest struct {
	
	FiscalYear *string `json:"fiscal_year,omitempty" validate:"omitempty,required"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty" validate:"omitempty,required"`
	
	// 	Status *string `json:"status,omitempty"`
	
	IsCurrent *bool `json:"is_current,omitempty"`
	
	ClosedBy *uuid.UUID `json:"closed_by,omitempty"`
	
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateFiscalYearsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FiscalYear != nil {
		hasUpdate = true
	}
	
	if r.StartDate != nil {
		hasUpdate = true
	}
	
	if r.EndDate != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.IsCurrent != nil {
		hasUpdate = true
	}
	
	if r.ClosedBy != nil {
		hasUpdate = true
	}
	
	if r.ClosedAt != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// FiscalYearsListResponse represents a paginated list of fiscal_years records
type FiscalYearsListResponse struct {
	Items      []*FiscalYearsResponse `json:"items"`
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
