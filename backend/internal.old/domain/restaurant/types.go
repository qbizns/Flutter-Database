package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// FloorPlan represents a restaurant floor/dining area layout
type FloorPlan struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	LocationID     uuid.UUID  `json:"location_id"`
	FloorName      string     `json:"floor_name"`
	FloorLevel     int        `json:"floor_level"`
	DisplayOrder   int        `json:"display_order"`
	LayoutConfig   any        `json:"layout_config,omitempty"` // JSONB
	IsActive       bool       `json:"is_active"`
	IsDefault      bool       `json:"is_default"`
	Description    *string    `json:"description,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	Metadata       any        `json:"metadata,omitempty"` // JSONB
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// TableSection represents a section/zone within a floor plan
type TableSection struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	LocationID     uuid.UUID  `json:"location_id"`
	FloorPlanID    *uuid.UUID `json:"floor_plan_id,omitempty"`
	SectionName    string     `json:"section_name"`
	SectionCode    *string    `json:"section_code,omitempty"`
	SectionType    string     `json:"section_type"` // regular, vip, outdoor, bar, smoking, non_smoking, private, other
	ColorCode      *string    `json:"color_code,omitempty"`
	Icon           *string    `json:"icon,omitempty"`
	DisplayOrder   int        `json:"display_order"`
	IsActive       bool       `json:"is_active"`
	Description    *string    `json:"description,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	Metadata       any        `json:"metadata,omitempty"` // JSONB
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// RestaurantTable represents a physical table in the restaurant
type RestaurantTable struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	LocationID          uuid.UUID  `json:"location_id"`
	FloorPlanID         *uuid.UUID `json:"floor_plan_id,omitempty"`
	SectionID           *uuid.UUID `json:"section_id,omitempty"`
	TableNumber         string     `json:"table_number"`
	TableName           *string    `json:"table_name,omitempty"`
	MinCapacity         int        `json:"min_capacity"`
	MaxCapacity         int        `json:"max_capacity"`
	TableShape          string     `json:"table_shape"` // square, round, rectangle, oval, custom
	IsCombinable        bool       `json:"is_combinable"`
	PositionX           *float64   `json:"position_x,omitempty"`
	PositionY           *float64   `json:"position_y,omitempty"`
	Rotation            int        `json:"rotation"`
	Status              string     `json:"status"` // available, occupied, reserved, cleaning, maintenance, unavailable
	CurrentCovers       int        `json:"current_covers"`
	SeatedAt            *time.Time `json:"seated_at,omitempty"`
	CurrentWaiterID     *uuid.UUID `json:"current_waiter_id,omitempty"`
	IsActive            bool       `json:"is_active"`
	AllowOnlineReservation bool    `json:"allow_online_reservation"`
	DisplayOrder        int        `json:"display_order"`
	ColorCode           *string    `json:"color_code,omitempty"`
	Icon                *string    `json:"icon,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	Metadata            any        `json:"metadata,omitempty"` // JSONB
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// Reservation represents a table reservation
type Reservation struct {
	ID                      uuid.UUID  `json:"id"`
	OrganizationID          uuid.UUID  `json:"organization_id"`
	LocationID              uuid.UUID  `json:"location_id"`
	TableID                 *uuid.UUID `json:"table_id,omitempty"`
	CustomerID              *uuid.UUID `json:"customer_id,omitempty"`
	ReservationNumber       string     `json:"reservation_number"`
	ReservationDate         time.Time  `json:"reservation_date"`
	ReservationTime         string     `json:"reservation_time"` // TIME format
	DurationMinutes         int        `json:"duration_minutes"`
	PartySize               int        `json:"party_size"`
	CustomerName            string     `json:"customer_name"`
	CustomerPhone           *string    `json:"customer_phone,omitempty"`
	CustomerEmail           *string    `json:"customer_email,omitempty"`
	Status                  string     `json:"status"` // pending, confirmed, seated, completed, cancelled, no_show
	AssignedWaiterID        *uuid.UUID `json:"assigned_waiter_id,omitempty"`
	AssignedAt              *time.Time `json:"assigned_at,omitempty"`
	SeatedAt                *time.Time `json:"seated_at,omitempty"`
	CompletedAt             *time.Time `json:"completed_at,omitempty"`
	SpecialRequests         *string    `json:"special_requests,omitempty"`
	Occasion                *string    `json:"occasion,omitempty"`
	DietaryRestrictions     *string    `json:"dietary_restrictions,omitempty"`
	ConfirmationCode        *string    `json:"confirmation_code,omitempty"`
	ConfirmedAt             *time.Time `json:"confirmed_at,omitempty"`
	ConfirmedBy             *uuid.UUID `json:"confirmed_by,omitempty"`
	ReminderSentAt          *time.Time `json:"reminder_sent_at,omitempty"`
	NotificationPreferences any        `json:"notification_preferences,omitempty"` // JSONB
	CancelledAt             *time.Time `json:"cancelled_at,omitempty"`
	CancelledBy             *uuid.UUID `json:"cancelled_by,omitempty"`
	CancellationReason      *string    `json:"cancellation_reason,omitempty"`
	Notes                   *string    `json:"notes,omitempty"`
	Metadata                any        `json:"metadata,omitempty"` // JSONB
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	CreatedBy               *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy               *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt               *time.Time `json:"deleted_at,omitempty"`
}

// ModifierGroup represents a group of modifiers (Size, Extras, Toppings, etc.)
type ModifierGroup struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	GroupName      string     `json:"group_name"`
	GroupCode      *string    `json:"group_code,omitempty"`
	DisplayName    *string    `json:"display_name,omitempty"`
	SelectionType  string     `json:"selection_type"` // single, multiple, exact
	MinSelections  int        `json:"min_selections"`
	MaxSelections  *int       `json:"max_selections,omitempty"`
	ExactSelections *int      `json:"exact_selections,omitempty"`
	IsRequired     bool       `json:"is_required"`
	AffectsPrice   bool       `json:"affects_price"`
	DisplayOrder   int        `json:"display_order"`
	IsActive       bool       `json:"is_active"`
	Description    *string    `json:"description,omitempty"`
	Notes          *string    `json:"notes,omitempty"`
	Metadata       any        `json:"metadata,omitempty"` // JSONB
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Modifier represents an individual modifier within a group
type Modifier struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	ModifierGroupID  *uuid.UUID `json:"modifier_group_id,omitempty"`
	ModifierName     string     `json:"modifier_name"`
	ModifierCode     *string    `json:"modifier_code,omitempty"`
	DisplayName      *string    `json:"display_name,omitempty"`
	PriceAdjustment  float64    `json:"price_adjustment"`
	PriceType        string     `json:"price_type"` // add, multiply, replace
	IsAvailable      bool       `json:"is_available"`
	IsDefault        bool       `json:"is_default"`
	TrackInventory   bool       `json:"track_inventory"`
	CurrentStock     float64    `json:"current_stock"`
	LowStockThreshold float64   `json:"low_stock_threshold"`
	DisplayOrder     int        `json:"display_order"`
	ImageURL         *string    `json:"image_url,omitempty"`
	Description      *string    `json:"description,omitempty"`
	AllergenInfo     *string    `json:"allergen_info,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
	Metadata         any        `json:"metadata,omitempty"` // JSONB
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// ProductModifierGroup represents the link between products and modifier groups
type ProductModifierGroup struct {
	ID                      uuid.UUID  `json:"id"`
	OrganizationID          uuid.UUID  `json:"organization_id"`
	ProductID               uuid.UUID  `json:"product_id"`
	ModifierGroupID         uuid.UUID  `json:"modifier_group_id"`
	IsRequired              bool       `json:"is_required"`
	DisplayOrder            int        `json:"display_order"`
	IsActive                bool       `json:"is_active"`
	OverrideMinSelections   *int       `json:"override_min_selections,omitempty"`
	OverrideMaxSelections   *int       `json:"override_max_selections,omitempty"`
	Notes                   *string    `json:"notes,omitempty"`
	Metadata                any        `json:"metadata,omitempty"` // JSONB
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// Course represents a course in a meal sequence (appetizer, main, dessert, etc.)
type Course struct {
	ID                     uuid.UUID  `json:"id"`
	OrganizationID         uuid.UUID  `json:"organization_id"`
	CourseName             string     `json:"course_name"`
	CourseCode             *string    `json:"course_code,omitempty"`
	CourseType             string     `json:"course_type"` // appetizer, soup, salad, main, side, dessert, beverage, other
	TypicalDurationMinutes int        `json:"typical_duration_minutes"`
	FireDelayMinutes       int        `json:"fire_delay_minutes"`
	DisplayOrder           int        `json:"display_order"`
	ColorCode              *string    `json:"color_code,omitempty"`
	Icon                   *string    `json:"icon,omitempty"`
	IsActive               bool       `json:"is_active"`
	IsDefault              bool       `json:"is_default"`
	Description            *string    `json:"description,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
	Metadata               any        `json:"metadata,omitempty"` // JSONB
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	CreatedBy              *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy              *uuid.UUID `json:"updated_by,omitempty"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

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

// Filters for list operations
type FloorPlanFilters struct {
	Search       string
	LocationID   *uuid.UUID
	IsActive     *bool
	IsDefault    *bool
	FloorLevel   *int
	Page         int
	PageSize     int
}

type TableSectionFilters struct {
	Search       string
	FloorPlanID  *uuid.UUID
	SectionType  *string
	IsActive     *bool
	LocationID   *uuid.UUID
	Page         int
	PageSize     int
}

type RestaurantTableFilters struct {
	Search       string
	FloorPlanID  *uuid.UUID
	SectionID    *uuid.UUID
	Status       *string
	IsActive     *bool
	WaiterID     *uuid.UUID
	LocationID   *uuid.UUID
	Page         int
	PageSize     int
}

type ReservationFilters struct {
	Search           string
	Status           *string
	ReservationDate  *time.Time
	CustomerName     *string
	CustomerPhone    *string
	WaiterID         *uuid.UUID
	LocationID       *uuid.UUID
	Page             int
	PageSize         int
}

type ModifierGroupFilters struct {
	Search       string
	IsActive     *bool
	SelectionType *string
	Page         int
	PageSize     int
}

type ModifierFilters struct {
	Search           string
	ModifierGroupID  *uuid.UUID
	IsAvailable      *bool
	IsActive         *bool
	Page             int
	PageSize         int
}

type CourseFilters struct {
	Search     string
	CourseType *string
	IsActive   *bool
	IsDefault  *bool
	Page       int
	PageSize   int
}

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
type FloorPlanRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters FloorPlanFilters) ([]FloorPlan, error)
	Count(ctx context.Context, orgID uuid.UUID, filters FloorPlanFilters) (int64, error)
	Create(ctx context.Context, floorPlan *FloorPlan) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*FloorPlan, error)
	Update(ctx context.Context, floorPlan *FloorPlan) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

type TableSectionRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters TableSectionFilters) ([]TableSection, error)
	Count(ctx context.Context, orgID uuid.UUID, filters TableSectionFilters) (int64, error)
	Create(ctx context.Context, section *TableSection) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*TableSection, error)
	Update(ctx context.Context, section *TableSection) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

type RestaurantTableRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters RestaurantTableFilters) ([]RestaurantTable, error)
	Count(ctx context.Context, orgID uuid.UUID, filters RestaurantTableFilters) (int64, error)
	Create(ctx context.Context, table *RestaurantTable) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*RestaurantTable, error)
	GetByNumber(ctx context.Context, orgID uuid.UUID, locationID uuid.UUID, tableNumber string) (*RestaurantTable, error)
	Update(ctx context.Context, table *RestaurantTable) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error
}

type ReservationRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters ReservationFilters) ([]Reservation, error)
	Count(ctx context.Context, orgID uuid.UUID, filters ReservationFilters) (int64, error)
	Create(ctx context.Context, reservation *Reservation) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Reservation, error)
	GetByNumber(ctx context.Context, orgID uuid.UUID, reservationNumber string) (*Reservation, error)
	Update(ctx context.Context, reservation *Reservation) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

type ModifierGroupRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters ModifierGroupFilters) ([]ModifierGroup, error)
	Count(ctx context.Context, orgID uuid.UUID, filters ModifierGroupFilters) (int64, error)
	Create(ctx context.Context, group *ModifierGroup) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ModifierGroup, error)
	Update(ctx context.Context, group *ModifierGroup) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

type ModifierRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters ModifierFilters) ([]Modifier, error)
	Count(ctx context.Context, orgID uuid.UUID, filters ModifierFilters) (int64, error)
	Create(ctx context.Context, modifier *Modifier) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Modifier, error)
	Update(ctx context.Context, modifier *Modifier) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	ListByGroup(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) ([]Modifier, error)
}

type ProductModifierGroupRepository interface {
	List(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) ([]ProductModifierGroup, error)
	Create(ctx context.Context, pmg *ProductModifierGroup) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*ProductModifierGroup, error)
	Update(ctx context.Context, pmg *ProductModifierGroup) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	DeleteByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) error
}

type CourseRepository interface {
	List(ctx context.Context, orgID uuid.UUID, filters CourseFilters) ([]Course, error)
	Count(ctx context.Context, orgID uuid.UUID, filters CourseFilters) (int64, error)
	Create(ctx context.Context, course *Course) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Course, error)
	Update(ctx context.Context, course *Course) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}

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
