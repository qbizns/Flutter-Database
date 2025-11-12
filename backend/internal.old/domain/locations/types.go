package locations

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Location represents a store location/branch for multi-location businesses
type Location struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	LocationCode   string     `json:"location_code"`
	Name           string     `json:"name"`
	LocationType   string     `json:"location_type"` // store, warehouse, headquarters, kiosk, online, other
	Phone          *string    `json:"phone,omitempty"`
	Email          *string    `json:"email,omitempty"`
	ManagerUserID  *uuid.UUID `json:"manager_user_id,omitempty"`
	AddressLine1   *string    `json:"address_line1,omitempty"`
	AddressLine2   *string    `json:"address_line2,omitempty"`
	City           *string    `json:"city,omitempty"`
	State          *string    `json:"state,omitempty"`
	Country        *string    `json:"country,omitempty"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	Timezone       string     `json:"timezone"`
	BusinessHours  any        `json:"business_hours,omitempty"` // JSONB
	IsActive       bool       `json:"is_active"`
	IsPrimary      bool       `json:"is_primary"`
	AllowSales     bool       `json:"allow_sales"`
	AllowPurchases bool       `json:"allow_purchases"`
	TaxRate        float64    `json:"tax_rate"`
	Notes          *string    `json:"notes,omitempty"`
	Settings       any        `json:"settings,omitempty"` // JSONB
	Metadata       any        `json:"metadata,omitempty"` // JSONB
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// LocationFilters represents filters for listing locations
type LocationFilters struct {
	Search       string
	LocationType *string
	IsActive     *bool
	IsPrimary    *bool
	City         *string
	State        *string
	Country      *string
	Page         int
	PageSize     int
}

// Repository defines the data access interface for locations
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters LocationFilters) ([]Location, error)
	Count(ctx context.Context, orgID uuid.UUID, filters LocationFilters) (int64, error)
	Create(ctx context.Context, location *Location) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Location, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*Location, error)
	Update(ctx context.Context, location *Location) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}
