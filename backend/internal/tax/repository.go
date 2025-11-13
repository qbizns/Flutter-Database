package tax

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

// Repository handles database operations for Taxes
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Taxes repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Taxes represents a taxes entity
type Taxes struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	TaxGroupId *uuid.UUID `json:"tax_group_id" db:"tax_group_id"`
	TaxCode string `json:"tax_code" db:"tax_code"`
	TaxName string `json:"tax_name" db:"tax_name"`
	TaxRate float64 `json:"tax_rate" db:"tax_rate"`
	TaxScope string `json:"tax_scope" db:"tax_scope"`
	IsPriceInclusive *bool `json:"is_price_inclusive" db:"is_price_inclusive"`
	TaxAccountId uuid.UUID `json:"tax_account_id" db:"tax_account_id"`
	TaxRefundAccountId *uuid.UUID `json:"tax_refund_account_id" db:"tax_refund_account_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Description *string `json:"description" db:"description"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new taxes record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Taxes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "taxes", duration, nil)
	}()

	query := `
		INSERT INTO taxes (
			, organization_id
			, tax_group_id
			, tax_code
			, tax_name
			, tax_rate
			, tax_scope
			, is_price_inclusive
			, tax_account_id
			, tax_refund_account_id
			, is_active
			, description
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
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.TaxGroupId,
		entity.TaxCode,
		entity.TaxName,
		entity.TaxRate,
		entity.TaxScope,
		entity.IsPriceInclusive,
		entity.TaxAccountId,
		entity.TaxRefundAccountId,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create taxes", zap.Error(err))
		return fmt.Errorf("failed to create taxes: %w", err)
	}

	r.logger.Info("created taxes",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a taxes by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Taxes, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "taxes", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, tax_group_id
			, tax_code
			, tax_name
			, tax_rate
			, tax_scope
			, is_price_inclusive
			, tax_account_id
			, tax_refund_account_id
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM taxes
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Taxes
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.TaxGroupId,
		&entity.TaxCode,
		&entity.TaxName,
		&entity.TaxRate,
		&entity.TaxScope,
		&entity.IsPriceInclusive,
		&entity.TaxAccountId,
		&entity.TaxRefundAccountId,
		&entity.IsActive,
		&entity.Description,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("taxes not found")
	}

	if err != nil {
		r.logger.Error("failed to get taxes", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get taxes: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of taxes records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Taxes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "taxes", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM taxes
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count taxes records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tax_group_id
			, tax_code
			, tax_name
			, tax_rate
			, tax_scope
			, is_price_inclusive
			, tax_account_id
			, tax_refund_account_id
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM taxes
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list taxes", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list taxes: %w", err)
	}
	defer rows.Close()

	var entities []*Taxes
	for rows.Next() {
		var entity Taxes
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TaxGroupId,
			&entity.TaxCode,
			&entity.TaxName,
			&entity.TaxRate,
			&entity.TaxScope,
			&entity.IsPriceInclusive,
			&entity.TaxAccountId,
			&entity.TaxRefundAccountId,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan taxes: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating taxes rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing taxes record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Taxes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "taxes", duration, nil)
	}()

	query := `
		UPDATE taxes
		SET
			, organization_id = $2
			, tax_group_id = $3
			, tax_code = $4
			, tax_name = $5
			, tax_rate = $6
			, tax_scope = $7
			, is_price_inclusive = $8
			, tax_account_id = $9
			, tax_refund_account_id = $10
			, is_active = $11
			, description = $12
			, created_by = $13
			, updated_by = $14
			, updated_at = $16
			, deleted_at = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.TaxGroupId,
		entity.TaxCode,
		entity.TaxName,
		entity.TaxRate,
		entity.TaxScope,
		entity.IsPriceInclusive,
		entity.TaxAccountId,
		entity.TaxRefundAccountId,
		entity.IsActive,
		entity.Description,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update taxes", zap.Error(err))
		return fmt.Errorf("failed to update taxes: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("taxes not found or already deleted")
	}

	r.logger.Info("updated taxes",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a taxes record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "taxes", duration, nil)
	}()

	query := `
		UPDATE taxes
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete taxes", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete taxes: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("taxes not found or already deleted")
	}

	r.logger.Info("deleted taxes", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves taxes records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Taxes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "taxes", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM taxes
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count taxes records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, tax_group_id
			, tax_code
			, tax_name
			, tax_rate
			, tax_scope
			, is_price_inclusive
			, tax_account_id
			, tax_refund_account_id
			, is_active
			, description
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM taxes
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list taxes by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list taxes: %w", err)
	}
	defer rows.Close()

	var entities []*Taxes
	for rows.Next() {
		var entity Taxes
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.TaxGroupId,
			&entity.TaxCode,
			&entity.TaxName,
			&entity.TaxRate,
			&entity.TaxScope,
			&entity.IsPriceInclusive,
			&entity.TaxAccountId,
			&entity.TaxRefundAccountId,
			&entity.IsActive,
			&entity.Description,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan taxes: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

