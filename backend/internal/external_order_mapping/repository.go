package external_order_mapping

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

// Repository handles database operations for ExternalOrderMappings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ExternalOrderMappings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ExternalOrderMappings represents a external_order_mappings entity
type ExternalOrderMappings struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SaleId uuid.UUID `json:"sale_id" db:"sale_id"`
	SalesChannelId uuid.UUID `json:"sales_channel_id" db:"sales_channel_id"`
	ExternalOrderId string `json:"external_order_id" db:"external_order_id"`
	ExternalOrderNumber *string `json:"external_order_number" db:"external_order_number"`
	SyncStatus *string `json:"sync_status" db:"sync_status"`
	// 	SyncStatus *string `json:"sync_status" db:"sync_status"`
	LastSyncAt *time.Time `json:"last_sync_at" db:"last_sync_at"`
	ExternalData json.RawMessage `json:"external_data" db:"external_data"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new external_order_mappings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ExternalOrderMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "external_order_mappings", duration, nil)
	}()

	query := `
		INSERT INTO external_order_mappings (
			, organization_id
			, sale_id
			, sales_channel_id
			, external_order_id
			, external_order_number
			, sync_status
			, sync_status
			, last_sync_at
			, external_data
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
			, $13
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SaleId,
		entity.SalesChannelId,
		entity.ExternalOrderId,
		entity.ExternalOrderNumber,
		entity.SyncStatus,
		entity.SyncStatus,
		entity.LastSyncAt,
		entity.ExternalData,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create external_order_mappings", zap.Error(err))
		return fmt.Errorf("failed to create external_order_mappings: %w", err)
	}

	r.logger.Info("created external_order_mappings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a external_order_mappings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ExternalOrderMappings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "external_order_mappings", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, sales_channel_id
			, external_order_id
			, external_order_number
			, sync_status
			, sync_status
			, last_sync_at
			, external_data
			, created_at
			, updated_at
			, deleted_at
		FROM external_order_mappings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ExternalOrderMappings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SaleId,
		&entity.SalesChannelId,
		&entity.ExternalOrderId,
		&entity.ExternalOrderNumber,
		&entity.SyncStatus,
		&entity.SyncStatus,
		&entity.LastSyncAt,
		&entity.ExternalData,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("external_order_mappings not found")
	}

	if err != nil {
		r.logger.Error("failed to get external_order_mappings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get external_order_mappings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of external_order_mappings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ExternalOrderMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "external_order_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM external_order_mappings
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count external_order_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, sales_channel_id
			, external_order_id
			, external_order_number
			, sync_status
			, sync_status
			, last_sync_at
			, external_data
			, created_at
			, updated_at
			, deleted_at
		FROM external_order_mappings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list external_order_mappings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list external_order_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*ExternalOrderMappings
	for rows.Next() {
		var entity ExternalOrderMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleId,
			&entity.SalesChannelId,
			&entity.ExternalOrderId,
			&entity.ExternalOrderNumber,
			&entity.SyncStatus,
			&entity.SyncStatus,
			&entity.LastSyncAt,
			&entity.ExternalData,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan external_order_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating external_order_mappings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing external_order_mappings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ExternalOrderMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "external_order_mappings", duration, nil)
	}()

	query := `
		UPDATE external_order_mappings
		SET
			, organization_id = $2
			, sale_id = $3
			, sales_channel_id = $4
			, external_order_id = $5
			, external_order_number = $6
			, sync_status = $7
			, sync_status = $8
			, last_sync_at = $9
			, external_data = $10
			, updated_at = $12
			, deleted_at = $13
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SaleId,
		entity.SalesChannelId,
		entity.ExternalOrderId,
		entity.ExternalOrderNumber,
		entity.SyncStatus,
		entity.SyncStatus,
		entity.LastSyncAt,
		entity.ExternalData,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update external_order_mappings", zap.Error(err))
		return fmt.Errorf("failed to update external_order_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("external_order_mappings not found or already deleted")
	}

	r.logger.Info("updated external_order_mappings",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a external_order_mappings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "external_order_mappings", duration, nil)
	}()

	query := `
		UPDATE external_order_mappings
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete external_order_mappings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete external_order_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("external_order_mappings not found or already deleted")
	}

	r.logger.Info("deleted external_order_mappings", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves external_order_mappings records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ExternalOrderMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "external_order_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM external_order_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count external_order_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_id
			, sales_channel_id
			, external_order_id
			, external_order_number
			, sync_status
			, sync_status
			, last_sync_at
			, external_data
			, created_at
			, updated_at
			, deleted_at
		FROM external_order_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list external_order_mappings by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list external_order_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*ExternalOrderMappings
	for rows.Next() {
		var entity ExternalOrderMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleId,
			&entity.SalesChannelId,
			&entity.ExternalOrderId,
			&entity.ExternalOrderNumber,
			&entity.SyncStatus,
			&entity.SyncStatus,
			&entity.LastSyncAt,
			&entity.ExternalData,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan external_order_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

