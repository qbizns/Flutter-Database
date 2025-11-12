package inventory

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TransferStatus represents the status of an inventory transfer
type TransferStatus string

const (
	TransferStatusDraft             TransferStatus = "draft"
	TransferStatusRequested         TransferStatus = "requested"
	TransferStatusApproved          TransferStatus = "approved"
	TransferStatusInTransit         TransferStatus = "in_transit"
	TransferStatusPartiallyReceived TransferStatus = "partially_received"
	TransferStatusReceived          TransferStatus = "received"
	TransferStatusCancelled         TransferStatus = "cancelled"
	TransferStatusRejected          TransferStatus = "rejected"
)

// Scan implements the sql.Scanner interface
func (s *TransferStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*s = TransferStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*s = TransferStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (s TransferStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// TransferItemStatus represents the status of a transfer item
type TransferItemStatus string

const (
	TransferItemStatusPending          TransferItemStatus = "pending"
	TransferItemStatusApproved         TransferItemStatus = "approved"
	TransferItemStatusShipped          TransferItemStatus = "shipped"
	TransferItemStatusPartiallyReceived TransferItemStatus = "partially_received"
	TransferItemStatusReceived         TransferItemStatus = "received"
	TransferItemStatusCancelled        TransferItemStatus = "cancelled"
)

// Scan implements the sql.Scanner interface
func (s *TransferItemStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*s = TransferItemStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*s = TransferItemStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (s TransferItemStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// InventoryTransfer represents an inter-location inventory transfer
type InventoryTransfer struct {
	ID                   uuid.UUID                `json:"id"`
	OrganizationID       uuid.UUID                `json:"organization_id"`
	TransferNumber       string                   `json:"transfer_number"`
	TransferDate         time.Time                `json:"transfer_date"`
	FromLocationID       uuid.UUID                `json:"from_location_id"`
	ToLocationID         uuid.UUID                `json:"to_location_id"`
	Status               TransferStatus           `json:"status"`
	RequestedDate        *time.Time               `json:"requested_date,omitempty"`
	ApprovedDate         *time.Time               `json:"approved_date,omitempty"`
	ShippedDate          *time.Time               `json:"shipped_date,omitempty"`
	ExpectedDeliveryDate *time.Time               `json:"expected_delivery_date,omitempty"`
	ReceivedDate         *time.Time               `json:"received_date,omitempty"`
	Carrier              *string                  `json:"carrier,omitempty"`
	TrackingNumber       *string                  `json:"tracking_number,omitempty"`
	ShippingCost         float64                  `json:"shipping_cost"`
	Reason               *string                  `json:"reason,omitempty"`
	Notes                *string                  `json:"notes,omitempty"`
	RejectionReason      *string                  `json:"rejection_reason,omitempty"`
	Metadata             json.RawMessage          `json:"metadata,omitempty"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
	DeletedAt            *time.Time               `json:"deleted_at,omitempty"`
	CreatedBy            *uuid.UUID               `json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID               `json:"updated_by,omitempty"`
	RequestedBy          *uuid.UUID               `json:"requested_by,omitempty"`
	ApprovedBy           *uuid.UUID               `json:"approved_by,omitempty"`
	ShippedBy            *uuid.UUID               `json:"shipped_by,omitempty"`
	ReceivedBy           *uuid.UUID               `json:"received_by,omitempty"`
}

// InventoryTransferItem represents a line item in an inventory transfer
type InventoryTransferItem struct {
	ID                   uuid.UUID             `json:"id"`
	OrganizationID       uuid.UUID             `json:"organization_id"`
	InventoryTransferID  uuid.UUID             `json:"inventory_transfer_id"`
	ProductID            *uuid.UUID            `json:"product_id,omitempty"`
	ProductVariantID     *uuid.UUID            `json:"product_variant_id,omitempty"`
	ProductName          string                `json:"product_name"`
	ProductSKU           *string               `json:"product_sku,omitempty"`
	QuantityRequested    float64               `json:"quantity_requested"`
	QuantityShipped      float64               `json:"quantity_shipped"`
	QuantityReceived     float64               `json:"quantity_received"`
	UnitOfMeasure        string                `json:"unit_of_measure"`
	UnitCost             *float64              `json:"unit_cost,omitempty"`
	TotalCost            *float64              `json:"total_cost,omitempty"`
	ItemStatus           TransferItemStatus    `json:"item_status"`
	VarianceQuantity     float64               `json:"variance_quantity"`
	VarianceReason       *string               `json:"variance_reason,omitempty"`
	Notes                *string               `json:"notes,omitempty"`
	Metadata             json.RawMessage       `json:"metadata,omitempty"`
	CreatedAt            time.Time             `json:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at"`
}

// TransferFilters represents filters for listing transfers
type TransferFilters struct {
	Status              *TransferStatus
	FromLocationID      *uuid.UUID
	ToLocationID        *uuid.UUID
	TransferNumberLike  *string
	StartDate           *time.Time
	EndDate             *time.Time
	CreatedBy           *uuid.UUID
	HasTracking         *bool
	Page                int
	PageSize            int
}

// TransferItemFilters represents filters for listing transfer items
type TransferItemFilters struct {
	TransferID    uuid.UUID
	ItemStatus    *TransferItemStatus
	ProductID     *uuid.UUID
	HasVariance   *bool
	Page          int
	PageSize      int
}

// TransferRepository defines the transfer data access interface
type TransferRepository interface {
	// Transfers
	ListTransfers(ctx context.Context, orgID uuid.UUID, filters TransferFilters) ([]InventoryTransfer, error)
	CountTransfers(ctx context.Context, orgID uuid.UUID, filters TransferFilters) (int64, error)
	GetTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransfer, error)
	GetTransferByNumber(ctx context.Context, orgID uuid.UUID, transferNumber string) (*InventoryTransfer, error)
	CreateTransfer(ctx context.Context, transfer *InventoryTransfer) error
	UpdateTransfer(ctx context.Context, transfer *InventoryTransfer) error
	DeleteTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateTransferStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status TransferStatus) error
	UpdateTransferApproval(ctx context.Context, orgID uuid.UUID, id uuid.UUID, approved bool, approvedBy uuid.UUID, approvalDate time.Time) error
	UpdateTransferShipping(ctx context.Context, orgID uuid.UUID, id uuid.UUID, carrier, trackingNumber *string, shippedBy uuid.UUID) error
	ReceiveTransfer(ctx context.Context, orgID uuid.UUID, id uuid.UUID, receivedBy uuid.UUID) error

	// Transfer Items
	ListTransferItems(ctx context.Context, orgID uuid.UUID, filters TransferItemFilters) ([]InventoryTransferItem, error)
	CountTransferItems(ctx context.Context, orgID uuid.UUID, filters TransferItemFilters) (int64, error)
	GetTransferItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransferItem, error)
	CreateTransferItem(ctx context.Context, item *InventoryTransferItem) error
	UpdateTransferItem(ctx context.Context, item *InventoryTransferItem) error
	DeleteTransferItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateTransferItemReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID, quantityReceived float64, varianceReason *string) error
	GetTransferItemsByTransfer(ctx context.Context, orgID uuid.UUID, transferID uuid.UUID) ([]InventoryTransferItem, error)
	CalculateTransferValue(ctx context.Context, orgID uuid.UUID, transferID uuid.UUID) (float64, error)
}
