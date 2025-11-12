package journal_entry_type

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

// Repository handles database operations for JournalEntryTypes
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new JournalEntryTypes repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// JournalEntryTypes represents a journal_entry_types entity
type JournalEntryTypes struct {
	Id *uuid.UUID `json:"id" db:"id"`
	TypeCode string `json:"type_code" db:"type_code"`
	TypeName string `json:"type_name" db:"type_name"`
	TypeCategory *string `json:"type_category" db:"type_category"`
	NumberPrefix *string `json:"number_prefix" db:"number_prefix"`
	Description *string `json:"description" db:"description"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new journal_entry_types record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *JournalEntryTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "journal_entry_types", duration, nil)
	}()

	query := `
		INSERT INTO journal_entry_types (
			, type_code
			, type_name
			, type_category
			, number_prefix
			, description
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.TypeCode,
		entity.TypeName,
		entity.TypeCategory,
		entity.NumberPrefix,
		entity.Description,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create journal_entry_types", zap.Error(err))
		return fmt.Errorf("failed to create journal_entry_types: %w", err)
	}

	r.logger.Info("created journal_entry_types",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a journal_entry_types by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*JournalEntryTypes, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entry_types", duration, nil)
	}()

	query := `
		SELECT
			id
			, type_code
			, type_name
			, type_category
			, number_prefix
			, description
			, created_at
			, updated_at
		FROM journal_entry_types
		WHERE id = $1
		
	`

	var entity JournalEntryTypes
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.TypeCode,
		&entity.TypeName,
		&entity.TypeCategory,
		&entity.NumberPrefix,
		&entity.Description,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("journal_entry_types not found")
	}

	if err != nil {
		r.logger.Error("failed to get journal_entry_types", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get journal_entry_types: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of journal_entry_types records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*JournalEntryTypes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "journal_entry_types", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM journal_entry_types
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count journal_entry_types records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, type_code
			, type_name
			, type_category
			, number_prefix
			, description
			, created_at
			, updated_at
		FROM journal_entry_types
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list journal_entry_types", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list journal_entry_types: %w", err)
	}
	defer rows.Close()

	var entities []*JournalEntryTypes
	for rows.Next() {
		var entity JournalEntryTypes
		err := rows.Scan(
			&entity.Id,
			&entity.TypeCode,
			&entity.TypeName,
			&entity.TypeCategory,
			&entity.NumberPrefix,
			&entity.Description,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan journal_entry_types: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating journal_entry_types rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing journal_entry_types record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *JournalEntryTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "journal_entry_types", duration, nil)
	}()

	query := `
		UPDATE journal_entry_types
		SET
			, type_code = $2
			, type_name = $3
			, type_category = $4
			, number_prefix = $5
			, description = $6
			, updated_at = $8
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.TypeCode,
		entity.TypeName,
		entity.TypeCategory,
		entity.NumberPrefix,
		entity.Description,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update journal_entry_types", zap.Error(err))
		return fmt.Errorf("failed to update journal_entry_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entry_types not found or already deleted")
	}

	r.logger.Info("updated journal_entry_types",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a journal_entry_types record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "journal_entry_types", duration, nil)
	}()

	query := `DELETE FROM journal_entry_types WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete journal_entry_types", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete journal_entry_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("journal_entry_types not found")
	}

	r.logger.Info("deleted journal_entry_types", zap.String("id", id.String()))
	return nil
}



