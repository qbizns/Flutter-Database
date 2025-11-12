package api_request_log

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

// Repository handles database operations for ApiRequestLogs
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ApiRequestLogs repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ApiRequestLogs represents a api_request_logs entity
type ApiRequestLogs struct {
	Id *uuid.UUID `json:"id" db:"id"`
	RequestId *string `json:"request_id" db:"request_id"`
	Method string `json:"method" db:"method"`
	Path string `json:"path" db:"path"`
	QueryParams json.RawMessage `json:"query_params" db:"query_params"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	ApiKeyId *uuid.UUID `json:"api_key_id" db:"api_key_id"`
	RequestHeaders json.RawMessage `json:"request_headers" db:"request_headers"`
	RequestBody json.RawMessage `json:"request_body" db:"request_body"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
	StatusCode int64 `json:"status_code" db:"status_code"`
	ResponseHeaders json.RawMessage `json:"response_headers" db:"response_headers"`
	ResponseBody json.RawMessage `json:"response_body" db:"response_body"`
	DurationMs *int64 `json:"duration_ms" db:"duration_ms"`
	ErrorMessage *string `json:"error_message" db:"error_message"`
	ErrorStack *string `json:"error_stack" db:"error_stack"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new api_request_logs record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ApiRequestLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "api_request_logs", duration, nil)
	}()

	query := `
		INSERT INTO api_request_logs (
			, request_id
			, method
			, path
			, query_params
			, user_id
			, organization_id
			, api_key_id
			, request_headers
			, request_body
			, ip_address
			, user_agent
			, status_code
			, response_headers
			, response_body
			, duration_ms
			, error_message
			, error_stack
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
		entity.RequestId,
		entity.Method,
		entity.Path,
		entity.QueryParams,
		entity.UserId,
		entity.OrganizationId,
		entity.ApiKeyId,
		entity.RequestHeaders,
		entity.RequestBody,
		entity.IpAddress,
		entity.UserAgent,
		entity.StatusCode,
		entity.ResponseHeaders,
		entity.ResponseBody,
		entity.DurationMs,
		entity.ErrorMessage,
		entity.ErrorStack,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create api_request_logs", zap.Error(err))
		return fmt.Errorf("failed to create api_request_logs: %w", err)
	}

	r.logger.Info("created api_request_logs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a api_request_logs by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ApiRequestLogs, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_request_logs", duration, nil)
	}()

	query := `
		SELECT
			id
			, request_id
			, method
			, path
			, query_params
			, user_id
			, organization_id
			, api_key_id
			, request_headers
			, request_body
			, ip_address
			, user_agent
			, status_code
			, response_headers
			, response_body
			, duration_ms
			, error_message
			, error_stack
			, created_at
		FROM api_request_logs
		WHERE id = $1
		
	`

	var entity ApiRequestLogs
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.RequestId,
		&entity.Method,
		&entity.Path,
		&entity.QueryParams,
		&entity.UserId,
		&entity.OrganizationId,
		&entity.ApiKeyId,
		&entity.RequestHeaders,
		&entity.RequestBody,
		&entity.IpAddress,
		&entity.UserAgent,
		&entity.StatusCode,
		&entity.ResponseHeaders,
		&entity.ResponseBody,
		&entity.DurationMs,
		&entity.ErrorMessage,
		&entity.ErrorStack,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("api_request_logs not found")
	}

	if err != nil {
		r.logger.Error("failed to get api_request_logs", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get api_request_logs: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of api_request_logs records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ApiRequestLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_request_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM api_request_logs
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count api_request_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, request_id
			, method
			, path
			, query_params
			, user_id
			, organization_id
			, api_key_id
			, request_headers
			, request_body
			, ip_address
			, user_agent
			, status_code
			, response_headers
			, response_body
			, duration_ms
			, error_message
			, error_stack
			, created_at
		FROM api_request_logs
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list api_request_logs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list api_request_logs: %w", err)
	}
	defer rows.Close()

	var entities []*ApiRequestLogs
	for rows.Next() {
		var entity ApiRequestLogs
		err := rows.Scan(
			&entity.Id,
			&entity.RequestId,
			&entity.Method,
			&entity.Path,
			&entity.QueryParams,
			&entity.UserId,
			&entity.OrganizationId,
			&entity.ApiKeyId,
			&entity.RequestHeaders,
			&entity.RequestBody,
			&entity.IpAddress,
			&entity.UserAgent,
			&entity.StatusCode,
			&entity.ResponseHeaders,
			&entity.ResponseBody,
			&entity.DurationMs,
			&entity.ErrorMessage,
			&entity.ErrorStack,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan api_request_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating api_request_logs rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing api_request_logs record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ApiRequestLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "api_request_logs", duration, nil)
	}()

	query := `
		UPDATE api_request_logs
		SET
			, request_id = $2
			, method = $3
			, path = $4
			, query_params = $5
			, user_id = $6
			, organization_id = $7
			, api_key_id = $8
			, request_headers = $9
			, request_body = $10
			, ip_address = $11
			, user_agent = $12
			, status_code = $13
			, response_headers = $14
			, response_body = $15
			, duration_ms = $16
			, error_message = $17
			, error_stack = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.RequestId,
		entity.Method,
		entity.Path,
		entity.QueryParams,
		entity.UserId,
		entity.OrganizationId,
		entity.ApiKeyId,
		entity.RequestHeaders,
		entity.RequestBody,
		entity.IpAddress,
		entity.UserAgent,
		entity.StatusCode,
		entity.ResponseHeaders,
		entity.ResponseBody,
		entity.DurationMs,
		entity.ErrorMessage,
		entity.ErrorStack,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update api_request_logs", zap.Error(err))
		return fmt.Errorf("failed to update api_request_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("api_request_logs not found or already deleted")
	}

	r.logger.Info("updated api_request_logs",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a api_request_logs record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "api_request_logs", duration, nil)
	}()

	query := `DELETE FROM api_request_logs WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete api_request_logs", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete api_request_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("api_request_logs not found")
	}

	r.logger.Info("deleted api_request_logs", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves api_request_logs records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ApiRequestLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_request_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM api_request_logs
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count api_request_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, request_id
			, method
			, path
			, query_params
			, user_id
			, organization_id
			, api_key_id
			, request_headers
			, request_body
			, ip_address
			, user_agent
			, status_code
			, response_headers
			, response_body
			, duration_ms
			, error_message
			, error_stack
			, created_at
		FROM api_request_logs
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list api_request_logs by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list api_request_logs: %w", err)
	}
	defer rows.Close()

	var entities []*ApiRequestLogs
	for rows.Next() {
		var entity ApiRequestLogs
		err := rows.Scan(
			&entity.Id,
			&entity.RequestId,
			&entity.Method,
			&entity.Path,
			&entity.QueryParams,
			&entity.UserId,
			&entity.OrganizationId,
			&entity.ApiKeyId,
			&entity.RequestHeaders,
			&entity.RequestBody,
			&entity.IpAddress,
			&entity.UserAgent,
			&entity.StatusCode,
			&entity.ResponseHeaders,
			&entity.ResponseBody,
			&entity.DurationMs,
			&entity.ErrorMessage,
			&entity.ErrorStack,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan api_request_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

