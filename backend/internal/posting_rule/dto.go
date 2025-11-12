package posting_rule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingRulesResponse represents a posting_rules response
type PostingRulesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PostingProfileDocumentId uuid.UUID `json:"posting_profile_document_id"`
	
	RuleCode string `json:"rule_code"`
	
	RuleName string `json:"rule_name"`
	
	Description *string `json:"description"`
	
	Event string `json:"event"`
	
	Level string `json:"level"`
	
	Priority int64 `json:"priority"`
	
	ConditionExpression *string `json:"condition_expression"`
	
	IsActive bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreatePostingRulesRequest represents a request to create a posting_rules
type CreatePostingRulesRequest struct {
	
	PostingProfileDocumentId uuid.UUID `json:"posting_profile_document_id" validate:"required"`
	
	RuleCode string `json:"rule_code" validate:"required"`
	
	RuleName string `json:"rule_name" validate:"required"`
	
	Description *string `json:"description"`
	
	Event string `json:"event" validate:"required"`
	
	Level string `json:"level" validate:"required"`
	
	Priority int64 `json:"priority" validate:"required"`
	
	ConditionExpression *string `json:"condition_expression"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePostingRulesRequest) Validate() error {
	
	if r.PostingProfileDocumentId == uuid.Nil {
		return fmt.Errorf("posting_profile_document_id is required")
	}
	
	if r.RuleCode == "" {
		return fmt.Errorf("rule_code is required")
	}
	
	if r.RuleName == "" {
		return fmt.Errorf("rule_name is required")
	}
	
	if r.Event == "" {
		return fmt.Errorf("event is required")
	}
	
	if r.Level == "" {
		return fmt.Errorf("level is required")
	}
	
	if r.Priority == 0 {
		return fmt.Errorf("priority is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingRulesRequest represents a request to update a posting_rules
type UpdatePostingRulesRequest struct {
	
	PostingProfileDocumentId *uuid.UUID `json:"posting_profile_document_id,omitempty" validate:"omitempty,required"`
	
	RuleCode *string `json:"rule_code,omitempty" validate:"omitempty,required"`
	
	RuleName *string `json:"rule_name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	Event *string `json:"event,omitempty" validate:"omitempty,required"`
	
	Level *string `json:"level,omitempty" validate:"omitempty,required"`
	
	Priority *int64 `json:"priority,omitempty" validate:"omitempty,required"`
	
	ConditionExpression *string `json:"condition_expression,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingRulesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PostingProfileDocumentId != nil {
		hasUpdate = true
	}
	
	if r.RuleCode != nil {
		hasUpdate = true
	}
	
	if r.RuleName != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Event != nil {
		hasUpdate = true
	}
	
	if r.Level != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.ConditionExpression != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PostingRulesListResponse represents a paginated list of posting_rules records
type PostingRulesListResponse struct {
	Items      []*PostingRulesResponse `json:"items"`
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
