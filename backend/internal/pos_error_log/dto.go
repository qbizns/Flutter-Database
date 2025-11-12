package pos_error_log

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PosErrorLogsResponse represents a pos_error_logs response
type PosErrorLogsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	ErrorLevel string `json:"error_level"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage string `json:"error_message"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	PosSessionId *uuid.UUID `json:"pos_session_id"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	StackTrace *string `json:"stack_trace"`
	
	RequestData json.RawMessage `json:"request_data"`
	
	ErrorData json.RawMessage `json:"error_data"`
	
	IsResolved *bool `json:"is_resolved"`
	
	ResolvedBy *uuid.UUID `json:"resolved_by"`
	
	ResolvedAt *time.Time `json:"resolved_at"`
	
	ResolutionNotes *string `json:"resolution_notes"`
	
	OccurredAt *time.Time `json:"occurred_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreatePosErrorLogsRequest represents a request to create a pos_error_logs
type CreatePosErrorLogsRequest struct {
	
	ErrorLevel string `json:"error_level" validate:"required"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage string `json:"error_message" validate:"required"`
	
	DeviceId *uuid.UUID `json:"device_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	PosSessionId *uuid.UUID `json:"pos_session_id"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	StackTrace *string `json:"stack_trace"`
	
	RequestData json.RawMessage `json:"request_data"`
	
	ErrorData json.RawMessage `json:"error_data"`
	
	IsResolved *bool `json:"is_resolved"`
	
	ResolvedBy *uuid.UUID `json:"resolved_by"`
	
	ResolvedAt *time.Time `json:"resolved_at"`
	
	ResolutionNotes *string `json:"resolution_notes"`
	
	OccurredAt *time.Time `json:"occurred_at"`
	
}

// Validate validates the create request
func (r *CreatePosErrorLogsRequest) Validate() error {
	
	if r.ErrorLevel == "" {
		return fmt.Errorf("error_level is required")
	}
	
	if r.ErrorMessage == "" {
		return fmt.Errorf("error_message is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePosErrorLogsRequest represents a request to update a pos_error_logs
type UpdatePosErrorLogsRequest struct {
	
	ErrorLevel *string `json:"error_level,omitempty" validate:"omitempty,required"`
	
	ErrorCode *string `json:"error_code,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty" validate:"omitempty,required"`
	
	DeviceId *uuid.UUID `json:"device_id,omitempty"`
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	PosSessionId *uuid.UUID `json:"pos_session_id,omitempty"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	StackTrace *string `json:"stack_trace,omitempty"`
	
	RequestData *json.RawMessage `json:"request_data,omitempty"`
	
	ErrorData *json.RawMessage `json:"error_data,omitempty"`
	
	IsResolved *bool `json:"is_resolved,omitempty"`
	
	ResolvedBy *uuid.UUID `json:"resolved_by,omitempty"`
	
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	
	ResolutionNotes *string `json:"resolution_notes,omitempty"`
	
	OccurredAt *time.Time `json:"occurred_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePosErrorLogsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ErrorLevel != nil {
		hasUpdate = true
	}
	
	if r.ErrorCode != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.DeviceId != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.PosSessionId != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.StackTrace != nil {
		hasUpdate = true
	}
	
	if r.RequestData != nil {
		hasUpdate = true
	}
	
	if r.ErrorData != nil {
		hasUpdate = true
	}
	
	if r.IsResolved != nil {
		hasUpdate = true
	}
	
	if r.ResolvedBy != nil {
		hasUpdate = true
	}
	
	if r.ResolvedAt != nil {
		hasUpdate = true
	}
	
	if r.ResolutionNotes != nil {
		hasUpdate = true
	}
	
	if r.OccurredAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PosErrorLogsListResponse represents a paginated list of pos_error_logs records
type PosErrorLogsListResponse struct {
	Items      []*PosErrorLogsResponse `json:"items"`
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
