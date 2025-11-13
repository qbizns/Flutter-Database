package employee_schedule

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

// Repository handles database operations for EmployeeSchedules
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new EmployeeSchedules repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// EmployeeSchedules represents a employee_schedules entity
type EmployeeSchedules struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	EmployeeId uuid.UUID `json:"employee_id" db:"employee_id"`
	ScheduleDate time.Time `json:"schedule_date" db:"schedule_date"`
	ShiftType *string `json:"shift_type" db:"shift_type"`
	Position *string `json:"position" db:"position"`
	ScheduledStartTime string `json:"scheduled_start_time" db:"scheduled_start_time"`
	ScheduledEndTime string `json:"scheduled_end_time" db:"scheduled_end_time"`
	BreakDurationMinutes *int64 `json:"break_duration_minutes" db:"break_duration_minutes"`
	Status *string `json:"status" db:"status"`
	RequiresApproval *bool `json:"requires_approval" db:"requires_approval"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Notes *string `json:"notes" db:"notes"`
	CancellationReason *string `json:"cancellation_reason" db:"cancellation_reason"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'scheduled', *string `json:"'scheduled'," db:"'scheduled',"`
	'regular', *string `json:"'regular'," db:"'regular',"`
}

// Create inserts a new employee_schedules record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *EmployeeSchedules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "employee_schedules", duration, nil)
	}()

	query := `
		INSERT INTO employee_schedules (
			, organization_id
			, location_id
			, employee_id
			, schedule_date
			, shift_type
			, position
			, scheduled_start_time
			, scheduled_end_time
			, break_duration_minutes
			, status
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, cancellation_reason
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
			, 'regular',
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
			, $20
			, $21
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.ScheduleDate,
		entity.ShiftType,
		entity.Position,
		entity.ScheduledStartTime,
		entity.ScheduledEndTime,
		entity.BreakDurationMinutes,
		entity.Status,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CancellationReason,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'scheduled',,
		entity.'regular',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create employee_schedules", zap.Error(err))
		return fmt.Errorf("failed to create employee_schedules: %w", err)
	}

	r.logger.Info("created employee_schedules",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a employee_schedules by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*EmployeeSchedules, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "employee_schedules", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_date
			, shift_type
			, position
			, scheduled_start_time
			, scheduled_end_time
			, break_duration_minutes
			, status
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, cancellation_reason
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
			, 'regular',
		FROM employee_schedules
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity EmployeeSchedules
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.EmployeeId,
		&entity.ScheduleDate,
		&entity.ShiftType,
		&entity.Position,
		&entity.ScheduledStartTime,
		&entity.ScheduledEndTime,
		&entity.BreakDurationMinutes,
		&entity.Status,
		&entity.RequiresApproval,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.Notes,
		&entity.CancellationReason,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'scheduled',,
		&entity.'regular',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("employee_schedules not found")
	}

	if err != nil {
		r.logger.Error("failed to get employee_schedules", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get employee_schedules: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of employee_schedules records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*EmployeeSchedules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "employee_schedules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM employee_schedules
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count employee_schedules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_date
			, shift_type
			, position
			, scheduled_start_time
			, scheduled_end_time
			, break_duration_minutes
			, status
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, cancellation_reason
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
			, 'regular',
		FROM employee_schedules
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list employee_schedules", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list employee_schedules: %w", err)
	}
	defer rows.Close()

	var entities []*EmployeeSchedules
	for rows.Next() {
		var entity EmployeeSchedules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.ScheduleDate,
			&entity.ShiftType,
			&entity.Position,
			&entity.ScheduledStartTime,
			&entity.ScheduledEndTime,
			&entity.BreakDurationMinutes,
			&entity.Status,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CancellationReason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'scheduled',,
			&entity.'regular',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan employee_schedules: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating employee_schedules rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing employee_schedules record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *EmployeeSchedules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "employee_schedules", duration, nil)
	}()

	query := `
		UPDATE employee_schedules
		SET
			, organization_id = $2
			, location_id = $3
			, employee_id = $4
			, schedule_date = $5
			, shift_type = $6
			, position = $7
			, scheduled_start_time = $8
			, scheduled_end_time = $9
			, break_duration_minutes = $10
			, status = $11
			, requires_approval = $12
			, approved_by = $13
			, approved_at = $14
			, notes = $15
			, cancellation_reason = $16
			, metadata = $17
			, updated_at = $19
			, created_by = $20
			, updated_by = $21
			, deleted_at = $22
			, 'scheduled', = $23
			, 'regular', = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.ScheduleDate,
		entity.ShiftType,
		entity.Position,
		entity.ScheduledStartTime,
		entity.ScheduledEndTime,
		entity.BreakDurationMinutes,
		entity.Status,
		entity.RequiresApproval,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CancellationReason,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'scheduled',,
		entity.'regular',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update employee_schedules", zap.Error(err))
		return fmt.Errorf("failed to update employee_schedules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("employee_schedules not found or already deleted")
	}

	r.logger.Info("updated employee_schedules",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a employee_schedules record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "employee_schedules", duration, nil)
	}()

	query := `
		UPDATE employee_schedules
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete employee_schedules", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete employee_schedules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("employee_schedules not found or already deleted")
	}

	r.logger.Info("deleted employee_schedules", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves employee_schedules records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*EmployeeSchedules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "employee_schedules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM employee_schedules
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count employee_schedules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, schedule_date
			, shift_type
			, position
			, scheduled_start_time
			, scheduled_end_time
			, break_duration_minutes
			, status
			, requires_approval
			, approved_by
			, approved_at
			, notes
			, cancellation_reason
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'scheduled',
			, 'regular',
		FROM employee_schedules
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list employee_schedules by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list employee_schedules: %w", err)
	}
	defer rows.Close()

	var entities []*EmployeeSchedules
	for rows.Next() {
		var entity EmployeeSchedules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.ScheduleDate,
			&entity.ShiftType,
			&entity.Position,
			&entity.ScheduledStartTime,
			&entity.ScheduledEndTime,
			&entity.BreakDurationMinutes,
			&entity.Status,
			&entity.RequiresApproval,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CancellationReason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'scheduled',,
			&entity.'regular',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan employee_schedules: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

