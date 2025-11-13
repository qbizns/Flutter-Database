package kitchen_station

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for KitchenStations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new KitchenStations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// KitchenStations represents a kitchen_stations entity
type KitchenStations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	StationName string `json:"station_name" db:"station_name"`
	StationCode string `json:"station_code" db:"station_code"`
	StationType *string `json:"station_type" db:"station_type"`
	Description *string `json:"description" db:"description"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	ColorCode *string `json:"color_code" db:"color_code"`
	PrinterId *uuid.UUID `json:"printer_id" db:"printer_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	AutoPrintTickets *bool `json:"auto_print_tickets" db:"auto_print_tickets"`
	AlertSoundEnabled *bool `json:"alert_sound_enabled" db:"alert_sound_enabled"`
	DisplayConfig json.RawMessage `json:"display_config" db:"display_config"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new kitchen_stations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *KitchenStations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "kitchen_stations", duration, nil)
	}()

	query := `
		INSERT INTO kitchen_stations (
			, organization_id
			, location_id
			, station_name
			, station_code
			, station_type
			, description
			, display_order
			, color_code
			, printer_id
			, is_active
			, auto_print_tickets
			, alert_sound_enabled
			, display_config
			, created_by
			, updated_by
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
			, $18
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.StationName,
		entity.StationCode,
		entity.StationType,
		entity.Description,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.PrinterId,
		entity.IsActive,
		entity.AutoPrintTickets,
		entity.AlertSoundEnabled,
		entity.DisplayConfig,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create kitchen_stations", zap.Error(err))
		return fmt.Errorf("failed to create kitchen_stations: %w", err)
	}

	r.logger.Info("created kitchen_stations",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a kitchen_stations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*KitchenStations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_stations", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, station_name
			, station_code
			, station_type
			, description
			, display_order
			, color_code
			, printer_id
			, is_active
			, auto_print_tickets
			, alert_sound_enabled
			, display_config
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM kitchen_stations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity KitchenStations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.StationName,
		&entity.StationCode,
		&entity.StationType,
		&entity.Description,
		&entity.DisplayOrder,
		&entity.ColorCode,
		&entity.PrinterId,
		&entity.IsActive,
		&entity.AutoPrintTickets,
		&entity.AlertSoundEnabled,
		&entity.DisplayConfig,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("kitchen_stations not found")
	}

	if err != nil {
		r.logger.Error("failed to get kitchen_stations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get kitchen_stations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of kitchen_stations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*KitchenStations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_stations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM kitchen_stations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count kitchen_stations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, station_name
			, station_code
			, station_type
			, description
			, display_order
			, color_code
			, printer_id
			, is_active
			, auto_print_tickets
			, alert_sound_enabled
			, display_config
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM kitchen_stations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list kitchen_stations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list kitchen_stations: %w", err)
	}
	defer rows.Close()

	var entities []*KitchenStations
	for rows.Next() {
		var entity KitchenStations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.StationName,
			&entity.StationCode,
			&entity.StationType,
			&entity.Description,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.PrinterId,
			&entity.IsActive,
			&entity.AutoPrintTickets,
			&entity.AlertSoundEnabled,
			&entity.DisplayConfig,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan kitchen_stations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating kitchen_stations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing kitchen_stations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *KitchenStations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "kitchen_stations", duration, nil)
	}()

	query := `
		UPDATE kitchen_stations
		SET
			, organization_id = $2
			, location_id = $3
			, station_name = $4
			, station_code = $5
			, station_type = $6
			, description = $7
			, display_order = $8
			, color_code = $9
			, printer_id = $10
			, is_active = $11
			, auto_print_tickets = $12
			, alert_sound_enabled = $13
			, display_config = $14
			, updated_at = $16
			, created_by = $17
			, updated_by = $18
			, deleted_at = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.StationName,
		entity.StationCode,
		entity.StationType,
		entity.Description,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.PrinterId,
		entity.IsActive,
		entity.AutoPrintTickets,
		entity.AlertSoundEnabled,
		entity.DisplayConfig,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update kitchen_stations", zap.Error(err))
		return fmt.Errorf("failed to update kitchen_stations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("kitchen_stations not found or already deleted")
	}

	r.logger.Info("updated kitchen_stations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a kitchen_stations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "kitchen_stations", duration, nil)
	}()

	query := `
		UPDATE kitchen_stations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete kitchen_stations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete kitchen_stations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("kitchen_stations not found or already deleted")
	}

	r.logger.Info("deleted kitchen_stations", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves kitchen_stations records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*KitchenStations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "kitchen_stations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM kitchen_stations
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count kitchen_stations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, station_name
			, station_code
			, station_type
			, description
			, display_order
			, color_code
			, printer_id
			, is_active
			, auto_print_tickets
			, alert_sound_enabled
			, display_config
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM kitchen_stations
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list kitchen_stations by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list kitchen_stations: %w", err)
	}
	defer rows.Close()

	var entities []*KitchenStations
	for rows.Next() {
		var entity KitchenStations
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.StationName,
			&entity.StationCode,
			&entity.StationType,
			&entity.Description,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.PrinterId,
			&entity.IsActive,
			&entity.AutoPrintTickets,
			&entity.AlertSoundEnabled,
			&entity.DisplayConfig,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan kitchen_stations: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

