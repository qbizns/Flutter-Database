package giftcards

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GiftCard represents a gift card entity
type GiftCard struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	CardNumber       string     `json:"card_number"`
	PINCode          *string    `json:"pin_code,omitempty"`
	CustomerID       *uuid.UUID `json:"customer_id"`
	OriginalValue    float64    `json:"original_value"`
	CurrentBalance   float64    `json:"current_balance"`
	IssuedDate       time.Time  `json:"issued_date"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	Status           string     `json:"status"` // active, inactive, blocked, expired, fully_redeemed
	IssuedByUserID   *uuid.UUID `json:"issued_by_user_id"`
	IssuedLocationID *uuid.UUID `json:"issued_location_id"`
	Notes            *string    `json:"notes"`
	CreatedBy        uuid.UUID  `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// GiftCardTransaction represents a transaction on a gift card
type GiftCardTransaction struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	GiftCardID      uuid.UUID  `json:"gift_card_id"`
	TransactionType string     `json:"transaction_type"` // issue, load, redemption, refund, adjustment, expiry
	Amount          float64    `json:"amount"`
	BalanceAfter    float64    `json:"balance_after"`
	SaleID          *uuid.UUID `json:"sale_id"`
	PaymentID       *uuid.UUID `json:"payment_id"`
	UserID          *uuid.UUID `json:"user_id"`
	LocationID      *uuid.UUID `json:"location_id"`
	Notes           *string    `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// StoreCreditAccount represents a customer's store credit account
type StoreCreditAccount struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	CustomerID     uuid.UUID  `json:"customer_id"`
	CurrentBalance float64    `json:"current_balance"`
	CreditLimit    *float64   `json:"credit_limit"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// StoreCreditTransaction represents a transaction on store credit
type StoreCreditTransaction struct {
	ID                   uuid.UUID  `json:"id"`
	OrganizationID       uuid.UUID  `json:"organization_id"`
	StoreCreditAccountID uuid.UUID  `json:"store_credit_account_id"`
	TransactionType      string     `json:"transaction_type"` // issue, redemption, refund, adjustment, expiry
	Amount               float64    `json:"amount"`
	BalanceAfter         float64    `json:"balance_after"`
	SaleID               *uuid.UUID `json:"sale_id"`
	PaymentID            *uuid.UUID `json:"payment_id"`
	UserID               *uuid.UUID `json:"user_id"`
	LocationID           *uuid.UUID `json:"location_id"`
	Notes                *string    `json:"notes"`
	CreatedAt            time.Time  `json:"created_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

// ReturnReason represents a return reason
type ReturnReason struct {
	ID               uuid.UUID  `json:"id"`
	OrganizationID   *uuid.UUID `json:"organization_id"`
	ReasonCode       string     `json:"reason_code"`
	ReasonName       string     `json:"reason_name"`
	RequiresApproval bool       `json:"requires_approval"`
	AffectsInventory bool       `json:"affects_inventory"`
	IsRestockable    bool       `json:"is_restockable"`
	IsActive         bool       `json:"is_active"`
	DisplayOrder     int        `json:"display_order"`
	CreatedAt        time.Time  `json:"created_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// SaleReturn represents a return/RMA header
type SaleReturn struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	ReturnNumber   string     `json:"return_number"`
	OriginalSaleID *uuid.UUID `json:"original_sale_id"`
	CustomerID     *uuid.UUID `json:"customer_id"`
	LocationID     uuid.UUID  `json:"location_id"`
	UserID         uuid.UUID  `json:"user_id"`
	ReturnDate     time.Time  `json:"return_date"`
	TotalAmount    float64    `json:"total_amount"`
	RefundAmount   float64    `json:"refund_amount"`
	RestockingFee  float64    `json:"restocking_fee"`
	RefundMethod   *string    `json:"refund_method"` // original_payment, cash, store_credit, gift_card, exchange
	Status         string     `json:"status"`        // pending, approved, refunded, rejected, completed
	ApprovedBy     *uuid.UUID `json:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at"`
	Notes          *string    `json:"notes"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// SaleReturnItem represents an item in a return
type SaleReturnItem struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	SaleReturnID       uuid.UUID  `json:"sale_return_id"`
	OriginalSaleItemID *uuid.UUID `json:"original_sale_item_id"`
	ProductID          uuid.UUID  `json:"product_id"`
	ProductVariantID   *uuid.UUID `json:"product_variant_id"`
	Quantity           float64    `json:"quantity"`
	UnitPrice          float64    `json:"unit_price"`
	Subtotal           float64    `json:"subtotal"`
	TaxAmount          float64    `json:"tax_amount"`
	DiscountAmount     float64    `json:"discount_amount"`
	TotalAmount        float64    `json:"total_amount"`
	ReturnReasonID     *uuid.UUID `json:"return_reason_id"`
	ReturnReasonNotes  *string    `json:"return_reason_notes"`
	ItemCondition      string     `json:"item_condition"` // resellable, damaged, defective, used, opened
	IsRestockable      bool       `json:"is_restockable"`
	CreatedAt          time.Time  `json:"created_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// GiftCardFilters for list operations
type GiftCardFilters struct {
	Status   *string
	Customer *uuid.UUID
	Search   string
	Page     int
	PageSize int
}

// StoreCreditFilters for list operations
type StoreCreditFilters struct {
	IsActive *bool
	Customer *uuid.UUID
	Page     int
	PageSize int
}

// ReturnFilters for list operations
type ReturnFilters struct {
	Status   *string
	Customer *uuid.UUID
	Location *uuid.UUID
	DateFrom *time.Time
	DateTo   *time.Time
	Page     int
	PageSize int
}

// GiftCardRepository defines the gift card data access interface
type GiftCardRepository interface {
	// Gift Card CRUD
	CreateGiftCard(ctx context.Context, card *GiftCard) error
	GetGiftCard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GiftCard, error)
	GetGiftCardByNumber(ctx context.Context, orgID uuid.UUID, cardNumber string) (*GiftCard, error)
	ListGiftCards(ctx context.Context, orgID uuid.UUID, filters GiftCardFilters) ([]GiftCard, error)
	UpdateGiftCard(ctx context.Context, card *GiftCard) error
	DeleteGiftCard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Gift Card Transactions
	CreateGiftCardTransaction(ctx context.Context, txn *GiftCardTransaction) error
	GetGiftCardTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GiftCardTransaction, error)
	ListGiftCardTransactions(ctx context.Context, orgID uuid.UUID, cardID uuid.UUID) ([]GiftCardTransaction, error)

	// Store Credit CRUD
	CreateStoreCreditAccount(ctx context.Context, account *StoreCreditAccount) error
	GetStoreCreditAccount(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*StoreCreditAccount, error)
	ListStoreCreditAccounts(ctx context.Context, orgID uuid.UUID, filters StoreCreditFilters) ([]StoreCreditAccount, error)
	UpdateStoreCreditAccount(ctx context.Context, account *StoreCreditAccount) error

	// Store Credit Transactions
	CreateStoreCreditTransaction(ctx context.Context, txn *StoreCreditTransaction) error
	ListStoreCreditTransactions(ctx context.Context, orgID uuid.UUID, accountID uuid.UUID) ([]StoreCreditTransaction, error)

	// Return Reasons
	ListReturnReasons(ctx context.Context, orgID *uuid.UUID) ([]ReturnReason, error)
	GetReturnReason(ctx context.Context, id uuid.UUID) (*ReturnReason, error)

	// Sale Returns CRUD
	CreateSaleReturn(ctx context.Context, return_ *SaleReturn) error
	GetSaleReturn(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SaleReturn, error)
	GetSaleReturnByNumber(ctx context.Context, orgID uuid.UUID, returnNumber string) (*SaleReturn, error)
	ListSaleReturns(ctx context.Context, orgID uuid.UUID, filters ReturnFilters) ([]SaleReturn, error)
	UpdateSaleReturn(ctx context.Context, return_ *SaleReturn) error

	// Sale Return Items
	CreateSaleReturnItem(ctx context.Context, item *SaleReturnItem) error
	ListSaleReturnItems(ctx context.Context, orgID uuid.UUID, returnID uuid.UUID) ([]SaleReturnItem, error)
}
