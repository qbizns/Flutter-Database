package integrations

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SalesChannel represents a sales channel (in-store, web, marketplace, etc.)
type SalesChannel struct {
	ID                  uuid.UUID   `json:"id"`
	OrganizationID      uuid.UUID   `json:"organization_id"`
	ChannelCode         string      `json:"channel_code"`
	ChannelName         string      `json:"channel_name"`
	ChannelType         string      `json:"channel_type"` // in_store, web, mobile_app, marketplace, phone, social, partner
	IsActive            bool        `json:"is_active"`
	SyncInventory       bool        `json:"sync_inventory"`
	SyncCustomers       bool        `json:"sync_customers"`
	ExternalSystemName  string      `json:"external_system_name"`
	APIEndpoint         string      `json:"api_endpoint"`
	Settings            interface{} `json:"settings"`
	CreatedBy           uuid.UUID   `json:"created_by"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	DeletedAt           *time.Time  `json:"deleted_at,omitempty"`
}

// ExternalOrderMapping maps internal sales to external order systems
type ExternalOrderMapping struct {
	ID                  uuid.UUID   `json:"id"`
	OrganizationID      uuid.UUID   `json:"organization_id"`
	SaleID              uuid.UUID   `json:"sale_id"`
	SalesChannelID      uuid.UUID   `json:"sales_channel_id"`
	ExternalOrderID     string      `json:"external_order_id"`
	ExternalOrderNumber string      `json:"external_order_number"`
	SyncStatus          string      `json:"sync_status"` // pending, synced, failed, conflict
	LastSyncAt          *time.Time  `json:"last_sync_at"`
	ExternalData        interface{} `json:"external_data"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	DeletedAt           *time.Time  `json:"deleted_at,omitempty"`
}

// SalesChannelFilters for list operations
type SalesChannelFilters struct {
	Search      string
	ChannelType *string
	IsActive    *bool
	Page        int
	PageSize    int
}

// ExternalOrderMappingFilters for list operations
type ExternalOrderMappingFilters struct {
	SalesChannelID *uuid.UUID
	SyncStatus     *string
	Page           int
	PageSize       int
}

// SalesChannelRepository defines data access interface
type SalesChannelRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters SalesChannelFilters) ([]SalesChannel, error)
	Count(ctx context.Context, orgID uuid.UUID, filters SalesChannelFilters) (int64, error)
	Create(ctx context.Context, channel *SalesChannel) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SalesChannel, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*SalesChannel, error)
	Update(ctx context.Context, channel *SalesChannel) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

// ExternalOrderMappingRepository defines data access interface
type ExternalOrderMappingRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters ExternalOrderMappingFilters) ([]ExternalOrderMapping, error)
	Count(ctx context.Context, orgID uuid.UUID, filters ExternalOrderMappingFilters) (int64, error)
	Create(ctx context.Context, mapping *ExternalOrderMapping) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ExternalOrderMapping, error)
	GetBySaleID(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*ExternalOrderMapping, error)
	GetByExternalOrderID(ctx context.Context, orgID uuid.UUID, channelID uuid.UUID, externalOrderID string) (*ExternalOrderMapping, error)
	Update(ctx context.Context, mapping *ExternalOrderMapping) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateSyncStatus(ctx context.Context, id uuid.UUID, status string, lastSyncAt time.Time) error
}
