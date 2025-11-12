package document_sequence

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

// Repository handles database operations for DocumentSequences
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new DocumentSequences repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// DocumentSequences represents a document_sequences entity
type DocumentSequences struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	DocumentType string `json:"document_type" db:"document_type"`
	Prefix *string `json:"prefix" db:"prefix"`
	Suffix *string `json:"suffix" db:"suffix"`
	NextNumber int64 `json:"next_number" db:"next_number"`
	Padding *int64 `json:"padding" db:"padding"`
	IncrementBy *int64 `json:"increment_by" db:"increment_by"`
	ResetFrequency *string `json:"reset_frequency" db:"reset_frequency"`
	'never', *string `json:"'never'," db:"'never',"`
	'daily', *string `json:"'daily'," db:"'daily',"`
	'monthly', *string `json:"'monthly'," db:"'monthly',"`
	'yearly', *string `json:"'yearly'," db:"'yearly',"`
	'manual' *string `json:"'manual'" db:"'manual'"`
	LastResetAt *time.Time `json:"last_reset_at" db:"last_reset_at"`
	LastResetValue *int64 `json:"last_reset_value" db:"last_reset_value"`
	IncludeDate *bool `json:"include_date" db:"include_date"`
	DateFormat *string `json:"date_format" db:"date_format"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	AllowManualOverride *bool `json:"allow_manual_override" db:"allow_manual_override"`
	ExampleNumber *string `json:"example_number" db:"example_number"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new document_sequences record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *DocumentSequences) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "document_sequences", duration, nil)
	}()

	query := `
		INSERT INTO document_sequences (
			, organization_id
			, document_type
			, prefix
			, suffix
			, next_number
			, padding
			, increment_by
			, reset_frequency
			, 'never',
			, 'daily',
			, 'monthly',
			, 'yearly',
			, 'manual'
			, last_reset_at
			, last_reset_value
			, include_date
			, date_format
			, location_id
			, is_active
			, allow_manual_override
			, example_number
			, description
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
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $28
			, $29
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.DocumentType,
		entity.Prefix,
		entity.Suffix,
		entity.NextNumber,
		entity.Padding,
		entity.IncrementBy,
		entity.ResetFrequency,
		entity.'never',,
		entity.'daily',,
		entity.'monthly',,
		entity.'yearly',,
		entity.'manual',
		entity.LastResetAt,
		entity.LastResetValue,
		entity.IncludeDate,
		entity.DateFormat,
		entity.LocationId,
		entity.IsActive,
		entity.AllowManualOverride,
		entity.ExampleNumber,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create document_sequences", zap.Error(err))
		return fmt.Errorf("failed to create document_sequences: %w", err)
	}

	r.logger.Info("created document_sequences",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a document_sequences by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*DocumentSequences, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "document_sequences", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, document_type
			, prefix
			, suffix
			, next_number
			, padding
			, increment_by
			, reset_frequency
			, 'never',
			, 'daily',
			, 'monthly',
			, 'yearly',
			, 'manual'
			, last_reset_at
			, last_reset_value
			, include_date
			, date_format
			, location_id
			, is_active
			, allow_manual_override
			, example_number
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM document_sequences
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity DocumentSequences
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.DocumentType,
		&entity.Prefix,
		&entity.Suffix,
		&entity.NextNumber,
		&entity.Padding,
		&entity.IncrementBy,
		&entity.ResetFrequency,
		&entity.'never',,
		&entity.'daily',,
		&entity.'monthly',,
		&entity.'yearly',,
		&entity.'manual',
		&entity.LastResetAt,
		&entity.LastResetValue,
		&entity.IncludeDate,
		&entity.DateFormat,
		&entity.LocationId,
		&entity.IsActive,
		&entity.AllowManualOverride,
		&entity.ExampleNumber,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("document_sequences not found")
	}

	if err != nil {
		r.logger.Error("failed to get document_sequences", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get document_sequences: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of document_sequences records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*DocumentSequences, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "document_sequences", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM document_sequences
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count document_sequences records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type
			, prefix
			, suffix
			, next_number
			, padding
			, increment_by
			, reset_frequency
			, 'never',
			, 'daily',
			, 'monthly',
			, 'yearly',
			, 'manual'
			, last_reset_at
			, last_reset_value
			, include_date
			, date_format
			, location_id
			, is_active
			, allow_manual_override
			, example_number
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM document_sequences
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list document_sequences", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list document_sequences: %w", err)
	}
	defer rows.Close()

	var entities []*DocumentSequences
	for rows.Next() {
		var entity DocumentSequences
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentType,
			&entity.Prefix,
			&entity.Suffix,
			&entity.NextNumber,
			&entity.Padding,
			&entity.IncrementBy,
			&entity.ResetFrequency,
			&entity.'never',,
			&entity.'daily',,
			&entity.'monthly',,
			&entity.'yearly',,
			&entity.'manual',
			&entity.LastResetAt,
			&entity.LastResetValue,
			&entity.IncludeDate,
			&entity.DateFormat,
			&entity.LocationId,
			&entity.IsActive,
			&entity.AllowManualOverride,
			&entity.ExampleNumber,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan document_sequences: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating document_sequences rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing document_sequences record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *DocumentSequences) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "document_sequences", duration, nil)
	}()

	query := `
		UPDATE document_sequences
		SET
			, organization_id = $2
			, document_type = $3
			, prefix = $4
			, suffix = $5
			, next_number = $6
			, padding = $7
			, increment_by = $8
			, reset_frequency = $9
			, 'never', = $10
			, 'daily', = $11
			, 'monthly', = $12
			, 'yearly', = $13
			, 'manual' = $14
			, last_reset_at = $15
			, last_reset_value = $16
			, include_date = $17
			, date_format = $18
			, location_id = $19
			, is_active = $20
			, allow_manual_override = $21
			, example_number = $22
			, description = $23
			, notes = $24
			, metadata = $25
			, updated_at = $27
			, deleted_at = $28
			, created_by = $29
			, updated_by = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.DocumentType,
		entity.Prefix,
		entity.Suffix,
		entity.NextNumber,
		entity.Padding,
		entity.IncrementBy,
		entity.ResetFrequency,
		entity.'never',,
		entity.'daily',,
		entity.'monthly',,
		entity.'yearly',,
		entity.'manual',
		entity.LastResetAt,
		entity.LastResetValue,
		entity.IncludeDate,
		entity.DateFormat,
		entity.LocationId,
		entity.IsActive,
		entity.AllowManualOverride,
		entity.ExampleNumber,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update document_sequences", zap.Error(err))
		return fmt.Errorf("failed to update document_sequences: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("document_sequences not found or already deleted")
	}

	r.logger.Info("updated document_sequences",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a document_sequences record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "document_sequences", duration, nil)
	}()

	query := `
		UPDATE document_sequences
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete document_sequences", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete document_sequences: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("document_sequences not found or already deleted")
	}

	r.logger.Info("deleted document_sequences", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves document_sequences records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*DocumentSequences, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "document_sequences", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM document_sequences
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count document_sequences records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, document_type
			, prefix
			, suffix
			, next_number
			, padding
			, increment_by
			, reset_frequency
			, 'never',
			, 'daily',
			, 'monthly',
			, 'yearly',
			, 'manual'
			, last_reset_at
			, last_reset_value
			, include_date
			, date_format
			, location_id
			, is_active
			, allow_manual_override
			, example_number
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM document_sequences
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list document_sequences by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list document_sequences: %w", err)
	}
	defer rows.Close()

	var entities []*DocumentSequences
	for rows.Next() {
		var entity DocumentSequences
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.DocumentType,
			&entity.Prefix,
			&entity.Suffix,
			&entity.NextNumber,
			&entity.Padding,
			&entity.IncrementBy,
			&entity.ResetFrequency,
			&entity.'never',,
			&entity.'daily',,
			&entity.'monthly',,
			&entity.'yearly',,
			&entity.'manual',
			&entity.LastResetAt,
			&entity.LastResetValue,
			&entity.IncludeDate,
			&entity.DateFormat,
			&entity.LocationId,
			&entity.IsActive,
			&entity.AllowManualOverride,
			&entity.ExampleNumber,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan document_sequences: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

