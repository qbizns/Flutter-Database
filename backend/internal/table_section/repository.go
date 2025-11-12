package table_section

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

// Repository handles database operations for TableSections
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TableSections repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TableSections represents a table_sections entity
type TableSections struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	FloorPlanId *uuid.UUID `json:"floor_plan_id" db:"floor_plan_id"`
	SectionName string `json:"section_name" db:"section_name"`
	SectionCode *string `json:"section_code" db:"section_code"`
	SectionType *string `json:"section_type" db:"section_type"`
	ColorCode *string `json:"color_code" db:"color_code"`
	Icon *string `json:"icon" db:"icon"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new table_sections record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TableSections) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "table_sections", duration, nil)
	}()

	query := `
		INSERT INTO table_sections (
			, organization_id
			, location_id
			, floor_plan_id
			, section_name
			, section_code
			, section_type
			, color_code
			, icon
			, display_order
			, is_active
			, description
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
			, $17
			, $18
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorPlanId,
		entity.SectionName,
		entity.SectionCode,
		entity.SectionType,
		entity.ColorCode,
		entity.Icon,
		entity.DisplayOrder,
		entity.IsActive,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create table_sections", zap.Error(err))
		return fmt.Errorf("failed to create table_sections: %w", err)
	}

	r.logger.Info("created table_sections",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a table_sections by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TableSections, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "table_sections", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_name
			, section_code
			, section_type
			, color_code
			, icon
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM table_sections
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TableSections
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.FloorPlanId,
		&entity.SectionName,
		&entity.SectionCode,
		&entity.SectionType,
		&entity.ColorCode,
		&entity.Icon,
		&entity.DisplayOrder,
		&entity.IsActive,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("table_sections not found")
	}

	if err != nil {
		r.logger.Error("failed to get table_sections", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get table_sections: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of table_sections records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TableSections, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "table_sections", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM table_sections
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count table_sections records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_name
			, section_code
			, section_type
			, color_code
			, icon
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM table_sections
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list table_sections", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list table_sections: %w", err)
	}
	defer rows.Close()

	var entities []*TableSections
	for rows.Next() {
		var entity TableSections
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorPlanId,
			&entity.SectionName,
			&entity.SectionCode,
			&entity.SectionType,
			&entity.ColorCode,
			&entity.Icon,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan table_sections: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating table_sections rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing table_sections record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TableSections) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "table_sections", duration, nil)
	}()

	query := `
		UPDATE table_sections
		SET
			, organization_id = $2
			, location_id = $3
			, floor_plan_id = $4
			, section_name = $5
			, section_code = $6
			, section_type = $7
			, color_code = $8
			, icon = $9
			, display_order = $10
			, is_active = $11
			, description = $12
			, notes = $13
			, metadata = $14
			, updated_at = $16
			, deleted_at = $17
			, created_by = $18
			, updated_by = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.FloorPlanId,
		entity.SectionName,
		entity.SectionCode,
		entity.SectionType,
		entity.ColorCode,
		entity.Icon,
		entity.DisplayOrder,
		entity.IsActive,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update table_sections", zap.Error(err))
		return fmt.Errorf("failed to update table_sections: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("table_sections not found or already deleted")
	}

	r.logger.Info("updated table_sections",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a table_sections record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "table_sections", duration, nil)
	}()

	query := `
		UPDATE table_sections
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete table_sections", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete table_sections: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("table_sections not found or already deleted")
	}

	r.logger.Info("deleted table_sections", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves table_sections records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TableSections, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "table_sections", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM table_sections
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count table_sections records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, floor_plan_id
			, section_name
			, section_code
			, section_type
			, color_code
			, icon
			, display_order
			, is_active
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM table_sections
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list table_sections by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list table_sections: %w", err)
	}
	defer rows.Close()

	var entities []*TableSections
	for rows.Next() {
		var entity TableSections
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.FloorPlanId,
			&entity.SectionName,
			&entity.SectionCode,
			&entity.SectionType,
			&entity.ColorCode,
			&entity.Icon,
			&entity.DisplayOrder,
			&entity.IsActive,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan table_sections: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

