package shift

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ShiftsResponse represents a shifts response
type ShiftsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	ShiftNumber string `json:"shift_number"`
	
	StartTime time.Time `json:"start_time"`
	
	EndTime *time.Time `json:"end_time"`
	
	Status *string `json:"status"`
	
	OpeningCash *float64 `json:"opening_cash"`
	
	OpeningNotes *string `json:"opening_notes"`
	
	ExpectedCash *float64 `json:"expected_cash"`
	
	ActualCash *float64 `json:"actual_cash"`
	
	CashDifference *float64 `json:"cash_difference"`
	
	ClosingNotes *string `json:"closing_notes"`
	
	TotalSales *float64 `json:"total_sales"`
	
	TotalTransactions *int64 `json:"total_transactions"`
	
	TotalRefunds *float64 `json:"total_refunds"`
	
	TotalDiscounts *float64 `json:"total_discounts"`
	
	PaymentBreakdown json.RawMessage `json:"payment_breakdown"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	ClosedBy *uuid.UUID `json:"closed_by"`

	ClosedAt *time.Time `json:"closed_at"`

}

// CreateShiftsRequest represents a request to create a shifts
type CreateShiftsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	ShiftNumber string `json:"shift_number" validate:"required"`
	
	StartTime time.Time `json:"start_time" validate:"required"`
	
	EndTime *time.Time `json:"end_time"`
	
	Status *string `json:"status"`
	
	OpeningCash *float64 `json:"opening_cash"`
	
	OpeningNotes *string `json:"opening_notes"`
	
	ExpectedCash *float64 `json:"expected_cash"`
	
	ActualCash *float64 `json:"actual_cash"`
	
	CashDifference *float64 `json:"cash_difference"`
	
	ClosingNotes *string `json:"closing_notes"`
	
	TotalSales *float64 `json:"total_sales"`
	
	TotalTransactions *int64 `json:"total_transactions"`
	
	TotalRefunds *float64 `json:"total_refunds"`
	
	TotalDiscounts *float64 `json:"total_discounts"`
	
	PaymentBreakdown json.RawMessage `json:"payment_breakdown"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`

	ClosedBy *uuid.UUID `json:"closed_by"`

	ClosedAt *time.Time `json:"closed_at"`

	CreatedBy *uuid.UUID `json:"created_by"`

	UpdatedBy *uuid.UUID `json:"updated_by"`

}

// Validate validates the create request
func (r *CreateShiftsRequest) Validate() error {
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.ShiftNumber == "" {
		return fmt.Errorf("shift_number is required")
	}

	if r.StartTime.IsZero() {
		return fmt.Errorf("start_time is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateShiftsRequest represents a request to update a shifts
type UpdateShiftsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	ShiftNumber *string `json:"shift_number,omitempty" validate:"omitempty,required"`
	
	StartTime *time.Time `json:"start_time,omitempty" validate:"omitempty,required"`
	
	EndTime *time.Time `json:"end_time,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	OpeningCash *float64 `json:"opening_cash,omitempty"`
	
	OpeningNotes *string `json:"opening_notes,omitempty"`
	
	ExpectedCash *float64 `json:"expected_cash,omitempty"`
	
	ActualCash *float64 `json:"actual_cash,omitempty"`
	
	CashDifference *float64 `json:"cash_difference,omitempty"`
	
	ClosingNotes *string `json:"closing_notes,omitempty"`
	
	TotalSales *float64 `json:"total_sales,omitempty"`
	
	TotalTransactions *int64 `json:"total_transactions,omitempty"`
	
	TotalRefunds *float64 `json:"total_refunds,omitempty"`
	
	TotalDiscounts *float64 `json:"total_discounts,omitempty"`
	
	PaymentBreakdown *json.RawMessage `json:"payment_breakdown,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	ClosedBy *uuid.UUID `json:"closed_by,omitempty"`
	
	ClosedAt *time.Time `json:"closed_at,omitempty"`

	CreatedBy *uuid.UUID `json:"created_by,omitempty"`

	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`

}

// Validate validates the update request
func (r *UpdateShiftsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.ShiftNumber != nil {
		hasUpdate = true
	}
	
	if r.StartTime != nil {
		hasUpdate = true
	}
	
	if r.EndTime != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.OpeningCash != nil {
		hasUpdate = true
	}
	
	if r.OpeningNotes != nil {
		hasUpdate = true
	}
	
	if r.ExpectedCash != nil {
		hasUpdate = true
	}
	
	if r.ActualCash != nil {
		hasUpdate = true
	}
	
	if r.CashDifference != nil {
		hasUpdate = true
	}
	
	if r.ClosingNotes != nil {
		hasUpdate = true
	}
	
	if r.TotalSales != nil {
		hasUpdate = true
	}
	
	if r.TotalTransactions != nil {
		hasUpdate = true
	}
	
	if r.TotalRefunds != nil {
		hasUpdate = true
	}
	
	if r.TotalDiscounts != nil {
		hasUpdate = true
	}
	
	if r.PaymentBreakdown != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.ClosedBy != nil {
		hasUpdate = true
	}
	
	if r.ClosedAt != nil {
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

// ShiftsListResponse represents a paginated list of shifts records
type ShiftsListResponse struct {
	Items      []*ShiftsResponse `json:"items"`
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
