package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankStatementsResponse represents a bank_statements response
type BankStatementsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	BankAccountId uuid.UUID `json:"bank_account_id"`
	
	StatementNumber *string `json:"statement_number"`
	
	StatementDate time.Time `json:"statement_date"`
	
	PeriodStartDate time.Time `json:"period_start_date"`
	
	PeriodEndDate time.Time `json:"period_end_date"`
	
	OpeningBalance float64 `json:"opening_balance"`
	
	ClosingBalance float64 `json:"closing_balance"`
	
	ImportSource *string `json:"import_source"`
	
	ImportSource *string `json:"import_source"`
	
	ImportFileName *string `json:"import_file_name"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBankStatementsRequest represents a request to create a bank_statements
type CreateBankStatementsRequest struct {
	
	BankAccountId uuid.UUID `json:"bank_account_id" validate:"required"`
	
	StatementNumber *string `json:"statement_number"`
	
	StatementDate time.Time `json:"statement_date" validate:"required"`
	
	PeriodStartDate time.Time `json:"period_start_date" validate:"required"`
	
	PeriodEndDate time.Time `json:"period_end_date" validate:"required"`
	
	OpeningBalance float64 `json:"opening_balance" validate:"required"`
	
	ClosingBalance float64 `json:"closing_balance" validate:"required"`
	
	ImportSource *string `json:"import_source"`
	
	ImportSource *string `json:"import_source"`
	
	ImportFileName *string `json:"import_file_name"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateBankStatementsRequest) Validate() error {
	
	if r.BankAccountId == uuid.Nil {
		return fmt.Errorf("bank_account_id is required")
	}
	
	if r.StatementDate == nil {
		return fmt.Errorf("statement_date is required")
	}
	
	if r.PeriodStartDate == nil {
		return fmt.Errorf("period_start_date is required")
	}
	
	if r.PeriodEndDate == nil {
		return fmt.Errorf("period_end_date is required")
	}
	
	if r.OpeningBalance == nil {
		return fmt.Errorf("opening_balance is required")
	}
	
	if r.ClosingBalance == nil {
		return fmt.Errorf("closing_balance is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankStatementsRequest represents a request to update a bank_statements
type UpdateBankStatementsRequest struct {
	
	BankAccountId *uuid.UUID `json:"bank_account_id,omitempty" validate:"omitempty,required"`
	
	StatementNumber *string `json:"statement_number,omitempty"`
	
	StatementDate *time.Time `json:"statement_date,omitempty" validate:"omitempty,required"`
	
	PeriodStartDate *time.Time `json:"period_start_date,omitempty" validate:"omitempty,required"`
	
	PeriodEndDate *time.Time `json:"period_end_date,omitempty" validate:"omitempty,required"`
	
	OpeningBalance *float64 `json:"opening_balance,omitempty" validate:"omitempty,required"`
	
	ClosingBalance *float64 `json:"closing_balance,omitempty" validate:"omitempty,required"`
	
	ImportSource *string `json:"import_source,omitempty"`
	
	ImportSource *string `json:"import_source,omitempty"`
	
	ImportFileName *string `json:"import_file_name,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankStatementsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BankAccountId != nil {
		hasUpdate = true
	}
	
	if r.StatementNumber != nil {
		hasUpdate = true
	}
	
	if r.StatementDate != nil {
		hasUpdate = true
	}
	
	if r.PeriodStartDate != nil {
		hasUpdate = true
	}
	
	if r.PeriodEndDate != nil {
		hasUpdate = true
	}
	
	if r.OpeningBalance != nil {
		hasUpdate = true
	}
	
	if r.ClosingBalance != nil {
		hasUpdate = true
	}
	
	if r.ImportSource != nil {
		hasUpdate = true
	}
	
	if r.ImportSource != nil {
		hasUpdate = true
	}
	
	if r.ImportFileName != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BankStatementsListResponse represents a paginated list of bank_statements records
type BankStatementsListResponse struct {
	Items      []*BankStatementsResponse `json:"items"`
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
