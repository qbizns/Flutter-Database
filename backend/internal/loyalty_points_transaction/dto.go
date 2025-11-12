package loyalty_points_transaction

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyPointsTransactionsResponse represents a loyalty_points_transactions response
type LoyaltyPointsTransactionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	TransactionType string `json:"transaction_type"`
	
	Points int64 `json:"points"`
	
	BalanceAfter int64 `json:"balance_after"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	RedemptionId *uuid.UUID `json:"redemption_id"`
	
	PointsRuleId *uuid.UUID `json:"points_rule_id"`
	
	Description *string `json:"description"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	TransactionDate *time.Time `json:"transaction_date"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreateLoyaltyPointsTransactionsRequest represents a request to create a loyalty_points_transactions
type CreateLoyaltyPointsTransactionsRequest struct {
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	TransactionType string `json:"transaction_type" validate:"required"`
	
	Points int64 `json:"points" validate:"required"`
	
	BalanceAfter int64 `json:"balance_after" validate:"required"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	RedemptionId *uuid.UUID `json:"redemption_id"`
	
	PointsRuleId *uuid.UUID `json:"points_rule_id"`
	
	Description *string `json:"description"`
	
	Reason *string `json:"reason"`
	
	Notes *string `json:"notes"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	TransactionDate *time.Time `json:"transaction_date"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyPointsTransactionsRequest) Validate() error {
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	
	if r.TransactionType == "" {
		return fmt.Errorf("transaction_type is required")
	}
	
	if r.Points == 0 {
		return fmt.Errorf("points is required")
	}
	
	if r.BalanceAfter == 0 {
		return fmt.Errorf("balance_after is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyPointsTransactionsRequest represents a request to update a loyalty_points_transactions
type UpdateLoyaltyPointsTransactionsRequest struct {
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	TransactionType *string `json:"transaction_type,omitempty" validate:"omitempty,required"`
	
	Points *int64 `json:"points,omitempty" validate:"omitempty,required"`
	
	BalanceAfter *int64 `json:"balance_after,omitempty" validate:"omitempty,required"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	RedemptionId *uuid.UUID `json:"redemption_id,omitempty"`
	
	PointsRuleId *uuid.UUID `json:"points_rule_id,omitempty"`
	
	Description *string `json:"description,omitempty"`
	
	Reason *string `json:"reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	TransactionDate *time.Time `json:"transaction_date,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyPointsTransactionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.TransactionType != nil {
		hasUpdate = true
	}
	
	if r.Points != nil {
		hasUpdate = true
	}
	
	if r.BalanceAfter != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.RedemptionId != nil {
		hasUpdate = true
	}
	
	if r.PointsRuleId != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.Reason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.ExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.TransactionDate != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// LoyaltyPointsTransactionsListResponse represents a paginated list of loyalty_points_transactions records
type LoyaltyPointsTransactionsListResponse struct {
	Items      []*LoyaltyPointsTransactionsResponse `json:"items"`
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
