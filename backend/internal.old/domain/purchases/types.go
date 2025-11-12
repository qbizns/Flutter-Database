package purchases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PurchaseOrder represents a purchase order entity
type PurchaseOrder struct {
	ID                   uuid.UUID             `json:"id"`
	OrganizationID       uuid.UUID             `json:"organization_id"`
	SupplierID           uuid.UUID             `json:"supplier_id"`
	PONumber             string                `json:"po_number"`
	PODate               time.Time             `json:"po_date"`
	ExpectedDeliveryDate *time.Time            `json:"expected_delivery_date"`
	ActualDeliveryDate   *time.Time            `json:"actual_delivery_date"`
	Status               string                `json:"status"`
	Subtotal             float64               `json:"subtotal"`
	TaxAmount            float64               `json:"tax_amount"`
	DiscountAmount       float64               `json:"discount_amount"`
	ShippingCost         float64               `json:"shipping_cost"`
	TotalAmount          float64               `json:"total_amount"`
	PaymentStatus        string                `json:"payment_status"`
	PaidAmount           float64               `json:"paid_amount"`
	Notes                string                `json:"notes"`
	TermsAndConditions   string                `json:"terms_and_conditions"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
	CreatedBy            uuid.UUID             `json:"created_by"`
	UpdatedBy            *uuid.UUID            `json:"updated_by"`
	ApprovedBy           *uuid.UUID            `json:"approved_by"`
	ApprovedAt           *time.Time            `json:"approved_at"`
	DeletedAt            *time.Time            `json:"deleted_at,omitempty"`
	Items                []PurchaseOrderItem   `json:"items,omitempty"`
}

// PurchaseOrderItem represents a line item in a purchase order
type PurchaseOrderItem struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	PurchaseOrderID  uuid.UUID  `json:"purchase_order_id"`
	ProductID        *uuid.UUID `json:"product_id"`
	ProductName      string     `json:"product_name"`
	ProductSKU       string     `json:"product_sku"`
	QuantityOrdered  float64    `json:"quantity_ordered"`
	QuantityReceived float64    `json:"quantity_received"`
	UnitOfMeasure    string     `json:"unit_of_measure"`
	UnitCost         float64    `json:"unit_cost"`
	DiscountPercent  float64    `json:"discount_percentage"`
	DiscountAmount   float64    `json:"discount_amount"`
	TaxPercent       float64    `json:"tax_percentage"`
	TaxAmount        float64    `json:"tax_amount"`
	LineTotal        float64    `json:"line_total"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CreatePurchaseOrderRequest represents request to create a purchase order
type CreatePurchaseOrderRequest struct {
	SupplierID           uuid.UUID                        `json:"supplier_id" binding:"required"`
	PONumber             string                           `json:"po_number" binding:"required"`
	PODate               time.Time                        `json:"po_date" binding:"required"`
	ExpectedDeliveryDate *time.Time                       `json:"expected_delivery_date"`
	Status               string                           `json:"status"`
	Notes                string                           `json:"notes"`
	TermsAndConditions   string                           `json:"terms_and_conditions"`
	Items                []CreatePurchaseOrderItemRequest `json:"items" binding:"required,min=1"`
}

// CreatePurchaseOrderItemRequest represents a line item in create request
type CreatePurchaseOrderItemRequest struct {
	ProductID       *uuid.UUID `json:"product_id"`
	ProductName     string     `json:"product_name" binding:"required"`
	ProductSKU      string     `json:"product_sku"`
	QuantityOrdered float64    `json:"quantity_ordered" binding:"required"`
	UnitOfMeasure   string     `json:"unit_of_measure"`
	UnitCost        float64    `json:"unit_cost" binding:"required"`
	DiscountPercent float64    `json:"discount_percentage"`
	TaxPercent      float64    `json:"tax_percentage"`
	Notes           string     `json:"notes"`
}

// UpdatePurchaseOrderRequest represents request to update a purchase order
type UpdatePurchaseOrderRequest struct {
	PONumber             string                           `json:"po_number"`
	ExpectedDeliveryDate *time.Time                       `json:"expected_delivery_date"`
	Status               string                           `json:"status"`
	Notes                string                           `json:"notes"`
	TermsAndConditions   string                           `json:"terms_and_conditions"`
	Items                []UpdatePurchaseOrderItemRequest `json:"items"`
}

// UpdatePurchaseOrderItemRequest represents item update request
type UpdatePurchaseOrderItemRequest struct {
	ID              uuid.UUID  `json:"id"`
	ProductName     string     `json:"product_name"`
	ProductSKU      string     `json:"product_sku"`
	QuantityOrdered float64    `json:"quantity_ordered"`
	QuantityReceived float64   `json:"quantity_received"`
	UnitOfMeasure   string     `json:"unit_of_measure"`
	UnitCost        float64    `json:"unit_cost"`
	DiscountPercent float64    `json:"discount_percentage"`
	TaxPercent      float64    `json:"tax_percentage"`
	Notes           string     `json:"notes"`
}

// PurchaseOrderFilters represents filters for listing purchase orders
type PurchaseOrderFilters struct {
	SupplierID *uuid.UUID
	Status     *string
	Search     string
	StartDate  *time.Time
	EndDate    *time.Time
	Page       int
	PageSize   int
}

// Repository defines the purchase order data access interface
type Repository interface {
	// Purchase Orders
	ListOrders(ctx context.Context, orgID uuid.UUID, filters PurchaseOrderFilters) ([]PurchaseOrder, error)
	CountOrders(ctx context.Context, orgID uuid.UUID, filters PurchaseOrderFilters) (int64, error)
	CreateOrder(ctx context.Context, po *PurchaseOrder) error
	GetOrder(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PurchaseOrder, error)
	GetOrderByNumber(ctx context.Context, orgID uuid.UUID, poNumber string) (*PurchaseOrder, error)
	UpdateOrder(ctx context.Context, po *PurchaseOrder) error
	DeleteOrder(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Purchase Order Items
	ListItems(ctx context.Context, orgID uuid.UUID, poID uuid.UUID) ([]PurchaseOrderItem, error)
	GetItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PurchaseOrderItem, error)
	CreateItem(ctx context.Context, item *PurchaseOrderItem) error
	UpdateItem(ctx context.Context, item *PurchaseOrderItem) error
	DeleteItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	DeleteItemsByOrderID(ctx context.Context, orgID uuid.UUID, poID uuid.UUID) error

	// Status and Calculations
	UpdateOrderStatus(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, status string) error
	UpdatePaymentStatus(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, paymentStatus string, paidAmount float64) error
	ApproveOrder(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, approvedBy uuid.UUID) error
	RecordReceipt(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, itemID uuid.UUID, quantityReceived float64) error
	GetReceivedQuantity(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID) (float64, error)
}
