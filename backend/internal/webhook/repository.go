package webhook

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

// Repository handles database operations for Webhooks
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Webhooks repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Webhooks represents a webhooks entity
type Webhooks struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	WebhookName string `json:"webhook_name" db:"webhook_name"`
	Url string `json:"url" db:"url"`
	Secret *string `json:"secret" db:"secret"`
	Events string `json:"events" db:"events"`
	HttpMethod *string `json:"http_method" db:"http_method"`
	Headers json.RawMessage `json:"headers" db:"headers"`
	TimeoutSeconds *int64 `json:"timeout_seconds" db:"timeout_seconds"`
	MaxRetries *int64 `json:"max_retries" db:"max_retries"`
	RetryBackoffSeconds *int64 `json:"retry_backoff_seconds" db:"retry_backoff_seconds"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsVerified *bool `json:"is_verified" db:"is_verified"`
	TotalDeliveries *int64 `json:"total_deliveries" db:"total_deliveries"`
	SuccessfulDeliveries *int64 `json:"successful_deliveries" db:"successful_deliveries"`
	FailedDeliveries *int64 `json:"failed_deliveries" db:"failed_deliveries"`
	LastDeliveryAt *time.Time `json:"last_delivery_at" db:"last_delivery_at"`
	LastSuccessAt *time.Time `json:"last_success_at" db:"last_success_at"`
	LastFailureAt *time.Time `json:"last_failure_at" db:"last_failure_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new webhooks record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Webhooks) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "webhooks", duration, nil)
	}()

	query := `
		INSERT INTO webhooks (
			, organization_id
			, webhook_name
			, url
			, secret
			, events
			, http_method
			, headers
			, timeout_seconds
			, max_retries
			, retry_backoff_seconds
			, is_active
			, is_verified
			, total_deliveries
			, successful_deliveries
			, failed_deliveries
			, last_delivery_at
			, last_success_at
			, last_failure_at
			, created_by
			, deleted_at
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
			, $23
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.WebhookName,
		entity.Url,
		entity.Secret,
		entity.Events,
		entity.HttpMethod,
		entity.Headers,
		entity.TimeoutSeconds,
		entity.MaxRetries,
		entity.RetryBackoffSeconds,
		entity.IsActive,
		entity.IsVerified,
		entity.TotalDeliveries,
		entity.SuccessfulDeliveries,
		entity.FailedDeliveries,
		entity.LastDeliveryAt,
		entity.LastSuccessAt,
		entity.LastFailureAt,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create webhooks", zap.Error(err))
		return fmt.Errorf("failed to create webhooks: %w", err)
	}

	r.logger.Info("created webhooks",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a webhooks by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Webhooks, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhooks", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, webhook_name
			, url
			, secret
			, events
			, http_method
			, headers
			, timeout_seconds
			, max_retries
			, retry_backoff_seconds
			, is_active
			, is_verified
			, total_deliveries
			, successful_deliveries
			, failed_deliveries
			, last_delivery_at
			, last_success_at
			, last_failure_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM webhooks
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Webhooks
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.WebhookName,
		&entity.Url,
		&entity.Secret,
		&entity.Events,
		&entity.HttpMethod,
		&entity.Headers,
		&entity.TimeoutSeconds,
		&entity.MaxRetries,
		&entity.RetryBackoffSeconds,
		&entity.IsActive,
		&entity.IsVerified,
		&entity.TotalDeliveries,
		&entity.SuccessfulDeliveries,
		&entity.FailedDeliveries,
		&entity.LastDeliveryAt,
		&entity.LastSuccessAt,
		&entity.LastFailureAt,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("webhooks not found")
	}

	if err != nil {
		r.logger.Error("failed to get webhooks", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get webhooks: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of webhooks records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Webhooks, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhooks", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM webhooks
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count webhooks records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, webhook_name
			, url
			, secret
			, events
			, http_method
			, headers
			, timeout_seconds
			, max_retries
			, retry_backoff_seconds
			, is_active
			, is_verified
			, total_deliveries
			, successful_deliveries
			, failed_deliveries
			, last_delivery_at
			, last_success_at
			, last_failure_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM webhooks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list webhooks", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var entities []*Webhooks
	for rows.Next() {
		var entity Webhooks
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.WebhookName,
			&entity.Url,
			&entity.Secret,
			&entity.Events,
			&entity.HttpMethod,
			&entity.Headers,
			&entity.TimeoutSeconds,
			&entity.MaxRetries,
			&entity.RetryBackoffSeconds,
			&entity.IsActive,
			&entity.IsVerified,
			&entity.TotalDeliveries,
			&entity.SuccessfulDeliveries,
			&entity.FailedDeliveries,
			&entity.LastDeliveryAt,
			&entity.LastSuccessAt,
			&entity.LastFailureAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan webhooks: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating webhooks rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing webhooks record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Webhooks) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "webhooks", duration, nil)
	}()

	query := `
		UPDATE webhooks
		SET
			, organization_id = $2
			, webhook_name = $3
			, url = $4
			, secret = $5
			, events = $6
			, http_method = $7
			, headers = $8
			, timeout_seconds = $9
			, max_retries = $10
			, retry_backoff_seconds = $11
			, is_active = $12
			, is_verified = $13
			, total_deliveries = $14
			, successful_deliveries = $15
			, failed_deliveries = $16
			, last_delivery_at = $17
			, last_success_at = $18
			, last_failure_at = $19
			, created_by = $20
			, updated_at = $22
			, deleted_at = $23
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $24
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.WebhookName,
		entity.Url,
		entity.Secret,
		entity.Events,
		entity.HttpMethod,
		entity.Headers,
		entity.TimeoutSeconds,
		entity.MaxRetries,
		entity.RetryBackoffSeconds,
		entity.IsActive,
		entity.IsVerified,
		entity.TotalDeliveries,
		entity.SuccessfulDeliveries,
		entity.FailedDeliveries,
		entity.LastDeliveryAt,
		entity.LastSuccessAt,
		entity.LastFailureAt,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update webhooks", zap.Error(err))
		return fmt.Errorf("failed to update webhooks: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("webhooks not found or already deleted")
	}

	r.logger.Info("updated webhooks",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a webhooks record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "webhooks", duration, nil)
	}()

	query := `
		UPDATE webhooks
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete webhooks", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete webhooks: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("webhooks not found or already deleted")
	}

	r.logger.Info("deleted webhooks", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves webhooks records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Webhooks, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "webhooks", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM webhooks
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count webhooks records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, webhook_name
			, url
			, secret
			, events
			, http_method
			, headers
			, timeout_seconds
			, max_retries
			, retry_backoff_seconds
			, is_active
			, is_verified
			, total_deliveries
			, successful_deliveries
			, failed_deliveries
			, last_delivery_at
			, last_success_at
			, last_failure_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM webhooks
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list webhooks by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list webhooks: %w", err)
	}
	defer rows.Close()

	var entities []*Webhooks
	for rows.Next() {
		var entity Webhooks
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.WebhookName,
			&entity.Url,
			&entity.Secret,
			&entity.Events,
			&entity.HttpMethod,
			&entity.Headers,
			&entity.TimeoutSeconds,
			&entity.MaxRetries,
			&entity.RetryBackoffSeconds,
			&entity.IsActive,
			&entity.IsVerified,
			&entity.TotalDeliveries,
			&entity.SuccessfulDeliveries,
			&entity.FailedDeliveries,
			&entity.LastDeliveryAt,
			&entity.LastSuccessAt,
			&entity.LastFailureAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan webhooks: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

