package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BankAccountsResponse represents a bank_accounts response
type BankAccountsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ChartAccountId uuid.UUID `json:"chart_account_id"`
	
	BankName string `json:"bank_name"`
	
	AccountNumber string `json:"account_number"`
	
	AccountType *string `json:"account_type"`
	
	RoutingNumber *string `json:"routing_number"`
	
	SwiftCode *string `json:"swift_code"`
	
	CurrencyCode *string `json:"currency_code"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	StatementBalance *float64 `json:"statement_balance"`
	
	LastStatementDate *time.Time `json:"last_statement_date"`
	
	IsActive *bool `json:"is_active"`
	
	OnlineBankingEnabled *bool `json:"online_banking_enabled"`
	
	LastSyncDate *time.Time `json:"last_sync_date"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateBankAccountsRequest represents a request to create a bank_accounts
type CreateBankAccountsRequest struct {
	
	ChartAccountId uuid.UUID `json:"chart_account_id" validate:"required"`
	
	BankName string `json:"bank_name" validate:"required"`
	
	AccountNumber string `json:"account_number" validate:"required"`
	
	AccountType *string `json:"account_type"`
	
	RoutingNumber *string `json:"routing_number"`
	
	SwiftCode *string `json:"swift_code"`
	
	CurrencyCode *string `json:"currency_code"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	StatementBalance *float64 `json:"statement_balance"`
	
	LastStatementDate *time.Time `json:"last_statement_date"`
	
	IsActive *bool `json:"is_active"`
	
	OnlineBankingEnabled *bool `json:"online_banking_enabled"`
	
	LastSyncDate *time.Time `json:"last_sync_date"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateBankAccountsRequest) Validate() error {
	
	if r.ChartAccountId == uuid.Nil {
		return fmt.Errorf("chart_account_id is required")
	}
	
	if r.BankName == "" {
		return fmt.Errorf("bank_name is required")
	}
	
	if r.AccountNumber == "" {
		return fmt.Errorf("account_number is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateBankAccountsRequest represents a request to update a bank_accounts
type UpdateBankAccountsRequest struct {
	
	ChartAccountId *uuid.UUID `json:"chart_account_id,omitempty" validate:"omitempty,required"`
	
	BankName *string `json:"bank_name,omitempty" validate:"omitempty,required"`
	
	AccountNumber *string `json:"account_number,omitempty" validate:"omitempty,required"`
	
	AccountType *string `json:"account_type,omitempty"`
	
	RoutingNumber *string `json:"routing_number,omitempty"`
	
	SwiftCode *string `json:"swift_code,omitempty"`
	
	CurrencyCode *string `json:"currency_code,omitempty"`
	
	CurrentBalance *float64 `json:"current_balance,omitempty"`
	
	StatementBalance *float64 `json:"statement_balance,omitempty"`
	
	LastStatementDate *time.Time `json:"last_statement_date,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	OnlineBankingEnabled *bool `json:"online_banking_enabled,omitempty"`
	
	LastSyncDate *time.Time `json:"last_sync_date,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateBankAccountsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ChartAccountId != nil {
		hasUpdate = true
	}
	
	if r.BankName != nil {
		hasUpdate = true
	}
	
	if r.AccountNumber != nil {
		hasUpdate = true
	}
	
	if r.AccountType != nil {
		hasUpdate = true
	}
	
	if r.RoutingNumber != nil {
		hasUpdate = true
	}
	
	if r.SwiftCode != nil {
		hasUpdate = true
	}
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.CurrentBalance != nil {
		hasUpdate = true
	}
	
	if r.StatementBalance != nil {
		hasUpdate = true
	}
	
	if r.LastStatementDate != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.OnlineBankingEnabled != nil {
		hasUpdate = true
	}
	
	if r.LastSyncDate != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// BankAccountsListResponse represents a paginated list of bank_accounts records
type BankAccountsListResponse struct {
	Items      []*BankAccountsResponse `json:"items"`
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
