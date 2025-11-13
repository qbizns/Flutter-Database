package webhook_delivery

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

// Repository handles database operations for WebhookDeliveries
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new WebhookDeliveries repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// WebhookDeliveries represents a webhook_deliveries entity
type WebhookDeliveries struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	WebhookId uuid.UUID `json:"webhook_id" db:"webhook_id"`
	EventType string `json:"event_type" db:"event_type"`
	EventId uuid.UUID `json:"event_id" db:"event_id"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	RequestUrl string `json:"request_url" db:"request_url"`
	RequestMethod string `json:"request_method" db:"request_method"`
	RequestHeaders json.RawMessage `json:"request_headers" db:"request_headers"`
	RequestBody json.RawMessage `json:"request_body" db:"request_body"`
	ResponseStatusCode *int64 `json:"response_status_code" db:"response_status_code"`
	ResponseHeaders json.RawMessage `json:"response_headers" db:"response_headers"`
	ResponseBody *string `json:"response_body" db:"response_body"`
	AttemptNumber *int64 `json:"attempt_number" db:"attempt_number"`
	DurationMs *int64 `json:"duration_ms" db:"duration_ms"`
	NextRetryAt *time.Time `json:"next_retry_at" db:"next_retry_at"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeliveredAt *time.Time `json:"delivered_at" db:"delivered_at"`
}

// Create inserts a new webhook_deliveries record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *WebhookDeliveries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "webhook_deliveries", duration, nil)
	}()

	query := `
		INSERT INTO webhook_deliveries (
			, organization_id
			, webhook_id
			, event_type
			, event_id
			, status
			, request_url
			, request_method
			, request_headers
			, request_body
			, response_status_code
			, response_headers
			, response_body
			, attempt_number
			, duration_ms
			, next_retry_at
			, error_message
			, delivered_at
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
			, $20
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.WebhookId,
		entity.EventType,
		entity.EventId,
		entity.Status,
		entity.Status,
		entity.RequestUrl,
		entity.RequestMethod,
		entity.RequestHeaders,
		entity.RequestBody,
		entity.ResponseStatusCode,
		entity.ResponseHeaders,
		entity.ResponseBody,
		entity.AttemptNumber,
		entity.DurationMs,
		entity.NextRetryAt,
		entity.ErrorMessage,
		entity.DeliveredAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create webhook_deliveries", zap.Error(err))
		return fmt.Errorf("failed to create webhook_deliveries: %w", err)
	}

	r.logger.Info("created webhook_deliveries",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a webhook_deliveries by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*WebhookDeliveries, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhook_deliveries", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, webhook_id
			, event_type
			, event_id
			, request_url
			, request_method
			, request_headers
			, request_body
			, response_status_code
			, response_headers
			, response_body
			, attempt_number
			, duration_ms
			, next_retry_at
			, error_message
			, created_at
			, delivered_at
		FROM webhook_deliveries
		WHERE id = $1
		
	`

	var entity WebhookDeliveries
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.WebhookId,
		&entity.EventType,
		&entity.EventId,
		&entity.Status,
		&entity.Status,
		&entity.RequestUrl,
		&entity.RequestMethod,
		&entity.RequestHeaders,
		&entity.RequestBody,
		&entity.ResponseStatusCode,
		&entity.ResponseHeaders,
		&entity.ResponseBody,
		&entity.AttemptNumber,
		&entity.DurationMs,
		&entity.NextRetryAt,
		&entity.ErrorMessage,
		&entity.CreatedAt,
		&entity.DeliveredAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("webhook_deliveries not found")
	}

	if err != nil {
		r.logger.Error("failed to get webhook_deliveries", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get webhook_deliveries: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of webhook_deliveries records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*WebhookDeliveries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhook_deliveries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM webhook_deliveries
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count webhook_deliveries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, webhook_id
			, event_type
			, event_id
			, request_url
			, request_method
			, request_headers
			, request_body
			, response_status_code
			, response_headers
			, response_body
			, attempt_number
			, duration_ms
			, next_retry_at
			, error_message
			, created_at
			, delivered_at
		FROM webhook_deliveries
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list webhook_deliveries", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list webhook_deliveries: %w", err)
	}
	defer rows.Close()

	var entities []*WebhookDeliveries
	for rows.Next() {
		var entity WebhookDeliveries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.WebhookId,
			&entity.EventType,
			&entity.EventId,
			&entity.Status,
			&entity.Status,
			&entity.RequestUrl,
			&entity.RequestMethod,
			&entity.RequestHeaders,
			&entity.RequestBody,
			&entity.ResponseStatusCode,
			&entity.ResponseHeaders,
			&entity.ResponseBody,
			&entity.AttemptNumber,
			&entity.DurationMs,
			&entity.NextRetryAt,
			&entity.ErrorMessage,
			&entity.CreatedAt,
			&entity.DeliveredAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan webhook_deliveries: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating webhook_deliveries rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing webhook_deliveries record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *WebhookDeliveries) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "webhook_deliveries", duration, nil)
	}()

	query := `
		UPDATE webhook_deliveries
		SET
			, organization_id = $2
			, webhook_id = $3
			, event_type = $4
			, event_id = $5
			, status = $6
			, request_url = $8
			, request_method = $9
			, request_headers = $10
			, request_body = $11
			, response_status_code = $12
			, response_headers = $13
			, response_body = $14
			, attempt_number = $15
			, duration_ms = $16
			, next_retry_at = $17
			, error_message = $18
			, delivered_at = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.WebhookId,
		entity.EventType,
		entity.EventId,
		entity.Status,
		entity.Status,
		entity.RequestUrl,
		entity.RequestMethod,
		entity.RequestHeaders,
		entity.RequestBody,
		entity.ResponseStatusCode,
		entity.ResponseHeaders,
		entity.ResponseBody,
		entity.AttemptNumber,
		entity.DurationMs,
		entity.NextRetryAt,
		entity.ErrorMessage,
		entity.DeliveredAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update webhook_deliveries", zap.Error(err))
		return fmt.Errorf("failed to update webhook_deliveries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("webhook_deliveries not found or already deleted")
	}

	r.logger.Info("updated webhook_deliveries",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a webhook_deliveries record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "webhook_deliveries", duration, nil)
	}()

	query := `DELETE FROM webhook_deliveries WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete webhook_deliveries", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete webhook_deliveries: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("webhook_deliveries not found")
	}

	r.logger.Info("deleted webhook_deliveries", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves webhook_deliveries records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*WebhookDeliveries, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhook_deliveries", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM webhook_deliveries
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count webhook_deliveries records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, webhook_id
			, event_type
			, event_id
			, request_url
			, request_method
			, request_headers
			, request_body
			, response_status_code
			, response_headers
			, response_body
			, attempt_number
			, duration_ms
			, next_retry_at
			, error_message
			, created_at
			, delivered_at
		FROM webhook_deliveries
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list webhook_deliveries by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list webhook_deliveries: %w", err)
	}
	defer rows.Close()

	var entities []*WebhookDeliveries
	for rows.Next() {
		var entity WebhookDeliveries
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.WebhookId,
			&entity.EventType,
			&entity.EventId,
			&entity.Status,
			&entity.Status,
			&entity.RequestUrl,
			&entity.RequestMethod,
			&entity.RequestHeaders,
			&entity.RequestBody,
			&entity.ResponseStatusCode,
			&entity.ResponseHeaders,
			&entity.ResponseBody,
			&entity.AttemptNumber,
			&entity.DurationMs,
			&entity.NextRetryAt,
			&entity.ErrorMessage,
			&entity.CreatedAt,
			&entity.DeliveredAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan webhook_deliveries: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

