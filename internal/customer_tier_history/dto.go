package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CustomerTierHistoryResponse represents a customer_tier_history response
type CustomerTierHistoryResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	CustomerId uuid.UUID `json:"customer_id"`
	
	TierId uuid.UUID `json:"tier_id"`
	
	PreviousTierId *uuid.UUID `json:"previous_tier_id"`
	
	ChangeType string `json:"change_type"`
	
	ChangeReason *string `json:"change_reason"`
	
	QualifyingPoints *int64 `json:"qualifying_points"`
	
	QualifyingSpend *float64 `json:"qualifying_spend"`
	
	QualifyingPurchases *int64 `json:"qualifying_purchases"`
	
	EffectiveDate time.Time `json:"effective_date"`
	
	ValidUntil *time.Time `json:"valid_until"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreateCustomerTierHistoryRequest represents a request to create a customer_tier_history
type CreateCustomerTierHistoryRequest struct {
	
	CustomerId uuid.UUID `json:"customer_id" validate:"required"`
	
	TierId uuid.UUID `json:"tier_id" validate:"required"`
	
	PreviousTierId *uuid.UUID `json:"previous_tier_id"`
	
	ChangeType string `json:"change_type" validate:"required"`
	
	ChangeReason *string `json:"change_reason"`
	
	QualifyingPoints *int64 `json:"qualifying_points"`
	
	QualifyingSpend *float64 `json:"qualifying_spend"`
	
	QualifyingPurchases *int64 `json:"qualifying_purchases"`
	
	EffectiveDate time.Time `json:"effective_date" validate:"required"`
	
	ValidUntil *time.Time `json:"valid_until"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateCustomerTierHistoryRequest) Validate() error {
	
	if r.CustomerId == uuid.Nil {
		return fmt.Errorf("customer_id is required")
	}
	
	if r.TierId == uuid.Nil {
		return fmt.Errorf("tier_id is required")
	}
	
	if r.ChangeType == "" {
		return fmt.Errorf("change_type is required")
	}
	
	if r.EffectiveDate == nil {
		return fmt.Errorf("effective_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateCustomerTierHistoryRequest represents a request to update a customer_tier_history
type UpdateCustomerTierHistoryRequest struct {
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty" validate:"omitempty,required"`
	
	TierId *uuid.UUID `json:"tier_id,omitempty" validate:"omitempty,required"`
	
	PreviousTierId *uuid.UUID `json:"previous_tier_id,omitempty"`
	
	ChangeType *string `json:"change_type,omitempty" validate:"omitempty,required"`
	
	ChangeReason *string `json:"change_reason,omitempty"`
	
	QualifyingPoints *int64 `json:"qualifying_points,omitempty"`
	
	QualifyingSpend *float64 `json:"qualifying_spend,omitempty"`
	
	QualifyingPurchases *int64 `json:"qualifying_purchases,omitempty"`
	
	EffectiveDate *time.Time `json:"effective_date,omitempty" validate:"omitempty,required"`
	
	ValidUntil *time.Time `json:"valid_until,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateCustomerTierHistoryRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.TierId != nil {
		hasUpdate = true
	}
	
	if r.PreviousTierId != nil {
		hasUpdate = true
	}
	
	if r.ChangeType != nil {
		hasUpdate = true
	}
	
	if r.ChangeReason != nil {
		hasUpdate = true
	}
	
	if r.QualifyingPoints != nil {
		hasUpdate = true
	}
	
	if r.QualifyingSpend != nil {
		hasUpdate = true
	}
	
	if r.QualifyingPurchases != nil {
		hasUpdate = true
	}
	
	if r.EffectiveDate != nil {
		hasUpdate = true
	}
	
	if r.ValidUntil != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// CustomerTierHistoryListResponse represents a paginated list of customer_tier_history records
type CustomerTierHistoryListResponse struct {
	Items      []*CustomerTierHistoryResponse `json:"items"`
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
