package permission

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

// Repository handles database operations for Permissions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Permissions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Permissions represents a permissions entity
type Permissions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Description *string `json:"description" db:"description"`
	Resource string `json:"resource" db:"resource"`
	Action string `json:"action" db:"action"`
	Category *string `json:"category" db:"category"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new permissions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Permissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "permissions", duration, nil)
	}()

	query := `
		INSERT INTO permissions (
			, name
			, slug
			, description
			, resource
			, action
			, category
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.Resource,
		entity.Action,
		entity.Category,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create permissions", zap.Error(err))
		return fmt.Errorf("failed to create permissions: %w", err)
	}

	r.logger.Info("created permissions",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a permissions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Permissions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "permissions", duration, nil)
	}()

	query := `
		SELECT
			id
			, name
			, slug
			, description
			, resource
			, action
			, category
			, created_at
			, updated_at
		FROM permissions
		WHERE id = $1
		
	`

	var entity Permissions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.Name,
		&entity.Slug,
		&entity.Description,
		&entity.Resource,
		&entity.Action,
		&entity.Category,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("permissions not found")
	}

	if err != nil {
		r.logger.Error("failed to get permissions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of permissions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Permissions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "permissions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM permissions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count permissions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, name
			, slug
			, description
			, resource
			, action
			, category
			, created_at
			, updated_at
		FROM permissions
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list permissions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	var entities []*Permissions
	for rows.Next() {
		var entity Permissions
		err := rows.Scan(
			&entity.Id,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.Resource,
			&entity.Action,
			&entity.Category,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan permissions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating permissions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing permissions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Permissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "permissions", duration, nil)
	}()

	query := `
		UPDATE permissions
		SET
			, name = $2
			, slug = $3
			, description = $4
			, resource = $5
			, action = $6
			, category = $7
			, updated_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.Resource,
		entity.Action,
		entity.Category,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update permissions", zap.Error(err))
		return fmt.Errorf("failed to update permissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("permissions not found or already deleted")
	}

	r.logger.Info("updated permissions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a permissions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "permissions", duration, nil)
	}()

	query := `DELETE FROM permissions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete permissions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete permissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("permissions not found")
	}

	r.logger.Info("deleted permissions", zap.String("id", id.String()))
	return nil
}



