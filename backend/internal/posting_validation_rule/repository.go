package posting_validation_rule

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

// Repository handles database operations for PostingValidationRules
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PostingValidationRules repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PostingValidationRules represents a posting_validation_rules entity
type PostingValidationRules struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId *uuid.UUID `json:"organization_id" db:"organization_id"`
	DocumentTypeCode *string `json:"document_type_code" db:"document_type_code"`
	Event *string `json:"event" db:"event"`
	Target string `json:"target" db:"target"`
	Code string `json:"code" db:"code"`
	Name string `json:"name" db:"name"`
	Description *string `json:"description" db:"description"`
	Expression string `json:"expression" db:"expression"`
	Severity string `json:"severity" db:"severity"`
	IsBlocking bool `json:"is_blocking" db:"is_blocking"`
	IsActive bool `json:"is_active" db:"is_active"`
	MessageTemplate *string `json:"message_template" db:"message_template"`
	Priority *int64 `json:"priority" db:"priority"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	COALESCE(organizationId, *uuid.UUID `json:"COALESCE(organization_id," db:"COALESCE(organization_id,"`
	COALESCE(documentTypeCode, *string `json:"COALESCE(document_type_code," db:"COALESCE(document_type_code,"`
	COALESCE(event, *string `json:"COALESCE(event," db:"COALESCE(event,"`
}

// Create inserts a new posting_validation_rules record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PostingValidationRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "posting_validation_rules", duration, nil)
	}()

	query := `
		INSERT INTO posting_validation_rules (
			, organization_id
			, document_type_code
			, event
			, target
			, code
			, name
			, description
			, expression
			, severity
			, is_blocking
			, is_active
			, message_template
			, priority
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, COALESCE(organization_id,
			, COALESCE(document_type_code,
			, COALESCE(event,
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
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.DocumentTypeCode,
		entity.Event,
		entity.Target,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.Expression,
		entity.Severity,
		entity.IsBlocking,
		entity.IsActive,
		entity.MessageTemplate,
		entity.Priority,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.COALESCE(organizationId,,
		entity.COALESCE(documentTypeCode,,
		entity.COALESCE(event,,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create posting_validation_rules", zap.Error(err))
		return fmt.Errorf("failed to create posting_validation_rules: %w", err)
	}

	r.logger.Info("created posting_validation_rules",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a posting_validation_rules by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PostingValidationRules, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_rules", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, event
			, target
			, code
			, name
			, description
			, expression
			, severity
			, is_blocking
			, is_active
			, message_template
			, priority
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, COALESCE(organization_id,
			, COALESCE(document_type_code,
			, COALESCE(event,
		FROM posting_validation_rules
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PostingValidationRules
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.DocumentTypeCode,
		&entity.Event,
		&entity.Target,
		&entity.Code,
		&entity.Name,
		&entity.Description,
		&entity.Expression,
		&entity.Severity,
		&entity.IsBlocking,
		&entity.IsActive,
		&entity.MessageTemplate,
		&entity.Priority,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.COALESCE(organizationId,,
		&entity.COALESCE(documentTypeCode,,
		&entity.COALESCE(event,,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("posting_validation_rules not found")
	}

	if err != nil {
		r.logger.Error("failed to get posting_validation_rules", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get posting_validation_rules: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of posting_validation_rules records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PostingValidationRules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_rules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_validation_rules
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_validation_rules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, event
			, target
			, code
			, name
			, description
			, expression
			, severity
			, is_blocking
			, is_active
			, message_template
			, priority
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, COALESCE(organization_id,
			, COALESCE(document_type_code,
			, COALESCE(event,
		FROM posting_validation_rules
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_validation_rules", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_validation_rules: %w", err)
	}
	defer rows.Close()

	var entities []*PostingValidationRules
	for rows.Next() {
		var entity PostingValidationRules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentTypeCode,
			&entity.Event,
			&entity.Target,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.Expression,
			&entity.Severity,
			&entity.IsBlocking,
			&entity.IsActive,
			&entity.MessageTemplate,
			&entity.Priority,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.COALESCE(organizationId,,
			&entity.COALESCE(documentTypeCode,,
			&entity.COALESCE(event,,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_validation_rules: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating posting_validation_rules rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing posting_validation_rules record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PostingValidationRules) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_validation_rules", duration, nil)
	}()

	query := `
		UPDATE posting_validation_rules
		SET
			, organization_id = $2
			, document_type_code = $3
			, event = $4
			, target = $5
			, code = $6
			, name = $7
			, description = $8
			, expression = $9
			, severity = $10
			, is_blocking = $11
			, is_active = $12
			, message_template = $13
			, priority = $14
			, notes = $15
			, metadata = $16
			, updated_at = $18
			, deleted_at = $19
			, created_by = $20
			, updated_by = $21
			, COALESCE(organization_id, = $22
			, COALESCE(document_type_code, = $23
			, COALESCE(event, = $24
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $25
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.DocumentTypeCode,
		entity.Event,
		entity.Target,
		entity.Code,
		entity.Name,
		entity.Description,
		entity.Expression,
		entity.Severity,
		entity.IsBlocking,
		entity.IsActive,
		entity.MessageTemplate,
		entity.Priority,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.COALESCE(organizationId,,
		entity.COALESCE(documentTypeCode,,
		entity.COALESCE(event,,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update posting_validation_rules", zap.Error(err))
		return fmt.Errorf("failed to update posting_validation_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_validation_rules not found or already deleted")
	}

	r.logger.Info("updated posting_validation_rules",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a posting_validation_rules record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "posting_validation_rules", duration, nil)
	}()

	query := `
		UPDATE posting_validation_rules
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete posting_validation_rules", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete posting_validation_rules: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("posting_validation_rules not found or already deleted")
	}

	r.logger.Info("deleted posting_validation_rules", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves posting_validation_rules records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PostingValidationRules, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "posting_validation_rules", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM posting_validation_rules
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count posting_validation_rules records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type_code
			, event
			, target
			, code
			, name
			, description
			, expression
			, severity
			, is_blocking
			, is_active
			, message_template
			, priority
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, COALESCE(organization_id,
			, COALESCE(document_type_code,
			, COALESCE(event,
		FROM posting_validation_rules
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list posting_validation_rules by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list posting_validation_rules: %w", err)
	}
	defer rows.Close()

	var entities []*PostingValidationRules
	for rows.Next() {
		var entity PostingValidationRules
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentTypeCode,
			&entity.Event,
			&entity.Target,
			&entity.Code,
			&entity.Name,
			&entity.Description,
			&entity.Expression,
			&entity.Severity,
			&entity.IsBlocking,
			&entity.IsActive,
			&entity.MessageTemplate,
			&entity.Priority,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.COALESCE(organizationId,,
			&entity.COALESCE(documentTypeCode,,
			&entity.COALESCE(event,,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan posting_validation_rules: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

