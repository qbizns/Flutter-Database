package posting_profile_document

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

// Repository handles database operations for PostingProfileDocuments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingProfileDocuments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingProfileDocuments represents a posting_profile_documents entity
type PostingProfileDocuments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PostingProfileId uuid.UUID `json:"posting_profile_id" db:"posting_profile_id"`
	PostingDocumentTypeId uuid.UUID `json:"posting_document_type_id" db:"posting_document_type_id"`
	IsActive bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new posting_profile_documents record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingProfileDocuments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_profile_documents", duration, nil)
	}()

	query := `
		INSERT INTO posting_profile_documents (
			, posting_profile_id
			, posting_document_type_id
			, is_active
			, notes
			, metadata
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $9
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PostingProfileId,
		entity.PostingDocumentTypeId,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_profile_documents", zap.Error(err))
		return fmt.Errorf("failed to create posting_profile_documents: %w", err)
	}

	r.logger.Info("created posting_profile_documents",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a posting_profile_documents by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingProfileDocuments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_profile_documents", duration, nil)
	}()

	query := `
		SELECT
			id
			, posting_profile_id
			, posting_document_type_id
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_profile_documents
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingProfileDocuments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PostingProfileId,
		&entity.PostingDocumentTypeId,
		&entity.IsActive,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_profile_documents not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_profile_documents", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_profile_documents: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_profile_documents records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingProfileDocuments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_profile_documents", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_profile_documents
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_profile_documents records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, posting_profile_id
			, posting_document_type_id
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM posting_profile_documents
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_profile_documents", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_profile_documents: %w", err)
	}
	defer rows.Close()

	var entities []*PostingProfileDocuments
	for rows.Next() {
		var entity PostingProfileDocuments
		err := rows.Scan(
			&entity.Id,
			&entity.PostingProfileId,
			&entity.PostingDocumentTypeId,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_profile_documents: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_profile_documents rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_profile_documents record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingProfileDocuments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_profile_documents", duration, nil)
	}()

	query := `
		UPDATE posting_profile_documents
		SET
			, posting_profile_id = $2
			, posting_document_type_id = $3
			, is_active = $4
			, notes = $5
			, metadata = $6
			, updated_at = $8
			, deleted_at = $9
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $10
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PostingProfileId,
		entity.PostingDocumentTypeId,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_profile_documents", zap.Error(err))
		return fmt.Errorf("failed to update posting_profile_documents: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_profile_documents not found or already deleted")
	}

	r.logger.Info("updated posting_profile_documents",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_profile_documents record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_profile_documents", duration, nil)
	}()

	query := `
		UPDATE posting_profile_documents
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_profile_documents", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_profile_documents: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_profile_documents not found or already deleted")
	}

	r.logger.Info("deleted posting_profile_documents", zap.String("id", id.String()))
	return nil
}



