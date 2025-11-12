package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles payment business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new payment service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of payments with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters PaymentFilters) ([]Payment, error) {
	payments, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list payments", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return payments, nil
}

// Count counts payments matching filters
func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters PaymentFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count payments", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new payment
func (s *Service) Create(ctx context.Context, payment *Payment) error {
	// Validate payment
	if err := s.validate(payment); err != nil {
		return err
	}

	// Set defaults
	payment.ID = uuid.New()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	if payment.PaymentStatus == "" {
		payment.PaymentStatus = PaymentStatusPending
	}
	if payment.PaymentDate.IsZero() {
		payment.PaymentDate = time.Now()
	}

	// Create payment
	if err := s.repo.Create(ctx, payment); err != nil {
		s.logger.Error("failed to create payment", zap.Error(err), zap.String("sale_id", payment.SaleID.String()))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves a payment by ID
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*Payment, error) {
	payment, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get payment", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if payment == nil {
		return nil, apperrors.NotFound("payment")
	}

	return payment, nil
}

// GetBySale retrieves all payments for a sale
func (s *Service) GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]Payment, error) {
	payments, err := s.repo.GetBySale(ctx, orgID, saleID)
	if err != nil {
		s.logger.Error("failed to get payments by sale", zap.Error(err), zap.String("sale_id", saleID.String()))
		return nil, apperrors.DatabaseError(err)
	}

	return payments, nil
}

// Update updates an existing payment
func (s *Service) Update(ctx context.Context, payment *Payment) error {
	// Validate payment
	if err := s.validate(payment); err != nil {
		return err
	}

	// Check if payment exists
	existing, err := s.repo.Get(ctx, payment.OrganizationID, payment.ID)
	if err != nil {
		s.logger.Error("failed to get payment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("payment")
	}

	// Update timestamp
	payment.UpdatedAt = time.Now()

	// Update payment
	if err := s.repo.Update(ctx, payment); err != nil {
		s.logger.Error("failed to update payment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes a payment
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if payment exists
	payment, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get payment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if payment == nil {
		return apperrors.NotFound("payment")
	}

	// Can't delete completed or refunded payments
	if payment.PaymentStatus == PaymentStatusCompleted || payment.PaymentStatus == PaymentStatusRefunded {
		return apperrors.ValidationFailed(fmt.Sprintf("cannot delete %s payment", payment.PaymentStatus))
	}

	// Delete payment
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete payment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// UpdateStatus updates the payment status
func (s *Service) UpdateStatus(ctx context.Context, orgID uuid.UUID, paymentID uuid.UUID, status PaymentStatus) error {
	// Validate status
	if !s.isValidStatus(status) {
		return apperrors.ValidationFailed(fmt.Sprintf("invalid payment status: %s", status))
	}

	// Check if payment exists
	payment, err := s.repo.Get(ctx, orgID, paymentID)
	if err != nil {
		s.logger.Error("failed to get payment", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if payment == nil {
		return apperrors.NotFound("payment")
	}

	// Update payment status with processed_at timestamp if completing
	if err := s.repo.UpdateStatus(ctx, orgID, paymentID, status); err != nil {
		s.logger.Error("failed to update payment status", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// validate validates a payment
func (s *Service) validate(payment *Payment) error {
	if payment.Amount <= 0 {
		return apperrors.ValidationFailed("amount must be greater than zero")
	}

	if !s.isValidMethod(payment.PaymentMethod) {
		return apperrors.ValidationFailed(fmt.Sprintf("invalid payment method: %s", payment.PaymentMethod))
	}

	if payment.SaleID == uuid.Nil {
		return apperrors.ValidationFailed("sale_id is required")
	}

	if payment.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}

	return nil
}

// isValidMethod validates if the payment method is valid
func (s *Service) isValidMethod(method PaymentMethod) bool {
	switch method {
	case PaymentMethodCash, PaymentMethodCard, PaymentMethodMobile, PaymentMethodBank, PaymentMethodCheck, PaymentMethodCredit:
		return true
	default:
		return false
	}
}

// isValidStatus validates if the payment status is valid
func (s *Service) isValidStatus(status PaymentStatus) bool {
	switch status {
	case PaymentStatusPending, PaymentStatusCompleted, PaymentStatusFailed, PaymentStatusRefunded, PaymentStatusCanceled:
		return true
	default:
		return false
	}
}
