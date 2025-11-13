package expens

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

// Repository handles database operations for Expenses
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Expenses repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Expenses represents a expenses entity
type Expenses struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	ExpenseNumber string `json:"expense_number" db:"expense_number"`
	ExpenseDate time.Time `json:"expense_date" db:"expense_date"`
	Category string `json:"category" db:"category"`
	Subcategory *string `json:"subcategory" db:"subcategory"`
	PayeeName string `json:"payee_name" db:"payee_name"`
	PaymentMethod *string `json:"payment_method" db:"payment_method"`
	Amount float64 `json:"amount" db:"amount"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	Currency *string `json:"currency" db:"currency"`
	Status *string `json:"status" db:"status"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	PurchaseOrderId *uuid.UUID `json:"purchase_order_id" db:"purchase_order_id"`
	ReceiptUrl *string `json:"receipt_url" db:"receipt_url"`
	AttachmentUrls json.RawMessage `json:"attachment_urls" db:"attachment_urls"`
	Description *string `json:"description" db:"description"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Amount *string `json:"amount" db:"amount"`
	TaxAmount *string `json:"tax_amount" db:"tax_amount"`
	TotalAmount *string `json:"total_amount" db:"total_amount"`
}

// Create inserts a new expenses record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Expenses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "expenses", duration, nil)
	}()

	query := `
		INSERT INTO expenses (
			, organization_id
			, location_id
			, expense_number
			, expense_date
			, category
			, subcategory
			, payee_name
			, payment_method
			, amount
			, tax_amount
			, total_amount
			, currency
			, status
			, reference_number
			, purchase_order_id
			, receipt_url
			, attachment_urls
			, description
			, notes
			, metadata
			, deleted_at
			, created_by
			, updated_by
			, approved_by
			, approved_at
			, amount
			, tax_amount
			, total_amount
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
			, $24
			, $25
			, $26
			, $27
			, $28
			, $29
			, $30
			, $31
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.ExpenseNumber,
		entity.ExpenseDate,
		entity.Category,
		entity.Subcategory,
		entity.PayeeName,
		entity.PaymentMethod,
		entity.Amount,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.Currency,
		entity.Status,
		entity.ReferenceNumber,
		entity.PurchaseOrderId,
		entity.ReceiptUrl,
		entity.AttachmentUrls,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Amount,
		entity.TaxAmount,
		entity.TotalAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create expenses", zap.Error(err))
		return fmt.Errorf("failed to create expenses: %w", err)
	}

	r.logger.Info("created expenses",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a expenses by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Expenses, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "expenses", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, expense_number
			, expense_date
			, category
			, subcategory
			, payee_name
			, payment_method
			, amount
			, tax_amount
			, total_amount
			, currency
			, reference_number
			, purchase_order_id
			, receipt_url
			, attachment_urls
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, approved_by
			, approved_at
			, amount
			, tax_amount
			, total_amount
		FROM expenses
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Expenses
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.ExpenseNumber,
		&entity.ExpenseDate,
		&entity.Category,
		&entity.Subcategory,
		&entity.PayeeName,
		&entity.PaymentMethod,
		&entity.Amount,
		&entity.TaxAmount,
		&entity.TotalAmount,
		&entity.Currency,
		&entity.Status,
		&entity.ReferenceNumber,
		&entity.PurchaseOrderId,
		&entity.ReceiptUrl,
		&entity.AttachmentUrls,
		&entity.Description,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.Amount,
		&entity.TaxAmount,
		&entity.TotalAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("expenses not found")
	}

	if err != nil {
		r.logger.Error("failed to get expenses", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of expenses records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Expenses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "expenses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM expenses
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count expenses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, expense_number
			, expense_date
			, category
			, subcategory
			, payee_name
			, payment_method
			, amount
			, tax_amount
			, total_amount
			, currency
			, reference_number
			, purchase_order_id
			, receipt_url
			, attachment_urls
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, approved_by
			, approved_at
			, amount
			, tax_amount
			, total_amount
		FROM expenses
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list expenses", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list expenses: %w", err)
	}
	defer rows.Close()

	var entities []*Expenses
	for rows.Next() {
		var entity Expenses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.ExpenseNumber,
			&entity.ExpenseDate,
			&entity.Category,
			&entity.Subcategory,
			&entity.PayeeName,
			&entity.PaymentMethod,
			&entity.Amount,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.Currency,
			&entity.Status,
			&entity.ReferenceNumber,
			&entity.PurchaseOrderId,
			&entity.ReceiptUrl,
			&entity.AttachmentUrls,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Amount,
			&entity.TaxAmount,
			&entity.TotalAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan expenses: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating expenses rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing expenses record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Expenses) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "expenses", duration, nil)
	}()

	query := `
		UPDATE expenses
		SET
			, organization_id = $2
			, location_id = $3
			, expense_number = $4
			, expense_date = $5
			, category = $6
			, subcategory = $7
			, payee_name = $8
			, payment_method = $9
			, amount = $10
			, tax_amount = $11
			, total_amount = $12
			, currency = $13
			, status = $14
			, reference_number = $15
			, purchase_order_id = $16
			, receipt_url = $17
			, attachment_urls = $18
			, description = $19
			, notes = $20
			, metadata = $21
			, updated_at = $23
			, deleted_at = $24
			, created_by = $25
			, updated_by = $26
			, approved_by = $27
			, approved_at = $28
			, amount = $29
			, tax_amount = $30
			, total_amount = $31
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $32
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.ExpenseNumber,
		entity.ExpenseDate,
		entity.Category,
		entity.Subcategory,
		entity.PayeeName,
		entity.PaymentMethod,
		entity.Amount,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.Currency,
		entity.Status,
		entity.ReferenceNumber,
		entity.PurchaseOrderId,
		entity.ReceiptUrl,
		entity.AttachmentUrls,
		entity.Description,
		entity.Notes,
		entity.Metadata,
		time.Now(),
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Amount,
		entity.TaxAmount,
		entity.TotalAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update expenses", zap.Error(err))
		return fmt.Errorf("failed to update expenses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("expenses not found or already deleted")
	}

	r.logger.Info("updated expenses",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a expenses record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "expenses", duration, nil)
	}()

	query := `
		UPDATE expenses
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete expenses", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete expenses: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("expenses not found or already deleted")
	}

	r.logger.Info("deleted expenses", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves expenses records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Expenses, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "expenses", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM expenses
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count expenses records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, expense_number
			, expense_date
			, category
			, subcategory
			, payee_name
			, payment_method
			, amount
			, tax_amount
			, total_amount
			, currency
			, status
			, reference_number
			, purchase_order_id
			, receipt_url
			, attachment_urls
			, description
			, notes
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
			, approved_by
			, approved_at
			, amount
			, tax_amount
			, total_amount
		FROM expenses
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list expenses by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list expenses: %w", err)
	}
	defer rows.Close()

	var entities []*Expenses
	for rows.Next() {
		var entity Expenses
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.ExpenseNumber,
			&entity.ExpenseDate,
			&entity.Category,
			&entity.Subcategory,
			&entity.PayeeName,
			&entity.PaymentMethod,
			&entity.Amount,
			&entity.TaxAmount,
			&entity.TotalAmount,
			&entity.Currency,
			&entity.Status,
			&entity.ReferenceNumber,
			&entity.PurchaseOrderId,
			&entity.ReceiptUrl,
			&entity.AttachmentUrls,
			&entity.Description,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Amount,
			&entity.TaxAmount,
			&entity.TotalAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan expenses: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

