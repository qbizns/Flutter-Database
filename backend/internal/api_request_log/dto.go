package api_request_log

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApiRequestLogsResponse represents a api_request_logs response
type ApiRequestLogsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	RequestId *string `json:"request_id"`
	
	Method string `json:"method"`
	
	Path string `json:"path"`
	
	QueryParams json.RawMessage `json:"query_params"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	ApiKeyId *uuid.UUID `json:"api_key_id"`
	
	RequestHeaders json.RawMessage `json:"request_headers"`
	
	RequestBody json.RawMessage `json:"request_body"`
	
	IpAddress *string `json:"ip_address"`
	
	UserAgent *string `json:"user_agent"`
	
	StatusCode int64 `json:"status_code"`
	
	ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ResponseBody json.RawMessage `json:"response_body"`
	
	DurationMs *int64 `json:"duration_ms"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorStack *string `json:"error_stack"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreateApiRequestLogsRequest represents a request to create a api_request_logs
type CreateApiRequestLogsRequest struct {
	
	RequestId *string `json:"request_id"`
	
	Method string `json:"method" validate:"required"`
	
	Path string `json:"path" validate:"required"`
	
	QueryParams json.RawMessage `json:"query_params"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	ApiKeyId *uuid.UUID `json:"api_key_id"`
	
	RequestHeaders json.RawMessage `json:"request_headers"`
	
	RequestBody json.RawMessage `json:"request_body"`
	
	IpAddress *string `json:"ip_address"`
	
	UserAgent *string `json:"user_agent"`
	
	StatusCode int64 `json:"status_code" validate:"required"`
	
	ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ResponseBody json.RawMessage `json:"response_body"`
	
	DurationMs *int64 `json:"duration_ms"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorStack *string `json:"error_stack"`
	
}

// Validate validates the create request
func (r *CreateApiRequestLogsRequest) Validate() error {
	
	if r.Method == "" {
		return fmt.Errorf("method is required")
	}
	
	if r.Path == "" {
		return fmt.Errorf("path is required")
	}
	
	if r.StatusCode == 0 {
		return fmt.Errorf("status_code is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateApiRequestLogsRequest represents a request to update a api_request_logs
type UpdateApiRequestLogsRequest struct {
	
	RequestId *string `json:"request_id,omitempty"`
	
	Method *string `json:"method,omitempty" validate:"omitempty,required"`
	
	Path *string `json:"path,omitempty" validate:"omitempty,required"`
	
	QueryParams *json.RawMessage `json:"query_params,omitempty"`
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	ApiKeyId *uuid.UUID `json:"api_key_id,omitempty"`
	
	RequestHeaders *json.RawMessage `json:"request_headers,omitempty"`
	
	RequestBody *json.RawMessage `json:"request_body,omitempty"`
	
	IpAddress *string `json:"ip_address,omitempty"`
	
	UserAgent *string `json:"user_agent,omitempty"`
	
	StatusCode *int64 `json:"status_code,omitempty" validate:"omitempty,required"`
	
	ResponseHeaders *json.RawMessage `json:"response_headers,omitempty"`
	
	ResponseBody *json.RawMessage `json:"response_body,omitempty"`
	
	DurationMs *int64 `json:"duration_ms,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	ErrorStack *string `json:"error_stack,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateApiRequestLogsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.RequestId != nil {
		hasUpdate = true
	}
	
	if r.Method != nil {
		hasUpdate = true
	}
	
	if r.Path != nil {
		hasUpdate = true
	}
	
	if r.QueryParams != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.ApiKeyId != nil {
		hasUpdate = true
	}
	
	if r.RequestHeaders != nil {
		hasUpdate = true
	}
	
	if r.RequestBody != nil {
		hasUpdate = true
	}
	
	if r.IpAddress != nil {
		hasUpdate = true
	}
	
	if r.UserAgent != nil {
		hasUpdate = true
	}
	
	if r.StatusCode != nil {
		hasUpdate = true
	}
	
	if r.ResponseHeaders != nil {
		hasUpdate = true
	}
	
	if r.ResponseBody != nil {
		hasUpdate = true
	}
	
	if r.DurationMs != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.ErrorStack != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ApiRequestLogsListResponse represents a paginated list of api_request_logs records
type ApiRequestLogsListResponse struct {
	Items      []*ApiRequestLogsResponse `json:"items"`
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
