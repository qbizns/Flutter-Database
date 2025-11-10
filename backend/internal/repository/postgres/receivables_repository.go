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
	"github.com/lib/pq"
	"github.com/your-org/pos-backend/internal/domain/receivables"
)

// ReceivablesRepository implements the receivables repository interface
type ReceivablesRepository struct {
	db *sql.DB
}

// NewReceivablesRepository creates a new receivables repository
func NewReceivablesRepository(db *sql.DB) receivables.Repository {
	return &ReceivablesRepository{db: db}
}

// CreateCustomerInvoice creates a new customer invoice
func (r *ReceivablesRepository) CreateCustomerInvoice(ctx context.Context, invoice *receivables.CustomerInvoice) error {
	attachments, err := json.Marshal(invoice.Attachments)
	if err != nil {
		attachments = []byte("[]")
	}

	metadata, err := json.Marshal(invoice.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO customer_invoices (
			id, organization_id, invoice_number, customer_id,
			invoice_date, due_date, payment_terms, accounting_period_id,
			subtotal, tax_amount, discount_amount, total_amount, paid_amount, balance_due,
			status, description, notes, memo, attachments, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	_, err = r.db.ExecContext(ctx, query,
		invoice.ID, invoice.OrganizationID, invoice.InvoiceNumber, invoice.CustomerID,
		invoice.InvoiceDate, invoice.DueDate, invoice.PaymentTerms, invoice.AccountingPeriodID,
		invoice.Subtotal, invoice.TaxAmount, invoice.DiscountAmount, invoice.TotalAmount, invoice.PaidAmount, invoice.BalanceDue,
		invoice.Status, invoice.Description, invoice.Notes, invoice.Memo, attachments, metadata,
		invoice.CreatedAt, invoice.UpdatedAt, invoice.CreatedBy, invoice.UpdatedBy,
	)

	return err
}

// GetCustomerInvoiceByID retrieves a customer invoice by ID
func (r *ReceivablesRepository) GetCustomerInvoiceByID(ctx context.Context, organizationID, invoiceID uuid.UUID) (*receivables.CustomerInvoice, error) {
	invoice := &receivables.CustomerInvoice{}

	query := `
		SELECT id, organization_id, invoice_number, customer_id,
		       invoice_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, discount_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, sale_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_invoices
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var attachments pq.JSONBArray
	var metadata pq.JSONBMap

	err := r.db.QueryRowContext(ctx, query, invoiceID, organizationID).Scan(
		&invoice.ID, &invoice.OrganizationID, &invoice.InvoiceNumber, &invoice.CustomerID,
		&invoice.InvoiceDate, &invoice.DueDate, &invoice.PaymentTerms, &invoice.AccountingPeriodID,
		&invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.PaidAmount, &invoice.BalanceDue,
		&invoice.Status, &invoice.JournalEntryID, &invoice.IsPosted, &invoice.SaleID,
		&invoice.Description, &invoice.Notes, &invoice.Memo, &attachments, &metadata,
		&invoice.CreatedAt, &invoice.UpdatedAt, &invoice.CreatedBy, &invoice.UpdatedBy, &invoice.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	invoice.Attachments = attachments
	invoice.Metadata = metadata
	return invoice, nil
}

// GetCustomerInvoiceByNumber retrieves a customer invoice by invoice number
func (r *ReceivablesRepository) GetCustomerInvoiceByNumber(ctx context.Context, organizationID uuid.UUID, invoiceNumber string) (*receivables.CustomerInvoice, error) {
	invoice := &receivables.CustomerInvoice{}

	query := `
		SELECT id, organization_id, invoice_number, customer_id,
		       invoice_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, discount_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, sale_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_invoices
		WHERE invoice_number = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var attachments pq.JSONBArray
	var metadata pq.JSONBMap

	err := r.db.QueryRowContext(ctx, query, invoiceNumber, organizationID).Scan(
		&invoice.ID, &invoice.OrganizationID, &invoice.InvoiceNumber, &invoice.CustomerID,
		&invoice.InvoiceDate, &invoice.DueDate, &invoice.PaymentTerms, &invoice.AccountingPeriodID,
		&invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.PaidAmount, &invoice.BalanceDue,
		&invoice.Status, &invoice.JournalEntryID, &invoice.IsPosted, &invoice.SaleID,
		&invoice.Description, &invoice.Notes, &invoice.Memo, &attachments, &metadata,
		&invoice.CreatedAt, &invoice.UpdatedAt, &invoice.CreatedBy, &invoice.UpdatedBy, &invoice.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invoice not found")
		}
		return nil, err
	}

	invoice.Attachments = attachments
	invoice.Metadata = metadata
	return invoice, nil
}

// ListCustomerInvoices retrieves a list of customer invoices with filtering
func (r *ReceivablesRepository) ListCustomerInvoices(ctx context.Context, organizationID uuid.UUID, opts *receivables.CustomerInvoiceFilterOptions) ([]receivables.CustomerInvoice, int64, error) {
	invoices := make([]receivables.CustomerInvoice, 0)

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

	if opts.CustomerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("customer_id = $%d", argCount))
		args = append(args, *opts.CustomerID)
		argCount++
	}

	if opts.InvoiceDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("invoice_date >= $%d", argCount))
		args = append(args, *opts.InvoiceDateFrom)
		argCount++
	}

	if opts.InvoiceDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("invoice_date <= $%d", argCount))
		args = append(args, *opts.InvoiceDateTo)
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customer_invoices WHERE %s", whereClause)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, organization_id, invoice_number, customer_id,
		       invoice_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, discount_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, sale_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_invoices
		WHERE %s
		ORDER BY invoice_date DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, opts.Limit, opts.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		invoice := receivables.CustomerInvoice{}
		var attachments pq.JSONBArray
		var metadata pq.JSONBMap

		err := rows.Scan(
			&invoice.ID, &invoice.OrganizationID, &invoice.InvoiceNumber, &invoice.CustomerID,
			&invoice.InvoiceDate, &invoice.DueDate, &invoice.PaymentTerms, &invoice.AccountingPeriodID,
			&invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.PaidAmount, &invoice.BalanceDue,
			&invoice.Status, &invoice.JournalEntryID, &invoice.IsPosted, &invoice.SaleID,
			&invoice.Description, &invoice.Notes, &invoice.Memo, &attachments, &metadata,
			&invoice.CreatedAt, &invoice.UpdatedAt, &invoice.CreatedBy, &invoice.UpdatedBy, &invoice.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		invoice.Attachments = attachments
		invoice.Metadata = metadata
		invoices = append(invoices, invoice)
	}

	return invoices, totalCount, nil
}

// UpdateCustomerInvoice updates a customer invoice
func (r *ReceivablesRepository) UpdateCustomerInvoice(ctx context.Context, invoice *receivables.CustomerInvoice) error {
	metadata, err := json.Marshal(invoice.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE customer_invoices
		SET invoice_date = $1, due_date = $2, payment_terms = $3,
		    subtotal = $4, tax_amount = $5, discount_amount = $6, total_amount = $7, paid_amount = $8, balance_due = $9,
		    status = $10, description = $11, notes = $12, memo = $13, metadata = $14,
		    updated_at = $15, updated_by = $16
		WHERE id = $17 AND organization_id = $18
	`

	result, err := r.db.ExecContext(ctx, query,
		invoice.InvoiceDate, invoice.DueDate, invoice.PaymentTerms,
		invoice.Subtotal, invoice.TaxAmount, invoice.DiscountAmount, invoice.TotalAmount, invoice.PaidAmount, invoice.BalanceDue,
		invoice.Status, invoice.Description, invoice.Notes, invoice.Memo, metadata,
		invoice.UpdatedAt, invoice.UpdatedBy,
		invoice.ID, invoice.OrganizationID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invoice not found")
	}

	return nil
}

// DeleteCustomerInvoice soft deletes a customer invoice
func (r *ReceivablesRepository) DeleteCustomerInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) error {
	query := `UPDATE customer_invoices SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
	result, err := r.db.ExecContext(ctx, query, invoiceID, organizationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invoice not found")
	}

	return nil
}

// CreateCustomerInvoiceLine creates a new customer invoice line
func (r *ReceivablesRepository) CreateCustomerInvoiceLine(ctx context.Context, line *receivables.CustomerInvoiceLine) error {
	metadata, err := json.Marshal(line.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO customer_invoice_lines (
			id, organization_id, customer_invoice_id, line_number, revenue_account_id,
			description, quantity, unit_price, amount, location_id, department,
			project_code, tax_code, tax_amount, product_id, metadata,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err = r.db.ExecContext(ctx, query,
		line.ID, line.OrganizationID, line.CustomerInvoiceID, line.LineNumber, line.RevenueAccountID,
		line.Description, line.Quantity, line.UnitPrice, line.Amount, line.LocationID, line.Department,
		line.ProjectCode, line.TaxCode, line.TaxAmount, line.ProductID, metadata,
		line.CreatedAt, line.UpdatedAt,
	)

	return err
}

// GetCustomerInvoiceLineByID retrieves a customer invoice line by ID
func (r *ReceivablesRepository) GetCustomerInvoiceLineByID(ctx context.Context, organizationID, lineID uuid.UUID) (*receivables.CustomerInvoiceLine, error) {
	line := &receivables.CustomerInvoiceLine{}

	query := `
		SELECT id, organization_id, customer_invoice_id, line_number, revenue_account_id,
		       description, quantity, unit_price, amount, location_id, department,
		       project_code, tax_code, tax_amount, product_id, metadata,
		       created_at, updated_at, deleted_at
		FROM customer_invoice_lines
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata pq.JSONBMap

	err := r.db.QueryRowContext(ctx, query, lineID, organizationID).Scan(
		&line.ID, &line.OrganizationID, &line.CustomerInvoiceID, &line.LineNumber, &line.RevenueAccountID,
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

// ListCustomerInvoiceLines retrieves all lines for an invoice
func (r *ReceivablesRepository) ListCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]receivables.CustomerInvoiceLine, error) {
	lines := make([]receivables.CustomerInvoiceLine, 0)

	query := `
		SELECT id, organization_id, customer_invoice_id, line_number, revenue_account_id,
		       description, quantity, unit_price, amount, location_id, department,
		       project_code, tax_code, tax_amount, product_id, metadata,
		       created_at, updated_at, deleted_at
		FROM customer_invoice_lines
		WHERE customer_invoice_id = $1 AND organization_id = $2 AND deleted_at IS NULL
		ORDER BY line_number ASC
	`

	rows, err := r.db.QueryContext(ctx, query, invoiceID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		line := receivables.CustomerInvoiceLine{}
		var metadata pq.JSONBMap

		err := rows.Scan(
			&line.ID, &line.OrganizationID, &line.CustomerInvoiceID, &line.LineNumber, &line.RevenueAccountID,
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

// UpdateCustomerInvoiceLine updates a customer invoice line
func (r *ReceivablesRepository) UpdateCustomerInvoiceLine(ctx context.Context, line *receivables.CustomerInvoiceLine) error {
	metadata, err := json.Marshal(line.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE customer_invoice_lines
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

// DeleteCustomerInvoiceLine soft deletes a customer invoice line
func (r *ReceivablesRepository) DeleteCustomerInvoiceLine(ctx context.Context, organizationID, lineID uuid.UUID) error {
	query := `UPDATE customer_invoice_lines SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
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

// DeleteCustomerInvoiceLines deletes all lines for an invoice
func (r *ReceivablesRepository) DeleteCustomerInvoiceLines(ctx context.Context, organizationID, invoiceID uuid.UUID) error {
	query := `UPDATE customer_invoice_lines SET deleted_at = NOW() WHERE customer_invoice_id = $1 AND organization_id = $2`
	_, err := r.db.ExecContext(ctx, query, invoiceID, organizationID)
	return err
}

// CreateCustomerPayment creates a new customer payment
func (r *ReceivablesRepository) CreateCustomerPayment(ctx context.Context, payment *receivables.CustomerPayment) error {
	metadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		INSERT INTO customer_payments (
			id, organization_id, payment_number, customer_id, payment_date,
			payment_method, reference_number, payment_amount, deposit_account_id,
			accounting_period_id, is_posted, memo, notes, metadata,
			created_at, updated_at, created_by, updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err = r.db.ExecContext(ctx, query,
		payment.ID, payment.OrganizationID, payment.PaymentNumber, payment.CustomerID, payment.PaymentDate,
		payment.PaymentMethod, payment.ReferenceNumber, payment.PaymentAmount, payment.DepositAccountID,
		payment.AccountingPeriodID, payment.IsPosted, payment.Memo, payment.Notes, metadata,
		payment.CreatedAt, payment.UpdatedAt, payment.CreatedBy, payment.UpdatedBy,
	)

	return err
}

// GetCustomerPaymentByID retrieves a customer payment by ID
func (r *ReceivablesRepository) GetCustomerPaymentByID(ctx context.Context, organizationID, paymentID uuid.UUID) (*receivables.CustomerPayment, error) {
	payment := &receivables.CustomerPayment{}

	query := `
		SELECT id, organization_id, payment_number, customer_id, payment_date,
		       payment_method, reference_number, payment_amount, deposit_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_payments
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata pq.JSONBMap

	err := r.db.QueryRowContext(ctx, query, paymentID, organizationID).Scan(
		&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.CustomerID, &payment.PaymentDate,
		&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.DepositAccountID,
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

// GetCustomerPaymentByNumber retrieves a customer payment by payment number
func (r *ReceivablesRepository) GetCustomerPaymentByNumber(ctx context.Context, organizationID uuid.UUID, paymentNumber string) (*receivables.CustomerPayment, error) {
	payment := &receivables.CustomerPayment{}

	query := `
		SELECT id, organization_id, payment_number, customer_id, payment_date,
		       payment_method, reference_number, payment_amount, deposit_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_payments
		WHERE payment_number = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var metadata pq.JSONBMap

	err := r.db.QueryRowContext(ctx, query, paymentNumber, organizationID).Scan(
		&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.CustomerID, &payment.PaymentDate,
		&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.DepositAccountID,
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

// ListCustomerPayments retrieves a list of customer payments with filtering
func (r *ReceivablesRepository) ListCustomerPayments(ctx context.Context, organizationID uuid.UUID, opts *receivables.CustomerPaymentFilterOptions) ([]receivables.CustomerPayment, int64, error) {
	payments := make([]receivables.CustomerPayment, 0)

	whereConditions := []string{"organization_id = $1"}
	args := []interface{}{organizationID}
	argCount := 2

	if !opts.DeletedIncluded {
		whereConditions = append(whereConditions, "deleted_at IS NULL")
	}

	if opts.CustomerID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("customer_id = $%d", argCount))
		args = append(args, *opts.CustomerID)
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM customer_payments WHERE %s", whereClause)
	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, organization_id, payment_number, customer_id, payment_date,
		       payment_method, reference_number, payment_amount, deposit_account_id,
		       accounting_period_id, journal_entry_id, is_posted, memo, notes, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_payments
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
		payment := receivables.CustomerPayment{}
		var metadata pq.JSONBMap

		err := rows.Scan(
			&payment.ID, &payment.OrganizationID, &payment.PaymentNumber, &payment.CustomerID, &payment.PaymentDate,
			&payment.PaymentMethod, &payment.ReferenceNumber, &payment.PaymentAmount, &payment.DepositAccountID,
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

// UpdateCustomerPayment updates a customer payment
func (r *ReceivablesRepository) UpdateCustomerPayment(ctx context.Context, payment *receivables.CustomerPayment) error {
	metadata, err := json.Marshal(payment.Metadata)
	if err != nil {
		metadata = []byte("{}")
	}

	query := `
		UPDATE customer_payments
		SET payment_date = $1, payment_method = $2, reference_number = $3, payment_amount = $4,
		    deposit_account_id = $5, memo = $6, notes = $7, metadata = $8,
		    updated_at = $9, updated_by = $10
		WHERE id = $11 AND organization_id = $12
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.PaymentDate, payment.PaymentMethod, payment.ReferenceNumber, payment.PaymentAmount,
		payment.DepositAccountID, payment.Memo, payment.Notes, metadata,
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

// DeleteCustomerPayment soft deletes a customer payment
func (r *ReceivablesRepository) DeleteCustomerPayment(ctx context.Context, organizationID, paymentID uuid.UUID) error {
	query := `UPDATE customer_payments SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
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

// CreateCustomerPaymentApplication creates a new payment application
func (r *ReceivablesRepository) CreateCustomerPaymentApplication(ctx context.Context, app *receivables.CustomerPaymentApplication) error {
	query := `
		INSERT INTO customer_payment_applications (
			id, organization_id, customer_payment_id, customer_invoice_id, applied_amount, created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		app.ID, app.OrganizationID, app.CustomerPaymentID, app.CustomerInvoiceID, app.AppliedAmount, app.CreatedAt,
	)

	return err
}

// GetCustomerPaymentApplicationByID retrieves a payment application by ID
func (r *ReceivablesRepository) GetCustomerPaymentApplicationByID(ctx context.Context, organizationID, appID uuid.UUID) (*receivables.CustomerPaymentApplication, error) {
	app := &receivables.CustomerPaymentApplication{}

	query := `
		SELECT id, organization_id, customer_payment_id, customer_invoice_id, applied_amount, created_at, deleted_at
		FROM customer_payment_applications
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, appID, organizationID).Scan(
		&app.ID, &app.OrganizationID, &app.CustomerPaymentID, &app.CustomerInvoiceID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment application not found")
		}
		return nil, err
	}

	return app, nil
}

// ListCustomerPaymentApplications retrieves all applications for a payment
func (r *ReceivablesRepository) ListCustomerPaymentApplications(ctx context.Context, organizationID, paymentID uuid.UUID) ([]receivables.CustomerPaymentApplication, error) {
	apps := make([]receivables.CustomerPaymentApplication, 0)

	query := `
		SELECT id, organization_id, customer_payment_id, customer_invoice_id, applied_amount, created_at, deleted_at
		FROM customer_payment_applications
		WHERE customer_payment_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, paymentID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		app := receivables.CustomerPaymentApplication{}
		err := rows.Scan(
			&app.ID, &app.OrganizationID, &app.CustomerPaymentID, &app.CustomerInvoiceID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// ListPaymentApplicationsByInvoice retrieves all applications for an invoice
func (r *ReceivablesRepository) ListPaymentApplicationsByInvoice(ctx context.Context, organizationID, invoiceID uuid.UUID) ([]receivables.CustomerPaymentApplication, error) {
	apps := make([]receivables.CustomerPaymentApplication, 0)

	query := `
		SELECT id, organization_id, customer_payment_id, customer_invoice_id, applied_amount, created_at, deleted_at
		FROM customer_payment_applications
		WHERE customer_invoice_id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, invoiceID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		app := receivables.CustomerPaymentApplication{}
		err := rows.Scan(
			&app.ID, &app.OrganizationID, &app.CustomerPaymentID, &app.CustomerInvoiceID, &app.AppliedAmount, &app.CreatedAt, &app.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// DeleteCustomerPaymentApplication soft deletes a payment application
func (r *ReceivablesRepository) DeleteCustomerPaymentApplication(ctx context.Context, organizationID, appID uuid.UUID) error {
	query := `UPDATE customer_payment_applications SET deleted_at = NOW() WHERE id = $1 AND organization_id = $2`
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

// GetCustomerAgingReport generates an aging report for AR
func (r *ReceivablesRepository) GetCustomerAgingReport(ctx context.Context, organizationID uuid.UUID) ([]receivables.CustomerAgingReport, error) {
	reports := make([]receivables.CustomerAgingReport, 0)

	query := `
		WITH customer_info AS (
			SELECT DISTINCT ci.customer_id, c.name as customer_name
			FROM customer_invoices ci
			JOIN customers c ON ci.customer_id = c.id
			WHERE ci.organization_id = $1 AND ci.deleted_at IS NULL AND ci.status != 'paid' AND ci.status != 'void'
		),
		aging_buckets AS (
			SELECT
				ci.customer_id,
				CASE
					WHEN DATE(ci.due_date) >= CURRENT_DATE THEN ci.balance_due ELSE 0
				END as current,
				CASE
					WHEN DATE(ci.due_date) < CURRENT_DATE AND DATE(ci.due_date) >= CURRENT_DATE - INTERVAL '30 days' THEN ci.balance_due ELSE 0
				END as days_30,
				CASE
					WHEN DATE(ci.due_date) < CURRENT_DATE - INTERVAL '30 days' AND DATE(ci.due_date) >= CURRENT_DATE - INTERVAL '60 days' THEN ci.balance_due ELSE 0
				END as days_60,
				CASE
					WHEN DATE(ci.due_date) < CURRENT_DATE - INTERVAL '60 days' THEN ci.balance_due ELSE 0
				END as days_90_plus
			FROM customer_invoices ci
			WHERE ci.organization_id = $1 AND ci.deleted_at IS NULL
		)
		SELECT
			ci.customer_id,
			ci.customer_name,
			COALESCE(SUM(ab.current), 0) as current,
			COALESCE(SUM(ab.days_30), 0) as days_30,
			COALESCE(SUM(ab.days_60), 0) as days_60,
			COALESCE(SUM(ab.days_90_plus), 0) as days_90_plus
		FROM customer_info ci
		LEFT JOIN aging_buckets ab ON ci.customer_id = ab.customer_id
		GROUP BY ci.customer_id, ci.customer_name
		ORDER BY (COALESCE(SUM(ab.current), 0) + COALESCE(SUM(ab.days_30), 0) + COALESCE(SUM(ab.days_60), 0) + COALESCE(SUM(ab.days_90_plus), 0)) DESC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		report := receivables.CustomerAgingReport{}
		err := rows.Scan(
			&report.CustomerID, &report.CustomerName, &report.Current, &report.Days30, &report.Days60, &report.Days90Plus,
		)

		if err != nil {
			return nil, err
		}

		report.TotalDue = report.Current + report.Days30 + report.Days60 + report.Days90Plus
		reports = append(reports, report)
	}

	return reports, nil
}

// GetInvoicesByDueDate retrieves invoices due on or before a specific date
func (r *ReceivablesRepository) GetInvoicesByDueDate(ctx context.Context, organizationID uuid.UUID, dueDate time.Time) ([]receivables.CustomerInvoice, error) {
	invoices := make([]receivables.CustomerInvoice, 0)

	query := `
		SELECT id, organization_id, invoice_number, customer_id,
		       invoice_date, due_date, payment_terms, accounting_period_id,
		       subtotal, tax_amount, discount_amount, total_amount, paid_amount, balance_due,
		       status, journal_entry_id, is_posted, sale_id,
		       description, notes, memo, attachments, metadata,
		       created_at, updated_at, created_by, updated_by, deleted_at
		FROM customer_invoices
		WHERE organization_id = $1 AND due_date < $2 AND deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue')
		ORDER BY due_date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		invoice := receivables.CustomerInvoice{}
		var attachments pq.JSONBArray
		var metadata pq.JSONBMap

		err := rows.Scan(
			&invoice.ID, &invoice.OrganizationID, &invoice.InvoiceNumber, &invoice.CustomerID,
			&invoice.InvoiceDate, &invoice.DueDate, &invoice.PaymentTerms, &invoice.AccountingPeriodID,
			&invoice.Subtotal, &invoice.TaxAmount, &invoice.DiscountAmount, &invoice.TotalAmount, &invoice.PaidAmount, &invoice.BalanceDue,
			&invoice.Status, &invoice.JournalEntryID, &invoice.IsPosted, &invoice.SaleID,
			&invoice.Description, &invoice.Notes, &invoice.Memo, &attachments, &metadata,
			&invoice.CreatedAt, &invoice.UpdatedAt, &invoice.CreatedBy, &invoice.UpdatedBy, &invoice.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		invoice.Attachments = attachments
		invoice.Metadata = metadata
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// GetCustomerBalance retrieves the total balance owed by a customer
func (r *ReceivablesRepository) GetCustomerBalance(ctx context.Context, organizationID, customerID uuid.UUID) (float64, error) {
	var balance float64

	query := `
		SELECT COALESCE(SUM(balance_due), 0)
		FROM customer_invoices
		WHERE organization_id = $1 AND customer_id = $2 AND deleted_at IS NULL AND status IN ('unpaid', 'partial', 'overdue')
	`

	err := r.db.QueryRowContext(ctx, query, organizationID, customerID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
