package general_ledger

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GeneralLedgerResponse represents a general_ledger response
type GeneralLedgerResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	JournalEntryId uuid.UUID `json:"journal_entry_id"`
	
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id"`
	
	AccountId uuid.UUID `json:"account_id"`
	
	TransactionDate time.Time `json:"transaction_date"`
	
	PostingDate time.Time `json:"posting_date"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	DebitAmount *float64 `json:"debit_amount"`
	
	CreditAmount *float64 `json:"credit_amount"`
	
	RunningDebitBalance *float64 `json:"running_debit_balance"`
	
	RunningCreditBalance *float64 `json:"running_credit_balance"`
	
	RunningBalance *float64 `json:"running_balance"`
	
	SourceModule *string `json:"source_module"`
	
	SourceDocumentType *string `json:"source_document_type"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	CostCenter *string `json:"cost_center"`
	
	Description *string `json:"description"`
	
	IsReversed *bool `json:"is_reversed"`
	
	ReversalGlId *uuid.UUID `json:"reversal_gl_id"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	(debitAmount *string `json:"(debit_amount"`
	
	(creditAmount *string `json:"(credit_amount"`
	
}

// CreateGeneralLedgerRequest represents a request to create a general_ledger
type CreateGeneralLedgerRequest struct {
	
	JournalEntryId uuid.UUID `json:"journal_entry_id" validate:"required"`
	
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id" validate:"required"`
	
	AccountId uuid.UUID `json:"account_id" validate:"required"`
	
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
	
	PostingDate time.Time `json:"posting_date" validate:"required"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id"`
	
	DebitAmount *float64 `json:"debit_amount"`
	
	CreditAmount *float64 `json:"credit_amount"`
	
	RunningDebitBalance *float64 `json:"running_debit_balance"`
	
	RunningCreditBalance *float64 `json:"running_credit_balance"`
	
	RunningBalance *float64 `json:"running_balance"`
	
	SourceModule *string `json:"source_module"`
	
	SourceDocumentType *string `json:"source_document_type"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	CostCenter *string `json:"cost_center"`
	
	Description *string `json:"description"`
	
	IsReversed *bool `json:"is_reversed"`
	
	ReversalGlId *uuid.UUID `json:"reversal_gl_id"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	(debitAmount *string `json:"(debit_amount"`
	
	(creditAmount *string `json:"(credit_amount"`
	
}

// Validate validates the create request
func (r *CreateGeneralLedgerRequest) Validate() error {
	
	if r.JournalEntryId == uuid.Nil {
		return fmt.Errorf("journal_entry_id is required")
	}
	
	if r.JournalEntryLineId == uuid.Nil {
		return fmt.Errorf("journal_entry_line_id is required")
	}
	
	if r.AccountId == uuid.Nil {
		return fmt.Errorf("account_id is required")
	}
	
	if r.TransactionDate == nil {
		return fmt.Errorf("transaction_date is required")
	}
	
	if r.PostingDate == nil {
		return fmt.Errorf("posting_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateGeneralLedgerRequest represents a request to update a general_ledger
type UpdateGeneralLedgerRequest struct {
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty" validate:"omitempty,required"`
	
	JournalEntryLineId *uuid.UUID `json:"journal_entry_line_id,omitempty" validate:"omitempty,required"`
	
	AccountId *uuid.UUID `json:"account_id,omitempty" validate:"omitempty,required"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty" validate:"omitempty,required"`
	
	PostingDate *time.Time `json:"posting_date,omitempty" validate:"omitempty,required"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	FiscalYearId *uuid.UUID `json:"fiscal_year_id,omitempty"`
	
	DebitAmount *float64 `json:"debit_amount,omitempty"`
	
	CreditAmount *float64 `json:"credit_amount,omitempty"`
	
	RunningDebitBalance *float64 `json:"running_debit_balance,omitempty"`
	
	RunningCreditBalance *float64 `json:"running_credit_balance,omitempty"`
	
	RunningBalance *float64 `json:"running_balance,omitempty"`
	
	SourceModule *string `json:"source_module,omitempty"`
	
	SourceDocumentType *string `json:"source_document_type,omitempty"`
	
	SourceDocumentId *uuid.UUID `json:"source_document_id,omitempty"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	Department *string `json:"department,omitempty"`
	
	ProjectCode *string `json:"project_code,omitempty"`
	
	CostCenter *string `json:"cost_center,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	IsReversed *bool `json:"is_reversed,omitempty"`
	
	ReversalGlId *uuid.UUID `json:"reversal_gl_id,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	(debitAmount *string `json:"(debit_amount,omitempty"`
	
	(creditAmount *string `json:"(credit_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateGeneralLedgerRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryLineId != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
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
	
	if r.DebitAmount != nil {
		hasUpdate = true
	}
	
	if r.CreditAmount != nil {
		hasUpdate = true
	}
	
	if r.RunningDebitBalance != nil {
		hasUpdate = true
	}
	
	if r.RunningCreditBalance != nil {
		hasUpdate = true
	}
	
	if r.RunningBalance != nil {
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
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.Department != nil {
		hasUpdate = true
	}
	
	if r.ProjectCode != nil {
		hasUpdate = true
	}
	
	if r.CostCenter != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.IsReversed != nil {
		hasUpdate = true
	}
	
	if r.ReversalGlId != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.(debitAmount != nil {
		hasUpdate = true
	}
	
	if r.(creditAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// GeneralLedgerListResponse represents a paginated list of general_ledger records
type GeneralLedgerListResponse struct {
	Items      []*GeneralLedgerResponse `json:"items"`
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
