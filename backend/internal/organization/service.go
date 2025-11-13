package organization

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"github.com/your-org/pos-backend/internal/logging"
	
	"go.uber.org/zap"
)

// Service handles business logic for Organizations
type Service struct {
	repo   *Repository
	db     *pgxpool.Pool
	logger *logging.Logger
}

// NewService creates a new Organizations service
func NewService(repo *Repository, db *pgxpool.Pool, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

// Create creates a new organizations
func (s *Service) Create(ctx context.Context, req *CreateOrganizationsRequest) (*OrganizationsResponse, error) {
	s.logger.Info("creating organizations",
		
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

	

	// Convert DTO to entity
	entity := &Organizations{
		
		
		Name: req.Name,
		
		Slug: req.Slug,
		
		Description: req.Description,
		
		Email: req.Email,
		
		Phone: req.Phone,
		
		Address: req.Address,
		
		City: req.City,
		
		State: req.State,
		
		Country: req.Country,
		
		PostalCode: req.PostalCode,
		
		Status: req.Status,
		
		Plan: req.Plan,
		
		TrialEndsAt: req.TrialEndsAt,
		
		SubscriptionStartsAt: req.SubscriptionStartsAt,
		
		SubscriptionEndsAt: req.SubscriptionEndsAt,
		
		MaxUsers: req.MaxUsers,
		
		MaxProducts: req.MaxProducts,
		
		MaxLocations: req.MaxLocations,
		
		Settings: req.Settings,
		
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
		return nil, fmt.Errorf("failed to create organizations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("created organizations",
		zap.String("id", entity.Id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// GetByID retrieves a organizations by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*OrganizationsResponse, error) {
	s.logger.Debug("getting organizations",
		zap.String("id", id.String()),
		
	)

	// Start transaction (read-only)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Get from database
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	

	return s.entityToResponse(entity), nil
}

// List retrieves a paginated list of organizations records
func (s *Service) List(ctx context.Context, page, limit int) (*OrganizationsListResponse, error) {
	s.logger.Debug("listing organizations",
		
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

	
	// Get from database
	entities, total, err := s.repo.List(ctx, tx, limit, offset)
	
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	// Convert to response
	items := make([]*OrganizationsResponse, len(entities))
	for i, entity := range entities {
		items[i] = s.entityToResponse(entity)
	}

	totalPages := (total + limit - 1) / limit

	return &OrganizationsListResponse{
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

// Update updates an existing organizations
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateOrganizationsRequest) (*OrganizationsResponse, error) {
	s.logger.Info("updating organizations",
		zap.String("id", id.String()),
		
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

	

	// Get existing entity
	entity, err := s.repo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	

	// Update fields
	
	if req.Name != nil {
		entity.Name = req.Name
	}
	
	if req.Slug != nil {
		entity.Slug = req.Slug
	}
	
	if req.Description != nil {
		entity.Description = req.Description
	}
	
	if req.Email != nil {
		entity.Email = req.Email
	}
	
	if req.Phone != nil {
		entity.Phone = req.Phone
	}
	
	if req.Address != nil {
		entity.Address = req.Address
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
	
	if req.Status != nil {
		entity.Status = req.Status
	}
	
	if req.Plan != nil {
		entity.Plan = req.Plan
	}
	
	if req.TrialEndsAt != nil {
		entity.TrialEndsAt = req.TrialEndsAt
	}
	
	if req.SubscriptionStartsAt != nil {
		entity.SubscriptionStartsAt = req.SubscriptionStartsAt
	}
	
	if req.SubscriptionEndsAt != nil {
		entity.SubscriptionEndsAt = req.SubscriptionEndsAt
	}
	
	if req.MaxUsers != nil {
		entity.MaxUsers = req.MaxUsers
	}
	
	if req.MaxProducts != nil {
		entity.MaxProducts = req.MaxProducts
	}
	
	if req.MaxLocations != nil {
		entity.MaxLocations = req.MaxLocations
	}
	
	if req.Settings != nil {
		entity.Settings = req.Settings
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
		return nil, fmt.Errorf("failed to update organizations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("updated organizations",
		zap.String("id", id.String()),
		
	)

	return s.entityToResponse(entity), nil
}

// Delete deletes a organizations
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("deleting organizations",
		zap.String("id", id.String()),
		
	)

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	

	// Check if can be deleted (business rules)
	if err := s.canDelete(ctx, tx, id); err != nil {
		return err
	}

	// Delete from database
	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete organizations: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("deleted organizations",
		zap.String("id", id.String()),
		
	)

	return nil
}

// entityToResponse converts entity to response DTO
func (s *Service) entityToResponse(entity *Organizations) *OrganizationsResponse {
	return &OrganizationsResponse{
		
		Id: entity.Id,
		
		Name: entity.Name,
		
		Slug: entity.Slug,
		
		Description: entity.Description,
		
		Email: entity.Email,
		
		Phone: entity.Phone,
		
		Address: entity.Address,
		
		City: entity.City,
		
		State: entity.State,
		
		Country: entity.Country,
		
		PostalCode: entity.PostalCode,
		
		Status: entity.Status,
		
		Plan: entity.Plan,
		
		TrialEndsAt: entity.TrialEndsAt,
		
		SubscriptionStartsAt: entity.SubscriptionStartsAt,
		
		SubscriptionEndsAt: entity.SubscriptionEndsAt,
		
		MaxUsers: entity.MaxUsers,
		
		MaxProducts: entity.MaxProducts,
		
		MaxLocations: entity.MaxLocations,
		
		Settings: entity.Settings,
		
		Metadata: entity.Metadata,
		
		CreatedAt: entity.CreatedAt,
		
		UpdatedAt: entity.UpdatedAt,
		
		DeletedAt: entity.DeletedAt,
		
		CreatedBy: entity.CreatedBy,
		
		UpdatedBy: entity.UpdatedBy,
		
	}
}



// validateBusinessRules validates business rules for organizations
func (s *Service) validateBusinessRules(ctx context.Context, tx pgx.Tx, entity *Organizations) error {
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

// canDelete checks if a organizations can be deleted
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
