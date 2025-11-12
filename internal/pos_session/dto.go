package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PosSessionsResponse represents a pos_sessions response
type PosSessionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SessionNumber string `json:"session_number"`
	
	SessionName *string `json:"session_name"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	OpenedAt time.Time `json:"opened_at"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	OpeningCash *float64 `json:"opening_cash"`
	
	OpeningCard *float64 `json:"opening_card"`
	
	OpeningOther *float64 `json:"opening_other"`
	
	ExpectedCash *float64 `json:"expected_cash"`
	
	ExpectedCard *float64 `json:"expected_card"`
	
	ExpectedOther *float64 `json:"expected_other"`
	
	CountedCash *float64 `json:"counted_cash"`
	
	CountedCard *float64 `json:"counted_card"`
	
	CountedOther *float64 `json:"counted_other"`
	
	DifferenceCash *float64 `json:"difference_cash"`
	
	DifferenceCard *float64 `json:"difference_card"`
	
	DifferenceOther *float64 `json:"difference_other"`
	
	Status *string `json:"status"`
	
	ZReportNumber *string `json:"z_report_number"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePosSessionsRequest represents a request to create a pos_sessions
type CreatePosSessionsRequest struct {
	
	SessionNumber string `json:"session_number" validate:"required"`
	
	SessionName *string `json:"session_name"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	ShiftId *uuid.UUID `json:"shift_id"`
	
	OpenedAt time.Time `json:"opened_at"`
	
	ClosedAt *time.Time `json:"closed_at"`
	
	OpeningCash *float64 `json:"opening_cash"`
	
	OpeningCard *float64 `json:"opening_card"`
	
	OpeningOther *float64 `json:"opening_other"`
	
	ExpectedCash *float64 `json:"expected_cash"`
	
	ExpectedCard *float64 `json:"expected_card"`
	
	ExpectedOther *float64 `json:"expected_other"`
	
	CountedCash *float64 `json:"counted_cash"`
	
	CountedCard *float64 `json:"counted_card"`
	
	CountedOther *float64 `json:"counted_other"`
	
	DifferenceCash *float64 `json:"difference_cash"`
	
	DifferenceCard *float64 `json:"difference_card"`
	
	DifferenceOther *float64 `json:"difference_other"`
	
	Status *string `json:"status"`
	
	ZReportNumber *string `json:"z_report_number"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePosSessionsRequest) Validate() error {
	
	if r.SessionNumber == "" {
		return fmt.Errorf("session_number is required")
	}
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.OpenedAt == nil {
		return fmt.Errorf("opened_at is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePosSessionsRequest represents a request to update a pos_sessions
type UpdatePosSessionsRequest struct {
	
	SessionNumber *string `json:"session_number,omitempty" validate:"omitempty,required"`
	
	SessionName *string `json:"session_name,omitempty"`
	
	DeviceId *uuid.UUID `json:"device_id,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	ShiftId *uuid.UUID `json:"shift_id,omitempty"`
	
	OpenedAt *time.Time `json:"opened_at,omitempty"`
	
	ClosedAt *time.Time `json:"closed_at,omitempty"`
	
	OpeningCash *float64 `json:"opening_cash,omitempty"`
	
	OpeningCard *float64 `json:"opening_card,omitempty"`
	
	OpeningOther *float64 `json:"opening_other,omitempty"`
	
	ExpectedCash *float64 `json:"expected_cash,omitempty"`
	
	ExpectedCard *float64 `json:"expected_card,omitempty"`
	
	ExpectedOther *float64 `json:"expected_other,omitempty"`
	
	CountedCash *float64 `json:"counted_cash,omitempty"`
	
	CountedCard *float64 `json:"counted_card,omitempty"`
	
	CountedOther *float64 `json:"counted_other,omitempty"`
	
	DifferenceCash *float64 `json:"difference_cash,omitempty"`
	
	DifferenceCard *float64 `json:"difference_card,omitempty"`
	
	DifferenceOther *float64 `json:"difference_other,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	ZReportNumber *string `json:"z_report_number,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePosSessionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SessionNumber != nil {
		hasUpdate = true
	}
	
	if r.SessionName != nil {
		hasUpdate = true
	}
	
	if r.DeviceId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.ShiftId != nil {
		hasUpdate = true
	}
	
	if r.OpenedAt != nil {
		hasUpdate = true
	}
	
	if r.ClosedAt != nil {
		hasUpdate = true
	}
	
	if r.OpeningCash != nil {
		hasUpdate = true
	}
	
	if r.OpeningCard != nil {
		hasUpdate = true
	}
	
	if r.OpeningOther != nil {
		hasUpdate = true
	}
	
	if r.ExpectedCash != nil {
		hasUpdate = true
	}
	
	if r.ExpectedCard != nil {
		hasUpdate = true
	}
	
	if r.ExpectedOther != nil {
		hasUpdate = true
	}
	
	if r.CountedCash != nil {
		hasUpdate = true
	}
	
	if r.CountedCard != nil {
		hasUpdate = true
	}
	
	if r.CountedOther != nil {
		hasUpdate = true
	}
	
	if r.DifferenceCash != nil {
		hasUpdate = true
	}
	
	if r.DifferenceCard != nil {
		hasUpdate = true
	}
	
	if r.DifferenceOther != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.ZReportNumber != nil {
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

// PosSessionsListResponse represents a paginated list of pos_sessions records
type PosSessionsListResponse struct {
	Items      []*PosSessionsResponse `json:"items"`
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
