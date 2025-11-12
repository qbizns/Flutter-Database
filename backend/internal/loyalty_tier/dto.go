package loyalty_tier

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyTiersResponse represents a loyalty_tiers response
type LoyaltyTiersResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	TierCode string `json:"tier_code"`
	
	TierName string `json:"tier_name"`
	
	TierLevel int64 `json:"tier_level"`
	
	Description *string `json:"description"`
	
	PointsThreshold int64 `json:"points_threshold"`
	
	AnnualSpendThreshold *float64 `json:"annual_spend_threshold"`
	
	PurchaseCountThreshold *int64 `json:"purchase_count_threshold"`
	
	PointsMultiplier *float64 `json:"points_multiplier"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	TierColor *string `json:"tier_color"`
	
	TierIcon *string `json:"tier_icon"`
	
	BadgeImageUrl *string `json:"badge_image_url"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	PointsThreshold *string `json:"points_threshold"`
	
	(annualSpendThreshold *string `json:"(annual_spend_threshold"`
	
	(purchaseCountThreshold *string `json:"(purchase_count_threshold"`
	
}

// CreateLoyaltyTiersRequest represents a request to create a loyalty_tiers
type CreateLoyaltyTiersRequest struct {
	
	TierCode string `json:"tier_code" validate:"required"`
	
	TierName string `json:"tier_name" validate:"required"`
	
	TierLevel int64 `json:"tier_level" validate:"required"`
	
	Description *string `json:"description"`
	
	PointsThreshold int64 `json:"points_threshold" validate:"required"`
	
	AnnualSpendThreshold *float64 `json:"annual_spend_threshold"`
	
	PurchaseCountThreshold *int64 `json:"purchase_count_threshold"`
	
	PointsMultiplier *float64 `json:"points_multiplier"`
	
	DiscountPercentage *float64 `json:"discount_percentage"`
	
	TierColor *string `json:"tier_color"`
	
	TierIcon *string `json:"tier_icon"`
	
	BadgeImageUrl *string `json:"badge_image_url" validate:"url"`
	
	IsActive *bool `json:"is_active"`
	
	IsDefault *bool `json:"is_default"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	PointsThreshold *string `json:"points_threshold"`
	
	(annualSpendThreshold *string `json:"(annual_spend_threshold"`
	
	(purchaseCountThreshold *string `json:"(purchase_count_threshold"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyTiersRequest) Validate() error {
	
	if r.TierCode == "" {
		return fmt.Errorf("tier_code is required")
	}
	
	if r.TierName == "" {
		return fmt.Errorf("tier_name is required")
	}
	
	if r.TierLevel == 0 {
		return fmt.Errorf("tier_level is required")
	}
	
	if r.PointsThreshold == 0 {
		return fmt.Errorf("points_threshold is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyTiersRequest represents a request to update a loyalty_tiers
type UpdateLoyaltyTiersRequest struct {
	
	TierCode *string `json:"tier_code,omitempty" validate:"omitempty,required"`
	
	TierName *string `json:"tier_name,omitempty" validate:"omitempty,required"`
	
	TierLevel *int64 `json:"tier_level,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	PointsThreshold *int64 `json:"points_threshold,omitempty" validate:"omitempty,required"`
	
	AnnualSpendThreshold *float64 `json:"annual_spend_threshold,omitempty"`
	
	PurchaseCountThreshold *int64 `json:"purchase_count_threshold,omitempty"`
	
	PointsMultiplier *float64 `json:"points_multiplier,omitempty"`
	
	DiscountPercentage *float64 `json:"discount_percentage,omitempty"`
	
	TierColor *string `json:"tier_color,omitempty"`
	
	TierIcon *string `json:"tier_icon,omitempty"`
	
	BadgeImageUrl *string `json:"badge_image_url,omitempty" validate:"omitempty,url"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsDefault *bool `json:"is_default,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	PointsThreshold *string `json:"points_threshold,omitempty"`
	
	(annualSpendThreshold *string `json:"(annual_spend_threshold,omitempty"`
	
	(purchaseCountThreshold *string `json:"(purchase_count_threshold,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyTiersRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TierCode != nil {
		hasUpdate = true
	}
	
	if r.TierName != nil {
		hasUpdate = true
	}
	
	if r.TierLevel != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.PointsThreshold != nil {
		hasUpdate = true
	}
	
	if r.AnnualSpendThreshold != nil {
		hasUpdate = true
	}
	
	if r.PurchaseCountThreshold != nil {
		hasUpdate = true
	}
	
	if r.PointsMultiplier != nil {
		hasUpdate = true
	}
	
	if r.DiscountPercentage != nil {
		hasUpdate = true
	}
	
	if r.TierColor != nil {
		hasUpdate = true
	}
	
	if r.TierIcon != nil {
		hasUpdate = true
	}
	
	if r.BadgeImageUrl != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsDefault != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
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
	
	if r.PointsThreshold != nil {
		hasUpdate = true
	}
	
	if r.(annualSpendThreshold != nil {
		hasUpdate = true
	}
	
	if r.(purchaseCountThreshold != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// LoyaltyTiersListResponse represents a paginated list of loyalty_tiers records
type LoyaltyTiersListResponse struct {
	Items      []*LoyaltyTiersResponse `json:"items"`
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
