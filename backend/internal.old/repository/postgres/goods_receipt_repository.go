package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/purchases"
)

// GoodsReceiptRepository implements purchases.GoodsReceiptRepository
type GoodsReceiptRepository struct {
	db *DB
}

// NewGoodsReceiptRepository creates a new goods receipt repository
func NewGoodsReceiptRepository(db *DB) *GoodsReceiptRepository {
	return &GoodsReceiptRepository{db: db}
}

// ListReceipts retrieves goods receipts with filters
func (r *GoodsReceiptRepository) ListReceipts(ctx context.Context, orgID uuid.UUID, filters purchases.GoodsReceiptFilters) ([]purchases.GoodsReceipt, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, receipt_number, purchase_order_id, supplier_id,
		       location_id, receipt_date, received_by, status, notes,
		       created_at, updated_at, created_by
		FROM goods_receipts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.PurchaseOrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND purchase_order_id = $%d", argCount)
		args = append(args, *filters.PurchaseOrderID)
	}

	if filters.SupplierID != nil {
		argCount++
		query += fmt.Sprintf(" AND supplier_id = $%d", argCount)
		args = append(args, *filters.SupplierID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND receipt_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND receipt_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND receipt_number ILIKE $%d", argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	query += " ORDER BY receipt_date DESC, receipt_number ASC"

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

	var receipts []purchases.GoodsReceipt
	for rows.Next() {
		var r purchases.GoodsReceipt
		if err := rows.Scan(
			&r.ID, &r.OrganizationID, &r.ReceiptNumber, &r.PurchaseOrderID, &r.SupplierID,
			&r.LocationID, &r.ReceiptDate, &r.ReceivedBy, &r.Status, &r.Notes,
			&r.CreatedAt, &r.UpdatedAt, &r.CreatedBy,
		); err != nil {
			return nil, err
		}
		receipts = append(receipts, r)
	}

	return receipts, rows.Err()
}

// CountReceipts counts goods receipts matching filters
func (r *GoodsReceiptRepository) CountReceipts(ctx context.Context, orgID uuid.UUID, filters purchases.GoodsReceiptFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM goods_receipts WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.PurchaseOrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND purchase_order_id = $%d", argCount)
		args = append(args, *filters.PurchaseOrderID)
	}

	if filters.SupplierID != nil {
		argCount++
		query += fmt.Sprintf(" AND supplier_id = $%d", argCount)
		args = append(args, *filters.SupplierID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND receipt_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND receipt_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND receipt_number ILIKE $%d", argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// CreateReceipt creates a new goods receipt with items
func (r *GoodsReceiptRepository) CreateReceipt(ctx context.Context, receipt *purchases.GoodsReceipt) error {
	if err := r.db.SetOrganizationContext(ctx, receipt.OrganizationID.String()); err != nil {
		return err
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert receipt
	insertQuery := `
		INSERT INTO goods_receipts (
			id, organization_id, receipt_number, purchase_order_id, supplier_id,
			location_id, receipt_date, received_by, status, notes, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(ctx, insertQuery,
		receipt.ID, receipt.OrganizationID, receipt.ReceiptNumber, receipt.PurchaseOrderID,
		receipt.SupplierID, receipt.LocationID, receipt.ReceiptDate, receipt.ReceivedBy,
		receipt.Status, receipt.Notes, receipt.CreatedBy, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert goods receipt: %w", err)
	}

	// Insert receipt items
	itemQuery := `
		INSERT INTO goods_receipt_items (
			id, organization_id, goods_receipt_id, purchase_order_item_id,
			product_id, product_variant_id, quantity_received,
			quantity_accepted, quantity_rejected, rejection_reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	for _, item := range receipt.Items {
		_, err = tx.Exec(ctx, itemQuery,
			item.ID, item.OrganizationID, item.GoodsReceiptID, item.PurchaseOrderItemID,
			item.ProductID, item.ProductVariantID, item.QuantityReceived,
			item.QuantityAccepted, item.QuantityRejected, item.RejectionReason, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to insert goods receipt item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// GetReceipt retrieves a single goods receipt by ID with items
func (r *GoodsReceiptRepository) GetReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*purchases.GoodsReceipt, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, receipt_number, purchase_order_id, supplier_id,
		       location_id, receipt_date, received_by, status, notes,
		       created_at, updated_at, created_by
		FROM goods_receipts
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var receipt purchases.GoodsReceipt
	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&receipt.ID, &receipt.OrganizationID, &receipt.ReceiptNumber, &receipt.PurchaseOrderID,
		&receipt.SupplierID, &receipt.LocationID, &receipt.ReceiptDate, &receipt.ReceivedBy,
		&receipt.Status, &receipt.Notes, &receipt.CreatedAt, &receipt.UpdatedAt, &receipt.CreatedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Fetch items
	itemsQuery := `
		SELECT id, organization_id, goods_receipt_id, purchase_order_item_id,
		       product_id, product_variant_id, quantity_received,
		       quantity_accepted, quantity_rejected, rejection_reason, created_at
		FROM goods_receipt_items
		WHERE goods_receipt_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.Pool.Query(ctx, itemsQuery, id, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item purchases.GoodsReceiptItem
		if err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.GoodsReceiptID, &item.PurchaseOrderItemID,
			&item.ProductID, &item.ProductVariantID, &item.QuantityReceived,
			&item.QuantityAccepted, &item.QuantityRejected, &item.RejectionReason, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		receipt.Items = append(receipt.Items, item)
	}

	return &receipt, rows.Err()
}

// GetReceiptByNumber retrieves a goods receipt by receipt number
func (r *GoodsReceiptRepository) GetReceiptByNumber(ctx context.Context, orgID uuid.UUID, receiptNumber string) (*purchases.GoodsReceipt, error) {
	query := `
		SELECT id FROM goods_receipts
		WHERE organization_id = $1 AND receipt_number = $2 AND deleted_at IS NULL
	`

	var id uuid.UUID
	err := r.db.Pool.QueryRow(ctx, query, orgID, receiptNumber).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.GetReceipt(ctx, orgID, id)
}

// UpdateReceipt updates a goods receipt and its items
func (r *GoodsReceiptRepository) UpdateReceipt(ctx context.Context, receipt *purchases.GoodsReceipt) error {
	if err := r.db.SetOrganizationContext(ctx, receipt.OrganizationID.String()); err != nil {
		return err
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE goods_receipts
		SET status = $1, notes = $2, updated_at = $3
		WHERE id = $4 AND organization_id = $5
	`

	_, err = tx.Exec(ctx, updateQuery, receipt.Status, receipt.Notes, time.Now(), receipt.ID, receipt.OrganizationID)
	if err != nil {
		return fmt.Errorf("failed to update goods receipt: %w", err)
	}

	// Update items
	for _, item := range receipt.Items {
		itemQuery := `
			UPDATE goods_receipt_items
			SET quantity_accepted = $1, quantity_rejected = $2, rejection_reason = $3
			WHERE id = $4 AND organization_id = $5
		`

		_, err = tx.Exec(ctx, itemQuery,
			item.QuantityAccepted, item.QuantityRejected, item.RejectionReason,
			item.ID, item.OrganizationID,
		)
		if err != nil {
			return fmt.Errorf("failed to update goods receipt item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

// DeleteReceipt soft-deletes a goods receipt
func (r *GoodsReceiptRepository) DeleteReceipt(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE goods_receipts
		SET deleted_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id, orgID)
	return err
}

// ListReceiptItems retrieves items for a goods receipt
func (r *GoodsReceiptRepository) ListReceiptItems(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID) ([]purchases.GoodsReceiptItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, goods_receipt_id, purchase_order_item_id,
		       product_id, product_variant_id, quantity_received,
		       quantity_accepted, quantity_rejected, rejection_reason, created_at
		FROM goods_receipt_items
		WHERE goods_receipt_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.Pool.Query(ctx, query, receiptID, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []purchases.GoodsReceiptItem
	for rows.Next() {
		var item purchases.GoodsReceiptItem
		if err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.GoodsReceiptID, &item.PurchaseOrderItemID,
			&item.ProductID, &item.ProductVariantID, &item.QuantityReceived,
			&item.QuantityAccepted, &item.QuantityRejected, &item.RejectionReason, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetReceiptItem retrieves a single goods receipt item
func (r *GoodsReceiptRepository) GetReceiptItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*purchases.GoodsReceiptItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, goods_receipt_id, purchase_order_item_id,
		       product_id, product_variant_id, quantity_received,
		       quantity_accepted, quantity_rejected, rejection_reason, created_at
		FROM goods_receipt_items
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var item purchases.GoodsReceiptItem
	err := r.db.Pool.QueryRow(ctx, query, id, orgID).Scan(
		&item.ID, &item.OrganizationID, &item.GoodsReceiptID, &item.PurchaseOrderItemID,
		&item.ProductID, &item.ProductVariantID, &item.QuantityReceived,
		&item.QuantityAccepted, &item.QuantityRejected, &item.RejectionReason, &item.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// CreateReceiptItem creates a new goods receipt item
func (r *GoodsReceiptRepository) CreateReceiptItem(ctx context.Context, item *purchases.GoodsReceiptItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO goods_receipt_items (
			id, organization_id, goods_receipt_id, purchase_order_item_id,
			product_id, product_variant_id, quantity_received,
			quantity_accepted, quantity_rejected, rejection_reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.GoodsReceiptID, item.PurchaseOrderItemID,
		item.ProductID, item.ProductVariantID, item.QuantityReceived,
		item.QuantityAccepted, item.QuantityRejected, item.RejectionReason, time.Now(),
	)
	return err
}

// UpdateReceiptItem updates a goods receipt item
func (r *GoodsReceiptRepository) UpdateReceiptItem(ctx context.Context, item *purchases.GoodsReceiptItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE goods_receipt_items
		SET quantity_accepted = $1, quantity_rejected = $2, rejection_reason = $3
		WHERE id = $4 AND organization_id = $5
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.QuantityAccepted, item.QuantityRejected, item.RejectionReason,
		item.ID, item.OrganizationID,
	)
	return err
}

// DeleteReceiptItem soft-deletes a goods receipt item
func (r *GoodsReceiptRepository) DeleteReceiptItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE goods_receipt_items
		SET deleted_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err := r.db.Pool.Exec(ctx, query, time.Now(), id, orgID)
	return err
}

// UpdateReceiptStatus updates the status of a goods receipt
func (r *GoodsReceiptRepository) UpdateReceiptStatus(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE goods_receipts
		SET status = $1, updated_at = $2
		WHERE id = $3 AND organization_id = $4
	`

	_, err := r.db.Pool.Exec(ctx, query, status, time.Now(), receiptID, orgID)
	return err
}

// InspectReceiptItems updates multiple receipt items with inspection results
func (r *GoodsReceiptRepository) InspectReceiptItems(ctx context.Context, orgID uuid.UUID, receiptID uuid.UUID, items []purchases.GoodsReceiptItem) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE goods_receipt_items
		SET quantity_accepted = $1, quantity_rejected = $2, rejection_reason = $3
		WHERE id = $4 AND organization_id = $5
	`

	for _, item := range items {
		_, err = tx.Exec(ctx, updateQuery,
			item.QuantityAccepted, item.QuantityRejected, item.RejectionReason,
			item.ID, orgID,
		)
		if err != nil {
			return fmt.Errorf("failed to update item inspection: %w", err)
		}
	}

	// Update receipt status to inspected
	statusQuery := `
		UPDATE goods_receipts
		SET status = 'inspected', updated_at = $1
		WHERE id = $2 AND organization_id = $3
	`

	_, err = tx.Exec(ctx, statusQuery, time.Now(), receiptID, orgID)
	if err != nil {
		return fmt.Errorf("failed to update receipt status: %w", err)
	}

	return tx.Commit(ctx)
}
