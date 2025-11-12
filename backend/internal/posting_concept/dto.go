package posting_concept

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingConceptsResponse represents a posting_concepts response
type PostingConceptsResponse struct {
	
	ConceptKey *string `json:"concept_key"`
	
	DefaultLabel string `json:"default_label"`
	
	DefaultDescription *string `json:"default_description"`
	
	ExpectedAccountTypeId *uuid.UUID `json:"expected_account_type_id"`
	
	NormalSide *string `json:"normal_side"`
	
	ExampleCode *string `json:"example_code"`
	
	ExampleAccountName *string `json:"example_account_name"`
	
	IsSystem bool `json:"is_system"`
	
	ConceptCategory *string `json:"concept_category"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreatePostingConceptsRequest represents a request to create a posting_concepts
type CreatePostingConceptsRequest struct {
	
	ConceptKey *string `json:"concept_key"`
	
	DefaultLabel string `json:"default_label" validate:"required"`
	
	DefaultDescription *string `json:"default_description"`
	
	ExpectedAccountTypeId *uuid.UUID `json:"expected_account_type_id"`
	
	NormalSide *string `json:"normal_side"`
	
	ExampleCode *string `json:"example_code"`
	
	ExampleAccountName *string `json:"example_account_name"`
	
	IsSystem bool `json:"is_system" validate:"required"`
	
	ConceptCategory *string `json:"concept_category"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreatePostingConceptsRequest) Validate() error {
	
	if r.DefaultLabel == "" {
		return fmt.Errorf("default_label is required")
	}
	
	if r.IsSystem == nil {
		return fmt.Errorf("is_system is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePostingConceptsRequest represents a request to update a posting_concepts
type UpdatePostingConceptsRequest struct {
	
	ConceptKey *string `json:"concept_key,omitempty"`
	
	DefaultLabel *string `json:"default_label,omitempty" validate:"omitempty,required"`
	
	DefaultDescription *string `json:"default_description,omitempty"`
	
	ExpectedAccountTypeId *uuid.UUID `json:"expected_account_type_id,omitempty"`
	
	NormalSide *string `json:"normal_side,omitempty"`
	
	ExampleCode *string `json:"example_code,omitempty"`
	
	ExampleAccountName *string `json:"example_account_name,omitempty"`
	
	IsSystem *bool `json:"is_system,omitempty" validate:"omitempty,required"`
	
	ConceptCategory *string `json:"concept_category,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePostingConceptsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ConceptKey != nil {
		hasUpdate = true
	}
	
	if r.DefaultLabel != nil {
		hasUpdate = true
	}
	
	if r.DefaultDescription != nil {
		hasUpdate = true
	}
	
	if r.ExpectedAccountTypeId != nil {
		hasUpdate = true
	}
	
	if r.NormalSide != nil {
		hasUpdate = true
	}
	
	if r.ExampleCode != nil {
		hasUpdate = true
	}
	
	if r.ExampleAccountName != nil {
		hasUpdate = true
	}
	
	if r.IsSystem != nil {
		hasUpdate = true
	}
	
	if r.ConceptCategory != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
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

// PostingConceptsListResponse represents a paginated list of posting_concepts records
type PostingConceptsListResponse struct {
	Items      []*PostingConceptsResponse `json:"items"`
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
