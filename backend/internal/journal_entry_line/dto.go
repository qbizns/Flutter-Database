package journal_entry_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JournalEntryLinesResponse represents a journal_entry_lines response
type JournalEntryLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	JournalEntryId uuid.UUID `json:"journal_entry_id"`
	
	LineNumber int64 `json:"line_number"`
	
	AccountId uuid.UUID `json:"account_id"`
	
	DebitAmount *float64 `json:"debit_amount"`
	
	CreditAmount *float64 `json:"credit_amount"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	CostCenter *string `json:"cost_center"`
	
	TaxCode *string `json:"tax_code"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	Description *string `json:"description"`
	
	Memo *string `json:"memo"`
	
	IsReconciled *bool `json:"is_reconciled"`
	
	ReconciledAt *time.Time `json:"reconciled_at"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	(debitAmount *string `json:"(debit_amount"`
	
	(creditAmount *string `json:"(credit_amount"`
	
	(debitAmount *string `json:"(debit_amount"`
	
}

// CreateJournalEntryLinesRequest represents a request to create a journal_entry_lines
type CreateJournalEntryLinesRequest struct {
	
	JournalEntryId uuid.UUID `json:"journal_entry_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	AccountId uuid.UUID `json:"account_id" validate:"required"`
	
	DebitAmount *float64 `json:"debit_amount"`
	
	CreditAmount *float64 `json:"credit_amount"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Department *string `json:"department"`
	
	ProjectCode *string `json:"project_code"`
	
	CostCenter *string `json:"cost_center"`
	
	TaxCode *string `json:"tax_code"`
	
	TaxAmount *float64 `json:"tax_amount"`
	
	Description *string `json:"description"`
	
	Memo *string `json:"memo"`
	
	IsReconciled *bool `json:"is_reconciled"`
	
	ReconciledAt *time.Time `json:"reconciled_at"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(debitAmount *string `json:"(debit_amount"`
	
	(creditAmount *string `json:"(credit_amount"`
	
	(debitAmount *string `json:"(debit_amount"`
	
}

// Validate validates the create request
func (r *CreateJournalEntryLinesRequest) Validate() error {
	
	if r.JournalEntryId == uuid.Nil {
		return fmt.Errorf("journal_entry_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.AccountId == uuid.Nil {
		return fmt.Errorf("account_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateJournalEntryLinesRequest represents a request to update a journal_entry_lines
type UpdateJournalEntryLinesRequest struct {
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	AccountId *uuid.UUID `json:"account_id,omitempty" validate:"omitempty,required"`
	
	DebitAmount *float64 `json:"debit_amount,omitempty"`
	
	CreditAmount *float64 `json:"credit_amount,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	Department *string `json:"department,omitempty"`
	
	ProjectCode *string `json:"project_code,omitempty"`
	
	CostCenter *string `json:"cost_center,omitempty"`
	
	TaxCode *string `json:"tax_code,omitempty"`
	
	TaxAmount *float64 `json:"tax_amount,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Memo *string `json:"memo,omitempty"`
	
	IsReconciled *bool `json:"is_reconciled,omitempty"`
	
	ReconciledAt *time.Time `json:"reconciled_at,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(debitAmount *string `json:"(debit_amount,omitempty"`
	
	(creditAmount *string `json:"(credit_amount,omitempty"`
	
	(debitAmount *string `json:"(debit_amount,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateJournalEntryLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.AccountId != nil {
		hasUpdate = true
	}
	
	if r.DebitAmount != nil {
		hasUpdate = true
	}
	
	if r.CreditAmount != nil {
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
	
	if r.TaxCode != nil {
		hasUpdate = true
	}
	
	if r.TaxAmount != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Memo != nil {
		hasUpdate = true
	}
	
	if r.IsReconciled != nil {
		hasUpdate = true
	}
	
	if r.ReconciledAt != nil {
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
	
	if r.(debitAmount != nil {
		hasUpdate = true
	}
	
	if r.(creditAmount != nil {
		hasUpdate = true
	}
	
	if r.(debitAmount != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// JournalEntryLinesListResponse represents a paginated list of journal_entry_lines records
type JournalEntryLinesListResponse struct {
	Items      []*JournalEntryLinesResponse `json:"items"`
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
