package sale_return

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

// Repository handles database operations for SaleReturns
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new SaleReturns repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// SaleReturns represents a sale_returns entity
type SaleReturns struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	ReturnNumber string `json:"return_number" db:"return_number"`
	OriginalSaleId *uuid.UUID `json:"original_sale_id" db:"original_sale_id"`
	CustomerId *uuid.UUID `json:"customer_id" db:"customer_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	ReturnDate time.Time `json:"return_date" db:"return_date"`
	TotalAmount *float64 `json:"total_amount" db:"total_amount"`
	RefundAmount *float64 `json:"refund_amount" db:"refund_amount"`
	RestockingFee *float64 `json:"restocking_fee" db:"restocking_fee"`
	RefundMethod *string `json:"refund_method" db:"refund_method"`
	Status *string `json:"status" db:"status"`
	// 	Status *string `json:"status" db:"status"`
	ApprovedBy *uuid.UUID `json:"approved_by" db:"approved_by"`
	ApprovedAt *time.Time `json:"approved_at" db:"approved_at"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new sale_returns record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *SaleReturns) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "sale_returns", duration, nil)
	}()

	query := `
		INSERT INTO sale_returns (
			, organization_id
			, return_number
			, original_sale_id
			, customer_id
			, location_id
			, user_id
			, return_date
			, total_amount
			, refund_amount
			, restocking_fee
			, refund_method
			, status
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
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
			, $21
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.ReturnNumber,
		entity.OriginalSaleId,
		entity.CustomerId,
		entity.LocationId,
		entity.UserId,
		entity.ReturnDate,
		entity.TotalAmount,
		entity.RefundAmount,
		entity.RestockingFee,
		entity.RefundMethod,
		entity.Status,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CreatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create sale_returns", zap.Error(err))
		return fmt.Errorf("failed to create sale_returns: %w", err)
	}

	r.logger.Info("created sale_returns",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a sale_returns by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*SaleReturns, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_returns", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, return_number
			, original_sale_id
			, customer_id
			, location_id
			, user_id
			, return_date
			, total_amount
			, refund_amount
			, restocking_fee
			, refund_method
			, status
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sale_returns
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity SaleReturns
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.ReturnNumber,
		&entity.OriginalSaleId,
		&entity.CustomerId,
		&entity.LocationId,
		&entity.UserId,
		&entity.ReturnDate,
		&entity.TotalAmount,
		&entity.RefundAmount,
		&entity.RestockingFee,
		&entity.RefundMethod,
		&entity.Status,
		&entity.Status,
		&entity.ApprovedBy,
		&entity.ApprovedAt,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("sale_returns not found")
	}

	if err != nil {
		r.logger.Error("failed to get sale_returns", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get sale_returns: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of sale_returns records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*SaleReturns, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_returns", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_returns
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_returns records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, return_number
			, original_sale_id
			, customer_id
			, location_id
			, user_id
			, return_date
			, total_amount
			, refund_amount
			, restocking_fee
			, refund_method
			, status
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sale_returns
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_returns", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_returns: %w", err)
	}
	defer rows.Close()

	var entities []*SaleReturns
	for rows.Next() {
		var entity SaleReturns
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReturnNumber,
			&entity.OriginalSaleId,
			&entity.CustomerId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ReturnDate,
			&entity.TotalAmount,
			&entity.RefundAmount,
			&entity.RestockingFee,
			&entity.RefundMethod,
			&entity.Status,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_returns: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating sale_returns rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing sale_returns record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *SaleReturns) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sale_returns", duration, nil)
	}()

	query := `
		UPDATE sale_returns
		SET
			, organization_id = $2
			, return_number = $3
			, original_sale_id = $4
			, customer_id = $5
			, location_id = $6
			, user_id = $7
			, return_date = $8
			, total_amount = $9
			, refund_amount = $10
			, restocking_fee = $11
			, refund_method = $12
			, status = $13
			, status = $14
			, approved_by = $15
			, approved_at = $16
			, notes = $17
			, created_by = $18
			, updated_at = $20
			, deleted_at = $21
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $22
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.ReturnNumber,
		entity.OriginalSaleId,
		entity.CustomerId,
		entity.LocationId,
		entity.UserId,
		entity.ReturnDate,
		entity.TotalAmount,
		entity.RefundAmount,
		entity.RestockingFee,
		entity.RefundMethod,
		entity.Status,
		entity.Status,
		entity.ApprovedBy,
		entity.ApprovedAt,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update sale_returns", zap.Error(err))
		return fmt.Errorf("failed to update sale_returns: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_returns not found or already deleted")
	}

	r.logger.Info("updated sale_returns",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a sale_returns record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "sale_returns", duration, nil)
	}()

	query := `
		UPDATE sale_returns
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete sale_returns", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete sale_returns: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("sale_returns not found or already deleted")
	}

	r.logger.Info("deleted sale_returns", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves sale_returns records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*SaleReturns, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "sale_returns", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM sale_returns
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sale_returns records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, return_number
			, original_sale_id
			, customer_id
			, location_id
			, user_id
			, return_date
			, total_amount
			, refund_amount
			, restocking_fee
			, refund_method
			, status
			, status
			, approved_by
			, approved_at
			, notes
			, created_by
			, created_at
			, updated_at
			, deleted_at
		FROM sale_returns
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list sale_returns by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list sale_returns: %w", err)
	}
	defer rows.Close()

	var entities []*SaleReturns
	for rows.Next() {
		var entity SaleReturns
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.ReturnNumber,
			&entity.OriginalSaleId,
			&entity.CustomerId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ReturnDate,
			&entity.TotalAmount,
			&entity.RefundAmount,
			&entity.RestockingFee,
			&entity.RefundMethod,
			&entity.Status,
			&entity.Status,
			&entity.ApprovedBy,
			&entity.ApprovedAt,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan sale_returns: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

