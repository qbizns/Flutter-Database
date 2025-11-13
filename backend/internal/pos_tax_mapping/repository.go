package pos_tax_mapping

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

// Repository handles database operations for PosTaxMappings
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PosTaxMappings repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PosTaxMappings represents a pos_tax_mappings entity
type PosTaxMappings struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PosTaxCode *string `json:"pos_tax_code" db:"pos_tax_code"`
	TaxCategoryCode *string `json:"tax_category_code" db:"tax_category_code"`
	PosTaxRate *float64 `json:"pos_tax_rate" db:"pos_tax_rate"`
	AccountingTaxId *uuid.UUID `json:"accounting_tax_id" db:"accounting_tax_id"`
	DefaultTaxAccountId *uuid.UUID `json:"default_tax_account_id" db:"default_tax_account_id"`
	DefaultTaxExpenseAccountId *uuid.UUID `json:"default_tax_expense_account_id" db:"default_tax_expense_account_id"`
	IsDefault *bool `json:"is_default" db:"is_default"`
	IsActive *bool `json:"is_active" db:"is_active"`
	Priority *int64 `json:"priority" db:"priority"`
	IsInclusive *bool `json:"is_inclusive" db:"is_inclusive"`
	AppliesToSales *bool `json:"applies_to_sales" db:"applies_to_sales"`
	AppliesToPurchases *bool `json:"applies_to_purchases" db:"applies_to_purchases"`
	EffectiveFrom *time.Time `json:"effective_from" db:"effective_from"`
	EffectiveTo *time.Time `json:"effective_to" db:"effective_to"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	EffectiveFrom *string `json:"effective_from" db:"effective_from"`
}

// Create inserts a new pos_tax_mappings record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PosTaxMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "pos_tax_mappings", duration, nil)
	}()

	query := `
		INSERT INTO pos_tax_mappings (
			, organization_id
			, pos_tax_code
			, tax_category_code
			, pos_tax_rate
			, accounting_tax_id
			, default_tax_account_id
			, default_tax_expense_account_id
			, is_default
			, is_active
			, priority
			, is_inclusive
			, applies_to_sales
			, applies_to_purchases
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, effective_from
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
			, $22
			, $23
			, $24
			, $25
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PosTaxCode,
		entity.TaxCategoryCode,
		entity.PosTaxRate,
		entity.AccountingTaxId,
		entity.DefaultTaxAccountId,
		entity.DefaultTaxExpenseAccountId,
		entity.IsDefault,
		entity.IsActive,
		entity.Priority,
		entity.IsInclusive,
		entity.AppliesToSales,
		entity.AppliesToPurchases,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.EffectiveFrom,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create pos_tax_mappings", zap.Error(err))
		return fmt.Errorf("failed to create pos_tax_mappings: %w", err)
	}

	r.logger.Info("created pos_tax_mappings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a pos_tax_mappings by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PosTaxMappings, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_tax_mappings", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, pos_tax_code
			, tax_category_code
			, pos_tax_rate
			, accounting_tax_id
			, default_tax_account_id
			, default_tax_expense_account_id
			, is_default
			, is_active
			, priority
			, is_inclusive
			, applies_to_sales
			, applies_to_purchases
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, effective_from
		FROM pos_tax_mappings
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PosTaxMappings
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PosTaxCode,
		&entity.TaxCategoryCode,
		&entity.PosTaxRate,
		&entity.AccountingTaxId,
		&entity.DefaultTaxAccountId,
		&entity.DefaultTaxExpenseAccountId,
		&entity.IsDefault,
		&entity.IsActive,
		&entity.Priority,
		&entity.IsInclusive,
		&entity.AppliesToSales,
		&entity.AppliesToPurchases,
		&entity.EffectiveFrom,
		&entity.EffectiveTo,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.EffectiveFrom,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("pos_tax_mappings not found")
	}

	if err != nil {
		r.logger.Error("failed to get pos_tax_mappings", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get pos_tax_mappings: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of pos_tax_mappings records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PosTaxMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_tax_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_tax_mappings
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_tax_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, pos_tax_code
			, tax_category_code
			, pos_tax_rate
			, accounting_tax_id
			, default_tax_account_id
			, default_tax_expense_account_id
			, is_default
			, is_active
			, priority
			, is_inclusive
			, applies_to_sales
			, applies_to_purchases
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, effective_from
		FROM pos_tax_mappings
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_tax_mappings", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_tax_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*PosTaxMappings
	for rows.Next() {
		var entity PosTaxMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PosTaxCode,
			&entity.TaxCategoryCode,
			&entity.PosTaxRate,
			&entity.AccountingTaxId,
			&entity.DefaultTaxAccountId,
			&entity.DefaultTaxExpenseAccountId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Priority,
			&entity.IsInclusive,
			&entity.AppliesToSales,
			&entity.AppliesToPurchases,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.EffectiveFrom,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_tax_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating pos_tax_mappings rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing pos_tax_mappings record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PosTaxMappings) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_tax_mappings", duration, nil)
	}()

	query := `
		UPDATE pos_tax_mappings
		SET
			, organization_id = $2
			, pos_tax_code = $3
			, tax_category_code = $4
			, pos_tax_rate = $5
			, accounting_tax_id = $6
			, default_tax_account_id = $7
			, default_tax_expense_account_id = $8
			, is_default = $9
			, is_active = $10
			, priority = $11
			, is_inclusive = $12
			, applies_to_sales = $13
			, applies_to_purchases = $14
			, effective_from = $15
			, effective_to = $16
			, description = $17
			, notes = $18
			, metadata = $19
			, updated_at = $21
			, deleted_at = $22
			, created_by = $23
			, updated_by = $24
			, effective_from = $25
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $26
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PosTaxCode,
		entity.TaxCategoryCode,
		entity.PosTaxRate,
		entity.AccountingTaxId,
		entity.DefaultTaxAccountId,
		entity.DefaultTaxExpenseAccountId,
		entity.IsDefault,
		entity.IsActive,
		entity.Priority,
		entity.IsInclusive,
		entity.AppliesToSales,
		entity.AppliesToPurchases,
		entity.EffectiveFrom,
		entity.EffectiveTo,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.EffectiveFrom,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update pos_tax_mappings", zap.Error(err))
		return fmt.Errorf("failed to update pos_tax_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_tax_mappings not found or already deleted")
	}

	r.logger.Info("updated pos_tax_mappings",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a pos_tax_mappings record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_tax_mappings", duration, nil)
	}()

	query := `
		UPDATE pos_tax_mappings
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete pos_tax_mappings", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete pos_tax_mappings: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_tax_mappings not found or already deleted")
	}

	r.logger.Info("deleted pos_tax_mappings", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves pos_tax_mappings records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PosTaxMappings, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_tax_mappings", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_tax_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_tax_mappings records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, pos_tax_code
			, tax_category_code
			, pos_tax_rate
			, accounting_tax_id
			, default_tax_account_id
			, default_tax_expense_account_id
			, is_default
			, is_active
			, priority
			, is_inclusive
			, applies_to_sales
			, applies_to_purchases
			, effective_from
			, effective_to
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, effective_from
		FROM pos_tax_mappings
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_tax_mappings by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_tax_mappings: %w", err)
	}
	defer rows.Close()

	var entities []*PosTaxMappings
	for rows.Next() {
		var entity PosTaxMappings
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PosTaxCode,
			&entity.TaxCategoryCode,
			&entity.PosTaxRate,
			&entity.AccountingTaxId,
			&entity.DefaultTaxAccountId,
			&entity.DefaultTaxExpenseAccountId,
			&entity.IsDefault,
			&entity.IsActive,
			&entity.Priority,
			&entity.IsInclusive,
			&entity.AppliesToSales,
			&entity.AppliesToPurchases,
			&entity.EffectiveFrom,
			&entity.EffectiveTo,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.EffectiveFrom,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_tax_mappings: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

