package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ReservationsResponse represents a reservations response
type ReservationsResponse struct {
	
	Id *uuid.UUID `json:"id"`
	
	OrganizationId uuid.UUID `json:"organization_id"`
	
	LocationId uuid.UUID `json:"location_id"`
	
	TableId *uuid.UUID `json:"table_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	ReservationNumber string `json:"reservation_number"`
	
	ReservationDate time.Time `json:"reservation_date"`
	
	ReservationTime string `json:"reservation_time"`
	
	DurationMinutes *int64 `json:"duration_minutes"`
	
	PartySize int64 `json:"party_size"`
	
	CustomerName string `json:"customer_name"`
	
	CustomerPhone *string `json:"customer_phone"`
	
	CustomerEmail *string `json:"customer_email"`
	
	Status *string `json:"status"`
	
	AssignedWaiterId *uuid.UUID `json:"assigned_waiter_id"`
	
	AssignedAt *time.Time `json:"assigned_at"`
	
	SeatedAt *time.Time `json:"seated_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	SpecialRequests *string `json:"special_requests"`
	
	Occasion *string `json:"occasion"`
	
	DietaryRestrictions *string `json:"dietary_restrictions"`
	
	ConfirmationCode *string `json:"confirmation_code"`
	
	ConfirmedAt *time.Time `json:"confirmed_at"`
	
	ConfirmedBy *uuid.UUID `json:"confirmed_by"`
	
	ReminderSentAt *time.Time `json:"reminder_sent_at"`
	
	NotificationPreferences json.RawMessage `json:"notification_preferences"`
	
	CancelledAt *time.Time `json:"cancelled_at"`
	
	CancelledBy *uuid.UUID `json:"cancelled_by"`
	
	CancellationReason *string `json:"cancellation_reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedAt *time.Time `json:"created_at"`
	
	UpdatedAt *time.Time `json:"updated_at"`
	
	DeletedAt *time.Time `json:"deleted_at"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// CreateReservationsRequest represents a request to create a reservations
type CreateReservationsRequest struct {
	
	LocationId uuid.UUID `json:"location_id" validate:"required"`
	
	TableId *uuid.UUID `json:"table_id"`
	
	CustomerId *uuid.UUID `json:"customer_id"`
	
	ReservationNumber string `json:"reservation_number" validate:"required"`
	
	ReservationDate time.Time `json:"reservation_date" validate:"required"`
	
	ReservationTime string `json:"reservation_time" validate:"required"`
	
	DurationMinutes *int64 `json:"duration_minutes"`
	
	PartySize int64 `json:"party_size" validate:"required"`
	
	CustomerName string `json:"customer_name" validate:"required"`
	
	CustomerPhone *string `json:"customer_phone" validate:"e164"`
	
	CustomerEmail *string `json:"customer_email" validate:"email"`
	
	Status *string `json:"status"`
	
	AssignedWaiterId *uuid.UUID `json:"assigned_waiter_id"`
	
	AssignedAt *time.Time `json:"assigned_at"`
	
	SeatedAt *time.Time `json:"seated_at"`
	
	CompletedAt *time.Time `json:"completed_at"`
	
	SpecialRequests *string `json:"special_requests"`
	
	Occasion *string `json:"occasion"`
	
	DietaryRestrictions *string `json:"dietary_restrictions"`
	
	ConfirmationCode *string `json:"confirmation_code"`
	
	ConfirmedAt *time.Time `json:"confirmed_at"`
	
	ConfirmedBy *uuid.UUID `json:"confirmed_by"`
	
	ReminderSentAt *time.Time `json:"reminder_sent_at"`
	
	NotificationPreferences json.RawMessage `json:"notification_preferences"`
	
	CancelledAt *time.Time `json:"cancelled_at"`
	
	CancelledBy *uuid.UUID `json:"cancelled_by"`
	
	CancellationReason *string `json:"cancellation_reason"`
	
	Notes *string `json:"notes"`
	
	Metadata json.RawMessage `json:"metadata"`
	
	CreatedBy *uuid.UUID `json:"created_by"`
	
	UpdatedBy *uuid.UUID `json:"updated_by"`
	
}

// Validate validates the create request
func (r *CreateReservationsRequest) Validate() error {
	
	if r.LocationId == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	
	if r.ReservationNumber == "" {
		return fmt.Errorf("reservation_number is required")
	}
	
	if r.ReservationDate == nil {
		return fmt.Errorf("reservation_date is required")
	}
	
	if r.ReservationTime == "" {
		return fmt.Errorf("reservation_time is required")
	}
	
	if r.PartySize == 0 {
		return fmt.Errorf("party_size is required")
	}
	
	if r.CustomerName == "" {
		return fmt.Errorf("customer_name is required")
	}
	

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// UpdateReservationsRequest represents a request to update a reservations
type UpdateReservationsRequest struct {
	
	LocationId *uuid.UUID `json:"location_id,omitempty" validate:"omitempty,required"`
	
	TableId *uuid.UUID `json:"table_id,omitempty"`
	
	CustomerId *uuid.UUID `json:"customer_id,omitempty"`
	
	ReservationNumber *string `json:"reservation_number,omitempty" validate:"omitempty,required"`
	
	ReservationDate *time.Time `json:"reservation_date,omitempty" validate:"omitempty,required"`
	
	ReservationTime *string `json:"reservation_time,omitempty" validate:"omitempty,required"`
	
	DurationMinutes *int64 `json:"duration_minutes,omitempty"`
	
	PartySize *int64 `json:"party_size,omitempty" validate:"omitempty,required"`
	
	CustomerName *string `json:"customer_name,omitempty" validate:"omitempty,required"`
	
	CustomerPhone *string `json:"customer_phone,omitempty" validate:"omitempty,e164"`
	
	CustomerEmail *string `json:"customer_email,omitempty" validate:"omitempty,email"`
	
	Status *string `json:"status,omitempty"`
	
	AssignedWaiterId *uuid.UUID `json:"assigned_waiter_id,omitempty"`
	
	AssignedAt *time.Time `json:"assigned_at,omitempty"`
	
	SeatedAt *time.Time `json:"seated_at,omitempty"`
	
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	SpecialRequests *string `json:"special_requests,omitempty"`
	
	Occasion *string `json:"occasion,omitempty"`
	
	DietaryRestrictions *string `json:"dietary_restrictions,omitempty"`
	
	ConfirmationCode *string `json:"confirmation_code,omitempty"`
	
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	
	ConfirmedBy *uuid.UUID `json:"confirmed_by,omitempty"`
	
	ReminderSentAt *time.Time `json:"reminder_sent_at,omitempty"`
	
	NotificationPreferences *json.RawMessage `json:"notification_preferences,omitempty"`
	
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`
	
	CancelledBy *uuid.UUID `json:"cancelled_by,omitempty"`
	
	CancellationReason *string `json:"cancellation_reason,omitempty"`
	
	Notes *string `json:"notes,omitempty"`
	
	Metadata *json.RawMessage `json:"metadata,omitempty"`
	
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
	
}

// Validate validates the update request
func (r *UpdateReservationsRequest) Validate() error {
	// At least one field must be provided
	hasUpdate := false
	
	if r.LocationId != nil {
		hasUpdate = true
	}
	
	if r.TableId != nil {
		hasUpdate = true
	}
	
	if r.CustomerId != nil {
		hasUpdate = true
	}
	
	if r.ReservationNumber != nil {
		hasUpdate = true
	}
	
	if r.ReservationDate != nil {
		hasUpdate = true
	}
	
	if r.ReservationTime != nil {
		hasUpdate = true
	}
	
	if r.DurationMinutes != nil {
		hasUpdate = true
	}
	
	if r.PartySize != nil {
		hasUpdate = true
	}
	
	if r.CustomerName != nil {
		hasUpdate = true
	}
	
	if r.CustomerPhone != nil {
		hasUpdate = true
	}
	
	if r.CustomerEmail != nil {
		hasUpdate = true
	}
	
	if r.Status != nil {
		hasUpdate = true
	}
	
	if r.AssignedWaiterId != nil {
		hasUpdate = true
	}
	
	if r.AssignedAt != nil {
		hasUpdate = true
	}
	
	if r.SeatedAt != nil {
		hasUpdate = true
	}
	
	if r.CompletedAt != nil {
		hasUpdate = true
	}
	
	if r.SpecialRequests != nil {
		hasUpdate = true
	}
	
	if r.Occasion != nil {
		hasUpdate = true
	}
	
	if r.DietaryRestrictions != nil {
		hasUpdate = true
	}
	
	if r.ConfirmationCode != nil {
		hasUpdate = true
	}
	
	if r.ConfirmedAt != nil {
		hasUpdate = true
	}
	
	if r.ConfirmedBy != nil {
		hasUpdate = true
	}
	
	if r.ReminderSentAt != nil {
		hasUpdate = true
	}
	
	if r.NotificationPreferences != nil {
		hasUpdate = true
	}
	
	if r.CancelledAt != nil {
		hasUpdate = true
	}
	
	if r.CancelledBy != nil {
		hasUpdate = true
	}
	
	if r.CancellationReason != nil {
		hasUpdate = true
	}
	
	if r.Notes != nil {
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
	

	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Additional validation
	// TODO: Add business-specific validation rules

	return nil
}

// ReservationsListResponse represents a paginated list of reservations records
type ReservationsListResponse struct {
	Items      []*ReservationsResponse `json:"items"`
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
