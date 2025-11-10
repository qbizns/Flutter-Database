package promotions

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Promotion represents a promotion/discount campaign
type Promotion struct {
	ID                        uuid.UUID       `json:"id"`
	OrganizationID            uuid.UUID       `json:"organization_id"`
	PromotionCode             string          `json:"promotion_code"`
	Name                      string          `json:"name"`
	Description               string          `json:"description"`
	PromotionType             string          `json:"promotion_type"` // percentage, fixed_amount, buy_x_get_y, bundle, quantity_discount
	DiscountValue             float64         `json:"discount_value"`
	AppliesToType             string          `json:"applies_to"` // all, specific_products, specific_categories, cart_total
	ApplicableProductIDs      []uuid.UUID     `json:"applicable_product_ids"`
	ApplicableCategoryIDs     []uuid.UUID     `json:"applicable_category_ids"`
	MinimumPurchaseAmount     float64         `json:"minimum_purchase_amount"`
	MinimumQuantity           int             `json:"minimum_quantity"`
	BuyQuantity               *int            `json:"buy_quantity"`
	GetQuantity               *int            `json:"get_quantity"`
	GetDiscountPercentage     *float64        `json:"get_discount_percentage"`
	MaxUsesTotal              *int            `json:"max_uses_total"`
	MaxUsesPerCustomer        *int            `json:"max_uses_per_customer"`
	CurrentUses               int             `json:"current_uses"`
	StartDate                 time.Time       `json:"start_date"`
	EndDate                   *time.Time      `json:"end_date"`
	IsActive                  bool            `json:"is_active"`
	IsCombinable              bool            `json:"is_combinable"`
	Priority                  int             `json:"priority"`
	TermsAndConditions        string          `json:"terms_and_conditions"`
	CreatedAt                 time.Time       `json:"created_at"`
	UpdatedAt                 time.Time       `json:"updated_at"`
	CreatedBy                 uuid.UUID       `json:"created_by"`
	UpdatedBy                 *uuid.UUID      `json:"updated_by"`
	DeletedAt                 *time.Time      `json:"deleted_at,omitempty"`
}

// PromotionUsage tracks promotion usage per transaction
type PromotionUsage struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	PromotionID    uuid.UUID  `json:"promotion_id"`
	SaleID         *uuid.UUID `json:"sale_id"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	DiscountAmount float64    `json:"discount_amount"`
	UsedAt         time.Time  `json:"used_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

// PromotionFilters for listing promotions
type PromotionFilters struct {
	Search       string
	PromotionType *string
	IsActive     *bool
	Page         int
	PageSize     int
}

// PromotionUsageFilters for listing promotion usage
type PromotionUsageFilters struct {
	PromotionID *uuid.UUID
	CustomerID  *uuid.UUID
	SaleID      *uuid.UUID
	DateFrom    *time.Time
	DateTo      *time.Time
	Page        int
	PageSize    int
}

// UUIDArray implements driver.Valuer and sql.Scanner for UUID arrays
type UUIDArray []uuid.UUID

func (a UUIDArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Repository defines the promotion data access interface
type Repository interface {
	// Promotions
	List(ctx context.Context, orgID uuid.UUID, filters PromotionFilters) ([]Promotion, error)
	Count(ctx context.Context, orgID uuid.UUID, filters PromotionFilters) (int64, error)
	Create(ctx context.Context, promotion *Promotion) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Promotion, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*Promotion, error)
	Update(ctx context.Context, promotion *Promotion) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Promotion Usage
	ListUsage(ctx context.Context, orgID uuid.UUID, filters PromotionUsageFilters) ([]PromotionUsage, error)
	CountUsage(ctx context.Context, orgID uuid.UUID, filters PromotionUsageFilters) (int64, error)
	RecordUsage(ctx context.Context, usage *PromotionUsage) error
	GetUsageByPromotion(ctx context.Context, orgID uuid.UUID, promotionID uuid.UUID, customerID *uuid.UUID) (int64, error)
}
