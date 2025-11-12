package inventory_transfer

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

// Validator handles InventoryTransfers validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new InventoryTransfers validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateInventoryTransfersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate TransferNumber
	
	if err := v.validateTransferNumber(req.TransferNumber); err != nil {
		return err
	}
	
	
	
	// Validate TransferDate
	
	
	
	// Validate FromLocationId
	
	
	if err := v.validateFromLocationIdExists(ctx, tx, req.FromLocationId); err != nil {
		return err
	}
	
	
	// Validate ToLocationId
	
	
	if err := v.validateToLocationIdExists(ctx, tx, req.ToLocationId); err != nil {
		return err
	}
	
	
	// Validate Status
	
	
	
	// Validate RequestedDate
	
	
	
	// Validate ApprovedDate
	
	
	
	// Validate ShippedDate
	
	
	
	// Validate ExpectedDeliveryDate
	
	
	
	// Validate ReceivedDate
	
	
	
	// Validate Carrier
	
	
	
	// Validate TrackingNumber
	
	
	
	// Validate ShippingCost
	
	
	
	// Validate Reason
	
	
	
	// Validate Notes
	
	
	
	// Validate RejectionReason
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	
	// Validate RequestedBy
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ShippedBy
	
	
	
	// Validate ReceivedBy
	
	
	
	// Validate (approvedDate
	
	
	
	// Validate (shippedDate
	
	
	
	// Validate (receivedDate
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateInventoryTransfersRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate TransferNumber if provided
	
	if req.TransferNumber != nil {
		if err := v.validateTransferNumber(*req.TransferNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate TransferDate if provided
	
	
	
	// Validate FromLocationId if provided
	
	
	if req.FromLocationId != nil {
		if err := v.validateFromLocationIdExists(ctx, tx, *req.FromLocationId); err != nil {
			return err
		}
	}
	
	
	// Validate ToLocationId if provided
	
	
	if req.ToLocationId != nil {
		if err := v.validateToLocationIdExists(ctx, tx, *req.ToLocationId); err != nil {
			return err
		}
	}
	
	
	// Validate Status if provided
	
	
	
	// Validate RequestedDate if provided
	
	
	
	// Validate ApprovedDate if provided
	
	
	
	// Validate ShippedDate if provided
	
	
	
	// Validate ExpectedDeliveryDate if provided
	
	
	
	// Validate ReceivedDate if provided
	
	
	
	// Validate Carrier if provided
	
	
	
	// Validate TrackingNumber if provided
	
	
	
	// Validate ShippingCost if provided
	
	
	
	// Validate Reason if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate RejectionReason if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	
	// Validate RequestedBy if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ShippedBy if provided
	
	
	
	// Validate ReceivedBy if provided
	
	
	
	// Validate (approvedDate if provided
	
	
	
	// Validate (shippedDate if provided
	
	
	
	// Validate (receivedDate if provided
	
	
	

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



// validateTransferNumber validates transfer_number field
func (v *Validator) validateTransferNumber(value string) error {
	
	// Add custom validation for transfer_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("transfer_number cannot be empty")
	}
	
	return nil
}











// validateFromLocationIdExists validates that from_location_id exists
func (v *Validator) validateFromLocationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for from_location
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM from_location WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check from_location existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("from_location with id %s does not exist", id)
	}
	return nil
}





// validateToLocationIdExists validates that to_location_id exists
func (v *Validator) validateToLocationIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for to_location
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM to_location WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check to_location existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("to_location with id %s does not exist", id)
	}
	return nil
}



























































































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateInventoryTransfersRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateInventoryTransfersRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *InventoryTransfers, req *dto.UpdateInventoryTransfersRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *InventoryTransfers) error {
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
