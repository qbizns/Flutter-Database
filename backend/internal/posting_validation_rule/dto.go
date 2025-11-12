package posting_validation_rule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingValidationRulesResponse represents a posting_validation_rules response
type PostingValidationRulesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	DocumentTypeCode *string `json:"document_type_code"`
	
	Event *string `json:"event"`
	
	Target string `json:"target"`
	
	Code string `json:"code"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	Expression string `json:"expression"`
	
	Severity string `json:"severity"`
	
	IsBlocking bool `json:"is_blocking"`
	
	IsActive bool `json:"is_active"`
	
	MessageTemplate *string `json:"message_template"`
	
	Priority *int64 `json:"priority"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	COALESCE(organizationId, *uuid.UUID `json:"COALESCE(organization_id,"`
	
	COALESCE(documentTypeCode, *string `json:"COALESCE(document_type_code,"`
	
	COALESCE(event, *string `json:"COALESCE(event,"`
	
}

// CreatePostingValidationRulesRequest represents a request to create a posting_validation_rules
type CreatePostingValidationRulesRequest struct {
	
	DocumentTypeCode *string `json:"document_type_code"`
	
	Event *string `json:"event"`
	
	Target string `json:"target" validate:"required"`
	
	Code string `json:"code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	Expression string `json:"expression" validate:"required"`
	
	Severity string `json:"severity" validate:"required"`
	
	IsBlocking bool `json:"is_blocking" validate:"required"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	MessageTemplate *string `json:"message_template"`
	
	Priority *int64 `json:"priority"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	COALESCE(organizationId, *uuid.UUID `json:"COALESCE(organization_id,"`
	
	COALESCE(documentTypeCode, *string `json:"COALESCE(document_type_code,"`
	
	COALESCE(event, *string `json:"COALESCE(event,"`
	
}

// Validate validates the create request
func (r *CreatePostingValidationRulesRequest) Validate() error {
	
	if r.Target == "" {
		return fmt.Errorf("target is required")
	}
	
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.Expression == "" {
		return fmt.Errorf("expression is required")
	}
	
	if r.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	
	if r.IsBlocking == nil {
		return fmt.Errorf("is_blocking is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingValidationRulesRequest represents a request to update a posting_validation_rules
type UpdatePostingValidationRulesRequest struct {
	
	DocumentTypeCode *string `json:"document_type_code,omitempty"`
	
	Event *string `json:"event,omitempty"`
	
	Target *string `json:"target,omitempty" validate:"omitempty,required"`
	
	Code *string `json:"code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Expression *string `json:"expression,omitempty" validate:"omitempty,required"`
	
	Severity *string `json:"severity,omitempty" validate:"omitempty,required"`
	
	IsBlocking *bool `json:"is_blocking,omitempty" validate:"omitempty,required"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	MessageTemplate *string `json:"message_template,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	COALESCE(organizationId, *uuid.UUID `json:"COALESCE(organization_id,,omitempty"`
	
	COALESCE(documentTypeCode, *string `json:"COALESCE(document_type_code,,omitempty"`
	
	COALESCE(event, *string `json:"COALESCE(event,,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingValidationRulesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.DocumentTypeCode != nil {
		hasUpdate = true
	}
	
	if r.Event != nil {
		hasUpdate = true
	}
	
	if r.Target != nil {
		hasUpdate = true
	}
	
	if r.Code != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Expression != nil {
		hasUpdate = true
	}
	
	if r.Severity != nil {
		hasUpdate = true
	}
	
	if r.IsBlocking != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.MessageTemplate != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	
	if r.COALESCE(organizationId, != nil {
		hasUpdate = true
	}
	
	if r.COALESCE(documentTypeCode, != nil {
		hasUpdate = true
	}
	
	if r.COALESCE(event, != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PostingValidationRulesListResponse represents a paginated list of posting_validation_rules records
type PostingValidationRulesListResponse struct {
	Items      []*PostingValidationRulesResponse `json:"items"`
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
