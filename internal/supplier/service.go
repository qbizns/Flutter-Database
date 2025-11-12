package supplier

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/your-org/pos-backend/internal/dto/supplier"
	"github.com/your-org/pos-backend/internal/logging"
	"github.com/your-org/pos-backend/internal/repository/supplier"
	"go.uber.org/zap"
)

// Service handles business logic for Suppliers
type Service struct {
	repo   *supplier.Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Suppliers service
func NewService(repo *supplier.Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new suppliers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *dto.CreateSuppliersRequest) (*dto.SuppliersResponse, error) {
	s.logger.Info("creating suppliers",
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
	entity := &supplier.Suppliers{
		OrganizationID: orgID,
		
		SupplierCode: req.SupplierCode,
		
		Name: req.Name,
		
		ContactPerson: req.ContactPerson,
		
		Email: req.Email,
		
		Phone: req.Phone,
		
		Address: req.Address,
		
		City: req.City,
		
		State: req.State,
		
		Country: req.Country,
		
		PostalCode: req.PostalCode,
		
		TaxNumber: req.TaxNumber,
		
		PaymentTerms: req.PaymentTerms,
		
		CreditLimit: req.CreditLimit,
		
		OutstandingBalance: req.OutstandingBalance,
		
		TotalPurchases: req.TotalPurchases,
		
		TotalOrders: req.TotalOrders,
		
		LastOrderDate: req.LastOrderDate,
		
		Status: req.Status,
		
		Notes: req.Notes,
		
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
		return nil, fmt.Errorf("failed to create suppliers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created suppliers",
		zap.String("id", entity.ID.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a suppliers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*dto.SuppliersResponse, error) {
	s.logger.Debug("getting suppliers",
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
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("suppliers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of suppliers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*dto.SuppliersListResponse, error) {
	s.logger.Debug("listing suppliers",
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
		return nil, fmt.Errorf("failed to list suppliers: %w", err)
	}

	// Convert to response
	items := make([]*dto.SuppliersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &dto.SuppliersListResponse{
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

// Update updates an existing suppliers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *dto.UpdateSuppliersRequest) (*dto.SuppliersResponse, error) {
	s.logger.Info("updating suppliers",
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
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationID != orgID {
		return nil, fmt.Errorf("suppliers not found or access denied")
	}
	

	// Update fields
	
	if req.SupplierCode != nil {
		entity.SupplierCode = *req.SupplierCode
	}
	
	if req.Name != nil {
		entity.Name = *req.Name
	}
	
	if req.ContactPerson != nil {
		entity.ContactPerson = *req.ContactPerson
	}
	
	if req.Email != nil {
		entity.Email = *req.Email
	}
	
	if req.Phone != nil {
		entity.Phone = *req.Phone
	}
	
	if req.Address != nil {
		entity.Address = *req.Address
	}
	
	if req.City != nil {
		entity.City = *req.City
	}
	
	if req.State != nil {
		entity.State = *req.State
	}
	
	if req.Country != nil {
		entity.Country = *req.Country
	}
	
	if req.PostalCode != nil {
		entity.PostalCode = *req.PostalCode
	}
	
	if req.TaxNumber != nil {
		entity.TaxNumber = *req.TaxNumber
	}
	
	if req.PaymentTerms != nil {
		entity.PaymentTerms = *req.PaymentTerms
	}
	
	if req.CreditLimit != nil {
		entity.CreditLimit = *req.CreditLimit
	}
	
	if req.OutstandingBalance != nil {
		entity.OutstandingBalance = *req.OutstandingBalance
	}
	
	if req.TotalPurchases != nil {
		entity.TotalPurchases = *req.TotalPurchases
	}
	
	if req.TotalOrders != nil {
		entity.TotalOrders = *req.TotalOrders
	}
	
	if req.LastOrderDate != nil {
		entity.LastOrderDate = *req.LastOrderDate
	}
	
	if req.Status != nil {
		entity.Status = *req.Status
	}
	
	if req.Notes != nil {
		entity.Notes = *req.Notes
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
		return nil, fmt.Errorf("failed to update suppliers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated suppliers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a suppliers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting suppliers",
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
		return fmt.Errorf("failed to get suppliers: %w", err)
	}

	if entity.OrganizationID != orgID {
		return fmt.Errorf("suppliers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete suppliers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted suppliers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *supplier.Suppliers) *dto.SuppliersResponse {
	return &dto.SuppliersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		SupplierCode: entity.SupplierCode,
		
		Name: entity.Name,
		
		ContactPerson: entity.ContactPerson,
		
		Email: entity.Email,
		
		Phone: entity.Phone,
		
		Address: entity.Address,
		
		City: entity.City,
		
		State: entity.State,
		
		Country: entity.Country,
		
		PostalCode: entity.PostalCode,
		
		TaxNumber: entity.TaxNumber,
		
		PaymentTerms: entity.PaymentTerms,
		
		CreditLimit: entity.CreditLimit,
		
		OutstandingBalance: entity.OutstandingBalance,
		
		TotalPurchases: entity.TotalPurchases,
		
		TotalOrders: entity.TotalOrders,
		
		LastOrderDate: entity.LastOrderDate,
		
		Status: entity.Status,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for suppliers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *supplier.Suppliers) error {
	// TODO: Add business rule validation
	// Example:
	// - Check for duplicate names within organization
	// - Validate foreign key references exist
	// - Check status transitions are valid
	// - Validate amounts are positive
	// - etc.

	return nil
}

// canDelete checks if a suppliers can be deleted
func (s *Service) canDelete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	// TODO: Add delete validation
	// Example:
	// - Check for dependent records
	// - Verify not referenced by other entities
	// - Check business rules allow deletion
	// - etc.

	return nil
}
