package receivables

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// InvoiceStatus represents the status of a customer invoice
type InvoiceStatus string

const (
	InvoiceStatusUnpaid   InvoiceStatus = "unpaid"
	InvoiceStatusPartial  InvoiceStatus = "partial"
	InvoiceStatusPaid     InvoiceStatus = "paid"
	InvoiceStatusOverdue  InvoiceStatus = "overdue"
	InvoiceStatusVoid     InvoiceStatus = "void"
)

// PaymentMethod represents payment method types
type PaymentMethod string

const (
	PaymentMethodCash  PaymentMethod = "cash"
	PaymentMethodCheck PaymentMethod = "check"
	PaymentMethodCard  PaymentMethod = "card"
	PaymentMethodWire  PaymentMethod = "wire"
	PaymentMethodACH   PaymentMethod = "ach"
	PaymentMethodOther PaymentMethod = "other"
)

// CustomerInvoice represents a customer invoice
type CustomerInvoice struct {
	ID                 uuid.UUID                 `json:"id"`
	OrganizationID     uuid.UUID                 `json:"organization_id"`
	InvoiceNumber      string                    `json:"invoice_number"`
	CustomerID         uuid.UUID                 `json:"customer_id"`
	InvoiceDate        time.Time                 `json:"invoice_date"`
	DueDate            time.Time                 `json:"due_date"`
	PaymentTerms       *string                   `json:"payment_terms"`
	AccountingPeriodID *uuid.UUID                `json:"accounting_period_id"`
	Subtotal           float64                   `json:"subtotal"`
	TaxAmount          float64                   `json:"tax_amount"`
	DiscountAmount     float64                   `json:"discount_amount"`
	TotalAmount        float64                   `json:"total_amount"`
	PaidAmount         float64                   `json:"paid_amount"`
	BalanceDue         float64                   `json:"balance_due"`
	Status             InvoiceStatus             `json:"status"`
	JournalEntryID     *uuid.UUID                `json:"journal_entry_id"`
	IsPosted           bool                      `json:"is_posted"`
	SaleID             *uuid.UUID                `json:"sale_id"`
	Description        *string                   `json:"description"`
	Notes              *string                   `json:"notes"`
	Memo               *string                   `json:"memo"`
	Attachments        json.RawMessage             `json:"attachments"`
	Metadata           json.RawMessage               `json:"metadata"`
	CreatedAt          time.Time                 `json:"created_at"`
	UpdatedAt          time.Time                 `json:"updated_at"`
	CreatedBy          *uuid.UUID                `json:"created_by"`
	UpdatedBy          *uuid.UUID                `json:"updated_by"`
	DeletedAt          *time.Time                `json:"deleted_at"`
}

// CustomerInvoiceLine represents a line item in an invoice
type CustomerInvoiceLine struct {
	ID              uuid.UUID      `json:"id"`
	OrganizationID  uuid.UUID      `json:"organization_id"`
	CustomerInvoiceID uuid.UUID    `json:"customer_invoice_id"`
	LineNumber      int            `json:"line_number"`
	RevenueAccountID uuid.UUID     `json:"revenue_account_id"`
	Description     string         `json:"description"`
	Quantity        float64        `json:"quantity"`
	UnitPrice       float64        `json:"unit_price"`
	Amount          float64        `json:"amount"`
	LocationID      *uuid.UUID     `json:"location_id"`
	Department      *string        `json:"department"`
	ProjectCode     *string        `json:"project_code"`
	TaxCode         *string        `json:"tax_code"`
	TaxAmount       float64        `json:"tax_amount"`
	ProductID       *uuid.UUID     `json:"product_id"`
	Metadata        json.RawMessage    `json:"metadata"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at"`
}

// CustomerPayment represents a payment received from a customer
type CustomerPayment struct {
	ID                 uuid.UUID      `json:"id"`
	OrganizationID     uuid.UUID      `json:"organization_id"`
	PaymentNumber      string         `json:"payment_number"`
	CustomerID         uuid.UUID      `json:"customer_id"`
	PaymentDate        time.Time      `json:"payment_date"`
	PaymentMethod      PaymentMethod  `json:"payment_method"`
	ReferenceNumber    *string        `json:"reference_number"`
	PaymentAmount      float64        `json:"payment_amount"`
	DepositAccountID   *uuid.UUID     `json:"deposit_account_id"`
	AccountingPeriodID *uuid.UUID     `json:"accounting_period_id"`
	JournalEntryID     *uuid.UUID     `json:"journal_entry_id"`
	IsPosted           bool           `json:"is_posted"`
	Memo               *string        `json:"memo"`
	Notes              *string        `json:"notes"`
	Metadata           json.RawMessage    `json:"metadata"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedBy          *uuid.UUID     `json:"created_by"`
	UpdatedBy          *uuid.UUID     `json:"updated_by"`
	DeletedAt          *time.Time     `json:"deleted_at"`
}

// CustomerPaymentApplication represents the application of a payment to an invoice
type CustomerPaymentApplication struct {
	ID                  uuid.UUID  `json:"id"`
	OrganizationID      uuid.UUID  `json:"organization_id"`
	CustomerPaymentID   uuid.UUID  `json:"customer_payment_id"`
	CustomerInvoiceID   uuid.UUID  `json:"customer_invoice_id"`
	AppliedAmount       float64    `json:"applied_amount"`
	CreatedAt           time.Time  `json:"created_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// CustomerInvoiceWithLines represents an invoice with its line items
type CustomerInvoiceWithLines struct {
	Invoice *CustomerInvoice   `json:"invoice"`
	Lines   []CustomerInvoiceLine `json:"lines"`
}

// CustomerPaymentWithApplications represents a payment with its applications
type CustomerPaymentWithApplications struct {
	Payment      *CustomerPayment               `json:"payment"`
	Applications []CustomerPaymentApplication   `json:"applications"`
}

// CreateCustomerInvoiceRequest represents the request to create a customer invoice
type CreateCustomerInvoiceRequest struct {
	InvoiceNumber      string                      `json:"invoice_number"`
	CustomerID         uuid.UUID                   `json:"customer_id"`
	InvoiceDate        time.Time                   `json:"invoice_date"`
	DueDate            time.Time                   `json:"due_date"`
	PaymentTerms       *string                     `json:"payment_terms"`
	AccountingPeriodID *uuid.UUID                  `json:"accounting_period_id"`
	Subtotal           float64                     `json:"subtotal"`
	TaxAmount          float64                     `json:"tax_amount"`
	DiscountAmount     float64                     `json:"discount_amount"`
	TotalAmount        float64                     `json:"total_amount"`
	SaleID             *uuid.UUID                  `json:"sale_id"`
	Description        *string                     `json:"description"`
	Notes              *string                     `json:"notes"`
	Memo               *string                     `json:"memo"`
	Attachments        []interface{}               `json:"attachments"`
	Metadata           map[string]interface{}      `json:"metadata"`
	Lines              []CreateCustomerInvoiceLineRequest `json:"lines"`
}

// CreateCustomerInvoiceLineRequest represents the request to create an invoice line
type CreateCustomerInvoiceLineRequest struct {
	LineNumber       int                   `json:"line_number"`
	RevenueAccountID uuid.UUID             `json:"revenue_account_id"`
	Description      string                `json:"description"`
	Quantity         float64               `json:"quantity"`
	UnitPrice        float64               `json:"unit_price"`
	Amount           float64               `json:"amount"`
	LocationID       *uuid.UUID            `json:"location_id"`
	Department       *string               `json:"department"`
	ProjectCode      *string               `json:"project_code"`
	TaxCode          *string               `json:"tax_code"`
	TaxAmount        float64               `json:"tax_amount"`
	ProductID        *uuid.UUID            `json:"product_id"`
	Metadata         map[string]interface{} `json:"metadata"`
}

// UpdateCustomerInvoiceRequest represents the request to update a customer invoice
type UpdateCustomerInvoiceRequest struct {
	InvoiceDate        *time.Time                    `json:"invoice_date"`
	DueDate            *time.Time                    `json:"due_date"`
	PaymentTerms       *string                       `json:"payment_terms"`
	Subtotal           *float64                      `json:"subtotal"`
	TaxAmount          *float64                      `json:"tax_amount"`
	DiscountAmount     *float64                      `json:"discount_amount"`
	TotalAmount        *float64                      `json:"total_amount"`
	Description        *string                       `json:"description"`
	Notes              *string                       `json:"notes"`
	Memo               *string                       `json:"memo"`
	Status             *InvoiceStatus                `json:"status"`
	Metadata           map[string]interface{}        `json:"metadata"`
	Lines              []UpdateCustomerInvoiceLineRequest `json:"lines"`
}

// UpdateCustomerInvoiceLineRequest represents the request to update an invoice line
type UpdateCustomerInvoiceLineRequest struct {
	ID               *uuid.UUID                `json:"id"`
	LineNumber       *int                      `json:"line_number"`
	Description      *string                   `json:"description"`
	Quantity         *float64                  `json:"quantity"`
	UnitPrice        *float64                  `json:"unit_price"`
	Amount           *float64                  `json:"amount"`
	TaxAmount        *float64                  `json:"tax_amount"`
	LocationID       *uuid.UUID                `json:"location_id"`
	Department       *string                   `json:"department"`
	ProjectCode      *string                   `json:"project_code"`
	TaxCode          *string                   `json:"tax_code"`
}

// CreateCustomerPaymentRequest represents the request to create a customer payment
type CreateCustomerPaymentRequest struct {
	PaymentNumber      string                      `json:"payment_number"`
	CustomerID         uuid.UUID                   `json:"customer_id"`
	PaymentDate        time.Time                   `json:"payment_date"`
	PaymentMethod      PaymentMethod               `json:"payment_method"`
	ReferenceNumber    *string                     `json:"reference_number"`
	PaymentAmount      float64                     `json:"payment_amount"`
	DepositAccountID   *uuid.UUID                  `json:"deposit_account_id"`
	AccountingPeriodID *uuid.UUID                  `json:"accounting_period_id"`
	Memo               *string                     `json:"memo"`
	Notes              *string                     `json:"notes"`
	Metadata           map[string]interface{}      `json:"metadata"`
	Applications       []CreatePaymentApplicationRequest `json:"applications"`
}

// CreatePaymentApplicationRequest represents the request to apply a payment to an invoice
type CreatePaymentApplicationRequest struct {
	CustomerInvoiceID uuid.UUID `json:"customer_invoice_id"`
	AppliedAmount     float64   `json:"applied_amount"`
}

// UpdateCustomerPaymentRequest represents the request to update a customer payment
type UpdateCustomerPaymentRequest struct {
	PaymentDate        *time.Time             `json:"payment_date"`
	PaymentMethod      *PaymentMethod         `json:"payment_method"`
	ReferenceNumber    *string                `json:"reference_number"`
	PaymentAmount      *float64               `json:"payment_amount"`
	DepositAccountID   *uuid.UUID             `json:"deposit_account_id"`
	Memo               *string                `json:"memo"`
	Notes              *string                `json:"notes"`
	Metadata           map[string]interface{} `json:"metadata"`
	Applications       []CreatePaymentApplicationRequest `json:"applications"`
}

// CustomerAgingReport represents an aging bucket for AR
type CustomerAgingReport struct {
	CustomerID   uuid.UUID  `json:"customer_id"`
	CustomerName string     `json:"customer_name"`
	Current      float64    `json:"current"`
	Days30       float64    `json:"days_30"`
	Days60       float64    `json:"days_60"`
	Days90Plus   float64    `json:"days_90_plus"`
	TotalDue     float64    `json:"total_due"`
}

// CustomerInvoiceFilterOptions represents filter options for invoice queries
type CustomerInvoiceFilterOptions struct {
	Status           *InvoiceStatus
	CustomerID       *uuid.UUID
	InvoiceDateFrom  *time.Time
	InvoiceDateTo    *time.Time
	DueDateFrom      *time.Time
	DueDateTo        *time.Time
	IsPosted         *bool
	PeriodID         *uuid.UUID
	DeletedIncluded  bool
	Limit            int
	Offset           int
}

// CustomerPaymentFilterOptions represents filter options for payment queries
type CustomerPaymentFilterOptions struct {
	CustomerID      *uuid.UUID
	PaymentMethod   *PaymentMethod
	PaymentDateFrom *time.Time
	PaymentDateTo   *time.Time
	IsPosted        *bool
	PeriodID        *uuid.UUID
	DeletedIncluded bool
	Limit           int
	Offset          int
}
