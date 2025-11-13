package account_type

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

// Repository handles database operations for AccountTypes
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new AccountTypes repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// AccountTypes represents a account_types entity
type AccountTypes struct {
	Id *uuid.UUID `json:"id" db:"id"`
	TypeCode string `json:"type_code" db:"type_code"`
	TypeName string `json:"type_name" db:"type_name"`
	TypeCategory string `json:"type_category" db:"type_category"`
	NormalBalance string `json:"normal_balance" db:"normal_balance"`
	IsBalanceSheet *bool `json:"is_balance_sheet" db:"is_balance_sheet"`
	IsIncomeStatement *bool `json:"is_income_statement" db:"is_income_statement"`
	DisplayOrder *int64 `json:"display_order" db:"display_order"`
	Description *string `json:"description" db:"description"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// Create inserts a new account_types record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *AccountTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "account_types", duration, nil)
	}()

	query := `
		INSERT INTO account_types (
			, type_code
			, type_name
			, type_category
			, normal_balance
			, is_balance_sheet
			, is_income_statement
			, display_order
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.TypeCode,
		entity.TypeName,
		entity.TypeCategory,
		entity.NormalBalance,
		entity.IsBalanceSheet,
		entity.IsIncomeStatement,
		entity.DisplayOrder,
		entity.Description,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create account_types", zap.Error(err))
		return fmt.Errorf("failed to create account_types: %w", err)
	}

	r.logger.Info("created account_types",
		zap.String("id", entity.Id.String()),
		
	)

	return nil
}

// GetByID retrieves a account_types by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*AccountTypes, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "account_types", duration, nil)
	}()

	query := `
		SELECT
			id
			, type_code
			, type_name
			, type_category
			, normal_balance
			, is_balance_sheet
			, is_income_statement
			, display_order
			, description
			, created_at
			, updated_at
		FROM account_types
		WHERE id = $1
		
	`

	var entity AccountTypes
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.TypeCode,
		&entity.TypeName,
		&entity.TypeCategory,
		&entity.NormalBalance,
		&entity.IsBalanceSheet,
		&entity.IsIncomeStatement,
		&entity.DisplayOrder,
		&entity.Description,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account_types not found")
	}

	if err != nil {
		r.logger.Error("failed to get account_types", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get account_types: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of account_types records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*AccountTypes, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "account_types", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM account_types
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count account_types records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, type_code
			, type_name
			, type_category
			, normal_balance
			, is_balance_sheet
			, is_income_statement
			, display_order
			, description
			, created_at
			, updated_at
		FROM account_types
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list account_types", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list account_types: %w", err)
	}
	defer rows.Close()

	var entities []*AccountTypes
	for rows.Next() {
		var entity AccountTypes
		err := rows.Scan(
			&entity.Id,
			&entity.TypeCode,
			&entity.TypeName,
			&entity.TypeCategory,
			&entity.NormalBalance,
			&entity.IsBalanceSheet,
			&entity.IsIncomeStatement,
			&entity.DisplayOrder,
			&entity.Description,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan account_types: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating account_types rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing account_types record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *AccountTypes) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "account_types", duration, nil)
	}()

	query := `
		UPDATE account_types
		SET
			, type_code = $2
			, type_name = $3
			, type_category = $4
			, normal_balance = $5
			, is_balance_sheet = $6
			, is_income_statement = $7
			, display_order = $8
			, description = $9
			, updated_at = $11
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $12
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.TypeCode,
		entity.TypeName,
		entity.TypeCategory,
		entity.NormalBalance,
		entity.IsBalanceSheet,
		entity.IsIncomeStatement,
		entity.DisplayOrder,
		entity.Description,
		time.Now(),
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update account_types", zap.Error(err))
		return fmt.Errorf("failed to update account_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account_types not found or already deleted")
	}

	r.logger.Info("updated account_types",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a account_types record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "account_types", duration, nil)
	}()

	query := `DELETE FROM account_types WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete account_types", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete account_types: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account_types not found")
	}

	r.logger.Info("deleted account_types", zap.String("id", id.String()))
	return nil
}



