package order

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

// Repository handles database operations for Orders
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Orders repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Orders represents a orders entity
type Orders struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	OrderNumber string `json:"order_number" db:"order_number"`
	DisplayNumber *int64 `json:"display_number" db:"display_number"`
	OrderType *string `json:"order_type" db:"order_type"`
	TableId *uuid.UUID `json:"table_id" db:"table_id"`
	ReservationId *uuid.UUID `json:"reservation_id" db:"reservation_id"`
	Covers *int64 `json:"covers" db:"covers"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	WaiterId *uuid.UUID `json:"waiter_id" db:"waiter_id"`
	Status *string `json:"status" db:"status"`
	OrderDate *time.Time `json:"order_date" db:"order_date"`
	SubmittedAt *time.Time `json:"submitted_at" db:"submitted_at"`
	KitchenReceivedAt *time.Time `json:"kitchen_received_at" db:"kitchen_received_at"`
	ReadyAt *time.Time `json:"ready_at" db:"ready_at"`
	ServedAt *time.Time `json:"served_at" db:"served_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Subtotal *float64 `json:"subtotal" db:"subtotal"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	DiscountAmount *float64 `json:"discount_amount" db:"discount_amount"`
	ServiceCharge *float64 `json:"service_charge" db:"service_charge"`
	TotalAmount *float64 `json:"total_amount" db:"total_amount"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	ShiftId *uuid.UUID `json:"shift_id" db:"shift_id"`
	CustomerNotes *string `json:"customer_notes" db:"customer_notes"`
	KitchenNotes *string `json:"kitchen_notes" db:"kitchen_notes"`
	InternalNotes *string `json:"internal_notes" db:"internal_notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'draft', *string `json:"'draft'," db:"'draft',"`
	'served', *string `json:"'served'," db:"'served',"`
	'dineIn', *string `json:"'dine_in'," db:"'dine_in',"`
	Subtotal *string `json:"subtotal" db:"subtotal"`
	DiscountAmount *string `json:"discount_amount" db:"discount_amount"`
}

// Create inserts a new orders record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Orders) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "orders", duration, nil)
	}()

	query := `
		INSERT INTO orders (
			, organization_id
			, location_id
			, order_number
			, display_number
			, order_type
			, table_id
			, reservation_id
			, covers
			, customer_id
			, waiter_id
			, status
			, order_date
			, submitted_at
			, kitchen_received_at
			, ready_at
			, served_at
			, completed_at
			, subtotal
			, tax_amount
			, discount_amount
			, service_charge
			, total_amount
			, sale_id
			, shift_id
			, customer_notes
			, kitchen_notes
			, internal_notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'draft',
			, 'served',
			, 'dine_in',
			, subtotal
			, discount_amount
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
			, $35
			, $36
			, $37
			, $38
			, $39
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.OrderNumber,
		entity.DisplayNumber,
		entity.OrderType,
		entity.TableId,
		entity.ReservationId,
		entity.Covers,
		entity.CustomerId,
		entity.WaiterId,
		entity.Status,
		entity.OrderDate,
		entity.SubmittedAt,
		entity.KitchenReceivedAt,
		entity.ReadyAt,
		entity.ServedAt,
		entity.CompletedAt,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.ServiceCharge,
		entity.TotalAmount,
		entity.SaleId,
		entity.ShiftId,
		entity.CustomerNotes,
		entity.KitchenNotes,
		entity.InternalNotes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'draft',,
		entity.'served',,
		entity.'dineIn',,
		entity.Subtotal,
		entity.DiscountAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create orders", zap.Error(err))
		return fmt.Errorf("failed to create orders: %w", err)
	}

	r.logger.Info("created orders",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a orders by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Orders, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "orders", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, order_number
			, display_number
			, order_type
			, table_id
			, reservation_id
			, covers
			, customer_id
			, waiter_id
			, order_date
			, submitted_at
			, kitchen_received_at
			, ready_at
			, served_at
			, completed_at
			, subtotal
			, tax_amount
			, discount_amount
			, service_charge
			, total_amount
			, sale_id
			, shift_id
			, customer_notes
			, kitchen_notes
			, internal_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'draft',
			, 'served',
			, 'dine_in',
			, subtotal
			, discount_amount
		FROM orders
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Orders
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.OrderNumber,
		&entity.DisplayNumber,
		&entity.OrderType,
		&entity.TableId,
		&entity.ReservationId,
		&entity.Covers,
		&entity.CustomerId,
		&entity.WaiterId,
		&entity.Status,
		&entity.OrderDate,
		&entity.SubmittedAt,
		&entity.KitchenReceivedAt,
		&entity.ReadyAt,
		&entity.ServedAt,
		&entity.CompletedAt,
		&entity.Subtotal,
		&entity.TaxAmount,
		&entity.DiscountAmount,
		&entity.ServiceCharge,
		&entity.TotalAmount,
		&entity.SaleId,
		&entity.ShiftId,
		&entity.CustomerNotes,
		&entity.KitchenNotes,
		&entity.InternalNotes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'draft',,
		&entity.'served',,
		&entity.'dineIn',,
		&entity.Subtotal,
		&entity.DiscountAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("orders not found")
	}

	if err != nil {
		r.logger.Error("failed to get orders", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of orders records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Orders, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "orders", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, order_number
			, display_number
			, order_type
			, table_id
			, reservation_id
			, covers
			, customer_id
			, waiter_id
			, order_date
			, submitted_at
			, kitchen_received_at
			, ready_at
			, served_at
			, completed_at
			, subtotal
			, tax_amount
			, discount_amount
			, service_charge
			, total_amount
			, sale_id
			, shift_id
			, customer_notes
			, kitchen_notes
			, internal_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'draft',
			, 'served',
			, 'dine_in',
			, subtotal
			, discount_amount
		FROM orders
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list orders", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var entities []*Orders
	for rows.Next() {
		var entity Orders
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.OrderNumber,
			&entity.DisplayNumber,
			&entity.OrderType,
			&entity.TableId,
			&entity.ReservationId,
			&entity.Covers,
			&entity.CustomerId,
			&entity.WaiterId,
			&entity.Status,
			&entity.OrderDate,
			&entity.SubmittedAt,
			&entity.KitchenReceivedAt,
			&entity.ReadyAt,
			&entity.ServedAt,
			&entity.CompletedAt,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.ServiceCharge,
			&entity.TotalAmount,
			&entity.SaleId,
			&entity.ShiftId,
			&entity.CustomerNotes,
			&entity.KitchenNotes,
			&entity.InternalNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'draft',,
			&entity.'served',,
			&entity.'dineIn',,
			&entity.Subtotal,
			&entity.DiscountAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan orders: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating orders rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing orders record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Orders) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "orders", duration, nil)
	}()

	query := `
		UPDATE orders
		SET
			, organization_id = $2
			, location_id = $3
			, order_number = $4
			, display_number = $5
			, order_type = $6
			, table_id = $7
			, reservation_id = $8
			, covers = $9
			, customer_id = $10
			, waiter_id = $11
			, status = $12
			, order_date = $13
			, submitted_at = $14
			, kitchen_received_at = $15
			, ready_at = $16
			, served_at = $17
			, completed_at = $18
			, subtotal = $19
			, tax_amount = $20
			, discount_amount = $21
			, service_charge = $22
			, total_amount = $23
			, sale_id = $24
			, shift_id = $25
			, customer_notes = $26
			, kitchen_notes = $27
			, internal_notes = $28
			, metadata = $29
			, updated_at = $31
			, created_by = $32
			, updated_by = $33
			, deleted_at = $34
			, 'draft', = $35
			, 'served', = $36
			, 'dine_in', = $37
			, subtotal = $38
			, discount_amount = $39
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $40
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.OrderNumber,
		entity.DisplayNumber,
		entity.OrderType,
		entity.TableId,
		entity.ReservationId,
		entity.Covers,
		entity.CustomerId,
		entity.WaiterId,
		entity.Status,
		entity.OrderDate,
		entity.SubmittedAt,
		entity.KitchenReceivedAt,
		entity.ReadyAt,
		entity.ServedAt,
		entity.CompletedAt,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.ServiceCharge,
		entity.TotalAmount,
		entity.SaleId,
		entity.ShiftId,
		entity.CustomerNotes,
		entity.KitchenNotes,
		entity.InternalNotes,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'draft',,
		entity.'served',,
		entity.'dineIn',,
		entity.Subtotal,
		entity.DiscountAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update orders", zap.Error(err))
		return fmt.Errorf("failed to update orders: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("orders not found or already deleted")
	}

	r.logger.Info("updated orders",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a orders record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "orders", duration, nil)
	}()

	query := `
		UPDATE orders
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete orders", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete orders: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("orders not found or already deleted")
	}

	r.logger.Info("deleted orders", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves orders records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Orders, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "orders", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count orders records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, order_number
			, display_number
			, order_type
			, table_id
			, reservation_id
			, covers
			, customer_id
			, waiter_id
			, status
			, order_date
			, submitted_at
			, kitchen_received_at
			, ready_at
			, served_at
			, completed_at
			, subtotal
			, tax_amount
			, discount_amount
			, service_charge
			, total_amount
			, sale_id
			, shift_id
			, customer_notes
			, kitchen_notes
			, internal_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'draft',
			, 'served',
			, 'dine_in',
			, subtotal
			, discount_amount
		FROM orders
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list orders by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var entities []*Orders
	for rows.Next() {
		var entity Orders
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.OrderNumber,
			&entity.DisplayNumber,
			&entity.OrderType,
			&entity.TableId,
			&entity.ReservationId,
			&entity.Covers,
			&entity.CustomerId,
			&entity.WaiterId,
			&entity.Status,
			&entity.OrderDate,
			&entity.SubmittedAt,
			&entity.KitchenReceivedAt,
			&entity.ReadyAt,
			&entity.ServedAt,
			&entity.CompletedAt,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.ServiceCharge,
			&entity.TotalAmount,
			&entity.SaleId,
			&entity.ShiftId,
			&entity.CustomerNotes,
			&entity.KitchenNotes,
			&entity.InternalNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'draft',,
			&entity.'served',,
			&entity.'dineIn',,
			&entity.Subtotal,
			&entity.DiscountAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan orders: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

