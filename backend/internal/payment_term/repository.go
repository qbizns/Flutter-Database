package payment_term

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

// Repository handles database operations for PaymentTerms
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PaymentTerms repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PaymentTerms represents a payment_terms entity
type PaymentTerms struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	TermCode string `json:"term_code" db:"term_code"`
	TermName string `json:"term_name" db:"term_name"`
	Note *string `json:"note" db:"note"`
	IsActive *bool `json:"is_active" db:"is_active"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new payment_terms record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PaymentTerms) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "payment_terms", duration, nil)
	}()

	query := `
		INSERT INTO payment_terms (
			, organization_id
			, term_code
			, term_name
			, note
			, is_active
			, created_by
			, updated_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $11
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TermCode,
		entity.TermName,
		entity.Note,
		entity.IsActive,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create payment_terms", zap.Error(err))
		return fmt.Errorf("failed to create payment_terms: %w", err)
	}

	r.logger.Info("created payment_terms",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a payment_terms by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PaymentTerms, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payment_terms", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, term_code
			, term_name
			, note
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM payment_terms
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PaymentTerms
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TermCode,
		&entity.TermName,
		&entity.Note,
		&entity.IsActive,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("payment_terms not found")
	}

	if err != nil {
		r.logger.Error("failed to get payment_terms", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get payment_terms: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of payment_terms records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PaymentTerms, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payment_terms", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM payment_terms
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payment_terms records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, term_code
			, term_name
			, note
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM payment_terms
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list payment_terms", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list payment_terms: %w", err)
	}
	defer rows.Close()

	var entities []*PaymentTerms
	for rows.Next() {
		var entity PaymentTerms
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TermCode,
			&entity.TermName,
			&entity.Note,
			&entity.IsActive,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment_terms: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating payment_terms rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing payment_terms record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PaymentTerms) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "payment_terms", duration, nil)
	}()

	query := `
		UPDATE payment_terms
		SET
			, organization_id = $2
			, term_code = $3
			, term_name = $4
			, note = $5
			, is_active = $6
			, created_by = $7
			, updated_by = $8
			, updated_at = $10
			, deleted_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TermCode,
		entity.TermName,
		entity.Note,
		entity.IsActive,
		entity.CreatedBy,
		entity.UpdatedBy,
		time.Now(),
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update payment_terms", zap.Error(err))
		return fmt.Errorf("failed to update payment_terms: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment_terms not found or already deleted")
	}

	r.logger.Info("updated payment_terms",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a payment_terms record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "payment_terms", duration, nil)
	}()

	query := `
		UPDATE payment_terms
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete payment_terms", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete payment_terms: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("payment_terms not found or already deleted")
	}

	r.logger.Info("deleted payment_terms", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves payment_terms records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PaymentTerms, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "payment_terms", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM payment_terms
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payment_terms records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, term_code
			, term_name
			, note
			, is_active
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM payment_terms
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list payment_terms by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list payment_terms: %w", err)
	}
	defer rows.Close()

	var entities []*PaymentTerms
	for rows.Next() {
		var entity PaymentTerms
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TermCode,
			&entity.TermName,
			&entity.Note,
			&entity.IsActive,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan payment_terms: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

