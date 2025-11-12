package staff_commission

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

// Service handles business logic for StaffCommissions
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new StaffCommissions service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new staff_commissions
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateStaffCommissionsRequest) (*StaffCommissionsResponse, error) {
	s.logger.Info("creating staff_commissions",
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
	entity := &StaffCommissions{
		OrganizationID: orgID,
		
		LocationId: req.LocationId,
		
		EmployeeId: req.EmployeeId,
		
		CommissionDate: req.CommissionDate,
		
		PeriodStart: req.PeriodStart,
		
		PeriodEnd: req.PeriodEnd,
		
		SourceType: req.SourceType,
		
		SourceSaleId: req.SourceSaleId,
		
		SourceOrderId: req.SourceOrderId,
		
		CommissionType: req.CommissionType,
		
		CommissionRate: req.CommissionRate,
		
		SalesAmount: req.SalesAmount,
		
		CommissionAmount: req.CommissionAmount,
		
		Status: req.Status,
		
		ApprovedBy: req.ApprovedBy,
		
		ApprovedAt: req.ApprovedAt,
		
		PaymentDate: req.PaymentDate,
		
		PaymentMethod: req.PaymentMethod,
		
		PaidBy: req.PaidBy,
		
		Notes: req.Notes,
		
		CalculationNotes: req.CalculationNotes,
		
		Metadata: req.Metadata,
		
		CreatedBy: req.CreatedBy,
		
		UpdatedBy: req.UpdatedBy,
		
		'sale',: req.'sale',,
		
		'percentage',: req.'percentage',,
		
		'pending',: req.'pending',,
		
		SalesAmount: req.SalesAmount,
		
	}

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to create staff_commissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created staff_commissions",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a staff_commissions by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*StaffCommissionsResponse, error) {
	s.logger.Debug("getting staff_commissions",
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
		return nil, fmt.Errorf("failed to get staff_commissions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("staff_commissions not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of staff_commissions records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*StaffCommissionsListResponse, error) {
	s.logger.Debug("listing staff_commissions",
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
		return nil, fmt.Errorf("failed to list staff_commissions: %w", err)
	}

	// Convert to response
	items := make([]*StaffCommissionsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &StaffCommissionsListResponse{
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

// Update updates an existing staff_commissions
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateStaffCommissionsRequest) (*StaffCommissionsResponse, error) {
	s.logger.Info("updating staff_commissions",
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
		return nil, fmt.Errorf("failed to get staff_commissions: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("staff_commissions not found or access denied")
	}
	

	// Update fields
	
	if req.LocationId != nil {
		entity.LocationId = *req.LocationId
	}
	
	if req.EmployeeId != nil {
		entity.EmployeeId = *req.EmployeeId
	}
	
	if req.CommissionDate != nil {
		entity.CommissionDate = *req.CommissionDate
	}
	
	if req.PeriodStart != nil {
		entity.PeriodStart = *req.PeriodStart
	}
	
	if req.PeriodEnd != nil {
		entity.PeriodEnd = *req.PeriodEnd
	}
	
	if req.SourceType != nil {
		entity.SourceType = *req.SourceType
	}
	
	if req.SourceSaleId != nil {
		entity.SourceSaleId = *req.SourceSaleId
	}
	
	if req.SourceOrderId != nil {
		entity.SourceOrderId = *req.SourceOrderId
	}
	
	if req.CommissionType != nil {
		entity.CommissionType = *req.CommissionType
	}
	
	if req.CommissionRate != nil {
		entity.CommissionRate = *req.CommissionRate
	}
	
	if req.SalesAmount != nil {
		entity.SalesAmount = *req.SalesAmount
	}
	
	if req.CommissionAmount != nil {
		entity.CommissionAmount = *req.CommissionAmount
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
	
	if req.PaymentDate != nil {
		entity.PaymentDate = *req.PaymentDate
	}
	
	if req.PaymentMethod != nil {
		entity.PaymentMethod = *req.PaymentMethod
	}
	
	if req.PaidBy != nil {
		entity.PaidBy = *req.PaidBy
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
	}
	
	if req.CalculationNotes != nil {
		entity.CalculationNotes = *req.CalculationNotes
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
	
	if req.'sale', != nil {
		entity.'sale', = *req.'sale',
	}
	
	if req.'percentage', != nil {
		entity.'percentage', = *req.'percentage',
	}
	
	if req.'pending', != nil {
		entity.'pending', = *req.'pending',
	}
	
	if req.SalesAmount != nil {
		entity.SalesAmount = *req.SalesAmount
	}
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update staff_commissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated staff_commissions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a staff_commissions
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting staff_commissions",
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
		return fmt.Errorf("failed to get staff_commissions: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("staff_commissions not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete staff_commissions: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted staff_commissions",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *StaffCommissions) *StaffCommissionsResponse {
	return &StaffCommissionsResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		LocationId: entity.LocationId,
		
		EmployeeId: entity.EmployeeId,
		
		CommissionDate: entity.CommissionDate,
		
		PeriodStart: entity.PeriodStart,
		
		PeriodEnd: entity.PeriodEnd,
		
		SourceType: entity.SourceType,
		
		SourceSaleId: entity.SourceSaleId,
		
		SourceOrderId: entity.SourceOrderId,
		
		CommissionType: entity.CommissionType,
		
		CommissionRate: entity.CommissionRate,
		
		SalesAmount: entity.SalesAmount,
		
		CommissionAmount: entity.CommissionAmount,
		
		Status: entity.Status,
		
		ApprovedBy: entity.ApprovedBy,
		
		ApprovedAt: entity.ApprovedAt,
		
		PaymentDate: entity.PaymentDate,
		
		PaymentMethod: entity.PaymentMethod,
		
		PaidBy: entity.PaidBy,
		
		Notes: entity.Notes,
		
		CalculationNotes: entity.CalculationNotes,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
		DeletedAt: entity.DeletedAt,
		
		'sale',: entity.'sale',,
		
		'percentage',: entity.'percentage',,
		
		'pending',: entity.'pending',,
		
		SalesAmount: entity.SalesAmount,
		
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


// validateBusinessRules validates business rules for staff_commissions
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *StaffCommissions) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a staff_commissions can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
