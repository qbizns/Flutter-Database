package purchases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GoodsReceipt represents a goods receiving document
type GoodsReceipt struct {
	ID              uuid.UUID                `json:"id"`
	OrganizationID  uuid.UUID                `json:"organization_id"`
	ReceiptNumber   string                   `json:"receipt_number"`
	PurchaseOrderID *uuid.UUID               `json:"purchase_order_id"`
	SupplierID      uuid.UUID                `json:"supplier_id"`
	LocationID      uuid.UUID                `json:"location_id"`
	ReceiptDate     time.Time                `json:"receipt_date"`
	ReceivedBy      uuid.UUID                `json:"received_by"`
	Status          string                   `json:"status"` // draft, received, inspected, completed
	Notes           string                   `json:"notes"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	DeletedAt       *time.Time               `json:"deleted_at,omitempty"`
	CreatedBy       *uuid.UUID               `json:"created_by"`
	Items           []GoodsReceiptItem       `json:"items,omitempty"`
}

// GoodsReceiptItem represents an item received in a goods receipt
type GoodsReceiptItem struct {
	ID                    uuid.UUID  `json:"id"`
	OrganizationID        uuid.UUID  `json:"organization_id"`
	GoodsReceiptID        uuid.UUID  `json:"goods_receipt_id"`
	PurchaseOrderItemID   *uuid.UUID `json:"purchase_order_item_id"`
	ProductID             uuid.UUID  `json:"product_id"`
	ProductVariantID      *uuid.UUID `json:"product_variant_id"`
	QuantityReceived      float64    `json:"quantity_received"`
	QuantityAccepted      *float64   `json:"quantity_accepted"`
	QuantityRejected      *float64   `json:"quantity_rejected"`
	RejectionReason       *string    `json:"rejection_reason"`
	CreatedAt             time.Time  `json:"created_at"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}

// CreateGoodsReceiptRequest represents request to create a goods receipt
type CreateGoodsReceiptRequest struct {
	ReceiptNumber   string                           `json:"receipt_number" binding:"required"`
	PurchaseOrderID *uuid.UUID                       `json:"purchase_order_id"`
	SupplierID      uuid.UUID                        `json:"supplier_id" binding:"required"`
	LocationID      uuid.UUID                        `json:"location_id" binding:"required"`
	ReceiptDate     time.Time                        `json:"receipt_date" binding:"required"`
	ReceivedBy      uuid.UUID                        `json:"received_by" binding:"required"`
	Notes           string                           `json:"notes"`
	Items           []CreateGoodsReceiptItemRequest  `json:"items" binding:"required,min=1"`
}

// CreateGoodsReceiptItemRequest represents an item in goods receipt creation
type CreateGoodsReceiptItemRequest struct {
	PurchaseOrderItemID *uuid.UUID `json:"purchase_order_item_id"`
	ProductID           uuid.UUID  `json:"product_id" binding:"required"`
	ProductVariantID    *uuid.UUID `json:"product_variant_id"`
	QuantityReceived    float64    `json:"quantity_received" binding:"required"`
	QuantityAccepted    *float64   `json:"quantity_accepted"`
	QuantityRejected    *float64   `json:"quantity_rejected"`
	RejectionReason     *string    `json:"rejection_reason"`
}

// UpdateGoodsReceiptRequest represents request to update a goods receipt
type UpdateGoodsReceiptRequest struct {
	Status string                          `json:"status"`
	Notes  string                          `json:"notes"`
	Items  []UpdateGoodsReceiptItemRequest `json:"items"`
}

// UpdateGoodsReceiptItemRequest represents item update in goods receipt
type UpdateGoodsReceiptItemRequest struct {
	ID               uuid.UUID `json:"id" binding:"required"`
	QuantityAccepted *float64  `json:"quantity_accepted"`
	QuantityRejected *float64  `json:"quantity_rejected"`
	RejectionReason  *string   `json:"rejection_reason"`
}

// GoodsReceiptFilters represents filters for listing goods receipts
type GoodsReceiptFilters struct {
	PurchaseOrderID *uuid.UUID
	SupplierID      *uuid.UUID
	LocationID      *uuid.UUID
	Status          *string
	StartDate       *time.Time
	EndDate         *time.Time
	Search          string
	Page            int
	PageSize        int
}

// GoodsReceiptRepository extends the main Repository with goods receipt specific methods
type GoodsReceiptRepository interface {
	// Goods Receipts
	ListReceipts(ctx context.Context, orgID uuid.UUID, filters GoodsReceiptFilters) ([]GoodsReceipt, error)
	CountReceipts(ctx context.Context, orgID uuid.UUID, filters GoodsReceiptFilters) (int64, error)
	CreateReceipt(ctx context.Context, receipt *GoodsReceipt) error
	GetReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GoodsReceipt, error)
	GetReceiptByNumber(ctx context.Context, orgID uuid.UUID, receiptNumber string) (*GoodsReceipt, error)
	UpdateReceipt(ctx context.Context, receipt *GoodsReceipt) error
	DeleteReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Goods Receipt Items
	ListReceiptItems(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID) ([]GoodsReceiptItem, error)
	GetReceiptItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GoodsReceiptItem, error)
	CreateReceiptItem(ctx context.Context, item *GoodsReceiptItem) error
	UpdateReceiptItem(ctx context.Context, item *GoodsReceiptItem) error
	DeleteReceiptItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Status and Calculations
	UpdateReceiptStatus(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID, status string) error
	InspectReceiptItems(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID, items []GoodsReceiptItem) error
}
