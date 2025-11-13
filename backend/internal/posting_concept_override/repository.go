package posting_concept_override

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

// Repository handles database operations for PostingConceptOverrides
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingConceptOverrides repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingConceptOverrides represents a posting_concept_overrides entity
type PostingConceptOverrides struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ConceptKey string `json:"concept_key" db:"concept_key"`
	Label string `json:"label" db:"label"`
	Description *string `json:"description" db:"description"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new posting_concept_overrides record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingConceptOverrides) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_concept_overrides", duration, nil)
	}()

	query := `
		INSERT INTO posting_concept_overrides (
			, organization_id
			, concept_key
			, label
			, description
			, is_active
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
			, $11
			, $12
			, $13
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ConceptKey,
		entity.Label,
		entity.Description,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_concept_overrides", zap.Error(err))
		return fmt.Errorf("failed to create posting_concept_overrides: %w", err)
	}

	r.logger.Info("created posting_concept_overrides",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a posting_concept_overrides by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingConceptOverrides, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_concept_overrides", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, concept_key
			, label
			, description
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_concept_overrides
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingConceptOverrides
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ConceptKey,
		&entity.Label,
		&entity.Description,
		&entity.IsActive,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_concept_overrides not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_concept_overrides", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_concept_overrides: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_concept_overrides records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingConceptOverrides, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_concept_overrides", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_concept_overrides
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_concept_overrides records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, concept_key
			, label
			, description
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_concept_overrides
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_concept_overrides", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_concept_overrides: %w", err)
	}
	defer rows.Close()

	var entities []*PostingConceptOverrides
	for rows.Next() {
		var entity PostingConceptOverrides
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ConceptKey,
			&entity.Label,
			&entity.Description,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_concept_overrides: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_concept_overrides rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_concept_overrides record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingConceptOverrides) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_concept_overrides", duration, nil)
	}()

	query := `
		UPDATE posting_concept_overrides
		SET
			, organization_id = $2
			, concept_key = $3
			, label = $4
			, description = $5
			, is_active = $6
			, notes = $7
			, metadata = $8
			, updated_at = $10
			, deleted_at = $11
			, created_by = $12
			, updated_by = $13
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $14
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ConceptKey,
		entity.Label,
		entity.Description,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_concept_overrides", zap.Error(err))
		return fmt.Errorf("failed to update posting_concept_overrides: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_concept_overrides not found or already deleted")
	}

	r.logger.Info("updated posting_concept_overrides",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_concept_overrides record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_concept_overrides", duration, nil)
	}()

	query := `
		UPDATE posting_concept_overrides
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_concept_overrides", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_concept_overrides: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_concept_overrides not found or already deleted")
	}

	r.logger.Info("deleted posting_concept_overrides", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves posting_concept_overrides records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PostingConceptOverrides, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_concept_overrides", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_concept_overrides
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_concept_overrides records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, concept_key
			, label
			, description
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_concept_overrides
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_concept_overrides by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_concept_overrides: %w", err)
	}
	defer rows.Close()

	var entities []*PostingConceptOverrides
	for rows.Next() {
		var entity PostingConceptOverrides
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ConceptKey,
			&entity.Label,
			&entity.Description,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_concept_overrides: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

