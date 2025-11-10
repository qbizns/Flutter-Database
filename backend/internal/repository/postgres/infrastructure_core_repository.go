package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"pos-backend/internal/domain/infrastructure"
)

// ============================================================================
// BACKGROUND JOBS REPOSITORY
// ============================================================================

// ListBackgroundJobs returns a list of background jobs
func (r *PostgresDB) ListBackgroundJobs(ctx context.Context, orgID uuid.UUID, status *string, limit int, offset int) ([]infrastructure.BackgroundJob, error) {
	query := `
		SELECT id, organization_id, job_type, job_name, queue_name, status,
		       payload, result, error_message, error_details, attempts, max_attempts,
		       priority, scheduled_at, started_at, completed_at, failed_at,
		       worker_id, processing_timeout, created_by, created_at, updated_at
		FROM background_jobs
		WHERE organization_id = $1
	`
	args := []interface{}{orgID}

	if status != nil {
		query += " AND status = $2"
		args = append(args, *status)
	}

	query += " ORDER BY scheduled_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []infrastructure.BackgroundJob
	for rows.Next() {
		var job infrastructure.BackgroundJob
		err := rows.Scan(
			&job.ID, &job.OrganizationID, &job.JobType, &job.JobName, &job.QueueName,
			&job.Status, &job.Payload, &job.Result, &job.ErrorMessage, &job.ErrorDetails,
			&job.Attempts, &job.MaxAttempts, &job.Priority, &job.ScheduledAt,
			&job.StartedAt, &job.CompletedAt, &job.FailedAt, &job.WorkerID,
			&job.ProcessingTimeout, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	return jobs, rows.Err()
}

// CreateBackgroundJob creates a new background job
func (r *PostgresDB) CreateBackgroundJob(ctx context.Context, job *infrastructure.BackgroundJob) error {
	job.ID = uuid.New()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	query := `
		INSERT INTO background_jobs (
			id, organization_id, job_type, job_name, queue_name, status,
			payload, result, error_message, error_details, attempts, max_attempts,
			priority, scheduled_at, started_at, completed_at, failed_at,
			worker_id, processing_timeout, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	return r.db.QueryRowContext(ctx, query,
		job.ID, job.OrganizationID, job.JobType, job.JobName, job.QueueName,
		job.Status, job.Payload, job.Result, job.ErrorMessage, job.ErrorDetails,
		job.Attempts, job.MaxAttempts, job.Priority, job.ScheduledAt,
		job.StartedAt, job.CompletedAt, job.FailedAt, job.WorkerID,
		job.ProcessingTimeout, job.CreatedBy, job.CreatedAt, job.UpdatedAt,
	).Scan(&job.ID)
}

// GetBackgroundJob retrieves a background job by ID
func (r *PostgresDB) GetBackgroundJob(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.BackgroundJob, error) {
	query := `
		SELECT id, organization_id, job_type, job_name, queue_name, status,
		       payload, result, error_message, error_details, attempts, max_attempts,
		       priority, scheduled_at, started_at, completed_at, failed_at,
		       worker_id, processing_timeout, created_by, created_at, updated_at
		FROM background_jobs
		WHERE id = $1 AND organization_id = $2
	`

	job := &infrastructure.BackgroundJob{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&job.ID, &job.OrganizationID, &job.JobType, &job.JobName, &job.QueueName,
		&job.Status, &job.Payload, &job.Result, &job.ErrorMessage, &job.ErrorDetails,
		&job.Attempts, &job.MaxAttempts, &job.Priority, &job.ScheduledAt,
		&job.StartedAt, &job.CompletedAt, &job.FailedAt, &job.WorkerID,
		&job.ProcessingTimeout, &job.CreatedBy, &job.CreatedAt, &job.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return job, nil
}

// UpdateBackgroundJob updates a background job
func (r *PostgresDB) UpdateBackgroundJob(ctx context.Context, job *infrastructure.BackgroundJob) error {
	job.UpdatedAt = time.Now()

	query := `
		UPDATE background_jobs
		SET job_type = $1, job_name = $2, queue_name = $3, status = $4,
		    payload = $5, result = $6, error_message = $7, error_details = $8,
		    attempts = $9, max_attempts = $10, priority = $11, scheduled_at = $12,
		    started_at = $13, completed_at = $14, failed_at = $15, worker_id = $16,
		    processing_timeout = $17, updated_at = $18
		WHERE id = $19 AND organization_id = $20
	`

	_, err := r.db.ExecContext(ctx, query,
		job.JobType, job.JobName, job.QueueName, job.Status, job.Payload, job.Result,
		job.ErrorMessage, job.ErrorDetails, job.Attempts, job.MaxAttempts, job.Priority,
		job.ScheduledAt, job.StartedAt, job.CompletedAt, job.FailedAt, job.WorkerID,
		job.ProcessingTimeout, job.UpdatedAt, job.ID, job.OrganizationID,
	)

	return err
}

// UpdateBackgroundJobStatus updates only the status of a background job
func (r *PostgresDB) UpdateBackgroundJobStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	query := `
		UPDATE background_jobs
		SET status = $1, updated_at = $2
		WHERE id = $3 AND organization_id = $4
	`

	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id, orgID)
	return err
}

// ============================================================================
// API KEYS REPOSITORY
// ============================================================================

// ListAPIKeys returns a list of API keys
func (r *PostgresDB) ListAPIKeys(ctx context.Context, orgID uuid.UUID, isActive *bool, limit int, offset int) ([]infrastructure.APIKey, error) {
	query := `
		SELECT id, organization_id, key_name, key_prefix, key_hash, scopes,
		       allowed_ips, is_active, last_used_at, usage_count, rate_limit_per_minute,
		       rate_limit_per_hour, expires_at, created_by, created_at, updated_at, deleted_at
		FROM api_keys
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}

	if isActive != nil {
		query += " AND is_active = $2"
		args = append(args, *isActive)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []infrastructure.APIKey
	for rows.Next() {
		var key infrastructure.APIKey
		err := rows.Scan(
			&key.ID, &key.OrganizationID, &key.KeyName, &key.KeyPrefix, &key.KeyHash,
			&key.Scopes, &key.AllowedIPs, &key.IsActive, &key.LastUsedAt, &key.UsageCount,
			&key.RateLimitPerMinute, &key.RateLimitPerHour, &key.ExpiresAt,
			&key.CreatedBy, &key.CreatedAt, &key.UpdatedAt, &key.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, rows.Err()
}

// CreateAPIKey creates a new API key
func (r *PostgresDB) CreateAPIKey(ctx context.Context, key *infrastructure.APIKey) error {
	key.ID = uuid.New()
	key.CreatedAt = time.Now()
	key.UpdatedAt = time.Now()

	query := `
		INSERT INTO api_keys (
			id, organization_id, key_name, key_prefix, key_hash, scopes,
			allowed_ips, is_active, last_used_at, usage_count, rate_limit_per_minute,
			rate_limit_per_hour, expires_at, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	return r.db.QueryRowContext(ctx, query,
		key.ID, key.OrganizationID, key.KeyName, key.KeyPrefix, key.KeyHash,
		pq.Array(key.Scopes), pq.Array(key.AllowedIPs), key.IsActive, key.LastUsedAt,
		key.UsageCount, key.RateLimitPerMinute, key.RateLimitPerHour, key.ExpiresAt,
		key.CreatedBy, key.CreatedAt, key.UpdatedAt,
	).Scan(&key.ID)
}

// GetAPIKey retrieves an API key by ID
func (r *PostgresDB) GetAPIKey(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.APIKey, error) {
	query := `
		SELECT id, organization_id, key_name, key_prefix, key_hash, scopes,
		       allowed_ips, is_active, last_used_at, usage_count, rate_limit_per_minute,
		       rate_limit_per_hour, expires_at, created_by, created_at, updated_at, deleted_at
		FROM api_keys
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	key := &infrastructure.APIKey{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&key.ID, &key.OrganizationID, &key.KeyName, &key.KeyPrefix, &key.KeyHash,
		&key.Scopes, &key.AllowedIPs, &key.IsActive, &key.LastUsedAt, &key.UsageCount,
		&key.RateLimitPerMinute, &key.RateLimitPerHour, &key.ExpiresAt,
		&key.CreatedBy, &key.CreatedAt, &key.UpdatedAt, &key.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GetAPIKeyByHash retrieves an API key by its hash
func (r *PostgresDB) GetAPIKeyByHash(ctx context.Context, keyHash string) (*infrastructure.APIKey, error) {
	query := `
		SELECT id, organization_id, key_name, key_prefix, key_hash, scopes,
		       allowed_ips, is_active, last_used_at, usage_count, rate_limit_per_minute,
		       rate_limit_per_hour, expires_at, created_by, created_at, updated_at, deleted_at
		FROM api_keys
		WHERE key_hash = $1 AND deleted_at IS NULL
	`

	key := &infrastructure.APIKey{}
	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&key.ID, &key.OrganizationID, &key.KeyName, &key.KeyPrefix, &key.KeyHash,
		&key.Scopes, &key.AllowedIPs, &key.IsActive, &key.LastUsedAt, &key.UsageCount,
		&key.RateLimitPerMinute, &key.RateLimitPerHour, &key.ExpiresAt,
		&key.CreatedBy, &key.CreatedAt, &key.UpdatedAt, &key.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return key, nil
}

// UpdateAPIKey updates an API key
func (r *PostgresDB) UpdateAPIKey(ctx context.Context, key *infrastructure.APIKey) error {
	key.UpdatedAt = time.Now()

	query := `
		UPDATE api_keys
		SET key_name = $1, scopes = $2, allowed_ips = $3, is_active = $4,
		    rate_limit_per_minute = $5, rate_limit_per_hour = $6, expires_at = $7, updated_at = $8
		WHERE id = $9 AND organization_id = $10
	`

	_, err := r.db.ExecContext(ctx, query,
		key.KeyName, pq.Array(key.Scopes), pq.Array(key.AllowedIPs), key.IsActive,
		key.RateLimitPerMinute, key.RateLimitPerHour, key.ExpiresAt, key.UpdatedAt,
		key.ID, key.OrganizationID,
	)

	return err
}

// UpdateAPIKeyUsage updates the usage count and last used time
func (r *PostgresDB) UpdateAPIKeyUsage(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := `
		UPDATE api_keys
		SET usage_count = usage_count + 1, last_used_at = $1, updated_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, orgID)
	return err
}

// DeleteAPIKey soft-deletes an API key
func (r *PostgresDB) DeleteAPIKey(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := `
		UPDATE api_keys
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, orgID)
	return err
}

// ============================================================================
// WEBHOOKS REPOSITORY
// ============================================================================

// ListWebhooks returns a list of webhooks
func (r *PostgresDB) ListWebhooks(ctx context.Context, orgID uuid.UUID, isActive *bool, limit int, offset int) ([]infrastructure.Webhook, error) {
	query := `
		SELECT id, organization_id, webhook_name, url, secret, events,
		       http_method, headers, timeout_seconds, max_retries, retry_backoff_seconds,
		       is_active, is_verified, total_deliveries, successful_deliveries, failed_deliveries,
		       last_delivery_at, last_success_at, last_failure_at, created_by,
		       created_at, updated_at, deleted_at
		FROM webhooks
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{orgID}

	if isActive != nil {
		query += " AND is_active = $2"
		args = append(args, *isActive)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var webhooks []infrastructure.Webhook
	for rows.Next() {
		var webhook infrastructure.Webhook
		err := rows.Scan(
			&webhook.ID, &webhook.OrganizationID, &webhook.WebhookName, &webhook.URL,
			&webhook.Secret, &webhook.Events, &webhook.HTTPMethod, &webhook.Headers,
			&webhook.TimeoutSeconds, &webhook.MaxRetries, &webhook.RetryBackoffSeconds,
			&webhook.IsActive, &webhook.IsVerified, &webhook.TotalDeliveries,
			&webhook.SuccessfulDeliveries, &webhook.FailedDeliveries, &webhook.LastDeliveryAt,
			&webhook.LastSuccessAt, &webhook.LastFailureAt, &webhook.CreatedBy,
			&webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}

	return webhooks, rows.Err()
}

// CreateWebhook creates a new webhook
func (r *PostgresDB) CreateWebhook(ctx context.Context, webhook *infrastructure.Webhook) error {
	webhook.ID = uuid.New()
	webhook.CreatedAt = time.Now()
	webhook.UpdatedAt = time.Now()

	query := `
		INSERT INTO webhooks (
			id, organization_id, webhook_name, url, secret, events,
			http_method, headers, timeout_seconds, max_retries, retry_backoff_seconds,
			is_active, is_verified, total_deliveries, successful_deliveries, failed_deliveries,
			last_delivery_at, last_success_at, last_failure_at, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	return r.db.QueryRowContext(ctx, query,
		webhook.ID, webhook.OrganizationID, webhook.WebhookName, webhook.URL, webhook.Secret,
		pq.Array(webhook.Events), webhook.HTTPMethod, webhook.Headers, webhook.TimeoutSeconds,
		webhook.MaxRetries, webhook.RetryBackoffSeconds, webhook.IsActive, webhook.IsVerified,
		webhook.TotalDeliveries, webhook.SuccessfulDeliveries, webhook.FailedDeliveries,
		webhook.LastDeliveryAt, webhook.LastSuccessAt, webhook.LastFailureAt,
		webhook.CreatedBy, webhook.CreatedAt, webhook.UpdatedAt,
	).Scan(&webhook.ID)
}

// GetWebhook retrieves a webhook by ID
func (r *PostgresDB) GetWebhook(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.Webhook, error) {
	query := `
		SELECT id, organization_id, webhook_name, url, secret, events,
		       http_method, headers, timeout_seconds, max_retries, retry_backoff_seconds,
		       is_active, is_verified, total_deliveries, successful_deliveries, failed_deliveries,
		       last_delivery_at, last_success_at, last_failure_at, created_by,
		       created_at, updated_at, deleted_at
		FROM webhooks
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	webhook := &infrastructure.Webhook{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&webhook.ID, &webhook.OrganizationID, &webhook.WebhookName, &webhook.URL,
		&webhook.Secret, &webhook.Events, &webhook.HTTPMethod, &webhook.Headers,
		&webhook.TimeoutSeconds, &webhook.MaxRetries, &webhook.RetryBackoffSeconds,
		&webhook.IsActive, &webhook.IsVerified, &webhook.TotalDeliveries,
		&webhook.SuccessfulDeliveries, &webhook.FailedDeliveries, &webhook.LastDeliveryAt,
		&webhook.LastSuccessAt, &webhook.LastFailureAt, &webhook.CreatedBy,
		&webhook.CreatedAt, &webhook.UpdatedAt, &webhook.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return webhook, nil
}

// UpdateWebhook updates a webhook
func (r *PostgresDB) UpdateWebhook(ctx context.Context, webhook *infrastructure.Webhook) error {
	webhook.UpdatedAt = time.Now()

	query := `
		UPDATE webhooks
		SET webhook_name = $1, url = $2, secret = $3, events = $4,
		    http_method = $5, headers = $6, timeout_seconds = $7, max_retries = $8,
		    retry_backoff_seconds = $9, is_active = $10, is_verified = $11, updated_at = $12
		WHERE id = $13 AND organization_id = $14
	`

	_, err := r.db.ExecContext(ctx, query,
		webhook.WebhookName, webhook.URL, webhook.Secret, pq.Array(webhook.Events),
		webhook.HTTPMethod, webhook.Headers, webhook.TimeoutSeconds, webhook.MaxRetries,
		webhook.RetryBackoffSeconds, webhook.IsActive, webhook.IsVerified, webhook.UpdatedAt,
		webhook.ID, webhook.OrganizationID,
	)

	return err
}

// UpdateWebhookStats updates webhook delivery statistics
func (r *PostgresDB) UpdateWebhookStats(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	query := `
		UPDATE webhooks
		SET total_deliveries = total_deliveries + 1,
		    last_delivery_at = $1,
		    updated_at = $1
	`

	if status == "success" {
		query += ", successful_deliveries = successful_deliveries + 1, last_success_at = $1"
	} else if status == "failed" {
		query += ", failed_deliveries = failed_deliveries + 1, last_failure_at = $1"
	}

	query += " WHERE id = $2 AND organization_id = $3"

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, orgID)
	return err
}

// DeleteWebhook soft-deletes a webhook
func (r *PostgresDB) DeleteWebhook(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := `
		UPDATE webhooks
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, orgID)
	return err
}

// ============================================================================
// WEBHOOK DELIVERIES REPOSITORY
// ============================================================================

// ListWebhookDeliveries returns a list of webhook deliveries
func (r *PostgresDB) ListWebhookDeliveries(ctx context.Context, orgID uuid.UUID, webhookID *uuid.UUID, status *string, limit int, offset int) ([]infrastructure.WebhookDelivery, error) {
	query := `
		SELECT id, organization_id, webhook_id, event_type, event_id, status,
		       request_url, request_method, request_headers, request_body,
		       response_status_code, response_headers, response_body, attempt_number,
		       duration_ms, next_retry_at, error_message, created_at, delivered_at
		FROM webhook_deliveries
		WHERE organization_id = $1
	`
	args := []interface{}{orgID}

	if webhookID != nil {
		args = append(args, *webhookID)
		query += fmt.Sprintf(" AND webhook_id = $%d", len(args))
	}

	if status != nil {
		args = append(args, *status)
		query += fmt.Sprintf(" AND status = $%d", len(args))
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []infrastructure.WebhookDelivery
	for rows.Next() {
		var delivery infrastructure.WebhookDelivery
		err := rows.Scan(
			&delivery.ID, &delivery.OrganizationID, &delivery.WebhookID, &delivery.EventType,
			&delivery.EventID, &delivery.Status, &delivery.RequestURL, &delivery.RequestMethod,
			&delivery.RequestHeaders, &delivery.RequestBody, &delivery.ResponseStatusCode,
			&delivery.ResponseHeaders, &delivery.ResponseBody, &delivery.AttemptNumber,
			&delivery.DurationMs, &delivery.NextRetryAt, &delivery.ErrorMessage,
			&delivery.CreatedAt, &delivery.DeliveredAt,
		)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, rows.Err()
}

// CreateWebhookDelivery creates a new webhook delivery record
func (r *PostgresDB) CreateWebhookDelivery(ctx context.Context, delivery *infrastructure.WebhookDelivery) error {
	delivery.ID = uuid.New()
	delivery.CreatedAt = time.Now()

	query := `
		INSERT INTO webhook_deliveries (
			id, organization_id, webhook_id, event_type, event_id, status,
			request_url, request_method, request_headers, request_body,
			response_status_code, response_headers, response_body, attempt_number,
			duration_ms, next_retry_at, error_message, created_at, delivered_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`

	return r.db.QueryRowContext(ctx, query,
		delivery.ID, delivery.OrganizationID, delivery.WebhookID, delivery.EventType,
		delivery.EventID, delivery.Status, delivery.RequestURL, delivery.RequestMethod,
		delivery.RequestHeaders, delivery.RequestBody, delivery.ResponseStatusCode,
		delivery.ResponseHeaders, delivery.ResponseBody, delivery.AttemptNumber,
		delivery.DurationMs, delivery.NextRetryAt, delivery.ErrorMessage,
		delivery.CreatedAt, delivery.DeliveredAt,
	).Scan(&delivery.ID)
}

// GetWebhookDelivery retrieves a webhook delivery by ID
func (r *PostgresDB) GetWebhookDelivery(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.WebhookDelivery, error) {
	query := `
		SELECT id, organization_id, webhook_id, event_type, event_id, status,
		       request_url, request_method, request_headers, request_body,
		       response_status_code, response_headers, response_body, attempt_number,
		       duration_ms, next_retry_at, error_message, created_at, delivered_at
		FROM webhook_deliveries
		WHERE id = $1 AND organization_id = $2
	`

	delivery := &infrastructure.WebhookDelivery{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&delivery.ID, &delivery.OrganizationID, &delivery.WebhookID, &delivery.EventType,
		&delivery.EventID, &delivery.Status, &delivery.RequestURL, &delivery.RequestMethod,
		&delivery.RequestHeaders, &delivery.RequestBody, &delivery.ResponseStatusCode,
		&delivery.ResponseHeaders, &delivery.ResponseBody, &delivery.AttemptNumber,
		&delivery.DurationMs, &delivery.NextRetryAt, &delivery.ErrorMessage,
		&delivery.CreatedAt, &delivery.DeliveredAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

// UpdateWebhookDelivery updates a webhook delivery
func (r *PostgresDB) UpdateWebhookDelivery(ctx context.Context, delivery *infrastructure.WebhookDelivery) error {
	query := `
		UPDATE webhook_deliveries
		SET status = $1, response_status_code = $2, response_headers = $3,
		    response_body = $4, duration_ms = $5, next_retry_at = $6,
		    error_message = $7, delivered_at = $8, attempt_number = $9
		WHERE id = $10 AND organization_id = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		delivery.Status, delivery.ResponseStatusCode, delivery.ResponseHeaders,
		delivery.ResponseBody, delivery.DurationMs, delivery.NextRetryAt,
		delivery.ErrorMessage, delivery.DeliveredAt, delivery.AttemptNumber,
		delivery.ID, delivery.OrganizationID,
	)

	return err
}

// GetPendingWebhookDeliveries retrieves pending webhook deliveries
func (r *PostgresDB) GetPendingWebhookDeliveries(ctx context.Context, limit int) ([]infrastructure.WebhookDelivery, error) {
	query := `
		SELECT id, organization_id, webhook_id, event_type, event_id, status,
		       request_url, request_method, request_headers, request_body,
		       response_status_code, response_headers, response_body, attempt_number,
		       duration_ms, next_retry_at, error_message, created_at, delivered_at
		FROM webhook_deliveries
		WHERE status IN ('pending', 'retrying')
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []infrastructure.WebhookDelivery
	for rows.Next() {
		var delivery infrastructure.WebhookDelivery
		err := rows.Scan(
			&delivery.ID, &delivery.OrganizationID, &delivery.WebhookID, &delivery.EventType,
			&delivery.EventID, &delivery.Status, &delivery.RequestURL, &delivery.RequestMethod,
			&delivery.RequestHeaders, &delivery.RequestBody, &delivery.ResponseStatusCode,
			&delivery.ResponseHeaders, &delivery.ResponseBody, &delivery.AttemptNumber,
			&delivery.DurationMs, &delivery.NextRetryAt, &delivery.ErrorMessage,
			&delivery.CreatedAt, &delivery.DeliveredAt,
		)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, rows.Err()
}

// ============================================================================
// NOTIFICATIONS REPOSITORY
// ============================================================================

// ListNotifications returns a list of notifications
func (r *PostgresDB) ListNotifications(ctx context.Context, userID uuid.UUID, isRead *bool, limit int, offset int) ([]infrastructure.Notification, error) {
	query := `
		SELECT id, organization_id, user_id, notification_type, category,
		       title, message, action_url, action_label, channels, is_read,
		       read_at, related_entity_type, related_entity_id, priority, expires_at, created_at
		FROM notifications
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if isRead != nil {
		query += " AND is_read = $2"
		args = append(args, *isRead)
	}

	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []infrastructure.Notification
	for rows.Next() {
		var notif infrastructure.Notification
		err := rows.Scan(
			&notif.ID, &notif.OrganizationID, &notif.UserID, &notif.NotificationType,
			&notif.Category, &notif.Title, &notif.Message, &notif.ActionURL,
			&notif.ActionLabel, &notif.Channels, &notif.IsRead, &notif.ReadAt,
			&notif.RelatedEntityType, &notif.RelatedEntityID, &notif.Priority,
			&notif.ExpiresAt, &notif.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}

	return notifications, rows.Err()
}

// CreateNotification creates a new notification
func (r *PostgresDB) CreateNotification(ctx context.Context, notification *infrastructure.Notification) error {
	notification.ID = uuid.New()
	notification.CreatedAt = time.Now()

	query := `
		INSERT INTO notifications (
			id, organization_id, user_id, notification_type, category,
			title, message, action_url, action_label, channels, is_read,
			read_at, related_entity_type, related_entity_id, priority, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	return r.db.QueryRowContext(ctx, query,
		notification.ID, notification.OrganizationID, notification.UserID,
		notification.NotificationType, notification.Category, notification.Title,
		notification.Message, notification.ActionURL, notification.ActionLabel,
		pq.Array(notification.Channels), notification.IsRead, notification.ReadAt,
		notification.RelatedEntityType, notification.RelatedEntityID, notification.Priority,
		notification.ExpiresAt, notification.CreatedAt,
	).Scan(&notification.ID)
}

// GetNotification retrieves a notification by ID
func (r *PostgresDB) GetNotification(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*infrastructure.Notification, error) {
	query := `
		SELECT id, organization_id, user_id, notification_type, category,
		       title, message, action_url, action_label, channels, is_read,
		       read_at, related_entity_type, related_entity_id, priority, expires_at, created_at
		FROM notifications
		WHERE id = $1 AND organization_id = $2
	`

	notif := &infrastructure.Notification{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&notif.ID, &notif.OrganizationID, &notif.UserID, &notif.NotificationType,
		&notif.Category, &notif.Title, &notif.Message, &notif.ActionURL,
		&notif.ActionLabel, &notif.Channels, &notif.IsRead, &notif.ReadAt,
		&notif.RelatedEntityType, &notif.RelatedEntityID, &notif.Priority,
		&notif.ExpiresAt, &notif.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return notif, nil
}

// UpdateNotification updates a notification
func (r *PostgresDB) UpdateNotification(ctx context.Context, notification *infrastructure.Notification) error {
	query := `
		UPDATE notifications
		SET notification_type = $1, category = $2, title = $3, message = $4,
		    action_url = $5, action_label = $6, channels = $7, is_read = $8,
		    read_at = $9, priority = $10, expires_at = $11
		WHERE id = $12 AND organization_id = $13
	`

	_, err := r.db.ExecContext(ctx, query,
		notification.NotificationType, notification.Category, notification.Title,
		notification.Message, notification.ActionURL, notification.ActionLabel,
		pq.Array(notification.Channels), notification.IsRead, notification.ReadAt,
		notification.Priority, notification.ExpiresAt, notification.ID, notification.OrganizationID,
	)

	return err
}

// MarkNotificationAsRead marks a single notification as read
func (r *PostgresDB) MarkNotificationAsRead(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, orgID)
	return err
}

// MarkAllNotificationsAsRead marks all notifications for a user as read
func (r *PostgresDB) MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE user_id = $2 AND is_read = FALSE
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

// ============================================================================
// NOTIFICATION PREFERENCES REPOSITORY
// ============================================================================

// GetNotificationPreference retrieves notification preference by user and category
func (r *PostgresDB) GetNotificationPreference(ctx context.Context, userID uuid.UUID, category string) (*infrastructure.NotificationPreference, error) {
	query := `
		SELECT id, organization_id, user_id, category, in_app_enabled, email_enabled,
		       sms_enabled, push_enabled, frequency, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1 AND category = $2
	`

	pref := &infrastructure.NotificationPreference{}
	err := r.db.QueryRowContext(ctx, query, userID, category).Scan(
		&pref.ID, &pref.OrganizationID, &pref.UserID, &pref.Category,
		&pref.InAppEnabled, &pref.EmailEnabled, &pref.SMSEnabled, &pref.PushEnabled,
		&pref.Frequency, &pref.CreatedAt, &pref.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return pref, nil
}

// GetAllNotificationPreferences retrieves all preferences for a user
func (r *PostgresDB) GetAllNotificationPreferences(ctx context.Context, userID uuid.UUID) ([]infrastructure.NotificationPreference, error) {
	query := `
		SELECT id, organization_id, user_id, category, in_app_enabled, email_enabled,
		       sms_enabled, push_enabled, frequency, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1
		ORDER BY category ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prefs []infrastructure.NotificationPreference
	for rows.Next() {
		var pref infrastructure.NotificationPreference
		err := rows.Scan(
			&pref.ID, &pref.OrganizationID, &pref.UserID, &pref.Category,
			&pref.InAppEnabled, &pref.EmailEnabled, &pref.SMSEnabled, &pref.PushEnabled,
			&pref.Frequency, &pref.CreatedAt, &pref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		prefs = append(prefs, pref)
	}

	return prefs, rows.Err()
}

// UpsertNotificationPreference creates or updates a notification preference
func (r *PostgresDB) UpsertNotificationPreference(ctx context.Context, pref *infrastructure.NotificationPreference) error {
	pref.UpdatedAt = time.Now()

	if pref.ID == uuid.Nil {
		pref.ID = uuid.New()
		pref.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO notification_preferences (
			id, organization_id, user_id, category, in_app_enabled, email_enabled,
			sms_enabled, push_enabled, frequency, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id, category) DO UPDATE SET
			in_app_enabled = $5, email_enabled = $6, sms_enabled = $7,
			push_enabled = $8, frequency = $9, updated_at = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		pref.ID, pref.OrganizationID, pref.UserID, pref.Category,
		pref.InAppEnabled, pref.EmailEnabled, pref.SMSEnabled, pref.PushEnabled,
		pref.Frequency, pref.CreatedAt, pref.UpdatedAt,
	)

	return err
}

// DeleteNotificationPreference deletes a notification preference
func (r *PostgresDB) DeleteNotificationPreference(ctx context.Context, userID uuid.UUID, category string) error {
	query := `
		DELETE FROM notification_preferences
		WHERE user_id = $1 AND category = $2
	`

	_, err := r.db.ExecContext(ctx, query, userID, category)
	return err
}

// ============================================================================
// USER SESSIONS REPOSITORY
// ============================================================================

// ListUserSessions returns a list of user sessions
func (r *PostgresDB) ListUserSessions(ctx context.Context, userID uuid.UUID, isActive *bool, limit int, offset int) ([]infrastructure.UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, session_token, refresh_token, user_agent,
		       ip_address, device_type, device_name, browser, os, country_code, city,
		       is_active, last_activity_at, expires_at, created_at, revoked_at
		FROM user_sessions
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if isActive != nil {
		query += " AND is_active = $2"
		args = append(args, *isActive)
	}

	query += " ORDER BY last_activity_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []infrastructure.UserSession
	for rows.Next() {
		var session infrastructure.UserSession
		err := rows.Scan(
			&session.ID, &session.UserID, &session.OrganizationID, &session.SessionToken,
			&session.RefreshToken, &session.UserAgent, &session.IPAddress, &session.DeviceType,
			&session.DeviceName, &session.Browser, &session.OS, &session.CountryCode,
			&session.City, &session.IsActive, &session.LastActivityAt, &session.ExpiresAt,
			&session.CreatedAt, &session.RevokedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// CreateUserSession creates a new user session
func (r *PostgresDB) CreateUserSession(ctx context.Context, session *infrastructure.UserSession) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now()
	session.LastActivityAt = time.Now()

	query := `
		INSERT INTO user_sessions (
			id, user_id, organization_id, session_token, refresh_token, user_agent,
			ip_address, device_type, device_name, browser, os, country_code, city,
			is_active, last_activity_at, expires_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	return r.db.QueryRowContext(ctx, query,
		session.ID, session.UserID, session.OrganizationID, session.SessionToken,
		session.RefreshToken, session.UserAgent, session.IPAddress, session.DeviceType,
		session.DeviceName, session.Browser, session.OS, session.CountryCode,
		session.City, session.IsActive, session.LastActivityAt, session.ExpiresAt,
		session.CreatedAt,
	).Scan(&session.ID)
}

// GetUserSession retrieves a user session by ID
func (r *PostgresDB) GetUserSession(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*infrastructure.UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, session_token, refresh_token, user_agent,
		       ip_address, device_type, device_name, browser, os, country_code, city,
		       is_active, last_activity_at, expires_at, created_at, revoked_at
		FROM user_sessions
		WHERE id = $1 AND user_id = $2
	`

	session := &infrastructure.UserSession{}
	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&session.ID, &session.UserID, &session.OrganizationID, &session.SessionToken,
		&session.RefreshToken, &session.UserAgent, &session.IPAddress, &session.DeviceType,
		&session.DeviceName, &session.Browser, &session.OS, &session.CountryCode,
		&session.City, &session.IsActive, &session.LastActivityAt, &session.ExpiresAt,
		&session.CreatedAt, &session.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetUserSessionByToken retrieves a user session by token
func (r *PostgresDB) GetUserSessionByToken(ctx context.Context, token string) (*infrastructure.UserSession, error) {
	query := `
		SELECT id, user_id, organization_id, session_token, refresh_token, user_agent,
		       ip_address, device_type, device_name, browser, os, country_code, city,
		       is_active, last_activity_at, expires_at, created_at, revoked_at
		FROM user_sessions
		WHERE session_token = $1 AND is_active = TRUE AND expires_at > NOW()
	`

	session := &infrastructure.UserSession{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.OrganizationID, &session.SessionToken,
		&session.RefreshToken, &session.UserAgent, &session.IPAddress, &session.DeviceType,
		&session.DeviceName, &session.Browser, &session.OS, &session.CountryCode,
		&session.City, &session.IsActive, &session.LastActivityAt, &session.ExpiresAt,
		&session.CreatedAt, &session.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

// UpdateUserSession updates a user session
func (r *PostgresDB) UpdateUserSession(ctx context.Context, session *infrastructure.UserSession) error {
	session.LastActivityAt = time.Now()

	query := `
		UPDATE user_sessions
		SET session_token = $1, refresh_token = $2, is_active = $3,
		    last_activity_at = $4, expires_at = $5
		WHERE id = $6 AND user_id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		session.SessionToken, session.RefreshToken, session.IsActive,
		session.LastActivityAt, session.ExpiresAt, session.ID, session.UserID,
	)

	return err
}

// RevokeUserSession revokes a single user session
func (r *PostgresDB) RevokeUserSession(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	query := `
		UPDATE user_sessions
		SET is_active = FALSE, revoked_at = $1
		WHERE id = $2 AND user_id = $3
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id, userID)
	return err
}

// RevokeAllUserSessions revokes all user sessions
func (r *PostgresDB) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE user_sessions
		SET is_active = FALSE, revoked_at = $1
		WHERE user_id = $2 AND is_active = TRUE
	`

	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

// CleanupExpiredUserSessions removes expired user sessions
func (r *PostgresDB) CleanupExpiredUserSessions(ctx context.Context) error {
	query := `
		DELETE FROM user_sessions
		WHERE expires_at < NOW()
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}
