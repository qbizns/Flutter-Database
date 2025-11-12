package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BackgroundJobsResponse represents a background_jobs response
type BackgroundJobsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	JobType string `json:"job_type"`
	
	JobName string `json:"job_name"`
	
	QueueName *string `json:"queue_name"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Payload json.RawMessage `json:"payload"`
	
	Result json.RawMessage `json:"result"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	Priority *int64 `json:"priority"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	WorkerId *string `json:"worker_id"`
	
	ProcessingTimeout *int64 `json:"processing_timeout"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateBackgroundJobsRequest represents a request to create a background_jobs
type CreateBackgroundJobsRequest struct {
	
	JobType string `json:"job_type" validate:"required"`
	
	JobName string `json:"job_name" validate:"required"`
	
	QueueName *string `json:"queue_name"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Payload json.RawMessage `json:"payload" validate:"required"`
	
	Result json.RawMessage `json:"result"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	Priority *int64 `json:"priority"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	StartedAt *time.Time `json:"started_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	WorkerId *string `json:"worker_id"`
	
	ProcessingTimeout *int64 `json:"processing_timeout"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateBackgroundJobsRequest) Validate() error {
	
	if r.JobType == "" {
		return fmt.Errorf("job_type is required")
	}
	
	if r.JobName == "" {
		return fmt.Errorf("job_name is required")
	}
	
	if r.Payload == nil {
		return fmt.Errorf("payload is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBackgroundJobsRequest represents a request to update a background_jobs
type UpdateBackgroundJobsRequest struct {
	
	JobType *string `json:"job_type,omitempty" validate:"omitempty,required"`
	
	JobName *string `json:"job_name,omitempty" validate:"omitempty,required"`
	
	QueueName *string `json:"queue_name,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Payload *json.RawMessage `json:"payload,omitempty" validate:"omitempty,required"`
	
	Result *json.RawMessage `json:"result,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	ErrorDetails *json.RawMessage `json:"error_details,omitempty"`
	
	Attempts *int64 `json:"attempts,omitempty"`
	
	MaxAttempts *int64 `json:"max_attempts,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	
	StartedAt *time.Time `json:"started_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	FailedAt *time.Time `json:"failed_at,omitempty"`
	
	WorkerId *string `json:"worker_id,omitempty"`
	
	ProcessingTimeout *int64 `json:"processing_timeout,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBackgroundJobsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.JobType != nil {
		hasUpdate = true
	}
	
	if r.JobName != nil {
		hasUpdate = true
	}
	
	if r.QueueName != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Payload != nil {
		hasUpdate = true
	}
	
	if r.Result != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.ErrorDetails != nil {
		hasUpdate = true
	}
	
	if r.Attempts != nil {
		hasUpdate = true
	}
	
	if r.MaxAttempts != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.ScheduledAt != nil {
		hasUpdate = true
	}
	
	if r.StartedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	
	if r.FailedAt != nil {
		hasUpdate = true
	}
	
	if r.WorkerId != nil {
		hasUpdate = true
	}
	
	if r.ProcessingTimeout != nil {
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

// BackgroundJobsListResponse represents a paginated list of background_jobs records
type BackgroundJobsListResponse struct {
	Items      []*BackgroundJobsResponse `json:"items"`
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
