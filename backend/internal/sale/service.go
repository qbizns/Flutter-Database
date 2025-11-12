package sale

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

// Service handles business logic for Sales
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Sales service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new sales
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateSalesRequest) (*SalesResponse, error) {
	s.logger.Info("creating sales",
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
	entity := &Sales{
		OrganizationID: orgID,
		
		SaleNumber: req.SaleNumber,
		
		ReferenceNumber: req.ReferenceNumber,
		
		TransactionType: req.TransactionType,
		
		CustomerId: req.CustomerId,
		
		CashierId: req.CashierId,
		
		Subtotal: req.Subtotal,
		
		TaxAmount: req.TaxAmount,
		
		DiscountAmount: req.DiscountAmount,
		
		TotalAmount: req.TotalAmount,
		
		PaidAmount: req.PaidAmount,
		
		ChangeAmount: req.ChangeAmount,
		
		OutstandingAmount: req.OutstandingAmount,
		
		PaymentStatus: req.PaymentStatus,
		
		DiscountType: req.DiscountType,
		
		DiscountValue: req.DiscountValue,
		
		DiscountReason: req.DiscountReason,
		
		TransactionDate: req.TransactionDate,
		
		CompletedAt: req.CompletedAt,
		
		Notes: req.Notes,
		
		InternalNotes: req.InternalNotes,
		
		CustomFields: req.CustomFields,
		
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
		return nil, fmt.Errorf("failed to create sales: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created sales",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a sales by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SalesResponse, error) {
	s.logger.Debug("getting sales",
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
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sales not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of sales records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*SalesListResponse, error) {
	s.logger.Debug("listing sales",
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
		return nil, fmt.Errorf("failed to list sales: %w", err)
	}

	// Convert to response
	items := make([]*SalesResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &SalesListResponse{
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

// Update updates an existing sales
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateSalesRequest) (*SalesResponse, error) {
	s.logger.Info("updating sales",
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
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("sales not found or access denied")
	}
	

	// Update fields
	
	if req.SaleNumber != nil {
		entity.SaleNumber = *req.SaleNumber
	}
	
	if req.ReferenceNumber != nil {
		entity.ReferenceNumber = *req.ReferenceNumber
	}
	
	if req.TransactionType != nil {
		entity.TransactionType = *req.TransactionType
	}
	
	if req.CustomerId != nil {
		entity.CustomerId = *req.CustomerId
	}
	
	if req.CashierId != nil {
		entity.CashierId = *req.CashierId
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
	
	if req.ChangeAmount != nil {
		entity.ChangeAmount = *req.ChangeAmount
	}
	
	if req.OutstandingAmount != nil {
		entity.OutstandingAmount = *req.OutstandingAmount
	}
	
	if req.PaymentStatus != nil {
		entity.PaymentStatus = *req.PaymentStatus
	}
	
	if req.DiscountType != nil {
		entity.DiscountType = *req.DiscountType
	}
	
	if req.DiscountValue != nil {
		entity.DiscountValue = *req.DiscountValue
	}
	
	if req.DiscountReason != nil {
		entity.DiscountReason = *req.DiscountReason
	}
	
	if req.TransactionDate != nil {
		entity.TransactionDate = *req.TransactionDate
	}
	
	if req.CompletedAt != nil {
		entity.CompletedAt = *req.CompletedAt
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.InternalNotes != nil {
		entity.InternalNotes = *req.InternalNotes
	}
	
	if req.CustomFields != nil {
		entity.CustomFields = *req.CustomFields
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
		return nil, fmt.Errorf("failed to update sales: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated sales",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a sales
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting sales",
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
		return fmt.Errorf("failed to get sales: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("sales not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete sales: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted sales",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Sales) *SalesResponse {
	return &SalesResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SaleNumber: entity.SaleNumber,
		
		ReferenceNumber: entity.ReferenceNumber,
		
		TransactionType: entity.TransactionType,
		
		CustomerId: entity.CustomerId,
		
		CashierId: entity.CashierId,
		
		Subtotal: entity.Subtotal,
		
		TaxAmount: entity.TaxAmount,
		
		DiscountAmount: entity.DiscountAmount,
		
		TotalAmount: entity.TotalAmount,
		
		PaidAmount: entity.PaidAmount,
		
		ChangeAmount: entity.ChangeAmount,
		
		OutstandingAmount: entity.OutstandingAmount,
		
		PaymentStatus: entity.PaymentStatus,
		
		DiscountType: entity.DiscountType,
		
		DiscountValue: entity.DiscountValue,
		
		DiscountReason: entity.DiscountReason,
		
		TransactionDate: entity.TransactionDate,
		
		CompletedAt: entity.CompletedAt,
		
		Notes: entity.Notes,
		
		InternalNotes: entity.InternalNotes,
		
		CustomFields: entity.CustomFields,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
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


// validateBusinessRules validates business rules for sales
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Sales) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a sales can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
