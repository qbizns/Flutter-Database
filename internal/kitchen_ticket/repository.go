package kitchen_ticket

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

// Repository handles database operations for KitchenTickets
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new KitchenTickets repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// KitchenTickets represents a kitchen_tickets entity
type KitchenTickets struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	TicketNumber string `json:"ticket_number" db:"ticket_number"`
	DisplaySequence *int64 `json:"display_sequence" db:"display_sequence"`
	OrderId uuid.UUID `json:"order_id" db:"order_id"`
	KitchenStationId uuid.UUID `json:"kitchen_station_id" db:"kitchen_station_id"`
	CourseId *uuid.UUID `json:"course_id" db:"course_id"`
	TicketType *string `json:"ticket_type" db:"ticket_type"`
	Priority *int64 `json:"priority" db:"priority"`
	Status *string `json:"status" db:"status"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	FiredAt *time.Time `json:"fired_at" db:"fired_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at" db:"acknowledged_at"`
	StartedAt *time.Time `json:"started_at" db:"started_at"`
	ReadyAt *time.Time `json:"ready_at" db:"ready_at"`
	BumpedAt *time.Time `json:"bumped_at" db:"bumped_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	PrepTimeMinutes *int64 `json:"prep_time_minutes" db:"prep_time_minutes"`
	TargetPrepTime *int64 `json:"target_prep_time" db:"target_prep_time"`
	TableNumber *string `json:"table_number" db:"table_number"`
	OrderType *string `json:"order_type" db:"order_type"`
	Covers *int64 `json:"covers" db:"covers"`
	WaiterName *string `json:"waiter_name" db:"waiter_name"`
	SpecialInstructions *string `json:"special_instructions" db:"special_instructions"`
	KitchenNotes *string `json:"kitchen_notes" db:"kitchen_notes"`
	DisplayConfig json.RawMessage `json:"display_config" db:"display_config"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'new', *string `json:"'new'," db:"'new',"`
	'completed', *string `json:"'completed'," db:"'completed',"`
	'normal', *string `json:"'normal'," db:"'normal',"`
}

// Create inserts a new kitchen_tickets record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *KitchenTickets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "kitchen_tickets", duration, nil)
	}()

	query := `
		INSERT INTO kitchen_tickets (
			, organization_id
			, location_id
			, ticket_number
			, display_sequence
			, order_id
			, kitchen_station_id
			, course_id
			, ticket_type
			, priority
			, status
			, fired_at
			, acknowledged_at
			, started_at
			, ready_at
			, bumped_at
			, completed_at
			, prep_time_minutes
			, target_prep_time
			, table_number
			, order_type
			, covers
			, waiter_name
			, special_instructions
			, kitchen_notes
			, display_config
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'new',
			, 'completed',
			, 'normal',
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
			, $30
			, $31
			, $32
			, $33
			, $34
			, $35
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TicketNumber,
		entity.DisplaySequence,
		entity.OrderId,
		entity.KitchenStationId,
		entity.CourseId,
		entity.TicketType,
		entity.Priority,
		entity.Status,
		entity.FiredAt,
		entity.AcknowledgedAt,
		entity.StartedAt,
		entity.ReadyAt,
		entity.BumpedAt,
		entity.CompletedAt,
		entity.PrepTimeMinutes,
		entity.TargetPrepTime,
		entity.TableNumber,
		entity.OrderType,
		entity.Covers,
		entity.WaiterName,
		entity.SpecialInstructions,
		entity.KitchenNotes,
		entity.DisplayConfig,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'new',,
		entity.'completed',,
		entity.'normal',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create kitchen_tickets", zap.Error(err))
		return fmt.Errorf("failed to create kitchen_tickets: %w", err)
	}

	r.logger.Info("created kitchen_tickets",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a kitchen_tickets by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*KitchenTickets, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_tickets", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, ticket_number
			, display_sequence
			, order_id
			, kitchen_station_id
			, course_id
			, ticket_type
			, priority
			, status
			, created_at
			, fired_at
			, acknowledged_at
			, started_at
			, ready_at
			, bumped_at
			, completed_at
			, prep_time_minutes
			, target_prep_time
			, table_number
			, order_type
			, covers
			, waiter_name
			, special_instructions
			, kitchen_notes
			, display_config
			, metadata
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'new',
			, 'completed',
			, 'normal',
		FROM kitchen_tickets
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity KitchenTickets
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.TicketNumber,
		&entity.DisplaySequence,
		&entity.OrderId,
		&entity.KitchenStationId,
		&entity.CourseId,
		&entity.TicketType,
		&entity.Priority,
		&entity.Status,
		&entity.CreatedAt,
		&entity.FiredAt,
		&entity.AcknowledgedAt,
		&entity.StartedAt,
		&entity.ReadyAt,
		&entity.BumpedAt,
		&entity.CompletedAt,
		&entity.PrepTimeMinutes,
		&entity.TargetPrepTime,
		&entity.TableNumber,
		&entity.OrderType,
		&entity.Covers,
		&entity.WaiterName,
		&entity.SpecialInstructions,
		&entity.KitchenNotes,
		&entity.DisplayConfig,
		&entity.Metadata,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'new',,
		&entity.'completed',,
		&entity.'normal',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("kitchen_tickets not found")
	}

	if err != nil {
		r.logger.Error("failed to get kitchen_tickets", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get kitchen_tickets: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of kitchen_tickets records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*KitchenTickets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_tickets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM kitchen_tickets
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count kitchen_tickets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, ticket_number
			, display_sequence
			, order_id
			, kitchen_station_id
			, course_id
			, ticket_type
			, priority
			, status
			, created_at
			, fired_at
			, acknowledged_at
			, started_at
			, ready_at
			, bumped_at
			, completed_at
			, prep_time_minutes
			, target_prep_time
			, table_number
			, order_type
			, covers
			, waiter_name
			, special_instructions
			, kitchen_notes
			, display_config
			, metadata
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'new',
			, 'completed',
			, 'normal',
		FROM kitchen_tickets
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list kitchen_tickets", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list kitchen_tickets: %w", err)
	}
	defer rows.Close()

	var entities []*KitchenTickets
	for rows.Next() {
		var entity KitchenTickets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TicketNumber,
			&entity.DisplaySequence,
			&entity.OrderId,
			&entity.KitchenStationId,
			&entity.CourseId,
			&entity.TicketType,
			&entity.Priority,
			&entity.Status,
			&entity.CreatedAt,
			&entity.FiredAt,
			&entity.AcknowledgedAt,
			&entity.StartedAt,
			&entity.ReadyAt,
			&entity.BumpedAt,
			&entity.CompletedAt,
			&entity.PrepTimeMinutes,
			&entity.TargetPrepTime,
			&entity.TableNumber,
			&entity.OrderType,
			&entity.Covers,
			&entity.WaiterName,
			&entity.SpecialInstructions,
			&entity.KitchenNotes,
			&entity.DisplayConfig,
			&entity.Metadata,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'new',,
			&entity.'completed',,
			&entity.'normal',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan kitchen_tickets: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating kitchen_tickets rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing kitchen_tickets record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *KitchenTickets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "kitchen_tickets", duration, nil)
	}()

	query := `
		UPDATE kitchen_tickets
		SET
			, organization_id = $2
			, location_id = $3
			, ticket_number = $4
			, display_sequence = $5
			, order_id = $6
			, kitchen_station_id = $7
			, course_id = $8
			, ticket_type = $9
			, priority = $10
			, status = $11
			, fired_at = $13
			, acknowledged_at = $14
			, started_at = $15
			, ready_at = $16
			, bumped_at = $17
			, completed_at = $18
			, prep_time_minutes = $19
			, target_prep_time = $20
			, table_number = $21
			, order_type = $22
			, covers = $23
			, waiter_name = $24
			, special_instructions = $25
			, kitchen_notes = $26
			, display_config = $27
			, metadata = $28
			, updated_at = $29
			, created_by = $30
			, updated_by = $31
			, deleted_at = $32
			, 'new', = $33
			, 'completed', = $34
			, 'normal', = $35
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $36
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TicketNumber,
		entity.DisplaySequence,
		entity.OrderId,
		entity.KitchenStationId,
		entity.CourseId,
		entity.TicketType,
		entity.Priority,
		entity.Status,
		entity.FiredAt,
		entity.AcknowledgedAt,
		entity.StartedAt,
		entity.ReadyAt,
		entity.BumpedAt,
		entity.CompletedAt,
		entity.PrepTimeMinutes,
		entity.TargetPrepTime,
		entity.TableNumber,
		entity.OrderType,
		entity.Covers,
		entity.WaiterName,
		entity.SpecialInstructions,
		entity.KitchenNotes,
		entity.DisplayConfig,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'new',,
		entity.'completed',,
		entity.'normal',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update kitchen_tickets", zap.Error(err))
		return fmt.Errorf("failed to update kitchen_tickets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("kitchen_tickets not found or already deleted")
	}

	r.logger.Info("updated kitchen_tickets",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a kitchen_tickets record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "kitchen_tickets", duration, nil)
	}()

	query := `
		UPDATE kitchen_tickets
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete kitchen_tickets", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete kitchen_tickets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("kitchen_tickets not found or already deleted")
	}

	r.logger.Info("deleted kitchen_tickets", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves kitchen_tickets records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*KitchenTickets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_tickets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM kitchen_tickets
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count kitchen_tickets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, ticket_number
			, display_sequence
			, order_id
			, kitchen_station_id
			, course_id
			, ticket_type
			, priority
			, status
			, created_at
			, fired_at
			, acknowledged_at
			, started_at
			, ready_at
			, bumped_at
			, completed_at
			, prep_time_minutes
			, target_prep_time
			, table_number
			, order_type
			, covers
			, waiter_name
			, special_instructions
			, kitchen_notes
			, display_config
			, metadata
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'new',
			, 'completed',
			, 'normal',
		FROM kitchen_tickets
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list kitchen_tickets by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list kitchen_tickets: %w", err)
	}
	defer rows.Close()

	var entities []*KitchenTickets
	for rows.Next() {
		var entity KitchenTickets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TicketNumber,
			&entity.DisplaySequence,
			&entity.OrderId,
			&entity.KitchenStationId,
			&entity.CourseId,
			&entity.TicketType,
			&entity.Priority,
			&entity.Status,
			&entity.CreatedAt,
			&entity.FiredAt,
			&entity.AcknowledgedAt,
			&entity.StartedAt,
			&entity.ReadyAt,
			&entity.BumpedAt,
			&entity.CompletedAt,
			&entity.PrepTimeMinutes,
			&entity.TargetPrepTime,
			&entity.TableNumber,
			&entity.OrderType,
			&entity.Covers,
			&entity.WaiterName,
			&entity.SpecialInstructions,
			&entity.KitchenNotes,
			&entity.DisplayConfig,
			&entity.Metadata,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'new',,
			&entity.'completed',,
			&entity.'normal',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan kitchen_tickets: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

