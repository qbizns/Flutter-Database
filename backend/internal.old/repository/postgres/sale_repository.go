package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/sales"
)

type SaleRepository struct {
	db *DB
}

func NewSaleRepository(db *DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) List(ctx context.Context, orgID uuid.UUID, filters sales.SaleFilters) ([]sales.Sale, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_number, reference_number, transaction_type,
		       customer_id, cashier_id, subtotal, tax_amount, discount_amount,
		       total_amount, paid_amount, change_amount, outstanding_amount,
		       payment_status, discount_type, discount_value, discount_reason,
		       transaction_date, completed_at, notes, internal_notes,
		       custom_fields, metadata, created_at, updated_at, created_by, updated_by
		FROM sales
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (sale_number ILIKE $%d OR reference_number ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.CashierID != nil {
		argCount++
		query += fmt.Sprintf(" AND cashier_id = $%d", argCount)
		args = append(args, *filters.CashierID)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, *filters.PaymentStatus)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, *filters.TransactionType)
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

	query += " ORDER BY transaction_date DESC, created_at DESC"

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

	var saleList []sales.Sale
	for rows.Next() {
		var s sales.Sale
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.SaleNumber, &s.ReferenceNumber, &s.TransactionType,
			&s.CustomerID, &s.CashierID, &s.Subtotal, &s.TaxAmount, &s.DiscountAmount,
			&s.TotalAmount, &s.PaidAmount, &s.ChangeAmount, &s.OutstandingAmount,
			&s.PaymentStatus, &s.DiscountType, &s.DiscountValue, &s.DiscountReason,
			&s.TransactionDate, &s.CompletedAt, &s.Notes, &s.InternalNotes,
			&s.CustomFields, &s.Metadata, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		saleList = append(saleList, s)
	}

	return saleList, rows.Err()
}

func (r *SaleRepository) Count(ctx context.Context, orgID uuid.UUID, filters sales.SaleFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM sales WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (sale_number ILIKE $%d OR reference_number ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.CashierID != nil {
		argCount++
		query += fmt.Sprintf(" AND cashier_id = $%d", argCount)
		args = append(args, *filters.CashierID)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, *filters.PaymentStatus)
	}

	if filters.TransactionType != nil {
		argCount++
		query += fmt.Sprintf(" AND transaction_type = $%d", argCount)
		args = append(args, *filters.TransactionType)
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

func (r *SaleRepository) Create(ctx context.Context, sale *sales.Sale) error {
	if err := r.db.SetOrganizationContext(ctx, sale.OrganizationID.String()); err != nil {
		return err
	}

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert sale
	query := `
		INSERT INTO sales (
			id, organization_id, sale_number, reference_number, transaction_type,
			customer_id, cashier_id, subtotal, tax_amount, discount_amount,
			total_amount, paid_amount, change_amount, outstanding_amount,
			payment_status, discount_type, discount_value, discount_reason,
			transaction_date, completed_at, notes, internal_notes,
			custom_fields, metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err = tx.Exec(ctx, query,
		sale.ID, sale.OrganizationID, sale.SaleNumber, sale.ReferenceNumber, sale.TransactionType,
		sale.CustomerID, sale.CashierID, sale.Subtotal, sale.TaxAmount, sale.DiscountAmount,
		sale.TotalAmount, sale.PaidAmount, sale.ChangeAmount, sale.OutstandingAmount,
		sale.PaymentStatus, sale.DiscountType, sale.DiscountValue, sale.DiscountReason,
		sale.TransactionDate, sale.CompletedAt, sale.Notes, sale.InternalNotes,
		sale.CustomFields, sale.Metadata, sale.CreatedAt, sale.UpdatedAt, sale.CreatedBy,
	)
	if err != nil {
		return err
	}

	// Insert sale items
	if len(sale.Items) > 0 {
		itemQuery := `
			INSERT INTO sale_items (
				id, sale_id, organization_id, product_id, product_name, product_sku,
				quantity, unit, unit_price, cost_price, subtotal, tax_rate, tax_amount,
				discount_amount, total, discount_type, discount_value, notes,
				custom_fields, metadata, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
				$15, $16, $17, $18, $19, $20, $21, $22
			)
		`

		for _, item := range sale.Items {
			_, err = tx.Exec(ctx, itemQuery,
				item.ID, item.SaleID, item.OrganizationID, item.ProductID, item.ProductName, item.ProductSKU,
				item.Quantity, item.Unit, item.UnitPrice, item.CostPrice, item.Subtotal, item.TaxRate, item.TaxAmount,
				item.DiscountAmount, item.Total, item.DiscountType, item.DiscountValue, item.Notes,
				item.CustomFields, item.Metadata, item.CreatedAt, item.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *SaleRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*sales.Sale, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_number, reference_number, transaction_type,
		       customer_id, cashier_id, subtotal, tax_amount, discount_amount,
		       total_amount, paid_amount, change_amount, outstanding_amount,
		       payment_status, discount_type, discount_value, discount_reason,
		       transaction_date, completed_at, notes, internal_notes,
		       custom_fields, metadata, created_at, updated_at, created_by, updated_by
		FROM sales
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var s sales.Sale
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&s.ID, &s.OrganizationID, &s.SaleNumber, &s.ReferenceNumber, &s.TransactionType,
		&s.CustomerID, &s.CashierID, &s.Subtotal, &s.TaxAmount, &s.DiscountAmount,
		&s.TotalAmount, &s.PaidAmount, &s.ChangeAmount, &s.OutstandingAmount,
		&s.PaymentStatus, &s.DiscountType, &s.DiscountValue, &s.DiscountReason,
		&s.TransactionDate, &s.CompletedAt, &s.Notes, &s.InternalNotes,
		&s.CustomFields, &s.Metadata, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SaleRepository) GetBySaleNumber(ctx context.Context, orgID uuid.UUID, saleNumber string) (*sales.Sale, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_number, reference_number, transaction_type,
		       customer_id, cashier_id, subtotal, tax_amount, discount_amount,
		       total_amount, paid_amount, change_amount, outstanding_amount,
		       payment_status, discount_type, discount_value, discount_reason,
		       transaction_date, completed_at, notes, internal_notes,
		       custom_fields, metadata, created_at, updated_at, created_by, updated_by
		FROM sales
		WHERE organization_id = $1 AND sale_number = $2 AND deleted_at IS NULL
	`

	var s sales.Sale
	err := r.db.Pool.QueryRow(ctx, query, orgID, saleNumber).Scan(
		&s.ID, &s.OrganizationID, &s.SaleNumber, &s.ReferenceNumber, &s.TransactionType,
		&s.CustomerID, &s.CashierID, &s.Subtotal, &s.TaxAmount, &s.DiscountAmount,
		&s.TotalAmount, &s.PaidAmount, &s.ChangeAmount, &s.OutstandingAmount,
		&s.PaymentStatus, &s.DiscountType, &s.DiscountValue, &s.DiscountReason,
		&s.TransactionDate, &s.CompletedAt, &s.Notes, &s.InternalNotes,
		&s.CustomFields, &s.Metadata, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SaleRepository) Update(ctx context.Context, sale *sales.Sale) error {
	if err := r.db.SetOrganizationContext(ctx, sale.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE sales SET
			reference_number = $3, transaction_type = $4, customer_id = $5,
			cashier_id = $6, subtotal = $7, tax_amount = $8, discount_amount = $9,
			total_amount = $10, paid_amount = $11, change_amount = $12,
			outstanding_amount = $13, payment_status = $14, discount_type = $15,
			discount_value = $16, discount_reason = $17, transaction_date = $18,
			completed_at = $19, notes = $20, internal_notes = $21,
			custom_fields = $22, metadata = $23, updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		sale.OrganizationID, sale.ID,
		sale.ReferenceNumber, sale.TransactionType, sale.CustomerID,
		sale.CashierID, sale.Subtotal, sale.TaxAmount, sale.DiscountAmount,
		sale.TotalAmount, sale.PaidAmount, sale.ChangeAmount,
		sale.OutstandingAmount, sale.PaymentStatus, sale.DiscountType,
		sale.DiscountValue, sale.DiscountReason, sale.TransactionDate,
		sale.CompletedAt, sale.Notes, sale.InternalNotes,
		sale.CustomFields, sale.Metadata, sale.UpdatedAt, sale.UpdatedBy,
	)
	return err
}

func (r *SaleRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE sales
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *SaleRepository) CreateItem(ctx context.Context, item *sales.SaleItem) error {
	query := `
		INSERT INTO sale_items (
			id, sale_id, organization_id, product_id, product_name, product_sku,
			quantity, unit, unit_price, cost_price, subtotal, tax_rate, tax_amount,
			discount_amount, total, discount_type, discount_value, notes,
			custom_fields, metadata, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.SaleID, item.OrganizationID, item.ProductID, item.ProductName, item.ProductSKU,
		item.Quantity, item.Unit, item.UnitPrice, item.CostPrice, item.Subtotal, item.TaxRate, item.TaxAmount,
		item.DiscountAmount, item.Total, item.DiscountType, item.DiscountValue, item.Notes,
		item.CustomFields, item.Metadata, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *SaleRepository) GetItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) (*sales.SaleItem, error) {
	query := `
		SELECT id, sale_id, organization_id, product_id, product_name, product_sku,
		       quantity, unit, unit_price, cost_price, subtotal, tax_rate, tax_amount,
		       discount_amount, total, discount_type, discount_value, notes,
		       custom_fields, metadata, created_at, updated_at
		FROM sale_items
		WHERE sale_id = $1 AND id = $2
	`

	var item sales.SaleItem
	err := r.db.Pool.QueryRow(ctx, query, saleID, itemID).Scan(
		&item.ID, &item.SaleID, &item.OrganizationID, &item.ProductID, &item.ProductName, &item.ProductSKU,
		&item.Quantity, &item.Unit, &item.UnitPrice, &item.CostPrice, &item.Subtotal, &item.TaxRate, &item.TaxAmount,
		&item.DiscountAmount, &item.Total, &item.DiscountType, &item.DiscountValue, &item.Notes,
		&item.CustomFields, &item.Metadata, &item.CreatedAt, &item.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SaleRepository) ListItems(ctx context.Context, saleID uuid.UUID, filters sales.SaleItemFilters) ([]sales.SaleItem, error) {
	query := `
		SELECT id, sale_id, organization_id, product_id, product_name, product_sku,
		       quantity, unit, unit_price, cost_price, subtotal, tax_rate, tax_amount,
		       discount_amount, total, discount_type, discount_value, notes,
		       custom_fields, metadata, created_at, updated_at
		FROM sale_items
		WHERE sale_id = $1
	`

	args := []interface{}{saleID}
	argCount := 1

	if filters.ProductID != nil {
		argCount++
		query += fmt.Sprintf(" AND product_id = $%d", argCount)
		args = append(args, *filters.ProductID)
	}

	query += " ORDER BY created_at ASC"

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

	var items []sales.SaleItem
	for rows.Next() {
		var item sales.SaleItem
		err := rows.Scan(
			&item.ID, &item.SaleID, &item.OrganizationID, &item.ProductID, &item.ProductName, &item.ProductSKU,
			&item.Quantity, &item.Unit, &item.UnitPrice, &item.CostPrice, &item.Subtotal, &item.TaxRate, &item.TaxAmount,
			&item.DiscountAmount, &item.Total, &item.DiscountType, &item.DiscountValue, &item.Notes,
			&item.CustomFields, &item.Metadata, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *SaleRepository) CountItems(ctx context.Context, saleID uuid.UUID) (int64, error) {
	query := "SELECT COUNT(*) FROM sale_items WHERE sale_id = $1"

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, saleID).Scan(&count)
	return count, err
}

func (r *SaleRepository) UpdateItem(ctx context.Context, item *sales.SaleItem) error {
	query := `
		UPDATE sale_items SET
			product_id = $3, product_name = $4, product_sku = $5,
			quantity = $6, unit = $7, unit_price = $8, cost_price = $9,
			subtotal = $10, tax_rate = $11, tax_amount = $12,
			discount_amount = $13, total = $14, discount_type = $15,
			discount_value = $16, notes = $17, custom_fields = $18,
			metadata = $19, updated_at = $20
		WHERE sale_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.SaleID, item.ID,
		item.ProductID, item.ProductName, item.ProductSKU,
		item.Quantity, item.Unit, item.UnitPrice, item.CostPrice,
		item.Subtotal, item.TaxRate, item.TaxAmount,
		item.DiscountAmount, item.Total, item.DiscountType,
		item.DiscountValue, item.Notes, item.CustomFields,
		item.Metadata, item.UpdatedAt,
	)
	return err
}

func (r *SaleRepository) DeleteItem(ctx context.Context, saleID uuid.UUID, itemID uuid.UUID) error {
	query := "DELETE FROM sale_items WHERE sale_id = $1 AND id = $2"
	_, err := r.db.Pool.Exec(ctx, query, saleID, itemID)
	return err
}

func (r *SaleRepository) GetWithItems(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*sales.Sale, error) {
	// Get the sale
	sale, err := r.Get(ctx, orgID, saleID)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, nil
	}

	// Get the items
	items, err := r.ListItems(ctx, saleID, sales.SaleItemFilters{})
	if err != nil {
		return nil, err
	}

	sale.Items = items
	return sale, nil
}
