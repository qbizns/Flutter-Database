package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankReconciliationItemsResponse represents a bank_reconciliation_items response
type BankReconciliationItemsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BankReconciliationId *uuid.UUID `json:"bank_reconciliation_id"`
	
	GeneralLedgerId uuid.UUID `json:"general_ledger_id"`
	
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id"`
	
	IsCleared *bool `json:"is_cleared"`
	
	ClearedDate *time.Time `json:"cleared_date"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	ClearedBy *uuid.UUID `json:"cleared_by"`
	
}

// CreateBankReconciliationItemsRequest represents a request to create a bank_reconciliation_items
type CreateBankReconciliationItemsRequest struct {
	
	BankReconciliationId *uuid.UUID `json:"bank_reconciliation_id"`
	
	GeneralLedgerId uuid.UUID `json:"general_ledger_id" validate:"required"`
	
	JournalEntryLineId uuid.UUID `json:"journal_entry_line_id" validate:"required"`
	
	IsCleared *bool `json:"is_cleared"`
	
	ClearedDate *time.Time `json:"cleared_date"`
	
	ClearedBy *uuid.UUID `json:"cleared_by"`
	
}

// Validate validates the create request
func (r *CreateBankReconciliationItemsRequest) Validate() error {
	
	if r.GeneralLedgerId == uuid.Nil {
		return fmt.Errorf("general_ledger_id is required")
	}
	
	if r.JournalEntryLineId == uuid.Nil {
		return fmt.Errorf("journal_entry_line_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankReconciliationItemsRequest represents a request to update a bank_reconciliation_items
type UpdateBankReconciliationItemsRequest struct {
	
	BankReconciliationId *uuid.UUID `json:"bank_reconciliation_id,omitempty"`
	
	GeneralLedgerId *uuid.UUID `json:"general_ledger_id,omitempty" validate:"omitempty,required"`
	
	JournalEntryLineId *uuid.UUID `json:"journal_entry_line_id,omitempty" validate:"omitempty,required"`
	
	IsCleared *bool `json:"is_cleared,omitempty"`
	
	ClearedDate *time.Time `json:"cleared_date,omitempty"`
	
	ClearedBy *uuid.UUID `json:"cleared_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankReconciliationItemsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BankReconciliationId != nil {
		hasUpdate = true
	}
	
	if r.GeneralLedgerId != nil {
		hasUpdate = true
	}
	
	if r.JournalEntryLineId != nil {
		hasUpdate = true
	}
	
	if r.IsCleared != nil {
		hasUpdate = true
	}
	
	if r.ClearedDate != nil {
		hasUpdate = true
	}
	
	if r.ClearedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BankReconciliationItemsListResponse represents a paginated list of bank_reconciliation_items records
type BankReconciliationItemsListResponse struct {
	Items      []*BankReconciliationItemsResponse `json:"items"`
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
