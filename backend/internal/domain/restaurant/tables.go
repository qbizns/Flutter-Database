package restaurant

import (
	"time"

	"github.com/google/uuid"
)

// Table-related business logic methods

// IsAvailable checks if table is available for new reservations
func (t *RestaurantTable) IsAvailable() bool {
	return t.Status == "available" && t.IsActive
}

// CanAccommodate checks if table can accommodate party size
func (t *RestaurantTable) CanAccommodate(partySize int) bool {
	return partySize >= t.MinCapacity && partySize <= t.MaxCapacity
}

// Seat marks table as occupied with covers
func (t *RestaurantTable) Seat(covers int, waiterID *uuid.UUID) {
	t.Status = "occupied"
	t.CurrentCovers = covers
	t.SeatedAt = &time.Time{}
	*t.SeatedAt = time.Now()
	t.CurrentWaiterID = waiterID
	t.UpdatedAt = time.Now()
}

// Unseat marks table as available
func (t *RestaurantTable) Unseat() {
	t.Status = "available"
	t.CurrentCovers = 0
	t.SeatedAt = nil
	t.CurrentWaiterID = nil
	t.UpdatedAt = time.Now()
}

// StartCleaning marks table for cleaning
func (t *RestaurantTable) StartCleaning() {
	t.Status = "cleaning"
	t.CurrentCovers = 0
	t.UpdatedAt = time.Now()
}

// FinishCleaning marks table as available after cleaning
func (t *RestaurantTable) FinishCleaning() {
	t.Status = "available"
	t.UpdatedAt = time.Now()
}

// Reservation-related business logic methods

// CanBeSeat checks if reservation can be seated
func (r *Reservation) CanBeSeated() bool {
	return r.Status == "confirmed"
}

// Seat marks reservation as seated
func (r *Reservation) Seat() {
	r.Status = "seated"
	now := time.Now()
	r.SeatedAt = &now
	r.UpdatedAt = now
}

// Complete marks reservation as completed
func (r *Reservation) Complete() {
	r.Status = "completed"
	now := time.Now()
	r.CompletedAt = &now
	r.UpdatedAt = now
}

// Cancel marks reservation as cancelled
func (r *Reservation) Cancel(cancelledBy *uuid.UUID, reason string) {
	r.Status = "cancelled"
	now := time.Now()
	r.CancelledAt = &now
	r.CancelledBy = cancelledBy
	r.CancellationReason = &reason
	r.UpdatedAt = now
}

// Confirm marks reservation as confirmed
func (r *Reservation) Confirm(confirmedBy *uuid.UUID) {
	r.Status = "confirmed"
	now := time.Now()
	r.ConfirmedAt = &now
	r.ConfirmedBy = confirmedBy
	r.UpdatedAt = now
}

// IsUpcoming checks if reservation is within next hour
func (r *Reservation) IsUpcoming() bool {
	resTime := time.Date(
		r.ReservationDate.Year(),
		r.ReservationDate.Month(),
		r.ReservationDate.Day(),
		0, 0, 0, 0,
		r.ReservationDate.Location(),
	)

	// Parse reservation_time (TIME format: HH:MM:SS)
	if len(r.ReservationTime) > 0 {
		// Simple parsing for TIME format
		// In real implementation, parse properly
		resTime = resTime.Add(1 * time.Hour)
	}

	now := time.Now()
	return now.Before(resTime.Add(1 * time.Hour)) && now.After(resTime.Add(-2*time.Hour))
}

// IsNoShow marks reservation as no-show
func (r *Reservation) MarkNoShow() {
	r.Status = "no_show"
	r.UpdatedAt = time.Now()
}

// SendReminder marks reminder as sent
func (r *Reservation) SendReminder() {
	now := time.Now()
	r.ReminderSentAt = &now
	r.UpdatedAt = now
}

// DurationUntilReservation calculates time until reservation
func (r *Reservation) DurationUntilReservation() time.Duration {
	resDateTime := time.Date(
		r.ReservationDate.Year(),
		r.ReservationDate.Month(),
		r.ReservationDate.Day(),
		0, 0, 0, 0,
		time.Local,
	)
	return resDateTime.Sub(time.Now())
}
