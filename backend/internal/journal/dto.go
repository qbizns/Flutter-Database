package journal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JournalsResponse represents a journals response
type JournalsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	JournalCode string `json:"journal_code"`
	
	JournalName string `json:"journal_name"`
	
	JournalType string `json:"journal_type"`
	
	'sale', *string `json:"'sale',"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id"`
	
	DefaultDebitAccountId *uuid.UUID `json:"default_debit_account_id"`
	
	DefaultCreditAccountId *uuid.UUID `json:"default_credit_account_id"`
	
	SequencePrefix *string `json:"sequence_prefix"`
	
	SequenceNumber *int64 `json:"sequence_number"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	(journalType string `json:"(journal_type"`
	
}

// CreateJournalsRequest represents a request to create a journals
type CreateJournalsRequest struct {
	
	JournalCode string `json:"journal_code" validate:"required"`
	
	JournalName string `json:"journal_name" validate:"required"`
	
	JournalType string `json:"journal_type" validate:"required"`
	
	'sale', *string `json:"'sale',"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id"`
	
	DefaultDebitAccountId *uuid.UUID `json:"default_debit_account_id"`
	
	DefaultCreditAccountId *uuid.UUID `json:"default_credit_account_id"`
	
	SequencePrefix *string `json:"sequence_prefix"`
	
	SequenceNumber *int64 `json:"sequence_number"`
	
	IsActive *bool `json:"is_active"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	(journalType string `json:"(journal_type" validate:"required"`
	
}

// Validate validates the create request
func (r *CreateJournalsRequest) Validate() error {
	
	if r.JournalCode == "" {
		return fmt.Errorf("journal_code is required")
	}
	
	if r.JournalName == "" {
		return fmt.Errorf("journal_name is required")
	}
	
	if r.JournalType == "" {
		return fmt.Errorf("journal_type is required")
	}
	
	if r.(journalType == "" {
		return fmt.Errorf("(journal_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateJournalsRequest represents a request to update a journals
type UpdateJournalsRequest struct {
	
	JournalCode *string `json:"journal_code,omitempty" validate:"omitempty,required"`
	
	JournalName *string `json:"journal_name,omitempty" validate:"omitempty,required"`
	
	JournalType *string `json:"journal_type,omitempty" validate:"omitempty,required"`
	
	'sale', *string `json:"'sale',,omitempty"`
	
	BankAccountId *uuid.UUID `json:"bank_account_id,omitempty"`
	
	DefaultDebitAccountId *uuid.UUID `json:"default_debit_account_id,omitempty"`
	
	DefaultCreditAccountId *uuid.UUID `json:"default_credit_account_id,omitempty"`
	
	SequencePrefix *string `json:"sequence_prefix,omitempty"`
	
	SequenceNumber *int64 `json:"sequence_number,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	(journalType *string `json:"(journal_type,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdateJournalsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.JournalCode != nil {
		hasUpdate = true
	}
	
	if r.JournalName != nil {
		hasUpdate = true
	}
	
	if r.JournalType != nil {
		hasUpdate = true
	}
	
	if r.'sale', != nil {
		hasUpdate = true
	}
	
	if r.BankAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultDebitAccountId != nil {
		hasUpdate = true
	}
	
	if r.DefaultCreditAccountId != nil {
		hasUpdate = true
	}
	
	if r.SequencePrefix != nil {
		hasUpdate = true
	}
	
	if r.SequenceNumber != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	
	if r.(journalType != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// JournalsListResponse represents a paginated list of journals records
type JournalsListResponse struct {
	Items      []*JournalsResponse `json:"items"`
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
