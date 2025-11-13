package customer

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Customers
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Customers service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new customers
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, req *CreateCustomersRequest) (*CustomersResponse, error) {
	s.logger.Info("creating customers",
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
	entity := &Customers{
		OrganizationId: orgID,
		
		CustomerCode: req.CustomerCode,
		
		FirstName: req.FirstName,
		
		LastName: req.LastName,
		
		CompanyName: req.CompanyName,
		
		Email: req.Email,
		
		Phone: req.Phone,
		
		AlternatePhone: req.AlternatePhone,
		
		AddressLine1: req.AddressLine1,
		
		AddressLine2: req.AddressLine2,
		
		City: req.City,
		
		State: req.State,
		
		Country: req.Country,
		
		PostalCode: req.PostalCode,
		
		DateOfBirth: req.DateOfBirth,
		
		Gender: req.Gender,
		
		TaxNumber: req.TaxNumber,
		
		LoyaltyPoints: req.LoyaltyPoints,
		
		LoyaltyTier: req.LoyaltyTier,
		
		CreditLimit: req.CreditLimit,
		
		OutstandingBalance: req.OutstandingBalance,
		
		TotalPurchases: req.TotalPurchases,
		
		TotalOrders: req.TotalOrders,
		
		LastPurchaseAt: req.LastPurchaseAt,
		
		IsActive: req.IsActive,
		
		Notes: req.Notes,
		
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
		return nil, fmt.Errorf("failed to create customers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created customers",
		zap.String("id", entity.Id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a customers by ID
func (s *Service) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomersResponse, error) {
	s.logger.Debug("getting customers",
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
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("customers not found or access denied")
	}
	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of customers records
func (s *Service) List(ctx context.Context, orgID uuid.UUID, page, limit int) (*CustomersListResponse, error) {
	s.logger.Debug("listing customers",
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
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	// Convert to response
	items := make([]*CustomersResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &CustomersListResponse{
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

// Update updates an existing customers
func (s *Service) Update(ctx context.Context, orgID uuid.UUID, id uuid.UUID, req *UpdateCustomersRequest) (*CustomersResponse, error) {
	s.logger.Info("updating customers",
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
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	
	// Verify ownership
	if entity.OrganizationId != orgID {
		return nil, fmt.Errorf("customers not found or access denied")
	}
	

	// Update fields
	
	if req.CustomerCode != nil {
		entity.CustomerCode = req.CustomerCode
	}
	
	if req.FirstName != nil {
		entity.FirstName = req.FirstName
	}
	
	if req.LastName != nil {
		entity.LastName = req.LastName
	}
	
	if req.CompanyName != nil {
		entity.CompanyName = req.CompanyName
	}
	
	if req.Email != nil {
		entity.Email = req.Email
	}
	
	if req.Phone != nil {
		entity.Phone = req.Phone
	}
	
	if req.AlternatePhone != nil {
		entity.AlternatePhone = req.AlternatePhone
	}
	
	if req.AddressLine1 != nil {
		entity.AddressLine1 = req.AddressLine1
	}
	
	if req.AddressLine2 != nil {
		entity.AddressLine2 = req.AddressLine2
	}
	
	if req.City != nil {
		entity.City = req.City
	}
	
	if req.State != nil {
		entity.State = req.State
	}
	
	if req.Country != nil {
		entity.Country = req.Country
	}
	
	if req.PostalCode != nil {
		entity.PostalCode = req.PostalCode
	}
	
	if req.DateOfBirth != nil {
		entity.DateOfBirth = req.DateOfBirth
	}
	
	if req.Gender != nil {
		entity.Gender = req.Gender
	}
	
	if req.TaxNumber != nil {
		entity.TaxNumber = req.TaxNumber
	}
	
	if req.LoyaltyPoints != nil {
		entity.LoyaltyPoints = req.LoyaltyPoints
	}
	
	if req.LoyaltyTier != nil {
		entity.LoyaltyTier = req.LoyaltyTier
	}
	
	if req.CreditLimit != nil {
		entity.CreditLimit = req.CreditLimit
	}
	
	if req.OutstandingBalance != nil {
		entity.OutstandingBalance = req.OutstandingBalance
	}
	
	if req.TotalPurchases != nil {
		entity.TotalPurchases = req.TotalPurchases
	}
	
	if req.TotalOrders != nil {
		entity.TotalOrders = req.TotalOrders
	}
	
	if req.LastPurchaseAt != nil {
		entity.LastPurchaseAt = req.LastPurchaseAt
	}
	
	if req.IsActive != nil {
		entity.IsActive = req.IsActive
	}
	
	if req.Notes != nil {
		entity.Notes = req.Notes
	}
	
	if req.CustomFields != nil {
		entity.CustomFields = req.CustomFields
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
	

	// Business logic validation
	if err := s.validateBusinessRules(ctx, tx, entity); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, tx, entity); err != nil {
		return nil, fmt.Errorf("failed to update customers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated customers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a customers
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	s.logger.Info("deleting customers",
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
		return fmt.Errorf("failed to get customers: %w", err)
	}

	if entity.OrganizationId != orgID {
		return fmt.Errorf("customers not found or access denied")
	}
	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete customers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted customers",
		zap.String("id", id.String()),
		zap.String("organization_id", orgID.String()),
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Customers) *CustomersResponse {
	return &CustomersResponse{
		
		Id: entity.Id,
		
		OrganizationId: entity.OrganizationId,
		
		CustomerCode: entity.CustomerCode,
		
		FirstName: entity.FirstName,
		
		LastName: entity.LastName,
		
		CompanyName: entity.CompanyName,
		
		Email: entity.Email,
		
		Phone: entity.Phone,
		
		AlternatePhone: entity.AlternatePhone,
		
		AddressLine1: entity.AddressLine1,
		
		AddressLine2: entity.AddressLine2,
		
		City: entity.City,
		
		State: entity.State,
		
		Country: entity.Country,
		
		PostalCode: entity.PostalCode,
		
		DateOfBirth: entity.DateOfBirth,
		
		Gender: entity.Gender,
		
		TaxNumber: entity.TaxNumber,
		
		LoyaltyPoints: entity.LoyaltyPoints,
		
		LoyaltyTier: entity.LoyaltyTier,
		
		CreditLimit: entity.CreditLimit,
		
		OutstandingBalance: entity.OutstandingBalance,
		
		TotalPurchases: entity.TotalPurchases,
		
		TotalOrders: entity.TotalOrders,
		
		LastPurchaseAt: entity.LastPurchaseAt,
		
		IsActive: entity.IsActive,
		
		Notes: entity.Notes,
		
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


// validateBusinessRules validates business rules for customers
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Customers) error {
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

// canDelete checks if a customers can be deleted
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
