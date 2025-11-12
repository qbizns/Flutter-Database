package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SystemHealthResponse represents a system_health response
type SystemHealthResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	Status *string `json:"status"`
	
	LastCheckAt *time.Time `json:"last_check_at"`
	
	LastSuccessAt *time.Time `json:"last_success_at"`
	
	LastFailureAt *time.Time `json:"last_failure_at"`
	
	MetricValue *float64 `json:"metric_value"`
	
	MetricUnit *string `json:"metric_unit"`
	
	ThresholdWarning *float64 `json:"threshold_warning"`
	
	ThresholdCritical *float64 `json:"threshold_critical"`
	
	Details json.RawMessage `json:"details"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateSystemHealthRequest represents a request to create a system_health
type CreateSystemHealthRequest struct {
	
	Status *string `json:"status"`
	
	LastCheckAt *time.Time `json:"last_check_at"`
	
	LastSuccessAt *time.Time `json:"last_success_at"`
	
	LastFailureAt *time.Time `json:"last_failure_at"`
	
	MetricValue *float64 `json:"metric_value"`
	
	MetricUnit *string `json:"metric_unit"`
	
	ThresholdWarning *float64 `json:"threshold_warning"`
	
	ThresholdCritical *float64 `json:"threshold_critical"`
	
	Details json.RawMessage `json:"details"`
	
}

// Validate validates the create request
func (r *CreateSystemHealthRequest) Validate() error {
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSystemHealthRequest represents a request to update a system_health
type UpdateSystemHealthRequest struct {
	
	Status *string `json:"status,omitempty"`
	
	LastCheckAt *time.Time `json:"last_check_at,omitempty"`
	
	LastSuccessAt *time.Time `json:"last_success_at,omitempty"`
	
	LastFailureAt *time.Time `json:"last_failure_at,omitempty"`
	
	MetricValue *float64 `json:"metric_value,omitempty"`
	
	MetricUnit *string `json:"metric_unit,omitempty"`
	
	ThresholdWarning *float64 `json:"threshold_warning,omitempty"`
	
	ThresholdCritical *float64 `json:"threshold_critical,omitempty"`
	
	Details *json.RawMessage `json:"details,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSystemHealthRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.LastCheckAt != nil {
		hasUpdate = true
	}
	
	if r.LastSuccessAt != nil {
		hasUpdate = true
	}
	
	if r.LastFailureAt != nil {
		hasUpdate = true
	}
	
	if r.MetricValue != nil {
		hasUpdate = true
	}
	
	if r.MetricUnit != nil {
		hasUpdate = true
	}
	
	if r.ThresholdWarning != nil {
		hasUpdate = true
	}
	
	if r.ThresholdCritical != nil {
		hasUpdate = true
	}
	
	if r.Details != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// SystemHealthListResponse represents a paginated list of system_health records
type SystemHealthListResponse struct {
	Items      []*SystemHealthResponse `json:"items"`
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
