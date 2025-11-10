package finance

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// ========================
// CURRENCIES
// ========================

type Currency struct {
	ID             uuid.UUID `db:"id" json:"id"`
	CurrencyCode   string    `db:"currency_code" json:"currency_code"`
	CurrencyName   string    `db:"currency_name" json:"currency_name"`
	CurrencySymbol *string   `db:"currency_symbol" json:"currency_symbol,omitempty"`
	DecimalPlaces  int       `db:"decimal_places" json:"decimal_places"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateCurrencyRequest struct {
	CurrencyCode   string  `json:"currency_code" validate:"required,len=3"`
	CurrencyName   string  `json:"currency_name" validate:"required,min=3,max=100"`
	CurrencySymbol *string `json:"currency_symbol"`
	DecimalPlaces  int     `json:"decimal_places" validate:"min=0,max=8"`
	IsActive       bool    `json:"is_active"`
}

type UpdateCurrencyRequest struct {
	CurrencyName   *string `json:"currency_name"`
	CurrencySymbol *string `json:"currency_symbol"`
	DecimalPlaces  *int    `json:"decimal_places"`
	IsActive       *bool   `json:"is_active"`
}

// ========================
// CURRENCY RATES
// ========================

type CurrencyRate struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	CurrencyCode   string    `db:"currency_code" json:"currency_code"`
	RateDate       time.Time `db:"rate_date" json:"rate_date"`
	Rate           string    `db:"rate" json:"rate"`
	Source         string    `db:"source" json:"source"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateCurrencyRateRequest struct {
	CurrencyCode string    `json:"currency_code" validate:"required,len=3"`
	RateDate     time.Time `json:"rate_date" validate:"required"`
	Rate         string    `json:"rate" validate:"required"`
	Source       string    `json:"source" validate:"required,min=1,max=50"`
}

type UpdateCurrencyRateRequest struct {
	Rate   *string `json:"rate"`
	Source *string `json:"source"`
}

type CurrencyConversionRequest struct {
	Amount           string    `json:"amount" validate:"required"`
	FromCurrency     string    `json:"from_currency" validate:"required,len=3"`
	ToCurrency       string    `json:"to_currency" validate:"required,len=3"`
	ConversionDate   time.Time `json:"conversion_date" validate:"required"`
}

type CurrencyConversionResult struct {
	Amount           string `json:"amount"`
	FromCurrency     string `json:"from_currency"`
	ToCurrency       string `json:"to_currency"`
	ConversionDate   time.Time `json:"conversion_date"`
	FromRate         string `json:"from_rate"`
	ToRate           string `json:"to_rate"`
	ConvertedAmount  string `json:"converted_amount"`
}

// ========================
// PAYMENT TERMS
// ========================

type PaymentTermValueType string

const (
	PaymentTermValueTypePercentage PaymentTermValueType = "percentage"
	PaymentTermValueTypeFixed      PaymentTermValueType = "fixed"
	PaymentTermValueTypeBalance    PaymentTermValueType = "balance"
)

func (t PaymentTermValueType) String() string {
	return string(t)
}

type PaymentTerm struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	TermCode       string    `db:"term_code" json:"term_code"`
	TermName       string    `db:"term_name" json:"term_name"`
	Note           *string   `db:"note" json:"note,omitempty"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedBy      *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
	UpdatedBy      *uuid.UUID `db:"updated_by" json:"updated_by,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`

	// Populated from child table
	Lines []*PaymentTermLine `db:"-" json:"lines,omitempty"`
}

type CreatePaymentTermRequest struct {
	TermCode string                      `json:"term_code" validate:"required,min=1,max=20"`
	TermName string                      `json:"term_name" validate:"required,min=1,max=255"`
	Note     *string                     `json:"note"`
	IsActive bool                        `json:"is_active"`
	Lines    []CreatePaymentTermLineRequest `json:"lines" validate:"required,min=1"`
}

type UpdatePaymentTermRequest struct {
	TermName *string                     `json:"term_name"`
	Note     *string                     `json:"note"`
	IsActive *bool                       `json:"is_active"`
	Lines    []CreatePaymentTermLineRequest `json:"lines"`
}

// ========================
// PAYMENT TERM LINES
// ========================

type PaymentTermLine struct {
	ID             uuid.UUID           `db:"id" json:"id"`
	PaymentTermID  uuid.UUID           `db:"payment_term_id" json:"payment_term_id"`
	Sequence       int                 `db:"sequence" json:"sequence"`
	ValueType      PaymentTermValueType `db:"value_type" json:"value_type"`
	ValueAmount    *string             `db:"value_amount" json:"value_amount,omitempty"`
	DaysAfter      int                 `db:"days_after" json:"days_after"`
	EndOfMonth     bool                `db:"end_of_month" json:"end_of_month"`
	DayOfMonth     *int                `db:"day_of_month" json:"day_of_month,omitempty"`
	CreatedAt      time.Time           `db:"created_at" json:"created_at"`
	DeletedAt      *time.Time          `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreatePaymentTermLineRequest struct {
	Sequence   int                 `json:"sequence" validate:"required,min=1"`
	ValueType  PaymentTermValueType `json:"value_type" validate:"required"`
	ValueAmount *string             `json:"value_amount"`
	DaysAfter  int                 `json:"days_after" validate:"min=0"`
	EndOfMonth bool                `json:"end_of_month"`
	DayOfMonth *int                `json:"day_of_month" validate:"omitempty,min=1,max=31"`
}

// ========================
// INVOICE PAYMENT SCHEDULES
// ========================

type InvoicePaymentScheduleStatus string

const (
	InvoicePaymentScheduleStatusPending   InvoicePaymentScheduleStatus = "pending"
	InvoicePaymentScheduleStatusPartial   InvoicePaymentScheduleStatus = "partial"
	InvoicePaymentScheduleStatusPaid      InvoicePaymentScheduleStatus = "paid"
	InvoicePaymentScheduleStatusOverdue   InvoicePaymentScheduleStatus = "overdue"
)

func (s InvoicePaymentScheduleStatus) String() string {
	return string(s)
}

type SourceType string

const (
	SourceTypeCustomerInvoice SourceType = "customer_invoice"
	SourceTypeVendorBill      SourceType = "vendor_bill"
)

func (s SourceType) String() string {
	return string(s)
}

type InvoicePaymentSchedule struct {
	ID             uuid.UUID                      `db:"id" json:"id"`
	OrganizationID uuid.UUID                      `db:"organization_id" json:"organization_id"`
	SourceType     SourceType                     `db:"source_type" json:"source_type"`
	SourceID       uuid.UUID                      `db:"source_id" json:"source_id"`
	LineNumber     int                            `db:"line_number" json:"line_number"`
	DueDate        time.Time                      `db:"due_date" json:"due_date"`
	AmountDue      string                         `db:"amount_due" json:"amount_due"`
	AmountPaid     string                         `db:"amount_paid" json:"amount_paid"`
	Status         InvoicePaymentScheduleStatus   `db:"status" json:"status"`
	CreatedAt      time.Time                      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time                      `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time                     `db:"deleted_at" json:"deleted_at,omitempty"`
}

type CreateInvoicePaymentScheduleRequest struct {
	SourceType string    `json:"source_type" validate:"required"`
	SourceID   uuid.UUID `json:"source_id" validate:"required"`
	Lines      []CreateInvoicePaymentScheduleLineRequest `json:"lines" validate:"required,min=1"`
}

type CreateInvoicePaymentScheduleLineRequest struct {
	LineNumber int       `json:"line_number" validate:"required,min=1"`
	DueDate    time.Time `json:"due_date" validate:"required"`
	AmountDue  string    `json:"amount_due" validate:"required"`
}

type UpdateInvoicePaymentScheduleRequest struct {
	AmountPaid *string `json:"amount_paid"`
	Status     *string `json:"status"`
}

// ========================
// CALCULATED PAYMENT SCHEDULE
// ========================

type PaymentScheduleCalculation struct {
	Payments      []*CalculatedPayment `json:"payments"`
	TotalAmount   string              `json:"total_amount"`
	RemainingDue  string              `json:"remaining_due"`
}

type CalculatedPayment struct {
	LineNumber  int       `json:"line_number"`
	DueDate     time.Time `json:"due_date"`
	AmountDue   string    `json:"amount_due"`
	PaymentType string    `json:"payment_type"`
}

// ========================
// VALIDATION ERRORS
// ========================

var (
	ErrInvalidCurrencyCode        = errors.New("invalid currency code format (must be 3 characters)")
	ErrCurrencyNotFound           = errors.New("currency not found")
	ErrCurrencyRateNotFound       = errors.New("currency rate not found")
	ErrPaymentTermNotFound        = errors.New("payment term not found")
	ErrPaymentTermLineNotFound    = errors.New("payment term line not found")
	ErrInvoicePaymentScheduleNotFound = errors.New("invoice payment schedule not found")
	ErrInvalidPaymentTermValue    = errors.New("invalid payment term value")
	ErrInvalidDueDate             = errors.New("due date cannot be in the past")
	ErrDuplicateCurrencyCode      = errors.New("currency code already exists")
	ErrDuplicatePaymentTermCode   = errors.New("payment term code already exists")
	ErrInvalidConversionDate      = errors.New("conversion date is in the future")
	ErrCurrencyConversionNotFound = errors.New("exchange rate not found for conversion")
	ErrPaymentTermLineValidation  = errors.New("payment term line validation failed")
)

// ========================
// HELPER FUNCTIONS
// ========================

// ParseDecimal parses a string to big.Decimal for precision arithmetic
func ParseDecimal(s string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(256)
	_, ok := f.SetString(s)
	if !ok {
		return nil, ErrInvalidPaymentTermValue
	}
	return f, nil
}

// FormatDecimal formats big.Decimal back to string with specified precision
func FormatDecimal(f *big.Float, precision int) string {
	format := "f"
	return f.Text(byte(format[0]), precision)
}

// StringValue implements driver.Valuer for string amounts
func (d string) Value() (driver.Value, error) {
	return d, nil
}

// Scan implements sql.Scanner for string amounts
func (d *string) Scan(value interface{}) error {
	if value == nil {
		*d = "0"
		return nil
	}

	switch v := value.(type) {
	case string:
		*d = v
	case []byte:
		*d = string(v)
	default:
		*d = "0"
	}
	return nil
}
