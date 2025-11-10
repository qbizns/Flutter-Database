package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/domain/restaurant"
)

// KITCHEN STATION REPOSITORY

type KitchenStationRepository struct {
	db *DB
}

func NewKitchenStationRepository(db *DB) *KitchenStationRepository {
	return &KitchenStationRepository{db: db}
}

func (r *KitchenStationRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.KitchenStationFilters) ([]restaurant.KitchenStation, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, station_name, station_code, station_type,
		       description, display_order, color_code, printer_id, is_active, auto_print_tickets,
		       alert_sound_enabled, display_config, created_at, updated_at, created_by, updated_by
		FROM kitchen_stations
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (station_name ILIKE $%d OR station_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.StationType != nil {
		argCount++
		query += fmt.Sprintf(" AND station_type = $%d", argCount)
		args = append(args, *filters.StationType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	query += " ORDER BY display_order ASC, station_name ASC"

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

	var stations []restaurant.KitchenStation
	for rows.Next() {
		var s restaurant.KitchenStation
		err := rows.Scan(&s.ID, &s.OrganizationID, &s.LocationID, &s.StationName, &s.StationCode,
			&s.StationType, &s.Description, &s.DisplayOrder, &s.ColorCode, &s.PrinterID,
			&s.IsActive, &s.AutoPrintTickets, &s.AlertSoundEnabled, &s.DisplayConfig,
			&s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy)
		if err != nil {
			return nil, err
		}
		stations = append(stations, s)
	}

	return stations, rows.Err()
}

func (r *KitchenStationRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.KitchenStationFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM kitchen_stations WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Search != "" {
		argCount++
		query += fmt.Sprintf(" AND (station_name ILIKE $%d OR station_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+filters.Search+"%")
	}

	if filters.StationType != nil {
		argCount++
		query += fmt.Sprintf(" AND station_type = $%d", argCount)
		args = append(args, *filters.StationType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *KitchenStationRepository) Create(ctx context.Context, station *restaurant.KitchenStation) error {
	if err := r.db.SetOrganizationContext(ctx, station.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO kitchen_stations (id, organization_id, location_id, station_name, station_code,
		    station_type, description, display_order, color_code, printer_id, is_active,
		    auto_print_tickets, alert_sound_enabled, display_config, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		station.ID, station.OrganizationID, station.LocationID, station.StationName, station.StationCode,
		station.StationType, station.Description, station.DisplayOrder, station.ColorCode, station.PrinterID,
		station.IsActive, station.AutoPrintTickets, station.AlertSoundEnabled, station.DisplayConfig,
		station.CreatedAt, station.UpdatedAt, station.CreatedBy)
	return err
}

func (r *KitchenStationRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.KitchenStation, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, station_name, station_code, station_type,
		       description, display_order, color_code, printer_id, is_active, auto_print_tickets,
		       alert_sound_enabled, display_config, created_at, updated_at, created_by, updated_by
		FROM kitchen_stations
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var s restaurant.KitchenStation
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&s.ID, &s.OrganizationID, &s.LocationID, &s.StationName, &s.StationCode, &s.StationType,
		&s.Description, &s.DisplayOrder, &s.ColorCode, &s.PrinterID, &s.IsActive, &s.AutoPrintTickets,
		&s.AlertSoundEnabled, &s.DisplayConfig, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *KitchenStationRepository) Update(ctx context.Context, station *restaurant.KitchenStation) error {
	if err := r.db.SetOrganizationContext(ctx, station.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE kitchen_stations SET
		    station_name = $3, station_code = $4, station_type = $5, description = $6,
		    display_order = $7, color_code = $8, printer_id = $9, is_active = $10,
		    auto_print_tickets = $11, alert_sound_enabled = $12, display_config = $13,
		    updated_at = $14, updated_by = $15
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		station.OrganizationID, station.ID, station.StationName, station.StationCode, station.StationType,
		station.Description, station.DisplayOrder, station.ColorCode, station.PrinterID, station.IsActive,
		station.AutoPrintTickets, station.AlertSoundEnabled, station.DisplayConfig, station.UpdatedAt, station.UpdatedBy)
	return err
}

func (r *KitchenStationRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE kitchen_stations SET deleted_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

// ORDER REPOSITORY

type OrderRepository struct {
	db *DB
}

func NewOrderRepository(db *DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.OrderFilters) ([]restaurant.Order, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, order_number, display_number, order_type, table_id,
		       reservation_id, covers, customer_id, waiter_id, status, order_date, submitted_at,
		       kitchen_received_at, ready_at, served_at, completed_at, subtotal, tax_amount,
		       discount_amount, service_charge, total_amount, sale_id, shift_id, customer_notes,
		       kitchen_notes, internal_notes, metadata, created_at, updated_at, created_by, updated_by
		FROM orders
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.OrderType != nil {
		argCount++
		query += fmt.Sprintf(" AND order_type = $%d", argCount)
		args = append(args, *filters.OrderType)
	}

	if filters.TableID != nil {
		argCount++
		query += fmt.Sprintf(" AND table_id = $%d", argCount)
		args = append(args, *filters.TableID)
	}

	if filters.WaiterID != nil {
		argCount++
		query += fmt.Sprintf(" AND waiter_id = $%d", argCount)
		args = append(args, *filters.WaiterID)
	}

	if filters.CustomerID != nil {
		argCount++
		query += fmt.Sprintf(" AND customer_id = $%d", argCount)
		args = append(args, *filters.CustomerID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND order_date >= $%d", argCount)
		args = append(args, *filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND order_date <= $%d", argCount)
		args = append(args, *filters.ToDate)
	}

	query += " ORDER BY order_date DESC"

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

	var orders []restaurant.Order
	for rows.Next() {
		var o restaurant.Order
		err := rows.Scan(&o.ID, &o.OrganizationID, &o.LocationID, &o.OrderNumber, &o.DisplayNumber,
			&o.OrderType, &o.TableID, &o.ReservationID, &o.Covers, &o.CustomerID, &o.WaiterID,
			&o.Status, &o.OrderDate, &o.SubmittedAt, &o.KitchenReceivedAt, &o.ReadyAt, &o.ServedAt,
			&o.CompletedAt, &o.Subtotal, &o.TaxAmount, &o.DiscountAmount, &o.ServiceCharge,
			&o.TotalAmount, &o.SaleID, &o.ShiftID, &o.CustomerNotes, &o.KitchenNotes, &o.InternalNotes,
			&o.Metadata, &o.CreatedAt, &o.UpdatedAt, &o.CreatedBy, &o.UpdatedBy)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, rows.Err()
}

func (r *OrderRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.OrderFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM orders WHERE organization_id = $1 AND deleted_at IS NULL"
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

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *OrderRepository) Create(ctx context.Context, order *restaurant.Order) error {
	if err := r.db.SetOrganizationContext(ctx, order.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO orders (id, organization_id, location_id, order_number, display_number, order_type,
		    table_id, reservation_id, covers, customer_id, waiter_id, status, order_date, submitted_at,
		    kitchen_received_at, ready_at, served_at, completed_at, subtotal, tax_amount, discount_amount,
		    service_charge, total_amount, sale_id, shift_id, customer_notes, kitchen_notes, internal_notes,
		    metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
		    $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		order.ID, order.OrganizationID, order.LocationID, order.OrderNumber, order.DisplayNumber,
		order.OrderType, order.TableID, order.ReservationID, order.Covers, order.CustomerID, order.WaiterID,
		order.Status, order.OrderDate, order.SubmittedAt, order.KitchenReceivedAt, order.ReadyAt,
		order.ServedAt, order.CompletedAt, order.Subtotal, order.TaxAmount, order.DiscountAmount,
		order.ServiceCharge, order.TotalAmount, order.SaleID, order.ShiftID, order.CustomerNotes,
		order.KitchenNotes, order.InternalNotes, order.Metadata, order.CreatedAt, order.UpdatedAt, order.CreatedBy)
	return err
}

func (r *OrderRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.Order, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, order_number, display_number, order_type, table_id,
		       reservation_id, covers, customer_id, waiter_id, status, order_date, submitted_at,
		       kitchen_received_at, ready_at, served_at, completed_at, subtotal, tax_amount,
		       discount_amount, service_charge, total_amount, sale_id, shift_id, customer_notes,
		       kitchen_notes, internal_notes, metadata, created_at, updated_at, created_by, updated_by
		FROM orders
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var o restaurant.Order
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&o.ID, &o.OrganizationID, &o.LocationID, &o.OrderNumber, &o.DisplayNumber, &o.OrderType,
		&o.TableID, &o.ReservationID, &o.Covers, &o.CustomerID, &o.WaiterID, &o.Status,
		&o.OrderDate, &o.SubmittedAt, &o.KitchenReceivedAt, &o.ReadyAt, &o.ServedAt, &o.CompletedAt,
		&o.Subtotal, &o.TaxAmount, &o.DiscountAmount, &o.ServiceCharge, &o.TotalAmount, &o.SaleID,
		&o.ShiftID, &o.CustomerNotes, &o.KitchenNotes, &o.InternalNotes, &o.Metadata,
		&o.CreatedAt, &o.UpdatedAt, &o.CreatedBy, &o.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) GetByNumber(ctx context.Context, orgID uuid.UUID, orderNumber string) (*restaurant.Order, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, order_number, display_number, order_type, table_id,
		       reservation_id, covers, customer_id, waiter_id, status, order_date, submitted_at,
		       kitchen_received_at, ready_at, served_at, completed_at, subtotal, tax_amount,
		       discount_amount, service_charge, total_amount, sale_id, shift_id, customer_notes,
		       kitchen_notes, internal_notes, metadata, created_at, updated_at, created_by, updated_by
		FROM orders
		WHERE organization_id = $1 AND order_number = $2 AND deleted_at IS NULL
	`

	var o restaurant.Order
	err := r.db.Pool.QueryRow(ctx, query, orgID, orderNumber).Scan(
		&o.ID, &o.OrganizationID, &o.LocationID, &o.OrderNumber, &o.DisplayNumber, &o.OrderType,
		&o.TableID, &o.ReservationID, &o.Covers, &o.CustomerID, &o.WaiterID, &o.Status,
		&o.OrderDate, &o.SubmittedAt, &o.KitchenReceivedAt, &o.ReadyAt, &o.ServedAt, &o.CompletedAt,
		&o.Subtotal, &o.TaxAmount, &o.DiscountAmount, &o.ServiceCharge, &o.TotalAmount, &o.SaleID,
		&o.ShiftID, &o.CustomerNotes, &o.KitchenNotes, &o.InternalNotes, &o.Metadata,
		&o.CreatedAt, &o.UpdatedAt, &o.CreatedBy, &o.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *restaurant.Order) error {
	if err := r.db.SetOrganizationContext(ctx, order.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE orders SET
		    order_number = $3, display_number = $4, order_type = $5, table_id = $6, reservation_id = $7,
		    covers = $8, customer_id = $9, waiter_id = $10, status = $11, submitted_at = $12,
		    kitchen_received_at = $13, ready_at = $14, served_at = $15, completed_at = $16,
		    subtotal = $17, tax_amount = $18, discount_amount = $19, service_charge = $20,
		    total_amount = $21, sale_id = $22, shift_id = $23, customer_notes = $24, kitchen_notes = $25,
		    internal_notes = $26, metadata = $27, updated_at = $28, updated_by = $29
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		order.OrganizationID, order.ID, order.OrderNumber, order.DisplayNumber, order.OrderType,
		order.TableID, order.ReservationID, order.Covers, order.CustomerID, order.WaiterID, order.Status,
		order.SubmittedAt, order.KitchenReceivedAt, order.ReadyAt, order.ServedAt, order.CompletedAt,
		order.Subtotal, order.TaxAmount, order.DiscountAmount, order.ServiceCharge, order.TotalAmount,
		order.SaleID, order.ShiftID, order.CustomerNotes, order.KitchenNotes, order.InternalNotes,
		order.Metadata, order.UpdatedAt, order.UpdatedBy)
	return err
}

func (r *OrderRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE orders SET deleted_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE orders SET status = $3, updated_at = $4 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, time.Now())
	return err
}

// ORDER ITEM REPOSITORY

type OrderItemRepository struct {
	db *DB
}

func NewOrderItemRepository(db *DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.OrderItemFilters) ([]restaurant.OrderItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, product_id, product_variant_id, item_name, quantity,
		       unit_price, course_id, course_position, fire_time, kitchen_station_id, kitchen_ticket_id,
		       status, fired_at, acknowledged_at, started_preparing_at, ready_at, served_at, modifiers_total,
		       discount_amount, line_total, special_instructions, customer_notes, kitchen_notes, seat_number,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM order_items
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.OrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *filters.OrderID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.CourseID != nil {
		argCount++
		query += fmt.Sprintf(" AND course_id = $%d", argCount)
		args = append(args, *filters.CourseID)
	}

	if filters.StationID != nil {
		argCount++
		query += fmt.Sprintf(" AND kitchen_station_id = $%d", argCount)
		args = append(args, *filters.StationID)
	}

	query += " ORDER BY course_position ASC, created_at ASC"

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

	var items []restaurant.OrderItem
	for rows.Next() {
		var i restaurant.OrderItem
		err := rows.Scan(&i.ID, &i.OrganizationID, &i.OrderID, &i.ProductID, &i.ProductVariantID,
			&i.ItemName, &i.Quantity, &i.UnitPrice, &i.CourseID, &i.CoursePosition, &i.FireTime,
			&i.KitchenStationID, &i.KitchenTicketID, &i.Status, &i.FiredAt, &i.AcknowledgedAt,
			&i.StartedPreparingAt, &i.ReadyAt, &i.ServedAt, &i.ModifiersTotal, &i.DiscountAmount,
			&i.LineTotal, &i.SpecialInstructions, &i.CustomerNotes, &i.KitchenNotes, &i.SeatNumber,
			&i.Metadata, &i.CreatedAt, &i.UpdatedAt, &i.CreatedBy, &i.UpdatedBy)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, rows.Err()
}

func (r *OrderItemRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.OrderItemFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM order_items WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.OrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *filters.OrderID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *OrderItemRepository) Create(ctx context.Context, item *restaurant.OrderItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO order_items (id, organization_id, order_id, product_id, product_variant_id, item_name,
		    quantity, unit_price, course_id, course_position, fire_time, kitchen_station_id, kitchen_ticket_id,
		    status, fired_at, acknowledged_at, started_preparing_at, ready_at, served_at, modifiers_total,
		    discount_amount, line_total, special_instructions, customer_notes, kitchen_notes, seat_number,
		    metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
		    $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.ID, item.OrganizationID, item.OrderID, item.ProductID, item.ProductVariantID, item.ItemName,
		item.Quantity, item.UnitPrice, item.CourseID, item.CoursePosition, item.FireTime, item.KitchenStationID,
		item.KitchenTicketID, item.Status, item.FiredAt, item.AcknowledgedAt, item.StartedPreparingAt,
		item.ReadyAt, item.ServedAt, item.ModifiersTotal, item.DiscountAmount, item.LineTotal,
		item.SpecialInstructions, item.CustomerNotes, item.KitchenNotes, item.SeatNumber, item.Metadata,
		item.CreatedAt, item.UpdatedAt, item.CreatedBy)
	return err
}

func (r *OrderItemRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.OrderItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, product_id, product_variant_id, item_name, quantity,
		       unit_price, course_id, course_position, fire_time, kitchen_station_id, kitchen_ticket_id,
		       status, fired_at, acknowledged_at, started_preparing_at, ready_at, served_at, modifiers_total,
		       discount_amount, line_total, special_instructions, customer_notes, kitchen_notes, seat_number,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM order_items
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var i restaurant.OrderItem
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&i.ID, &i.OrganizationID, &i.OrderID, &i.ProductID, &i.ProductVariantID, &i.ItemName,
		&i.Quantity, &i.UnitPrice, &i.CourseID, &i.CoursePosition, &i.FireTime, &i.KitchenStationID,
		&i.KitchenTicketID, &i.Status, &i.FiredAt, &i.AcknowledgedAt, &i.StartedPreparingAt,
		&i.ReadyAt, &i.ServedAt, &i.ModifiersTotal, &i.DiscountAmount, &i.LineTotal,
		&i.SpecialInstructions, &i.CustomerNotes, &i.KitchenNotes, &i.SeatNumber, &i.Metadata,
		&i.CreatedAt, &i.UpdatedAt, &i.CreatedBy, &i.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *OrderItemRepository) Update(ctx context.Context, item *restaurant.OrderItem) error {
	if err := r.db.SetOrganizationContext(ctx, item.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE order_items SET
		    status = $3, fired_at = $4, acknowledged_at = $5, started_preparing_at = $6, ready_at = $7,
		    served_at = $8, kitchen_station_id = $9, kitchen_ticket_id = $10, fire_time = $11,
		    modifiers_total = $12, discount_amount = $13, line_total = $14, special_instructions = $15,
		    customer_notes = $16, kitchen_notes = $17, updated_at = $18, updated_by = $19
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		item.OrganizationID, item.ID, item.Status, item.FiredAt, item.AcknowledgedAt, item.StartedPreparingAt,
		item.ReadyAt, item.ServedAt, item.KitchenStationID, item.KitchenTicketID, item.FireTime,
		item.ModifiersTotal, item.DiscountAmount, item.LineTotal, item.SpecialInstructions,
		item.CustomerNotes, item.KitchenNotes, item.UpdatedAt, item.UpdatedBy)
	return err
}

func (r *OrderItemRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE order_items SET deleted_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *OrderItemRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE order_items SET status = $3, updated_at = $4 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, time.Now())
	return err
}

func (r *OrderItemRepository) ListByOrder(ctx context.Context, orgID uuid.UUID, orderID uuid.UUID) ([]restaurant.OrderItem, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_id, product_id, product_variant_id, item_name, quantity,
		       unit_price, course_id, course_position, fire_time, kitchen_station_id, kitchen_ticket_id,
		       status, fired_at, acknowledged_at, started_preparing_at, ready_at, served_at, modifiers_total,
		       discount_amount, line_total, special_instructions, customer_notes, kitchen_notes, seat_number,
		       metadata, created_at, updated_at, created_by, updated_by
		FROM order_items
		WHERE organization_id = $1 AND order_id = $2 AND deleted_at IS NULL
		ORDER BY course_position ASC, created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []restaurant.OrderItem
	for rows.Next() {
		var i restaurant.OrderItem
		err := rows.Scan(&i.ID, &i.OrganizationID, &i.OrderID, &i.ProductID, &i.ProductVariantID,
			&i.ItemName, &i.Quantity, &i.UnitPrice, &i.CourseID, &i.CoursePosition, &i.FireTime,
			&i.KitchenStationID, &i.KitchenTicketID, &i.Status, &i.FiredAt, &i.AcknowledgedAt,
			&i.StartedPreparingAt, &i.ReadyAt, &i.ServedAt, &i.ModifiersTotal, &i.DiscountAmount,
			&i.LineTotal, &i.SpecialInstructions, &i.CustomerNotes, &i.KitchenNotes, &i.SeatNumber,
			&i.Metadata, &i.CreatedAt, &i.UpdatedAt, &i.CreatedBy, &i.UpdatedBy)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, rows.Err()
}

// ORDER ITEM MODIFIER REPOSITORY

type OrderItemModifierRepository struct {
	db *DB
}

func NewOrderItemModifierRepository(db *DB) *OrderItemModifierRepository {
	return &OrderItemModifierRepository{db: db}
}

func (r *OrderItemModifierRepository) List(ctx context.Context, orgID uuid.UUID, orderItemID uuid.UUID) ([]restaurant.OrderItemModifier, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_item_id, modifier_id, modifier_group_id, modifier_name,
		       quantity, price_adjustment, display_order, metadata, created_at, updated_at, created_by, updated_by
		FROM order_item_modifiers
		WHERE organization_id = $1 AND order_item_id = $2 AND deleted_at IS NULL
		ORDER BY display_order ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orgID, orderItemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modifiers []restaurant.OrderItemModifier
	for rows.Next() {
		var m restaurant.OrderItemModifier
		err := rows.Scan(&m.ID, &m.OrganizationID, &m.OrderItemID, &m.ModifierID, &m.ModifierGroupID,
			&m.ModifierName, &m.Quantity, &m.PriceAdjustment, &m.DisplayOrder, &m.Metadata,
			&m.CreatedAt, &m.UpdatedAt, &m.CreatedBy, &m.UpdatedBy)
		if err != nil {
			return nil, err
		}
		modifiers = append(modifiers, m)
	}

	return modifiers, rows.Err()
}

func (r *OrderItemModifierRepository) Create(ctx context.Context, modifier *restaurant.OrderItemModifier) error {
	if err := r.db.SetOrganizationContext(ctx, modifier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO order_item_modifiers (id, organization_id, order_item_id, modifier_id, modifier_group_id,
		    modifier_name, quantity, price_adjustment, display_order, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		modifier.ID, modifier.OrganizationID, modifier.OrderItemID, modifier.ModifierID, modifier.ModifierGroupID,
		modifier.ModifierName, modifier.Quantity, modifier.PriceAdjustment, modifier.DisplayOrder,
		modifier.Metadata, modifier.CreatedAt, modifier.UpdatedAt, modifier.CreatedBy)
	return err
}

func (r *OrderItemModifierRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.OrderItemModifier, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, order_item_id, modifier_id, modifier_group_id, modifier_name,
		       quantity, price_adjustment, display_order, metadata, created_at, updated_at, created_by, updated_by
		FROM order_item_modifiers
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var m restaurant.OrderItemModifier
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&m.ID, &m.OrganizationID, &m.OrderItemID, &m.ModifierID, &m.ModifierGroupID, &m.ModifierName,
		&m.Quantity, &m.PriceAdjustment, &m.DisplayOrder, &m.Metadata, &m.CreatedAt, &m.UpdatedAt, &m.CreatedBy, &m.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *OrderItemModifierRepository) Update(ctx context.Context, modifier *restaurant.OrderItemModifier) error {
	if err := r.db.SetOrganizationContext(ctx, modifier.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE order_item_modifiers SET
		    modifier_name = $3, quantity = $4, price_adjustment = $5, display_order = $6,
		    metadata = $7, updated_at = $8, updated_by = $9
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		modifier.OrganizationID, modifier.ID, modifier.ModifierName, modifier.Quantity, modifier.PriceAdjustment,
		modifier.DisplayOrder, modifier.Metadata, modifier.UpdatedAt, modifier.UpdatedBy)
	return err
}

func (r *OrderItemModifierRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE order_item_modifiers SET deleted_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *OrderItemModifierRepository) DeleteByOrderItem(ctx context.Context, orgID uuid.UUID, orderItemID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE order_item_modifiers SET deleted_at = $3 WHERE organization_id = $1 AND order_item_id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, orderItemID, time.Now())
	return err
}

// KITCHEN TICKET REPOSITORY

type KitchenTicketRepository struct {
	db *DB
}

func NewKitchenTicketRepository(db *DB) *KitchenTicketRepository {
	return &KitchenTicketRepository{db: db}
}

func (r *KitchenTicketRepository) List(ctx context.Context, orgID uuid.UUID, filters restaurant.KitchenTicketFilters) ([]restaurant.KitchenTicket, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, ticket_number, display_sequence, order_id,
		       kitchen_station_id, course_id, ticket_type, priority, status, created_at, fired_at,
		       acknowledged_at, started_at, ready_at, bumped_at, completed_at, prep_time_minutes,
		       target_prep_time, table_number, order_type, covers, waiter_name, special_instructions,
		       kitchen_notes, display_config, metadata, updated_at, created_by, updated_by
		FROM kitchen_tickets
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
	}

	if filters.TicketType != nil {
		argCount++
		query += fmt.Sprintf(" AND ticket_type = $%d", argCount)
		args = append(args, *filters.TicketType)
	}

	if filters.StationID != nil {
		argCount++
		query += fmt.Sprintf(" AND kitchen_station_id = $%d", argCount)
		args = append(args, *filters.StationID)
	}

	if filters.OrderID != nil {
		argCount++
		query += fmt.Sprintf(" AND order_id = $%d", argCount)
		args = append(args, *filters.OrderID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, *filters.LocationID)
	}

	if filters.FromDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filters.FromDate)
	}

	if filters.ToDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filters.ToDate)
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

	var tickets []restaurant.KitchenTicket
	for rows.Next() {
		var t restaurant.KitchenTicket
		err := rows.Scan(&t.ID, &t.OrganizationID, &t.LocationID, &t.TicketNumber, &t.DisplaySequence,
			&t.OrderID, &t.KitchenStationID, &t.CourseID, &t.TicketType, &t.Priority, &t.Status,
			&t.CreatedAt, &t.FiredAt, &t.AcknowledgedAt, &t.StartedAt, &t.ReadyAt, &t.BumpedAt,
			&t.CompletedAt, &t.PrepTimeMinutes, &t.TargetPrepTime, &t.TableNumber, &t.OrderType,
			&t.Covers, &t.WaiterName, &t.SpecialInstructions, &t.KitchenNotes, &t.DisplayConfig,
			&t.Metadata, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}

func (r *KitchenTicketRepository) Count(ctx context.Context, orgID uuid.UUID, filters restaurant.KitchenTicketFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM kitchen_tickets WHERE organization_id = $1 AND deleted_at IS NULL"
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

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *KitchenTicketRepository) Create(ctx context.Context, ticket *restaurant.KitchenTicket) error {
	if err := r.db.SetOrganizationContext(ctx, ticket.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO kitchen_tickets (id, organization_id, location_id, ticket_number, display_sequence,
		    order_id, kitchen_station_id, course_id, ticket_type, priority, status, created_at, fired_at,
		    acknowledged_at, started_at, ready_at, bumped_at, completed_at, prep_time_minutes,
		    target_prep_time, table_number, order_type, covers, waiter_name, special_instructions,
		    kitchen_notes, display_config, metadata, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18,
		    $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		ticket.ID, ticket.OrganizationID, ticket.LocationID, ticket.TicketNumber, ticket.DisplaySequence,
		ticket.OrderID, ticket.KitchenStationID, ticket.CourseID, ticket.TicketType, ticket.Priority, ticket.Status,
		ticket.CreatedAt, ticket.FiredAt, ticket.AcknowledgedAt, ticket.StartedAt, ticket.ReadyAt, ticket.BumpedAt,
		ticket.CompletedAt, ticket.PrepTimeMinutes, ticket.TargetPrepTime, ticket.TableNumber, ticket.OrderType,
		ticket.Covers, ticket.WaiterName, ticket.SpecialInstructions, ticket.KitchenNotes, ticket.DisplayConfig,
		ticket.Metadata, ticket.UpdatedAt, ticket.CreatedBy)
	return err
}

func (r *KitchenTicketRepository) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*restaurant.KitchenTicket, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, ticket_number, display_sequence, order_id,
		       kitchen_station_id, course_id, ticket_type, priority, status, created_at, fired_at,
		       acknowledged_at, started_at, ready_at, bumped_at, completed_at, prep_time_minutes,
		       target_prep_time, table_number, order_type, covers, waiter_name, special_instructions,
		       kitchen_notes, display_config, metadata, updated_at, created_by, updated_by
		FROM kitchen_tickets
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	var t restaurant.KitchenTicket
	err := r.db.Pool.QueryRow(ctx, query, orgID, id).Scan(
		&t.ID, &t.OrganizationID, &t.LocationID, &t.TicketNumber, &t.DisplaySequence, &t.OrderID,
		&t.KitchenStationID, &t.CourseID, &t.TicketType, &t.Priority, &t.Status, &t.CreatedAt,
		&t.FiredAt, &t.AcknowledgedAt, &t.StartedAt, &t.ReadyAt, &t.BumpedAt, &t.CompletedAt,
		&t.PrepTimeMinutes, &t.TargetPrepTime, &t.TableNumber, &t.OrderType, &t.Covers, &t.WaiterName,
		&t.SpecialInstructions, &t.KitchenNotes, &t.DisplayConfig, &t.Metadata, &t.UpdatedAt,
		&t.CreatedBy, &t.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *KitchenTicketRepository) GetByNumber(ctx context.Context, orgID uuid.UUID, ticketNumber string) (*restaurant.KitchenTicket, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, ticket_number, display_sequence, order_id,
		       kitchen_station_id, course_id, ticket_type, priority, status, created_at, fired_at,
		       acknowledged_at, started_at, ready_at, bumped_at, completed_at, prep_time_minutes,
		       target_prep_time, table_number, order_type, covers, waiter_name, special_instructions,
		       kitchen_notes, display_config, metadata, updated_at, created_by, updated_by
		FROM kitchen_tickets
		WHERE organization_id = $1 AND ticket_number = $2 AND deleted_at IS NULL
	`

	var t restaurant.KitchenTicket
	err := r.db.Pool.QueryRow(ctx, query, orgID, ticketNumber).Scan(
		&t.ID, &t.OrganizationID, &t.LocationID, &t.TicketNumber, &t.DisplaySequence, &t.OrderID,
		&t.KitchenStationID, &t.CourseID, &t.TicketType, &t.Priority, &t.Status, &t.CreatedAt,
		&t.FiredAt, &t.AcknowledgedAt, &t.StartedAt, &t.ReadyAt, &t.BumpedAt, &t.CompletedAt,
		&t.PrepTimeMinutes, &t.TargetPrepTime, &t.TableNumber, &t.OrderType, &t.Covers, &t.WaiterName,
		&t.SpecialInstructions, &t.KitchenNotes, &t.DisplayConfig, &t.Metadata, &t.UpdatedAt,
		&t.CreatedBy, &t.UpdatedBy)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *KitchenTicketRepository) Update(ctx context.Context, ticket *restaurant.KitchenTicket) error {
	if err := r.db.SetOrganizationContext(ctx, ticket.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE kitchen_tickets SET
		    ticket_number = $3, display_sequence = $4, ticket_type = $5, priority = $6, status = $7,
		    fired_at = $8, acknowledged_at = $9, started_at = $10, ready_at = $11, bumped_at = $12,
		    completed_at = $13, prep_time_minutes = $14, target_prep_time = $15, table_number = $16,
		    order_type = $17, covers = $18, waiter_name = $19, special_instructions = $20,
		    kitchen_notes = $21, display_config = $22, metadata = $23, updated_at = $24, updated_by = $25
		WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		ticket.OrganizationID, ticket.ID, ticket.TicketNumber, ticket.DisplaySequence, ticket.TicketType,
		ticket.Priority, ticket.Status, ticket.FiredAt, ticket.AcknowledgedAt, ticket.StartedAt, ticket.ReadyAt,
		ticket.BumpedAt, ticket.CompletedAt, ticket.PrepTimeMinutes, ticket.TargetPrepTime, ticket.TableNumber,
		ticket.OrderType, ticket.Covers, ticket.WaiterName, ticket.SpecialInstructions, ticket.KitchenNotes,
		ticket.DisplayConfig, ticket.Metadata, ticket.UpdatedAt, ticket.UpdatedBy)
	return err
}

func (r *KitchenTicketRepository) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := "UPDATE kitchen_tickets SET deleted_at = $3 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL"
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, time.Now())
	return err
}

func (r *KitchenTicketRepository) UpdateStatus(ctx context.Context, orgID uuid.UUID, id uuid.UUID, status string) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `UPDATE kitchen_tickets SET status = $3, updated_at = $4 WHERE organization_id = $1 AND id = $2 AND deleted_at IS NULL`
	_, err := r.db.Pool.Exec(ctx, query, orgID, id, status, time.Now())
	return err
}

func (r *KitchenTicketRepository) ListByStation(ctx context.Context, orgID uuid.UUID, stationID uuid.UUID, statuses []string) ([]restaurant.KitchenTicket, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, ticket_number, display_sequence, order_id,
		       kitchen_station_id, course_id, ticket_type, priority, status, created_at, fired_at,
		       acknowledged_at, started_at, ready_at, bumped_at, completed_at, prep_time_minutes,
		       target_prep_time, table_number, order_type, covers, waiter_name, special_instructions,
		       kitchen_notes, display_config, metadata, updated_at, created_by, updated_by
		FROM kitchen_tickets
		WHERE organization_id = $1 AND kitchen_station_id = $2 AND deleted_at IS NULL
	`

	args := []interface{}{orgID, stationID}

	if len(statuses) > 0 {
		query += " AND status = ANY($3)"
		args = append(args, statuses)
	}

	query += " ORDER BY priority DESC, created_at ASC"

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []restaurant.KitchenTicket
	for rows.Next() {
		var t restaurant.KitchenTicket
		err := rows.Scan(&t.ID, &t.OrganizationID, &t.LocationID, &t.TicketNumber, &t.DisplaySequence,
			&t.OrderID, &t.KitchenStationID, &t.CourseID, &t.TicketType, &t.Priority, &t.Status,
			&t.CreatedAt, &t.FiredAt, &t.AcknowledgedAt, &t.StartedAt, &t.ReadyAt, &t.BumpedAt,
			&t.CompletedAt, &t.PrepTimeMinutes, &t.TargetPrepTime, &t.TableNumber, &t.OrderType,
			&t.Covers, &t.WaiterName, &t.SpecialInstructions, &t.KitchenNotes, &t.DisplayConfig,
			&t.Metadata, &t.UpdatedAt, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, rows.Err()
}
