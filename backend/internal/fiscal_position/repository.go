package fiscal_position

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

// Repository handles database operations for FiscalPositions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FiscalPositions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FiscalPositions represents a fiscal_positions entity
type FiscalPositions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PositionCode string `json:"position_code" db:"position_code"`
	PositionName string `json:"position_name" db:"position_name"`
	AutoApply *bool `json:"auto_apply" db:"auto_apply"`
	CountryId *string `json:"country_id" db:"country_id"`
	StateProvince *string `json:"state_province" db:"state_province"`
	ZipPostalCodeRange *string `json:"zip_postal_code_range" db:"zip_postal_code_range"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new fiscal_positions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FiscalPositions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "fiscal_positions", duration, nil)
	}()

	query := `
		INSERT INTO fiscal_positions (
			, organization_id
			, position_code
			, position_name
			, auto_apply
			, country_id
			, state_province
			, zip_postal_code_range
			, is_active
			, notes
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
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PositionCode,
		entity.PositionName,
		entity.AutoApply,
		entity.CountryId,
		entity.StateProvince,
		entity.ZipPostalCodeRange,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create fiscal_positions", zap.Error(err))
		return fmt.Errorf("failed to create fiscal_positions: %w", err)
	}

	r.logger.Info("created fiscal_positions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a fiscal_positions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FiscalPositions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_positions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, position_code
			, position_name
			, auto_apply
			, country_id
			, state_province
			, zip_postal_code_range
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM fiscal_positions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FiscalPositions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PositionCode,
		&entity.PositionName,
		&entity.AutoApply,
		&entity.CountryId,
		&entity.StateProvince,
		&entity.ZipPostalCodeRange,
		&entity.IsActive,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fiscal_positions not found")
	}

	if err != nil {
		r.logger.Error("failed to get fiscal_positions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get fiscal_positions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of fiscal_positions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FiscalPositions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_positions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fiscal_positions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fiscal_positions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, position_code
			, position_name
			, auto_apply
			, country_id
			, state_province
			, zip_postal_code_range
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM fiscal_positions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fiscal_positions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fiscal_positions: %w", err)
	}
	defer rows.Close()

	var entities []*FiscalPositions
	for rows.Next() {
		var entity FiscalPositions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PositionCode,
			&entity.PositionName,
			&entity.AutoApply,
			&entity.CountryId,
			&entity.StateProvince,
			&entity.ZipPostalCodeRange,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fiscal_positions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating fiscal_positions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing fiscal_positions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FiscalPositions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_positions", duration, nil)
	}()

	query := `
		UPDATE fiscal_positions
		SET
			, organization_id = $2
			, position_code = $3
			, position_name = $4
			, auto_apply = $5
			, country_id = $6
			, state_province = $7
			, zip_postal_code_range = $8
			, is_active = $9
			, notes = $10
			, created_by = $11
			, updated_by = $12
			, updated_at = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PositionCode,
		entity.PositionName,
		entity.AutoApply,
		entity.CountryId,
		entity.StateProvince,
		entity.ZipPostalCodeRange,
		entity.IsActive,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update fiscal_positions", zap.Error(err))
		return fmt.Errorf("failed to update fiscal_positions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_positions not found or already deleted")
	}

	r.logger.Info("updated fiscal_positions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a fiscal_positions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fiscal_positions", duration, nil)
	}()

	query := `
		UPDATE fiscal_positions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete fiscal_positions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete fiscal_positions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fiscal_positions not found or already deleted")
	}

	r.logger.Info("deleted fiscal_positions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves fiscal_positions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*FiscalPositions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fiscal_positions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fiscal_positions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fiscal_positions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, position_code
			, position_name
			, auto_apply
			, country_id
			, state_province
			, zip_postal_code_range
			, is_active
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM fiscal_positions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fiscal_positions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fiscal_positions: %w", err)
	}
	defer rows.Close()

	var entities []*FiscalPositions
	for rows.Next() {
		var entity FiscalPositions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PositionCode,
			&entity.PositionName,
			&entity.AutoApply,
			&entity.CountryId,
			&entity.StateProvince,
			&entity.ZipPostalCodeRange,
			&entity.IsActive,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fiscal_positions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

