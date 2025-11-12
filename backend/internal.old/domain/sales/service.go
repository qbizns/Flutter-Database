package sales

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

type Service struct {
	repo   Repository
	logger *logging.Logger
}

func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters SaleFilters) ([]Sale, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	return s.repo.List(ctx, orgID, filters)
}

func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters SaleFilters) (int64, error) {
	return s.repo.Count(ctx, orgID, filters)
}

func (s *Service) Create(ctx context.Context, sale *Sale) error {
	if err := s.validate(sale); err != nil {
		return err
	}

	// Generate ID if not set
	if sale.ID == uuid.Nil {
		sale.ID = uuid.New()
	}

	// Set timestamps
	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now

	// Validate items
	if len(sale.Items) == 0 {
		return apperrors.ValidationFailed("sale must have at least one item")
	}

	for i := range sale.Items {
		if sale.Items[i].ID == uuid.Nil {
			sale.Items[i].ID = uuid.New()
		}
		sale.Items[i].SaleID = sale.ID
		sale.Items[i].OrganizationID = sale.OrganizationID
		sale.Items[i].CreatedAt = now
		sale.Items[i].UpdatedAt = now
	}

	if err := s.repo.Create(ctx, sale); err != nil {
		s.logger.Error("failed to create sale", zap.Error(err))
		return apperrors.Internal("failed to create sale")
	}

	return nil
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Sale, error) {
	sale, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, apperrors.NotFound("sale")
	}
	return sale, nil
}

func (s *Service) GetBySaleNumber(ctx context.Context, orgID uuid.UUID, saleNumber string) (*Sale, error) {
	sale, err := s.repo.GetBySaleNumber(ctx, orgID, saleNumber)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, apperrors.NotFound("sale")
	}
	return sale, nil
}

func (s *Service) GetWithItems(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) (*Sale, error) {
	sale, err := s.repo.GetWithItems(ctx, orgID, saleID)
	if err != nil {
		return nil, err
	}
	if sale == nil {
		return nil, apperrors.NotFound("sale")
	}
	return sale, nil
}

func (s *Service) Update(ctx context.Context, sale *Sale) error {
	if err := s.validate(sale); err != nil {
		return err
	}

	// Check if sale exists
	existing, err := s.repo.Get(ctx, sale.OrganizationID, sale.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.NotFound("sale")
	}

	// Update timestamp
	sale.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, sale); err != nil {
		s.logger.Error("failed to update sale", zap.Error(err))
		return apperrors.Internal("failed to update sale")
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if sale exists
	existing, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return apperrors.NotFound("sale")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete sale", zap.Error(err))
		return apperrors.Internal("failed to delete sale")
	}

	return nil
}

func (s *Service) ListItems(ctx context.Context, saleID uuid.UUID, filters SaleItemFilters) ([]SaleItem, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 20
	}

	return s.repo.ListItems(ctx, saleID, filters)
}

func (s *Service) validate(sale *Sale) error {
	if sale.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}

	if sale.SaleNumber == "" {
		return apperrors.ValidationFailed("sale_number is required")
	}

	if sale.TransactionType == "" {
		return apperrors.ValidationFailed("transaction_type is required")
	}

	if sale.TotalAmount < 0 {
		return apperrors.ValidationFailed("total_amount cannot be negative")
	}

	if sale.Subtotal < 0 {
		return apperrors.ValidationFailed("subtotal cannot be negative")
	}

	if sale.TaxAmount < 0 {
		return apperrors.ValidationFailed("tax_amount cannot be negative")
	}

	if sale.DiscountAmount < 0 {
		return apperrors.ValidationFailed("discount_amount cannot be negative")
	}

	return nil
}
