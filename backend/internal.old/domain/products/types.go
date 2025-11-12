package products

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	Name                string     `json:"name"`
	SKU                 string     `json:"sku"`
	Barcode             string     `json:"barcode"`
	Description         string     `json:"description"`
	CategoryID          *uuid.UUID `json:"category_id"`
	UnitPrice           float64    `json:"unit_price"`
	Cost                float64    `json:"cost"`
	TaxRate             float64    `json:"tax_rate"`
	Unit                string     `json:"unit"`
	CurrentStock        float64    `json:"current_stock"`
	MinStockLevel       float64    `json:"min_stock_level"`
	MaxStockLevel       float64    `json:"max_stock_level"`
	IsActive            bool       `json:"is_active"`
	IsTrackInventory    bool       `json:"is_track_inventory"`
	AllowNegativeStock  bool       `json:"allow_negative_stock"`
	ImageURL            string     `json:"image_url"`
	ProductType         string     `json:"product_type"`
	AccountingAccountID *uuid.UUID `json:"accounting_account_id"`
	TaxCodeID           *uuid.UUID `json:"tax_code_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CreatedBy           uuid.UUID  `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// ProductFilters represents filters for listing products
type ProductFilters struct {
	Search     string
	CategoryID *uuid.UUID
	IsActive   *bool
	Page       int
	PageSize   int
}

// Repository defines the product data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters ProductFilters) ([]Product, error)
	Count(ctx context.Context, orgID uuid.UUID, filters ProductFilters) (int64, error)
	Create(ctx context.Context, product *Product) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Product, error)
	GetBySKU(ctx context.Context, orgID uuid.UUID, sku string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStock(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, quantity float64) error
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}
