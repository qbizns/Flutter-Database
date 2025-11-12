package posting_concept_override

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingConceptOverridesResponse represents a posting_concept_overrides response
type PostingConceptOverridesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ConceptKey string `json:"concept_key"`
	
	Label string `json:"label"`
	
	Description *string `json:"description"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreatePostingConceptOverridesRequest represents a request to create a posting_concept_overrides
type CreatePostingConceptOverridesRequest struct {
	
	ConceptKey string `json:"concept_key" validate:"required"`
	
	Label string `json:"label" validate:"required"`
	
	Description *string `json:"description"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePostingConceptOverridesRequest) Validate() error {
	
	if r.ConceptKey == "" {
		return fmt.Errorf("concept_key is required")
	}
	
	if r.Label == "" {
		return fmt.Errorf("label is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingConceptOverridesRequest represents a request to update a posting_concept_overrides
type UpdatePostingConceptOverridesRequest struct {
	
	ConceptKey *string `json:"concept_key,omitempty" validate:"omitempty,required"`
	
	Label *string `json:"label,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingConceptOverridesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ConceptKey != nil {
		hasUpdate = true
	}
	
	if r.Label != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// PostingConceptOverridesListResponse represents a paginated list of posting_concept_overrides records
type PostingConceptOverridesListResponse struct {
	Items      []*PostingConceptOverridesResponse `json:"items"`
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
