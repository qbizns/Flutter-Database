package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// KitchenStation represents a kitchen station/prep area
type KitchenStation struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	LocationID       *uuid.UUID `json:"location_id,omitempty"`
	StationName      string     `json:"station_name"`
	StationCode      string     `json:"station_code"`
	StationType      string     `json:"station_type"` // kitchen, bar, dessert, prep
	Description      *string    `json:"description,omitempty"`
	DisplayOrder     int        `json:"display_order"`
	ColorCode        *string    `json:"color_code,omitempty"`
	PrinterID        *uuid.UUID `json:"printer_id,omitempty"`
	IsActive         bool       `json:"is_active"`
	AutoPrintTickets bool       `json:"auto_print_tickets"`
	AlertSoundEnabled bool      `json:"alert_sound_enabled"`
	DisplayConfig    any        `json:"display_config,omitempty"` // JSONB
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// Order represents a restaurant order
type Order struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	LocationID      *uuid.UUID `json:"location_id,omitempty"`
	OrderNumber     string     `json:"order_number"`
	DisplayNumber   *int       `json:"display_number,omitempty"`
	OrderType       string     `json:"order_type"` // dine_in, takeout, delivery, online
	TableID         *uuid.UUID `json:"table_id,omitempty"`
	ReservationID   *uuid.UUID `json:"reservation_id,omitempty"`
	Covers          int        `json:"covers"`
	CustomerID      *uuid.UUID `json:"customer_id,omitempty"`
	WaiterID        *uuid.UUID `json:"waiter_id,omitempty"`
	Status          string     `json:"status"` // draft, submitted, sent_to_kitchen, preparing, ready, served, completed, cancelled, on_hold
	OrderDate       time.Time  `json:"order_date"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	KitchenReceivedAt *time.Time `json:"kitchen_received_at,omitempty"`
	ReadyAt         *time.Time `json:"ready_at,omitempty"`
	ServedAt        *time.Time `json:"served_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Subtotal        float64    `json:"subtotal"`
	TaxAmount       float64    `json:"tax_amount"`
	DiscountAmount  float64    `json:"discount_amount"`
	ServiceCharge   float64    `json:"service_charge"`
	TotalAmount     float64    `json:"total_amount"`
	SaleID          *uuid.UUID `json:"sale_id,omitempty"`
	ShiftID         *uuid.UUID `json:"shift_id,omitempty"`
	CustomerNotes   *string    `json:"customer_notes,omitempty"`
	KitchenNotes    *string    `json:"kitchen_notes,omitempty"`
	InternalNotes   *string    `json:"internal_notes,omitempty"`
	Metadata        any        `json:"metadata,omitempty"` // JSONB
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// OrderItem represents a line item in an order
type OrderItem struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	OrderID            uuid.UUID  `json:"order_id"`
	ProductID          uuid.UUID  `json:"product_id"`
	ProductVariantID   *uuid.UUID `json:"product_variant_id,omitempty"`
	ItemName           string     `json:"item_name"`
	Quantity           float64    `json:"quantity"`
	UnitPrice          float64    `json:"unit_price"`
	CourseID           *uuid.UUID `json:"course_id,omitempty"`
	CoursePosition     int        `json:"course_position"`
	FireTime           *time.Time `json:"fire_time,omitempty"`
	KitchenStationID   *uuid.UUID `json:"kitchen_station_id,omitempty"`
	KitchenTicketID    *uuid.UUID `json:"kitchen_ticket_id,omitempty"`
	Status             string     `json:"status"` // pending, fired, acknowledged, preparing, ready, served, cancelled, on_hold, voided
	FiredAt            *time.Time `json:"fired_at,omitempty"`
	AcknowledgedAt     *time.Time `json:"acknowledged_at,omitempty"`
	StartedPreparingAt *time.Time `json:"started_preparing_at,omitempty"`
	ReadyAt            *time.Time `json:"ready_at,omitempty"`
	ServedAt           *time.Time `json:"served_at,omitempty"`
	ModifiersTotal     float64    `json:"modifiers_total"`
	DiscountAmount     float64    `json:"discount_amount"`
	LineTotal          float64    `json:"line_total"`
	SpecialInstructions *string   `json:"special_instructions,omitempty"`
	CustomerNotes      *string    `json:"customer_notes,omitempty"`
	KitchenNotes       *string    `json:"kitchen_notes,omitempty"`
	SeatNumber         *int       `json:"seat_number,omitempty"`
	Metadata           any        `json:"metadata,omitempty"` // JSONB
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// OrderItemModifier represents a selected modifier for an order item
type OrderItemModifier struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	OrderItemID     uuid.UUID  `json:"order_item_id"`
	ModifierID      uuid.UUID  `json:"modifier_id"`
	ModifierGroupID *uuid.UUID `json:"modifier_group_id,omitempty"`
	ModifierName    string     `json:"modifier_name"`
	Quantity        int        `json:"quantity"`
	PriceAdjustment float64    `json:"price_adjustment"`
	DisplayOrder    int        `json:"display_order"`
	Metadata        any        `json:"metadata,omitempty"` // JSONB
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// KitchenTicket represents a KDS ticket for a kitchen station
type KitchenTicket struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	LocationID         *uuid.UUID `json:"location_id,omitempty"`
	TicketNumber       string     `json:"ticket_number"`
	DisplaySequence    *int       `json:"display_sequence,omitempty"`
	OrderID            uuid.UUID  `json:"order_id"`
	KitchenStationID   uuid.UUID  `json:"kitchen_station_id"`
	CourseID           *uuid.UUID `json:"course_id,omitempty"`
	TicketType         string     `json:"ticket_type"` // normal, rush, remake, special
	Priority           int        `json:"priority"`
	Status             string     `json:"status"` // new, acknowledged, preparing, ready, served, completed, cancelled, bumped
	CreatedAt          time.Time  `json:"created_at"`
	FiredAt            *time.Time `json:"fired_at,omitempty"`
	AcknowledgedAt     *time.Time `json:"acknowledged_at,omitempty"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	ReadyAt            *time.Time `json:"ready_at,omitempty"`
	BumpedAt           *time.Time `json:"bumped_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	PrepTimeMinutes    *int       `json:"prep_time_minutes,omitempty"`
	TargetPrepTime     *int       `json:"target_prep_time,omitempty"`
	TableNumber        *string    `json:"table_number,omitempty"`
	OrderType          *string    `json:"order_type,omitempty"`
	Covers             *int       `json:"covers,omitempty"`
	WaiterName         *string    `json:"waiter_name,omitempty"`
	SpecialInstructions *string   `json:"special_instructions,omitempty"`
	KitchenNotes       *string    `json:"kitchen_notes,omitempty"`
	DisplayConfig      any        `json:"display_config,omitempty"` // JSONB
	Metadata           any        `json:"metadata,omitempty"`       // JSONB
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy          *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// Filter types
type KitchenStationFilters struct {
	Search     string
	StationType *string
	IsActive   *bool
	LocationID *uuid.UUID
	Page       int
	PageSize   int
}

type OrderFilters struct {
	Search     string
	Status     *string
	OrderType  *string
	TableID    *uuid.UUID
	WaiterID   *uuid.UUID
	CustomerID *uuid.UUID
	LocationID *uuid.UUID
	FromDate   *time.Time
	ToDate     *time.Time
	Page       int
	PageSize   int
}

type OrderItemFilters struct {
	Search       string
	OrderID      *uuid.UUID
	Status       *string
	CourseID     *uuid.UUID
	StationID    *uuid.UUID
	Page         int
	PageSize     int
}

type KitchenTicketFilters struct {
	Search      string
	Status      *string
	TicketType  *string
	StationID   *uuid.UUID
	OrderID     *uuid.UUID
	LocationID  *uuid.UUID
	FromDate    *time.Time
	ToDate      *time.Time
	Page        int
	PageSize    int
}

// Repository interfaces
type KitchenStationRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters KitchenStationFilters) ([]KitchenStation, error)
	Count(ctx context.Context, orgID uuid.UUID, filters KitchenStationFilters) (int64, error)
	Create(ctx context.Context, station *KitchenStation) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenStation, error)
	Update(ctx context.Context, station *KitchenStation) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

type OrderRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters OrderFilters) ([]Order, error)
	Count(ctx context.Context, orgID uuid.UUID, filters OrderFilters) (int64, error)
	Create(ctx context.Context, order *Order) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Order, error)
	GetByNumber(ctx context.Context, orgID uuid.UUID, orderNumber string) (*Order, error)
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
}

type OrderItemRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters OrderItemFilters) ([]OrderItem, error)
	Count(ctx context.Context, orgID uuid.UUID, filters OrderItemFilters) (int64, error)
	Create(ctx context.Context, item *OrderItem) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderItem, error)
	Update(ctx context.Context, item *OrderItem) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
	ListByOrder(ctx context.Context, orgID uuid.UUID, orderID uuid.UUID) ([]OrderItem, error)
}

type OrderItemModifierRepository interface {
	List(ctx context.Context, orgID uuid.UUID, orderItemID uuid.UUID) ([]OrderItemModifier, error)
	Create(ctx context.Context, modifier *OrderItemModifier) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderItemModifier, error)
	Update(ctx context.Context, modifier *OrderItemModifier) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	DeleteByOrderItem(ctx context.Context, orgID uuid.UUID, orderItemID uuid.UUID) error
}

type KitchenTicketRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters KitchenTicketFilters) ([]KitchenTicket, error)
	Count(ctx context.Context, orgID uuid.UUID, filters KitchenTicketFilters) (int64, error)
	Create(ctx context.Context, ticket *KitchenTicket) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*KitchenTicket, error)
	GetByNumber(ctx context.Context, orgID uuid.UUID, ticketNumber string) (*KitchenTicket, error)
	Update(ctx context.Context, ticket *KitchenTicket) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
	ListByStation(ctx context.Context, orgID uuid.UUID, stationID uuid.UUID, statuses []string) ([]KitchenTicket, error)
}
