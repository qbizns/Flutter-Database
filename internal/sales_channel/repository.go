package sales_channel

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

// Repository handles database operations for SalesChannels
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SalesChannels repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SalesChannels represents a sales_channels entity
type SalesChannels struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ChannelCode string `json:"channel_code" db:"channel_code"`
	ChannelName string `json:"channel_name" db:"channel_name"`
	ChannelType string `json:"channel_type" db:"channel_type"`
	ChannelType *string `json:"channel_type" db:"channel_type"`
	IsActive *bool `json:"is_active" db:"is_active"`
	SyncInventory *bool `json:"sync_inventory" db:"sync_inventory"`
	SyncCustomers *bool `json:"sync_customers" db:"sync_customers"`
	ExternalSystemName *string `json:"external_system_name" db:"external_system_name"`
	ApiEndpoint *string `json:"api_endpoint" db:"api_endpoint"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new sales_channels record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SalesChannels) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sales_channels", duration, nil)
	}()

	query := `
		INSERT INTO sales_channels (
			, organization_id
			, channel_code
			, channel_name
			, channel_type
			, channel_type
			, is_active
			, sync_inventory
			, sync_customers
			, external_system_name
			, api_endpoint
			, settings
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
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ChannelCode,
		entity.ChannelName,
		entity.ChannelType,
		entity.ChannelType,
		entity.IsActive,
		entity.SyncInventory,
		entity.SyncCustomers,
		entity.ExternalSystemName,
		entity.ApiEndpoint,
		entity.Settings,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create sales_channels", zap.Error(err))
		return fmt.Errorf("failed to create sales_channels: %w", err)
	}

	r.logger.Info("created sales_channels",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a sales_channels by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SalesChannels, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales_channels", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, channel_code
			, channel_name
			, channel_type
			, channel_type
			, is_active
			, sync_inventory
			, sync_customers
			, external_system_name
			, api_endpoint
			, settings
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sales_channels
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity SalesChannels
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ChannelCode,
		&entity.ChannelName,
		&entity.ChannelType,
		&entity.ChannelType,
		&entity.IsActive,
		&entity.SyncInventory,
		&entity.SyncCustomers,
		&entity.ExternalSystemName,
		&entity.ApiEndpoint,
		&entity.Settings,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sales_channels not found")
	}

	if err != nil {
		r.logger.Error("failed to get sales_channels", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sales_channels: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sales_channels records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SalesChannels, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales_channels", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sales_channels
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sales_channels records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, channel_code
			, channel_name
			, channel_type
			, channel_type
			, is_active
			, sync_inventory
			, sync_customers
			, external_system_name
			, api_endpoint
			, settings
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sales_channels
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sales_channels", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sales_channels: %w", err)
	}
	defer rows.Close()

	var entities []*SalesChannels
	for rows.Next() {
		var entity SalesChannels
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ChannelCode,
			&entity.ChannelName,
			&entity.ChannelType,
			&entity.ChannelType,
			&entity.IsActive,
			&entity.SyncInventory,
			&entity.SyncCustomers,
			&entity.ExternalSystemName,
			&entity.ApiEndpoint,
			&entity.Settings,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sales_channels: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sales_channels rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sales_channels record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SalesChannels) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sales_channels", duration, nil)
	}()

	query := `
		UPDATE sales_channels
		SET
			, organization_id = $2
			, channel_code = $3
			, channel_name = $4
			, channel_type = $5
			, channel_type = $6
			, is_active = $7
			, sync_inventory = $8
			, sync_customers = $9
			, external_system_name = $10
			, api_endpoint = $11
			, settings = $12
			, created_by = $13
			, updated_at = $15
			, deleted_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ChannelCode,
		entity.ChannelName,
		entity.ChannelType,
		entity.ChannelType,
		entity.IsActive,
		entity.SyncInventory,
		entity.SyncCustomers,
		entity.ExternalSystemName,
		entity.ApiEndpoint,
		entity.Settings,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sales_channels", zap.Error(err))
		return fmt.Errorf("failed to update sales_channels: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sales_channels not found or already deleted")
	}

	r.logger.Info("updated sales_channels",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a sales_channels record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sales_channels", duration, nil)
	}()

	query := `
		UPDATE sales_channels
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sales_channels", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sales_channels: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sales_channels not found or already deleted")
	}

	r.logger.Info("deleted sales_channels", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sales_channels records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SalesChannels, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales_channels", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sales_channels
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sales_channels records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, channel_code
			, channel_name
			, channel_type
			, channel_type
			, is_active
			, sync_inventory
			, sync_customers
			, external_system_name
			, api_endpoint
			, settings
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sales_channels
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sales_channels by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sales_channels: %w", err)
	}
	defer rows.Close()

	var entities []*SalesChannels
	for rows.Next() {
		var entity SalesChannels
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ChannelCode,
			&entity.ChannelName,
			&entity.ChannelType,
			&entity.ChannelType,
			&entity.IsActive,
			&entity.SyncInventory,
			&entity.SyncCustomers,
			&entity.ExternalSystemName,
			&entity.ApiEndpoint,
			&entity.Settings,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sales_channels: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

