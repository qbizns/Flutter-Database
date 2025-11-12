package reconciliation_rule_model

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

// Repository handles database operations for ReconciliationRuleModels
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new ReconciliationRuleModels repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// ReconciliationRuleModels represents a reconciliation_rule_models entity
type ReconciliationRuleModels struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	RuleName string `json:"rule_name" db:"rule_name"`
	RuleCode *string `json:"rule_code" db:"rule_code"`
	Sequence *int64 `json:"sequence" db:"sequence"`
	AmountMin *float64 `json:"amount_min" db:"amount_min"`
	AmountMax *float64 `json:"amount_max" db:"amount_max"`
	DescriptionPattern *string `json:"description_pattern" db:"description_pattern"`
	CounterpartyPattern *string `json:"counterparty_pattern" db:"counterparty_pattern"`
	ReferencePattern *string `json:"reference_pattern" db:"reference_pattern"`
	JournalId *uuid.UUID `json:"journal_id" db:"journal_id"`
	AccountId *uuid.UUID `json:"account_id" db:"account_id"`
	AnalyticAccountId *uuid.UUID `json:"analytic_account_id" db:"analytic_account_id"`
	TaxId *uuid.UUID `json:"tax_id" db:"tax_id"`
	IsActive *bool `json:"is_active" db:"is_active"`
	AutoApply *bool `json:"auto_apply" db:"auto_apply"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new reconciliation_rule_models record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *ReconciliationRuleModels) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "reconciliation_rule_models", duration, nil)
	}()

	query := `
		INSERT INTO reconciliation_rule_models (
			, organization_id
			, rule_name
			, rule_code
			, sequence
			, amount_min
			, amount_max
			, description_pattern
			, counterparty_pattern
			, reference_pattern
			, journal_id
			, account_id
			, analytic_account_id
			, tax_id
			, is_active
			, auto_apply
			, notes
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
			, $15
			, $16
			, $17
			, $18
			, $19
			, $22
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.RuleName,
		entity.RuleCode,
		entity.Sequence,
		entity.AmountMin,
		entity.AmountMax,
		entity.DescriptionPattern,
		entity.CounterpartyPattern,
		entity.ReferencePattern,
		entity.JournalId,
		entity.AccountId,
		entity.AnalyticAccountId,
		entity.TaxId,
		entity.IsActive,
		entity.AutoApply,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create reconciliation_rule_models", zap.Error(err))
		return fmt.Errorf("failed to create reconciliation_rule_models: %w", err)
	}

	r.logger.Info("created reconciliation_rule_models",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a reconciliation_rule_models by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*ReconciliationRuleModels, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reconciliation_rule_models", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, rule_name
			, rule_code
			, sequence
			, amount_min
			, amount_max
			, description_pattern
			, counterparty_pattern
			, reference_pattern
			, journal_id
			, account_id
			, analytic_account_id
			, tax_id
			, is_active
			, auto_apply
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM reconciliation_rule_models
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity ReconciliationRuleModels
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.RuleName,
		&entity.RuleCode,
		&entity.Sequence,
		&entity.AmountMin,
		&entity.AmountMax,
		&entity.DescriptionPattern,
		&entity.CounterpartyPattern,
		&entity.ReferencePattern,
		&entity.JournalId,
		&entity.AccountId,
		&entity.AnalyticAccountId,
		&entity.TaxId,
		&entity.IsActive,
		&entity.AutoApply,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("reconciliation_rule_models not found")
	}

	if err != nil {
		r.logger.Error("failed to get reconciliation_rule_models", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get reconciliation_rule_models: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of reconciliation_rule_models records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*ReconciliationRuleModels, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reconciliation_rule_models", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM reconciliation_rule_models
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reconciliation_rule_models records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, rule_name
			, rule_code
			, sequence
			, amount_min
			, amount_max
			, description_pattern
			, counterparty_pattern
			, reference_pattern
			, journal_id
			, account_id
			, analytic_account_id
			, tax_id
			, is_active
			, auto_apply
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM reconciliation_rule_models
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list reconciliation_rule_models", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list reconciliation_rule_models: %w", err)
	}
	defer rows.Close()

	var entities []*ReconciliationRuleModels
	for rows.Next() {
		var entity ReconciliationRuleModels
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RuleName,
			&entity.RuleCode,
			&entity.Sequence,
			&entity.AmountMin,
			&entity.AmountMax,
			&entity.DescriptionPattern,
			&entity.CounterpartyPattern,
			&entity.ReferencePattern,
			&entity.JournalId,
			&entity.AccountId,
			&entity.AnalyticAccountId,
			&entity.TaxId,
			&entity.IsActive,
			&entity.AutoApply,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reconciliation_rule_models: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating reconciliation_rule_models rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing reconciliation_rule_models record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *ReconciliationRuleModels) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "reconciliation_rule_models", duration, nil)
	}()

	query := `
		UPDATE reconciliation_rule_models
		SET
			, organization_id = $2
			, rule_name = $3
			, rule_code = $4
			, sequence = $5
			, amount_min = $6
			, amount_max = $7
			, description_pattern = $8
			, counterparty_pattern = $9
			, reference_pattern = $10
			, journal_id = $11
			, account_id = $12
			, analytic_account_id = $13
			, tax_id = $14
			, is_active = $15
			, auto_apply = $16
			, notes = $17
			, created_by = $18
			, updated_by = $19
			, updated_at = $21
			, deleted_at = $22
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $23
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.RuleName,
		entity.RuleCode,
		entity.Sequence,
		entity.AmountMin,
		entity.AmountMax,
		entity.DescriptionPattern,
		entity.CounterpartyPattern,
		entity.ReferencePattern,
		entity.JournalId,
		entity.AccountId,
		entity.AnalyticAccountId,
		entity.TaxId,
		entity.IsActive,
		entity.AutoApply,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update reconciliation_rule_models", zap.Error(err))
		return fmt.Errorf("failed to update reconciliation_rule_models: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reconciliation_rule_models not found or already deleted")
	}

	r.logger.Info("updated reconciliation_rule_models",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a reconciliation_rule_models record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "reconciliation_rule_models", duration, nil)
	}()

	query := `
		UPDATE reconciliation_rule_models
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete reconciliation_rule_models", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete reconciliation_rule_models: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("reconciliation_rule_models not found or already deleted")
	}

	r.logger.Info("deleted reconciliation_rule_models", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves reconciliation_rule_models records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*ReconciliationRuleModels, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "reconciliation_rule_models", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM reconciliation_rule_models
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reconciliation_rule_models records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, rule_name
			, rule_code
			, sequence
			, amount_min
			, amount_max
			, description_pattern
			, counterparty_pattern
			, reference_pattern
			, journal_id
			, account_id
			, analytic_account_id
			, tax_id
			, is_active
			, auto_apply
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM reconciliation_rule_models
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list reconciliation_rule_models by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list reconciliation_rule_models: %w", err)
	}
	defer rows.Close()

	var entities []*ReconciliationRuleModels
	for rows.Next() {
		var entity ReconciliationRuleModels
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.RuleName,
			&entity.RuleCode,
			&entity.Sequence,
			&entity.AmountMin,
			&entity.AmountMax,
			&entity.DescriptionPattern,
			&entity.CounterpartyPattern,
			&entity.ReferencePattern,
			&entity.JournalId,
			&entity.AccountId,
			&entity.AnalyticAccountId,
			&entity.TaxId,
			&entity.IsActive,
			&entity.AutoApply,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan reconciliation_rule_models: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

