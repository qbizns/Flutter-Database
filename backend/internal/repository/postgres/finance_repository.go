package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/your-org/pos-backend/internal/domain/finance"
)

// FinanceRepository implements finance.FinanceRepository
type FinanceRepository struct {
	db *sqlx.DB
}

// NewFinanceRepository creates a new finance repository
func NewFinanceRepository(db *sqlx.DB) finance.FinanceRepository {
	return &FinanceRepository{db: db}
}

// ========================
// CURRENCIES
// ========================

// CreateCurrency creates a new currency
func (r *FinanceRepository) CreateCurrency(ctx context.Context, currency *finance.Currency) error {
	query := `
		INSERT INTO currencies (
			id, currency_code, currency_name, currency_symbol,
			decimal_places, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		currency.ID, currency.CurrencyCode, currency.CurrencyName,
		currency.CurrencySymbol, currency.DecimalPlaces, currency.IsActive,
		currency.CreatedAt, currency.UpdatedAt,
	)

	return err
}

// GetCurrency retrieves a currency by ID
func (r *FinanceRepository) GetCurrency(ctx context.Context, id uuid.UUID) (*finance.Currency, error) {
	query := `
		SELECT id, currency_code, currency_name, currency_symbol,
		       decimal_places, is_active, created_at, updated_at, deleted_at
		FROM currencies
		WHERE id = $1 AND deleted_at IS NULL
	`

	c := &finance.Currency{}
	err := r.db.GetContext(ctx, c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrCurrencyNotFound
		}
		return nil, err
	}

	return c, nil
}

// GetCurrencyByCode retrieves a currency by code
func (r *FinanceRepository) GetCurrencyByCode(ctx context.Context, code string) (*finance.Currency, error) {
	query := `
		SELECT id, currency_code, currency_name, currency_symbol,
		       decimal_places, is_active, created_at, updated_at, deleted_at
		FROM currencies
		WHERE currency_code = $1 AND deleted_at IS NULL
	`

	c := &finance.Currency{}
	err := r.db.GetContext(ctx, c, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrCurrencyNotFound
		}
		return nil, err
	}

	return c, nil
}

// ListCurrencies lists currencies with filters
func (r *FinanceRepository) ListCurrencies(ctx context.Context, filter map[string]interface{}) ([]*finance.Currency, error) {
	query := `
		SELECT id, currency_code, currency_name, currency_symbol,
		       decimal_places, is_active, created_at, updated_at, deleted_at
		FROM currencies
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	query += " ORDER BY currency_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []*finance.Currency
	for rows.Next() {
		c := &finance.Currency{}
		if err := rows.StructScan(c); err != nil {
			return nil, err
		}
		currencies = append(currencies, c)
	}

	return currencies, rows.Err()
}

// UpdateCurrency updates a currency
func (r *FinanceRepository) UpdateCurrency(ctx context.Context, currency *finance.Currency) error {
	query := `
		UPDATE currencies
		SET currency_name = $1, currency_symbol = $2, decimal_places = $3,
		    is_active = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		currency.CurrencyName, currency.CurrencySymbol, currency.DecimalPlaces,
		currency.IsActive, currency.UpdatedAt, currency.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrCurrencyNotFound
	}

	return nil
}

// DeleteCurrency soft deletes a currency
func (r *FinanceRepository) DeleteCurrency(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE currencies
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrCurrencyNotFound
	}

	return nil
}

// ========================
// CURRENCY RATES
// ========================

// CreateCurrencyRate creates a new exchange rate
func (r *FinanceRepository) CreateCurrencyRate(ctx context.Context, rate *finance.CurrencyRate) error {
	query := `
		INSERT INTO currency_rates (
			id, organization_id, currency_code, rate_date, rate,
			source, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		rate.ID, rate.OrganizationID, rate.CurrencyCode, rate.RateDate,
		rate.Rate, rate.Source, rate.CreatedBy, rate.CreatedAt,
	)

	return err
}

// GetCurrencyRate retrieves a currency rate by ID
func (r *FinanceRepository) GetCurrencyRate(ctx context.Context, id uuid.UUID) (*finance.CurrencyRate, error) {
	query := `
		SELECT id, organization_id, currency_code, rate_date, rate,
		       source, created_by, created_at, deleted_at
		FROM currency_rates
		WHERE id = $1 AND deleted_at IS NULL
	`

	cr := &finance.CurrencyRate{}
	err := r.db.GetContext(ctx, cr, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrCurrencyRateNotFound
		}
		return nil, err
	}

	return cr, nil
}

// GetCurrencyRateByDate retrieves the exchange rate for a specific date
func (r *FinanceRepository) GetCurrencyRateByDate(ctx context.Context, organizationID uuid.UUID, currencyCode string, rateDate time.Time) (*finance.CurrencyRate, error) {
	query := `
		SELECT id, organization_id, currency_code, rate_date, rate,
		       source, created_by, created_at, deleted_at
		FROM currency_rates
		WHERE organization_id = $1 AND currency_code = $2 AND rate_date = $3 AND deleted_at IS NULL
		LIMIT 1
	`

	cr := &finance.CurrencyRate{}
	err := r.db.GetContext(ctx, cr, query, organizationID, currencyCode, rateDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrCurrencyRateNotFound
		}
		return nil, err
	}

	return cr, nil
}

// GetLatestCurrencyRate retrieves the most recent exchange rate
func (r *FinanceRepository) GetLatestCurrencyRate(ctx context.Context, organizationID uuid.UUID, currencyCode string) (*finance.CurrencyRate, error) {
	query := `
		SELECT id, organization_id, currency_code, rate_date, rate,
		       source, created_by, created_at, deleted_at
		FROM currency_rates
		WHERE organization_id = $1 AND currency_code = $2 AND deleted_at IS NULL
		ORDER BY rate_date DESC
		LIMIT 1
	`

	cr := &finance.CurrencyRate{}
	err := r.db.GetContext(ctx, cr, query, organizationID, currencyCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrCurrencyRateNotFound
		}
		return nil, err
	}

	return cr, nil
}

// ListCurrencyRates lists exchange rates with filters
func (r *FinanceRepository) ListCurrencyRates(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*finance.CurrencyRate, error) {
	query := `
		SELECT id, organization_id, currency_code, rate_date, rate,
		       source, created_by, created_at, deleted_at
		FROM currency_rates
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if currencyCode, ok := filter["currency_code"].(string); ok {
		query += fmt.Sprintf(" AND currency_code = $%d", argIndex)
		args = append(args, currencyCode)
		argIndex++
	}

	if startDate, ok := filter["start_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND rate_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate, ok := filter["end_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND rate_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	query += " ORDER BY rate_date DESC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rates []*finance.CurrencyRate
	for rows.Next() {
		cr := &finance.CurrencyRate{}
		if err := rows.StructScan(cr); err != nil {
			return nil, err
		}
		rates = append(rates, cr)
	}

	return rates, rows.Err()
}

// UpdateCurrencyRate updates an exchange rate
func (r *FinanceRepository) UpdateCurrencyRate(ctx context.Context, rate *finance.CurrencyRate) error {
	query := `
		UPDATE currency_rates
		SET rate = $1, source = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		rate.Rate, rate.Source, rate.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrCurrencyRateNotFound
	}

	return nil
}

// DeleteCurrencyRate soft deletes an exchange rate
func (r *FinanceRepository) DeleteCurrencyRate(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE currency_rates
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrCurrencyRateNotFound
	}

	return nil
}

// ========================
// PAYMENT TERMS
// ========================

// CreatePaymentTerm creates a new payment term
func (r *FinanceRepository) CreatePaymentTerm(ctx context.Context, term *finance.PaymentTerm) error {
	query := `
		INSERT INTO payment_terms (
			id, organization_id, term_code, term_name, note,
			is_active, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		term.ID, term.OrganizationID, term.TermCode, term.TermName, term.Note,
		term.IsActive, term.CreatedBy, term.UpdatedBy, term.CreatedAt, term.UpdatedAt,
	)

	return err
}

// GetPaymentTerm retrieves a payment term by ID
func (r *FinanceRepository) GetPaymentTerm(ctx context.Context, id, organizationID uuid.UUID) (*finance.PaymentTerm, error) {
	query := `
		SELECT id, organization_id, term_code, term_name, note,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM payment_terms
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	pt := &finance.PaymentTerm{}
	err := r.db.GetContext(ctx, pt, query, id, organizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrPaymentTermNotFound
		}
		return nil, err
	}

	return pt, nil
}

// GetPaymentTermByCode retrieves a payment term by code
func (r *FinanceRepository) GetPaymentTermByCode(ctx context.Context, organizationID uuid.UUID, code string) (*finance.PaymentTerm, error) {
	query := `
		SELECT id, organization_id, term_code, term_name, note,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM payment_terms
		WHERE organization_id = $1 AND term_code = $2 AND deleted_at IS NULL
	`

	pt := &finance.PaymentTerm{}
	err := r.db.GetContext(ctx, pt, query, organizationID, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrPaymentTermNotFound
		}
		return nil, err
	}

	return pt, nil
}

// ListPaymentTerms lists payment terms
func (r *FinanceRepository) ListPaymentTerms(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*finance.PaymentTerm, error) {
	query := `
		SELECT id, organization_id, term_code, term_name, note,
		       is_active, created_by, updated_by, created_at, updated_at, deleted_at
		FROM payment_terms
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if isActive, ok := filter["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, isActive)
		argIndex++
	}

	query += " ORDER BY term_code ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terms []*finance.PaymentTerm
	for rows.Next() {
		pt := &finance.PaymentTerm{}
		if err := rows.StructScan(pt); err != nil {
			return nil, err
		}
		terms = append(terms, pt)
	}

	return terms, rows.Err()
}

// UpdatePaymentTerm updates a payment term
func (r *FinanceRepository) UpdatePaymentTerm(ctx context.Context, term *finance.PaymentTerm) error {
	query := `
		UPDATE payment_terms
		SET term_name = $1, note = $2, is_active = $3,
		    updated_by = $4, updated_at = $5
		WHERE id = $6 AND organization_id = $7 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		term.TermName, term.Note, term.IsActive,
		term.UpdatedBy, term.UpdatedAt, term.ID, term.OrganizationID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrPaymentTermNotFound
	}

	return nil
}

// DeletePaymentTerm soft deletes a payment term
func (r *FinanceRepository) DeletePaymentTerm(ctx context.Context, id, organizationID uuid.UUID) error {
	query := `
		UPDATE payment_terms
		SET deleted_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrPaymentTermNotFound
	}

	return nil
}

// ========================
// PAYMENT TERM LINES
// ========================

// CreatePaymentTermLine creates a payment term line
func (r *FinanceRepository) CreatePaymentTermLine(ctx context.Context, line *finance.PaymentTermLine) error {
	query := `
		INSERT INTO payment_term_lines (
			id, payment_term_id, sequence, value_type, value_amount,
			days_after, end_of_month, day_of_month, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		line.ID, line.PaymentTermID, line.Sequence, line.ValueType, line.ValueAmount,
		line.DaysAfter, line.EndOfMonth, line.DayOfMonth, line.CreatedAt,
	)

	return err
}

// GetPaymentTermLine retrieves a payment term line
func (r *FinanceRepository) GetPaymentTermLine(ctx context.Context, id uuid.UUID) (*finance.PaymentTermLine, error) {
	query := `
		SELECT id, payment_term_id, sequence, value_type, value_amount,
		       days_after, end_of_month, day_of_month, created_at, deleted_at
		FROM payment_term_lines
		WHERE id = $1 AND deleted_at IS NULL
	`

	ptl := &finance.PaymentTermLine{}
	err := r.db.GetContext(ctx, ptl, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrPaymentTermLineNotFound
		}
		return nil, err
	}

	return ptl, nil
}

// ListPaymentTermLines lists lines for a payment term
func (r *FinanceRepository) ListPaymentTermLines(ctx context.Context, paymentTermID uuid.UUID) ([]*finance.PaymentTermLine, error) {
	query := `
		SELECT id, payment_term_id, sequence, value_type, value_amount,
		       days_after, end_of_month, day_of_month, created_at, deleted_at
		FROM payment_term_lines
		WHERE payment_term_id = $1 AND deleted_at IS NULL
		ORDER BY sequence ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, paymentTermID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []*finance.PaymentTermLine
	for rows.Next() {
		ptl := &finance.PaymentTermLine{}
		if err := rows.StructScan(ptl); err != nil {
			return nil, err
		}
		lines = append(lines, ptl)
	}

	return lines, rows.Err()
}

// DeletePaymentTermLine soft deletes a payment term line
func (r *FinanceRepository) DeletePaymentTermLine(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE payment_term_lines
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ========================
// INVOICE PAYMENT SCHEDULES
// ========================

// CreateInvoicePaymentSchedules creates payment schedules
func (r *FinanceRepository) CreateInvoicePaymentSchedules(ctx context.Context, schedules []*finance.InvoicePaymentSchedule) error {
	query := `
		INSERT INTO invoice_payment_schedules (
			id, organization_id, source_type, source_id, line_number,
			due_date, amount_due, amount_paid, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	for _, schedule := range schedules {
		_, err := r.db.ExecContext(ctx, query,
			schedule.ID, schedule.OrganizationID, schedule.SourceType, schedule.SourceID,
			schedule.LineNumber, schedule.DueDate, schedule.AmountDue, schedule.AmountPaid,
			schedule.Status, schedule.CreatedAt, schedule.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetInvoicePaymentSchedule retrieves a payment schedule
func (r *FinanceRepository) GetInvoicePaymentSchedule(ctx context.Context, id uuid.UUID) (*finance.InvoicePaymentSchedule, error) {
	query := `
		SELECT id, organization_id, source_type, source_id, line_number,
		       due_date, amount_due, amount_paid, status, created_at, updated_at, deleted_at
		FROM invoice_payment_schedules
		WHERE id = $1 AND deleted_at IS NULL
	`

	ips := &finance.InvoicePaymentSchedule{}
	err := r.db.GetContext(ctx, ips, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrInvoicePaymentScheduleNotFound
		}
		return nil, err
	}

	return ips, nil
}

// ListInvoicePaymentSchedules lists payment schedules with filters
func (r *FinanceRepository) ListInvoicePaymentSchedules(ctx context.Context, organizationID uuid.UUID, filter map[string]interface{}) ([]*finance.InvoicePaymentSchedule, error) {
	query := `
		SELECT id, organization_id, source_type, source_id, line_number,
		       due_date, amount_due, amount_paid, status, created_at, updated_at, deleted_at
		FROM invoice_payment_schedules
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{organizationID}
	argIndex := 2

	if status, ok := filter["status"].(string); ok {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if sourceType, ok := filter["source_type"].(string); ok {
		query += fmt.Sprintf(" AND source_type = $%d", argIndex)
		args = append(args, sourceType)
		argIndex++
	}

	if sourceID, ok := filter["source_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND source_id = $%d", argIndex)
		args = append(args, sourceID)
		argIndex++
	}

	if startDate, ok := filter["start_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND due_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate, ok := filter["end_date"].(time.Time); ok {
		query += fmt.Sprintf(" AND due_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	query += " ORDER BY due_date ASC"

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*finance.InvoicePaymentSchedule
	for rows.Next() {
		ips := &finance.InvoicePaymentSchedule{}
		if err := rows.StructScan(ips); err != nil {
			return nil, err
		}
		schedules = append(schedules, ips)
	}

	return schedules, rows.Err()
}

// ListInvoicePaymentSchedulesBySource lists schedules for a specific invoice
func (r *FinanceRepository) ListInvoicePaymentSchedulesBySource(ctx context.Context, organizationID uuid.UUID, sourceType finance.SourceType, sourceID uuid.UUID) ([]*finance.InvoicePaymentSchedule, error) {
	query := `
		SELECT id, organization_id, source_type, source_id, line_number,
		       due_date, amount_due, amount_paid, status, created_at, updated_at, deleted_at
		FROM invoice_payment_schedules
		WHERE organization_id = $1 AND source_type = $2 AND source_id = $3 AND deleted_at IS NULL
		ORDER BY line_number ASC
	`

	rows, err := r.db.QueryxContext(ctx, query, organizationID, sourceType, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*finance.InvoicePaymentSchedule
	for rows.Next() {
		ips := &finance.InvoicePaymentSchedule{}
		if err := rows.StructScan(ips); err != nil {
			return nil, err
		}
		schedules = append(schedules, ips)
	}

	return schedules, rows.Err()
}

// GetInvoicePaymentScheduleBySourceAndLine retrieves a specific schedule line
func (r *FinanceRepository) GetInvoicePaymentScheduleBySourceAndLine(ctx context.Context, organizationID uuid.UUID, sourceType finance.SourceType, sourceID uuid.UUID, lineNumber int) (*finance.InvoicePaymentSchedule, error) {
	query := `
		SELECT id, organization_id, source_type, source_id, line_number,
		       due_date, amount_due, amount_paid, status, created_at, updated_at, deleted_at
		FROM invoice_payment_schedules
		WHERE organization_id = $1 AND source_type = $2 AND source_id = $3 AND line_number = $4 AND deleted_at IS NULL
		LIMIT 1
	`

	ips := &finance.InvoicePaymentSchedule{}
	err := r.db.GetContext(ctx, ips, query, organizationID, sourceType, sourceID, lineNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finance.ErrInvoicePaymentScheduleNotFound
		}
		return nil, err
	}

	return ips, nil
}

// UpdateInvoicePaymentSchedule updates a payment schedule
func (r *FinanceRepository) UpdateInvoicePaymentSchedule(ctx context.Context, schedule *finance.InvoicePaymentSchedule) error {
	query := `
		UPDATE invoice_payment_schedules
		SET amount_paid = $1, status = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		schedule.AmountPaid, schedule.Status, schedule.UpdatedAt, schedule.ID,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrInvoicePaymentScheduleNotFound
	}

	return nil
}

// DeleteInvoicePaymentSchedule soft deletes a payment schedule
func (r *FinanceRepository) DeleteInvoicePaymentSchedule(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE invoice_payment_schedules
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return finance.ErrInvoicePaymentScheduleNotFound
	}

	return nil
}
