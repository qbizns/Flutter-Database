package inventory_transfer

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

// Service handles business logic for InventoryTransfers
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new InventoryTransfers service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new inventory_transfers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateInventoryTransfersRequest) (*InventoryTransfersResponse, error) {
	s.logger.Info("creating inventory_transfers",
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
	entity := &InventoryTransfers{
		OrganizationID: orgID,
		
		TransferNumber: req.TransferNumber,
		
		TransferDate: req.TransferDate,
		
		FromLocationId: req.FromLocationId,
		
		ToLocationId: req.ToLocationId,
		
		Status: req.Status,
		
		RequestedDate: req.RequestedDate,
		
		ApprovedDate: req.ApprovedDate,
		
		ShippedDate: req.ShippedDate,
		
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		
		ReceivedDate: req.ReceivedDate,
		
		Carrier: req.Carrier,
		
		TrackingNumber: req.TrackingNumber,
		
		ShippingCost: req.ShippingCost,
		
		Reason: req.Reason,
		
		Notes: req.Notes,
		
		RejectionReason: req.RejectionReason,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		RequestedBy: req.RequestedBy,
		
		ApprovedBy: req.ApprovedBy,
		
		ShippedBy: req.ShippedBy,
		
		ReceivedBy: req.ReceivedBy,
		
		(approvedDate: req.(approvedDate,
		
		(shippedDate: req.(shippedDate,
		
		(receivedDate: req.(receivedDate,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create inventory_transfers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created inventory_transfers",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a inventory_transfers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransfersResponse, error) {
	s.logger.Debug("getting inventory_transfers",
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
		return nil, fmt.Errorf("failed to get inventory_transfers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_transfers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of inventory_transfers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*InventoryTransfersListResponse, error) {
	s.logger.Debug("listing inventory_transfers",
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
		return nil, fmt.Errorf("failed to list inventory_transfers: %w", err)
	}

	// Convert to response
	items := make([]*InventoryTransfersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &InventoryTransfersListResponse{
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

// Update updates an existing inventory_transfers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateInventoryTransfersRequest) (*InventoryTransfersResponse, error) {
	s.logger.Info("updating inventory_transfers",
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
		return nil, fmt.Errorf("failed to get inventory_transfers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("inventory_transfers not found or access denied")
	}
	

	// Update fields
	
	if req.TransferNumber != nil {
		entity.TransferNumber = *req.TransferNumber
	}
	
	if req.TransferDate != nil {
		entity.TransferDate = *req.TransferDate
	}
	
	if req.FromLocationId != nil {
		entity.FromLocationId = *req.FromLocationId
	}
	
	if req.ToLocationId != nil {
		entity.ToLocationId = *req.ToLocationId
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.RequestedDate != nil {
		entity.RequestedDate = *req.RequestedDate
	}
	
	if req.ApprovedDate != nil {
		entity.ApprovedDate = *req.ApprovedDate
	}
	
	if req.ShippedDate != nil {
		entity.ShippedDate = *req.ShippedDate
	}
	
	if req.ExpectedDeliveryDate != nil {
		entity.ExpectedDeliveryDate = *req.ExpectedDeliveryDate
	}
	
	if req.ReceivedDate != nil {
		entity.ReceivedDate = *req.ReceivedDate
	}
	
	if req.Carrier != nil {
		entity.Carrier = *req.Carrier
	}
	
	if req.TrackingNumber != nil {
		entity.TrackingNumber = *req.TrackingNumber
	}
	
	if req.ShippingCost != nil {
		entity.ShippingCost = *req.ShippingCost
	}
	
	if req.Reason != nil {
		entity.Reason = *req.Reason
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.RejectionReason != nil {
		entity.RejectionReason = *req.RejectionReason
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
	
	if req.RequestedBy != nil {
		entity.RequestedBy = *req.RequestedBy
	}
	
	if req.ApprovedBy != nil {
		entity.ApprovedBy = *req.ApprovedBy
	}
	
	if req.ShippedBy != nil {
		entity.ShippedBy = *req.ShippedBy
	}
	
	if req.ReceivedBy != nil {
		entity.ReceivedBy = *req.ReceivedBy
	}
	
	if req.(approvedDate != nil {
		entity.(approvedDate = *req.(approvedDate
	}
	
	if req.(shippedDate != nil {
		entity.(shippedDate = *req.(shippedDate
	}
	
	if req.(receivedDate != nil {
		entity.(receivedDate = *req.(receivedDate
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update inventory_transfers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated inventory_transfers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a inventory_transfers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting inventory_transfers",
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
		return fmt.Errorf("failed to get inventory_transfers: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("inventory_transfers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete inventory_transfers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted inventory_transfers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *InventoryTransfers) *InventoryTransfersResponse {
	return &InventoryTransfersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		TransferNumber: entity.TransferNumber,
		
		TransferDate: entity.TransferDate,
		
		FromLocationId: entity.FromLocationId,
		
		ToLocationId: entity.ToLocationId,
		
		Status: entity.Status,
		
		RequestedDate: entity.RequestedDate,
		
		ApprovedDate: entity.ApprovedDate,
		
		ShippedDate: entity.ShippedDate,
		
		ExpectedDeliveryDate: entity.ExpectedDeliveryDate,
		
		ReceivedDate: entity.ReceivedDate,
		
		Carrier: entity.Carrier,
		
		TrackingNumber: entity.TrackingNumber,
		
		ShippingCost: entity.ShippingCost,
		
		Reason: entity.Reason,
		
		Notes: entity.Notes,
		
		RejectionReason: entity.RejectionReason,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		RequestedBy: entity.RequestedBy,
		
		ApprovedBy: entity.ApprovedBy,
		
		ShippedBy: entity.ShippedBy,
		
		ReceivedBy: entity.ReceivedBy,
		
		(approvedDate: entity.(approvedDate,
		
		(shippedDate: entity.(shippedDate,
		
		(receivedDate: entity.(receivedDate,
		
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


// validateBusinessRules validates business rules for inventory_transfers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *InventoryTransfers) error {
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

// canDelete checks if a inventory_transfers can be deleted
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
