package posting_validation_result

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

// Repository handles database operations for PostingValidationResults
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingValidationResults repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingValidationResults represents a posting_validation_results entity
type PostingValidationResults struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	DocumentTypeCode string `json:"document_type_code" db:"document_type_code"`
	DocumentId uuid.UUID `json:"document_id" db:"document_id"`
	Event string `json:"event" db:"event"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	ValidationRuleId *uuid.UUID `json:"validation_rule_id" db:"validation_rule_id"`
	Severity string `json:"severity" db:"severity"`
	MessageCode string `json:"message_code" db:"message_code"`
	Message string `json:"message" db:"message"`
	IsBlocking bool `json:"is_blocking" db:"is_blocking"`
	Context json.RawMessage `json:"context" db:"context"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new posting_validation_results record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingValidationResults) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_validation_results", duration, nil)
	}()

	query := `
		INSERT INTO posting_validation_results (
			, organization_id
			, document_type_code
			, document_id
			, event
			, journal_entry_id
			, validation_rule_id
			, severity
			, message_code
			, message
			, is_blocking
			, context
			, metadata
			, created_by
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
			, $15
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.DocumentTypeCode,
		entity.DocumentId,
		entity.Event,
		entity.JournalEntryId,
		entity.ValidationRuleId,
		entity.Severity,
		entity.MessageCode,
		entity.Message,
		entity.IsBlocking,
		entity.Context,
		entity.Metadata,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_validation_results", zap.Error(err))
		return fmt.Errorf("failed to create posting_validation_results: %w", err)
	}

	r.logger.Info("created posting_validation_results",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a posting_validation_results by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingValidationResults, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_results", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, document_id
			, event
			, journal_entry_id
			, validation_rule_id
			, severity
			, message_code
			, message
			, is_blocking
			, context
			, metadata
			, created_at
			, created_by
		FROM posting_validation_results
		WHERE id = $1
		
	`

	var entity PostingValidationResults
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.DocumentTypeCode,
		&entity.DocumentId,
		&entity.Event,
		&entity.JournalEntryId,
		&entity.ValidationRuleId,
		&entity.Severity,
		&entity.MessageCode,
		&entity.Message,
		&entity.IsBlocking,
		&entity.Context,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_validation_results not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_validation_results", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_validation_results: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_validation_results records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingValidationResults, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_results", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_validation_results
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_validation_results records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, document_id
			, event
			, journal_entry_id
			, validation_rule_id
			, severity
			, message_code
			, message
			, is_blocking
			, context
			, metadata
			, created_at
			, created_by
		FROM posting_validation_results
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_validation_results", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_validation_results: %w", err)
	}
	defer rows.Close()

	var entities []*PostingValidationResults
	for rows.Next() {
		var entity PostingValidationResults
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentTypeCode,
			&entity.DocumentId,
			&entity.Event,
			&entity.JournalEntryId,
			&entity.ValidationRuleId,
			&entity.Severity,
			&entity.MessageCode,
			&entity.Message,
			&entity.IsBlocking,
			&entity.Context,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_validation_results: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_validation_results rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_validation_results record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingValidationResults) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_validation_results", duration, nil)
	}()

	query := `
		UPDATE posting_validation_results
		SET
			, organization_id = $2
			, document_type_code = $3
			, document_id = $4
			, event = $5
			, journal_entry_id = $6
			, validation_rule_id = $7
			, severity = $8
			, message_code = $9
			, message = $10
			, is_blocking = $11
			, context = $12
			, metadata = $13
			, created_by = $15
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $16
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.DocumentTypeCode,
		entity.DocumentId,
		entity.Event,
		entity.JournalEntryId,
		entity.ValidationRuleId,
		entity.Severity,
		entity.MessageCode,
		entity.Message,
		entity.IsBlocking,
		entity.Context,
		entity.Metadata,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_validation_results", zap.Error(err))
		return fmt.Errorf("failed to update posting_validation_results: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_validation_results not found or already deleted")
	}

	r.logger.Info("updated posting_validation_results",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a posting_validation_results record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "posting_validation_results", duration, nil)
	}()

	query := `DELETE FROM posting_validation_results WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_validation_results", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_validation_results: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_validation_results not found")
	}

	r.logger.Info("deleted posting_validation_results", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves posting_validation_results records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PostingValidationResults, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_results", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_validation_results
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_validation_results records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, document_id
			, event
			, journal_entry_id
			, validation_rule_id
			, severity
			, message_code
			, message
			, is_blocking
			, context
			, metadata
			, created_at
			, created_by
		FROM posting_validation_results
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_validation_results by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_validation_results: %w", err)
	}
	defer rows.Close()

	var entities []*PostingValidationResults
	for rows.Next() {
		var entity PostingValidationResults
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentTypeCode,
			&entity.DocumentId,
			&entity.Event,
			&entity.JournalEntryId,
			&entity.ValidationRuleId,
			&entity.Severity,
			&entity.MessageCode,
			&entity.Message,
			&entity.IsBlocking,
			&entity.Context,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_validation_results: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

