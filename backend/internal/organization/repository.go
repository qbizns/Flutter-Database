package organization

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

// Repository handles database operations for Organizations
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Organizations repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Organizations represents a organizations entity
type Organizations struct {
	Id *uuid.UUID `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Slug string `json:"slug" db:"slug"`
	Description *string `json:"description" db:"description"`
	Email *string `json:"email" db:"email"`
	Phone *string `json:"phone" db:"phone"`
	Address *string `json:"address" db:"address"`
	City *string `json:"city" db:"city"`
	State *string `json:"state" db:"state"`
	Country *string `json:"country" db:"country"`
	PostalCode *string `json:"postal_code" db:"postal_code"`
	Status string `json:"status" db:"status"`
	Plan *string `json:"plan" db:"plan"`
	TrialEndsAt *time.Time `json:"trial_ends_at" db:"trial_ends_at"`
	SubscriptionStartsAt *time.Time `json:"subscription_starts_at" db:"subscription_starts_at"`
	SubscriptionEndsAt *time.Time `json:"subscription_ends_at" db:"subscription_ends_at"`
	MaxUsers *int64 `json:"max_users" db:"max_users"`
	MaxProducts *int64 `json:"max_products" db:"max_products"`
	MaxLocations *int64 `json:"max_locations" db:"max_locations"`
	Settings json.RawMessage `json:"settings" db:"settings"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new organizations record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Organizations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "organizations", duration, nil)
	}()

	query := `
		INSERT INTO organizations (
			, name
			, slug
			, description
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, status
			, plan
			, trial_ends_at
			, subscription_starts_at
			, subscription_ends_at
			, max_users
			, max_products
			, max_locations
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
			, $17
			, $18
			, $19
			, $20
			, $21
			, $24
			, $25
			, $26
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.Email,
		entity.Phone,
		entity.Address,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.Status,
		entity.Plan,
		entity.TrialEndsAt,
		entity.SubscriptionStartsAt,
		entity.SubscriptionEndsAt,
		entity.MaxUsers,
		entity.MaxProducts,
		entity.MaxLocations,
		entity.Settings,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create organizations", zap.Error(err))
		return fmt.Errorf("failed to create organizations: %w", err)
	}

	r.logger.Info("created organizations",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a organizations by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Organizations, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organizations", duration, nil)
	}()

	query := `
		SELECT
			id
			, name
			, slug
			, description
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, status
			, plan
			, trial_ends_at
			, subscription_starts_at
			, subscription_ends_at
			, max_users
			, max_products
			, max_locations
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM organizations
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Organizations
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.Name,
		&entity.Slug,
		&entity.Description,
		&entity.Email,
		&entity.Phone,
		&entity.Address,
		&entity.City,
		&entity.State,
		&entity.Country,
		&entity.PostalCode,
		&entity.Status,
		&entity.Plan,
		&entity.TrialEndsAt,
		&entity.SubscriptionStartsAt,
		&entity.SubscriptionEndsAt,
		&entity.MaxUsers,
		&entity.MaxProducts,
		&entity.MaxLocations,
		&entity.Settings,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("organizations not found")
	}

	if err != nil {
		r.logger.Error("failed to get organizations", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of organizations records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Organizations, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "organizations", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM organizations
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count organizations records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, name
			, slug
			, description
			, email
			, phone
			, address
			, city
			, state
			, country
			, postal_code
			, status
			, plan
			, trial_ends_at
			, subscription_starts_at
			, subscription_ends_at
			, max_users
			, max_products
			, max_locations
			, settings
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM organizations
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list organizations", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var entities []*Organizations
	for rows.Next() {
		var entity Organizations
		err := rows.Scan(
			&entity.Id,
			&entity.Name,
			&entity.Slug,
			&entity.Description,
			&entity.Email,
			&entity.Phone,
			&entity.Address,
			&entity.City,
			&entity.State,
			&entity.Country,
			&entity.PostalCode,
			&entity.Status,
			&entity.Plan,
			&entity.TrialEndsAt,
			&entity.SubscriptionStartsAt,
			&entity.SubscriptionEndsAt,
			&entity.MaxUsers,
			&entity.MaxProducts,
			&entity.MaxLocations,
			&entity.Settings,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan organizations: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating organizations rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing organizations record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Organizations) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "organizations", duration, nil)
	}()

	query := `
		UPDATE organizations
		SET
			, name = $2
			, slug = $3
			, description = $4
			, email = $5
			, phone = $6
			, address = $7
			, city = $8
			, state = $9
			, country = $10
			, postal_code = $11
			, status = $12
			, plan = $13
			, trial_ends_at = $14
			, subscription_starts_at = $15
			, subscription_ends_at = $16
			, max_users = $17
			, max_products = $18
			, max_locations = $19
			, settings = $20
			, metadata = $21
			, updated_at = $23
			, deleted_at = $24
			, created_by = $25
			, updated_by = $26
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $27
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.Name,
		entity.Slug,
		entity.Description,
		entity.Email,
		entity.Phone,
		entity.Address,
		entity.City,
		entity.State,
		entity.Country,
		entity.PostalCode,
		entity.Status,
		entity.Plan,
		entity.TrialEndsAt,
		entity.SubscriptionStartsAt,
		entity.SubscriptionEndsAt,
		entity.MaxUsers,
		entity.MaxProducts,
		entity.MaxLocations,
		entity.Settings,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update organizations", zap.Error(err))
		return fmt.Errorf("failed to update organizations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organizations not found or already deleted")
	}

	r.logger.Info("updated organizations",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a organizations record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "organizations", duration, nil)
	}()

	query := `
		UPDATE organizations
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete organizations", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete organizations: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("organizations not found or already deleted")
	}

	r.logger.Info("deleted organizations", zap.String("id", id.String()))
	return nil
}



