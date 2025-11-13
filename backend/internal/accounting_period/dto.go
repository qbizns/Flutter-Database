package accounting_period

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountingPeriodsResponse represents a accounting_periods response
type AccountingPeriodsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	FiscalYearId uuid.UUID `json:"fiscal_year_id"`
	
	PeriodNumber int64 `json:"period_number"`
	
	PeriodName string `json:"period_name"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate time.Time `json:"end_date"`
	
	Status *string `json:"status"`
	
	ClosedBy *uuid.UUID `json:"closed_by"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateAccountingPeriodsRequest represents a request to create a accounting_periods
type CreateAccountingPeriodsRequest struct {
	
	FiscalYearId uuid.UUID `json:"fiscal_year_id" validate:"required"`
	
	PeriodNumber int64 `json:"period_number" validate:"required"`
	
	PeriodName string `json:"period_name" validate:"required"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate time.Time `json:"end_date" validate:"required"`
	
	// 	Status *string `json:"status"`
	
	ClosedBy *uuid.UUID `json:"closed_by"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateAccountingPeriodsRequest) Validate() error {
	
	if r.FiscalYearId == uuid.Nil {
		return fmt.Errorf("fiscal_year_id is required")
	}
	
	if r.PeriodNumber == 0 {
		return fmt.Errorf("period_number is required")
	}
	
	if r.PeriodName == "" {
		return fmt.Errorf("period_name is required")
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

// UpdateAccountingPeriodsRequest represents a request to update a accounting_periods
type UpdateAccountingPeriodsRequest struct {
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id,omitempty" validate:"omitempty,required"`
	
	PeriodNumber *int64 `json:"period_number,omitempty" validate:"omitempty,required"`
	
	PeriodName *string `json:"period_name,omitempty" validate:"omitempty,required"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty" validate:"omitempty,required"`
	
	// 	Status *string `json:"status,omitempty"`
	
	ClosedBy *uuid.UUID `json:"closed_by,omitempty"`
	
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAccountingPeriodsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.FiscalYearId != nil {
		hasUpdate = true
	}
	
	if r.PeriodNumber != nil {
		hasUpdate = true
	}
	
	if r.PeriodName != nil {
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
	
	if r.ClosedBy != nil {
		hasUpdate = true
	}
	
	if r.ClosedAt != nil {
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

// AccountingPeriodsListResponse represents a paginated list of accounting_periods records
type AccountingPeriodsListResponse struct {
	Items      []*AccountingPeriodsResponse `json:"items"`
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
