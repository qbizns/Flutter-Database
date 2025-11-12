package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GiftCardTransactionsResponse represents a gift_card_transactions response
type GiftCardTransactionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	GiftCardId uuid.UUID `json:"gift_card_id"`
	
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

// CreateGiftCardTransactionsRequest represents a request to create a gift_card_transactions
type CreateGiftCardTransactionsRequest struct {
	
	GiftCardId uuid.UUID `json:"gift_card_id" validate:"required"`
	
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
func (r *CreateGiftCardTransactionsRequest) Validate() error {
	
	if r.GiftCardId == uuid.Nil {
		return fmt.Errorf("gift_card_id is required")
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

// UpdateGiftCardTransactionsRequest represents a request to update a gift_card_transactions
type UpdateGiftCardTransactionsRequest struct {
	
	GiftCardId *uuid.UUID `json:"gift_card_id,omitempty" validate:"omitempty,required"`
	
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
func (r *UpdateGiftCardTransactionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.GiftCardId != nil {
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

// GiftCardTransactionsListResponse represents a paginated list of gift_card_transactions records
type GiftCardTransactionsListResponse struct {
	Items      []*GiftCardTransactionsResponse `json:"items"`
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
