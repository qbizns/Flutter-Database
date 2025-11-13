package posting_profile

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

// Repository handles database operations for PostingProfiles
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingProfiles repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingProfiles represents a posting_profiles entity
type PostingProfiles struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	IsDefault bool `json:"is_default" db:"is_default"`
	IsActive bool `json:"is_active" db:"is_active"`
	DefaultFiscalYearId *uuid.UUID `json:"default_fiscal_year_id" db:"default_fiscal_year_id"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new posting_profiles record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingProfiles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_profiles", duration, nil)
	}()

	query := `
		INSERT INTO posting_profiles (
			, organization_id
			, code
			, name
			, description
			, is_default
			, is_active
			, default_fiscal_year_id
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
			, $13
			, $14
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.IsDefault,
		entity.IsActive,
		entity.DefaultFiscalYearId,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_profiles", zap.Error(err))
		return fmt.Errorf("failed to create posting_profiles: %w", err)
	}

	r.logger.Info("created posting_profiles",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a posting_profiles by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingProfiles, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_profiles", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, is_default
			, is_active
			, default_fiscal_year_id
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_profiles
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingProfiles
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.Code,
		&entity.Name,
		&entity.Description,
		&entity.IsDefault,
		&entity.IsActive,
		&entity.DefaultFiscalYearId,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_profiles not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_profiles", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_profiles: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_profiles records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingProfiles, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_profiles", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_profiles
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_profiles records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, is_default
			, is_active
			, default_fiscal_year_id
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_profiles
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_profiles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_profiles: %w", err)
	}
	defer rows.Close()

	var entities []*PostingProfiles
	for rows.Next() {
		var entity PostingProfiles
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.DefaultFiscalYearId,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_profiles: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_profiles rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_profiles record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingProfiles) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_profiles", duration, nil)
	}()

	query := `
		UPDATE posting_profiles
		SET
			, organization_id = $2
			, code = $3
			, name = $4
			, description = $5
			, is_default = $6
			, is_active = $7
			, default_fiscal_year_id = $8
			, notes = $9
			, metadata = $10
			, updated_at = $12
			, deleted_at = $13
			, created_by = $14
			, updated_by = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.IsDefault,
		entity.IsActive,
		entity.DefaultFiscalYearId,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_profiles", zap.Error(err))
		return fmt.Errorf("failed to update posting_profiles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_profiles not found or already deleted")
	}

	r.logger.Info("updated posting_profiles",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_profiles record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_profiles", duration, nil)
	}()

	query := `
		UPDATE posting_profiles
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_profiles", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_profiles: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_profiles not found or already deleted")
	}

	r.logger.Info("deleted posting_profiles", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves posting_profiles records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PostingProfiles, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_profiles", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_profiles
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_profiles records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, code
			, name
			, description
			, is_default
			, is_active
			, default_fiscal_year_id
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_profiles
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_profiles by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_profiles: %w", err)
	}
	defer rows.Close()

	var entities []*PostingProfiles
	for rows.Next() {
		var entity PostingProfiles
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.DefaultFiscalYearId,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_profiles: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

