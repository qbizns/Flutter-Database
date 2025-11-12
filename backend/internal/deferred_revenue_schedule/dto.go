package deferred_revenue_schedule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeferredRevenueScheduleResponse represents a deferred_revenue_schedule response
type DeferredRevenueScheduleResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	ContractId uuid.UUID `json:"contract_id"`
	
	LineNumber int64 `json:"line_number"`
	
	RecognitionDate time.Time `json:"recognition_date"`
	
	RecognitionAmount float64 `json:"recognition_amount"`
	
	Status *string `json:"status"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateDeferredRevenueScheduleRequest represents a request to create a deferred_revenue_schedule
type CreateDeferredRevenueScheduleRequest struct {
	
	ContractId uuid.UUID `json:"contract_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	RecognitionDate time.Time `json:"recognition_date" validate:"required"`
	
	RecognitionAmount float64 `json:"recognition_amount" validate:"required"`
	
	Status *string `json:"status"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	PostedAt *time.Time `json:"posted_at"`
	
}

// Validate validates the create request
func (r *CreateDeferredRevenueScheduleRequest) Validate() error {
	
	if r.ContractId == uuid.Nil {
		return fmt.Errorf("contract_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.RecognitionDate == nil {
		return fmt.Errorf("recognition_date is required")
	}
	
	if r.RecognitionAmount == nil {
		return fmt.Errorf("recognition_amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeferredRevenueScheduleRequest represents a request to update a deferred_revenue_schedule
type UpdateDeferredRevenueScheduleRequest struct {
	
	ContractId *uuid.UUID `json:"contract_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	RecognitionDate *time.Time `json:"recognition_date,omitempty" validate:"omitempty,required"`
	
	RecognitionAmount *float64 `json:"recognition_amount,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	PostedAt *time.Time `json:"posted_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeferredRevenueScheduleRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ContractId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.RecognitionDate != nil {
		hasUpdate = true
	}
	
	if r.RecognitionAmount != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.PostedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DeferredRevenueScheduleListResponse represents a paginated list of deferred_revenue_schedule records
type DeferredRevenueScheduleListResponse struct {
	Items      []*DeferredRevenueScheduleResponse `json:"items"`
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
