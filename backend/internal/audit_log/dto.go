package audit_log

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AuditLogsResponse represents a audit_logs response
type AuditLogsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	Action string `json:"action"`
	
	ResourceType string `json:"resource_type"`
	
	ResourceId *uuid.UUID `json:"resource_id"`
	
	OldValues json.RawMessage `json:"old_values"`
	
	NewValues json.RawMessage `json:"new_values"`
	
	Changes json.RawMessage `json:"changes"`
	
	IpAddress *string `json:"ip_address"`
	
	UserAgent *string `json:"user_agent"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
}

// CreateAuditLogsRequest represents a request to create a audit_logs
type CreateAuditLogsRequest struct {
	
	UserId *uuid.UUID `json:"user_id"`
	
	Action string `json:"action" validate:"required"`
	
	ResourceType string `json:"resource_type" validate:"required"`
	
	ResourceId *uuid.UUID `json:"resource_id"`
	
	OldValues json.RawMessage `json:"old_values"`
	
	NewValues json.RawMessage `json:"new_values"`
	
	Changes json.RawMessage `json:"changes"`
	
	IpAddress *string `json:"ip_address"`
	
	UserAgent *string `json:"user_agent"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateAuditLogsRequest) Validate() error {
	
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}
	
	if r.ResourceType == "" {
		return fmt.Errorf("resource_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateAuditLogsRequest represents a request to update a audit_logs
type UpdateAuditLogsRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	Action *string `json:"action,omitempty" validate:"omitempty,required"`
	
	ResourceType *string `json:"resource_type,omitempty" validate:"omitempty,required"`
	
	ResourceId *uuid.UUID `json:"resource_id,omitempty"`
	
	OldValues *json.RawMessage `json:"old_values,omitempty"`
	
	NewValues *json.RawMessage `json:"new_values,omitempty"`
	
	Changes *json.RawMessage `json:"changes,omitempty"`
	
	IpAddress *string `json:"ip_address,omitempty"`
	
	UserAgent *string `json:"user_agent,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateAuditLogsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.Action != nil {
		hasUpdate = true
	}
	
	if r.ResourceType != nil {
		hasUpdate = true
	}
	
	if r.ResourceId != nil {
		hasUpdate = true
	}
	
	if r.OldValues != nil {
		hasUpdate = true
	}
	
	if r.NewValues != nil {
		hasUpdate = true
	}
	
	if r.Changes != nil {
		hasUpdate = true
	}
	
	if r.IpAddress != nil {
		hasUpdate = true
	}
	
	if r.UserAgent != nil {
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

// AuditLogsListResponse represents a paginated list of audit_logs records
type AuditLogsListResponse struct {
	Items      []*AuditLogsResponse `json:"items"`
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
