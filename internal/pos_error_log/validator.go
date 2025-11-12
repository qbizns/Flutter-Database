package pos_error_log

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

// Validator handles PosErrorLogs validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new PosErrorLogs validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreatePosErrorLogsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ErrorLevel
	
	if err := v.validateErrorLevel(req.ErrorLevel); err != nil {
		return err
	}
	
	
	
	// Validate ErrorCode
	
	
	
	// Validate ErrorMessage
	
	if err := v.validateErrorMessage(req.ErrorMessage); err != nil {
		return err
	}
	
	
	
	// Validate DeviceId
	
	
	if err := v.validateDeviceIdExists(ctx, tx, req.DeviceId); err != nil {
		return err
	}
	
	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate PosSessionId
	
	
	if err := v.validatePosSessionIdExists(ctx, tx, req.PosSessionId); err != nil {
		return err
	}
	
	
	// Validate SaleId
	
	
	if err := v.validateSaleIdExists(ctx, tx, req.SaleId); err != nil {
		return err
	}
	
	
	// Validate StackTrace
	
	
	
	// Validate RequestData
	
	
	
	// Validate ErrorData
	
	
	
	// Validate IsResolved
	
	
	
	// Validate ResolvedBy
	
	
	
	// Validate ResolvedAt
	
	
	
	// Validate ResolutionNotes
	
	
	
	// Validate OccurredAt
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdatePosErrorLogsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ErrorLevel if provided
	
	if req.ErrorLevel != nil {
		if err := v.validateErrorLevel(*req.ErrorLevel); err != nil {
			return err
		}
	}
	
	
	
	// Validate ErrorCode if provided
	
	
	
	// Validate ErrorMessage if provided
	
	if req.ErrorMessage != nil {
		if err := v.validateErrorMessage(*req.ErrorMessage); err != nil {
			return err
		}
	}
	
	
	
	// Validate DeviceId if provided
	
	
	if req.DeviceId != nil {
		if err := v.validateDeviceIdExists(ctx, tx, *req.DeviceId); err != nil {
			return err
		}
	}
	
	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	
	// Validate PosSessionId if provided
	
	
	if req.PosSessionId != nil {
		if err := v.validatePosSessionIdExists(ctx, tx, *req.PosSessionId); err != nil {
			return err
		}
	}
	
	
	// Validate SaleId if provided
	
	
	if req.SaleId != nil {
		if err := v.validateSaleIdExists(ctx, tx, *req.SaleId); err != nil {
			return err
		}
	}
	
	
	// Validate StackTrace if provided
	
	
	
	// Validate RequestData if provided
	
	
	
	// Validate ErrorData if provided
	
	
	
	// Validate IsResolved if provided
	
	
	
	// Validate ResolvedBy if provided
	
	
	
	// Validate ResolvedAt if provided
	
	
	
	// Validate ResolutionNotes if provided
	
	
	
	// Validate OccurredAt if provided
	
	
	

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



// validateErrorLevel validates error_level field
func (v *Validator) validateErrorLevel(value string) error {
	
	// Add custom validation for error_level
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("error_level cannot be empty")
	}
	
	return nil
}









// validateErrorMessage validates error_message field
func (v *Validator) validateErrorMessage(value string) error {
	
	// Add custom validation for error_message
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("error_message cannot be empty")
	}
	
	return nil
}







// validateDeviceIdExists validates that device_id exists
func (v *Validator) validateDeviceIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for device
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM device WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check device existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("device with id %s does not exist", id)
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



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreatePosErrorLogsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreatePosErrorLogsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *PosErrorLogs, req *dto.UpdatePosErrorLogsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *PosErrorLogs) error {
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
