package products

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProductVariant represents a product variant with SKU, attributes, and pricing
type ProductVariant struct {
	ID              uuid.UUID            `json:"id"`
	OrganizationID  uuid.UUID            `json:"organization_id"`
	ProductID       uuid.UUID            `json:"product_id"`
	VariantName     string               `json:"variant_name"`
	SKU             string               `json:"sku"`
	Barcode         string               `json:"barcode"`
	Attributes      VariantAttributes    `json:"attributes"`
	CostPrice       *float64             `json:"cost_price"`
	SellingPrice    *float64             `json:"selling_price"`
	CompareAtPrice  *float64             `json:"compare_at_price"`
	CurrentStock    float64              `json:"current_stock"`
	ReorderLevel    float64              `json:"reorder_level"`
	ReorderQuantity float64              `json:"reorder_quantity"`
	Weight          *float64             `json:"weight"`
	WeightUnit      string               `json:"weight_unit"`
	Dimensions      *VariantDimensions   `json:"dimensions"`
	IsActive        bool                 `json:"is_active"`
	IsDefault       bool                 `json:"is_default"`
	SortOrder       int                  `json:"sort_order"`
	ImageURL        string               `json:"image_url"`
	Notes           string               `json:"notes"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	DeletedAt       *time.Time           `json:"deleted_at,omitempty"`
	CreatedBy       uuid.UUID            `json:"created_by"`
	UpdatedBy       *uuid.UUID           `json:"updated_by"`
}

// VariantAttributes contains flexible attributes like size, color, etc.
type VariantAttributes map[string]interface{}

// VariantDimensions contains physical dimensions
type VariantDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"`
}

// Value implements driver.Valuer for JSONB storage
func (va VariantAttributes) Value() (driver.Value, error) {
	return json.Marshal(va)
}

// Scan implements sql.Scanner for JSONB retrieval
func (va *VariantAttributes) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(data, va)
}

// Value implements driver.Valuer for JSONB storage
func (vd *VariantDimensions) Value() (driver.Value, error) {
	if vd == nil {
		return nil, nil
	}
	return json.Marshal(vd)
}

// Scan implements sql.Scanner for JSONB retrieval
func (vd *VariantDimensions) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(data, vd)
}

// VariantFilters for querying variants
type VariantFilters struct {
	ProductID *uuid.UUID
	IsActive  *bool
	Search    string
	Page      int
	PageSize  int
}

// VariantRepository defines variant data access interface
type VariantRepository interface {
	// CRUD operations
	Create(ctx context.Context, variant *ProductVariant) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductVariant, error)
	GetBySKU(ctx context.Context, orgID uuid.UUID, sku string) (*ProductVariant, error)
	Update(ctx context.Context, variant *ProductVariant) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// List/Filter operations
	List(ctx context.Context, orgID uuid.UUID, filters VariantFilters) ([]ProductVariant, error)
	Count(ctx context.Context, orgID uuid.UUID, filters VariantFilters) (int64, error)

	// Batch operations
	ListByProductID(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]ProductVariant, error)
	DeleteByProductID(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) error

	// Stock management
	UpdateStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error
	AdjustStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error

	// Default variant
	SetDefault(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID) error
}

// VariantService handles variant business logic
type VariantService struct {
	repo VariantRepository
}

// NewVariantService creates a new variant service
func NewVariantService(repo VariantRepository) *VariantService {
	return &VariantService{repo: repo}
}

// Create creates a new variant
func (s *VariantService) Create(ctx context.Context, variant *ProductVariant) error {
	if err := s.validateVariant(variant); err != nil {
		return err
	}

	variant.ID = uuid.New()
	variant.CreatedAt = time.Now()
	variant.UpdatedAt = time.Now()

	return s.repo.Create(ctx, variant)
}

// Get retrieves a variant by ID
func (s *VariantService) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductVariant, error) {
	return s.repo.Get(ctx, orgID, id)
}

// GetBySKU retrieves a variant by SKU
func (s *VariantService) GetBySKU(ctx context.Context, orgID uuid.UUID, sku string) (*ProductVariant, error) {
	return s.repo.GetBySKU(ctx, orgID, sku)
}

// Update updates an existing variant
func (s *VariantService) Update(ctx context.Context, variant *ProductVariant) error {
	if err := s.validateVariant(variant); err != nil {
		return err
	}

	variant.UpdatedAt = time.Now()
	return s.repo.Update(ctx, variant)
}

// Delete deletes a variant
func (s *VariantService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.Delete(ctx, orgID, id)
}

// List retrieves variants with filters
func (s *VariantService) List(ctx context.Context, orgID uuid.UUID, filters VariantFilters) ([]ProductVariant, error) {
	return s.repo.List(ctx, orgID, filters)
}

// ListByProductID retrieves all variants for a product
func (s *VariantService) ListByProductID(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]ProductVariant, error) {
	return s.repo.ListByProductID(ctx, orgID, productID)
}

// UpdateStock sets the stock quantity
func (s *VariantService) UpdateStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error {
	if quantity < 0 {
		return NewValidationError("stock quantity cannot be negative")
	}
	return s.repo.UpdateStock(ctx, orgID, variantID, quantity)
}

// AdjustStock adjusts stock by a delta
func (s *VariantService) AdjustStock(ctx context.Context, orgID uuid.UUID, variantID uuid.UUID, quantity float64) error {
	return s.repo.AdjustStock(ctx, orgID, variantID, quantity)
}

// validateVariant validates variant data
func (s *VariantService) validateVariant(variant *ProductVariant) error {
	if variant.ProductID == uuid.Nil {
		return NewValidationError("product_id is required")
	}
	if variant.VariantName == "" {
		return NewValidationError("variant_name is required")
	}
	if variant.SKU == "" {
		return NewValidationError("sku is required")
	}
	if variant.CostPrice != nil && *variant.CostPrice < 0 {
		return NewValidationError("cost_price cannot be negative")
	}
	if variant.SellingPrice != nil && *variant.SellingPrice < 0 {
		return NewValidationError("selling_price cannot be negative")
	}
	if variant.CompareAtPrice != nil && *variant.CompareAtPrice < 0 {
		return NewValidationError("compare_at_price cannot be negative")
	}
	if variant.Weight != nil && *variant.Weight < 0 {
		return NewValidationError("weight cannot be negative")
	}
	if variant.ReorderLevel < 0 {
		return NewValidationError("reorder_level cannot be negative")
	}
	return nil
}
