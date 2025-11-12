package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingProfileDocumentsResponse represents a posting_profile_documents response
type PostingProfileDocumentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PostingProfileId uuid.UUID `json:"posting_profile_id"`
	
	PostingDocumentTypeId uuid.UUID `json:"posting_document_type_id"`
	
	IsActive bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePostingProfileDocumentsRequest represents a request to create a posting_profile_documents
type CreatePostingProfileDocumentsRequest struct {
	
	PostingProfileId uuid.UUID `json:"posting_profile_id" validate:"required"`
	
	PostingDocumentTypeId uuid.UUID `json:"posting_document_type_id" validate:"required"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreatePostingProfileDocumentsRequest) Validate() error {
	
	if r.PostingProfileId == uuid.Nil {
		return fmt.Errorf("posting_profile_id is required")
	}
	
	if r.PostingDocumentTypeId == uuid.Nil {
		return fmt.Errorf("posting_document_type_id is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingProfileDocumentsRequest represents a request to update a posting_profile_documents
type UpdatePostingProfileDocumentsRequest struct {
	
	PostingProfileId *uuid.UUID `json:"posting_profile_id,omitempty" validate:"omitempty,required"`
	
	PostingDocumentTypeId *uuid.UUID `json:"posting_document_type_id,omitempty" validate:"omitempty,required"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingProfileDocumentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PostingProfileId != nil {
		hasUpdate = true
	}
	
	if r.PostingDocumentTypeId != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PostingProfileDocumentsListResponse represents a paginated list of posting_profile_documents records
type PostingProfileDocumentsListResponse struct {
	Items      []*PostingProfileDocumentsResponse `json:"items"`
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
