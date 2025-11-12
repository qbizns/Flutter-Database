package customers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles customer business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new customer service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of customers with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters CustomerFilters) ([]Customer, error) {
	customers, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list customers", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return customers, nil
}

// Create creates a new customer
func (s *Service) Create(ctx context.Context, customer *Customer) error {
	// Validate customer
	if err := s.validate(customer); err != nil {
		return err
	}

	// Check for duplicate email
	existing, err := s.repo.GetByEmail(ctx, customer.OrganizationID, customer.Email)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate email", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("customer", "Email already exists")
	}

	// Set defaults
	customer.ID = uuid.New()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.CurrentBalance = 0
	customer.TotalSpent = 0
	customer.TotalVisits = 0
	customer.LoyaltyPoints = 0

	// Create customer
	if err := s.repo.Create(ctx, customer); err != nil {
		s.logger.Error("failed to create customer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a customer by ID
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Customer, error) {
	customer, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get customer", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if customer == nil {
		return nil, apperrors.NotFound("customer")
	}

	return customer, nil
}

// Update updates an existing customer
func (s *Service) Update(ctx context.Context, customer *Customer) error {
	// Validate customer
	if err := s.validate(customer); err != nil {
		return err
	}

	// Check if customer exists
	existing, err := s.repo.Get(ctx, customer.OrganizationID, customer.ID)
	if err != nil {
		s.logger.Error("failed to get customer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("customer")
	}

	// Check for duplicate email (if email changed)
	if customer.Email != existing.Email {
		duplicate, err := s.repo.GetByEmail(ctx, customer.OrganizationID, customer.Email)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate email", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != customer.ID {
			return apperrors.AlreadyExists("customer", "Email already exists")
		}
	}

	// Update timestamp
	customer.UpdatedAt = time.Now()

	// Update customer
	if err := s.repo.Update(ctx, customer); err != nil {
		s.logger.Error("failed to update customer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a customer
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if customer exists
	customer, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get customer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if customer == nil {
		return apperrors.NotFound("customer")
	}

	// Soft delete customer
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete customer", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validate validates a customer
func (s *Service) validate(customer *Customer) error {
	if customer.FirstName == "" {
		return apperrors.ValidationFailed("first name is required")
	}
	if customer.Email == "" {
		return apperrors.ValidationFailed("email is required")
	}
	if customer.CreditLimit < 0 {
		return apperrors.ValidationFailed("credit limit must be non-negative")
	}

	return nil
}
