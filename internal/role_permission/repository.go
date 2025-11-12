package role_permission

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

// Repository handles database operations for RolePermissions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new RolePermissions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// RolePermissions represents a role_permissions entity
type RolePermissions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	RoleId uuid.UUID `json:"role_id" db:"role_id"`
	PermissionId uuid.UUID `json:"permission_id" db:"permission_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new role_permissions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *RolePermissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "role_permissions", duration, nil)
	}()

	query := `
		INSERT INTO role_permissions (
			, role_id
			, permission_id
		) VALUES (
			, $2
			, $3
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.RoleId,
		entity.PermissionId,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create role_permissions", zap.Error(err))
		return fmt.Errorf("failed to create role_permissions: %w", err)
	}

	r.logger.Info("created role_permissions",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a role_permissions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*RolePermissions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "role_permissions", duration, nil)
	}()

	query := `
		SELECT
			id
			, role_id
			, permission_id
			, created_at
		FROM role_permissions
		WHERE id = $1
		
	`

	var entity RolePermissions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.RoleId,
		&entity.PermissionId,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("role_permissions not found")
	}

	if err != nil {
		r.logger.Error("failed to get role_permissions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get role_permissions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of role_permissions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*RolePermissions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "role_permissions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM role_permissions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count role_permissions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, role_id
			, permission_id
			, created_at
		FROM role_permissions
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list role_permissions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list role_permissions: %w", err)
	}
	defer rows.Close()

	var entities []*RolePermissions
	for rows.Next() {
		var entity RolePermissions
		err := rows.Scan(
			&entity.Id,
			&entity.RoleId,
			&entity.PermissionId,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan role_permissions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating role_permissions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing role_permissions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *RolePermissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "role_permissions", duration, nil)
	}()

	query := `
		UPDATE role_permissions
		SET
			, role_id = $2
			, permission_id = $3
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.RoleId,
		entity.PermissionId,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update role_permissions", zap.Error(err))
		return fmt.Errorf("failed to update role_permissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("role_permissions not found or already deleted")
	}

	r.logger.Info("updated role_permissions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a role_permissions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "role_permissions", duration, nil)
	}()

	query := `DELETE FROM role_permissions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete role_permissions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete role_permissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("role_permissions not found")
	}

	r.logger.Info("deleted role_permissions", zap.String("id", id.String()))
	return nil
}



