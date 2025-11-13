package webhook_delivery

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WebhookDeliveriesResponse represents a webhook_deliveries response
type WebhookDeliveriesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	WebhookId uuid.UUID `json:"webhook_id"`
	
	EventType string `json:"event_type"`
	
	EventId uuid.UUID `json:"event_id"`
	
	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	RequestUrl string `json:"request_url"`
	
	RequestMethod string `json:"request_method"`
	
	RequestHeaders json.RawMessage `json:"request_headers"`
	
	RequestBody json.RawMessage `json:"request_body"`
	
	ResponseStatusCode *int64 `json:"response_status_code"`
	
	ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ResponseBody *string `json:"response_body"`
	
	AttemptNumber *int64 `json:"attempt_number"`
	
	DurationMs *int64 `json:"duration_ms"`
	
	NextRetryAt *time.Time `json:"next_retry_at"`
	
	ErrorMessage *string `json:"error_message"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeliveredAt *time.Time `json:"delivered_at"`
	
}

// CreateWebhookDeliveriesRequest represents a request to create a webhook_deliveries
type CreateWebhookDeliveriesRequest struct {
	
	WebhookId uuid.UUID `json:"webhook_id" validate:"required"`
	
	EventType string `json:"event_type" validate:"required"`
	
	EventId uuid.UUID `json:"event_id" validate:"required"`
	
	// 	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	RequestUrl string `json:"request_url" validate:"required,url"`
	
	RequestMethod string `json:"request_method" validate:"required"`
	
	// Duplicate removed: RequestHeaders json.RawMessage `json:"request_headers"`
	
	// Duplicate removed: RequestBody json.RawMessage `json:"request_body"`
	
	ResponseStatusCode *int64 `json:"response_status_code"`
	
	// Duplicate removed: ResponseHeaders json.RawMessage `json:"response_headers"`
	
	ResponseBody *string `json:"response_body"`
	
	AttemptNumber *int64 `json:"attempt_number"`
	
	DurationMs *int64 `json:"duration_ms"`
	
	NextRetryAt *time.Time `json:"next_retry_at"`
	
	ErrorMessage *string `json:"error_message"`
	
	DeliveredAt *time.Time `json:"delivered_at"`
	
}

// Validate validates the create request
func (r *CreateWebhookDeliveriesRequest) Validate() error {
	
	if r.WebhookId == uuid.Nil {
		return fmt.Errorf("webhook_id is required")
	}
	
	if r.EventType == "" {
		return fmt.Errorf("event_type is required")
	}
	
	if r.EventId == uuid.Nil {
		return fmt.Errorf("event_id is required")
	}
	
	if r.RequestUrl == "" {
		return fmt.Errorf("request_url is required")
	}
	
	if r.RequestMethod == "" {
		return fmt.Errorf("request_method is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateWebhookDeliveriesRequest represents a request to update a webhook_deliveries
type UpdateWebhookDeliveriesRequest struct {
	
	WebhookId *uuid.UUID `json:"webhook_id,omitempty" validate:"omitempty,required"`
	
	EventType *string `json:"event_type,omitempty" validate:"omitempty,required"`
	
	EventId *uuid.UUID `json:"event_id,omitempty" validate:"omitempty,required"`
	
	// 	Status *string `json:"status,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	RequestUrl *string `json:"request_url,omitempty" validate:"omitempty,required,url"`
	
	RequestMethod *string `json:"request_method,omitempty" validate:"omitempty,required"`
	
	RequestHeaders *json.RawMessage `json:"request_headers,omitempty"`
	
	RequestBody *json.RawMessage `json:"request_body,omitempty"`
	
	ResponseStatusCode *int64 `json:"response_status_code,omitempty"`
	
	ResponseHeaders *json.RawMessage `json:"response_headers,omitempty"`
	
	ResponseBody *string `json:"response_body,omitempty"`
	
	AttemptNumber *int64 `json:"attempt_number,omitempty"`
	
	DurationMs *int64 `json:"duration_ms,omitempty"`
	
	NextRetryAt *time.Time `json:"next_retry_at,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	DeliveredAt *time.Time `json:"delivered_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateWebhookDeliveriesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.WebhookId != nil {
		hasUpdate = true
	}
	
	if r.EventType != nil {
		hasUpdate = true
	}
	
	if r.EventId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.RequestUrl != nil {
		hasUpdate = true
	}
	
	if r.RequestMethod != nil {
		hasUpdate = true
	}
	
	if r.RequestHeaders != nil {
		hasUpdate = true
	}
	
	if r.RequestBody != nil {
		hasUpdate = true
	}
	
	if r.ResponseStatusCode != nil {
		hasUpdate = true
	}
	
	if r.ResponseHeaders != nil {
		hasUpdate = true
	}
	
	if r.ResponseBody != nil {
		hasUpdate = true
	}
	
	if r.AttemptNumber != nil {
		hasUpdate = true
	}
	
	if r.DurationMs != nil {
		hasUpdate = true
	}
	
	if r.NextRetryAt != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.DeliveredAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// WebhookDeliveriesListResponse represents a paginated list of webhook_deliveries records
type WebhookDeliveriesListResponse struct {
	Items      []*WebhookDeliveriesResponse `json:"items"`
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
