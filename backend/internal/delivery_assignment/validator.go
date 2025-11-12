package delivery_assignment

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto"
)

// Validator handles DeliveryAssignments validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new DeliveryAssignments validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryAssignmentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate OrderId
	
	
	if err := v.validateOrderIdExists(ctx, tx, req.OrderId); err != nil {
		return err
	}
	
	
	// Validate DriverId
	
	
	if err := v.validateDriverIdExists(ctx, tx, req.DriverId); err != nil {
		return err
	}
	
	
	// Validate DriverShiftId
	
	
	if err := v.validateDriverShiftIdExists(ctx, tx, req.DriverShiftId); err != nil {
		return err
	}
	
	
	// Validate DeliveryZoneId
	
	
	if err := v.validateDeliveryZoneIdExists(ctx, tx, req.DeliveryZoneId); err != nil {
		return err
	}
	
	
	// Validate CustomerAddressId
	
	
	if err := v.validateCustomerAddressIdExists(ctx, tx, req.CustomerAddressId); err != nil {
		return err
	}
	
	
	// Validate DeliveryAddress
	
	if err := v.validateDeliveryAddress(req.DeliveryAddress); err != nil {
		return err
	}
	
	
	
	// Validate DeliveryLocation
	
	
	
	// Validate AssignedAt
	
	
	
	// Validate AssignedBy
	
	
	
	// Validate Status
	
	
	
	// Validate AcceptedAt
	
	
	
	// Validate PickedUpAt
	
	
	
	// Validate DispatchedAt
	
	
	
	// Validate ArrivedAt
	
	
	
	// Validate DeliveredAt
	
	
	
	// Validate FailedAt
	
	
	
	// Validate EstimatedPickupTime
	
	
	
	// Validate EstimatedDeliveryTime
	
	
	
	// Validate DistanceKm
	
	
	
	// Validate RouteInfo
	
	
	
	// Validate DeliveryFee
	
	
	
	// Validate DriverCommission
	
	
	
	// Validate PaymentMethod
	
	
	
	// Validate CashCollected
	
	
	
	// Validate SignatureImageUrl
	
	
	
	// Validate DeliveryPhotoUrl
	
	
	
	// Validate RecipientName
	
	
	
	// Validate DeliveryNotes
	
	
	
	// Validate FailureReason
	
	
	
	// Validate FailureNotes
	
	
	
	// Validate RetryCount
	
	
	
	// Validate CustomerRating
	
	
	
	// Validate CustomerFeedback
	
	
	
	// Validate DriverNotes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'assigned',
	
	
	
	// Validate 'delivered',
	
	
	
	// Validate CustomerRating
	
	
	
	// Validate DeliveryFee
	
	
	

	// Cross-field validation
	if err := v.validateCrossFields(ctx, tx, req); err != nil {
		return err
	}

	// Business rules validation
	if err := v.validateBusinessRules(ctx, tx, req); err != nil {
		return err
	}

	return nil
}

// ValidateUpdate validates an update request
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateDeliveryAssignmentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate OrderId if provided
	
	
	if req.OrderId != nil {
		if err := v.validateOrderIdExists(ctx, tx, *req.OrderId); err != nil {
			return err
		}
	}
	
	
	// Validate DriverId if provided
	
	
	if req.DriverId != nil {
		if err := v.validateDriverIdExists(ctx, tx, *req.DriverId); err != nil {
			return err
		}
	}
	
	
	// Validate DriverShiftId if provided
	
	
	if req.DriverShiftId != nil {
		if err := v.validateDriverShiftIdExists(ctx, tx, *req.DriverShiftId); err != nil {
			return err
		}
	}
	
	
	// Validate DeliveryZoneId if provided
	
	
	if req.DeliveryZoneId != nil {
		if err := v.validateDeliveryZoneIdExists(ctx, tx, *req.DeliveryZoneId); err != nil {
			return err
		}
	}
	
	
	// Validate CustomerAddressId if provided
	
	
	if req.CustomerAddressId != nil {
		if err := v.validateCustomerAddressIdExists(ctx, tx, *req.CustomerAddressId); err != nil {
			return err
		}
	}
	
	
	// Validate DeliveryAddress if provided
	
	if req.DeliveryAddress != nil {
		if err := v.validateDeliveryAddress(*req.DeliveryAddress); err != nil {
			return err
		}
	}
	
	
	
	// Validate DeliveryLocation if provided
	
	
	
	// Validate AssignedAt if provided
	
	
	
	// Validate AssignedBy if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate AcceptedAt if provided
	
	
	
	// Validate PickedUpAt if provided
	
	
	
	// Validate DispatchedAt if provided
	
	
	
	// Validate ArrivedAt if provided
	
	
	
	// Validate DeliveredAt if provided
	
	
	
	// Validate FailedAt if provided
	
	
	
	// Validate EstimatedPickupTime if provided
	
	
	
	// Validate EstimatedDeliveryTime if provided
	
	
	
	// Validate DistanceKm if provided
	
	
	
	// Validate RouteInfo if provided
	
	
	
	// Validate DeliveryFee if provided
	
	
	
	// Validate DriverCommission if provided
	
	
	
	// Validate PaymentMethod if provided
	
	
	
	// Validate CashCollected if provided
	
	
	
	// Validate SignatureImageUrl if provided
	
	
	
	// Validate DeliveryPhotoUrl if provided
	
	
	
	// Validate RecipientName if provided
	
	
	
	// Validate DeliveryNotes if provided
	
	
	
	// Validate FailureReason if provided
	
	
	
	// Validate FailureNotes if provided
	
	
	
	// Validate RetryCount if provided
	
	
	
	// Validate CustomerRating if provided
	
	
	
	// Validate CustomerFeedback if provided
	
	
	
	// Validate DriverNotes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'assigned', if provided
	
	
	
	// Validate 'delivered', if provided
	
	
	
	// Validate CustomerRating if provided
	
	
	
	// Validate DeliveryFee if provided
	
	
	

	// Business rules validation
	if err := v.validateUpdateBusinessRules(ctx, tx, existing, req); err != nil {
		return err
	}

	return nil
}

// ValidateDelete validates a delete request
func (v *Validator) ValidateDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	// Check if entity can be deleted (no foreign key constraints)
	if err := v.validateCanDelete(ctx, tx, existing); err != nil {
		return err
	}

	return nil
}





// validateOrderIdExists validates that order_id exists
func (v *Validator) validateOrderIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for order
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM order WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("order with id %s does not exist", id)
	}
	return nil
}





// validateDriverIdExists validates that driver_id exists
func (v *Validator) validateDriverIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for driver
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM driver WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check driver existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("driver with id %s does not exist", id)
	}
	return nil
}





// validateDriverShiftIdExists validates that driver_shift_id exists
func (v *Validator) validateDriverShiftIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for driver_shift
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM driver_shift WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check driver_shift existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("driver_shift with id %s does not exist", id)
	}
	return nil
}





// validateDeliveryZoneIdExists validates that delivery_zone_id exists
func (v *Validator) validateDeliveryZoneIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for delivery_zone
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM delivery_zone WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check delivery_zone existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("delivery_zone with id %s does not exist", id)
	}
	return nil
}





// validateCustomerAddressIdExists validates that customer_address_id exists
func (v *Validator) validateCustomerAddressIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for customer_address
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer_address WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check customer_address existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("customer_address with id %s does not exist", id)
	}
	return nil
}



// validateDeliveryAddress validates delivery_address field
func (v *Validator) validateDeliveryAddress(value string) error {
	
	// Add custom validation for delivery_address
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("delivery_address cannot be empty")
	}
	
	return nil
}

















































































































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryAssignmentsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateDeliveryAssignmentsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *DeliveryAssignments, req *dto.UpdateDeliveryAssignmentsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *DeliveryAssignments) error {
	// Add delete validation here
	// Example: check for dependent records in other tables
	// Example: prevent deletion of active/in-use entities
	return nil
}

// Helper validation functions

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`) // E.164 format
	urlRegex   = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// isValidPhone validates phone number format (E.164)
func isValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// isValidURL validates URL format
func isValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

// isValidUUID validates UUID format
func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// isValidDateRange validates date range
func isValidDateRange(start, end time.Time) bool {
	return start.Before(end)
}

// isPositive validates positive numbers
func isPositive(value float64) bool {
	return value > 0
}

// isNonNegative validates non-negative numbers
func isNonNegative(value float64) bool {
	return value >= 0
}

// isWithinRange validates value is within range
func isWithinRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// isValidLength validates string length
func isValidLength(value string, min, max int) bool {
	length := len(value)
	return length >= min && length <= max
}
