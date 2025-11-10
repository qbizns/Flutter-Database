package payments

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PaymentMethod represents the payment method type
type PaymentMethod string

const (
	PaymentMethodCash   PaymentMethod = "cash"
	PaymentMethodCard   PaymentMethod = "card"
	PaymentMethodMobile PaymentMethod = "mobile"
	PaymentMethodBank   PaymentMethod = "bank"
	PaymentMethodCheck  PaymentMethod = "check"
	PaymentMethodCredit PaymentMethod = "credit"
)

// PaymentStatus represents the payment status
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCanceled  PaymentStatus = "canceled"
)

// Scan implements the sql.Scanner interface
func (p *PaymentMethod) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*p = PaymentMethod(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*p = PaymentMethod(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (p PaymentMethod) Value() (driver.Value, error) {
	return string(p), nil
}

// Scan implements the sql.Scanner interface
func (p *PaymentStatus) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(string); ok {
		*p = PaymentStatus(v)
		return nil
	}
	if v, ok := value.([]byte); ok {
		*p = PaymentStatus(v)
		return nil
	}
	return nil
}

// Value implements the driver.Valuer interface
func (p PaymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}

// Payment represents a payment entity
type Payment struct {
	ID                uuid.UUID     `json:"id"`
	OrganizationID    uuid.UUID     `json:"organization_id"`
	SaleID            uuid.UUID     `json:"sale_id"`
	PaymentMethod     PaymentMethod `json:"payment_method"`
	PaymentStatus     PaymentStatus `json:"payment_status"`
	Amount            float64       `json:"amount"`
	CardLastFour      *string       `json:"card_last_four,omitempty"`
	CardType          *string       `json:"card_type,omitempty"`
	TransactionID     *string       `json:"transaction_id,omitempty"`
	ReferenceNumber   *string       `json:"reference_number,omitempty"`
	AccountNumber     *string       `json:"account_number,omitempty"`
	AccountName       *string       `json:"account_name,omitempty"`
	PaymentDate       time.Time     `json:"payment_date"`
	ProcessedAt       *time.Time    `json:"processed_at,omitempty"`
	Notes             string        `json:"notes"`
	Metadata          interface{}   `json:"metadata"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	CreatedBy         uuid.UUID     `json:"created_by"`
}

// PaymentFilters represents filters for listing payments
type PaymentFilters struct {
	SaleID        *uuid.UUID
	PaymentMethod *PaymentMethod
	PaymentStatus *PaymentStatus
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	PageSize      int
}

// Repository defines the payment data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters PaymentFilters) ([]Payment, error)
	Count(ctx context.Context, orgID uuid.UUID, filters PaymentFilters) (int64, error)
	Create(ctx context.Context, payment *Payment) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Payment, error)
	GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]Payment, error)
	Update(ctx context.Context, payment *Payment) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateStatus(ctx context.Context, orgID uuid.UUID, paymentID uuid.UUID, status PaymentStatus) error
}
