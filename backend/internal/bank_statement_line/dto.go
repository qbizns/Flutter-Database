package bank_statement_line

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankStatementLinesResponse represents a bank_statement_lines response
type BankStatementLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	BankStatementId uuid.UUID `json:"bank_statement_id"`
	
	LineNumber int64 `json:"line_number"`
	
	TransactionDate time.Time `json:"transaction_date"`
	
	ValueDate *time.Time `json:"value_date"`
	
	Amount float64 `json:"amount"`
	
	CurrencyCode *string `json:"currency_code"`
	
	Description *string `json:"description"`
	
	Reference *string `json:"reference"`
	
	CounterpartyName *string `json:"counterparty_name"`
	
	CounterpartyAccount *string `json:"counterparty_account"`
	
	BankReference *string `json:"bank_reference"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBankStatementLinesRequest represents a request to create a bank_statement_lines
type CreateBankStatementLinesRequest struct {
	
	BankStatementId uuid.UUID `json:"bank_statement_id" validate:"required"`
	
	LineNumber int64 `json:"line_number" validate:"required"`
	
	TransactionDate time.Time `json:"transaction_date" validate:"required"`
	
	ValueDate *time.Time `json:"value_date"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	CurrencyCode *string `json:"currency_code"`
	
	Description *string `json:"description"`
	
	Reference *string `json:"reference"`
	
	CounterpartyName *string `json:"counterparty_name"`
	
	CounterpartyAccount *string `json:"counterparty_account"`
	
	BankReference *string `json:"bank_reference"`
	
	Status *string `json:"status"`
	
	Status *string `json:"status"`
	
	Notes *string `json:"notes"`
	
}

// Validate validates the create request
func (r *CreateBankStatementLinesRequest) Validate() error {
	
	if r.BankStatementId == uuid.Nil {
		return fmt.Errorf("bank_statement_id is required")
	}
	
	if r.LineNumber == 0 {
		return fmt.Errorf("line_number is required")
	}
	
	if r.TransactionDate == nil {
		return fmt.Errorf("transaction_date is required")
	}
	
	if r.Amount == nil {
		return fmt.Errorf("amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankStatementLinesRequest represents a request to update a bank_statement_lines
type UpdateBankStatementLinesRequest struct {
	
	BankStatementId *uuid.UUID `json:"bank_statement_id,omitempty" validate:"omitempty,required"`
	
	LineNumber *int64 `json:"line_number,omitempty" validate:"omitempty,required"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty" validate:"omitempty,required"`
	
	ValueDate *time.Time `json:"value_date,omitempty"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	CurrencyCode *string `json:"currency_code,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Reference *string `json:"reference,omitempty"`
	
	CounterpartyName *string `json:"counterparty_name,omitempty"`
	
	CounterpartyAccount *string `json:"counterparty_account,omitempty"`
	
	BankReference *string `json:"bank_reference,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankStatementLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.BankStatementId != nil {
		hasUpdate = true
	}
	
	if r.LineNumber != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
		hasUpdate = true
	}
	
	if r.ValueDate != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Reference != nil {
		hasUpdate = true
	}
	
	if r.CounterpartyName != nil {
		hasUpdate = true
	}
	
	if r.CounterpartyAccount != nil {
		hasUpdate = true
	}
	
	if r.BankReference != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BankStatementLinesListResponse represents a paginated list of bank_statement_lines records
type BankStatementLinesListResponse struct {
	Items      []*BankStatementLinesResponse `json:"items"`
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
