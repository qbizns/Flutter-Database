package integration_config

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

// Repository handles database operations for IntegrationConfigs
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new IntegrationConfigs repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// IntegrationConfigs represents a integration_configs entity
type IntegrationConfigs struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	IntegrationType string `json:"integration_type" db:"integration_type"`
	ProviderName string `json:"provider_name" db:"provider_name"`
	Credentials json.RawMessage `json:"credentials" db:"credentials"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsConnected *bool `json:"is_connected" db:"is_connected"`
	ConnectionStatus *string `json:"connection_status" db:"connection_status"`
	LastSyncAt *time.Time `json:"last_sync_at" db:"last_sync_at"`
	LastSyncStatus *string `json:"last_sync_status" db:"last_sync_status"`
	SyncFrequency *string `json:"sync_frequency" db:"sync_frequency"`
	WebhookUrl *string `json:"webhook_url" db:"webhook_url"`
	WebhookSecret *string `json:"webhook_secret" db:"webhook_secret"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new integration_configs record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *IntegrationConfigs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "integration_configs", duration, nil)
	}()

	query := `
		INSERT INTO integration_configs (
			, organization_id
			, integration_type
			, provider_name
			, credentials
			, settings
			, is_active
			, is_connected
			, connection_status
			, last_sync_at
			, last_sync_status
			, sync_frequency
			, webhook_url
			, webhook_secret
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
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.IntegrationType,
		entity.ProviderName,
		entity.Credentials,
		entity.Settings,
		entity.IsActive,
		entity.IsConnected,
		entity.ConnectionStatus,
		entity.LastSyncAt,
		entity.LastSyncStatus,
		entity.SyncFrequency,
		entity.WebhookUrl,
		entity.WebhookSecret,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create integration_configs", zap.Error(err))
		return fmt.Errorf("failed to create integration_configs: %w", err)
	}

	r.logger.Info("created integration_configs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a integration_configs by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*IntegrationConfigs, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "integration_configs", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, integration_type
			, provider_name
			, credentials
			, settings
			, is_active
			, is_connected
			, connection_status
			, last_sync_at
			, last_sync_status
			, sync_frequency
			, webhook_url
			, webhook_secret
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM integration_configs
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity IntegrationConfigs
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.IntegrationType,
		&entity.ProviderName,
		&entity.Credentials,
		&entity.Settings,
		&entity.IsActive,
		&entity.IsConnected,
		&entity.ConnectionStatus,
		&entity.LastSyncAt,
		&entity.LastSyncStatus,
		&entity.SyncFrequency,
		&entity.WebhookUrl,
		&entity.WebhookSecret,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("integration_configs not found")
	}

	if err != nil {
		r.logger.Error("failed to get integration_configs", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get integration_configs: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of integration_configs records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*IntegrationConfigs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "integration_configs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM integration_configs
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count integration_configs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, integration_type
			, provider_name
			, credentials
			, settings
			, is_active
			, is_connected
			, connection_status
			, last_sync_at
			, last_sync_status
			, sync_frequency
			, webhook_url
			, webhook_secret
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM integration_configs
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list integration_configs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list integration_configs: %w", err)
	}
	defer rows.Close()

	var entities []*IntegrationConfigs
	for rows.Next() {
		var entity IntegrationConfigs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.IntegrationType,
			&entity.ProviderName,
			&entity.Credentials,
			&entity.Settings,
			&entity.IsActive,
			&entity.IsConnected,
			&entity.ConnectionStatus,
			&entity.LastSyncAt,
			&entity.LastSyncStatus,
			&entity.SyncFrequency,
			&entity.WebhookUrl,
			&entity.WebhookSecret,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan integration_configs: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating integration_configs rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing integration_configs record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *IntegrationConfigs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "integration_configs", duration, nil)
	}()

	query := `
		UPDATE integration_configs
		SET
			, organization_id = $2
			, integration_type = $3
			, provider_name = $4
			, credentials = $5
			, settings = $6
			, is_active = $7
			, is_connected = $8
			, connection_status = $9
			, last_sync_at = $10
			, last_sync_status = $11
			, sync_frequency = $12
			, webhook_url = $13
			, webhook_secret = $14
			, created_by = $15
			, updated_at = $17
			, deleted_at = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.IntegrationType,
		entity.ProviderName,
		entity.Credentials,
		entity.Settings,
		entity.IsActive,
		entity.IsConnected,
		entity.ConnectionStatus,
		entity.LastSyncAt,
		entity.LastSyncStatus,
		entity.SyncFrequency,
		entity.WebhookUrl,
		entity.WebhookSecret,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update integration_configs", zap.Error(err))
		return fmt.Errorf("failed to update integration_configs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("integration_configs not found or already deleted")
	}

	r.logger.Info("updated integration_configs",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a integration_configs record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "integration_configs", duration, nil)
	}()

	query := `
		UPDATE integration_configs
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete integration_configs", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete integration_configs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("integration_configs not found or already deleted")
	}

	r.logger.Info("deleted integration_configs", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves integration_configs records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*IntegrationConfigs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "integration_configs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM integration_configs
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count integration_configs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, integration_type
			, provider_name
			, credentials
			, settings
			, is_active
			, is_connected
			, connection_status
			, last_sync_at
			, last_sync_status
			, sync_frequency
			, webhook_url
			, webhook_secret
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM integration_configs
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list integration_configs by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list integration_configs: %w", err)
	}
	defer rows.Close()

	var entities []*IntegrationConfigs
	for rows.Next() {
		var entity IntegrationConfigs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.IntegrationType,
			&entity.ProviderName,
			&entity.Credentials,
			&entity.Settings,
			&entity.IsActive,
			&entity.IsConnected,
			&entity.ConnectionStatus,
			&entity.LastSyncAt,
			&entity.LastSyncStatus,
			&entity.SyncFrequency,
			&entity.WebhookUrl,
			&entity.WebhookSecret,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan integration_configs: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

