package giftcards

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// Service handles gift card, store credit, and return business logic
type Service struct {
	repo   GiftCardRepository
	logger *logging.Logger
}

// NewService creates a new service
func NewService(repo GiftCardRepository, logger *logging.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// ============= GIFT CARDS =============

// CreateGiftCard creates a new gift card
func (s *Service) CreateGiftCard(ctx context.Context, card *GiftCard) error {
	if err := s.validateGiftCard(card); err != nil {
		return err
	}

	card.ID = uuid.New()
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()
	card.Status = "active"
	card.CurrentBalance = card.OriginalValue

	if err := s.repo.CreateGiftCard(ctx, card); err != nil {
		s.logger.Error("failed to create gift card", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetGiftCard retrieves a gift card by ID
func (s *Service) GetGiftCard(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*GiftCard, error) {
	card, err := s.repo.GetGiftCard(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get gift card", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return card, nil
}

// LookupGiftCard retrieves a gift card by card number
func (s *Service) LookupGiftCard(ctx context.Context, orgID uuid.UUID, cardNumber string) (*GiftCard, error) {
	card, err := s.repo.GetGiftCardByNumber(ctx, orgID, cardNumber)
	if err != nil {
		s.logger.Error("failed to lookup gift card", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return card, nil
}

// ListGiftCards lists gift cards with filters
func (s *Service) ListGiftCards(ctx context.Context, orgID uuid.UUID, filters GiftCardFilters) ([]GiftCard, error) {
	cards, err := s.repo.ListGiftCards(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list gift cards", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return cards, nil
}

// RedeemGiftCard redeems a gift card and records transaction
func (s *Service) RedeemGiftCard(ctx context.Context, orgID uuid.UUID, cardID uuid.UUID, amount float64, saleID *uuid.UUID, userID *uuid.UUID, locationID *uuid.UUID) error {
	card, err := s.repo.GetGiftCard(ctx, orgID, cardID)
	if err != nil || card == nil {
		return apperrors.NotFound("gift_card", "not found")
	}

	if card.Status != "active" {
		return apperrors.Invalid("gift_card", fmt.Sprintf("card status is %s", card.Status))
	}

	if card.CurrentBalance < amount {
		return apperrors.Invalid("gift_card", "insufficient balance")
	}

	card.CurrentBalance -= amount
	if card.CurrentBalance == 0 {
		card.Status = "fully_redeemed"
	}
	card.UpdatedAt = time.Now()

	if err := s.repo.UpdateGiftCard(ctx, card); err != nil {
		s.logger.Error("failed to update gift card balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	// Record transaction
	txn := &GiftCardTransaction{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		GiftCardID:      cardID,
		TransactionType: "redemption",
		Amount:          amount,
		BalanceAfter:    card.CurrentBalance,
		SaleID:          saleID,
		UserID:          userID,
		LocationID:      locationID,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreateGiftCardTransaction(ctx, txn); err != nil {
		s.logger.Error("failed to create gift card transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// LoadGiftCard adds balance to a gift card
func (s *Service) LoadGiftCard(ctx context.Context, orgID uuid.UUID, cardID uuid.UUID, amount float64, userID *uuid.UUID) error {
	card, err := s.repo.GetGiftCard(ctx, orgID, cardID)
	if err != nil || card == nil {
		return apperrors.NotFound("gift_card", "not found")
	}

	card.CurrentBalance += amount
	card.UpdatedAt = time.Now()

	if err := s.repo.UpdateGiftCard(ctx, card); err != nil {
		s.logger.Error("failed to update gift card balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	txn := &GiftCardTransaction{
		ID:              uuid.New(),
		OrganizationID:  orgID,
		GiftCardID:      cardID,
		TransactionType: "load",
		Amount:          amount,
		BalanceAfter:    card.CurrentBalance,
		UserID:          userID,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.CreateGiftCardTransaction(ctx, txn); err != nil {
		s.logger.Error("failed to create gift card transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ============= STORE CREDIT =============

// GetOrCreateStoreCreditAccount gets or creates a store credit account
func (s *Service) GetOrCreateStoreCreditAccount(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*StoreCreditAccount, error) {
	account, err := s.repo.GetStoreCreditAccount(ctx, orgID, customerID)
	if err != nil {
		s.logger.Error("failed to get store credit account", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if account != nil {
		return account, nil
	}

	account = &StoreCreditAccount{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CustomerID:     customerID,
		CurrentBalance: 0,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.CreateStoreCreditAccount(ctx, account); err != nil {
		s.logger.Error("failed to create store credit account", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return account, nil
}

// RedeemStoreCredit redeems store credit
func (s *Service) RedeemStoreCredit(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, amount float64, saleID *uuid.UUID, userID *uuid.UUID) error {
	account, err := s.repo.GetStoreCreditAccount(ctx, orgID, customerID)
	if err != nil || account == nil {
		return apperrors.NotFound("store_credit_account", "not found")
	}

	if !account.IsActive {
		return apperrors.Invalid("store_credit_account", "account is inactive")
	}

	if account.CurrentBalance < amount {
		return apperrors.Invalid("store_credit_account", "insufficient balance")
	}

	account.CurrentBalance -= amount
	account.UpdatedAt = time.Now()

	if err := s.repo.UpdateStoreCreditAccount(ctx, account); err != nil {
		s.logger.Error("failed to update store credit balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	txn := &StoreCreditTransaction{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		StoreCreditAccountID: account.ID,
		TransactionType:      "redemption",
		Amount:               amount,
		BalanceAfter:         account.CurrentBalance,
		SaleID:               saleID,
		UserID:               userID,
		CreatedAt:            time.Now(),
	}

	if err := s.repo.CreateStoreCreditTransaction(ctx, txn); err != nil {
		s.logger.Error("failed to create store credit transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// IssueStoreCredit issues store credit to a customer
func (s *Service) IssueStoreCredit(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, amount float64, userID *uuid.UUID) error {
	account, err := s.GetOrCreateStoreCreditAccount(ctx, orgID, customerID)
	if err != nil {
		return err
	}

	account.CurrentBalance += amount
	account.UpdatedAt = time.Now()

	if err := s.repo.UpdateStoreCreditAccount(ctx, account); err != nil {
		s.logger.Error("failed to update store credit balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	txn := &StoreCreditTransaction{
		ID:                   uuid.New(),
		OrganizationID:       orgID,
		StoreCreditAccountID: account.ID,
		TransactionType:      "issue",
		Amount:               amount,
		BalanceAfter:         account.CurrentBalance,
		UserID:               userID,
		CreatedAt:            time.Now(),
	}

	if err := s.repo.CreateStoreCreditTransaction(ctx, txn); err != nil {
		s.logger.Error("failed to create store credit transaction", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ============= SALE RETURNS =============

// CreateSaleReturn creates a return/RMA
func (s *Service) CreateSaleReturn(ctx context.Context, return_ *SaleReturn) error {
	if err := s.validateSaleReturn(return_); err != nil {
		return err
	}

	return_.ID = uuid.New()
	return_.CreatedAt = time.Now()
	return_.UpdatedAt = time.Now()
	return_.Status = "pending"

	if err := s.repo.CreateSaleReturn(ctx, return_); err != nil {
		s.logger.Error("failed to create sale return", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetSaleReturn retrieves a return
func (s *Service) GetSaleReturn(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*SaleReturn, error) {
	return_, err := s.repo.GetSaleReturn(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get sale return", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return return_, nil
}

// ListSaleReturns lists returns with filters
func (s *Service) ListSaleReturns(ctx context.Context, orgID uuid.UUID, filters ReturnFilters) ([]SaleReturn, error) {
	returns, err := s.repo.ListSaleReturns(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list sale returns", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}
	return returns, nil
}

// ApproveSaleReturn approves a return
func (s *Service) ApproveSaleReturn(ctx context.Context, orgID uuid.UUID, returnID uuid.UUID, approvedBy uuid.UUID) error {
	return_, err := s.repo.GetSaleReturn(ctx, orgID, returnID)
	if err != nil || return_ == nil {
		return apperrors.NotFound("sale_return", "not found")
	}

	return_.Status = "approved"
	return_.ApprovedBy = &approvedBy
	now := time.Now()
	return_.ApprovedAt = &now
	return_.UpdatedAt = now

	if err := s.repo.UpdateSaleReturn(ctx, return_); err != nil {
		s.logger.Error("failed to approve sale return", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// CompleteSaleReturn marks a return as completed
func (s *Service) CompleteSaleReturn(ctx context.Context, orgID uuid.UUID, returnID uuid.UUID) error {
	return_, err := s.repo.GetSaleReturn(ctx, orgID, returnID)
	if err != nil || return_ == nil {
		return apperrors.NotFound("sale_return", "not found")
	}

	return_.Status = "completed"
	return_.UpdatedAt = time.Now()

	if err := s.repo.UpdateSaleReturn(ctx, return_); err != nil {
		s.logger.Error("failed to complete sale return", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// AddReturnItem adds an item to a return
func (s *Service) AddReturnItem(ctx context.Context, item *SaleReturnItem) error {
	item.ID = uuid.New()
	item.CreatedAt = time.Now()

	if err := s.repo.CreateSaleReturnItem(ctx, item); err != nil {
		s.logger.Error("failed to create sale return item", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ============= VALIDATION =============

func (s *Service) validateGiftCard(card *GiftCard) error {
	if card.OrganizationID == uuid.Nil {
		return apperrors.Invalid("gift_card", "organization_id required")
	}
	if card.CardNumber == "" {
		return apperrors.Invalid("gift_card", "card_number required")
	}
	if card.OriginalValue <= 0 {
		return apperrors.Invalid("gift_card", "original_value must be positive")
	}
	return nil
}

func (s *Service) validateSaleReturn(return_ *SaleReturn) error {
	if return_.OrganizationID == uuid.Nil {
		return apperrors.Invalid("sale_return", "organization_id required")
	}
	if return_.LocationID == uuid.Nil {
		return apperrors.Invalid("sale_return", "location_id required")
	}
	if return_.UserID == uuid.Nil {
		return apperrors.Invalid("sale_return", "user_id required")
	}
	return nil
}
