package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/inventory"
)

type InventoryRepository struct {
	db *DB
}

func NewInventoryRepository(db *DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) List(ctx context.Context, orgID uuid.UUID, filters inventory.InventoryTransactionFilters) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, filters.TransactionType)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY transaction_date DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) Count(ctx context.Context, orgID uuid.UUID, filters inventory.InventoryTransactionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM inventory_transactions WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, filters.TransactionType)
	}

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *InventoryRepository) Create(ctx context.Context, transaction *inventory.InventoryTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, transaction.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO inventory_transactions (
			id, organization_id, product_id, transaction_type, quantity, unit,
			balance_after, sale_id, reference_number, unit_cost, total_cost,
			transaction_date, notes, reason, metadata, created_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transaction.ID, transaction.OrganizationID, transaction.ProductID, transaction.TransactionType, transaction.Quantity, transaction.Unit,
		transaction.BalanceAfter, transaction.SaleID, transaction.ReferenceNumber, transaction.UnitCost, transaction.TotalCost,
		transaction.TransactionDate, transaction.Notes, transaction.Reason, transaction.Metadata, transaction.CreatedAt, transaction.CreatedBy,
	)
	return err
}

func (r *InventoryRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND id = $2
	`

	var tx inventory.InventoryTransaction
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
		&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
		&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.Metadata = metadata
	return &tx, nil
}

func (r *InventoryRepository) GetByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, limit int, offset int) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND product_id = $2
		ORDER BY transaction_date DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND sale_id = $2
		ORDER BY transaction_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txList []inventory.InventoryTransaction
	for rows.Next() {
		var tx inventory.InventoryTransaction
		var metadata interface{}
		err := rows.Scan(
			&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
			&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
			&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		tx.Metadata = metadata
		txList = append(txList, tx)
	}

	return txList, rows.Err()
}

func (r *InventoryRepository) Update(ctx context.Context, transaction *inventory.InventoryTransaction) error {
	if err := r.db.SetOrganizationContext(ctx, transaction.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE inventory_transactions SET
			transaction_type = $3, quantity = $4, unit = $5, balance_after = $6,
			sale_id = $7, reference_number = $8, unit_cost = $9, total_cost = $10,
			transaction_date = $11, notes = $12, reason = $13, metadata = $14
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		transaction.OrganizationID, transaction.ID,
		transaction.TransactionType, transaction.Quantity, transaction.Unit, transaction.BalanceAfter,
		transaction.SaleID, transaction.ReferenceNumber, transaction.UnitCost, transaction.TotalCost,
		transaction.TransactionDate, transaction.Notes, transaction.Reason, transaction.Metadata,
	)
	return err
}

func (r *InventoryRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		DELETE FROM inventory_transactions
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

func (r *InventoryRepository) GetLatestBalanceForProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) (*inventory.InventoryTransaction, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, product_id, transaction_type, quantity, unit,
		       balance_after, sale_id, reference_number, unit_cost, total_cost,
		       transaction_date, notes, reason, metadata, created_at, created_by
		FROM inventory_transactions
		WHERE organization_id = $1 AND product_id = $2
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT 1
	`

	var tx inventory.InventoryTransaction
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, productID).Scan(
		&tx.ID, &tx.OrganizationID, &tx.ProductID, &tx.TransactionType, &tx.Quantity, &tx.Unit,
		&tx.BalanceAfter, &tx.SaleID, &tx.ReferenceNumber, &tx.UnitCost, &tx.TotalCost,
		&tx.TransactionDate, &tx.Notes, &tx.Reason, &metadata, &tx.CreatedAt, &tx.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tx.Metadata = metadata
	return &tx, nil
}
