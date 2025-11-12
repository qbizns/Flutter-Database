package store_credit_transaction

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

// Repository handles database operations for StoreCreditTransactions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new StoreCreditTransactions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// StoreCreditTransactions represents a store_credit_transactions entity
type StoreCreditTransactions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	StoreCreditAccountId uuid.UUID `json:"store_credit_account_id" db:"store_credit_account_id"`
	TransactionType string `json:"transaction_type" db:"transaction_type"`
	TransactionType *string `json:"transaction_type" db:"transaction_type"`
	Amount float64 `json:"amount" db:"amount"`
	BalanceAfter float64 `json:"balance_after" db:"balance_after"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	PaymentId *uuid.UUID `json:"payment_id" db:"payment_id"`
	UserId *uuid.UUID `json:"user_id" db:"user_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	Notes *string `json:"notes" db:"notes"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new store_credit_transactions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *StoreCreditTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "store_credit_transactions", duration, nil)
	}()

	query := `
		INSERT INTO store_credit_transactions (
			, organization_id
			, store_credit_account_id
			, transaction_type
			, transaction_type
			, amount
			, balance_after
			, sale_id
			, payment_id
			, user_id
			, location_id
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
			, $14
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.StoreCreditAccountId,
		entity.TransactionType,
		entity.TransactionType,
		entity.Amount,
		entity.BalanceAfter,
		entity.SaleId,
		entity.PaymentId,
		entity.UserId,
		entity.LocationId,
		entity.Notes,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create store_credit_transactions", zap.Error(err))
		return fmt.Errorf("failed to create store_credit_transactions: %w", err)
	}

	r.logger.Info("created store_credit_transactions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a store_credit_transactions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*StoreCreditTransactions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "store_credit_transactions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, store_credit_account_id
			, transaction_type
			, transaction_type
			, amount
			, balance_after
			, sale_id
			, payment_id
			, user_id
			, location_id
			, notes
			, created_at
			, deleted_at
		FROM store_credit_transactions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity StoreCreditTransactions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.StoreCreditAccountId,
		&entity.TransactionType,
		&entity.TransactionType,
		&entity.Amount,
		&entity.BalanceAfter,
		&entity.SaleId,
		&entity.PaymentId,
		&entity.UserId,
		&entity.LocationId,
		&entity.Notes,
		&entity.CreatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("store_credit_transactions not found")
	}

	if err != nil {
		r.logger.Error("failed to get store_credit_transactions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get store_credit_transactions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of store_credit_transactions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*StoreCreditTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "store_credit_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM store_credit_transactions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count store_credit_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, store_credit_account_id
			, transaction_type
			, transaction_type
			, amount
			, balance_after
			, sale_id
			, payment_id
			, user_id
			, location_id
			, notes
			, created_at
			, deleted_at
		FROM store_credit_transactions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list store_credit_transactions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list store_credit_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*StoreCreditTransactions
	for rows.Next() {
		var entity StoreCreditTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.StoreCreditAccountId,
			&entity.TransactionType,
			&entity.TransactionType,
			&entity.Amount,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.PaymentId,
			&entity.UserId,
			&entity.LocationId,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan store_credit_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating store_credit_transactions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing store_credit_transactions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *StoreCreditTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "store_credit_transactions", duration, nil)
	}()

	query := `
		UPDATE store_credit_transactions
		SET
			, organization_id = $2
			, store_credit_account_id = $3
			, transaction_type = $4
			, transaction_type = $5
			, amount = $6
			, balance_after = $7
			, sale_id = $8
			, payment_id = $9
			, user_id = $10
			, location_id = $11
			, notes = $12
			, deleted_at = $14
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $15
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.StoreCreditAccountId,
		entity.TransactionType,
		entity.TransactionType,
		entity.Amount,
		entity.BalanceAfter,
		entity.SaleId,
		entity.PaymentId,
		entity.UserId,
		entity.LocationId,
		entity.Notes,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update store_credit_transactions", zap.Error(err))
		return fmt.Errorf("failed to update store_credit_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("store_credit_transactions not found or already deleted")
	}

	r.logger.Info("updated store_credit_transactions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a store_credit_transactions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "store_credit_transactions", duration, nil)
	}()

	query := `
		UPDATE store_credit_transactions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete store_credit_transactions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete store_credit_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("store_credit_transactions not found or already deleted")
	}

	r.logger.Info("deleted store_credit_transactions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves store_credit_transactions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*StoreCreditTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "store_credit_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM store_credit_transactions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count store_credit_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, store_credit_account_id
			, transaction_type
			, transaction_type
			, amount
			, balance_after
			, sale_id
			, payment_id
			, user_id
			, location_id
			, notes
			, created_at
			, deleted_at
		FROM store_credit_transactions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list store_credit_transactions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list store_credit_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*StoreCreditTransactions
	for rows.Next() {
		var entity StoreCreditTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.StoreCreditAccountId,
			&entity.TransactionType,
			&entity.TransactionType,
			&entity.Amount,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.PaymentId,
			&entity.UserId,
			&entity.LocationId,
			&entity.Notes,
			&entity.CreatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan store_credit_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

