package device

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

// Repository handles database operations for Devices
type Repository struct {
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewRepository creates a new Devices repository
func NewRepository(db *pgxpool.Pool, logger *logging.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Devices represents a devices entity
type Devices struct {
	Id *uuid.UUID `json:"id" db:"id"`
	OrganizationId uuid.UUID `json:"organization_id" db:"organization_id"`
	LocationId *uuid.UUID `json:"location_id" db:"location_id"`
	DeviceCode string `json:"device_code" db:"device_code"`
	DeviceName string `json:"device_name" db:"device_name"`
	DeviceType string `json:"device_type" db:"device_type"`
	Manufacturer *string `json:"manufacturer" db:"manufacturer"`
	Model *string `json:"model" db:"model"`
	SerialNumber *string `json:"serial_number" db:"serial_number"`
	MacAddress *string `json:"mac_address" db:"mac_address"`
	IpAddress *string `json:"ip_address" db:"ip_address"`
	DeviceConfig json.RawMessage `json:"device_config" db:"device_config"`
	ScreenResolution *string `json:"screen_resolution" db:"screen_resolution"`
	OsVersion *string `json:"os_version" db:"os_version"`
	ConnectionType *string `json:"connection_type" db:"connection_type"`
	ConnectionString *string `json:"connection_string" db:"connection_string"`
	Status *string `json:"status" db:"status"`
	LastOnlineAt *time.Time `json:"last_online_at" db:"last_online_at"`
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at" db:"last_heartbeat_at"`
	AssignedToUserId *uuid.UUID `json:"assigned_to_user_id" db:"assigned_to_user_id"`
	AssignedToStationId *uuid.UUID `json:"assigned_to_station_id" db:"assigned_to_station_id"`
	PurchaseDate *time.Time `json:"purchase_date" db:"purchase_date"`
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date" db:"warranty_expiry_date"`
	LicenseKey *string `json:"license_key" db:"license_key"`
	LicenseExpiryDate *time.Time `json:"license_expiry_date" db:"license_expiry_date"`
	InstallationNotes *string `json:"installation_notes" db:"installation_notes"`
	MaintenanceNotes *string `json:"maintenance_notes" db:"maintenance_notes"`
	Metadata json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by" db:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by" db:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	'active', *string `json:"'active'," db:"'active',"`
	'posTerminal', *string `json:"'pos_terminal'," db:"'pos_terminal',"`
	'printer', *string `json:"'printer'," db:"'printer',"`
	'kitchenPrinter', *int64 `json:"'kitchen_printer'," db:"'kitchen_printer',"`
	'network', *string `json:"'network'," db:"'network',"`
}

// Create inserts a new devices record
func (r *Repository) Create(ctx context.Context, tx pgx.Tx, entity *Devices) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("INSERT", "devices", duration, nil)
	}()

	query := `
		INSERT INTO devices (
			, organization_id
			, location_id
			, device_code
			, device_name
			, device_type
			, manufacturer
			, model
			, serial_number
			, mac_address
			, ip_address
			, device_config
			, screen_resolution
			, os_version
			, connection_type
			, connection_string
			, status
			, last_online_at
			, last_heartbeat_at
			, assigned_to_user_id
			, assigned_to_station_id
			, purchase_date
			, warranty_expiry_date
			, license_key
			, license_expiry_date
			, installation_notes
			, maintenance_notes
			, metadata
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'pos_terminal',
			, 'printer',
			, 'kitchen_printer',
			, 'network',
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
			, $28
			, $31
			, $32
			, $33
			, $34
			, $35
			, $36
			, $37
			, $38
		)
		RETURNING id, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.DeviceCode,
		entity.DeviceName,
		entity.DeviceType,
		entity.Manufacturer,
		entity.Model,
		entity.SerialNumber,
		entity.MacAddress,
		entity.IpAddress,
		entity.DeviceConfig,
		entity.ScreenResolution,
		entity.OsVersion,
		entity.ConnectionType,
		entity.ConnectionString,
		entity.Status,
		entity.LastOnlineAt,
		entity.LastHeartbeatAt,
		entity.AssignedToUserId,
		entity.AssignedToStationId,
		entity.PurchaseDate,
		entity.WarrantyExpiryDate,
		entity.LicenseKey,
		entity.LicenseExpiryDate,
		entity.InstallationNotes,
		entity.MaintenanceNotes,
		entity.Metadata,
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'active',,
		entity.'posTerminal',,
		entity.'printer',,
		entity.'kitchenPrinter',,
		entity.'network',,
	)

	
	err := row.Scan(&entity.Id, &entity.CreatedAt, func() *time.Time { t := time.Now(); return &t }())
	

	if err != nil {
		r.logger.Error("failed to create devices", zap.Error(err))
		return fmt.Errorf("failed to create devices: %w", err)
	}

	r.logger.Info("created devices",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", entity.OrganizationId.String()),
	)

	return nil
}

// GetByID retrieves a devices by ID
func (r *Repository) GetByID(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Devices, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "devices", duration, nil)
	}()

	query := `
		SELECT
			id
			, organization_id
			, location_id
			, device_code
			, device_name
			, device_type
			, manufacturer
			, model
			, serial_number
			, mac_address
			, ip_address
			, device_config
			, screen_resolution
			, os_version
			, connection_type
			, connection_string
			, last_online_at
			, last_heartbeat_at
			, assigned_to_user_id
			, assigned_to_station_id
			, purchase_date
			, warranty_expiry_date
			, license_key
			, license_expiry_date
			, installation_notes
			, maintenance_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'pos_terminal',
			, 'printer',
			, 'kitchen_printer',
			, 'network',
		FROM devices
		WHERE id = $1
		AND deleted_at IS NULL
	`

	var entity Devices
	err := tx.QueryRow(ctx, query, id).Scan(
		&entity.Id,
		&entity.OrganizationId,
		&entity.LocationId,
		&entity.DeviceCode,
		&entity.DeviceName,
		&entity.DeviceType,
		&entity.Manufacturer,
		&entity.Model,
		&entity.SerialNumber,
		&entity.MacAddress,
		&entity.IpAddress,
		&entity.DeviceConfig,
		&entity.ScreenResolution,
		&entity.OsVersion,
		&entity.ConnectionType,
		&entity.ConnectionString,
		&entity.Status,
		&entity.LastOnlineAt,
		&entity.LastHeartbeatAt,
		&entity.AssignedToUserId,
		&entity.AssignedToStationId,
		&entity.PurchaseDate,
		&entity.WarrantyExpiryDate,
		&entity.LicenseKey,
		&entity.LicenseExpiryDate,
		&entity.InstallationNotes,
		&entity.MaintenanceNotes,
		&entity.Metadata,
		&entity.CreatedAt,
		func() *time.Time { t := time.Now(); return &t }(),
		&entity.CreatedBy,
		&entity.UpdatedBy,
		&entity.DeletedAt,
		&entity.'active',,
		&entity.'posTerminal',,
		&entity.'printer',,
		&entity.'kitchenPrinter',,
		&entity.'network',,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("devices not found")
	}

	if err != nil {
		r.logger.Error("failed to get devices", zap.Error(err), zap.String("id", id.String()))
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}

	return &entity, nil
}

// List retrieves a paginated list of devices records
func (r *Repository) List(ctx context.Context, tx pgx.Tx, limit, offset int) ([]*Devices, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "devices", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM devices
		WHERE deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count devices records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, device_code
			, device_name
			, device_type
			, manufacturer
			, model
			, serial_number
			, mac_address
			, ip_address
			, device_config
			, screen_resolution
			, os_version
			, connection_type
			, connection_string
			, last_online_at
			, last_heartbeat_at
			, assigned_to_user_id
			, assigned_to_station_id
			, purchase_date
			, warranty_expiry_date
			, license_key
			, license_expiry_date
			, installation_notes
			, maintenance_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'pos_terminal',
			, 'printer',
			, 'kitchen_printer',
			, 'network',
		FROM devices
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list devices", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var entities []*Devices
	for rows.Next() {
		var entity Devices
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.DeviceCode,
			&entity.DeviceName,
			&entity.DeviceType,
			&entity.Manufacturer,
			&entity.Model,
			&entity.SerialNumber,
			&entity.MacAddress,
			&entity.IpAddress,
			&entity.DeviceConfig,
			&entity.ScreenResolution,
			&entity.OsVersion,
			&entity.ConnectionType,
			&entity.ConnectionString,
			&entity.Status,
			&entity.LastOnlineAt,
			&entity.LastHeartbeatAt,
			&entity.AssignedToUserId,
			&entity.AssignedToStationId,
			&entity.PurchaseDate,
			&entity.WarrantyExpiryDate,
			&entity.LicenseKey,
			&entity.LicenseExpiryDate,
			&entity.InstallationNotes,
			&entity.MaintenanceNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'active',,
			&entity.'posTerminal',,
			&entity.'printer',,
			&entity.'kitchenPrinter',,
			&entity.'network',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan devices: %w", err)
		}
		entities = append(entities, &entity)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error iterating devices rows: %w", rows.Err())
	}

	return entities, total, nil
}

// Update updates an existing devices record
func (r *Repository) Update(ctx context.Context, tx pgx.Tx, entity *Devices) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "devices", duration, nil)
	}()

	query := `
		UPDATE devices
		SET
			, organization_id = $2
			, location_id = $3
			, device_code = $4
			, device_name = $5
			, device_type = $6
			, manufacturer = $7
			, model = $8
			, serial_number = $9
			, mac_address = $10
			, ip_address = $11
			, device_config = $12
			, screen_resolution = $13
			, os_version = $14
			, connection_type = $15
			, connection_string = $16
			, status = $17
			, last_online_at = $18
			, last_heartbeat_at = $19
			, assigned_to_user_id = $20
			, assigned_to_station_id = $21
			, purchase_date = $22
			, warranty_expiry_date = $23
			, license_key = $24
			, license_expiry_date = $25
			, installation_notes = $26
			, maintenance_notes = $27
			, metadata = $28
			, updated_at = $30
			, created_by = $31
			, updated_by = $32
			, deleted_at = $33
			, 'active', = $34
			, 'pos_terminal', = $35
			, 'printer', = $36
			, 'kitchen_printer', = $37
			, 'network', = $38
			, updated_at = CURRENT_TIMESTAMP
		WHERE id = $39
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query,
		entity.OrganizationId,
		entity.LocationId,
		entity.DeviceCode,
		entity.DeviceName,
		entity.DeviceType,
		entity.Manufacturer,
		entity.Model,
		entity.SerialNumber,
		entity.MacAddress,
		entity.IpAddress,
		entity.DeviceConfig,
		entity.ScreenResolution,
		entity.OsVersion,
		entity.ConnectionType,
		entity.ConnectionString,
		entity.Status,
		entity.LastOnlineAt,
		entity.LastHeartbeatAt,
		entity.AssignedToUserId,
		entity.AssignedToStationId,
		entity.PurchaseDate,
		entity.WarrantyExpiryDate,
		entity.LicenseKey,
		entity.LicenseExpiryDate,
		entity.InstallationNotes,
		entity.MaintenanceNotes,
		entity.Metadata,
		time.Now(),
		entity.CreatedBy,
		entity.UpdatedBy,
		entity.DeletedAt,
		entity.'active',,
		entity.'posTerminal',,
		entity.'printer',,
		entity.'kitchenPrinter',,
		entity.'network',,
		entity.Id,
	)

	if err != nil {
		r.logger.Error("failed to update devices", zap.Error(err))
		return fmt.Errorf("failed to update devices: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("devices not found or already deleted")
	}

	r.logger.Info("updated devices",
		zap.String("id", entity.Id.String()),
	)

	return nil
}


// Delete soft-deletes a devices record
func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("UPDATE", "devices", duration, nil)
	}()

	query := `
		UPDATE devices
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1
		AND deleted_at IS NULL
	`

	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete devices", zap.Error(err), zap.String("id", id.String()))
		return fmt.Errorf("failed to delete devices: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("devices not found or already deleted")
	}

	r.logger.Info("deleted devices", zap.String("id", id.String()))
	return nil
}



// ListByOrganization retrieves devices records for a specific organization
func (r *Repository) ListByOrganization(ctx context.Context, tx pgx.Tx, orgID uuid.UUID, limit, offset int) ([]*Devices, int, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		metrics.RecordDatabaseQuery("SELECT", "devices", duration, nil)
	}()

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM devices
		WHERE organization_id = $1
		AND deleted_at IS NULL
	`

	err := tx.QueryRow(ctx, countQuery, orgID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count devices records: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id
			, organization_id
			, location_id
			, device_code
			, device_name
			, device_type
			, manufacturer
			, model
			, serial_number
			, mac_address
			, ip_address
			, device_config
			, screen_resolution
			, os_version
			, connection_type
			, connection_string
			, status
			, last_online_at
			, last_heartbeat_at
			, assigned_to_user_id
			, assigned_to_station_id
			, purchase_date
			, warranty_expiry_date
			, license_key
			, license_expiry_date
			, installation_notes
			, maintenance_notes
			, metadata
			, created_at
			, updated_at
			, created_by
			, updated_by
			, deleted_at
			, 'active',
			, 'pos_terminal',
			, 'printer',
			, 'kitchen_printer',
			, 'network',
		FROM devices
		WHERE organization_id = $1
		AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := tx.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		r.logger.Error("failed to list devices by organization", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list devices: %w", err)
	}
	defer rows.Close()

	var entities []*Devices
	for rows.Next() {
		var entity Devices
		err := rows.Scan(
			&entity.Id,
			&entity.OrganizationId,
			&entity.LocationId,
			&entity.DeviceCode,
			&entity.DeviceName,
			&entity.DeviceType,
			&entity.Manufacturer,
			&entity.Model,
			&entity.SerialNumber,
			&entity.MacAddress,
			&entity.IpAddress,
			&entity.DeviceConfig,
			&entity.ScreenResolution,
			&entity.OsVersion,
			&entity.ConnectionType,
			&entity.ConnectionString,
			&entity.Status,
			&entity.LastOnlineAt,
			&entity.LastHeartbeatAt,
			&entity.AssignedToUserId,
			&entity.AssignedToStationId,
			&entity.PurchaseDate,
			&entity.WarrantyExpiryDate,
			&entity.LicenseKey,
			&entity.LicenseExpiryDate,
			&entity.InstallationNotes,
			&entity.MaintenanceNotes,
			&entity.Metadata,
			&entity.CreatedAt,
			func() *time.Time { t := time.Now(); return &t }(),
			&entity.CreatedBy,
			&entity.UpdatedBy,
			&entity.DeletedAt,
			&entity.'active',,
			&entity.'posTerminal',,
			&entity.'printer',,
			&entity.'kitchenPrinter',,
			&entity.'network',,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan devices: %w", err)
		}
		entities = append(entities, &entity)
	}

	return entities, total, nil
}

