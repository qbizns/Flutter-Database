package vendor_bill

import (
	"encoding/json"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/metrics"
	"go.uber.org/zap"
)

// Repository handles database operations for VendorBills
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new VendorBills repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// VendorBills represents a vendor_bills entity
type VendorBills struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	BillNumber string `json:"bill_number" db:"bill_number"`
	VendorBillNumber *string `json:"vendor_bill_number" db:"vendor_bill_number"`
	SupplierId uuid.UUID `json:"supplier_id" db:"supplier_id"`
	BillDate time.Time `json:"bill_date" db:"bill_date"`
	DueDate time.Time `json:"due_date" db:"due_date"`
	PaymentTerms *string `json:"payment_terms" db:"payment_terms"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	Subtotal *float64 `json:"subtotal" db:"subtotal"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	PaidAmount *float64 `json:"paid_amount" db:"paid_amount"`
	BalanceDue *float64 `json:"balance_due" db:"balance_due"`
	Status *string `json:"status" db:"status"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	IsPosted *bool `json:"is_posted" db:"is_posted"`
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Memo *string `json:"memo" db:"memo"`
	Attachments json.RawMessage `json:"attachments" db:"attachments"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new vendor_bills record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *VendorBills) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "vendor_bills", duration, nil)
	}()

	query := `
		INSERT INTO vendor_bills (
			, organization_id
			, bill_number
			, vendor_bill_number
			, supplier_id
			, bill_date
			, due_date
			, payment_terms
			, accounting_period_id
			, subtotal
			, tax_amount
			, total_amount
			, paid_amount
			, balance_due
			, status
			, journal_entry_id
			, is_posted
			, purchase_order_id
			, description
			, notes
			, memo
			, attachments
			, metadata
			, created_by
			, updated_by
			, deleted_at
		) VALUES (
			, $2
			, $3
			, $4
			, $5
			, $6
			, $7
			, $8
			, $9
			, $10
			, $11
			, $12
			, $13
			, $14
			, $15
			, $16
			, $17
			, $18
			, $19
			, $20
			, $21
			, $22
			, $23
			, $26
			, $27
			, $28
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.BillNumber,
		entity.VendorBillNumber,
		entity.SupplierId,
		entity.BillDate,
		entity.DueDate,
		entity.PaymentTerms,
		entity.AccountingPeriodId,
		entity.Subtotal,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.PaidAmount,
		entity.BalanceDue,
		entity.Status,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.PurchaseOrderId,
		entity.Description,
		entity.Notes,
		entity.Memo,
		entity.Attachments,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create vendor_bills", zap.Error(err))
		return fmt.Errorf("failed to create vendor_bills: %w", err)
	}

	r.logger.Info("created vendor_bills",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a vendor_bills by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*VendorBills, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bills", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, bill_number
			, vendor_bill_number
			, supplier_id
			, bill_date
			, due_date
			, payment_terms
			, accounting_period_id
			, subtotal
			, tax_amount
			, total_amount
			, paid_amount
			, balance_due
			, status
			, journal_entry_id
			, is_posted
			, purchase_order_id
			, description
			, notes
			, memo
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_bills
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity VendorBills
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.BillNumber,
		&entity.VendorBillNumber,
		&entity.SupplierId,
		&entity.BillDate,
		&entity.DueDate,
		&entity.PaymentTerms,
		&entity.AccountingPeriodId,
		&entity.Subtotal,
		&entity.TaxAmount,
		&entity.TotalAmount,
		&entity.PaidAmount,
		&entity.BalanceDue,
		&entity.Status,
		&entity.JournalEntryId,
		&entity.IsPosted,
		&entity.PurchaseOrderId,
		&entity.Description,
		&entity.Notes,
		&entity.Memo,
		&entity.Attachments,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vendor_bills not found")
	}

	if err != nil {
		r.logger.Error("failed to get vendor_bills", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get vendor_bills: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of vendor_bills records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*VendorBills, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bills", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_bills
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_bills records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bill_number
			, vendor_bill_number
			, supplier_id
			, bill_date
			, due_date
			, payment_terms
			, accounting_period_id
			, subtotal
			, tax_amount
			, total_amount
			, paid_amount
			, balance_due
			, status
			, journal_entry_id
			, is_posted
			, purchase_order_id
			, description
			, notes
			, memo
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_bills
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_bills", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_bills: %w", err)
	}
	defer rows.Close()

	var entities []*VendorBills
	for rows.Next() {
		var entity VendorBills
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BillNumber,
			&entity.VendorBillNumber,
			&entity.SupplierId,
			&entity.BillDate,
			&entity.DueDate,
			&entity.PaymentTerms,
			&entity.AccountingPeriodId,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.PaidAmount,
			&entity.BalanceDue,
			&entity.Status,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.PurchaseOrderId,
			&entity.Description,
			&entity.Notes,
			&entity.Memo,
			&entity.Attachments,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_bills: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating vendor_bills rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing vendor_bills record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *VendorBills) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_bills", duration, nil)
	}()

	query := `
		UPDATE vendor_bills
		SET
			, organization_id = $2
			, bill_number = $3
			, vendor_bill_number = $4
			, supplier_id = $5
			, bill_date = $6
			, due_date = $7
			, payment_terms = $8
			, accounting_period_id = $9
			, subtotal = $10
			, tax_amount = $11
			, total_amount = $12
			, paid_amount = $13
			, balance_due = $14
			, status = $15
			, journal_entry_id = $16
			, is_posted = $17
			, purchase_order_id = $18
			, description = $19
			, notes = $20
			, memo = $21
			, attachments = $22
			, metadata = $23
			, updated_at = $25
			, created_by = $26
			, updated_by = $27
			, deleted_at = $28
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $29
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.BillNumber,
		entity.VendorBillNumber,
		entity.SupplierId,
		entity.BillDate,
		entity.DueDate,
		entity.PaymentTerms,
		entity.AccountingPeriodId,
		entity.Subtotal,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.PaidAmount,
		entity.BalanceDue,
		entity.Status,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.PurchaseOrderId,
		entity.Description,
		entity.Notes,
		entity.Memo,
		entity.Attachments,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update vendor_bills", zap.Error(err))
		return fmt.Errorf("failed to update vendor_bills: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_bills not found or already deleted")
	}

	r.logger.Info("updated vendor_bills",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a vendor_bills record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_bills", duration, nil)
	}()

	query := `
		UPDATE vendor_bills
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete vendor_bills", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete vendor_bills: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_bills not found or already deleted")
	}

	r.logger.Info("deleted vendor_bills", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves vendor_bills records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*VendorBills, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bills", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_bills
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_bills records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, bill_number
			, vendor_bill_number
			, supplier_id
			, bill_date
			, due_date
			, payment_terms
			, accounting_period_id
			, subtotal
			, tax_amount
			, total_amount
			, paid_amount
			, balance_due
			, status
			, journal_entry_id
			, is_posted
			, purchase_order_id
			, description
			, notes
			, memo
			, attachments
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_bills
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_bills by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_bills: %w", err)
	}
	defer rows.Close()

	var entities []*VendorBills
	for rows.Next() {
		var entity VendorBills
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.BillNumber,
			&entity.VendorBillNumber,
			&entity.SupplierId,
			&entity.BillDate,
			&entity.DueDate,
			&entity.PaymentTerms,
			&entity.AccountingPeriodId,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.PaidAmount,
			&entity.BalanceDue,
			&entity.Status,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.PurchaseOrderId,
			&entity.Description,
			&entity.Notes,
			&entity.Memo,
			&entity.Attachments,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_bills: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

