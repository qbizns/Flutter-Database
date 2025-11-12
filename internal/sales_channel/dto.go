package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SalesChannelsResponse represents a sales_channels response
type SalesChannelsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	ChannelCode string `json:"channel_code"`
	
	ChannelName string `json:"channel_name"`
	
	ChannelType string `json:"channel_type"`
	
	ChannelType *string `json:"channel_type"`
	
	IsActive *bool `json:"is_active"`
	
	SyncInventory *bool `json:"sync_inventory"`
	
	SyncCustomers *bool `json:"sync_customers"`
	
	ExternalSystemName *string `json:"external_system_name"`
	
	ApiEndpoint *string `json:"api_endpoint"`
	
	Settings json.RawMessage `json:"settings"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateSalesChannelsRequest represents a request to create a sales_channels
type CreateSalesChannelsRequest struct {
	
	ChannelCode string `json:"channel_code" validate:"required"`
	
	ChannelName string `json:"channel_name" validate:"required"`
	
	ChannelType string `json:"channel_type" validate:"required"`
	
	ChannelType *string `json:"channel_type"`
	
	IsActive *bool `json:"is_active"`
	
	SyncInventory *bool `json:"sync_inventory"`
	
	SyncCustomers *bool `json:"sync_customers"`
	
	ExternalSystemName *string `json:"external_system_name"`
	
	ApiEndpoint *string `json:"api_endpoint"`
	
	Settings json.RawMessage `json:"settings"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreateSalesChannelsRequest) Validate() error {
	
	if r.ChannelCode == "" {
		return fmt.Errorf("channel_code is required")
	}
	
	if r.ChannelName == "" {
		return fmt.Errorf("channel_name is required")
	}
	
	if r.ChannelType == "" {
		return fmt.Errorf("channel_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateSalesChannelsRequest represents a request to update a sales_channels
type UpdateSalesChannelsRequest struct {
	
	ChannelCode *string `json:"channel_code,omitempty" validate:"omitempty,required"`
	
	ChannelName *string `json:"channel_name,omitempty" validate:"omitempty,required"`
	
	ChannelType *string `json:"channel_type,omitempty" validate:"omitempty,required"`
	
	ChannelType *string `json:"channel_type,omitempty"`
	
	IsActive *bool `json:"is_active,omitempty"`
	
	SyncInventory *bool `json:"sync_inventory,omitempty"`
	
	SyncCustomers *bool `json:"sync_customers,omitempty"`
	
	ExternalSystemName *string `json:"external_system_name,omitempty"`
	
	ApiEndpoint *string `json:"api_endpoint,omitempty"`
	
	Settings *json.RawMessage `json:"settings,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateSalesChannelsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.ChannelCode != nil {
		hasUpdate = true
	}
	
	if r.ChannelName != nil {
		hasUpdate = true
	}
	
	if r.ChannelType != nil {
		hasUpdate = true
	}
	
	if r.ChannelType != nil {
		hasUpdate = true
	}
	
	if r.IsActive != nil {
		hasUpdate = true
	}
	
	if r.SyncInventory != nil {
		hasUpdate = true
	}
	
	if r.SyncCustomers != nil {
		hasUpdate = true
	}
	
	if r.ExternalSystemName != nil {
		hasUpdate = true
	}
	
	if r.ApiEndpoint != nil {
		hasUpdate = true
	}
	
	if r.Settings != nil {
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

// SalesChannelsListResponse represents a paginated list of sales_channels records
type SalesChannelsListResponse struct {
	Items      []*SalesChannelsResponse `json:"items"`
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
