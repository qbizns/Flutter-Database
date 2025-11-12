package external_order_mapping

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ExternalOrderMappingsResponse represents a external_order_mappings response
type ExternalOrderMappingsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SaleId uuid.UUID `json:"sale_id"`
	
	SalesChannelId uuid.UUID `json:"sales_channel_id"`
	
	ExternalOrderId string `json:"external_order_id"`
	
	ExternalOrderNumber *string `json:"external_order_number"`
	
	SyncStatus *string `json:"sync_status"`
	
	SyncStatus *string `json:"sync_status"`
	
	LastSyncAt *time.Time `json:"last_sync_at"`
	
	ExternalData json.RawMessage `json:"external_data"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
}

// CreateExternalOrderMappingsRequest represents a request to create a external_order_mappings
type CreateExternalOrderMappingsRequest struct {
	
	SaleId uuid.UUID `json:"sale_id" validate:"required"`
	
	SalesChannelId uuid.UUID `json:"sales_channel_id" validate:"required"`
	
	ExternalOrderId string `json:"external_order_id" validate:"required"`
	
	ExternalOrderNumber *string `json:"external_order_number"`
	
	SyncStatus *string `json:"sync_status"`
	
	SyncStatus *string `json:"sync_status"`
	
	LastSyncAt *time.Time `json:"last_sync_at"`
	
	ExternalData json.RawMessage `json:"external_data"`
	
}

// Validate validates the create request
func (r *CreateExternalOrderMappingsRequest) Validate() error {
	
	if r.SaleId == uuid.Nil {
		return fmt.Errorf("sale_id is required")
	}
	
	if r.SalesChannelId == uuid.Nil {
		return fmt.Errorf("sales_channel_id is required")
	}
	
	if r.ExternalOrderId == "" {
		return fmt.Errorf("external_order_id is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateExternalOrderMappingsRequest represents a request to update a external_order_mappings
type UpdateExternalOrderMappingsRequest struct {
	
	SaleId *uuid.UUID `json:"sale_id,omitempty" validate:"omitempty,required"`
	
	SalesChannelId *uuid.UUID `json:"sales_channel_id,omitempty" validate:"omitempty,required"`
	
	ExternalOrderId *string `json:"external_order_id,omitempty" validate:"omitempty,required"`
	
	ExternalOrderNumber *string `json:"external_order_number,omitempty"`
	
	SyncStatus *string `json:"sync_status,omitempty"`
	
	SyncStatus *string `json:"sync_status,omitempty"`
	
	LastSyncAt *time.Time `json:"last_sync_at,omitempty"`
	
	ExternalData *json.RawMessage `json:"external_data,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateExternalOrderMappingsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.SalesChannelId != nil {
		hasUpdate = true
	}
	
	if r.ExternalOrderId != nil {
		hasUpdate = true
	}
	
	if r.ExternalOrderNumber != nil {
		hasUpdate = true
	}
	
	if r.SyncStatus != nil {
		hasUpdate = true
	}
	
	if r.SyncStatus != nil {
		hasUpdate = true
	}
	
	if r.LastSyncAt != nil {
		hasUpdate = true
	}
	
	if r.ExternalData != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ExternalOrderMappingsListResponse represents a paginated list of external_order_mappings records
type ExternalOrderMappingsListResponse struct {
	Items      []*ExternalOrderMappingsResponse `json:"items"`
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
