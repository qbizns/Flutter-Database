package staff_commission

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

// Repository handles database operations for StaffCommissions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new StaffCommissions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// StaffCommissions represents a staff_commissions entity
type StaffCommissions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	EmployeeId uuid.UUID `json:"employee_id" db:"employee_id"`
	CommissionDate time.Time `json:"commission_date" db:"commission_date"`
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd time.Time `json:"period_end" db:"period_end"`
	SourceType string `json:"source_type" db:"source_type"`
	SourceSaleId *uuid.UUID `json:"source_sale_id" db:"source_sale_id"`
	SourceOrderId *uuid.UUID `json:"source_order_id" db:"source_order_id"`
	CommissionType *string `json:"commission_type" db:"commission_type"`
	CommissionRate *float64 `json:"commission_rate" db:"commission_rate"`
	SalesAmount *float64 `json:"sales_amount" db:"sales_amount"`
	CommissionAmount float64 `json:"commission_amount" db:"commission_amount"`
	Status *string `json:"status" db:"status"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	PaymentDate *time.Time `json:"payment_date" db:"payment_date"`
	PaymentMethod *string `json:"payment_method" db:"payment_method"`
	PaidBy *uuid.UUID `json:"paid_by" db:"paid_by"`
	Notes *string `json:"notes" db:"notes"`
	CalculationNotes *string `json:"calculation_notes" db:"calculation_notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'sale', *string `json:"'sale'," db:"'sale',"`
	'percentage', *string `json:"'percentage'," db:"'percentage',"`
	'pending', *string `json:"'pending'," db:"'pending',"`
	SalesAmount *string `json:"sales_amount" db:"sales_amount"`
}

// Create inserts a new staff_commissions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *StaffCommissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "staff_commissions", duration, nil)
	}()

	query := `
		INSERT INTO staff_commissions (
			, organization_id
			, location_id
			, employee_id
			, commission_date
			, period_start
			, period_end
			, source_type
			, source_sale_id
			, source_order_id
			, commission_type
			, commission_rate
			, sales_amount
			, commission_amount
			, status
			, approved_by
			, approved_at
			, payment_date
			, payment_method
			, paid_by
			, notes
			, calculation_notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'percentage',
			, 'pending',
			, sales_amount
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
			, $29
			, $30
			, $31
			, $32
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.CommissionDate,
		entity.PeriodStart,
		entity.PeriodEnd,
		entity.SourceType,
		entity.SourceSaleId,
		entity.SourceOrderId,
		entity.CommissionType,
		entity.CommissionRate,
		entity.SalesAmount,
		entity.CommissionAmount,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.PaymentDate,
		entity.PaymentMethod,
		entity.PaidBy,
		entity.Notes,
		entity.CalculationNotes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'sale',,
		entity.'percentage',,
		entity.'pending',,
		entity.SalesAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create staff_commissions", zap.Error(err))
		return fmt.Errorf("failed to create staff_commissions: %w", err)
	}

	r.logger.Info("created staff_commissions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a staff_commissions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*StaffCommissions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "staff_commissions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, commission_date
			, period_start
			, period_end
			, source_type
			, source_sale_id
			, source_order_id
			, commission_type
			, commission_rate
			, sales_amount
			, commission_amount
			, status
			, approved_by
			, approved_at
			, payment_date
			, payment_method
			, paid_by
			, notes
			, calculation_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'percentage',
			, 'pending',
			, sales_amount
		FROM staff_commissions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity StaffCommissions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.EmployeeId,
		&entity.CommissionDate,
		&entity.PeriodStart,
		&entity.PeriodEnd,
		&entity.SourceType,
		&entity.SourceSaleId,
		&entity.SourceOrderId,
		&entity.CommissionType,
		&entity.CommissionRate,
		&entity.SalesAmount,
		&entity.CommissionAmount,
		&entity.Status,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.PaymentDate,
		&entity.PaymentMethod,
		&entity.PaidBy,
		&entity.Notes,
		&entity.CalculationNotes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'sale',,
		&entity.'percentage',,
		&entity.'pending',,
		&entity.SalesAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("staff_commissions not found")
	}

	if err != nil {
		r.logger.Error("failed to get staff_commissions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get staff_commissions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of staff_commissions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*StaffCommissions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "staff_commissions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM staff_commissions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count staff_commissions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, commission_date
			, period_start
			, period_end
			, source_type
			, source_sale_id
			, source_order_id
			, commission_type
			, commission_rate
			, sales_amount
			, commission_amount
			, status
			, approved_by
			, approved_at
			, payment_date
			, payment_method
			, paid_by
			, notes
			, calculation_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'percentage',
			, 'pending',
			, sales_amount
		FROM staff_commissions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list staff_commissions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list staff_commissions: %w", err)
	}
	defer rows.Close()

	var entities []*StaffCommissions
	for rows.Next() {
		var entity StaffCommissions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.CommissionDate,
			&entity.PeriodStart,
			&entity.PeriodEnd,
			&entity.SourceType,
			&entity.SourceSaleId,
			&entity.SourceOrderId,
			&entity.CommissionType,
			&entity.CommissionRate,
			&entity.SalesAmount,
			&entity.CommissionAmount,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.PaymentDate,
			&entity.PaymentMethod,
			&entity.PaidBy,
			&entity.Notes,
			&entity.CalculationNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'sale',,
			&entity.'percentage',,
			&entity.'pending',,
			&entity.SalesAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan staff_commissions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating staff_commissions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing staff_commissions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *StaffCommissions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "staff_commissions", duration, nil)
	}()

	query := `
		UPDATE staff_commissions
		SET
			, organization_id = $2
			, location_id = $3
			, employee_id = $4
			, commission_date = $5
			, period_start = $6
			, period_end = $7
			, source_type = $8
			, source_sale_id = $9
			, source_order_id = $10
			, commission_type = $11
			, commission_rate = $12
			, sales_amount = $13
			, commission_amount = $14
			, status = $15
			, approved_by = $16
			, approved_at = $17
			, payment_date = $18
			, payment_method = $19
			, paid_by = $20
			, notes = $21
			, calculation_notes = $22
			, metadata = $23
			, updated_at = $25
			, created_by = $26
			, updated_by = $27
			, deleted_at = $28
			, 'sale', = $29
			, 'percentage', = $30
			, 'pending', = $31
			, sales_amount = $32
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $33
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.EmployeeId,
		entity.CommissionDate,
		entity.PeriodStart,
		entity.PeriodEnd,
		entity.SourceType,
		entity.SourceSaleId,
		entity.SourceOrderId,
		entity.CommissionType,
		entity.CommissionRate,
		entity.SalesAmount,
		entity.CommissionAmount,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.PaymentDate,
		entity.PaymentMethod,
		entity.PaidBy,
		entity.Notes,
		entity.CalculationNotes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'sale',,
		entity.'percentage',,
		entity.'pending',,
		entity.SalesAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update staff_commissions", zap.Error(err))
		return fmt.Errorf("failed to update staff_commissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("staff_commissions not found or already deleted")
	}

	r.logger.Info("updated staff_commissions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a staff_commissions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "staff_commissions", duration, nil)
	}()

	query := `
		UPDATE staff_commissions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete staff_commissions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete staff_commissions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("staff_commissions not found or already deleted")
	}

	r.logger.Info("deleted staff_commissions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves staff_commissions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*StaffCommissions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "staff_commissions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM staff_commissions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count staff_commissions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, employee_id
			, commission_date
			, period_start
			, period_end
			, source_type
			, source_sale_id
			, source_order_id
			, commission_type
			, commission_rate
			, sales_amount
			, commission_amount
			, status
			, approved_by
			, approved_at
			, payment_date
			, payment_method
			, paid_by
			, notes
			, calculation_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'percentage',
			, 'pending',
			, sales_amount
		FROM staff_commissions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list staff_commissions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list staff_commissions: %w", err)
	}
	defer rows.Close()

	var entities []*StaffCommissions
	for rows.Next() {
		var entity StaffCommissions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.EmployeeId,
			&entity.CommissionDate,
			&entity.PeriodStart,
			&entity.PeriodEnd,
			&entity.SourceType,
			&entity.SourceSaleId,
			&entity.SourceOrderId,
			&entity.CommissionType,
			&entity.CommissionRate,
			&entity.SalesAmount,
			&entity.CommissionAmount,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.PaymentDate,
			&entity.PaymentMethod,
			&entity.PaidBy,
			&entity.Notes,
			&entity.CalculationNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'sale',,
			&entity.'percentage',,
			&entity.'pending',,
			&entity.SalesAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan staff_commissions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

