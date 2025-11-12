package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankReconciliationsResponse represents a bank_reconciliations response
type BankReconciliationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BankAccountId uuid.UUID `json:"bank_account_id"`
	
	StatementDate time.Time `json:"statement_date"`
	
	StatementBalance float64 `json:"statement_balance"`
	
	ReconciliationDate *time.Time `json:"reconciliation_date"`
	
	BookBalance *float64 `json:"book_balance"`
	
	ClearedBalance *float64 `json:"cleared_balance"`
	
	Difference *float64 `json:"difference"`
	
	Status *string `json:"status"`
	
	IsReconciled *bool `json:"is_reconciled"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	ReconciledBy *uuid.UUID `json:"reconciled_by"`
	
	ReconciledAt *time.Time `json:"reconciled_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBankReconciliationsRequest represents a request to create a bank_reconciliations
type CreateBankReconciliationsRequest struct {
	
	BankAccountId uuid.UUID `json:"bank_account_id" validate:"required"`
	
	StatementDate time.Time `json:"statement_date" validate:"required"`
	
	StatementBalance float64 `json:"statement_balance" validate:"required"`
	
	ReconciliationDate *time.Time `json:"reconciliation_date"`
	
	BookBalance *float64 `json:"book_balance"`
	
	ClearedBalance *float64 `json:"cleared_balance"`
	
	Difference *float64 `json:"difference"`
	
	Status *string `json:"status"`
	
	IsReconciled *bool `json:"is_reconciled"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	ReconciledBy *uuid.UUID `json:"reconciled_by"`
	
	ReconciledAt *time.Time `json:"reconciled_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateBankReconciliationsRequest) Validate() error {
	
	if r.BankAccountId == uuid.Nil {
		return fmt.Errorf("bank_account_id is required")
	}
	
	if r.StatementDate == nil {
		return fmt.Errorf("statement_date is required")
	}
	
	if r.StatementBalance == nil {
		return fmt.Errorf("statement_balance is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankReconciliationsRequest represents a request to update a bank_reconciliations
type UpdateBankReconciliationsRequest struct {
	
	BankAccountId *uuid.UUID `json:"bank_account_id,omitempty" validate:"omitempty,required"`
	
	StatementDate *time.Time `json:"statement_date,omitempty" validate:"omitempty,required"`
	
	StatementBalance *float64 `json:"statement_balance,omitempty" validate:"omitempty,required"`
	
	ReconciliationDate *time.Time `json:"reconciliation_date,omitempty"`
	
	BookBalance *float64 `json:"book_balance,omitempty"`
	
	ClearedBalance *float64 `json:"cleared_balance,omitempty"`
	
	Difference *float64 `json:"difference,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	IsReconciled *bool `json:"is_reconciled,omitempty"`
	
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	ReconciledBy *uuid.UUID `json:"reconciled_by,omitempty"`
	
	ReconciledAt *time.Time `json:"reconciled_at,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankReconciliationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BankAccountId != nil {
		hasUpdate = true
	}
	
	if r.StatementDate != nil {
		hasUpdate = true
	}
	
	if r.StatementBalance != nil {
		hasUpdate = true
	}
	
	if r.ReconciliationDate != nil {
		hasUpdate = true
	}
	
	if r.BookBalance != nil {
		hasUpdate = true
	}
	
	if r.ClearedBalance != nil {
		hasUpdate = true
	}
	
	if r.Difference != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.IsReconciled != nil {
		hasUpdate = true
	}
	
	if r.AccountingPeriodId != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.ReconciledBy != nil {
		hasUpdate = true
	}
	
	if r.ReconciledAt != nil {
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

// BankReconciliationsListResponse represents a paginated list of bank_reconciliations records
type BankReconciliationsListResponse struct {
	Items      []*BankReconciliationsResponse `json:"items"`
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
