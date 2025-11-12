package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/purchases"
)

// PurchaseRepository handles purchase order database operations
type PurchaseRepository struct {
	db *DB
}

// NewPurchaseRepository creates a new purchase repository
func NewPurchaseRepository(db *DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

// ListOrders retrieves all purchase orders matching filters
func (r *PurchaseRepository) ListOrders(ctx context.Context, orgID uuid.UUID, filters purchases.PurchaseOrderFilters) ([]purchases.PurchaseOrder, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_id, po_number, po_date, expected_delivery_date,
		       actual_delivery_date, status, subtotal, tax_amount, discount_amount, shipping_cost,
		       total_amount, payment_status, paid_amount, notes, terms_and_conditions,
		       created_at, updated_at, created_by, updated_by, approved_by, approved_at
		FROM purchase_orders
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.SupplierID != nil {
		argCount++
		query += fmt.Sprintf(" AND supplier_id = $%d", argCount)
		args = append(args, *filters.SupplierID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (po_number ILIKE $%d OR notes ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND po_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND po_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY po_date DESC, po_number DESC"

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

	var orders []purchases.PurchaseOrder
	for rows.Next() {
		var po purchases.PurchaseOrder
		err := rows.Scan(
			&po.ID, &po.OrganizationID, &po.SupplierID, &po.PONumber, &po.PODate,
			&po.ExpectedDeliveryDate, &po.ActualDeliveryDate, &po.Status, &po.Subtotal,
			&po.TaxAmount, &po.DiscountAmount, &po.ShippingCost, &po.TotalAmount,
			&po.PaymentStatus, &po.PaidAmount, &po.Notes, &po.TermsAndConditions,
			&po.CreatedAt, &po.UpdatedAt, &po.CreatedBy, &po.UpdatedBy,
			&po.ApprovedBy, &po.ApprovedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, po)
	}

	return orders, rows.Err()
}

// CountOrders returns total count of purchase orders matching filters
func (r *PurchaseRepository) CountOrders(ctx context.Context, orgID uuid.UUID, filters purchases.PurchaseOrderFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM purchase_orders WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.SupplierID != nil {
		argCount++
		query += fmt.Sprintf(" AND supplier_id = $%d", argCount)
		args = append(args, *filters.SupplierID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (po_number ILIKE $%d OR notes ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND po_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND po_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// CreateOrder creates a new purchase order
func (r *PurchaseRepository) CreateOrder(ctx context.Context, po *purchases.PurchaseOrder) error {
	if err := r.db.SetOrganizationContext(ctx, po.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO purchase_orders (
			id, organization_id, supplier_id, po_number, po_date, expected_delivery_date,
			actual_delivery_date, status, subtotal, tax_amount, discount_amount, shipping_cost,
			total_amount, payment_status, paid_amount, notes, terms_and_conditions,
			created_at, updated_at, created_by, approved_by, approved_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		po.ID, po.OrganizationID, po.SupplierID, po.PONumber, po.PODate,
		po.ExpectedDeliveryDate, po.ActualDeliveryDate, po.Status, po.Subtotal,
		po.TaxAmount, po.DiscountAmount, po.ShippingCost, po.TotalAmount,
		po.PaymentStatus, po.PaidAmount, po.Notes, po.TermsAndConditions,
		po.CreatedAt, po.UpdatedAt, po.CreatedBy, po.ApprovedBy, po.ApprovedAt,
	)
	return err
}

// GetOrder retrieves a purchase order by ID
func (r *PurchaseRepository) GetOrder(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*purchases.PurchaseOrder, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_id, po_number, po_date, expected_delivery_date,
		       actual_delivery_date, status, subtotal, tax_amount, discount_amount, shipping_cost,
		       total_amount, payment_status, paid_amount, notes, terms_and_conditions,
		       created_at, updated_at, created_by, updated_by, approved_by, approved_at
		FROM purchase_orders
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var po purchases.PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&po.ID, &po.OrganizationID, &po.SupplierID, &po.PONumber, &po.PODate,
		&po.ExpectedDeliveryDate, &po.ActualDeliveryDate, &po.Status, &po.Subtotal,
		&po.TaxAmount, &po.DiscountAmount, &po.ShippingCost, &po.TotalAmount,
		&po.PaymentStatus, &po.PaidAmount, &po.Notes, &po.TermsAndConditions,
		&po.CreatedAt, &po.UpdatedAt, &po.CreatedBy, &po.UpdatedBy,
		&po.ApprovedBy, &po.ApprovedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &po, nil
}

// GetOrderByNumber retrieves a purchase order by PO number
func (r *PurchaseRepository) GetOrderByNumber(ctx context.Context, orgID uuid.UUID, poNumber string) (*purchases.PurchaseOrder, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_id, po_number, po_date, expected_delivery_date,
		       actual_delivery_date, status, subtotal, tax_amount, discount_amount, shipping_cost,
		       total_amount, payment_status, paid_amount, notes, terms_and_conditions,
		       created_at, updated_at, created_by, updated_by, approved_by, approved_at
		FROM purchase_orders
		WHERE organization_id = $1 AND po_number = $2 AND deleted_at IS NULL
	`

	var po purchases.PurchaseOrder
	err := r.db.Pool.QueryRow(ctx, query, orgID, poNumber).Scan(
		&po.ID, &po.OrganizationID, &po.SupplierID, &po.PONumber, &po.PODate,
		&po.ExpectedDeliveryDate, &po.ActualDeliveryDate, &po.Status, &po.Subtotal,
		&po.TaxAmount, &po.DiscountAmount, &po.ShippingCost, &po.TotalAmount,
		&po.PaymentStatus, &po.PaidAmount, &po.Notes, &po.TermsAndConditions,
		&po.CreatedAt, &po.UpdatedAt, &po.CreatedBy, &po.UpdatedBy,
		&po.ApprovedBy, &po.ApprovedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &po, nil
}

// UpdateOrder updates a purchase order
func (r *PurchaseRepository) UpdateOrder(ctx context.Context, po *purchases.PurchaseOrder) error {
	if err := r.db.SetOrganizationContext(ctx, po.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_orders SET
			po_number = $3, expected_delivery_date = $4, actual_delivery_date = $5,
			status = $6, subtotal = $7, tax_amount = $8, discount_amount = $9,
			shipping_cost = $10, total_amount = $11, payment_status = $12, paid_amount = $13,
			notes = $14, terms_and_conditions = $15, updated_at = $16, updated_by = $17,
			approved_by = $18, approved_at = $19
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		po.OrganizationID, po.ID,
		po.PONumber, po.ExpectedDeliveryDate, po.ActualDeliveryDate,
		po.Status, po.Subtotal, po.TaxAmount, po.DiscountAmount,
		po.ShippingCost, po.TotalAmount, po.PaymentStatus, po.PaidAmount,
		po.Notes, po.TermsAndConditions, po.UpdatedAt, po.UpdatedBy,
		po.ApprovedBy, po.ApprovedAt,
	)
	return err
}

// DeleteOrder soft-deletes a purchase order
func (r *PurchaseRepository) DeleteOrder(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_orders
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ListItems retrieves all items for a purchase order
func (r *PurchaseRepository) ListItems(ctx context.Context, orgID uuid.UUID, poID uuid.UUID) ([]purchases.PurchaseOrderItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, purchase_order_id, product_id, product_name, product_sku,
		       quantity_ordered, quantity_received, unit_of_measure, unit_cost, discount_percentage,
		       discount_amount, tax_percentage, tax_amount, line_total, notes, created_at, updated_at
		FROM purchase_order_items
		WHERE organization_id = $1 AND purchase_order_id = $2
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, poID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []purchases.PurchaseOrderItem
	for rows.Next() {
		var item purchases.PurchaseOrderItem
		err := rows.Scan(
			&item.ID, &item.OrganizationID, &item.PurchaseOrderID, &item.ProductID,
			&item.ProductName, &item.ProductSKU, &item.QuantityOrdered, &item.QuantityReceived,
			&item.UnitOfMeasure, &item.UnitCost, &item.DiscountPercent, &item.DiscountAmount,
			&item.TaxPercent, &item.TaxAmount, &item.LineTotal, &item.Notes,
			&item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetItem retrieves a single purchase order item
func (r *PurchaseRepository) GetItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*purchases.PurchaseOrderItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, purchase_order_id, product_id, product_name, product_sku,
		       quantity_ordered, quantity_received, unit_of_measure, unit_cost, discount_percentage,
		       discount_amount, tax_percentage, tax_amount, line_total, notes, created_at, updated_at
		FROM purchase_order_items
		WHERE organization_id = $1 AND id = $2
	`

	var item purchases.PurchaseOrderItem
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&item.ID, &item.OrganizationID, &item.PurchaseOrderID, &item.ProductID,
		&item.ProductName, &item.ProductSKU, &item.QuantityOrdered, &item.QuantityReceived,
		&item.UnitOfMeasure, &item.UnitCost, &item.DiscountPercent, &item.DiscountAmount,
		&item.TaxPercent, &item.TaxAmount, &item.LineTotal, &item.Notes,
		&item.CreatedAt, &item.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateItem creates a new purchase order item
func (r *PurchaseRepository) CreateItem(ctx context.Context, item *purchases.PurchaseOrderItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO purchase_order_items (
			id, organization_id, purchase_order_id, product_id, product_name, product_sku,
			quantity_ordered, quantity_received, unit_of_measure, unit_cost, discount_percentage,
			discount_amount, tax_percentage, tax_amount, line_total, notes, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.PurchaseOrderID, item.ProductID,
		item.ProductName, item.ProductSKU, item.QuantityOrdered, item.QuantityReceived,
		item.UnitOfMeasure, item.UnitCost, item.DiscountPercent, item.DiscountAmount,
		item.TaxPercent, item.TaxAmount, item.LineTotal, item.Notes,
		item.CreatedAt, item.UpdatedAt,
	)
	return err
}

// UpdateItem updates a purchase order item
func (r *PurchaseRepository) UpdateItem(ctx context.Context, item *purchases.PurchaseOrderItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_order_items SET
			product_name = $3, product_sku = $4, quantity_ordered = $5, quantity_received = $6,
			unit_of_measure = $7, unit_cost = $8, discount_percentage = $9, discount_amount = $10,
			tax_percentage = $11, tax_amount = $12, line_total = $13, notes = $14, updated_at = $15
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.OrganizationID, item.ID,
		item.ProductName, item.ProductSKU, item.QuantityOrdered, item.QuantityReceived,
		item.UnitOfMeasure, item.UnitCost, item.DiscountPercent, item.DiscountAmount,
		item.TaxPercent, item.TaxAmount, item.LineTotal, item.Notes, item.UpdatedAt,
	)
	return err
}

// DeleteItem deletes a purchase order item
func (r *PurchaseRepository) DeleteItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM purchase_order_items WHERE organization_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

// DeleteItemsByOrderID deletes all items for a purchase order
func (r *PurchaseRepository) DeleteItemsByOrderID(ctx context.Context, orgID uuid.UUID, poID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "DELETE FROM purchase_order_items WHERE organization_id = $1 AND purchase_order_id = $2"
	_, err := r.db.Pool.Exec(ctx, query, orgID, poID)
	return err
}

// UpdateOrderStatus updates the status of a purchase order
func (r *PurchaseRepository) UpdateOrderStatus(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_orders
		SET status = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, poID, status, time.Now())
	return err
}

// UpdatePaymentStatus updates the payment status of a purchase order
func (r *PurchaseRepository) UpdatePaymentStatus(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, paymentStatus string, paidAmount float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_orders
		SET payment_status = $3, paid_amount = $4, updated_at = $5
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, poID, paymentStatus, paidAmount, time.Now())
	return err
}

// ApproveOrder marks a purchase order as approved
func (r *PurchaseRepository) ApproveOrder(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, approvedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	now := time.Now()
	query := `
		UPDATE purchase_orders
		SET status = $3, approved_by = $4, approved_at = $5, updated_at = $5
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, poID, "approved", approvedBy, now)
	return err
}

// RecordReceipt records a quantity received for an item
func (r *PurchaseRepository) RecordReceipt(ctx context.Context, orgID uuid.UUID, poID uuid.UUID, itemID uuid.UUID, quantityReceived float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE purchase_order_items
		SET quantity_received = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND purchase_order_id = $5
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, itemID, quantityReceived, time.Now(), poID)
	return err
}

// GetReceivedQuantity returns total received quantity for an item
func (r *PurchaseRepository) GetReceivedQuantity(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID) (float64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := `
		SELECT COALESCE(quantity_received, 0)
		FROM purchase_order_items
		WHERE organization_id = $1 AND id = $2
	`

	var received float64
	err := r.db.Pool.QueryRow(ctx, query, orgID, itemID).Scan(&received)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	return received, err
}
