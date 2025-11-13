package background_job

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

// Repository handles database operations for BackgroundJobs
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BackgroundJobs repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BackgroundJobs represents a background_jobs entity
type BackgroundJobs struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	JobType string `json:"job_type" db:"job_type"`
	JobName string `json:"job_name" db:"job_name"`
	QueueName *string `json:"queue_name" db:"queue_name"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	Payload json.RawMessage `json:"payload" db:"payload"`
	Result json.RawMessage `json:"result" db:"result"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	ErrorDetails json.RawMessage `json:"error_details" db:"error_details"`
	Attempts *int64 `json:"attempts" db:"attempts"`
	MaxAttempts *int64 `json:"max_attempts" db:"max_attempts"`
	Priority *int64 `json:"priority" db:"priority"`
	ScheduledAt *time.Time `json:"scheduled_at" db:"scheduled_at"`
	StartedAt *time.Time `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	FailedAt *time.Time `json:"failed_at" db:"failed_at"`
	WorkerId *string `json:"worker_id" db:"worker_id"`
	ProcessingTimeout *int64 `json:"processing_timeout" db:"processing_timeout"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new background_jobs record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BackgroundJobs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "background_jobs", duration, nil)
	}()

	query := `
		INSERT INTO background_jobs (
			, organization_id
			, job_type
			, job_name
			, queue_name
			, status
			, payload
			, result
			, error_message
			, error_details
			, attempts
			, max_attempts
			, priority
			, scheduled_at
			, started_at
			, completed_at
			, failed_at
			, worker_id
			, processing_timeout
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
			, $19
			, $20
			, $21
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.JobType,
		entity.JobName,
		entity.QueueName,
		entity.Status,
		entity.Status,
		entity.Payload,
		entity.Result,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.Attempts,
		entity.MaxAttempts,
		entity.Priority,
		entity.ScheduledAt,
		entity.StartedAt,
		entity.CompletedAt,
		entity.FailedAt,
		entity.WorkerId,
		entity.ProcessingTimeout,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create background_jobs", zap.Error(err))
		return fmt.Errorf("failed to create background_jobs: %w", err)
	}

	r.logger.Info("created background_jobs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a background_jobs by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BackgroundJobs, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "background_jobs", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, job_type
			, job_name
			, queue_name
			, payload
			, result
			, error_message
			, error_details
			, attempts
			, max_attempts
			, priority
			, scheduled_at
			, started_at
			, completed_at
			, failed_at
			, worker_id
			, processing_timeout
			, created_by
			, created_at
			, updated_at
		FROM background_jobs
		WHERE id = $1
		
	`

	var entity BackgroundJobs
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.JobType,
		&entity.JobName,
		&entity.QueueName,
		&entity.Status,
		&entity.Status,
		&entity.Payload,
		&entity.Result,
		&entity.ErrorMessage,
		&entity.ErrorDetails,
		&entity.Attempts,
		&entity.MaxAttempts,
		&entity.Priority,
		&entity.ScheduledAt,
		&entity.StartedAt,
		&entity.CompletedAt,
		&entity.FailedAt,
		&entity.WorkerId,
		&entity.ProcessingTimeout,
		&entity.CreatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("background_jobs not found")
	}

	if err != nil {
		r.logger.Error("failed to get background_jobs", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get background_jobs: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of background_jobs records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BackgroundJobs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "background_jobs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM background_jobs
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count background_jobs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, job_type
			, job_name
			, queue_name
			, payload
			, result
			, error_message
			, error_details
			, attempts
			, max_attempts
			, priority
			, scheduled_at
			, started_at
			, completed_at
			, failed_at
			, worker_id
			, processing_timeout
			, created_by
			, created_at
			, updated_at
		FROM background_jobs
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list background_jobs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list background_jobs: %w", err)
	}
	defer rows.Close()

	var entities []*BackgroundJobs
	for rows.Next() {
		var entity BackgroundJobs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JobType,
			&entity.JobName,
			&entity.QueueName,
			&entity.Status,
			&entity.Status,
			&entity.Payload,
			&entity.Result,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.Priority,
			&entity.ScheduledAt,
			&entity.StartedAt,
			&entity.CompletedAt,
			&entity.FailedAt,
			&entity.WorkerId,
			&entity.ProcessingTimeout,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan background_jobs: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating background_jobs rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing background_jobs record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BackgroundJobs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "background_jobs", duration, nil)
	}()

	query := `
		UPDATE background_jobs
		SET
			, organization_id = $2
			, job_type = $3
			, job_name = $4
			, queue_name = $5
			, status = $6
			, payload = $8
			, result = $9
			, error_message = $10
			, error_details = $11
			, attempts = $12
			, max_attempts = $13
			, priority = $14
			, scheduled_at = $15
			, started_at = $16
			, completed_at = $17
			, failed_at = $18
			, worker_id = $19
			, processing_timeout = $20
			, created_by = $21
			, updated_at = $23
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $24
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.JobType,
		entity.JobName,
		entity.QueueName,
		entity.Status,
		entity.Status,
		entity.Payload,
		entity.Result,
		entity.ErrorMessage,
		entity.ErrorDetails,
		entity.Attempts,
		entity.MaxAttempts,
		entity.Priority,
		entity.ScheduledAt,
		entity.StartedAt,
		entity.CompletedAt,
		entity.FailedAt,
		entity.WorkerId,
		entity.ProcessingTimeout,
		entity.CreatedBy,
		time.Now(),
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update background_jobs", zap.Error(err))
		return fmt.Errorf("failed to update background_jobs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("background_jobs not found or already deleted")
	}

	r.logger.Info("updated background_jobs",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a background_jobs record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "background_jobs", duration, nil)
	}()

	query := `DELETE FROM background_jobs WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete background_jobs", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete background_jobs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("background_jobs not found")
	}

	r.logger.Info("deleted background_jobs", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves background_jobs records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BackgroundJobs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "background_jobs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM background_jobs
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count background_jobs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, job_type
			, job_name
			, queue_name
			, payload
			, result
			, error_message
			, error_details
			, attempts
			, max_attempts
			, priority
			, scheduled_at
			, started_at
			, completed_at
			, failed_at
			, worker_id
			, processing_timeout
			, created_by
			, created_at
			, updated_at
		FROM background_jobs
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list background_jobs by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list background_jobs: %w", err)
	}
	defer rows.Close()

	var entities []*BackgroundJobs
	for rows.Next() {
		var entity BackgroundJobs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.JobType,
			&entity.JobName,
			&entity.QueueName,
			&entity.Status,
			&entity.Status,
			&entity.Payload,
			&entity.Result,
			&entity.ErrorMessage,
			&entity.ErrorDetails,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.Priority,
			&entity.ScheduledAt,
			&entity.StartedAt,
			&entity.CompletedAt,
			&entity.FailedAt,
			&entity.WorkerId,
			&entity.ProcessingTimeout,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan background_jobs: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

