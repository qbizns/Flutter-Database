package bank_statement_reconciliation

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankStatementReconciliationsResponse represents a bank_statement_reconciliations response
type BankStatementReconciliationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BankStatementLineId uuid.UUID `json:"bank_statement_line_id"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	PaymentId *uuid.UUID `json:"payment_id"`
	
	MatchedAmount float64 `json:"matched_amount"`
	
	MatchedBy *uuid.UUID `json:"matched_by"`
	
	MatchedAt *time.Time `json:"matched_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBankStatementReconciliationsRequest represents a request to create a bank_statement_reconciliations
type CreateBankStatementReconciliationsRequest struct {
	
	BankStatementLineId uuid.UUID `json:"bank_statement_line_id" validate:"required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id"`
	
	PaymentId *uuid.UUID `json:"payment_id"`
	
	MatchedAmount float64 `json:"matched_amount" validate:"required"`
	
	MatchedBy *uuid.UUID `json:"matched_by"`
	
	MatchedAt *time.Time `json:"matched_at"`
	
}

// Validate validates the create request
func (r *CreateBankStatementReconciliationsRequest) Validate() error {
	
	if r.BankStatementLineId == uuid.Nil {
		return fmt.Errorf("bank_statement_line_id is required")
	}
	
	// Numeric field validation
	// TODO: Add validation for numeric fields
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankStatementReconciliationsRequest represents a request to update a bank_statement_reconciliations
type UpdateBankStatementReconciliationsRequest struct {
	
	BankStatementLineId *uuid.UUID `json:"bank_statement_line_id,omitempty" validate:"omitempty,required"`
	
	JournalEntryId *uuid.UUID `json:"journal_entry_id,omitempty"`
	
	PaymentId *uuid.UUID `json:"payment_id,omitempty"`
	
	MatchedAmount *float64 `json:"matched_amount,omitempty" validate:"omitempty,required"`
	
	MatchedBy *uuid.UUID `json:"matched_by,omitempty"`
	
	MatchedAt *time.Time `json:"matched_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankStatementReconciliationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BankStatementLineId != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryId != nil {
		hasUpdate = true
	}
	
	if r.PaymentId != nil {
		hasUpdate = true
	}
	
	if r.MatchedAmount != nil {
		hasUpdate = true
	}
	
	if r.MatchedBy != nil {
		hasUpdate = true
	}
	
	if r.MatchedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BankStatementReconciliationsListResponse represents a paginated list of bank_statement_reconciliations records
type BankStatementReconciliationsListResponse struct {
	Items      []*BankStatementReconciliationsResponse `json:"items"`
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
