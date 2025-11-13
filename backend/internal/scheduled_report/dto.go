package scheduled_report

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ScheduledReportsResponse represents a scheduled_reports response
type ScheduledReportsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ReportName string `json:"report_name"`
	
	ReportType string `json:"report_type"`
	
	ScheduleFrequency string `json:"schedule_frequency"`
	
	ScheduleDayOfWeek *int64 `json:"schedule_day_of_week"`
	
	ScheduleDayOfMonth *int64 `json:"schedule_day_of_month"`
	
	ScheduleTime string `json:"schedule_time"`
	
	ScheduleTimezone *string `json:"schedule_timezone"`
	
	ReportParameters json.RawMessage `json:"report_parameters"`
	
	DeliveryMethod *string `json:"delivery_method"`
	
	DeliveryRecipients *string `json:"delivery_recipients"`
	
	OutputFormat *string `json:"output_format"`
	
	IsActive *bool `json:"is_active"`
	
	LastRunAt *time.Time `json:"last_run_at"`
	
	LastRunStatus *string `json:"last_run_status"`
	
	NextRunAt *time.Time `json:"next_run_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateScheduledReportsRequest represents a request to create a scheduled_reports
type CreateScheduledReportsRequest struct {
	
	ReportName string `json:"report_name" validate:"required"`
	
	ReportType string `json:"report_type" validate:"required"`
	
	ScheduleFrequency string `json:"schedule_frequency" validate:"required"`
	
	ScheduleDayOfWeek *int64 `json:"schedule_day_of_week"`
	
	ScheduleDayOfMonth *int64 `json:"schedule_day_of_month"`
	
	ScheduleTime string `json:"schedule_time" validate:"required"`
	
	ScheduleTimezone *string `json:"schedule_timezone"`
	
	// Duplicate removed: ReportParameters json.RawMessage `json:"report_parameters"`
	
	DeliveryMethod *string `json:"delivery_method"`
	
	DeliveryRecipients *string `json:"delivery_recipients"`
	
	OutputFormat *string `json:"output_format"`
	
	IsActive *bool `json:"is_active"`
	
	LastRunAt *time.Time `json:"last_run_at"`
	
	// 	LastRunStatus *string `json:"last_run_status"`
	
	NextRunAt *time.Time `json:"next_run_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateScheduledReportsRequest) Validate() error {
	
	if r.ReportName == "" {
		return fmt.Errorf("report_name is required")
	}
	
	if r.ReportType == "" {
		return fmt.Errorf("report_type is required")
	}
	
	if r.ScheduleFrequency == "" {
		return fmt.Errorf("schedule_frequency is required")
	}
	
	if r.ScheduleTime == "" {
		return fmt.Errorf("schedule_time is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateScheduledReportsRequest represents a request to update a scheduled_reports
type UpdateScheduledReportsRequest struct {
	
	ReportName *string `json:"report_name,omitempty" validate:"omitempty,required"`
	
	ReportType *string `json:"report_type,omitempty" validate:"omitempty,required"`
	
	ScheduleFrequency *string `json:"schedule_frequency,omitempty" validate:"omitempty,required"`
	
	ScheduleDayOfWeek *int64 `json:"schedule_day_of_week,omitempty"`
	
	ScheduleDayOfMonth *int64 `json:"schedule_day_of_month,omitempty"`
	
	ScheduleTime *string `json:"schedule_time,omitempty" validate:"omitempty,required"`
	
	ScheduleTimezone *string `json:"schedule_timezone,omitempty"`
	
	ReportParameters *json.RawMessage `json:"report_parameters,omitempty"`
	
	DeliveryMethod *string `json:"delivery_method,omitempty"`
	
	DeliveryRecipients *string `json:"delivery_recipients,omitempty"`
	
	OutputFormat *string `json:"output_format,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	
	// 	LastRunStatus *string `json:"last_run_status,omitempty"`
	
	NextRunAt *time.Time `json:"next_run_at,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateScheduledReportsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ReportName != nil {
		hasUpdate = true
	}
	
	if r.ReportType != nil {
		hasUpdate = true
	}
	
	if r.ScheduleFrequency != nil {
		hasUpdate = true
	}
	
	if r.ScheduleDayOfWeek != nil {
		hasUpdate = true
	}
	
	if r.ScheduleDayOfMonth != nil {
		hasUpdate = true
	}
	
	if r.ScheduleTime != nil {
		hasUpdate = true
	}
	
	if r.ScheduleTimezone != nil {
		hasUpdate = true
	}
	
	if r.ReportParameters != nil {
		hasUpdate = true
	}
	
	if r.DeliveryMethod != nil {
		hasUpdate = true
	}
	
	if r.DeliveryRecipients != nil {
		hasUpdate = true
	}
	
	if r.OutputFormat != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.LastRunAt != nil {
		hasUpdate = true
	}
	
	if r.LastRunStatus != nil {
		hasUpdate = true
	}
	
	if r.NextRunAt != nil {
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

// ScheduledReportsListResponse represents a paginated list of scheduled_reports records
type ScheduledReportsListResponse struct {
	Items      []*ScheduledReportsResponse `json:"items"`
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
