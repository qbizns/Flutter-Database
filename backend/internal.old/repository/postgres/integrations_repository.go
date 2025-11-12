package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/integrations"
)

type SalesChannelRepository struct {
	db *DB
}

func NewSalesChannelRepository(db *DB) *SalesChannelRepository {
	return &SalesChannelRepository{db: db}
}

func (r *SalesChannelRepository) List(ctx context.Context, orgID uuid.UUID, filters integrations.SalesChannelFilters) ([]integrations.SalesChannel, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, channel_code, channel_name, channel_type,
		       is_active, sync_inventory, sync_customers, external_system_name,
		       api_endpoint, settings, created_by, created_at, updated_at
		FROM sales_channels
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (channel_name ILIKE $%d OR channel_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.ChannelType != nil {
		argCount++
		query += fmt.Sprintf(" AND channel_type = $%d", argCount)
		args = append(args, *filters.ChannelType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY channel_name ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []integrations.SalesChannel
	for rows.Next() {
		var c integrations.SalesChannel
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.ChannelCode, &c.ChannelName, &c.ChannelType,
			&c.IsActive, &c.SyncInventory, &c.SyncCustomers, &c.ExternalSystemName,
			&c.APIEndpoint, &c.Settings, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}

	return channels, rows.Err()
}

func (r *SalesChannelRepository) Count(ctx context.Context, orgID uuid.UUID, filters integrations.SalesChannelFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM sales_channels WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (channel_name ILIKE $%d OR channel_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.ChannelType != nil {
		argCount++
		query += fmt.Sprintf(" AND channel_type = $%d", argCount)
		args = append(args, *filters.ChannelType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SalesChannelRepository) Create(ctx context.Context, channel *integrations.SalesChannel) error {
	if err := r.db.SetOrganizationContext(ctx, channel.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO sales_channels (
			id, organization_id, channel_code, channel_name, channel_type,
			is_active, sync_inventory, sync_customers, external_system_name,
			api_endpoint, settings, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		channel.ID, channel.OrganizationID, channel.ChannelCode, channel.ChannelName, channel.ChannelType,
		channel.IsActive, channel.SyncInventory, channel.SyncCustomers, channel.ExternalSystemName,
		channel.APIEndpoint, channel.Settings, channel.CreatedBy, channel.CreatedAt, channel.UpdatedAt,
	)
	return err
}

func (r *SalesChannelRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*integrations.SalesChannel, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, channel_code, channel_name, channel_type,
		       is_active, sync_inventory, sync_customers, external_system_name,
		       api_endpoint, settings, created_by, created_at, updated_at
		FROM sales_channels
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var c integrations.SalesChannel
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&c.ID, &c.OrganizationID, &c.ChannelCode, &c.ChannelName, &c.ChannelType,
		&c.IsActive, &c.SyncInventory, &c.SyncCustomers, &c.ExternalSystemName,
		&c.APIEndpoint, &c.Settings, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SalesChannelRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*integrations.SalesChannel, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, channel_code, channel_name, channel_type,
		       is_active, sync_inventory, sync_customers, external_system_name,
		       api_endpoint, settings, created_by, created_at, updated_at
		FROM sales_channels
		WHERE organization_id = $1 AND channel_code = $2 AND deleted_at IS NULL
	`

	var c integrations.SalesChannel
	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&c.ID, &c.OrganizationID, &c.ChannelCode, &c.ChannelName, &c.ChannelType,
		&c.IsActive, &c.SyncInventory, &c.SyncCustomers, &c.ExternalSystemName,
		&c.APIEndpoint, &c.Settings, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SalesChannelRepository) Update(ctx context.Context, channel *integrations.SalesChannel) error {
	if err := r.db.SetOrganizationContext(ctx, channel.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE sales_channels SET
			channel_code = $3, channel_name = $4, channel_type = $5,
			is_active = $6, sync_inventory = $7, sync_customers = $8,
			external_system_name = $9, api_endpoint = $10, settings = $11,
			updated_at = $12
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		channel.OrganizationID, channel.ID,
		channel.ChannelCode, channel.ChannelName, channel.ChannelType,
		channel.IsActive, channel.SyncInventory, channel.SyncCustomers,
		channel.ExternalSystemName, channel.APIEndpoint, channel.Settings,
		channel.UpdatedAt,
	)
	return err
}

func (r *SalesChannelRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE sales_channels
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ExternalOrderMappingRepository implementation

type ExternalOrderMappingRepository struct {
	db *DB
}

func NewExternalOrderMappingRepository(db *DB) *ExternalOrderMappingRepository {
	return &ExternalOrderMappingRepository{db: db}
}

func (r *ExternalOrderMappingRepository) List(ctx context.Context, orgID uuid.UUID, filters integrations.ExternalOrderMappingFilters) ([]integrations.ExternalOrderMapping, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, sales_channel_id, external_order_id,
		       external_order_number, sync_status, last_sync_at, external_data,
		       created_at, updated_at
		FROM external_order_mappings
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.SalesChannelID != nil {
		argCount++
		query += fmt.Sprintf(" AND sales_channel_id = $%d", argCount)
		args = append(args, *filters.SalesChannelID)
	}

	if filters.SyncStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND sync_status = $%d", argCount)
		args = append(args, *filters.SyncStatus)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []integrations.ExternalOrderMapping
	for rows.Next() {
		var m integrations.ExternalOrderMapping
		err := rows.Scan(
			&m.ID, &m.OrganizationID, &m.SaleID, &m.SalesChannelID, &m.ExternalOrderID,
			&m.ExternalOrderNumber, &m.SyncStatus, &m.LastSyncAt, &m.ExternalData,
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, m)
	}

	return mappings, rows.Err()
}

func (r *ExternalOrderMappingRepository) Count(ctx context.Context, orgID uuid.UUID, filters integrations.ExternalOrderMappingFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM external_order_mappings WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.SalesChannelID != nil {
		argCount++
		query += fmt.Sprintf(" AND sales_channel_id = $%d", argCount)
		args = append(args, *filters.SalesChannelID)
	}

	if filters.SyncStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND sync_status = $%d", argCount)
		args = append(args, *filters.SyncStatus)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *ExternalOrderMappingRepository) Create(ctx context.Context, mapping *integrations.ExternalOrderMapping) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO external_order_mappings (
			id, organization_id, sale_id, sales_channel_id, external_order_id,
			external_order_number, sync_status, last_sync_at, external_data,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.ID, mapping.OrganizationID, mapping.SaleID, mapping.SalesChannelID, mapping.ExternalOrderID,
		mapping.ExternalOrderNumber, mapping.SyncStatus, mapping.LastSyncAt, mapping.ExternalData,
		mapping.CreatedAt, mapping.UpdatedAt,
	)
	return err
}

func (r *ExternalOrderMappingRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*integrations.ExternalOrderMapping, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, sales_channel_id, external_order_id,
		       external_order_number, sync_status, last_sync_at, external_data,
		       created_at, updated_at
		FROM external_order_mappings
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var m integrations.ExternalOrderMapping
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&m.ID, &m.OrganizationID, &m.SaleID, &m.SalesChannelID, &m.ExternalOrderID,
		&m.ExternalOrderNumber, &m.SyncStatus, &m.LastSyncAt, &m.ExternalData,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ExternalOrderMappingRepository) GetBySaleID(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*integrations.ExternalOrderMapping, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, sales_channel_id, external_order_id,
		       external_order_number, sync_status, last_sync_at, external_data,
		       created_at, updated_at
		FROM external_order_mappings
		WHERE organization_id = $1 AND sale_id = $2 AND deleted_at IS NULL
	`

	var m integrations.ExternalOrderMapping
	err := r.db.Pool.QueryRow(ctx, query, orgID, saleID).Scan(
		&m.ID, &m.OrganizationID, &m.SaleID, &m.SalesChannelID, &m.ExternalOrderID,
		&m.ExternalOrderNumber, &m.SyncStatus, &m.LastSyncAt, &m.ExternalData,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ExternalOrderMappingRepository) GetByExternalOrderID(ctx context.Context, orgID uuid.UUID, channelID uuid.UUID, externalOrderID string) (*integrations.ExternalOrderMapping, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, sales_channel_id, external_order_id,
		       external_order_number, sync_status, last_sync_at, external_data,
		       created_at, updated_at
		FROM external_order_mappings
		WHERE organization_id = $1 AND sales_channel_id = $2 AND external_order_id = $3 AND deleted_at IS NULL
	`

	var m integrations.ExternalOrderMapping
	err := r.db.Pool.QueryRow(ctx, query, orgID, channelID, externalOrderID).Scan(
		&m.ID, &m.OrganizationID, &m.SaleID, &m.SalesChannelID, &m.ExternalOrderID,
		&m.ExternalOrderNumber, &m.SyncStatus, &m.LastSyncAt, &m.ExternalData,
		&m.CreatedAt, &m.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ExternalOrderMappingRepository) Update(ctx context.Context, mapping *integrations.ExternalOrderMapping) error {
	if err := r.db.SetOrganizationContext(ctx, mapping.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE external_order_mappings SET
			sale_id = $3, sales_channel_id = $4, external_order_id = $5,
			external_order_number = $6, sync_status = $7, last_sync_at = $8,
			external_data = $9, updated_at = $10
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		mapping.OrganizationID, mapping.ID,
		mapping.SaleID, mapping.SalesChannelID, mapping.ExternalOrderID,
		mapping.ExternalOrderNumber, mapping.SyncStatus, mapping.LastSyncAt,
		mapping.ExternalData, mapping.UpdatedAt,
	)
	return err
}

func (r *ExternalOrderMappingRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE external_order_mappings
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *ExternalOrderMappingRepository) UpdateSyncStatus(ctx context.Context, id uuid.UUID, status string, lastSyncAt time.Time) error {
	query := `
		UPDATE external_order_mappings
		SET sync_status = $2, last_sync_at = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, status, lastSyncAt, time.Now())
	return err
}
