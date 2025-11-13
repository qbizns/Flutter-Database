package webhook

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// WebhooksResponse represents a webhooks response
type WebhooksResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	WebhookName string `json:"webhook_name"`
	
	Url string `json:"url"`
	
	Secret *string `json:"secret"`
	
	Events string `json:"events"`
	
	HttpMethod *string `json:"http_method"`
	
	Headers json.RawMessage `json:"headers"`
	
	TimeoutSeconds *int64 `json:"timeout_seconds"`
	
	MaxRetries *int64 `json:"max_retries"`
	
	RetryBackoffSeconds *int64 `json:"retry_backoff_seconds"`
	
	IsActive *bool `json:"is_active"`
	
	IsVerified *bool `json:"is_verified"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries"`
	
	FailedDeliveries *int64 `json:"failed_deliveries"`
	
	LastDeliveryAt *time.Time `json:"last_delivery_at"`
	
	LastSuccessAt *time.Time `json:"last_success_at"`
	
	LastFailureAt *time.Time `json:"last_failure_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateWebhooksRequest represents a request to create a webhooks
type CreateWebhooksRequest struct {
	
	WebhookName string `json:"webhook_name" validate:"required"`
	
	Url string `json:"url" validate:"required,url"`
	
	Secret *string `json:"secret"`
	
	Events string `json:"events" validate:"required"`
	
	HttpMethod *string `json:"http_method"`
	
	// Duplicate removed: Headers json.RawMessage `json:"headers"`
	
	TimeoutSeconds *int64 `json:"timeout_seconds"`
	
	MaxRetries *int64 `json:"max_retries"`
	
	RetryBackoffSeconds *int64 `json:"retry_backoff_seconds"`
	
	IsActive *bool `json:"is_active"`
	
	IsVerified *bool `json:"is_verified"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries"`
	
	FailedDeliveries *int64 `json:"failed_deliveries"`
	
	LastDeliveryAt *time.Time `json:"last_delivery_at"`
	
	LastSuccessAt *time.Time `json:"last_success_at"`
	
	LastFailureAt *time.Time `json:"last_failure_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateWebhooksRequest) Validate() error {
	
	if r.WebhookName == "" {
		return fmt.Errorf("webhook_name is required")
	}
	
	if r.Url == "" {
		return fmt.Errorf("url is required")
	}
	
	if r.Events == "" {
		return fmt.Errorf("events is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateWebhooksRequest represents a request to update a webhooks
type UpdateWebhooksRequest struct {
	
	WebhookName *string `json:"webhook_name,omitempty" validate:"omitempty,required"`
	
	Url *string `json:"url,omitempty" validate:"omitempty,required,url"`
	
	Secret *string `json:"secret,omitempty"`
	
	Events *string `json:"events,omitempty" validate:"omitempty,required"`
	
	HttpMethod *string `json:"http_method,omitempty"`
	
	Headers *json.RawMessage `json:"headers,omitempty"`
	
	TimeoutSeconds *int64 `json:"timeout_seconds,omitempty"`
	
	MaxRetries *int64 `json:"max_retries,omitempty"`
	
	RetryBackoffSeconds *int64 `json:"retry_backoff_seconds,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsVerified *bool `json:"is_verified,omitempty"`
	
	TotalDeliveries *int64 `json:"total_deliveries,omitempty"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries,omitempty"`
	
	FailedDeliveries *int64 `json:"failed_deliveries,omitempty"`
	
	LastDeliveryAt *time.Time `json:"last_delivery_at,omitempty"`
	
	LastSuccessAt *time.Time `json:"last_success_at,omitempty"`
	
	LastFailureAt *time.Time `json:"last_failure_at,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateWebhooksRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.WebhookName != nil {
		hasUpdate = true
	}
	
	if r.Url != nil {
		hasUpdate = true
	}
	
	if r.Secret != nil {
		hasUpdate = true
	}
	
	if r.Events != nil {
		hasUpdate = true
	}
	
	if r.HttpMethod != nil {
		hasUpdate = true
	}
	
	if r.Headers != nil {
		hasUpdate = true
	}
	
	if r.TimeoutSeconds != nil {
		hasUpdate = true
	}
	
	if r.MaxRetries != nil {
		hasUpdate = true
	}
	
	if r.RetryBackoffSeconds != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsVerified != nil {
		hasUpdate = true
	}
	
	if r.TotalDeliveries != nil {
		hasUpdate = true
	}
	
	if r.SuccessfulDeliveries != nil {
		hasUpdate = true
	}
	
	if r.FailedDeliveries != nil {
		hasUpdate = true
	}
	
	if r.LastDeliveryAt != nil {
		hasUpdate = true
	}
	
	if r.LastSuccessAt != nil {
		hasUpdate = true
	}
	
	if r.LastFailureAt != nil {
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

// WebhooksListResponse represents a paginated list of webhooks records
type WebhooksListResponse struct {
	Items      []*WebhooksResponse `json:"items"`
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
