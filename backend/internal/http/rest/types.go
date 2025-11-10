package rest

import (
	"time"

	"github.com/google/uuid"
)

// Common response types
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationParams struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// ============================================================================
// AUTHENTICATION
// ============================================================================

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type RegisterRequest struct {
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password" validate:"required,min=8"`
	FirstName      string    `json:"first_name" validate:"required"`
	LastName       string    `json:"last_name" validate:"required"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

type UserInfo struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Roles          []string  `json:"roles"`
}

// ============================================================================
// ORGANIZATIONS
// ============================================================================

type CreateOrganizationRequest struct {
	Name            string `json:"name" validate:"required,min=3,max=255"`
	TaxID           string `json:"tax_id"`
	RegistrationNo  string `json:"registration_no"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone"`
	Website         string `json:"website"`
	AddressLine1    string `json:"address_line_1"`
	AddressLine2    string `json:"address_line_2"`
	City            string `json:"city"`
	State           string `json:"state"`
	PostalCode      string `json:"postal_code"`
	Country         string `json:"country"`
	DefaultCurrency string `json:"default_currency"`
	DefaultLanguage string `json:"default_language"`
}

type UpdateOrganizationRequest struct {
	Name            *string `json:"name" validate:"omitempty,min=3,max=255"`
	TaxID           *string `json:"tax_id"`
	RegistrationNo  *string `json:"registration_no"`
	Email           *string `json:"email" validate:"omitempty,email"`
	Phone           *string `json:"phone"`
	Website         *string `json:"website"`
	AddressLine1    *string `json:"address_line_1"`
	AddressLine2    *string `json:"address_line_2"`
	City            *string `json:"city"`
	State           *string `json:"state"`
	PostalCode      *string `json:"postal_code"`
	Country         *string `json:"country"`
	DefaultCurrency *string `json:"default_currency"`
	DefaultLanguage *string `json:"default_language"`
	IsActive        *bool   `json:"is_active"`
}

// ============================================================================
// PRODUCTS
// ============================================================================

type CreateProductRequest struct {
	Name                string     `json:"name" validate:"required,min=3,max=255"`
	SKU                 string     `json:"sku" validate:"required"`
	Barcode             string     `json:"barcode"`
	Description         string     `json:"description"`
	CategoryID          *uuid.UUID `json:"category_id"`
	UnitPrice           float64    `json:"unit_price" validate:"required,gte=0"`
	Cost                float64    `json:"cost" validate:"gte=0"`
	TaxRate             float64    `json:"tax_rate" validate:"gte=0,lte=100"`
	Unit                string     `json:"unit"`
	MinStockLevel       float64    `json:"min_stock_level" validate:"gte=0"`
	MaxStockLevel       float64    `json:"max_stock_level" validate:"gte=0"`
	IsActive            bool       `json:"is_active"`
	IsTrackInventory    bool       `json:"is_track_inventory"`
	AllowNegativeStock  bool       `json:"allow_negative_stock"`
	ImageURL            string     `json:"image_url"`
	ProductType         string     `json:"product_type"`
	AccountingAccountID *uuid.UUID `json:"accounting_account_id"`
}

type UpdateProductRequest struct {
	Name                *string     `json:"name" validate:"omitempty,min=3,max=255"`
	SKU                 *string     `json:"sku"`
	Barcode             *string     `json:"barcode"`
	Description         *string     `json:"description"`
	CategoryID          *uuid.UUID  `json:"category_id"`
	UnitPrice           *float64    `json:"unit_price" validate:"omitempty,gte=0"`
	Cost                *float64    `json:"cost" validate:"omitempty,gte=0"`
	TaxRate             *float64    `json:"tax_rate" validate:"omitempty,gte=0,lte=100"`
	Unit                *string     `json:"unit"`
	MinStockLevel       *float64    `json:"min_stock_level" validate:"omitempty,gte=0"`
	MaxStockLevel       *float64    `json:"max_stock_level" validate:"omitempty,gte=0"`
	IsActive            *bool       `json:"is_active"`
	IsTrackInventory    *bool       `json:"is_track_inventory"`
	AllowNegativeStock  *bool       `json:"allow_negative_stock"`
	ImageURL            *string     `json:"image_url"`
	ProductType         *string     `json:"product_type"`
	AccountingAccountID *uuid.UUID  `json:"accounting_account_id"`
}

// ============================================================================
// CUSTOMERS
// ============================================================================

type CreateCustomerRequest struct {
	FirstName       string     `json:"first_name" validate:"required"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email" validate:"required,email"`
	Phone           string     `json:"phone"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	Gender          string     `json:"gender"`
	AddressLine1    string     `json:"address_line_1"`
	AddressLine2    string     `json:"address_line_2"`
	City            string     `json:"city"`
	State           string     `json:"state"`
	PostalCode      string     `json:"postal_code"`
	Country         string     `json:"country"`
	CustomerType    string     `json:"customer_type"`
	TaxID           string     `json:"tax_id"`
	PaymentTermDays int        `json:"payment_term_days"`
	CreditLimit     float64    `json:"credit_limit" validate:"gte=0"`
	IsActive        bool       `json:"is_active"`
}

type UpdateCustomerRequest struct {
	FirstName       *string    `json:"first_name"`
	LastName        *string    `json:"last_name"`
	Email           *string    `json:"email" validate:"omitempty,email"`
	Phone           *string    `json:"phone"`
	DateOfBirth     *time.Time `json:"date_of_birth"`
	Gender          *string    `json:"gender"`
	AddressLine1    *string    `json:"address_line_1"`
	AddressLine2    *string    `json:"address_line_2"`
	City            *string    `json:"city"`
	State           *string    `json:"state"`
	PostalCode      *string    `json:"postal_code"`
	Country         *string    `json:"country"`
	CustomerType    *string    `json:"customer_type"`
	TaxID           *string    `json:"tax_id"`
	PaymentTermDays *int       `json:"payment_term_days"`
	CreditLimit     *float64   `json:"credit_limit" validate:"omitempty,gte=0"`
	IsActive        *bool      `json:"is_active"`
}

// ============================================================================
// SUPPLIERS
// ============================================================================

type CreateSupplierRequest struct {
	Name            string  `json:"name" validate:"required,min=3,max=255"`
	ContactPerson   string  `json:"contact_person"`
	Email           string  `json:"email" validate:"required,email"`
	Phone           string  `json:"phone"`
	Website         string  `json:"website"`
	TaxID           string  `json:"tax_id"`
	RegistrationNo  string  `json:"registration_no"`
	AddressLine1    string  `json:"address_line_1"`
	AddressLine2    string  `json:"address_line_2"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	PostalCode      string  `json:"postal_code"`
	Country         string  `json:"country"`
	PaymentTermDays int     `json:"payment_term_days"`
	CreditLimit     float64 `json:"credit_limit" validate:"gte=0"`
	IsActive        bool    `json:"is_active"`
}

type UpdateSupplierRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=3,max=255"`
	ContactPerson   *string  `json:"contact_person"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Phone           *string  `json:"phone"`
	Website         *string  `json:"website"`
	TaxID           *string  `json:"tax_id"`
	RegistrationNo  *string  `json:"registration_no"`
	AddressLine1    *string  `json:"address_line_1"`
	AddressLine2    *string  `json:"address_line_2"`
	City            *string  `json:"city"`
	State           *string  `json:"state"`
	PostalCode      *string  `json:"postal_code"`
	Country         *string  `json:"country"`
	PaymentTermDays *int     `json:"payment_term_days"`
	CreditLimit     *float64 `json:"credit_limit" validate:"omitempty,gte=0"`
	IsActive        *bool    `json:"is_active"`
}

// ============================================================================
// CATEGORIES
// ============================================================================

type CreateCategoryRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id"`
	Color       string     `json:"color"`
	IconName    string     `json:"icon_name"`
	DisplayOrder int       `json:"display_order"`
	IsActive    bool       `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name         *string    `json:"name" validate:"omitempty,min=3,max=255"`
	Description  *string    `json:"description"`
	ParentID     *uuid.UUID `json:"parent_id"`
	Color        *string    `json:"color"`
	IconName     *string    `json:"icon_name"`
	DisplayOrder *int       `json:"display_order"`
	IsActive     *bool      `json:"is_active"`
}

// ============================================================================
// LOCATIONS
// ============================================================================

type CreateLocationRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=255"`
	Code         string `json:"code" validate:"required"`
	LocationType string `json:"location_type"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	Phone        string `json:"phone"`
	Email        string `json:"email" validate:"omitempty,email"`
	IsActive     bool   `json:"is_active"`
}

type UpdateLocationRequest struct {
	Name         *string `json:"name" validate:"omitempty,min=3,max=255"`
	Code         *string `json:"code"`
	LocationType *string `json:"location_type"`
	AddressLine1 *string `json:"address_line_1"`
	AddressLine2 *string `json:"address_line_2"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	PostalCode   *string `json:"postal_code"`
	Country      *string `json:"country"`
	Phone        *string `json:"phone"`
	Email        *string `json:"email" validate:"omitempty,email"`
	IsActive     *bool   `json:"is_active"`
}

// ============================================================================
// SALES
// ============================================================================

type CreateSaleRequest struct {
	SaleNumber      string          `json:"sale_number"`
	CustomerID      *uuid.UUID      `json:"customer_id"`
	LocationID      uuid.UUID       `json:"location_id" validate:"required"`
	SaleDate        time.Time       `json:"sale_date" validate:"required"`
	Subtotal        float64         `json:"subtotal" validate:"required,gte=0"`
	TaxAmount       float64         `json:"tax_amount" validate:"gte=0"`
	DiscountAmount  float64         `json:"discount_amount" validate:"gte=0"`
	TotalAmount     float64         `json:"total_amount" validate:"required,gte=0"`
	Status          string          `json:"status"`
	PaymentMethod   string          `json:"payment_method"`
	Notes           string          `json:"notes"`
	Items           []SaleItemInput `json:"items" validate:"required,min=1,dive"`
}

type SaleItemInput struct {
	ProductID      uuid.UUID `json:"product_id" validate:"required"`
	Quantity       float64   `json:"quantity" validate:"required,gt=0"`
	UnitPrice      float64   `json:"unit_price" validate:"required,gte=0"`
	TaxRate        float64   `json:"tax_rate" validate:"gte=0,lte=100"`
	DiscountAmount float64   `json:"discount_amount" validate:"gte=0"`
	LineTotal      float64   `json:"line_total" validate:"required,gte=0"`
}

// ============================================================================
// POSTING ENGINE
// ============================================================================

type PostDocumentRequest struct {
	DocumentType string    `json:"document_type" validate:"required"`
	DocumentID   uuid.UUID `json:"document_id" validate:"required"`
	Event        string    `json:"event" validate:"required"`
}

type PostDocumentResponse struct {
	JournalEntryID uuid.UUID `json:"journal_entry_id"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
}
