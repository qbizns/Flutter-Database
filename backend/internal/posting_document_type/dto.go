package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingDocumentTypesResponse represents a posting_document_types response
type PostingDocumentTypesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	Code string `json:"code"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	SourceSchema string `json:"source_schema"`
	
	SourceTable string `json:"source_table"`
	
	SourcePkColumn string `json:"source_pk_column"`
	
	Category *string `json:"category"`
	
	IsActive bool `json:"is_active"`
	
	IsSystem bool `json:"is_system"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePostingDocumentTypesRequest represents a request to create a posting_document_types
type CreatePostingDocumentTypesRequest struct {
	
	Code string `json:"code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	SourceSchema string `json:"source_schema" validate:"required"`
	
	SourceTable string `json:"source_table" validate:"required"`
	
	SourcePkColumn string `json:"source_pk_column" validate:"required"`
	
	Category *string `json:"category"`
	
	IsActive bool `json:"is_active" validate:"required"`
	
	IsSystem bool `json:"is_system" validate:"required"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreatePostingDocumentTypesRequest) Validate() error {
	
	if r.Code == "" {
		return fmt.Errorf("code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.SourceSchema == "" {
		return fmt.Errorf("source_schema is required")
	}
	
	if r.SourceTable == "" {
		return fmt.Errorf("source_table is required")
	}
	
	if r.SourcePkColumn == "" {
		return fmt.Errorf("source_pk_column is required")
	}
	
	if r.IsActive == nil {
		return fmt.Errorf("is_active is required")
	}
	
	if r.IsSystem == nil {
		return fmt.Errorf("is_system is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingDocumentTypesRequest represents a request to update a posting_document_types
type UpdatePostingDocumentTypesRequest struct {
	
	Code *string `json:"code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	SourceSchema *string `json:"source_schema,omitempty" validate:"omitempty,required"`
	
	SourceTable *string `json:"source_table,omitempty" validate:"omitempty,required"`
	
	SourcePkColumn *string `json:"source_pk_column,omitempty" validate:"omitempty,required"`
	
	Category *string `json:"category,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty" validate:"omitempty,required"`
	
	IsSystem *bool `json:"is_system,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingDocumentTypesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.Code != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.SourceSchema != nil {
		hasUpdate = true
	}
	
	if r.SourceTable != nil {
		hasUpdate = true
	}
	
	if r.SourcePkColumn != nil {
		hasUpdate = true
	}
	
	if r.Category != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsSystem != nil {
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

// PostingDocumentTypesListResponse represents a paginated list of posting_document_types records
type PostingDocumentTypesListResponse struct {
	Items      []*PostingDocumentTypesResponse `json:"items"`
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
