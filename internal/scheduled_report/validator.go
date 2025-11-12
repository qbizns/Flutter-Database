package scheduled_report

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

// Validator handles ScheduledReports validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new ScheduledReports validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateScheduledReportsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate ReportName
	
	if err := v.validateReportName(req.ReportName); err != nil {
		return err
	}
	
	
	
	// Validate ReportType
	
	if err := v.validateReportType(req.ReportType); err != nil {
		return err
	}
	
	
	
	// Validate ScheduleFrequency
	
	if err := v.validateScheduleFrequency(req.ScheduleFrequency); err != nil {
		return err
	}
	
	
	
	// Validate ScheduleDayOfWeek
	
	
	
	// Validate ScheduleDayOfMonth
	
	
	
	// Validate ScheduleTime
	
	if err := v.validateScheduleTime(req.ScheduleTime); err != nil {
		return err
	}
	
	
	
	// Validate ScheduleTimezone
	
	
	
	// Validate ReportParameters
	
	
	
	// Validate DeliveryMethod
	
	
	
	// Validate DeliveryRecipients
	
	
	
	// Validate OutputFormat
	
	
	
	// Validate IsActive
	
	
	
	// Validate LastRunAt
	
	
	
	// Validate LastRunStatus
	
	
	
	// Validate NextRunAt
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateScheduledReportsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate ReportName if provided
	
	if req.ReportName != nil {
		if err := v.validateReportName(*req.ReportName); err != nil {
			return err
		}
	}
	
	
	
	// Validate ReportType if provided
	
	if req.ReportType != nil {
		if err := v.validateReportType(*req.ReportType); err != nil {
			return err
		}
	}
	
	
	
	// Validate ScheduleFrequency if provided
	
	if req.ScheduleFrequency != nil {
		if err := v.validateScheduleFrequency(*req.ScheduleFrequency); err != nil {
			return err
		}
	}
	
	
	
	// Validate ScheduleDayOfWeek if provided
	
	
	
	// Validate ScheduleDayOfMonth if provided
	
	
	
	// Validate ScheduleTime if provided
	
	if req.ScheduleTime != nil {
		if err := v.validateScheduleTime(*req.ScheduleTime); err != nil {
			return err
		}
	}
	
	
	
	// Validate ScheduleTimezone if provided
	
	
	
	// Validate ReportParameters if provided
	
	
	
	// Validate DeliveryMethod if provided
	
	
	
	// Validate DeliveryRecipients if provided
	
	
	
	// Validate OutputFormat if provided
	
	
	
	// Validate IsActive if provided
	
	
	
	// Validate LastRunAt if provided
	
	
	
	// Validate LastRunStatus if provided
	
	
	
	// Validate NextRunAt if provided
	
	
	
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



// validateReportName validates report_name field
func (v *Validator) validateReportName(value string) error {
	
	// Add custom validation for report_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("report_name cannot be empty")
	}
	
	return nil
}





// validateReportType validates report_type field
func (v *Validator) validateReportType(value string) error {
	
	// Add custom validation for report_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("report_type cannot be empty")
	}
	
	return nil
}





// validateScheduleFrequency validates schedule_frequency field
func (v *Validator) validateScheduleFrequency(value string) error {
	
	// Add custom validation for schedule_frequency
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("schedule_frequency cannot be empty")
	}
	
	return nil
}













// validateScheduleTime validates schedule_time field
func (v *Validator) validateScheduleTime(value string) error {
	
	// Add custom validation for schedule_time
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("schedule_time cannot be empty")
	}
	
	return nil
}













































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateScheduledReportsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateScheduledReportsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *ScheduledReports, req *dto.UpdateScheduledReportsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *ScheduledReports) error {
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
