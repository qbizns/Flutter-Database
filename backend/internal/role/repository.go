package role

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

// Repository handles database operations for Roles
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Roles repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Roles represents a roles entity
type Roles struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Description *string `json:"description" db:"description"`
	IsSystemRole *bool `json:"is_system_role" db:"is_system_role"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new roles record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Roles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "roles", duration, nil)
	}()

	query := `
		INSERT INTO roles (
			, organization_id
			, name
			, slug
			, description
			, is_system_role
			, is_default
			, settings
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsDefault,
		entity.Settings,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create roles", zap.Error(err))
		return fmt.Errorf("failed to create roles: %w", err)
	}

	r.logger.Info("created roles",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a roles by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Roles, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "roles", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, is_system_role
			, is_default
			, settings
			, created_at
			, updated_at
			, deleted_at
		FROM roles
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Roles
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Name,
		&entity.Slug,
		&entity.Description,
		&entity.IsSystemRole,
		&entity.IsDefault,
		&entity.Settings,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("roles not found")
	}

	if err != nil {
		r.logger.Error("failed to get roles", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of roles records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Roles, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "roles", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM roles
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, is_system_role
			, is_default
			, settings
			, created_at
			, updated_at
			, deleted_at
		FROM roles
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list roles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var entities []*Roles
	for rows.Next() {
		var entity Roles
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.IsSystemRole,
			&entity.IsDefault,
			&entity.Settings,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan roles: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating roles rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing roles record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Roles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "roles", duration, nil)
	}()

	query := `
		UPDATE roles
		SET
			, organization_id = $2
			, name = $3
			, slug = $4
			, description = $5
			, is_system_role = $6
			, is_default = $7
			, settings = $8
			, updated_at = $10
			, deleted_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.IsSystemRole,
		entity.IsDefault,
		entity.Settings,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update roles", zap.Error(err))
		return fmt.Errorf("failed to update roles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("roles not found or already deleted")
	}

	r.logger.Info("updated roles",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a roles record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "roles", duration, nil)
	}()

	query := `
		UPDATE roles
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete roles", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete roles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("roles not found or already deleted")
	}

	r.logger.Info("deleted roles", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves roles records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Roles, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "roles", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM roles
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, is_system_role
			, is_default
			, settings
			, created_at
			, updated_at
			, deleted_at
		FROM roles
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list roles by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var entities []*Roles
	for rows.Next() {
		var entity Roles
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.IsSystemRole,
			&entity.IsDefault,
			&entity.Settings,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan roles: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

