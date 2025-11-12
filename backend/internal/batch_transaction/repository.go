package batch_transaction

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

// Repository handles database operations for BatchTransactions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new BatchTransactions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// BatchTransactions represents a batch_transactions entity
type BatchTransactions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BatchId uuid.UUID `json:"batch_id" db:"batch_id"`
	TransactionType string `json:"transaction_type" db:"transaction_type"`
	Quantity float64 `json:"quantity" db:"quantity"`
	BalanceAfter float64 `json:"balance_after" db:"balance_after"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	InventoryTransferId *uuid.UUID `json:"inventory_transfer_id" db:"inventory_transfer_id"`
	Reason *string `json:"reason" db:"reason"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	TransactionDate *time.Time `json:"transaction_date" db:"transaction_date"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new batch_transactions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *BatchTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "batch_transactions", duration, nil)
	}()

	query := `
		INSERT INTO batch_transactions (
			, organization_id
			, batch_id
			, transaction_type
			, quantity
			, balance_after
			, sale_id
			, inventory_transfer_id
			, reason
			, notes
			, metadata
			, transaction_date
			, created_by
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
		)
		RETURNING id
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BatchId,
		entity.TransactionType,
		entity.Quantity,
		entity.BalanceAfter,
		entity.SaleId,
		entity.InventoryTransferId,
		entity.Reason,
		entity.Notes,
		entity.Metadata,
		entity.TransactionDate,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id)
	

	if err != nil {
		r.logger.Error("failed to create batch_transactions", zap.Error(err))
		return fmt.Errorf("failed to create batch_transactions: %w", err)
	}

	r.logger.Info("created batch_transactions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a batch_transactions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*BatchTransactions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "batch_transactions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, batch_id
			, transaction_type
			, quantity
			, balance_after
			, sale_id
			, inventory_transfer_id
			, reason
			, notes
			, metadata
			, transaction_date
			, created_by
		FROM batch_transactions
		WHERE id = $1
		
	`

	var entity BatchTransactions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BatchId,
		&entity.TransactionType,
		&entity.Quantity,
		&entity.BalanceAfter,
		&entity.SaleId,
		&entity.InventoryTransferId,
		&entity.Reason,
		&entity.Notes,
		&entity.Metadata,
		&entity.TransactionDate,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("batch_transactions not found")
	}

	if err != nil {
		r.logger.Error("failed to get batch_transactions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get batch_transactions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of batch_transactions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*BatchTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "batch_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM batch_transactions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count batch_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, batch_id
			, transaction_type
			, quantity
			, balance_after
			, sale_id
			, inventory_transfer_id
			, reason
			, notes
			, metadata
			, transaction_date
			, created_by
		FROM batch_transactions
		
		
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list batch_transactions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list batch_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*BatchTransactions
	for rows.Next() {
		var entity BatchTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BatchId,
			&entity.TransactionType,
			&entity.Quantity,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.InventoryTransferId,
			&entity.Reason,
			&entity.Notes,
			&entity.Metadata,
			&entity.TransactionDate,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan batch_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating batch_transactions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing batch_transactions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *BatchTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "batch_transactions", duration, nil)
	}()

	query := `
		UPDATE batch_transactions
		SET
			, organization_id = $2
			, batch_id = $3
			, transaction_type = $4
			, quantity = $5
			, balance_after = $6
			, sale_id = $7
			, inventory_transfer_id = $8
			, reason = $9
			, notes = $10
			, metadata = $11
			, transaction_date = $12
			, created_by = $13
			
		WHERE id = $14
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BatchId,
		entity.TransactionType,
		entity.Quantity,
		entity.BalanceAfter,
		entity.SaleId,
		entity.InventoryTransferId,
		entity.Reason,
		entity.Notes,
		entity.Metadata,
		entity.TransactionDate,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update batch_transactions", zap.Error(err))
		return fmt.Errorf("failed to update batch_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("batch_transactions not found or already deleted")
	}

	r.logger.Info("updated batch_transactions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a batch_transactions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "batch_transactions", duration, nil)
	}()

	query := `DELETE FROM batch_transactions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete batch_transactions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete batch_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("batch_transactions not found")
	}

	r.logger.Info("deleted batch_transactions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves batch_transactions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*BatchTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "batch_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM batch_transactions
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count batch_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, batch_id
			, transaction_type
			, quantity
			, balance_after
			, sale_id
			, inventory_transfer_id
			, reason
			, notes
			, metadata
			, transaction_date
			, created_by
		FROM batch_transactions
		WHERE organization_id = $1
		
		
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list batch_transactions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list batch_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*BatchTransactions
	for rows.Next() {
		var entity BatchTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BatchId,
			&entity.TransactionType,
			&entity.Quantity,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.InventoryTransferId,
			&entity.Reason,
			&entity.Notes,
			&entity.Metadata,
			&entity.TransactionDate,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan batch_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

