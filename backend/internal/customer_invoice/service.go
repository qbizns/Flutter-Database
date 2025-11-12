package customer_invoice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for CustomerInvoices
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new CustomerInvoices service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new customer_invoices
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCustomerInvoicesRequest) (*CustomerInvoicesResponse, error) {
	s.logger.Info("creating customer_invoices",
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Convert DTO to entity
	entity := &CustomerInvoices{
		OrganizationID: orgID,
		
		InvoiceNumber: req.InvoiceNumber,
		
		CustomerId: req.CustomerId,
		
		InvoiceDate: req.InvoiceDate,
		
		DueDate: req.DueDate,
		
		PaymentTerms: req.PaymentTerms,
		
		AccountingPeriodId: req.AccountingPeriodId,
		
		Subtotal: req.Subtotal,
		
		TaxAmount: req.TaxAmount,
		
		DiscountAmount: req.DiscountAmount,
		
		TotalAmount: req.TotalAmount,
		
		PaidAmount: req.PaidAmount,
		
		BalanceDue: req.BalanceDue,
		
		Status: req.Status,
		
		JournalEntryId: req.JournalEntryId,
		
		IsPosted: req.IsPosted,
		
		SaleId: req.SaleId,
		
		Description: req.Description,
		
		Notes: req.Notes,
		
		Memo: req.Memo,
		
		Attachments: req.Attachments,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create customer_invoices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created customer_invoices",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a customer_invoices by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomerInvoicesResponse, error) {
	s.logger.Debug("getting customer_invoices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer_invoices: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("customer_invoices not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of customer_invoices records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CustomerInvoicesListResponse, error) {
	s.logger.Debug("listing customer_invoices",
		zap.String("organization_id", orgID.String()),
		zap.Int("page", page),
		zap.Int("limit", limit),
	)

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}

	// Get from database (organization-scoped)
	entities, total, err := s.repo.ListByOrganization(ctx, tx, orgID, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list customer_invoices: %w", err)
	}

	// Convert to response
	items := make([]*CustomerInvoicesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CustomerInvoicesListResponse{
		Items: items,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing customer_invoices
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCustomerInvoicesRequest) (*CustomerInvoicesResponse, error) {
	s.logger.Info("updating customer_invoices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return nil, err
	}
	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer_invoices: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("customer_invoices not found or access denied")
	}
	

	// Update fields
	
	if req.InvoiceNumber != nil {
		entity.InvoiceNumber = *req.InvoiceNumber
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.InvoiceDate != nil {
		entity.InvoiceDate = *req.InvoiceDate
	}
	
	if req.DueDate != nil {
		entity.DueDate = *req.DueDate
	}
	
	if req.PaymentTerms != nil {
		entity.PaymentTerms = *req.PaymentTerms
	}
	
	if req.AccountingPeriodId != nil {
		entity.AccountingPeriodId = *req.AccountingPeriodId
	}
	
	if req.Subtotal != nil {
		entity.Subtotal = *req.Subtotal
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = *req.TaxAmount
	}
	
	if req.DiscountAmount != nil {
		entity.DiscountAmount = *req.DiscountAmount
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = *req.TotalAmount
	}
	
	if req.PaidAmount != nil {
		entity.PaidAmount = *req.PaidAmount
	}
	
	if req.BalanceDue != nil {
		entity.BalanceDue = *req.BalanceDue
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.JournalEntryId != nil {
		entity.JournalEntryId = *req.JournalEntryId
	}
	
	if req.IsPosted != nil {
		entity.IsPosted = *req.IsPosted
	}
	
	if req.SaleId != nil {
		entity.SaleId = *req.SaleId
	}
	
	if req.Description != nil {
		entity.Description = *req.Description
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.Memo != nil {
		entity.Memo = *req.Memo
	}
	
	if req.Attachments != nil {
		entity.Attachments = *req.Attachments
	}
	
	if req.Metadata != nil {
		entity.Metadata = *req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = *req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update customer_invoices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated customer_invoices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a customer_invoices
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting customer_invoices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	
	// Set organization context for RLS
	if err := s.setOrganizationContext(ctx, tx, orgID); err != nil {
		return err
	}

	// Verify ownership
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer_invoices: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("customer_invoices not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete customer_invoices: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted customer_invoices",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *CustomerInvoices) *CustomerInvoicesResponse {
	return &CustomerInvoicesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		InvoiceNumber: entity.InvoiceNumber,
		
		CustomerId: entity.CustomerId,
		
		InvoiceDate: entity.InvoiceDate,
		
		DueDate: entity.DueDate,
		
		PaymentTerms: entity.PaymentTerms,
		
		AccountingPeriodId: entity.AccountingPeriodId,
		
		Subtotal: entity.Subtotal,
		
		TaxAmount: entity.TaxAmount,
		
		DiscountAmount: entity.DiscountAmount,
		
		TotalAmount: entity.TotalAmount,
		
		PaidAmount: entity.PaidAmount,
		
		BalanceDue: entity.BalanceDue,
		
		Status: entity.Status,
		
		JournalEntryId: entity.JournalEntryId,
		
		IsPosted: entity.IsPosted,
		
		SaleId: entity.SaleId,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		Memo: entity.Memo,
		
		Attachments: entity.Attachments,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
	}
}


// setOrganizationContext sets the organization context for RLS
func (s *Service) setOrganizationContext(ctx context.Context, tx pgx.Tx, orgID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_organization_id = $1", orgID)
	if err != nil {
		return fmt.Errorf("failed to set organization context: %w", err)
	}
	return nil
}


// validateBusinessRules validates business rules for customer_invoices
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *CustomerInvoices) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a customer_invoices can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
