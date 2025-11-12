package vendor_bill_line

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

// Repository handles database operations for VendorBillLines
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new VendorBillLines repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// VendorBillLines represents a vendor_bill_lines entity
type VendorBillLines struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	VendorBillId uuid.UUID `json:"vendor_bill_id" db:"vendor_bill_id"`
	LineNumber int64 `json:"line_number" db:"line_number"`
	ExpenseAccountId uuid.UUID `json:"expense_account_id" db:"expense_account_id"`
	Description string `json:"description" db:"description"`
	Quantity *float64 `json:"quantity" db:"quantity"`
	UnitPrice float64 `json:"unit_price" db:"unit_price"`
	Amount float64 `json:"amount" db:"amount"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	Department *string `json:"department" db:"department"`
	ProjectCode *string `json:"project_code" db:"project_code"`
	TaxCode *string `json:"tax_code" db:"tax_code"`
	TaxAmount *float64 `json:"tax_amount" db:"tax_amount"`
	ProductId *uuid.UUID `json:"product_id" db:"product_id"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new vendor_bill_lines record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *VendorBillLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "vendor_bill_lines", duration, nil)
	}()

	query := `
		INSERT INTO vendor_bill_lines (
			, organization_id
			, vendor_bill_id
			, line_number
			, expense_account_id
			, description
			, quantity
			, unit_price
			, amount
			, location_id
			, department
			, project_code
			, tax_code
			, tax_amount
			, product_id
			, metadata
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
			, $19
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.VendorBillId,
		entity.LineNumber,
		entity.ExpenseAccountId,
		entity.Description,
		entity.Quantity,
		entity.UnitPrice,
		entity.Amount,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.TaxCode,
		entity.TaxAmount,
		entity.ProductId,
		entity.Metadata,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create vendor_bill_lines", zap.Error(err))
		return fmt.Errorf("failed to create vendor_bill_lines: %w", err)
	}

	r.logger.Info("created vendor_bill_lines",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a vendor_bill_lines by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*VendorBillLines, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bill_lines", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, line_number
			, expense_account_id
			, description
			, quantity
			, unit_price
			, amount
			, location_id
			, department
			, project_code
			, tax_code
			, tax_amount
			, product_id
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM vendor_bill_lines
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity VendorBillLines
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.VendorBillId,
		&entity.LineNumber,
		&entity.ExpenseAccountId,
		&entity.Description,
		&entity.Quantity,
		&entity.UnitPrice,
		&entity.Amount,
		&entity.LocationId,
		&entity.Department,
		&entity.ProjectCode,
		&entity.TaxCode,
		&entity.TaxAmount,
		&entity.ProductId,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("vendor_bill_lines not found")
	}

	if err != nil {
		r.logger.Error("failed to get vendor_bill_lines", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get vendor_bill_lines: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of vendor_bill_lines records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*VendorBillLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bill_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_bill_lines
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_bill_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, line_number
			, expense_account_id
			, description
			, quantity
			, unit_price
			, amount
			, location_id
			, department
			, project_code
			, tax_code
			, tax_amount
			, product_id
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM vendor_bill_lines
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_bill_lines", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_bill_lines: %w", err)
	}
	defer rows.Close()

	var entities []*VendorBillLines
	for rows.Next() {
		var entity VendorBillLines
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.VendorBillId,
			&entity.LineNumber,
			&entity.ExpenseAccountId,
			&entity.Description,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.Amount,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.TaxCode,
			&entity.TaxAmount,
			&entity.ProductId,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_bill_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating vendor_bill_lines rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing vendor_bill_lines record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *VendorBillLines) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_bill_lines", duration, nil)
	}()

	query := `
		UPDATE vendor_bill_lines
		SET
			, organization_id = $2
			, vendor_bill_id = $3
			, line_number = $4
			, expense_account_id = $5
			, description = $6
			, quantity = $7
			, unit_price = $8
			, amount = $9
			, location_id = $10
			, department = $11
			, project_code = $12
			, tax_code = $13
			, tax_amount = $14
			, product_id = $15
			, metadata = $16
			, updated_at = $18
			, deleted_at = $19
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $20
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.VendorBillId,
		entity.LineNumber,
		entity.ExpenseAccountId,
		entity.Description,
		entity.Quantity,
		entity.UnitPrice,
		entity.Amount,
		entity.LocationId,
		entity.Department,
		entity.ProjectCode,
		entity.TaxCode,
		entity.TaxAmount,
		entity.ProductId,
		entity.Metadata,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update vendor_bill_lines", zap.Error(err))
		return fmt.Errorf("failed to update vendor_bill_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_bill_lines not found or already deleted")
	}

	r.logger.Info("updated vendor_bill_lines",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a vendor_bill_lines record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "vendor_bill_lines", duration, nil)
	}()

	query := `
		UPDATE vendor_bill_lines
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete vendor_bill_lines", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete vendor_bill_lines: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("vendor_bill_lines not found or already deleted")
	}

	r.logger.Info("deleted vendor_bill_lines", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves vendor_bill_lines records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*VendorBillLines, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "vendor_bill_lines", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM vendor_bill_lines
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count vendor_bill_lines records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, vendor_bill_id
			, line_number
			, expense_account_id
			, description
			, quantity
			, unit_price
			, amount
			, location_id
			, department
			, project_code
			, tax_code
			, tax_amount
			, product_id
			, metadata
			, created_at
			, updated_at
			, deleted_at
		FROM vendor_bill_lines
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list vendor_bill_lines by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vendor_bill_lines: %w", err)
	}
	defer rows.Close()

	var entities []*VendorBillLines
	for rows.Next() {
		var entity VendorBillLines
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.VendorBillId,
			&entity.LineNumber,
			&entity.ExpenseAccountId,
			&entity.Description,
			&entity.Quantity,
			&entity.UnitPrice,
			&entity.Amount,
			&entity.LocationId,
			&entity.Department,
			&entity.ProjectCode,
			&entity.TaxCode,
			&entity.TaxAmount,
			&entity.ProductId,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan vendor_bill_lines: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

