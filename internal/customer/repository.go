package customer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for Customers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Customers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Customers represents a customers entity
type Customers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CustomerCode *string `json:"customer_code" db:"customer_code"`
	FirstName *string `json:"first_name" db:"first_name"`
	LastName *string `json:"last_name" db:"last_name"`
	CompanyName *string `json:"company_name" db:"company_name"`
	Email *string `json:"email" db:"email"`
	Phone *string `json:"phone" db:"phone"`
	AlternatePhone *string `json:"alternate_phone" db:"alternate_phone"`
	AddressLine1 *string `json:"address_line1" db:"address_line1"`
	AddressLine2 *string `json:"address_line2" db:"address_line2"`
	City *string `json:"city" db:"city"`
	State *string `json:"state" db:"state"`
	Country *string `json:"country" db:"country"`
	PostalCode *string `json:"postal_code" db:"postal_code"`
	DateOfBirth *time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender *string `json:"gender" db:"gender"`
	TaxNumber *string `json:"tax_number" db:"tax_number"`
	LoyaltyPoints *int64 `json:"loyalty_points" db:"loyalty_points"`
	LoyaltyTier *string `json:"loyalty_tier" db:"loyalty_tier"`
	CreditLimit *float64 `json:"credit_limit" db:"credit_limit"`
	OutstandingBalance *float64 `json:"outstanding_balance" db:"outstanding_balance"`
	TotalPurchases *float64 `json:"total_purchases" db:"total_purchases"`
	TotalOrders *int64 `json:"total_orders" db:"total_orders"`
	LastPurchaseAt *time.Time `json:"last_purchase_at" db:"last_purchase_at"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	CustomFields json.RawMessage `json:"custom_fields" db:"custom_fields"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new customers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Customers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "customers", duration, nil)
	}()

	query := `
		INSERT INTO customers (
			, organization_id
			, customer_code
			, first_name
			, last_name
			, company_name
			, email
			, phone
			, alternate_phone
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, date_of_birth
			, gender
			, tax_number
			, loyalty_points
			, loyalty_tier
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_purchase_at
			, is_active
			, notes
			, custom_fields
			, metadata
			, deleted_at
			, created_by
			, updated_by
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $32
			, $33
			, $34
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CustomerCode,
		entity.FirstName,
		entity.LastName,
		entity.CompanyName,
		entity.Email,
		entity.Phone,
		entity.AlternatePhone,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.DateOfBirth,
		entity.Gender,
		entity.TaxNumber,
		entity.LoyaltyPoints,
		entity.LoyaltyTier,
		entity.CreditLimit,
		entity.OutstandingBalance,
		entity.TotalPurchases,
		entity.TotalOrders,
		entity.LastPurchaseAt,
		entity.IsActive,
		entity.Notes,
		entity.CustomFields,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create customers", zap.Error(err))
		return fmt.Errorf("failed to create customers: %w", err)
	}

	r.logger.Info("created customers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a customers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Customers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, customer_code
			, first_name
			, last_name
			, company_name
			, email
			, phone
			, alternate_phone
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, date_of_birth
			, gender
			, tax_number
			, loyalty_points
			, loyalty_tier
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_purchase_at
			, is_active
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM customers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Customers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CustomerCode,
		&entity.FirstName,
		&entity.LastName,
		&entity.CompanyName,
		&entity.Email,
		&entity.Phone,
		&entity.AlternatePhone,
		&entity.AddressLine1,
		&entity.AddressLine2,
		&entity.City,
		&entity.State,
		&entity.Country,
		&entity.PostalCode,
		&entity.DateOfBirth,
		&entity.Gender,
		&entity.TaxNumber,
		&entity.LoyaltyPoints,
		&entity.LoyaltyTier,
		&entity.CreditLimit,
		&entity.OutstandingBalance,
		&entity.TotalPurchases,
		&entity.TotalOrders,
		&entity.LastPurchaseAt,
		&entity.IsActive,
		&entity.Notes,
		&entity.CustomFields,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("customers not found")
	}

	if err != nil {
		r.logger.Error("failed to get customers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of customers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Customers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_code
			, first_name
			, last_name
			, company_name
			, email
			, phone
			, alternate_phone
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, date_of_birth
			, gender
			, tax_number
			, loyalty_points
			, loyalty_tier
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_purchase_at
			, is_active
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM customers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var entities []*Customers
	for rows.Next() {
		var entity Customers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerCode,
			&entity.FirstName,
			&entity.LastName,
			&entity.CompanyName,
			&entity.Email,
			&entity.Phone,
			&entity.AlternatePhone,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.DateOfBirth,
			&entity.Gender,
			&entity.TaxNumber,
			&entity.LoyaltyPoints,
			&entity.LoyaltyTier,
			&entity.CreditLimit,
			&entity.OutstandingBalance,
			&entity.TotalPurchases,
			&entity.TotalOrders,
			&entity.LastPurchaseAt,
			&entity.IsActive,
			&entity.Notes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating customers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing customers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Customers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customers", duration, nil)
	}()

	query := `
		UPDATE customers
		SET
			, organization_id = $2
			, customer_code = $3
			, first_name = $4
			, last_name = $5
			, company_name = $6
			, email = $7
			, phone = $8
			, alternate_phone = $9
			, address_line1 = $10
			, address_line2 = $11
			, city = $12
			, state = $13
			, country = $14
			, postal_code = $15
			, date_of_birth = $16
			, gender = $17
			, tax_number = $18
			, loyalty_points = $19
			, loyalty_tier = $20
			, credit_limit = $21
			, outstanding_balance = $22
			, total_purchases = $23
			, total_orders = $24
			, last_purchase_at = $25
			, is_active = $26
			, notes = $27
			, custom_fields = $28
			, metadata = $29
			, updated_at = $31
			, deleted_at = $32
			, created_by = $33
			, updated_by = $34
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $35
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CustomerCode,
		entity.FirstName,
		entity.LastName,
		entity.CompanyName,
		entity.Email,
		entity.Phone,
		entity.AlternatePhone,
		entity.AddressLine1,
		entity.AddressLine2,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.DateOfBirth,
		entity.Gender,
		entity.TaxNumber,
		entity.LoyaltyPoints,
		entity.LoyaltyTier,
		entity.CreditLimit,
		entity.OutstandingBalance,
		entity.TotalPurchases,
		entity.TotalOrders,
		entity.LastPurchaseAt,
		entity.IsActive,
		entity.Notes,
		entity.CustomFields,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update customers", zap.Error(err))
		return fmt.Errorf("failed to update customers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customers not found or already deleted")
	}

	r.logger.Info("updated customers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a customers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "customers", duration, nil)
	}()

	query := `
		UPDATE customers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete customers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete customers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("customers not found or already deleted")
	}

	r.logger.Info("deleted customers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves customers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Customers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "customers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM customers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, customer_code
			, first_name
			, last_name
			, company_name
			, email
			, phone
			, alternate_phone
			, address_line1
			, address_line2
			, city
			, state
			, country
			, postal_code
			, date_of_birth
			, gender
			, tax_number
			, loyalty_points
			, loyalty_tier
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_purchase_at
			, is_active
			, notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM customers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list customers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var entities []*Customers
	for rows.Next() {
		var entity Customers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CustomerCode,
			&entity.FirstName,
			&entity.LastName,
			&entity.CompanyName,
			&entity.Email,
			&entity.Phone,
			&entity.AlternatePhone,
			&entity.AddressLine1,
			&entity.AddressLine2,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.DateOfBirth,
			&entity.Gender,
			&entity.TaxNumber,
			&entity.LoyaltyPoints,
			&entity.LoyaltyTier,
			&entity.CreditLimit,
			&entity.OutstandingBalance,
			&entity.TotalPurchases,
			&entity.TotalOrders,
			&entity.LastPurchaseAt,
			&entity.IsActive,
			&entity.Notes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

