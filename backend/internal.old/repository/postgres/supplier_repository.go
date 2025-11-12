package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/suppliers"
)

type SupplierRepository struct {
	db *DB
}

func NewSupplierRepository(db *DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) List(ctx context.Context, orgID uuid.UUID, filters suppliers.SupplierFilters) ([]suppliers.Supplier, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_code, name, contact_person, email, phone,
		       address, city, state, country, postal_code, tax_number, payment_terms,
		       credit_limit, outstanding_balance, total_purchases, total_orders, last_order_date,
		       status, notes, created_at, updated_at, created_by, updated_by
		FROM suppliers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR supplier_code ILIKE $%d OR email ILIKE $%d OR contact_person ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	query += " ORDER BY name ASC"

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

	var supplierList []suppliers.Supplier
	for rows.Next() {
		var s suppliers.Supplier
		err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.SupplierCode, &s.Name, &s.ContactPerson, &s.Email, &s.Phone,
			&s.Address, &s.City, &s.State, &s.Country, &s.PostalCode, &s.TaxNumber, &s.PaymentTerms,
			&s.CreditLimit, &s.OutstandingBalance, &s.TotalPurchases, &s.TotalOrders, &s.LastOrderDate,
			&s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		supplierList = append(supplierList, s)
	}

	return supplierList, rows.Err()
}

func (r *SupplierRepository) Count(ctx context.Context, orgID uuid.UUID, filters suppliers.SupplierFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM suppliers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR supplier_code ILIKE $%d OR email ILIKE $%d OR contact_person ILIKE $%d)", argCount, argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *SupplierRepository) Create(ctx context.Context, supplier *suppliers.Supplier) error {
	if err := r.db.SetOrganizationContext(ctx, supplier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO suppliers (
			id, organization_id, supplier_code, name, contact_person, email, phone,
			address, city, state, country, postal_code, tax_number, payment_terms,
			credit_limit, outstanding_balance, total_purchases, total_orders, last_order_date,
			status, notes, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		supplier.ID, supplier.OrganizationID, supplier.SupplierCode, supplier.Name,
		supplier.ContactPerson, supplier.Email, supplier.Phone, supplier.Address,
		supplier.City, supplier.State, supplier.Country, supplier.PostalCode,
		supplier.TaxNumber, supplier.PaymentTerms, supplier.CreditLimit,
		supplier.OutstandingBalance, supplier.TotalPurchases, supplier.TotalOrders,
		supplier.LastOrderDate, supplier.Status, supplier.Notes,
		supplier.CreatedAt, supplier.UpdatedAt, supplier.CreatedBy,
	)
	return err
}

func (r *SupplierRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*suppliers.Supplier, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_code, name, contact_person, email, phone,
		       address, city, state, country, postal_code, tax_number, payment_terms,
		       credit_limit, outstanding_balance, total_purchases, total_orders, last_order_date,
		       status, notes, created_at, updated_at, created_by, updated_by
		FROM suppliers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var s suppliers.Supplier
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&s.ID, &s.OrganizationID, &s.SupplierCode, &s.Name, &s.ContactPerson, &s.Email, &s.Phone,
		&s.Address, &s.City, &s.State, &s.Country, &s.PostalCode, &s.TaxNumber, &s.PaymentTerms,
		&s.CreditLimit, &s.OutstandingBalance, &s.TotalPurchases, &s.TotalOrders, &s.LastOrderDate,
		&s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) GetByEmail(ctx context.Context, orgID uuid.UUID, email string) (*suppliers.Supplier, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, supplier_code, name, contact_person, email, phone,
		       address, city, state, country, postal_code, tax_number, payment_terms,
		       credit_limit, outstanding_balance, total_purchases, total_orders, last_order_date,
		       status, notes, created_at, updated_at, created_by, updated_by
		FROM suppliers
		WHERE organization_id = $1 AND LOWER(email) = LOWER($2) AND deleted_at IS NULL
	`

	var s suppliers.Supplier
	err := r.db.Pool.QueryRow(ctx, query, orgID, email).Scan(
		&s.ID, &s.OrganizationID, &s.SupplierCode, &s.Name, &s.ContactPerson, &s.Email, &s.Phone,
		&s.Address, &s.City, &s.State, &s.Country, &s.PostalCode, &s.TaxNumber, &s.PaymentTerms,
		&s.CreditLimit, &s.OutstandingBalance, &s.TotalPurchases, &s.TotalOrders, &s.LastOrderDate,
		&s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SupplierRepository) Update(ctx context.Context, supplier *suppliers.Supplier) error {
	if err := r.db.SetOrganizationContext(ctx, supplier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE suppliers SET
			supplier_code = $3, name = $4, contact_person = $5, email = $6, phone = $7,
			address = $8, city = $9, state = $10, country = $11, postal_code = $12,
			tax_number = $13, payment_terms = $14, credit_limit = $15,
			status = $16, notes = $17, updated_at = $18, updated_by = $19
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		supplier.OrganizationID, supplier.ID,
		supplier.SupplierCode, supplier.Name, supplier.ContactPerson, supplier.Email, supplier.Phone,
		supplier.Address, supplier.City, supplier.State, supplier.Country, supplier.PostalCode,
		supplier.TaxNumber, supplier.PaymentTerms, supplier.CreditLimit,
		supplier.Status, supplier.Notes, supplier.UpdatedAt, supplier.UpdatedBy,
	)
	return err
}

func (r *SupplierRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE suppliers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}
