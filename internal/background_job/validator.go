package background_job

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

// Validator handles BackgroundJobs validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new BackgroundJobs validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateBackgroundJobsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate JobType
	
	if err := v.validateJobType(req.JobType); err != nil {
		return err
	}
	
	
	
	// Validate JobName
	
	if err := v.validateJobName(req.JobName); err != nil {
		return err
	}
	
	
	
	// Validate QueueName
	
	
	
	// Validate Status
	
	
	
	// Validate Status
	
	
	
	// Validate Payload
	
	
	
	// Validate Result
	
	
	
	// Validate ErrorMessage
	
	
	
	// Validate ErrorDetails
	
	
	
	// Validate Attempts
	
	
	
	// Validate MaxAttempts
	
	
	
	// Validate Priority
	
	
	
	// Validate ScheduledAt
	
	
	
	// Validate StartedAt
	
	
	
	// Validate CompletedAt
	
	
	
	// Validate FailedAt
	
	
	
	// Validate WorkerId
	
	
	if err := v.validateWorkerIdExists(ctx, tx, req.WorkerId); err != nil {
		return err
	}
	
	
	// Validate ProcessingTimeout
	
	
	
	// Validate CreatedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateBackgroundJobsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate JobType if provided
	
	if req.JobType != nil {
		if err := v.validateJobType(*req.JobType); err != nil {
			return err
		}
	}
	
	
	
	// Validate JobName if provided
	
	if req.JobName != nil {
		if err := v.validateJobName(*req.JobName); err != nil {
			return err
		}
	}
	
	
	
	// Validate QueueName if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate Payload if provided
	
	
	
	// Validate Result if provided
	
	
	
	// Validate ErrorMessage if provided
	
	
	
	// Validate ErrorDetails if provided
	
	
	
	// Validate Attempts if provided
	
	
	
	// Validate MaxAttempts if provided
	
	
	
	// Validate Priority if provided
	
	
	
	// Validate ScheduledAt if provided
	
	
	
	// Validate StartedAt if provided
	
	
	
	// Validate CompletedAt if provided
	
	
	
	// Validate FailedAt if provided
	
	
	
	// Validate WorkerId if provided
	
	
	if req.WorkerId != nil {
		if err := v.validateWorkerIdExists(ctx, tx, *req.WorkerId); err != nil {
			return err
		}
	}
	
	
	// Validate ProcessingTimeout if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	

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



// validateJobType validates job_type field
func (v *Validator) validateJobType(value string) error {
	
	// Add custom validation for job_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("job_type cannot be empty")
	}
	
	return nil
}





// validateJobName validates job_name field
func (v *Validator) validateJobName(value string) error {
	
	// Add custom validation for job_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("job_name cannot be empty")
	}
	
	return nil
}































































// validateWorkerIdExists validates that worker_id exists
func (v *Validator) validateWorkerIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for worker
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM worker WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check worker existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("worker with id %s does not exist", id)
	}
	return nil
}











// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateBackgroundJobsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateBackgroundJobsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *BackgroundJobs, req *dto.UpdateBackgroundJobsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *BackgroundJobs) error {
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
