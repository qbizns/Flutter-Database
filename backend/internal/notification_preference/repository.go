package notification_preference

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

// Repository handles database operations for NotificationPreferences
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new NotificationPreferences repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// NotificationPreferences represents a notification_preferences entity
type NotificationPreferences struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Category string `json:"category" db:"category"`
	InAppEnabled *bool `json:"in_app_enabled" db:"in_app_enabled"`
	EmailEnabled *bool `json:"email_enabled" db:"email_enabled"`
	SmsEnabled *bool `json:"sms_enabled" db:"sms_enabled"`
	PushEnabled *bool `json:"push_enabled" db:"push_enabled"`
	Frequency *string `json:"frequency" db:"frequency"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new notification_preferences record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *NotificationPreferences) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "notification_preferences", duration, nil)
	}()

	query := `
		INSERT INTO notification_preferences (
			, organization_id
			, user_id
			, category
			, in_app_enabled
			, email_enabled
			, sms_enabled
			, push_enabled
			, frequency
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.Category,
		entity.InAppEnabled,
		entity.EmailEnabled,
		entity.SmsEnabled,
		entity.PushEnabled,
		entity.Frequency,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create notification_preferences", zap.Error(err))
		return fmt.Errorf("failed to create notification_preferences: %w", err)
	}

	r.logger.Info("created notification_preferences",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a notification_preferences by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*NotificationPreferences, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notification_preferences", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, user_id
			, category
			, in_app_enabled
			, email_enabled
			, sms_enabled
			, push_enabled
			, frequency
			, created_at
			, updated_at
		FROM notification_preferences
		WHERE id = $1
		
	`

	var entity NotificationPreferences
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.UserId,
		&entity.Category,
		&entity.InAppEnabled,
		&entity.EmailEnabled,
		&entity.SmsEnabled,
		&entity.PushEnabled,
		&entity.Frequency,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("notification_preferences not found")
	}

	if err != nil {
		r.logger.Error("failed to get notification_preferences", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get notification_preferences: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of notification_preferences records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*NotificationPreferences, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notification_preferences", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM notification_preferences
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notification_preferences records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, category
			, in_app_enabled
			, email_enabled
			, sms_enabled
			, push_enabled
			, frequency
			, created_at
			, updated_at
		FROM notification_preferences
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list notification_preferences", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list notification_preferences: %w", err)
	}
	defer rows.Close()

	var entities []*NotificationPreferences
	for rows.Next() {
		var entity NotificationPreferences
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.Category,
			&entity.InAppEnabled,
			&entity.EmailEnabled,
			&entity.SmsEnabled,
			&entity.PushEnabled,
			&entity.Frequency,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification_preferences: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating notification_preferences rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing notification_preferences record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *NotificationPreferences) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "notification_preferences", duration, nil)
	}()

	query := `
		UPDATE notification_preferences
		SET
			, organization_id = $2
			, user_id = $3
			, category = $4
			, in_app_enabled = $5
			, email_enabled = $6
			, sms_enabled = $7
			, push_enabled = $8
			, frequency = $9
			, updated_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.Category,
		entity.InAppEnabled,
		entity.EmailEnabled,
		entity.SmsEnabled,
		entity.PushEnabled,
		entity.Frequency,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update notification_preferences", zap.Error(err))
		return fmt.Errorf("failed to update notification_preferences: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification_preferences not found or already deleted")
	}

	r.logger.Info("updated notification_preferences",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a notification_preferences record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "notification_preferences", duration, nil)
	}()

	query := `DELETE FROM notification_preferences WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete notification_preferences", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete notification_preferences: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification_preferences not found")
	}

	r.logger.Info("deleted notification_preferences", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves notification_preferences records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*NotificationPreferences, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notification_preferences", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM notification_preferences
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notification_preferences records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, category
			, in_app_enabled
			, email_enabled
			, sms_enabled
			, push_enabled
			, frequency
			, created_at
			, updated_at
		FROM notification_preferences
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list notification_preferences by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list notification_preferences: %w", err)
	}
	defer rows.Close()

	var entities []*NotificationPreferences
	for rows.Next() {
		var entity NotificationPreferences
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.Category,
			&entity.InAppEnabled,
			&entity.EmailEnabled,
			&entity.SmsEnabled,
			&entity.PushEnabled,
			&entity.Frequency,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification_preferences: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

