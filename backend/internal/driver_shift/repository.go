package driver_shift

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

// Repository handles database operations for DriverShifts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DriverShifts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DriverShifts represents a driver_shifts entity
type DriverShifts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	DriverId uuid.UUID `json:"driver_id" db:"driver_id"`
	ShiftDate time.Time `json:"shift_date" db:"shift_date"`
	ScheduledStartTime *string `json:"scheduled_start_time" db:"scheduled_start_time"`
	ScheduledEndTime *string `json:"scheduled_end_time" db:"scheduled_end_time"`
	ActualStartTime *time.Time `json:"actual_start_time" db:"actual_start_time"`
	ActualEndTime *time.Time `json:"actual_end_time" db:"actual_end_time"`
	Status *string `json:"status" db:"status"`
	TotalBreakMinutes *int64 `json:"total_break_minutes" db:"total_break_minutes"`
	TotalDeliveries *int64 `json:"total_deliveries" db:"total_deliveries"`
	TotalDistanceKm *float64 `json:"total_distance_km" db:"total_distance_km"`
	TotalEarnings *float64 `json:"total_earnings" db:"total_earnings"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'scheduled', *string `json:"'scheduled'," db:"'scheduled',"`
}

// Create inserts a new driver_shifts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DriverShifts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "driver_shifts", duration, nil)
	}()

	query := `
		INSERT INTO driver_shifts (
			, organization_id
			, location_id
			, driver_id
			, shift_date
			, scheduled_start_time
			, scheduled_end_time
			, actual_start_time
			, actual_end_time
			, status
			, total_break_minutes
			, total_deliveries
			, total_distance_km
			, total_earnings
			, notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
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
			, $19
			, $20
			, $21
			, $22
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.DriverId,
		entity.ShiftDate,
		entity.ScheduledStartTime,
		entity.ScheduledEndTime,
		entity.ActualStartTime,
		entity.ActualEndTime,
		entity.Status,
		entity.TotalBreakMinutes,
		entity.TotalDeliveries,
		entity.TotalDistanceKm,
		entity.TotalEarnings,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'scheduled',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create driver_shifts", zap.Error(err))
		return fmt.Errorf("failed to create driver_shifts: %w", err)
	}

	r.logger.Info("created driver_shifts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a driver_shifts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DriverShifts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "driver_shifts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, driver_id
			, shift_date
			, scheduled_start_time
			, scheduled_end_time
			, actual_start_time
			, actual_end_time
			, status
			, total_break_minutes
			, total_deliveries
			, total_distance_km
			, total_earnings
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
		FROM driver_shifts
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DriverShifts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.DriverId,
		&entity.ShiftDate,
		&entity.ScheduledStartTime,
		&entity.ScheduledEndTime,
		&entity.ActualStartTime,
		&entity.ActualEndTime,
		&entity.Status,
		&entity.TotalBreakMinutes,
		&entity.TotalDeliveries,
		&entity.TotalDistanceKm,
		&entity.TotalEarnings,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'scheduled',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("driver_shifts not found")
	}

	if err != nil {
		r.logger.Error("failed to get driver_shifts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get driver_shifts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of driver_shifts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DriverShifts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "driver_shifts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM driver_shifts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count driver_shifts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, driver_id
			, shift_date
			, scheduled_start_time
			, scheduled_end_time
			, actual_start_time
			, actual_end_time
			, status
			, total_break_minutes
			, total_deliveries
			, total_distance_km
			, total_earnings
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
		FROM driver_shifts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list driver_shifts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list driver_shifts: %w", err)
	}
	defer rows.Close()

	var entities []*DriverShifts
	for rows.Next() {
		var entity DriverShifts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.DriverId,
			&entity.ShiftDate,
			&entity.ScheduledStartTime,
			&entity.ScheduledEndTime,
			&entity.ActualStartTime,
			&entity.ActualEndTime,
			&entity.Status,
			&entity.TotalBreakMinutes,
			&entity.TotalDeliveries,
			&entity.TotalDistanceKm,
			&entity.TotalEarnings,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'scheduled',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan driver_shifts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating driver_shifts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing driver_shifts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DriverShifts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "driver_shifts", duration, nil)
	}()

	query := `
		UPDATE driver_shifts
		SET
			, organization_id = $2
			, location_id = $3
			, driver_id = $4
			, shift_date = $5
			, scheduled_start_time = $6
			, scheduled_end_time = $7
			, actual_start_time = $8
			, actual_end_time = $9
			, status = $10
			, total_break_minutes = $11
			, total_deliveries = $12
			, total_distance_km = $13
			, total_earnings = $14
			, notes = $15
			, metadata = $16
			, updated_at = $18
			, created_by = $19
			, updated_by = $20
			, deleted_at = $21
			, 'scheduled', = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.DriverId,
		entity.ShiftDate,
		entity.ScheduledStartTime,
		entity.ScheduledEndTime,
		entity.ActualStartTime,
		entity.ActualEndTime,
		entity.Status,
		entity.TotalBreakMinutes,
		entity.TotalDeliveries,
		entity.TotalDistanceKm,
		entity.TotalEarnings,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'scheduled',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update driver_shifts", zap.Error(err))
		return fmt.Errorf("failed to update driver_shifts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("driver_shifts not found or already deleted")
	}

	r.logger.Info("updated driver_shifts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a driver_shifts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "driver_shifts", duration, nil)
	}()

	query := `
		UPDATE driver_shifts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete driver_shifts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete driver_shifts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("driver_shifts not found or already deleted")
	}

	r.logger.Info("deleted driver_shifts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves driver_shifts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DriverShifts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "driver_shifts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM driver_shifts
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count driver_shifts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, driver_id
			, shift_date
			, scheduled_start_time
			, scheduled_end_time
			, actual_start_time
			, actual_end_time
			, status
			, total_break_minutes
			, total_deliveries
			, total_distance_km
			, total_earnings
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
		FROM driver_shifts
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list driver_shifts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list driver_shifts: %w", err)
	}
	defer rows.Close()

	var entities []*DriverShifts
	for rows.Next() {
		var entity DriverShifts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.DriverId,
			&entity.ShiftDate,
			&entity.ScheduledStartTime,
			&entity.ScheduledEndTime,
			&entity.ActualStartTime,
			&entity.ActualEndTime,
			&entity.Status,
			&entity.TotalBreakMinutes,
			&entity.TotalDeliveries,
			&entity.TotalDistanceKm,
			&entity.TotalEarnings,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'scheduled',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan driver_shifts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

