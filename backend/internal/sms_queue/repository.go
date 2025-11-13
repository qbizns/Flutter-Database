package sms_queue

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

// Repository handles database operations for SmsQueue
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SmsQueue repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SmsQueue represents a sms_queue entity
type SmsQueue struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	ToPhone string `json:"to_phone" db:"to_phone"`
	FromPhone *string `json:"from_phone" db:"from_phone"`
	Message string `json:"message" db:"message"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	Provider *string `json:"provider" db:"provider"`
	ProviderMessageId *string `json:"provider_message_id" db:"provider_message_id"`
	Attempts *int64 `json:"attempts" db:"attempts"`
	MaxAttempts *int64 `json:"max_attempts" db:"max_attempts"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	CostAmount *float64 `json:"cost_amount" db:"cost_amount"`
	CostCurrency *string `json:"cost_currency" db:"cost_currency"`
	ScheduledAt *time.Time `json:"scheduled_at" db:"scheduled_at"`
	SentAt *time.Time `json:"sent_at" db:"sent_at"`
	FailedAt *time.Time `json:"failed_at" db:"failed_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new sms_queue record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SmsQueue) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sms_queue", duration, nil)
	}()

	query := `
		INSERT INTO sms_queue (
			, organization_id
			, to_phone
			, from_phone
			, message
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, cost_amount
			, cost_currency
			, scheduled_at
			, sent_at
			, failed_at
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ToPhone,
		entity.FromPhone,
		entity.Message,
		entity.Status,
		entity.Status,
		entity.Provider,
		entity.ProviderMessageId,
		entity.Attempts,
		entity.MaxAttempts,
		entity.ErrorMessage,
		entity.CostAmount,
		entity.CostCurrency,
		entity.ScheduledAt,
		entity.SentAt,
		entity.FailedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create sms_queue", zap.Error(err))
		return fmt.Errorf("failed to create sms_queue: %w", err)
	}

	r.logger.Info("created sms_queue",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a sms_queue by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SmsQueue, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sms_queue", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, to_phone
			, from_phone
			, message
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, cost_amount
			, cost_currency
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM sms_queue
		WHERE id = $1
		
	`

	var entity SmsQueue
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ToPhone,
		&entity.FromPhone,
		&entity.Message,
		&entity.Status,
		&entity.Status,
		&entity.Provider,
		&entity.ProviderMessageId,
		&entity.Attempts,
		&entity.MaxAttempts,
		&entity.ErrorMessage,
		&entity.CostAmount,
		&entity.CostCurrency,
		&entity.ScheduledAt,
		&entity.SentAt,
		&entity.FailedAt,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sms_queue not found")
	}

	if err != nil {
		r.logger.Error("failed to get sms_queue", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sms_queue: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sms_queue records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SmsQueue, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sms_queue", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sms_queue
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sms_queue records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, to_phone
			, from_phone
			, message
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, cost_amount
			, cost_currency
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM sms_queue
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sms_queue", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sms_queue: %w", err)
	}
	defer rows.Close()

	var entities []*SmsQueue
	for rows.Next() {
		var entity SmsQueue
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ToPhone,
			&entity.FromPhone,
			&entity.Message,
			&entity.Status,
			&entity.Status,
			&entity.Provider,
			&entity.ProviderMessageId,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.ErrorMessage,
			&entity.CostAmount,
			&entity.CostCurrency,
			&entity.ScheduledAt,
			&entity.SentAt,
			&entity.FailedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sms_queue: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sms_queue rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sms_queue record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SmsQueue) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sms_queue", duration, nil)
	}()

	query := `
		UPDATE sms_queue
		SET
			, organization_id = $2
			, to_phone = $3
			, from_phone = $4
			, message = $5
			, status = $6
			, status = $7
			, provider = $8
			, provider_message_id = $9
			, attempts = $10
			, max_attempts = $11
			, error_message = $12
			, cost_amount = $13
			, cost_currency = $14
			, scheduled_at = $15
			, sent_at = $16
			, failed_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ToPhone,
		entity.FromPhone,
		entity.Message,
		entity.Status,
		entity.Status,
		entity.Provider,
		entity.ProviderMessageId,
		entity.Attempts,
		entity.MaxAttempts,
		entity.ErrorMessage,
		entity.CostAmount,
		entity.CostCurrency,
		entity.ScheduledAt,
		entity.SentAt,
		entity.FailedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sms_queue", zap.Error(err))
		return fmt.Errorf("failed to update sms_queue: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sms_queue not found or already deleted")
	}

	r.logger.Info("updated sms_queue",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a sms_queue record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "sms_queue", duration, nil)
	}()

	query := `DELETE FROM sms_queue WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sms_queue", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sms_queue: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sms_queue not found")
	}

	r.logger.Info("deleted sms_queue", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sms_queue records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SmsQueue, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sms_queue", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sms_queue
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sms_queue records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, to_phone
			, from_phone
			, message
			, status
			, status
			, provider
			, provider_message_id
			, attempts
			, max_attempts
			, error_message
			, cost_amount
			, cost_currency
			, scheduled_at
			, sent_at
			, failed_at
			, created_at
		FROM sms_queue
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sms_queue by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sms_queue: %w", err)
	}
	defer rows.Close()

	var entities []*SmsQueue
	for rows.Next() {
		var entity SmsQueue
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ToPhone,
			&entity.FromPhone,
			&entity.Message,
			&entity.Status,
			&entity.Status,
			&entity.Provider,
			&entity.ProviderMessageId,
			&entity.Attempts,
			&entity.MaxAttempts,
			&entity.ErrorMessage,
			&entity.CostAmount,
			&entity.CostCurrency,
			&entity.ScheduledAt,
			&entity.SentAt,
			&entity.FailedAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sms_queue: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

