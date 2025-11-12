package delivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeliveryZone represents a geographic area for delivery with fees and timing
type DeliveryZone struct {
	ID                        uuid.UUID       `json:"id"`
	OrganizationID            uuid.UUID       `json:"organization_id"`
	LocationID                *uuid.UUID      `json:"location_id"`
	ZoneName                  string          `json:"zone_name"`
	ZoneCode                  string          `json:"zone_code"`
	Description               string          `json:"description"`
	Geofence                  json.RawMessage `json:"geofence"`
	PostalCodes               []string        `json:"postal_codes"`
	CoverageNotes             string          `json:"coverage_notes"`
	BaseDeliveryFee           float64         `json:"base_delivery_fee"`
	FeeType                   string          `json:"fee_type"`
	MinimumOrderAmount        float64         `json:"minimum_order_amount"`
	FreeDeliveryThreshold     *float64        `json:"free_delivery_threshold"`
	EstimatedDeliveryMinutes  int             `json:"estimated_delivery_time_minutes"`
	MaxDeliveryTimeMinutes    int             `json:"max_delivery_time_minutes"`
	Priority                  int             `json:"priority"`
	IsActive                  bool            `json:"is_active"`
	ActiveHours               json.RawMessage `json:"active_hours"`
	Metadata                  json.RawMessage `json:"metadata"`
	CreatedAt                 time.Time       `json:"created_at"`
	UpdatedAt                 time.Time       `json:"updated_at"`
	CreatedBy                 *uuid.UUID      `json:"created_by"`
	UpdatedBy                 *uuid.UUID      `json:"updated_by"`
	DeletedAt                 *time.Time      `json:"deleted_at,omitempty"`
}

// DeliveryDriver represents a delivery driver
type DeliveryDriver struct {
	ID                      uuid.UUID       `json:"id"`
	OrganizationID          uuid.UUID       `json:"organization_id"`
	UserID                  *uuid.UUID      `json:"user_id"`
	DriverCode              string          `json:"driver_code"`
	FullName                string          `json:"full_name"`
	Phone                   string          `json:"phone"`
	Email                   string          `json:"email"`
	EmergencyContactName    string          `json:"emergency_contact_name"`
	EmergencyContactPhone   string          `json:"emergency_contact_phone"`
	VehicleType             string          `json:"vehicle_type"`
	VehicleMake             string          `json:"vehicle_make"`
	VehicleModel            string          `json:"vehicle_model"`
	VehicleYear             *int            `json:"vehicle_year"`
	VehicleColor            string          `json:"vehicle_color"`
	LicensePlate            string          `json:"license_plate"`
	DriversLicenseNumber    string          `json:"drivers_license_number"`
	LicenseExpiryDate       *time.Time      `json:"license_expiry_date"`
	InsurancePolicyNumber   string          `json:"insurance_policy_number"`
	InsuranceExpiryDate     *time.Time      `json:"insurance_expiry_date"`
	HireDate                *time.Time      `json:"hire_date"`
	EmploymentType          string          `json:"employment_type"`
	Status                  string          `json:"status"`
	TotalDeliveries         int             `json:"total_deliveries"`
	SuccessfulDeliveries    int             `json:"successful_deliveries"`
	Rating                  *float64        `json:"rating"`
	RatingCount             int             `json:"rating_count"`
	CurrentLocation         json.RawMessage `json:"current_location"`
	IsAvailable             bool            `json:"is_available"`
	LastLocationUpdate      *time.Time      `json:"last_location_update"`
	CommissionRate          *float64        `json:"commission_rate"`
	PaymentMethod           string          `json:"payment_method"`
	Documents               json.RawMessage `json:"documents"`
	Metadata                json.RawMessage `json:"metadata"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
	CreatedBy               *uuid.UUID      `json:"created_by"`
	UpdatedBy               *uuid.UUID      `json:"updated_by"`
	DeletedAt               *time.Time      `json:"deleted_at,omitempty"`
}

// DriverShift represents a driver work schedule
type DriverShift struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	LocationID         *uuid.UUID `json:"location_id"`
	DriverID           uuid.UUID  `json:"driver_id"`
	ShiftDate          time.Time  `json:"shift_date"`
	ScheduledStartTime *time.Time `json:"scheduled_start_time"`
	ScheduledEndTime   *time.Time `json:"scheduled_end_time"`
	ActualStartTime    *time.Time `json:"actual_start_time"`
	ActualEndTime      *time.Time `json:"actual_end_time"`
	Status             string     `json:"status"`
	TotalBreakMinutes  int        `json:"total_break_minutes"`
	TotalDeliveries    int        `json:"total_deliveries"`
	TotalDistanceKm    float64    `json:"total_distance_km"`
	TotalEarnings      float64    `json:"total_earnings"`
	Notes              string     `json:"notes"`
	Metadata           json.RawMessage `json:"metadata"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          *uuid.UUID `json:"created_by"`
	UpdatedBy          *uuid.UUID `json:"updated_by"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// CustomerAddress represents a customer delivery address
type CustomerAddress struct {
	ID              uuid.UUID       `json:"id"`
	OrganizationID  uuid.UUID       `json:"organization_id"`
	CustomerID      uuid.UUID       `json:"customer_id"`
	AddressLabel    string          `json:"address_label"`
	AddressLine1    string          `json:"address_line_1"`
	AddressLine2    string          `json:"address_line_2"`
	City            string          `json:"city"`
	StateProvince   string          `json:"state_province"`
	PostalCode      string          `json:"postal_code"`
	Country         string          `json:"country"`
	Latitude        *float64        `json:"latitude"`
	Longitude       *float64        `json:"longitude"`
	LocationNotes   string          `json:"location_notes"`
	DeliveryZoneID  *uuid.UUID      `json:"delivery_zone_id"`
	IsDefault       bool            `json:"is_default"`
	IsActive        bool            `json:"is_active"`
	Metadata        json.RawMessage `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	CreatedBy       *uuid.UUID      `json:"created_by"`
	UpdatedBy       *uuid.UUID      `json:"updated_by"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
}

// DeliveryAssignment represents the assignment of a delivery to a driver
type DeliveryAssignment struct {
	ID                      uuid.UUID       `json:"id"`
	OrganizationID          uuid.UUID       `json:"organization_id"`
	OrderID                 uuid.UUID       `json:"order_id"`
	DriverID                uuid.UUID       `json:"driver_id"`
	DriverShiftID           *uuid.UUID      `json:"driver_shift_id"`
	DeliveryZoneID          *uuid.UUID      `json:"delivery_zone_id"`
	CustomerAddressID       *uuid.UUID      `json:"customer_address_id"`
	DeliveryAddress         string          `json:"delivery_address"`
	DeliveryLocation        json.RawMessage `json:"delivery_location"`
	AssignedAt              time.Time       `json:"assigned_at"`
	AssignedBy              *uuid.UUID      `json:"assigned_by"`
	Status                  string          `json:"status"`
	AcceptedAt              *time.Time      `json:"accepted_at"`
	PickedUpAt              *time.Time      `json:"picked_up_at"`
	DispatchedAt            *time.Time      `json:"dispatched_at"`
	ArrivedAt               *time.Time      `json:"arrived_at"`
	DeliveredAt             *time.Time      `json:"delivered_at"`
	FailedAt                *time.Time      `json:"failed_at"`
	EstimatedPickupTime     *time.Time      `json:"estimated_pickup_time"`
	EstimatedDeliveryTime   *time.Time      `json:"estimated_delivery_time"`
	DistanceKm              *float64        `json:"distance_km"`
	RouteInfo               json.RawMessage `json:"route_info"`
	DeliveryFee             float64         `json:"delivery_fee"`
	DriverCommission        float64         `json:"driver_commission"`
	PaymentMethod           string          `json:"payment_method"`
	CashCollected           float64         `json:"cash_collected"`
	SignatureImageURL       string          `json:"signature_image_url"`
	DeliveryPhotoURL        string          `json:"delivery_photo_url"`
	RecipientName           string          `json:"recipient_name"`
	DeliveryNotes           string          `json:"delivery_notes"`
	FailureReason           string          `json:"failure_reason"`
	FailureNotes            string          `json:"failure_notes"`
	RetryCount              int             `json:"retry_count"`
	CustomerRating          *int            `json:"customer_rating"`
	CustomerFeedback        string          `json:"customer_feedback"`
	DriverNotes             string          `json:"driver_notes"`
	Metadata                json.RawMessage `json:"metadata"`
	CreatedAt               time.Time       `json:"created_at"`
	UpdatedAt               time.Time       `json:"updated_at"`
	CreatedBy               *uuid.UUID      `json:"created_by"`
	UpdatedBy               *uuid.UUID      `json:"updated_by"`
	DeletedAt               *time.Time      `json:"deleted_at,omitempty"`
}

// OrderTrackingEvent represents a real-time tracking event
type OrderTrackingEvent struct {
	ID                    uuid.UUID       `json:"id"`
	OrganizationID        uuid.UUID       `json:"organization_id"`
	OrderID               uuid.UUID       `json:"order_id"`
	DeliveryAssignmentID  *uuid.UUID      `json:"delivery_assignment_id"`
	EventType             string          `json:"event_type"`
	EventTimestamp        time.Time       `json:"event_timestamp"`
	EventMessage          string          `json:"event_message"`
	Location              json.RawMessage `json:"location"`
	LocationName          string          `json:"location_name"`
	ActorType             string          `json:"actor_type"`
	ActorID               *uuid.UUID      `json:"actor_id"`
	ActorName             string          `json:"actor_name"`
	Metadata              json.RawMessage `json:"metadata"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedBy             *uuid.UUID      `json:"created_by"`
}

// Filters for queries
type DeliveryZoneFilters struct {
	Search   string
	IsActive *bool
	Page     int
	PageSize int
}

type DeliveryDriverFilters struct {
	Search      string
	Status      *string
	IsAvailable *bool
	Page        int
	PageSize    int
}

type DriverShiftFilters struct {
	Status   *string
	FromDate *time.Time
	ToDate   *time.Time
	DriverID *uuid.UUID
	Page     int
	PageSize int
}

type CustomerAddressFilters struct {
	Search     string
	CustomerID *uuid.UUID
	IsActive   *bool
	Page       int
	PageSize   int
}

type DeliveryAssignmentFilters struct {
	Status     *string
	DriverID   *uuid.UUID
	OrderID    *uuid.UUID
	FromDate   *time.Time
	ToDate     *time.Time
	Page       int
	PageSize   int
}

type OrderTrackingEventFilters struct {
	EventType *string
	FromDate  *time.Time
	ToDate    *time.Time
	Page      int
	PageSize  int
}

// Repository defines the delivery data access interface
type Repository interface {
	// Delivery Zones
	ListDeliveryZones(ctx context.Context, orgID uuid.UUID, filters DeliveryZoneFilters) ([]DeliveryZone, error)
	CountDeliveryZones(ctx context.Context, orgID uuid.UUID, filters DeliveryZoneFilters) (int64, error)
	CreateDeliveryZone(ctx context.Context, zone *DeliveryZone) error
	GetDeliveryZone(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DeliveryZone, error)
	UpdateDeliveryZone(ctx context.Context, zone *DeliveryZone) error
	DeleteDeliveryZone(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Delivery Drivers
	ListDeliveryDrivers(ctx context.Context, orgID uuid.UUID, filters DeliveryDriverFilters) ([]DeliveryDriver, error)
	CountDeliveryDrivers(ctx context.Context, orgID uuid.UUID, filters DeliveryDriverFilters) (int64, error)
	CreateDeliveryDriver(ctx context.Context, driver *DeliveryDriver) error
	GetDeliveryDriver(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DeliveryDriver, error)
	GetDeliveryDriverByCode(ctx context.Context, orgID uuid.UUID, code string) (*DeliveryDriver, error)
	UpdateDeliveryDriver(ctx context.Context, driver *DeliveryDriver) error
	DeleteDeliveryDriver(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateDriverAvailability(ctx context.Context, orgID uuid.UUID, driverID uuid.UUID, available bool, location json.RawMessage) error
	UpdateDriverRating(ctx context.Context, orgID uuid.UUID, driverID uuid.UUID, rating float64) error

	// Driver Shifts
	ListDriverShifts(ctx context.Context, orgID uuid.UUID, filters DriverShiftFilters) ([]DriverShift, error)
	CountDriverShifts(ctx context.Context, orgID uuid.UUID, filters DriverShiftFilters) (int64, error)
	CreateDriverShift(ctx context.Context, shift *DriverShift) error
	GetDriverShift(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DriverShift, error)
	UpdateDriverShift(ctx context.Context, shift *DriverShift) error
	DeleteDriverShift(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Customer Addresses
	ListCustomerAddresses(ctx context.Context, orgID uuid.UUID, filters CustomerAddressFilters) ([]CustomerAddress, error)
	CountCustomerAddresses(ctx context.Context, orgID uuid.UUID, filters CustomerAddressFilters) (int64, error)
	CreateCustomerAddress(ctx context.Context, address *CustomerAddress) error
	GetCustomerAddress(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomerAddress, error)
	UpdateCustomerAddress(ctx context.Context, address *CustomerAddress) error
	DeleteCustomerAddress(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	SetDefaultCustomerAddress(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, addressID uuid.UUID) error

	// Delivery Assignments
	ListDeliveryAssignments(ctx context.Context, orgID uuid.UUID, filters DeliveryAssignmentFilters) ([]DeliveryAssignment, error)
	CountDeliveryAssignments(ctx context.Context, orgID uuid.UUID, filters DeliveryAssignmentFilters) (int64, error)
	CreateDeliveryAssignment(ctx context.Context, assignment *DeliveryAssignment) error
	GetDeliveryAssignment(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DeliveryAssignment, error)
	GetAssignmentByOrderID(ctx context.Context, orgID uuid.UUID, orderID uuid.UUID) (*DeliveryAssignment, error)
	UpdateDeliveryAssignment(ctx context.Context, assignment *DeliveryAssignment) error
	DeleteDeliveryAssignment(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateAssignmentStatus(ctx context.Context, orgID uuid.UUID, assignmentID uuid.UUID, status string) error

	// Order Tracking Events
	ListOrderTrackingEvents(ctx context.Context, orgID uuid.UUID, filters OrderTrackingEventFilters) ([]OrderTrackingEvent, error)
	CountOrderTrackingEvents(ctx context.Context, orgID uuid.UUID, filters OrderTrackingEventFilters) (int64, error)
	CreateOrderTrackingEvent(ctx context.Context, event *OrderTrackingEvent) error
	GetOrderTrackingEvent(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*OrderTrackingEvent, error)
}
