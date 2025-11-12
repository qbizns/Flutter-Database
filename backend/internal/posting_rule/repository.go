package posting_rule

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

// Repository handles database operations for PostingRules
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingRules repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingRules represents a posting_rules entity
type PostingRules struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PostingProfileDocumentId uuid.UUID `json:"posting_profile_document_id" db:"posting_profile_document_id"`
	RuleCode string `json:"rule_code" db:"rule_code"`
	RuleName string `json:"rule_name" db:"rule_name"`
	Description *string `json:"description" db:"description"`
	Event string `json:"event" db:"event"`
	Level string `json:"level" db:"level"`
	Priority int64 `json:"priority" db:"priority"`
	ConditionExpression *string `json:"condition_expression" db:"condition_expression"`
	IsActive bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new posting_rules record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_rules", duration, nil)
	}()

	query := `
		INSERT INTO posting_rules (
			, posting_profile_document_id
			, rule_code
			, rule_name
			, description
			, event
			, level
			, priority
			, condition_expression
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
			, $9
			, $10
			, $11
			, $12
			, $15
			, $16
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PostingProfileDocumentId,
		entity.RuleCode,
		entity.RuleName,
		entity.Description,
		entity.Event,
		entity.Level,
		entity.Priority,
		entity.ConditionExpression,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_rules", zap.Error(err))
		return fmt.Errorf("failed to create posting_rules: %w", err)
	}

	r.logger.Info("created posting_rules",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a posting_rules by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingRules, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_rules", duration, nil)
	}()

	query := `
		SELECT
			id
			, posting_profile_document_id
			, rule_code
			, rule_name
			, description
			, event
			, level
			, priority
			, condition_expression
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_rules
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingRules
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PostingProfileDocumentId,
		&entity.RuleCode,
		&entity.RuleName,
		&entity.Description,
		&entity.Event,
		&entity.Level,
		&entity.Priority,
		&entity.ConditionExpression,
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
		return nil, fmt.Errorf("posting_rules not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_rules", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_rules: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_rules records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingRules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_rules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_rules
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_rules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, posting_profile_document_id
			, rule_code
			, rule_name
			, description
			, event
			, level
			, priority
			, condition_expression
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM posting_rules
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_rules", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_rules: %w", err)
	}
	defer rows.Close()

	var entities []*PostingRules
	for rows.Next() {
		var entity PostingRules
		err := rows.Scan(
			&entity.Id,
			&entity.PostingProfileDocumentId,
			&entity.RuleCode,
			&entity.RuleName,
			&entity.Description,
			&entity.Event,
			&entity.Level,
			&entity.Priority,
			&entity.ConditionExpression,
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
			return nil, 0, fmt.Errorf("failed to scan posting_rules: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_rules rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_rules record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_rules", duration, nil)
	}()

	query := `
		UPDATE posting_rules
		SET
			, posting_profile_document_id = $2
			, rule_code = $3
			, rule_name = $4
			, description = $5
			, event = $6
			, level = $7
			, priority = $8
			, condition_expression = $9
			, is_active = $10
			, notes = $11
			, metadata = $12
			, updated_at = $14
			, deleted_at = $15
			, created_by = $16
			, updated_by = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PostingProfileDocumentId,
		entity.RuleCode,
		entity.RuleName,
		entity.Description,
		entity.Event,
		entity.Level,
		entity.Priority,
		entity.ConditionExpression,
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
		r.logger.Error("failed to update posting_rules", zap.Error(err))
		return fmt.Errorf("failed to update posting_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_rules not found or already deleted")
	}

	r.logger.Info("updated posting_rules",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_rules record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_rules", duration, nil)
	}()

	query := `
		UPDATE posting_rules
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_rules", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_rules not found or already deleted")
	}

	r.logger.Info("deleted posting_rules", zap.String("id", id.String()))
	return nil
}



