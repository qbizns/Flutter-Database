package user_setting

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

// Repository handles database operations for UserSettings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new UserSettings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// UserSettings represents a user_settings entity
type UserSettings struct {
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	Theme *string `json:"theme" db:"theme"`
	Language *string `json:"language" db:"language"`
	Timezone *string `json:"timezone" db:"timezone"`
	DefaultDashboard *string `json:"default_dashboard" db:"default_dashboard"`
	DashboardLayout json.RawMessage `json:"dashboard_layout" db:"dashboard_layout"`
	ItemsPerPage *int64 `json:"items_per_page" db:"items_per_page"`
	DefaultView *string `json:"default_view" db:"default_view"`
	DesktopNotifications *bool `json:"desktop_notifications" db:"desktop_notifications"`
	SoundNotifications *bool `json:"sound_notifications" db:"sound_notifications"`
	DefaultLocationId *uuid.UUID `json:"default_location_id" db:"default_location_id"`
	QuickActions json.RawMessage `json:"quick_actions" db:"quick_actions"`
	CustomPreferences json.RawMessage `json:"custom_preferences" db:"custom_preferences"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new user_settings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *UserSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "user_settings", duration, nil)
	}()

	query := `
		INSERT INTO user_settings (
			user_id
			, theme
			, language
			, timezone
			, default_dashboard
			, dashboard_layout
			, items_per_page
			, default_view
			, desktop_notifications
			, sound_notifications
			, default_location_id
			, quick_actions
			, custom_preferences
		) VALUES (
			$1
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
		)
		RETURNING user_id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.UserId,
		entity.Theme,
		entity.Language,
		entity.Timezone,
		entity.DefaultDashboard,
		entity.DashboardLayout,
		entity.ItemsPerPage,
		entity.DefaultView,
		entity.DesktopNotifications,
		entity.SoundNotifications,
		entity.DefaultLocationId,
		entity.QuickActions,
		entity.CustomPreferences,
	)

	
	err := row.Scan(&entity.UserId, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create user_settings", zap.Error(err))
		return fmt.Errorf("failed to create user_settings: %w", err)
	}

	r.logger.Info("created user_settings",
		zap.String("id", entity.UserId.String()),
		
	)

	return nil
}

// GetByID retrieves a user_settings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*UserSettings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_settings", duration, nil)
	}()

	query := `
		SELECT
			user_id
			, theme
			, language
			, timezone
			, default_dashboard
			, dashboard_layout
			, items_per_page
			, default_view
			, desktop_notifications
			, sound_notifications
			, default_location_id
			, quick_actions
			, custom_preferences
			, updated_at
		FROM user_settings
		WHERE user_id = $1
		
	`

	var entity UserSettings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.UserId,
		&entity.Theme,
		&entity.Language,
		&entity.Timezone,
		&entity.DefaultDashboard,
		&entity.DashboardLayout,
		&entity.ItemsPerPage,
		&entity.DefaultView,
		&entity.DesktopNotifications,
		&entity.SoundNotifications,
		&entity.DefaultLocationId,
		&entity.QuickActions,
		&entity.CustomPreferences,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("user_settings not found")
	}

	if err != nil {
		r.logger.Error("failed to get user_settings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get user_settings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of user_settings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*UserSettings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "user_settings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM user_settings
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count user_settings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			user_id
			, theme
			, language
			, timezone
			, default_dashboard
			, dashboard_layout
			, items_per_page
			, default_view
			, desktop_notifications
			, sound_notifications
			, default_location_id
			, quick_actions
			, custom_preferences
			, updated_at
		FROM user_settings
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list user_settings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list user_settings: %w", err)
	}
	defer rows.Close()

	var entities []*UserSettings
	for rows.Next() {
		var entity UserSettings
		err := rows.Scan(
			&entity.UserId,
			&entity.Theme,
			&entity.Language,
			&entity.Timezone,
			&entity.DefaultDashboard,
			&entity.DashboardLayout,
			&entity.ItemsPerPage,
			&entity.DefaultView,
			&entity.DesktopNotifications,
			&entity.SoundNotifications,
			&entity.DefaultLocationId,
			&entity.QuickActions,
			&entity.CustomPreferences,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user_settings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating user_settings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing user_settings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *UserSettings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "user_settings", duration, nil)
	}()

	query := `
		UPDATE user_settings
		SET
			user_id = $1
			, theme = $2
			, language = $3
			, timezone = $4
			, default_dashboard = $5
			, dashboard_layout = $6
			, items_per_page = $7
			, default_view = $8
			, desktop_notifications = $9
			, sound_notifications = $10
			, default_location_id = $11
			, quick_actions = $12
			, custom_preferences = $13
			, updated_at = $14
			, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $15
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.UserId,
		entity.Theme,
		entity.Language,
		entity.Timezone,
		entity.DefaultDashboard,
		entity.DashboardLayout,
		entity.ItemsPerPage,
		entity.DefaultView,
		entity.DesktopNotifications,
		entity.SoundNotifications,
		entity.DefaultLocationId,
		entity.QuickActions,
		entity.CustomPreferences,
		time.Now(),
		entity.UserId,
	)

	if err != nil {
		r.logger.Error("failed to update user_settings", zap.Error(err))
		return fmt.Errorf("failed to update user_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_settings not found or already deleted")
	}

	r.logger.Info("updated user_settings",
		zap.String("id", entity.UserId.String()),
	)

	return nil
}


// Delete permanently deletes a user_settings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "user_settings", duration, nil)
	}()

	query := `DELETE FROM user_settings WHERE user_id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user_settings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete user_settings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user_settings not found")
	}

	r.logger.Info("deleted user_settings", zap.String("id", id.String()))
	return nil
}



