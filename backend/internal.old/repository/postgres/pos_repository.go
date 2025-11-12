package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/pos"
)

// POSSessionRepository implements pos.POSSessionRepository
type POSSessionRepository struct {
	db *DB
}

// NewPOSSessionRepository creates a new POS session repository
func NewPOSSessionRepository(db *DB) *POSSessionRepository {
	return &POSSessionRepository{db: db}
}

// Create creates a new POS session
func (r *POSSessionRepository) Create(ctx context.Context, session *pos.POSSession) error {
	if err := r.db.SetOrganizationContext(ctx, session.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO pos_sessions (
			id, organization_id, session_number, session_name, device_id,
			location_id, user_id, shift_id, opened_at, opening_cash,
			opening_card, opening_other, expected_cash, expected_card,
			expected_other, difference_cash, difference_card, difference_other,
			status, notes, created_by, updated_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		session.ID, session.OrganizationID, session.SessionNumber, session.SessionName,
		session.DeviceID, session.LocationID, session.UserID, session.ShiftID, session.OpenedAt,
		session.OpeningCash, session.OpeningCard, session.OpeningOther,
		session.ExpectedCash, session.ExpectedCard, session.ExpectedOther,
		session.DifferenceCash, session.DifferenceCard, session.DifferenceOther,
		session.Status, session.Notes, session.CreatedBy, session.UpdatedBy,
		session.CreatedAt, session.UpdatedAt,
	)
	return err
}

// Get retrieves a POS session by ID
func (r *POSSessionRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*pos.POSSession, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, session_number, session_name, device_id, location_id, user_id,
		       shift_id, opened_at, closed_at, opening_cash, opening_card, opening_other,
		       expected_cash, expected_card, expected_other, counted_cash, counted_card,
		       counted_other, difference_cash, difference_card, difference_other, status,
		       z_report_number, notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM pos_sessions
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var session pos.POSSession
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&session.ID, &session.OrganizationID, &session.SessionNumber, &session.SessionName,
		&session.DeviceID, &session.LocationID, &session.UserID, &session.ShiftID, &session.OpenedAt,
		&session.ClosedAt, &session.OpeningCash, &session.OpeningCard, &session.OpeningOther,
		&session.ExpectedCash, &session.ExpectedCard, &session.ExpectedOther,
		&session.CountedCash, &session.CountedCard, &session.CountedOther,
		&session.DifferenceCash, &session.DifferenceCard, &session.DifferenceOther,
		&session.Status, &session.ZReportNumber, &session.Notes,
		&session.CreatedBy, &session.UpdatedBy, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetBySessionNumber retrieves a session by session number
func (r *POSSessionRepository) GetBySessionNumber(ctx context.Context, orgID uuid.UUID, sessionNumber string) (*pos.POSSession, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, session_number, session_name, device_id, location_id, user_id,
		       shift_id, opened_at, closed_at, opening_cash, opening_card, opening_other,
		       expected_cash, expected_card, expected_other, counted_cash, counted_card,
		       counted_other, difference_cash, difference_card, difference_other, status,
		       z_report_number, notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM pos_sessions
		WHERE organization_id = $1 AND session_number = $2 AND deleted_at IS NULL
	`

	var session pos.POSSession
	err := r.db.Pool.QueryRow(ctx, query, orgID, sessionNumber).Scan(
		&session.ID, &session.OrganizationID, &session.SessionNumber, &session.SessionName,
		&session.DeviceID, &session.LocationID, &session.UserID, &session.ShiftID, &session.OpenedAt,
		&session.ClosedAt, &session.OpeningCash, &session.OpeningCard, &session.OpeningOther,
		&session.ExpectedCash, &session.ExpectedCard, &session.ExpectedOther,
		&session.CountedCash, &session.CountedCard, &session.CountedOther,
		&session.DifferenceCash, &session.DifferenceCard, &session.DifferenceOther,
		&session.Status, &session.ZReportNumber, &session.Notes,
		&session.CreatedBy, &session.UpdatedBy, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// List retrieves POS sessions with filters
func (r *POSSessionRepository) List(ctx context.Context, orgID uuid.UUID, filters pos.POSSessionFilters) ([]pos.POSSession, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, session_number, session_name, device_id, location_id, user_id,
		       shift_id, opened_at, closed_at, opening_cash, opening_card, opening_other,
		       expected_cash, expected_card, expected_other, counted_cash, counted_card,
		       counted_other, difference_cash, difference_card, difference_other, status,
		       z_report_number, notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM pos_sessions
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
	}

	if filters.DeviceID != nil {
		argCount++
		query += fmt.Sprintf(" AND device_id = $%d", argCount)
		args = append(args, *filters.DeviceID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND opened_at >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND opened_at <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	if filters.SearchTerm != nil && *filters.SearchTerm != "" {
		argCount++
		query += fmt.Sprintf(" AND (session_number ILIKE $%d OR session_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+*filters.SearchTerm+"%")
	}

	query += " ORDER BY opened_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []pos.POSSession
	for rows.Next() {
		var session pos.POSSession
		err := rows.Scan(
			&session.ID, &session.OrganizationID, &session.SessionNumber, &session.SessionName,
			&session.DeviceID, &session.LocationID, &session.UserID, &session.ShiftID, &session.OpenedAt,
			&session.ClosedAt, &session.OpeningCash, &session.OpeningCard, &session.OpeningOther,
			&session.ExpectedCash, &session.ExpectedCard, &session.ExpectedOther,
			&session.CountedCash, &session.CountedCard, &session.CountedOther,
			&session.DifferenceCash, &session.DifferenceCard, &session.DifferenceOther,
			&session.Status, &session.ZReportNumber, &session.Notes,
			&session.CreatedBy, &session.UpdatedBy, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// Count counts sessions matching filters
func (r *POSSessionRepository) Count(ctx context.Context, orgID uuid.UUID, filters pos.POSSessionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM pos_sessions WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
	}

	if filters.DeviceID != nil {
		argCount++
		query += fmt.Sprintf(" AND device_id = $%d", argCount)
		args = append(args, *filters.DeviceID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND opened_at >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND opened_at <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// Update updates a POS session
func (r *POSSessionRepository) Update(ctx context.Context, session *pos.POSSession) error {
	if err := r.db.SetOrganizationContext(ctx, session.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions SET
			session_name = $3, closed_at = $4, expected_cash = $5, expected_card = $6,
			expected_other = $7, counted_cash = $8, counted_card = $9, counted_other = $10,
			difference_cash = $11, difference_card = $12, difference_other = $13,
			status = $14, z_report_number = $15, notes = $16, updated_by = $17, updated_at = $18
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		session.OrganizationID, session.ID, session.SessionName, session.ClosedAt,
		session.ExpectedCash, session.ExpectedCard, session.ExpectedOther,
		session.CountedCash, session.CountedCard, session.CountedOther,
		session.DifferenceCash, session.DifferenceCard, session.DifferenceOther,
		session.Status, session.ZReportNumber, session.Notes,
		session.UpdatedBy, session.UpdatedAt,
	)
	return err
}

// Delete soft-deletes a POS session
func (r *POSSessionRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetOpenSession retrieves the open session for a device
func (r *POSSessionRepository) GetOpenSession(ctx context.Context, orgID uuid.UUID, deviceID uuid.UUID) (*pos.POSSession, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, session_number, session_name, device_id, location_id, user_id,
		       shift_id, opened_at, closed_at, opening_cash, opening_card, opening_other,
		       expected_cash, expected_card, expected_other, counted_cash, counted_card,
		       counted_other, difference_cash, difference_card, difference_other, status,
		       z_report_number, notes, created_by, updated_by, created_at, updated_at, deleted_at
		FROM pos_sessions
		WHERE organization_id = $1 AND device_id = $2 AND status = 'open' AND deleted_at IS NULL
		LIMIT 1
	`

	var session pos.POSSession
	err := r.db.Pool.QueryRow(ctx, query, orgID, deviceID).Scan(
		&session.ID, &session.OrganizationID, &session.SessionNumber, &session.SessionName,
		&session.DeviceID, &session.LocationID, &session.UserID, &session.ShiftID, &session.OpenedAt,
		&session.ClosedAt, &session.OpeningCash, &session.OpeningCard, &session.OpeningOther,
		&session.ExpectedCash, &session.ExpectedCard, &session.ExpectedOther,
		&session.CountedCash, &session.CountedCard, &session.CountedOther,
		&session.DifferenceCash, &session.DifferenceCard, &session.DifferenceOther,
		&session.Status, &session.ZReportNumber, &session.Notes,
		&session.CreatedBy, &session.UpdatedBy, &session.CreatedAt, &session.UpdatedAt, &session.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// UpdateSessionStatus updates the status of a session
func (r *POSSessionRepository) UpdateSessionStatus(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, status pos.SessionStatus) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions
		SET status = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, sessionID, status, time.Now())
	return err
}

// UpdateExpectedAmounts updates the expected amounts
func (r *POSSessionRepository) UpdateExpectedAmounts(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, cash, card, other float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions
		SET expected_cash = $3, expected_card = $4, expected_other = $5, updated_at = $6
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, sessionID, cash, card, other, time.Now())
	return err
}

// UpdateCountedAmounts updates the counted amounts
func (r *POSSessionRepository) UpdateCountedAmounts(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID, cash, card, other float64) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions
		SET counted_cash = $3, counted_card = $4, counted_other = $5, updated_at = $6
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, sessionID, cash, card, other, time.Now())
	return err
}

// RecalculateDifferences recalculates the differences
func (r *POSSessionRepository) RecalculateDifferences(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE pos_sessions
		SET
			difference_cash = COALESCE(counted_cash, 0) - expected_cash,
			difference_card = COALESCE(counted_card, 0) - expected_card,
			difference_other = COALESCE(counted_other, 0) - expected_other,
			updated_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, sessionID, time.Now())
	return err
}

// CashDrawerRepository implements pos.CashDrawerRepository
type CashDrawerRepository struct {
	db *DB
}

// NewCashDrawerRepository creates a new cash drawer repository
func NewCashDrawerRepository(db *DB) *CashDrawerRepository {
	return &CashDrawerRepository{db: db}
}

// Create creates a new cash drawer
func (r *CashDrawerRepository) Create(ctx context.Context, drawer *pos.CashDrawer) error {
	if err := r.db.SetOrganizationContext(ctx, drawer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO cash_drawers (
			id, organization_id, drawer_code, drawer_name, location_id, device_id,
			is_active, notes, created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		drawer.ID, drawer.OrganizationID, drawer.DrawerCode, drawer.DrawerName,
		drawer.LocationID, drawer.DeviceID, drawer.IsActive, drawer.Notes,
		drawer.CreatedBy, drawer.CreatedAt, drawer.UpdatedAt,
	)
	return err
}

// Get retrieves a cash drawer by ID
func (r *CashDrawerRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*pos.CashDrawer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, drawer_code, drawer_name, location_id, device_id,
		       is_active, notes, created_by, created_at, updated_at, deleted_at
		FROM cash_drawers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var drawer pos.CashDrawer
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&drawer.ID, &drawer.OrganizationID, &drawer.DrawerCode, &drawer.DrawerName,
		&drawer.LocationID, &drawer.DeviceID, &drawer.IsActive, &drawer.Notes,
		&drawer.CreatedBy, &drawer.CreatedAt, &drawer.UpdatedAt, &drawer.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

// GetByCode retrieves a cash drawer by code
func (r *CashDrawerRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*pos.CashDrawer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, drawer_code, drawer_name, location_id, device_id,
		       is_active, notes, created_by, created_at, updated_at, deleted_at
		FROM cash_drawers
		WHERE organization_id = $1 AND LOWER(drawer_code) = LOWER($2) AND deleted_at IS NULL
	`

	var drawer pos.CashDrawer
	err := r.db.Pool.QueryRow(ctx, query, orgID, code).Scan(
		&drawer.ID, &drawer.OrganizationID, &drawer.DrawerCode, &drawer.DrawerName,
		&drawer.LocationID, &drawer.DeviceID, &drawer.IsActive, &drawer.Notes,
		&drawer.CreatedBy, &drawer.CreatedAt, &drawer.UpdatedAt, &drawer.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

// List retrieves cash drawers with filters
func (r *CashDrawerRepository) List(ctx context.Context, orgID uuid.UUID, filters pos.CashDrawerFilters) ([]pos.CashDrawer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, drawer_code, drawer_name, location_id, device_id,
		       is_active, notes, created_by, created_at, updated_at, deleted_at
		FROM cash_drawers
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.SearchTerm != nil && *filters.SearchTerm != "" {
		argCount++
		query += fmt.Sprintf(" AND (drawer_code ILIKE $%d OR drawer_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+*filters.SearchTerm+"%")
	}

	query += " ORDER BY drawer_code ASC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drawers []pos.CashDrawer
	for rows.Next() {
		var drawer pos.CashDrawer
		err := rows.Scan(
			&drawer.ID, &drawer.OrganizationID, &drawer.DrawerCode, &drawer.DrawerName,
			&drawer.LocationID, &drawer.DeviceID, &drawer.IsActive, &drawer.Notes,
			&drawer.CreatedBy, &drawer.CreatedAt, &drawer.UpdatedAt, &drawer.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		drawers = append(drawers, drawer)
	}

	return drawers, rows.Err()
}

// Count counts cash drawers matching filters
func (r *CashDrawerRepository) Count(ctx context.Context, orgID uuid.UUID, filters pos.CashDrawerFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM cash_drawers WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// Update updates a cash drawer
func (r *CashDrawerRepository) Update(ctx context.Context, drawer *pos.CashDrawer) error {
	if err := r.db.SetOrganizationContext(ctx, drawer.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cash_drawers SET
			drawer_name = $3, device_id = $4, is_active = $5, notes = $6, updated_at = $7
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		drawer.OrganizationID, drawer.ID, drawer.DrawerName, drawer.DeviceID,
		drawer.IsActive, drawer.Notes, drawer.UpdatedAt,
	)
	return err
}

// Delete soft-deletes a cash drawer
func (r *CashDrawerRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cash_drawers
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetByLocation retrieves all cash drawers for a location
func (r *CashDrawerRepository) GetByLocation(ctx context.Context, orgID uuid.UUID, locationID uuid.UUID) ([]pos.CashDrawer, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, drawer_code, drawer_name, location_id, device_id,
		       is_active, notes, created_by, created_at, updated_at, deleted_at
		FROM cash_drawers
		WHERE organization_id = $1 AND location_id = $2 AND deleted_at IS NULL
		ORDER BY drawer_code ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, locationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drawers []pos.CashDrawer
	for rows.Next() {
		var drawer pos.CashDrawer
		err := rows.Scan(
			&drawer.ID, &drawer.OrganizationID, &drawer.DrawerCode, &drawer.DrawerName,
			&drawer.LocationID, &drawer.DeviceID, &drawer.IsActive, &drawer.Notes,
			&drawer.CreatedBy, &drawer.CreatedAt, &drawer.UpdatedAt, &drawer.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		drawers = append(drawers, drawer)
	}

	return drawers, rows.Err()
}

// CashMovementRepository implements pos.CashMovementRepository
type CashMovementRepository struct {
	db *DB
}

// NewCashMovementRepository creates a new cash movement repository
func NewCashMovementRepository(db *DB) *CashMovementRepository {
	return &CashMovementRepository{db: db}
}

// Create creates a new cash movement
func (r *CashMovementRepository) Create(ctx context.Context, movement *pos.CashMovement) error {
	if err := r.db.SetOrganizationContext(ctx, movement.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO cash_movements (
			id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
			reason_code, reason_description, user_id, requires_approval, approved_by,
			approved_at, notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		movement.ID, movement.OrganizationID, movement.POSSessionID, movement.CashDrawerID,
		movement.MovementType, movement.Amount, movement.ReasonCode, movement.ReasonDescription,
		movement.UserID, movement.RequiresApproval, movement.ApprovedBy, movement.ApprovedAt,
		movement.Notes, movement.CreatedAt,
	)
	return err
}

// Get retrieves a cash movement by ID
func (r *CashMovementRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*pos.CashMovement, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
		       reason_code, reason_description, user_id, requires_approval, approved_by,
		       approved_at, notes, created_at, deleted_at
		FROM cash_movements
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var movement pos.CashMovement
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&movement.ID, &movement.OrganizationID, &movement.POSSessionID, &movement.CashDrawerID,
		&movement.MovementType, &movement.Amount, &movement.ReasonCode, &movement.ReasonDescription,
		&movement.UserID, &movement.RequiresApproval, &movement.ApprovedBy, &movement.ApprovedAt,
		&movement.Notes, &movement.CreatedAt, &movement.DeletedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

// List retrieves cash movements with filters
func (r *CashMovementRepository) List(ctx context.Context, orgID uuid.UUID, filters pos.CashMovementFilters) ([]pos.CashMovement, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
		       reason_code, reason_description, user_id, requires_approval, approved_by,
		       approved_at, notes, created_at, deleted_at
		FROM cash_movements
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.POSSessionID != nil {
		argCount++
		query += fmt.Sprintf(" AND pos_session_id = $%d", argCount)
		args = append(args, *filters.POSSessionID)
	}

	if filters.CashDrawerID != nil {
		argCount++
		query += fmt.Sprintf(" AND cash_drawer_id = $%d", argCount)
		args = append(args, *filters.CashDrawerID)
	}

	if filters.MovementType != nil {
		argCount++
		query += fmt.Sprintf(" AND movement_type = $%d", argCount)
		args = append(args, *filters.MovementType)
	}

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filters.DateTo)
	}

	if filters.RequiresApproval != nil {
		argCount++
		query += fmt.Sprintf(" AND requires_approval = $%d", argCount)
		args = append(args, *filters.RequiresApproval)
	}

	query += " ORDER BY created_at DESC"

	if filters.PageSize > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)

		if filters.Page > 1 {
			argCount++
			offset := (filters.Page - 1) * filters.PageSize
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, offset)
		}
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []pos.CashMovement
	for rows.Next() {
		var movement pos.CashMovement
		err := rows.Scan(
			&movement.ID, &movement.OrganizationID, &movement.POSSessionID, &movement.CashDrawerID,
			&movement.MovementType, &movement.Amount, &movement.ReasonCode, &movement.ReasonDescription,
			&movement.UserID, &movement.RequiresApproval, &movement.ApprovedBy, &movement.ApprovedAt,
			&movement.Notes, &movement.CreatedAt, &movement.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}

	return movements, rows.Err()
}

// Count counts cash movements matching filters
func (r *CashMovementRepository) Count(ctx context.Context, orgID uuid.UUID, filters pos.CashMovementFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM cash_movements WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.POSSessionID != nil {
		argCount++
		query += fmt.Sprintf(" AND pos_session_id = $%d", argCount)
		args = append(args, *filters.POSSessionID)
	}

	if filters.CashDrawerID != nil {
		argCount++
		query += fmt.Sprintf(" AND cash_drawer_id = $%d", argCount)
		args = append(args, *filters.CashDrawerID)
	}

	if filters.MovementType != nil {
		argCount++
		query += fmt.Sprintf(" AND movement_type = $%d", argCount)
		args = append(args, *filters.MovementType)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

// Update updates a cash movement
func (r *CashMovementRepository) Update(ctx context.Context, movement *pos.CashMovement) error {
	if err := r.db.SetOrganizationContext(ctx, movement.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cash_movements SET
			reason_description = $3, notes = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		movement.OrganizationID, movement.ID, movement.ReasonDescription, movement.Notes,
	)
	return err
}

// Delete soft-deletes a cash movement
func (r *CashMovementRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cash_movements
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// GetForSession retrieves all cash movements for a session
func (r *CashMovementRepository) GetForSession(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) ([]pos.CashMovement, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
		       reason_code, reason_description, user_id, requires_approval, approved_by,
		       approved_at, notes, created_at, deleted_at
		FROM cash_movements
		WHERE organization_id = $1 AND pos_session_id = $2 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []pos.CashMovement
	for rows.Next() {
		var movement pos.CashMovement
		err := rows.Scan(
			&movement.ID, &movement.OrganizationID, &movement.POSSessionID, &movement.CashDrawerID,
			&movement.MovementType, &movement.Amount, &movement.ReasonCode, &movement.ReasonDescription,
			&movement.UserID, &movement.RequiresApproval, &movement.ApprovedBy, &movement.ApprovedAt,
			&movement.Notes, &movement.CreatedAt, &movement.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}

	return movements, rows.Err()
}

// GetForDrawer retrieves all cash movements for a drawer
func (r *CashMovementRepository) GetForDrawer(ctx context.Context, orgID uuid.UUID, drawerID uuid.UUID) ([]pos.CashMovement, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
		       reason_code, reason_description, user_id, requires_approval, approved_by,
		       approved_at, notes, created_at, deleted_at
		FROM cash_movements
		WHERE organization_id = $1 AND cash_drawer_id = $2 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, drawerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []pos.CashMovement
	for rows.Next() {
		var movement pos.CashMovement
		err := rows.Scan(
			&movement.ID, &movement.OrganizationID, &movement.POSSessionID, &movement.CashDrawerID,
			&movement.MovementType, &movement.Amount, &movement.ReasonCode, &movement.ReasonDescription,
			&movement.UserID, &movement.RequiresApproval, &movement.ApprovedBy, &movement.ApprovedAt,
			&movement.Notes, &movement.CreatedAt, &movement.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}

	return movements, rows.Err()
}

// GetPendingApprovals retrieves all pending approval movements
func (r *CashMovementRepository) GetPendingApprovals(ctx context.Context, orgID uuid.UUID) ([]pos.CashMovement, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, pos_session_id, cash_drawer_id, movement_type, amount,
		       reason_code, reason_description, user_id, requires_approval, approved_by,
		       approved_at, notes, created_at, deleted_at
		FROM cash_movements
		WHERE organization_id = $1 AND requires_approval = true AND approved_by IS NULL AND deleted_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []pos.CashMovement
	for rows.Next() {
		var movement pos.CashMovement
		err := rows.Scan(
			&movement.ID, &movement.OrganizationID, &movement.POSSessionID, &movement.CashDrawerID,
			&movement.MovementType, &movement.Amount, &movement.ReasonCode, &movement.ReasonDescription,
			&movement.UserID, &movement.RequiresApproval, &movement.ApprovedBy, &movement.ApprovedAt,
			&movement.Notes, &movement.CreatedAt, &movement.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		movements = append(movements, movement)
	}

	return movements, rows.Err()
}

// ApproveMovement approves a cash movement
func (r *CashMovementRepository) ApproveMovement(ctx context.Context, orgID uuid.UUID, movementID uuid.UUID, approvedBy uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE cash_movements
		SET approved_by = $3, approved_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, movementID, approvedBy, time.Now())
	return err
}

// CalculateSessionBalance calculates the total balance impact for a session
func (r *CashMovementRepository) CalculateSessionBalance(ctx context.Context, orgID uuid.UUID, sessionID uuid.UUID) (float64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := `
		SELECT COALESCE(SUM(amount), 0) as total_balance
		FROM cash_movements
		WHERE organization_id = $1 AND pos_session_id = $2 AND deleted_at IS NULL
		  AND (approved_by IS NOT NULL OR requires_approval = false)
	`

	var balance float64
	err := r.db.Pool.QueryRow(ctx, query, orgID, sessionID).Scan(&balance)
	return balance, err
}
