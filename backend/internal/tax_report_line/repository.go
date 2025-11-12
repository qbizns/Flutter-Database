package tax_report_line

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

// Repository handles database operations for TaxReportLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TaxReportLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TaxReportLines represents a tax_report_lines entity
type TaxReportLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	TaxReportDefinitionId uuid.UUID `json:"tax_report_definition_id" db:"tax_report_definition_id"`
	LineCode string `json:"line_code" db:"line_code"`
	LineName string `json:"line_name" db:"line_name"`
	Sequence *int64 `json:"sequence" db:"sequence"`
	ParentLineId *uuid.UUID `json:"parent_line_id" db:"parent_line_id"`
	FormulaType *string `json:"formula_type" db:"formula_type"`
	Formula *string `json:"formula" db:"formula"`
	TaxGroupIds *uuid.UUID `json:"tax_group_ids" db:"tax_group_ids"`
	AccountIds *uuid.UUID `json:"account_ids" db:"account_ids"`
	TaxIds *uuid.UUID `json:"tax_ids" db:"tax_ids"`
	IsSubtotal *bool `json:"is_subtotal" db:"is_subtotal"`
	IsTotal *bool `json:"is_total" db:"is_total"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new tax_report_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TaxReportLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "tax_report_lines", duration, nil)
	}()

	query := `
		INSERT INTO tax_report_lines (
			, tax_report_definition_id
			, line_code
			, line_name
			, sequence
			, parent_line_id
			, formula_type
			, formula
			, tax_group_ids
			, account_ids
			, tax_ids
			, is_subtotal
			, is_total
			, notes
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
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.TaxReportDefinitionId,
		entity.LineCode,
		entity.LineName,
		entity.Sequence,
		entity.ParentLineId,
		entity.FormulaType,
		entity.Formula,
		entity.TaxGroupIds,
		entity.AccountIds,
		entity.TaxIds,
		entity.IsSubtotal,
		entity.IsTotal,
		entity.Notes,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create tax_report_lines", zap.Error(err))
		return fmt.Errorf("failed to create tax_report_lines: %w", err)
	}

	r.logger.Info("created tax_report_lines",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a tax_report_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TaxReportLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_report_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, tax_report_definition_id
			, line_code
			, line_name
			, sequence
			, parent_line_id
			, formula_type
			, formula
			, tax_group_ids
			, account_ids
			, tax_ids
			, is_subtotal
			, is_total
			, notes
			, created_at
			, deleted_at
		FROM tax_report_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TaxReportLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.TaxReportDefinitionId,
		&entity.LineCode,
		&entity.LineName,
		&entity.Sequence,
		&entity.ParentLineId,
		&entity.FormulaType,
		&entity.Formula,
		&entity.TaxGroupIds,
		&entity.AccountIds,
		&entity.TaxIds,
		&entity.IsSubtotal,
		&entity.IsTotal,
		&entity.Notes,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tax_report_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get tax_report_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get tax_report_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of tax_report_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TaxReportLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tax_report_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tax_report_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax_report_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, tax_report_definition_id
			, line_code
			, line_name
			, sequence
			, parent_line_id
			, formula_type
			, formula
			, tax_group_ids
			, account_ids
			, tax_ids
			, is_subtotal
			, is_total
			, notes
			, created_at
			, deleted_at
		FROM tax_report_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tax_report_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tax_report_lines: %w", err)
	}
	defer rows.Close()

	var entities []*TaxReportLines
	for rows.Next() {
		var entity TaxReportLines
		err := rows.Scan(
			&entity.Id,
			&entity.TaxReportDefinitionId,
			&entity.LineCode,
			&entity.LineName,
			&entity.Sequence,
			&entity.ParentLineId,
			&entity.FormulaType,
			&entity.Formula,
			&entity.TaxGroupIds,
			&entity.AccountIds,
			&entity.TaxIds,
			&entity.IsSubtotal,
			&entity.IsTotal,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tax_report_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating tax_report_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing tax_report_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TaxReportLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_report_lines", duration, nil)
	}()

	query := `
		UPDATE tax_report_lines
		SET
			, tax_report_definition_id = $2
			, line_code = $3
			, line_name = $4
			, sequence = $5
			, parent_line_id = $6
			, formula_type = $7
			, formula = $8
			, tax_group_ids = $9
			, account_ids = $10
			, tax_ids = $11
			, is_subtotal = $12
			, is_total = $13
			, notes = $14
			, deleted_at = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.TaxReportDefinitionId,
		entity.LineCode,
		entity.LineName,
		entity.Sequence,
		entity.ParentLineId,
		entity.FormulaType,
		entity.Formula,
		entity.TaxGroupIds,
		entity.AccountIds,
		entity.TaxIds,
		entity.IsSubtotal,
		entity.IsTotal,
		entity.Notes,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update tax_report_lines", zap.Error(err))
		return fmt.Errorf("failed to update tax_report_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_report_lines not found or already deleted")
	}

	r.logger.Info("updated tax_report_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a tax_report_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tax_report_lines", duration, nil)
	}()

	query := `
		UPDATE tax_report_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete tax_report_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete tax_report_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tax_report_lines not found or already deleted")
	}

	r.logger.Info("deleted tax_report_lines", zap.String("id", id.String()))
	return nil
}



