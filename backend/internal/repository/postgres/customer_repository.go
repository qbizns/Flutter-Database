package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/customers"
)

type CustomerRepository struct {
	db *DB
}

func NewCustomerRepository(db *DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) List(ctx context.Context, orgID uuid.UUID, filters customers.CustomerFilters) ([]customers.Customer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, first_name, last_name, email, phone,
		       date_of_birth, gender, address_line_1, address_line_2, city,
		       state, postal_code, country, customer_type, tax_id,
		       payment_term_days, credit_limit, current_balance, total_spent,
		       total_visits, last_visit_date, is_active, loyalty_member_number,
		       loyalty_points, loyalty_tier_id, accounting_customer_id,
		       preferred_payment_method, notes, created_at, updated_at,
		       created_by, updated_by
		FROM customers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CustomerType != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_type = $%d", argCount)
		args = append(args, *filters.CustomerType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	query += " ORDER BY first_name, last_name ASC"

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

	var customerList []customers.Customer
	for rows.Next() {
		var c customers.Customer
		err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
			&c.DateOfBirth, &c.Gender, &c.AddressLine1, &c.AddressLine2, &c.City,
			&c.State, &c.PostalCode, &c.Country, &c.CustomerType, &c.TaxID,
			&c.PaymentTermDays, &c.CreditLimit, &c.CurrentBalance, &c.TotalSpent,
			&c.TotalVisits, &c.LastVisitDate, &c.IsActive, &c.LoyaltyMemberNumber,
			&c.LoyaltyPoints, &c.LoyaltyTierID, &c.AccountingCustomerID,
			&c.PreferredPaymentMethod, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
			&c.CreatedBy, &c.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		customerList = append(customerList, c)
	}

	return customerList, rows.Err()
}

func (r *CustomerRepository) Count(ctx context.Context, orgID uuid.UUID, filters customers.CustomerFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM customers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR phone ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.CustomerType != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_type = $%d", argCount)
		args = append(args, *filters.CustomerType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *CustomerRepository) Create(ctx context.Context, customer *customers.Customer) error {
	if err := r.db.SetOrganizationContext(ctx, customer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO customers (
			id, organization_id, first_name, last_name, email, phone,
			date_of_birth, gender, address_line_1, address_line_2, city,
			state, postal_code, country, customer_type, tax_id,
			payment_term_days, credit_limit, current_balance, total_spent,
			total_visits, is_active, loyalty_member_number, loyalty_points,
			loyalty_tier_id, accounting_customer_id, preferred_payment_method,
			notes, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26,
			$27, $28, $29, $30, $31
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		customer.ID, customer.OrganizationID, customer.FirstName, customer.LastName,
		customer.Email, customer.Phone, customer.DateOfBirth, customer.Gender,
		customer.AddressLine1, customer.AddressLine2, customer.City, customer.State,
		customer.PostalCode, customer.Country, customer.CustomerType, customer.TaxID,
		customer.PaymentTermDays, customer.CreditLimit, customer.CurrentBalance,
		customer.TotalSpent, customer.TotalVisits, customer.IsActive,
		customer.LoyaltyMemberNumber, customer.LoyaltyPoints, customer.LoyaltyTierID,
		customer.AccountingCustomerID, customer.PreferredPaymentMethod, customer.Notes,
		customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy,
	)
	return err
}

func (r *CustomerRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*customers.Customer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, first_name, last_name, email, phone,
		       date_of_birth, gender, address_line_1, address_line_2, city,
		       state, postal_code, country, customer_type, tax_id,
		       payment_term_days, credit_limit, current_balance, total_spent,
		       total_visits, last_visit_date, is_active, loyalty_member_number,
		       loyalty_points, loyalty_tier_id, accounting_customer_id,
		       preferred_payment_method, notes, created_at, updated_at,
		       created_by, updated_by
		FROM customers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var c customers.Customer
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&c.ID, &c.OrganizationID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
		&c.DateOfBirth, &c.Gender, &c.AddressLine1, &c.AddressLine2, &c.City,
		&c.State, &c.PostalCode, &c.Country, &c.CustomerType, &c.TaxID,
		&c.PaymentTermDays, &c.CreditLimit, &c.CurrentBalance, &c.TotalSpent,
		&c.TotalVisits, &c.LastVisitDate, &c.IsActive, &c.LoyaltyMemberNumber,
		&c.LoyaltyPoints, &c.LoyaltyTierID, &c.AccountingCustomerID,
		&c.PreferredPaymentMethod, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
		&c.CreatedBy, &c.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*customers.Customer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, first_name, last_name, email, phone,
		       date_of_birth, gender, address_line_1, address_line_2, city,
		       state, postal_code, country, customer_type, tax_id,
		       payment_term_days, credit_limit, current_balance, total_spent,
		       total_visits, last_visit_date, is_active, loyalty_member_number,
		       loyalty_points, loyalty_tier_id, accounting_customer_id,
		       preferred_payment_method, notes, created_at, updated_at,
		       created_by, updated_by
		FROM customers
		WHERE organization_id = $1 AND LOWER(email) = LOWER($2) AND deleted_at IS NULL
	`

	var c customers.Customer
	err := r.db.Pool.QueryRow(ctx, query, orgID, email).Scan(
		&c.ID, &c.OrganizationID, &c.FirstName, &c.LastName, &c.Email, &c.Phone,
		&c.DateOfBirth, &c.Gender, &c.AddressLine1, &c.AddressLine2, &c.City,
		&c.State, &c.PostalCode, &c.Country, &c.CustomerType, &c.TaxID,
		&c.PaymentTermDays, &c.CreditLimit, &c.CurrentBalance, &c.TotalSpent,
		&c.TotalVisits, &c.LastVisitDate, &c.IsActive, &c.LoyaltyMemberNumber,
		&c.LoyaltyPoints, &c.LoyaltyTierID, &c.AccountingCustomerID,
		&c.PreferredPaymentMethod, &c.Notes, &c.CreatedAt, &c.UpdatedAt,
		&c.CreatedBy, &c.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepository) Update(ctx context.Context, customer *customers.Customer) error {
	if err := r.db.SetOrganizationContext(ctx, customer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customers SET
			first_name = $3, last_name = $4, email = $5, phone = $6,
			date_of_birth = $7, gender = $8, address_line_1 = $9,
			address_line_2 = $10, city = $11, state = $12, postal_code = $13,
			country = $14, customer_type = $15, tax_id = $16,
			payment_term_days = $17, credit_limit = $18, is_active = $19,
			loyalty_tier_id = $20, accounting_customer_id = $21,
			preferred_payment_method = $22, notes = $23,
			updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		customer.OrganizationID, customer.ID,
		customer.FirstName, customer.LastName, customer.Email, customer.Phone,
		customer.DateOfBirth, customer.Gender, customer.AddressLine1,
		customer.AddressLine2, customer.City, customer.State, customer.PostalCode,
		customer.Country, customer.CustomerType, customer.TaxID,
		customer.PaymentTermDays, customer.CreditLimit, customer.IsActive,
		customer.LoyaltyTierID, customer.AccountingCustomerID,
		customer.PreferredPaymentMethod, customer.Notes,
		customer.UpdatedAt, customer.UpdatedBy,
	)
	return err
}

func (r *CustomerRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *CustomerRepository) UpdateBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, amount float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customers
		SET current_balance = current_balance + $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, customerID, amount, time.Now())
	return err
}

func (r *CustomerRepository) UpdateLoyaltyPoints(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, points int) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE customers
		SET loyalty_points = loyalty_points + $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, customerID, points, time.Now())
	return err
}
