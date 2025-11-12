package organizations

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FeatureKey represents available feature keys
type FeatureKey string

const (
	FeatureAccounting           FeatureKey = "accounting"
	FeatureEInvoicing           FeatureKey = "e_invoicing"
	FeatureAdvancedInventory    FeatureKey = "advanced_inventory"
	FeatureMultiLocation        FeatureKey = "multi_location"
	FeatureMultiCurrency        FeatureKey = "multi_currency"
	FeatureLoyaltyProgram       FeatureKey = "loyalty_program"
	FeatureDeliveryManagement   FeatureKey = "delivery_management"
	FeatureRestaurantMode       FeatureKey = "restaurant_mode"
	FeatureTableManagement      FeatureKey = "table_management"
	FeatureKitchenDisplay       FeatureKey = "kitchen_display"
	FeatureOnlineOrdering       FeatureKey = "online_ordering"
	FeatureAPIAccess            FeatureKey = "api_access"
)

// OrganizationFeature represents a feature flag for an organization
type OrganizationFeature struct {
	ID              uuid.UUID       `json:"id"`
	OrganizationID  uuid.UUID       `json:"organization_id"`
	FeatureKey      FeatureKey      `json:"feature_key"`
	IsEnabled       bool            `json:"is_enabled"`
	IsAvailable     bool            `json:"is_available"`
	Configuration   json.RawMessage `json:"configuration"`
	Limits          json.RawMessage `json:"limits"`
	EnabledAt       *time.Time      `json:"enabled_at"`
	DisabledAt      *time.Time      `json:"disabled_at"`
	ExpiresAt       *time.Time      `json:"expires_at"`
	Notes           string          `json:"notes"`
	Metadata        json.RawMessage `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
	CreatedBy       *uuid.UUID      `json:"created_by"`
	UpdatedBy       *uuid.UUID      `json:"updated_by"`
}

// CreateFeatureRequest represents request to create/enable a feature
type CreateFeatureRequest struct {
	FeatureKey    FeatureKey      `json:"feature_key" binding:"required"`
	IsEnabled     bool            `json:"is_enabled"`
	Configuration json.RawMessage `json:"configuration"`
	Limits        json.RawMessage `json:"limits"`
	ExpiresAt     *time.Time      `json:"expires_at"`
	Notes         string          `json:"notes"`
}

// UpdateFeatureRequest represents request to update a feature
type UpdateFeatureRequest struct {
	IsEnabled     *bool           `json:"is_enabled"`
	Configuration json.RawMessage `json:"configuration"`
	Limits        json.RawMessage `json:"limits"`
	ExpiresAt     *time.Time      `json:"expires_at"`
	Notes         string          `json:"notes"`
}

// FeatureFilters represents filters for listing organization features
type FeatureFilters struct {
	IsEnabled *bool
	Search    string
	Page      int
	PageSize  int
}

// OrganizationFeatureRepository extends the Organization Repository with feature management
type OrganizationFeatureRepository interface {
	// Features
	ListFeatures(ctx context.Context, orgID uuid.UUID, filters FeatureFilters) ([]OrganizationFeature, error)
	CountFeatures(ctx context.Context, orgID uuid.UUID, filters FeatureFilters) (int64, error)
	GetFeature(ctx context.Context, orgID uuid.UUID, featureKey FeatureKey) (*OrganizationFeature, error)
	CreateFeature(ctx context.Context, feature *OrganizationFeature) error
	UpdateFeature(ctx context.Context, feature *OrganizationFeature) error
	DeleteFeature(ctx context.Context, orgID uuid.UUID, featureKey FeatureKey) error

	// Feature Status
	IsFeatureEnabled(ctx context.Context, orgID uuid.UUID, featureKey FeatureKey) (bool, error)
	SetFeatureEnabled(ctx context.Context, orgID uuid.UUID, featureKey FeatureKey, isEnabled bool) error

	// Bulk operations
	ListEnabledFeatures(ctx context.Context, orgID uuid.UUID) ([]OrganizationFeature, error)
	GetFeaturesByKeys(ctx context.Context, orgID uuid.UUID, featureKeys []FeatureKey) ([]OrganizationFeature, error)
}

// Scan implements the Scanner interface for FeatureKey
func (f *FeatureKey) Scan(value interface{}) error {
	*f = FeatureKey(value.(string))
	return nil
}

// Value implements the driver Valuer interface for FeatureKey
func (f FeatureKey) Value() (driver.Value, error) {
	return string(f), nil
}
