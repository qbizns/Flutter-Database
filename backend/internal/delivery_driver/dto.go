package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DeliveryDriversResponse represents a delivery_drivers response
type DeliveryDriversResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	UserId *uuid.UUID `json:"user_id"`
	
	DriverCode string `json:"driver_code"`
	
	FullName string `json:"full_name"`
	
	Phone *string `json:"phone"`
	
	Email *string `json:"email"`
	
	EmergencyContactName *string `json:"emergency_contact_name"`
	
	EmergencyContactPhone *string `json:"emergency_contact_phone"`
	
	VehicleType *string `json:"vehicle_type"`
	
	VehicleMake *string `json:"vehicle_make"`
	
	VehicleModel *string `json:"vehicle_model"`
	
	VehicleYear *int64 `json:"vehicle_year"`
	
	VehicleColor *string `json:"vehicle_color"`
	
	LicensePlate *string `json:"license_plate"`
	
	DriversLicenseNumber *string `json:"drivers_license_number"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date"`
	
	InsurancePolicyNumber *string `json:"insurance_policy_number"`
	
	InsuranceExpiryDate *time.Time `json:"insurance_expiry_date"`
	
	HireDate *time.Time `json:"hire_date"`
	
	EmploymentType *string `json:"employment_type"`
	
	Status *string `json:"status"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries"`
	
	Rating *float64 `json:"rating"`
	
	RatingCount *int64 `json:"rating_count"`
	
	CurrentLocation json.RawMessage `json:"current_location"`
	
	IsAvailable *bool `json:"is_available"`
	
	LastLocationUpdate *time.Time `json:"last_location_update"`
	
	CommissionRate *float64 `json:"commission_rate"`
	
	PaymentMethod *string `json:"payment_method"`
	
	Documents json.RawMessage `json:"documents"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	'active', *string `json:"'active',"`
	
	'fullTime', *string `json:"'full_time',"`
	
	'bike', *string `json:"'bike',"`
	
}

// CreateDeliveryDriversRequest represents a request to create a delivery_drivers
type CreateDeliveryDriversRequest struct {
	
	UserId *uuid.UUID `json:"user_id"`
	
	DriverCode string `json:"driver_code" validate:"required"`
	
	FullName string `json:"full_name" validate:"required"`
	
	Phone *string `json:"phone" validate:"e164"`
	
	Email *string `json:"email" validate:"email"`
	
	EmergencyContactName *string `json:"emergency_contact_name"`
	
	EmergencyContactPhone *string `json:"emergency_contact_phone" validate:"e164"`
	
	VehicleType *string `json:"vehicle_type"`
	
	VehicleMake *string `json:"vehicle_make"`
	
	VehicleModel *string `json:"vehicle_model"`
	
	VehicleYear *int64 `json:"vehicle_year"`
	
	VehicleColor *string `json:"vehicle_color"`
	
	LicensePlate *string `json:"license_plate"`
	
	DriversLicenseNumber *string `json:"drivers_license_number"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date"`
	
	InsurancePolicyNumber *string `json:"insurance_policy_number"`
	
	InsuranceExpiryDate *time.Time `json:"insurance_expiry_date"`
	
	HireDate *time.Time `json:"hire_date"`
	
	EmploymentType *string `json:"employment_type"`
	
	Status *string `json:"status"`
	
	TotalDeliveries *int64 `json:"total_deliveries"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries"`
	
	Rating *float64 `json:"rating"`
	
	RatingCount *int64 `json:"rating_count"`
	
	CurrentLocation json.RawMessage `json:"current_location"`
	
	IsAvailable *bool `json:"is_available"`
	
	LastLocationUpdate *time.Time `json:"last_location_update"`
	
	CommissionRate *float64 `json:"commission_rate"`
	
	PaymentMethod *string `json:"payment_method"`
	
	Documents json.RawMessage `json:"documents"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
	'active', *string `json:"'active',"`
	
	'fullTime', *string `json:"'full_time',"`
	
	'bike', *string `json:"'bike',"`
	
}

// Validate validates the create request
func (r *CreateDeliveryDriversRequest) Validate() error {
	
	if r.DriverCode == "" {
		return fmt.Errorf("driver_code is required")
	}
	
	if r.FullName == "" {
		return fmt.Errorf("full_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateDeliveryDriversRequest represents a request to update a delivery_drivers
type UpdateDeliveryDriversRequest struct {
	
	UserId *uuid.UUID `json:"user_id,omitempty"`
	
	DriverCode *string `json:"driver_code,omitempty" validate:"omitempty,required"`
	
	FullName *string `json:"full_name,omitempty" validate:"omitempty,required"`
	
	Phone *string `json:"phone,omitempty" validate:"omitempty,e164"`
	
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
	
	EmergencyContactName *string `json:"emergency_contact_name,omitempty"`
	
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty" validate:"omitempty,e164"`
	
	VehicleType *string `json:"vehicle_type,omitempty"`
	
	VehicleMake *string `json:"vehicle_make,omitempty"`
	
	VehicleModel *string `json:"vehicle_model,omitempty"`
	
	VehicleYear *int64 `json:"vehicle_year,omitempty"`
	
	VehicleColor *string `json:"vehicle_color,omitempty"`
	
	LicensePlate *string `json:"license_plate,omitempty"`
	
	DriversLicenseNumber *string `json:"drivers_license_number,omitempty"`
	
	LicenseExpiryDate *time.Time `json:"license_expiry_date,omitempty"`
	
	InsurancePolicyNumber *string `json:"insurance_policy_number,omitempty"`
	
	InsuranceExpiryDate *time.Time `json:"insurance_expiry_date,omitempty"`
	
	HireDate *time.Time `json:"hire_date,omitempty"`
	
	EmploymentType *string `json:"employment_type,omitempty"`
	
	Status *string `json:"status,omitempty"`
	
	TotalDeliveries *int64 `json:"total_deliveries,omitempty"`
	
	SuccessfulDeliveries *int64 `json:"successful_deliveries,omitempty"`
	
	Rating *float64 `json:"rating,omitempty"`
	
	RatingCount *int64 `json:"rating_count,omitempty"`
	
	CurrentLocation *json.RawMessage `json:"current_location,omitempty"`
	
	IsAvailable *bool `json:"is_available,omitempty"`
	
	LastLocationUpdate *time.Time `json:"last_location_update,omitempty"`
	
	CommissionRate *float64 `json:"commission_rate,omitempty"`
	
	PaymentMethod *string `json:"payment_method,omitempty"`
	
	Documents *json.RawMessage `json:"documents,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
	'active', *string `json:"'active',,omitempty"`
	
	'fullTime', *string `json:"'full_time',,omitempty"`
	
	'bike', *string `json:"'bike',,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateDeliveryDriversRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.UserId != nil {
		hasUpdate = true
	}
	
	if r.DriverCode != nil {
		hasUpdate = true
	}
	
	if r.FullName != nil {
		hasUpdate = true
	}
	
	if r.Phone != nil {
		hasUpdate = true
	}
	
	if r.Email != nil {
		hasUpdate = true
	}
	
	if r.EmergencyContactName != nil {
		hasUpdate = true
	}
	
	if r.EmergencyContactPhone != nil {
		hasUpdate = true
	}
	
	if r.VehicleType != nil {
		hasUpdate = true
	}
	
	if r.VehicleMake != nil {
		hasUpdate = true
	}
	
	if r.VehicleModel != nil {
		hasUpdate = true
	}
	
	if r.VehicleYear != nil {
		hasUpdate = true
	}
	
	if r.VehicleColor != nil {
		hasUpdate = true
	}
	
	if r.LicensePlate != nil {
		hasUpdate = true
	}
	
	if r.DriversLicenseNumber != nil {
		hasUpdate = true
	}
	
	if r.LicenseExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.InsurancePolicyNumber != nil {
		hasUpdate = true
	}
	
	if r.InsuranceExpiryDate != nil {
		hasUpdate = true
	}
	
	if r.HireDate != nil {
		hasUpdate = true
	}
	
	if r.EmploymentType != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.TotalDeliveries != nil {
		hasUpdate = true
	}
	
	if r.SuccessfulDeliveries != nil {
		hasUpdate = true
	}
	
	if r.Rating != nil {
		hasUpdate = true
	}
	
	if r.RatingCount != nil {
		hasUpdate = true
	}
	
	if r.CurrentLocation != nil {
		hasUpdate = true
	}
	
	if r.IsAvailable != nil {
		hasUpdate = true
	}
	
	if r.LastLocationUpdate != nil {
		hasUpdate = true
	}
	
	if r.CommissionRate != nil {
		hasUpdate = true
	}
	
	if r.PaymentMethod != nil {
		hasUpdate = true
	}
	
	if r.Documents != nil {
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
	
	if r.'fullTime', != nil {
		hasUpdate = true
	}
	
	if r.'bike', != nil {
		hasUpdate = true
	}
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// DeliveryDriversListResponse represents a paginated list of delivery_drivers records
type DeliveryDriversListResponse struct {
	Items      []*DeliveryDriversResponse `json:"items"`
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
