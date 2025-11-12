package vendor_payment

import (
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

// Repository handles database operations for VendorPayments
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new VendorPayments repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// VendorPayments represents a vendor_payments entity
type VendorPayments struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	PaymentNumber string `json:"payment_number" db:"payment_number"`
	SupplierId uuid.UUID `json:"supplier_id" db:"supplier_id"`
	PaymentDate time.Time `json:"payment_date" db:"payment_date"`
	PaymentMethod string `json:"payment_method" db:"payment_method"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	PaymentAmount float64 `json:"payment_amount" db:"payment_amount"`
	BankAccountId *uuid.UUID `json:"bank_account_id" db:"bank_account_id"`
	AccountingPeriodId *uuid.UUID `json:"accounting_period_id" db:"accounting_period_id"`
	JournalEntryId *uuid.UUID `json:"journal_entry_id" db:"journal_entry_id"`
	IsPosted *bool `json:"is_posted" db:"is_posted"`
	Memo *string `json:"memo" db:"memo"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new vendor_payments record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *VendorPayments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "vendor_payments", duration, nil)
	}()

	query := `
		INSERT INTO vendor_payments (
			, organization_id
			, payment_number
			, supplier_id
			, payment_date
			, payment_method
			, reference_number
			, payment_amount
			, bank_account_id
			, accounting_period_id
			, journal_entry_id
			, is_posted
			, memo
			, notes
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
			, $18
			, $19
			, $20
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.PaymentNumber,
		entity.SupplierId,
		entity.PaymentDate,
		entity.PaymentMethod,
		entity.ReferenceNumber,
		entity.PaymentAmount,
		entity.BankAccountId,
		entity.AccountingPeriodId,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.Memo,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create vendor_payments", zap.Error(err))
		return fmt.Errorf("failed to create vendor_payments: %w", err)
	}

	r.logger.Info("created vendor_payments",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a vendor_payments by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*VendorPayments, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_payments", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, payment_number
			, supplier_id
			, payment_date
			, payment_method
			, reference_number
			, payment_amount
			, bank_account_id
			, accounting_period_id
			, journal_entry_id
			, is_posted
			, memo
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_payments
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity VendorPayments
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.PaymentNumber,
		&entity.SupplierId,
		&entity.PaymentDate,
		&entity.PaymentMethod,
		&entity.ReferenceNumber,
		&entity.PaymentAmount,
		&entity.BankAccountId,
		&entity.AccountingPeriodId,
		&entity.JournalEntryId,
		&entity.IsPosted,
		&entity.Memo,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vendor_payments not found")
	}

	if err != nil {
		r.logger.Error("failed to get vendor_payments", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get vendor_payments: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of vendor_payments records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*VendorPayments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_payments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_payments
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_payments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, payment_number
			, supplier_id
			, payment_date
			, payment_method
			, reference_number
			, payment_amount
			, bank_account_id
			, accounting_period_id
			, journal_entry_id
			, is_posted
			, memo
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_payments
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_payments", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_payments: %w", err)
	}
	defer rows.Close()

	var entities []*VendorPayments
	for rows.Next() {
		var entity VendorPayments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PaymentNumber,
			&entity.SupplierId,
			&entity.PaymentDate,
			&entity.PaymentMethod,
			&entity.ReferenceNumber,
			&entity.PaymentAmount,
			&entity.BankAccountId,
			&entity.AccountingPeriodId,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.Memo,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_payments: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating vendor_payments rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing vendor_payments record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *VendorPayments) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_payments", duration, nil)
	}()

	query := `
		UPDATE vendor_payments
		SET
			, organization_id = $2
			, payment_number = $3
			, supplier_id = $4
			, payment_date = $5
			, payment_method = $6
			, reference_number = $7
			, payment_amount = $8
			, bank_account_id = $9
			, accounting_period_id = $10
			, journal_entry_id = $11
			, is_posted = $12
			, memo = $13
			, notes = $14
			, metadata = $15
			, updated_at = $17
			, created_by = $18
			, updated_by = $19
			, deleted_at = $20
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.PaymentNumber,
		entity.SupplierId,
		entity.PaymentDate,
		entity.PaymentMethod,
		entity.ReferenceNumber,
		entity.PaymentAmount,
		entity.BankAccountId,
		entity.AccountingPeriodId,
		entity.JournalEntryId,
		entity.IsPosted,
		entity.Memo,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update vendor_payments", zap.Error(err))
		return fmt.Errorf("failed to update vendor_payments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_payments not found or already deleted")
	}

	r.logger.Info("updated vendor_payments",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a vendor_payments record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_payments", duration, nil)
	}()

	query := `
		UPDATE vendor_payments
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete vendor_payments", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete vendor_payments: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_payments not found or already deleted")
	}

	r.logger.Info("deleted vendor_payments", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves vendor_payments records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*VendorPayments, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_payments", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_payments
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_payments records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, payment_number
			, supplier_id
			, payment_date
			, payment_method
			, reference_number
			, payment_amount
			, bank_account_id
			, accounting_period_id
			, journal_entry_id
			, is_posted
			, memo
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
		FROM vendor_payments
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_payments by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_payments: %w", err)
	}
	defer rows.Close()

	var entities []*VendorPayments
	for rows.Next() {
		var entity VendorPayments
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.PaymentNumber,
			&entity.SupplierId,
			&entity.PaymentDate,
			&entity.PaymentMethod,
			&entity.ReferenceNumber,
			&entity.PaymentAmount,
			&entity.BankAccountId,
			&entity.AccountingPeriodId,
			&entity.JournalEntryId,
			&entity.IsPosted,
			&entity.Memo,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_payments: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

