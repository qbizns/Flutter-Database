package notification

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

// Repository handles database operations for Notifications
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Notifications repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Notifications represents a notifications entity
type Notifications struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	NotificationType string `json:"notification_type" db:"notification_type"`
	Category string `json:"category" db:"category"`
	Title string `json:"title" db:"title"`
	Message string `json:"message" db:"message"`
	ActionUrl *string `json:"action_url" db:"action_url"`
	ActionLabel *string `json:"action_label" db:"action_label"`
	Channels *string `json:"channels" db:"channels"`
	IsRead *bool `json:"is_read" db:"is_read"`
	ReadAt *time.Time `json:"read_at" db:"read_at"`
	RelatedEntityType *string `json:"related_entity_type" db:"related_entity_type"`
	RelatedEntityId *uuid.UUID `json:"related_entity_id" db:"related_entity_id"`
	Priority *string `json:"priority" db:"priority"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new notifications record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Notifications) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "notifications", duration, nil)
	}()

	query := `
		INSERT INTO notifications (
			, organization_id
			, user_id
			, notification_type
			, category
			, title
			, message
			, action_url
			, action_label
			, channels
			, is_read
			, read_at
			, related_entity_type
			, related_entity_id
			, priority
			, expires_at
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.NotificationType,
		entity.Category,
		entity.Title,
		entity.Message,
		entity.ActionUrl,
		entity.ActionLabel,
		entity.Channels,
		entity.IsRead,
		entity.ReadAt,
		entity.RelatedEntityType,
		entity.RelatedEntityId,
		entity.Priority,
		entity.ExpiresAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create notifications", zap.Error(err))
		return fmt.Errorf("failed to create notifications: %w", err)
	}

	r.logger.Info("created notifications",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a notifications by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Notifications, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notifications", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, user_id
			, notification_type
			, category
			, title
			, message
			, action_url
			, action_label
			, channels
			, is_read
			, read_at
			, related_entity_type
			, related_entity_id
			, priority
			, expires_at
			, created_at
		FROM notifications
		WHERE id = $1
		
	`

	var entity Notifications
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.UserId,
		&entity.NotificationType,
		&entity.Category,
		&entity.Title,
		&entity.Message,
		&entity.ActionUrl,
		&entity.ActionLabel,
		&entity.Channels,
		&entity.IsRead,
		&entity.ReadAt,
		&entity.RelatedEntityType,
		&entity.RelatedEntityId,
		&entity.Priority,
		&entity.ExpiresAt,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("notifications not found")
	}

	if err != nil {
		r.logger.Error("failed to get notifications", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of notifications records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Notifications, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notifications", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM notifications
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, notification_type
			, category
			, title
			, message
			, action_url
			, action_label
			, channels
			, is_read
			, read_at
			, related_entity_type
			, related_entity_id
			, priority
			, expires_at
			, created_at
		FROM notifications
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list notifications", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var entities []*Notifications
	for rows.Next() {
		var entity Notifications
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.NotificationType,
			&entity.Category,
			&entity.Title,
			&entity.Message,
			&entity.ActionUrl,
			&entity.ActionLabel,
			&entity.Channels,
			&entity.IsRead,
			&entity.ReadAt,
			&entity.RelatedEntityType,
			&entity.RelatedEntityId,
			&entity.Priority,
			&entity.ExpiresAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notifications: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating notifications rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing notifications record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Notifications) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "notifications", duration, nil)
	}()

	query := `
		UPDATE notifications
		SET
			, organization_id = $2
			, user_id = $3
			, notification_type = $4
			, category = $5
			, title = $6
			, message = $7
			, action_url = $8
			, action_label = $9
			, channels = $10
			, is_read = $11
			, read_at = $12
			, related_entity_type = $13
			, related_entity_id = $14
			, priority = $15
			, expires_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.NotificationType,
		entity.Category,
		entity.Title,
		entity.Message,
		entity.ActionUrl,
		entity.ActionLabel,
		entity.Channels,
		entity.IsRead,
		entity.ReadAt,
		entity.RelatedEntityType,
		entity.RelatedEntityId,
		entity.Priority,
		entity.ExpiresAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update notifications", zap.Error(err))
		return fmt.Errorf("failed to update notifications: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notifications not found or already deleted")
	}

	r.logger.Info("updated notifications",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a notifications record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "notifications", duration, nil)
	}()

	query := `DELETE FROM notifications WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete notifications", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete notifications: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notifications not found")
	}

	r.logger.Info("deleted notifications", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves notifications records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Notifications, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "notifications", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM notifications
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, notification_type
			, category
			, title
			, message
			, action_url
			, action_label
			, channels
			, is_read
			, read_at
			, related_entity_type
			, related_entity_id
			, priority
			, expires_at
			, created_at
		FROM notifications
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list notifications by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var entities []*Notifications
	for rows.Next() {
		var entity Notifications
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.NotificationType,
			&entity.Category,
			&entity.Title,
			&entity.Message,
			&entity.ActionUrl,
			&entity.ActionLabel,
			&entity.Channels,
			&entity.IsRead,
			&entity.ReadAt,
			&entity.RelatedEntityType,
			&entity.RelatedEntityId,
			&entity.Priority,
			&entity.ExpiresAt,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notifications: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

