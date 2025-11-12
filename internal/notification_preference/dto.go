package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NotificationPreferencesResponse represents a notification_preferences response
type NotificationPreferencesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	Category string `json:"category"`
	
	InAppEnabled *bool `json:"in_app_enabled"`
	
	EmailEnabled *bool `json:"email_enabled"`
	
	SmsEnabled *bool `json:"sms_enabled"`
	
	PushEnabled *bool `json:"push_enabled"`
	
	Frequency *string `json:"frequency"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateNotificationPreferencesRequest represents a request to create a notification_preferences
type CreateNotificationPreferencesRequest struct {
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	Category string `json:"category" validate:"required"`
	
	InAppEnabled *bool `json:"in_app_enabled"`
	
	EmailEnabled *bool `json:"email_enabled"`
	
	SmsEnabled *bool `json:"sms_enabled"`
	
	PushEnabled *bool `json:"push_enabled"`
	
	Frequency *string `json:"frequency"`
	
}

// Validate validates the create request
func (r *CreateNotificationPreferencesRequest) Validate() error {
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.Category == "" {
		return fmt.Errorf("category is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateNotificationPreferencesRequest represents a request to update a notification_preferences
type UpdateNotificationPreferencesRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	Category *string `json:"category,omitempty" validate:"omitempty,required"`
	
	InAppEnabled *bool `json:"in_app_enabled,omitempty"`
	
	EmailEnabled *bool `json:"email_enabled,omitempty"`
	
	SmsEnabled *bool `json:"sms_enabled,omitempty"`
	
	PushEnabled *bool `json:"push_enabled,omitempty"`
	
	Frequency *string `json:"frequency,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateNotificationPreferencesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.Category != nil {
		hasUpdate = true
	}
	
	if r.InAppEnabled != nil {
		hasUpdate = true
	}
	
	if r.EmailEnabled != nil {
		hasUpdate = true
	}
	
	if r.SmsEnabled != nil {
		hasUpdate = true
	}
	
	if r.PushEnabled != nil {
		hasUpdate = true
	}
	
	if r.Frequency != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// NotificationPreferencesListResponse represents a paginated list of notification_preferences records
type NotificationPreferencesListResponse struct {
	Items      []*NotificationPreferencesResponse `json:"items"`
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
