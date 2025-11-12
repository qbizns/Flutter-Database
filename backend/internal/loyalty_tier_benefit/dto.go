package loyalty_tier_benefit

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LoyaltyTierBenefitsResponse represents a loyalty_tier_benefits response
type LoyaltyTierBenefitsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	TierId uuid.UUID `json:"tier_id"`
	
	BenefitCode string `json:"benefit_code"`
	
	BenefitName string `json:"benefit_name"`
	
	BenefitDescription *string `json:"benefit_description"`
	
	BenefitType string `json:"benefit_type"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	DiscountType *string `json:"discount_type"`
	
	IsActive *bool `json:"is_active"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Icon *string `json:"icon"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateLoyaltyTierBenefitsRequest represents a request to create a loyalty_tier_benefits
type CreateLoyaltyTierBenefitsRequest struct {
	
	TierId uuid.UUID `json:"tier_id" validate:"required"`
	
	BenefitCode string `json:"benefit_code" validate:"required"`
	
	BenefitName string `json:"benefit_name" validate:"required"`
	
	BenefitDescription *string `json:"benefit_description"`
	
	BenefitType string `json:"benefit_type" validate:"required"`
	
	DiscountValue *float64 `json:"discount_value"`
	
	DiscountType *string `json:"discount_type"`
	
	IsActive *bool `json:"is_active"`
	
	SortOrder *int64 `json:"sort_order"`
	
	Icon *string `json:"icon"`
	
	TermsAndConditions *string `json:"terms_and_conditions"`
	
	Metadata json.RawMessage `json:"metadata"`
	
}

// Validate validates the create request
func (r *CreateLoyaltyTierBenefitsRequest) Validate() error {
	
	if r.TierId == uuid.Nil {
		return fmt.Errorf("tier_id is required")
	}
	
	if r.BenefitCode == "" {
		return fmt.Errorf("benefit_code is required")
	}
	
	if r.BenefitName == "" {
		return fmt.Errorf("benefit_name is required")
	}
	
	if r.BenefitType == "" {
		return fmt.Errorf("benefit_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateLoyaltyTierBenefitsRequest represents a request to update a loyalty_tier_benefits
type UpdateLoyaltyTierBenefitsRequest struct {
	
	TierId *uuid.UUID `json:"tier_id,omitempty" validate:"omitempty,required"`
	
	BenefitCode *string `json:"benefit_code,omitempty" validate:"omitempty,required"`
	
	BenefitName *string `json:"benefit_name,omitempty" validate:"omitempty,required"`
	
	BenefitDescription *string `json:"benefit_description,omitempty"`
	
	BenefitType *string `json:"benefit_type,omitempty" validate:"omitempty,required"`
	
	DiscountValue *float64 `json:"discount_value,omitempty"`
	
	DiscountType *string `json:"discount_type,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	SortOrder *int64 `json:"sort_order,omitempty"`
	
	Icon *string `json:"icon,omitempty"`
	
	TermsAndConditions *string `json:"terms_and_conditions,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateLoyaltyTierBenefitsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.TierId != nil {
		hasUpdate = true
	}
	
	if r.BenefitCode != nil {
		hasUpdate = true
	}
	
	if r.BenefitName != nil {
		hasUpdate = true
	}
	
	if r.BenefitDescription != nil {
		hasUpdate = true
	}
	
	if r.BenefitType != nil {
		hasUpdate = true
	}
	
	if r.DiscountValue != nil {
		hasUpdate = true
	}
	
	if r.DiscountType != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.SortOrder != nil {
		hasUpdate = true
	}
	
	if r.Icon != nil {
		hasUpdate = true
	}
	
	if r.TermsAndConditions != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// LoyaltyTierBenefitsListResponse represents a paginated list of loyalty_tier_benefits records
type LoyaltyTierBenefitsListResponse struct {
	Items      []*LoyaltyTierBenefitsResponse `json:"items"`
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
