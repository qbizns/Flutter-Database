package chart_of_account

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ChartOfAccountsResponse represents a chart_of_accounts response
type ChartOfAccountsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	AccountCode string `json:"account_code"`
	
	AccountNumber string `json:"account_number"`
	
	AccountName string `json:"account_name"`
	
	AccountTypeId uuid.UUID `json:"account_type_id"`
	
	AccountSubtypeId *uuid.UUID `json:"account_subtype_id"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id"`
	
	AccountLevel *int64 `json:"account_level"`
	
	AccountPath *string `json:"account_path"`
	
	IsActive *bool `json:"is_active"`
	
	IsSystemAccount *bool `json:"is_system_account"`
	
	IsHeaderAccount *bool `json:"is_header_account"`
	
	IsBankAccount *bool `json:"is_bank_account"`
	
	IsReconcilable *bool `json:"is_reconcilable"`
	
	DefaultTaxCode *string `json:"default_tax_code"`
	
	CurrencyCode *string `json:"currency_code"`
	
	OpeningBalance *float64 `json:"opening_balance"`
	
	OpeningBalanceDate *time.Time `json:"opening_balance_date"`
	
	CurrentDebitBalance *float64 `json:"current_debit_balance"`
	
	CurrentCreditBalance *float64 `json:"current_credit_balance"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	LastBalanceUpdate *time.Time `json:"last_balance_update"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateChartOfAccountsRequest represents a request to create a chart_of_accounts
type CreateChartOfAccountsRequest struct {
	
	AccountCode string `json:"account_code" validate:"required"`
	
	AccountNumber string `json:"account_number" validate:"required"`
	
	AccountName string `json:"account_name" validate:"required"`
	
	AccountTypeId uuid.UUID `json:"account_type_id" validate:"required"`
	
	AccountSubtypeId *uuid.UUID `json:"account_subtype_id"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id"`
	
	AccountLevel *int64 `json:"account_level"`
	
	AccountPath *string `json:"account_path"`
	
	IsActive *bool `json:"is_active"`
	
	IsSystemAccount *bool `json:"is_system_account"`
	
	IsHeaderAccount *bool `json:"is_header_account"`
	
	IsBankAccount *bool `json:"is_bank_account"`
	
	IsReconcilable *bool `json:"is_reconcilable"`
	
	DefaultTaxCode *string `json:"default_tax_code"`
	
	CurrencyCode *string `json:"currency_code"`
	
	OpeningBalance *float64 `json:"opening_balance"`
	
	OpeningBalanceDate *time.Time `json:"opening_balance_date"`
	
	CurrentDebitBalance *float64 `json:"current_debit_balance"`
	
	CurrentCreditBalance *float64 `json:"current_credit_balance"`
	
	CurrentBalance *float64 `json:"current_balance"`
	
	LastBalanceUpdate *time.Time `json:"last_balance_update"`
	
	Description *string `json:"description"`
	
	Notes *string `json:"notes"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateChartOfAccountsRequest) Validate() error {
	
	if r.AccountCode == "" {
		return fmt.Errorf("account_code is required")
	}
	
	if r.AccountNumber == "" {
		return fmt.Errorf("account_number is required")
	}
	
	if r.AccountName == "" {
		return fmt.Errorf("account_name is required")
	}
	
	if r.AccountTypeId == uuid.Nil {
		return fmt.Errorf("account_type_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateChartOfAccountsRequest represents a request to update a chart_of_accounts
type UpdateChartOfAccountsRequest struct {
	
	AccountCode *string `json:"account_code,omitempty" validate:"omitempty,required"`
	
	AccountNumber *string `json:"account_number,omitempty" validate:"omitempty,required"`
	
	AccountName *string `json:"account_name,omitempty" validate:"omitempty,required"`
	
	AccountTypeId *uuid.UUID `json:"account_type_id,omitempty" validate:"omitempty,required"`
	
	AccountSubtypeId *uuid.UUID `json:"account_subtype_id,omitempty"`
	
	ParentAccountId *uuid.UUID `json:"parent_account_id,omitempty"`
	
	AccountLevel *int64 `json:"account_level,omitempty"`
	
	AccountPath *string `json:"account_path,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsSystemAccount *bool `json:"is_system_account,omitempty"`
	
	IsHeaderAccount *bool `json:"is_header_account,omitempty"`
	
	IsBankAccount *bool `json:"is_bank_account,omitempty"`
	
	IsReconcilable *bool `json:"is_reconcilable,omitempty"`
	
	DefaultTaxCode *string `json:"default_tax_code,omitempty"`
	
	CurrencyCode *string `json:"currency_code,omitempty"`
	
	OpeningBalance *float64 `json:"opening_balance,omitempty"`
	
	OpeningBalanceDate *time.Time `json:"opening_balance_date,omitempty"`
	
	CurrentDebitBalance *float64 `json:"current_debit_balance,omitempty"`
	
	CurrentCreditBalance *float64 `json:"current_credit_balance,omitempty"`
	
	CurrentBalance *float64 `json:"current_balance,omitempty"`
	
	LastBalanceUpdate *time.Time `json:"last_balance_update,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateChartOfAccountsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.AccountCode != nil {
		hasUpdate = true
	}
	
	if r.AccountNumber != nil {
		hasUpdate = true
	}
	
	if r.AccountName != nil {
		hasUpdate = true
	}
	
	if r.AccountTypeId != nil {
		hasUpdate = true
	}
	
	if r.AccountSubtypeId != nil {
		hasUpdate = true
	}
	
	if r.ParentAccountId != nil {
		hasUpdate = true
	}
	
	if r.AccountLevel != nil {
		hasUpdate = true
	}
	
	if r.AccountPath != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsSystemAccount != nil {
		hasUpdate = true
	}
	
	if r.IsHeaderAccount != nil {
		hasUpdate = true
	}
	
	if r.IsBankAccount != nil {
		hasUpdate = true
	}
	
	if r.IsReconcilable != nil {
		hasUpdate = true
	}
	
	if r.DefaultTaxCode != nil {
		hasUpdate = true
	}
	
	if r.CurrencyCode != nil {
		hasUpdate = true
	}
	
	if r.OpeningBalance != nil {
		hasUpdate = true
	}
	
	if r.OpeningBalanceDate != nil {
		hasUpdate = true
	}
	
	if r.CurrentDebitBalance != nil {
		hasUpdate = true
	}
	
	if r.CurrentCreditBalance != nil {
		hasUpdate = true
	}
	
	if r.CurrentBalance != nil {
		hasUpdate = true
	}
	
	if r.LastBalanceUpdate != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
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

// ChartOfAccountsListResponse represents a paginated list of chart_of_accounts records
type ChartOfAccountsListResponse struct {
	Items      []*ChartOfAccountsResponse `json:"items"`
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
