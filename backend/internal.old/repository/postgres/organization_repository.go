package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/organizations"
)

// OrganizationRepository implements organizations.Repository
type OrganizationRepository struct {
	db *DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

// List retrieves organizations with filters
func (r *OrganizationRepository) List(ctx context.Context, filters organizations.OrganizationFilters) ([]organizations.Organization, error) {
	// Build query
	query := `
		SELECT
			id, name, slug, description, email, phone, address, city, state,
			country, postal_code, status, plan, trial_ends_at, subscription_starts_at,
			subscription_ends_at, max_users, max_products, max_locations,
			settings, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM organizations
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argCount := 0

	// Apply filters
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(*filters.Status))
	}

	// Add ordering
	query += " ORDER BY created_at DESC"

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
	var orgList []organizations.Organization
	for rows.Next() {
		var org organizations.Organization
		err := rows.Scan(
			&org.ID, &org.Name, &org.Slug, &org.Description, &org.Email, &org.Phone, &org.Address,
			&org.City, &org.State, &org.Country, &org.PostalCode, &org.Status, &org.Plan,
			&org.TrialEndsAt, &org.SubscriptionStartsAt, &org.SubscriptionEndsAt,
			&org.MaxUsers, &org.MaxProducts, &org.MaxLocations,
			&org.Settings, &org.Metadata, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
			&org.CreatedBy, &org.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		orgList = append(orgList, org)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgList, nil
}

// Count counts organizations matching filters
func (r *OrganizationRepository) Count(ctx context.Context, filters organizations.OrganizationFilters) (int64, error) {
	query := "SELECT COUNT(*) FROM organizations WHERE deleted_at IS NULL"
	args := []interface{}{}
	argCount := 0

	// Apply filters (same as List)
	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, string(*filters.Status))
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *organizations.Organization) error {
	query := `
		INSERT INTO organizations (
			id, name, slug, description, email, phone, address, city, state,
			country, postal_code, status, plan, trial_ends_at, subscription_starts_at,
			subscription_ends_at, max_users, max_products, max_locations,
			settings, metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		org.ID, org.Name, org.Slug, org.Description, org.Email, org.Phone, org.Address,
		org.City, org.State, org.Country, org.PostalCode, org.Status, org.Plan,
		org.TrialEndsAt, org.SubscriptionStartsAt, org.SubscriptionEndsAt,
		org.MaxUsers, org.MaxProducts, org.MaxLocations,
		org.Settings, org.Metadata, org.CreatedAt, org.UpdatedAt, org.CreatedBy,
	)

	return err
}

// Get retrieves an organization by ID
func (r *OrganizationRepository) Get(ctx context.Context, id uuid.UUID) (*organizations.Organization, error) {
	query := `
		SELECT
			id, name, slug, description, email, phone, address, city, state,
			country, postal_code, status, plan, trial_ends_at, subscription_starts_at,
			subscription_ends_at, max_users, max_products, max_locations,
			settings, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM organizations
		WHERE id = $1 AND deleted_at IS NULL
	`

	var org organizations.Organization
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.Slug, &org.Description, &org.Email, &org.Phone, &org.Address,
		&org.City, &org.State, &org.Country, &org.PostalCode, &org.Status, &org.Plan,
		&org.TrialEndsAt, &org.SubscriptionStartsAt, &org.SubscriptionEndsAt,
		&org.MaxUsers, &org.MaxProducts, &org.MaxLocations,
		&org.Settings, &org.Metadata, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
		&org.CreatedBy, &org.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// GetBySlug retrieves an organization by slug
func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*organizations.Organization, error) {
	query := `
		SELECT
			id, name, slug, description, email, phone, address, city, state,
			country, postal_code, status, plan, trial_ends_at, subscription_starts_at,
			subscription_ends_at, max_users, max_products, max_locations,
			settings, metadata, created_at, updated_at, deleted_at, created_by, updated_by
		FROM organizations
		WHERE LOWER(slug) = LOWER($1) AND deleted_at IS NULL
	`

	var org organizations.Organization
	err := r.db.Pool.QueryRow(ctx, query, slug).Scan(
		&org.ID, &org.Name, &org.Slug, &org.Description, &org.Email, &org.Phone, &org.Address,
		&org.City, &org.State, &org.Country, &org.PostalCode, &org.Status, &org.Plan,
		&org.TrialEndsAt, &org.SubscriptionStartsAt, &org.SubscriptionEndsAt,
		&org.MaxUsers, &org.MaxProducts, &org.MaxLocations,
		&org.Settings, &org.Metadata, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt,
		&org.CreatedBy, &org.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// Update updates an existing organization
func (r *OrganizationRepository) Update(ctx context.Context, org *organizations.Organization) error {
	query := `
		UPDATE organizations SET
			name = $3, slug = $4, description = $5, email = $6, phone = $7, address = $8,
			city = $9, state = $10, country = $11, postal_code = $12, status = $13,
			plan = $14, trial_ends_at = $15, subscription_starts_at = $16,
			subscription_ends_at = $17, max_users = $18, max_products = $19, max_locations = $20,
			settings = $21, metadata = $22, updated_at = $23, updated_by = $24
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		org.ID, org.CreatedAt,
		org.Name, org.Slug, org.Description, org.Email, org.Phone, org.Address,
		org.City, org.State, org.Country, org.PostalCode, org.Status,
		org.Plan, org.TrialEndsAt, org.SubscriptionStartsAt,
		org.SubscriptionEndsAt, org.MaxUsers, org.MaxProducts, org.MaxLocations,
		org.Settings, org.Metadata, org.UpdatedAt, org.UpdatedBy,
	)

	return err
}

// Delete soft-deletes an organization
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE organizations
		SET deleted_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, id, time.Now())
	return err
}
