package loyalty

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// LoyaltyService handles loyalty program business logic
type LoyaltyService struct {
	repo   LoyaltyTierRepository
	logger *logging.Logger
}

// NewLoyaltyService creates a new loyalty service
func NewLoyaltyService(repo LoyaltyTierRepository, logger *logging.Logger) *LoyaltyService {
	return &LoyaltyService{
		repo:   repo,
		logger: logger,
	}
}

// ========================================================================
// TIER MANAGEMENT
// ========================================================================

// ListTiers retrieves a list of loyalty tiers with filters
func (s *LoyaltyService) ListTiers(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierFilters) ([]LoyaltyTier, error) {
	tiers, err := s.repo.ListTiers(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty tiers", zap.Error(err), zap.String("org_id", orgID.String()))
		return nil, apperrors.DatabaseError(err)
	}

	return tiers, nil
}

// CreateTier creates a new loyalty tier
func (s *LoyaltyService) CreateTier(ctx context.Context, tier *LoyaltyTier) error {
	// Validate tier
	if err := s.validateTier(tier); err != nil {
		return err
	}

	// Check for duplicate tier code
	existing, err := s.repo.GetTierByCode(ctx, tier.OrganizationID, tier.TierCode)
	if err != nil && !apperrors.IsNotFound(err) {
		s.logger.Error("failed to check duplicate tier code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing != nil {
		return apperrors.AlreadyExists("tier", "tier code already exists")
	}

	// Check for duplicate tier level
	if tier.TierLevel > 0 {
		// This will be checked by the database constraint
	}

	// Set defaults
	tier.ID = uuid.New()
	tier.CreatedAt = time.Now()
	tier.UpdatedAt = time.Now()
	if tier.Metadata == nil {
		tier.Metadata = []byte("{}")
	}
	if tier.PointsMultiplier == 0 {
		tier.PointsMultiplier = 1.0
	}

	// Create tier
	if err := s.repo.CreateTier(ctx, tier); err != nil {
		s.logger.Error("failed to create loyalty tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetTier retrieves a loyalty tier by ID
func (s *LoyaltyService) GetTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTier, error) {
	tier, err := s.repo.GetTier(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty tier", zap.Error(err), zap.String("id", id.String()))
		return nil, apperrors.DatabaseError(err)
	}

	if tier == nil {
		return nil, apperrors.NotFound("loyalty tier")
	}

	return tier, nil
}

// UpdateTier updates an existing loyalty tier
func (s *LoyaltyService) UpdateTier(ctx context.Context, tier *LoyaltyTier) error {
	// Validate tier
	if err := s.validateTier(tier); err != nil {
		return err
	}

	// Check if tier exists
	existing, err := s.repo.GetTier(ctx, tier.OrganizationID, tier.ID)
	if err != nil {
		s.logger.Error("failed to get loyalty tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("loyalty tier")
	}

	// Check for duplicate tier code if changed
	if tier.TierCode != existing.TierCode {
		duplicate, err := s.repo.GetTierByCode(ctx, tier.OrganizationID, tier.TierCode)
		if err != nil && !apperrors.IsNotFound(err) {
			s.logger.Error("failed to check duplicate tier code", zap.Error(err))
			return apperrors.DatabaseError(err)
		}
		if duplicate != nil && duplicate.ID != tier.ID {
			return apperrors.AlreadyExists("tier", "tier code already exists")
		}
	}

	// Update timestamp
	tier.UpdatedAt = time.Now()

	// Update tier
	if err := s.repo.UpdateTier(ctx, tier); err != nil {
		s.logger.Error("failed to update loyalty tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteTier soft-deletes a loyalty tier
func (s *LoyaltyService) DeleteTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if tier exists
	tier, err := s.repo.GetTier(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if tier == nil {
		return apperrors.NotFound("loyalty tier")
	}

	// Prevent deletion of default tier
	if tier.IsDefault {
		return apperrors.ValidationFailed("cannot delete the default loyalty tier")
	}

	// Delete tier
	if err := s.repo.DeleteTier(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete loyalty tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ========================================================================
// BENEFIT MANAGEMENT
// ========================================================================

// ListBenefits retrieves benefits with filters
func (s *LoyaltyService) ListBenefits(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierBenefitFilters) ([]LoyaltyTierBenefit, error) {
	benefits, err := s.repo.ListBenefits(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty tier benefits", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return benefits, nil
}

// CreateBenefit creates a new tier benefit
func (s *LoyaltyService) CreateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error {
	// Validate benefit
	if err := s.validateBenefit(benefit); err != nil {
		return err
	}

	// Verify tier exists
	tier, err := s.repo.GetTier(ctx, benefit.OrganizationID, benefit.TierID)
	if err != nil {
		s.logger.Error("failed to get tier for benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if tier == nil {
		return apperrors.NotFound("loyalty tier")
	}

	// Set defaults
	benefit.ID = uuid.New()
	benefit.CreatedAt = time.Now()
	benefit.UpdatedAt = time.Now()
	if benefit.Metadata == nil {
		benefit.Metadata = []byte("{}")
	}

	// Create benefit
	if err := s.repo.CreateBenefit(ctx, benefit); err != nil {
		s.logger.Error("failed to create loyalty tier benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetBenefit retrieves a benefit by ID
func (s *LoyaltyService) GetBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTierBenefit, error) {
	benefit, err := s.repo.GetBenefit(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty tier benefit", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if benefit == nil {
		return nil, apperrors.NotFound("loyalty tier benefit")
	}

	return benefit, nil
}

// UpdateBenefit updates an existing benefit
func (s *LoyaltyService) UpdateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error {
	// Validate benefit
	if err := s.validateBenefit(benefit); err != nil {
		return err
	}

	// Check if benefit exists
	existing, err := s.repo.GetBenefit(ctx, benefit.OrganizationID, benefit.ID)
	if err != nil {
		s.logger.Error("failed to get loyalty tier benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("loyalty tier benefit")
	}

	// Update timestamp
	benefit.UpdatedAt = time.Now()

	// Update benefit
	if err := s.repo.UpdateBenefit(ctx, benefit); err != nil {
		s.logger.Error("failed to update loyalty tier benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteBenefit deletes a benefit
func (s *LoyaltyService) DeleteBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if benefit exists
	benefit, err := s.repo.GetBenefit(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty tier benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if benefit == nil {
		return apperrors.NotFound("loyalty tier benefit")
	}

	// Delete benefit
	if err := s.repo.DeleteBenefit(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete loyalty tier benefit", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetTierBenefits retrieves all benefits for a tier
func (s *LoyaltyService) GetTierBenefits(ctx context.Context, orgID uuid.UUID, tierID uuid.UUID) ([]LoyaltyTierBenefit, error) {
	benefits, err := s.repo.GetTierBenefits(ctx, orgID, tierID)
	if err != nil {
		s.logger.Error("failed to get tier benefits", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return benefits, nil
}

// ========================================================================
// POINTS RULE MANAGEMENT
// ========================================================================

// ListPointsRules retrieves points rules with filters
func (s *LoyaltyService) ListPointsRules(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsRuleFilters) ([]LoyaltyPointsRule, error) {
	rules, err := s.repo.ListPointsRules(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty points rules", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return rules, nil
}

// CreatePointsRule creates a new points rule
func (s *LoyaltyService) CreatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error {
	// Validate rule
	if err := s.validatePointsRule(rule); err != nil {
		return err
	}

	// Check for duplicate rule code
	existing, err := s.repo.ListPointsRules(ctx, rule.OrganizationID, LoyaltyPointsRuleFilters{})
	if err != nil {
		s.logger.Error("failed to check duplicate rule code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	for _, r := range existing {
		if r.RuleCode == rule.RuleCode && r.ID != rule.ID {
			return apperrors.AlreadyExists("rule", "rule code already exists")
		}
	}

	// Set defaults
	rule.ID = uuid.New()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()
	if rule.Metadata == nil {
		rule.Metadata = []byte("{}")
	}
	if rule.Multiplier == 0 {
		rule.Multiplier = 1.0
	}
	if rule.ApplicableProductIDs == nil {
		rule.ApplicableProductIDs = StringArray{}
	}
	if rule.ApplicableCategoryIDs == nil {
		rule.ApplicableCategoryIDs = StringArray{}
	}
	if rule.ApplicableTierIDs == nil {
		rule.ApplicableTierIDs = StringArray{}
	}

	// Create rule
	if err := s.repo.CreatePointsRule(ctx, rule); err != nil {
		s.logger.Error("failed to create loyalty points rule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetPointsRule retrieves a points rule by ID
func (s *LoyaltyService) GetPointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyPointsRule, error) {
	rule, err := s.repo.GetPointsRule(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty points rule", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if rule == nil {
		return nil, apperrors.NotFound("loyalty points rule")
	}

	return rule, nil
}

// UpdatePointsRule updates an existing points rule
func (s *LoyaltyService) UpdatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error {
	// Validate rule
	if err := s.validatePointsRule(rule); err != nil {
		return err
	}

	// Check if rule exists
	existing, err := s.repo.GetPointsRule(ctx, rule.OrganizationID, rule.ID)
	if err != nil {
		s.logger.Error("failed to get loyalty points rule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("loyalty points rule")
	}

	// Update timestamp
	rule.UpdatedAt = time.Now()

	// Update rule
	if err := s.repo.UpdatePointsRule(ctx, rule); err != nil {
		s.logger.Error("failed to update loyalty points rule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeletePointsRule deletes a points rule
func (s *LoyaltyService) DeletePointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if rule exists
	rule, err := s.repo.GetPointsRule(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty points rule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if rule == nil {
		return apperrors.NotFound("loyalty points rule")
	}

	// Delete rule
	if err := s.repo.DeletePointsRule(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete loyalty points rule", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ========================================================================
// REWARD MANAGEMENT
// ========================================================================

// ListRewards retrieves rewards with filters
func (s *LoyaltyService) ListRewards(ctx context.Context, orgID uuid.UUID, filters LoyaltyRewardFilters) ([]LoyaltyReward, error) {
	rewards, err := s.repo.ListRewards(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty rewards", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return rewards, nil
}

// CreateReward creates a new reward
func (s *LoyaltyService) CreateReward(ctx context.Context, reward *LoyaltyReward) error {
	// Validate reward
	if err := s.validateReward(reward); err != nil {
		return err
	}

	// Check for duplicate reward code
	existing, err := s.repo.ListRewards(ctx, reward.OrganizationID, LoyaltyRewardFilters{})
	if err != nil {
		s.logger.Error("failed to check duplicate reward code", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	for _, r := range existing {
		if r.RewardCode == reward.RewardCode && r.ID != reward.ID {
			return apperrors.AlreadyExists("reward", "reward code already exists")
		}
	}

	// Set defaults
	reward.ID = uuid.New()
	reward.CreatedAt = time.Now()
	reward.UpdatedAt = time.Now()
	reward.TotalRedeemed = 0
	if reward.Metadata == nil {
		reward.Metadata = []byte("{}")
	}
	if reward.TierIDs == nil {
		reward.TierIDs = StringArray{}
	}

	// Create reward
	if err := s.repo.CreateReward(ctx, reward); err != nil {
		s.logger.Error("failed to create loyalty reward", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetReward retrieves a reward by ID
func (s *LoyaltyService) GetReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyReward, error) {
	reward, err := s.repo.GetReward(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty reward", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if reward == nil {
		return nil, apperrors.NotFound("loyalty reward")
	}

	return reward, nil
}

// UpdateReward updates an existing reward
func (s *LoyaltyService) UpdateReward(ctx context.Context, reward *LoyaltyReward) error {
	// Validate reward
	if err := s.validateReward(reward); err != nil {
		return err
	}

	// Check if reward exists
	existing, err := s.repo.GetReward(ctx, reward.OrganizationID, reward.ID)
	if err != nil {
		s.logger.Error("failed to get loyalty reward", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("loyalty reward")
	}

	// Update timestamp
	reward.UpdatedAt = time.Now()

	// Update reward
	if err := s.repo.UpdateReward(ctx, reward); err != nil {
		s.logger.Error("failed to update loyalty reward", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteReward deletes a reward
func (s *LoyaltyService) DeleteReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if reward exists
	reward, err := s.repo.GetReward(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty reward", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if reward == nil {
		return apperrors.NotFound("loyalty reward")
	}

	// Delete reward
	if err := s.repo.DeleteReward(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete loyalty reward", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetFeaturedRewards retrieves featured rewards
func (s *LoyaltyService) GetFeaturedRewards(ctx context.Context, orgID uuid.UUID, limit int) ([]LoyaltyReward, error) {
	rewards, err := s.repo.GetFeaturedRewards(ctx, orgID, limit)
	if err != nil {
		s.logger.Error("failed to get featured rewards", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return rewards, nil
}

// ========================================================================
// REDEMPTION MANAGEMENT
// ========================================================================

// ListRedemptions retrieves redemptions with filters
func (s *LoyaltyService) ListRedemptions(ctx context.Context, orgID uuid.UUID, filters LoyaltyRedemptionFilters) ([]LoyaltyRedemption, error) {
	redemptions, err := s.repo.ListRedemptions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty redemptions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return redemptions, nil
}

// CreateRedemption creates a new redemption
func (s *LoyaltyService) CreateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error {
	// Validate redemption
	if err := s.validateRedemption(redemption); err != nil {
		return err
	}

	// Verify reward exists
	reward, err := s.repo.GetReward(ctx, redemption.OrganizationID, redemption.RewardID)
	if err != nil {
		s.logger.Error("failed to get reward for redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if reward == nil {
		return apperrors.NotFound("loyalty reward")
	}

	// Check reward availability
	if !reward.IsActive {
		return apperrors.ValidationFailed("reward is not active")
	}

	now := time.Now()
	if reward.AvailableFrom != nil && now.Before(*reward.AvailableFrom) {
		return apperrors.ValidationFailed("reward is not yet available")
	}

	if reward.AvailableTo != nil && now.After(*reward.AvailableTo) {
		return apperrors.ValidationFailed("reward is no longer available")
	}

	// Check reward availability count
	if reward.TotalAvailable != nil && reward.TotalRedeemed >= *reward.TotalAvailable {
		return apperrors.ValidationFailed("reward is no longer available (sold out)")
	}

	// Check customer points balance
	balance, err := s.repo.GetCustomerPointsBalance(ctx, redemption.OrganizationID, redemption.CustomerID)
	if err != nil {
		s.logger.Error("failed to get customer points balance", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	if balance < reward.PointsCost {
		return apperrors.ValidationFailed(fmt.Sprintf("insufficient points (need %d, have %d)", reward.PointsCost, balance))
	}

	// Set defaults
	redemption.ID = uuid.New()
	redemption.RedemptionNumber = generateRedemptionNumber()
	redemption.RedemptionDate = time.Now()
	redemption.PointsRedeemed = reward.PointsCost
	redemption.Status = RedemptionStatusPending
	redemption.FulfillmentStatus = FulfillmentStatusPending
	if redemption.Metadata == nil {
		redemption.Metadata = []byte("{}")
	}

	// Create redemption
	if err := s.repo.CreateRedemption(ctx, redemption); err != nil {
		s.logger.Error("failed to create loyalty redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// GetRedemption retrieves a redemption by ID
func (s *LoyaltyService) GetRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyRedemption, error) {
	redemption, err := s.repo.GetRedemption(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty redemption", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if redemption == nil {
		return nil, apperrors.NotFound("loyalty redemption")
	}

	return redemption, nil
}

// UpdateRedemption updates an existing redemption
func (s *LoyaltyService) UpdateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error {
	// Validate redemption
	if err := s.validateRedemption(redemption); err != nil {
		return err
	}

	// Check if redemption exists
	existing, err := s.repo.GetRedemption(ctx, redemption.OrganizationID, redemption.ID)
	if err != nil {
		s.logger.Error("failed to get loyalty redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if existing == nil {
		return apperrors.NotFound("loyalty redemption")
	}

	// Update timestamp
	redemption.UpdatedAt = time.Now()

	// Update redemption
	if err := s.repo.UpdateRedemption(ctx, redemption); err != nil {
		s.logger.Error("failed to update loyalty redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// DeleteRedemption deletes a redemption
func (s *LoyaltyService) DeleteRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	// Check if redemption exists
	redemption, err := s.repo.GetRedemption(ctx, orgID, id)
	if err != nil {
		s.logger.Error("failed to get loyalty redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if redemption == nil {
		return apperrors.NotFound("loyalty redemption")
	}

	// Delete redemption
	if err := s.repo.DeleteRedemption(ctx, orgID, id); err != nil {
		s.logger.Error("failed to delete loyalty redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	return nil
}

// ApproveRedemption approves a redemption
func (s *LoyaltyService) ApproveRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error {
	redemption, err := s.repo.GetRedemption(ctx, orgID, redemptionID)
	if err != nil {
		s.logger.Error("failed to get redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if redemption == nil {
		return apperrors.NotFound("loyalty redemption")
	}

	if redemption.Status != RedemptionStatusPending {
		return apperrors.ValidationFailed("only pending redemptions can be approved")
	}

	redemption.Status = RedemptionStatusApproved
	redemption.UpdatedAt = time.Now()
	redemption.UpdatedBy = &userID

	return s.repo.UpdateRedemption(ctx, redemption)
}

// FulfillRedemption marks a redemption as fulfilled
func (s *LoyaltyService) FulfillRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error {
	redemption, err := s.repo.GetRedemption(ctx, orgID, redemptionID)
	if err != nil {
		s.logger.Error("failed to get redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if redemption == nil {
		return apperrors.NotFound("loyalty redemption")
	}

	if redemption.Status != RedemptionStatusApproved {
		return apperrors.ValidationFailed("only approved redemptions can be fulfilled")
	}

	now := time.Now()
	redemption.Status = RedemptionStatusFulfilled
	redemption.FulfillmentStatus = FulfillmentStatusCompleted
	redemption.FulfilledBy = &userID
	redemption.FulfilledAt = &now
	redemption.UpdatedAt = now
	redemption.UpdatedBy = &userID

	return s.repo.UpdateRedemption(ctx, redemption)
}

// CancelRedemption cancels a redemption
func (s *LoyaltyService) CancelRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error {
	redemption, err := s.repo.GetRedemption(ctx, orgID, redemptionID)
	if err != nil {
		s.logger.Error("failed to get redemption", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if redemption == nil {
		return apperrors.NotFound("loyalty redemption")
	}

	if redemption.Status == RedemptionStatusCancelled {
		return apperrors.ValidationFailed("redemption is already cancelled")
	}

	if redemption.Status == RedemptionStatusFulfilled {
		return apperrors.ValidationFailed("cannot cancel a fulfilled redemption")
	}

	redemption.Status = RedemptionStatusCancelled
	redemption.FulfillmentStatus = FulfillmentStatusCancelled
	redemption.UpdatedAt = time.Now()
	redemption.UpdatedBy = &userID

	return s.repo.UpdateRedemption(ctx, redemption)
}

// ========================================================================
// POINTS TRANSACTION MANAGEMENT
// ========================================================================

// GetCustomerPointsBalance retrieves a customer's current points balance
func (s *LoyaltyService) GetCustomerPointsBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (int, error) {
	balance, err := s.repo.GetCustomerPointsBalance(ctx, orgID, customerID)
	if err != nil {
		s.logger.Error("failed to get customer points balance", zap.Error(err))
		return 0, apperrors.DatabaseError(err)
	}

	return balance, nil
}

// ListPointsTransactions retrieves points transactions with filters
func (s *LoyaltyService) ListPointsTransactions(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsTransactionFilters) ([]LoyaltyPointsTransaction, error) {
	transactions, err := s.repo.ListPointsTransactions(ctx, orgID, filters)
	if err != nil {
		s.logger.Error("failed to list loyalty points transactions", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	return transactions, nil
}

// ========================================================================
// CUSTOMER TIER MANAGEMENT
// ========================================================================

// GetCustomerCurrentTier retrieves a customer's current loyalty tier
func (s *LoyaltyService) GetCustomerCurrentTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*LoyaltyTier, error) {
	tier, err := s.repo.GetCustomerCurrentTier(ctx, orgID, customerID)
	if err != nil {
		s.logger.Error("failed to get customer current tier", zap.Error(err))
		return nil, apperrors.DatabaseError(err)
	}

	if tier == nil {
		// Return default tier if customer has no tier
		defaultTier, err := s.repo.GetDefaultTier(ctx, orgID)
		if err != nil {
			s.logger.Error("failed to get default tier", zap.Error(err))
			return nil, apperrors.DatabaseError(err)
		}
		return defaultTier, nil
	}

	return tier, nil
}

// PromoteCustomerTier promotes a customer to a new tier
func (s *LoyaltyService) PromoteCustomerTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, tierID uuid.UUID, reason string, userID uuid.UUID) error {
	// Get new tier
	newTier, err := s.repo.GetTier(ctx, orgID, tierID)
	if err != nil {
		s.logger.Error("failed to get tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if newTier == nil {
		return apperrors.NotFound("loyalty tier")
	}

	// Get current tier
	currentTier, err := s.repo.GetCustomerCurrentTier(ctx, orgID, customerID)
	if err != nil {
		s.logger.Error("failed to get current tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	// Check tier level
	if currentTier != nil && newTier.TierLevel <= currentTier.TierLevel {
		return apperrors.ValidationFailed("cannot promote to a lower or equal tier level")
	}

	// Create tier history entry
	history := &CustomerTierHistory{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CustomerID:     customerID,
		TierID:         tierID,
		ChangeType:     ChangeTypeUpgrade,
		ChangeReason:   reason,
		EffectiveDate:  time.Now(),
		CreatedAt:      time.Now(),
		CreatedBy:      userID,
	}

	if currentTier != nil {
		history.PreviousTierID = &currentTier.ID
	}

	return s.repo.CreateTierHistory(ctx, history)
}

// DemoteCustomerTier demotes a customer to a new tier
func (s *LoyaltyService) DemoteCustomerTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, tierID uuid.UUID, reason string, userID uuid.UUID) error {
	// Get new tier
	newTier, err := s.repo.GetTier(ctx, orgID, tierID)
	if err != nil {
		s.logger.Error("failed to get tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}
	if newTier == nil {
		return apperrors.NotFound("loyalty tier")
	}

	// Get current tier
	currentTier, err := s.repo.GetCustomerCurrentTier(ctx, orgID, customerID)
	if err != nil {
		s.logger.Error("failed to get current tier", zap.Error(err))
		return apperrors.DatabaseError(err)
	}

	// Check tier level
	if currentTier != nil && newTier.TierLevel >= currentTier.TierLevel {
		return apperrors.ValidationFailed("cannot demote to a higher or equal tier level")
	}

	// Create tier history entry
	history := &CustomerTierHistory{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CustomerID:     customerID,
		TierID:         tierID,
		ChangeType:     ChangeTypeDowngrade,
		ChangeReason:   reason,
		EffectiveDate:  time.Now(),
		CreatedAt:      time.Now(),
		CreatedBy:      userID,
	}

	if currentTier != nil {
		history.PreviousTierID = &currentTier.ID
	}

	return s.repo.CreateTierHistory(ctx, history)
}

// ========================================================================
// VALIDATION HELPERS
// ========================================================================

func (s *LoyaltyService) validateTier(tier *LoyaltyTier) error {
	if tier.TierCode == "" {
		return apperrors.ValidationFailed("tier code is required")
	}
	if len(tier.TierCode) > 50 {
		return apperrors.ValidationFailed("tier code must not exceed 50 characters")
	}
	if tier.TierName == "" {
		return apperrors.ValidationFailed("tier name is required")
	}
	if len(tier.TierName) > 100 {
		return apperrors.ValidationFailed("tier name must not exceed 100 characters")
	}
	if tier.TierLevel < 0 {
		return apperrors.ValidationFailed("tier level must be non-negative")
	}
	if tier.PointsThreshold < 0 {
		return apperrors.ValidationFailed("points threshold must be non-negative")
	}
	if tier.PointsMultiplier <= 0 {
		return apperrors.ValidationFailed("points multiplier must be positive")
	}
	if tier.DiscountPercentage < 0 || tier.DiscountPercentage > 100 {
		return apperrors.ValidationFailed("discount percentage must be between 0 and 100")
	}

	return nil
}

func (s *LoyaltyService) validateBenefit(benefit *LoyaltyTierBenefit) error {
	if benefit.BenefitCode == "" {
		return apperrors.ValidationFailed("benefit code is required")
	}
	if benefit.BenefitName == "" {
		return apperrors.ValidationFailed("benefit name is required")
	}
	if benefit.BenefitType == "" {
		return apperrors.ValidationFailed("benefit type is required")
	}

	return nil
}

func (s *LoyaltyService) validatePointsRule(rule *LoyaltyPointsRule) error {
	if rule.RuleCode == "" {
		return apperrors.ValidationFailed("rule code is required")
	}
	if rule.RuleName == "" {
		return apperrors.ValidationFailed("rule name is required")
	}
	if rule.RuleType == "" {
		return apperrors.ValidationFailed("rule type is required")
	}
	if rule.PointsPerAmount == nil && rule.FixedPoints == nil {
		return apperrors.ValidationFailed("either points_per_amount or fixed_points is required")
	}

	return nil
}

func (s *LoyaltyService) validateReward(reward *LoyaltyReward) error {
	if reward.RewardCode == "" {
		return apperrors.ValidationFailed("reward code is required")
	}
	if reward.RewardName == "" {
		return apperrors.ValidationFailed("reward name is required")
	}
	if reward.RewardType == "" {
		return apperrors.ValidationFailed("reward type is required")
	}
	if reward.PointsCost <= 0 {
		return apperrors.ValidationFailed("points cost must be positive")
	}

	return nil
}

func (s *LoyaltyService) validateRedemption(redemption *LoyaltyRedemption) error {
	if redemption.CustomerID == uuid.Nil {
		return apperrors.ValidationFailed("customer id is required")
	}
	if redemption.RewardID == uuid.Nil {
		return apperrors.ValidationFailed("reward id is required")
	}

	return nil
}

// ========================================================================
// HELPER FUNCTIONS
// ========================================================================

func generateRedemptionNumber() string {
	return fmt.Sprintf("RED-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
}
