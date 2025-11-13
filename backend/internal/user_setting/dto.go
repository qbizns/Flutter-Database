package user_setting

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserSettingsResponse represents a user_settings response
type UserSettingsResponse struct {
	
	UserId *uuid.UUID `json:"user_id"`
	
	Theme *string `json:"theme"`
	
	Language *string `json:"language"`
	
	Timezone *string `json:"timezone"`
	
	DefaultDashboard *string `json:"default_dashboard"`
	
	DashboardLayout json.RawMessage `json:"dashboard_layout"`
	
	ItemsPerPage *int64 `json:"items_per_page"`
	
	DefaultView *string `json:"default_view"`
	
	DesktopNotifications *bool `json:"desktop_notifications"`
	
	SoundNotifications *bool `json:"sound_notifications"`
	
	DefaultLocationId *uuid.UUID `json:"default_location_id"`
	
	QuickActions json.RawMessage `json:"quick_actions"`
	
	CustomPreferences json.RawMessage `json:"custom_preferences"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
}

// CreateUserSettingsRequest represents a request to create a user_settings
type CreateUserSettingsRequest struct {
	
	UserId *uuid.UUID `json:"user_id"`
	
	Theme *string `json:"theme"`
	
	Language *string `json:"language"`
	
	Timezone *string `json:"timezone"`
	
	DefaultDashboard *string `json:"default_dashboard"`
	
	// Duplicate removed: DashboardLayout json.RawMessage `json:"dashboard_layout"`
	
	ItemsPerPage *int64 `json:"items_per_page"`
	
	DefaultView *string `json:"default_view"`
	
	DesktopNotifications *bool `json:"desktop_notifications"`
	
	SoundNotifications *bool `json:"sound_notifications"`
	
	DefaultLocationId *uuid.UUID `json:"default_location_id"`
	
	// Duplicate removed: QuickActions json.RawMessage `json:"quick_actions"`
	
	// Duplicate removed: CustomPreferences json.RawMessage `json:"custom_preferences"`
	
}

// Validate validates the create request
func (r *CreateUserSettingsRequest) Validate() error {
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateUserSettingsRequest represents a request to update a user_settings
type UpdateUserSettingsRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	Theme *string `json:"theme,omitempty"`
	
	Language *string `json:"language,omitempty"`
	
	Timezone *string `json:"timezone,omitempty"`
	
	DefaultDashboard *string `json:"default_dashboard,omitempty"`
	
	DashboardLayout *json.RawMessage `json:"dashboard_layout,omitempty"`
	
	ItemsPerPage *int64 `json:"items_per_page,omitempty"`
	
	DefaultView *string `json:"default_view,omitempty"`
	
	DesktopNotifications *bool `json:"desktop_notifications,omitempty"`
	
	SoundNotifications *bool `json:"sound_notifications,omitempty"`
	
	DefaultLocationId *uuid.UUID `json:"default_location_id,omitempty"`
	
	QuickActions *json.RawMessage `json:"quick_actions,omitempty"`
	
	CustomPreferences *json.RawMessage `json:"custom_preferences,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateUserSettingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.Theme != nil {
		hasUpdate = true
	}
	
	if r.Language != nil {
		hasUpdate = true
	}
	
	if r.Timezone != nil {
		hasUpdate = true
	}
	
	if r.DefaultDashboard != nil {
		hasUpdate = true
	}
	
	if r.DashboardLayout != nil {
		hasUpdate = true
	}
	
	if r.ItemsPerPage != nil {
		hasUpdate = true
	}
	
	if r.DefaultView != nil {
		hasUpdate = true
	}
	
	if r.DesktopNotifications != nil {
		hasUpdate = true
	}
	
	if r.SoundNotifications != nil {
		hasUpdate = true
	}
	
	if r.DefaultLocationId != nil {
		hasUpdate = true
	}
	
	if r.QuickActions != nil {
		hasUpdate = true
	}
	
	if r.CustomPreferences != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UserSettingsListResponse represents a paginated list of user_settings records
type UserSettingsListResponse struct {
	Items      []*UserSettingsResponse `json:"items"`
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
