package api_key

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

// Repository handles database operations for ApiKeys
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ApiKeys repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ApiKeys represents a api_keys entity
type ApiKeys struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	KeyName string `json:"key_name" db:"key_name"`
	KeyPrefix string `json:"key_prefix" db:"key_prefix"`
	KeyHash string `json:"key_hash" db:"key_hash"`
	Scopes json.RawMessage `json:"scopes" db:"scopes"`
	AllowedIps *string `json:"allowed_ips" db:"allowed_ips"`
	IsActive *bool `json:"is_active" db:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at" db:"last_used_at"`
	UsageCount *int64 `json:"usage_count" db:"usage_count"`
	RateLimitPerMinute *int64 `json:"rate_limit_per_minute" db:"rate_limit_per_minute"`
	RateLimitPerHour *int64 `json:"rate_limit_per_hour" db:"rate_limit_per_hour"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new api_keys record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ApiKeys) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "api_keys", duration, nil)
	}()

	query := `
		INSERT INTO api_keys (
			, organization_id
			, key_name
			, key_prefix
			, key_hash
			, scopes
			, allowed_ips
			, is_active
			, last_used_at
			, usage_count
			, rate_limit_per_minute
			, rate_limit_per_hour
			, expires_at
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
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.KeyName,
		entity.KeyPrefix,
		entity.KeyHash,
		entity.Scopes,
		entity.AllowedIps,
		entity.IsActive,
		entity.LastUsedAt,
		entity.UsageCount,
		entity.RateLimitPerMinute,
		entity.RateLimitPerHour,
		entity.ExpiresAt,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create api_keys", zap.Error(err))
		return fmt.Errorf("failed to create api_keys: %w", err)
	}

	r.logger.Info("created api_keys",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a api_keys by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ApiKeys, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_keys", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, key_name
			, key_prefix
			, key_hash
			, scopes
			, allowed_ips
			, is_active
			, last_used_at
			, usage_count
			, rate_limit_per_minute
			, rate_limit_per_hour
			, expires_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM api_keys
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ApiKeys
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.KeyName,
		&entity.KeyPrefix,
		&entity.KeyHash,
		&entity.Scopes,
		&entity.AllowedIps,
		&entity.IsActive,
		&entity.LastUsedAt,
		&entity.UsageCount,
		&entity.RateLimitPerMinute,
		&entity.RateLimitPerHour,
		&entity.ExpiresAt,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("api_keys not found")
	}

	if err != nil {
		r.logger.Error("failed to get api_keys", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get api_keys: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of api_keys records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ApiKeys, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_keys", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM api_keys
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count api_keys records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, key_name
			, key_prefix
			, key_hash
			, scopes
			, allowed_ips
			, is_active
			, last_used_at
			, usage_count
			, rate_limit_per_minute
			, rate_limit_per_hour
			, expires_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM api_keys
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list api_keys", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list api_keys: %w", err)
	}
	defer rows.Close()

	var entities []*ApiKeys
	for rows.Next() {
		var entity ApiKeys
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.KeyName,
			&entity.KeyPrefix,
			&entity.KeyHash,
			&entity.Scopes,
			&entity.AllowedIps,
			&entity.IsActive,
			&entity.LastUsedAt,
			&entity.UsageCount,
			&entity.RateLimitPerMinute,
			&entity.RateLimitPerHour,
			&entity.ExpiresAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan api_keys: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating api_keys rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing api_keys record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ApiKeys) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "api_keys", duration, nil)
	}()

	query := `
		UPDATE api_keys
		SET
			, organization_id = $2
			, key_name = $3
			, key_prefix = $4
			, key_hash = $5
			, scopes = $6
			, allowed_ips = $7
			, is_active = $8
			, last_used_at = $9
			, usage_count = $10
			, rate_limit_per_minute = $11
			, rate_limit_per_hour = $12
			, expires_at = $13
			, created_by = $14
			, updated_at = $16
			, deleted_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.KeyName,
		entity.KeyPrefix,
		entity.KeyHash,
		entity.Scopes,
		entity.AllowedIps,
		entity.IsActive,
		entity.LastUsedAt,
		entity.UsageCount,
		entity.RateLimitPerMinute,
		entity.RateLimitPerHour,
		entity.ExpiresAt,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update api_keys", zap.Error(err))
		return fmt.Errorf("failed to update api_keys: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("api_keys not found or already deleted")
	}

	r.logger.Info("updated api_keys",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a api_keys record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "api_keys", duration, nil)
	}()

	query := `
		UPDATE api_keys
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete api_keys", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete api_keys: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("api_keys not found or already deleted")
	}

	r.logger.Info("deleted api_keys", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves api_keys records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ApiKeys, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "api_keys", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM api_keys
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count api_keys records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, key_name
			, key_prefix
			, key_hash
			, scopes
			, allowed_ips
			, is_active
			, last_used_at
			, usage_count
			, rate_limit_per_minute
			, rate_limit_per_hour
			, expires_at
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM api_keys
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list api_keys by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list api_keys: %w", err)
	}
	defer rows.Close()

	var entities []*ApiKeys
	for rows.Next() {
		var entity ApiKeys
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.KeyName,
			&entity.KeyPrefix,
			&entity.KeyHash,
			&entity.Scopes,
			&entity.AllowedIps,
			&entity.IsActive,
			&entity.LastUsedAt,
			&entity.UsageCount,
			&entity.RateLimitPerMinute,
			&entity.RateLimitPerHour,
			&entity.ExpiresAt,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan api_keys: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

