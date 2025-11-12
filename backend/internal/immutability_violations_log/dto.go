package immutability_violations_log

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ImmutabilityViolationsLogResponse represents a immutability_violations_log response
type ImmutabilityViolationsLogResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	TableName string `json:"table_name"`
	
	RecordId *uuid.UUID `json:"record_id"`
	
	Operation string `json:"operation"`
	
	AttemptedBy *uuid.UUID `json:"attempted_by"`
	
	AttemptedAt time.Time `json:"attempted_at"`
	
	ErrorMessage *string `json:"error_message"`
	
	BlockedData json.RawMessage `json:"blocked_data"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// CreateImmutabilityViolationsLogRequest represents a request to create a immutability_violations_log
type CreateImmutabilityViolationsLogRequest struct {
	
	TableName string `json:"table_name" validate:"required"`
	
	RecordId *uuid.UUID `json:"record_id"`
	
	Operation string `json:"operation" validate:"required"`
	
	AttemptedBy *uuid.UUID `json:"attempted_by"`
	
	AttemptedAt time.Time `json:"attempted_at"`
	
	ErrorMessage *string `json:"error_message"`
	
	BlockedData json.RawMessage `json:"blocked_data"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateImmutabilityViolationsLogRequest) Validate() error {
	
	if r.TableName == "" {
		return fmt.Errorf("table_name is required")
	}
	
	if r.Operation == "" {
		return fmt.Errorf("operation is required")
	}
	
	if r.AttemptedAt == nil {
		return fmt.Errorf("attempted_at is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateImmutabilityViolationsLogRequest represents a request to update a immutability_violations_log
type UpdateImmutabilityViolationsLogRequest struct {
	
	TableName *string `json:"table_name,omitempty" validate:"omitempty,required"`
	
	RecordId *uuid.UUID `json:"record_id,omitempty"`
	
	Operation *string `json:"operation,omitempty" validate:"omitempty,required"`
	
	AttemptedBy *uuid.UUID `json:"attempted_by,omitempty"`
	
	AttemptedAt *time.Time `json:"attempted_at,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	BlockedData *json.RawMessage `json:"blocked_data,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateImmutabilityViolationsLogRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TableName != nil {
		hasUpdate = true
	}
	
	if r.RecordId != nil {
		hasUpdate = true
	}
	
	if r.Operation != nil {
		hasUpdate = true
	}
	
	if r.AttemptedBy != nil {
		hasUpdate = true
	}
	
	if r.AttemptedAt != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.BlockedData != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ImmutabilityViolationsLogListResponse represents a paginated list of immutability_violations_log records
type ImmutabilityViolationsLogListResponse struct {
	Items      []*ImmutabilityViolationsLogResponse `json:"items"`
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
