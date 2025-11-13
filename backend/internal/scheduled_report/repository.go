package scheduled_report

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

// Repository handles database operations for ScheduledReports
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ScheduledReports repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ScheduledReports represents a scheduled_reports entity
type ScheduledReports struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ReportName string `json:"report_name" db:"report_name"`
	ReportType string `json:"report_type" db:"report_type"`
	ScheduleFrequency string `json:"schedule_frequency" db:"schedule_frequency"`
	ScheduleDayOfWeek *int64 `json:"schedule_day_of_week" db:"schedule_day_of_week"`
	ScheduleDayOfMonth *int64 `json:"schedule_day_of_month" db:"schedule_day_of_month"`
	ScheduleTime string `json:"schedule_time" db:"schedule_time"`
	ScheduleTimezone *string `json:"schedule_timezone" db:"schedule_timezone"`
	ReportParameters json.RawMessage `json:"report_parameters" db:"report_parameters"`
	DeliveryMethod *string `json:"delivery_method" db:"delivery_method"`
	DeliveryRecipients *string `json:"delivery_recipients" db:"delivery_recipients"`
	OutputFormat *string `json:"output_format" db:"output_format"`
	IsActive *bool `json:"is_active" db:"is_active"`
	LastRunAt *time.Time `json:"last_run_at" db:"last_run_at"`
	LastRunStatus *string `json:"last_run_status" db:"last_run_status"`
	NextRunAt *time.Time `json:"next_run_at" db:"next_run_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new scheduled_reports record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ScheduledReports) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "scheduled_reports", duration, nil)
	}()

	query := `
		INSERT INTO scheduled_reports (
			, organization_id
			, report_name
			, report_type
			, schedule_frequency
			, schedule_day_of_week
			, schedule_day_of_month
			, schedule_time
			, schedule_timezone
			, report_parameters
			, delivery_method
			, delivery_recipients
			, output_format
			, is_active
			, last_run_at
			, last_run_status
			, next_run_at
			, created_by
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ReportName,
		entity.ReportType,
		entity.ScheduleFrequency,
		entity.ScheduleDayOfWeek,
		entity.ScheduleDayOfMonth,
		entity.ScheduleTime,
		entity.ScheduleTimezone,
		entity.ReportParameters,
		entity.DeliveryMethod,
		entity.DeliveryRecipients,
		entity.OutputFormat,
		entity.IsActive,
		entity.LastRunAt,
		entity.LastRunStatus,
		entity.NextRunAt,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create scheduled_reports", zap.Error(err))
		return fmt.Errorf("failed to create scheduled_reports: %w", err)
	}

	r.logger.Info("created scheduled_reports",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a scheduled_reports by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ScheduledReports, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "scheduled_reports", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, report_name
			, report_type
			, schedule_frequency
			, schedule_day_of_week
			, schedule_day_of_month
			, schedule_time
			, schedule_timezone
			, report_parameters
			, delivery_method
			, delivery_recipients
			, output_format
			, is_active
			, last_run_at
			, last_run_status
			, next_run_at
			, created_by
			, created_at
			, updated_at
		FROM scheduled_reports
		WHERE id = $1
		
	`

	var entity ScheduledReports
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ReportName,
		&entity.ReportType,
		&entity.ScheduleFrequency,
		&entity.ScheduleDayOfWeek,
		&entity.ScheduleDayOfMonth,
		&entity.ScheduleTime,
		&entity.ScheduleTimezone,
		&entity.ReportParameters,
		&entity.DeliveryMethod,
		&entity.DeliveryRecipients,
		&entity.OutputFormat,
		&entity.IsActive,
		&entity.LastRunAt,
		&entity.LastRunStatus,
		&entity.NextRunAt,
		&entity.CreatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("scheduled_reports not found")
	}

	if err != nil {
		r.logger.Error("failed to get scheduled_reports", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get scheduled_reports: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of scheduled_reports records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ScheduledReports, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "scheduled_reports", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM scheduled_reports
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduled_reports records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, report_name
			, report_type
			, schedule_frequency
			, schedule_day_of_week
			, schedule_day_of_month
			, schedule_time
			, schedule_timezone
			, report_parameters
			, delivery_method
			, delivery_recipients
			, output_format
			, is_active
			, last_run_at
			, last_run_status
			, next_run_at
			, created_by
			, created_at
			, updated_at
		FROM scheduled_reports
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list scheduled_reports", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list scheduled_reports: %w", err)
	}
	defer rows.Close()

	var entities []*ScheduledReports
	for rows.Next() {
		var entity ScheduledReports
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReportName,
			&entity.ReportType,
			&entity.ScheduleFrequency,
			&entity.ScheduleDayOfWeek,
			&entity.ScheduleDayOfMonth,
			&entity.ScheduleTime,
			&entity.ScheduleTimezone,
			&entity.ReportParameters,
			&entity.DeliveryMethod,
			&entity.DeliveryRecipients,
			&entity.OutputFormat,
			&entity.IsActive,
			&entity.LastRunAt,
			&entity.LastRunStatus,
			&entity.NextRunAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan scheduled_reports: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating scheduled_reports rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing scheduled_reports record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ScheduledReports) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "scheduled_reports", duration, nil)
	}()

	query := `
		UPDATE scheduled_reports
		SET
			, organization_id = $2
			, report_name = $3
			, report_type = $4
			, schedule_frequency = $5
			, schedule_day_of_week = $6
			, schedule_day_of_month = $7
			, schedule_time = $8
			, schedule_timezone = $9
			, report_parameters = $10
			, delivery_method = $11
			, delivery_recipients = $12
			, output_format = $13
			, is_active = $14
			, last_run_at = $15
			, last_run_status = $16
			, next_run_at = $17
			, created_by = $18
			, updated_at = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ReportName,
		entity.ReportType,
		entity.ScheduleFrequency,
		entity.ScheduleDayOfWeek,
		entity.ScheduleDayOfMonth,
		entity.ScheduleTime,
		entity.ScheduleTimezone,
		entity.ReportParameters,
		entity.DeliveryMethod,
		entity.DeliveryRecipients,
		entity.OutputFormat,
		entity.IsActive,
		entity.LastRunAt,
		entity.LastRunStatus,
		entity.NextRunAt,
		entity.CreatedBy,
		time.Now(),
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update scheduled_reports", zap.Error(err))
		return fmt.Errorf("failed to update scheduled_reports: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("scheduled_reports not found or already deleted")
	}

	r.logger.Info("updated scheduled_reports",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a scheduled_reports record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "scheduled_reports", duration, nil)
	}()

	query := `DELETE FROM scheduled_reports WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete scheduled_reports", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete scheduled_reports: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("scheduled_reports not found")
	}

	r.logger.Info("deleted scheduled_reports", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves scheduled_reports records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ScheduledReports, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "scheduled_reports", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM scheduled_reports
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count scheduled_reports records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, report_name
			, report_type
			, schedule_frequency
			, schedule_day_of_week
			, schedule_day_of_month
			, schedule_time
			, schedule_timezone
			, report_parameters
			, delivery_method
			, delivery_recipients
			, output_format
			, is_active
			, last_run_at
			, last_run_status
			, next_run_at
			, created_by
			, created_at
			, updated_at
		FROM scheduled_reports
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list scheduled_reports by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list scheduled_reports: %w", err)
	}
	defer rows.Close()

	var entities []*ScheduledReports
	for rows.Next() {
		var entity ScheduledReports
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReportName,
			&entity.ReportType,
			&entity.ScheduleFrequency,
			&entity.ScheduleDayOfWeek,
			&entity.ScheduleDayOfMonth,
			&entity.ScheduleTime,
			&entity.ScheduleTimezone,
			&entity.ReportParameters,
			&entity.DeliveryMethod,
			&entity.DeliveryRecipients,
			&entity.OutputFormat,
			&entity.IsActive,
			&entity.LastRunAt,
			&entity.LastRunStatus,
			&entity.NextRunAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan scheduled_reports: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

