package pos_session

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

// Repository handles database operations for PosSessions
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new PosSessions repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// PosSessions represents a pos_sessions entity
type PosSessions struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	SessionNumber string `json:"session_number" db:"session_number"`
	SessionName *string `json:"session_name" db:"session_name"`
	DeviceId *uuid.UUID `json:"device_id" db:"device_id"`
	LocationId uuid.UUID `json:"location_id" db:"location_id"`
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	ShiftId *uuid.UUID `json:"shift_id" db:"shift_id"`
	OpenedAt time.Time `json:"opened_at" db:"opened_at"`
	ClosedAt *time.Time `json:"closed_at" db:"closed_at"`
	OpeningCash *float64 `json:"opening_cash" db:"opening_cash"`
	OpeningCard *float64 `json:"opening_card" db:"opening_card"`
	OpeningOther *float64 `json:"opening_other" db:"opening_other"`
	ExpectedCash *float64 `json:"expected_cash" db:"expected_cash"`
	ExpectedCard *float64 `json:"expected_card" db:"expected_card"`
	ExpectedOther *float64 `json:"expected_other" db:"expected_other"`
	CountedCash *float64 `json:"counted_cash" db:"counted_cash"`
	CountedCard *float64 `json:"counted_card" db:"counted_card"`
	CountedOther *float64 `json:"counted_other" db:"counted_other"`
	DifferenceCash *float64 `json:"difference_cash" db:"difference_cash"`
	DifferenceCard *float64 `json:"difference_card" db:"difference_card"`
	DifferenceOther *float64 `json:"difference_other" db:"difference_other"`
	Status *string `json:"status" db:"status"`
	ZReportNumber *string `json:"z_report_number" db:"z_report_number"`
	Notes *string `json:"notes" db:"notes"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Create inserts a new pos_sessions record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *PosSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "pos_sessions", duration, nil)
	}()

	query := `
		INSERT INTO pos_sessions (
			, organization_id
			, session_number
			, session_name
			, device_id
			, location_id
			, user_id
			, shift_id
			, opened_at
			, closed_at
			, opening_cash
			, opening_card
			, opening_other
			, expected_cash
			, expected_card
			, expected_other
			, counted_cash
			, counted_card
			, counted_other
			, difference_cash
			, difference_card
			, difference_other
			, status
			, z_report_number
			, notes
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
			, $24
			, $25
			, $26
			, $27
			, $30
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.SessionNumber,
		entity.SessionName,
		entity.DeviceId,
		entity.LocationId,
		entity.UserId,
		entity.ShiftId,
		entity.OpenedAt,
		entity.ClosedAt,
		entity.OpeningCash,
		entity.OpeningCard,
		entity.OpeningOther,
		entity.ExpectedCash,
		entity.ExpectedCard,
		entity.ExpectedOther,
		entity.CountedCash,
		entity.CountedCard,
		entity.CountedOther,
		entity.DifferenceCash,
		entity.DifferenceCard,
		entity.DifferenceOther,
		entity.Status,
		entity.ZReportNumber,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, &entity.UpdatedAt)
	

	if err != nil {
		r.logger.Error("failed to create pos_sessions", zap.Error(err))
		return fmt.Errorf("failed to create pos_sessions: %w", err)
	}

	r.logger.Info("created pos_sessions",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a pos_sessions by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*PosSessions, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_sessions", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, session_number
			, session_name
			, device_id
			, location_id
			, user_id
			, shift_id
			, opened_at
			, closed_at
			, opening_cash
			, opening_card
			, opening_other
			, expected_cash
			, expected_card
			, expected_other
			, counted_cash
			, counted_card
			, counted_other
			, difference_cash
			, difference_card
			, difference_other
			, status
			, z_report_number
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM pos_sessions
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity PosSessions
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.SessionNumber,
		&entity.SessionName,
		&entity.DeviceId,
		&entity.LocationId,
		&entity.UserId,
		&entity.ShiftId,
		&entity.OpenedAt,
		&entity.ClosedAt,
		&entity.OpeningCash,
		&entity.OpeningCard,
		&entity.OpeningOther,
		&entity.ExpectedCash,
		&entity.ExpectedCard,
		&entity.ExpectedOther,
		&entity.CountedCash,
		&entity.CountedCard,
		&entity.CountedOther,
		&entity.DifferenceCash,
		&entity.DifferenceCard,
		&entity.DifferenceOther,
		&entity.Status,
		&entity.ZReportNumber,
		&entity.Notes,
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("pos_sessions not found")
	}

	if err != nil {
		r.logger.Error("failed to get pos_sessions", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get pos_sessions: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of pos_sessions records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*PosSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_sessions
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, session_number
			, session_name
			, device_id
			, location_id
			, user_id
			, shift_id
			, opened_at
			, closed_at
			, opening_cash
			, opening_card
			, opening_other
			, expected_cash
			, expected_card
			, expected_other
			, counted_cash
			, counted_card
			, counted_other
			, difference_cash
			, difference_card
			, difference_other
			, status
			, z_report_number
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM pos_sessions
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_sessions", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*PosSessions
	for rows.Next() {
		var entity PosSessions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SessionNumber,
			&entity.SessionName,
			&entity.DeviceId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ShiftId,
			&entity.OpenedAt,
			&entity.ClosedAt,
			&entity.OpeningCash,
			&entity.OpeningCard,
			&entity.OpeningOther,
			&entity.ExpectedCash,
			&entity.ExpectedCard,
			&entity.ExpectedOther,
			&entity.CountedCash,
			&entity.CountedCard,
			&entity.CountedOther,
			&entity.DifferenceCash,
			&entity.DifferenceCard,
			&entity.DifferenceOther,
			&entity.Status,
			&entity.ZReportNumber,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating pos_sessions rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing pos_sessions record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *PosSessions) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_sessions", duration, nil)
	}()

	query := `
		UPDATE pos_sessions
		SET
			, organization_id = $2
			, session_number = $3
			, session_name = $4
			, device_id = $5
			, location_id = $6
			, user_id = $7
			, shift_id = $8
			, opened_at = $9
			, closed_at = $10
			, opening_cash = $11
			, opening_card = $12
			, opening_other = $13
			, expected_cash = $14
			, expected_card = $15
			, expected_other = $16
			, counted_cash = $17
			, counted_card = $18
			, counted_other = $19
			, difference_cash = $20
			, difference_card = $21
			, difference_other = $22
			, status = $23
			, z_report_number = $24
			, notes = $25
			, created_by = $26
			, updated_by = $27
			, updated_at = $29
			, deleted_at = $30
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $31
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.SessionNumber,
		entity.SessionName,
		entity.DeviceId,
		entity.LocationId,
		entity.UserId,
		entity.ShiftId,
		entity.OpenedAt,
		entity.ClosedAt,
		entity.OpeningCash,
		entity.OpeningCard,
		entity.OpeningOther,
		entity.ExpectedCash,
		entity.ExpectedCard,
		entity.ExpectedOther,
		entity.CountedCash,
		entity.CountedCard,
		entity.CountedOther,
		entity.DifferenceCash,
		entity.DifferenceCard,
		entity.DifferenceOther,
		entity.Status,
		entity.ZReportNumber,
		entity.Notes,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.UpdatedAt,
		entity.DeletedAt,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update pos_sessions", zap.Error(err))
		return fmt.Errorf("failed to update pos_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_sessions not found or already deleted")
	}

	r.logger.Info("updated pos_sessions",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a pos_sessions record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "pos_sessions", duration, nil)
	}()

	query := `
		UPDATE pos_sessions
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete pos_sessions", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete pos_sessions: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("pos_sessions not found or already deleted")
	}

	r.logger.Info("deleted pos_sessions", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves pos_sessions records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*PosSessions, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "pos_sessions", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pos_sessions
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pos_sessions records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, session_number
			, session_name
			, device_id
			, location_id
			, user_id
			, shift_id
			, opened_at
			, closed_at
			, opening_cash
			, opening_card
			, opening_other
			, expected_cash
			, expected_card
			, expected_other
			, counted_cash
			, counted_card
			, counted_other
			, difference_cash
			, difference_card
			, difference_other
			, status
			, z_report_number
			, notes
			, created_by
			, updated_by
			, created_at
			, updated_at
			, deleted_at
		FROM pos_sessions
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list pos_sessions by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list pos_sessions: %w", err)
	}
	defer rows.Close()

	var entities []*PosSessions
	for rows.Next() {
		var entity PosSessions
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.SessionNumber,
			&entity.SessionName,
			&entity.DeviceId,
			&entity.LocationId,
			&entity.UserId,
			&entity.ShiftId,
			&entity.OpenedAt,
			&entity.ClosedAt,
			&entity.OpeningCash,
			&entity.OpeningCard,
			&entity.OpeningOther,
			&entity.ExpectedCash,
			&entity.ExpectedCard,
			&entity.ExpectedOther,
			&entity.CountedCash,
			&entity.CountedCard,
			&entity.CountedOther,
			&entity.DifferenceCash,
			&entity.DifferenceCard,
			&entity.DifferenceOther,
			&entity.Status,
			&entity.ZReportNumber,
			&entity.Notes,
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.CreatedAt,
			&entity.UpdatedAt,
			&entity.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pos_sessions: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

