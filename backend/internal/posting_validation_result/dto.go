package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingValidationResultsResponse represents a posting_validation_results response
type PostingValidationResultsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	DocumentTypeCode string `json:"document_type_code"`
	
	DocumentId uuid.UUID `json:"document_id"`
	
	Event string `json:"event"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	ValidationRuleId *uuid.UUID `json:"validation_rule_id"`
	
	Severity string `json:"severity"`
	
	MessageCode string `json:"message_code"`
	
	Message string `json:"message"`
	
	IsBlocking bool `json:"is_blocking"`
	
	Context json.RawMessage `json:"context"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreatePostingValidationResultsRequest represents a request to create a posting_validation_results
type CreatePostingValidationResultsRequest struct {
	
	DocumentTypeCode string `json:"document_type_code" validate:"required"`
	
	DocumentId uuid.UUID `json:"document_id" validate:"required"`
	
	Event string `json:"event" validate:"required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	ValidationRuleId *uuid.UUID `json:"validation_rule_id"`
	
	Severity string `json:"severity" validate:"required"`
	
	MessageCode string `json:"message_code" validate:"required"`
	
	Message string `json:"message" validate:"required"`
	
	IsBlocking bool `json:"is_blocking" validate:"required"`
	
	Context json.RawMessage `json:"context"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreatePostingValidationResultsRequest) Validate() error {
	
	if r.DocumentTypeCode == "" {
		return fmt.Errorf("document_type_code is required")
	}
	
	if r.DocumentId == uuid.Nil {
		return fmt.Errorf("document_id is required")
	}
	
	if r.Event == "" {
		return fmt.Errorf("event is required")
	}
	
	if r.Severity == "" {
		return fmt.Errorf("severity is required")
	}
	
	if r.MessageCode == "" {
		return fmt.Errorf("message_code is required")
	}
	
	if r.Message == "" {
		return fmt.Errorf("message is required")
	}
	
	if r.IsBlocking == nil {
		return fmt.Errorf("is_blocking is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingValidationResultsRequest represents a request to update a posting_validation_results
type UpdatePostingValidationResultsRequest struct {
	
	DocumentTypeCode *string `json:"document_type_code,omitempty" validate:"omitempty,required"`
	
	DocumentId *uuid.UUID `json:"document_id,omitempty" validate:"omitempty,required"`
	
	Event *string `json:"event,omitempty" validate:"omitempty,required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	ValidationRuleId *uuid.UUID `json:"validation_rule_id,omitempty"`
	
	Severity *string `json:"severity,omitempty" validate:"omitempty,required"`
	
	MessageCode *string `json:"message_code,omitempty" validate:"omitempty,required"`
	
	Message *string `json:"message,omitempty" validate:"omitempty,required"`
	
	IsBlocking *bool `json:"is_blocking,omitempty" validate:"omitempty,required"`
	
	Context *json.RawMessage `json:"context,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingValidationResultsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.DocumentTypeCode != nil {
		hasUpdate = true
	}
	
	if r.DocumentId != nil {
		hasUpdate = true
	}
	
	if r.Event != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.ValidationRuleId != nil {
		hasUpdate = true
	}
	
	if r.Severity != nil {
		hasUpdate = true
	}
	
	if r.MessageCode != nil {
		hasUpdate = true
	}
	
	if r.Message != nil {
		hasUpdate = true
	}
	
	if r.IsBlocking != nil {
		hasUpdate = true
	}
	
	if r.Context != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
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

// PostingValidationResultsListResponse represents a paginated list of posting_validation_results records
type PostingValidationResultsListResponse struct {
	Items      []*PostingValidationResultsResponse `json:"items"`
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
