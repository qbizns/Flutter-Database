package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/restaurant"
)

type RestaurantTableRepository struct {
	db *DB
}

func NewRestaurantTableRepository(db *DB) *RestaurantTableRepository {
	return &RestaurantTableRepository{db: db}
}

// ============ RESTAURANT TABLE CRUD ============

func (r *RestaurantTableRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.RestaurantTableFilters) ([]restaurant.RestaurantTable, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, floor_plan_id, section_id,
		       table_number, table_name, min_capacity, max_capacity, table_shape,
		       is_combinable, position_x, position_y, rotation, status,
		       current_covers, seated_at, current_waiter_id, is_active,
		       allow_online_reservation, display_order, color_code, icon, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM restaurant_tables
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (table_number ILIKE $%d OR table_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.FloorPlanID != nil {
		argCount++
		query += fmt.Sprintf(" AND floor_plan_id = $%d", argCount)
		args = append(args, *filters.FloorPlanID)
	}

	if filters.SectionID != nil {
		argCount++
		query += fmt.Sprintf(" AND section_id = $%d", argCount)
		args = append(args, *filters.SectionID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.WaiterID != nil {
		argCount++
		query += fmt.Sprintf(" AND current_waiter_id = $%d", argCount)
		args = append(args, *filters.WaiterID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	query += " ORDER BY display_order, table_number ASC"

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

	var tableList []restaurant.RestaurantTable
	for rows.Next() {
		var t restaurant.RestaurantTable
		var metadata sql.NullString

		err := rows.Scan(
			&t.ID, &t.OrganizationID, &t.LocationID, &t.FloorPlanID, &t.SectionID,
			&t.TableNumber, &t.TableName, &t.MinCapacity, &t.MaxCapacity, &t.TableShape,
			&t.IsCombinable, &t.PositionX, &t.PositionY, &t.Rotation, &t.Status,
			&t.CurrentCovers, &t.SeatedAt, &t.CurrentWaiterID, &t.IsActive,
			&t.AllowOnlineReservation, &t.DisplayOrder, &t.ColorCode, &t.Icon, &t.Notes,
			&metadata, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		if metadata.Valid {
			json.Unmarshal([]byte(metadata.String), &t.Metadata)
		}

		tableList = append(tableList, t)
	}

	return tableList, rows.Err()
}

func (r *RestaurantTableRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.RestaurantTableFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM restaurant_tables WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (table_number ILIKE $%d OR table_name ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
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

func (r *RestaurantTableRepository) Create(ctx context.Context, table *restaurant.RestaurantTable) error {
	if err := r.db.SetOrganizationContext(ctx, table.OrganizationID.String()); err != nil {
		return err
	}

	metadata := sql.NullString{}
	if table.Metadata != nil {
		if b, err := json.Marshal(table.Metadata); err == nil {
			metadata = sql.NullString{String: string(b), Valid: true}
		}
	}

	query := `
		INSERT INTO restaurant_tables (
			id, organization_id, location_id, floor_plan_id, section_id,
			table_number, table_name, min_capacity, max_capacity, table_shape,
			is_combinable, position_x, position_y, rotation, status,
			current_covers, seated_at, current_waiter_id, is_active,
			allow_online_reservation, display_order, color_code, icon, notes,
			metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		table.ID, table.OrganizationID, table.LocationID, table.FloorPlanID, table.SectionID,
		table.TableNumber, table.TableName, table.MinCapacity, table.MaxCapacity, table.TableShape,
		table.IsCombinable, table.PositionX, table.PositionY, table.Rotation, table.Status,
		table.CurrentCovers, table.SeatedAt, table.CurrentWaiterID, table.IsActive,
		table.AllowOnlineReservation, table.DisplayOrder, table.ColorCode, table.Icon, table.Notes,
		metadata, table.CreatedAt, table.UpdatedAt, table.CreatedBy,
	)
	return err
}

func (r *RestaurantTableRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.RestaurantTable, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, floor_plan_id, section_id,
		       table_number, table_name, min_capacity, max_capacity, table_shape,
		       is_combinable, position_x, position_y, rotation, status,
		       current_covers, seated_at, current_waiter_id, is_active,
		       allow_online_reservation, display_order, color_code, icon, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM restaurant_tables
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var t restaurant.RestaurantTable
	var metadata sql.NullString

	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&t.ID, &t.OrganizationID, &t.LocationID, &t.FloorPlanID, &t.SectionID,
		&t.TableNumber, &t.TableName, &t.MinCapacity, &t.MaxCapacity, &t.TableShape,
		&t.IsCombinable, &t.PositionX, &t.PositionY, &t.Rotation, &t.Status,
		&t.CurrentCovers, &t.SeatedAt, &t.CurrentWaiterID, &t.IsActive,
		&t.AllowOnlineReservation, &t.DisplayOrder, &t.ColorCode, &t.Icon, &t.Notes,
		&metadata, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata.Valid {
		json.Unmarshal([]byte(metadata.String), &t.Metadata)
	}

	return &t, nil
}

func (r *RestaurantTableRepository) GetByNumber(ctx context.Context, orgID uuid.UUID, locationID uuid.UUID, tableNumber string) (*restaurant.RestaurantTable, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, floor_plan_id, section_id,
		       table_number, table_name, min_capacity, max_capacity, table_shape,
		       is_combinable, position_x, position_y, rotation, status,
		       current_covers, seated_at, current_waiter_id, is_active,
		       allow_online_reservation, display_order, color_code, icon, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM restaurant_tables
		WHERE organization_id = $1 AND location_id = $2 AND table_number = $3 AND deleted_at IS NULL
	`

	var t restaurant.RestaurantTable
	var metadata sql.NullString

	err := r.db.Pool.QueryRow(ctx, query, orgID, locationID, tableNumber).Scan(
		&t.ID, &t.OrganizationID, &t.LocationID, &t.FloorPlanID, &t.SectionID,
		&t.TableNumber, &t.TableName, &t.MinCapacity, &t.MaxCapacity, &t.TableShape,
		&t.IsCombinable, &t.PositionX, &t.PositionY, &t.Rotation, &t.Status,
		&t.CurrentCovers, &t.SeatedAt, &t.CurrentWaiterID, &t.IsActive,
		&t.AllowOnlineReservation, &t.DisplayOrder, &t.ColorCode, &t.Icon, &t.Notes,
		&metadata, &t.CreatedAt, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata.Valid {
		json.Unmarshal([]byte(metadata.String), &t.Metadata)
	}

	return &t, nil
}

func (r *RestaurantTableRepository) Update(ctx context.Context, table *restaurant.RestaurantTable) error {
	if err := r.db.SetOrganizationContext(ctx, table.OrganizationID.String()); err != nil {
		return err
	}

	metadata := sql.NullString{}
	if table.Metadata != nil {
		if b, err := json.Marshal(table.Metadata); err == nil {
			metadata = sql.NullString{String: string(b), Valid: true}
		}
	}

	query := `
		UPDATE restaurant_tables SET
			floor_plan_id = $3, section_id = $4, table_name = $5,
			min_capacity = $6, max_capacity = $7, table_shape = $8,
			is_combinable = $9, position_x = $10, position_y = $11, rotation = $12,
			status = $13, current_covers = $14, seated_at = $15, current_waiter_id = $16,
			is_active = $17, allow_online_reservation = $18, display_order = $19,
			color_code = $20, icon = $21, notes = $22, metadata = $23,
			updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		table.OrganizationID, table.ID, table.FloorPlanID, table.SectionID, table.TableName,
		table.MinCapacity, table.MaxCapacity, table.TableShape, table.IsCombinable,
		table.PositionX, table.PositionY, table.Rotation, table.Status, table.CurrentCovers,
		table.SeatedAt, table.CurrentWaiterID, table.IsActive, table.AllowOnlineReservation,
		table.DisplayOrder, table.ColorCode, table.Icon, table.Notes, metadata,
		table.UpdatedAt, table.UpdatedBy,
	)
	return err
}

func (r *RestaurantTableRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE restaurant_tables
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *RestaurantTableRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE restaurant_tables
		SET status = $3, updated_at = $4
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, time.Now())
	return err
}

// ============ RESERVATION CRUD ============

type ReservationRepository struct {
	db *DB
}

func NewReservationRepository(db *DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.ReservationFilters) ([]restaurant.Reservation, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, table_id, customer_id,
		       reservation_number, reservation_date, reservation_time, duration_minutes,
		       party_size, customer_name, customer_phone, customer_email, status,
		       assigned_waiter_id, assigned_at, seated_at, completed_at,
		       special_requests, occasion, dietary_restrictions, confirmation_code,
		       confirmed_at, confirmed_by, reminder_sent_at, notification_preferences,
		       cancelled_at, cancelled_by, cancellation_reason, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM reservations
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (customer_name ILIKE $%d OR reservation_number ILIKE $%d OR customer_phone ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ReservationDate != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(reservation_date) = $%d", argCount)
		args = append(args, filters.ReservationDate)
	}

	if filters.WaiterID != nil {
		argCount++
		query += fmt.Sprintf(" AND assigned_waiter_id = $%d", argCount)
		args = append(args, *filters.WaiterID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	query += " ORDER BY reservation_date, reservation_time ASC"

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

	var resList []restaurant.Reservation
	for rows.Next() {
		var res restaurant.Reservation
		var metadata, notificationPrefs sql.NullString

		err := rows.Scan(
			&res.ID, &res.OrganizationID, &res.LocationID, &res.TableID, &res.CustomerID,
			&res.ReservationNumber, &res.ReservationDate, &res.ReservationTime, &res.DurationMinutes,
			&res.PartySize, &res.CustomerName, &res.CustomerPhone, &res.CustomerEmail, &res.Status,
			&res.AssignedWaiterID, &res.AssignedAt, &res.SeatedAt, &res.CompletedAt,
			&res.SpecialRequests, &res.Occasion, &res.DietaryRestrictions, &res.ConfirmationCode,
			&res.ConfirmedAt, &res.ConfirmedBy, &res.ReminderSentAt, &notificationPrefs,
			&res.CancelledAt, &res.CancelledBy, &res.CancellationReason, &res.Notes,
			&metadata, &res.CreatedAt, &res.UpdatedAt, &res.CreatedBy, &res.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}

		if metadata.Valid {
			json.Unmarshal([]byte(metadata.String), &res.Metadata)
		}
		if notificationPrefs.Valid {
			json.Unmarshal([]byte(notificationPrefs.String), &res.NotificationPreferences)
		}

		resList = append(resList, res)
	}

	return resList, rows.Err()
}

func (r *ReservationRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.ReservationFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM reservations WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.ReservationDate != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(reservation_date) = $%d", argCount)
		args = append(args, filters.ReservationDate)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *ReservationRepository) Create(ctx context.Context, reservation *restaurant.Reservation) error {
	if err := r.db.SetOrganizationContext(ctx, reservation.OrganizationID.String()); err != nil {
		return err
	}

	metadata := sql.NullString{}
	if reservation.Metadata != nil {
		if b, err := json.Marshal(reservation.Metadata); err == nil {
			metadata = sql.NullString{String: string(b), Valid: true}
		}
	}

	notificationPrefs := sql.NullString{}
	if reservation.NotificationPreferences != nil {
		if b, err := json.Marshal(reservation.NotificationPreferences); err == nil {
			notificationPrefs = sql.NullString{String: string(b), Valid: true}
		}
	}

	query := `
		INSERT INTO reservations (
			id, organization_id, location_id, table_id, customer_id,
			reservation_number, reservation_date, reservation_time, duration_minutes,
			party_size, customer_name, customer_phone, customer_email, status,
			assigned_waiter_id, assigned_at, seated_at, completed_at,
			special_requests, occasion, dietary_restrictions, confirmation_code,
			confirmed_at, confirmed_by, reminder_sent_at, notification_preferences,
			cancelled_at, cancelled_by, cancellation_reason, notes,
			metadata, created_at, updated_at, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32, $33
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		reservation.ID, reservation.OrganizationID, reservation.LocationID, reservation.TableID, reservation.CustomerID,
		reservation.ReservationNumber, reservation.ReservationDate, reservation.ReservationTime, reservation.DurationMinutes,
		reservation.PartySize, reservation.CustomerName, reservation.CustomerPhone, reservation.CustomerEmail, reservation.Status,
		reservation.AssignedWaiterID, reservation.AssignedAt, reservation.SeatedAt, reservation.CompletedAt,
		reservation.SpecialRequests, reservation.Occasion, reservation.DietaryRestrictions, reservation.ConfirmationCode,
		reservation.ConfirmedAt, reservation.ConfirmedBy, reservation.ReminderSentAt, notificationPrefs,
		reservation.CancelledAt, reservation.CancelledBy, reservation.CancellationReason, reservation.Notes,
		metadata, reservation.CreatedAt, reservation.UpdatedAt, reservation.CreatedBy,
	)
	return err
}

func (r *ReservationRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.Reservation, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, table_id, customer_id,
		       reservation_number, reservation_date, reservation_time, duration_minutes,
		       party_size, customer_name, customer_phone, customer_email, status,
		       assigned_waiter_id, assigned_at, seated_at, completed_at,
		       special_requests, occasion, dietary_restrictions, confirmation_code,
		       confirmed_at, confirmed_by, reminder_sent_at, notification_preferences,
		       cancelled_at, cancelled_by, cancellation_reason, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM reservations
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var res restaurant.Reservation
	var metadata, notificationPrefs sql.NullString

	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&res.ID, &res.OrganizationID, &res.LocationID, &res.TableID, &res.CustomerID,
		&res.ReservationNumber, &res.ReservationDate, &res.ReservationTime, &res.DurationMinutes,
		&res.PartySize, &res.CustomerName, &res.CustomerPhone, &res.CustomerEmail, &res.Status,
		&res.AssignedWaiterID, &res.AssignedAt, &res.SeatedAt, &res.CompletedAt,
		&res.SpecialRequests, &res.Occasion, &res.DietaryRestrictions, &res.ConfirmationCode,
		&res.ConfirmedAt, &res.ConfirmedBy, &res.ReminderSentAt, &notificationPrefs,
		&res.CancelledAt, &res.CancelledBy, &res.CancellationReason, &res.Notes,
		&metadata, &res.CreatedAt, &res.UpdatedAt, &res.CreatedBy, &res.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata.Valid {
		json.Unmarshal([]byte(metadata.String), &res.Metadata)
	}
	if notificationPrefs.Valid {
		json.Unmarshal([]byte(notificationPrefs.String), &res.NotificationPreferences)
	}

	return &res, nil
}

func (r *ReservationRepository) GetByNumber(ctx context.Context, orgID uuid.UUID, reservationNumber string) (*restaurant.Reservation, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, table_id, customer_id,
		       reservation_number, reservation_date, reservation_time, duration_minutes,
		       party_size, customer_name, customer_phone, customer_email, status,
		       assigned_waiter_id, assigned_at, seated_at, completed_at,
		       special_requests, occasion, dietary_restrictions, confirmation_code,
		       confirmed_at, confirmed_by, reminder_sent_at, notification_preferences,
		       cancelled_at, cancelled_by, cancellation_reason, notes,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM reservations
		WHERE organization_id = $1 AND reservation_number = $2 AND deleted_at IS NULL
	`

	var res restaurant.Reservation
	var metadata, notificationPrefs sql.NullString

	err := r.db.Pool.QueryRow(ctx, query, orgID, reservationNumber).Scan(
		&res.ID, &res.OrganizationID, &res.LocationID, &res.TableID, &res.CustomerID,
		&res.ReservationNumber, &res.ReservationDate, &res.ReservationTime, &res.DurationMinutes,
		&res.PartySize, &res.CustomerName, &res.CustomerPhone, &res.CustomerEmail, &res.Status,
		&res.AssignedWaiterID, &res.AssignedAt, &res.SeatedAt, &res.CompletedAt,
		&res.SpecialRequests, &res.Occasion, &res.DietaryRestrictions, &res.ConfirmationCode,
		&res.ConfirmedAt, &res.ConfirmedBy, &res.ReminderSentAt, &notificationPrefs,
		&res.CancelledAt, &res.CancelledBy, &res.CancellationReason, &res.Notes,
		&metadata, &res.CreatedAt, &res.UpdatedAt, &res.CreatedBy, &res.UpdatedBy,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if metadata.Valid {
		json.Unmarshal([]byte(metadata.String), &res.Metadata)
	}
	if notificationPrefs.Valid {
		json.Unmarshal([]byte(notificationPrefs.String), &res.NotificationPreferences)
	}

	return &res, nil
}

func (r *ReservationRepository) Update(ctx context.Context, reservation *restaurant.Reservation) error {
	if err := r.db.SetOrganizationContext(ctx, reservation.OrganizationID.String()); err != nil {
		return err
	}

	metadata := sql.NullString{}
	if reservation.Metadata != nil {
		if b, err := json.Marshal(reservation.Metadata); err == nil {
			metadata = sql.NullString{String: string(b), Valid: true}
		}
	}

	notificationPrefs := sql.NullString{}
	if reservation.NotificationPreferences != nil {
		if b, err := json.Marshal(reservation.NotificationPreferences); err == nil {
			notificationPrefs = sql.NullString{String: string(b), Valid: true}
		}
	}

	query := `
		UPDATE reservations SET
			table_id = $3, customer_id = $4, reservation_date = $5, reservation_time = $6,
			duration_minutes = $7, party_size = $8, customer_name = $9, customer_phone = $10,
			customer_email = $11, status = $12, assigned_waiter_id = $13, assigned_at = $14,
			seated_at = $15, completed_at = $16, special_requests = $17, occasion = $18,
			dietary_restrictions = $19, confirmation_code = $20, confirmed_at = $21,
			confirmed_by = $22, reminder_sent_at = $23, notification_preferences = $24,
			cancelled_at = $25, cancelled_by = $26, cancellation_reason = $27, notes = $28,
			metadata = $29, updated_at = $30, updated_by = $31
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		reservation.OrganizationID, reservation.ID, reservation.TableID, reservation.CustomerID,
		reservation.ReservationDate, reservation.ReservationTime, reservation.DurationMinutes,
		reservation.PartySize, reservation.CustomerName, reservation.CustomerPhone,
		reservation.CustomerEmail, reservation.Status, reservation.AssignedWaiterID,
		reservation.AssignedAt, reservation.SeatedAt, reservation.CompletedAt,
		reservation.SpecialRequests, reservation.Occasion, reservation.DietaryRestrictions,
		reservation.ConfirmationCode, reservation.ConfirmedAt, reservation.ConfirmedBy,
		reservation.ReminderSentAt, notificationPrefs, reservation.CancelledAt,
		reservation.CancelledBy, reservation.CancellationReason, reservation.Notes,
		metadata, reservation.UpdatedAt, reservation.UpdatedBy,
	)
	return err
}

func (r *ReservationRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE reservations
		SET deleted_at = $3
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}
