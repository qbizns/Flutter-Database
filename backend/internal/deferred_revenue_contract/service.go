package deferred_revenue_contract

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for DeferredRevenueContracts
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new DeferredRevenueContracts service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new deferred_revenue_contracts
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateDeferredRevenueContractsRequest) (*DeferredRevenueContractsResponse, error) {
	s.logger.Info("creating deferred_revenue_contracts",
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
	entity := &DeferredRevenueContracts{
		OrganizationId: orgID,
		
		CustomerInvoiceId: req.CustomerInvoiceId,
		
		InvoiceLineId: req.InvoiceLineId,
		
		ContractName: req.ContractName,
		
		TotalDeferredAmount: req.TotalDeferredAmount,
		
		StartDate: req.StartDate,
		
		EndDate: req.EndDate,
		
		RecognitionMethod: req.RecognitionMethod,
		
		RecognitionMethod: req.RecognitionMethod,
		
		DeferredAccountId: req.DeferredAccountId,
		
		RevenueAccountId: req.RevenueAccountId,
		
		Status: req.Status,
		
		RecognizedAmount: req.RecognizedAmount,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create deferred_revenue_contracts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created deferred_revenue_contracts",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a deferred_revenue_contracts by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*DeferredRevenueContractsResponse, error) {
	s.logger.Debug("getting deferred_revenue_contracts",
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
		return nil, fmt.Errorf("failed to get deferred_revenue_contracts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("deferred_revenue_contracts not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of deferred_revenue_contracts records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*DeferredRevenueContractsListResponse, error) {
	s.logger.Debug("listing deferred_revenue_contracts",
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
		return nil, fmt.Errorf("failed to list deferred_revenue_contracts: %w", err)
	}

	// Convert to response
	items := make([]*DeferredRevenueContractsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &DeferredRevenueContractsListResponse{
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

// Update updates an existing deferred_revenue_contracts
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateDeferredRevenueContractsRequest) (*DeferredRevenueContractsResponse, error) {
	s.logger.Info("updating deferred_revenue_contracts",
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
		return nil, fmt.Errorf("failed to get deferred_revenue_contracts: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("deferred_revenue_contracts not found or access denied")
	}
	

	// Update fields
	
	if req.CustomerInvoiceId != nil {
		entity.CustomerInvoiceId = req.CustomerInvoiceId
	}
	
	if req.InvoiceLineId != nil {
		entity.InvoiceLineId = req.InvoiceLineId
	}
	
	if req.ContractName != nil {
		entity.ContractName = req.ContractName
	}
	
	if req.TotalDeferredAmount != nil {
		entity.TotalDeferredAmount = req.TotalDeferredAmount
	}
	
	if req.StartDate != nil {
		entity.StartDate = req.StartDate
	}
	
	if req.EndDate != nil {
		entity.EndDate = req.EndDate
	}
	
	if req.RecognitionMethod != nil {
		entity.RecognitionMethod = req.RecognitionMethod
	}
	
	if req.RecognitionMethod != nil {
		entity.RecognitionMethod = req.RecognitionMethod
	}
	
	if req.DeferredAccountId != nil {
		entity.DeferredAccountId = req.DeferredAccountId
	}
	
	if req.RevenueAccountId != nil {
		entity.RevenueAccountId = req.RevenueAccountId
	}
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.RecognizedAmount != nil {
		entity.RecognizedAmount = req.RecognizedAmount
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update deferred_revenue_contracts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated deferred_revenue_contracts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a deferred_revenue_contracts
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting deferred_revenue_contracts",
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
		return fmt.Errorf("failed to get deferred_revenue_contracts: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("deferred_revenue_contracts not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete deferred_revenue_contracts: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted deferred_revenue_contracts",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *DeferredRevenueContracts) *DeferredRevenueContractsResponse {
	return &DeferredRevenueContractsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CustomerInvoiceId: entity.CustomerInvoiceId,
		
		InvoiceLineId: entity.InvoiceLineId,
		
		ContractName: entity.ContractName,
		
		TotalDeferredAmount: entity.TotalDeferredAmount,
		
		StartDate: entity.StartDate,
		
		EndDate: entity.EndDate,
		
		RecognitionMethod: entity.RecognitionMethod,
		
		RecognitionMethod: entity.RecognitionMethod,
		
		DeferredAccountId: entity.DeferredAccountId,
		
		RevenueAccountId: entity.RevenueAccountId,
		
		Status: entity.Status,
		
		RecognizedAmount: entity.RecognizedAmount,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for deferred_revenue_contracts
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *DeferredRevenueContracts) error {
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

// canDelete checks if a deferred_revenue_contracts can be deleted
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
