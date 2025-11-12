package suppliers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Supplier represents a supplier entity
type Supplier struct {
	ID                 uuid.UUID  `json:"id"`
	OrganizationID     uuid.UUID  `json:"organization_id"`
	SupplierCode       string     `json:"supplier_code"`
	Name               string     `json:"name"`
	ContactPerson      string     `json:"contact_person"`
	Email              string     `json:"email"`
	Phone              string     `json:"phone"`
	Address            string     `json:"address"`
	City               string     `json:"city"`
	State              string     `json:"state"`
	Country            string     `json:"country"`
	PostalCode         string     `json:"postal_code"`
	TaxNumber          string     `json:"tax_number"`
	PaymentTerms       string     `json:"payment_terms"`
	CreditLimit        float64    `json:"credit_limit"`
	OutstandingBalance float64    `json:"outstanding_balance"`
	TotalPurchases     float64    `json:"total_purchases"`
	TotalOrders        int        `json:"total_orders"`
	LastOrderDate      *time.Time `json:"last_order_date"`
	Status             string     `json:"status"`
	Notes              string     `json:"notes"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedBy          uuid.UUID  `json:"created_by"`
	UpdatedBy          *uuid.UUID `json:"updated_by"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// SupplierFilters represents filters for listing suppliers
type SupplierFilters struct {
	Search   string
	Status   *string
	Page     int
	PageSize int
}

// Repository defines the supplier data access interface
type Repository interface {
	List(ctx context.Context, orgID uuid.UUID, filters SupplierFilters) ([]Supplier, error)
	Count(ctx context.Context, orgID uuid.UUID, filters SupplierFilters) (int64, error)
	Create(ctx context.Context, supplier *Supplier) error
	Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Supplier, error)
	GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*Supplier, error)
	Update(ctx context.Context, supplier *Supplier) error
	Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
}
