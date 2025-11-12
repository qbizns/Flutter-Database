package rate_limit

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

// Repository handles database operations for RateLimits
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new RateLimits repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// RateLimits represents a rate_limits entity
type RateLimits struct {
	Id *uuid.UUID `json:"id" db:"id"`
	IdentifierType string `json:"identifier_type" db:"identifier_type"`
	IdentifierValue string `json:"identifier_value" db:"identifier_value"`
	EndpointPath *string `json:"endpoint_path" db:"endpoint_path"`
	HttpMethod *string `json:"http_method" db:"http_method"`
	WindowStart time.Time `json:"window_start" db:"window_start"`
	WindowDurationSeconds int64 `json:"window_duration_seconds" db:"window_duration_seconds"`
	RequestCount *int64 `json:"request_count" db:"request_count"`
	AllowedCount int64 `json:"allowed_count" db:"allowed_count"`
	IsBlocked *bool `json:"is_blocked" db:"is_blocked"`
	BlockedUntil *time.Time `json:"blocked_until" db:"blocked_until"`
	FirstRequestAt *time.Time `json:"first_request_at" db:"first_request_at"`
	LastRequestAt *time.Time `json:"last_request_at" db:"last_request_at"`
	IdentifierType, *string `json:"identifier_type," db:"identifier_type,"`
}

// Create inserts a new rate_limits record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *RateLimits) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "rate_limits", duration, nil)
	}()

	query := `
		INSERT INTO rate_limits (
			, identifier_type
			, identifier_value
			, endpoint_path
			, http_method
			, window_start
			, window_duration_seconds
			, request_count
			, allowed_count
			, is_blocked
			, blocked_until
			, first_request_at
			, last_request_at
			, identifier_type,
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
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.IdentifierType,
		entity.IdentifierValue,
		entity.EndpointPath,
		entity.HttpMethod,
		entity.WindowStart,
		entity.WindowDurationSeconds,
		entity.RequestCount,
		entity.AllowedCount,
		entity.IsBlocked,
		entity.BlockedUntil,
		entity.FirstRequestAt,
		entity.LastRequestAt,
		entity.IdentifierType,,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create rate_limits", zap.Error(err))
		return fmt.Errorf("failed to create rate_limits: %w", err)
	}

	r.logger.Info("created rate_limits",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a rate_limits by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*RateLimits, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "rate_limits", duration, nil)
	}()

	query := `
		SELECT
			id
			, identifier_type
			, identifier_value
			, endpoint_path
			, http_method
			, window_start
			, window_duration_seconds
			, request_count
			, allowed_count
			, is_blocked
			, blocked_until
			, first_request_at
			, last_request_at
			, identifier_type,
		FROM rate_limits
		WHERE id = $1
		
	`

	var entity RateLimits
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.IdentifierType,
		&entity.IdentifierValue,
		&entity.EndpointPath,
		&entity.HttpMethod,
		&entity.WindowStart,
		&entity.WindowDurationSeconds,
		&entity.RequestCount,
		&entity.AllowedCount,
		&entity.IsBlocked,
		&entity.BlockedUntil,
		&entity.FirstRequestAt,
		&entity.LastRequestAt,
		&entity.IdentifierType,,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("rate_limits not found")
	}

	if err != nil {
		r.logger.Error("failed to get rate_limits", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get rate_limits: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of rate_limits records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*RateLimits, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "rate_limits", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM rate_limits
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rate_limits records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, identifier_type
			, identifier_value
			, endpoint_path
			, http_method
			, window_start
			, window_duration_seconds
			, request_count
			, allowed_count
			, is_blocked
			, blocked_until
			, first_request_at
			, last_request_at
			, identifier_type,
		FROM rate_limits
		
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list rate_limits", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list rate_limits: %w", err)
	}
	defer rows.Close()

	var entities []*RateLimits
	for rows.Next() {
		var entity RateLimits
		err := rows.Scan(
			&entity.Id,
			&entity.IdentifierType,
			&entity.IdentifierValue,
			&entity.EndpointPath,
			&entity.HttpMethod,
			&entity.WindowStart,
			&entity.WindowDurationSeconds,
			&entity.RequestCount,
			&entity.AllowedCount,
			&entity.IsBlocked,
			&entity.BlockedUntil,
			&entity.FirstRequestAt,
			&entity.LastRequestAt,
			&entity.IdentifierType,,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rate_limits: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating rate_limits rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing rate_limits record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *RateLimits) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "rate_limits", duration, nil)
	}()

	query := `
		UPDATE rate_limits
		SET
			, identifier_type = $2
			, identifier_value = $3
			, endpoint_path = $4
			, http_method = $5
			, window_start = $6
			, window_duration_seconds = $7
			, request_count = $8
			, allowed_count = $9
			, is_blocked = $10
			, blocked_until = $11
			, first_request_at = $12
			, last_request_at = $13
			, identifier_type, = $14
			
		WHERE id = $15
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.IdentifierType,
		entity.IdentifierValue,
		entity.EndpointPath,
		entity.HttpMethod,
		entity.WindowStart,
		entity.WindowDurationSeconds,
		entity.RequestCount,
		entity.AllowedCount,
		entity.IsBlocked,
		entity.BlockedUntil,
		entity.FirstRequestAt,
		entity.LastRequestAt,
		entity.IdentifierType,,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update rate_limits", zap.Error(err))
		return fmt.Errorf("failed to update rate_limits: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("rate_limits not found or already deleted")
	}

	r.logger.Info("updated rate_limits",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a rate_limits record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "rate_limits", duration, nil)
	}()

	query := `DELETE FROM rate_limits WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete rate_limits", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete rate_limits: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("rate_limits not found")
	}

	r.logger.Info("deleted rate_limits", zap.String("id", id.String()))
	return nil
}



