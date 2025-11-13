package user

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

// Repository handles database operations for Users
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Users repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Users represents a users entity
type Users struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	Email string `json:"email" db:"email"`
	PasswordHash *string `json:"password_hash" db:"password_hash"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName string `json:"last_name" db:"last_name"`
	Phone *string `json:"phone" db:"phone"`
	AvatarUrl *string `json:"avatar_url" db:"avatar_url"`
	Status string `json:"status" db:"status"`
	EmailVerified *bool `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" db:"email_verified_at"`
	LastLoginAt *time.Time `json:"last_login_at" db:"last_login_at"`
	LastLoginIp *string `json:"last_login_ip" db:"last_login_ip"`
	FailedLoginAttempts *int64 `json:"failed_login_attempts" db:"failed_login_attempts"`
	LockedUntil *time.Time `json:"locked_until" db:"locked_until"`
	TwoFactorEnabled *bool `json:"two_factor_enabled" db:"two_factor_enabled"`
	TwoFactorSecret *string `json:"two_factor_secret" db:"two_factor_secret"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new users record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Users) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "users", duration, nil)
	}()

	query := `
		INSERT INTO users (
			, organization_id
			, email
			, password_hash
			, first_name
			, last_name
			, phone
			, avatar_url
			, status
			, email_verified
			, email_verified_at
			, last_login_at
			, last_login_ip
			, failed_login_attempts
			, locked_until
			, two_factor_enabled
			, two_factor_secret
			, settings
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
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Email,
		entity.PasswordHash,
		entity.FirstName,
		entity.LastName,
		entity.Phone,
		entity.AvatarUrl,
		entity.Status,
		entity.EmailVerified,
		entity.EmailVerifiedAt,
		entity.LastLoginAt,
		entity.LastLoginIp,
		entity.FailedLoginAttempts,
		entity.LockedUntil,
		entity.TwoFactorEnabled,
		entity.TwoFactorSecret,
		entity.Settings,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create users", zap.Error(err))
		return fmt.Errorf("failed to create users: %w", err)
	}

	r.logger.Info("created users",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a users by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Users, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "users", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, email
			, password_hash
			, first_name
			, last_name
			, phone
			, avatar_url
			, status
			, email_verified
			, email_verified_at
			, last_login_at
			, last_login_ip
			, failed_login_attempts
			, locked_until
			, two_factor_enabled
			, two_factor_secret
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM users
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Users
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Email,
		&entity.PasswordHash,
		&entity.FirstName,
		&entity.LastName,
		&entity.Phone,
		&entity.AvatarUrl,
		&entity.Status,
		&entity.EmailVerified,
		&entity.EmailVerifiedAt,
		&entity.LastLoginAt,
		&entity.LastLoginIp,
		&entity.FailedLoginAttempts,
		&entity.LockedUntil,
		&entity.TwoFactorEnabled,
		&entity.TwoFactorSecret,
		&entity.Settings,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("users not found")
	}

	if err != nil {
		r.logger.Error("failed to get users", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of users records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Users, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "users", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, email
			, password_hash
			, first_name
			, last_name
			, phone
			, avatar_url
			, status
			, email_verified
			, email_verified_at
			, last_login_at
			, last_login_ip
			, failed_login_attempts
			, locked_until
			, two_factor_enabled
			, two_factor_secret
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var entities []*Users
	for rows.Next() {
		var entity Users
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Email,
			&entity.PasswordHash,
			&entity.FirstName,
			&entity.LastName,
			&entity.Phone,
			&entity.AvatarUrl,
			&entity.Status,
			&entity.EmailVerified,
			&entity.EmailVerifiedAt,
			&entity.LastLoginAt,
			&entity.LastLoginIp,
			&entity.FailedLoginAttempts,
			&entity.LockedUntil,
			&entity.TwoFactorEnabled,
			&entity.TwoFactorSecret,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating users rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing users record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Users) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "users", duration, nil)
	}()

	query := `
		UPDATE users
		SET
			, organization_id = $2
			, email = $3
			, password_hash = $4
			, first_name = $5
			, last_name = $6
			, phone = $7
			, avatar_url = $8
			, status = $9
			, email_verified = $10
			, email_verified_at = $11
			, last_login_at = $12
			, last_login_ip = $13
			, failed_login_attempts = $14
			, locked_until = $15
			, two_factor_enabled = $16
			, two_factor_secret = $17
			, settings = $18
			, metadata = $19
			, updated_at = $21
			, deleted_at = $22
			, created_by = $23
			, updated_by = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Email,
		entity.PasswordHash,
		entity.FirstName,
		entity.LastName,
		entity.Phone,
		entity.AvatarUrl,
		entity.Status,
		entity.EmailVerified,
		entity.EmailVerifiedAt,
		entity.LastLoginAt,
		entity.LastLoginIp,
		entity.FailedLoginAttempts,
		entity.LockedUntil,
		entity.TwoFactorEnabled,
		entity.TwoFactorSecret,
		entity.Settings,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update users", zap.Error(err))
		return fmt.Errorf("failed to update users: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("users not found or already deleted")
	}

	r.logger.Info("updated users",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a users record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "users", duration, nil)
	}()

	query := `
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete users", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete users: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("users not found or already deleted")
	}

	r.logger.Info("deleted users", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves users records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Users, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "users", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, email
			, password_hash
			, first_name
			, last_name
			, phone
			, avatar_url
			, status
			, email_verified
			, email_verified_at
			, last_login_at
			, last_login_ip
			, failed_login_attempts
			, locked_until
			, two_factor_enabled
			, two_factor_secret
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM users
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list users by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var entities []*Users
	for rows.Next() {
		var entity Users
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Email,
			&entity.PasswordHash,
			&entity.FirstName,
			&entity.LastName,
			&entity.Phone,
			&entity.AvatarUrl,
			&entity.Status,
			&entity.EmailVerified,
			&entity.EmailVerifiedAt,
			&entity.LastLoginAt,
			&entity.LastLoginIp,
			&entity.FailedLoginAttempts,
			&entity.LockedUntil,
			&entity.TwoFactorEnabled,
			&entity.TwoFactorSecret,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan users: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

