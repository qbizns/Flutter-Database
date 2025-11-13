package asset_depreciation_schedule

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

// Repository handles database operations for AssetDepreciationSchedule
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AssetDepreciationSchedule repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AssetDepreciationSchedule represents a asset_depreciation_schedule entity
type AssetDepreciationSchedule struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	FixedAssetId uuid.UUID `json:"fixed_asset_id" db:"fixed_asset_id"`
	FiscalYearId *uuid.UUID `json:"fiscal_year_id" db:"fiscal_year_id"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	DepreciationDate time.Time `json:"depreciation_date" db:"depreciation_date"`
	DepreciationAmount float64 `json:"depreciation_amount" db:"depreciation_amount"`
	AccumulatedDepreciationBeginning float64 `json:"accumulated_depreciation_beginning" db:"accumulated_depreciation_beginning"`
	AccumulatedDepreciationEnding float64 `json:"accumulated_depreciation_ending" db:"accumulated_depreciation_ending"`
	BookValueBeginning float64 `json:"book_value_beginning" db:"book_value_beginning"`
	BookValueEnding float64 `json:"book_value_ending" db:"book_value_ending"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	IsPosted *bool `json:"is_posted" db:"is_posted"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	PostedAt *time.Time `json:"posted_at" db:"posted_at"`
	PostedBy *uuid.UUID `json:"posted_by" db:"posted_by"`
}

// Create inserts a new asset_depreciation_schedule record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AssetDepreciationSchedule) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "asset_depreciation_schedule", duration, nil)
	}()

	query := `
		INSERT INTO asset_depreciation_schedule (
			, organization_id
			, fixed_asset_id
			, fiscal_year_id
			, accounting_period_id
			, depreciation_date
			, depreciation_amount
			, accumulated_depreciation_beginning
			, accumulated_depreciation_ending
			, book_value_beginning
			, book_value_ending
			, journal_entry_id
			, is_posted
			, posted_at
			, posted_by
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
			, $16
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.FixedAssetId,
		entity.FiscalYearId,
		entity.AccountingPeriodId,
		entity.DepreciationDate,
		entity.DepreciationAmount,
		entity.AccumulatedDepreciationBeginning,
		entity.AccumulatedDepreciationEnding,
		entity.BookValueBeginning,
		entity.BookValueEnding,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.PostedAt,
		entity.PostedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create asset_depreciation_schedule", zap.Error(err))
		return fmt.Errorf("failed to create asset_depreciation_schedule: %w", err)
	}

	r.logger.Info("created asset_depreciation_schedule",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a asset_depreciation_schedule by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AssetDepreciationSchedule, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "asset_depreciation_schedule", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, fixed_asset_id
			, fiscal_year_id
			, accounting_period_id
			, depreciation_date
			, depreciation_amount
			, accumulated_depreciation_beginning
			, accumulated_depreciation_ending
			, book_value_beginning
			, book_value_ending
			, journal_entry_id
			, is_posted
			, created_at
			, posted_at
			, posted_by
		FROM asset_depreciation_schedule
		WHERE id = $1
		
	`

	var entity AssetDepreciationSchedule
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.FixedAssetId,
		&entity.FiscalYearId,
		&entity.AccountingPeriodId,
		&entity.DepreciationDate,
		&entity.DepreciationAmount,
		&entity.AccumulatedDepreciationBeginning,
		&entity.AccumulatedDepreciationEnding,
		&entity.BookValueBeginning,
		&entity.BookValueEnding,
		&entity.JournalEntryId,
		&entity.IsPosted,
		&entity.CreatedAt,
		&entity.PostedAt,
		&entity.PostedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("asset_depreciation_schedule not found")
	}

	if err != nil {
		r.logger.Error("failed to get asset_depreciation_schedule", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get asset_depreciation_schedule: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of asset_depreciation_schedule records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AssetDepreciationSchedule, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "asset_depreciation_schedule", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM asset_depreciation_schedule
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count asset_depreciation_schedule records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fixed_asset_id
			, fiscal_year_id
			, accounting_period_id
			, depreciation_date
			, depreciation_amount
			, accumulated_depreciation_beginning
			, accumulated_depreciation_ending
			, book_value_beginning
			, book_value_ending
			, journal_entry_id
			, is_posted
			, created_at
			, posted_at
			, posted_by
		FROM asset_depreciation_schedule
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list asset_depreciation_schedule", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list asset_depreciation_schedule: %w", err)
	}
	defer rows.Close()

	var entities []*AssetDepreciationSchedule
	for rows.Next() {
		var entity AssetDepreciationSchedule
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FixedAssetId,
			&entity.FiscalYearId,
			&entity.AccountingPeriodId,
			&entity.DepreciationDate,
			&entity.DepreciationAmount,
			&entity.AccumulatedDepreciationBeginning,
			&entity.AccumulatedDepreciationEnding,
			&entity.BookValueBeginning,
			&entity.BookValueEnding,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.CreatedAt,
			&entity.PostedAt,
			&entity.PostedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset_depreciation_schedule: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating asset_depreciation_schedule rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing asset_depreciation_schedule record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AssetDepreciationSchedule) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "asset_depreciation_schedule", duration, nil)
	}()

	query := `
		UPDATE asset_depreciation_schedule
		SET
			, organization_id = $2
			, fixed_asset_id = $3
			, fiscal_year_id = $4
			, accounting_period_id = $5
			, depreciation_date = $6
			, depreciation_amount = $7
			, accumulated_depreciation_beginning = $8
			, accumulated_depreciation_ending = $9
			, book_value_beginning = $10
			, book_value_ending = $11
			, journal_entry_id = $12
			, is_posted = $13
			, posted_at = $15
			, posted_by = $16
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $17
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.FixedAssetId,
		entity.FiscalYearId,
		entity.AccountingPeriodId,
		entity.DepreciationDate,
		entity.DepreciationAmount,
		entity.AccumulatedDepreciationBeginning,
		entity.AccumulatedDepreciationEnding,
		entity.BookValueBeginning,
		entity.BookValueEnding,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.PostedAt,
		entity.PostedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update asset_depreciation_schedule", zap.Error(err))
		return fmt.Errorf("failed to update asset_depreciation_schedule: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("asset_depreciation_schedule not found or already deleted")
	}

	r.logger.Info("updated asset_depreciation_schedule",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a asset_depreciation_schedule record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "asset_depreciation_schedule", duration, nil)
	}()

	query := `DELETE FROM asset_depreciation_schedule WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete asset_depreciation_schedule", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete asset_depreciation_schedule: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("asset_depreciation_schedule not found")
	}

	r.logger.Info("deleted asset_depreciation_schedule", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves asset_depreciation_schedule records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*AssetDepreciationSchedule, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "asset_depreciation_schedule", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM asset_depreciation_schedule
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count asset_depreciation_schedule records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, fixed_asset_id
			, fiscal_year_id
			, accounting_period_id
			, depreciation_date
			, depreciation_amount
			, accumulated_depreciation_beginning
			, accumulated_depreciation_ending
			, book_value_beginning
			, book_value_ending
			, journal_entry_id
			, is_posted
			, created_at
			, posted_at
			, posted_by
		FROM asset_depreciation_schedule
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list asset_depreciation_schedule by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list asset_depreciation_schedule: %w", err)
	}
	defer rows.Close()

	var entities []*AssetDepreciationSchedule
	for rows.Next() {
		var entity AssetDepreciationSchedule
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.FixedAssetId,
			&entity.FiscalYearId,
			&entity.AccountingPeriodId,
			&entity.DepreciationDate,
			&entity.DepreciationAmount,
			&entity.AccumulatedDepreciationBeginning,
			&entity.AccumulatedDepreciationEnding,
			&entity.BookValueBeginning,
			&entity.BookValueEnding,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.CreatedAt,
			&entity.PostedAt,
			&entity.PostedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset_depreciation_schedule: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

