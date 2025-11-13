package loyalty_points_rule

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyPointsRulesResponse represents a loyalty_points_rules response
type LoyaltyPointsRulesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	RuleCode string `json:"rule_code"`
	
	RuleName string `json:"rule_name"`
	
	Description *string `json:"description"`
	
	RuleType string `json:"rule_type"`
	
	PointsPerAmount *float64 `json:"points_per_amount"`
	
	FixedPoints *int64 `json:"fixed_points"`
	
	Multiplier *float64 `json:"multiplier"`
	
	AppliesTo *string `json:"applies_to"`
	
	ApplicableProductIds json.RawMessage `json:"applicable_product_ids"`
	
	ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids"`
	
	ApplicableTierIds json.RawMessage `json:"applicable_tier_ids"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount"`
	
	MaximumPointsPerTransaction *int64 `json:"maximum_points_per_transaction"`
	
	MaximumPointsPerDay *int64 `json:"maximum_points_per_day"`
	
	MaximumPointsPerMonth *int64 `json:"maximum_points_per_month"`
	
	StartDate *time.Time `json:"start_date"`
	
	EndDate *time.Time `json:"end_date"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateLoyaltyPointsRulesRequest represents a request to create a loyalty_points_rules
type CreateLoyaltyPointsRulesRequest struct {
	
	RuleCode string `json:"rule_code" validate:"required"`
	
	RuleName string `json:"rule_name" validate:"required"`
	
	Description *string `json:"description"`
	
	RuleType string `json:"rule_type" validate:"required"`
	
	PointsPerAmount *float64 `json:"points_per_amount"`
	
	FixedPoints *int64 `json:"fixed_points"`
	
	Multiplier *float64 `json:"multiplier"`
	
	AppliesTo *string `json:"applies_to"`
	
	// Duplicate removed: ApplicableProductIds json.RawMessage `json:"applicable_product_ids"`
	
	// Duplicate removed: ApplicableCategoryIds json.RawMessage `json:"applicable_category_ids"`
	
	// Duplicate removed: ApplicableTierIds json.RawMessage `json:"applicable_tier_ids"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount"`
	
	MaximumPointsPerTransaction *int64 `json:"maximum_points_per_transaction"`
	
	MaximumPointsPerDay *int64 `json:"maximum_points_per_day"`
	
	MaximumPointsPerMonth *int64 `json:"maximum_points_per_month"`
	
	StartDate *time.Time `json:"start_date"`
	
	EndDate *time.Time `json:"end_date"`
	
	IsActive *bool `json:"is_active"`
	
	Priority *int64 `json:"priority"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	// Duplicate removed: Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyPointsRulesRequest) Validate() error {
	
	if r.RuleCode == "" {
		return fmt.Errorf("rule_code is required")
	}
	
	if r.RuleName == "" {
		return fmt.Errorf("rule_name is required")
	}
	
	if r.RuleType == "" {
		return fmt.Errorf("rule_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyPointsRulesRequest represents a request to update a loyalty_points_rules
type UpdateLoyaltyPointsRulesRequest struct {
	
	RuleCode *string `json:"rule_code,omitempty" validate:"omitempty,required"`
	
	RuleName *string `json:"rule_name,omitempty" validate:"omitempty,required"`
	
	Description *string `json:"description,omitempty"`
	
	RuleType *string `json:"rule_type,omitempty" validate:"omitempty,required"`
	
	PointsPerAmount *float64 `json:"points_per_amount,omitempty"`
	
	FixedPoints *int64 `json:"fixed_points,omitempty"`
	
	Multiplier *float64 `json:"multiplier,omitempty"`
	
	AppliesTo *string `json:"applies_to,omitempty"`
	
	ApplicableProductIds *json.RawMessage `json:"applicable_product_ids,omitempty"`
	
	ApplicableCategoryIds *json.RawMessage `json:"applicable_category_ids,omitempty"`
	
	ApplicableTierIds *json.RawMessage `json:"applicable_tier_ids,omitempty"`
	
	MinimumPurchaseAmount *float64 `json:"minimum_purchase_amount,omitempty"`
	
	MaximumPointsPerTransaction *int64 `json:"maximum_points_per_transaction,omitempty"`
	
	MaximumPointsPerDay *int64 `json:"maximum_points_per_day,omitempty"`
	
	MaximumPointsPerMonth *int64 `json:"maximum_points_per_month,omitempty"`
	
	StartDate *time.Time `json:"start_date,omitempty"`
	
	EndDate *time.Time `json:"end_date,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	Priority *int64 `json:"priority,omitempty"`
	
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyPointsRulesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.RuleCode != nil {
		hasUpdate = true
	}
	
	if r.RuleName != nil {
		hasUpdate = true
	}
	
	if r.Description != nil {
		hasUpdate = true
	}
	
	if r.RuleType != nil {
		hasUpdate = true
	}
	
	if r.PointsPerAmount != nil {
		hasUpdate = true
	}
	
	if r.FixedPoints != nil {
		hasUpdate = true
	}
	
	if r.Multiplier != nil {
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
	
	if r.ApplicableTierIds != nil {
		hasUpdate = true
	}
	
	if r.MinimumPurchaseAmount != nil {
		hasUpdate = true
	}
	
	if r.MaximumPointsPerTransaction != nil {
		hasUpdate = true
	}
	
	if r.MaximumPointsPerDay != nil {
		hasUpdate = true
	}
	
	if r.MaximumPointsPerMonth != nil {
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

// LoyaltyPointsRulesListResponse represents a paginated list of loyalty_points_rules records
type LoyaltyPointsRulesListResponse struct {
	Items      []*LoyaltyPointsRulesResponse `json:"items"`
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
