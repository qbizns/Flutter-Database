package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RateLimitsResponse represents a rate_limits response
type RateLimitsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	IdentifierType string `json:"identifier_type"`
	
	IdentifierValue string `json:"identifier_value"`
	
	EndpointPath *string `json:"endpoint_path"`
	
	HttpMethod *string `json:"http_method"`
	
	WindowStart time.Time `json:"window_start"`
	
	WindowDurationSeconds int64 `json:"window_duration_seconds"`
	
	RequestCount *int64 `json:"request_count"`
	
	AllowedCount int64 `json:"allowed_count"`
	
	IsBlocked *bool `json:"is_blocked"`
	
	BlockedUntil *time.Time `json:"blocked_until"`
	
	FirstRequestAt *time.Time `json:"first_request_at"`
	
	LastRequestAt *time.Time `json:"last_request_at"`
	
	IdentifierType, *string `json:"identifier_type,"`
	
}

// CreateRateLimitsRequest represents a request to create a rate_limits
type CreateRateLimitsRequest struct {
	
	IdentifierType string `json:"identifier_type" validate:"required"`
	
	IdentifierValue string `json:"identifier_value" validate:"required"`
	
	EndpointPath *string `json:"endpoint_path"`
	
	HttpMethod *string `json:"http_method"`
	
	WindowStart time.Time `json:"window_start" validate:"required"`
	
	WindowDurationSeconds int64 `json:"window_duration_seconds" validate:"required"`
	
	RequestCount *int64 `json:"request_count"`
	
	AllowedCount int64 `json:"allowed_count" validate:"required"`
	
	IsBlocked *bool `json:"is_blocked"`
	
	BlockedUntil *time.Time `json:"blocked_until"`
	
	FirstRequestAt *time.Time `json:"first_request_at"`
	
	LastRequestAt *time.Time `json:"last_request_at"`
	
	IdentifierType, *string `json:"identifier_type,"`
	
}

// Validate validates the create request
func (r *CreateRateLimitsRequest) Validate() error {
	
	if r.IdentifierType == "" {
		return fmt.Errorf("identifier_type is required")
	}
	
	if r.IdentifierValue == "" {
		return fmt.Errorf("identifier_value is required")
	}
	
	if r.WindowStart == nil {
		return fmt.Errorf("window_start is required")
	}
	
	if r.WindowDurationSeconds == 0 {
		return fmt.Errorf("window_duration_seconds is required")
	}
	
	if r.AllowedCount == 0 {
		return fmt.Errorf("allowed_count is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateRateLimitsRequest represents a request to update a rate_limits
type UpdateRateLimitsRequest struct {
	
	IdentifierType *string `json:"identifier_type,omitempty" validate:"omitempty,required"`
	
	IdentifierValue *string `json:"identifier_value,omitempty" validate:"omitempty,required"`
	
	EndpointPath *string `json:"endpoint_path,omitempty"`
	
	HttpMethod *string `json:"http_method,omitempty"`
	
	WindowStart *time.Time `json:"window_start,omitempty" validate:"omitempty,required"`
	
	WindowDurationSeconds *int64 `json:"window_duration_seconds,omitempty" validate:"omitempty,required"`
	
	RequestCount *int64 `json:"request_count,omitempty"`
	
	AllowedCount *int64 `json:"allowed_count,omitempty" validate:"omitempty,required"`
	
	IsBlocked *bool `json:"is_blocked,omitempty"`
	
	BlockedUntil *time.Time `json:"blocked_until,omitempty"`
	
	FirstRequestAt *time.Time `json:"first_request_at,omitempty"`
	
	LastRequestAt *time.Time `json:"last_request_at,omitempty"`
	
	IdentifierType, *string `json:"identifier_type,,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateRateLimitsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.IdentifierType != nil {
		hasUpdate = true
	}
	
	if r.IdentifierValue != nil {
		hasUpdate = true
	}
	
	if r.EndpointPath != nil {
		hasUpdate = true
	}
	
	if r.HttpMethod != nil {
		hasUpdate = true
	}
	
	if r.WindowStart != nil {
		hasUpdate = true
	}
	
	if r.WindowDurationSeconds != nil {
		hasUpdate = true
	}
	
	if r.RequestCount != nil {
		hasUpdate = true
	}
	
	if r.AllowedCount != nil {
		hasUpdate = true
	}
	
	if r.IsBlocked != nil {
		hasUpdate = true
	}
	
	if r.BlockedUntil != nil {
		hasUpdate = true
	}
	
	if r.FirstRequestAt != nil {
		hasUpdate = true
	}
	
	if r.LastRequestAt != nil {
		hasUpdate = true
	}
	
	if r.IdentifierType, != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// RateLimitsListResponse represents a paginated list of rate_limits records
type RateLimitsListResponse struct {
	Items      []*RateLimitsResponse `json:"items"`
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
