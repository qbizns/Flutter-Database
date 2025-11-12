package order

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

// Validator handles Orders validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new Orders validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateOrdersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate OrderNumber
	
	if err := v.validateOrderNumber(req.OrderNumber); err != nil {
		return err
	}
	
	
	
	// Validate DisplayNumber
	
	
	
	// Validate OrderType
	
	
	
	// Validate TableId
	
	
	if err := v.validateTableIdExists(ctx, tx, req.TableId); err != nil {
		return err
	}
	
	
	// Validate ReservationId
	
	
	if err := v.validateReservationIdExists(ctx, tx, req.ReservationId); err != nil {
		return err
	}
	
	
	// Validate Covers
	
	
	
	// Validate CustomerId
	
	
	if err := v.validateCustomerIdExists(ctx, tx, req.CustomerId); err != nil {
		return err
	}
	
	
	// Validate WaiterId
	
	
	if err := v.validateWaiterIdExists(ctx, tx, req.WaiterId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate OrderDate
	
	
	
	// Validate SubmittedAt
	
	
	
	// Validate KitchenReceivedAt
	
	
	
	// Validate ReadyAt
	
	
	
	// Validate ServedAt
	
	
	
	// Validate CompletedAt
	
	
	
	// Validate Subtotal
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate DiscountAmount
	
	
	
	// Validate ServiceCharge
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate SaleId
	
	
	if err := v.validateSaleIdExists(ctx, tx, req.SaleId); err != nil {
		return err
	}
	
	
	// Validate ShiftId
	
	
	if err := v.validateShiftIdExists(ctx, tx, req.ShiftId); err != nil {
		return err
	}
	
	
	// Validate CustomerNotes
	
	
	
	// Validate KitchenNotes
	
	
	
	// Validate InternalNotes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate 'draft',
	
	
	
	// Validate 'served',
	
	
	
	// Validate 'dineIn',
	
	
	
	// Validate Subtotal
	
	
	
	// Validate DiscountAmount
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateOrdersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate LocationId if provided
	
	
	if req.LocationId != nil {
		if err := v.validateLocationIdExists(ctx, tx, *req.LocationId); err != nil {
			return err
		}
	}
	
	
	// Validate OrderNumber if provided
	
	if req.OrderNumber != nil {
		if err := v.validateOrderNumber(*req.OrderNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate DisplayNumber if provided
	
	
	
	// Validate OrderType if provided
	
	
	
	// Validate TableId if provided
	
	
	if req.TableId != nil {
		if err := v.validateTableIdExists(ctx, tx, *req.TableId); err != nil {
			return err
		}
	}
	
	
	// Validate ReservationId if provided
	
	
	if req.ReservationId != nil {
		if err := v.validateReservationIdExists(ctx, tx, *req.ReservationId); err != nil {
			return err
		}
	}
	
	
	// Validate Covers if provided
	
	
	
	// Validate CustomerId if provided
	
	
	if req.CustomerId != nil {
		if err := v.validateCustomerIdExists(ctx, tx, *req.CustomerId); err != nil {
			return err
		}
	}
	
	
	// Validate WaiterId if provided
	
	
	if req.WaiterId != nil {
		if err := v.validateWaiterIdExists(ctx, tx, *req.WaiterId); err != nil {
			return err
		}
	}
	
	
	// Validate Status if provided
	
	
	
	// Validate OrderDate if provided
	
	
	
	// Validate SubmittedAt if provided
	
	
	
	// Validate KitchenReceivedAt if provided
	
	
	
	// Validate ReadyAt if provided
	
	
	
	// Validate ServedAt if provided
	
	
	
	// Validate CompletedAt if provided
	
	
	
	// Validate Subtotal if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	
	// Validate ServiceCharge if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate SaleId if provided
	
	
	if req.SaleId != nil {
		if err := v.validateSaleIdExists(ctx, tx, *req.SaleId); err != nil {
			return err
		}
	}
	
	
	// Validate ShiftId if provided
	
	
	if req.ShiftId != nil {
		if err := v.validateShiftIdExists(ctx, tx, *req.ShiftId); err != nil {
			return err
		}
	}
	
	
	// Validate CustomerNotes if provided
	
	
	
	// Validate KitchenNotes if provided
	
	
	
	// Validate InternalNotes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate 'draft', if provided
	
	
	
	// Validate 'served', if provided
	
	
	
	// Validate 'dineIn', if provided
	
	
	
	// Validate Subtotal if provided
	
	
	
	// Validate DiscountAmount if provided
	
	
	

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





// validateLocationIdExists validates that location_id exists
func (v *Validator) validateLocationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for location
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM location WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check location existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("location with id %s does not exist", id)
	}
	return nil
}



// validateOrderNumber validates order_number field
func (v *Validator) validateOrderNumber(value string) error {
	
	// Add custom validation for order_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("order_number cannot be empty")
	}
	
	return nil
}















// validateTableIdExists validates that table_id exists
func (v *Validator) validateTableIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for table
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM table WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("table with id %s does not exist", id)
	}
	return nil
}





// validateReservationIdExists validates that reservation_id exists
func (v *Validator) validateReservationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for reservation
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM reservation WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check reservation existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("reservation with id %s does not exist", id)
	}
	return nil
}









// validateCustomerIdExists validates that customer_id exists
func (v *Validator) validateCustomerIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for customer
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM customer WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check customer existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("customer with id %s does not exist", id)
	}
	return nil
}





// validateWaiterIdExists validates that waiter_id exists
func (v *Validator) validateWaiterIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for waiter
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM waiter WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check waiter existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("waiter with id %s does not exist", id)
	}
	return nil
}





















































// validateSaleIdExists validates that sale_id exists
func (v *Validator) validateSaleIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for sale
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM sale WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check sale existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("sale with id %s does not exist", id)
	}
	return nil
}





// validateShiftIdExists validates that shift_id exists
func (v *Validator) validateShiftIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for shift
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM shift WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check shift existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("shift with id %s does not exist", id)
	}
	return nil
}















































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateOrdersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateOrdersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *Orders, req *dto.UpdateOrdersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *Orders) error {
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
