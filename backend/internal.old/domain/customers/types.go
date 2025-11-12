package customers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Customer represents a customer entity
type Customer struct {
	ID                     uuid.UUID  `json:"id"`
	OrganizationID         uuid.UUID  `json:"organization_id"`
	FirstName              string     `json:"first_name"`
	LastName               string     `json:"last_name"`
	Email                  string     `json:"email"`
	Phone                  string     `json:"phone"`
	DateOfBirth            *time.Time `json:"date_of_birth"`
	Gender                 string     `json:"gender"`
	AddressLine1           string     `json:"address_line_1"`
	AddressLine2           string     `json:"address_line_2"`
	City                   string     `json:"city"`
	State                  string     `json:"state"`
	PostalCode             string     `json:"postal_code"`
	Country                string     `json:"country"`
	CustomerType           string     `json:"customer_type"`
	TaxID                  string     `json:"tax_id"`
	PaymentTermDays        int        `json:"payment_term_days"`
	CreditLimit            float64    `json:"credit_limit"`
	CurrentBalance         float64    `json:"current_balance"`
	TotalSpent             float64    `json:"total_spent"`
	TotalVisits            int        `json:"total_visits"`
	LastVisitDate          *time.Time `json:"last_visit_date"`
	IsActive               bool       `json:"is_active"`
	LoyaltyMemberNumber    string     `json:"loyalty_member_number"`
	LoyaltyPoints          int        `json:"loyalty_points"`
	LoyaltyTierID          *uuid.UUID `json:"loyalty_tier_id"`
	AccountingCustomerID   *uuid.UUID `json:"accounting_customer_id"`
	PreferredPaymentMethod string     `json:"preferred_payment_method"`
	Notes                  string     `json:"notes"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	CreatedBy              uuid.UUID  `json:"created_by"`
	UpdatedBy              *uuid.UUID `json:"updated_by"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

// CustomerFilters represents filters for listing customers
type CustomerFilters struct {
	Search       string
	CustomerType *string
	IsActive     *bool
	Page         int
	PageSize     int
}

// Repository defines the customer data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters CustomerFilters) ([]Customer, error)
	Count(ctx context.Context, orgID uuid.UUID, filters CustomerFilters) (int64, error)
	Create(ctx context.Context, customer *Customer) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Customer, error)
	GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*Customer, error)
	Update(ctx context.Context, customer *Customer) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	UpdateBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, amount float64) error
	UpdateLoyaltyPoints(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, points int) error
}
