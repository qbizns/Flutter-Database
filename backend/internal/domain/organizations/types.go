package organizations

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// OrganizationStatus represents the status of an organization
type OrganizationStatus string

const (
	StatusTrial      OrganizationStatus = "trial"
	StatusActive     OrganizationStatus = "active"
	StatusSuspended  OrganizationStatus = "suspended"
	StatusCancelled  OrganizationStatus = "cancelled"
)

// Organization represents a tenant organization
type Organization struct {
	ID                   uuid.UUID              `json:"id"`
	Name                 string                 `json:"name"`
	Slug                 string                 `json:"slug"`
	Description          string                 `json:"description"`
	Email                string                 `json:"email"`
	Phone                string                 `json:"phone"`
	Address              string                 `json:"address"`
	City                 string                 `json:"city"`
	State                string                 `json:"state"`
	Country              string                 `json:"country"`
	PostalCode           string                 `json:"postal_code"`
	Status               OrganizationStatus     `json:"status"`
	Plan                 string                 `json:"plan"`
	TrialEndsAt          *time.Time             `json:"trial_ends_at"`
	SubscriptionStartsAt *time.Time             `json:"subscription_starts_at"`
	SubscriptionEndsAt   *time.Time             `json:"subscription_ends_at"`
	MaxUsers             int                    `json:"max_users"`
	MaxProducts          int                    `json:"max_products"`
	MaxLocations         int                    `json:"max_locations"`
	Settings             json.RawMessage        `json:"settings"`
	Metadata             json.RawMessage        `json:"metadata"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DeletedAt            *time.Time             `json:"deleted_at,omitempty"`
	CreatedBy            *uuid.UUID             `json:"created_by"`
	UpdatedBy            *uuid.UUID             `json:"updated_by"`
}

// OrganizationFilters represents filters for listing organizations
type OrganizationFilters struct {
	Search   string
	Status   *OrganizationStatus
	Page     int
	PageSize int
}

// Repository defines the organization data access interface
type Repository interface {
	List(ctx context.Context, filters OrganizationFilters) ([]Organization, error)
	Count(ctx context.Context, filters OrganizationFilters) (int64, error)
	Create(ctx context.Context, org *Organization) error
	Get(ctx context.Context, id uuid.UUID) (*Organization, error)
	GetBySlug(ctx context.Context, slug string) (*Organization, error)
	Update(ctx context.Context, org *Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Scan implements the Scanner interface for OrganizationStatus
func (s *OrganizationStatus) Scan(value interface{}) error {
	*s = OrganizationStatus(value.(string))
	return nil
}

// Value implements the driver Valuer interface for OrganizationStatus
func (s OrganizationStatus) Value() (driver.Value, error) {
	return string(s), nil
}
