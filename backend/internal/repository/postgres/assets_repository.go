package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/domain/assets"
)

// AssetsRepository implements assets.Repository for PostgreSQL
type AssetsRepository struct {
	pool *pgxpool.Pool
}

// NewAssetsRepository creates a new assets repository
func NewAssetsRepository(pool *pgxpool.Pool) *AssetsRepository {
	return &AssetsRepository{pool: pool}
}

// ==================== ASSET CATEGORIES ====================

// CreateAssetCategory creates a new asset category
func (r *AssetsRepository) CreateAssetCategory(ctx context.Context, category *assets.AssetCategory) error {
	query := `
		INSERT INTO asset_categories (
			id, category_code, category_name,
			default_depreciation_method, default_useful_life_years,
			default_salvage_value_percent, asset_account_id,
			accumulated_depreciation_account_id, depreciation_expense_account_id,
			description, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.pool.Exec(ctx, query,
		category.ID, category.CategoryCode, category.CategoryName,
		category.DefaultDepreciationMethod, category.DefaultUsefulLifeYears,
		category.DefaultSalvageValuePercent, category.AssetAccountID,
		category.AccumulatedDepreciationAccountID, category.DepreciationExpenseAccountID,
		category.Description, category.CreatedAt, category.UpdatedAt,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return assets.ErrDuplicateAssetNumber
		}
		return fmt.Errorf("failed to create asset category: %w", err)
	}

	return nil
}

// GetAssetCategory retrieves an asset category by ID
func (r *AssetsRepository) GetAssetCategory(ctx context.Context, id uuid.UUID) (*assets.AssetCategory, error) {
	query := `
		SELECT id, category_code, category_name,
		       default_depreciation_method, default_useful_life_years,
		       default_salvage_value_percent, asset_account_id,
		       accumulated_depreciation_account_id, depreciation_expense_account_id,
		       description, created_at, updated_at
		FROM asset_categories
		WHERE id = $1
	`

	var category assets.AssetCategory
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&category.ID, &category.CategoryCode, &category.CategoryName,
		&category.DefaultDepreciationMethod, &category.DefaultUsefulLifeYears,
		&category.DefaultSalvageValuePercent, &category.AssetAccountID,
		&category.AccumulatedDepreciationAccountID, &category.DepreciationExpenseAccountID,
		&category.Description, &category.CreatedAt, &category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset category: %w", err)
	}

	return &category, nil
}

// GetAssetCategoryByCode retrieves an asset category by code
func (r *AssetsRepository) GetAssetCategoryByCode(ctx context.Context, code string) (*assets.AssetCategory, error) {
	query := `
		SELECT id, category_code, category_name,
		       default_depreciation_method, default_useful_life_years,
		       default_salvage_value_percent, asset_account_id,
		       accumulated_depreciation_account_id, depreciation_expense_account_id,
		       description, created_at, updated_at
		FROM asset_categories
		WHERE category_code = $1
	`

	var category assets.AssetCategory
	err := r.pool.QueryRow(ctx, query, code).Scan(
		&category.ID, &category.CategoryCode, &category.CategoryName,
		&category.DefaultDepreciationMethod, &category.DefaultUsefulLifeYears,
		&category.DefaultSalvageValuePercent, &category.AssetAccountID,
		&category.AccumulatedDepreciationAccountID, &category.DepreciationExpenseAccountID,
		&category.Description, &category.CreatedAt, &category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get asset category by code: %w", err)
	}

	return &category, nil
}

// ListAssetCategories retrieves all asset categories
func (r *AssetsRepository) ListAssetCategories(ctx context.Context, organizationID *uuid.UUID) ([]*assets.AssetCategory, error) {
	query := `
		SELECT id, category_code, category_name,
		       default_depreciation_method, default_useful_life_years,
		       default_salvage_value_percent, asset_account_id,
		       accumulated_depreciation_account_id, depreciation_expense_account_id,
		       description, created_at, updated_at
		FROM asset_categories
		ORDER BY category_code ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list asset categories: %w", err)
	}
	defer rows.Close()

	var categories []*assets.AssetCategory
	for rows.Next() {
		var category assets.AssetCategory
		err := rows.Scan(
			&category.ID, &category.CategoryCode, &category.CategoryName,
			&category.DefaultDepreciationMethod, &category.DefaultUsefulLifeYears,
			&category.DefaultSalvageValuePercent, &category.AssetAccountID,
			&category.AccumulatedDepreciationAccountID, &category.DepreciationExpenseAccountID,
			&category.Description, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, rows.Err()
}

// UpdateAssetCategory updates an asset category
func (r *AssetsRepository) UpdateAssetCategory(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE asset_categories SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += fmt.Sprintf(" WHERE id = $1")

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteAssetCategory deletes an asset category
func (r *AssetsRepository) DeleteAssetCategory(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM asset_categories WHERE id = $1", id)
	return err
}

// ==================== FIXED ASSETS ====================

// CreateFixedAsset creates a new fixed asset
func (r *AssetsRepository) CreateFixedAsset(ctx context.Context, asset *assets.FixedAsset) error {
	query := `
		INSERT INTO fixed_assets (
			id, organization_id, asset_number, asset_name, asset_category_id,
			acquisition_date, acquisition_cost, salvage_value, supplier_id,
			depreciation_method, useful_life_years, depreciation_start_date,
			asset_account_id, accumulated_depreciation_account_id,
			depreciation_expense_account_id, current_book_value,
			accumulated_depreciation, location_id, department,
			is_disposed, description, serial_number, notes, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27, $28
		)
	`

	_, err := r.pool.Exec(ctx, query,
		asset.ID, asset.OrganizationID, asset.AssetNumber, asset.AssetName,
		asset.AssetCategoryID, asset.AcquisitionDate, asset.AcquisitionCost,
		asset.SalvageValue, asset.SupplierID, asset.DepreciationMethod,
		asset.UsefulLifeYears, asset.DepreciationStartDate,
		asset.AssetAccountID, asset.AccumulatedDepreciationAccountID,
		asset.DepreciationExpenseAccountID, asset.CurrentBookValue,
		asset.AccumulatedDepreciation, asset.LocationID, asset.Department,
		asset.IsDisposed, asset.Description, asset.SerialNumber, asset.Notes,
		asset.Metadata, asset.CreatedAt, asset.UpdatedAt, asset.CreatedBy,
		asset.UpdatedBy,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return assets.ErrDuplicateAssetNumber
		}
		return fmt.Errorf("failed to create fixed asset: %w", err)
	}

	return nil
}

// GetFixedAsset retrieves a fixed asset by ID
func (r *AssetsRepository) GetFixedAsset(ctx context.Context, id uuid.UUID) (*assets.FixedAsset, error) {
	query := `
		SELECT id, organization_id, asset_number, asset_name, asset_category_id,
		       acquisition_date, acquisition_cost, salvage_value, supplier_id,
		       depreciation_method, useful_life_years, depreciation_start_date,
		       asset_account_id, accumulated_depreciation_account_id,
		       depreciation_expense_account_id, current_book_value,
		       accumulated_depreciation, last_depreciation_date, location_id,
		       department, is_disposed, disposal_date, disposal_proceeds,
		       disposal_journal_entry_id, description, serial_number, notes,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM fixed_assets
		WHERE id = $1 AND deleted_at IS NULL
	`

	var asset assets.FixedAsset
	var metadata sql.NullString
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&asset.ID, &asset.OrganizationID, &asset.AssetNumber, &asset.AssetName,
		&asset.AssetCategoryID, &asset.AcquisitionDate, &asset.AcquisitionCost,
		&asset.SalvageValue, &asset.SupplierID, &asset.DepreciationMethod,
		&asset.UsefulLifeYears, &asset.DepreciationStartDate,
		&asset.AssetAccountID, &asset.AccumulatedDepreciationAccountID,
		&asset.DepreciationExpenseAccountID, &asset.CurrentBookValue,
		&asset.AccumulatedDepreciation, &asset.LastDepreciationDate,
		&asset.LocationID, &asset.Department, &asset.IsDisposed,
		&asset.DisposalDate, &asset.DisposalProceeds,
		&asset.DisposalJournalEntryID, &asset.Description, &asset.SerialNumber,
		&asset.Notes, &metadata, &asset.CreatedAt, &asset.UpdatedAt,
		&asset.CreatedBy, &asset.UpdatedBy, &asset.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed asset: %w", err)
	}

	if metadata.Valid {
		asset.Metadata = json.RawMessage(metadata.String)
	}

	return &asset, nil
}

// GetFixedAssetByNumber retrieves a fixed asset by asset number
func (r *AssetsRepository) GetFixedAssetByNumber(ctx context.Context, organizationID uuid.UUID, assetNumber string) (*assets.FixedAsset, error) {
	query := `
		SELECT id, organization_id, asset_number, asset_name, asset_category_id,
		       acquisition_date, acquisition_cost, salvage_value, supplier_id,
		       depreciation_method, useful_life_years, depreciation_start_date,
		       asset_account_id, accumulated_depreciation_account_id,
		       depreciation_expense_account_id, current_book_value,
		       accumulated_depreciation, last_depreciation_date, location_id,
		       department, is_disposed, disposal_date, disposal_proceeds,
		       disposal_journal_entry_id, description, serial_number, notes,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM fixed_assets
		WHERE organization_id = $1 AND asset_number = $2 AND deleted_at IS NULL
	`

	var asset assets.FixedAsset
	var metadata sql.NullString
	err := r.pool.QueryRow(ctx, query, organizationID, assetNumber).Scan(
		&asset.ID, &asset.OrganizationID, &asset.AssetNumber, &asset.AssetName,
		&asset.AssetCategoryID, &asset.AcquisitionDate, &asset.AcquisitionCost,
		&asset.SalvageValue, &asset.SupplierID, &asset.DepreciationMethod,
		&asset.UsefulLifeYears, &asset.DepreciationStartDate,
		&asset.AssetAccountID, &asset.AccumulatedDepreciationAccountID,
		&asset.DepreciationExpenseAccountID, &asset.CurrentBookValue,
		&asset.AccumulatedDepreciation, &asset.LastDepreciationDate,
		&asset.LocationID, &asset.Department, &asset.IsDisposed,
		&asset.DisposalDate, &asset.DisposalProceeds,
		&asset.DisposalJournalEntryID, &asset.Description, &asset.SerialNumber,
		&asset.Notes, &metadata, &asset.CreatedAt, &asset.UpdatedAt,
		&asset.CreatedBy, &asset.UpdatedBy, &asset.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get fixed asset by number: %w", err)
	}

	if metadata.Valid {
		asset.Metadata = json.RawMessage(metadata.String)
	}

	return &asset, nil
}

// ListFixedAssets retrieves fixed assets based on query
func (r *AssetsRepository) ListFixedAssets(ctx context.Context, query assets.AssetQuery) ([]*assets.FixedAsset, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.AssetCategoryID != nil {
		where = append(where, fmt.Sprintf("asset_category_id = $%d", idx))
		args = append(args, *query.AssetCategoryID)
		idx++
	}

	if query.LocationID != nil {
		where = append(where, fmt.Sprintf("location_id = $%d", idx))
		args = append(args, *query.LocationID)
		idx++
	}

	if query.IsDisposed != nil {
		where = append(where, fmt.Sprintf("is_disposed = $%d", idx))
		args = append(args, *query.IsDisposed)
		idx++
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM fixed_assets WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fixed assets: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, asset_number, asset_name, asset_category_id,
		       acquisition_date, acquisition_cost, salvage_value, supplier_id,
		       depreciation_method, useful_life_years, depreciation_start_date,
		       asset_account_id, accumulated_depreciation_account_id,
		       depreciation_expense_account_id, current_book_value,
		       accumulated_depreciation, last_depreciation_date, location_id,
		       department, is_disposed, disposal_date, disposal_proceeds,
		       disposal_journal_entry_id, description, serial_number, notes,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM fixed_assets
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY asset_number ASC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list fixed assets: %w", err)
	}
	defer rows.Close()

	var assetsList []*assets.FixedAsset
	for rows.Next() {
		var asset assets.FixedAsset
		var metadata sql.NullString
		err := rows.Scan(
			&asset.ID, &asset.OrganizationID, &asset.AssetNumber, &asset.AssetName,
			&asset.AssetCategoryID, &asset.AcquisitionDate, &asset.AcquisitionCost,
			&asset.SalvageValue, &asset.SupplierID, &asset.DepreciationMethod,
			&asset.UsefulLifeYears, &asset.DepreciationStartDate,
			&asset.AssetAccountID, &asset.AccumulatedDepreciationAccountID,
			&asset.DepreciationExpenseAccountID, &asset.CurrentBookValue,
			&asset.AccumulatedDepreciation, &asset.LastDepreciationDate,
			&asset.LocationID, &asset.Department, &asset.IsDisposed,
			&asset.DisposalDate, &asset.DisposalProceeds,
			&asset.DisposalJournalEntryID, &asset.Description, &asset.SerialNumber,
			&asset.Notes, &metadata, &asset.CreatedAt, &asset.UpdatedAt,
			&asset.CreatedBy, &asset.UpdatedBy, &asset.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fixed asset: %w", err)
		}

		if metadata.Valid {
			asset.Metadata = json.RawMessage(metadata.String)
		}

		assetsList = append(assetsList, &asset)
	}

	return assetsList, total, rows.Err()
}

// UpdateFixedAsset updates a fixed asset
func (r *AssetsRepository) UpdateFixedAsset(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE fixed_assets SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteFixedAsset soft deletes a fixed asset
func (r *AssetsRepository) DeleteFixedAsset(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE fixed_assets SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// ==================== DEPRECIATION SCHEDULES ====================

// CreateDepreciationSchedule creates a depreciation schedule entry
func (r *AssetsRepository) CreateDepreciationSchedule(ctx context.Context, schedule *assets.AssetDepreciationSchedule) error {
	query := `
		INSERT INTO asset_depreciation_schedule (
			id, organization_id, fixed_asset_id, fiscal_year_id,
			accounting_period_id, depreciation_date, depreciation_amount,
			accumulated_depreciation_beginning, accumulated_depreciation_ending,
			book_value_beginning, book_value_ending, is_posted, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		schedule.ID, schedule.OrganizationID, schedule.FixedAssetID,
		schedule.FiscalYearID, schedule.AccountingPeriodID,
		schedule.DepreciationDate, schedule.DepreciationAmount,
		schedule.AccumulatedDepreciationBeginning,
		schedule.AccumulatedDepreciationEnding,
		schedule.BookValueBeginning, schedule.BookValueEnding,
		schedule.IsPosted, schedule.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create depreciation schedule: %w", err)
	}

	return nil
}

// GetDepreciationSchedule retrieves a depreciation schedule
func (r *AssetsRepository) GetDepreciationSchedule(ctx context.Context, id uuid.UUID) (*assets.AssetDepreciationSchedule, error) {
	query := `
		SELECT id, organization_id, fixed_asset_id, fiscal_year_id,
		       accounting_period_id, depreciation_date, depreciation_amount,
		       accumulated_depreciation_beginning, accumulated_depreciation_ending,
		       book_value_beginning, book_value_ending, journal_entry_id,
		       is_posted, created_at, posted_at, posted_by, deleted_at
		FROM asset_depreciation_schedule
		WHERE id = $1 AND deleted_at IS NULL
	`

	var schedule assets.AssetDepreciationSchedule
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&schedule.ID, &schedule.OrganizationID, &schedule.FixedAssetID,
		&schedule.FiscalYearID, &schedule.AccountingPeriodID,
		&schedule.DepreciationDate, &schedule.DepreciationAmount,
		&schedule.AccumulatedDepreciationBeginning,
		&schedule.AccumulatedDepreciationEnding,
		&schedule.BookValueBeginning, &schedule.BookValueEnding,
		&schedule.JournalEntryID, &schedule.IsPosted,
		&schedule.CreatedAt, &schedule.PostedAt, &schedule.PostedBy,
		&schedule.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get depreciation schedule: %w", err)
	}

	return &schedule, nil
}

// ListDepreciationSchedules retrieves depreciation schedules
func (r *AssetsRepository) ListDepreciationSchedules(ctx context.Context, query assets.DepreciationQuery) ([]*assets.AssetDepreciationSchedule, int64, error) {
	where := []string{"deleted_at IS NULL", "organization_id = $1"}
	args := []interface{}{query.OrganizationID}
	idx := 2

	if query.FixedAssetID != nil {
		where = append(where, fmt.Sprintf("fixed_asset_id = $%d", idx))
		args = append(args, *query.FixedAssetID)
		idx++
	}

	if query.IsPosted != nil {
		where = append(where, fmt.Sprintf("is_posted = $%d", idx))
		args = append(args, *query.IsPosted)
		idx++
	}

	if query.DateFrom != nil {
		where = append(where, fmt.Sprintf("depreciation_date >= $%d", idx))
		args = append(args, *query.DateFrom)
		idx++
	}

	if query.DateTo != nil {
		where = append(where, fmt.Sprintf("depreciation_date <= $%d", idx))
		args = append(args, *query.DateTo)
		idx++
	}

	// Get count
	countQuery := "SELECT COUNT(*) FROM asset_depreciation_schedule WHERE " + strings.Join(where, " AND ")
	var total int64
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count depreciation schedules: %w", err)
	}

	// Get records
	listQuery := `
		SELECT id, organization_id, fixed_asset_id, fiscal_year_id,
		       accounting_period_id, depreciation_date, depreciation_amount,
		       accumulated_depreciation_beginning, accumulated_depreciation_ending,
		       book_value_beginning, book_value_ending, journal_entry_id,
		       is_posted, created_at, posted_at, posted_by, deleted_at
		FROM asset_depreciation_schedule
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY depreciation_date DESC LIMIT $` + fmt.Sprintf("%d", idx) +
		` OFFSET $` + fmt.Sprintf("%d", idx+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list depreciation schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*assets.AssetDepreciationSchedule
	for rows.Next() {
		var schedule assets.AssetDepreciationSchedule
		err := rows.Scan(
			&schedule.ID, &schedule.OrganizationID, &schedule.FixedAssetID,
			&schedule.FiscalYearID, &schedule.AccountingPeriodID,
			&schedule.DepreciationDate, &schedule.DepreciationAmount,
			&schedule.AccumulatedDepreciationBeginning,
			&schedule.AccumulatedDepreciationEnding,
			&schedule.BookValueBeginning, &schedule.BookValueEnding,
			&schedule.JournalEntryID, &schedule.IsPosted,
			&schedule.CreatedAt, &schedule.PostedAt, &schedule.PostedBy,
			&schedule.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan depreciation schedule: %w", err)
		}
		schedules = append(schedules, &schedule)
	}

	return schedules, total, rows.Err()
}

// UpdateDepreciationSchedule updates a depreciation schedule
func (r *AssetsRepository) UpdateDepreciationSchedule(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	query := "UPDATE asset_depreciation_schedule SET "
	args := []interface{}{id}
	idx := 2

	columns := []string{}
	for key := range updates {
		columns = append(columns, key)
	}

	for i, col := range columns {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", col, idx)
		args = append(args, updates[col])
		idx++
	}

	query += " WHERE id = $1"

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

// DeleteDepreciationSchedule soft deletes a depreciation schedule
func (r *AssetsRepository) DeleteDepreciationSchedule(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "UPDATE asset_depreciation_schedule SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// GetDepreciationScheduleByAssetAndPeriod retrieves schedule by asset and period
func (r *AssetsRepository) GetDepreciationScheduleByAssetAndPeriod(ctx context.Context, assetID, periodID uuid.UUID) (*assets.AssetDepreciationSchedule, error) {
	query := `
		SELECT id, organization_id, fixed_asset_id, fiscal_year_id,
		       accounting_period_id, depreciation_date, depreciation_amount,
		       accumulated_depreciation_beginning, accumulated_depreciation_ending,
		       book_value_beginning, book_value_ending, journal_entry_id,
		       is_posted, created_at, posted_at, posted_by, deleted_at
		FROM asset_depreciation_schedule
		WHERE fixed_asset_id = $1 AND accounting_period_id = $2 AND deleted_at IS NULL
	`

	var schedule assets.AssetDepreciationSchedule
	err := r.pool.QueryRow(ctx, query, assetID, periodID).Scan(
		&schedule.ID, &schedule.OrganizationID, &schedule.FixedAssetID,
		&schedule.FiscalYearID, &schedule.AccountingPeriodID,
		&schedule.DepreciationDate, &schedule.DepreciationAmount,
		&schedule.AccumulatedDepreciationBeginning,
		&schedule.AccumulatedDepreciationEnding,
		&schedule.BookValueBeginning, &schedule.BookValueEnding,
		&schedule.JournalEntryID, &schedule.IsPosted,
		&schedule.CreatedAt, &schedule.PostedAt, &schedule.PostedBy,
		&schedule.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get depreciation schedule by asset and period: %w", err)
	}

	return &schedule, nil
}

// GetAssetStatistics retrieves asset statistics
func (r *AssetsRepository) GetAssetStatistics(ctx context.Context, organizationID uuid.UUID) (*assets.AssetStatistics, error) {
	query := `
		SELECT
			COUNT(*) as total_assets,
			COALESCE(SUM(acquisition_cost), 0) as total_acquisition_cost,
			COALESCE(SUM(accumulated_depreciation), 0) as total_accumulated_depreciation,
			COALESCE(SUM(current_book_value), 0) as total_book_value,
			COUNT(CASE WHEN is_disposed = false THEN 1 END) as active_assets,
			COUNT(CASE WHEN is_disposed = true THEN 1 END) as disposed_assets
		FROM fixed_assets
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	var stats assets.AssetStatistics
	var depreciatingAssets int64

	err := r.pool.QueryRow(ctx, query, organizationID).Scan(
		&stats.TotalAssets,
		&stats.TotalAcquisitionCost,
		&stats.TotalAccumulatedDepreciation,
		&stats.TotalBookValue,
		&stats.ActiveAssets,
		&stats.DisposedAssets,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get asset statistics: %w", err)
	}

	// Get depreciating assets count
	depQuery := `
		SELECT COUNT(DISTINCT fixed_asset_id)
		FROM asset_depreciation_schedule
		WHERE organization_id = $1 AND deleted_at IS NULL
	`
	err = r.pool.QueryRow(ctx, depQuery, organizationID).Scan(&depreciatingAssets)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get depreciating assets count: %w", err)
	}

	stats.DepreciatingAssets = depreciatingAssets

	return &stats, nil
}

// GetAssetByCategory retrieves assets by category
func (r *AssetsRepository) GetAssetByCategory(ctx context.Context, organizationID uuid.UUID, categoryID uuid.UUID) ([]*assets.FixedAsset, error) {
	query := `
		SELECT id, organization_id, asset_number, asset_name, asset_category_id,
		       acquisition_date, acquisition_cost, salvage_value, supplier_id,
		       depreciation_method, useful_life_years, depreciation_start_date,
		       asset_account_id, accumulated_depreciation_account_id,
		       depreciation_expense_account_id, current_book_value,
		       accumulated_depreciation, last_depreciation_date, location_id,
		       department, is_disposed, disposal_date, disposal_proceeds,
		       disposal_journal_entry_id, description, serial_number, notes,
		       metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM fixed_assets
		WHERE organization_id = $1 AND asset_category_id = $2 AND deleted_at IS NULL
		ORDER BY asset_number ASC
	`

	rows, err := r.pool.Query(ctx, query, organizationID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets by category: %w", err)
	}
	defer rows.Close()

	var assetList []*assets.FixedAsset
	for rows.Next() {
		var asset assets.FixedAsset
		var metadata sql.NullString
		err := rows.Scan(
			&asset.ID, &asset.OrganizationID, &asset.AssetNumber, &asset.AssetName,
			&asset.AssetCategoryID, &asset.AcquisitionDate, &asset.AcquisitionCost,
			&asset.SalvageValue, &asset.SupplierID, &asset.DepreciationMethod,
			&asset.UsefulLifeYears, &asset.DepreciationStartDate,
			&asset.AssetAccountID, &asset.AccumulatedDepreciationAccountID,
			&asset.DepreciationExpenseAccountID, &asset.CurrentBookValue,
			&asset.AccumulatedDepreciation, &asset.LastDepreciationDate,
			&asset.LocationID, &asset.Department, &asset.IsDisposed,
			&asset.DisposalDate, &asset.DisposalProceeds,
			&asset.DisposalJournalEntryID, &asset.Description, &asset.SerialNumber,
			&asset.Notes, &metadata, &asset.CreatedAt, &asset.UpdatedAt,
			&asset.CreatedBy, &asset.UpdatedBy, &asset.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}

		if metadata.Valid {
			asset.Metadata = json.RawMessage(metadata.String)
		}

		assetList = append(assetList, &asset)
	}

	return assetList, rows.Err()
}
