package staff

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// EmployeeSchedule represents an employee's work schedule
type EmployeeSchedule struct {
	ID                uuid.UUID  `json:"id"`
	OrganizationID    uuid.UUID  `json:"organization_id"`
	LocationID        *uuid.UUID `json:"location_id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	ScheduleDate      time.Time  `json:"schedule_date"`
	ShiftType         string     `json:"shift_type"`
	Position          string     `json:"position"`
	ScheduledStartTime string     `json:"scheduled_start_time"`
	ScheduledEndTime  string     `json:"scheduled_end_time"`
	BreakDurationMins int        `json:"break_duration_minutes"`
	Status            string     `json:"status"`
	RequiresApproval  bool       `json:"requires_approval"`
	ApprovedBy        *uuid.UUID `json:"approved_by"`
	ApprovedAt        *time.Time `json:"approved_at"`
	Notes             string     `json:"notes"`
	CancellationReason string    `json:"cancellation_reason"`
	Metadata          JSONB      `json:"metadata"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	CreatedBy         *uuid.UUID `json:"created_by"`
	UpdatedBy         *uuid.UUID `json:"updated_by"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// TimeClockEntry represents a time clock entry (clock in/out)
type TimeClockEntry struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	LocationID          *uuid.UUID `json:"location_id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	ScheduleID          *uuid.UUID `json:"schedule_id"`
	EntryType           string     `json:"entry_type"`
	EntryTimestamp      time.Time  `json:"entry_timestamp"`
	ScheduledTimestamp  *time.Time `json:"scheduled_timestamp"`
	DeviceID            *uuid.UUID `json:"device_id"`
	GPSLocation         JSONB      `json:"gps_location"`
	IPAddress           string     `json:"ip_address"`
	IsLate              bool       `json:"is_late"`
	IsEarly             bool       `json:"is_early"`
	VarianceMinutes     *int       `json:"variance_minutes"`
	RequiresApproval    bool       `json:"requires_approval"`
	ApprovedBy          *uuid.UUID `json:"approved_by"`
	ApprovedAt          *time.Time `json:"approved_at"`
	IsManualEntry       bool       `json:"is_manual_entry"`
	CorrectionNotes     string     `json:"correction_notes"`
	PhotoURL            string     `json:"photo_url"`
	Notes               string     `json:"notes"`
	Metadata            JSONB      `json:"metadata"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CreatedBy           *uuid.UUID `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// Device represents a POS device or terminal
type Device struct {
	ID                   uuid.UUID  `json:"id"`
	OrganizationID       uuid.UUID  `json:"organization_id"`
	LocationID           *uuid.UUID `json:"location_id"`
	DeviceCode           string     `json:"device_code"`
	DeviceName           string     `json:"device_name"`
	DeviceType           string     `json:"device_type"`
	Manufacturer         string     `json:"manufacturer"`
	Model                string     `json:"model"`
	SerialNumber         string     `json:"serial_number"`
	MACAddress           string     `json:"mac_address"`
	IPAddress            string     `json:"ip_address"`
	DeviceConfig         JSONB      `json:"device_config"`
	ScreenResolution     string     `json:"screen_resolution"`
	OSVersion            string     `json:"os_version"`
	ConnectionType       string     `json:"connection_type"`
	ConnectionString     string     `json:"connection_string"`
	Status               string     `json:"status"`
	LastOnlineAt         *time.Time `json:"last_online_at"`
	LastHeartbeatAt      *time.Time `json:"last_heartbeat_at"`
	AssignedToUserID     *uuid.UUID `json:"assigned_to_user_id"`
	AssignedToStationID  *uuid.UUID `json:"assigned_to_station_id"`
	PurchaseDate         *time.Time `json:"purchase_date"`
	WarrantyExpiryDate   *time.Time `json:"warranty_expiry_date"`
	LicenseKey           string     `json:"license_key"`
	LicenseExpiryDate    *time.Time `json:"license_expiry_date"`
	InstallationNotes    string     `json:"installation_notes"`
	MaintenanceNotes     string     `json:"maintenance_notes"`
	Metadata             JSONB      `json:"metadata"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	CreatedBy            *uuid.UUID `json:"created_by"`
	UpdatedBy            *uuid.UUID `json:"updated_by"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

// PrinterConfiguration represents printer routing rules
type PrinterConfiguration struct {
	ID                    uuid.UUID  `json:"id"`
	OrganizationID        uuid.UUID  `json:"organization_id"`
	LocationID            *uuid.UUID `json:"location_id"`
	PrinterDeviceID       uuid.UUID  `json:"printer_device_id"`
	DocumentType          string     `json:"document_type"`
	FilterOrderType       *string    `json:"filter_order_type"`
	FilterKitchenStationID *uuid.UUID `json:"filter_kitchen_station_id"`
	FilterProductCategoryID *uuid.UUID `json:"filter_product_category_id"`
	FilterCourseID        *uuid.UUID `json:"filter_course_id"`
	NumberOfCopies        int        `json:"number_of_copies"`
	AutoPrint             bool       `json:"auto_print"`
	PrintPriority         int        `json:"print_priority"`
	TemplateConfig        JSONB      `json:"template_config"`
	PaperSize             string     `json:"paper_size"`
	PrintOrientation      string     `json:"print_orientation"`
	IsActive              bool       `json:"is_active"`
	Notes                 string     `json:"notes"`
	Metadata              JSONB      `json:"metadata"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             *uuid.UUID `json:"created_by"`
	UpdatedBy             *uuid.UUID `json:"updated_by"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}

// TipPool represents a tip pooling configuration
type TipPool struct {
	ID                   uuid.UUID `json:"id"`
	OrganizationID       uuid.UUID `json:"organization_id"`
	LocationID           *uuid.UUID `json:"location_id"`
	PoolName             string    `json:"pool_name"`
	PoolType             string    `json:"pool_type"`
	Description          string    `json:"description"`
	DistributionMethod   string    `json:"distribution_method"`
	DistributionConfig   JSONB     `json:"distribution_config"`
	EligiblePositions    pq.StringArray `json:"eligible_positions"`
	IsActive             bool      `json:"is_active"`
	Metadata             JSONB     `json:"metadata"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedBy            *uuid.UUID `json:"created_by"`
	UpdatedBy            *uuid.UUID `json:"updated_by"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

// TipDistribution represents a tip distribution record
type TipDistribution struct {
	ID                    uuid.UUID  `json:"id"`
	OrganizationID        uuid.UUID  `json:"organization_id"`
	LocationID            *uuid.UUID `json:"location_id"`
	TipPoolID             *uuid.UUID `json:"tip_pool_id"`
	DistributionDate      time.Time  `json:"distribution_date"`
	PeriodStart           *time.Time `json:"period_start"`
	PeriodEnd             *time.Time `json:"period_end"`
	ShiftID               *uuid.UUID `json:"shift_id"`
	EmployeeID            uuid.UUID  `json:"employee_id"`
	SourceType            string     `json:"source_type"`
	SourceSaleID          *uuid.UUID `json:"source_sale_id"`
	SourceOrderID         *uuid.UUID `json:"source_order_id"`
	TipAmount             float64    `json:"tip_amount"`
	DistributionAmount    float64    `json:"distribution_amount"`
	DistributionPercentage *float64  `json:"distribution_percentage"`
	PaymentStatus         string     `json:"payment_status"`
	PaymentMethod         *string    `json:"payment_method"`
	PaidAt                *time.Time `json:"paid_at"`
	PaidBy                *uuid.UUID `json:"paid_by"`
	Notes                 string     `json:"notes"`
	Metadata              JSONB      `json:"metadata"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	CreatedBy             *uuid.UUID `json:"created_by"`
	UpdatedBy             *uuid.UUID `json:"updated_by"`
	DeletedAt             *time.Time `json:"deleted_at,omitempty"`
}

// StaffCommission represents a staff commission record
type StaffCommission struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	LocationID          *uuid.UUID `json:"location_id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	CommissionDate      time.Time  `json:"commission_date"`
	PeriodStart         time.Time  `json:"period_start"`
	PeriodEnd           time.Time  `json:"period_end"`
	SourceType          string     `json:"source_type"`
	SourceSaleID        *uuid.UUID `json:"source_sale_id"`
	SourceOrderID       *uuid.UUID `json:"source_order_id"`
	CommissionType      string     `json:"commission_type"`
	CommissionRate      *float64   `json:"commission_rate"`
	SalesAmount         float64    `json:"sales_amount"`
	CommissionAmount    float64    `json:"commission_amount"`
	Status              string     `json:"status"`
	ApprovedBy          *uuid.UUID `json:"approved_by"`
	ApprovedAt          *time.Time `json:"approved_at"`
	PaymentDate         *time.Time `json:"payment_date"`
	PaymentMethod       *string    `json:"payment_method"`
	PaidBy              *uuid.UUID `json:"paid_by"`
	Notes               string     `json:"notes"`
	CalculationNotes    string     `json:"calculation_notes"`
	Metadata            JSONB      `json:"metadata"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CreatedBy           *uuid.UUID `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
}

// Shift represents a cashier shift
type Shift struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	LocationID         *uuid.UUID `json:"location_id"`
	UserID             uuid.UUID  `json:"user_id"`
	ShiftNumber        string     `json:"shift_number"`
	StartTime          time.Time  `json:"start_time"`
	EndTime            *time.Time `json:"end_time"`
	Status             string     `json:"status"`
	OpeningCash        float64    `json:"opening_cash"`
	OpeningNotes       string     `json:"opening_notes"`
	ExpectedCash       *float64   `json:"expected_cash"`
	ActualCash         *float64   `json:"actual_cash"`
	CashDifference     *float64   `json:"cash_difference"`
	ClosingNotes       string     `json:"closing_notes"`
	TotalSales         float64    `json:"total_sales"`
	TotalTransactions  int        `json:"total_transactions"`
	TotalRefunds       float64    `json:"total_refunds"`
	TotalDiscounts     float64    `json:"total_discounts"`
	PaymentBreakdown   JSONB      `json:"payment_breakdown"`
	Notes              string     `json:"notes"`
	Metadata           JSONB      `json:"metadata"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ClosedBy           *uuid.UUID `json:"closed_by"`
	ClosedAt           *time.Time `json:"closed_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// Expense represents a business expense
type Expense struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	LocationID       *uuid.UUID `json:"location_id"`
	ExpenseNumber    string     `json:"expense_number"`
	ExpenseDate      time.Time  `json:"expense_date"`
	Category         string     `json:"category"`
	Subcategory      string     `json:"subcategory"`
	PayeeName        string     `json:"payee_name"`
	PaymentMethod    *string    `json:"payment_method"`
	Amount           float64    `json:"amount"`
	TaxAmount        float64    `json:"tax_amount"`
	TotalAmount      float64    `json:"total_amount"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"`
	ReferenceNumber  string     `json:"reference_number"`
	PurchaseOrderID  *uuid.UUID `json:"purchase_order_id"`
	ReceiptURL       string     `json:"receipt_url"`
	AttachmentURLs   JSONB      `json:"attachment_urls"`
	Description      string     `json:"description"`
	Notes            string     `json:"notes"`
	Metadata         JSONB      `json:"metadata"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	CreatedBy        *uuid.UUID `json:"created_by"`
	UpdatedBy        *uuid.UUID `json:"updated_by"`
	ApprovedBy       *uuid.UUID `json:"approved_by"`
	ApprovedAt       *time.Time `json:"approved_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// JSONB represents a PostgreSQL JSONB value
type JSONB map[string]interface{}

// Value implements the database/sql/driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the database/sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = JSONB(result)
	return err
}

// Filters

// EmployeeScheduleFilters represents filters for listing employee schedules
type EmployeeScheduleFilters struct {
	EmployeeID *uuid.UUID
	LocationID *uuid.UUID
	Status     *string
	ShiftType  *string
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PageSize   int
}

// TimeClockEntryFilters represents filters for listing time clock entries
type TimeClockEntryFilters struct {
	EmployeeID   *uuid.UUID
	LocationID   *uuid.UUID
	EntryType    *string
	DateFrom     *time.Time
	DateTo       *time.Time
	IsApproved   *bool
	Page         int
	PageSize     int
}

// DeviceFilters represents filters for listing devices
type DeviceFilters struct {
	DeviceType  *string
	Status      *string
	LocationID  *uuid.UUID
	Search      *string
	Page        int
	PageSize    int
}

// PrinterConfigurationFilters represents filters for listing printer configurations
type PrinterConfigurationFilters struct {
	DocumentType *string
	IsActive     *bool
	LocationID   *uuid.UUID
	Page         int
	PageSize     int
}

// TipPoolFilters represents filters for listing tip pools
type TipPoolFilters struct {
	IsActive   *bool
	LocationID *uuid.UUID
	Page       int
	PageSize   int
}

// TipDistributionFilters represents filters for listing tip distributions
type TipDistributionFilters struct {
	EmployeeID      *uuid.UUID
	PaymentStatus   *string
	LocationID      *uuid.UUID
	DistributionDateFrom *time.Time
	DistributionDateTo   *time.Time
	Page            int
	PageSize        int
}

// StaffCommissionFilters represents filters for listing staff commissions
type StaffCommissionFilters struct {
	EmployeeID    *uuid.UUID
	Status        *string
	LocationID    *uuid.UUID
	CommissionDateFrom *time.Time
	CommissionDateTo   *time.Time
	Page          int
	PageSize      int
}

// ShiftFilters represents filters for listing shifts
type ShiftFilters struct {
	UserID     *uuid.UUID
	Status     *string
	LocationID *uuid.UUID
	DateFrom   *time.Time
	DateTo     *time.Time
	Page       int
	PageSize   int
}

// ExpenseFilters represents filters for listing expenses
type ExpenseFilters struct {
	Category    *string
	Status      *string
	LocationID  *uuid.UUID
	ExpenseDateFrom *time.Time
	ExpenseDateTo   *time.Time
	Page        int
	PageSize    int
}

// Repository defines the staff data access interface
type Repository interface {
	// EmployeeSchedule CRUD
	CreateSchedule(ctx context.Context, schedule *EmployeeSchedule) error
	GetSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) (*EmployeeSchedule, error)
	ListSchedules(ctx context.Context, orgID uuid.UUID, filters EmployeeScheduleFilters) ([]EmployeeSchedule, error)
	CountSchedules(ctx context.Context, orgID uuid.UUID, filters EmployeeScheduleFilters) (int64, error)
	UpdateSchedule(ctx context.Context, schedule *EmployeeSchedule) error
	DeleteSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) error

	// TimeClockEntry CRUD
	CreateClockEntry(ctx context.Context, entry *TimeClockEntry) error
	GetClockEntry(ctx context.Context, orgID, entryID uuid.UUID) (*TimeClockEntry, error)
	ListClockEntries(ctx context.Context, orgID uuid.UUID, filters TimeClockEntryFilters) ([]TimeClockEntry, error)
	CountClockEntries(ctx context.Context, orgID uuid.UUID, filters TimeClockEntryFilters) (int64, error)
	UpdateClockEntry(ctx context.Context, entry *TimeClockEntry) error
	DeleteClockEntry(ctx context.Context, orgID, entryID uuid.UUID) error

	// Device CRUD
	CreateDevice(ctx context.Context, device *Device) error
	GetDevice(ctx context.Context, orgID, deviceID uuid.UUID) (*Device, error)
	ListDevices(ctx context.Context, orgID uuid.UUID, filters DeviceFilters) ([]Device, error)
	CountDevices(ctx context.Context, orgID uuid.UUID, filters DeviceFilters) (int64, error)
	UpdateDevice(ctx context.Context, device *Device) error
	DeleteDevice(ctx context.Context, orgID, deviceID uuid.UUID) error
	UpdateDeviceHeartbeat(ctx context.Context, deviceID uuid.UUID) error

	// PrinterConfiguration CRUD
	CreatePrinterConfig(ctx context.Context, config *PrinterConfiguration) error
	GetPrinterConfig(ctx context.Context, orgID, configID uuid.UUID) (*PrinterConfiguration, error)
	ListPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters PrinterConfigurationFilters) ([]PrinterConfiguration, error)
	CountPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters PrinterConfigurationFilters) (int64, error)
	UpdatePrinterConfig(ctx context.Context, config *PrinterConfiguration) error
	DeletePrinterConfig(ctx context.Context, orgID, configID uuid.UUID) error

	// TipPool CRUD
	CreateTipPool(ctx context.Context, pool *TipPool) error
	GetTipPool(ctx context.Context, orgID, poolID uuid.UUID) (*TipPool, error)
	ListTipPools(ctx context.Context, orgID uuid.UUID, filters TipPoolFilters) ([]TipPool, error)
	CountTipPools(ctx context.Context, orgID uuid.UUID, filters TipPoolFilters) (int64, error)
	UpdateTipPool(ctx context.Context, pool *TipPool) error
	DeleteTipPool(ctx context.Context, orgID, poolID uuid.UUID) error

	// TipDistribution CRUD
	CreateTipDistribution(ctx context.Context, dist *TipDistribution) error
	GetTipDistribution(ctx context.Context, orgID, distID uuid.UUID) (*TipDistribution, error)
	ListTipDistributions(ctx context.Context, orgID uuid.UUID, filters TipDistributionFilters) ([]TipDistribution, error)
	CountTipDistributions(ctx context.Context, orgID uuid.UUID, filters TipDistributionFilters) (int64, error)
	UpdateTipDistribution(ctx context.Context, dist *TipDistribution) error
	DeleteTipDistribution(ctx context.Context, orgID, distID uuid.UUID) error

	// StaffCommission CRUD
	CreateCommission(ctx context.Context, comm *StaffCommission) error
	GetCommission(ctx context.Context, orgID, commID uuid.UUID) (*StaffCommission, error)
	ListCommissions(ctx context.Context, orgID uuid.UUID, filters StaffCommissionFilters) ([]StaffCommission, error)
	CountCommissions(ctx context.Context, orgID uuid.UUID, filters StaffCommissionFilters) (int64, error)
	UpdateCommission(ctx context.Context, comm *StaffCommission) error
	DeleteCommission(ctx context.Context, orgID, commID uuid.UUID) error

	// Shift CRUD
	CreateShift(ctx context.Context, shift *Shift) error
	GetShift(ctx context.Context, orgID, shiftID uuid.UUID) (*Shift, error)
	ListShifts(ctx context.Context, orgID uuid.UUID, filters ShiftFilters) ([]Shift, error)
	CountShifts(ctx context.Context, orgID uuid.UUID, filters ShiftFilters) (int64, error)
	UpdateShift(ctx context.Context, shift *Shift) error
	DeleteShift(ctx context.Context, orgID, shiftID uuid.UUID) error

	// Expense CRUD
	CreateExpense(ctx context.Context, expense *Expense) error
	GetExpense(ctx context.Context, orgID, expenseID uuid.UUID) (*Expense, error)
	ListExpenses(ctx context.Context, orgID uuid.UUID, filters ExpenseFilters) ([]Expense, error)
	CountExpenses(ctx context.Context, orgID uuid.UUID, filters ExpenseFilters) (int64, error)
	UpdateExpense(ctx context.Context, expense *Expense) error
	DeleteExpense(ctx context.Context, orgID, expenseID uuid.UUID) error
}
