package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyRewardsResponse represents a loyalty_rewards response
type LoyaltyRewardsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	RewardCode string `json:"reward_code"`
	
	RewardName string `json:"reward_name"`
	
	Description *string `json:"description"`
	
	RewardType string `json:"reward_type"`
	
	PointsCost int64 `json:"points_cost"`
	
	RewardValue *float64 `json:"reward_value"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	IsActive *bool `json:"is_active"`
	
	AvailableFrom *time.Time `json:"available_from"`
	
	AvailableTo *time.Time `json:"available_to"`
	
	TotalAvailable *int64 `json:"total_available"`
	
	TotalRedeemed *int64 `json:"total_redeemed"`
	
	MaxRedemptionsPerCustomer *int64 `json:"max_redemptions_per_customer"`
	
	MinimumTierLevel *int64 `json:"minimum_tier_level"`
	
	TierIds json.RawMessage `json:"tier_ids"`
	
	ImageUrl *string `json:"image_url"`
	
	ThumbnailUrl *string `json:"thumbnail_url"`
	
	Featured *bool `json:"featured"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsFeatured *bool `json:"is_featured"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	RedemptionInstructions *string `json:"redemption_instructions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	AvailableTo *string `json:"available_to"`
	
	TotalRedeemed *string `json:"total_redeemed"`
	
	(totalAvailable *string `json:"(total_available"`
	
}

// CreateLoyaltyRewardsRequest represents a request to create a loyalty_rewards
type CreateLoyaltyRewardsRequest struct {
	
	RewardCode string `json:"reward_code" validate:"required"`
	
	RewardName string `json:"reward_name" validate:"required"`
	
	Description *string `json:"description"`
	
	RewardType string `json:"reward_type" validate:"required"`
	
	PointsCost int64 `json:"points_cost" validate:"required"`
	
	RewardValue *float64 `json:"reward_value"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	DiscountAmount *float64 `json:"discount_amount"`
	
	ProductId *uuid.UUID `json:"product_id"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id"`
	
	IsActive *bool `json:"is_active"`
	
	AvailableFrom *time.Time `json:"available_from"`
	
	AvailableTo *time.Time `json:"available_to"`
	
	TotalAvailable *int64 `json:"total_available"`
	
	TotalRedeemed *int64 `json:"total_redeemed"`
	
	MaxRedemptionsPerCustomer *int64 `json:"max_redemptions_per_customer"`
	
	MinimumTierLevel *int64 `json:"minimum_tier_level"`
	
	TierIds json.RawMessage `json:"tier_ids"`
	
	ImageUrl *string `json:"image_url" validate:"url"`
	
	ThumbnailUrl *string `json:"thumbnail_url" validate:"url"`
	
	Featured *bool `json:"featured"`
	
	SortOrder *int64 `json:"sort_order"`
	
	IsFeatured *bool `json:"is_featured"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	RedemptionInstructions *string `json:"redemption_instructions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	AvailableTo *string `json:"available_to"`
	
	TotalRedeemed *string `json:"total_redeemed"`
	
	(totalAvailable *string `json:"(total_available"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyRewardsRequest) Validate() error {
	
	if r.RewardCode == "" {
		return fmt.Errorf("reward_code is required")
	}
	
	if r.RewardName == "" {
		return fmt.Errorf("reward_name is required")
	}
	
	if r.RewardType == "" {
		return fmt.Errorf("reward_type is required")
	}
	
	if r.PointsCost == 0 {
		return fmt.Errorf("points_cost is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyRewardsRequest represents a request to update a loyalty_rewards
type UpdateLoyaltyRewardsRequest struct {
	
	RewardCode *string `json:"reward_code,omitempty" validate:"omitempty,required"`
	
	RewardName *string `json:"reward_name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	RewardType *string `json:"reward_type,omitempty" validate:"omitempty,required"`
	
	PointsCost *int64 `json:"points_cost,omitempty" validate:"omitempty,required"`
	
	RewardValue *float64 `json:"reward_value,omitempty"`
	
	DiscountPercentage *float64 `json:"discount_percentage,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty"`
	
	ProductId *uuid.UUID `json:"product_id,omitempty"`
	
	ProductVariantId *uuid.UUID `json:"product_variant_id,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	AvailableFrom *time.Time `json:"available_from,omitempty"`
	
	AvailableTo *time.Time `json:"available_to,omitempty"`
	
	TotalAvailable *int64 `json:"total_available,omitempty"`
	
	TotalRedeemed *int64 `json:"total_redeemed,omitempty"`
	
	MaxRedemptionsPerCustomer *int64 `json:"max_redemptions_per_customer,omitempty"`
	
	MinimumTierLevel *int64 `json:"minimum_tier_level,omitempty"`
	
	TierIds *json.RawMessage `json:"tier_ids,omitempty"`
	
	ImageUrl *string `json:"image_url,omitempty" validate:"omitempty,url"`
	
	ThumbnailUrl *string `json:"thumbnail_url,omitempty" validate:"omitempty,url"`
	
	Featured *bool `json:"featured,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	IsFeatured *bool `json:"is_featured,omitempty"`
	
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`
	
	RedemptionInstructions *string `json:"redemption_instructions,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	AvailableTo *string `json:"available_to,omitempty"`
	
	TotalRedeemed *string `json:"total_redeemed,omitempty"`
	
	(totalAvailable *string `json:"(total_available,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyRewardsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.RewardCode != nil {
		hasUpdate = true
	}
	
	if r.RewardName != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.RewardType != nil {
		hasUpdate = true
	}
	
	if r.PointsCost != nil {
		hasUpdate = true
	}
	
	if r.RewardValue != nil {
		hasUpdate = true
	}
	
	if r.DiscountPercentage != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.ProductId != nil {
		hasUpdate = true
	}
	
	if r.ProductVariantId != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.AvailableFrom != nil {
		hasUpdate = true
	}
	
	if r.AvailableTo != nil {
		hasUpdate = true
	}
	
	if r.TotalAvailable != nil {
		hasUpdate = true
	}
	
	if r.TotalRedeemed != nil {
		hasUpdate = true
	}
	
	if r.MaxRedemptionsPerCustomer != nil {
		hasUpdate = true
	}
	
	if r.MinimumTierLevel != nil {
		hasUpdate = true
	}
	
	if r.TierIds != nil {
		hasUpdate = true
	}
	
	if r.ImageUrl != nil {
		hasUpdate = true
	}
	
	if r.ThumbnailUrl != nil {
		hasUpdate = true
	}
	
	if r.Featured != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.IsFeatured != nil {
		hasUpdate = true
	}
	
	if r.TermsAndConditions != nil {
		hasUpdate = true
	}
	
	if r.RedemptionInstructions != nil {
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
	
	if r.AvailableTo != nil {
		hasUpdate = true
	}
	
	if r.TotalRedeemed != nil {
		hasUpdate = true
	}
	
	if r.(totalAvailable != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// LoyaltyRewardsListResponse represents a paginated list of loyalty_rewards records
type LoyaltyRewardsListResponse struct {
	Items      []*LoyaltyRewardsResponse `json:"items"`
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
