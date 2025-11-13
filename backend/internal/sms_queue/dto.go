package sms_queue

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SmsQueueResponse represents a sms_queue response
type SmsQueueResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId *uuid.UUID `json:"organization_id"`
	
	ToPhone string `json:"to_phone"`
	
	FromPhone *string `json:"from_phone"`
	
	Message string `json:"message"`
	
	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	Provider *string `json:"provider"`
	
	ProviderMessageId *string `json:"provider_message_id"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	ErrorMessage *string `json:"error_message"`
	
	CostAmount *float64 `json:"cost_amount"`
	
	CostCurrency *string `json:"cost_currency"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	SentAt *time.Time `json:"sent_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreateSmsQueueRequest represents a request to create a sms_queue
type CreateSmsQueueRequest struct {
	
	ToPhone string `json:"to_phone" validate:"required,e164"`
	
	FromPhone *string `json:"from_phone" validate:"e164"`
	
	Message string `json:"message" validate:"required"`
	
	// 	Status *string `json:"status"`
	
	// 	Status *string `json:"status"`
	
	Provider *string `json:"provider"`
	
	ProviderMessageId *string `json:"provider_message_id"`
	
	Attempts *int64 `json:"attempts"`
	
	MaxAttempts *int64 `json:"max_attempts"`
	
	ErrorMessage *string `json:"error_message"`
	
	CostAmount *float64 `json:"cost_amount"`
	
	CostCurrency *string `json:"cost_currency"`
	
	ScheduledAt *time.Time `json:"scheduled_at"`
	
	SentAt *time.Time `json:"sent_at"`
	
	FailedAt *time.Time `json:"failed_at"`
	
}

// Validate validates the create request
func (r *CreateSmsQueueRequest) Validate() error {
	
	if r.ToPhone == "" {
		return fmt.Errorf("to_phone is required")
	}
	
	if r.Message == "" {
		return fmt.Errorf("message is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSmsQueueRequest represents a request to update a sms_queue
type UpdateSmsQueueRequest struct {
	
	ToPhone *string `json:"to_phone,omitempty" validate:"omitempty,required,e164"`
	
	FromPhone *string `json:"from_phone,omitempty" validate:"omitempty,e164"`
	
	Message *string `json:"message,omitempty" validate:"omitempty,required"`
	
	// 	Status *string `json:"status,omitempty"`
	
	// 	Status *string `json:"status,omitempty"`
	
	Provider *string `json:"provider,omitempty"`
	
	ProviderMessageId *string `json:"provider_message_id,omitempty"`
	
	Attempts *int64 `json:"attempts,omitempty"`
	
	MaxAttempts *int64 `json:"max_attempts,omitempty"`
	
	ErrorMessage *string `json:"error_message,omitempty"`
	
	CostAmount *float64 `json:"cost_amount,omitempty"`
	
	CostCurrency *string `json:"cost_currency,omitempty"`
	
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	
	SentAt *time.Time `json:"sent_at,omitempty"`
	
	FailedAt *time.Time `json:"failed_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSmsQueueRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ToPhone != nil {
		hasUpdate = true
	}
	
	if r.FromPhone != nil {
		hasUpdate = true
	}
	
	if r.Message != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.Provider != nil {
		hasUpdate = true
	}
	
	if r.ProviderMessageId != nil {
		hasUpdate = true
	}
	
	if r.Attempts != nil {
		hasUpdate = true
	}
	
	if r.MaxAttempts != nil {
		hasUpdate = true
	}
	
	if r.ErrorMessage != nil {
		hasUpdate = true
	}
	
	if r.CostAmount != nil {
		hasUpdate = true
	}
	
	if r.CostCurrency != nil {
		hasUpdate = true
	}
	
	if r.ScheduledAt != nil {
		hasUpdate = true
	}
	
	if r.SentAt != nil {
		hasUpdate = true
	}
	
	if r.FailedAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// SmsQueueListResponse represents a paginated list of sms_queue records
type SmsQueueListResponse struct {
	Items      []*SmsQueueResponse `json:"items"`
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
