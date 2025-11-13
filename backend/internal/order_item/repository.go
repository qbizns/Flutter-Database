package order_item

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

// Repository handles database operations for OrderItems
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new OrderItems repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// OrderItems represents a order_items entity
type OrderItems struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	OrderId uuid.UUID `json:"order_id" db:"order_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	ProductVariantId *uuid.UUID `json:"product_variant_id" db:"product_variant_id"`
	ItemName string `json:"item_name" db:"item_name"`
	Quantity float64 `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	CourseId *uuid.UUID `json:"course_id" db:"course_id"`
	CoursePosition *int64 `json:"course_position" db:"course_position"`
	FireTime *time.Time `json:"fire_time" db:"fire_time"`
	KitchenStationId *uuid.UUID `json:"kitchen_station_id" db:"kitchen_station_id"`
	KitchenTicketId *uuid.UUID `json:"kitchen_ticket_id" db:"kitchen_ticket_id"`
	Status *string `json:"status" db:"status"`
	FiredAt *time.Time `json:"fired_at" db:"fired_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at" db:"acknowledged_at"`
	StartedPreparingAt *time.Time `json:"started_preparing_at" db:"started_preparing_at"`
	ReadyAt *time.Time `json:"ready_at" db:"ready_at"`
	ServedAt *time.Time `json:"served_at" db:"served_at"`
	ModifiersTotal *float64 `json:"modifiers_total" db:"modifiers_total"`
	DiscountAmount *float64 `json:"discount_amount" db:"discount_amount"`
	LineTotal float64 `json:"line_total" db:"line_total"`
	SpecialInstructions *string `json:"special_instructions" db:"special_instructions"`
	CustomerNotes *string `json:"customer_notes" db:"customer_notes"`
	KitchenNotes *string `json:"kitchen_notes" db:"kitchen_notes"`
	SeatNumber *int64 `json:"seat_number" db:"seat_number"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'pending', *string `json:"'pending'," db:"'pending',"`
	'served', *string `json:"'served'," db:"'served',"`
	UnitPrice *string `json:"unit_price" db:"unit_price"`
	DiscountAmount *string `json:"discount_amount" db:"discount_amount"`
}

// Create inserts a new order_items record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *OrderItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "order_items", duration, nil)
	}()

	query := `
		INSERT INTO order_items (
			, organization_id
			, order_id
			, product_id
			, product_variant_id
			, item_name
			, quantity
			, unit_price
			, course_id
			, course_position
			, fire_time
			, kitchen_station_id
			, kitchen_ticket_id
			, status
			, fired_at
			, acknowledged_at
			, started_preparing_at
			, ready_at
			, served_at
			, modifiers_total
			, discount_amount
			, line_total
			, special_instructions
			, customer_notes
			, kitchen_notes
			, seat_number
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'pending',
			, 'served',
			, unit_price
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
			, $30
			, $31
			, $32
			, $33
			, $34
			, $35
			, $36
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ItemName,
		entity.Quantity,
		entity.UnitPrice,
		entity.CourseId,
		entity.CoursePosition,
		entity.FireTime,
		entity.KitchenStationId,
		entity.KitchenTicketId,
		entity.Status,
		entity.FiredAt,
		entity.AcknowledgedAt,
		entity.StartedPreparingAt,
		entity.ReadyAt,
		entity.ServedAt,
		entity.ModifiersTotal,
		entity.DiscountAmount,
		entity.LineTotal,
		entity.SpecialInstructions,
		entity.CustomerNotes,
		entity.KitchenNotes,
		entity.SeatNumber,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'pending',,
		entity.'served',,
		entity.UnitPrice,
		entity.DiscountAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create order_items", zap.Error(err))
		return fmt.Errorf("failed to create order_items: %w", err)
	}

	r.logger.Info("created order_items",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a order_items by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*OrderItems, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_items", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, order_id
			, product_id
			, product_variant_id
			, item_name
			, quantity
			, unit_price
			, course_id
			, course_position
			, fire_time
			, kitchen_station_id
			, kitchen_ticket_id
			, status
			, fired_at
			, acknowledged_at
			, started_preparing_at
			, ready_at
			, served_at
			, modifiers_total
			, discount_amount
			, line_total
			, special_instructions
			, customer_notes
			, kitchen_notes
			, seat_number
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'pending',
			, 'served',
			, unit_price
			, discount_amount
		FROM order_items
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity OrderItems
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.OrderId,
		&entity.ProductId,
		&entity.ProductVariantId,
		&entity.ItemName,
		&entity.Quantity,
		&entity.UnitPrice,
		&entity.CourseId,
		&entity.CoursePosition,
		&entity.FireTime,
		&entity.KitchenStationId,
		&entity.KitchenTicketId,
		&entity.Status,
		&entity.FiredAt,
		&entity.AcknowledgedAt,
		&entity.StartedPreparingAt,
		&entity.ReadyAt,
		&entity.ServedAt,
		&entity.ModifiersTotal,
		&entity.DiscountAmount,
		&entity.LineTotal,
		&entity.SpecialInstructions,
		&entity.CustomerNotes,
		&entity.KitchenNotes,
		&entity.SeatNumber,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'pending',,
		&entity.'served',,
		&entity.UnitPrice,
		&entity.DiscountAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("order_items not found")
	}

	if err != nil {
		r.logger.Error("failed to get order_items", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get order_items: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of order_items records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*OrderItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_items
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, product_id
			, product_variant_id
			, item_name
			, quantity
			, unit_price
			, course_id
			, course_position
			, fire_time
			, kitchen_station_id
			, kitchen_ticket_id
			, status
			, fired_at
			, acknowledged_at
			, started_preparing_at
			, ready_at
			, served_at
			, modifiers_total
			, discount_amount
			, line_total
			, special_instructions
			, customer_notes
			, kitchen_notes
			, seat_number
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'pending',
			, 'served',
			, unit_price
			, discount_amount
		FROM order_items
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_items", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_items: %w", err)
	}
	defer rows.Close()

	var entities []*OrderItems
	for rows.Next() {
		var entity OrderItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ItemName,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.CourseId,
			&entity.CoursePosition,
			&entity.FireTime,
			&entity.KitchenStationId,
			&entity.KitchenTicketId,
			&entity.Status,
			&entity.FiredAt,
			&entity.AcknowledgedAt,
			&entity.StartedPreparingAt,
			&entity.ReadyAt,
			&entity.ServedAt,
			&entity.ModifiersTotal,
			&entity.DiscountAmount,
			&entity.LineTotal,
			&entity.SpecialInstructions,
			&entity.CustomerNotes,
			&entity.KitchenNotes,
			&entity.SeatNumber,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'pending',,
			&entity.'served',,
			&entity.UnitPrice,
			&entity.DiscountAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating order_items rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing order_items record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *OrderItems) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "order_items", duration, nil)
	}()

	query := `
		UPDATE order_items
		SET
			, organization_id = $2
			, order_id = $3
			, product_id = $4
			, product_variant_id = $5
			, item_name = $6
			, quantity = $7
			, unit_price = $8
			, course_id = $9
			, course_position = $10
			, fire_time = $11
			, kitchen_station_id = $12
			, kitchen_ticket_id = $13
			, status = $14
			, fired_at = $15
			, acknowledged_at = $16
			, started_preparing_at = $17
			, ready_at = $18
			, served_at = $19
			, modifiers_total = $20
			, discount_amount = $21
			, line_total = $22
			, special_instructions = $23
			, customer_notes = $24
			, kitchen_notes = $25
			, seat_number = $26
			, metadata = $27
			, updated_at = $29
			, created_by = $30
			, updated_by = $31
			, deleted_at = $32
			, 'pending', = $33
			, 'served', = $34
			, unit_price = $35
			, discount_amount = $36
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $37
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.ProductId,
		entity.ProductVariantId,
		entity.ItemName,
		entity.Quantity,
		entity.UnitPrice,
		entity.CourseId,
		entity.CoursePosition,
		entity.FireTime,
		entity.KitchenStationId,
		entity.KitchenTicketId,
		entity.Status,
		entity.FiredAt,
		entity.AcknowledgedAt,
		entity.StartedPreparingAt,
		entity.ReadyAt,
		entity.ServedAt,
		entity.ModifiersTotal,
		entity.DiscountAmount,
		entity.LineTotal,
		entity.SpecialInstructions,
		entity.CustomerNotes,
		entity.KitchenNotes,
		entity.SeatNumber,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'pending',,
		entity.'served',,
		entity.UnitPrice,
		entity.DiscountAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update order_items", zap.Error(err))
		return fmt.Errorf("failed to update order_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_items not found or already deleted")
	}

	r.logger.Info("updated order_items",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a order_items record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "order_items", duration, nil)
	}()

	query := `
		UPDATE order_items
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete order_items", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete order_items: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_items not found or already deleted")
	}

	r.logger.Info("deleted order_items", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves order_items records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*OrderItems, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_items", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_items records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, product_id
			, product_variant_id
			, item_name
			, quantity
			, unit_price
			, course_id
			, course_position
			, fire_time
			, kitchen_station_id
			, kitchen_ticket_id
			, status
			, fired_at
			, acknowledged_at
			, started_preparing_at
			, ready_at
			, served_at
			, modifiers_total
			, discount_amount
			, line_total
			, special_instructions
			, customer_notes
			, kitchen_notes
			, seat_number
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'pending',
			, 'served',
			, unit_price
			, discount_amount
		FROM order_items
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_items by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_items: %w", err)
	}
	defer rows.Close()

	var entities []*OrderItems
	for rows.Next() {
		var entity OrderItems
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.ProductId,
			&entity.ProductVariantId,
			&entity.ItemName,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.CourseId,
			&entity.CoursePosition,
			&entity.FireTime,
			&entity.KitchenStationId,
			&entity.KitchenTicketId,
			&entity.Status,
			&entity.FiredAt,
			&entity.AcknowledgedAt,
			&entity.StartedPreparingAt,
			&entity.ReadyAt,
			&entity.ServedAt,
			&entity.ModifiersTotal,
			&entity.DiscountAmount,
			&entity.LineTotal,
			&entity.SpecialInstructions,
			&entity.CustomerNotes,
			&entity.KitchenNotes,
			&entity.SeatNumber,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'pending',,
			&entity.'served',,
			&entity.UnitPrice,
			&entity.DiscountAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_items: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

