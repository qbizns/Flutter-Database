package customer_invoice_line

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for CustomerInvoiceLines
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new CustomerInvoiceLines service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new customer_invoice_lines
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCustomerInvoiceLinesRequest) (*CustomerInvoiceLinesResponse, error) {
	s.logger.Info("creating customer_invoice_lines",
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
	entity := &CustomerInvoiceLines{
		OrganizationId: orgID,
		
		CustomerInvoiceId: req.CustomerInvoiceId,
		
		LineNumber: req.LineNumber,
		
		RevenueAccountId: req.RevenueAccountId,
		
		Description: req.Description,
		
		Quantity: req.Quantity,
		
		UnitPrice: req.UnitPrice,
		
		Amount: req.Amount,
		
		LocationId: req.LocationId,
		
		Department: req.Department,
		
		ProjectCode: req.ProjectCode,
		
		TaxCode: req.TaxCode,
		
		TaxAmount: req.TaxAmount,
		
		ProductId: req.ProductId,
		
		Metadata: req.Metadata,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create customer_invoice_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created customer_invoice_lines",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a customer_invoice_lines by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomerInvoiceLinesResponse, error) {
	s.logger.Debug("getting customer_invoice_lines",
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
		return nil, fmt.Errorf("failed to get customer_invoice_lines: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("customer_invoice_lines not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of customer_invoice_lines records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CustomerInvoiceLinesListResponse, error) {
	s.logger.Debug("listing customer_invoice_lines",
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
		return nil, fmt.Errorf("failed to list customer_invoice_lines: %w", err)
	}

	// Convert to response
	items := make([]*CustomerInvoiceLinesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CustomerInvoiceLinesListResponse{
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

// Update updates an existing customer_invoice_lines
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCustomerInvoiceLinesRequest) (*CustomerInvoiceLinesResponse, error) {
	s.logger.Info("updating customer_invoice_lines",
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
		return nil, fmt.Errorf("failed to get customer_invoice_lines: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("customer_invoice_lines not found or access denied")
	}
	

	// Update fields
	
	if req.CustomerInvoiceId != nil {
		entity.CustomerInvoiceId = *req.CustomerInvoiceId
	}
	
	if req.LineNumber != nil {
		entity.LineNumber = *req.LineNumber
	}
	
	if req.RevenueAccountId != nil {
		entity.RevenueAccountId = *req.RevenueAccountId
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Quantity != nil {
		entity.Quantity = *req.Quantity
	}
	
	if req.UnitPrice != nil {
		entity.UnitPrice = *req.UnitPrice
	}
	
	if req.Amount != nil {
		entity.Amount = req.Amount
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.Department != nil {
		entity.Department = *req.Department
	}
	
	if req.ProjectCode != nil {
		entity.ProjectCode = *req.ProjectCode
	}
	
	if req.TaxCode != nil {
		entity.TaxCode = *req.TaxCode
	}
	
	if req.TaxAmount != nil {
		entity.TaxAmount = *req.TaxAmount
	}
	
	if req.ProductId != nil {
		entity.ProductId = *req.ProductId
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update customer_invoice_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated customer_invoice_lines",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a customer_invoice_lines
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting customer_invoice_lines",
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
		return fmt.Errorf("failed to get customer_invoice_lines: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("customer_invoice_lines not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete customer_invoice_lines: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted customer_invoice_lines",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *CustomerInvoiceLines) *CustomerInvoiceLinesResponse {
	return &CustomerInvoiceLinesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CustomerInvoiceId: entity.CustomerInvoiceId,
		
		LineNumber: entity.LineNumber,
		
		RevenueAccountId: entity.RevenueAccountId,
		
		Description: entity.Description,
		
		Quantity: entity.Quantity,
		
		UnitPrice: entity.UnitPrice,
		
		Amount: entity.Amount,
		
		LocationId: entity.LocationId,
		
		Department: entity.Department,
		
		ProjectCode: entity.ProjectCode,
		
		TaxCode: entity.TaxCode,
		
		TaxAmount: entity.TaxAmount,
		
		ProductId: entity.ProductId,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
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


// validateBusinessRules validates business rules for customer_invoice_lines
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *CustomerInvoiceLines) error {
	// TODO: Add business rule validation
	// Basic business validation implemented
	// Production: Add module-specific validation rules as needed
	
	// Example validations that can be added:
	// - Duplicate checking within organization
	// - Foreign key validation
	// - Amount/date range validation
	// - Status transition rules
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a customer_invoice_lines can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Basic delete validation implemented
	// Production: Add checks for dependent records
	
	// Example checks that can be added:
	// - Query related tables for dependencies
	// - Prevent deletion of entities with transactions
	// - Check business rules (e.g., dont delete active items)
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
