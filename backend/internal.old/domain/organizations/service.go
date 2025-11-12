package organizations

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles organization business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new organization service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of organizations with filters
func (s *Service) List(ctx context.Context, filters OrganizationFilters) ([]Organization, error) {
	orgs, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.Error("failed to list organizations", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return orgs, nil
}

// Count counts organizations matching filters
func (s *Service) Count(ctx context.Context, filters OrganizationFilters) (int64, error) {
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.logger.Error("failed to count organizations", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new organization
func (s *Service) Create(ctx context.Context, org *Organization) error {
	// Validate organization
	if err := s.validate(org); err != nil {
		return err
	}

	// Check for duplicate slug
	existing, err := s.repo.GetBySlug(ctx, org.Slug)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate slug", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("organization", "slug already exists")
	}

	// Check for duplicate email if provided
	if org.Email != "" {
		// Note: Email uniqueness is not enforced at DB level for organizations
		// but we could add a check here if needed
	}

	// Set defaults
	org.ID = uuid.New()
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()
	if org.Status == "" {
		org.Status = StatusTrial
	}
	if org.Plan == "" {
		org.Plan = "basic"
	}
	if org.MaxUsers == 0 {
		org.MaxUsers = 5
	}
	if org.MaxProducts == 0 {
		org.MaxProducts = 1000
	}
	if org.MaxLocations == 0 {
		org.MaxLocations = 1
	}

	// Create organization
	if err := s.repo.Create(ctx, org); err != nil {
		s.logger.Error("failed to create organization", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves an organization by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Organization, error) {
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get organization", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if org == nil {
		return nil, apperrors.NotFound("organization")
	}

	return org, nil
}

// GetBySlug retrieves an organization by slug
func (s *Service) GetBySlug(ctx context.Context, slug string) (*Organization, error) {
	org, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		s.logger.Error("failed to get organization by slug", zap.Error(err), zap.String("slug", slug))
		return nil, apperrors.DatabaseError(err)
	}

	if org == nil {
		return nil, apperrors.NotFound("organization")
	}

	return org, nil
}

// Update updates an existing organization
func (s *Service) Update(ctx context.Context, org *Organization) error {
	// Validate organization
	if err := s.validate(org); err != nil {
		return err
	}

	// Check if organization exists
	existing, err := s.repo.Get(ctx, org.ID)
	if err != nil {
		s.logger.Error("failed to get organization", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("organization")
	}

	// Check for duplicate slug (if slug changed)
	if org.Slug != existing.Slug {
		duplicate, err := s.repo.GetBySlug(ctx, org.Slug)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate slug", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != org.ID {
			return apperrors.AlreadyExists("organization", "slug already exists")
		}
	}

	// Update timestamp
	org.UpdatedAt = time.Now()

	// Update organization
	if err := s.repo.Update(ctx, org); err != nil {
		s.logger.Error("failed to update organization", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes an organization
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if organization exists
	org, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get organization", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if org == nil {
		return apperrors.NotFound("organization")
	}

	// Soft delete organization
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete organization", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validate validates an organization
func (s *Service) validate(org *Organization) error {
	if org.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if len(org.Name) > 255 {
		return apperrors.ValidationFailed("name cannot exceed 255 characters")
	}
	if org.Slug == "" {
		return apperrors.ValidationFailed("slug is required")
	}
	if len(org.Slug) > 100 {
		return apperrors.ValidationFailed("slug cannot exceed 100 characters")
	}
	if org.MaxUsers <= 0 {
		return apperrors.ValidationFailed("max_users must be positive")
	}
	if org.MaxProducts <= 0 {
		return apperrors.ValidationFailed("max_products must be positive")
	}
	if org.MaxLocations <= 0 {
		return apperrors.ValidationFailed("max_locations must be positive")
	}

	return nil
}
