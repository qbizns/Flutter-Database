package promotions

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
	return &Service{repo: repo, logger: logger}
}

// Promotions

func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters PromotionFilters) ([]Promotion, error) {
	promotions, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list promotions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return promotions, nil
}

func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters PromotionFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count promotions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *Service) Create(ctx context.Context, promotion *Promotion) error {
	if err := s.validatePromotion(promotion); err != nil {
		return err
	}

	// Check for duplicate code
	existing, err := s.repo.GetByCode(ctx, promotion.OrganizationID, promotion.PromotionCode)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("promotion", "Code already exists")
	}

	promotion.ID = uuid.New()
	promotion.CreatedAt = time.Now()
	promotion.UpdatedAt = time.Now()
	promotion.CurrentUses = 0

	if promotion.IsActive == false && promotion.PromotionType == "" {
		promotion.IsActive = true
	}

	if err := s.repo.Create(ctx, promotion); err != nil {
		s.logger.Error("failed to create promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Promotion, error) {
	promotion, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get promotion", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if promotion == nil {
		return nil, apperrors.NotFound("promotion")
	}
	return promotion, nil
}

func (s *Service) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*Promotion, error) {
	promotion, err := s.repo.GetByCode(ctx, orgID, code)
	if err != nil {
		s.logger.Error("failed to get promotion by code", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	if promotion == nil {
		return nil, apperrors.NotFound("promotion")
	}
	return promotion, nil
}

func (s *Service) Update(ctx context.Context, promotion *Promotion) error {
	if err := s.validatePromotion(promotion); err != nil {
		return err
	}

	existing, err := s.repo.Get(ctx, promotion.OrganizationID, promotion.ID)
	if err != nil {
		s.logger.Error("failed to get promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("promotion")
	}

	if promotion.PromotionCode != existing.PromotionCode {
		duplicate, err := s.repo.GetByCode(ctx, promotion.OrganizationID, promotion.PromotionCode)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate code", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != promotion.ID {
			return apperrors.AlreadyExists("promotion", "Code already exists")
		}
	}

	promotion.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, promotion); err != nil {
		s.logger.Error("failed to update promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	promotion, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if promotion == nil {
		return apperrors.NotFound("promotion")
	}

	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// Promotion Usage

func (s *Service) ListUsage(ctx context.Context, orgID uuid.UUID, filters PromotionUsageFilters) ([]PromotionUsage, error) {
	usage, err := s.repo.ListUsage(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list promotion usage", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return usage, nil
}

func (s *Service) CountUsage(ctx context.Context, orgID uuid.UUID, filters PromotionUsageFilters) (int64, error) {
	count, err := s.repo.CountUsage(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count promotion usage", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}
	return count, nil
}

func (s *Service) RecordUsage(ctx context.Context, orgID uuid.UUID, promotionID uuid.UUID, saleID *uuid.UUID, customerID *uuid.UUID, discountAmount float64) error {
	promotion, err := s.repo.Get(ctx, orgID, promotionID)
	if err != nil {
		s.logger.Error("failed to get promotion", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if promotion == nil {
		return apperrors.NotFound("promotion")
	}

	// Check usage limits
	if promotion.MaxUsesTotal != nil && promotion.CurrentUses >= *promotion.MaxUsesTotal {
		return apperrors.ValidationFailed("promotion has reached maximum usage limit")
	}

	if promotion.MaxUsesPerCustomer != nil && customerID != nil {
		usageCount, err := s.repo.GetUsageByPromotion(ctx, orgID, promotionID, customerID)
		if err != nil {
			s.logger.Error("failed to check customer usage", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if usageCount >= int64(*promotion.MaxUsesPerCustomer) {
			return apperrors.ValidationFailed("customer has reached maximum usage limit for this promotion")
		}
	}

	usage := &PromotionUsage{
		ID:             uuid.New(),
		OrganizationID: orgID,
		PromotionID:    promotionID,
		SaleID:         saleID,
		CustomerID:     customerID,
		DiscountAmount: discountAmount,
		UsedAt:         time.Now(),
		CreatedAt:      time.Now(),
	}

	if err := s.repo.RecordUsage(ctx, usage); err != nil {
		s.logger.Error("failed to record promotion usage", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	return nil
}

// Validation

func (s *Service) validatePromotion(p *Promotion) error {
	if p.PromotionCode == "" {
		return apperrors.ValidationFailed("promotion code is required")
	}
	if p.Name == "" {
		return apperrors.ValidationFailed("name is required")
	}
	if p.PromotionType == "" {
		return apperrors.ValidationFailed("promotion type is required")
	}

	validTypes := map[string]bool{
		"percentage":          true,
		"fixed_amount":        true,
		"buy_x_get_y":         true,
		"bundle":              true,
		"quantity_discount":   true,
	}
	if !validTypes[p.PromotionType] {
		return apperrors.ValidationFailed("invalid promotion type")
	}

	if p.DiscountValue < 0 {
		return apperrors.ValidationFailed("discount value must be non-negative")
	}

	validAppliesToTypes := map[string]bool{
		"all":                  true,
		"specific_products":    true,
		"specific_categories": true,
		"cart_total":           true,
	}
	if !validAppliesToTypes[p.AppliesToType] {
		return apperrors.ValidationFailed("invalid applies_to type")
	}

	if p.MinimumPurchaseAmount < 0 {
		return apperrors.ValidationFailed("minimum purchase amount must be non-negative")
	}

	if p.MinimumQuantity < 0 {
		return apperrors.ValidationFailed("minimum quantity must be non-negative")
	}

	if p.EndDate != nil && p.EndDate.Before(p.StartDate) {
		return apperrors.ValidationFailed("end date must be after start date")
	}

	if p.PromotionType == "buy_x_get_y" {
		if p.BuyQuantity == nil || *p.BuyQuantity <= 0 {
			return apperrors.ValidationFailed("buy quantity is required for buy_x_get_y type")
		}
		if p.GetQuantity == nil || *p.GetQuantity <= 0 {
			return apperrors.ValidationFailed("get quantity is required for buy_x_get_y type")
		}
		if p.GetDiscountPercentage == nil || *p.GetDiscountPercentage < 0 {
			return apperrors.ValidationFailed("get discount percentage is required for buy_x_get_y type")
		}
	}

	return nil
}
