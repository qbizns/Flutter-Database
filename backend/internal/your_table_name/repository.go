package your_table_name

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

// Repository handles database operations for YourTableName
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new YourTableName repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// YourTableName represents a your_table_name entity
type YourTableName struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new your_table_name record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *YourTableName) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "your_table_name", duration, nil)
	}()

	query := `
		INSERT INTO your_table_name (
			, organization_id
			, name
			, description
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
			, $9
			, $10
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Description,
		entity.Settings,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create your_table_name", zap.Error(err))
		return fmt.Errorf("failed to create your_table_name: %w", err)
	}

	r.logger.Info("created your_table_name",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a your_table_name by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*YourTableName, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "your_table_name", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, name
			, description
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM your_table_name
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity YourTableName
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Name,
		&entity.Description,
		&entity.Settings,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("your_table_name not found")
	}

	if err != nil {
		r.logger.Error("failed to get your_table_name", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get your_table_name: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of your_table_name records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*YourTableName, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "your_table_name", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM your_table_name
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count your_table_name records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, description
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM your_table_name
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list your_table_name", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list your_table_name: %w", err)
	}
	defer rows.Close()

	var entities []*YourTableName
	for rows.Next() {
		var entity YourTableName
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Description,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan your_table_name: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating your_table_name rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing your_table_name record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *YourTableName) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "your_table_name", duration, nil)
	}()

	query := `
		UPDATE your_table_name
		SET
			, organization_id = $2
			, name = $3
			, description = $4
			, settings = $5
			, metadata = $6
			, updated_at = $8
			, deleted_at = $9
			, created_by = $10
			, updated_by = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Description,
		entity.Settings,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update your_table_name", zap.Error(err))
		return fmt.Errorf("failed to update your_table_name: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("your_table_name not found or already deleted")
	}

	r.logger.Info("updated your_table_name",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a your_table_name record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "your_table_name", duration, nil)
	}()

	query := `
		UPDATE your_table_name
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete your_table_name", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete your_table_name: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("your_table_name not found or already deleted")
	}

	r.logger.Info("deleted your_table_name", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves your_table_name records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*YourTableName, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "your_table_name", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM your_table_name
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count your_table_name records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, description
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM your_table_name
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list your_table_name by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list your_table_name: %w", err)
	}
	defer rows.Close()

	var entities []*YourTableName
	for rows.Next() {
		var entity YourTableName
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Description,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan your_table_name: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

