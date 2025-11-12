package products

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles product business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new product service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of products with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters ProductFilters) ([]Product, error) {
	products, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list products", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return products, nil
}

// Create creates a new product
func (s *Service) Create(ctx context.Context, product *Product) error {
	// Validate product
	if err := s.validate(product); err != nil {
		return err
	}

	// Check for duplicate SKU
	existing, err := s.repo.GetBySKU(ctx, product.OrganizationID, product.SKU)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate SKU", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("product", "SKU already exists")
	}

	// Set defaults
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if product.CurrentStock == 0 {
		product.CurrentStock = 0
	}

	// Create product
	if err := s.repo.Create(ctx, product); err != nil {
		s.logger.Error("failed to create product", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a product by ID
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Product, error) {
	product, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if product == nil {
		return nil, apperrors.NotFound("product")
	}

	return product, nil
}

// Update updates an existing product
func (s *Service) Update(ctx context.Context, product *Product) error {
	// Validate product
	if err := s.validate(product); err != nil {
		return err
	}

	// Check if product exists
	existing, err := s.repo.Get(ctx, product.OrganizationID, product.ID)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("product")
	}

	// Check for duplicate SKU (if SKU changed)
	if product.SKU != existing.SKU {
		duplicate, err := s.repo.GetBySKU(ctx, product.OrganizationID, product.SKU)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate SKU", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != product.ID {
			return apperrors.AlreadyExists("product", "SKU already exists")
		}
	}

	// Update timestamp
	product.UpdatedAt = time.Now()

	// Update product
	if err := s.repo.Update(ctx, product); err != nil {
		s.logger.Error("failed to update product", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete soft-deletes a product
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if product exists
	product, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get product", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if product == nil {
		return apperrors.NotFound("product")
	}

	// Soft delete product
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete product", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validate validates a product
func (s *Service) validate(product *Product) error {
	if product.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if product.SKU == "" {
		return apperrors.ValidationFailed("SKU is required")
	}
	if product.UnitPrice < 0 {
		return apperrors.ValidationFailed("unit price must be non-negative")
	}
	if product.Cost < 0 {
		return apperrors.ValidationFailed("cost must be non-negative")
	}
	if product.TaxRate < 0 || product.TaxRate > 100 {
		return apperrors.ValidationFailed("tax rate must be between 0 and 100")
	}

	return nil
}
