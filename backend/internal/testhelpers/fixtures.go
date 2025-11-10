package testhelpers

import (
	"time"

	"github.com/google/uuid"
)

// Fixtures provides common test data
type Fixtures struct{}

// NewFixtures creates a new fixtures instance
func NewFixtures() *Fixtures {
	return &Fixtures{}
}

// OrganizationData returns test organization data
func (f *Fixtures) OrganizationData() map[string]interface{} {
	return map[string]interface{}{
		"id":         uuid.New(),
		"name":       "Test Organization",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
}

// UserData returns test user data
func (f *Fixtures) UserData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"id":              uuid.New(),
		"organization_id": orgID,
		"email":           "test@example.com",
		"password_hash":   "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password123
		"first_name":      "Test",
		"last_name":       "User",
		"is_active":       true,
		"created_at":      time.Now(),
		"updated_at":      time.Now(),
	}
}

// CustomerData returns test customer data
func (f *Fixtures) CustomerData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": orgID,
		"first_name":      "John",
		"last_name":       "Doe",
		"email":           "john.doe@example.com",
		"phone":           "+1234567890",
		"address_line1":   "123 Main St",
		"city":            "New York",
		"state":           "NY",
		"postal_code":     "10001",
		"country":         "US",
		"customer_type":   "regular",
		"is_active":       true,
	}
}

// ProductData returns test product data
func (f *Fixtures) ProductData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":       orgID,
		"name":                  "Test Product",
		"sku":                   "TEST-001",
		"barcode":               "1234567890",
		"description":           "A test product",
		"unit_price":            9.99,
		"cost":                  5.00,
		"tax_rate":              0.08,
		"unit":                  "piece",
		"min_stock_level":       10,
		"max_stock_level":       100,
		"is_active":             true,
		"is_track_inventory":    true,
		"allow_negative_stock":  false,
		"product_type":          "simple",
	}
}

// CategoryData returns test category data
func (f *Fixtures) CategoryData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": orgID,
		"name":            "Test Category",
		"description":     "A test category",
		"color":           "#FF5733",
		"icon_name":       "category",
		"display_order":   1,
		"is_active":       true,
	}
}

// LocationData returns test location data
func (f *Fixtures) LocationData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": orgID,
		"name":            "Main Warehouse",
		"code":            "WH-001",
		"location_type":   "warehouse",
		"address_line1":   "456 Storage St",
		"city":            "Los Angeles",
		"state":           "CA",
		"postal_code":     "90001",
		"country":         "US",
		"phone":           "+1234567891",
		"email":           "warehouse@example.com",
		"is_active":       true,
	}
}

// SupplierData returns test supplier data
func (f *Fixtures) SupplierData(orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id": orgID,
		"name":            "Test Supplier Inc",
		"contact_person":  "Jane Smith",
		"email":           "supplier@example.com",
		"phone":           "+1234567892",
		"address":         "789 Supply Ave",
		"city":            "Chicago",
		"state":           "IL",
		"country":         "US",
		"postal_code":     "60601",
		"tax_number":      "TAX123456",
		"credit_limit":    10000.00,
		"status":          "active",
	}
}

// SaleData returns test sale data
func (f *Fixtures) SaleData(orgID, customerID, cashierID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":    orgID,
		"sale_number":        "SALE-TEST-001",
		"transaction_type":   "sale",
		"customer_id":        customerID,
		"cashier_id":         cashierID,
		"subtotal":           100.00,
		"tax_amount":         8.00,
		"discount_amount":    0.00,
		"total_amount":       108.00,
		"paid_amount":        108.00,
		"change_amount":      0.00,
		"outstanding_amount": 0.00,
		"payment_status":     "paid",
		"transaction_date":   time.Now(),
		"notes":              "Test sale",
	}
}

// JWTClaims returns test JWT claims
func (f *Fixtures) JWTClaims(userID, orgID uuid.UUID) map[string]interface{} {
	return map[string]interface{}{
		"user_id":         userID.String(),
		"organization_id": orgID.String(),
		"email":           "test@example.com",
		"exp":             time.Now().Add(1 * time.Hour).Unix(),
		"iat":             time.Now().Unix(),
	}
}
