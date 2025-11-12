package api_request_log

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

// Validator handles ApiRequestLogs validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ApiRequestLogs validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateApiRequestLogsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate RequestId
	
	
	if err := v.validateRequestIdExists(ctx, tx, req.RequestId); err != nil {
		return err
	}
	
	
	// Validate Method
	
	if err := v.validateMethod(req.Method); err != nil {
		return err
	}
	
	
	
	// Validate Path
	
	if err := v.validatePath(req.Path); err != nil {
		return err
	}
	
	
	
	// Validate QueryParams
	
	
	
	// Validate UserId
	
	
	if err := v.validateUserIdExists(ctx, tx, req.UserId); err != nil {
		return err
	}
	
	
	// Validate ApiKeyId
	
	
	if err := v.validateApiKeyIdExists(ctx, tx, req.ApiKeyId); err != nil {
		return err
	}
	
	
	// Validate RequestHeaders
	
	
	
	// Validate RequestBody
	
	
	
	// Validate IpAddress
	
	
	
	// Validate UserAgent
	
	
	
	// Validate StatusCode
	
	
	
	// Validate ResponseHeaders
	
	
	
	// Validate ResponseBody
	
	
	
	// Validate DurationMs
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate ErrorStack
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateApiRequestLogsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate RequestId if provided
	
	
	if req.RequestId != nil {
		if err := v.validateRequestIdExists(ctx, tx, *req.RequestId); err != nil {
			return err
		}
	}
	
	
	// Validate Method if provided
	
	if req.Method != nil {
		if err := v.validateMethod(*req.Method); err != nil {
			return err
		}
	}
	
	
	
	// Validate Path if provided
	
	if req.Path != nil {
		if err := v.validatePath(*req.Path); err != nil {
			return err
		}
	}
	
	
	
	// Validate QueryParams if provided
	
	
	
	// Validate UserId if provided
	
	
	if req.UserId != nil {
		if err := v.validateUserIdExists(ctx, tx, *req.UserId); err != nil {
			return err
		}
	}
	
	
	// Validate ApiKeyId if provided
	
	
	if req.ApiKeyId != nil {
		if err := v.validateApiKeyIdExists(ctx, tx, *req.ApiKeyId); err != nil {
			return err
		}
	}
	
	
	// Validate RequestHeaders if provided
	
	
	
	// Validate RequestBody if provided
	
	
	
	// Validate IpAddress if provided
	
	
	
	// Validate UserAgent if provided
	
	
	
	// Validate StatusCode if provided
	
	
	
	// Validate ResponseHeaders if provided
	
	
	
	// Validate ResponseBody if provided
	
	
	
	// Validate DurationMs if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate ErrorStack if provided
	
	
	

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





// validateRequestIdExists validates that request_id exists
func (v *Validator) validateRequestIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for request
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM request WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check request existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("request with id %s does not exist", id)
	}
	return nil
}



// validateMethod validates method field
func (v *Validator) validateMethod(value string) error {
	
	// Add custom validation for method
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("method cannot be empty")
	}
	
	return nil
}





// validatePath validates path field
func (v *Validator) validatePath(value string) error {
	
	// Add custom validation for path
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("path cannot be empty")
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





// validateApiKeyIdExists validates that api_key_id exists
func (v *Validator) validateApiKeyIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for api_key
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM api_key WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check api_key existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("api_key with id %s does not exist", id)
	}
	return nil
}











































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateApiRequestLogsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateApiRequestLogsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ApiRequestLogs, req *dto.UpdateApiRequestLogsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ApiRequestLogs) error {
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
