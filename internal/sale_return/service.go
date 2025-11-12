package sale_return

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/sale_return"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/sale_return"
	"go.uber.org/zap"
)

// Service handles business logic for SaleReturns
type Service struct {
	repo   *sale_return.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new SaleReturns service
func NewService(repo *sale_return.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new sale_returns
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateSaleReturnsRequest) (*dto.SaleReturnsResponse, error) {
	s.logger.Info("creating sale_returns",
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
	entity := &sale_return.SaleReturns{
		OrganizationID: orgID,
		
		ReturnNumber: req.ReturnNumber,
		
		OriginalSaleId: req.OriginalSaleId,
		
		CustomerId: req.CustomerId,
		
		LocationId: req.LocationId,
		
		UserId: req.UserId,
		
		ReturnDate: req.ReturnDate,
		
		TotalAmount: req.TotalAmount,
		
		RefundAmount: req.RefundAmount,
		
		RestockingFee: req.RestockingFee,
		
		RefundMethod: req.RefundMethod,
		
		Status: req.Status,
		
		Status: req.Status,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		Notes: req.Notes,
		
		CreatedBy: req.CreatedBy,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create sale_returns: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created sale_returns",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a sale_returns by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.SaleReturnsResponse, error) {
	s.logger.Debug("getting sale_returns",
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
		return nil, fmt.Errorf("failed to get sale_returns: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sale_returns not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of sale_returns records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.SaleReturnsListResponse, error) {
	s.logger.Debug("listing sale_returns",
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
		return nil, fmt.Errorf("failed to list sale_returns: %w", err)
	}

	// Convert to response
	items := make([]*dto.SaleReturnsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.SaleReturnsListResponse{
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

// Update updates an existing sale_returns
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateSaleReturnsRequest) (*dto.SaleReturnsResponse, error) {
	s.logger.Info("updating sale_returns",
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
		return nil, fmt.Errorf("failed to get sale_returns: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sale_returns not found or access denied")
	}
	

	// Update fields
	
	if req.ReturnNumber != nil {
		entity.ReturnNumber = *req.ReturnNumber
	}
	
	if req.OriginalSaleId != nil {
		entity.OriginalSaleId = *req.OriginalSaleId
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.UserId != nil {
		entity.UserId = *req.UserId
	}
	
	if req.ReturnDate != nil {
		entity.ReturnDate = *req.ReturnDate
	}
	
	if req.TotalAmount != nil {
		entity.TotalAmount = *req.TotalAmount
	}
	
	if req.RefundAmount != nil {
		entity.RefundAmount = *req.RefundAmount
	}
	
	if req.RestockingFee != nil {
		entity.RestockingFee = *req.RestockingFee
	}
	
	if req.RefundMethod != nil {
		entity.RefundMethod = *req.RefundMethod
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = *req.ApprovedBy
	}
	
	if req.ApprovedAt != nil {
		entity.ApprovedAt = *req.ApprovedAt
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.CreatedBy != nil {
		entity.CreatedBy = *req.CreatedBy
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update sale_returns: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated sale_returns",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a sale_returns
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting sale_returns",
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
		return fmt.Errorf("failed to get sale_returns: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("sale_returns not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete sale_returns: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted sale_returns",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *sale_return.SaleReturns) *dto.SaleReturnsResponse {
	return &dto.SaleReturnsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		ReturnNumber: entity.ReturnNumber,
		
		OriginalSaleId: entity.OriginalSaleId,
		
		CustomerId: entity.CustomerId,
		
		LocationId: entity.LocationId,
		
		UserId: entity.UserId,
		
		ReturnDate: entity.ReturnDate,
		
		TotalAmount: entity.TotalAmount,
		
		RefundAmount: entity.RefundAmount,
		
		RestockingFee: entity.RestockingFee,
		
		RefundMethod: entity.RefundMethod,
		
		Status: entity.Status,
		
		Status: entity.Status,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		Notes: entity.Notes,
		
		CreatedBy: entity.CreatedBy,
		
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


// validateBusinessRules validates business rules for sale_returns
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *sale_return.SaleReturns) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a sale_returns can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
