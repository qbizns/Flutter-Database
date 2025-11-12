package shift

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

// Repository handles database operations for Shifts
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Shifts repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Shifts represents a shifts entity
type Shifts struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	ShiftNumber string `json:"shift_number" db:"shift_number"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime *time.Time `json:"end_time" db:"end_time"`
	Status *string `json:"status" db:"status"`
	OpeningCash *float64 `json:"opening_cash" db:"opening_cash"`
	OpeningNotes *string `json:"opening_notes" db:"opening_notes"`
	ExpectedCash *float64 `json:"expected_cash" db:"expected_cash"`
	ActualCash *float64 `json:"actual_cash" db:"actual_cash"`
	CashDifference *float64 `json:"cash_difference" db:"cash_difference"`
	ClosingNotes *string `json:"closing_notes" db:"closing_notes"`
	TotalSales *float64 `json:"total_sales" db:"total_sales"`
	TotalTransactions *int64 `json:"total_transactions" db:"total_transactions"`
	TotalRefunds *float64 `json:"total_refunds" db:"total_refunds"`
	TotalDiscounts *float64 `json:"total_discounts" db:"total_discounts"`
	PaymentBreakdown json.RawMessage `json:"payment_breakdown" db:"payment_breakdown"`
	Notes *string `json:"notes" db:"notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	ClosedBy *uuid.UUID `json:"closed_by" db:"closed_by"`
	ClosedAt *time.Time `json:"closed_at" db:"closed_at"`
	OpeningCash *string `json:"opening_cash" db:"opening_cash"`
	(expectedCash *string `json:"(expected_cash" db:"(expected_cash"`
	(actualCash *string `json:"(actual_cash" db:"(actual_cash"`
}

// Create inserts a new shifts record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Shifts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "shifts", duration, nil)
	}()

	query := `
		INSERT INTO shifts (
			, organization_id
			, location_id
			, user_id
			, shift_number
			, start_time
			, end_time
			, status
			, opening_cash
			, opening_notes
			, expected_cash
			, actual_cash
			, cash_difference
			, closing_notes
			, total_sales
			, total_transactions
			, total_refunds
			, total_discounts
			, payment_breakdown
			, notes
			, metadata
			, closed_by
			, closed_at
			, opening_cash
			, (expected_cash
			, (actual_cash
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
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.UserId,
		entity.ShiftNumber,
		entity.StartTime,
		entity.EndTime,
		entity.Status,
		entity.OpeningCash,
		entity.OpeningNotes,
		entity.ExpectedCash,
		entity.ActualCash,
		entity.CashDifference,
		entity.ClosingNotes,
		entity.TotalSales,
		entity.TotalTransactions,
		entity.TotalRefunds,
		entity.TotalDiscounts,
		entity.PaymentBreakdown,
		entity.Notes,
		entity.Metadata,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.OpeningCash,
		entity.(expectedCash,
		entity.(actualCash,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create shifts", zap.Error(err))
		return fmt.Errorf("failed to create shifts: %w", err)
	}

	r.logger.Info("created shifts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationID.String()),
	)

	return nil
}

// GetByID retrieves a shifts by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Shifts, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "shifts", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, user_id
			, shift_number
			, start_time
			, end_time
			, status
			, opening_cash
			, opening_notes
			, expected_cash
			, actual_cash
			, cash_difference
			, closing_notes
			, total_sales
			, total_transactions
			, total_refunds
			, total_discounts
			, payment_breakdown
			, notes
			, metadata
			, created_at
			, updated_at
			, closed_by
			, closed_at
			, opening_cash
			, (expected_cash
			, (actual_cash
		FROM shifts
		WHERE id = $1
		
	`

	var entity Shifts
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.UserId,
		&entity.ShiftNumber,
		&entity.StartTime,
		&entity.EndTime,
		&entity.Status,
		&entity.OpeningCash,
		&entity.OpeningNotes,
		&entity.ExpectedCash,
		&entity.ActualCash,
		&entity.CashDifference,
		&entity.ClosingNotes,
		&entity.TotalSales,
		&entity.TotalTransactions,
		&entity.TotalRefunds,
		&entity.TotalDiscounts,
		&entity.PaymentBreakdown,
		&entity.Notes,
		&entity.Metadata,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.ClosedBy,
		&entity.ClosedAt,
		&entity.OpeningCash,
		&entity.(expectedCash,
		&entity.(actualCash,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("shifts not found")
	}

	if err != nil {
		r.logger.Error("failed to get shifts", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of shifts records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Shifts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "shifts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM shifts
		
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count shifts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, user_id
			, shift_number
			, start_time
			, end_time
			, status
			, opening_cash
			, opening_notes
			, expected_cash
			, actual_cash
			, cash_difference
			, closing_notes
			, total_sales
			, total_transactions
			, total_refunds
			, total_discounts
			, payment_breakdown
			, notes
			, metadata
			, created_at
			, updated_at
			, closed_by
			, closed_at
			, opening_cash
			, (expected_cash
			, (actual_cash
		FROM shifts
		
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list shifts", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list shifts: %w", err)
	}
	defer rows.Close()

	var entities []*Shifts
	for rows.Next() {
		var entity Shifts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ShiftNumber,
			&entity.StartTime,
			&entity.EndTime,
			&entity.Status,
			&entity.OpeningCash,
			&entity.OpeningNotes,
			&entity.ExpectedCash,
			&entity.ActualCash,
			&entity.CashDifference,
			&entity.ClosingNotes,
			&entity.TotalSales,
			&entity.TotalTransactions,
			&entity.TotalRefunds,
			&entity.TotalDiscounts,
			&entity.PaymentBreakdown,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.OpeningCash,
			&entity.(expectedCash,
			&entity.(actualCash,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan shifts: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating shifts rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing shifts record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Shifts) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "shifts", duration, nil)
	}()

	query := `
		UPDATE shifts
		SET
			, organization_id = $2
			, location_id = $3
			, user_id = $4
			, shift_number = $5
			, start_time = $6
			, end_time = $7
			, status = $8
			, opening_cash = $9
			, opening_notes = $10
			, expected_cash = $11
			, actual_cash = $12
			, cash_difference = $13
			, closing_notes = $14
			, total_sales = $15
			, total_transactions = $16
			, total_refunds = $17
			, total_discounts = $18
			, payment_breakdown = $19
			, notes = $20
			, metadata = $21
			, updated_at = $23
			, closed_by = $24
			, closed_at = $25
			, opening_cash = $26
			, (expected_cash = $27
			, (actual_cash = $28
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $29
		
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.UserId,
		entity.ShiftNumber,
		entity.StartTime,
		entity.EndTime,
		entity.Status,
		entity.OpeningCash,
		entity.OpeningNotes,
		entity.ExpectedCash,
		entity.ActualCash,
		entity.CashDifference,
		entity.ClosingNotes,
		entity.TotalSales,
		entity.TotalTransactions,
		entity.TotalRefunds,
		entity.TotalDiscounts,
		entity.PaymentBreakdown,
		entity.Notes,
		entity.Metadata,
		entity.UpdatedAt,
		entity.ClosedBy,
		entity.ClosedAt,
		entity.OpeningCash,
		entity.(expectedCash,
		entity.(actualCash,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update shifts", zap.Error(err))
		return fmt.Errorf("failed to update shifts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("shifts not found or already deleted")
	}

	r.logger.Info("updated shifts",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete permanently deletes a shifts record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("DELETE", "shifts", duration, nil)
	}()

	query := `DELETE FROM shifts WHERE id = $1`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete shifts", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete shifts: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("shifts not found")
	}

	r.logger.Info("deleted shifts", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves shifts records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Shifts, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "shifts", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM shifts
		WHERE organization_id = $1
		
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count shifts records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, user_id
			, shift_number
			, start_time
			, end_time
			, status
			, opening_cash
			, opening_notes
			, expected_cash
			, actual_cash
			, cash_difference
			, closing_notes
			, total_sales
			, total_transactions
			, total_refunds
			, total_discounts
			, payment_breakdown
			, notes
			, metadata
			, created_at
			, updated_at
			, closed_by
			, closed_at
			, opening_cash
			, (expected_cash
			, (actual_cash
		FROM shifts
		WHERE organization_id = $1
		
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list shifts by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list shifts: %w", err)
	}
	defer rows.Close()

	var entities []*Shifts
	for rows.Next() {
		var entity Shifts
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ShiftNumber,
			&entity.StartTime,
			&entity.EndTime,
			&entity.Status,
			&entity.OpeningCash,
			&entity.OpeningNotes,
			&entity.ExpectedCash,
			&entity.ActualCash,
			&entity.CashDifference,
			&entity.ClosingNotes,
			&entity.TotalSales,
			&entity.TotalTransactions,
			&entity.TotalRefunds,
			&entity.TotalDiscounts,
			&entity.PaymentBreakdown,
			&entity.Notes,
			&entity.Metadata,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.ClosedBy,
			&entity.ClosedAt,
			&entity.OpeningCash,
			&entity.(expectedCash,
			&entity.(actualCash,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan shifts: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

