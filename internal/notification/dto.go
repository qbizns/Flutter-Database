package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NotificationsResponse represents a notifications response
type NotificationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	UserId uuid.UUID `json:"user_id"`
	
	NotificationType string `json:"notification_type"`
	
	Category string `json:"category"`
	
	Title string `json:"title"`
	
	Message string `json:"message"`
	
	ActionUrl *string `json:"action_url"`
	
	ActionLabel *string `json:"action_label"`
	
	Channels *string `json:"channels"`
	
	IsRead *bool `json:"is_read"`
	
	ReadAt *time.Time `json:"read_at"`
	
	RelatedEntityType *string `json:"related_entity_type"`
	
	RelatedEntityId *uuid.UUID `json:"related_entity_id"`
	
	Priority *string `json:"priority"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
	CreatedAt *time.Time `json:"created_at"`
	
}

// CreateNotificationsRequest represents a request to create a notifications
type CreateNotificationsRequest struct {
	
	UserId uuid.UUID `json:"user_id" validate:"required"`
	
	NotificationType string `json:"notification_type" validate:"required"`
	
	Category string `json:"category" validate:"required"`
	
	Title string `json:"title" validate:"required"`
	
	Message string `json:"message" validate:"required"`
	
	ActionUrl *string `json:"action_url" validate:"url"`
	
	ActionLabel *string `json:"action_label"`
	
	Channels *string `json:"channels"`
	
	IsRead *bool `json:"is_read"`
	
	ReadAt *time.Time `json:"read_at"`
	
	RelatedEntityType *string `json:"related_entity_type"`
	
	RelatedEntityId *uuid.UUID `json:"related_entity_id"`
	
	Priority *string `json:"priority"`
	
	ExpiresAt *time.Time `json:"expires_at"`
	
}

// Validate validates the create request
func (r *CreateNotificationsRequest) Validate() error {
	
	if r.UserId == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	
	if r.NotificationType == "" {
		return fmt.Errorf("notification_type is required")
	}
	
	if r.Category == "" {
		return fmt.Errorf("category is required")
	}
	
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	
	if r.Message == "" {
		return fmt.Errorf("message is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateNotificationsRequest represents a request to update a notifications
type UpdateNotificationsRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,required"`
	
	NotificationType *string `json:"notification_type,omitempty" validate:"omitempty,required"`
	
	Category *string `json:"category,omitempty" validate:"omitempty,required"`
	
	Title *string `json:"title,omitempty" validate:"omitempty,required"`
	
	Message *string `json:"message,omitempty" validate:"omitempty,required"`
	
	ActionUrl *string `json:"action_url,omitempty" validate:"omitempty,url"`
	
	ActionLabel *string `json:"action_label,omitempty"`
	
	Channels *string `json:"channels,omitempty"`
	
	IsRead *bool `json:"is_read,omitempty"`
	
	ReadAt *time.Time `json:"read_at,omitempty"`
	
	RelatedEntityType *string `json:"related_entity_type,omitempty"`
	
	RelatedEntityId *uuid.UUID `json:"related_entity_id,omitempty"`
	
	Priority *string `json:"priority,omitempty"`
	
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateNotificationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.NotificationType != nil {
		hasUpdate = true
	}
	
	if r.Category != nil {
		hasUpdate = true
	}
	
	if r.Title != nil {
		hasUpdate = true
	}
	
	if r.Message != nil {
		hasUpdate = true
	}
	
	if r.ActionUrl != nil {
		hasUpdate = true
	}
	
	if r.ActionLabel != nil {
		hasUpdate = true
	}
	
	if r.Channels != nil {
		hasUpdate = true
	}
	
	if r.IsRead != nil {
		hasUpdate = true
	}
	
	if r.ReadAt != nil {
		hasUpdate = true
	}
	
	if r.RelatedEntityType != nil {
		hasUpdate = true
	}
	
	if r.RelatedEntityId != nil {
		hasUpdate = true
	}
	
	if r.Priority != nil {
		hasUpdate = true
	}
	
	if r.ExpiresAt != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// NotificationsListResponse represents a paginated list of notifications records
type NotificationsListResponse struct {
	Items      []*NotificationsResponse `json:"items"`
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
