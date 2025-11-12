package integration_config

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// IntegrationConfigsResponse represents a integration_configs response
type IntegrationConfigsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	IntegrationType string `json:"integration_type"`
	
	ProviderName string `json:"provider_name"`
	
	Credentials json.RawMessage `json:"credentials"`
	
	Settings json.RawMessage `json:"settings"`
	
	IsActive *bool `json:"is_active"`
	
	IsConnected *bool `json:"is_connected"`
	
	ConnectionStatus *string `json:"connection_status"`
	
	LastSyncAt *time.Time `json:"last_sync_at"`
	
	LastSyncStatus *string `json:"last_sync_status"`
	
	SyncFrequency *string `json:"sync_frequency"`
	
	WebhookUrl *string `json:"webhook_url"`
	
	WebhookSecret *string `json:"webhook_secret"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateIntegrationConfigsRequest represents a request to create a integration_configs
type CreateIntegrationConfigsRequest struct {
	
	IntegrationType string `json:"integration_type" validate:"required"`
	
	ProviderName string `json:"provider_name" validate:"required"`
	
	Credentials json.RawMessage `json:"credentials" validate:"required"`
	
	Settings json.RawMessage `json:"settings"`
	
	IsActive *bool `json:"is_active"`
	
	IsConnected *bool `json:"is_connected"`
	
	ConnectionStatus *string `json:"connection_status"`
	
	LastSyncAt *time.Time `json:"last_sync_at"`
	
	LastSyncStatus *string `json:"last_sync_status"`
	
	SyncFrequency *string `json:"sync_frequency"`
	
	WebhookUrl *string `json:"webhook_url" validate:"url"`
	
	WebhookSecret *string `json:"webhook_secret"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateIntegrationConfigsRequest) Validate() error {
	
	if r.IntegrationType == "" {
		return fmt.Errorf("integration_type is required")
	}
	
	if r.ProviderName == "" {
		return fmt.Errorf("provider_name is required")
	}
	
	if r.Credentials == nil {
		return fmt.Errorf("credentials is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateIntegrationConfigsRequest represents a request to update a integration_configs
type UpdateIntegrationConfigsRequest struct {
	
	IntegrationType *string `json:"integration_type,omitempty" validate:"omitempty,required"`
	
	ProviderName *string `json:"provider_name,omitempty" validate:"omitempty,required"`
	
	Credentials *json.RawMessage `json:"credentials,omitempty" validate:"omitempty,required"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	IsConnected *bool `json:"is_connected,omitempty"`
	
	ConnectionStatus *string `json:"connection_status,omitempty"`
	
	LastSyncAt *time.Time `json:"last_sync_at,omitempty"`
	
	LastSyncStatus *string `json:"last_sync_status,omitempty"`
	
	SyncFrequency *string `json:"sync_frequency,omitempty"`
	
	WebhookUrl *string `json:"webhook_url,omitempty" validate:"omitempty,url"`
	
	WebhookSecret *string `json:"webhook_secret,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateIntegrationConfigsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.IntegrationType != nil {
		hasUpdate = true
	}
	
	if r.ProviderName != nil {
		hasUpdate = true
	}
	
	if r.Credentials != nil {
		hasUpdate = true
	}
	
	if r.Settings != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.IsConnected != nil {
		hasUpdate = true
	}
	
	if r.ConnectionStatus != nil {
		hasUpdate = true
	}
	
	if r.LastSyncAt != nil {
		hasUpdate = true
	}
	
	if r.LastSyncStatus != nil {
		hasUpdate = true
	}
	
	if r.SyncFrequency != nil {
		hasUpdate = true
	}
	
	if r.WebhookUrl != nil {
		hasUpdate = true
	}
	
	if r.WebhookSecret != nil {
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

// IntegrationConfigsListResponse represents a paginated list of integration_configs records
type IntegrationConfigsListResponse struct {
	Items      []*IntegrationConfigsResponse `json:"items"`
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
