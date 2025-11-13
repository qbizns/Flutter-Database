package time_clock_entry

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

// Repository handles database operations for TimeClockEntries
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TimeClockEntries repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TimeClockEntries represents a time_clock_entries entity
type TimeClockEntries struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	EmployeeId uuid.UUID `json:"employee_id" db:"employee_id"`
	ScheduleId *uuid.UUID `json:"schedule_id" db:"schedule_id"`
	EntryType string `json:"entry_type" db:"entry_type"`
	EntryTimestamp *time.Time `json:"entry_timestamp" db:"entry_timestamp"`
	ScheduledTimestamp *time.Time `json:"scheduled_timestamp" db:"scheduled_timestamp"`
	DeviceId *uuid.UUID `json:"device_id" db:"device_id"`
	GpsLocation json.RawMessage `json:"gps_location" db:"gps_location"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	IsLate *bool `json:"is_late" db:"is_late"`
	IsEarly *bool `json:"is_early" db:"is_early"`
	VarianceMinutes *int64 `json:"variance_minutes" db:"variance_minutes"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	IsManualEntry *bool `json:"is_manual_entry" db:"is_manual_entry"`
	CorrectionNotes *string `json:"correction_notes" db:"correction_notes"`
	PhotoUrl *string `json:"photo_url" db:"photo_url"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'clockIn', *string `json:"'clock_in'," db:"'clock_in',"`
	'mealStart', *string `json:"'meal_start'," db:"'meal_start',"`
}

// Create inserts a new time_clock_entries record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TimeClockEntries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "time_clock_entries", duration, nil)
	}()

	query := `
		INSERT INTO time_clock_entries (
			, organization_id
			, location_id
			, employee_id
			, schedule_id
			, entry_type
			, entry_timestamp
			, scheduled_timestamp
			, device_id
			, gps_location
			, ip_address
			, is_late
			, is_early
			, variance_minutes
			, requires_approval
			, approved_by
			, approved_at
			, is_manual_entry
			, correction_notes
			, photo_url
			, notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'clock_in',
			, 'meal_start',
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
			, $28
			, $29
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.ScheduleId,
		entity.EntryType,
		entity.EntryTimestamp,
		entity.ScheduledTimestamp,
		entity.DeviceId,
		entity.GpsLocation,
		entity.IpAddress,
		entity.IsLate,
		entity.IsEarly,
		entity.VarianceMinutes,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.IsManualEntry,
		entity.CorrectionNotes,
		entity.PhotoUrl,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'clockIn',,
		entity.'mealStart',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create time_clock_entries", zap.Error(err))
		return fmt.Errorf("failed to create time_clock_entries: %w", err)
	}

	r.logger.Info("created time_clock_entries",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a time_clock_entries by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TimeClockEntries, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "time_clock_entries", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_id
			, entry_type
			, entry_timestamp
			, scheduled_timestamp
			, device_id
			, gps_location
			, ip_address
			, is_late
			, is_early
			, variance_minutes
			, requires_approval
			, approved_by
			, approved_at
			, is_manual_entry
			, correction_notes
			, photo_url
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'clock_in',
			, 'meal_start',
		FROM time_clock_entries
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TimeClockEntries
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.EmployeeId,
		&entity.ScheduleId,
		&entity.EntryType,
		&entity.EntryTimestamp,
		&entity.ScheduledTimestamp,
		&entity.DeviceId,
		&entity.GpsLocation,
		&entity.IpAddress,
		&entity.IsLate,
		&entity.IsEarly,
		&entity.VarianceMinutes,
		&entity.RequiresApproval,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.IsManualEntry,
		&entity.CorrectionNotes,
		&entity.PhotoUrl,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'clockIn',,
		&entity.'mealStart',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("time_clock_entries not found")
	}

	if err != nil {
		r.logger.Error("failed to get time_clock_entries", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get time_clock_entries: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of time_clock_entries records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TimeClockEntries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "time_clock_entries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM time_clock_entries
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count time_clock_entries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_id
			, entry_type
			, entry_timestamp
			, scheduled_timestamp
			, device_id
			, gps_location
			, ip_address
			, is_late
			, is_early
			, variance_minutes
			, requires_approval
			, approved_by
			, approved_at
			, is_manual_entry
			, correction_notes
			, photo_url
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'clock_in',
			, 'meal_start',
		FROM time_clock_entries
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list time_clock_entries", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list time_clock_entries: %w", err)
	}
	defer rows.Close()

	var entities []*TimeClockEntries
	for rows.Next() {
		var entity TimeClockEntries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.ScheduleId,
			&entity.EntryType,
			&entity.EntryTimestamp,
			&entity.ScheduledTimestamp,
			&entity.DeviceId,
			&entity.GpsLocation,
			&entity.IpAddress,
			&entity.IsLate,
			&entity.IsEarly,
			&entity.VarianceMinutes,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.IsManualEntry,
			&entity.CorrectionNotes,
			&entity.PhotoUrl,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'clockIn',,
			&entity.'mealStart',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan time_clock_entries: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating time_clock_entries rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing time_clock_entries record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TimeClockEntries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "time_clock_entries", duration, nil)
	}()

	query := `
		UPDATE time_clock_entries
		SET
			, organization_id = $2
			, location_id = $3
			, employee_id = $4
			, schedule_id = $5
			, entry_type = $6
			, entry_timestamp = $7
			, scheduled_timestamp = $8
			, device_id = $9
			, gps_location = $10
			, ip_address = $11
			, is_late = $12
			, is_early = $13
			, variance_minutes = $14
			, requires_approval = $15
			, approved_by = $16
			, approved_at = $17
			, is_manual_entry = $18
			, correction_notes = $19
			, photo_url = $20
			, notes = $21
			, metadata = $22
			, updated_at = $24
			, created_by = $25
			, updated_by = $26
			, deleted_at = $27
			, 'clock_in', = $28
			, 'meal_start', = $29
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $30
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.ScheduleId,
		entity.EntryType,
		entity.EntryTimestamp,
		entity.ScheduledTimestamp,
		entity.DeviceId,
		entity.GpsLocation,
		entity.IpAddress,
		entity.IsLate,
		entity.IsEarly,
		entity.VarianceMinutes,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.IsManualEntry,
		entity.CorrectionNotes,
		entity.PhotoUrl,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'clockIn',,
		entity.'mealStart',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update time_clock_entries", zap.Error(err))
		return fmt.Errorf("failed to update time_clock_entries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("time_clock_entries not found or already deleted")
	}

	r.logger.Info("updated time_clock_entries",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a time_clock_entries record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "time_clock_entries", duration, nil)
	}()

	query := `
		UPDATE time_clock_entries
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete time_clock_entries", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete time_clock_entries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("time_clock_entries not found or already deleted")
	}

	r.logger.Info("deleted time_clock_entries", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves time_clock_entries records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TimeClockEntries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "time_clock_entries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM time_clock_entries
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count time_clock_entries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_id
			, entry_type
			, entry_timestamp
			, scheduled_timestamp
			, device_id
			, gps_location
			, ip_address
			, is_late
			, is_early
			, variance_minutes
			, requires_approval
			, approved_by
			, approved_at
			, is_manual_entry
			, correction_notes
			, photo_url
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'clock_in',
			, 'meal_start',
		FROM time_clock_entries
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list time_clock_entries by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list time_clock_entries: %w", err)
	}
	defer rows.Close()

	var entities []*TimeClockEntries
	for rows.Next() {
		var entity TimeClockEntries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.ScheduleId,
			&entity.EntryType,
			&entity.EntryTimestamp,
			&entity.ScheduledTimestamp,
			&entity.DeviceId,
			&entity.GpsLocation,
			&entity.IpAddress,
			&entity.IsLate,
			&entity.IsEarly,
			&entity.VarianceMinutes,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.IsManualEntry,
			&entity.CorrectionNotes,
			&entity.PhotoUrl,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'clockIn',,
			&entity.'mealStart',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan time_clock_entries: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

