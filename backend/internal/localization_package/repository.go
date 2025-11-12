package localization_package

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

// Repository handles database operations for LocalizationPackages
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new LocalizationPackages repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// LocalizationPackages represents a localization_packages entity
type LocalizationPackages struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PackageCode string `json:"package_code" db:"package_code"`
	PackageName string `json:"package_name" db:"package_name"`
	CountryCode *string `json:"country_code" db:"country_code"`
	Region *string `json:"region" db:"region"`
	Description *string `json:"description" db:"description"`
	Version *string `json:"version" db:"version"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new localization_packages record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *LocalizationPackages) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "localization_packages", duration, nil)
	}()

	query := `
		INSERT INTO localization_packages (
			, package_code
			, package_name
			, country_code
			, region
			, description
			, version
			, is_active
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PackageCode,
		entity.PackageName,
		entity.CountryCode,
		entity.Region,
		entity.Description,
		entity.Version,
		entity.IsActive,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create localization_packages", zap.Error(err))
		return fmt.Errorf("failed to create localization_packages: %w", err)
	}

	r.logger.Info("created localization_packages",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a localization_packages by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*LocalizationPackages, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "localization_packages", duration, nil)
	}()

	query := `
		SELECT
			id
			, package_code
			, package_name
			, country_code
			, region
			, description
			, version
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM localization_packages
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity LocalizationPackages
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PackageCode,
		&entity.PackageName,
		&entity.CountryCode,
		&entity.Region,
		&entity.Description,
		&entity.Version,
		&entity.IsActive,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("localization_packages not found")
	}

	if err != nil {
		r.logger.Error("failed to get localization_packages", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get localization_packages: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of localization_packages records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*LocalizationPackages, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "localization_packages", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM localization_packages
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count localization_packages records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, package_code
			, package_name
			, country_code
			, region
			, description
			, version
			, is_active
			, created_at
			, updated_at
			, deleted_at
		FROM localization_packages
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list localization_packages", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list localization_packages: %w", err)
	}
	defer rows.Close()

	var entities []*LocalizationPackages
	for rows.Next() {
		var entity LocalizationPackages
		err := rows.Scan(
			&entity.Id,
			&entity.PackageCode,
			&entity.PackageName,
			&entity.CountryCode,
			&entity.Region,
			&entity.Description,
			&entity.Version,
			&entity.IsActive,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan localization_packages: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating localization_packages rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing localization_packages record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *LocalizationPackages) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "localization_packages", duration, nil)
	}()

	query := `
		UPDATE localization_packages
		SET
			, package_code = $2
			, package_name = $3
			, country_code = $4
			, region = $5
			, description = $6
			, version = $7
			, is_active = $8
			, updated_at = $10
			, deleted_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PackageCode,
		entity.PackageName,
		entity.CountryCode,
		entity.Region,
		entity.Description,
		entity.Version,
		entity.IsActive,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update localization_packages", zap.Error(err))
		return fmt.Errorf("failed to update localization_packages: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("localization_packages not found or already deleted")
	}

	r.logger.Info("updated localization_packages",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a localization_packages record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "localization_packages", duration, nil)
	}()

	query := `
		UPDATE localization_packages
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete localization_packages", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete localization_packages: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("localization_packages not found or already deleted")
	}

	r.logger.Info("deleted localization_packages", zap.String("id", id.String()))
	return nil
}



