package loyalty

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ========================================================================
// LOYALTY TIERS
// ========================================================================

// LoyaltyTier represents a customer loyalty tier/level
type LoyaltyTier struct {
	ID                       uuid.UUID      `json:"id"`
	OrganizationID           uuid.UUID      `json:"organization_id"`
	TierCode                 string         `json:"tier_code"`
	TierName                 string         `json:"tier_name"`
	TierLevel                int            `json:"tier_level"`
	Description              string         `json:"description"`
	PointsThreshold          int            `json:"points_threshold"`
	AnnualSpendThreshold     *float64       `json:"annual_spend_threshold"`
	PurchaseCountThreshold   *int           `json:"purchase_count_threshold"`
	PointsMultiplier         float64        `json:"points_multiplier"`
	DiscountPercentage       float64        `json:"discount_percentage"`
	TierColor                string         `json:"tier_color"`
	TierIcon                 string         `json:"tier_icon"`
	BadgeImageURL            string         `json:"badge_image_url"`
	IsActive                 bool           `json:"is_active"`
	IsDefault                bool           `json:"is_default"`
	SortOrder                int            `json:"sort_order"`
	Metadata                 []byte         `json:"metadata"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                *time.Time     `json:"deleted_at,omitempty"`
	CreatedBy                uuid.UUID      `json:"created_by"`
	UpdatedBy                *uuid.UUID     `json:"updated_by"`
}

// LoyaltyTierFilters represents filters for listing loyalty tiers
type LoyaltyTierFilters struct {
	Search   string
	IsActive *bool
	Page     int
	PageSize int
}

// ========================================================================
// LOYALTY TIER BENEFITS
// ========================================================================

// BenefitType represents the type of benefit
type BenefitType string

const (
	BenefitTypeDiscount         BenefitType = "discount"
	BenefitTypeFreeShipping     BenefitType = "free_shipping"
	BenefitTypeBirthdayBonus    BenefitType = "birthday_bonus"
	BenefitTypeEarlyAccess      BenefitType = "early_access"
	BenefitTypePrioritySupport  BenefitType = "priority_support"
	BenefitTypeExclusiveProducts BenefitType = "exclusive_products"
	BenefitTypeFreeProduct      BenefitType = "free_product"
	BenefitTypeOther            BenefitType = "other"
)

// DiscountType represents how discount is applied
type DiscountType string

const (
	DiscountTypePercentage  DiscountType = "percentage"
	DiscountTypeFixedAmount DiscountType = "fixed_amount"
)

// LoyaltyTierBenefit represents a specific benefit for a loyalty tier
type LoyaltyTierBenefit struct {
	ID                   uuid.UUID  `json:"id"`
	OrganizationID       uuid.UUID  `json:"organization_id"`
	TierID               uuid.UUID  `json:"tier_id"`
	BenefitCode          string     `json:"benefit_code"`
	BenefitName          string     `json:"benefit_name"`
	BenefitDescription   string     `json:"benefit_description"`
	BenefitType          BenefitType `json:"benefit_type"`
	DiscountValue        *float64   `json:"discount_value"`
	DiscountType         *DiscountType `json:"discount_type"`
	IsActive             bool       `json:"is_active"`
	SortOrder            int        `json:"sort_order"`
	Icon                 string     `json:"icon"`
	TermsAndConditions   string     `json:"terms_and_conditions"`
	Metadata             []byte     `json:"metadata"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// LoyaltyTierBenefitFilters represents filters for listing benefits
type LoyaltyTierBenefitFilters struct {
	TierID   *uuid.UUID
	IsActive *bool
	Page     int
	PageSize int
}

// ========================================================================
// LOYALTY POINTS RULES
// ========================================================================

// RuleType represents the type of points earning rule
type RuleType string

const (
	RuleTypePurchase     RuleType = "purchase"
	RuleTypeSignup       RuleType = "signup"
	RuleTypeBirthday     RuleType = "birthday"
	RuleTypeReferral     RuleType = "referral"
	RuleTypeReview       RuleType = "review"
	RuleTypeSocialShare  RuleType = "social_share"
	RuleTypeManual       RuleType = "manual"
	RuleTypeOther        RuleType = "other"
)

// AppliesToType represents what the rule applies to
type AppliesToType string

const (
	AppliesToAll                AppliesToType = "all"
	AppliesToSpecificProducts   AppliesToType = "specific_products"
	AppliesToSpecificCategories AppliesToType = "specific_categories"
	AppliesToSpecificTiers      AppliesToType = "specific_tiers"
)

// StringArray is a custom type for storing UUID arrays in JSONB
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return sql.ErrScanFailure
	}
	return json.Unmarshal(bytes, &a)
}

// LoyaltyPointsRule represents a rule for earning loyalty points
type LoyaltyPointsRule struct {
	ID                          uuid.UUID      `json:"id"`
	OrganizationID              uuid.UUID      `json:"organization_id"`
	RuleCode                    string         `json:"rule_code"`
	RuleName                    string         `json:"rule_name"`
	Description                 string         `json:"description"`
	RuleType                    RuleType       `json:"rule_type"`
	PointsPerAmount             *float64       `json:"points_per_amount"`
	FixedPoints                 *int           `json:"fixed_points"`
	Multiplier                  float64        `json:"multiplier"`
	AppliesToType               AppliesToType  `json:"applies_to"`
	ApplicableProductIDs        StringArray    `json:"applicable_product_ids"`
	ApplicableCategoryIDs       StringArray    `json:"applicable_category_ids"`
	ApplicableTierIDs           StringArray    `json:"applicable_tier_ids"`
	MinimumPurchaseAmount       *float64       `json:"minimum_purchase_amount"`
	MaximumPointsPerTransaction *int           `json:"maximum_points_per_transaction"`
	MaximumPointsPerDay         *int           `json:"maximum_points_per_day"`
	MaximumPointsPerMonth       *int           `json:"maximum_points_per_month"`
	StartDate                   *time.Time     `json:"start_date"`
	EndDate                     *time.Time     `json:"end_date"`
	IsActive                    bool           `json:"is_active"`
	Priority                    int            `json:"priority"`
	TermsAndConditions          string         `json:"terms_and_conditions"`
	Metadata                    []byte         `json:"metadata"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   *time.Time     `json:"deleted_at,omitempty"`
	CreatedBy                   uuid.UUID      `json:"created_by"`
	UpdatedBy                   *uuid.UUID     `json:"updated_by"`
}

// LoyaltyPointsRuleFilters represents filters for listing points rules
type LoyaltyPointsRuleFilters struct {
	Search    string
	RuleType  *RuleType
	IsActive  *bool
	Page      int
	PageSize  int
}

// ========================================================================
// LOYALTY REWARDS
// ========================================================================

// RewardType represents the type of reward
type RewardType string

const (
	RewardTypeDiscountPercentage RewardType = "discount_percentage"
	RewardTypeDiscountFixed      RewardType = "discount_fixed"
	RewardTypeFreeProduct        RewardType = "free_product"
	RewardTypeFreeShipping       RewardType = "free_shipping"
	RewardTypeGiftCard           RewardType = "gift_card"
	RewardTypeExperience         RewardType = "experience"
	RewardTypeOther              RewardType = "other"
)

// LoyaltyReward represents a reward that can be redeemed with points
type LoyaltyReward struct {
	ID                        uuid.UUID  `json:"id"`
	OrganizationID            uuid.UUID  `json:"organization_id"`
	RewardCode                string     `json:"reward_code"`
	RewardName                string     `json:"reward_name"`
	Description               string     `json:"description"`
	RewardType                RewardType `json:"reward_type"`
	PointsCost                int        `json:"points_cost"`
	RewardValue               *float64   `json:"reward_value"`
	DiscountPercentage        *float64   `json:"discount_percentage"`
	DiscountAmount            *float64   `json:"discount_amount"`
	ProductID                 *uuid.UUID `json:"product_id"`
	ProductVariantID          *uuid.UUID `json:"product_variant_id"`
	IsActive                  bool       `json:"is_active"`
	AvailableFrom             *time.Time `json:"available_from"`
	AvailableTo               *time.Time `json:"available_to"`
	TotalAvailable            *int       `json:"total_available"`
	TotalRedeemed             int        `json:"total_redeemed"`
	MaxRedemptionsPerCustomer *int       `json:"max_redemptions_per_customer"`
	MinimumTierLevel          *int       `json:"minimum_tier_level"`
	TierIDs                   StringArray `json:"tier_ids"`
	ImageURL                  string     `json:"image_url"`
	ThumbnailURL              string     `json:"thumbnail_url"`
	SortOrder                 int        `json:"sort_order"`
	IsFeatured                bool       `json:"is_featured"`
	TermsAndConditions        string     `json:"terms_and_conditions"`
	RedemptionInstructions    string     `json:"redemption_instructions"`
	Metadata                  []byte     `json:"metadata"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
	DeletedAt                 *time.Time `json:"deleted_at,omitempty"`
	CreatedBy                 uuid.UUID  `json:"created_by"`
	UpdatedBy                 *uuid.UUID `json:"updated_by"`
}

// LoyaltyRewardFilters represents filters for listing rewards
type LoyaltyRewardFilters struct {
	Search      string
	RewardType  *RewardType
	IsActive    *bool
	IsFeatured  *bool
	MinPoints   *int
	MaxPoints   *int
	Page        int
	PageSize    int
}

// ========================================================================
// LOYALTY REDEMPTIONS
// ========================================================================

// RedemptionStatus represents the status of a redemption
type RedemptionStatus string

const (
	RedemptionStatusPending   RedemptionStatus = "pending"
	RedemptionStatusApproved  RedemptionStatus = "approved"
	RedemptionStatusFulfilled RedemptionStatus = "fulfilled"
	RedemptionStatusCancelled RedemptionStatus = "cancelled"
	RedemptionStatusExpired   RedemptionStatus = "expired"
)

// FulfillmentStatus represents the fulfillment status
type FulfillmentStatus string

const (
	FulfillmentStatusPending    FulfillmentStatus = "pending"
	FulfillmentStatusInProgress FulfillmentStatus = "in_progress"
	FulfillmentStatusCompleted  FulfillmentStatus = "completed"
	FulfillmentStatusCancelled  FulfillmentStatus = "cancelled"
)

// LoyaltyRedemption represents a reward redemption
type LoyaltyRedemption struct {
	ID                 uuid.UUID         `json:"id"`
	OrganizationID     uuid.UUID         `json:"organization_id"`
	CustomerID         uuid.UUID         `json:"customer_id"`
	RewardID           uuid.UUID         `json:"reward_id"`
	RedemptionNumber   string            `json:"redemption_number"`
	RedemptionDate     time.Time         `json:"redemption_date"`
	PointsRedeemed     int               `json:"points_redeemed"`
	Status             RedemptionStatus  `json:"status"`
	SaleID             *uuid.UUID        `json:"sale_id"`
	UsedDate           *time.Time        `json:"used_date"`
	ExpiryDate         *time.Time        `json:"expiry_date"`
	FulfillmentStatus  FulfillmentStatus `json:"fulfillment_status"`
	FulfillmentNotes   string            `json:"fulfillment_notes"`
	FulfilledBy        *uuid.UUID        `json:"fulfilled_by"`
	FulfilledAt        *time.Time        `json:"fulfilled_at"`
	Notes              string            `json:"notes"`
	Metadata           []byte            `json:"metadata"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	CreatedBy          uuid.UUID         `json:"created_by"`
	UpdatedBy          *uuid.UUID        `json:"updated_by"`
}

// LoyaltyRedemptionFilters represents filters for listing redemptions
type LoyaltyRedemptionFilters struct {
	CustomerID        *uuid.UUID
	Status            *RedemptionStatus
	FulfillmentStatus *FulfillmentStatus
	Page              int
	PageSize          int
}

// ========================================================================
// LOYALTY POINTS TRANSACTIONS
// ========================================================================

// TransactionType represents the type of points transaction
type TransactionType string

const (
	TransactionTypeEarned    TransactionType = "earned"
	TransactionTypeRedeemed  TransactionType = "redeemed"
	TransactionTypeExpired   TransactionType = "expired"
	TransactionTypeAdjusted  TransactionType = "adjusted"
	TransactionTypeBonus     TransactionType = "bonus"
	TransactionTypeRefunded  TransactionType = "refunded"
)

// LoyaltyPointsTransaction represents a points transaction
type LoyaltyPointsTransaction struct {
	ID              uuid.UUID       `json:"id"`
	OrganizationID  uuid.UUID       `json:"organization_id"`
	CustomerID      uuid.UUID       `json:"customer_id"`
	TransactionType TransactionType `json:"transaction_type"`
	Points          int             `json:"points"`
	BalanceAfter    int             `json:"balance_after"`
	SaleID          *uuid.UUID      `json:"sale_id"`
	RedemptionID    *uuid.UUID      `json:"redemption_id"`
	PointsRuleID    *uuid.UUID      `json:"points_rule_id"`
	Description     string          `json:"description"`
	Reason          string          `json:"reason"`
	Notes           string          `json:"notes"`
	ExpiryDate      *time.Time      `json:"expiry_date"`
	Metadata        []byte          `json:"metadata"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedBy       uuid.UUID       `json:"created_by"`
}

// LoyaltyPointsTransactionFilters represents filters for listing transactions
type LoyaltyPointsTransactionFilters struct {
	CustomerID      *uuid.UUID
	TransactionType *TransactionType
	StartDate       *time.Time
	EndDate         *time.Time
	Page            int
	PageSize        int
}

// ========================================================================
// CUSTOMER TIER HISTORY
// ========================================================================

// ChangeType represents the type of tier change
type ChangeType string

const (
	ChangeTypeUpgrade   ChangeType = "upgrade"
	ChangeTypeDowngrade ChangeType = "downgrade"
	ChangeTypeInitial   ChangeType = "initial"
	ChangeTypeManual    ChangeType = "manual"
)

// CustomerTierHistory represents a customer's tier change history
type CustomerTierHistory struct {
	ID                 uuid.UUID   `json:"id"`
	OrganizationID     uuid.UUID   `json:"organization_id"`
	CustomerID         uuid.UUID   `json:"customer_id"`
	TierID             uuid.UUID   `json:"tier_id"`
	PreviousTierID     *uuid.UUID  `json:"previous_tier_id"`
	ChangeType         ChangeType  `json:"change_type"`
	ChangeReason       string      `json:"change_reason"`
	QualifyingPoints   *int        `json:"qualifying_points"`
	QualifyingSpend    *float64    `json:"qualifying_spend"`
	QualifyingPurchases *int       `json:"qualifying_purchases"`
	EffectiveDate      time.Time   `json:"effective_date"`
	ValidUntil         *time.Time  `json:"valid_until"`
	Notes              string      `json:"notes"`
	Metadata           []byte      `json:"metadata"`
	CreatedAt          time.Time   `json:"created_at"`
	CreatedBy          uuid.UUID   `json:"created_by"`
}

// CustomerTierHistoryFilters represents filters for listing tier history
type CustomerTierHistoryFilters struct {
	CustomerID *uuid.UUID
	TierID     *uuid.UUID
	ChangeType *ChangeType
	Page       int
	PageSize   int
}

// ========================================================================
// REPOSITORY INTERFACES
// ========================================================================

// LoyaltyTierRepository defines the loyalty tier data access interface
type LoyaltyTierRepository interface {
	// Tier operations
	ListTiers(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierFilters) ([]LoyaltyTier, error)
	CountTiers(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierFilters) (int64, error)
	CreateTier(ctx context.Context, tier *LoyaltyTier) error
	GetTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTier, error)
	GetTierByCode(ctx context.Context, orgID uuid.UUID, code string) (*LoyaltyTier, error)
	UpdateTier(ctx context.Context, tier *LoyaltyTier) error
	DeleteTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetDefaultTier(ctx context.Context, orgID uuid.UUID) (*LoyaltyTier, error)

	// Benefit operations
	ListBenefits(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierBenefitFilters) ([]LoyaltyTierBenefit, error)
	CreateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error
	GetBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTierBenefit, error)
	UpdateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error
	DeleteBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetTierBenefits(ctx context.Context, orgID uuid.UUID, tierID uuid.UUID) ([]LoyaltyTierBenefit, error)

	// Points rule operations
	ListPointsRules(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsRuleFilters) ([]LoyaltyPointsRule, error)
	CreatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error
	GetPointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyPointsRule, error)
	UpdatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error
	DeletePointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetActivePointsRules(ctx context.Context, orgID uuid.UUID) ([]LoyaltyPointsRule, error)

	// Reward operations
	ListRewards(ctx context.Context, orgID uuid.UUID, filters LoyaltyRewardFilters) ([]LoyaltyReward, error)
	CreateReward(ctx context.Context, reward *LoyaltyReward) error
	GetReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyReward, error)
	UpdateReward(ctx context.Context, reward *LoyaltyReward) error
	DeleteReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetFeaturedRewards(ctx context.Context, orgID uuid.UUID, limit int) ([]LoyaltyReward, error)

	// Redemption operations
	ListRedemptions(ctx context.Context, orgID uuid.UUID, filters LoyaltyRedemptionFilters) ([]LoyaltyRedemption, error)
	CreateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error
	GetRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyRedemption, error)
	UpdateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error
	DeleteRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetCustomerRedemptions(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]LoyaltyRedemption, error)

	// Points transaction operations
	ListPointsTransactions(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsTransactionFilters) ([]LoyaltyPointsTransaction, error)
	CreatePointsTransaction(ctx context.Context, transaction *LoyaltyPointsTransaction) error
	GetPointsTransaction(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyPointsTransaction, error)
	GetCustomerPointsBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (int, error)
	GetCustomerPointsTransactions(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]LoyaltyPointsTransaction, error)

	// Tier history operations
	ListTierHistory(ctx context.Context, orgID uuid.UUID, filters CustomerTierHistoryFilters) ([]CustomerTierHistory, error)
	CreateTierHistory(ctx context.Context, history *CustomerTierHistory) error
	GetTierHistory(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*CustomerTierHistory, error)
	GetCustomerCurrentTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*LoyaltyTier, error)
	GetCustomerTierHistory(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) ([]CustomerTierHistory, error)
}

// Service defines the loyalty business logic interface
type Service interface {
	// Tier management
	ListTiers(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierFilters) ([]LoyaltyTier, error)
	CreateTier(ctx context.Context, tier *LoyaltyTier) error
	GetTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTier, error)
	UpdateTier(ctx context.Context, tier *LoyaltyTier) error
	DeleteTier(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Benefit management
	ListBenefits(ctx context.Context, orgID uuid.UUID, filters LoyaltyTierBenefitFilters) ([]LoyaltyTierBenefit, error)
	CreateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error
	GetBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyTierBenefit, error)
	UpdateBenefit(ctx context.Context, benefit *LoyaltyTierBenefit) error
	DeleteBenefit(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetTierBenefits(ctx context.Context, orgID uuid.UUID, tierID uuid.UUID) ([]LoyaltyTierBenefit, error)

	// Points rule management
	ListPointsRules(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsRuleFilters) ([]LoyaltyPointsRule, error)
	CreatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error
	GetPointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyPointsRule, error)
	UpdatePointsRule(ctx context.Context, rule *LoyaltyPointsRule) error
	DeletePointsRule(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error

	// Reward management
	ListRewards(ctx context.Context, orgID uuid.UUID, filters LoyaltyRewardFilters) ([]LoyaltyReward, error)
	CreateReward(ctx context.Context, reward *LoyaltyReward) error
	GetReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyReward, error)
	UpdateReward(ctx context.Context, reward *LoyaltyReward) error
	DeleteReward(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	GetFeaturedRewards(ctx context.Context, orgID uuid.UUID, limit int) ([]LoyaltyReward, error)

	// Redemption management
	ListRedemptions(ctx context.Context, orgID uuid.UUID, filters LoyaltyRedemptionFilters) ([]LoyaltyRedemption, error)
	CreateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error
	GetRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*LoyaltyRedemption, error)
	UpdateRedemption(ctx context.Context, redemption *LoyaltyRedemption) error
	DeleteRedemption(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error
	ApproveRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error
	FulfillRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error
	CancelRedemption(ctx context.Context, orgID uuid.UUID, redemptionID uuid.UUID, userID uuid.UUID) error

	// Points transaction management
	GetCustomerPointsBalance(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (int, error)
	ListPointsTransactions(ctx context.Context, orgID uuid.UUID, filters LoyaltyPointsTransactionFilters) ([]LoyaltyPointsTransaction, error)

	// Tier management
	GetCustomerCurrentTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID) (*LoyaltyTier, error)
	PromoteCustomerTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, tierID uuid.UUID, reason string, userID uuid.UUID) error
	DemoteCustomerTier(ctx context.Context, orgID uuid.UUID, customerID uuid.UUID, tierID uuid.UUID, reason string, userID uuid.UUID) error
}
