package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/payables"
)

// PayablesRepository implements the payables repository interface
type PayablesRepository struct {
	db *sql.DB
}

// NewPayablesRepository creates a new payables repository
func NewPayablesRepository(db *sql.DB) payables.Repository {
	return &PayablesRepository{db: db}
}

// CreateVendorBill creates a new vendor bill
func (r *PayablesRepository) CreateVendorBill(ctx context.Context, bill *payables.VendorBill) error {
	attachments, err := json.Marshal(bill.Attachments)
	if err != nil {
		attachments = []byte("[]")
	}

	metadata, err := json.Marshal(bill.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO vendor_bills (
			id, organization_id, bill_number, vendor_bill_number, supplier_id,
			bill_date, due_date, payment_terms, accounting_period_id,
			subtotal, tax_amount, total_amount, paid_amount, balance_due,
			status, description, notes, memo, attachments, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	_, err = r.db.ExecContext(ctx, query,
		bill.ID, bill.OrganizationID, bill.BillNumber, bill.VendorBillNumber, bill.SupplierID,
		bill.BillDate, bill.DueDate, bill.PaymentTerms, bill.AccountingPeriodID,
		bill.Subtotal, bill.TaxAmount, bill.TotalAmount, bill.PaidAmount, bill.BalanceDue,
		bill.Status, bill.Description, bill.Notes, bill.Memo, attachments, metadata,
		bill.CreatedAt, bill.UpdatedAt, bill.CreatedBy, bill.UpdatedBy,
	)

	return err
}

// GetVendorBillByID retrieves a vendor bill by ID
func (r *PayablesRepository) GetVendorBillByID(ctx context.Context, organizationID, billID uuid.UUID) (*payables.VendorBill, error) {
	bill := &payables.VendorBill{}

	query := `
		SELECT id, organization_id, bill_number, vendor_bill_number, supplier_id,
		       bill_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, purchase_order_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_bills
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var attachments json.RawMessage
	var metadata json.RawMessage

	err := r.db.QueryRowContext(ctx, query, billID, organizationID).Scan(
		&bill.ID, &bill.OrganizationID, &bill.BillNumber, &bill.VendorBillNumber, &bill.SupplierID,
		&bill.BillDate, &bill.DueDate, &bill.PaymentTerms, &bill.AccountingPeriodID,
		&bill.Subtotal, &bill.TaxAmount, &bill.TotalAmount, &bill.PaidAmount, &bill.BalanceDue,
		&bill.Status, &bill.JournalEntryID, &bill.IsPosted, &bill.PurchaseOrderID,
		&bill.Description, &bill.Notes, &bill.Memo, &attachments, &metadata,
		&bill.CreatedAt, &bill.UpdatedAt, &bill.CreatedBy, &bill.UpdatedBy, &bill.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("bill not found")
		}
		return nil, err
	}

	bill.Attachments = attachments
	bill.Metadata = metadata
	return bill, nil
}

// GetVendorBillByNumber retrieves a vendor bill by bill number
func (r *PayablesRepository) GetVendorBillByNumber(ctx context.Context, organizationID uuid.UUID, billNumber string) (*payables.VendorBill, error) {
	bill := &payables.VendorBill{}

	query := `
		SELECT id, organization_id, bill_number, vendor_bill_number, supplier_id,
		       bill_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, purchase_order_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_bills
		WHERE bill_number = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var attachments json.RawMessage
	var metadata json.RawMessage

	err := r.db.QueryRowContext(ctx, query, billNumber, organizationID).Scan(
		&bill.ID, &bill.OrganizationID, &bill.BillNumber, &bill.VendorBillNumber, &bill.SupplierID,
		&bill.BillDate, &bill.DueDate, &bill.PaymentTerms, &bill.AccountingPeriodID,
		&bill.Subtotal, &bill.TaxAmount, &bill.TotalAmount, &bill.PaidAmount, &bill.BalanceDue,
		&bill.Status, &bill.JournalEntryID, &bill.IsPosted, &bill.PurchaseOrderID,
		&bill.Description, &bill.Notes, &bill.Memo, &attachments, &metadata,
		&bill.CreatedAt, &bill.UpdatedAt, &bill.CreatedBy, &bill.UpdatedBy, &bill.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("bill not found")
		}
		return nil, err
	}

	bill.Attachments = attachments
	bill.Metadata = metadata
	return bill, nil
}

// ListVendorBills retrieves a list of vendor bills with filtering
func (r *PayablesRepository) ListVendorBills(ctx context.Context, organizationID uuid.UUID, opts *payables.VendorBillFilterOptions) ([]payables.VendorBill, int64, error) {
	bills := make([]payables.VendorBill, 0)

	whereConditions := []string{"organization_id = $1"}
	args := []interface{}{organizationID}
	argCount := 2

	if !opts.DeletedIncluded {
		whereConditions = append(whereConditions, "deleted_at IS NULL")
	}

	if opts.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *opts.Status)
		argCount++
	}

	if opts.SupplierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("supplier_id = $%d", argCount))
		args = append(args, *opts.SupplierID)
		argCount++
	}

	if opts.BillDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("bill_date >= $%d", argCount))
		args = append(args, *opts.BillDateFrom)
		argCount++
	}

	if opts.BillDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("bill_date <= $%d", argCount))
		args = append(args, *opts.BillDateTo)
		argCount++
	}

	if opts.DueDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("due_date >= $%d", argCount))
		args = append(args, *opts.DueDateFrom)
		argCount++
	}

	if opts.DueDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("due_date <= $%d", argCount))
		args = append(args, *opts.DueDateTo)
		argCount++
	}

	if opts.IsPosted != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_posted = $%d", argCount))
		args = append(args, *opts.IsPosted)
		argCount++
	}

	if opts.PeriodID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("accounting_period_id = $%d", argCount))
		args = append(args, *opts.PeriodID)
		argCount++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vendor_bills WHERE %s", whereClause)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, organization_id, bill_number, vendor_bill_number, supplier_id,
		       bill_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, purchase_order_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_bills
		WHERE %s
		ORDER BY bill_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		bill := payables.VendorBill{}
		var attachments json.RawMessage
		var metadata json.RawMessage

		err := rows.Scan(
			&bill.ID, &bill.OrganizationID, &bill.BillNumber, &bill.VendorBillNumber, &bill.SupplierID,
			&bill.BillDate, &bill.DueDate, &bill.PaymentTerms, &bill.AccountingPeriodID,
			&bill.Subtotal, &bill.TaxAmount, &bill.TotalAmount, &bill.PaidAmount, &bill.BalanceDue,
			&bill.Status, &bill.JournalEntryID, &bill.IsPosted, &bill.PurchaseOrderID,
			&bill.Description, &bill.Notes, &bill.Memo, &attachments, &metadata,
			&bill.CreatedAt, &bill.UpdatedAt, &bill.CreatedBy, &bill.UpdatedBy, &bill.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		bill.Attachments = attachments
		bill.Metadata = metadata
		bills = append(bills, bill)
	}

	return bills, totalCount, nil
}

// UpdateVendorBill updates a vendor bill
func (r *PayablesRepository) UpdateVendorBill(ctx context.Context, bill *payables.VendorBill) error {
	metadata, err := json.Marshal(bill.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE vendor_bills
		SET vendor_bill_number = $1, bill_date = $2, due_date = $3, payment_terms = $4,
		    subtotal = $5, tax_amount = $6, total_amount = $7, paid_amount = $8, balance_due = $9,
		    status = $10, description = $11, notes = $12, memo = $13, metadata = $14,
		    updated_at = $15, updated_by = $16
		WHERE id = $17 AND organization_id = $18
	`

	result, err := r.db.ExecContext(ctx, query,
		bill.VendorBillNumber, bill.BillDate, bill.DueDate, bill.PaymentTerms,
		bill.Subtotal, bill.TaxAmount, bill.TotalAmount, bill.PaidAmount, bill.BalanceDue,
		bill.Status, bill.Description, bill.Notes, bill.Memo, metadata,
		bill.UpdatedAt, bill.UpdatedBy,
		bill.ID, bill.OrganizationID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("bill not found")
	}

	return nil
}

// DeleteVendorBill soft deletes a vendor bill
func (r *PayablesRepository) DeleteVendorBill(ctx context.Context, organizationID, billID uuid.UUID) error {
	query := `UPDATE vendor_bills SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
	result, err := r.db.ExecContext(ctx, query, billID, organizationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("bill not found")
	}

	return nil
}

// CreateVendorBillLine creates a new vendor bill line
func (r *PayablesRepository) CreateVendorBillLine(ctx context.Context, line *payables.VendorBillLine) error {
	metadata, err := json.Marshal(line.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO vendor_bill_lines (
			id, organization_id, vendor_bill_id, line_number, expense_account_id,
			description, quantity, unit_price, amount, location_id, department,
			project_code, tax_code, tax_amount, product_id, metadata,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err = r.db.ExecContext(ctx, query,
		line.ID, line.OrganizationID, line.VendorBillID, line.LineNumber, line.ExpenseAccountID,
		line.Description, line.Quantity, line.UnitPrice, line.Amount, line.LocationID, line.Department,
		line.ProjectCode, line.TaxCode, line.TaxAmount, line.ProductID, metadata,
		line.CreatedAt, line.UpdatedAt,
	)

	return err
}

// GetVendorBillLineByID retrieves a vendor bill line by ID
func (r *PayablesRepository) GetVendorBillLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*payables.VendorBillLine, error) {
	line := &payables.VendorBillLine{}

	query := `
		SELECT id, organization_id, vendor_bill_id, line_number, expense_account_id,
		       description, quantity, unit_price, amount, location_id, department,
		       project_code, tax_code, tax_amount, product_id, metadata,
		       created_at, updated_at, deleted_at
		FROM vendor_bill_lines
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata json.RawMessage

	err := r.db.QueryRowContext(ctx, query, lineID, organizationID).Scan(
		&line.ID, &line.OrganizationID, &line.VendorBillID, &line.LineNumber, &line.ExpenseAccountID,
		&line.Description, &line.Quantity, &line.UnitPrice, &line.Amount, &line.LocationID, &line.Department,
		&line.ProjectCode, &line.TaxCode, &line.TaxAmount, &line.ProductID, &metadata,
		&line.CreatedAt, &line.UpdatedAt, &line.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("line not found")
		}
		return nil, err
	}

	line.Metadata = metadata
	return line, nil
}

// ListVendorBillLines retrieves all lines for a bill
func (r *PayablesRepository) ListVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) ([]payables.VendorBillLine, error) {
	lines := make([]payables.VendorBillLine, 0)

	query := `
		SELECT id, organization_id, vendor_bill_id, line_number, expense_account_id,
		       description, quantity, unit_price, amount, location_id, department,
		       project_code, tax_code, tax_amount, product_id, metadata,
		       created_at, updated_at, deleted_at
		FROM vendor_bill_lines
		WHERE vendor_bill_id = $1 AND organization_id = $2 AND deleted_at IS NULL
		ORDER BY line_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, billID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		line := payables.VendorBillLine{}
		var metadata json.RawMessage

		err := rows.Scan(
			&line.ID, &line.OrganizationID, &line.VendorBillID, &line.LineNumber, &line.ExpenseAccountID,
			&line.Description, &line.Quantity, &line.UnitPrice, &line.Amount, &line.LocationID, &line.Department,
			&line.ProjectCode, &line.TaxCode, &line.TaxAmount, &line.ProductID, &metadata,
			&line.CreatedAt, &line.UpdatedAt, &line.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		line.Metadata = metadata
		lines = append(lines, line)
	}

	return lines, nil
}

// UpdateVendorBillLine updates a vendor bill line
func (r *PayablesRepository) UpdateVendorBillLine(ctx context.Context, line *payables.VendorBillLine) error {
	metadata, err := json.Marshal(line.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE vendor_bill_lines
		SET line_number = $1, description = $2, quantity = $3, unit_price = $4, amount = $5,
		    location_id = $6, department = $7, project_code = $8, tax_code = $9, tax_amount = $10,
		    metadata = $11, updated_at = $12
		WHERE id = $13 AND organization_id = $14
	`

	result, err := r.db.ExecContext(ctx, query,
		line.LineNumber, line.Description, line.Quantity, line.UnitPrice, line.Amount,
		line.LocationID, line.Department, line.ProjectCode, line.TaxCode, line.TaxAmount,
		metadata, line.UpdatedAt,
		line.ID, line.OrganizationID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("line not found")
	}

	return nil
}

// DeleteVendorBillLine soft deletes a vendor bill line
func (r *PayablesRepository) DeleteVendorBillLine(ctx context.Context, organizationID, lineID uuid.UUID) error {
	query := `UPDATE vendor_bill_lines SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
	result, err := r.db.ExecContext(ctx, query, lineID, organizationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("line not found")
	}

	return nil
}

// DeleteVendorBillLines deletes all lines for a bill
func (r *PayablesRepository) DeleteVendorBillLines(ctx context.Context, organizationID, billID uuid.UUID) error {
	query := `UPDATE vendor_bill_lines SET deleted_at = NOW() WHERE vendor_bill_id = $1 AND organization_id = $2`
	_, err := r.db.ExecContext(ctx, query, billID, organizationID)
	return err
}

// CreateVendorPayment creates a new vendor payment
func (r *PayablesRepository) CreateVendorPayment(ctx context.Context, payment *payables.VendorPayment) error {
	metadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO vendor_payments (
			id, organization_id, payment_number, supplier_id, payment_date,
			payment_method, reference_number, payment_amount, bank_account_id,
			accounting_period_id, is_posted, memo, notes, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err = r.db.ExecContext(ctx, query,
		payment.ID, payment.OrganizationID, payment.PaymentNumber, payment.SupplierID, payment.PaymentDate,
		payment.PaymentMethod, payment.ReferenceNumber, payment.PaymentAmount, payment.BankAccountID,
		payment.AccountingPeriodID, payment.IsPosted, payment.Memo, payment.Notes, metadata,
		payment.CreatedAt, payment.UpdatedAt, payment.CreatedBy, payment.UpdatedBy,
	)

	return err
}

// GetVendorPaymentByID retrieves a vendor payment by ID
func (r *PayablesRepository) GetVendorPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*payables.VendorPayment, error) {
	payment := &payables.VendorPayment{}

	query := `
		SELECT id, organization_id, payment_number, supplier_id, payment_date,
		       payment_method, reference_number, payment_amount, bank_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_payments
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata json.RawMessage

	err := r.db.QueryRowContext(ctx, query, paymentID, organizationID).Scan(
		&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.SupplierID, &payment.PaymentDate,
		&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.BankAccountID,
		&payment.AccountingPeriodID, &payment.JournalEntryID, &payment.IsPosted, &payment.Memo, &payment.Notes, &metadata,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	payment.Metadata = metadata
	return payment, nil
}

// GetVendorPaymentByNumber retrieves a vendor payment by payment number
func (r *PayablesRepository) GetVendorPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*payables.VendorPayment, error) {
	payment := &payables.VendorPayment{}

	query := `
		SELECT id, organization_id, payment_number, supplier_id, payment_date,
		       payment_method, reference_number, payment_amount, bank_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_payments
		WHERE payment_number = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata json.RawMessage

	err := r.db.QueryRowContext(ctx, query, paymentNumber, organizationID).Scan(
		&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.SupplierID, &payment.PaymentDate,
		&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.BankAccountID,
		&payment.AccountingPeriodID, &payment.JournalEntryID, &payment.IsPosted, &payment.Memo, &payment.Notes, &metadata,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	payment.Metadata = metadata
	return payment, nil
}

// ListVendorPayments retrieves a list of vendor payments with filtering
func (r *PayablesRepository) ListVendorPayments(ctx context.Context, organizationID uuid.UUID, opts *payables.VendorPaymentFilterOptions) ([]payables.VendorPayment, int64, error) {
	payments := make([]payables.VendorPayment, 0)

	whereConditions := []string{"organization_id = $1"}
	args := []interface{}{organizationID}
	argCount := 2

	if !opts.DeletedIncluded {
		whereConditions = append(whereConditions, "deleted_at IS NULL")
	}

	if opts.SupplierID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("supplier_id = $%d", argCount))
		args = append(args, *opts.SupplierID)
		argCount++
	}

	if opts.PaymentMethod != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("payment_method = $%d", argCount))
		args = append(args, *opts.PaymentMethod)
		argCount++
	}

	if opts.PaymentDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("payment_date >= $%d", argCount))
		args = append(args, *opts.PaymentDateFrom)
		argCount++
	}

	if opts.PaymentDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("payment_date <= $%d", argCount))
		args = append(args, *opts.PaymentDateTo)
		argCount++
	}

	if opts.IsPosted != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("is_posted = $%d", argCount))
		args = append(args, *opts.IsPosted)
		argCount++
	}

	if opts.PeriodID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("accounting_period_id = $%d", argCount))
		args = append(args, *opts.PeriodID)
		argCount++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vendor_payments WHERE %s", whereClause)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, organization_id, payment_number, supplier_id, payment_date,
		       payment_method, reference_number, payment_amount, bank_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_payments
		WHERE %s
		ORDER BY payment_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		payment := payables.VendorPayment{}
		var metadata json.RawMessage

		err := rows.Scan(
			&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.SupplierID, &payment.PaymentDate,
			&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.BankAccountID,
			&payment.AccountingPeriodID, &payment.JournalEntryID, &payment.IsPosted, &payment.Memo, &payment.Notes, &metadata,
			&payment.CreatedAt, &payment.UpdatedAt, &payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		payment.Metadata = metadata
		payments = append(payments, payment)
	}

	return payments, totalCount, nil
}

// UpdateVendorPayment updates a vendor payment
func (r *PayablesRepository) UpdateVendorPayment(ctx context.Context, payment *payables.VendorPayment) error {
	metadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE vendor_payments
		SET payment_date = $1, payment_method = $2, reference_number = $3, payment_amount = $4,
		    bank_account_id = $5, memo = $6, notes = $7, metadata = $8,
		    updated_at = $9, updated_by = $10
		WHERE id = $11 AND organization_id = $12
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.PaymentDate, payment.PaymentMethod, payment.ReferenceNumber, payment.PaymentAmount,
		payment.BankAccountID, payment.Memo, payment.Notes, metadata,
		payment.UpdatedAt, payment.UpdatedBy,
		payment.ID, payment.OrganizationID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("payment not found")
	}

	return nil
}

// DeleteVendorPayment soft deletes a vendor payment
func (r *PayablesRepository) DeleteVendorPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	query := `UPDATE vendor_payments SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
	result, err := r.db.ExecContext(ctx, query, paymentID, organizationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("payment not found")
	}

	return nil
}

// CreateVendorPaymentApplication creates a new payment application
func (r *PayablesRepository) CreateVendorPaymentApplication(ctx context.Context, app *payables.VendorPaymentApplication) error {
	query := `
		INSERT INTO vendor_payment_applications (
			id, organization_id, vendor_payment_id, vendor_bill_id, applied_amount, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		app.ID, app.OrganizationID, app.VendorPaymentID, app.VendorBillID, app.AppliedAmount, app.CreatedAt,
	)

	return err
}

// GetVendorPaymentApplicationByID retrieves a payment application by ID
func (r *PayablesRepository) GetVendorPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*payables.VendorPaymentApplication, error) {
	app := &payables.VendorPaymentApplication{}

	query := `
		SELECT id, organization_id, vendor_payment_id, vendor_bill_id, applied_amount, created_at, deleted_at
		FROM vendor_payment_applications
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, appID, organizationID).Scan(
		&app.ID, &app.OrganizationID, &app.VendorPaymentID, &app.VendorBillID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment application not found")
		}
		return nil, err
	}

	return app, nil
}

// ListVendorPaymentApplications retrieves all applications for a payment
func (r *PayablesRepository) ListVendorPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]payables.VendorPaymentApplication, error) {
	apps := make([]payables.VendorPaymentApplication, 0)

	query := `
		SELECT id, organization_id, vendor_payment_id, vendor_bill_id, applied_amount, created_at, deleted_at
		FROM vendor_payment_applications
		WHERE vendor_payment_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, paymentID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		app := payables.VendorPaymentApplication{}
		err := rows.Scan(
			&app.ID, &app.OrganizationID, &app.VendorPaymentID, &app.VendorBillID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// ListPaymentApplicationsByBill retrieves all applications for a bill
func (r *PayablesRepository) ListPaymentApplicationsByBill(ctx context.Context, organizationID, billID uuid.UUID) ([]payables.VendorPaymentApplication, error) {
	apps := make([]payables.VendorPaymentApplication, 0)

	query := `
		SELECT id, organization_id, vendor_payment_id, vendor_bill_id, applied_amount, created_at, deleted_at
		FROM vendor_payment_applications
		WHERE vendor_bill_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, billID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		app := payables.VendorPaymentApplication{}
		err := rows.Scan(
			&app.ID, &app.OrganizationID, &app.VendorPaymentID, &app.VendorBillID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// DeleteVendorPaymentApplication soft deletes a payment application
func (r *PayablesRepository) DeleteVendorPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error {
	query := `UPDATE vendor_payment_applications SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
	result, err := r.db.ExecContext(ctx, query, appID, organizationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("payment application not found")
	}

	return nil
}

// GetVendorAgingReport generates an aging report for AP
func (r *PayablesRepository) GetVendorAgingReport(ctx context.Context, organizationID uuid.UUID) ([]payables.VendorAgingReport, error) {
	reports := make([]payables.VendorAgingReport, 0)

	query := `
		WITH supplier_info AS (
			SELECT DISTINCT vb.supplier_id, s.name as supplier_name
			FROM vendor_bills vb
			JOIN suppliers s ON vb.supplier_id = s.id
			WHERE vb.organization_id = $1 AND vb.deleted_at IS NULL AND vb.status != 'paid' AND vb.status != 'void'
		),
		aging_buckets AS (
			SELECT
				vb.supplier_id,
				CASE
					WHEN DATE(vb.due_date) >= CURRENT_DATE THEN vb.balance_due ELSE 0
				END as current,
				CASE
					WHEN DATE(vb.due_date) < CURRENT_DATE AND DATE(vb.due_date) >= CURRENT_DATE - INTERVAL '30 days' THEN vb.balance_due ELSE 0
				END as days_30,
				CASE
					WHEN DATE(vb.due_date) < CURRENT_DATE - INTERVAL '30 days' AND DATE(vb.due_date) >= CURRENT_DATE - INTERVAL '60 days' THEN vb.balance_due ELSE 0
				END as days_60,
				CASE
					WHEN DATE(vb.due_date) < CURRENT_DATE - INTERVAL '60 days' THEN vb.balance_due ELSE 0
				END as days_90_plus
			FROM vendor_bills vb
			WHERE vb.organization_id = $1 AND vb.deleted_at IS NULL
		)
		SELECT
			si.supplier_id,
			si.supplier_name,
			COALESCE(SUM(ab.current), 0) as current,
			COALESCE(SUM(ab.days_30), 0) as days_30,
			COALESCE(SUM(ab.days_60), 0) as days_60,
			COALESCE(SUM(ab.days_90_plus), 0) as days_90_plus
		FROM supplier_info si
		LEFT JOIN aging_buckets ab ON si.supplier_id = ab.supplier_id
		GROUP BY si.supplier_id, si.supplier_name
		ORDER BY (COALESCE(SUM(ab.current), 0) + COALESCE(SUM(ab.days_30), 0) + COALESCE(SUM(ab.days_60), 0) + COALESCE(SUM(ab.days_90_plus), 0)) DESC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		report := payables.VendorAgingReport{}
		err := rows.Scan(
			&report.SupplierID, &report.SupplierName, &report.Current, &report.Days30, &report.Days60, &report.Days90Plus,
		)

		if err != nil {
			return nil, err
		}

		report.TotalDue = report.Current + report.Days30 + report.Days60 + report.Days90Plus
		reports = append(reports, report)
	}

	return reports, nil
}

// GetBillsByDueDate retrieves bills due on or before a specific date
func (r *PayablesRepository) GetBillsByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]payables.VendorBill, error) {
	bills := make([]payables.VendorBill, 0)

	query := `
		SELECT id, organization_id, bill_number, vendor_bill_number, supplier_id,
		       bill_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, purchase_order_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM vendor_bills
		WHERE organization_id = $1 AND due_date < $2 AND deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue')
		ORDER BY due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		bill := payables.VendorBill{}
		var attachments json.RawMessage
		var metadata json.RawMessage

		err := rows.Scan(
			&bill.ID, &bill.OrganizationID, &bill.BillNumber, &bill.VendorBillNumber, &bill.SupplierID,
			&bill.BillDate, &bill.DueDate, &bill.PaymentTerms, &bill.AccountingPeriodID,
			&bill.Subtotal, &bill.TaxAmount, &bill.TotalAmount, &bill.PaidAmount, &bill.BalanceDue,
			&bill.Status, &bill.JournalEntryID, &bill.IsPosted, &bill.PurchaseOrderID,
			&bill.Description, &bill.Notes, &bill.Memo, &attachments, &metadata,
			&bill.CreatedAt, &bill.UpdatedAt, &bill.CreatedBy, &bill.UpdatedBy, &bill.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		bill.Attachments = attachments
		bill.Metadata = metadata
		bills = append(bills, bill)
	}

	return bills, nil
}

// GetSupplierBalance retrieves the total balance owed to a supplier
func (r *PayablesRepository) GetSupplierBalance(ctx context.Context, organizationID, supplierID uuid.UUID) (float64, error) {
	var balance float64

	query := `
		SELECT COALESCE(SUM(balance_due), 0)
		FROM vendor_bills
		WHERE organization_id = $1 AND supplier_id = $2 AND deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue')
	`

	err := r.db.QueryRowContext(ctx, query, organizationID, supplierID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
