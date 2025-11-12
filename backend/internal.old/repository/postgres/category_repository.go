package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/categories"
)

// CategoryRepository implements categories.Repository
type CategoryRepository struct {
	db *DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// List retrieves categories with filters
func (r *CategoryRepository) List(ctx context.Context, orgID uuid.UUID, filters categories.CategoryFilters) ([]categories.Category, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	// Build query
	query := `
		SELECT
			id, organization_id, name, slug, description,
			parent_id, level, path, image_url, icon, color,
			sort_order, is_active, metadata,
			created_at, updated_at, created_by, updated_by
		FROM categories
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.ParentID != nil {
		argCount++
		query += fmt.Sprintf(" AND parent_id = $%d", argCount)
		args = append(args, *filters.ParentID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	// Add ordering by sort_order then name
	query += " ORDER BY sort_order ASC, name ASC"

	// Add pagination
	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan results
	var categoryList []categories.Category
	for rows.Next() {
		var c categories.Category
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.Name, &c.Slug, &c.Description,
			&c.ParentID, &c.Level, &c.Path, &c.ImageURL, &c.Icon, &c.Color,
			&c.SortOrder, &c.IsActive, &c.Metadata,
			&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		categoryList = append(categoryList, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categoryList, nil
}

// Count counts categories matching filters
func (r *CategoryRepository) Count(ctx context.Context, orgID uuid.UUID, filters categories.CategoryFilters) (int64, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM categories WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	// Apply filters (same as List)
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.ParentID != nil {
		argCount++
		query += fmt.Sprintf(" AND parent_id = $%d", argCount)
		args = append(args, *filters.ParentID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *categories.Category) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, category.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO categories (
			id, organization_id, name, slug, description,
			parent_id, level, path, image_url, icon, color,
			sort_order, is_active, metadata,
			created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		category.ID, category.OrganizationID, category.Name, category.Slug, category.Description,
		category.ParentID, category.Level, category.Path, category.ImageURL, category.Icon, category.Color,
		category.SortOrder, category.IsActive, category.Metadata,
		category.CreatedAt, category.UpdatedAt, category.CreatedBy,
	)

	return err
}

// Get retrieves a category by ID
func (r *CategoryRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*categories.Category, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, name, slug, description,
			parent_id, level, path, image_url, icon, color,
			sort_order, is_active, metadata,
			created_at, updated_at, created_by, updated_by
		FROM categories
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var c categories.Category
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&c.ID, &c.OrganizationID, &c.Name, &c.Slug, &c.Description,
		&c.ParentID, &c.Level, &c.Path, &c.ImageURL, &c.Icon, &c.Color,
		&c.SortOrder, &c.IsActive, &c.Metadata,
		&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(ctx context.Context, orgID uuid.UUID, slug string) (*categories.Category, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, name, slug, description,
			parent_id, level, path, image_url, icon, color,
			sort_order, is_active, metadata,
			created_at, updated_at, created_by, updated_by
		FROM categories
		WHERE organization_id = $1 AND slug = $2 AND deleted_at IS NULL
	`

	var c categories.Category
	err := r.db.Pool.QueryRow(ctx, query, orgID, slug).Scan(
		&c.ID, &c.OrganizationID, &c.Name, &c.Slug, &c.Description,
		&c.ParentID, &c.Level, &c.Path, &c.ImageURL, &c.Icon, &c.Color,
		&c.SortOrder, &c.IsActive, &c.Metadata,
		&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetChildren retrieves all direct child categories of a parent
func (r *CategoryRepository) GetChildren(ctx context.Context, orgID uuid.UUID, parentID uuid.UUID) ([]categories.Category, error) {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT
			id, organization_id, name, slug, description,
			parent_id, level, path, image_url, icon, color,
			sort_order, is_active, metadata,
			created_at, updated_at, created_by, updated_by
		FROM categories
		WHERE organization_id = $1 AND parent_id = $2 AND deleted_at IS NULL
		ORDER BY sort_order ASC, name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []categories.Category
	for rows.Next() {
		var c categories.Category
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.Name, &c.Slug, &c.Description,
			&c.ParentID, &c.Level, &c.Path, &c.ImageURL, &c.Icon, &c.Color,
			&c.SortOrder, &c.IsActive, &c.Metadata,
			&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return children, nil
}

// Update updates an existing category
func (r *CategoryRepository) Update(ctx context.Context, category *categories.Category) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, category.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE categories SET
			name = $3, slug = $4, description = $5,
			parent_id = $6, level = $7, path = $8,
			image_url = $9, icon = $10, color = $11,
			sort_order = $12, is_active = $13, metadata = $14,
			updated_at = $15, updated_by = $16
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		category.OrganizationID, category.ID,
		category.Name, category.Slug, category.Description,
		category.ParentID, category.Level, category.Path,
		category.ImageURL, category.Icon, category.Color,
		category.SortOrder, category.IsActive, category.Metadata,
		category.UpdatedAt, category.UpdatedBy,
	)

	return err
}

// Delete soft-deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Set organization context for RLS
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE categories
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}
