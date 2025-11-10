package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles inventory transaction business logic
type Service struct {
	repo   Repository
	logger *logging.Logger
}

// NewService creates a new inventory transaction service
func NewService(repo Repository, logger *logging.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves a list of inventory transactions with filters
func (s *Service) List(ctx context.Context, orgID uuid.UUID, filters InventoryTransactionFilters) ([]InventoryTransaction, error) {
	transactions, err := s.repo.List(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list inventory transactions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return transactions, nil
}

// Count counts inventory transactions matching filters
func (s *Service) Count(ctx context.Context, orgID uuid.UUID, filters InventoryTransactionFilters) (int64, error) {
	count, err := s.repo.Count(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to count inventory transactions", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return count, nil
}

// Create creates a new inventory transaction
func (s *Service) Create(ctx context.Context, transaction *InventoryTransaction) error {
	// Validate transaction
	if err := s.validate(transaction); err != nil {
		return err
	}

	// Set defaults
	transaction.ID = uuid.New()
	transaction.CreatedAt = time.Now()
	if transaction.TransactionDate.IsZero() {
		transaction.TransactionDate = time.Now()
	}

	// Calculate balance after transaction
	// Get the latest balance for the product
	latestTx, err := s.repo.GetLatestBalanceForProduct(ctx, transaction.OrganizationID, transaction.ProductID)
	if err != nil {
		s.logger.Error("failed to get latest balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	var currentBalance float64 = 0
	if latestTx != nil {
		currentBalance = latestTx.BalanceAfter
	}

	transaction.BalanceAfter = currentBalance + transaction.Quantity

	// Create transaction
	if err := s.repo.Create(ctx, transaction); err != nil {
		s.logger.Error("failed to create inventory transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Get retrieves an inventory transaction by ID
func (s *Service) Get(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*InventoryTransaction, error) {
	transaction, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get inventory transaction", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if transaction == nil {
		return nil, apperrors.NotFound("inventory transaction")
	}

	return transaction, nil
}

// GetByProduct retrieves all transactions for a product
func (s *Service) GetByProduct(ctx context.Context, orgID uuid.UUID, productID uuid.UUID, limit int, offset int) ([]InventoryTransaction, error) {
	transactions, err := s.repo.GetByProduct(ctx, orgID, productID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get transactions by product", zap.Error(err), zap.String("product_id", productID.String()))
		return nil, apperrors.DatabaseError(err)
	}

	return transactions, nil
}

// GetBySale retrieves all transactions for a sale
func (s *Service) GetBySale(ctx context.Context, orgID uuid.UUID, saleID uuid.UUID) ([]InventoryTransaction, error) {
	transactions, err := s.repo.GetBySale(ctx, orgID, saleID)
	if err != nil {
		s.logger.Error("failed to get transactions by sale", zap.Error(err), zap.String("sale_id", saleID.String()))
		return nil, apperrors.DatabaseError(err)
	}

	return transactions, nil
}

// Update updates an existing inventory transaction
func (s *Service) Update(ctx context.Context, transaction *InventoryTransaction) error {
	// Validate transaction
	if err := s.validate(transaction); err != nil {
		return err
	}

	// Check if transaction exists
	existing, err := s.repo.Get(ctx, transaction.OrganizationID, transaction.ID)
	if err != nil {
		s.logger.Error("failed to get inventory transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("inventory transaction")
	}

	// Update transaction
	if err := s.repo.Update(ctx, transaction); err != nil {
		s.logger.Error("failed to update inventory transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// Delete deletes an inventory transaction
func (s *Service) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if transaction exists
	transaction, err := s.repo.Get(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get inventory transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if transaction == nil {
		return apperrors.NotFound("inventory transaction")
	}

	// Delete transaction
	if err := s.repo.Delete(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete inventory transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetCurrentStock retrieves the current stock for a product
func (s *Service) GetCurrentStock(ctx context.Context, orgID uuid.UUID, productID uuid.UUID) (float64, error) {
	latestTx, err := s.repo.GetLatestBalanceForProduct(ctx, orgID, productID)
	if err != nil {
		s.logger.Error("failed to get latest balance", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	if latestTx == nil {
		return 0, nil
	}

	return latestTx.BalanceAfter, nil
}

// validate validates an inventory transaction
func (s *Service) validate(transaction *InventoryTransaction) error {
	if transaction.ProductID == uuid.Nil {
		return apperrors.ValidationFailed("product_id is required")
	}

	if !s.isValidType(transaction.TransactionType) {
		return apperrors.ValidationFailed(fmt.Sprintf("invalid transaction type: %s", transaction.TransactionType))
	}

	if transaction.Quantity == 0 {
		return apperrors.ValidationFailed("quantity cannot be zero")
	}

	if transaction.OrganizationID == uuid.Nil {
		return apperrors.ValidationFailed("organization_id is required")
	}

	// Validate quantity for certain transaction types
	if transaction.TransactionType == TransactionTypeSale && transaction.Quantity > 0 {
		return apperrors.ValidationFailed("sale transactions must have negative quantity")
	}

	return nil
}

// isValidType validates if the transaction type is valid
func (s *Service) isValidType(txType TransactionType) bool {
	switch txType {
	case TransactionTypeSale, TransactionTypePurchase, TransactionTypeAdjustment, TransactionTypeTransfer, TransactionTypeReturn, TransactionTypeWriteOff:
		return true
	default:
		return false
	}
}
