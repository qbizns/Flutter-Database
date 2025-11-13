package user_role

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

// Repository handles database operations for UserRoles
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new UserRoles repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// UserRoles represents a user_roles entity
type UserRoles struct {
	Id *uuid.UUID `json:"id" db:"id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	RoleId uuid.UUID `json:"role_id" db:"role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	AssignedBy *uuid.UUID `json:"assigned_by" db:"assigned_by"`
}

// Create inserts a new user_roles record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *UserRoles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "user_roles", duration, nil)
	}()

	query := `
		INSERT INTO user_roles (
			, user_id
			, role_id
			, assigned_by
		) VALUES (
			, $2
			, $3
			, $5
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.UserId,
		entity.RoleId,
		entity.AssignedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create user_roles", zap.Error(err))
		return fmt.Errorf("failed to create user_roles: %w", err)
	}

	r.logger.Info("created user_roles",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a user_roles by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*UserRoles, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_roles", duration, nil)
	}()

	query := `
		SELECT
			id
			, user_id
			, role_id
			, created_at
			, assigned_by
		FROM user_roles
		WHERE id = $1
		
	`

	var entity UserRoles
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.UserId,
		&entity.RoleId,
		&entity.CreatedAt,
		&entity.AssignedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user_roles not found")
	}

	if err != nil {
		r.logger.Error("failed to get user_roles", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get user_roles: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of user_roles records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*UserRoles, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_roles", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM user_roles
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user_roles records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, user_id
			, role_id
			, created_at
			, assigned_by
		FROM user_roles
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list user_roles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list user_roles: %w", err)
	}
	defer rows.Close()

	var entities []*UserRoles
	for rows.Next() {
		var entity UserRoles
		err := rows.Scan(
			&entity.Id,
			&entity.UserId,
			&entity.RoleId,
			&entity.CreatedAt,
			&entity.AssignedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user_roles: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating user_roles rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing user_roles record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *UserRoles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "user_roles", duration, nil)
	}()

	query := `
		UPDATE user_roles
		SET
			, user_id = $2
			, role_id = $3
			, assigned_by = $5
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.UserId,
		entity.RoleId,
		entity.AssignedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update user_roles", zap.Error(err))
		return fmt.Errorf("failed to update user_roles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_roles not found or already deleted")
	}

	r.logger.Info("updated user_roles",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a user_roles record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "user_roles", duration, nil)
	}()

	query := `DELETE FROM user_roles WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user_roles", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete user_roles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_roles not found")
	}

	r.logger.Info("deleted user_roles", zap.String("id", id.String()))
	return nil
}



