package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JournalEntriesResponse represents a journal_entries response
type JournalEntriesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	EntryNumber string `json:"entry_number"`
	
	EntryTypeId uuid.UUID `json:"entry_type_id"`
	
	EntryDate time.Time `json:"entry_date"`
	
	PostingDate time.Time `json:"posting_date"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	Status *string `json:"status"`
	
	IsPosted *bool `json:"is_posted"`
	
	IsReversed *bool `json:"is_reversed"`
	
	ReversalEntryId *uuid.UUID `json:"reversal_entry_id"`
	
	SourceModule *string `json:"source_module"`
	
	SourceDocumentType *string `json:"source_document_type"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	TotalDebit *float64 `json:"total_debit"`
	
	TotalCredit *float64 `json:"total_credit"`
	
	Description string `json:"description"`
	
	Notes *string `json:"notes"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	Attachments json.RawMessage `json:"attachments"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	(isPosted *string `json:"(is_posted"`
	
	(ABS(totalDebit *string `json:"(ABS(total_debit"`
	
}

// CreateJournalEntriesRequest represents a request to create a journal_entries
type CreateJournalEntriesRequest struct {
	
	EntryNumber string `json:"entry_number" validate:"required"`
	
	EntryTypeId uuid.UUID `json:"entry_type_id" validate:"required"`
	
	EntryDate time.Time `json:"entry_date" validate:"required"`
	
	PostingDate time.Time `json:"posting_date" validate:"required"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	Status *string `json:"status"`
	
	IsPosted *bool `json:"is_posted"`
	
	IsReversed *bool `json:"is_reversed"`
	
	ReversalEntryId *uuid.UUID `json:"reversal_entry_id"`
	
	SourceModule *string `json:"source_module"`
	
	SourceDocumentType *string `json:"source_document_type"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	TotalDebit *float64 `json:"total_debit"`
	
	TotalCredit *float64 `json:"total_credit"`
	
	Description string `json:"description" validate:"required"`
	
	Notes *string `json:"notes"`
	
	RequiresApproval *bool `json:"requires_approval"`
	
	ApprovedBy *uuid.UUID `json:"approved_by"`
	
	ApprovedAt *time.Time `json:"approved_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	Attachments json.RawMessage `json:"attachments"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(isPosted *string `json:"(is_posted"`
	
	(ABS(totalDebit *string `json:"(ABS(total_debit"`
	
}

// Validate validates the create request
func (r *CreateJournalEntriesRequest) Validate() error {
	
	if r.EntryNumber == "" {
		return fmt.Errorf("entry_number is required")
	}
	
	if r.EntryTypeId == uuid.Nil {
		return fmt.Errorf("entry_type_id is required")
	}
	
	if r.EntryDate == nil {
		return fmt.Errorf("entry_date is required")
	}
	
	if r.PostingDate == nil {
		return fmt.Errorf("posting_date is required")
	}
	
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateJournalEntriesRequest represents a request to update a journal_entries
type UpdateJournalEntriesRequest struct {
	
	EntryNumber *string `json:"entry_number,omitempty" validate:"omitempty,required"`
	
	EntryTypeId *uuid.UUID `json:"entry_type_id,omitempty" validate:"omitempty,required"`
	
	EntryDate *time.Time `json:"entry_date,omitempty" validate:"omitempty,required"`
	
	PostingDate *time.Time `json:"posting_date,omitempty" validate:"omitempty,required"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	IsPosted *bool `json:"is_posted,omitempty"`
	
	IsReversed *bool `json:"is_reversed,omitempty"`
	
	ReversalEntryId *uuid.UUID `json:"reversal_entry_id,omitempty"`
	
	SourceModule *string `json:"source_module,omitempty"`
	
	SourceDocumentType *string `json:"source_document_type,omitempty"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id,omitempty"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	TotalDebit *float64 `json:"total_debit,omitempty"`
	
	TotalCredit *float64 `json:"total_credit,omitempty"`
	
	Description *string `json:"description,omitempty" validate:"omitempty,required"`
	
	Notes *string `json:"notes,omitempty"`
	
	RequiresApproval *bool `json:"requires_approval,omitempty"`
	
	ApprovedBy *uuid.UUID `json:"approved_by,omitempty"`
	
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	
	PostedBy *uuid.UUID `json:"posted_by,omitempty"`
	
	PostedAt *time.Time `json:"posted_at,omitempty"`
	
	Attachments *json.RawMessage `json:"attachments,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(isPosted *string `json:"(is_posted,omitempty"`
	
	(ABS(totalDebit *string `json:"(ABS(total_debit,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateJournalEntriesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.EntryNumber != nil {
		hasUpdate = true
	}
	
	if r.EntryTypeId != nil {
		hasUpdate = true
	}
	
	if r.EntryDate != nil {
		hasUpdate = true
	}
	
	if r.PostingDate != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
		hasUpdate = true
	}
	
	if r.FiscalYearId != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.IsPosted != nil {
		hasUpdate = true
	}
	
	if r.IsReversed != nil {
		hasUpdate = true
	}
	
	if r.ReversalEntryId != nil {
		hasUpdate = true
	}
	
	if r.SourceModule != nil {
		hasUpdate = true
	}
	
	if r.SourceDocumentType != nil {
		hasUpdate = true
	}
	
	if r.SourceDocumentId != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.TotalDebit != nil {
		hasUpdate = true
	}
	
	if r.TotalCredit != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.RequiresApproval != nil {
		hasUpdate = true
	}
	
	if r.ApprovedBy != nil {
		hasUpdate = true
	}
	
	if r.ApprovedAt != nil {
		hasUpdate = true
	}
	
	if r.PostedBy != nil {
		hasUpdate = true
	}
	
	if r.PostedAt != nil {
		hasUpdate = true
	}
	
	if r.Attachments != nil {
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
	
	if r.(isPosted != nil {
		hasUpdate = true
	}
	
	if r.(ABS(totalDebit != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// JournalEntriesListResponse represents a paginated list of journal_entries records
type JournalEntriesListResponse struct {
	Items      []*JournalEntriesResponse `json:"items"`
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
