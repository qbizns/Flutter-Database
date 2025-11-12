package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DriverShiftsResponse represents a driver_shifts response
type DriverShiftsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	DriverId uuid.UUID `json:"driver_id"`
	
	ShiftDate time.Time `json:"shift_date"`
	
	ScheduledStartTime *string `json:"scheduled_start_time"`
	
	ScheduledEndTime *string `json:"scheduled_end_time"`
	
	ActualStartTime *time.Time `json:"actual_start_time"`
	
	ActualEndTime *time.Time `json:"actual_end_time"`
	
	Status *string `json:"status"`
	
	TotalBreakMinutes *int64 `json:"total_break_minutes"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	TotalDistanceKm *float64 `json:"total_distance_km"`
	
	TotalEarnings *float64 `json:"total_earnings"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'scheduled', *string `json:"'scheduled',"`
	
}

// CreateDriverShiftsRequest represents a request to create a driver_shifts
type CreateDriverShiftsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	DriverId uuid.UUID `json:"driver_id" validate:"required"`
	
	ShiftDate time.Time `json:"shift_date" validate:"required"`
	
	ScheduledStartTime *string `json:"scheduled_start_time"`
	
	ScheduledEndTime *string `json:"scheduled_end_time"`
	
	ActualStartTime *time.Time `json:"actual_start_time"`
	
	ActualEndTime *time.Time `json:"actual_end_time"`
	
	Status *string `json:"status"`
	
	TotalBreakMinutes *int64 `json:"total_break_minutes"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	TotalDistanceKm *float64 `json:"total_distance_km"`
	
	TotalEarnings *float64 `json:"total_earnings"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'scheduled', *string `json:"'scheduled',"`
	
}

// Validate validates the create request
func (r *CreateDriverShiftsRequest) Validate() error {
	
	if r.DriverId == uuid.Nil {
		return fmt.Errorf("driver_id is required")
	}
	
	if r.ShiftDate == nil {
		return fmt.Errorf("shift_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDriverShiftsRequest represents a request to update a driver_shifts
type UpdateDriverShiftsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	DriverId *uuid.UUID `json:"driver_id,omitempty" validate:"omitempty,required"`
	
	ShiftDate *time.Time `json:"shift_date,omitempty" validate:"omitempty,required"`
	
	ScheduledStartTime *string `json:"scheduled_start_time,omitempty"`
	
	ScheduledEndTime *string `json:"scheduled_end_time,omitempty"`
	
	ActualStartTime *time.Time `json:"actual_start_time,omitempty"`
	
	ActualEndTime *time.Time `json:"actual_end_time,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	TotalBreakMinutes *int64 `json:"total_break_minutes,omitempty"`
	
	TotalDeliveries *int64 `json:"total_deliveries,omitempty"`
	
	TotalDistanceKm *float64 `json:"total_distance_km,omitempty"`
	
	TotalEarnings *float64 `json:"total_earnings,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'scheduled', *string `json:"'scheduled',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDriverShiftsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.DriverId != nil {
		hasUpdate = true
	}
	
	if r.ShiftDate != nil {
		hasUpdate = true
	}
	
	if r.ScheduledStartTime != nil {
		hasUpdate = true
	}
	
	if r.ScheduledEndTime != nil {
		hasUpdate = true
	}
	
	if r.ActualStartTime != nil {
		hasUpdate = true
	}
	
	if r.ActualEndTime != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.TotalBreakMinutes != nil {
		hasUpdate = true
	}
	
	if r.TotalDeliveries != nil {
		hasUpdate = true
	}
	
	if r.TotalDistanceKm != nil {
		hasUpdate = true
	}
	
	if r.TotalEarnings != nil {
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
	
	if r.'scheduled', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DriverShiftsListResponse represents a paginated list of driver_shifts records
type DriverShiftsListResponse struct {
	Items      []*DriverShiftsResponse `json:"items"`
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
