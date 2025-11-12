package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyRedemptionsResponse represents a loyalty_redemptions response
type LoyaltyRedemptionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	RewardId uuid.UUID `json:"reward_id"`
	
	RedemptionNumber string `json:"redemption_number"`
	
	RedemptionDate *time.Time `json:"redemption_date"`
	
	PointsRedeemed int64 `json:"points_redeemed"`
	
	Status *string `json:"status"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	UsedDate *time.Time `json:"used_date"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	FulfillmentStatus *string `json:"fulfillment_status"`
	
	FulfillmentNotes *string `json:"fulfillment_notes"`
	
	FulfilledBy *uuid.UUID `json:"fulfilled_by"`
	
	FulfilledAt *time.Time `json:"fulfilled_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateLoyaltyRedemptionsRequest represents a request to create a loyalty_redemptions
type CreateLoyaltyRedemptionsRequest struct {
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	RewardId uuid.UUID `json:"reward_id" validate:"required"`
	
	RedemptionNumber string `json:"redemption_number" validate:"required"`
	
	RedemptionDate *time.Time `json:"redemption_date"`
	
	PointsRedeemed int64 `json:"points_redeemed" validate:"required"`
	
	Status *string `json:"status"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	UsedDate *time.Time `json:"used_date"`
	
	ExpiryDate *time.Time `json:"expiry_date"`
	
	FulfillmentStatus *string `json:"fulfillment_status"`
	
	FulfillmentNotes *string `json:"fulfillment_notes"`
	
	FulfilledBy *uuid.UUID `json:"fulfilled_by"`
	
	FulfilledAt *time.Time `json:"fulfilled_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyRedemptionsRequest) Validate() error {
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	
	if r.RewardId == uuid.Nil {
		return fmt.Errorf("reward_id is required")
	}
	
	if r.RedemptionNumber == "" {
		return fmt.Errorf("redemption_number is required")
	}
	
	if r.PointsRedeemed == 0 {
		return fmt.Errorf("points_redeemed is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyRedemptionsRequest represents a request to update a loyalty_redemptions
type UpdateLoyaltyRedemptionsRequest struct {
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	RewardId *uuid.UUID `json:"reward_id,omitempty" validate:"omitempty,required"`
	
	RedemptionNumber *string `json:"redemption_number,omitempty" validate:"omitempty,required"`
	
	RedemptionDate *time.Time `json:"redemption_date,omitempty"`
	
	PointsRedeemed *int64 `json:"points_redeemed,omitempty" validate:"omitempty,required"`
	
	Status *string `json:"status,omitempty"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	UsedDate *time.Time `json:"used_date,omitempty"`
	
	ExpiryDate *time.Time `json:"expiry_date,omitempty"`
	
	FulfillmentStatus *string `json:"fulfillment_status,omitempty"`
	
	FulfillmentNotes *string `json:"fulfillment_notes,omitempty"`
	
	FulfilledBy *uuid.UUID `json:"fulfilled_by,omitempty"`
	
	FulfilledAt *time.Time `json:"fulfilled_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyRedemptionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.RewardId != nil {
		hasUpdate = true
	}
	
	if r.RedemptionNumber != nil {
		hasUpdate = true
	}
	
	if r.RedemptionDate != nil {
		hasUpdate = true
	}
	
	if r.PointsRedeemed != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.UsedDate != nil {
		hasUpdate = true
	}
	
	if r.ExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.FulfillmentStatus != nil {
		hasUpdate = true
	}
	
	if r.FulfillmentNotes != nil {
		hasUpdate = true
	}
	
	if r.FulfilledBy != nil {
		hasUpdate = true
	}
	
	if r.FulfilledAt != nil {
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

// LoyaltyRedemptionsListResponse represents a paginated list of loyalty_redemptions records
type LoyaltyRedemptionsListResponse struct {
	Items      []*LoyaltyRedemptionsResponse `json:"items"`
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
