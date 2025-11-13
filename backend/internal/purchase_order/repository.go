package purchase_order

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

// Repository handles database operations for PurchaseOrders
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PurchaseOrders repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PurchaseOrders represents a purchase_orders entity
type PurchaseOrders struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PoNumber string `json:"po_number" db:"po_number"`
	SupplierId uuid.UUID `json:"supplier_id" db:"supplier_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	OrderDate time.Time `json:"order_date" db:"order_date"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" db:"expected_delivery_date"`
	ActualDeliveryDate *time.Time `json:"actual_delivery_date" db:"actual_delivery_date"`
	SubtotalAmount *float64 `json:"subtotal_amount" db:"subtotal_amount"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	ShippingAmount *float64 `json:"shipping_amount" db:"shipping_amount"`
	TotalAmount *float64 `json:"total_amount" db:"total_amount"`
	PaymentTerms *string `json:"payment_terms" db:"payment_terms"`
	PaymentDueDate *time.Time `json:"payment_due_date" db:"payment_due_date"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new purchase_orders record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PurchaseOrders) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "purchase_orders", duration, nil)
	}()

	query := `
		INSERT INTO purchase_orders (
			, organization_id
			, po_number
			, supplier_id
			, location_id
			, order_date
			, expected_delivery_date
			, actual_delivery_date
			, subtotal_amount
			, tax_amount
			, shipping_amount
			, total_amount
			, payment_terms
			, payment_due_date
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
			, updated_by
			, deleted_at
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
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PoNumber,
		entity.SupplierId,
		entity.LocationId,
		entity.OrderDate,
		entity.ExpectedDeliveryDate,
		entity.ActualDeliveryDate,
		entity.SubtotalAmount,
		entity.TaxAmount,
		entity.ShippingAmount,
		entity.TotalAmount,
		entity.PaymentTerms,
		entity.PaymentDueDate,
		entity.Status,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create purchase_orders", zap.Error(err))
		return fmt.Errorf("failed to create purchase_orders: %w", err)
	}

	r.logger.Info("created purchase_orders",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a purchase_orders by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PurchaseOrders, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_orders", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, po_number
			, supplier_id
			, location_id
			, order_date
			, expected_delivery_date
			, actual_delivery_date
			, subtotal_amount
			, tax_amount
			, shipping_amount
			, total_amount
			, payment_terms
			, payment_due_date
			, approved_by
			, approved_at
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_orders
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PurchaseOrders
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PoNumber,
		&entity.SupplierId,
		&entity.LocationId,
		&entity.OrderDate,
		&entity.ExpectedDeliveryDate,
		&entity.ActualDeliveryDate,
		&entity.SubtotalAmount,
		&entity.TaxAmount,
		&entity.ShippingAmount,
		&entity.TotalAmount,
		&entity.PaymentTerms,
		&entity.PaymentDueDate,
		&entity.Status,
		&entity.Status,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("purchase_orders not found")
	}

	if err != nil {
		r.logger.Error("failed to get purchase_orders", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get purchase_orders: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of purchase_orders records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PurchaseOrders, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_orders", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM purchase_orders
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase_orders records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, po_number
			, supplier_id
			, location_id
			, order_date
			, expected_delivery_date
			, actual_delivery_date
			, subtotal_amount
			, tax_amount
			, shipping_amount
			, total_amount
			, payment_terms
			, payment_due_date
			, approved_by
			, approved_at
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_orders
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list purchase_orders", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list purchase_orders: %w", err)
	}
	defer rows.Close()

	var entities []*PurchaseOrders
	for rows.Next() {
		var entity PurchaseOrders
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PoNumber,
			&entity.SupplierId,
			&entity.LocationId,
			&entity.OrderDate,
			&entity.ExpectedDeliveryDate,
			&entity.ActualDeliveryDate,
			&entity.SubtotalAmount,
			&entity.TaxAmount,
			&entity.ShippingAmount,
			&entity.TotalAmount,
			&entity.PaymentTerms,
			&entity.PaymentDueDate,
			&entity.Status,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase_orders: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating purchase_orders rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing purchase_orders record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PurchaseOrders) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "purchase_orders", duration, nil)
	}()

	query := `
		UPDATE purchase_orders
		SET
			, organization_id = $2
			, po_number = $3
			, supplier_id = $4
			, location_id = $5
			, order_date = $6
			, expected_delivery_date = $7
			, actual_delivery_date = $8
			, subtotal_amount = $9
			, tax_amount = $10
			, shipping_amount = $11
			, total_amount = $12
			, payment_terms = $13
			, payment_due_date = $14
			, status = $15
			, approved_by = $17
			, approved_at = $18
			, notes = $19
			, created_by = $20
			, updated_by = $21
			, updated_at = $23
			, deleted_at = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PoNumber,
		entity.SupplierId,
		entity.LocationId,
		entity.OrderDate,
		entity.ExpectedDeliveryDate,
		entity.ActualDeliveryDate,
		entity.SubtotalAmount,
		entity.TaxAmount,
		entity.ShippingAmount,
		entity.TotalAmount,
		entity.PaymentTerms,
		entity.PaymentDueDate,
		entity.Status,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update purchase_orders", zap.Error(err))
		return fmt.Errorf("failed to update purchase_orders: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("purchase_orders not found or already deleted")
	}

	r.logger.Info("updated purchase_orders",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a purchase_orders record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "purchase_orders", duration, nil)
	}()

	query := `
		UPDATE purchase_orders
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete purchase_orders", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete purchase_orders: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("purchase_orders not found or already deleted")
	}

	r.logger.Info("deleted purchase_orders", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves purchase_orders records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PurchaseOrders, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "purchase_orders", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM purchase_orders
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count purchase_orders records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, po_number
			, supplier_id
			, location_id
			, order_date
			, expected_delivery_date
			, actual_delivery_date
			, subtotal_amount
			, tax_amount
			, shipping_amount
			, total_amount
			, payment_terms
			, payment_due_date
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM purchase_orders
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list purchase_orders by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list purchase_orders: %w", err)
	}
	defer rows.Close()

	var entities []*PurchaseOrders
	for rows.Next() {
		var entity PurchaseOrders
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PoNumber,
			&entity.SupplierId,
			&entity.LocationId,
			&entity.OrderDate,
			&entity.ExpectedDeliveryDate,
			&entity.ActualDeliveryDate,
			&entity.SubtotalAmount,
			&entity.TaxAmount,
			&entity.ShippingAmount,
			&entity.TotalAmount,
			&entity.PaymentTerms,
			&entity.PaymentDueDate,
			&entity.Status,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase_orders: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

