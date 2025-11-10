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
