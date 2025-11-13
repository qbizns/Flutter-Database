package order_tracking_event

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

// Repository handles database operations for OrderTrackingEvents
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new OrderTrackingEvents repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// OrderTrackingEvents represents a order_tracking_events entity
type OrderTrackingEvents struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	OrderId uuid.UUID `json:"order_id" db:"order_id"`
	DeliveryAssignmentId *uuid.UUID `json:"delivery_assignment_id" db:"delivery_assignment_id"`
	EventType string `json:"event_type" db:"event_type"`
	EventTimestamp *time.Time `json:"event_timestamp" db:"event_timestamp"`
	EventMessage *string `json:"event_message" db:"event_message"`
	Location json.RawMessage `json:"location" db:"location"`
	LocationName *string `json:"location_name" db:"location_name"`
	ActorType *string `json:"actor_type" db:"actor_type"`
	ActorId *uuid.UUID `json:"actor_id" db:"actor_id"`
	ActorName *string `json:"actor_name" db:"actor_name"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	'orderPlaced', *string `json:"'order_placed'," db:"'order_placed',"`
	'readyForPickup', *string `json:"'ready_for_pickup'," db:"'ready_for_pickup',"`
	'arrived', *string `json:"'arrived'," db:"'arrived',"`
	'rescheduled', *string `json:"'rescheduled'," db:"'rescheduled',"`
	'system', *string `json:"'system'," db:"'system',"`
}

// Create inserts a new order_tracking_events record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *OrderTrackingEvents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "order_tracking_events", duration, nil)
	}()

	query := `
		INSERT INTO order_tracking_events (
			, organization_id
			, order_id
			, delivery_assignment_id
			, event_type
			, event_timestamp
			, event_message
			, location
			, location_name
			, actor_type
			, actor_id
			, actor_name
			, metadata
			, created_by
			, 'order_placed',
			, 'ready_for_pickup',
			, 'arrived',
			, 'rescheduled',
			, 'system',
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
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.DeliveryAssignmentId,
		entity.EventType,
		entity.EventTimestamp,
		entity.EventMessage,
		entity.Location,
		entity.LocationName,
		entity.ActorType,
		entity.ActorId,
		entity.ActorName,
		entity.Metadata,
		entity.CreatedBy,
		entity.'orderPlaced',,
		entity.'readyForPickup',,
		entity.'arrived',,
		entity.'rescheduled',,
		entity.'system',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create order_tracking_events", zap.Error(err))
		return fmt.Errorf("failed to create order_tracking_events: %w", err)
	}

	r.logger.Info("created order_tracking_events",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a order_tracking_events by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*OrderTrackingEvents, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_tracking_events", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, order_id
			, delivery_assignment_id
			, event_type
			, event_timestamp
			, event_message
			, location
			, location_name
			, actor_type
			, actor_id
			, actor_name
			, metadata
			, created_at
			, created_by
			, 'order_placed',
			, 'ready_for_pickup',
			, 'arrived',
			, 'rescheduled',
			, 'system',
		FROM order_tracking_events
		WHERE id = $1
		
	`

	var entity OrderTrackingEvents
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.OrderId,
		&entity.DeliveryAssignmentId,
		&entity.EventType,
		&entity.EventTimestamp,
		&entity.EventMessage,
		&entity.Location,
		&entity.LocationName,
		&entity.ActorType,
		&entity.ActorId,
		&entity.ActorName,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.CreatedBy,
		&entity.'orderPlaced',,
		&entity.'readyForPickup',,
		&entity.'arrived',,
		&entity.'rescheduled',,
		&entity.'system',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("order_tracking_events not found")
	}

	if err != nil {
		r.logger.Error("failed to get order_tracking_events", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get order_tracking_events: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of order_tracking_events records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*OrderTrackingEvents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_tracking_events", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_tracking_events
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_tracking_events records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, delivery_assignment_id
			, event_type
			, event_timestamp
			, event_message
			, location
			, location_name
			, actor_type
			, actor_id
			, actor_name
			, metadata
			, created_at
			, created_by
			, 'order_placed',
			, 'ready_for_pickup',
			, 'arrived',
			, 'rescheduled',
			, 'system',
		FROM order_tracking_events
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_tracking_events", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_tracking_events: %w", err)
	}
	defer rows.Close()

	var entities []*OrderTrackingEvents
	for rows.Next() {
		var entity OrderTrackingEvents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.DeliveryAssignmentId,
			&entity.EventType,
			&entity.EventTimestamp,
			&entity.EventMessage,
			&entity.Location,
			&entity.LocationName,
			&entity.ActorType,
			&entity.ActorId,
			&entity.ActorName,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.'orderPlaced',,
			&entity.'readyForPickup',,
			&entity.'arrived',,
			&entity.'rescheduled',,
			&entity.'system',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_tracking_events: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating order_tracking_events rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing order_tracking_events record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *OrderTrackingEvents) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "order_tracking_events", duration, nil)
	}()

	query := `
		UPDATE order_tracking_events
		SET
			, organization_id = $2
			, order_id = $3
			, delivery_assignment_id = $4
			, event_type = $5
			, event_timestamp = $6
			, event_message = $7
			, location = $8
			, location_name = $9
			, actor_type = $10
			, actor_id = $11
			, actor_name = $12
			, metadata = $13
			, created_by = $15
			, 'order_placed', = $16
			, 'ready_for_pickup', = $17
			, 'arrived', = $18
			, 'rescheduled', = $19
			, 'system', = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.OrderId,
		entity.DeliveryAssignmentId,
		entity.EventType,
		entity.EventTimestamp,
		entity.EventMessage,
		entity.Location,
		entity.LocationName,
		entity.ActorType,
		entity.ActorId,
		entity.ActorName,
		entity.Metadata,
		entity.CreatedBy,
		entity.'orderPlaced',,
		entity.'readyForPickup',,
		entity.'arrived',,
		entity.'rescheduled',,
		entity.'system',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update order_tracking_events", zap.Error(err))
		return fmt.Errorf("failed to update order_tracking_events: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_tracking_events not found or already deleted")
	}

	r.logger.Info("updated order_tracking_events",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a order_tracking_events record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "order_tracking_events", duration, nil)
	}()

	query := `DELETE FROM order_tracking_events WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete order_tracking_events", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete order_tracking_events: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order_tracking_events not found")
	}

	r.logger.Info("deleted order_tracking_events", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves order_tracking_events records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*OrderTrackingEvents, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "order_tracking_events", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM order_tracking_events
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count order_tracking_events records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, order_id
			, delivery_assignment_id
			, event_type
			, event_timestamp
			, event_message
			, location
			, location_name
			, actor_type
			, actor_id
			, actor_name
			, metadata
			, created_at
			, created_by
			, 'order_placed',
			, 'ready_for_pickup',
			, 'arrived',
			, 'rescheduled',
			, 'system',
		FROM order_tracking_events
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list order_tracking_events by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list order_tracking_events: %w", err)
	}
	defer rows.Close()

	var entities []*OrderTrackingEvents
	for rows.Next() {
		var entity OrderTrackingEvents
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.OrderId,
			&entity.DeliveryAssignmentId,
			&entity.EventType,
			&entity.EventTimestamp,
			&entity.EventMessage,
			&entity.Location,
			&entity.LocationName,
			&entity.ActorType,
			&entity.ActorId,
			&entity.ActorName,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
			&entity.'orderPlaced',,
			&entity.'readyForPickup',,
			&entity.'arrived',,
			&entity.'rescheduled',,
			&entity.'system',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan order_tracking_events: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

