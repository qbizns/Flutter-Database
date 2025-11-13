package restaurant_table

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

// Repository handles database operations for RestaurantTables
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new RestaurantTables repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// RestaurantTables represents a restaurant_tables entity
type RestaurantTables struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	FloorPlanId *uuid.UUID `json:"floor_plan_id" db:"floor_plan_id"`
	SectionId *uuid.UUID `json:"section_id" db:"section_id"`
	TableNumber string `json:"table_number" db:"table_number"`
	TableName *string `json:"table_name" db:"table_name"`
	MinCapacity *int64 `json:"min_capacity" db:"min_capacity"`
	MaxCapacity int64 `json:"max_capacity" db:"max_capacity"`
	TableShape *string `json:"table_shape" db:"table_shape"`
	IsCombinable *bool `json:"is_combinable" db:"is_combinable"`
	PositionX *float64 `json:"position_x" db:"position_x"`
	PositionY *float64 `json:"position_y" db:"position_y"`
	Rotation *int64 `json:"rotation" db:"rotation"`
	Status *string `json:"status" db:"status"`
	CurrentCovers *int64 `json:"current_covers" db:"current_covers"`
	SeatedAt *time.Time `json:"seated_at" db:"seated_at"`
	CurrentWaiterId *uuid.UUID `json:"current_waiter_id" db:"current_waiter_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	AllowOnlineReservation *bool `json:"allow_online_reservation" db:"allow_online_reservation"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	ColorCode *string `json:"color_code" db:"color_code"`
	Icon *string `json:"icon" db:"icon"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new restaurant_tables record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *RestaurantTables) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "restaurant_tables", duration, nil)
	}()

	query := `
		INSERT INTO restaurant_tables (
			, organization_id
			, location_id
			, floor_plan_id
			, section_id
			, table_number
			, table_name
			, min_capacity
			, max_capacity
			, table_shape
			, is_combinable
			, position_x
			, position_y
			, rotation
			, status
			, current_covers
			, seated_at
			, current_waiter_id
			, is_active
			, allow_online_reservation
			, display_order
			, color_code
			, icon
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
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
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $28
			, $29
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorPlanId,
		entity.SectionId,
		entity.TableNumber,
		entity.TableName,
		entity.MinCapacity,
		entity.MaxCapacity,
		entity.TableShape,
		entity.IsCombinable,
		entity.PositionX,
		entity.PositionY,
		entity.Rotation,
		entity.Status,
		entity.CurrentCovers,
		entity.SeatedAt,
		entity.CurrentWaiterId,
		entity.IsActive,
		entity.AllowOnlineReservation,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.Icon,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create restaurant_tables", zap.Error(err))
		return fmt.Errorf("failed to create restaurant_tables: %w", err)
	}

	r.logger.Info("created restaurant_tables",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a restaurant_tables by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*RestaurantTables, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "restaurant_tables", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_id
			, table_number
			, table_name
			, min_capacity
			, max_capacity
			, table_shape
			, is_combinable
			, position_x
			, position_y
			, rotation
			, status
			, current_covers
			, seated_at
			, current_waiter_id
			, is_active
			, allow_online_reservation
			, display_order
			, color_code
			, icon
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM restaurant_tables
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity RestaurantTables
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.FloorPlanId,
		&entity.SectionId,
		&entity.TableNumber,
		&entity.TableName,
		&entity.MinCapacity,
		&entity.MaxCapacity,
		&entity.TableShape,
		&entity.IsCombinable,
		&entity.PositionX,
		&entity.PositionY,
		&entity.Rotation,
		&entity.Status,
		&entity.CurrentCovers,
		&entity.SeatedAt,
		&entity.CurrentWaiterId,
		&entity.IsActive,
		&entity.AllowOnlineReservation,
		&entity.DisplayOrder,
		&entity.ColorCode,
		&entity.Icon,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("restaurant_tables not found")
	}

	if err != nil {
		r.logger.Error("failed to get restaurant_tables", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get restaurant_tables: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of restaurant_tables records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*RestaurantTables, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "restaurant_tables", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM restaurant_tables
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count restaurant_tables records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_id
			, table_number
			, table_name
			, min_capacity
			, max_capacity
			, table_shape
			, is_combinable
			, position_x
			, position_y
			, rotation
			, status
			, current_covers
			, seated_at
			, current_waiter_id
			, is_active
			, allow_online_reservation
			, display_order
			, color_code
			, icon
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM restaurant_tables
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list restaurant_tables", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list restaurant_tables: %w", err)
	}
	defer rows.Close()

	var entities []*RestaurantTables
	for rows.Next() {
		var entity RestaurantTables
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorPlanId,
			&entity.SectionId,
			&entity.TableNumber,
			&entity.TableName,
			&entity.MinCapacity,
			&entity.MaxCapacity,
			&entity.TableShape,
			&entity.IsCombinable,
			&entity.PositionX,
			&entity.PositionY,
			&entity.Rotation,
			&entity.Status,
			&entity.CurrentCovers,
			&entity.SeatedAt,
			&entity.CurrentWaiterId,
			&entity.IsActive,
			&entity.AllowOnlineReservation,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.Icon,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan restaurant_tables: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating restaurant_tables rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing restaurant_tables record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *RestaurantTables) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "restaurant_tables", duration, nil)
	}()

	query := `
		UPDATE restaurant_tables
		SET
			, organization_id = $2
			, location_id = $3
			, floor_plan_id = $4
			, section_id = $5
			, table_number = $6
			, table_name = $7
			, min_capacity = $8
			, max_capacity = $9
			, table_shape = $10
			, is_combinable = $11
			, position_x = $12
			, position_y = $13
			, rotation = $14
			, status = $15
			, current_covers = $16
			, seated_at = $17
			, current_waiter_id = $18
			, is_active = $19
			, allow_online_reservation = $20
			, display_order = $21
			, color_code = $22
			, icon = $23
			, notes = $24
			, metadata = $25
			, updated_at = $27
			, deleted_at = $28
			, created_by = $29
			, updated_by = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorPlanId,
		entity.SectionId,
		entity.TableNumber,
		entity.TableName,
		entity.MinCapacity,
		entity.MaxCapacity,
		entity.TableShape,
		entity.IsCombinable,
		entity.PositionX,
		entity.PositionY,
		entity.Rotation,
		entity.Status,
		entity.CurrentCovers,
		entity.SeatedAt,
		entity.CurrentWaiterId,
		entity.IsActive,
		entity.AllowOnlineReservation,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.Icon,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update restaurant_tables", zap.Error(err))
		return fmt.Errorf("failed to update restaurant_tables: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("restaurant_tables not found or already deleted")
	}

	r.logger.Info("updated restaurant_tables",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a restaurant_tables record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "restaurant_tables", duration, nil)
	}()

	query := `
		UPDATE restaurant_tables
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete restaurant_tables", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete restaurant_tables: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("restaurant_tables not found or already deleted")
	}

	r.logger.Info("deleted restaurant_tables", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves restaurant_tables records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*RestaurantTables, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "restaurant_tables", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM restaurant_tables
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count restaurant_tables records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_id
			, table_number
			, table_name
			, min_capacity
			, max_capacity
			, table_shape
			, is_combinable
			, position_x
			, position_y
			, rotation
			, status
			, current_covers
			, seated_at
			, current_waiter_id
			, is_active
			, allow_online_reservation
			, display_order
			, color_code
			, icon
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM restaurant_tables
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list restaurant_tables by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list restaurant_tables: %w", err)
	}
	defer rows.Close()

	var entities []*RestaurantTables
	for rows.Next() {
		var entity RestaurantTables
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorPlanId,
			&entity.SectionId,
			&entity.TableNumber,
			&entity.TableName,
			&entity.MinCapacity,
			&entity.MaxCapacity,
			&entity.TableShape,
			&entity.IsCombinable,
			&entity.PositionX,
			&entity.PositionY,
			&entity.Rotation,
			&entity.Status,
			&entity.CurrentCovers,
			&entity.SeatedAt,
			&entity.CurrentWaiterId,
			&entity.IsActive,
			&entity.AllowOnlineReservation,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.Icon,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan restaurant_tables: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

