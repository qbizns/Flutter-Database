package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentTermLinesResponse represents a payment_term_lines response
type PaymentTermLinesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	PaymentTermId uuid.UUID `json:"payment_term_id"`
	
	Sequence int64 `json:"sequence"`
	
	ValueType string `json:"value_type"`
	
	ValueAmount *float64 `json:"value_amount"`
	
	DaysAfter *int64 `json:"days_after"`
	
	EndOfMonth *bool `json:"end_of_month"`
	
	DayOfMonth *int64 `json:"day_of_month"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	(valueType string `json:"(value_type"`
	
}

// CreatePaymentTermLinesRequest represents a request to create a payment_term_lines
type CreatePaymentTermLinesRequest struct {
	
	PaymentTermId uuid.UUID `json:"payment_term_id" validate:"required"`
	
	Sequence int64 `json:"sequence" validate:"required"`
	
	ValueType string `json:"value_type" validate:"required"`
	
	ValueAmount *float64 `json:"value_amount"`
	
	DaysAfter *int64 `json:"days_after"`
	
	EndOfMonth *bool `json:"end_of_month"`
	
	DayOfMonth *int64 `json:"day_of_month"`
	
	(valueType string `json:"(value_type" validate:"required"`
	
}

// Validate validates the create request
func (r *CreatePaymentTermLinesRequest) Validate() error {
	
	if r.PaymentTermId == uuid.Nil {
		return fmt.Errorf("payment_term_id is required")
	}
	
	if r.Sequence == 0 {
		return fmt.Errorf("sequence is required")
	}
	
	if r.ValueType == "" {
		return fmt.Errorf("value_type is required")
	}
	
	if r.(valueType == "" {
		return fmt.Errorf("(value_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePaymentTermLinesRequest represents a request to update a payment_term_lines
type UpdatePaymentTermLinesRequest struct {
	
	PaymentTermId *uuid.UUID `json:"payment_term_id,omitempty" validate:"omitempty,required"`
	
	Sequence *int64 `json:"sequence,omitempty" validate:"omitempty,required"`
	
	ValueType *string `json:"value_type,omitempty" validate:"omitempty,required"`
	
	ValueAmount *float64 `json:"value_amount,omitempty"`
	
	DaysAfter *int64 `json:"days_after,omitempty"`
	
	EndOfMonth *bool `json:"end_of_month,omitempty"`
	
	DayOfMonth *int64 `json:"day_of_month,omitempty"`
	
	(valueType *string `json:"(value_type,omitempty" validate:"omitempty,required"`
	
}

// Validate validates the update request
func (r *UpdatePaymentTermLinesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.PaymentTermId != nil {
		hasUpdate = true
	}
	
	if r.Sequence != nil {
		hasUpdate = true
	}
	
	if r.ValueType != nil {
		hasUpdate = true
	}
	
	if r.ValueAmount != nil {
		hasUpdate = true
	}
	
	if r.DaysAfter != nil {
		hasUpdate = true
	}
	
	if r.EndOfMonth != nil {
		hasUpdate = true
	}
	
	if r.DayOfMonth != nil {
		hasUpdate = true
	}
	
	if r.(valueType != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PaymentTermLinesListResponse represents a paginated list of payment_term_lines records
type PaymentTermLinesListResponse struct {
	Items      []*PaymentTermLinesResponse `json:"items"`
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
