package suppliers

import (
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Service struct {
	repo   Repository
	logger *logging.Logger
}

func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters SupplierFilters) ([]Supplier, error) {
	suppliers, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list suppliers", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return suppliers, nil
}

func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters SupplierFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count suppliers", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *Service) Create(ctx context.Context, supplier *Supplier) error {
	if err := s.validate(supplier); err != nil {
		return err
	}

	if supplier.Email != "" {
		existing, err := s.repo.GetByEmail(ctx, supplier.OrganizationID, supplier.Email)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate email", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if existing != nil {
			return apperrors.AlreadyExists("supplier", "Email already exists")
		}
	}

	supplier.ID = uuid.New()
	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()
	supplier.OutstandingBalance = 0
	supplier.TotalPurchases = 0
	supplier.TotalOrders = 0

	if supplier.Status == "" {
		supplier.Status = "active"
	}

	if err := s.repo.Create(ctx, supplier); err != nil {
		s.logger.Error("failed to create supplier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Supplier, error) {
	supplier, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get supplier", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}
	if supplier == nil {
		return nil, apperrors.NotFound("supplier")
	}
	return supplier, nil
}

func (s *Service) Update(ctx context.Context, supplier *Supplier) error {
	if err := s.validate(supplier); err != nil {
		return err
	}

	existing, err := s.repo.Get(ctx, supplier.OrganizationID, supplier.ID)
	if err != nil {
		s.logger.Error("failed to get supplier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("supplier")
	}

	if supplier.Email != "" && supplier.Email != existing.Email {
		duplicate, err := s.repo.GetByEmail(ctx, supplier.OrganizationID, supplier.Email)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate email", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != supplier.ID {
			return apperrors.AlreadyExists("supplier", "Email already exists")
		}
	}

	supplier.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, supplier); err != nil {
		s.logger.Error("failed to update supplier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	supplier, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get supplier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if supplier == nil {
		return apperrors.NotFound("supplier")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete supplier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) validate(supplier *Supplier) error {
	if supplier.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if supplier.Email == "" {
		return apperrors.ValidationFailed("email is required")
	}
	if !emailRegex.MatchString(supplier.Email) {
		return apperrors.ValidationFailed("invalid email format")
	}
	if supplier.CreditLimit < 0 {
		return apperrors.ValidationFailed("credit limit must be non-negative")
	}
	if supplier.Status != "" && supplier.Status != "active" && supplier.Status != "inactive" && supplier.Status != "suspended" {
		return apperrors.ValidationFailed("status must be active, inactive, or suspended")
	}
	return nil
}
