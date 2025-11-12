package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JournalEntryTypesResponse represents a journal_entry_types response
type JournalEntryTypesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	TypeCode string `json:"type_code"`
	
	TypeName string `json:"type_name"`
	
	TypeCategory *string `json:"type_category"`
	
	NumberPrefix *string `json:"number_prefix"`
	
	Description *string `json:"description"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateJournalEntryTypesRequest represents a request to create a journal_entry_types
type CreateJournalEntryTypesRequest struct {
	
	TypeCode string `json:"type_code" validate:"required"`
	
	TypeName string `json:"type_name" validate:"required"`
	
	TypeCategory *string `json:"type_category"`
	
	NumberPrefix *string `json:"number_prefix"`
	
	Description *string `json:"description"`
	
}

// Validate validates the create request
func (r *CreateJournalEntryTypesRequest) Validate() error {
	
	if r.TypeCode == "" {
		return fmt.Errorf("type_code is required")
	}
	
	if r.TypeName == "" {
		return fmt.Errorf("type_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateJournalEntryTypesRequest represents a request to update a journal_entry_types
type UpdateJournalEntryTypesRequest struct {
	
	TypeCode *string `json:"type_code,omitempty" validate:"omitempty,required"`
	
	TypeName *string `json:"type_name,omitempty" validate:"omitempty,required"`
	
	TypeCategory *string `json:"type_category,omitempty"`
	
	NumberPrefix *string `json:"number_prefix,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateJournalEntryTypesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TypeCode != nil {
		hasUpdate = true
	}
	
	if r.TypeName != nil {
		hasUpdate = true
	}
	
	if r.TypeCategory != nil {
		hasUpdate = true
	}
	
	if r.NumberPrefix != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// JournalEntryTypesListResponse represents a paginated list of journal_entry_types records
type JournalEntryTypesListResponse struct {
	Items      []*JournalEntryTypesResponse `json:"items"`
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
