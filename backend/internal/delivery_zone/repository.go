package delivery_zone

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

// Repository handles database operations for DeliveryZones
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DeliveryZones repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DeliveryZones represents a delivery_zones entity
type DeliveryZones struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	ZoneName string `json:"zone_name" db:"zone_name"`
	ZoneCode *string `json:"zone_code" db:"zone_code"`
	Description *string `json:"description" db:"description"`
	Geofence json.RawMessage `json:"geofence" db:"geofence"`
	PostalCodes *string `json:"postal_codes" db:"postal_codes"`
	CoverageNotes *string `json:"coverage_notes" db:"coverage_notes"`
	BaseDeliveryFee *float64 `json:"base_delivery_fee" db:"base_delivery_fee"`
	FeeType *string `json:"fee_type" db:"fee_type"`
	MinimumOrderAmount *float64 `json:"minimum_order_amount" db:"minimum_order_amount"`
	FreeDeliveryThreshold *float64 `json:"free_delivery_threshold" db:"free_delivery_threshold"`
	EstimatedDeliveryTimeMinutes *int64 `json:"estimated_delivery_time_minutes" db:"estimated_delivery_time_minutes"`
	MaxDeliveryTimeMinutes *int64 `json:"max_delivery_time_minutes" db:"max_delivery_time_minutes"`
	Priority *int64 `json:"priority" db:"priority"`
	IsActive *bool `json:"is_active" db:"is_active"`
	ActiveHours json.RawMessage `json:"active_hours" db:"active_hours"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	BaseDeliveryFee *string `json:"base_delivery_fee" db:"base_delivery_fee"`
}

// Create inserts a new delivery_zones record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DeliveryZones) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "delivery_zones", duration, nil)
	}()

	query := `
		INSERT INTO delivery_zones (
			, organization_id
			, location_id
			, zone_name
			, zone_code
			, description
			, geofence
			, postal_codes
			, coverage_notes
			, base_delivery_fee
			, fee_type
			, minimum_order_amount
			, free_delivery_threshold
			, estimated_delivery_time_minutes
			, max_delivery_time_minutes
			, priority
			, is_active
			, active_hours
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, base_delivery_fee
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
			, $22
			, $23
			, $24
			, $25
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.ZoneName,
		entity.ZoneCode,
		entity.Description,
		entity.Geofence,
		entity.PostalCodes,
		entity.CoverageNotes,
		entity.BaseDeliveryFee,
		entity.FeeType,
		entity.MinimumOrderAmount,
		entity.FreeDeliveryThreshold,
		entity.EstimatedDeliveryTimeMinutes,
		entity.MaxDeliveryTimeMinutes,
		entity.Priority,
		entity.IsActive,
		entity.ActiveHours,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.BaseDeliveryFee,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create delivery_zones", zap.Error(err))
		return fmt.Errorf("failed to create delivery_zones: %w", err)
	}

	r.logger.Info("created delivery_zones",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a delivery_zones by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DeliveryZones, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_zones", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, zone_name
			, zone_code
			, description
			, geofence
			, postal_codes
			, coverage_notes
			, base_delivery_fee
			, fee_type
			, minimum_order_amount
			, free_delivery_threshold
			, estimated_delivery_time_minutes
			, max_delivery_time_minutes
			, priority
			, is_active
			, active_hours
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, base_delivery_fee
		FROM delivery_zones
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DeliveryZones
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.ZoneName,
		&entity.ZoneCode,
		&entity.Description,
		&entity.Geofence,
		&entity.PostalCodes,
		&entity.CoverageNotes,
		&entity.BaseDeliveryFee,
		&entity.FeeType,
		&entity.MinimumOrderAmount,
		&entity.FreeDeliveryThreshold,
		&entity.EstimatedDeliveryTimeMinutes,
		&entity.MaxDeliveryTimeMinutes,
		&entity.Priority,
		&entity.IsActive,
		&entity.ActiveHours,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.BaseDeliveryFee,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("delivery_zones not found")
	}

	if err != nil {
		r.logger.Error("failed to get delivery_zones", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get delivery_zones: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of delivery_zones records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DeliveryZones, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_zones", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_zones
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_zones records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, zone_name
			, zone_code
			, description
			, geofence
			, postal_codes
			, coverage_notes
			, base_delivery_fee
			, fee_type
			, minimum_order_amount
			, free_delivery_threshold
			, estimated_delivery_time_minutes
			, max_delivery_time_minutes
			, priority
			, is_active
			, active_hours
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, base_delivery_fee
		FROM delivery_zones
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_zones", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_zones: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryZones
	for rows.Next() {
		var entity DeliveryZones
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.ZoneName,
			&entity.ZoneCode,
			&entity.Description,
			&entity.Geofence,
			&entity.PostalCodes,
			&entity.CoverageNotes,
			&entity.BaseDeliveryFee,
			&entity.FeeType,
			&entity.MinimumOrderAmount,
			&entity.FreeDeliveryThreshold,
			&entity.EstimatedDeliveryTimeMinutes,
			&entity.MaxDeliveryTimeMinutes,
			&entity.Priority,
			&entity.IsActive,
			&entity.ActiveHours,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.BaseDeliveryFee,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_zones: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating delivery_zones rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing delivery_zones record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DeliveryZones) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_zones", duration, nil)
	}()

	query := `
		UPDATE delivery_zones
		SET
			, organization_id = $2
			, location_id = $3
			, zone_name = $4
			, zone_code = $5
			, description = $6
			, geofence = $7
			, postal_codes = $8
			, coverage_notes = $9
			, base_delivery_fee = $10
			, fee_type = $11
			, minimum_order_amount = $12
			, free_delivery_threshold = $13
			, estimated_delivery_time_minutes = $14
			, max_delivery_time_minutes = $15
			, priority = $16
			, is_active = $17
			, active_hours = $18
			, metadata = $19
			, updated_at = $21
			, created_by = $22
			, updated_by = $23
			, deleted_at = $24
			, base_delivery_fee = $25
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $26
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.ZoneName,
		entity.ZoneCode,
		entity.Description,
		entity.Geofence,
		entity.PostalCodes,
		entity.CoverageNotes,
		entity.BaseDeliveryFee,
		entity.FeeType,
		entity.MinimumOrderAmount,
		entity.FreeDeliveryThreshold,
		entity.EstimatedDeliveryTimeMinutes,
		entity.MaxDeliveryTimeMinutes,
		entity.Priority,
		entity.IsActive,
		entity.ActiveHours,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.BaseDeliveryFee,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update delivery_zones", zap.Error(err))
		return fmt.Errorf("failed to update delivery_zones: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_zones not found or already deleted")
	}

	r.logger.Info("updated delivery_zones",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a delivery_zones record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "delivery_zones", duration, nil)
	}()

	query := `
		UPDATE delivery_zones
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete delivery_zones", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete delivery_zones: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery_zones not found or already deleted")
	}

	r.logger.Info("deleted delivery_zones", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves delivery_zones records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DeliveryZones, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "delivery_zones", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM delivery_zones
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count delivery_zones records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, zone_name
			, zone_code
			, description
			, geofence
			, postal_codes
			, coverage_notes
			, base_delivery_fee
			, fee_type
			, minimum_order_amount
			, free_delivery_threshold
			, estimated_delivery_time_minutes
			, max_delivery_time_minutes
			, priority
			, is_active
			, active_hours
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, base_delivery_fee
		FROM delivery_zones
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list delivery_zones by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list delivery_zones: %w", err)
	}
	defer rows.Close()

	var entities []*DeliveryZones
	for rows.Next() {
		var entity DeliveryZones
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.ZoneName,
			&entity.ZoneCode,
			&entity.Description,
			&entity.Geofence,
			&entity.PostalCodes,
			&entity.CoverageNotes,
			&entity.BaseDeliveryFee,
			&entity.FeeType,
			&entity.MinimumOrderAmount,
			&entity.FreeDeliveryThreshold,
			&entity.EstimatedDeliveryTimeMinutes,
			&entity.MaxDeliveryTimeMinutes,
			&entity.Priority,
			&entity.IsActive,
			&entity.ActiveHours,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.BaseDeliveryFee,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan delivery_zones: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

