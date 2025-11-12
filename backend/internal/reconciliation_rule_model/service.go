package reconciliation_rule_model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/reconciliation_rule_model"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/reconciliation_rule_model"
	"go.uber.org/zap"
)

// Service handles business logic for ReconciliationRuleModels
type Service struct {
	repo   *reconciliation_rule_model.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new ReconciliationRuleModels service
func NewService(repo *reconciliation_rule_model.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new reconciliation_rule_models
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateReconciliationRuleModelsRequest) (*dto.ReconciliationRuleModelsResponse, error) {
	s.logger.Info("creating reconciliation_rule_models",
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
	entity := &reconciliation_rule_model.ReconciliationRuleModels{
		OrganizationID: orgID,
		
		RuleName: req.RuleName,
		
		RuleCode: req.RuleCode,
		
		Sequence: req.Sequence,
		
		AmountMin: req.AmountMin,
		
		AmountMax: req.AmountMax,
		
		DescriptionPattern: req.DescriptionPattern,
		
		CounterpartyPattern: req.CounterpartyPattern,
		
		ReferencePattern: req.ReferencePattern,
		
		JournalId: req.JournalId,
		
		AccountId: req.AccountId,
		
		AnalyticAccountId: req.AnalyticAccountId,
		
		TaxId: req.TaxId,
		
		IsActive: req.IsActive,
		
		AutoApply: req.AutoApply,
		
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
		return nil, fmt.Errorf("failed to create reconciliation_rule_models: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created reconciliation_rule_models",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a reconciliation_rule_models by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.ReconciliationRuleModelsResponse, error) {
	s.logger.Debug("getting reconciliation_rule_models",
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
		return nil, fmt.Errorf("failed to get reconciliation_rule_models: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("reconciliation_rule_models not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of reconciliation_rule_models records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.ReconciliationRuleModelsListResponse, error) {
	s.logger.Debug("listing reconciliation_rule_models",
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
		return nil, fmt.Errorf("failed to list reconciliation_rule_models: %w", err)
	}

	// Convert to response
	items := make([]*dto.ReconciliationRuleModelsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.ReconciliationRuleModelsListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

// Update updates an existing reconciliation_rule_models
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateReconciliationRuleModelsRequest) (*dto.ReconciliationRuleModelsResponse, error) {
	s.logger.Info("updating reconciliation_rule_models",
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
		return nil, fmt.Errorf("failed to get reconciliation_rule_models: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("reconciliation_rule_models not found or access denied")
	}
	

	// Update fields
	
	if req.RuleName != nil {
		entity.RuleName = *req.RuleName
	}
	
	if req.RuleCode != nil {
		entity.RuleCode = *req.RuleCode
	}
	
	if req.Sequence != nil {
		entity.Sequence = *req.Sequence
	}
	
	if req.AmountMin != nil {
		entity.AmountMin = *req.AmountMin
	}
	
	if req.AmountMax != nil {
		entity.AmountMax = *req.AmountMax
	}
	
	if req.DescriptionPattern != nil {
		entity.DescriptionPattern = *req.DescriptionPattern
	}
	
	if req.CounterpartyPattern != nil {
		entity.CounterpartyPattern = *req.CounterpartyPattern
	}
	
	if req.ReferencePattern != nil {
		entity.ReferencePattern = *req.ReferencePattern
	}
	
	if req.JournalId != nil {
		entity.JournalId = *req.JournalId
	}
	
	if req.AccountId != nil {
		entity.AccountId = *req.AccountId
	}
	
	if req.AnalyticAccountId != nil {
		entity.AnalyticAccountId = *req.AnalyticAccountId
	}
	
	if req.TaxId != nil {
		entity.TaxId = *req.TaxId
	}
	
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	
	if req.AutoApply != nil {
		entity.AutoApply = *req.AutoApply
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update reconciliation_rule_models: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated reconciliation_rule_models",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a reconciliation_rule_models
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting reconciliation_rule_models",
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
		return fmt.Errorf("failed to get reconciliation_rule_models: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("reconciliation_rule_models not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete reconciliation_rule_models: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted reconciliation_rule_models",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *reconciliation_rule_model.ReconciliationRuleModels) *dto.ReconciliationRuleModelsResponse {
	return &dto.ReconciliationRuleModelsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		RuleName: entity.RuleName,
		
		RuleCode: entity.RuleCode,
		
		Sequence: entity.Sequence,
		
		AmountMin: entity.AmountMin,
		
		AmountMax: entity.AmountMax,
		
		DescriptionPattern: entity.DescriptionPattern,
		
		CounterpartyPattern: entity.CounterpartyPattern,
		
		ReferencePattern: entity.ReferencePattern,
		
		JournalId: entity.JournalId,
		
		AccountId: entity.AccountId,
		
		AnalyticAccountId: entity.AnalyticAccountId,
		
		TaxId: entity.TaxId,
		
		IsActive: entity.IsActive,
		
		AutoApply: entity.AutoApply,
		
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


// validateBusinessRules validates business rules for reconciliation_rule_models
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *reconciliation_rule_model.ReconciliationRuleModels) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a reconciliation_rule_models can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
