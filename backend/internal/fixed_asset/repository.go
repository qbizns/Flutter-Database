package fixed_asset

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

// Repository handles database operations for FixedAssets
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new FixedAssets repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FixedAssets represents a fixed_assets entity
type FixedAssets struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	AssetNumber string `json:"asset_number" db:"asset_number"`
	AssetName string `json:"asset_name" db:"asset_name"`
	AssetCategoryId *uuid.UUID `json:"asset_category_id" db:"asset_category_id"`
	AcquisitionDate time.Time `json:"acquisition_date" db:"acquisition_date"`
	AcquisitionCost float64 `json:"acquisition_cost" db:"acquisition_cost"`
	SalvageValue *float64 `json:"salvage_value" db:"salvage_value"`
	SupplierId *uuid.UUID `json:"supplier_id" db:"supplier_id"`
	VendorBillId *uuid.UUID `json:"vendor_bill_id" db:"vendor_bill_id"`
	DepreciationMethod string `json:"depreciation_method" db:"depreciation_method"`
	UsefulLifeYears int64 `json:"useful_life_years" db:"useful_life_years"`
	DepreciationStartDate time.Time `json:"depreciation_start_date" db:"depreciation_start_date"`
	AssetAccountId uuid.UUID `json:"asset_account_id" db:"asset_account_id"`
	AccumulatedDepreciationAccountId uuid.UUID `json:"accumulated_depreciation_account_id" db:"accumulated_depreciation_account_id"`
	DepreciationExpenseAccountId uuid.UUID `json:"depreciation_expense_account_id" db:"depreciation_expense_account_id"`
	CurrentBookValue *float64 `json:"current_book_value" db:"current_book_value"`
	AccumulatedDepreciation *float64 `json:"accumulated_depreciation" db:"accumulated_depreciation"`
	LastDepreciationDate *time.Time `json:"last_depreciation_date" db:"last_depreciation_date"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	Department *string `json:"department" db:"department"`
	IsDisposed *bool `json:"is_disposed" db:"is_disposed"`
	DisposalDate *time.Time `json:"disposal_date" db:"disposal_date"`
	DisposalProceeds *float64 `json:"disposal_proceeds" db:"disposal_proceeds"`
	DisposalJournalEntryId *uuid.UUID `json:"disposal_journal_entry_id" db:"disposal_journal_entry_id"`
	Description *string `json:"description" db:"description"`
	SerialNumber *string `json:"serial_number" db:"serial_number"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new fixed_assets record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *FixedAssets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "fixed_assets", duration, nil)
	}()

	query := `
		INSERT INTO fixed_assets (
			, organization_id
			, asset_number
			, asset_name
			, asset_category_id
			, acquisition_date
			, acquisition_cost
			, salvage_value
			, supplier_id
			, vendor_bill_id
			, depreciation_method
			, useful_life_years
			, depreciation_start_date
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, current_book_value
			, accumulated_depreciation
			, last_depreciation_date
			, location_id
			, department
			, is_disposed
			, disposal_date
			, disposal_proceeds
			, disposal_journal_entry_id
			, description
			, serial_number
			, notes
			, metadata
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
			, $20
			, $21
			, $22
			, $23
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $32
			, $33
			, $34
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.AssetNumber,
		entity.AssetName,
		entity.AssetCategoryId,
		entity.AcquisitionDate,
		entity.AcquisitionCost,
		entity.SalvageValue,
		entity.SupplierId,
		entity.VendorBillId,
		entity.DepreciationMethod,
		entity.UsefulLifeYears,
		entity.DepreciationStartDate,
		entity.AssetAccountId,
		entity.AccumulatedDepreciationAccountId,
		entity.DepreciationExpenseAccountId,
		entity.CurrentBookValue,
		entity.AccumulatedDepreciation,
		entity.LastDepreciationDate,
		entity.LocationId,
		entity.Department,
		entity.IsDisposed,
		entity.DisposalDate,
		entity.DisposalProceeds,
		entity.DisposalJournalEntryId,
		entity.Description,
		entity.SerialNumber,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create fixed_assets", zap.Error(err))
		return fmt.Errorf("failed to create fixed_assets: %w", err)
	}

	r.logger.Info("created fixed_assets",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a fixed_assets by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*FixedAssets, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fixed_assets", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, asset_number
			, asset_name
			, asset_category_id
			, acquisition_date
			, acquisition_cost
			, salvage_value
			, supplier_id
			, vendor_bill_id
			, depreciation_method
			, useful_life_years
			, depreciation_start_date
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, current_book_value
			, accumulated_depreciation
			, last_depreciation_date
			, location_id
			, department
			, is_disposed
			, disposal_date
			, disposal_proceeds
			, disposal_journal_entry_id
			, description
			, serial_number
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fixed_assets
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity FixedAssets
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.AssetNumber,
		&entity.AssetName,
		&entity.AssetCategoryId,
		&entity.AcquisitionDate,
		&entity.AcquisitionCost,
		&entity.SalvageValue,
		&entity.SupplierId,
		&entity.VendorBillId,
		&entity.DepreciationMethod,
		&entity.UsefulLifeYears,
		&entity.DepreciationStartDate,
		&entity.AssetAccountId,
		&entity.AccumulatedDepreciationAccountId,
		&entity.DepreciationExpenseAccountId,
		&entity.CurrentBookValue,
		&entity.AccumulatedDepreciation,
		&entity.LastDepreciationDate,
		&entity.LocationId,
		&entity.Department,
		&entity.IsDisposed,
		&entity.DisposalDate,
		&entity.DisposalProceeds,
		&entity.DisposalJournalEntryId,
		&entity.Description,
		&entity.SerialNumber,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("fixed_assets not found")
	}

	if err != nil {
		r.logger.Error("failed to get fixed_assets", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get fixed_assets: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of fixed_assets records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*FixedAssets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fixed_assets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fixed_assets
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fixed_assets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, asset_number
			, asset_name
			, asset_category_id
			, acquisition_date
			, acquisition_cost
			, salvage_value
			, supplier_id
			, vendor_bill_id
			, depreciation_method
			, useful_life_years
			, depreciation_start_date
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, current_book_value
			, accumulated_depreciation
			, last_depreciation_date
			, location_id
			, department
			, is_disposed
			, disposal_date
			, disposal_proceeds
			, disposal_journal_entry_id
			, description
			, serial_number
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fixed_assets
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fixed_assets", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fixed_assets: %w", err)
	}
	defer rows.Close()

	var entities []*FixedAssets
	for rows.Next() {
		var entity FixedAssets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AssetNumber,
			&entity.AssetName,
			&entity.AssetCategoryId,
			&entity.AcquisitionDate,
			&entity.AcquisitionCost,
			&entity.SalvageValue,
			&entity.SupplierId,
			&entity.VendorBillId,
			&entity.DepreciationMethod,
			&entity.UsefulLifeYears,
			&entity.DepreciationStartDate,
			&entity.AssetAccountId,
			&entity.AccumulatedDepreciationAccountId,
			&entity.DepreciationExpenseAccountId,
			&entity.CurrentBookValue,
			&entity.AccumulatedDepreciation,
			&entity.LastDepreciationDate,
			&entity.LocationId,
			&entity.Department,
			&entity.IsDisposed,
			&entity.DisposalDate,
			&entity.DisposalProceeds,
			&entity.DisposalJournalEntryId,
			&entity.Description,
			&entity.SerialNumber,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fixed_assets: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating fixed_assets rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing fixed_assets record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *FixedAssets) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fixed_assets", duration, nil)
	}()

	query := `
		UPDATE fixed_assets
		SET
			, organization_id = $2
			, asset_number = $3
			, asset_name = $4
			, asset_category_id = $5
			, acquisition_date = $6
			, acquisition_cost = $7
			, salvage_value = $8
			, supplier_id = $9
			, vendor_bill_id = $10
			, depreciation_method = $11
			, useful_life_years = $12
			, depreciation_start_date = $13
			, asset_account_id = $14
			, accumulated_depreciation_account_id = $15
			, depreciation_expense_account_id = $16
			, current_book_value = $17
			, accumulated_depreciation = $18
			, last_depreciation_date = $19
			, location_id = $20
			, department = $21
			, is_disposed = $22
			, disposal_date = $23
			, disposal_proceeds = $24
			, disposal_journal_entry_id = $25
			, description = $26
			, serial_number = $27
			, notes = $28
			, metadata = $29
			, updated_at = $31
			, created_by = $32
			, updated_by = $33
			, deleted_at = $34
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $35
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.AssetNumber,
		entity.AssetName,
		entity.AssetCategoryId,
		entity.AcquisitionDate,
		entity.AcquisitionCost,
		entity.SalvageValue,
		entity.SupplierId,
		entity.VendorBillId,
		entity.DepreciationMethod,
		entity.UsefulLifeYears,
		entity.DepreciationStartDate,
		entity.AssetAccountId,
		entity.AccumulatedDepreciationAccountId,
		entity.DepreciationExpenseAccountId,
		entity.CurrentBookValue,
		entity.AccumulatedDepreciation,
		entity.LastDepreciationDate,
		entity.LocationId,
		entity.Department,
		entity.IsDisposed,
		entity.DisposalDate,
		entity.DisposalProceeds,
		entity.DisposalJournalEntryId,
		entity.Description,
		entity.SerialNumber,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update fixed_assets", zap.Error(err))
		return fmt.Errorf("failed to update fixed_assets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fixed_assets not found or already deleted")
	}

	r.logger.Info("updated fixed_assets",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a fixed_assets record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "fixed_assets", duration, nil)
	}()

	query := `
		UPDATE fixed_assets
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete fixed_assets", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete fixed_assets: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("fixed_assets not found or already deleted")
	}

	r.logger.Info("deleted fixed_assets", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves fixed_assets records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*FixedAssets, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "fixed_assets", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM fixed_assets
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fixed_assets records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, asset_number
			, asset_name
			, asset_category_id
			, acquisition_date
			, acquisition_cost
			, salvage_value
			, supplier_id
			, vendor_bill_id
			, depreciation_method
			, useful_life_years
			, depreciation_start_date
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, current_book_value
			, accumulated_depreciation
			, last_depreciation_date
			, location_id
			, department
			, is_disposed
			, disposal_date
			, disposal_proceeds
			, disposal_journal_entry_id
			, description
			, serial_number
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM fixed_assets
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list fixed_assets by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list fixed_assets: %w", err)
	}
	defer rows.Close()

	var entities []*FixedAssets
	for rows.Next() {
		var entity FixedAssets
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.AssetNumber,
			&entity.AssetName,
			&entity.AssetCategoryId,
			&entity.AcquisitionDate,
			&entity.AcquisitionCost,
			&entity.SalvageValue,
			&entity.SupplierId,
			&entity.VendorBillId,
			&entity.DepreciationMethod,
			&entity.UsefulLifeYears,
			&entity.DepreciationStartDate,
			&entity.AssetAccountId,
			&entity.AccumulatedDepreciationAccountId,
			&entity.DepreciationExpenseAccountId,
			&entity.CurrentBookValue,
			&entity.AccumulatedDepreciation,
			&entity.LastDepreciationDate,
			&entity.LocationId,
			&entity.Department,
			&entity.IsDisposed,
			&entity.DisposalDate,
			&entity.DisposalProceeds,
			&entity.DisposalJournalEntryId,
			&entity.Description,
			&entity.SerialNumber,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan fixed_assets: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

