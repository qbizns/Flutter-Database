package supplier

import (
	"encoding/json"
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

// Repository handles database operations for Suppliers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Suppliers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Suppliers represents a suppliers entity
type Suppliers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SupplierCode string `json:"supplier_code" db:"supplier_code"`
	Name string `json:"name" db:"name"`
	ContactPerson *string `json:"contact_person" db:"contact_person"`
	Email *string `json:"email" db:"email"`
	Phone *string `json:"phone" db:"phone"`
	Address *string `json:"address" db:"address"`
	City *string `json:"city" db:"city"`
	State *string `json:"state" db:"state"`
	Country *string `json:"country" db:"country"`
	PostalCode *string `json:"postal_code" db:"postal_code"`
	TaxNumber *string `json:"tax_number" db:"tax_number"`
	PaymentTerms *string `json:"payment_terms" db:"payment_terms"`
	CreditLimit *float64 `json:"credit_limit" db:"credit_limit"`
	OutstandingBalance *float64 `json:"outstanding_balance" db:"outstanding_balance"`
	TotalPurchases *float64 `json:"total_purchases" db:"total_purchases"`
	TotalOrders *int64 `json:"total_orders" db:"total_orders"`
	LastOrderDate *time.Time `json:"last_order_date" db:"last_order_date"`
	Status *string `json:"status" db:"status"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new suppliers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Suppliers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "suppliers", duration, nil)
	}()

	query := `
		INSERT INTO suppliers (
			, organization_id
			, supplier_code
			, name
			, contact_person
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, tax_number
			, payment_terms
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_order_date
			, status
			, notes
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
			, $25
			, $26
			, $27
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SupplierCode,
		entity.Name,
		entity.ContactPerson,
		entity.Email,
		entity.Phone,
		entity.Address,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.TaxNumber,
		entity.PaymentTerms,
		entity.CreditLimit,
		entity.OutstandingBalance,
		entity.TotalPurchases,
		entity.TotalOrders,
		entity.LastOrderDate,
		entity.Status,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create suppliers", zap.Error(err))
		return fmt.Errorf("failed to create suppliers: %w", err)
	}

	r.logger.Info("created suppliers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a suppliers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Suppliers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "suppliers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, supplier_code
			, name
			, contact_person
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, tax_number
			, payment_terms
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_order_date
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM suppliers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Suppliers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SupplierCode,
		&entity.Name,
		&entity.ContactPerson,
		&entity.Email,
		&entity.Phone,
		&entity.Address,
		&entity.City,
		&entity.State,
		&entity.Country,
		&entity.PostalCode,
		&entity.TaxNumber,
		&entity.PaymentTerms,
		&entity.CreditLimit,
		&entity.OutstandingBalance,
		&entity.TotalPurchases,
		&entity.TotalOrders,
		&entity.LastOrderDate,
		&entity.Status,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("suppliers not found")
	}

	if err != nil {
		r.logger.Error("failed to get suppliers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of suppliers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Suppliers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "suppliers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM suppliers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, supplier_code
			, name
			, contact_person
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, tax_number
			, payment_terms
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_order_date
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM suppliers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list suppliers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	var entities []*Suppliers
	for rows.Next() {
		var entity Suppliers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SupplierCode,
			&entity.Name,
			&entity.ContactPerson,
			&entity.Email,
			&entity.Phone,
			&entity.Address,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.TaxNumber,
			&entity.PaymentTerms,
			&entity.CreditLimit,
			&entity.OutstandingBalance,
			&entity.TotalPurchases,
			&entity.TotalOrders,
			&entity.LastOrderDate,
			&entity.Status,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan suppliers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating suppliers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing suppliers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Suppliers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "suppliers", duration, nil)
	}()

	query := `
		UPDATE suppliers
		SET
			, organization_id = $2
			, supplier_code = $3
			, name = $4
			, contact_person = $5
			, email = $6
			, phone = $7
			, address = $8
			, city = $9
			, state = $10
			, country = $11
			, postal_code = $12
			, tax_number = $13
			, payment_terms = $14
			, credit_limit = $15
			, outstanding_balance = $16
			, total_purchases = $17
			, total_orders = $18
			, last_order_date = $19
			, status = $20
			, notes = $21
			, metadata = $22
			, updated_at = $24
			, deleted_at = $25
			, created_by = $26
			, updated_by = $27
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $28
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SupplierCode,
		entity.Name,
		entity.ContactPerson,
		entity.Email,
		entity.Phone,
		entity.Address,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.TaxNumber,
		entity.PaymentTerms,
		entity.CreditLimit,
		entity.OutstandingBalance,
		entity.TotalPurchases,
		entity.TotalOrders,
		entity.LastOrderDate,
		entity.Status,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update suppliers", zap.Error(err))
		return fmt.Errorf("failed to update suppliers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("suppliers not found or already deleted")
	}

	r.logger.Info("updated suppliers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a suppliers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "suppliers", duration, nil)
	}()

	query := `
		UPDATE suppliers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete suppliers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete suppliers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("suppliers not found or already deleted")
	}

	r.logger.Info("deleted suppliers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves suppliers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Suppliers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "suppliers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM suppliers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count suppliers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, supplier_code
			, name
			, contact_person
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, tax_number
			, payment_terms
			, credit_limit
			, outstanding_balance
			, total_purchases
			, total_orders
			, last_order_date
			, status
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM suppliers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list suppliers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list suppliers: %w", err)
	}
	defer rows.Close()

	var entities []*Suppliers
	for rows.Next() {
		var entity Suppliers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SupplierCode,
			&entity.Name,
			&entity.ContactPerson,
			&entity.Email,
			&entity.Phone,
			&entity.Address,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.TaxNumber,
			&entity.PaymentTerms,
			&entity.CreditLimit,
			&entity.OutstandingBalance,
			&entity.TotalPurchases,
			&entity.TotalOrders,
			&entity.LastOrderDate,
			&entity.Status,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan suppliers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

