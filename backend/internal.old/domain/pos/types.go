package pos

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionStatus represents the status of a POS session
type SessionStatus string

const (
	SessionStatusOpen      SessionStatus = "open"
	SessionStatusClosing   SessionStatus = "closing"
	SessionStatusClosed    SessionStatus = "closed"
	SessionStatusReconciled SessionStatus = "reconciled"
)

// MovementType represents the type of cash movement
type MovementType string

const (
	MovementTypePayIn    MovementType = "pay_in"
	MovementTypePayOut   MovementType = "pay_out"
	MovementTypeFloatAdd MovementType = "float_add"
	MovementTypeDropToSafe MovementType = "drop_to_safe"
)

// POSSession represents a POS register session
type POSSession struct {
	ID              uuid.UUID      `json:"id"`
	OrganizationID  uuid.UUID      `json:"organization_id"`
	SessionNumber   string         `json:"session_number"`
	SessionName     *string        `json:"session_name"`
	DeviceID        *uuid.UUID     `json:"device_id"`
	LocationID      uuid.UUID      `json:"location_id"`
	UserID          uuid.UUID      `json:"user_id"` // Primary cashier
	ShiftID         *uuid.UUID     `json:"shift_id"`

	// Session timing
	OpenedAt  time.Time  `json:"opened_at"`
	ClosedAt  *time.Time `json:"closed_at"`

	// Opening float by payment method
	OpeningCash  float64 `json:"opening_cash"`
	OpeningCard  float64 `json:"opening_card"`
	OpeningOther float64 `json:"opening_other"`

	// Expected closing (system calculated)
	ExpectedCash  float64 `json:"expected_cash"`
	ExpectedCard  float64 `json:"expected_card"`
	ExpectedOther float64 `json:"expected_other"`

	// Actual counted closing
	CountedCash  *float64 `json:"counted_cash"`
	CountedCard  *float64 `json:"counted_card"`
	CountedOther *float64 `json:"counted_other"`

	// Differences (counted - expected)
	DifferenceCash  float64 `json:"difference_cash"`
	DifferenceCard  float64 `json:"difference_card"`
	DifferenceOther float64 `json:"difference_other"`

	// Status
	Status SessionStatus `json:"status"`

	// Z-report reference
	ZReportNumber *string `json:"z_report_number"`

	// Metadata
	Notes     *string    `json:"notes"`
	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CashDrawer represents a physical cash drawer
type CashDrawer struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	DrawerCode      string     `json:"drawer_code"`
	DrawerName      string     `json:"drawer_name"`
	LocationID      uuid.UUID  `json:"location_id"`
	DeviceID        *uuid.UUID `json:"device_id"`
	IsActive        bool       `json:"is_active"`
	Notes           *string    `json:"notes"`
	CreatedBy       *uuid.UUID `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// CashDrawerSession links a cash drawer to a POS session
type CashDrawerSession struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	CashDrawerID    uuid.UUID  `json:"cash_drawer_id"`
	POSSessionID    uuid.UUID  `json:"pos_session_id"`
	OpeningAmount   float64    `json:"opening_amount"`
	ClosingAmount   *float64   `json:"closing_amount"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// CashMovement represents a non-sale cash movement
type CashMovement struct {
	ID                uuid.UUID    `json:"id"`
	OrganizationID    uuid.UUID    `json:"organization_id"`
	POSSessionID      uuid.UUID    `json:"pos_session_id"`
	CashDrawerID      *uuid.UUID   `json:"cash_drawer_id"`
	MovementType      MovementType `json:"movement_type"`
	Amount            float64      `json:"amount"`
	ReasonCode        *string      `json:"reason_code"`
	ReasonDescription string       `json:"reason_description"`
	UserID            uuid.UUID    `json:"user_id"`
	RequiresApproval  bool         `json:"requires_approval"`
	ApprovedBy        *uuid.UUID   `json:"approved_by"`
	ApprovedAt        *time.Time   `json:"approved_at"`
	Notes             *string      `json:"notes"`
	CreatedAt         time.Time    `json:"created_at"`
	DeletedAt         *time.Time   `json:"deleted_at,omitempty"`
}

// Filters for listing operations

// POSSessionFilters represents filters for listing POS sessions
type POSSessionFilters struct {
	Status       *SessionStatus
	LocationID   *uuid.UUID
	UserID       *uuid.UUID
	DeviceID     *uuid.UUID
	DateFrom     *time.Time
	DateTo       *time.Time
	SearchTerm   *string
	Page         int
	PageSize     int
}

// CashDrawerFilters represents filters for listing cash drawers
type CashDrawerFilters struct {
	LocationID *uuid.UUID
	IsActive   *bool
	SearchTerm *string
	Page       int
	PageSize   int
}

// CashMovementFilters represents filters for listing cash movements
type CashMovementFilters struct {
	POSSessionID     *uuid.UUID
	CashDrawerID     *uuid.UUID
	MovementType     *MovementType
	UserID           *uuid.UUID
	DateFrom         *time.Time
	DateTo           *time.Time
	RequiresApproval *bool
	Page             int
	PageSize         int
}

// Repository interfaces

// POSSessionRepository defines data access for POS sessions
type POSSessionRepository interface {
	// CRUD operations
	Create(ctx context.Context, session *POSSession) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*POSSession, error)
	GetBySessionNumber(ctx context.Context, orgID uuid.UUID, sessionNumber string) (*POSSession, error)
	List(ctx context.Context, orgID uuid.UUID, filters POSSessionFilters) ([]POSSession, error)
	Count(ctx context.Context, orgID uuid.UUID, filters POSSessionFilters) (int64, error)
	Update(ctx context.Context, session *POSSession) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Session-specific operations
	GetOpenSession(ctx context.Context, orgID uuid.UUID, deviceID uuid.UUID) (*POSSession, error)
	UpdateSessionStatus(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, status SessionStatus) error
	UpdateExpectedAmounts(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, cash, card, other float64) error
	UpdateCountedAmounts(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, cash, card, other float64) error
	RecalculateDifferences(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) error
}

// CashDrawerRepository defines data access for cash drawers
type CashDrawerRepository interface {
	// CRUD operations
	Create(ctx context.Context, drawer *CashDrawer) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CashDrawer, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*CashDrawer, error)
	List(ctx context.Context, orgID uuid.UUID, filters CashDrawerFilters) ([]CashDrawer, error)
	Count(ctx context.Context, orgID uuid.UUID, filters CashDrawerFilters) (int64, error)
	Update(ctx context.Context, drawer *CashDrawer) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Drawer-specific operations
	GetByLocation(ctx context.Context, orgID uuid.UUID, locationID uuid.UUID) ([]CashDrawer, error)
}

// CashDrawerSessionRepository defines data access for cash drawer sessions
type CashDrawerSessionRepository interface {
	// CRUD operations
	Create(ctx context.Context, session *CashDrawerSession) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CashDrawerSession, error)
	List(ctx context.Context, orgID uuid.UUID, drawerID uuid.UUID) ([]CashDrawerSession, error)
	Update(ctx context.Context, session *CashDrawerSession) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Query operations
	GetByDrawerAndSession(ctx context.Context, orgID uuid.UUID, drawerID, sessionID uuid.UUID) (*CashDrawerSession, error)
	GetForSession(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) ([]CashDrawerSession, error)
}

// CashMovementRepository defines data access for cash movements
type CashMovementRepository interface {
	// CRUD operations
	Create(ctx context.Context, movement *CashMovement) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CashMovement, error)
	List(ctx context.Context, orgID uuid.UUID, filters CashMovementFilters) ([]CashMovement, error)
	Count(ctx context.Context, orgID uuid.UUID, filters CashMovementFilters) (int64, error)
	Update(ctx context.Context, movement *CashMovement) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Movement-specific operations
	GetForSession(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) ([]CashMovement, error)
	GetForDrawer(ctx context.Context, orgID uuid.UUID, drawerID uuid.UUID) ([]CashMovement, error)
	GetPendingApprovals(ctx context.Context, orgID uuid.UUID) ([]CashMovement, error)
	ApproveMovement(ctx context.Context, orgID uuid.UUID, movementID uuid.UUID, approvedBy uuid.UUID) error
	CalculateSessionBalance(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) (float64, error)
}
