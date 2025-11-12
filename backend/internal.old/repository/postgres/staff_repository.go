package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/staff"
)

type StaffRepository struct {
	db *DB
}

func NewStaffRepository(db *DB) *StaffRepository {
	return &StaffRepository{db: db}
}

// ==================== EmployeeSchedule Operations ====================

func (r *StaffRepository) CreateSchedule(ctx context.Context, schedule *staff.EmployeeSchedule) error {
	if err := r.db.SetOrganizationContext(ctx, schedule.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO employee_schedules (
			id, organization_id, location_id, employee_id, schedule_date,
			shift_type, position, scheduled_start_time, scheduled_end_time,
			break_duration_minutes, status, requires_approval, approved_by,
			approved_at, notes, cancellation_reason, metadata, created_at,
			updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		schedule.ID, schedule.OrganizationID, schedule.LocationID, schedule.EmployeeID,
		schedule.ScheduleDate, schedule.ShiftType, schedule.Position,
		schedule.ScheduledStartTime, schedule.ScheduledEndTime, schedule.BreakDurationMins,
		schedule.Status, schedule.RequiresApproval, schedule.ApprovedBy, schedule.ApprovedAt,
		schedule.Notes, schedule.CancellationReason, schedule.Metadata, schedule.CreatedAt,
		schedule.UpdatedAt, schedule.CreatedBy, schedule.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) (*staff.EmployeeSchedule, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, schedule_date,
			   shift_type, position, scheduled_start_time, scheduled_end_time,
			   break_duration_minutes, status, requires_approval, approved_by,
			   approved_at, notes, cancellation_reason, metadata, created_at,
			   updated_at, created_by, updated_by, deleted_at
		FROM employee_schedules
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var schedule staff.EmployeeSchedule
	err := r.db.Pool.QueryRow(ctx, query, scheduleID, orgID).Scan(
		&schedule.ID, &schedule.OrganizationID, &schedule.LocationID, &schedule.EmployeeID,
		&schedule.ScheduleDate, &schedule.ShiftType, &schedule.Position,
		&schedule.ScheduledStartTime, &schedule.ScheduledEndTime, &schedule.BreakDurationMins,
		&schedule.Status, &schedule.RequiresApproval, &schedule.ApprovedBy, &schedule.ApprovedAt,
		&schedule.Notes, &schedule.CancellationReason, &schedule.Metadata, &schedule.CreatedAt,
		&schedule.UpdatedAt, &schedule.CreatedBy, &schedule.UpdatedBy, &schedule.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &schedule, err
}

func (r *StaffRepository) ListSchedules(ctx context.Context, orgID uuid.UUID, filters staff.EmployeeScheduleFilters) ([]staff.EmployeeSchedule, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, schedule_date,
			   shift_type, position, scheduled_start_time, scheduled_end_time,
			   break_duration_minutes, status, requires_approval, approved_by,
			   approved_at, notes, cancellation_reason, metadata, created_at,
			   updated_at, created_by, updated_by, deleted_at
		FROM employee_schedules
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.ShiftType != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_type = $%d", argCount)
		args = append(args, filters.ShiftType)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND schedule_date >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND schedule_date <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	query += " ORDER BY schedule_date DESC, scheduled_start_time ASC"

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

	var schedules []staff.EmployeeSchedule
	for rows.Next() {
		var s staff.EmployeeSchedule
		if err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.LocationID, &s.EmployeeID,
			&s.ScheduleDate, &s.ShiftType, &s.Position,
			&s.ScheduledStartTime, &s.ScheduledEndTime, &s.BreakDurationMins,
			&s.Status, &s.RequiresApproval, &s.ApprovedBy, &s.ApprovedAt,
			&s.Notes, &s.CancellationReason, &s.Metadata, &s.CreatedAt,
			&s.UpdatedAt, &s.CreatedBy, &s.UpdatedBy, &s.DeletedAt,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	return schedules, rows.Err()
}

func (r *StaffRepository) CountSchedules(ctx context.Context, orgID uuid.UUID, filters staff.EmployeeScheduleFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM employee_schedules WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.ShiftType != nil {
		argCount++
		query += fmt.Sprintf(" AND shift_type = $%d", argCount)
		args = append(args, filters.ShiftType)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND schedule_date >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND schedule_date <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateSchedule(ctx context.Context, schedule *staff.EmployeeSchedule) error {
	if err := r.db.SetOrganizationContext(ctx, schedule.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE employee_schedules
		SET location_id = $2, schedule_date = $3, shift_type = $4, position = $5,
			scheduled_start_time = $6, scheduled_end_time = $7, break_duration_minutes = $8,
			status = $9, requires_approval = $10, approved_by = $11, approved_at = $12,
			notes = $13, cancellation_reason = $14, metadata = $15, updated_at = $16,
			updated_by = $17
		WHERE id = $1 AND organization_id = $18 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		schedule.ID, schedule.LocationID, schedule.ScheduleDate, schedule.ShiftType,
		schedule.Position, schedule.ScheduledStartTime, schedule.ScheduledEndTime,
		schedule.BreakDurationMins, schedule.Status, schedule.RequiresApproval,
		schedule.ApprovedBy, schedule.ApprovedAt, schedule.Notes, schedule.CancellationReason,
		schedule.Metadata, schedule.UpdatedAt, schedule.UpdatedBy, schedule.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteSchedule(ctx context.Context, orgID, scheduleID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE employee_schedules
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, scheduleID, orgID)
	return err
}

// ==================== TimeClockEntry Operations ====================

func (r *StaffRepository) CreateClockEntry(ctx context.Context, entry *staff.TimeClockEntry) error {
	if err := r.db.SetOrganizationContext(ctx, entry.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO time_clock_entries (
			id, organization_id, location_id, employee_id, schedule_id,
			entry_type, entry_timestamp, scheduled_timestamp, device_id,
			gps_location, ip_address, is_late, is_early, variance_minutes,
			requires_approval, approved_by, approved_at, is_manual_entry,
			correction_notes, photo_url, notes, metadata, created_at,
			updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		entry.ID, entry.OrganizationID, entry.LocationID, entry.EmployeeID,
		entry.ScheduleID, entry.EntryType, entry.EntryTimestamp,
		entry.ScheduledTimestamp, entry.DeviceID, entry.GPSLocation,
		entry.IPAddress, entry.IsLate, entry.IsEarly, entry.VarianceMinutes,
		entry.RequiresApproval, entry.ApprovedBy, entry.ApprovedAt,
		entry.IsManualEntry, entry.CorrectionNotes, entry.PhotoURL,
		entry.Notes, entry.Metadata, entry.CreatedAt, entry.UpdatedAt,
		entry.CreatedBy, entry.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetClockEntry(ctx context.Context, orgID, entryID uuid.UUID) (*staff.TimeClockEntry, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, schedule_id,
			   entry_type, entry_timestamp, scheduled_timestamp, device_id,
			   gps_location, ip_address, is_late, is_early, variance_minutes,
			   requires_approval, approved_by, approved_at, is_manual_entry,
			   correction_notes, photo_url, notes, metadata, created_at,
			   updated_at, created_by, updated_by, deleted_at
		FROM time_clock_entries
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var entry staff.TimeClockEntry
	err := r.db.Pool.QueryRow(ctx, query, entryID, orgID).Scan(
		&entry.ID, &entry.OrganizationID, &entry.LocationID, &entry.EmployeeID,
		&entry.ScheduleID, &entry.EntryType, &entry.EntryTimestamp,
		&entry.ScheduledTimestamp, &entry.DeviceID, &entry.GPSLocation,
		&entry.IPAddress, &entry.IsLate, &entry.IsEarly, &entry.VarianceMinutes,
		&entry.RequiresApproval, &entry.ApprovedBy, &entry.ApprovedAt,
		&entry.IsManualEntry, &entry.CorrectionNotes, &entry.PhotoURL,
		&entry.Notes, &entry.Metadata, &entry.CreatedAt, &entry.UpdatedAt,
		&entry.CreatedBy, &entry.UpdatedBy, &entry.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &entry, err
}

func (r *StaffRepository) ListClockEntries(ctx context.Context, orgID uuid.UUID, filters staff.TimeClockEntryFilters) ([]staff.TimeClockEntry, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, schedule_id,
			   entry_type, entry_timestamp, scheduled_timestamp, device_id,
			   gps_location, ip_address, is_late, is_early, variance_minutes,
			   requires_approval, approved_by, approved_at, is_manual_entry,
			   correction_notes, photo_url, notes, metadata, created_at,
			   updated_at, created_by, updated_by, deleted_at
		FROM time_clock_entries
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.EntryType != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_type = $%d", argCount)
		args = append(args, filters.EntryType)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_timestamp >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_timestamp <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	query += " ORDER BY entry_timestamp DESC"

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

	var entries []staff.TimeClockEntry
	for rows.Next() {
		var e staff.TimeClockEntry
		if err := rows.Scan(
			&e.ID, &e.OrganizationID, &e.LocationID, &e.EmployeeID,
			&e.ScheduleID, &e.EntryType, &e.EntryTimestamp,
			&e.ScheduledTimestamp, &e.DeviceID, &e.GPSLocation,
			&e.IPAddress, &e.IsLate, &e.IsEarly, &e.VarianceMinutes,
			&e.RequiresApproval, &e.ApprovedBy, &e.ApprovedAt,
			&e.IsManualEntry, &e.CorrectionNotes, &e.PhotoURL,
			&e.Notes, &e.Metadata, &e.CreatedAt, &e.UpdatedAt,
			&e.CreatedBy, &e.UpdatedBy, &e.DeletedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

func (r *StaffRepository) CountClockEntries(ctx context.Context, orgID uuid.UUID, filters staff.TimeClockEntryFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM time_clock_entries WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.EntryType != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_type = $%d", argCount)
		args = append(args, filters.EntryType)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_timestamp >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND entry_timestamp <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateClockEntry(ctx context.Context, entry *staff.TimeClockEntry) error {
	if err := r.db.SetOrganizationContext(ctx, entry.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE time_clock_entries
		SET location_id = $2, schedule_id = $3, entry_type = $4,
			entry_timestamp = $5, scheduled_timestamp = $6, device_id = $7,
			gps_location = $8, ip_address = $9, is_late = $10, is_early = $11,
			variance_minutes = $12, requires_approval = $13, approved_by = $14,
			approved_at = $15, is_manual_entry = $16, correction_notes = $17,
			photo_url = $18, notes = $19, metadata = $20, updated_at = $21,
			updated_by = $22
		WHERE id = $1 AND organization_id = $23 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		entry.ID, entry.LocationID, entry.ScheduleID, entry.EntryType,
		entry.EntryTimestamp, entry.ScheduledTimestamp, entry.DeviceID,
		entry.GPSLocation, entry.IPAddress, entry.IsLate, entry.IsEarly,
		entry.VarianceMinutes, entry.RequiresApproval, entry.ApprovedBy,
		entry.ApprovedAt, entry.IsManualEntry, entry.CorrectionNotes,
		entry.PhotoURL, entry.Notes, entry.Metadata, entry.UpdatedAt,
		entry.UpdatedBy, entry.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteClockEntry(ctx context.Context, orgID, entryID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE time_clock_entries
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, entryID, orgID)
	return err
}

// ==================== Device Operations ====================

func (r *StaffRepository) CreateDevice(ctx context.Context, device *staff.Device) error {
	if err := r.db.SetOrganizationContext(ctx, device.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO devices (
			id, organization_id, location_id, device_code, device_name, device_type,
			manufacturer, model, serial_number, mac_address, ip_address, device_config,
			screen_resolution, os_version, connection_type, connection_string, status,
			last_online_at, last_heartbeat_at, assigned_to_user_id, assigned_to_station_id,
			purchase_date, warranty_expiry_date, license_key, license_expiry_date,
			installation_notes, maintenance_notes, metadata, created_at, updated_at,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		device.ID, device.OrganizationID, device.LocationID, device.DeviceCode,
		device.DeviceName, device.DeviceType, device.Manufacturer, device.Model,
		device.SerialNumber, device.MACAddress, device.IPAddress, device.DeviceConfig,
		device.ScreenResolution, device.OSVersion, device.ConnectionType,
		device.ConnectionString, device.Status, device.LastOnlineAt, device.LastHeartbeatAt,
		device.AssignedToUserID, device.AssignedToStationID, device.PurchaseDate,
		device.WarrantyExpiryDate, device.LicenseKey, device.LicenseExpiryDate,
		device.InstallationNotes, device.MaintenanceNotes, device.Metadata,
		device.CreatedAt, device.UpdatedAt, device.CreatedBy, device.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetDevice(ctx context.Context, orgID, deviceID uuid.UUID) (*staff.Device, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, device_code, device_name, device_type,
			   manufacturer, model, serial_number, mac_address, ip_address, device_config,
			   screen_resolution, os_version, connection_type, connection_string, status,
			   last_online_at, last_heartbeat_at, assigned_to_user_id, assigned_to_station_id,
			   purchase_date, warranty_expiry_date, license_key, license_expiry_date,
			   installation_notes, maintenance_notes, metadata, created_at, updated_at,
			   created_by, updated_by, deleted_at
		FROM devices
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var device staff.Device
	err := r.db.Pool.QueryRow(ctx, query, deviceID, orgID).Scan(
		&device.ID, &device.OrganizationID, &device.LocationID, &device.DeviceCode,
		&device.DeviceName, &device.DeviceType, &device.Manufacturer, &device.Model,
		&device.SerialNumber, &device.MACAddress, &device.IPAddress, &device.DeviceConfig,
		&device.ScreenResolution, &device.OSVersion, &device.ConnectionType,
		&device.ConnectionString, &device.Status, &device.LastOnlineAt, &device.LastHeartbeatAt,
		&device.AssignedToUserID, &device.AssignedToStationID, &device.PurchaseDate,
		&device.WarrantyExpiryDate, &device.LicenseKey, &device.LicenseExpiryDate,
		&device.InstallationNotes, &device.MaintenanceNotes, &device.Metadata,
		&device.CreatedAt, &device.UpdatedAt, &device.CreatedBy, &device.UpdatedBy, &device.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &device, err
}

func (r *StaffRepository) ListDevices(ctx context.Context, orgID uuid.UUID, filters staff.DeviceFilters) ([]staff.Device, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, device_code, device_name, device_type,
			   manufacturer, model, serial_number, mac_address, ip_address, device_config,
			   screen_resolution, os_version, connection_type, connection_string, status,
			   last_online_at, last_heartbeat_at, assigned_to_user_id, assigned_to_station_id,
			   purchase_date, warranty_expiry_date, license_key, license_expiry_date,
			   installation_notes, maintenance_notes, metadata, created_at, updated_at,
			   created_by, updated_by, deleted_at
		FROM devices
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.DeviceType != nil {
		argCount++
		query += fmt.Sprintf(" AND device_type = $%d", argCount)
		args = append(args, filters.DeviceType)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.Search != nil {
		argCount++
		query += fmt.Sprintf(" AND (device_name ILIKE $%d OR device_code ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+*filters.Search+"%")
	}

	query += " ORDER BY device_name ASC"

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

	var devices []staff.Device
	for rows.Next() {
		var d staff.Device
		if err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.LocationID, &d.DeviceCode,
			&d.DeviceName, &d.DeviceType, &d.Manufacturer, &d.Model,
			&d.SerialNumber, &d.MACAddress, &d.IPAddress, &d.DeviceConfig,
			&d.ScreenResolution, &d.OSVersion, &d.ConnectionType,
			&d.ConnectionString, &d.Status, &d.LastOnlineAt, &d.LastHeartbeatAt,
			&d.AssignedToUserID, &d.AssignedToStationID, &d.PurchaseDate,
			&d.WarrantyExpiryDate, &d.LicenseKey, &d.LicenseExpiryDate,
			&d.InstallationNotes, &d.MaintenanceNotes, &d.Metadata,
			&d.CreatedAt, &d.UpdatedAt, &d.CreatedBy, &d.UpdatedBy, &d.DeletedAt,
		); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, rows.Err()
}

func (r *StaffRepository) CountDevices(ctx context.Context, orgID uuid.UUID, filters staff.DeviceFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM devices WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.DeviceType != nil {
		argCount++
		query += fmt.Sprintf(" AND device_type = $%d", argCount)
		args = append(args, filters.DeviceType)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateDevice(ctx context.Context, device *staff.Device) error {
	if err := r.db.SetOrganizationContext(ctx, device.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE devices
		SET location_id = $2, device_code = $3, device_name = $4, device_type = $5,
			manufacturer = $6, model = $7, serial_number = $8, mac_address = $9,
			ip_address = $10, device_config = $11, screen_resolution = $12,
			os_version = $13, connection_type = $14, connection_string = $15,
			status = $16, last_online_at = $17, last_heartbeat_at = $18,
			assigned_to_user_id = $19, assigned_to_station_id = $20,
			purchase_date = $21, warranty_expiry_date = $22, license_key = $23,
			license_expiry_date = $24, installation_notes = $25, maintenance_notes = $26,
			metadata = $27, updated_at = $28, updated_by = $29
		WHERE id = $1 AND organization_id = $30 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		device.ID, device.LocationID, device.DeviceCode, device.DeviceName,
		device.DeviceType, device.Manufacturer, device.Model, device.SerialNumber,
		device.MACAddress, device.IPAddress, device.DeviceConfig, device.ScreenResolution,
		device.OSVersion, device.ConnectionType, device.ConnectionString, device.Status,
		device.LastOnlineAt, device.LastHeartbeatAt, device.AssignedToUserID,
		device.AssignedToStationID, device.PurchaseDate, device.WarrantyExpiryDate,
		device.LicenseKey, device.LicenseExpiryDate, device.InstallationNotes,
		device.MaintenanceNotes, device.Metadata, device.UpdatedAt, device.UpdatedBy,
		device.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteDevice(ctx context.Context, orgID, deviceID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE devices
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, deviceID, orgID)
	return err
}

func (r *StaffRepository) UpdateDeviceHeartbeat(ctx context.Context, deviceID uuid.UUID) error {
	query := `
		UPDATE devices
		SET last_heartbeat_at = NOW(), last_online_at = NOW(), status = 'active', updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, deviceID)
	return err
}

// ==================== PrinterConfiguration Operations ====================

func (r *StaffRepository) CreatePrinterConfig(ctx context.Context, config *staff.PrinterConfiguration) error {
	if err := r.db.SetOrganizationContext(ctx, config.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO printer_configurations (
			id, organization_id, location_id, printer_device_id, document_type,
			filter_order_type, filter_kitchen_station_id, filter_product_category_id,
			filter_course_id, number_of_copies, auto_print, print_priority,
			template_config, paper_size, print_orientation, is_active, notes,
			metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		config.ID, config.OrganizationID, config.LocationID, config.PrinterDeviceID,
		config.DocumentType, config.FilterOrderType, config.FilterKitchenStationID,
		config.FilterProductCategoryID, config.FilterCourseID, config.NumberOfCopies,
		config.AutoPrint, config.PrintPriority, config.TemplateConfig, config.PaperSize,
		config.PrintOrientation, config.IsActive, config.Notes, config.Metadata,
		config.CreatedAt, config.UpdatedAt, config.CreatedBy, config.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetPrinterConfig(ctx context.Context, orgID, configID uuid.UUID) (*staff.PrinterConfiguration, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, printer_device_id, document_type,
			   filter_order_type, filter_kitchen_station_id, filter_product_category_id,
			   filter_course_id, number_of_copies, auto_print, print_priority,
			   template_config, paper_size, print_orientation, is_active, notes,
			   metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM printer_configurations
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var config staff.PrinterConfiguration
	err := r.db.Pool.QueryRow(ctx, query, configID, orgID).Scan(
		&config.ID, &config.OrganizationID, &config.LocationID, &config.PrinterDeviceID,
		&config.DocumentType, &config.FilterOrderType, &config.FilterKitchenStationID,
		&config.FilterProductCategoryID, &config.FilterCourseID, &config.NumberOfCopies,
		&config.AutoPrint, &config.PrintPriority, &config.TemplateConfig, &config.PaperSize,
		&config.PrintOrientation, &config.IsActive, &config.Notes, &config.Metadata,
		&config.CreatedAt, &config.UpdatedAt, &config.CreatedBy, &config.UpdatedBy, &config.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &config, err
}

func (r *StaffRepository) ListPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters staff.PrinterConfigurationFilters) ([]staff.PrinterConfiguration, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, printer_device_id, document_type,
			   filter_order_type, filter_kitchen_station_id, filter_product_category_id,
			   filter_course_id, number_of_copies, auto_print, print_priority,
			   template_config, paper_size, print_orientation, is_active, notes,
			   metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM printer_configurations
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.DocumentType != nil {
		argCount++
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, filters.DocumentType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	query += " ORDER BY document_type, print_priority DESC"

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

	var configs []staff.PrinterConfiguration
	for rows.Next() {
		var c staff.PrinterConfiguration
		if err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.LocationID, &c.PrinterDeviceID,
			&c.DocumentType, &c.FilterOrderType, &c.FilterKitchenStationID,
			&c.FilterProductCategoryID, &c.FilterCourseID, &c.NumberOfCopies,
			&c.AutoPrint, &c.PrintPriority, &c.TemplateConfig, &c.PaperSize,
			&c.PrintOrientation, &c.IsActive, &c.Notes, &c.Metadata,
			&c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy, &c.DeletedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}

	return configs, rows.Err()
}

func (r *StaffRepository) CountPrinterConfigs(ctx context.Context, orgID uuid.UUID, filters staff.PrinterConfigurationFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM printer_configurations WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.DocumentType != nil {
		argCount++
		query += fmt.Sprintf(" AND document_type = $%d", argCount)
		args = append(args, filters.DocumentType)
	}

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdatePrinterConfig(ctx context.Context, config *staff.PrinterConfiguration) error {
	if err := r.db.SetOrganizationContext(ctx, config.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE printer_configurations
		SET location_id = $2, printer_device_id = $3, document_type = $4,
			filter_order_type = $5, filter_kitchen_station_id = $6,
			filter_product_category_id = $7, filter_course_id = $8,
			number_of_copies = $9, auto_print = $10, print_priority = $11,
			template_config = $12, paper_size = $13, print_orientation = $14,
			is_active = $15, notes = $16, metadata = $17, updated_at = $18, updated_by = $19
		WHERE id = $1 AND organization_id = $20 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		config.ID, config.LocationID, config.PrinterDeviceID, config.DocumentType,
		config.FilterOrderType, config.FilterKitchenStationID, config.FilterProductCategoryID,
		config.FilterCourseID, config.NumberOfCopies, config.AutoPrint, config.PrintPriority,
		config.TemplateConfig, config.PaperSize, config.PrintOrientation, config.IsActive,
		config.Notes, config.Metadata, config.UpdatedAt, config.UpdatedBy, config.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeletePrinterConfig(ctx context.Context, orgID, configID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE printer_configurations
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, configID, orgID)
	return err
}

// ==================== TipPool Operations ====================

func (r *StaffRepository) CreateTipPool(ctx context.Context, pool *staff.TipPool) error {
	if err := r.db.SetOrganizationContext(ctx, pool.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO tip_pools (
			id, organization_id, location_id, pool_name, pool_type, description,
			distribution_method, distribution_config, eligible_positions, is_active,
			metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		pool.ID, pool.OrganizationID, pool.LocationID, pool.PoolName,
		pool.PoolType, pool.Description, pool.DistributionMethod,
		pool.DistributionConfig, pool.EligiblePositions, pool.IsActive,
		pool.Metadata, pool.CreatedAt, pool.UpdatedAt, pool.CreatedBy, pool.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetTipPool(ctx context.Context, orgID, poolID uuid.UUID) (*staff.TipPool, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, pool_name, pool_type, description,
			   distribution_method, distribution_config, eligible_positions, is_active,
			   metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM tip_pools
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var pool staff.TipPool
	err := r.db.Pool.QueryRow(ctx, query, poolID, orgID).Scan(
		&pool.ID, &pool.OrganizationID, &pool.LocationID, &pool.PoolName,
		&pool.PoolType, &pool.Description, &pool.DistributionMethod,
		&pool.DistributionConfig, &pool.EligiblePositions, &pool.IsActive,
		&pool.Metadata, &pool.CreatedAt, &pool.UpdatedAt, &pool.CreatedBy, &pool.UpdatedBy,
		&pool.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &pool, err
}

func (r *StaffRepository) ListTipPools(ctx context.Context, orgID uuid.UUID, filters staff.TipPoolFilters) ([]staff.TipPool, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, pool_name, pool_type, description,
			   distribution_method, distribution_config, eligible_positions, is_active,
			   metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM tip_pools
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	query += " ORDER BY pool_name ASC"

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

	var pools []staff.TipPool
	for rows.Next() {
		var p staff.TipPool
		if err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.LocationID, &p.PoolName,
			&p.PoolType, &p.Description, &p.DistributionMethod,
			&p.DistributionConfig, &p.EligiblePositions, &p.IsActive,
			&p.Metadata, &p.CreatedAt, &p.UpdatedAt, &p.CreatedBy, &p.UpdatedBy,
			&p.DeletedAt,
		); err != nil {
			return nil, err
		}
		pools = append(pools, p)
	}

	return pools, rows.Err()
}

func (r *StaffRepository) CountTipPools(ctx context.Context, orgID uuid.UUID, filters staff.TipPoolFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM tip_pools WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, filters.IsActive)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateTipPool(ctx context.Context, pool *staff.TipPool) error {
	if err := r.db.SetOrganizationContext(ctx, pool.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE tip_pools
		SET location_id = $2, pool_name = $3, pool_type = $4, description = $5,
			distribution_method = $6, distribution_config = $7, eligible_positions = $8,
			is_active = $9, metadata = $10, updated_at = $11, updated_by = $12
		WHERE id = $1 AND organization_id = $13 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		pool.ID, pool.LocationID, pool.PoolName, pool.PoolType,
		pool.Description, pool.DistributionMethod, pool.DistributionConfig,
		pool.EligiblePositions, pool.IsActive, pool.Metadata, pool.UpdatedAt,
		pool.UpdatedBy, pool.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteTipPool(ctx context.Context, orgID, poolID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE tip_pools
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, poolID, orgID)
	return err
}

// ==================== TipDistribution Operations ====================

func (r *StaffRepository) CreateTipDistribution(ctx context.Context, dist *staff.TipDistribution) error {
	if err := r.db.SetOrganizationContext(ctx, dist.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO tip_distributions (
			id, organization_id, location_id, tip_pool_id, distribution_date,
			period_start, period_end, shift_id, employee_id, source_type,
			source_sale_id, source_order_id, tip_amount, distribution_amount,
			distribution_percentage, payment_status, payment_method, paid_at,
			paid_by, notes, metadata, created_at, updated_at, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		dist.ID, dist.OrganizationID, dist.LocationID, dist.TipPoolID,
		dist.DistributionDate, dist.PeriodStart, dist.PeriodEnd, dist.ShiftID,
		dist.EmployeeID, dist.SourceType, dist.SourceSaleID, dist.SourceOrderID,
		dist.TipAmount, dist.DistributionAmount, dist.DistributionPercentage,
		dist.PaymentStatus, dist.PaymentMethod, dist.PaidAt, dist.PaidBy,
		dist.Notes, dist.Metadata, dist.CreatedAt, dist.UpdatedAt,
		dist.CreatedBy, dist.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetTipDistribution(ctx context.Context, orgID, distID uuid.UUID) (*staff.TipDistribution, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, tip_pool_id, distribution_date,
			   period_start, period_end, shift_id, employee_id, source_type,
			   source_sale_id, source_order_id, tip_amount, distribution_amount,
			   distribution_percentage, payment_status, payment_method, paid_at,
			   paid_by, notes, metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM tip_distributions
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var dist staff.TipDistribution
	err := r.db.Pool.QueryRow(ctx, query, distID, orgID).Scan(
		&dist.ID, &dist.OrganizationID, &dist.LocationID, &dist.TipPoolID,
		&dist.DistributionDate, &dist.PeriodStart, &dist.PeriodEnd, &dist.ShiftID,
		&dist.EmployeeID, &dist.SourceType, &dist.SourceSaleID, &dist.SourceOrderID,
		&dist.TipAmount, &dist.DistributionAmount, &dist.DistributionPercentage,
		&dist.PaymentStatus, &dist.PaymentMethod, &dist.PaidAt, &dist.PaidBy,
		&dist.Notes, &dist.Metadata, &dist.CreatedAt, &dist.UpdatedAt,
		&dist.CreatedBy, &dist.UpdatedBy, &dist.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &dist, err
}

func (r *StaffRepository) ListTipDistributions(ctx context.Context, orgID uuid.UUID, filters staff.TipDistributionFilters) ([]staff.TipDistribution, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, tip_pool_id, distribution_date,
			   period_start, period_end, shift_id, employee_id, source_type,
			   source_sale_id, source_order_id, tip_amount, distribution_amount,
			   distribution_percentage, payment_status, payment_method, paid_at,
			   paid_by, notes, metadata, created_at, updated_at, created_by, updated_by, deleted_at
		FROM tip_distributions
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, filters.PaymentStatus)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.DistributionDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND distribution_date >= $%d", argCount)
		args = append(args, filters.DistributionDateFrom)
	}

	if filters.DistributionDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND distribution_date <= $%d", argCount)
		args = append(args, filters.DistributionDateTo)
	}

	query += " ORDER BY distribution_date DESC"

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

	var dists []staff.TipDistribution
	for rows.Next() {
		var d staff.TipDistribution
		if err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.LocationID, &d.TipPoolID,
			&d.DistributionDate, &d.PeriodStart, &d.PeriodEnd, &d.ShiftID,
			&d.EmployeeID, &d.SourceType, &d.SourceSaleID, &d.SourceOrderID,
			&d.TipAmount, &d.DistributionAmount, &d.DistributionPercentage,
			&d.PaymentStatus, &d.PaymentMethod, &d.PaidAt, &d.PaidBy,
			&d.Notes, &d.Metadata, &d.CreatedAt, &d.UpdatedAt,
			&d.CreatedBy, &d.UpdatedBy, &d.DeletedAt,
		); err != nil {
			return nil, err
		}
		dists = append(dists, d)
	}

	return dists, rows.Err()
}

func (r *StaffRepository) CountTipDistributions(ctx context.Context, orgID uuid.UUID, filters staff.TipDistributionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM tip_distributions WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.PaymentStatus != nil {
		argCount++
		query += fmt.Sprintf(" AND payment_status = $%d", argCount)
		args = append(args, filters.PaymentStatus)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.DistributionDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND distribution_date >= $%d", argCount)
		args = append(args, filters.DistributionDateFrom)
	}

	if filters.DistributionDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND distribution_date <= $%d", argCount)
		args = append(args, filters.DistributionDateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateTipDistribution(ctx context.Context, dist *staff.TipDistribution) error {
	if err := r.db.SetOrganizationContext(ctx, dist.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE tip_distributions
		SET location_id = $2, tip_pool_id = $3, distribution_date = $4,
			period_start = $5, period_end = $6, shift_id = $7, source_type = $8,
			source_sale_id = $9, source_order_id = $10, tip_amount = $11,
			distribution_amount = $12, distribution_percentage = $13,
			payment_status = $14, payment_method = $15, paid_at = $16,
			paid_by = $17, notes = $18, metadata = $19, updated_at = $20, updated_by = $21
		WHERE id = $1 AND organization_id = $22 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		dist.ID, dist.LocationID, dist.TipPoolID, dist.DistributionDate,
		dist.PeriodStart, dist.PeriodEnd, dist.ShiftID, dist.SourceType,
		dist.SourceSaleID, dist.SourceOrderID, dist.TipAmount,
		dist.DistributionAmount, dist.DistributionPercentage,
		dist.PaymentStatus, dist.PaymentMethod, dist.PaidAt,
		dist.PaidBy, dist.Notes, dist.Metadata, dist.UpdatedAt, dist.UpdatedBy,
		dist.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteTipDistribution(ctx context.Context, orgID, distID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE tip_distributions
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, distID, orgID)
	return err
}

// ==================== StaffCommission Operations ====================

func (r *StaffRepository) CreateCommission(ctx context.Context, comm *staff.StaffCommission) error {
	if err := r.db.SetOrganizationContext(ctx, comm.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO staff_commissions (
			id, organization_id, location_id, employee_id, commission_date,
			period_start, period_end, source_type, source_sale_id, source_order_id,
			commission_type, commission_rate, sales_amount, commission_amount,
			status, approved_by, approved_at, payment_date, payment_method,
			paid_by, notes, calculation_notes, metadata, created_at, updated_at,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		comm.ID, comm.OrganizationID, comm.LocationID, comm.EmployeeID,
		comm.CommissionDate, comm.PeriodStart, comm.PeriodEnd, comm.SourceType,
		comm.SourceSaleID, comm.SourceOrderID, comm.CommissionType,
		comm.CommissionRate, comm.SalesAmount, comm.CommissionAmount,
		comm.Status, comm.ApprovedBy, comm.ApprovedAt, comm.PaymentDate,
		comm.PaymentMethod, comm.PaidBy, comm.Notes, comm.CalculationNotes,
		comm.Metadata, comm.CreatedAt, comm.UpdatedAt, comm.CreatedBy, comm.UpdatedBy,
	)

	return err
}

func (r *StaffRepository) GetCommission(ctx context.Context, orgID, commID uuid.UUID) (*staff.StaffCommission, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, commission_date,
			   period_start, period_end, source_type, source_sale_id, source_order_id,
			   commission_type, commission_rate, sales_amount, commission_amount,
			   status, approved_by, approved_at, payment_date, payment_method,
			   paid_by, notes, calculation_notes, metadata, created_at, updated_at,
			   created_by, updated_by, deleted_at
		FROM staff_commissions
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var comm staff.StaffCommission
	err := r.db.Pool.QueryRow(ctx, query, commID, orgID).Scan(
		&comm.ID, &comm.OrganizationID, &comm.LocationID, &comm.EmployeeID,
		&comm.CommissionDate, &comm.PeriodStart, &comm.PeriodEnd, &comm.SourceType,
		&comm.SourceSaleID, &comm.SourceOrderID, &comm.CommissionType,
		&comm.CommissionRate, &comm.SalesAmount, &comm.CommissionAmount,
		&comm.Status, &comm.ApprovedBy, &comm.ApprovedAt, &comm.PaymentDate,
		&comm.PaymentMethod, &comm.PaidBy, &comm.Notes, &comm.CalculationNotes,
		&comm.Metadata, &comm.CreatedAt, &comm.UpdatedAt, &comm.CreatedBy, &comm.UpdatedBy,
		&comm.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &comm, err
}

func (r *StaffRepository) ListCommissions(ctx context.Context, orgID uuid.UUID, filters staff.StaffCommissionFilters) ([]staff.StaffCommission, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, employee_id, commission_date,
			   period_start, period_end, source_type, source_sale_id, source_order_id,
			   commission_type, commission_rate, sales_amount, commission_amount,
			   status, approved_by, approved_at, payment_date, payment_method,
			   paid_by, notes, calculation_notes, metadata, created_at, updated_at,
			   created_by, updated_by, deleted_at
		FROM staff_commissions
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.CommissionDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND commission_date >= $%d", argCount)
		args = append(args, filters.CommissionDateFrom)
	}

	if filters.CommissionDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND commission_date <= $%d", argCount)
		args = append(args, filters.CommissionDateTo)
	}

	query += " ORDER BY commission_date DESC"

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

	var comms []staff.StaffCommission
	for rows.Next() {
		var c staff.StaffCommission
		if err := rows.Scan(
			&c.ID, &c.OrganizationID, &c.LocationID, &c.EmployeeID,
			&c.CommissionDate, &c.PeriodStart, &c.PeriodEnd, &c.SourceType,
			&c.SourceSaleID, &c.SourceOrderID, &c.CommissionType,
			&c.CommissionRate, &c.SalesAmount, &c.CommissionAmount,
			&c.Status, &c.ApprovedBy, &c.ApprovedAt, &c.PaymentDate,
			&c.PaymentMethod, &c.PaidBy, &c.Notes, &c.CalculationNotes,
			&c.Metadata, &c.CreatedAt, &c.UpdatedAt, &c.CreatedBy, &c.UpdatedBy,
			&c.DeletedAt,
		); err != nil {
			return nil, err
		}
		comms = append(comms, c)
	}

	return comms, rows.Err()
}

func (r *StaffRepository) CountCommissions(ctx context.Context, orgID uuid.UUID, filters staff.StaffCommissionFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM staff_commissions WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.EmployeeID != nil {
		argCount++
		query += fmt.Sprintf(" AND employee_id = $%d", argCount)
		args = append(args, filters.EmployeeID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.CommissionDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND commission_date >= $%d", argCount)
		args = append(args, filters.CommissionDateFrom)
	}

	if filters.CommissionDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND commission_date <= $%d", argCount)
		args = append(args, filters.CommissionDateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateCommission(ctx context.Context, comm *staff.StaffCommission) error {
	if err := r.db.SetOrganizationContext(ctx, comm.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE staff_commissions
		SET location_id = $2, commission_date = $3, period_start = $4,
			period_end = $5, source_type = $6, source_sale_id = $7,
			source_order_id = $8, commission_type = $9, commission_rate = $10,
			sales_amount = $11, commission_amount = $12, status = $13,
			approved_by = $14, approved_at = $15, payment_date = $16,
			payment_method = $17, paid_by = $18, notes = $19,
			calculation_notes = $20, metadata = $21, updated_at = $22, updated_by = $23
		WHERE id = $1 AND organization_id = $24 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		comm.ID, comm.LocationID, comm.CommissionDate, comm.PeriodStart,
		comm.PeriodEnd, comm.SourceType, comm.SourceSaleID,
		comm.SourceOrderID, comm.CommissionType, comm.CommissionRate,
		comm.SalesAmount, comm.CommissionAmount, comm.Status,
		comm.ApprovedBy, comm.ApprovedAt, comm.PaymentDate,
		comm.PaymentMethod, comm.PaidBy, comm.Notes,
		comm.CalculationNotes, comm.Metadata, comm.UpdatedAt, comm.UpdatedBy,
		comm.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteCommission(ctx context.Context, orgID, commID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE staff_commissions
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, commID, orgID)
	return err
}

// ==================== Shift Operations ====================

func (r *StaffRepository) CreateShift(ctx context.Context, shift *staff.Shift) error {
	if err := r.db.SetOrganizationContext(ctx, shift.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO shifts (
			id, organization_id, location_id, user_id, shift_number, start_time,
			end_time, status, opening_cash, opening_notes, expected_cash,
			actual_cash, cash_difference, closing_notes, total_sales,
			total_transactions, total_refunds, total_discounts, payment_breakdown,
			notes, metadata, created_at, updated_at, closed_by, closed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		shift.ID, shift.OrganizationID, shift.LocationID, shift.UserID,
		shift.ShiftNumber, shift.StartTime, shift.EndTime, shift.Status,
		shift.OpeningCash, shift.OpeningNotes, shift.ExpectedCash,
		shift.ActualCash, shift.CashDifference, shift.ClosingNotes,
		shift.TotalSales, shift.TotalTransactions, shift.TotalRefunds,
		shift.TotalDiscounts, shift.PaymentBreakdown, shift.Notes,
		shift.Metadata, shift.CreatedAt, shift.UpdatedAt, shift.ClosedBy,
		shift.ClosedAt,
	)

	return err
}

func (r *StaffRepository) GetShift(ctx context.Context, orgID, shiftID uuid.UUID) (*staff.Shift, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, user_id, shift_number, start_time,
			   end_time, status, opening_cash, opening_notes, expected_cash,
			   actual_cash, cash_difference, closing_notes, total_sales,
			   total_transactions, total_refunds, total_discounts, payment_breakdown,
			   notes, metadata, created_at, updated_at, closed_by, closed_at, deleted_at
		FROM shifts
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var shift staff.Shift
	err := r.db.Pool.QueryRow(ctx, query, shiftID, orgID).Scan(
		&shift.ID, &shift.OrganizationID, &shift.LocationID, &shift.UserID,
		&shift.ShiftNumber, &shift.StartTime, &shift.EndTime, &shift.Status,
		&shift.OpeningCash, &shift.OpeningNotes, &shift.ExpectedCash,
		&shift.ActualCash, &shift.CashDifference, &shift.ClosingNotes,
		&shift.TotalSales, &shift.TotalTransactions, &shift.TotalRefunds,
		&shift.TotalDiscounts, &shift.PaymentBreakdown, &shift.Notes,
		&shift.Metadata, &shift.CreatedAt, &shift.UpdatedAt, &shift.ClosedBy,
		&shift.ClosedAt, &shift.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &shift, err
}

func (r *StaffRepository) ListShifts(ctx context.Context, orgID uuid.UUID, filters staff.ShiftFilters) ([]staff.Shift, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, user_id, shift_number, start_time,
			   end_time, status, opening_cash, opening_notes, expected_cash,
			   actual_cash, cash_difference, closing_notes, total_sales,
			   total_transactions, total_refunds, total_discounts, payment_breakdown,
			   notes, metadata, created_at, updated_at, closed_by, closed_at, deleted_at
		FROM shifts
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, filters.UserID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(start_time) >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(start_time) <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	query += " ORDER BY start_time DESC"

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

	var shifts []staff.Shift
	for rows.Next() {
		var s staff.Shift
		if err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.LocationID, &s.UserID,
			&s.ShiftNumber, &s.StartTime, &s.EndTime, &s.Status,
			&s.OpeningCash, &s.OpeningNotes, &s.ExpectedCash,
			&s.ActualCash, &s.CashDifference, &s.ClosingNotes,
			&s.TotalSales, &s.TotalTransactions, &s.TotalRefunds,
			&s.TotalDiscounts, &s.PaymentBreakdown, &s.Notes,
			&s.Metadata, &s.CreatedAt, &s.UpdatedAt, &s.ClosedBy,
			&s.ClosedAt, &s.DeletedAt,
		); err != nil {
			return nil, err
		}
		shifts = append(shifts, s)
	}

	return shifts, rows.Err()
}

func (r *StaffRepository) CountShifts(ctx context.Context, orgID uuid.UUID, filters staff.ShiftFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM shifts WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, filters.UserID)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(start_time) >= $%d", argCount)
		args = append(args, filters.DateFrom)
	}

	if filters.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND DATE(start_time) <= $%d", argCount)
		args = append(args, filters.DateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateShift(ctx context.Context, shift *staff.Shift) error {
	if err := r.db.SetOrganizationContext(ctx, shift.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE shifts
		SET location_id = $2, shift_number = $3, start_time = $4,
			end_time = $5, status = $6, opening_cash = $7, opening_notes = $8,
			expected_cash = $9, actual_cash = $10, cash_difference = $11,
			closing_notes = $12, total_sales = $13, total_transactions = $14,
			total_refunds = $15, total_discounts = $16, payment_breakdown = $17,
			notes = $18, metadata = $19, updated_at = $20, closed_by = $21, closed_at = $22
		WHERE id = $1 AND organization_id = $23 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		shift.ID, shift.LocationID, shift.ShiftNumber, shift.StartTime,
		shift.EndTime, shift.Status, shift.OpeningCash, shift.OpeningNotes,
		shift.ExpectedCash, shift.ActualCash, shift.CashDifference,
		shift.ClosingNotes, shift.TotalSales, shift.TotalTransactions,
		shift.TotalRefunds, shift.TotalDiscounts, shift.PaymentBreakdown,
		shift.Notes, shift.Metadata, shift.UpdatedAt, shift.ClosedBy, shift.ClosedAt,
		shift.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteShift(ctx context.Context, orgID, shiftID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE shifts
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, shiftID, orgID)
	return err
}

// ==================== Expense Operations ====================

func (r *StaffRepository) CreateExpense(ctx context.Context, expense *staff.Expense) error {
	if err := r.db.SetOrganizationContext(ctx, expense.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		INSERT INTO expenses (
			id, organization_id, location_id, expense_number, expense_date,
			category, subcategory, payee_name, payment_method, amount,
			tax_amount, total_amount, currency, status, reference_number,
			purchase_order_id, receipt_url, attachment_urls, description,
			notes, metadata, created_at, updated_at, created_by, updated_by,
			approved_by, approved_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27
		)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		expense.ID, expense.OrganizationID, expense.LocationID, expense.ExpenseNumber,
		expense.ExpenseDate, expense.Category, expense.Subcategory, expense.PayeeName,
		expense.PaymentMethod, expense.Amount, expense.TaxAmount, expense.TotalAmount,
		expense.Currency, expense.Status, expense.ReferenceNumber,
		expense.PurchaseOrderID, expense.ReceiptURL, expense.AttachmentURLs,
		expense.Description, expense.Notes, expense.Metadata, expense.CreatedAt,
		expense.UpdatedAt, expense.CreatedBy, expense.UpdatedBy, expense.ApprovedBy,
		expense.ApprovedAt,
	)

	return err
}

func (r *StaffRepository) GetExpense(ctx context.Context, orgID, expenseID uuid.UUID) (*staff.Expense, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, expense_number, expense_date,
			   category, subcategory, payee_name, payment_method, amount,
			   tax_amount, total_amount, currency, status, reference_number,
			   purchase_order_id, receipt_url, attachment_urls, description,
			   notes, metadata, created_at, updated_at, created_by, updated_by,
			   approved_by, approved_at, deleted_at
		FROM expenses
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	var expense staff.Expense
	err := r.db.Pool.QueryRow(ctx, query, expenseID, orgID).Scan(
		&expense.ID, &expense.OrganizationID, &expense.LocationID, &expense.ExpenseNumber,
		&expense.ExpenseDate, &expense.Category, &expense.Subcategory, &expense.PayeeName,
		&expense.PaymentMethod, &expense.Amount, &expense.TaxAmount, &expense.TotalAmount,
		&expense.Currency, &expense.Status, &expense.ReferenceNumber,
		&expense.PurchaseOrderID, &expense.ReceiptURL, &expense.AttachmentURLs,
		&expense.Description, &expense.Notes, &expense.Metadata, &expense.CreatedAt,
		&expense.UpdatedAt, &expense.CreatedBy, &expense.UpdatedBy, &expense.ApprovedBy,
		&expense.ApprovedAt, &expense.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &expense, err
}

func (r *StaffRepository) ListExpenses(ctx context.Context, orgID uuid.UUID, filters staff.ExpenseFilters) ([]staff.Expense, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return nil, err
	}

	query := `
		SELECT id, organization_id, location_id, expense_number, expense_date,
			   category, subcategory, payee_name, payment_method, amount,
			   tax_amount, total_amount, currency, status, reference_number,
			   purchase_order_id, receipt_url, attachment_urls, description,
			   notes, metadata, created_at, updated_at, created_by, updated_by,
			   approved_by, approved_at, deleted_at
		FROM expenses
		WHERE organization_id = $1 AND deleted_at IS NULL
	`

	args := []interface{}{orgID}
	argCount := 1

	if filters.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filters.Category)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.ExpenseDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND expense_date >= $%d", argCount)
		args = append(args, filters.ExpenseDateFrom)
	}

	if filters.ExpenseDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND expense_date <= $%d", argCount)
		args = append(args, filters.ExpenseDateTo)
	}

	query += " ORDER BY expense_date DESC"

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

	var expenses []staff.Expense
	for rows.Next() {
		var e staff.Expense
		if err := rows.Scan(
			&e.ID, &e.OrganizationID, &e.LocationID, &e.ExpenseNumber,
			&e.ExpenseDate, &e.Category, &e.Subcategory, &e.PayeeName,
			&e.PaymentMethod, &e.Amount, &e.TaxAmount, &e.TotalAmount,
			&e.Currency, &e.Status, &e.ReferenceNumber,
			&e.PurchaseOrderID, &e.ReceiptURL, &e.AttachmentURLs,
			&e.Description, &e.Notes, &e.Metadata, &e.CreatedAt,
			&e.UpdatedAt, &e.CreatedBy, &e.UpdatedBy, &e.ApprovedBy,
			&e.ApprovedAt, &e.DeletedAt,
		); err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, rows.Err()
}

func (r *StaffRepository) CountExpenses(ctx context.Context, orgID uuid.UUID, filters staff.ExpenseFilters) (int64, error) {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return 0, err
	}

	query := "SELECT COUNT(*) FROM expenses WHERE organization_id = $1 AND deleted_at IS NULL"
	args := []interface{}{orgID}
	argCount := 1

	if filters.Category != nil {
		argCount++
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filters.Category)
	}

	if filters.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status)
	}

	if filters.LocationID != nil {
		argCount++
		query += fmt.Sprintf(" AND location_id = $%d", argCount)
		args = append(args, filters.LocationID)
	}

	if filters.ExpenseDateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND expense_date >= $%d", argCount)
		args = append(args, filters.ExpenseDateFrom)
	}

	if filters.ExpenseDateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND expense_date <= $%d", argCount)
		args = append(args, filters.ExpenseDateTo)
	}

	var count int64
	err := r.db.Pool.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *StaffRepository) UpdateExpense(ctx context.Context, expense *staff.Expense) error {
	if err := r.db.SetOrganizationContext(ctx, expense.OrganizationID.String()); err != nil {
		return err
	}

	query := `
		UPDATE expenses
		SET location_id = $2, expense_number = $3, expense_date = $4,
			category = $5, subcategory = $6, payee_name = $7, payment_method = $8,
			amount = $9, tax_amount = $10, total_amount = $11, currency = $12,
			status = $13, reference_number = $14, purchase_order_id = $15,
			receipt_url = $16, attachment_urls = $17, description = $18, notes = $19,
			metadata = $20, updated_at = $21, updated_by = $22, approved_by = $23,
			approved_at = $24
		WHERE id = $1 AND organization_id = $25 AND deleted_at IS NULL
	`

	_, err := r.db.Pool.Exec(ctx, query,
		expense.ID, expense.LocationID, expense.ExpenseNumber, expense.ExpenseDate,
		expense.Category, expense.Subcategory, expense.PayeeName, expense.PaymentMethod,
		expense.Amount, expense.TaxAmount, expense.TotalAmount, expense.Currency,
		expense.Status, expense.ReferenceNumber, expense.PurchaseOrderID,
		expense.ReceiptURL, expense.AttachmentURLs, expense.Description, expense.Notes,
		expense.Metadata, expense.UpdatedAt, expense.UpdatedBy, expense.ApprovedBy,
		expense.ApprovedAt, expense.OrganizationID,
	)

	return err
}

func (r *StaffRepository) DeleteExpense(ctx context.Context, orgID, expenseID uuid.UUID) error {
	if err := r.db.SetOrganizationContext(ctx, orgID.String()); err != nil {
		return err
	}

	query := `
		UPDATE expenses
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, expenseID, orgID)
	return err
}
