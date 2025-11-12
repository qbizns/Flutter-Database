package promotion_usage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PromotionUsageResponse represents a promotion_usage response
type PromotionUsageResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	PromotionId uuid.UUID `json:"promotion_id"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	DiscountAmount float64 `json:"discount_amount"`
	
	UsedAt *time.Time `json:"used_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreatePromotionUsageRequest represents a request to create a promotion_usage
type CreatePromotionUsageRequest struct {
	
	PromotionId uuid.UUID `json:"promotion_id" validate:"required"`
	
	SaleId *uuid.UUID `json:"sale_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	DiscountAmount float64 `json:"discount_amount" validate:"required"`
	
	UsedAt *time.Time `json:"used_at"`
	
}

// Validate validates the create request
func (r *CreatePromotionUsageRequest) Validate() error {
	
	if r.PromotionId == uuid.Nil {
		return fmt.Errorf("promotion_id is required")
	}
	
	if r.DiscountAmount == nil {
		return fmt.Errorf("discount_amount is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePromotionUsageRequest represents a request to update a promotion_usage
type UpdatePromotionUsageRequest struct {
	
	PromotionId *uuid.UUID `json:"promotion_id,omitempty" validate:"omitempty,required"`
	
	SaleId *uuid.UUID `json:"sale_id,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	DiscountAmount *float64 `json:"discount_amount,omitempty" validate:"omitempty,required"`
	
	UsedAt *time.Time `json:"used_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePromotionUsageRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PromotionId != nil {
		hasUpdate = true
	}
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.DiscountAmount != nil {
		hasUpdate = true
	}
	
	if r.UsedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PromotionUsageListResponse represents a paginated list of promotion_usage records
type PromotionUsageListResponse struct {
	Items      []*PromotionUsageResponse `json:"items"`
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
