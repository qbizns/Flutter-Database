package payment

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PaymentsResponse represents a payments response
type PaymentsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	SaleId uuid.UUID `json:"sale_id"`
	
	PaymentMethod string `json:"payment_method"`
	
	PaymentStatus string `json:"payment_status"`
	
	Amount float64 `json:"amount"`
	
	CardLastFour *string `json:"card_last_four"`
	
	CardType *string `json:"card_type"`
	
	TransactionId *string `json:"transaction_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	AccountNumber *string `json:"account_number"`
	
	AccountName *string `json:"account_name"`
	
	PaymentDate time.Time `json:"payment_date"`
	
	ProcessedAt *time.Time `json:"processed_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt time.Time `json:"created_at"`
	
	UpdatedAt time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// CreatePaymentsRequest represents a request to create a payments
type CreatePaymentsRequest struct {
	
	SaleId uuid.UUID `json:"sale_id" validate:"required"`
	
	PaymentMethod string `json:"payment_method" validate:"required"`
	
	PaymentStatus string `json:"payment_status" validate:"required"`
	
	Amount float64 `json:"amount" validate:"required"`
	
	CardLastFour *string `json:"card_last_four"`
	
	CardType *string `json:"card_type"`
	
	TransactionId *string `json:"transaction_id"`
	
	ReferenceNumber *string `json:"reference_number"`
	
	AccountNumber *string `json:"account_number"`
	
	AccountName *string `json:"account_name"`
	
	PaymentDate time.Time `json:"payment_date" validate:"required"`
	
	ProcessedAt *time.Time `json:"processed_at"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
}

// Validate validates the create request
func (r *CreatePaymentsRequest) Validate() error {
	
	if r.SaleId == uuid.Nil {
		return fmt.Errorf("sale_id is required")
	}
	
	if r.PaymentMethod == "" {
		return fmt.Errorf("payment_method is required")
	}
	
	if r.PaymentStatus == "" {
		return fmt.Errorf("payment_status is required")
	}
	
	if r.Amount == 0 {
		return fmt.Errorf("amount is required")
	}

	if r.PaymentDate.IsZero() {
		return fmt.Errorf("payment_date is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdatePaymentsRequest represents a request to update a payments
type UpdatePaymentsRequest struct {
	
	SaleId *uuid.UUID `json:"sale_id,omitempty" validate:"omitempty,required"`
	
	PaymentMethod *string `json:"payment_method,omitempty" validate:"omitempty,required"`
	
	PaymentStatus *string `json:"payment_status,omitempty" validate:"omitempty,required"`
	
	Amount *float64 `json:"amount,omitempty" validate:"omitempty,required"`
	
	CardLastFour *string `json:"card_last_four,omitempty"`
	
	CardType *string `json:"card_type,omitempty"`
	
	TransactionId *string `json:"transaction_id,omitempty"`
	
	ReferenceNumber *string `json:"reference_number,omitempty"`
	
	AccountNumber *string `json:"account_number,omitempty"`
	
	AccountName *string `json:"account_name,omitempty"`
	
	PaymentDate *time.Time `json:"payment_date,omitempty" validate:"omitempty,required"`
	
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdatePaymentsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.SaleId != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.PaymentStatus != nil {
		hasUpdate = true
	}
	
	if r.Amount != nil {
		hasUpdate = true
	}
	
	if r.CardLastFour != nil {
		hasUpdate = true
	}
	
	if r.CardType != nil {
		hasUpdate = true
	}
	
	if r.TransactionId != nil {
		hasUpdate = true
	}
	
	if r.ReferenceNumber != nil {
		hasUpdate = true
	}
	
	if r.AccountNumber != nil {
		hasUpdate = true
	}
	
	if r.AccountName != nil {
		hasUpdate = true
	}
	
	if r.PaymentDate != nil {
		hasUpdate = true
	}
	
	if r.ProcessedAt != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// PaymentsListResponse represents a paginated list of payments records
type PaymentsListResponse struct {
	Items      []*PaymentsResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

// CancelPaymentRequest represents a request to cancel a payment
type CancelPaymentRequest struct {
	Reason string `json:"reason"`
}

// PaymentStatisticsResponse represents payment statistics
type PaymentStatisticsResponse struct {
	TotalAmount             float64            `json:"total_amount"`
	PaymentCount            int                `json:"payment_count"`
	AveragePayment          float64            `json:"average_payment"`
	PaymentMethodBreakdown  map[string]float64 `json:"payment_method_breakdown"`
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentID   uuid.UUID  `json:"payment_id" validate:"required"`
	OrderID     *uuid.UUID `json:"order_id"`
	Amount      float64    `json:"amount" validate:"required"`
	Reason      string     `json:"reason" validate:"required"`
	Status      string     `json:"status"`
	RequestedBy *uuid.UUID `json:"requested_by"`
	Notes       *string    `json:"notes"`
}

// RefundResponse represents a refund response
type RefundResponse struct {
	ID          uuid.UUID  `json:"id"`
	PaymentID   uuid.UUID  `json:"payment_id"`
	OrderID     *uuid.UUID `json:"order_id"`
	Amount      float64    `json:"amount"`
	Reason      string     `json:"reason"`
	Status      string     `json:"status"`
	RequestedBy *uuid.UUID `json:"requested_by"`
	RequestedAt *time.Time `json:"requested_at"`
	ProcessedAt *time.Time `json:"processed_at"`
	ProcessedBy *uuid.UUID `json:"processed_by"`
	Notes       *string    `json:"notes"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
