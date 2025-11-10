package payables

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// VendorBillStatus represents the status of a vendor bill
type VendorBillStatus string

const (
	VendorBillStatusUnpaid   VendorBillStatus = "unpaid"
	VendorBillStatusPartial  VendorBillStatus = "partial"
	VendorBillStatusPaid     VendorBillStatus = "paid"
	VendorBillStatusOverdue  VendorBillStatus = "overdue"
	VendorBillStatusVoid     VendorBillStatus = "void"
)

// PaymentMethod represents payment method types
type PaymentMethod string

const (
	PaymentMethodCheck PaymentMethod = "check"
	PaymentMethodCash  PaymentMethod = "cash"
	PaymentMethodWire  PaymentMethod = "wire"
	PaymentMethodACH   PaymentMethod = "ach"
	PaymentMethodCard  PaymentMethod = "card"
	PaymentMethodOther PaymentMethod = "other"
)

// VendorBill represents a vendor bill/invoice from a supplier
type VendorBill struct {
	ID                   uuid.UUID                  `json:"id"`
	OrganizationID       uuid.UUID                  `json:"organization_id"`
	BillNumber           string                     `json:"bill_number"`
	VendorBillNumber     *string                    `json:"vendor_bill_number"`
	SupplierID           uuid.UUID                  `json:"supplier_id"`
	BillDate             time.Time                  `json:"bill_date"`
	DueDate              time.Time                  `json:"due_date"`
	PaymentTerms         *string                    `json:"payment_terms"`
	AccountingPeriodID   *uuid.UUID                 `json:"accounting_period_id"`
	Subtotal             float64                    `json:"subtotal"`
	TaxAmount            float64                    `json:"tax_amount"`
	TotalAmount          float64                    `json:"total_amount"`
	PaidAmount           float64                    `json:"paid_amount"`
	BalanceDue           float64                    `json:"balance_due"`
	Status               VendorBillStatus           `json:"status"`
	JournalEntryID       *uuid.UUID                 `json:"journal_entry_id"`
	IsPosted             bool                       `json:"is_posted"`
	PurchaseOrderID      *uuid.UUID                 `json:"purchase_order_id"`
	Description          *string                    `json:"description"`
	Notes                *string                    `json:"notes"`
	Memo                 *string                    `json:"memo"`
	Attachments          json.RawMessage              `json:"attachments"`
	Metadata             json.RawMessage                `json:"metadata"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
	CreatedBy            *uuid.UUID                 `json:"created_by"`
	UpdatedBy            *uuid.UUID                 `json:"updated_by"`
	DeletedAt            *time.Time                 `json:"deleted_at"`
}

// VendorBillLine represents a line item in a vendor bill
type VendorBillLine struct {
	ID                 uuid.UUID      `json:"id"`
	OrganizationID     uuid.UUID      `json:"organization_id"`
	VendorBillID       uuid.UUID      `json:"vendor_bill_id"`
	LineNumber         int            `json:"line_number"`
	ExpenseAccountID   uuid.UUID      `json:"expense_account_id"`
	Description        string         `json:"description"`
	Quantity           float64        `json:"quantity"`
	UnitPrice          float64        `json:"unit_price"`
	Amount             float64        `json:"amount"`
	LocationID         *uuid.UUID     `json:"location_id"`
	Department         *string        `json:"department"`
	ProjectCode        *string        `json:"project_code"`
	TaxCode            *string        `json:"tax_code"`
	TaxAmount          float64        `json:"tax_amount"`
	ProductID          *uuid.UUID     `json:"product_id"`
	Metadata           json.RawMessage    `json:"metadata"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          *time.Time     `json:"deleted_at"`
}

// VendorPayment represents a payment made to a vendor
type VendorPayment struct {
	ID                 uuid.UUID      `json:"id"`
	OrganizationID     uuid.UUID      `json:"organization_id"`
	PaymentNumber      string         `json:"payment_number"`
	SupplierID         uuid.UUID      `json:"supplier_id"`
	PaymentDate        time.Time      `json:"payment_date"`
	PaymentMethod      PaymentMethod  `json:"payment_method"`
	ReferenceNumber    *string        `json:"reference_number"`
	PaymentAmount      float64        `json:"payment_amount"`
	BankAccountID      *uuid.UUID     `json:"bank_account_id"`
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

// VendorPaymentApplication represents the application of a payment to a bill
type VendorPaymentApplication struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	VendorPaymentID  uuid.UUID  `json:"vendor_payment_id"`
	VendorBillID     uuid.UUID  `json:"vendor_bill_id"`
	AppliedAmount    float64    `json:"applied_amount"`
	CreatedAt        time.Time  `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

// VendorBillWithLines represents a bill with its line items
type VendorBillWithLines struct {
	Bill  *VendorBill   `json:"bill"`
	Lines []VendorBillLine `json:"lines"`
}

// VendorPaymentWithApplications represents a payment with its applications
type VendorPaymentWithApplications struct {
	Payment      *VendorPayment                `json:"payment"`
	Applications []VendorPaymentApplication    `json:"applications"`
}

// CreateVendorBillRequest represents the request to create a vendor bill
type CreateVendorBillRequest struct {
	BillNumber         string                    `json:"bill_number"`
	VendorBillNumber   *string                   `json:"vendor_bill_number"`
	SupplierID         uuid.UUID                 `json:"supplier_id"`
	BillDate           time.Time                 `json:"bill_date"`
	DueDate            time.Time                 `json:"due_date"`
	PaymentTerms       *string                   `json:"payment_terms"`
	AccountingPeriodID *uuid.UUID                `json:"accounting_period_id"`
	Subtotal           float64                   `json:"subtotal"`
	TaxAmount          float64                   `json:"tax_amount"`
	TotalAmount        float64                   `json:"total_amount"`
	PurchaseOrderID    *uuid.UUID                `json:"purchase_order_id"`
	Description        *string                   `json:"description"`
	Notes              *string                   `json:"notes"`
	Memo               *string                   `json:"memo"`
	Attachments        []interface{}             `json:"attachments"`
	Metadata           map[string]interface{}    `json:"metadata"`
	Lines              []CreateVendorBillLineRequest `json:"lines"`
}

// CreateVendorBillLineRequest represents the request to create a bill line
type CreateVendorBillLineRequest struct {
	LineNumber       int                   `json:"line_number"`
	ExpenseAccountID uuid.UUID             `json:"expense_account_id"`
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

// UpdateVendorBillRequest represents the request to update a vendor bill
type UpdateVendorBillRequest struct {
	VendorBillNumber   *string                   `json:"vendor_bill_number"`
	BillDate           *time.Time                `json:"bill_date"`
	DueDate            *time.Time                `json:"due_date"`
	PaymentTerms       *string                   `json:"payment_terms"`
	Subtotal           *float64                  `json:"subtotal"`
	TaxAmount          *float64                  `json:"tax_amount"`
	TotalAmount        *float64                  `json:"total_amount"`
	Description        *string                   `json:"description"`
	Notes              *string                   `json:"notes"`
	Memo               *string                   `json:"memo"`
	Status             *VendorBillStatus         `json:"status"`
	Metadata           map[string]interface{}    `json:"metadata"`
	Lines              []UpdateVendorBillLineRequest `json:"lines"`
}

// UpdateVendorBillLineRequest represents the request to update a bill line
type UpdateVendorBillLineRequest struct {
	ID           *uuid.UUID                `json:"id"`
	LineNumber   *int                      `json:"line_number"`
	Description  *string                   `json:"description"`
	Quantity     *float64                  `json:"quantity"`
	UnitPrice    *float64                  `json:"unit_price"`
	Amount       *float64                  `json:"amount"`
	TaxAmount    *float64                  `json:"tax_amount"`
	LocationID   *uuid.UUID                `json:"location_id"`
	Department   *string                   `json:"department"`
	ProjectCode  *string                   `json:"project_code"`
	TaxCode      *string                   `json:"tax_code"`
}

// CreateVendorPaymentRequest represents the request to create a vendor payment
type CreateVendorPaymentRequest struct {
	PaymentNumber      string                    `json:"payment_number"`
	SupplierID         uuid.UUID                 `json:"supplier_id"`
	PaymentDate        time.Time                 `json:"payment_date"`
	PaymentMethod      PaymentMethod             `json:"payment_method"`
	ReferenceNumber    *string                   `json:"reference_number"`
	PaymentAmount      float64                   `json:"payment_amount"`
	BankAccountID      *uuid.UUID                `json:"bank_account_id"`
	AccountingPeriodID *uuid.UUID                `json:"accounting_period_id"`
	Memo               *string                   `json:"memo"`
	Notes              *string                   `json:"notes"`
	Metadata           map[string]interface{}    `json:"metadata"`
	Applications       []CreatePaymentApplicationRequest `json:"applications"`
}

// CreatePaymentApplicationRequest represents the request to apply a payment to a bill
type CreatePaymentApplicationRequest struct {
	VendorBillID  uuid.UUID `json:"vendor_bill_id"`
	AppliedAmount float64   `json:"applied_amount"`
}

// UpdateVendorPaymentRequest represents the request to update a vendor payment
type UpdateVendorPaymentRequest struct {
	PaymentDate        *time.Time             `json:"payment_date"`
	PaymentMethod      *PaymentMethod         `json:"payment_method"`
	ReferenceNumber    *string                `json:"reference_number"`
	PaymentAmount      *float64               `json:"payment_amount"`
	BankAccountID      *uuid.UUID             `json:"bank_account_id"`
	Memo               *string                `json:"memo"`
	Notes              *string                `json:"notes"`
	Metadata           map[string]interface{} `json:"metadata"`
	Applications       []CreatePaymentApplicationRequest `json:"applications"`
}

// VendorAgingReport represents an aging bucket for AP
type VendorAgingReport struct {
	SupplierID   uuid.UUID  `json:"supplier_id"`
	SupplierName string     `json:"supplier_name"`
	Current      float64    `json:"current"`
	Days30       float64    `json:"days_30"`
	Days60       float64    `json:"days_60"`
	Days90Plus   float64    `json:"days_90_plus"`
	TotalDue     float64    `json:"total_due"`
}

// VendorBillFilterOptions represents filter options for bill queries
type VendorBillFilterOptions struct {
	Status          *VendorBillStatus
	SupplierID      *uuid.UUID
	BillDateFrom    *time.Time
	BillDateTo      *time.Time
	DueDateFrom     *time.Time
	DueDateTo       *time.Time
	IsPosted        *bool
	PeriodID        *uuid.UUID
	DeletedIncluded bool
	Limit           int
	Offset          int
}

// VendorPaymentFilterOptions represents filter options for payment queries
type VendorPaymentFilterOptions struct {
	SupplierID      *uuid.UUID
	PaymentMethod   *PaymentMethod
	PaymentDateFrom *time.Time
	PaymentDateTo   *time.Time
	IsPosted        *bool
	PeriodID        *uuid.UUID
	DeletedIncluded bool
	Limit           int
	Offset          int
}
