package inventory

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SerialNumberStatus represents the status of a serial number
type SerialNumberStatus string

const (
	SerialStatusInStock   SerialNumberStatus = "in_stock"
	SerialStatusSold      SerialNumberStatus = "sold"
	SerialStatusReserved  SerialNumberStatus = "reserved"
	SerialStatusDamaged   SerialNumberStatus = "damaged"
	SerialStatusReturned  SerialNumberStatus = "returned"
	SerialStatusInRepair  SerialNumberStatus = "in_repair"
	SerialStatusScrapped  SerialNumberStatus = "scrapped"
)

// Scan implements the sql.Scanner interface
func (s *SerialNumberStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*s = SerialNumberStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*s = SerialNumberStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (s SerialNumberStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// ProductSerialNumber represents an individual product tracked by serial number
type ProductSerialNumber struct {
	ID                 uuid.UUID            `json:"id"`
	OrganizationID     uuid.UUID            `json:"organization_id"`
	ProductID          uuid.UUID            `json:"product_id"`
	ProductVariantID   *uuid.UUID           `json:"product_variant_id,omitempty"`
	LocationID         *uuid.UUID           `json:"location_id,omitempty"`
	SerialNumber       string               `json:"serial_number"`
	Status             SerialNumberStatus   `json:"status"`
	PurchaseOrderID    *uuid.UUID           `json:"purchase_order_id,omitempty"`
	PurchaseDate       *time.Time           `json:"purchase_date,omitempty"`
	PurchaseCost       *float64             `json:"purchase_cost,omitempty"`
	SupplierID         *uuid.UUID           `json:"supplier_id,omitempty"`
	SaleID             *uuid.UUID           `json:"sale_id,omitempty"`
	SaleDate           *time.Time           `json:"sale_date,omitempty"`
	SalePrice          *float64             `json:"sale_price,omitempty"`
	CustomerID         *uuid.UUID           `json:"customer_id,omitempty"`
	WarrantyStartDate  *time.Time           `json:"warranty_start_date,omitempty"`
	WarrantyEndDate    *time.Time           `json:"warranty_end_date,omitempty"`
	WarrantyProvider   *string              `json:"warranty_provider,omitempty"`
	WarrantyTerms      *string              `json:"warranty_terms,omitempty"`
	Notes              *string              `json:"notes,omitempty"`
	Metadata           json.RawMessage      `json:"metadata,omitempty"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
	DeletedAt          *time.Time           `json:"deleted_at,omitempty"`
	CreatedBy          *uuid.UUID           `json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID           `json:"updated_by,omitempty"`
}

// BatchStatus represents the status of a batch/lot
type BatchStatus string

const (
	BatchStatusActive    BatchStatus = "active"
	BatchStatusExpired   BatchStatus = "expired"
	BatchStatusRecalled  BatchStatus = "recalled"
	BatchStatusQuarantine BatchStatus = "quarantine"
	BatchStatusDepleted  BatchStatus = "depleted"
)

// Scan implements the sql.Scanner interface
func (b *BatchStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*b = BatchStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*b = BatchStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (b BatchStatus) Value() (driver.Value, error) {
	return string(b), nil
}

// QualityStatus represents the quality control status of a batch
type QualityStatus string

const (
	QualityStatusPending    QualityStatus = "pending"
	QualityStatusPassed     QualityStatus = "passed"
	QualityStatusFailed     QualityStatus = "failed"
	QualityStatusQuarantine QualityStatus = "quarantine"
)

// Scan implements the sql.Scanner interface
func (q *QualityStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*q = QualityStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*q = QualityStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (q QualityStatus) Value() (driver.Value, error) {
	return string(q), nil
}

// ProductBatch represents a batch or lot of products
type ProductBatch struct {
	ID                   uuid.UUID           `json:"id"`
	OrganizationID       uuid.UUID           `json:"organization_id"`
	ProductID            uuid.UUID           `json:"product_id"`
	ProductVariantID     *uuid.UUID          `json:"product_variant_id,omitempty"`
	LocationID           *uuid.UUID          `json:"location_id,omitempty"`
	BatchNumber          string              `json:"batch_number"`
	LotNumber            *string             `json:"lot_number,omitempty"`
	Status               BatchStatus         `json:"status"`
	InitialQuantity      float64             `json:"initial_quantity"`
	CurrentQuantity      float64             `json:"current_quantity"`
	UnitOfMeasure        string              `json:"unit_of_measure"`
	ManufacturingDate    *time.Time          `json:"manufacturing_date,omitempty"`
	ExpirationDate       *time.Time          `json:"expiration_date,omitempty"`
	ReceivedDate         time.Time           `json:"received_date"`
	PurchaseOrderID      *uuid.UUID          `json:"purchase_order_id,omitempty"`
	SupplierID           *uuid.UUID          `json:"supplier_id,omitempty"`
	SupplierBatchNumber  *string             `json:"supplier_batch_number,omitempty"`
	UnitCost             *float64            `json:"unit_cost,omitempty"`
	TotalCost            *float64            `json:"total_cost,omitempty"`
	QualityStatus        QualityStatus       `json:"quality_status"`
	QualityCheckDate     *time.Time          `json:"quality_check_date,omitempty"`
	QualityCheckedBy     *uuid.UUID          `json:"quality_checked_by,omitempty"`
	QualityNotes         *string             `json:"quality_notes,omitempty"`
	Notes                *string             `json:"notes,omitempty"`
	Metadata             json.RawMessage     `json:"metadata,omitempty"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	DeletedAt            *time.Time          `json:"deleted_at,omitempty"`
	CreatedBy            *uuid.UUID          `json:"created_by,omitempty"`
	UpdatedBy            *uuid.UUID          `json:"updated_by,omitempty"`
}

// BatchTransactionType represents the type of batch transaction
type BatchTransactionType string

const (
	BatchTxTypeSale       BatchTransactionType = "sale"
	BatchTxTypeAdjustment BatchTransactionType = "adjustment"
	BatchTxTypeReturn     BatchTransactionType = "return"
	BatchTxTypeWaste      BatchTransactionType = "waste"
	BatchTxTypeTransfer   BatchTransactionType = "transfer"
	BatchTxTypeExpiration BatchTransactionType = "expiration"
)

// Scan implements the sql.Scanner interface
func (t *BatchTransactionType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*t = BatchTransactionType(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*t = BatchTransactionType(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (t BatchTransactionType) Value() (driver.Value, error) {
	return string(t), nil
}

// BatchTransaction represents a transaction against a batch
type BatchTransaction struct {
	ID                  uuid.UUID             `json:"id"`
	OrganizationID      uuid.UUID             `json:"organization_id"`
	BatchID             uuid.UUID             `json:"batch_id"`
	TransactionType     BatchTransactionType  `json:"transaction_type"`
	Quantity            float64               `json:"quantity"`
	BalanceAfter        float64               `json:"balance_after"`
	SaleID              *uuid.UUID            `json:"sale_id,omitempty"`
	InventoryTransferID *uuid.UUID            `json:"inventory_transfer_id,omitempty"`
	Reason              *string               `json:"reason,omitempty"`
	Notes               *string               `json:"notes,omitempty"`
	Metadata            json.RawMessage       `json:"metadata,omitempty"`
	TransactionDate     time.Time             `json:"transaction_date"`
	CreatedBy           *uuid.UUID            `json:"created_by,omitempty"`
}

// SerialNumberFilters represents filters for listing serial numbers
type SerialNumberFilters struct {
	Status           *SerialNumberStatus
	ProductID        *uuid.UUID
	LocationID       *uuid.UUID
	SerialLike       *string
	SaleID           *uuid.UUID
	HasWarranty      *bool
	WarrantyExpiring *bool
	Page             int
	PageSize         int
}

// BatchFilters represents filters for listing batches
type BatchFilters struct {
	Status              *BatchStatus
	ProductID           *uuid.UUID
	LocationID          *uuid.UUID
	BatchNumberLike     *string
	QualityStatus       *QualityStatus
	ExpiringBefore      *time.Time
	ExpiredAfter        *time.Time
	HasQualityIssues    *bool
	MinStockLevel       *float64
	Page                int
	PageSize            int
}

// BatchTransactionFilters represents filters for listing batch transactions
type BatchTransactionFilters struct {
	BatchID         uuid.UUID
	TransactionType *BatchTransactionType
	StartDate       *time.Time
	EndDate         *time.Time
	Page            int
	PageSize        int
}

// SerialNumberRepository defines the serial number data access interface
type SerialNumberRepository interface {
	ListSerialNumbers(ctx context.Context, orgID uuid.UUID, filters SerialNumberFilters) ([]ProductSerialNumber, error)
	CountSerialNumbers(ctx context.Context, orgID uuid.UUID, filters SerialNumberFilters) (int64, error)
	GetSerialNumber(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductSerialNumber, error)
	GetBySerialNumber(ctx context.Context, orgID uuid.UUID, serialNumber string) (*ProductSerialNumber, error)
	CreateSerialNumber(ctx context.Context, serial *ProductSerialNumber) error
	UpdateSerialNumber(ctx context.Context, serial *ProductSerialNumber) error
	DeleteSerialNumber(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateSerialStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status SerialNumberStatus) error
	GetSerialsByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]ProductSerialNumber, error)
	GetExpiredWarrantySerials(ctx context.Context, orgID uuid.UUID) ([]ProductSerialNumber, error)
}

// BatchRepository defines the batch data access interface
type BatchRepository interface {
	// Batches
	ListBatches(ctx context.Context, orgID uuid.UUID, filters BatchFilters) ([]ProductBatch, error)
	CountBatches(ctx context.Context, orgID uuid.UUID, filters BatchFilters) (int64, error)
	GetBatch(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductBatch, error)
	GetBatchByNumber(ctx context.Context, orgID uuid.UUID, batchNumber string) (*ProductBatch, error)
	CreateBatch(ctx context.Context, batch *ProductBatch) error
	UpdateBatch(ctx context.Context, batch *ProductBatch) error
	DeleteBatch(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateBatchQuantity(ctx context.Context, orgID uuid.UUID, id uuid.UUID, quantity float64) error
	UpdateBatchStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status BatchStatus) error
	UpdateBatchQuality(ctx context.Context, orgID uuid.UUID, id uuid.UUID, qualityStatus QualityStatus, checkedBy uuid.UUID) error
	GetExpiringBatches(ctx context.Context, orgID uuid.UUID, daysUntilExpiry int) ([]ProductBatch, error)
	GetBatchesByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]ProductBatch, error)

	// Batch Transactions
	ListBatchTransactions(ctx context.Context, orgID uuid.UUID, filters BatchTransactionFilters) ([]BatchTransaction, error)
	CountBatchTransactions(ctx context.Context, orgID uuid.UUID, filters BatchTransactionFilters) (int64, error)
	GetBatchTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*BatchTransaction, error)
	CreateBatchTransaction(ctx context.Context, tx *BatchTransaction) error
	GetBatchTransactionHistory(ctx context.Context, orgID uuid.UUID, batchID uuid.UUID) ([]BatchTransaction, error)
	CalculateBatchCurrentQuantity(ctx context.Context, orgID uuid.UUID, batchID uuid.UUID) (float64, error)
}
