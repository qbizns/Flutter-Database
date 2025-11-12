package audit_log

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

// Repository handles database operations for AuditLogs
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AuditLogs repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AuditLogs represents a audit_logs entity
type AuditLogs struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	Action string `json:"action" db:"action"`
	ResourceType string `json:"resource_type" db:"resource_type"`
	ResourceId *uuid.UUID `json:"resource_id" db:"resource_id"`
	OldValues json.RawMessage `json:"old_values" db:"old_values"`
	NewValues json.RawMessage `json:"new_values" db:"new_values"`
	Changes json.RawMessage `json:"changes" db:"changes"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	UserAgent *string `json:"user_agent" db:"user_agent"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Create inserts a new audit_logs record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AuditLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "audit_logs", duration, nil)
	}()

	query := `
		INSERT INTO audit_logs (
			, organization_id
			, user_id
			, action
			, resource_type
			, resource_id
			, old_values
			, new_values
			, changes
			, ip_address
			, user_agent
			, metadata
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.Action,
		entity.ResourceType,
		entity.ResourceId,
		entity.OldValues,
		entity.NewValues,
		entity.Changes,
		entity.IpAddress,
		entity.UserAgent,
		entity.Metadata,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create audit_logs", zap.Error(err))
		return fmt.Errorf("failed to create audit_logs: %w", err)
	}

	r.logger.Info("created audit_logs",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a audit_logs by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AuditLogs, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "audit_logs", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, user_id
			, action
			, resource_type
			, resource_id
			, old_values
			, new_values
			, changes
			, ip_address
			, user_agent
			, metadata
			, created_at
		FROM audit_logs
		WHERE id = $1
		
	`

	var entity AuditLogs
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.UserId,
		&entity.Action,
		&entity.ResourceType,
		&entity.ResourceId,
		&entity.OldValues,
		&entity.NewValues,
		&entity.Changes,
		&entity.IpAddress,
		&entity.UserAgent,
		&entity.Metadata,
		&entity.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("audit_logs not found")
	}

	if err != nil {
		r.logger.Error("failed to get audit_logs", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get audit_logs: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of audit_logs records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AuditLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "audit_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM audit_logs
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, action
			, resource_type
			, resource_id
			, old_values
			, new_values
			, changes
			, ip_address
			, user_agent
			, metadata
			, created_at
		FROM audit_logs
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list audit_logs", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list audit_logs: %w", err)
	}
	defer rows.Close()

	var entities []*AuditLogs
	for rows.Next() {
		var entity AuditLogs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.Action,
			&entity.ResourceType,
			&entity.ResourceId,
			&entity.OldValues,
			&entity.NewValues,
			&entity.Changes,
			&entity.IpAddress,
			&entity.UserAgent,
			&entity.Metadata,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating audit_logs rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing audit_logs record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AuditLogs) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "audit_logs", duration, nil)
	}()

	query := `
		UPDATE audit_logs
		SET
			, organization_id = $2
			, user_id = $3
			, action = $4
			, resource_type = $5
			, resource_id = $6
			, old_values = $7
			, new_values = $8
			, changes = $9
			, ip_address = $10
			, user_agent = $11
			, metadata = $12
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.UserId,
		entity.Action,
		entity.ResourceType,
		entity.ResourceId,
		entity.OldValues,
		entity.NewValues,
		entity.Changes,
		entity.IpAddress,
		entity.UserAgent,
		entity.Metadata,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update audit_logs", zap.Error(err))
		return fmt.Errorf("failed to update audit_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("audit_logs not found or already deleted")
	}

	r.logger.Info("updated audit_logs",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a audit_logs record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "audit_logs", duration, nil)
	}()

	query := `DELETE FROM audit_logs WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete audit_logs", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete audit_logs: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("audit_logs not found")
	}

	r.logger.Info("deleted audit_logs", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves audit_logs records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*AuditLogs, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "audit_logs", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM audit_logs
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit_logs records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, user_id
			, action
			, resource_type
			, resource_id
			, old_values
			, new_values
			, changes
			, ip_address
			, user_agent
			, metadata
			, created_at
		FROM audit_logs
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list audit_logs by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list audit_logs: %w", err)
	}
	defer rows.Close()

	var entities []*AuditLogs
	for rows.Next() {
		var entity AuditLogs
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.UserId,
			&entity.Action,
			&entity.ResourceType,
			&entity.ResourceId,
			&entity.OldValues,
			&entity.NewValues,
			&entity.Changes,
			&entity.IpAddress,
			&entity.UserAgent,
			&entity.Metadata,
			&entity.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit_logs: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

