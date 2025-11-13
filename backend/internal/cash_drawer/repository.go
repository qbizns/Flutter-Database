package cash_drawer

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

// Repository handles database operations for CashDrawers
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new CashDrawers repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// CashDrawers represents a cash_drawers entity
type CashDrawers struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	DrawerCode string `json:"drawer_code" db:"drawer_code"`
	DrawerName string `json:"drawer_name" db:"drawer_name"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	DeviceId *uuid.UUID `json:"device_id" db:"device_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new cash_drawers record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *CashDrawers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "cash_drawers", duration, nil)
	}()

	query := `
		INSERT INTO cash_drawers (
			, organization_id
			, drawer_code
			, drawer_name
			, location_id
			, device_id
			, is_active
			, notes
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
			, $12
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.DrawerCode,
		entity.DrawerName,
		entity.LocationId,
		entity.DeviceId,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create cash_drawers", zap.Error(err))
		return fmt.Errorf("failed to create cash_drawers: %w", err)
	}

	r.logger.Info("created cash_drawers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a cash_drawers by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*CashDrawers, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawers", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, drawer_code
			, drawer_name
			, location_id
			, device_id
			, is_active
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM cash_drawers
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity CashDrawers
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.DrawerCode,
		&entity.DrawerName,
		&entity.LocationId,
		&entity.DeviceId,
		&entity.IsActive,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("cash_drawers not found")
	}

	if err != nil {
		r.logger.Error("failed to get cash_drawers", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get cash_drawers: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of cash_drawers records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*CashDrawers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_drawers
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_drawers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, drawer_code
			, drawer_name
			, location_id
			, device_id
			, is_active
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM cash_drawers
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_drawers", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_drawers: %w", err)
	}
	defer rows.Close()

	var entities []*CashDrawers
	for rows.Next() {
		var entity CashDrawers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DrawerCode,
			&entity.DrawerName,
			&entity.LocationId,
			&entity.DeviceId,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_drawers: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating cash_drawers rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing cash_drawers record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *CashDrawers) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_drawers", duration, nil)
	}()

	query := `
		UPDATE cash_drawers
		SET
			, organization_id = $2
			, drawer_code = $3
			, drawer_name = $4
			, location_id = $5
			, device_id = $6
			, is_active = $7
			, notes = $8
			, created_by = $9
			, updated_at = $11
			, deleted_at = $12
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.DrawerCode,
		entity.DrawerName,
		entity.LocationId,
		entity.DeviceId,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update cash_drawers", zap.Error(err))
		return fmt.Errorf("failed to update cash_drawers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_drawers not found or already deleted")
	}

	r.logger.Info("updated cash_drawers",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a cash_drawers record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "cash_drawers", duration, nil)
	}()

	query := `
		UPDATE cash_drawers
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete cash_drawers", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete cash_drawers: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("cash_drawers not found or already deleted")
	}

	r.logger.Info("deleted cash_drawers", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves cash_drawers records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*CashDrawers, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "cash_drawers", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM cash_drawers
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cash_drawers records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, drawer_code
			, drawer_name
			, location_id
			, device_id
			, is_active
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM cash_drawers
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list cash_drawers by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list cash_drawers: %w", err)
	}
	defer rows.Close()

	var entities []*CashDrawers
	for rows.Next() {
		var entity CashDrawers
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DrawerCode,
			&entity.DrawerName,
			&entity.LocationId,
			&entity.DeviceId,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cash_drawers: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

