package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/payments"
)

type PaymentRepository struct {
	db *DB
}

func NewPaymentRepository(db *DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) List(ctx context.Context, orgID uuid.UUID, filters payments.PaymentFilters) ([]payments.Payment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, payment_method, payment_status,
		       amount, card_last_four, card_type, transaction_id, reference_number,
		       account_number, account_name, payment_date, processed_at, notes,
		       metadata, created_at, updated_at, created_by
		FROM payments
		WHERE organization_id = $1
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.PaymentMethod != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_method = $%d", argCount)
		args = append(args, filters.PaymentMethod)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, filters.PaymentStatus)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	query += " ORDER BY payment_date DESC"

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

	var paymentList []payments.Payment
	for rows.Next() {
		var p payments.Payment
		var metadata interface{}
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.SaleID, &p.PaymentMethod, &p.PaymentStatus,
			&p.Amount, &p.CardLastFour, &p.CardType, &p.TransactionID, &p.ReferenceNumber,
			&p.AccountNumber, &p.AccountName, &p.PaymentDate, &p.ProcessedAt, &p.Notes,
			&metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.Metadata = metadata
		paymentList = append(paymentList, p)
	}

	return paymentList, rows.Err()
}

func (r *PaymentRepository) Count(ctx context.Context, orgID uuid.UUID, filters payments.PaymentFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM payments WHERE organization_id = $1"
	args := []interface{}{orgID}
	argCount := 1

	if filters.SaleID != nil {
		argCount++
		query += fmt.Sprintf(" AND sale_id = $%d", argCount)
		args = append(args, *filters.SaleID)
	}

	if filters.PaymentMethod != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_method = $%d", argCount)
		args = append(args, filters.PaymentMethod)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, filters.PaymentStatus)
	}

	if filters.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_date >= $%d", argCount)
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_date <= $%d", argCount)
		args = append(args, *filters.EndDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *PaymentRepository) Create(ctx context.Context, payment *payments.Payment) error {
	if err := r.db.SetOrganizationContext(ctx, payment.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO payments (
			id, organization_id, sale_id, payment_method, payment_status,
			amount, card_last_four, card_type, transaction_id, reference_number,
			account_number, account_name, payment_date, processed_at, notes,
			metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		payment.ID, payment.OrganizationID, payment.SaleID, payment.PaymentMethod, payment.PaymentStatus,
		payment.Amount, payment.CardLastFour, payment.CardType, payment.TransactionID, payment.ReferenceNumber,
		payment.AccountNumber, payment.AccountName, payment.PaymentDate, payment.ProcessedAt, payment.Notes,
		payment.Metadata, payment.CreatedAt, payment.UpdatedAt, payment.CreatedBy,
	)
	return err
}

func (r *PaymentRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*payments.Payment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, payment_method, payment_status,
		       amount, card_last_four, card_type, transaction_id, reference_number,
		       account_number, account_name, payment_date, processed_at, notes,
		       metadata, created_at, updated_at, created_by
		FROM payments
		WHERE organization_id = $1 AND id = $2
	`

	var p payments.Payment
	var metadata interface{}
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&p.ID, &p.OrganizationID, &p.SaleID, &p.PaymentMethod, &p.PaymentStatus,
		&p.Amount, &p.CardLastFour, &p.CardType, &p.TransactionID, &p.ReferenceNumber,
		&p.AccountNumber, &p.AccountName, &p.PaymentDate, &p.ProcessedAt, &p.Notes,
		&metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	p.Metadata = metadata
	return &p, nil
}

func (r *PaymentRepository) GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]payments.Payment, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, sale_id, payment_method, payment_status,
		       amount, card_last_four, card_type, transaction_id, reference_number,
		       account_number, account_name, payment_date, processed_at, notes,
		       metadata, created_at, updated_at, created_by
		FROM payments
		WHERE organization_id = $1 AND sale_id = $2
		ORDER BY payment_date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentList []payments.Payment
	for rows.Next() {
		var p payments.Payment
		var metadata interface{}
		err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.SaleID, &p.PaymentMethod, &p.PaymentStatus,
			&p.Amount, &p.CardLastFour, &p.CardType, &p.TransactionID, &p.ReferenceNumber,
			&p.AccountNumber, &p.AccountName, &p.PaymentDate, &p.ProcessedAt, &p.Notes,
			&metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		p.Metadata = metadata
		paymentList = append(paymentList, p)
	}

	return paymentList, rows.Err()
}

func (r *PaymentRepository) Update(ctx context.Context, payment *payments.Payment) error {
	if err := r.db.SetOrganizationContext(ctx, payment.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE payments SET
			payment_method = $3, payment_status = $4, amount = $5,
			card_last_four = $6, card_type = $7, transaction_id = $8,
			reference_number = $9, account_number = $10, account_name = $11,
			payment_date = $12, processed_at = $13, notes = $14,
			metadata = $15, updated_at = $16
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query,
		payment.OrganizationID, payment.ID,
		payment.PaymentMethod, payment.PaymentStatus, payment.Amount,
		payment.CardLastFour, payment.CardType, payment.TransactionID,
		payment.ReferenceNumber, payment.AccountNumber, payment.AccountName,
		payment.PaymentDate, payment.ProcessedAt, payment.Notes,
		payment.Metadata, payment.UpdatedAt,
	)
	return err
}

func (r *PaymentRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		DELETE FROM payments
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id)
	return err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, paymentID uuid.UUID, status payments.PaymentStatus) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	var processedAt *time.Time
	if status == payments.PaymentStatusCompleted {
		now := time.Now()
		processedAt = &now
	}

	query := `
		UPDATE payments SET
			payment_status = $3, processed_at = $4, updated_at = $5
		WHERE organization_id = $1 AND id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, paymentID, status, processedAt, time.Now())
	return err
}
