package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DevicesResponse represents a devices response
type DevicesResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId *uuid.UUID `json:"location_id"`
	
	DeviceCode string `json:"device_code"`
	
	DeviceName string `json:"device_name"`
	
	DeviceType string `json:"device_type"`
	
	Manufacturer *string `json:"manufacturer"`
	
	Model *string `json:"model"`
	
	SerialNumber *string `json:"serial_number"`
	
	MacAddress *string `json:"mac_address"`
	
	IpAddress *string `json:"ip_address"`
	
	DeviceConfig json.RawMessage `json:"device_config"`
	
	ScreenResolution *string `json:"screen_resolution"`
	
	OsVersion *string `json:"os_version"`
	
	ConnectionType *string `json:"connection_type"`
	
	ConnectionString *string `json:"connection_string"`
	
	Status *string `json:"status"`
	
	LastOnlineAt *time.Time `json:"last_online_at"`
	
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at"`
	
	AssignedToUserId *uuid.UUID `json:"assigned_to_user_id"`
	
	AssignedToStationId *uuid.UUID `json:"assigned_to_station_id"`
	
	PurchaseDate *time.Time `json:"purchase_date"`
	
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date"`
	
	LicenseKey *string `json:"license_key"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date"`
	
	InstallationNotes *string `json:"installation_notes"`
	
	MaintenanceNotes *string `json:"maintenance_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'active', *string `json:"'active',"`
	
	'posTerminal', *string `json:"'pos_terminal',"`
	
	'printer', *string `json:"'printer',"`
	
	'kitchenPrinter', *int64 `json:"'kitchen_printer',"`
	
	'network', *string `json:"'network',"`
	
}

// CreateDevicesRequest represents a request to create a devices
type CreateDevicesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id"`
	
	DeviceCode string `json:"device_code" validate:"required"`
	
	DeviceName string `json:"device_name" validate:"required"`
	
	DeviceType string `json:"device_type" validate:"required"`
	
	Manufacturer *string `json:"manufacturer"`
	
	Model *string `json:"model"`
	
	SerialNumber *string `json:"serial_number"`
	
	MacAddress *string `json:"mac_address"`
	
	IpAddress *string `json:"ip_address"`
	
	DeviceConfig json.RawMessage `json:"device_config"`
	
	ScreenResolution *string `json:"screen_resolution"`
	
	OsVersion *string `json:"os_version"`
	
	ConnectionType *string `json:"connection_type"`
	
	ConnectionString *string `json:"connection_string"`
	
	Status *string `json:"status"`
	
	LastOnlineAt *time.Time `json:"last_online_at"`
	
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at"`
	
	AssignedToUserId *uuid.UUID `json:"assigned_to_user_id"`
	
	AssignedToStationId *uuid.UUID `json:"assigned_to_station_id"`
	
	PurchaseDate *time.Time `json:"purchase_date"`
	
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date"`
	
	LicenseKey *string `json:"license_key"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date"`
	
	InstallationNotes *string `json:"installation_notes"`
	
	MaintenanceNotes *string `json:"maintenance_notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'active', *string `json:"'active',"`
	
	'posTerminal', *string `json:"'pos_terminal',"`
	
	'printer', *string `json:"'printer',"`
	
	'kitchenPrinter', *int64 `json:"'kitchen_printer',"`
	
	'network', *string `json:"'network',"`
	
}

// Validate validates the create request
func (r *CreateDevicesRequest) Validate() error {
	
	if r.DeviceCode == "" {
		return fmt.Errorf("device_code is required")
	}
	
	if r.DeviceName == "" {
		return fmt.Errorf("device_name is required")
	}
	
	if r.DeviceType == "" {
		return fmt.Errorf("device_type is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDevicesRequest represents a request to update a devices
type UpdateDevicesRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty"`
	
	DeviceCode *string `json:"device_code,omitempty" validate:"omitempty,required"`
	
	DeviceName *string `json:"device_name,omitempty" validate:"omitempty,required"`
	
	DeviceType *string `json:"device_type,omitempty" validate:"omitempty,required"`
	
	Manufacturer *string `json:"manufacturer,omitempty"`
	
	Model *string `json:"model,omitempty"`
	
	SerialNumber *string `json:"serial_number,omitempty"`
	
	MacAddress *string `json:"mac_address,omitempty"`
	
	IpAddress *string `json:"ip_address,omitempty"`
	
	DeviceConfig *json.RawMessage `json:"device_config,omitempty"`
	
	ScreenResolution *string `json:"screen_resolution,omitempty"`
	
	OsVersion *string `json:"os_version,omitempty"`
	
	ConnectionType *string `json:"connection_type,omitempty"`
	
	ConnectionString *string `json:"connection_string,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	LastOnlineAt *time.Time `json:"last_online_at,omitempty"`
	
	LastHeartbeatAt *time.Time `json:"last_heartbeat_at,omitempty"`
	
	AssignedToUserId *uuid.UUID `json:"assigned_to_user_id,omitempty"`
	
	AssignedToStationId *uuid.UUID `json:"assigned_to_station_id,omitempty"`
	
	PurchaseDate *time.Time `json:"purchase_date,omitempty"`
	
	WarrantyExpiryDate *time.Time `json:"warranty_expiry_date,omitempty"`
	
	LicenseKey *string `json:"license_key,omitempty"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date,omitempty"`
	
	InstallationNotes *string `json:"installation_notes,omitempty"`
	
	MaintenanceNotes *string `json:"maintenance_notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'active', *string `json:"'active',,omitempty"`
	
	'posTerminal', *string `json:"'pos_terminal',,omitempty"`
	
	'printer', *string `json:"'printer',,omitempty"`
	
	'kitchenPrinter', *int64 `json:"'kitchen_printer',,omitempty"`
	
	'network', *string `json:"'network',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDevicesRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.DeviceCode != nil {
		hasUpdate = true
	}
	
	if r.DeviceName != nil {
		hasUpdate = true
	}
	
	if r.DeviceType != nil {
		hasUpdate = true
	}
	
	if r.Manufacturer != nil {
		hasUpdate = true
	}
	
	if r.Model != nil {
		hasUpdate = true
	}
	
	if r.SerialNumber != nil {
		hasUpdate = true
	}
	
	if r.MacAddress != nil {
		hasUpdate = true
	}
	
	if r.IpAddress != nil {
		hasUpdate = true
	}
	
	if r.DeviceConfig != nil {
		hasUpdate = true
	}
	
	if r.ScreenResolution != nil {
		hasUpdate = true
	}
	
	if r.OsVersion != nil {
		hasUpdate = true
	}
	
	if r.ConnectionType != nil {
		hasUpdate = true
	}
	
	if r.ConnectionString != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.LastOnlineAt != nil {
		hasUpdate = true
	}
	
	if r.LastHeartbeatAt != nil {
		hasUpdate = true
	}
	
	if r.AssignedToUserId != nil {
		hasUpdate = true
	}
	
	if r.AssignedToStationId != nil {
		hasUpdate = true
	}
	
	if r.PurchaseDate != nil {
		hasUpdate = true
	}
	
	if r.WarrantyExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.LicenseKey != nil {
		hasUpdate = true
	}
	
	if r.LicenseExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.InstallationNotes != nil {
		hasUpdate = true
	}
	
	if r.MaintenanceNotes != nil {
		hasUpdate = true
	}
	
	if r.Metadata != nil {
		hasUpdate = true
	}
	
	if r.CreatedBy != nil {
		hasUpdate = true
	}
	
	if r.UpdatedBy != nil {
		hasUpdate = true
	}
	
	if r.'active', != nil {
		hasUpdate = true
	}
	
	if r.'posTerminal', != nil {
		hasUpdate = true
	}
	
	if r.'printer', != nil {
		hasUpdate = true
	}
	
	if r.'kitchenPrinter', != nil {
		hasUpdate = true
	}
	
	if r.'network', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DevicesListResponse represents a paginated list of devices records
type DevicesListResponse struct {
	Items      []*DevicesResponse `json:"items"`
	Pagination Pagination             `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}
