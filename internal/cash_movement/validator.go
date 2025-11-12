package cash_movement

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

// Validator handles CashMovements validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new CashMovements validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateCashMovementsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate PosSessionId
	
	
	if err := v.validatePosSessionIdExists(ctx, tx, req.PosSessionId); err != nil {
		return err
	}
	
	
	// Validate CashDrawerId
	
	
	if err := v.validateCashDrawerIdExists(ctx, tx, req.CashDrawerId); err != nil {
		return err
	}
	
	
	// Validate MovementType
	
	if err := v.validateMovementType(req.MovementType); err != nil {
		return err
	}
	
	
	
	// Validate Amount
	
	
	
	// Validate ReasonCode
	
	
	
	// Validate ReasonDescription
	
	if err := v.validateReasonDescription(req.ReasonDescription); err != nil {
		return err
	}
	
	
	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate RequiresApproval
	
	
	
	// Validate ApprovedBy
	
	
	
	// Validate ApprovedAt
	
	
	
	// Validate Notes
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateCashMovementsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate PosSessionId if provided
	
	
	if req.PosSessionId != nil {
		if err := v.validatePosSessionIdExists(ctx, tx, *req.PosSessionId); err != nil {
			return err
		}
	}
	
	
	// Validate CashDrawerId if provided
	
	
	if req.CashDrawerId != nil {
		if err := v.validateCashDrawerIdExists(ctx, tx, *req.CashDrawerId); err != nil {
			return err
		}
	}
	
	
	// Validate MovementType if provided
	
	if req.MovementType != nil {
		if err := v.validateMovementType(*req.MovementType); err != nil {
			return err
		}
	}
	
	
	
	// Validate Amount if provided
	
	
	
	// Validate ReasonCode if provided
	
	
	
	// Validate ReasonDescription if provided
	
	if req.ReasonDescription != nil {
		if err := v.validateReasonDescription(*req.ReasonDescription); err != nil {
			return err
		}
	}
	
	
	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	
	// Validate RequiresApproval if provided
	
	
	
	// Validate ApprovedBy if provided
	
	
	
	// Validate ApprovedAt if provided
	
	
	
	// Validate Notes if provided
	
	
	

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





// validatePosSessionIdExists validates that pos_session_id exists
func (v *Validator) validatePosSessionIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for pos_session
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pos_session WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check pos_session existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("pos_session with id %s does not exist", id)
	}
	return nil
}





// validateCashDrawerIdExists validates that cash_drawer_id exists
func (v *Validator) validateCashDrawerIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for cash_drawer
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM cash_drawer WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check cash_drawer existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("cash_drawer with id %s does not exist", id)
	}
	return nil
}



// validateMovementType validates movement_type field
func (v *Validator) validateMovementType(value string) error {
	
	// Add custom validation for movement_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("movement_type cannot be empty")
	}
	
	return nil
}













// validateReasonDescription validates reason_description field
func (v *Validator) validateReasonDescription(value string) error {
	
	// Add custom validation for reason_description
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("reason_description cannot be empty")
	}
	
	return nil
}







// validateUserIdExists validates that user_id exists
func (v *Validator) validateUserIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for user
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM user WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("user with id %s does not exist", id)
	}
	return nil
}



















// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateCashMovementsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateCashMovementsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *CashMovements, req *dto.UpdateCashMovementsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *CashMovements) error {
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
