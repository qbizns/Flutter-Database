package store_credit_transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StoreCreditTransactionsResponse represents a store_credit_transactions response
type StoreCreditTransactionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	StoreCreditAccountId uuid.UUID `json:"store_credit_account_id"`
	
	TransactionType string `json:"transaction_type"`
	
	TransactionType *string `json:"transaction_type"`
	
	Amount float64 `json:"amount"`
	
	BalanceAfter float64 `json:"balance_after"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	PaymentId *uuid.UUID `json:"payment_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Notes *string `json:"notes"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateStoreCreditTransactionsRequest represents a request to create a store_credit_transactions
type CreateStoreCreditTransactionsRequest struct {
	
	StoreCreditAccountId uuid.UUID `json:"store_credit_account_id" validate:"required"`
	
	TransactionType string `json:"transaction_type" validate:"required"`
	
	TransactionType *string `json:"transaction_type"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	BalanceAfter float64 `json:"balance_after" validate:"required"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	PaymentId *uuid.UUID `json:"payment_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	Notes *string `json:"notes"`
	
}

// Validate validates the create request
func (r *CreateStoreCreditTransactionsRequest) Validate() error {
	
	if r.StoreCreditAccountId == uuid.Nil {
		return fmt.Errorf("store_credit_account_id is required")
	}
	
	if r.TransactionType == "" {
		return fmt.Errorf("transaction_type is required")
	}
	
	if r.Amount == nil {
		return fmt.Errorf("amount is required")
	}
	
	if r.BalanceAfter == nil {
		return fmt.Errorf("balance_after is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateStoreCreditTransactionsRequest represents a request to update a store_credit_transactions
type UpdateStoreCreditTransactionsRequest struct {
	
	StoreCreditAccountId *uuid.UUID `json:"store_credit_account_id,omitempty" validate:"omitempty,required"`
	
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,required"`
	
	TransactionType *string `json:"transaction_type,omitempty"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	BalanceAfter *float64 `json:"balance_after,omitempty" validate:"omitempty,required"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	PaymentId *uuid.UUID `json:"payment_id,omitempty"`
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateStoreCreditTransactionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.StoreCreditAccountId != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.BalanceAfter != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.PaymentId != nil {
		hasUpdate = true
	}
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.LocationId != nil {
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

// StoreCreditTransactionsListResponse represents a paginated list of store_credit_transactions records
type StoreCreditTransactionsListResponse struct {
	Items      []*StoreCreditTransactionsResponse `json:"items"`
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
