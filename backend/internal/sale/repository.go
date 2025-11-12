package sale

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

// Repository handles database operations for Sales
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Sales repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Sales represents a sales entity
type Sales struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SaleNumber string `json:"sale_number" db:"sale_number"`
	ReferenceNumber *string `json:"reference_number" db:"reference_number"`
	TransactionType string `json:"transaction_type" db:"transaction_type"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	CashierId *uuid.UUID `json:"cashier_id" db:"cashier_id"`
	Subtotal float64 `json:"subtotal" db:"subtotal"`
	TaxAmount float64 `json:"tax_amount" db:"tax_amount"`
	DiscountAmount float64 `json:"discount_amount" db:"discount_amount"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	PaidAmount float64 `json:"paid_amount" db:"paid_amount"`
	ChangeAmount float64 `json:"change_amount" db:"change_amount"`
	OutstandingAmount float64 `json:"outstanding_amount" db:"outstanding_amount"`
	PaymentStatus string `json:"payment_status" db:"payment_status"`
	DiscountType *string `json:"discount_type" db:"discount_type"`
	DiscountValue *float64 `json:"discount_value" db:"discount_value"`
	DiscountReason *string `json:"discount_reason" db:"discount_reason"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
	Notes *string `json:"notes" db:"notes"`
	InternalNotes *string `json:"internal_notes" db:"internal_notes"`
	CustomFields json.RawMessage `json:"custom_fields" db:"custom_fields"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
}

// Create inserts a new sales record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Sales) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sales", duration, nil)
	}()

	query := `
		INSERT INTO sales (
			, organization_id
			, sale_number
			, reference_number
			, transaction_type
			, customer_id
			, cashier_id
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, paid_amount
			, change_amount
			, outstanding_amount
			, payment_status
			, discount_type
			, discount_value
			, discount_reason
			, transaction_date
			, completed_at
			, notes
			, internal_notes
			, custom_fields
			, metadata
			, deleted_at
			, created_by
			, updated_by
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
			, $24
			, $27
			, $28
			, $29
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SaleNumber,
		entity.ReferenceNumber,
		entity.TransactionType,
		entity.CustomerId,
		entity.CashierId,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.TotalAmount,
		entity.PaidAmount,
		entity.ChangeAmount,
		entity.OutstandingAmount,
		entity.PaymentStatus,
		entity.DiscountType,
		entity.DiscountValue,
		entity.DiscountReason,
		entity.TransactionDate,
		entity.CompletedAt,
		entity.Notes,
		entity.InternalNotes,
		entity.CustomFields,
		entity.Metadata,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create sales", zap.Error(err))
		return fmt.Errorf("failed to create sales: %w", err)
	}

	r.logger.Info("created sales",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a sales by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Sales, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, sale_number
			, reference_number
			, transaction_type
			, customer_id
			, cashier_id
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, paid_amount
			, change_amount
			, outstanding_amount
			, payment_status
			, discount_type
			, discount_value
			, discount_reason
			, transaction_date
			, completed_at
			, notes
			, internal_notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM sales
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Sales
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SaleNumber,
		&entity.ReferenceNumber,
		&entity.TransactionType,
		&entity.CustomerId,
		&entity.CashierId,
		&entity.Subtotal,
		&entity.TaxAmount,
		&entity.DiscountAmount,
		&entity.TotalAmount,
		&entity.PaidAmount,
		&entity.ChangeAmount,
		&entity.OutstandingAmount,
		&entity.PaymentStatus,
		&entity.DiscountType,
		&entity.DiscountValue,
		&entity.DiscountReason,
		&entity.TransactionDate,
		&entity.CompletedAt,
		&entity.Notes,
		&entity.InternalNotes,
		&entity.CustomFields,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sales not found")
	}

	if err != nil {
		r.logger.Error("failed to get sales", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sales records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Sales, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sales
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sales records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_number
			, reference_number
			, transaction_type
			, customer_id
			, cashier_id
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, paid_amount
			, change_amount
			, outstanding_amount
			, payment_status
			, discount_type
			, discount_value
			, discount_reason
			, transaction_date
			, completed_at
			, notes
			, internal_notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM sales
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sales", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sales: %w", err)
	}
	defer rows.Close()

	var entities []*Sales
	for rows.Next() {
		var entity Sales
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleNumber,
			&entity.ReferenceNumber,
			&entity.TransactionType,
			&entity.CustomerId,
			&entity.CashierId,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.TotalAmount,
			&entity.PaidAmount,
			&entity.ChangeAmount,
			&entity.OutstandingAmount,
			&entity.PaymentStatus,
			&entity.DiscountType,
			&entity.DiscountValue,
			&entity.DiscountReason,
			&entity.TransactionDate,
			&entity.CompletedAt,
			&entity.Notes,
			&entity.InternalNotes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sales: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sales rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sales record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Sales) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sales", duration, nil)
	}()

	query := `
		UPDATE sales
		SET
			, organization_id = $2
			, sale_number = $3
			, reference_number = $4
			, transaction_type = $5
			, customer_id = $6
			, cashier_id = $7
			, subtotal = $8
			, tax_amount = $9
			, discount_amount = $10
			, total_amount = $11
			, paid_amount = $12
			, change_amount = $13
			, outstanding_amount = $14
			, payment_status = $15
			, discount_type = $16
			, discount_value = $17
			, discount_reason = $18
			, transaction_date = $19
			, completed_at = $20
			, notes = $21
			, internal_notes = $22
			, custom_fields = $23
			, metadata = $24
			, updated_at = $26
			, deleted_at = $27
			, created_by = $28
			, updated_by = $29
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $30
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SaleNumber,
		entity.ReferenceNumber,
		entity.TransactionType,
		entity.CustomerId,
		entity.CashierId,
		entity.Subtotal,
		entity.TaxAmount,
		entity.DiscountAmount,
		entity.TotalAmount,
		entity.PaidAmount,
		entity.ChangeAmount,
		entity.OutstandingAmount,
		entity.PaymentStatus,
		entity.DiscountType,
		entity.DiscountValue,
		entity.DiscountReason,
		entity.TransactionDate,
		entity.CompletedAt,
		entity.Notes,
		entity.InternalNotes,
		entity.CustomFields,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sales", zap.Error(err))
		return fmt.Errorf("failed to update sales: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sales not found or already deleted")
	}

	r.logger.Info("updated sales",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a sales record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sales", duration, nil)
	}()

	query := `
		UPDATE sales
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sales", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sales: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sales not found or already deleted")
	}

	r.logger.Info("deleted sales", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sales records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Sales, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sales", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sales
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sales records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, sale_number
			, reference_number
			, transaction_type
			, customer_id
			, cashier_id
			, subtotal
			, tax_amount
			, discount_amount
			, total_amount
			, paid_amount
			, change_amount
			, outstanding_amount
			, payment_status
			, discount_type
			, discount_value
			, discount_reason
			, transaction_date
			, completed_at
			, notes
			, internal_notes
			, custom_fields
			, metadata
			, created_at
			, updated_at
			, deleted_at
			, created_by
			, updated_by
		FROM sales
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sales by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sales: %w", err)
	}
	defer rows.Close()

	var entities []*Sales
	for rows.Next() {
		var entity Sales
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SaleNumber,
			&entity.ReferenceNumber,
			&entity.TransactionType,
			&entity.CustomerId,
			&entity.CashierId,
			&entity.Subtotal,
			&entity.TaxAmount,
			&entity.DiscountAmount,
			&entity.TotalAmount,
			&entity.PaidAmount,
			&entity.ChangeAmount,
			&entity.OutstandingAmount,
			&entity.PaymentStatus,
			&entity.DiscountType,
			&entity.DiscountValue,
			&entity.DiscountReason,
			&entity.TransactionDate,
			&entity.CompletedAt,
			&entity.Notes,
			&entity.InternalNotes,
			&entity.CustomFields,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sales: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

