package time_clock_entry

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TimeClockEntriesResponse represents a time_clock_entries response
type TimeClockEntriesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id"`
	
	ScheduleId *uuid.UUID `json:"schedule_id"`
	
	EntryType string `json:"entry_type"`
	
	EntryTimestamp *time.Time `json:"entry_timestamp"`
	
	ScheduledTimestamp *time.Time `json:"scheduled_timestamp"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	GpsLocation json.RawMessage `json:"gps_location"`
	
	IpAddress *string `json:"ip_address"`
	
	IsLate *bool `json:"is_late"`
	
	IsEarly *bool `json:"is_early"`
	
	VarianceMinutes *int64 `json:"variance_minutes"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	IsManualEntry *bool `json:"is_manual_entry"`
	
	CorrectionNotes *string `json:"correction_notes"`
	
	PhotoUrl *string `json:"photo_url"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'clockIn', *string `json:"'clock_in',"`
	
	'mealStart', *string `json:"'meal_start',"`
	
}

// CreateTimeClockEntriesRequest represents a request to create a time_clock_entries
type CreateTimeClockEntriesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	EmployeeId uuid.UUID `json:"employee_id" validate:"required"`
	
	ScheduleId *uuid.UUID `json:"schedule_id"`
	
	EntryType string `json:"entry_type" validate:"required"`
	
	EntryTimestamp *time.Time `json:"entry_timestamp"`
	
	ScheduledTimestamp *time.Time `json:"scheduled_timestamp"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	// Duplicate removed: GpsLocation json.RawMessage `json:"gps_location"`
	
	IpAddress *string `json:"ip_address"`
	
	IsLate *bool `json:"is_late"`
	
	IsEarly *bool `json:"is_early"`
	
	VarianceMinutes *int64 `json:"variance_minutes"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	IsManualEntry *bool `json:"is_manual_entry"`
	
	CorrectionNotes *string `json:"correction_notes"`
	
	PhotoUrl *string `json:"photo_url" validate:"url"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'clockIn', *string `json:"'clock_in',"`
	
	'mealStart', *string `json:"'meal_start',"`
	
}

// Validate validates the create request
func (r *CreateTimeClockEntriesRequest) Validate() error {
	
	if r.EmployeeId == uuid.Nil {
		return fmt.Errorf("employee_id is required")
	}
	
	if r.EntryType == "" {
		return fmt.Errorf("entry_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateTimeClockEntriesRequest represents a request to update a time_clock_entries
type UpdateTimeClockEntriesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	EmployeeId *uuid.UUID `json:"employee_id,omitempty" validate:"omitempty,required"`
	
	ScheduleId *uuid.UUID `json:"schedule_id,omitempty"`
	
	EntryType *string `json:"entry_type,omitempty" validate:"omitempty,required"`
	
	EntryTimestamp *time.Time `json:"entry_timestamp,omitempty"`
	
	ScheduledTimestamp *time.Time `json:"scheduled_timestamp,omitempty"`
	
	DeviceId *uuid.UUID `json:"device_id,omitempty"`
	
	GpsLocation *json.RawMessage `json:"gps_location,omitempty"`
	
	IpAddress *string `json:"ip_address,omitempty"`
	
	IsLate *bool `json:"is_late,omitempty"`
	
	IsEarly *bool `json:"is_early,omitempty"`
	
	VarianceMinutes *int64 `json:"variance_minutes,omitempty"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	IsManualEntry *bool `json:"is_manual_entry,omitempty"`
	
	CorrectionNotes *string `json:"correction_notes,omitempty"`
	
	PhotoUrl *string `json:"photo_url,omitempty" validate:"omitempty,url"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'clockIn', *string `json:"'clock_in',,omitempty"`
	
	'mealStart', *string `json:"'meal_start',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateTimeClockEntriesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.EmployeeId != nil {
		hasUpdate = true
	}
	
	if r.ScheduleId != nil {
		hasUpdate = true
	}
	
	if r.EntryType != nil {
		hasUpdate = true
	}
	
	if r.EntryTimestamp != nil {
		hasUpdate = true
	}
	
	if r.ScheduledTimestamp != nil {
		hasUpdate = true
	}
	
	if r.DeviceId != nil {
		hasUpdate = true
	}
	
	if r.GpsLocation != nil {
		hasUpdate = true
	}
	
	if r.IpAddress != nil {
		hasUpdate = true
	}
	
	if r.IsLate != nil {
		hasUpdate = true
	}
	
	if r.IsEarly != nil {
		hasUpdate = true
	}
	
	if r.VarianceMinutes != nil {
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
	
	if r.IsManualEntry != nil {
		hasUpdate = true
	}
	
	if r.CorrectionNotes != nil {
		hasUpdate = true
	}
	
	if r.PhotoUrl != nil {
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
	
	if r.'clockIn', != nil {
		hasUpdate = true
	}
	
	if r.'mealStart', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// TimeClockEntriesListResponse represents a paginated list of time_clock_entries records
type TimeClockEntriesListResponse struct {
	Items      []*TimeClockEntriesResponse `json:"items"`
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
