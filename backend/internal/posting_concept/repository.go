package posting_concept

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

// Repository handles database operations for PostingConcepts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingConcepts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingConcepts represents a posting_concepts entity
type PostingConcepts struct {
	ConceptKey *string `json:"concept_key" db:"concept_key"`
	DefaultLabel string `json:"default_label" db:"default_label"`
	DefaultDescription *string `json:"default_description" db:"default_description"`
	ExpectedAccountTypeId *uuid.UUID `json:"expected_account_type_id" db:"expected_account_type_id"`
	NormalSide *string `json:"normal_side" db:"normal_side"`
	ExampleCode *string `json:"example_code" db:"example_code"`
	ExampleAccountName *string `json:"example_account_name" db:"example_account_name"`
	IsSystem bool `json:"is_system" db:"is_system"`
	ConceptCategory *string `json:"concept_category" db:"concept_category"`
	SortOrder *int64 `json:"sort_order" db:"sort_order"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new posting_concepts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingConcepts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_concepts", duration, nil)
	}()

	query := `
		INSERT INTO posting_concepts (
			concept_key
			, default_label
			, default_description
			, expected_account_type_id
			, normal_side
			, example_code
			, example_account_name
			, is_system
			, concept_category
			, sort_order
			, notes
			, metadata
			, deleted_at
		) VALUES (
			$1
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
		)
		RETURNING concept_key, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.ConceptKey,
		entity.DefaultLabel,
		entity.DefaultDescription,
		entity.ExpectedAccountTypeId,
		entity.NormalSide,
		entity.ExampleCode,
		entity.ExampleAccountName,
		entity.IsSystem,
		entity.ConceptCategory,
		entity.SortOrder,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.ConceptKey, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create posting_concepts", zap.Error(err))
		return fmt.Errorf("failed to create posting_concepts: %w", err)
	}

	r.logger.Info("created posting_concepts",
		zap.String("id", entity.ConceptKey.String()),
		
	)

	return nil
}

// GetByID retrieves a posting_concepts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingConcepts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_concepts", duration, nil)
	}()

	query := `
		SELECT
			concept_key
			, default_label
			, default_description
			, expected_account_type_id
			, normal_side
			, example_code
			, example_account_name
			, is_system
			, concept_category
			, sort_order
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_concepts
		WHERE concept_key = $1
		AND deleted_at IS NULL
	`

	var entity PostingConcepts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.ConceptKey,
		&entity.DefaultLabel,
		&entity.DefaultDescription,
		&entity.ExpectedAccountTypeId,
		&entity.NormalSide,
		&entity.ExampleCode,
		&entity.ExampleAccountName,
		&entity.IsSystem,
		&entity.ConceptCategory,
		&entity.SortOrder,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_concepts not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_concepts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_concepts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_concepts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingConcepts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_concepts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_concepts
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_concepts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			concept_key
			, default_label
			, default_description
			, expected_account_type_id
			, normal_side
			, example_code
			, example_account_name
			, is_system
			, concept_category
			, sort_order
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_concepts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_concepts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_concepts: %w", err)
	}
	defer rows.Close()

	var entities []*PostingConcepts
	for rows.Next() {
		var entity PostingConcepts
		err := rows.Scan(
			&entity.ConceptKey,
			&entity.DefaultLabel,
			&entity.DefaultDescription,
			&entity.ExpectedAccountTypeId,
			&entity.NormalSide,
			&entity.ExampleCode,
			&entity.ExampleAccountName,
			&entity.IsSystem,
			&entity.ConceptCategory,
			&entity.SortOrder,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_concepts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_concepts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_concepts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingConcepts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_concepts", duration, nil)
	}()

	query := `
		UPDATE posting_concepts
		SET
			concept_key = $1
			, default_label = $2
			, default_description = $3
			, expected_account_type_id = $4
			, normal_side = $5
			, example_code = $6
			, example_account_name = $7
			, is_system = $8
			, concept_category = $9
			, sort_order = $10
			, notes = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE concept_key = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.ConceptKey,
		entity.DefaultLabel,
		entity.DefaultDescription,
		entity.ExpectedAccountTypeId,
		entity.NormalSide,
		entity.ExampleCode,
		entity.ExampleAccountName,
		entity.IsSystem,
		entity.ConceptCategory,
		entity.SortOrder,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.ConceptKey,
	)

	if err != nil {
		r.logger.Error("failed to update posting_concepts", zap.Error(err))
		return fmt.Errorf("failed to update posting_concepts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_concepts not found or already deleted")
	}

	r.logger.Info("updated posting_concepts",
		zap.String("id", entity.ConceptKey.String()),
	)

	return nil
}


// Delete soft-deletes a posting_concepts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_concepts", duration, nil)
	}()

	query := `
		UPDATE posting_concepts
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE concept_key = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_concepts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_concepts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_concepts not found or already deleted")
	}

	r.logger.Info("deleted posting_concepts", zap.String("id", id.String()))
	return nil
}



