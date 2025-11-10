package products

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PriceList represents a price list for tiered pricing (retail, wholesale, VIP, etc.)
type PriceList struct {
	ID                        uuid.UUID  `json:"id"`
	OrganizationID            uuid.UUID  `json:"organization_id"`
	Code                      string     `json:"code"`
	Name                      string     `json:"name"`
	Type                      string     `json:"type"` // standard, retail, wholesale, vip, seasonal, location, customer_group
	EffectiveFrom             *time.Time `json:"effective_from"`
	EffectiveTo               *time.Time `json:"effective_to"`
	BaseAdjustmentType        string     `json:"base_adjustment_type"` // percentage, fixed, none
	BaseAdjustmentValue       float64    `json:"base_adjustment_value"`
	Priority                  int        `json:"priority"` // lower = higher priority
	IsActive                  bool       `json:"is_active"`
	Description               string     `json:"description"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at,omitempty"`
	CreatedBy                 uuid.UUID  `json:"created_by"`
	UpdatedBy                 *uuid.UUID `json:"updated_by"`
}

// PriceListItem represents a price override within a price list
type PriceListItem struct {
	ID                 uuid.UUID  `json:"id"`
	PriceListID        uuid.UUID  `json:"price_list_id"`
	ProductID          *uuid.UUID `json:"product_id"` // Can be NULL for category pricing
	ProductVariantID   *uuid.UUID `json:"product_variant_id"`
	CategoryID         *uuid.UUID `json:"category_id"`
	OverridePrice      *float64   `json:"override_price"`
	DiscountPercentage *float64   `json:"discount_percentage"`
	MarkupPercentage   *float64   `json:"markup_percentage"`
	MinPrice           *float64   `json:"min_price"`
	MaxPrice           *float64   `json:"max_price"`
	MinQuantity        *float64   `json:"min_quantity"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// ProductComponent represents a component in a bundle/kit product
type ProductComponent struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	ParentProductID  uuid.UUID  `json:"parent_product_id"`
	ComponentProductID *uuid.UUID `json:"component_product_id"`
	ComponentVariantID *uuid.UUID `json:"component_variant_id"`
	Quantity         float64    `json:"quantity"`
	InheritPrice     bool       `json:"inherit_price"`
	PriceOverride    *float64   `json:"price_override"`
	DisplayOrder     int        `json:"display_order"`
	IsOptional       bool       `json:"is_optional"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// PricingData contains pricing details (used in computed pricing endpoints)
type PricingData struct {
	ProductID     uuid.UUID  `json:"product_id"`
	VariantID     *uuid.UUID `json:"variant_id"`
	BasePrice     float64    `json:"base_price"`
	CostPrice     float64    `json:"cost_price"`
	SellingPrice  float64    `json:"selling_price"`
	ApplicableLists []PriceListApplication `json:"applicable_lists"`
	FinalPrice    float64    `json:"final_price"`
	Margin        float64    `json:"margin"`
	MarginPercent float64    `json:"margin_percent"`
}

// PriceListApplication shows which price list applies
type PriceListApplication struct {
	PriceListID uuid.UUID `json:"price_list_id"`
	Name        string    `json:"name"`
	AdjustmentType string `json:"adjustment_type"`
	AdjustmentValue float64 `json:"adjustment_value"`
	FinalPrice   float64   `json:"final_price"`
	Priority     int       `json:"priority"`
}

// PriceListFilters for querying price lists
type PriceListFilters struct {
	Type      string
	IsActive  *bool
	Search    string
	Page      int
	PageSize  int
}

// PriceListItemFilters for querying price list items
type PriceListItemFilters struct {
	PriceListID  *uuid.UUID
	ProductID    *uuid.UUID
	VariantID    *uuid.UUID
	CategoryID   *uuid.UUID
	Page         int
	PageSize     int
}

// ComponentFilters for querying product components
type ComponentFilters struct {
	ParentProductID *uuid.UUID
	Page            int
	PageSize        int
}

// PriceListRepository defines price list data access interface
type PriceListRepository interface {
	// PriceList CRUD
	CreatePriceList(ctx context.Context, pl *PriceList) error
	GetPriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PriceList, error)
	UpdatePriceList(ctx context.Context, pl *PriceList) error
	DeletePriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	ListPriceLists(ctx context.Context, orgID uuid.UUID, filters PriceListFilters) ([]PriceList, error)
	CountPriceLists(ctx context.Context, orgID uuid.UUID, filters PriceListFilters) (int64, error)

	// PriceListItem CRUD
	CreatePriceListItem(ctx context.Context, item *PriceListItem) error
	GetPriceListItem(ctx context.Context, id uuid.UUID) (*PriceListItem, error)
	UpdatePriceListItem(ctx context.Context, item *PriceListItem) error
	DeletePriceListItem(ctx context.Context, id uuid.UUID) error
	ListPriceListItems(ctx context.Context, filters PriceListItemFilters) ([]PriceListItem, error)
	CountPriceListItems(ctx context.Context, filters PriceListItemFilters) (int64, error)

	// Pricing calculations
	GetApplicablePriceLists(ctx context.Context, orgID uuid.UUID, now time.Time) ([]PriceList, error)
	CalculatePrice(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, basePrice float64) (float64, error)

	// ProductComponent CRUD
	CreateComponent(ctx context.Context, component *ProductComponent) error
	GetComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductComponent, error)
	UpdateComponent(ctx context.Context, component *ProductComponent) error
	DeleteComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	ListComponents(ctx context.Context, orgID uuid.UUID, filters ComponentFilters) ([]ProductComponent, error)
	ListComponentsByParent(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]ProductComponent, error)
	DeleteComponentsByParent(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) error
}

// PriceListService handles price list business logic
type PriceListService struct {
	repo PriceListRepository
}

// NewPriceListService creates a new price list service
func NewPriceListService(repo PriceListRepository) *PriceListService {
	return &PriceListService{repo: repo}
}

// CreatePriceList creates a new price list
func (s *PriceListService) CreatePriceList(ctx context.Context, pl *PriceList) error {
	if err := s.validatePriceList(pl); err != nil {
		return err
	}

	pl.ID = uuid.New()
	pl.CreatedAt = time.Now()
	pl.UpdatedAt = time.Now()

	return s.repo.CreatePriceList(ctx, pl)
}

// GetPriceList retrieves a price list by ID
func (s *PriceListService) GetPriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PriceList, error) {
	return s.repo.GetPriceList(ctx, orgID, id)
}

// UpdatePriceList updates an existing price list
func (s *PriceListService) UpdatePriceList(ctx context.Context, pl *PriceList) error {
	if err := s.validatePriceList(pl); err != nil {
		return err
	}

	pl.UpdatedAt = time.Now()
	return s.repo.UpdatePriceList(ctx, pl)
}

// DeletePriceList deletes a price list
func (s *PriceListService) DeletePriceList(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeletePriceList(ctx, orgID, id)
}

// ListPriceLists retrieves price lists
func (s *PriceListService) ListPriceLists(ctx context.Context, orgID uuid.UUID, filters PriceListFilters) ([]PriceList, error) {
	return s.repo.ListPriceLists(ctx, orgID, filters)
}

// CreatePriceListItem creates a price override
func (s *PriceListService) CreatePriceListItem(ctx context.Context, item *PriceListItem) error {
	if err := s.validatePriceListItem(item); err != nil {
		return err
	}

	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	return s.repo.CreatePriceListItem(ctx, item)
}

// GetPriceListItem retrieves a price list item
func (s *PriceListService) GetPriceListItem(ctx context.Context, id uuid.UUID) (*PriceListItem, error) {
	return s.repo.GetPriceListItem(ctx, id)
}

// UpdatePriceListItem updates a price list item
func (s *PriceListService) UpdatePriceListItem(ctx context.Context, item *PriceListItem) error {
	if err := s.validatePriceListItem(item); err != nil {
		return err
	}

	item.UpdatedAt = time.Now()
	return s.repo.UpdatePriceListItem(ctx, item)
}

// DeletePriceListItem deletes a price list item
func (s *PriceListService) DeletePriceListItem(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePriceListItem(ctx, id)
}

// CalculatePrice calculates the final price for a product
func (s *PriceListService) CalculatePrice(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, basePrice float64) (float64, error) {
	if basePrice < 0 {
		return 0, NewValidationError("base price cannot be negative")
	}

	return s.repo.CalculatePrice(ctx, orgID, productID, variantID, basePrice)
}

// CreateComponent creates a product component
func (s *PriceListService) CreateComponent(ctx context.Context, component *ProductComponent) error {
	if err := s.validateComponent(component); err != nil {
		return err
	}

	component.ID = uuid.New()
	component.CreatedAt = time.Now()
	component.UpdatedAt = time.Now()

	return s.repo.CreateComponent(ctx, component)
}

// GetComponent retrieves a component
func (s *PriceListService) GetComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductComponent, error) {
	return s.repo.GetComponent(ctx, orgID, id)
}

// UpdateComponent updates a component
func (s *PriceListService) UpdateComponent(ctx context.Context, component *ProductComponent) error {
	if err := s.validateComponent(component); err != nil {
		return err
	}

	component.UpdatedAt = time.Now()
	return s.repo.UpdateComponent(ctx, component)
}

// DeleteComponent deletes a component
func (s *PriceListService) DeleteComponent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteComponent(ctx, orgID, id)
}

// ListComponentsByParent retrieves all components of a bundle
func (s *PriceListService) ListComponentsByParent(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]ProductComponent, error) {
	return s.repo.ListComponentsByParent(ctx, orgID, parentID)
}

// validatePriceList validates price list data
func (s *PriceListService) validatePriceList(pl *PriceList) error {
	if pl.Code == "" {
		return NewValidationError("code is required")
	}
	if pl.Name == "" {
		return NewValidationError("name is required")
	}
	if pl.Type == "" {
		return NewValidationError("type is required")
	}

	validTypes := map[string]bool{
		"standard":       true,
		"retail":         true,
		"wholesale":      true,
		"vip":            true,
		"seasonal":       true,
		"location":       true,
		"customer_group": true,
	}
	if !validTypes[pl.Type] {
		return NewValidationError("invalid price list type")
	}

	if pl.EffectiveFrom != nil && pl.EffectiveTo != nil && pl.EffectiveFrom.After(*pl.EffectiveTo) {
		return NewValidationError("effective_from must be before effective_to")
	}

	if pl.BaseAdjustmentType == "percentage" {
		if pl.BaseAdjustmentValue < -100 || pl.BaseAdjustmentValue > 100 {
			return NewValidationError("percentage adjustment must be between -100 and 100")
		}
	}

	return nil
}

// validatePriceListItem validates price list item data
func (s *PriceListService) validatePriceListItem(item *PriceListItem) error {
	if item.PriceListID == uuid.Nil {
		return NewValidationError("price_list_id is required")
	}

	// At least one target must be specified
	if item.ProductID == nil && item.ProductVariantID == nil && item.CategoryID == nil {
		return NewValidationError("at least one of product_id, product_variant_id, or category_id is required")
	}

	// At least one price modifier must be specified
	if item.OverridePrice == nil && item.DiscountPercentage == nil && item.MarkupPercentage == nil {
		return NewValidationError("at least one price modifier must be specified")
	}

	// Validate percentage values
	if item.DiscountPercentage != nil && (*item.DiscountPercentage < 0 || *item.DiscountPercentage > 100) {
		return NewValidationError("discount_percentage must be between 0 and 100")
	}
	if item.MarkupPercentage != nil && (*item.MarkupPercentage < 0 || *item.MarkupPercentage > 100) {
		return NewValidationError("markup_percentage must be between 0 and 100")
	}

	// Validate prices
	if item.OverridePrice != nil && *item.OverridePrice < 0 {
		return NewValidationError("override_price cannot be negative")
	}
	if item.MinPrice != nil && *item.MinPrice < 0 {
		return NewValidationError("min_price cannot be negative")
	}
	if item.MaxPrice != nil && *item.MaxPrice < 0 {
		return NewValidationError("max_price cannot be negative")
	}
	if item.MinQuantity != nil && *item.MinQuantity <= 0 {
		return NewValidationError("min_quantity must be positive")
	}

	return nil
}

// validateComponent validates component data
func (s *PriceListService) validateComponent(component *ProductComponent) error {
	if component.ParentProductID == uuid.Nil {
		return NewValidationError("parent_product_id is required")
	}

	// At least one component target must be specified
	if component.ComponentProductID == nil && component.ComponentVariantID == nil {
		return NewValidationError("either component_product_id or component_variant_id is required")
	}

	if component.Quantity <= 0 {
		return NewValidationError("quantity must be positive")
	}

	if component.PriceOverride != nil && *component.PriceOverride < 0 {
		return NewValidationError("price_override cannot be negative")
	}

	return nil
}
