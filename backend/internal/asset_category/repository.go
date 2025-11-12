package asset_category

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

// Repository handles database operations for AssetCategories
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AssetCategories repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AssetCategories represents a asset_categories entity
type AssetCategories struct {
	Id *uuid.UUID `json:"id" db:"id"`
	CategoryCode string `json:"category_code" db:"category_code"`
	CategoryName string `json:"category_name" db:"category_name"`
	DefaultDepreciationMethod *string `json:"default_depreciation_method" db:"default_depreciation_method"`
	DefaultUsefulLifeYears *int64 `json:"default_useful_life_years" db:"default_useful_life_years"`
	DefaultSalvageValuePercent *float64 `json:"default_salvage_value_percent" db:"default_salvage_value_percent"`
	AssetAccountId *uuid.UUID `json:"asset_account_id" db:"asset_account_id"`
	AccumulatedDepreciationAccountId *uuid.UUID `json:"accumulated_depreciation_account_id" db:"accumulated_depreciation_account_id"`
	DepreciationExpenseAccountId *uuid.UUID `json:"depreciation_expense_account_id" db:"depreciation_expense_account_id"`
	Description *string `json:"description" db:"description"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new asset_categories record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AssetCategories) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "asset_categories", duration, nil)
	}()

	query := `
		INSERT INTO asset_categories (
			, category_code
			, category_name
			, default_depreciation_method
			, default_useful_life_years
			, default_salvage_value_percent
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, description
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.CategoryCode,
		entity.CategoryName,
		entity.DefaultDepreciationMethod,
		entity.DefaultUsefulLifeYears,
		entity.DefaultSalvageValuePercent,
		entity.AssetAccountId,
		entity.AccumulatedDepreciationAccountId,
		entity.DepreciationExpenseAccountId,
		entity.Description,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create asset_categories", zap.Error(err))
		return fmt.Errorf("failed to create asset_categories: %w", err)
	}

	r.logger.Info("created asset_categories",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a asset_categories by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AssetCategories, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "asset_categories", duration, nil)
	}()

	query := `
		SELECT
			id
			, category_code
			, category_name
			, default_depreciation_method
			, default_useful_life_years
			, default_salvage_value_percent
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, description
			, created_at
			, updated_at
		FROM asset_categories
		WHERE id = $1
		
	`

	var entity AssetCategories
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.CategoryCode,
		&entity.CategoryName,
		&entity.DefaultDepreciationMethod,
		&entity.DefaultUsefulLifeYears,
		&entity.DefaultSalvageValuePercent,
		&entity.AssetAccountId,
		&entity.AccumulatedDepreciationAccountId,
		&entity.DepreciationExpenseAccountId,
		&entity.Description,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("asset_categories not found")
	}

	if err != nil {
		r.logger.Error("failed to get asset_categories", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get asset_categories: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of asset_categories records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AssetCategories, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "asset_categories", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM asset_categories
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count asset_categories records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, category_code
			, category_name
			, default_depreciation_method
			, default_useful_life_years
			, default_salvage_value_percent
			, asset_account_id
			, accumulated_depreciation_account_id
			, depreciation_expense_account_id
			, description
			, created_at
			, updated_at
		FROM asset_categories
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list asset_categories", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list asset_categories: %w", err)
	}
	defer rows.Close()

	var entities []*AssetCategories
	for rows.Next() {
		var entity AssetCategories
		err := rows.Scan(
			&entity.Id,
			&entity.CategoryCode,
			&entity.CategoryName,
			&entity.DefaultDepreciationMethod,
			&entity.DefaultUsefulLifeYears,
			&entity.DefaultSalvageValuePercent,
			&entity.AssetAccountId,
			&entity.AccumulatedDepreciationAccountId,
			&entity.DepreciationExpenseAccountId,
			&entity.Description,
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan asset_categories: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating asset_categories rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing asset_categories record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AssetCategories) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "asset_categories", duration, nil)
	}()

	query := `
		UPDATE asset_categories
		SET
			, category_code = $2
			, category_name = $3
			, default_depreciation_method = $4
			, default_useful_life_years = $5
			, default_salvage_value_percent = $6
			, asset_account_id = $7
			, accumulated_depreciation_account_id = $8
			, depreciation_expense_account_id = $9
			, description = $10
			, updated_at = $12
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.CategoryCode,
		entity.CategoryName,
		entity.DefaultDepreciationMethod,
		entity.DefaultUsefulLifeYears,
		entity.DefaultSalvageValuePercent,
		entity.AssetAccountId,
		entity.AccumulatedDepreciationAccountId,
		entity.DepreciationExpenseAccountId,
		entity.Description,
		entity.UpdatedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update asset_categories", zap.Error(err))
		return fmt.Errorf("failed to update asset_categories: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("asset_categories not found or already deleted")
	}

	r.logger.Info("updated asset_categories",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a asset_categories record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "asset_categories", duration, nil)
	}()

	query := `DELETE FROM asset_categories WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete asset_categories", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete asset_categories: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("asset_categories not found")
	}

	r.logger.Info("deleted asset_categories", zap.String("id", id.String()))
	return nil
}



