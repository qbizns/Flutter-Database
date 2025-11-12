package sales

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Sale struct {
	ID                 uuid.UUID   `json:"id"`
	OrganizationID     uuid.UUID   `json:"organization_id"`
	SaleNumber         string      `json:"sale_number"`
	ReferenceNumber    string      `json:"reference_number"`
	TransactionType    string      `json:"transaction_type"`
	CustomerID         *uuid.UUID  `json:"customer_id"`
	CashierID          *uuid.UUID  `json:"cashier_id"`
	Subtotal           float64     `json:"subtotal"`
	TaxAmount          float64     `json:"tax_amount"`
	DiscountAmount     float64     `json:"discount_amount"`
	TotalAmount        float64     `json:"total_amount"`
	PaidAmount         float64     `json:"paid_amount"`
	ChangeAmount       float64     `json:"change_amount"`
	OutstandingAmount  float64     `json:"outstanding_amount"`
	PaymentStatus      string      `json:"payment_status"`
	DiscountType       string      `json:"discount_type"`
	DiscountValue      float64     `json:"discount_value"`
	DiscountReason     string      `json:"discount_reason"`
	TransactionDate    time.Time   `json:"transaction_date"`
	CompletedAt        *time.Time  `json:"completed_at"`
	Notes              string      `json:"notes"`
	InternalNotes      string      `json:"internal_notes"`
	CustomFields       interface{} `json:"custom_fields"`
	Metadata           interface{} `json:"metadata"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	CreatedBy          uuid.UUID   `json:"created_by"`
	UpdatedBy          *uuid.UUID  `json:"updated_by"`
	DeletedAt          *time.Time  `json:"deleted_at,omitempty"`
	Items              []SaleItem  `json:"items,omitempty"`
}

type SaleItem struct {
	ID             uuid.UUID   `json:"id"`
	SaleID         uuid.UUID   `json:"sale_id"`
	OrganizationID uuid.UUID   `json:"organization_id"`
	ProductID      *uuid.UUID  `json:"product_id"`
	ProductName    string      `json:"product_name"`
	ProductSKU     string      `json:"product_sku"`
	Quantity       float64     `json:"quantity"`
	Unit           string      `json:"unit"`
	UnitPrice      float64     `json:"unit_price"`
	CostPrice      float64     `json:"cost_price"`
	Subtotal       float64     `json:"subtotal"`
	TaxRate        float64     `json:"tax_rate"`
	TaxAmount      float64     `json:"tax_amount"`
	DiscountAmount float64     `json:"discount_amount"`
	Total          float64     `json:"total"`
	DiscountType   string      `json:"discount_type"`
	DiscountValue  float64     `json:"discount_value"`
	Notes          string      `json:"notes"`
	CustomFields   interface{} `json:"custom_fields"`
	Metadata       interface{} `json:"metadata"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type SaleFilters struct {
	Search          string
	CustomerID      *uuid.UUID
	CashierID       *uuid.UUID
	PaymentStatus   *string
	TransactionType *string
	StartDate       *time.Time
	EndDate         *time.Time
	Page            int
	PageSize        int
}

type SaleItemFilters struct {
	ProductID *uuid.UUID
	Page      int
	PageSize  int
}

type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error)
	Count(ctx context.Context, orgID uuid.UUID, filters SaleFilters) (int64, error)
	Create(ctx context.Context, sale *Sale) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Sale, error)
	GetBySaleNumber(ctx context.Context, orgID uuid.UUID, saleNumber string) (*Sale, error)
	Update(ctx context.Context, sale *Sale) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	CreateItem(ctx context.Context, item *SaleItem) error
	GetItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) (*SaleItem, error)
	ListItems(ctx context.Context, saleID uuid.UUID, filters SaleItemFilters) ([]SaleItem, error)
	CountItems(ctx context.Context, saleID uuid.UUID) (int64, error)
	UpdateItem(ctx context.Context, item *SaleItem) error
	DeleteItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) error
	GetWithItems(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*Sale, error)
}
