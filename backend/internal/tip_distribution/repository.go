package tip_distribution

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

// Repository handles database operations for TipDistributions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new TipDistributions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// TipDistributions represents a tip_distributions entity
type TipDistributions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	TipPoolId *uuid.UUID `json:"tip_pool_id" db:"tip_pool_id"`
	DistributionDate time.Time `json:"distribution_date" db:"distribution_date"`
	PeriodStart *time.Time `json:"period_start" db:"period_start"`
	PeriodEnd *time.Time `json:"period_end" db:"period_end"`
	ShiftId *uuid.UUID `json:"shift_id" db:"shift_id"`
	EmployeeId uuid.UUID `json:"employee_id" db:"employee_id"`
	SourceType string `json:"source_type" db:"source_type"`
	SourceSaleId *uuid.UUID `json:"source_sale_id" db:"source_sale_id"`
	SourceOrderId *uuid.UUID `json:"source_order_id" db:"source_order_id"`
	TipAmount float64 `json:"tip_amount" db:"tip_amount"`
	DistributionAmount float64 `json:"distribution_amount" db:"distribution_amount"`
	DistributionPercentage *float64 `json:"distribution_percentage" db:"distribution_percentage"`
	PaymentStatus *string `json:"payment_status" db:"payment_status"`
	PaymentMethod *string `json:"payment_method" db:"payment_method"`
	PaidAt *time.Time `json:"paid_at" db:"paid_at"`
	PaidBy *uuid.UUID `json:"paid_by" db:"paid_by"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'sale', *string `json:"'sale'," db:"'sale',"`
	'pending', *string `json:"'pending'," db:"'pending',"`
	TipAmount *string `json:"tip_amount" db:"tip_amount"`
}

// Create inserts a new tip_distributions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *TipDistributions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "tip_distributions", duration, nil)
	}()

	query := `
		INSERT INTO tip_distributions (
			, organization_id
			, location_id
			, tip_pool_id
			, distribution_date
			, period_start
			, period_end
			, shift_id
			, employee_id
			, source_type
			, source_sale_id
			, source_order_id
			, tip_amount
			, distribution_amount
			, distribution_percentage
			, payment_status
			, payment_method
			, paid_at
			, paid_by
			, notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'pending',
			, tip_amount
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TipPoolId,
		entity.DistributionDate,
		entity.PeriodStart,
		entity.PeriodEnd,
		entity.ShiftId,
		entity.EmployeeId,
		entity.SourceType,
		entity.SourceSaleId,
		entity.SourceOrderId,
		entity.TipAmount,
		entity.DistributionAmount,
		entity.DistributionPercentage,
		entity.PaymentStatus,
		entity.PaymentMethod,
		entity.PaidAt,
		entity.PaidBy,
		entity.Notes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'sale',,
		entity.'pending',,
		entity.TipAmount,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create tip_distributions", zap.Error(err))
		return fmt.Errorf("failed to create tip_distributions: %w", err)
	}

	r.logger.Info("created tip_distributions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a tip_distributions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*TipDistributions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_distributions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, tip_pool_id
			, distribution_date
			, period_start
			, period_end
			, shift_id
			, employee_id
			, source_type
			, source_sale_id
			, source_order_id
			, tip_amount
			, distribution_amount
			, distribution_percentage
			, payment_status
			, payment_method
			, paid_at
			, paid_by
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'pending',
			, tip_amount
		FROM tip_distributions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity TipDistributions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.TipPoolId,
		&entity.DistributionDate,
		&entity.PeriodStart,
		&entity.PeriodEnd,
		&entity.ShiftId,
		&entity.EmployeeId,
		&entity.SourceType,
		&entity.SourceSaleId,
		&entity.SourceOrderId,
		&entity.TipAmount,
		&entity.DistributionAmount,
		&entity.DistributionPercentage,
		&entity.PaymentStatus,
		&entity.PaymentMethod,
		&entity.PaidAt,
		&entity.PaidBy,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'sale',,
		&entity.'pending',,
		&entity.TipAmount,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("tip_distributions not found")
	}

	if err != nil {
		r.logger.Error("failed to get tip_distributions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get tip_distributions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of tip_distributions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*TipDistributions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_distributions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tip_distributions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tip_distributions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, tip_pool_id
			, distribution_date
			, period_start
			, period_end
			, shift_id
			, employee_id
			, source_type
			, source_sale_id
			, source_order_id
			, tip_amount
			, distribution_amount
			, distribution_percentage
			, payment_status
			, payment_method
			, paid_at
			, paid_by
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'pending',
			, tip_amount
		FROM tip_distributions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tip_distributions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tip_distributions: %w", err)
	}
	defer rows.Close()

	var entities []*TipDistributions
	for rows.Next() {
		var entity TipDistributions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TipPoolId,
			&entity.DistributionDate,
			&entity.PeriodStart,
			&entity.PeriodEnd,
			&entity.ShiftId,
			&entity.EmployeeId,
			&entity.SourceType,
			&entity.SourceSaleId,
			&entity.SourceOrderId,
			&entity.TipAmount,
			&entity.DistributionAmount,
			&entity.DistributionPercentage,
			&entity.PaymentStatus,
			&entity.PaymentMethod,
			&entity.PaidAt,
			&entity.PaidBy,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'sale',,
			&entity.'pending',,
			&entity.TipAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tip_distributions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating tip_distributions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing tip_distributions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *TipDistributions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tip_distributions", duration, nil)
	}()

	query := `
		UPDATE tip_distributions
		SET
			, organization_id = $2
			, location_id = $3
			, tip_pool_id = $4
			, distribution_date = $5
			, period_start = $6
			, period_end = $7
			, shift_id = $8
			, employee_id = $9
			, source_type = $10
			, source_sale_id = $11
			, source_order_id = $12
			, tip_amount = $13
			, distribution_amount = $14
			, distribution_percentage = $15
			, payment_status = $16
			, payment_method = $17
			, paid_at = $18
			, paid_by = $19
			, notes = $20
			, metadata = $21
			, updated_at = $23
			, created_by = $24
			, updated_by = $25
			, deleted_at = $26
			, 'sale', = $27
			, 'pending', = $28
			, tip_amount = $29
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $30
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.TipPoolId,
		entity.DistributionDate,
		entity.PeriodStart,
		entity.PeriodEnd,
		entity.ShiftId,
		entity.EmployeeId,
		entity.SourceType,
		entity.SourceSaleId,
		entity.SourceOrderId,
		entity.TipAmount,
		entity.DistributionAmount,
		entity.DistributionPercentage,
		entity.PaymentStatus,
		entity.PaymentMethod,
		entity.PaidAt,
		entity.PaidBy,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'sale',,
		entity.'pending',,
		entity.TipAmount,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update tip_distributions", zap.Error(err))
		return fmt.Errorf("failed to update tip_distributions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tip_distributions not found or already deleted")
	}

	r.logger.Info("updated tip_distributions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a tip_distributions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "tip_distributions", duration, nil)
	}()

	query := `
		UPDATE tip_distributions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete tip_distributions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete tip_distributions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("tip_distributions not found or already deleted")
	}

	r.logger.Info("deleted tip_distributions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves tip_distributions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*TipDistributions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "tip_distributions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM tip_distributions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tip_distributions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, tip_pool_id
			, distribution_date
			, period_start
			, period_end
			, shift_id
			, employee_id
			, source_type
			, source_sale_id
			, source_order_id
			, tip_amount
			, distribution_amount
			, distribution_percentage
			, payment_status
			, payment_method
			, paid_at
			, paid_by
			, notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'sale',
			, 'pending',
			, tip_amount
		FROM tip_distributions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list tip_distributions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list tip_distributions: %w", err)
	}
	defer rows.Close()

	var entities []*TipDistributions
	for rows.Next() {
		var entity TipDistributions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.TipPoolId,
			&entity.DistributionDate,
			&entity.PeriodStart,
			&entity.PeriodEnd,
			&entity.ShiftId,
			&entity.EmployeeId,
			&entity.SourceType,
			&entity.SourceSaleId,
			&entity.SourceOrderId,
			&entity.TipAmount,
			&entity.DistributionAmount,
			&entity.DistributionPercentage,
			&entity.PaymentStatus,
			&entity.PaymentMethod,
			&entity.PaidAt,
			&entity.PaidBy,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'sale',,
			&entity.'pending',,
			&entity.TipAmount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan tip_distributions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

