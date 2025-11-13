package posting_rule_line

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

// Repository handles database operations for PostingRuleLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingRuleLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingRuleLines represents a posting_rule_lines entity
type PostingRuleLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	PostingRuleId uuid.UUID `json:"posting_rule_id" db:"posting_rule_id"`
	LineNo int64 `json:"line_no" db:"line_no"`
	Side string `json:"side" db:"side"`
	ConceptKey *string `json:"concept_key" db:"concept_key"`
	AccountSource string `json:"account_source" db:"account_source"`
	FixedAccountId *uuid.UUID `json:"fixed_account_id" db:"fixed_account_id"`
	AccountFieldPath *string `json:"account_field_path" db:"account_field_path"`
	AccountExpression *string `json:"account_expression" db:"account_expression"`
	AmountSource string `json:"amount_source" db:"amount_source"`
	AmountFieldPath *string `json:"amount_field_path" db:"amount_field_path"`
	AmountExpression *string `json:"amount_expression" db:"amount_expression"`
	MappingContext json.RawMessage `json:"mapping_context" db:"mapping_context"`
	DescriptionTemplate *string `json:"description_template" db:"description_template"`
	IsActive bool `json:"is_active" db:"is_active"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	(accountSource string `json:"(account_source" db:"(account_source"`
	(accountSource string `json:"(account_source" db:"(account_source"`
	(accountSource string `json:"(account_source" db:"(account_source"`
	(accountSource string `json:"(account_source" db:"(account_source"`
	(amountSource string `json:"(amount_source" db:"(amount_source"`
	(amountSource string `json:"(amount_source" db:"(amount_source"`
	(amountSource string `json:"(amount_source" db:"(amount_source"`
}

// Create inserts a new posting_rule_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingRuleLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_rule_lines", duration, nil)
	}()

	query := `
		INSERT INTO posting_rule_lines (
			, posting_rule_id
			, line_no
			, side
			, concept_key
			, account_source
			, fixed_account_id
			, account_field_path
			, account_expression
			, amount_source
			, amount_field_path
			, amount_expression
			, mapping_context
			, description_template
			, is_active
			, notes
			, metadata
			, deleted_at
			, (account_source
			, (account_source
			, (account_source
			, (account_source
			, (amount_source
			, (amount_source
			, (amount_source
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
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.PostingRuleId,
		entity.LineNo,
		entity.Side,
		entity.ConceptKey,
		entity.AccountSource,
		entity.FixedAccountId,
		entity.AccountFieldPath,
		entity.AccountExpression,
		entity.AmountSource,
		entity.AmountFieldPath,
		entity.AmountExpression,
		entity.MappingContext,
		entity.DescriptionTemplate,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.(accountSource,
		entity.(accountSource,
		entity.(accountSource,
		entity.(accountSource,
		entity.(amountSource,
		entity.(amountSource,
		entity.(amountSource,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create posting_rule_lines", zap.Error(err))
		return fmt.Errorf("failed to create posting_rule_lines: %w", err)
	}

	r.logger.Info("created posting_rule_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a posting_rule_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingRuleLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_rule_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, posting_rule_id
			, line_no
			, side
			, concept_key
			, account_source
			, fixed_account_id
			, account_field_path
			, account_expression
			, amount_source
			, amount_field_path
			, amount_expression
			, mapping_context
			, description_template
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, (account_source
			, (account_source
			, (account_source
			, (account_source
			, (amount_source
			, (amount_source
			, (amount_source
		FROM posting_rule_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingRuleLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.PostingRuleId,
		&entity.LineNo,
		&entity.Side,
		&entity.ConceptKey,
		&entity.AccountSource,
		&entity.FixedAccountId,
		&entity.AccountFieldPath,
		&entity.AccountExpression,
		&entity.AmountSource,
		&entity.AmountFieldPath,
		&entity.AmountExpression,
		&entity.MappingContext,
		&entity.DescriptionTemplate,
		&entity.IsActive,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.(accountSource,
		&entity.(accountSource,
		&entity.(accountSource,
		&entity.(accountSource,
		&entity.(amountSource,
		&entity.(amountSource,
		&entity.(amountSource,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_rule_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_rule_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_rule_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_rule_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingRuleLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_rule_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_rule_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_rule_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, posting_rule_id
			, line_no
			, side
			, concept_key
			, account_source
			, fixed_account_id
			, account_field_path
			, account_expression
			, amount_source
			, amount_field_path
			, amount_expression
			, mapping_context
			, description_template
			, is_active
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, (account_source
			, (account_source
			, (account_source
			, (account_source
			, (amount_source
			, (amount_source
			, (amount_source
		FROM posting_rule_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_rule_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_rule_lines: %w", err)
	}
	defer rows.Close()

	var entities []*PostingRuleLines
	for rows.Next() {
		var entity PostingRuleLines
		err := rows.Scan(
			&entity.Id,
			&entity.PostingRuleId,
			&entity.LineNo,
			&entity.Side,
			&entity.ConceptKey,
			&entity.AccountSource,
			&entity.FixedAccountId,
			&entity.AccountFieldPath,
			&entity.AccountExpression,
			&entity.AmountSource,
			&entity.AmountFieldPath,
			&entity.AmountExpression,
			&entity.MappingContext,
			&entity.DescriptionTemplate,
			&entity.IsActive,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.(accountSource,
			&entity.(accountSource,
			&entity.(accountSource,
			&entity.(accountSource,
			&entity.(amountSource,
			&entity.(amountSource,
			&entity.(amountSource,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_rule_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_rule_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_rule_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingRuleLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_rule_lines", duration, nil)
	}()

	query := `
		UPDATE posting_rule_lines
		SET
			, posting_rule_id = $2
			, line_no = $3
			, side = $4
			, concept_key = $5
			, account_source = $6
			, fixed_account_id = $7
			, account_field_path = $8
			, account_expression = $9
			, amount_source = $10
			, amount_field_path = $11
			, amount_expression = $12
			, mapping_context = $13
			, description_template = $14
			, is_active = $15
			, notes = $16
			, metadata = $17
			, updated_at = $19
			, deleted_at = $20
			, (account_source = $21
			, (account_source = $22
			, (account_source = $23
			, (account_source = $24
			, (amount_source = $25
			, (amount_source = $26
			, (amount_source = $27
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $28
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.PostingRuleId,
		entity.LineNo,
		entity.Side,
		entity.ConceptKey,
		entity.AccountSource,
		entity.FixedAccountId,
		entity.AccountFieldPath,
		entity.AccountExpression,
		entity.AmountSource,
		entity.AmountFieldPath,
		entity.AmountExpression,
		entity.MappingContext,
		entity.DescriptionTemplate,
		entity.IsActive,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.(accountSource,
		entity.(accountSource,
		entity.(accountSource,
		entity.(accountSource,
		entity.(amountSource,
		entity.(amountSource,
		entity.(amountSource,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_rule_lines", zap.Error(err))
		return fmt.Errorf("failed to update posting_rule_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_rule_lines not found or already deleted")
	}

	r.logger.Info("updated posting_rule_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_rule_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_rule_lines", duration, nil)
	}()

	query := `
		UPDATE posting_rule_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_rule_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_rule_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_rule_lines not found or already deleted")
	}

	r.logger.Info("deleted posting_rule_lines", zap.String("id", id.String()))
	return nil
}



