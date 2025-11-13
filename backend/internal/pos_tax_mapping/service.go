package pos_tax_mapping

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for PosTaxMappings
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new PosTaxMappings service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new pos_tax_mappings
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreatePosTaxMappingsRequest) (*PosTaxMappingsResponse, error) {
	s.logger.Info("creating pos_tax_mappings",
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
	entity := &PosTaxMappings{
		OrganizationId: orgID,
		
		PosTaxCode: req.PosTaxCode,
		
		TaxCategoryCode: req.TaxCategoryCode,
		
		PosTaxRate: req.PosTaxRate,
		
		AccountingTaxId: req.AccountingTaxId,
		
		DefaultTaxAccountId: req.DefaultTaxAccountId,
		
		DefaultTaxExpenseAccountId: req.DefaultTaxExpenseAccountId,
		
		IsDefault: req.IsDefault,
		
		IsActive: req.IsActive,
		
		Priority: req.Priority,
		
		IsInclusive: req.IsInclusive,
		
		AppliesToSales: req.AppliesToSales,
		
		AppliesToPurchases: req.AppliesToPurchases,
		
		EffectiveFrom: req.EffectiveFrom,
		
		EffectiveTo: req.EffectiveTo,
		
		Description: req.Description,
		
		Notes: req.Notes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		EffectiveFrom: req.EffectiveFrom,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create pos_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created pos_tax_mappings",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a pos_tax_mappings by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*PosTaxMappingsResponse, error) {
	s.logger.Debug("getting pos_tax_mappings",
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
		return nil, fmt.Errorf("failed to get pos_tax_mappings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("pos_tax_mappings not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of pos_tax_mappings records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*PosTaxMappingsListResponse, error) {
	s.logger.Debug("listing pos_tax_mappings",
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
		return nil, fmt.Errorf("failed to list pos_tax_mappings: %w", err)
	}

	// Convert to response
	items := make([]*PosTaxMappingsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &PosTaxMappingsListResponse{
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

// Update updates an existing pos_tax_mappings
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdatePosTaxMappingsRequest) (*PosTaxMappingsResponse, error) {
	s.logger.Info("updating pos_tax_mappings",
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
		return nil, fmt.Errorf("failed to get pos_tax_mappings: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("pos_tax_mappings not found or access denied")
	}
	

	// Update fields
	
	if req.PosTaxCode != nil {
		entity.PosTaxCode = *req.PosTaxCode
	}
	
	if req.TaxCategoryCode != nil {
		entity.TaxCategoryCode = *req.TaxCategoryCode
	}
	
	if req.PosTaxRate != nil {
		entity.PosTaxRate = *req.PosTaxRate
	}
	
	if req.AccountingTaxId != nil {
		entity.AccountingTaxId = *req.AccountingTaxId
	}
	
	if req.DefaultTaxAccountId != nil {
		entity.DefaultTaxAccountId = *req.DefaultTaxAccountId
	}
	
	if req.DefaultTaxExpenseAccountId != nil {
		entity.DefaultTaxExpenseAccountId = *req.DefaultTaxExpenseAccountId
	}
	
	if req.IsDefault != nil {
		entity.IsDefault = req.IsDefault
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.Priority != nil {
		entity.Priority = req.Priority
	}
	
	if req.IsInclusive != nil {
		entity.IsInclusive = req.IsInclusive
	}
	
	if req.AppliesToSales != nil {
		entity.AppliesToSales = *req.AppliesToSales
	}
	
	if req.AppliesToPurchases != nil {
		entity.AppliesToPurchases = *req.AppliesToPurchases
	}
	
	if req.EffectiveFrom != nil {
		entity.EffectiveFrom = *req.EffectiveFrom
	}
	
	if req.EffectiveTo != nil {
		entity.EffectiveTo = *req.EffectiveTo
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.Metadata != nil {
		entity.Metadata = req.Metadata
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = req.CreatedBy
	}
	
	if req.UpdatedBy != nil {
		entity.UpdatedBy = req.UpdatedBy
	}
	
	if req.EffectiveFrom != nil {
		entity.EffectiveFrom = *req.EffectiveFrom
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update pos_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated pos_tax_mappings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a pos_tax_mappings
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting pos_tax_mappings",
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
		return fmt.Errorf("failed to get pos_tax_mappings: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("pos_tax_mappings not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete pos_tax_mappings: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted pos_tax_mappings",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *PosTaxMappings) *PosTaxMappingsResponse {
	return &PosTaxMappingsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		PosTaxCode: entity.PosTaxCode,
		
		TaxCategoryCode: entity.TaxCategoryCode,
		
		PosTaxRate: entity.PosTaxRate,
		
		AccountingTaxId: entity.AccountingTaxId,
		
		DefaultTaxAccountId: entity.DefaultTaxAccountId,
		
		DefaultTaxExpenseAccountId: entity.DefaultTaxExpenseAccountId,
		
		IsDefault: entity.IsDefault,
		
		IsActive: entity.IsActive,
		
		Priority: entity.Priority,
		
		IsInclusive: entity.IsInclusive,
		
		AppliesToSales: entity.AppliesToSales,
		
		AppliesToPurchases: entity.AppliesToPurchases,
		
		EffectiveFrom: entity.EffectiveFrom,
		
		EffectiveTo: entity.EffectiveTo,
		
		Description: entity.Description,
		
		Notes: entity.Notes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		EffectiveFrom: entity.EffectiveFrom,
		
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


// validateBusinessRules validates business rules for pos_tax_mappings
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *PosTaxMappings) error {
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

// canDelete checks if a pos_tax_mappings can be deleted
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
