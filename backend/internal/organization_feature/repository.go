package organization_feature

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

// Repository handles database operations for OrganizationFeatures
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new OrganizationFeatures repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// OrganizationFeatures represents a organization_features entity
type OrganizationFeatures struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	FeatureKey string `json:"feature_key" db:"feature_key"`
	IsEnabled *bool `json:"is_enabled" db:"is_enabled"`
	IsAvailable *bool `json:"is_available" db:"is_available"`
	Configuration json.RawMessage `json:"configuration" db:"configuration"`
	Limits json.RawMessage `json:"limits" db:"limits"`
	EnabledAt *time.Time `json:"enabled_at" db:"enabled_at"`
	DisabledAt *time.Time `json:"disabled_at" db:"disabled_at"`
	ExpiresAt *time.Time `json:"expires_at" db:"expires_at"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new organization_features record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *OrganizationFeatures) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "organization_features", duration, nil)
	}()

	query := `
		INSERT INTO organization_features (
			, organization_id
			, feature_key
			, is_enabled
			, is_available
			, configuration
			, limits
			, enabled_at
			, disabled_at
			, expires_at
			, notes
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
			, $15
			, $16
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.FeatureKey,
		entity.IsEnabled,
		entity.IsAvailable,
		entity.Configuration,
		entity.Limits,
		entity.EnabledAt,
		entity.DisabledAt,
		entity.ExpiresAt,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create organization_features", zap.Error(err))
		return fmt.Errorf("failed to create organization_features: %w", err)
	}

	r.logger.Info("created organization_features",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a organization_features by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*OrganizationFeatures, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_features", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, feature_key
			, is_enabled
			, is_available
			, configuration
			, limits
			, enabled_at
			, disabled_at
			, expires_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM organization_features
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity OrganizationFeatures
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.FeatureKey,
		&entity.IsEnabled,
		&entity.IsAvailable,
		&entity.Configuration,
		&entity.Limits,
		&entity.EnabledAt,
		&entity.DisabledAt,
		&entity.ExpiresAt,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("organization_features not found")
	}

	if err != nil {
		r.logger.Error("failed to get organization_features", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get organization_features: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of organization_features records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*OrganizationFeatures, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_features", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM organization_features
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization_features records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, feature_key
			, is_enabled
			, is_available
			, configuration
			, limits
			, enabled_at
			, disabled_at
			, expires_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM organization_features
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list organization_features", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list organization_features: %w", err)
	}
	defer rows.Close()

	var entities []*OrganizationFeatures
	for rows.Next() {
		var entity OrganizationFeatures
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FeatureKey,
			&entity.IsEnabled,
			&entity.IsAvailable,
			&entity.Configuration,
			&entity.Limits,
			&entity.EnabledAt,
			&entity.DisabledAt,
			&entity.ExpiresAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization_features: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating organization_features rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing organization_features record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *OrganizationFeatures) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "organization_features", duration, nil)
	}()

	query := `
		UPDATE organization_features
		SET
			, organization_id = $2
			, feature_key = $3
			, is_enabled = $4
			, is_available = $5
			, configuration = $6
			, limits = $7
			, enabled_at = $8
			, disabled_at = $9
			, expires_at = $10
			, notes = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, created_by = $16
			, updated_by = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.FeatureKey,
		entity.IsEnabled,
		entity.IsAvailable,
		entity.Configuration,
		entity.Limits,
		entity.EnabledAt,
		entity.DisabledAt,
		entity.ExpiresAt,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update organization_features", zap.Error(err))
		return fmt.Errorf("failed to update organization_features: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organization_features not found or already deleted")
	}

	r.logger.Info("updated organization_features",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a organization_features record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "organization_features", duration, nil)
	}()

	query := `
		UPDATE organization_features
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete organization_features", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete organization_features: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organization_features not found or already deleted")
	}

	r.logger.Info("deleted organization_features", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves organization_features records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*OrganizationFeatures, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organization_features", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM organization_features
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organization_features records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, feature_key
			, is_enabled
			, is_available
			, configuration
			, limits
			, enabled_at
			, disabled_at
			, expires_at
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM organization_features
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list organization_features by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list organization_features: %w", err)
	}
	defer rows.Close()

	var entities []*OrganizationFeatures
	for rows.Next() {
		var entity OrganizationFeatures
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FeatureKey,
			&entity.IsEnabled,
			&entity.IsAvailable,
			&entity.Configuration,
			&entity.Limits,
			&entity.EnabledAt,
			&entity.DisabledAt,
			&entity.ExpiresAt,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organization_features: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

