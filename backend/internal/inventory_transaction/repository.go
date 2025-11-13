package inventory_transaction

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

// Repository handles database operations for InventoryTransactions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new InventoryTransactions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// InventoryTransactions represents a inventory_transactions entity
type InventoryTransactions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ProductId uuid.UUID `json:"product_id" db:"product_id"`
	TransactionType string `json:"transaction_type" db:"transaction_type"`
	Quantity float64 `json:"quantity" db:"quantity"`
	Unit *string `json:"unit" db:"unit"`
	BalanceAfter float64 `json:"balance_after" db:"balance_after"`
	SaleId *uuid.UUID `json:"sale_id" db:"sale_id"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	UnitCost *float64 `json:"unit_cost" db:"unit_cost"`
	TotalCost *float64 `json:"total_cost" db:"total_cost"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	Notes *string `json:"notes" db:"notes"`
	Reason *string `json:"reason" db:"reason"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
}

// Create inserts a new inventory_transactions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *InventoryTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "inventory_transactions", duration, nil)
	}()

	query := `
		INSERT INTO inventory_transactions (
			, organization_id
			, product_id
			, transaction_type
			, quantity
			, unit
			, balance_after
			, sale_id
			, reference_number
			, unit_cost
			, total_cost
			, transaction_date
			, notes
			, reason
			, metadata
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
			, $14
			, $15
			, $17
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.TransactionType,
		entity.Quantity,
		entity.Unit,
		entity.BalanceAfter,
		entity.SaleId,
		entity.ReferenceNumber,
		entity.UnitCost,
		entity.TotalCost,
		entity.TransactionDate,
		entity.Notes,
		entity.Reason,
		entity.Metadata,
		entity.CreatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create inventory_transactions", zap.Error(err))
		return fmt.Errorf("failed to create inventory_transactions: %w", err)
	}

	r.logger.Info("created inventory_transactions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a inventory_transactions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*InventoryTransactions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transactions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, product_id
			, transaction_type
			, quantity
			, unit
			, balance_after
			, sale_id
			, reference_number
			, unit_cost
			, total_cost
			, transaction_date
			, notes
			, reason
			, metadata
			, created_at
			, created_by
		FROM inventory_transactions
		WHERE id = $1
		
	`

	var entity InventoryTransactions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ProductId,
		&entity.TransactionType,
		&entity.Quantity,
		&entity.Unit,
		&entity.BalanceAfter,
		&entity.SaleId,
		&entity.ReferenceNumber,
		&entity.UnitCost,
		&entity.TotalCost,
		&entity.TransactionDate,
		&entity.Notes,
		&entity.Reason,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("inventory_transactions not found")
	}

	if err != nil {
		r.logger.Error("failed to get inventory_transactions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get inventory_transactions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of inventory_transactions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*InventoryTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transactions
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, transaction_type
			, quantity
			, unit
			, balance_after
			, sale_id
			, reference_number
			, unit_cost
			, total_cost
			, transaction_date
			, notes
			, reason
			, metadata
			, created_at
			, created_by
		FROM inventory_transactions
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transactions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransactions
	for rows.Next() {
		var entity InventoryTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.TransactionType,
			&entity.Quantity,
			&entity.Unit,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.ReferenceNumber,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.TransactionDate,
			&entity.Notes,
			&entity.Reason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating inventory_transactions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing inventory_transactions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *InventoryTransactions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "inventory_transactions", duration, nil)
	}()

	query := `
		UPDATE inventory_transactions
		SET
			, organization_id = $2
			, product_id = $3
			, transaction_type = $4
			, quantity = $5
			, unit = $6
			, balance_after = $7
			, sale_id = $8
			, reference_number = $9
			, unit_cost = $10
			, total_cost = $11
			, transaction_date = $12
			, notes = $13
			, reason = $14
			, metadata = $15
			, created_by = $17
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ProductId,
		entity.TransactionType,
		entity.Quantity,
		entity.Unit,
		entity.BalanceAfter,
		entity.SaleId,
		entity.ReferenceNumber,
		entity.UnitCost,
		entity.TotalCost,
		entity.TransactionDate,
		entity.Notes,
		entity.Reason,
		entity.Metadata,
		entity.CreatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update inventory_transactions", zap.Error(err))
		return fmt.Errorf("failed to update inventory_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transactions not found or already deleted")
	}

	r.logger.Info("updated inventory_transactions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a inventory_transactions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "inventory_transactions", duration, nil)
	}()

	query := `DELETE FROM inventory_transactions WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete inventory_transactions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete inventory_transactions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("inventory_transactions not found")
	}

	r.logger.Info("deleted inventory_transactions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves inventory_transactions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*InventoryTransactions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "inventory_transactions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_transactions
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory_transactions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, product_id
			, transaction_type
			, quantity
			, unit
			, balance_after
			, sale_id
			, reference_number
			, unit_cost
			, total_cost
			, transaction_date
			, notes
			, reason
			, metadata
			, created_at
			, created_by
		FROM inventory_transactions
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list inventory_transactions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list inventory_transactions: %w", err)
	}
	defer rows.Close()

	var entities []*InventoryTransactions
	for rows.Next() {
		var entity InventoryTransactions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ProductId,
			&entity.TransactionType,
			&entity.Quantity,
			&entity.Unit,
			&entity.BalanceAfter,
			&entity.SaleId,
			&entity.ReferenceNumber,
			&entity.UnitCost,
			&entity.TotalCost,
			&entity.TransactionDate,
			&entity.Notes,
			&entity.Reason,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.CreatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan inventory_transactions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

