package promotion

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PromotionsResponse represents a promotions response
type PromotionsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PromotionCode string `json:"promotion_code"`
	
	Name string `json:"name"`
	
	Description *string `json:"description"`
	
	PromotionType string `json:"promotion_type"`
	
	DiscountValue float64 `json:"discount_value"`
	
	AppliesTo *string `json:"applies_to"`
	
	ApplicableProductIds json.RawMessage `json:"applicable_product_ids"`
	
	ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount"`
	
	MinimumQuantity *int64 `json:"minimum_quantity"`
	
	BuyQuantity *int64 `json:"buy_quantity"`
	
	GetQuantity *int64 `json:"get_quantity"`
	
	GetDiscountPercentage *float64 `json:"get_discount_percentage"`
	
	MaxUsesTotal *int64 `json:"max_uses_total"`
	
	MaxUsesPerCustomer *int64 `json:"max_uses_per_customer"`
	
	CurrentUses *int64 `json:"current_uses"`
	
	StartDate time.Time `json:"start_date"`
	
	EndDate *time.Time `json:"end_date"`
	
	IsActive *bool `json:"is_active"`
	
	IsCombinable *bool `json:"is_combinable"`
	
	Priority *int64 `json:"priority"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreatePromotionsRequest represents a request to create a promotions
type CreatePromotionsRequest struct {
	
	PromotionCode string `json:"promotion_code" validate:"required"`
	
	Name string `json:"name" validate:"required"`
	
	Description *string `json:"description"`
	
	PromotionType string `json:"promotion_type" validate:"required"`
	
	DiscountValue float64 `json:"discount_value" validate:"required"`
	
	AppliesTo *string `json:"applies_to"`
	
	ApplicableProductIds json.RawMessage `json:"applicable_product_ids"`
	
	ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount"`
	
	MinimumQuantity *int64 `json:"minimum_quantity"`
	
	BuyQuantity *int64 `json:"buy_quantity"`
	
	GetQuantity *int64 `json:"get_quantity"`
	
	GetDiscountPercentage *float64 `json:"get_discount_percentage"`
	
	MaxUsesTotal *int64 `json:"max_uses_total"`
	
	MaxUsesPerCustomer *int64 `json:"max_uses_per_customer"`
	
	CurrentUses *int64 `json:"current_uses"`
	
	StartDate time.Time `json:"start_date" validate:"required"`
	
	EndDate *time.Time `json:"end_date"`
	
	IsActive *bool `json:"is_active"`
	
	IsCombinable *bool `json:"is_combinable"`
	
	Priority *int64 `json:"priority"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreatePromotionsRequest) Validate() error {
	
	if r.PromotionCode == "" {
		return fmt.Errorf("promotion_code is required")
	}
	
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	
	if r.PromotionType == "" {
		return fmt.Errorf("promotion_type is required")
	}
	
	if r.DiscountValue == nil {
		return fmt.Errorf("discount_value is required")
	}
	
	if r.StartDate == nil {
		return fmt.Errorf("start_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePromotionsRequest represents a request to update a promotions
type UpdatePromotionsRequest struct {
	
	PromotionCode *string `json:"promotion_code,omitempty" validate:"omitempty,required"`
	
	Name *string `json:"name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	PromotionType *string `json:"promotion_type,omitempty" validate:"omitempty,required"`
	
	DiscountValue *float64 `json:"discount_value,omitempty" validate:"omitempty,required"`
	
	AppliesTo *string `json:"applies_to,omitempty"`
	
	ApplicableProductIds *json.RawMessage `json:"applicable_product_ids,omitempty"`
	
	ApplicableCategoryIds *json.RawMessage `json:"applicable_category_ids,omitempty"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount,omitempty"`
	
	MinimumQuantity *int64 `json:"minimum_quantity,omitempty"`
	
	BuyQuantity *int64 `json:"buy_quantity,omitempty"`
	
	GetQuantity *int64 `json:"get_quantity,omitempty"`
	
	GetDiscountPercentage *float64 `json:"get_discount_percentage,omitempty"`
	
	MaxUsesTotal *int64 `json:"max_uses_total,omitempty"`
	
	MaxUsesPerCustomer *int64 `json:"max_uses_per_customer,omitempty"`
	
	CurrentUses *int64 `json:"current_uses,omitempty"`
	
	StartDate *time.Time `json:"start_date,omitempty" validate:"omitempty,required"`
	
	EndDate *time.Time `json:"end_date,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsCombinable *bool `json:"is_combinable,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePromotionsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PromotionCode != nil {
		hasUpdate = true
	}
	
	if r.Name != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.PromotionType != nil {
		hasUpdate = true
	}
	
	if r.DiscountValue != nil {
		hasUpdate = true
	}
	
	if r.AppliesTo != nil {
		hasUpdate = true
	}
	
	if r.ApplicableProductIds != nil {
		hasUpdate = true
	}
	
	if r.ApplicableCategoryIds != nil {
		hasUpdate = true
	}
	
	if r.MinimumPurchaseAmount != nil {
		hasUpdate = true
	}
	
	if r.MinimumQuantity != nil {
		hasUpdate = true
	}
	
	if r.BuyQuantity != nil {
		hasUpdate = true
	}
	
	if r.GetQuantity != nil {
		hasUpdate = true
	}
	
	if r.GetDiscountPercentage != nil {
		hasUpdate = true
	}
	
	if r.MaxUsesTotal != nil {
		hasUpdate = true
	}
	
	if r.MaxUsesPerCustomer != nil {
		hasUpdate = true
	}
	
	if r.CurrentUses != nil {
		hasUpdate = true
	}
	
	if r.StartDate != nil {
		hasUpdate = true
	}
	
	if r.EndDate != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsCombinable != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.TermsAndConditions != nil {
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

// PromotionsListResponse represents a paginated list of promotions records
type PromotionsListResponse struct {
	Items      []*PromotionsResponse `json:"items"`
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
