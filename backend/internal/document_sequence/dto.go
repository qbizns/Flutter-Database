package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DocumentSequencesResponse represents a document_sequences response
type DocumentSequencesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	DocumentType string `json:"document_type"`
	
	Prefix *string `json:"prefix"`
	
	Suffix *string `json:"suffix"`
	
	NextNumber int64 `json:"next_number"`
	
	Padding *int64 `json:"padding"`
	
	IncrementBy *int64 `json:"increment_by"`
	
	ResetFrequency *string `json:"reset_frequency"`
	
	'never', *string `json:"'never',"`
	
	'daily', *string `json:"'daily',"`
	
	'monthly', *string `json:"'monthly',"`
	
	'yearly', *string `json:"'yearly',"`
	
	'manual' *string `json:"'manual'"`
	
	LastResetAt *time.Time `json:"last_reset_at"`
	
	LastResetValue *int64 `json:"last_reset_value"`
	
	IncludeDate *bool `json:"include_date"`
	
	DateFormat *string `json:"date_format"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	IsActive *bool `json:"is_active"`
	
	AllowManualOverride *bool `json:"allow_manual_override"`
	
	ExampleNumber *string `json:"example_number"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateDocumentSequencesRequest represents a request to create a document_sequences
type CreateDocumentSequencesRequest struct {
	
	DocumentType string `json:"document_type" validate:"required"`
	
	Prefix *string `json:"prefix"`
	
	Suffix *string `json:"suffix"`
	
	NextNumber int64 `json:"next_number" validate:"required"`
	
	Padding *int64 `json:"padding"`
	
	IncrementBy *int64 `json:"increment_by"`
	
	ResetFrequency *string `json:"reset_frequency"`
	
	'never', *string `json:"'never',"`
	
	'daily', *string `json:"'daily',"`
	
	'monthly', *string `json:"'monthly',"`
	
	'yearly', *string `json:"'yearly',"`
	
	'manual' *string `json:"'manual'"`
	
	LastResetAt *time.Time `json:"last_reset_at"`
	
	LastResetValue *int64 `json:"last_reset_value"`
	
	IncludeDate *bool `json:"include_date"`
	
	DateFormat *string `json:"date_format"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	IsActive *bool `json:"is_active"`
	
	AllowManualOverride *bool `json:"allow_manual_override"`
	
	ExampleNumber *string `json:"example_number"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateDocumentSequencesRequest) Validate() error {
	
	if r.DocumentType == "" {
		return fmt.Errorf("document_type is required")
	}
	
	if r.NextNumber == 0 {
		return fmt.Errorf("next_number is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDocumentSequencesRequest represents a request to update a document_sequences
type UpdateDocumentSequencesRequest struct {
	
	DocumentType *string `json:"document_type,omitempty" validate:"omitempty,required"`
	
	Prefix *string `json:"prefix,omitempty"`
	
	Suffix *string `json:"suffix,omitempty"`
	
	NextNumber *int64 `json:"next_number,omitempty" validate:"omitempty,required"`
	
	Padding *int64 `json:"padding,omitempty"`
	
	IncrementBy *int64 `json:"increment_by,omitempty"`
	
	ResetFrequency *string `json:"reset_frequency,omitempty"`
	
	'never', *string `json:"'never',,omitempty"`
	
	'daily', *string `json:"'daily',,omitempty"`
	
	'monthly', *string `json:"'monthly',,omitempty"`
	
	'yearly', *string `json:"'yearly',,omitempty"`
	
	'manual' *string `json:"'manual',omitempty"`
	
	LastResetAt *time.Time `json:"last_reset_at,omitempty"`
	
	LastResetValue *int64 `json:"last_reset_value,omitempty"`
	
	IncludeDate *bool `json:"include_date,omitempty"`
	
	DateFormat *string `json:"date_format,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	AllowManualOverride *bool `json:"allow_manual_override,omitempty"`
	
	ExampleNumber *string `json:"example_number,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDocumentSequencesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.DocumentType != nil {
		hasUpdate = true
	}
	
	if r.Prefix != nil {
		hasUpdate = true
	}
	
	if r.Suffix != nil {
		hasUpdate = true
	}
	
	if r.NextNumber != nil {
		hasUpdate = true
	}
	
	if r.Padding != nil {
		hasUpdate = true
	}
	
	if r.IncrementBy != nil {
		hasUpdate = true
	}
	
	if r.ResetFrequency != nil {
		hasUpdate = true
	}
	
	if r.'never', != nil {
		hasUpdate = true
	}
	
	if r.'daily', != nil {
		hasUpdate = true
	}
	
	if r.'monthly', != nil {
		hasUpdate = true
	}
	
	if r.'yearly', != nil {
		hasUpdate = true
	}
	
	if r.'manual' != nil {
		hasUpdate = true
	}
	
	if r.LastResetAt != nil {
		hasUpdate = true
	}
	
	if r.LastResetValue != nil {
		hasUpdate = true
	}
	
	if r.IncludeDate != nil {
		hasUpdate = true
	}
	
	if r.DateFormat != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.AllowManualOverride != nil {
		hasUpdate = true
	}
	
	if r.ExampleNumber != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// DocumentSequencesListResponse represents a paginated list of document_sequences records
type DocumentSequencesListResponse struct {
	Items      []*DocumentSequencesResponse `json:"items"`
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
