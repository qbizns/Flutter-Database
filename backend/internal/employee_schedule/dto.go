package employee_schedule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmployeeSchedulesResponse represents a employee_schedules response
type EmployeeSchedulesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id"`
	
	ScheduleDate time.Time `json:"schedule_date"`
	
	ShiftType *string `json:"shift_type"`
	
	Position *string `json:"position"`
	
	ScheduledStartTime string `json:"scheduled_start_time"`
	
	ScheduledEndTime string `json:"scheduled_end_time"`
	
	BreakDurationMinutes *int64 `json:"break_duration_minutes"`
	
	Status *string `json:"status"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CancellationReason *string `json:"cancellation_reason"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'scheduled', *string `json:"'scheduled',"`
	
	'regular', *string `json:"'regular',"`
	
}

// CreateEmployeeSchedulesRequest represents a request to create a employee_schedules
type CreateEmployeeSchedulesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id" validate:"required"`
	
	ScheduleDate time.Time `json:"schedule_date" validate:"required"`
	
	ShiftType *string `json:"shift_type"`
	
	Position *string `json:"position"`
	
	ScheduledStartTime string `json:"scheduled_start_time" validate:"required"`
	
	ScheduledEndTime string `json:"scheduled_end_time" validate:"required"`
	
	BreakDurationMinutes *int64 `json:"break_duration_minutes"`
	
	Status *string `json:"status"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	Notes *string `json:"notes"`
	
	CancellationReason *string `json:"cancellation_reason"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'scheduled', *string `json:"'scheduled',"`
	
	'regular', *string `json:"'regular',"`
	
}

// Validate validates the create request
func (r *CreateEmployeeSchedulesRequest) Validate() error {
	
	if r.EmployeeId == uuid.Nil {
		return fmt.Errorf("employee_id is required")
	}
	
	if r.ScheduleDate == nil {
		return fmt.Errorf("schedule_date is required")
	}
	
	if r.ScheduledStartTime == "" {
		return fmt.Errorf("scheduled_start_time is required")
	}
	
	if r.ScheduledEndTime == "" {
		return fmt.Errorf("scheduled_end_time is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateEmployeeSchedulesRequest represents a request to update a employee_schedules
type UpdateEmployeeSchedulesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	EmployeeId *uuid.UUID `json:"employee_id,omitempty" validate:"omitempty,required"`
	
	ScheduleDate *time.Time `json:"schedule_date,omitempty" validate:"omitempty,required"`
	
	ShiftType *string `json:"shift_type,omitempty"`
	
	Position *string `json:"position,omitempty"`
	
	ScheduledStartTime *string `json:"scheduled_start_time,omitempty" validate:"omitempty,required"`
	
	ScheduledEndTime *string `json:"scheduled_end_time,omitempty" validate:"omitempty,required"`
	
	BreakDurationMinutes *int64 `json:"break_duration_minutes,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CancellationReason *string `json:"cancellation_reason,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'scheduled', *string `json:"'scheduled',,omitempty"`
	
	'regular', *string `json:"'regular',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateEmployeeSchedulesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.EmployeeId != nil {
		hasUpdate = true
	}
	
	if r.ScheduleDate != nil {
		hasUpdate = true
	}
	
	if r.ShiftType != nil {
		hasUpdate = true
	}
	
	if r.Position != nil {
		hasUpdate = true
	}
	
	if r.ScheduledStartTime != nil {
		hasUpdate = true
	}
	
	if r.ScheduledEndTime != nil {
		hasUpdate = true
	}
	
	if r.BreakDurationMinutes != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RequiresApproval != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CancellationReason != nil {
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
	
	if r.'regular', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// EmployeeSchedulesListResponse represents a paginated list of employee_schedules records
type EmployeeSchedulesListResponse struct {
	Items      []*EmployeeSchedulesResponse `json:"items"`
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
