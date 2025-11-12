package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PosPostingAuditResponse represents a pos_posting_audit response
type PosPostingAuditResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SourceTable string `json:"source_table"`
	
	SourceId uuid.UUID `json:"source_id"`
	
	SourceReference *string `json:"source_reference"`
	
	PostingStatus string `json:"posting_status"`
	
	'pending', *string `json:"'pending',"`
	
	'processing', *string `json:"'processing',"`
	
	'posted', *string `json:"'posted',"`
	
	'failed', *string `json:"'failed',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'reversed' *string `json:"'reversed'"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	ReversalJournalEntryId *uuid.UUID `json:"reversal_journal_entry_id"`
	
	PostingDate *time.Time `json:"posting_date"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
	PostingMethod *string `json:"posting_method"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	RetryCount *int64 `json:"retry_count"`
	
	LastRetryAt *time.Time `json:"last_retry_at"`
	
	MaxRetries *int64 `json:"max_retries"`
	
	ReversedAt *time.Time `json:"reversed_at"`
	
	ReversedBy *uuid.UUID `json:"reversed_by"`
	
	ReversalReason *string `json:"reversal_reason"`
	
	TotalDebit *float64 `json:"total_debit"`
	
	TotalCredit *float64 `json:"total_credit"`
	
	LineCount *int64 `json:"line_count"`
	
	CurrencyCode *string `json:"currency_code"`
	
	PostingContext json.RawMessage `json:"posting_context"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(postingStatus string `json:"(posting_status"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(postingStatus string `json:"(posting_status"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(ABS(COALESCE(totalDebit, *string `json:"(ABS(COALESCE(total_debit,"`
	
}

// CreatePosPostingAuditRequest represents a request to create a pos_posting_audit
type CreatePosPostingAuditRequest struct {
	
	SourceTable string `json:"source_table" validate:"required"`
	
	SourceId uuid.UUID `json:"source_id" validate:"required"`
	
	SourceReference *string `json:"source_reference"`
	
	PostingStatus string `json:"posting_status" validate:"required"`
	
	'pending', *string `json:"'pending',"`
	
	'processing', *string `json:"'processing',"`
	
	'posted', *string `json:"'posted',"`
	
	'failed', *string `json:"'failed',"`
	
	'cancelled', *string `json:"'cancelled',"`
	
	'reversed' *string `json:"'reversed'"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	ReversalJournalEntryId *uuid.UUID `json:"reversal_journal_entry_id"`
	
	PostingDate *time.Time `json:"posting_date"`
	
	PostedAt *time.Time `json:"posted_at"`
	
	PostedBy *uuid.UUID `json:"posted_by"`
	
	PostingMethod *string `json:"posting_method"`
	
	ErrorCode *string `json:"error_code"`
	
	ErrorMessage *string `json:"error_message"`
	
	ErrorDetails json.RawMessage `json:"error_details"`
	
	RetryCount *int64 `json:"retry_count"`
	
	LastRetryAt *time.Time `json:"last_retry_at"`
	
	MaxRetries *int64 `json:"max_retries"`
	
	ReversedAt *time.Time `json:"reversed_at"`
	
	ReversedBy *uuid.UUID `json:"reversed_by"`
	
	ReversalReason *string `json:"reversal_reason"`
	
	TotalDebit *float64 `json:"total_debit"`
	
	TotalCredit *float64 `json:"total_credit"`
	
	LineCount *int64 `json:"line_count"`
	
	CurrencyCode *string `json:"currency_code"`
	
	PostingContext json.RawMessage `json:"posting_context"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(postingStatus string `json:"(posting_status" validate:"required"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(postingStatus string `json:"(posting_status" validate:"required"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(postingStatus *string `json:"(posting_status"`
	
	(ABS(COALESCE(totalDebit, *string `json:"(ABS(COALESCE(total_debit,"`
	
}

// Validate validates the create request
func (r *CreatePosPostingAuditRequest) Validate() error {
	
	if r.SourceTable == "" {
		return fmt.Errorf("source_table is required")
	}
	
	if r.SourceId == uuid.Nil {
		return fmt.Errorf("source_id is required")
	}
	
	if r.PostingStatus == "" {
		return fmt.Errorf("posting_status is required")
	}
	
	if r.(postingStatus == "" {
		return fmt.Errorf("(posting_status is required")
	}
	
	if r.(postingStatus == "" {
		return fmt.Errorf("(posting_status is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePosPostingAuditRequest represents a request to update a pos_posting_audit
type UpdatePosPostingAuditRequest struct {
	
	SourceTable *string `json:"source_table,omitempty" validate:"omitempty,required"`
	
	SourceId *uuid.UUID `json:"source_id,omitempty" validate:"omitempty,required"`
	
	SourceReference *string `json:"source_reference,omitempty"`
	
	PostingStatus *string `json:"posting_status,omitempty" validate:"omitempty,required"`
	
	'pending', *string `json:"'pending',,omitempty"`
	
	'processing', *string `json:"'processing',,omitempty"`
	
	'posted', *string `json:"'posted',,omitempty"`
	
	'failed', *string `json:"'failed',,omitempty"`
	
	'cancelled', *string `json:"'cancelled',,omitempty"`
	
	'reversed' *string `json:"'reversed',omitempty"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	ReversalJournalEntryId *uuid.UUID `json:"reversal_journal_entry_id,omitempty"`
	
	PostingDate *time.Time `json:"posting_date,omitempty"`
	
	PostedAt *time.Time `json:"posted_at,omitempty"`
	
	PostedBy *uuid.UUID `json:"posted_by,omitempty"`
	
	PostingMethod *string `json:"posting_method,omitempty"`
	
	ErrorCode *string `json:"error_code,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	ErrorDetails *json.RawMessage `json:"error_details,omitempty"`
	
	RetryCount *int64 `json:"retry_count,omitempty"`
	
	LastRetryAt *time.Time `json:"last_retry_at,omitempty"`
	
	MaxRetries *int64 `json:"max_retries,omitempty"`
	
	ReversedAt *time.Time `json:"reversed_at,omitempty"`
	
	ReversedBy *uuid.UUID `json:"reversed_by,omitempty"`
	
	ReversalReason *string `json:"reversal_reason,omitempty"`
	
	TotalDebit *float64 `json:"total_debit,omitempty"`
	
	TotalCredit *float64 `json:"total_credit,omitempty"`
	
	LineCount *int64 `json:"line_count,omitempty"`
	
	CurrencyCode *string `json:"currency_code,omitempty"`
	
	PostingContext *json.RawMessage `json:"posting_context,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(postingStatus *string `json:"(posting_status,omitempty" validate:"omitempty,required"`
	
	(postingStatus *string `json:"(posting_status,omitempty"`
	
	(postingStatus *string `json:"(posting_status,omitempty" validate:"omitempty,required"`
	
	(postingStatus *string `json:"(posting_status,omitempty"`
	
	(postingStatus *string `json:"(posting_status,omitempty"`
	
	(ABS(COALESCE(totalDebit, *string `json:"(ABS(COALESCE(total_debit,,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePosPostingAuditRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SourceTable != nil {
		hasUpdate = true
	}
	
	if r.SourceId != nil {
		hasUpdate = true
	}
	
	if r.SourceReference != nil {
		hasUpdate = true
	}
	
	if r.PostingStatus != nil {
		hasUpdate = true
	}
	
	if r.'pending', != nil {
		hasUpdate = true
	}
	
	if r.'processing', != nil {
		hasUpdate = true
	}
	
	if r.'posted', != nil {
		hasUpdate = true
	}
	
	if r.'failed', != nil {
		hasUpdate = true
	}
	
	if r.'cancelled', != nil {
		hasUpdate = true
	}
	
	if r.'reversed' != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.ReversalJournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.PostingDate != nil {
		hasUpdate = true
	}
	
	if r.PostedAt != nil {
		hasUpdate = true
	}
	
	if r.PostedBy != nil {
		hasUpdate = true
	}
	
	if r.PostingMethod != nil {
		hasUpdate = true
	}
	
	if r.ErrorCode != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.ErrorDetails != nil {
		hasUpdate = true
	}
	
	if r.RetryCount != nil {
		hasUpdate = true
	}
	
	if r.LastRetryAt != nil {
		hasUpdate = true
	}
	
	if r.MaxRetries != nil {
		hasUpdate = true
	}
	
	if r.ReversedAt != nil {
		hasUpdate = true
	}
	
	if r.ReversedBy != nil {
		hasUpdate = true
	}
	
	if r.ReversalReason != nil {
		hasUpdate = true
	}
	
	if r.TotalDebit != nil {
		hasUpdate = true
	}
	
	if r.TotalCredit != nil {
		hasUpdate = true
	}
	
	if r.LineCount != nil {
		hasUpdate = true
	}
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.PostingContext != nil {
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
	
	if r.(postingStatus != nil {
		hasUpdate = true
	}
	
	if r.(postingStatus != nil {
		hasUpdate = true
	}
	
	if r.(postingStatus != nil {
		hasUpdate = true
	}
	
	if r.(postingStatus != nil {
		hasUpdate = true
	}
	
	if r.(postingStatus != nil {
		hasUpdate = true
	}
	
	if r.(ABS(COALESCE(totalDebit, != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PosPostingAuditListResponse represents a paginated list of pos_posting_audit records
type PosPostingAuditListResponse struct {
	Items      []*PosPostingAuditResponse `json:"items"`
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
