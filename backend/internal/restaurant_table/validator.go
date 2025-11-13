package restaurant_table

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	
)

// Validator handles RestaurantTables validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new RestaurantTables validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *CreateRestaurantTablesRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate LocationId
	
	
	if err := v.validateLocationIdExists(ctx, tx, req.LocationId); err != nil {
		return err
	}
	
	
	// Validate FloorPlanId if provided
	if req.FloorPlanId != nil {
		if err := v.validateFloorPlanIdExists(ctx, tx, *req.FloorPlanId); err != nil {
			return err
		}
	}


	// Validate SectionId if provided
	if req.SectionId != nil {
		if err := v.validateSectionIdExists(ctx, tx, *req.SectionId); err != nil {
			return err
		}
	}
	
	
	// Validate TableNumber
	
	if err := v.validateTableNumber(req.TableNumber); err != nil {
		return err
	}
	
	
	
	// Validate TableName
	
	
	
	// Validate MinCapacity
	
	
	
	// Validate MaxCapacity
	
	
	
	// Validate TableShape
	
	
	
	// Validate IsCombinable
	
	
	
	// Validate PositionX
	
	
	
	// Validate PositionY
	
	
	
	// Validate Rotation
	
	
	
	// Validate Status
	
	
	
	// Validate CurrentCovers
	
	
	
	// Validate SeatedAt
	
	
	
	// Validate CurrentWaiterId if provided
	if req.CurrentWaiterId != nil {
		if err := v.validateCurrentWaiterIdExists(ctx, tx, *req.CurrentWaiterId); err != nil {
			return err
		}
	}
	
	
	// Validate IsActive
	
	
	
	// Validate AllowOnlineReservation
	
	
	
	// Validate DisplayOrder
	
	
	
	// Validate ColorCode
	
	
	
	// Validate Icon
	
	
	
	// Validate Notes
	
	
	
	// Validate Metadata
	
	
	
	// Validate CreatedBy
	
	
	
	// Validate UpdatedBy
	
	
	

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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *UpdateRestaurantTablesRequest) error {
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
	
	
	// Validate FloorPlanId if provided
	
	
	if req.FloorPlanId != nil {
		if err := v.validateFloorPlanIdExists(ctx, tx, *req.FloorPlanId); err != nil {
			return err
		}
	}
	
	
	// Validate SectionId if provided
	
	
	if req.SectionId != nil {
		if err := v.validateSectionIdExists(ctx, tx, *req.SectionId); err != nil {
			return err
		}
	}
	
	
	// Validate TableNumber if provided
	
	if req.TableNumber != nil {
		if err := v.validateTableNumber(*req.TableNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate TableName if provided
	
	
	
	// Validate MinCapacity if provided
	
	
	
	// Validate MaxCapacity if provided
	
	
	
	// Validate TableShape if provided
	
	
	
	// Validate IsCombinable if provided
	
	
	
	// Validate PositionX if provided
	
	
	
	// Validate PositionY if provided
	
	
	
	// Validate Rotation if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate CurrentCovers if provided
	
	
	
	// Validate SeatedAt if provided
	
	
	
	// Validate CurrentWaiterId if provided
	
	
	if req.CurrentWaiterId != nil {
		if err := v.validateCurrentWaiterIdExists(ctx, tx, *req.CurrentWaiterId); err != nil {
			return err
		}
	}
	
	
	// Validate IsActive if provided
	
	
	
	// Validate AllowOnlineReservation if provided
	
	
	
	// Validate DisplayOrder if provided
	
	
	
	// Validate ColorCode if provided
	
	
	
	// Validate Icon if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Metadata if provided
	
	
	
	// Validate CreatedBy if provided
	
	
	
	// Validate UpdatedBy if provided
	
	
	

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





// validateFloorPlanIdExists validates that floor_plan_id exists
func (v *Validator) validateFloorPlanIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for floor_plan
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM floor_plan WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check floor_plan existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("floor_plan with id %s does not exist", id)
	}
	return nil
}





// validateSectionIdExists validates that section_id exists
func (v *Validator) validateSectionIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for section
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM section WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check section existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("section with id %s does not exist", id)
	}
	return nil
}



// validateTableNumber validates table_number field
func (v *Validator) validateTableNumber(value string) error {
	
	// Add custom validation for table_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("table_number cannot be empty")
	}
	
	return nil
}



















































// validateCurrentWaiterIdExists validates that current_waiter_id exists
func (v *Validator) validateCurrentWaiterIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for current_waiter
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM current_waiter WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check current_waiter existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("current_waiter with id %s does not exist", id)
	}
	return nil
}







































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *CreateRestaurantTablesRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *CreateRestaurantTablesRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *RestaurantTables, req *UpdateRestaurantTablesRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *RestaurantTables) error {
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
