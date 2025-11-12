package vendor_bill

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

// Validator handles VendorBills validation logic
type Validator struct {
	repo *Repository
}

// NewValidator creates a new VendorBills validator
func NewValidator(repo *Repository) *Validator {
	return &Validator{
		repo: repo,
	}
}

// ValidateCreate validates a create request
func (v *Validator) ValidateCreate(ctx context.Context, tx pgx.Tx, req *dto.CreateVendorBillsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	
	// Validate BillNumber
	
	if err := v.validateBillNumber(req.BillNumber); err != nil {
		return err
	}
	
	
	
	// Validate VendorBillNumber
	
	
	
	// Validate SupplierId
	
	
	if err := v.validateSupplierIdExists(ctx, tx, req.SupplierId); err != nil {
		return err
	}
	
	
	// Validate BillDate
	
	
	
	// Validate DueDate
	
	
	
	// Validate PaymentTerms
	
	
	
	// Validate AccountingPeriodId
	
	
	if err := v.validateAccountingPeriodIdExists(ctx, tx, req.AccountingPeriodId); err != nil {
		return err
	}
	
	
	// Validate Subtotal
	
	
	
	// Validate TaxAmount
	
	
	
	// Validate TotalAmount
	
	
	
	// Validate PaidAmount
	
	
	
	// Validate BalanceDue
	
	
	
	// Validate Status
	
	
	
	// Validate JournalEntryId
	
	
	if err := v.validateJournalEntryIdExists(ctx, tx, req.JournalEntryId); err != nil {
		return err
	}
	
	
	// Validate IsPosted
	
	
	
	// Validate PurchaseOrderId
	
	
	if err := v.validatePurchaseOrderIdExists(ctx, tx, req.PurchaseOrderId); err != nil {
		return err
	}
	
	
	// Validate Description
	
	
	
	// Validate Notes
	
	
	
	// Validate Memo
	
	
	
	// Validate Attachments
	
	
	
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
func (v *Validator) ValidateUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID, req *dto.UpdateVendorBillsRequest) error {
	// Basic validation
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if entity exists
	existing, err := v.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("entity not found: %w", err)
	}

	
	// Validate BillNumber if provided
	
	if req.BillNumber != nil {
		if err := v.validateBillNumber(*req.BillNumber); err != nil {
			return err
		}
	}
	
	
	
	// Validate VendorBillNumber if provided
	
	
	
	// Validate SupplierId if provided
	
	
	if req.SupplierId != nil {
		if err := v.validateSupplierIdExists(ctx, tx, *req.SupplierId); err != nil {
			return err
		}
	}
	
	
	// Validate BillDate if provided
	
	
	
	// Validate DueDate if provided
	
	
	
	// Validate PaymentTerms if provided
	
	
	
	// Validate AccountingPeriodId if provided
	
	
	if req.AccountingPeriodId != nil {
		if err := v.validateAccountingPeriodIdExists(ctx, tx, *req.AccountingPeriodId); err != nil {
			return err
		}
	}
	
	
	// Validate Subtotal if provided
	
	
	
	// Validate TaxAmount if provided
	
	
	
	// Validate TotalAmount if provided
	
	
	
	// Validate PaidAmount if provided
	
	
	
	// Validate BalanceDue if provided
	
	
	
	// Validate Status if provided
	
	
	
	// Validate JournalEntryId if provided
	
	
	if req.JournalEntryId != nil {
		if err := v.validateJournalEntryIdExists(ctx, tx, *req.JournalEntryId); err != nil {
			return err
		}
	}
	
	
	// Validate IsPosted if provided
	
	
	
	// Validate PurchaseOrderId if provided
	
	
	if req.PurchaseOrderId != nil {
		if err := v.validatePurchaseOrderIdExists(ctx, tx, *req.PurchaseOrderId); err != nil {
			return err
		}
	}
	
	
	// Validate Description if provided
	
	
	
	// Validate Notes if provided
	
	
	
	// Validate Memo if provided
	
	
	
	// Validate Attachments if provided
	
	
	
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



// validateBillNumber validates bill_number field
func (v *Validator) validateBillNumber(value string) error {
	
	// Add custom validation for bill_number
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("bill_number cannot be empty")
	}
	
	return nil
}











// validateSupplierIdExists validates that supplier_id exists
func (v *Validator) validateSupplierIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for supplier
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM supplier WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check supplier existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("supplier with id %s does not exist", id)
	}
	return nil
}

















// validateAccountingPeriodIdExists validates that accounting_period_id exists
func (v *Validator) validateAccountingPeriodIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for accounting_period
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM accounting_period WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check accounting_period existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("accounting_period with id %s does not exist", id)
	}
	return nil
}





























// validateJournalEntryIdExists validates that journal_entry_id exists
func (v *Validator) validateJournalEntryIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for journal_entry
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM journal_entry WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check journal_entry existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("journal_entry with id %s does not exist", id)
	}
	return nil
}









// validatePurchaseOrderIdExists validates that purchase_order_id exists
func (v *Validator) validatePurchaseOrderIdExists(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Implement existence check for purchase_order
	// This should query the referenced table to ensure the ID exists
	var exists bool
	err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM purchase_order WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check purchase_order existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("purchase_order with id %s does not exist", id)
	}
	return nil
}































// validateCrossFields validates relationships between fields
func (v *Validator) validateCrossFields(ctx context.Context, tx pgx.Tx, req *dto.CreateVendorBillsRequest) error {
	// Add cross-field validation logic here
	// Example: start_date must be before end_date
	// Example: price must be less than max_price
	return nil
}

// validateBusinessRules validates business-specific rules
func (v *Validator) validateBusinessRules(ctx context.Context, tx pgx.Tx, req *dto.CreateVendorBillsRequest) error {
	// Add business rule validation here
	// Example: check inventory levels
	// Example: validate credit limits
	// Example: check permission constraints
	return nil
}

// validateUpdateBusinessRules validates business rules for updates
func (v *Validator) validateUpdateBusinessRules(ctx context.Context, tx pgx.Tx, existing *VendorBills, req *dto.UpdateVendorBillsRequest) error {
	// Add update-specific business rule validation here
	// Example: can't change status from 'completed' to 'pending'
	// Example: can't reduce quantity below reserved amount
	return nil
}

// validateCanDelete checks if entity can be safely deleted
func (v *Validator) validateCanDelete(ctx context.Context, tx pgx.Tx, entity *VendorBills) error {
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
