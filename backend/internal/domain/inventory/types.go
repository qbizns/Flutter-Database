package inventory

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the inventory transaction type
type TransactionType string

const (
	TransactionTypeSale       TransactionType = "sale"
	TransactionTypePurchase   TransactionType = "purchase"
	TransactionTypeAdjustment TransactionType = "adjustment"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeReturn     TransactionType = "return"
	TransactionTypeWriteOff   TransactionType = "write_off"
)

// Scan implements the sql.Scanner interface
func (t *TransactionType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*t = TransactionType(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*t = TransactionType(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (t TransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

// InventoryTransaction represents an inventory transaction
type InventoryTransaction struct {
	ID                uuid.UUID       `json:"id"`
	OrganizationID    uuid.UUID       `json:"organization_id"`
	ProductID         uuid.UUID       `json:"product_id"`
	TransactionType   TransactionType `json:"transaction_type"`
	Quantity          float64         `json:"quantity"`
	Unit              string          `json:"unit"`
	BalanceAfter      float64         `json:"balance_after"`
	SaleID            *uuid.UUID      `json:"sale_id,omitempty"`
	ReferenceNumber   *string         `json:"reference_number,omitempty"`
	UnitCost          *float64        `json:"unit_cost,omitempty"`
	TotalCost         *float64        `json:"total_cost,omitempty"`
	TransactionDate   time.Time       `json:"transaction_date"`
	Notes             string          `json:"notes"`
	Reason            string          `json:"reason"`
	Metadata          interface{}     `json:"metadata"`
	CreatedAt         time.Time       `json:"created_at"`
	CreatedBy         uuid.UUID       `json:"created_by"`
}

// InventoryTransactionFilters represents filters for listing inventory transactions
type InventoryTransactionFilters struct {
	ProductID       *uuid.UUID
	TransactionType *TransactionType
	SaleID          *uuid.UUID
	StartDate       *time.Time
	EndDate         *time.Time
	Page            int
	PageSize        int
}

// Repository defines the inventory transaction data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters InventoryTransactionFilters) ([]InventoryTransaction, error)
	Count(ctx context.Context, orgID uuid.UUID, filters InventoryTransactionFilters) (int64, error)
	Create(ctx context.Context, transaction *InventoryTransaction) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransaction, error)
	GetByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, limit int, offset int) ([]InventoryTransaction, error)
	GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]InventoryTransaction, error)
	Update(ctx context.Context, transaction *InventoryTransaction) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetLatestBalanceForProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) (*InventoryTransaction, error)
}

// ============================================================================
// INVENTORY VALUATION
// ============================================================================

// InventoryValuationSettings represents per-organization inventory valuation configuration
type InventoryValuationSettings struct {
	ID                                     uuid.UUID   `json:"id"`
	OrganizationID                         uuid.UUID   `json:"organization_id"`
	ValuationMethod                        string      `json:"valuation_method"` // fifo, lifo, weighted_average, moving_average, standard_cost, specific_id
	CostLayerGranularity                   string      `json:"cost_layer_granularity"` // product, product_location, product_location_lot, serial_number
	DefaultInventoryAccountID              *uuid.UUID  `json:"default_inventory_account_id"`
	DefaultCOGSAccountID                   *uuid.UUID  `json:"default_cogs_account_id"`
	DefaultInventoryAdjustmentAccountID    *uuid.UUID  `json:"default_inventory_adjustment_account_id"`
	DefaultInventoryVarianceAccountID      *uuid.UUID  `json:"default_inventory_variance_account_id"`
	COGSRecognitionTiming                  string      `json:"cogs_recognition_timing"` // on_sale, on_delivery, on_payment
	AllowNegativeInventory                 bool        `json:"allow_negative_inventory"`
	RevalueOnPurchase                      bool        `json:"revalue_on_purchase"`
	RoundUnitCostToDecimals                int         `json:"round_unit_cost_to_decimals"`
	AutoCreateCostLayersOnPurchase         bool        `json:"auto_create_cost_layers_on_purchase"`
	AutoConsumeCostLayersOnSale            bool        `json:"auto_consume_cost_layers_on_sale"`
	RecalculateInventoryValueOnAdjustment  bool        `json:"recalculate_inventory_value_on_adjustment"`
	CreatedAt                              time.Time   `json:"created_at"`
	UpdatedAt                              time.Time   `json:"updated_at"`
}

// InventoryCostLayer represents a cost layer for FIFO/LIFO inventory valuation
type InventoryCostLayer struct {
	ID                    uuid.UUID       `json:"id"`
	OrganizationID        uuid.UUID       `json:"organization_id"`
	ProductID             uuid.UUID       `json:"product_id"`
	LocationID            *uuid.UUID      `json:"location_id"`
	LotNumber             string          `json:"lot_number"`
	SerialNumber          string          `json:"serial_number"`
	LayerDate             time.Time       `json:"layer_date"`
	UnitCost              float64         `json:"unit_cost"`
	OriginalQuantity      float64         `json:"original_quantity"`
	RemainingQuantity     float64         `json:"remaining_quantity"`
	UOMCode               string          `json:"uom_code"`
	SourceTransactionType string          `json:"source_transaction_type"` // purchase, transfer_in, adjustment, production
	SourceTransactionID   *uuid.UUID      `json:"source_transaction_id"`
	SourceReference       string          `json:"source_reference"`
	IsFullyConsumed       bool            `json:"is_fully_consumed"`
	ConsumedAt            *time.Time      `json:"consumed_at"`
	Metadata              json.RawMessage `json:"metadata"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	DeletedAt             *time.Time      `json:"deleted_at"`
}

// InventoryValuationSettingsRepository defines data access interface
type InventoryValuationSettingsRepository interface {
	Get(ctx context.Context, orgID uuid.UUID) (*InventoryValuationSettings, error)
	Create(ctx context.Context, settings *InventoryValuationSettings) error
	Update(ctx context.Context, settings *InventoryValuationSettings) error
}

// InventoryCostLayerRepository defines data access interface
type InventoryCostLayerRepository interface {
	List(ctx context.Context, orgID uuid.UUID, productID *uuid.UUID, locationID *uuid.UUID) ([]InventoryCostLayer, error)
	Create(ctx context.Context, layer *InventoryCostLayer) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryCostLayer, error)
	Update(ctx context.Context, layer *InventoryCostLayer) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetAvailableLayers(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) ([]InventoryCostLayer, error)
	ConsumeQuantity(ctx context.Context, layerID uuid.UUID, quantity float64) error
	GetOldestLayer(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) (*InventoryCostLayer, error)
	GetNewestLayer(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, locationID *uuid.UUID) (*InventoryCostLayer, error)
}
