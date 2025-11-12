package categories

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Category represents a product category entity with hierarchical support
type Category struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Description    string     `json:"description"`
	ParentID       *uuid.UUID `json:"parent_id"`
	Level          int        `json:"level"`
	Path           string     `json:"path"`
	ImageURL       string     `json:"image_url"`
	Icon           string     `json:"icon"`
	Color          string     `json:"color"`
	SortOrder      int        `json:"sort_order"`
	IsActive       bool       `json:"is_active"`
	Metadata       []byte     `json:"metadata"` // JSONB stored as bytes
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	UpdatedBy      *uuid.UUID `json:"updated_by"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// CategoryFilters represents filters for listing categories
type CategoryFilters struct {
	Search   string
	ParentID *uuid.UUID // Filter by parent category (nil = root categories)
	IsActive *bool
	Page     int
	PageSize int
}

// Repository defines the category data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters CategoryFilters) ([]Category, error)
	Count(ctx context.Context, orgID uuid.UUID, filters CategoryFilters) (int64, error)
	Create(ctx context.Context, category *Category) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Category, error)
	GetBySlug(ctx context.Context, orgID uuid.UUID, slug string) (*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetChildren(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]Category, error)
}
