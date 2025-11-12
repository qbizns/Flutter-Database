package user_session

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

// Repository handles database operations for UserSessions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new UserSessions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// UserSessions represents a user_sessions entity
type UserSessions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	SessionToken string `json:"session_token" db:"session_token"`
	RefreshToken *string `json:"refresh_token" db:"refresh_token"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	DeviceType *string `json:"device_type" db:"device_type"`
	DeviceName *string `json:"device_name" db:"device_name"`
	Browser *string `json:"browser" db:"browser"`
	Os *string `json:"os" db:"os"`
	CountryCode *string `json:"country_code" db:"country_code"`
	City *string `json:"city" db:"city"`
	IsActive *bool `json:"is_active" db:"is_active"`
	LastActivityAt *time.Time `json:"last_activity_at" db:"last_activity_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
}

// Create inserts a new user_sessions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *UserSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "user_sessions", duration, nil)
	}()

	query := `
		INSERT INTO user_sessions (
			, user_id
			, organization_id
			, session_token
			, refresh_token
			, user_agent
			, ip_address
			, device_type
			, device_name
			, browser
			, os
			, country_code
			, city
			, is_active
			, last_activity_at
			, expires_at
			, revoked_at
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
			, $18
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.UserId,
		entity.OrganizationId,
		entity.SessionToken,
		entity.RefreshToken,
		entity.UserAgent,
		entity.IpAddress,
		entity.DeviceType,
		entity.DeviceName,
		entity.Browser,
		entity.Os,
		entity.CountryCode,
		entity.City,
		entity.IsActive,
		entity.LastActivityAt,
		entity.ExpiresAt,
		entity.RevokedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create user_sessions", zap.Error(err))
		return fmt.Errorf("failed to create user_sessions: %w", err)
	}

	r.logger.Info("created user_sessions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a user_sessions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*UserSessions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_sessions", duration, nil)
	}()

	query := `
		SELECT
			id
			, user_id
			, organization_id
			, session_token
			, refresh_token
			, user_agent
			, ip_address
			, device_type
			, device_name
			, browser
			, os
			, country_code
			, city
			, is_active
			, last_activity_at
			, expires_at
			, created_at
			, revoked_at
		FROM user_sessions
		WHERE id = $1
		
	`

	var entity UserSessions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.UserId,
		&entity.OrganizationId,
		&entity.SessionToken,
		&entity.RefreshToken,
		&entity.UserAgent,
		&entity.IpAddress,
		&entity.DeviceType,
		&entity.DeviceName,
		&entity.Browser,
		&entity.Os,
		&entity.CountryCode,
		&entity.City,
		&entity.IsActive,
		&entity.LastActivityAt,
		&entity.ExpiresAt,
		&entity.CreatedAt,
		&entity.RevokedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user_sessions not found")
	}

	if err != nil {
		r.logger.Error("failed to get user_sessions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get user_sessions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of user_sessions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*UserSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM user_sessions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, user_id
			, organization_id
			, session_token
			, refresh_token
			, user_agent
			, ip_address
			, device_type
			, device_name
			, browser
			, os
			, country_code
			, city
			, is_active
			, last_activity_at
			, expires_at
			, created_at
			, revoked_at
		FROM user_sessions
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list user_sessions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list user_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*UserSessions
	for rows.Next() {
		var entity UserSessions
		err := rows.Scan(
			&entity.Id,
			&entity.UserId,
			&entity.OrganizationId,
			&entity.SessionToken,
			&entity.RefreshToken,
			&entity.UserAgent,
			&entity.IpAddress,
			&entity.DeviceType,
			&entity.DeviceName,
			&entity.Browser,
			&entity.Os,
			&entity.CountryCode,
			&entity.City,
			&entity.IsActive,
			&entity.LastActivityAt,
			&entity.ExpiresAt,
			&entity.CreatedAt,
			&entity.RevokedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating user_sessions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing user_sessions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *UserSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "user_sessions", duration, nil)
	}()

	query := `
		UPDATE user_sessions
		SET
			, user_id = $2
			, organization_id = $3
			, session_token = $4
			, refresh_token = $5
			, user_agent = $6
			, ip_address = $7
			, device_type = $8
			, device_name = $9
			, browser = $10
			, os = $11
			, country_code = $12
			, city = $13
			, is_active = $14
			, last_activity_at = $15
			, expires_at = $16
			, revoked_at = $18
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $19
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.UserId,
		entity.OrganizationId,
		entity.SessionToken,
		entity.RefreshToken,
		entity.UserAgent,
		entity.IpAddress,
		entity.DeviceType,
		entity.DeviceName,
		entity.Browser,
		entity.Os,
		entity.CountryCode,
		entity.City,
		entity.IsActive,
		entity.LastActivityAt,
		entity.ExpiresAt,
		entity.RevokedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update user_sessions", zap.Error(err))
		return fmt.Errorf("failed to update user_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_sessions not found or already deleted")
	}

	r.logger.Info("updated user_sessions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a user_sessions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "user_sessions", duration, nil)
	}()

	query := `DELETE FROM user_sessions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user_sessions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete user_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_sessions not found")
	}

	r.logger.Info("deleted user_sessions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves user_sessions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*UserSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM user_sessions
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, user_id
			, organization_id
			, session_token
			, refresh_token
			, user_agent
			, ip_address
			, device_type
			, device_name
			, browser
			, os
			, country_code
			, city
			, is_active
			, last_activity_at
			, expires_at
			, created_at
			, revoked_at
		FROM user_sessions
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list user_sessions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list user_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*UserSessions
	for rows.Next() {
		var entity UserSessions
		err := rows.Scan(
			&entity.Id,
			&entity.UserId,
			&entity.OrganizationId,
			&entity.SessionToken,
			&entity.RefreshToken,
			&entity.UserAgent,
			&entity.IpAddress,
			&entity.DeviceType,
			&entity.DeviceName,
			&entity.Browser,
			&entity.Os,
			&entity.CountryCode,
			&entity.City,
			&entity.IsActive,
			&entity.LastActivityAt,
			&entity.ExpiresAt,
			&entity.CreatedAt,
			&entity.RevokedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

