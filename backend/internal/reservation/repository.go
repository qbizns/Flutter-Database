package reservation

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

// Repository handles database operations for Reservations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Reservations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Reservations represents a reservations entity
type Reservations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	TableId *uuid.UUID `json:"table_id" db:"table_id"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	ReservationNumber string `json:"reservation_number" db:"reservation_number"`
	ReservationDate time.Time `json:"reservation_date" db:"reservation_date"`
	ReservationTime string `json:"reservation_time" db:"reservation_time"`
	DurationMinutes *int64 `json:"duration_minutes" db:"duration_minutes"`
	PartySize int64 `json:"party_size" db:"party_size"`
	CustomerName string `json:"customer_name" db:"customer_name"`
	CustomerPhone *string `json:"customer_phone" db:"customer_phone"`
	CustomerEmail *string `json:"customer_email" db:"customer_email"`
	Status *string `json:"status" db:"status"`
	AssignedWaiterId *uuid.UUID `json:"assigned_waiter_id" db:"assigned_waiter_id"`
	AssignedAt *time.Time `json:"assigned_at" db:"assigned_at"`
	SeatedAt *time.Time `json:"seated_at" db:"seated_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	SpecialRequests *string `json:"special_requests" db:"special_requests"`
	Occasion *string `json:"occasion" db:"occasion"`
	DietaryRestrictions *string `json:"dietary_restrictions" db:"dietary_restrictions"`
	ConfirmationCode *string `json:"confirmation_code" db:"confirmation_code"`
	ConfirmedAt *time.Time `json:"confirmed_at" db:"confirmed_at"`
	ConfirmedBy *uuid.UUID `json:"confirmed_by" db:"confirmed_by"`
	ReminderSentAt *time.Time `json:"reminder_sent_at" db:"reminder_sent_at"`
	NotificationPreferences json.RawMessage `json:"notification_preferences" db:"notification_preferences"`
	CancelledAt *time.Time `json:"cancelled_at" db:"cancelled_at"`
	CancelledBy *uuid.UUID `json:"cancelled_by" db:"cancelled_by"`
	CancellationReason *string `json:"cancellation_reason" db:"cancellation_reason"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new reservations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Reservations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "reservations", duration, nil)
	}()

	query := `
		INSERT INTO reservations (
			, organization_id
			, location_id
			, table_id
			, customer_id
			, reservation_number
			, reservation_date
			, reservation_time
			, duration_minutes
			, party_size
			, customer_name
			, customer_phone
			, customer_email
			, status
			, assigned_waiter_id
			, assigned_at
			, seated_at
			, completed_at
			, special_requests
			, occasion
			, dietary_restrictions
			, confirmation_code
			, confirmed_at
			, confirmed_by
			, reminder_sent_at
			, notification_preferences
			, cancelled_at
			, cancelled_by
			, cancellation_reason
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
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
			, $34
			, $35
			, $36
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TableId,
		entity.CustomerId,
		entity.ReservationNumber,
		entity.ReservationDate,
		entity.ReservationTime,
		entity.DurationMinutes,
		entity.PartySize,
		entity.CustomerName,
		entity.CustomerPhone,
		entity.CustomerEmail,
		entity.Status,
		entity.AssignedWaiterId,
		entity.AssignedAt,
		entity.SeatedAt,
		entity.CompletedAt,
		entity.SpecialRequests,
		entity.Occasion,
		entity.DietaryRestrictions,
		entity.ConfirmationCode,
		entity.ConfirmedAt,
		entity.ConfirmedBy,
		entity.ReminderSentAt,
		entity.NotificationPreferences,
		entity.CancelledAt,
		entity.CancelledBy,
		entity.CancellationReason,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create reservations", zap.Error(err))
		return fmt.Errorf("failed to create reservations: %w", err)
	}

	r.logger.Info("created reservations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a reservations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Reservations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reservations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, table_id
			, customer_id
			, reservation_number
			, reservation_date
			, reservation_time
			, duration_minutes
			, party_size
			, customer_name
			, customer_phone
			, customer_email
			, status
			, assigned_waiter_id
			, assigned_at
			, seated_at
			, completed_at
			, special_requests
			, occasion
			, dietary_restrictions
			, confirmation_code
			, confirmed_at
			, confirmed_by
			, reminder_sent_at
			, notification_preferences
			, cancelled_at
			, cancelled_by
			, cancellation_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM reservations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Reservations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.TableId,
		&entity.CustomerId,
		&entity.ReservationNumber,
		&entity.ReservationDate,
		&entity.ReservationTime,
		&entity.DurationMinutes,
		&entity.PartySize,
		&entity.CustomerName,
		&entity.CustomerPhone,
		&entity.CustomerEmail,
		&entity.Status,
		&entity.AssignedWaiterId,
		&entity.AssignedAt,
		&entity.SeatedAt,
		&entity.CompletedAt,
		&entity.SpecialRequests,
		&entity.Occasion,
		&entity.DietaryRestrictions,
		&entity.ConfirmationCode,
		&entity.ConfirmedAt,
		&entity.ConfirmedBy,
		&entity.ReminderSentAt,
		&entity.NotificationPreferences,
		&entity.CancelledAt,
		&entity.CancelledBy,
		&entity.CancellationReason,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("reservations not found")
	}

	if err != nil {
		r.logger.Error("failed to get reservations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of reservations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Reservations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reservations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM reservations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, table_id
			, customer_id
			, reservation_number
			, reservation_date
			, reservation_time
			, duration_minutes
			, party_size
			, customer_name
			, customer_phone
			, customer_email
			, status
			, assigned_waiter_id
			, assigned_at
			, seated_at
			, completed_at
			, special_requests
			, occasion
			, dietary_restrictions
			, confirmation_code
			, confirmed_at
			, confirmed_by
			, reminder_sent_at
			, notification_preferences
			, cancelled_at
			, cancelled_by
			, cancellation_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM reservations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list reservations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
	}
	defer rows.Close()

	var entities []*Reservations
	for rows.Next() {
		var entity Reservations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TableId,
			&entity.CustomerId,
			&entity.ReservationNumber,
			&entity.ReservationDate,
			&entity.ReservationTime,
			&entity.DurationMinutes,
			&entity.PartySize,
			&entity.CustomerName,
			&entity.CustomerPhone,
			&entity.CustomerEmail,
			&entity.Status,
			&entity.AssignedWaiterId,
			&entity.AssignedAt,
			&entity.SeatedAt,
			&entity.CompletedAt,
			&entity.SpecialRequests,
			&entity.Occasion,
			&entity.DietaryRestrictions,
			&entity.ConfirmationCode,
			&entity.ConfirmedAt,
			&entity.ConfirmedBy,
			&entity.ReminderSentAt,
			&entity.NotificationPreferences,
			&entity.CancelledAt,
			&entity.CancelledBy,
			&entity.CancellationReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating reservations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing reservations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Reservations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "reservations", duration, nil)
	}()

	query := `
		UPDATE reservations
		SET
			, organization_id = $2
			, location_id = $3
			, table_id = $4
			, customer_id = $5
			, reservation_number = $6
			, reservation_date = $7
			, reservation_time = $8
			, duration_minutes = $9
			, party_size = $10
			, customer_name = $11
			, customer_phone = $12
			, customer_email = $13
			, status = $14
			, assigned_waiter_id = $15
			, assigned_at = $16
			, seated_at = $17
			, completed_at = $18
			, special_requests = $19
			, occasion = $20
			, dietary_restrictions = $21
			, confirmation_code = $22
			, confirmed_at = $23
			, confirmed_by = $24
			, reminder_sent_at = $25
			, notification_preferences = $26
			, cancelled_at = $27
			, cancelled_by = $28
			, cancellation_reason = $29
			, notes = $30
			, metadata = $31
			, updated_at = $33
			, deleted_at = $34
			, created_by = $35
			, updated_by = $36
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $37
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TableId,
		entity.CustomerId,
		entity.ReservationNumber,
		entity.ReservationDate,
		entity.ReservationTime,
		entity.DurationMinutes,
		entity.PartySize,
		entity.CustomerName,
		entity.CustomerPhone,
		entity.CustomerEmail,
		entity.Status,
		entity.AssignedWaiterId,
		entity.AssignedAt,
		entity.SeatedAt,
		entity.CompletedAt,
		entity.SpecialRequests,
		entity.Occasion,
		entity.DietaryRestrictions,
		entity.ConfirmationCode,
		entity.ConfirmedAt,
		entity.ConfirmedBy,
		entity.ReminderSentAt,
		entity.NotificationPreferences,
		entity.CancelledAt,
		entity.CancelledBy,
		entity.CancellationReason,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update reservations", zap.Error(err))
		return fmt.Errorf("failed to update reservations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reservations not found or already deleted")
	}

	r.logger.Info("updated reservations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a reservations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "reservations", duration, nil)
	}()

	query := `
		UPDATE reservations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete reservations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete reservations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reservations not found or already deleted")
	}

	r.logger.Info("deleted reservations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves reservations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Reservations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reservations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM reservations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reservations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, table_id
			, customer_id
			, reservation_number
			, reservation_date
			, reservation_time
			, duration_minutes
			, party_size
			, customer_name
			, customer_phone
			, customer_email
			, status
			, assigned_waiter_id
			, assigned_at
			, seated_at
			, completed_at
			, special_requests
			, occasion
			, dietary_restrictions
			, confirmation_code
			, confirmed_at
			, confirmed_by
			, reminder_sent_at
			, notification_preferences
			, cancelled_at
			, cancelled_by
			, cancellation_reason
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM reservations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list reservations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
	}
	defer rows.Close()

	var entities []*Reservations
	for rows.Next() {
		var entity Reservations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TableId,
			&entity.CustomerId,
			&entity.ReservationNumber,
			&entity.ReservationDate,
			&entity.ReservationTime,
			&entity.DurationMinutes,
			&entity.PartySize,
			&entity.CustomerName,
			&entity.CustomerPhone,
			&entity.CustomerEmail,
			&entity.Status,
			&entity.AssignedWaiterId,
			&entity.AssignedAt,
			&entity.SeatedAt,
			&entity.CompletedAt,
			&entity.SpecialRequests,
			&entity.Occasion,
			&entity.DietaryRestrictions,
			&entity.ConfirmationCode,
			&entity.ConfirmedAt,
			&entity.ConfirmedBy,
			&entity.ReminderSentAt,
			&entity.NotificationPreferences,
			&entity.CancelledAt,
			&entity.CancelledBy,
			&entity.CancellationReason,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reservations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

