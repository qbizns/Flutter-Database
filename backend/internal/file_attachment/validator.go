package file_attachment

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

// Validator handles FileAttachments validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new FileAttachments validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateFileAttachmentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate FileName
	
	if err := v.validateFileName(req.FileName); err != nil {
		return err
	}
	
	
	
	// Validate FileSize
	
	
	
	// Validate MimeType
	
	if err := v.validateMimeType(req.MimeType); err != nil {
		return err
	}
	
	
	
	// Validate FileExtension
	
	
	
	// Validate StorageProvider
	
	
	
	// Validate StoragePath
	
	if err := v.validateStoragePath(req.StoragePath); err != nil {
		return err
	}
	
	
	
	// Validate StorageUrl
	
	
	
	// Validate FileHash
	
	
	
	// Validate EntityType
	
	if err := v.validateEntityType(req.EntityType); err != nil {
		return err
	}
	
	
	
	// Validate EntityId
	
	
	if err := v.validateEntityIdExists(ctx, tx, req.EntityId); err != nil {
		return err
	}
	
	
	// Validate Description
	
	
	
	// Validate Tags
	
	
	
	// Validate IsPublic
	
	
	
	// Validate ImageWidth
	
	
	
	// Validate ImageHeight
	
	
	
	// Validate VirusScanStatus
	
	
	
	// Validate VirusScanAt
	
	
	
	// Validate UploadedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateFileAttachmentsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate FileName if provided
	
	if req.FileName != nil {
		if err := v.validateFileName(*req.FileName); err != nil {
			return err
		}
	}
	
	
	
	// Validate FileSize if provided
	
	
	
	// Validate MimeType if provided
	
	if req.MimeType != nil {
		if err := v.validateMimeType(*req.MimeType); err != nil {
			return err
		}
	}
	
	
	
	// Validate FileExtension if provided
	
	
	
	// Validate StorageProvider if provided
	
	
	
	// Validate StoragePath if provided
	
	if req.StoragePath != nil {
		if err := v.validateStoragePath(*req.StoragePath); err != nil {
			return err
		}
	}
	
	
	
	// Validate StorageUrl if provided
	
	
	
	// Validate FileHash if provided
	
	
	
	// Validate EntityType if provided
	
	if req.EntityType != nil {
		if err := v.validateEntityType(*req.EntityType); err != nil {
			return err
		}
	}
	
	
	
	// Validate EntityId if provided
	
	
	if req.EntityId != nil {
		if err := v.validateEntityIdExists(ctx, tx, *req.EntityId); err != nil {
			return err
		}
	}
	
	
	// Validate Description if provided
	
	
	
	// Validate Tags if provided
	
	
	
	// Validate IsPublic if provided
	
	
	
	// Validate ImageWidth if provided
	
	
	
	// Validate ImageHeight if provided
	
	
	
	// Validate VirusScanStatus if provided
	
	
	
	// Validate VirusScanAt if provided
	
	
	
	// Validate UploadedBy if provided
	
	
	

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



// validateFileName validates file_name field
func (v *Validator) validateFileName(value string) error {
	
	// Add custom validation for file_name
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("file_name cannot be empty")
	}
	
	return nil
}









// validateMimeType validates mime_type field
func (v *Validator) validateMimeType(value string) error {
	
	// Add custom validation for mime_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("mime_type cannot be empty")
	}
	
	return nil
}













// validateStoragePath validates storage_path field
func (v *Validator) validateStoragePath(value string) error {
	
	// Add custom validation for storage_path
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("storage_path cannot be empty")
	}
	
	return nil
}













// validateEntityType validates entity_type field
func (v *Validator) validateEntityType(value string) error {
	
	// Add custom validation for entity_type
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("entity_type cannot be empty")
	}
	
	return nil
}







// validateEntityIdExists validates that entity_id exists
func (v *Validator) validateEntityIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for entity
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM entity WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check entity existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("entity with id %s does not exist", id)
	}
	return nil
}



































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateFileAttachmentsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateFileAttachmentsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *FileAttachments, req *dto.UpdateFileAttachmentsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *FileAttachments) error {
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
