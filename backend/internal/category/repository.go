package category

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for Categories
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Categories repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Categories represents a categories entity
type Categories struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Description *string `json:"description" db:"description"`
	ParentId *uuid.UUID `json:"parent_id" db:"parent_id"`
	Level *int64 `json:"level" db:"level"`
	Path *string `json:"path" db:"path"`
	ImageUrl *string `json:"image_url" db:"image_url"`
	Icon *string `json:"icon" db:"icon"`
	Color *string `json:"color" db:"color"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new categories record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Categories) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "categories", duration, nil)
	}()

	query := `
		INSERT INTO categories (
			, organization_id
			, name
			, slug
			, description
			, parent_id
			, level
			, path
			, image_url
			, icon
			, color
			, sort_order
			, is_active
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
			, $17
			, $18
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.ParentId,
		entity.Level,
		entity.Path,
		entity.ImageUrl,
		entity.Icon,
		entity.Color,
		entity.SortOrder,
		entity.IsActive,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create categories", zap.Error(err))
		return fmt.Errorf("failed to create categories: %w", err)
	}

	r.logger.Info("created categories",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a categories by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Categories, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "categories", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, parent_id
			, level
			, path
			, image_url
			, icon
			, color
			, sort_order
			, is_active
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM categories
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Categories
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Name,
		&entity.Slug,
		&entity.Description,
		&entity.ParentId,
		&entity.Level,
		&entity.Path,
		&entity.ImageUrl,
		&entity.Icon,
		&entity.Color,
		&entity.SortOrder,
		&entity.IsActive,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("categories not found")
	}

	if err != nil {
		r.logger.Error("failed to get categories", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of categories records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Categories, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "categories", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM categories
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, parent_id
			, level
			, path
			, image_url
			, icon
			, color
			, sort_order
			, is_active
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM categories
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list categories", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var entities []*Categories
	for rows.Next() {
		var entity Categories
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.ParentId,
			&entity.Level,
			&entity.Path,
			&entity.ImageUrl,
			&entity.Icon,
			&entity.Color,
			&entity.SortOrder,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan categories: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating categories rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing categories record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Categories) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "categories", duration, nil)
	}()

	query := `
		UPDATE categories
		SET
			, organization_id = $2
			, name = $3
			, slug = $4
			, description = $5
			, parent_id = $6
			, level = $7
			, path = $8
			, image_url = $9
			, icon = $10
			, color = $11
			, sort_order = $12
			, is_active = $13
			, metadata = $14
			, updated_at = $16
			, deleted_at = $17
			, created_by = $18
			, updated_by = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.ParentId,
		entity.Level,
		entity.Path,
		entity.ImageUrl,
		entity.Icon,
		entity.Color,
		entity.SortOrder,
		entity.IsActive,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update categories", zap.Error(err))
		return fmt.Errorf("failed to update categories: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("categories not found or already deleted")
	}

	r.logger.Info("updated categories",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a categories record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "categories", duration, nil)
	}()

	query := `
		UPDATE categories
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete categories", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete categories: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("categories not found or already deleted")
	}

	r.logger.Info("deleted categories", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves categories records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Categories, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "categories", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM categories
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, name
			, slug
			, description
			, parent_id
			, level
			, path
			, image_url
			, icon
			, color
			, sort_order
			, is_active
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM categories
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list categories by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var entities []*Categories
	for rows.Next() {
		var entity Categories
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.ParentId,
			&entity.Level,
			&entity.Path,
			&entity.ImageUrl,
			&entity.Icon,
			&entity.Color,
			&entity.SortOrder,
			&entity.IsActive,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan categories: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}


// GetProductsCountByCategory retrieves the count of products per category
func (r *Repository) GetProductsCountByCategory(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) (map[string]int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "categories", duration, nil)
	}()

	query := `
		SELECT c.id, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id AND p.deleted_at IS NULL
		WHERE c.organization_id = $1 AND c.deleted_at IS NULL
		GROUP BY c.id`

	rows, err := tx.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products count by category: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var categoryID uuid.UUID
		var count int
		if err := rows.Scan(&categoryID, &count); err != nil {
			return nil, fmt.Errorf("failed to scan count: %w", err)
		}
		counts[categoryID.String()] = count
	}

	r.logger.Debug("retrieved products count by category",
		zap.String("organization_id", orgID.String()),
		zap.Int("categories", len(counts)),
	)

	return counts, nil
}
