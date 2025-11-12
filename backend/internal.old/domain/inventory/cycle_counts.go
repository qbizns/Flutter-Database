package inventory

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CycleCountType represents the type of cycle count
type CycleCountType string

const (
	CycleCountTypeCycle CycleCountType = "cycle"
	CycleCountTypeFull  CycleCountType = "full"
	CycleCountTypeSpot  CycleCountType = "spot"
	CycleCountTypeBlind CycleCountType = "blind"
)

// Scan implements the sql.Scanner interface
func (c *CycleCountType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*c = CycleCountType(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*c = CycleCountType(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (c CycleCountType) Value() (driver.Value, error) {
	return string(c), nil
}

// CycleCountStatus represents the status of a cycle count
type CycleCountStatus string

const (
	CycleCountStatusPlanned    CycleCountStatus = "planned"
	CycleCountStatusInProgress CycleCountStatus = "in_progress"
	CycleCountStatusCompleted  CycleCountStatus = "completed"
	CycleCountStatusCancelled  CycleCountStatus = "cancelled"
)

// Scan implements the sql.Scanner interface
func (c *CycleCountStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*c = CycleCountStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*c = CycleCountStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (c CycleCountStatus) Value() (driver.Value, error) {
	return string(c), nil
}

// CycleCountItemStatus represents the status of a cycle count item
type CycleCountItemStatus string

const (
	CycleCountItemStatusPending   CycleCountItemStatus = "pending"
	CycleCountItemStatusCounted   CycleCountItemStatus = "counted"
	CycleCountItemStatusRecounted CycleCountItemStatus = "recounted"
	CycleCountItemStatusAdjusted  CycleCountItemStatus = "adjusted"
	CycleCountItemStatusApproved  CycleCountItemStatus = "approved"
)

// Scan implements the sql.Scanner interface
func (c *CycleCountItemStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*c = CycleCountItemStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*c = CycleCountItemStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (c CycleCountItemStatus) Value() (driver.Value, error) {
	return string(c), nil
}

// AdjustmentReasonType represents the type of adjustment reason
type AdjustmentReasonType string

const (
	AdjustmentReasonTypeIncrease AdjustmentReasonType = "increase"
	AdjustmentReasonTypeDecrease AdjustmentReasonType = "decrease"
	AdjustmentReasonTypeBoth     AdjustmentReasonType = "both"
)

// Scan implements the sql.Scanner interface
func (a *AdjustmentReasonType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*a = AdjustmentReasonType(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*a = AdjustmentReasonType(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (a AdjustmentReasonType) Value() (driver.Value, error) {
	return string(a), nil
}

// CycleCount represents a physical inventory cycle count
type CycleCount struct {
	ID                 uuid.UUID              `json:"id"`
	OrganizationID     uuid.UUID              `json:"organization_id"`
	LocationID         *uuid.UUID             `json:"location_id,omitempty"`
	CountNumber        string                 `json:"count_number"`
	CountDate          time.Time              `json:"count_date"`
	CountType          CycleCountType         `json:"count_type"`
	Status             CycleCountStatus       `json:"status"`
	CategoryID         *uuid.UUID             `json:"category_id,omitempty"`
	IncludeZeroStock   bool                   `json:"include_zero_stock"`
	TotalItemsPlanned  int                    `json:"total_items_planned"`
	TotalItemsCounted  int                    `json:"total_items_counted"`
	ItemsWithVariance  int                    `json:"items_with_variance"`
	TotalVarianceValue float64                `json:"total_variance_value"`
	ScheduledDate      *time.Time             `json:"scheduled_date,omitempty"`
	StartedAt          *time.Time             `json:"started_at,omitempty"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
	Notes              *string                `json:"notes,omitempty"`
	Metadata           json.RawMessage        `json:"metadata,omitempty"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	DeletedAt          *time.Time             `json:"deleted_at,omitempty"`
	CreatedBy          *uuid.UUID             `json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID             `json:"updated_by,omitempty"`
	CountedBy          *uuid.UUID             `json:"counted_by,omitempty"`
	ApprovedBy         *uuid.UUID             `json:"approved_by,omitempty"`
}

// CycleCountItem represents an item in a cycle count
type CycleCountItem struct {
	ID                  uuid.UUID                `json:"id"`
	OrganizationID      uuid.UUID                `json:"organization_id"`
	CycleCountID        uuid.UUID                `json:"cycle_count_id"`
	ProductID           uuid.UUID                `json:"product_id"`
	ProductVariantID    *uuid.UUID               `json:"product_variant_id,omitempty"`
	ProductName         string                   `json:"product_name"`
	ProductSKU          *string                  `json:"product_sku,omitempty"`
	SystemQuantity      float64                  `json:"system_quantity"`
	CountedQuantity     *float64                 `json:"counted_quantity,omitempty"`
	VarianceQuantity    *float64                 `json:"variance_quantity,omitempty"`
	VariancePercentage  *float64                 `json:"variance_percentage,omitempty"`
	UnitCost            *float64                 `json:"unit_cost,omitempty"`
	VarianceValue       *float64                 `json:"variance_value,omitempty"`
	Status              CycleCountItemStatus     `json:"status"`
	RecountRequired     bool                     `json:"recount_required"`
	RecountQuantity     *float64                 `json:"recount_quantity,omitempty"`
	RecountReason       *string                  `json:"recount_reason,omitempty"`
	AdjustmentApplied   bool                     `json:"adjustment_applied"`
	AdjustmentDate      *time.Time               `json:"adjustment_date,omitempty"`
	AdjustmentReason    *string                  `json:"adjustment_reason,omitempty"`
	Notes               *string                  `json:"notes,omitempty"`
	Metadata            json.RawMessage          `json:"metadata,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	CountedAt           *time.Time               `json:"counted_at,omitempty"`
	CountedBy           *uuid.UUID               `json:"counted_by,omitempty"`
}

// StockAdjustmentReason represents a pre-defined reason for stock adjustment
type StockAdjustmentReason struct {
	ID               uuid.UUID             `json:"id"`
	OrganizationID   *uuid.UUID            `json:"organization_id,omitempty"`
	Code             string                `json:"code"`
	Name             string                `json:"name"`
	Description      *string               `json:"description,omitempty"`
	ReasonType       AdjustmentReasonType  `json:"reason_type"`
	IsSystemReason   bool                  `json:"is_system_reason"`
	IsActive         bool                  `json:"is_active"`
	RequiresApproval bool                  `json:"requires_approval"`
	RequiresNotes    bool                  `json:"requires_notes"`
	SortOrder        int                   `json:"sort_order"`
	Metadata         json.RawMessage       `json:"metadata,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	DeletedAt        *time.Time            `json:"deleted_at,omitempty"`
}

// CycleCountFilters represents filters for listing cycle counts
type CycleCountFilters struct {
	Status         *CycleCountStatus
	CountType      *CycleCountType
	LocationID     *uuid.UUID
	CategoryID     *uuid.UUID
	CountNumberLike *string
	StartDate      *time.Time
	EndDate        *time.Time
	HasVariance    *bool
	Page           int
	PageSize       int
}

// CycleCountItemFilters represents filters for listing cycle count items
type CycleCountItemFilters struct {
	CycleCountID    uuid.UUID
	Status          *CycleCountItemStatus
	ProductID       *uuid.UUID
	HasVariance     *bool
	RecountRequired *bool
	Page            int
	PageSize        int
}

// AdjustmentReasonFilters represents filters for listing adjustment reasons
type AdjustmentReasonFilters struct {
	Code           *string
	ReasonType     *AdjustmentReasonType
	IsActive       *bool
	IsSystemReason *bool
	Page           int
	PageSize       int
}

// CycleCountRepository defines the cycle count data access interface
type CycleCountRepository interface {
	// Cycle Counts
	ListCycleCounts(ctx context.Context, orgID uuid.UUID, filters CycleCountFilters) ([]CycleCount, error)
	CountCycleCounts(ctx context.Context, orgID uuid.UUID, filters CycleCountFilters) (int64, error)
	GetCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CycleCount, error)
	GetCycleCountByNumber(ctx context.Context, orgID uuid.UUID, countNumber string) (*CycleCount, error)
	CreateCycleCount(ctx context.Context, count *CycleCount) error
	UpdateCycleCount(ctx context.Context, count *CycleCount) error
	DeleteCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateCycleCountStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status CycleCountStatus, countedBy *uuid.UUID) error
	CompleteCycleCount(ctx context.Context, orgID uuid.UUID, id uuid.UUID, approvedBy uuid.UUID) error
	RecalculateCycleCountStats(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Cycle Count Items
	ListCycleCountItems(ctx context.Context, orgID uuid.UUID, filters CycleCountItemFilters) ([]CycleCountItem, error)
	CountCycleCountItems(ctx context.Context, orgID uuid.UUID, filters CycleCountItemFilters) (int64, error)
	GetCycleCountItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CycleCountItem, error)
	CreateCycleCountItem(ctx context.Context, item *CycleCountItem) error
	UpdateCycleCountItem(ctx context.Context, item *CycleCountItem) error
	DeleteCycleCountItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateCountedQuantity(ctx context.Context, orgID uuid.UUID, id uuid.UUID, countedQuantity float64, countedBy uuid.UUID) error
	GetCycleCountItemsByCount(ctx context.Context, orgID uuid.UUID, countID uuid.UUID) ([]CycleCountItem, error)
	GetItemsWithVariance(ctx context.Context, orgID uuid.UUID, countID uuid.UUID) ([]CycleCountItem, error)
	ApplyAdjustment(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID, adjustmentReason string, approvedBy uuid.UUID) error

	// Stock Adjustment Reasons
	ListAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, filters AdjustmentReasonFilters) ([]StockAdjustmentReason, error)
	CountAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, filters AdjustmentReasonFilters) (int64, error)
	GetAdjustmentReason(ctx context.Context, id uuid.UUID) (*StockAdjustmentReason, error)
	GetAdjustmentReasonByCode(ctx context.Context, orgID *uuid.UUID, code string) (*StockAdjustmentReason, error)
	CreateAdjustmentReason(ctx context.Context, reason *StockAdjustmentReason) error
	UpdateAdjustmentReason(ctx context.Context, reason *StockAdjustmentReason) error
	DeleteAdjustmentReason(ctx context.Context, id uuid.UUID) error
	GetActiveAdjustmentReasons(ctx context.Context, orgID *uuid.UUID, reasonType AdjustmentReasonType) ([]StockAdjustmentReason, error)
}
