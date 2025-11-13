package posting_document_type

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

// Repository handles database operations for PostingDocumentTypes
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingDocumentTypes repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingDocumentTypes represents a posting_document_types entity
type PostingDocumentTypes struct {
	Id *uuid.UUID `json:"id" db:"id"`
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	SourceSchema string `json:"source_schema" db:"source_schema"`
	SourceTable string `json:"source_table" db:"source_table"`
	SourcePkColumn string `json:"source_pk_column" db:"source_pk_column"`
	Category *string `json:"category" db:"category"`
	IsActive bool `json:"is_active" db:"is_active"`
	IsSystem bool `json:"is_system" db:"is_system"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new posting_document_types record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingDocumentTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_document_types", duration, nil)
	}()

	query := `
		INSERT INTO posting_document_types (
			, code
			, name
			, description
			, source_schema
			, source_table
			, source_pk_column
			, category
			, is_active
			, is_system
			, notes
			, metadata
			, deleted_at
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.SourceSchema,
		entity.SourceTable,
		entity.SourcePkColumn,
		entity.Category,
		entity.IsActive,
		entity.IsSystem,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create posting_document_types", zap.Error(err))
		return fmt.Errorf("failed to create posting_document_types: %w", err)
	}

	r.logger.Info("created posting_document_types",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a posting_document_types by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingDocumentTypes, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_document_types", duration, nil)
	}()

	query := `
		SELECT
			id
			, code
			, name
			, description
			, source_schema
			, source_table
			, source_pk_column
			, category
			, is_active
			, is_system
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_document_types
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingDocumentTypes
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.Code,
		&entity.Name,
		&entity.Description,
		&entity.SourceSchema,
		&entity.SourceTable,
		&entity.SourcePkColumn,
		&entity.Category,
		&entity.IsActive,
		&entity.IsSystem,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_document_types not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_document_types", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_document_types: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_document_types records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingDocumentTypes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_document_types", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_document_types
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_document_types records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, code
			, name
			, description
			, source_schema
			, source_table
			, source_pk_column
			, category
			, is_active
			, is_system
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_document_types
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_document_types", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_document_types: %w", err)
	}
	defer rows.Close()

	var entities []*PostingDocumentTypes
	for rows.Next() {
		var entity PostingDocumentTypes
		err := rows.Scan(
			&entity.Id,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.SourceSchema,
			&entity.SourceTable,
			&entity.SourcePkColumn,
			&entity.Category,
			&entity.IsActive,
			&entity.IsSystem,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_document_types: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_document_types rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_document_types record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingDocumentTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_document_types", duration, nil)
	}()

	query := `
		UPDATE posting_document_types
		SET
			, code = $2
			, name = $3
			, description = $4
			, source_schema = $5
			, source_table = $6
			, source_pk_column = $7
			, category = $8
			, is_active = $9
			, is_system = $10
			, notes = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.SourceSchema,
		entity.SourceTable,
		entity.SourcePkColumn,
		entity.Category,
		entity.IsActive,
		entity.IsSystem,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_document_types", zap.Error(err))
		return fmt.Errorf("failed to update posting_document_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_document_types not found or already deleted")
	}

	r.logger.Info("updated posting_document_types",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_document_types record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_document_types", duration, nil)
	}()

	query := `
		UPDATE posting_document_types
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_document_types", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_document_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_document_types not found or already deleted")
	}

	r.logger.Info("deleted posting_document_types", zap.String("id", id.String()))
	return nil
}



