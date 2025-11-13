package cours

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

// Repository handles database operations for Courses
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Courses repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Courses represents a courses entity
type Courses struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	CourseName string `json:"course_name" db:"course_name"`
	CourseCode *string `json:"course_code" db:"course_code"`
	CourseType *string `json:"course_type" db:"course_type"`
	TypicalDurationMinutes *int64 `json:"typical_duration_minutes" db:"typical_duration_minutes"`
	FireDelayMinutes *int64 `json:"fire_delay_minutes" db:"fire_delay_minutes"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	ColorCode *string `json:"color_code" db:"color_code"`
	Icon *string `json:"icon" db:"icon"`
	IsActive *bool `json:"is_active" db:"is_active"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new courses record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Courses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "courses", duration, nil)
	}()

	query := `
		INSERT INTO courses (
			, organization_id
			, course_name
			, course_code
			, course_type
			, typical_duration_minutes
			, fire_delay_minutes
			, display_order
			, color_code
			, icon
			, is_active
			, is_default
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
			, $15
			, $18
			, $19
			, $20
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.CourseName,
		entity.CourseCode,
		entity.CourseType,
		entity.TypicalDurationMinutes,
		entity.FireDelayMinutes,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.Icon,
		entity.IsActive,
		entity.IsDefault,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create courses", zap.Error(err))
		return fmt.Errorf("failed to create courses: %w", err)
	}

	r.logger.Info("created courses",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a courses by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Courses, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "courses", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, course_name
			, course_code
			, course_type
			, typical_duration_minutes
			, fire_delay_minutes
			, display_order
			, color_code
			, icon
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM courses
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Courses
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.CourseName,
		&entity.CourseCode,
		&entity.CourseType,
		&entity.TypicalDurationMinutes,
		&entity.FireDelayMinutes,
		&entity.DisplayOrder,
		&entity.ColorCode,
		&entity.Icon,
		&entity.IsActive,
		&entity.IsDefault,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("courses not found")
	}

	if err != nil {
		r.logger.Error("failed to get courses", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of courses records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Courses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "courses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM courses
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, course_name
			, course_code
			, course_type
			, typical_duration_minutes
			, fire_delay_minutes
			, display_order
			, color_code
			, icon
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM courses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list courses", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list courses: %w", err)
	}
	defer rows.Close()

	var entities []*Courses
	for rows.Next() {
		var entity Courses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CourseName,
			&entity.CourseCode,
			&entity.CourseType,
			&entity.TypicalDurationMinutes,
			&entity.FireDelayMinutes,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.Icon,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan courses: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating courses rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing courses record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Courses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "courses", duration, nil)
	}()

	query := `
		UPDATE courses
		SET
			, organization_id = $2
			, course_name = $3
			, course_code = $4
			, course_type = $5
			, typical_duration_minutes = $6
			, fire_delay_minutes = $7
			, display_order = $8
			, color_code = $9
			, icon = $10
			, is_active = $11
			, is_default = $12
			, description = $13
			, notes = $14
			, metadata = $15
			, updated_at = $17
			, deleted_at = $18
			, created_by = $19
			, updated_by = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.CourseName,
		entity.CourseCode,
		entity.CourseType,
		entity.TypicalDurationMinutes,
		entity.FireDelayMinutes,
		entity.DisplayOrder,
		entity.ColorCode,
		entity.Icon,
		entity.IsActive,
		entity.IsDefault,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update courses", zap.Error(err))
		return fmt.Errorf("failed to update courses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("courses not found or already deleted")
	}

	r.logger.Info("updated courses",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a courses record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "courses", duration, nil)
	}()

	query := `
		UPDATE courses
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete courses", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete courses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("courses not found or already deleted")
	}

	r.logger.Info("deleted courses", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves courses records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Courses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "courses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM courses
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, course_name
			, course_code
			, course_type
			, typical_duration_minutes
			, fire_delay_minutes
			, display_order
			, color_code
			, icon
			, is_active
			, is_default
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM courses
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list courses by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list courses: %w", err)
	}
	defer rows.Close()

	var entities []*Courses
	for rows.Next() {
		var entity Courses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.CourseName,
			&entity.CourseCode,
			&entity.CourseType,
			&entity.TypicalDurationMinutes,
			&entity.FireDelayMinutes,
			&entity.DisplayOrder,
			&entity.ColorCode,
			&entity.Icon,
			&entity.IsActive,
			&entity.IsDefault,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan courses: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

