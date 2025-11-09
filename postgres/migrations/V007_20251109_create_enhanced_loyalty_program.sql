-- Migration V007: Enhanced Loyalty Program
-- Created: 2025-11-09
-- Description: Adds loyalty tiers, rewards, redemption tracking, and points rules
-- Dependencies: V001, V002, V003, V004

-- ============================================================================
-- LOYALTY TIERS
-- ============================================================================

-- Customer loyalty tiers/levels (Bronze, Silver, Gold, Platinum, etc.)
CREATE TABLE IF NOT EXISTS loyalty_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Tier Information
    tier_code VARCHAR(50) NOT NULL,
    tier_name VARCHAR(100) NOT NULL,
    tier_level INTEGER NOT NULL, -- 1 = lowest, higher = better
    description TEXT,

    -- Requirements
    points_threshold INTEGER NOT NULL, -- Minimum points needed to reach this tier
    annual_spend_threshold NUMERIC(15, 2), -- Minimum annual spend
    purchase_count_threshold INTEGER, -- Minimum number of purchases

    -- Tier Settings
    points_multiplier NUMERIC(5, 2) DEFAULT 1.00, -- e.g., 1.5x points for Gold tier
    discount_percentage NUMERIC(5, 2) DEFAULT 0, -- Automatic discount for tier members

    -- Visual Settings
    tier_color VARCHAR(20),
    tier_icon VARCHAR(50),
    badge_image_url TEXT,

    -- Status
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false, -- Starting tier for new customers

    -- Additional Settings
    sort_order INTEGER DEFAULT 0,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_tier_code_per_org UNIQUE (organization_id, tier_code),
    CONSTRAINT unique_tier_level_per_org UNIQUE (organization_id, tier_level),
    CONSTRAINT positive_thresholds CHECK (
        points_threshold >= 0 AND
        (annual_spend_threshold IS NULL OR annual_spend_threshold >= 0) AND
        (purchase_count_threshold IS NULL OR purchase_count_threshold >= 0)
    ),
    CONSTRAINT valid_multiplier CHECK (points_multiplier > 0)
);

-- Indexes
CREATE INDEX idx_loyalty_tiers_org_id ON loyalty_tiers(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_tiers_tier_code ON loyalty_tiers(tier_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_tiers_tier_level ON loyalty_tiers(tier_level) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_tiers_active ON loyalty_tiers(is_active) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_loyalty_tiers_updated_at
    BEFORE UPDATE ON loyalty_tiers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE loyalty_tiers IS 'Customer loyalty tier definitions with thresholds and benefits';

-- ============================================================================
-- LOYALTY TIER BENEFITS
-- ============================================================================

-- Specific benefits for each loyalty tier
CREATE TABLE IF NOT EXISTS loyalty_tier_benefits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    tier_id UUID NOT NULL REFERENCES loyalty_tiers(id) ON DELETE CASCADE,

    -- Benefit Information
    benefit_code VARCHAR(50) NOT NULL,
    benefit_name VARCHAR(200) NOT NULL,
    benefit_description TEXT,
    benefit_type VARCHAR(50) NOT NULL CHECK (benefit_type IN ('discount', 'free_shipping', 'birthday_bonus', 'early_access', 'priority_support', 'exclusive_products', 'free_product', 'other')),

    -- Benefit Value
    discount_value NUMERIC(15, 2), -- Fixed amount or percentage
    discount_type VARCHAR(20) CHECK (discount_type IN ('percentage', 'fixed_amount')),

    -- Status
    is_active BOOLEAN DEFAULT true,

    -- Display
    sort_order INTEGER DEFAULT 0,
    icon VARCHAR(50),

    -- Additional Information
    terms_and_conditions TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_benefit_per_tier UNIQUE (tier_id, benefit_code)
);

-- Indexes
CREATE INDEX idx_loyalty_tier_benefits_org_id ON loyalty_tier_benefits(organization_id);
CREATE INDEX idx_loyalty_tier_benefits_tier_id ON loyalty_tier_benefits(tier_id);
CREATE INDEX idx_loyalty_tier_benefits_active ON loyalty_tier_benefits(is_active);

-- Auto-update trigger
CREATE TRIGGER update_loyalty_tier_benefits_updated_at
    BEFORE UPDATE ON loyalty_tier_benefits
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE loyalty_tier_benefits IS 'Specific benefits associated with each loyalty tier';

-- ============================================================================
-- LOYALTY POINTS RULES
-- ============================================================================

-- Rules for earning loyalty points
CREATE TABLE IF NOT EXISTS loyalty_points_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Rule Information
    rule_code VARCHAR(50) NOT NULL,
    rule_name VARCHAR(200) NOT NULL,
    description TEXT,
    rule_type VARCHAR(50) NOT NULL CHECK (rule_type IN ('purchase', 'signup', 'birthday', 'referral', 'review', 'social_share', 'manual', 'other')),

    -- Points Calculation
    points_per_amount NUMERIC(10, 2), -- Points per currency unit spent (e.g., 1 point per $1)
    fixed_points INTEGER, -- Fixed points awarded
    multiplier NUMERIC(5, 2) DEFAULT 1.00,

    -- Applicability
    applies_to VARCHAR(50) DEFAULT 'all' CHECK (applies_to IN ('all', 'specific_products', 'specific_categories', 'specific_tiers')),
    applicable_product_ids JSONB DEFAULT '[]',
    applicable_category_ids JSONB DEFAULT '[]',
    applicable_tier_ids JSONB DEFAULT '[]',

    -- Conditions
    minimum_purchase_amount NUMERIC(15, 2),
    maximum_points_per_transaction INTEGER,
    maximum_points_per_day INTEGER,
    maximum_points_per_month INTEGER,

    -- Date Range
    start_date DATE,
    end_date DATE,

    -- Status
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0, -- Higher priority rules apply first

    -- Additional Information
    terms_and_conditions TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_rule_code_per_org UNIQUE (organization_id, rule_code),
    CONSTRAINT valid_date_range CHECK (end_date IS NULL OR start_date IS NULL OR end_date >= start_date)
);

-- Indexes
CREATE INDEX idx_loyalty_points_rules_org_id ON loyalty_points_rules(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_points_rules_rule_code ON loyalty_points_rules(rule_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_points_rules_active ON loyalty_points_rules(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_points_rules_type ON loyalty_points_rules(rule_type) WHERE deleted_at IS NULL;

-- Auto-update trigger
CREATE TRIGGER update_loyalty_points_rules_updated_at
    BEFORE UPDATE ON loyalty_points_rules
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE loyalty_points_rules IS 'Rules for earning loyalty points across different activities';

-- ============================================================================
-- LOYALTY REWARDS
-- ============================================================================

-- Catalog of rewards that customers can redeem with points
CREATE TABLE IF NOT EXISTS loyalty_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Reward Information
    reward_code VARCHAR(50) NOT NULL,
    reward_name VARCHAR(200) NOT NULL,
    description TEXT,
    reward_type VARCHAR(50) NOT NULL CHECK (reward_type IN ('discount_percentage', 'discount_fixed', 'free_product', 'free_shipping', 'gift_card', 'experience', 'other')),

    -- Points Cost
    points_cost INTEGER NOT NULL,

    -- Reward Value
    reward_value NUMERIC(15, 2), -- Monetary value of reward
    discount_percentage NUMERIC(5, 2),
    discount_amount NUMERIC(15, 2),

    -- Product Linkage (for free product rewards)
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    product_variant_id UUID REFERENCES product_variants(id) ON DELETE SET NULL,

    -- Availability
    is_active BOOLEAN DEFAULT true,
    available_from DATE,
    available_to DATE,
    total_available INTEGER, -- NULL = unlimited
    total_redeemed INTEGER DEFAULT 0,
    max_redemptions_per_customer INTEGER,

    -- Tier Restrictions
    minimum_tier_level INTEGER, -- Minimum tier required to redeem
    tier_ids JSONB DEFAULT '[]', -- Specific tiers that can redeem (empty = all)

    -- Visual Settings
    image_url TEXT,
    thumbnail_url TEXT,
    featured BOOLEAN DEFAULT false,

    -- Display
    sort_order INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT false,

    -- Additional Information
    terms_and_conditions TEXT,
    redemption_instructions TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_reward_code_per_org UNIQUE (organization_id, reward_code),
    CONSTRAINT positive_points_cost CHECK (points_cost > 0),
    CONSTRAINT valid_availability CHECK (
        available_to IS NULL OR available_from IS NULL OR available_to >= available_from
    ),
    CONSTRAINT valid_redemptions CHECK (
        total_redeemed >= 0 AND
        (total_available IS NULL OR total_redeemed <= total_available)
    )
);

-- Indexes
CREATE INDEX idx_loyalty_rewards_org_id ON loyalty_rewards(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_rewards_reward_code ON loyalty_rewards(reward_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_rewards_active ON loyalty_rewards(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_loyalty_rewards_featured ON loyalty_rewards(is_featured) WHERE deleted_at IS NULL AND is_active = true;
CREATE INDEX idx_loyalty_rewards_points_cost ON loyalty_rewards(points_cost) WHERE deleted_at IS NULL AND is_active = true;

-- Auto-update trigger
CREATE TRIGGER update_loyalty_rewards_updated_at
    BEFORE UPDATE ON loyalty_rewards
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE loyalty_rewards IS 'Catalog of rewards that customers can redeem with loyalty points';

-- ============================================================================
-- LOYALTY REDEMPTIONS
-- ============================================================================

-- Track customer reward redemptions
CREATE TABLE IF NOT EXISTS loyalty_redemptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,
    reward_id UUID NOT NULL REFERENCES loyalty_rewards(id) ON DELETE RESTRICT,

    -- Redemption Information
    redemption_number VARCHAR(50) NOT NULL,
    redemption_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    points_redeemed INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'fulfilled', 'cancelled', 'expired')),

    -- Usage
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL, -- If reward was used in a transaction
    used_date TIMESTAMP WITH TIME ZONE,
    expiry_date DATE,

    -- Fulfillment
    fulfillment_status VARCHAR(20) DEFAULT 'pending' CHECK (fulfillment_status IN ('pending', 'in_progress', 'completed', 'cancelled')),
    fulfillment_notes TEXT,
    fulfilled_by UUID REFERENCES users(id),
    fulfilled_at TIMESTAMP WITH TIME ZONE,

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT unique_redemption_number_per_org UNIQUE (organization_id, redemption_number),
    CONSTRAINT positive_points CHECK (points_redeemed > 0)
);

-- Indexes
CREATE INDEX idx_loyalty_redemptions_org_id ON loyalty_redemptions(organization_id);
CREATE INDEX idx_loyalty_redemptions_customer_id ON loyalty_redemptions(customer_id);
CREATE INDEX idx_loyalty_redemptions_reward_id ON loyalty_redemptions(reward_id);
CREATE INDEX idx_loyalty_redemptions_number ON loyalty_redemptions(redemption_number);
CREATE INDEX idx_loyalty_redemptions_status ON loyalty_redemptions(status);
CREATE INDEX idx_loyalty_redemptions_date ON loyalty_redemptions(redemption_date);

-- Auto-update trigger
CREATE TRIGGER update_loyalty_redemptions_updated_at
    BEFORE UPDATE ON loyalty_redemptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE loyalty_redemptions IS 'Track customer loyalty reward redemptions and fulfillment';

-- ============================================================================
-- LOYALTY POINTS TRANSACTIONS
-- ============================================================================

-- Detailed log of all loyalty points earned and spent
CREATE TABLE IF NOT EXISTS loyalty_points_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE RESTRICT,

    -- Transaction Information
    transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('earned', 'redeemed', 'expired', 'adjusted', 'bonus', 'refunded')),
    points INTEGER NOT NULL, -- Positive for earned, negative for spent
    balance_after INTEGER NOT NULL,

    -- References
    sale_id UUID REFERENCES sales(id) ON DELETE SET NULL,
    redemption_id UUID REFERENCES loyalty_redemptions(id) ON DELETE SET NULL,
    points_rule_id UUID REFERENCES loyalty_points_rules(id) ON DELETE SET NULL,

    -- Description
    description TEXT,
    reason TEXT,
    notes TEXT,

    -- Expiration (for earned points)
    expiry_date DATE,

    -- Additional Information
    metadata JSONB DEFAULT '{}',

    -- Audit
    transaction_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id),

    -- Constraints
    CONSTRAINT valid_balance CHECK (balance_after >= 0)
);

-- Indexes
CREATE INDEX idx_loyalty_points_transactions_org_id ON loyalty_points_transactions(organization_id);
CREATE INDEX idx_loyalty_points_transactions_customer_id ON loyalty_points_transactions(customer_id);
CREATE INDEX idx_loyalty_points_transactions_type ON loyalty_points_transactions(transaction_type);
CREATE INDEX idx_loyalty_points_transactions_date ON loyalty_points_transactions(transaction_date);
CREATE INDEX idx_loyalty_points_transactions_expiry ON loyalty_points_transactions(expiry_date) WHERE expiry_date IS NOT NULL;

COMMENT ON TABLE loyalty_points_transactions IS 'Complete audit trail of all loyalty points transactions';

-- ============================================================================
-- CUSTOMER TIER HISTORY
-- ============================================================================

-- Track customer tier changes over time
CREATE TABLE IF NOT EXISTS customer_tier_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    tier_id UUID NOT NULL REFERENCES loyalty_tiers(id) ON DELETE RESTRICT,

    -- Tier Change Information
    previous_tier_id UUID REFERENCES loyalty_tiers(id) ON DELETE SET NULL,
    change_type VARCHAR(20) NOT NULL CHECK (change_type IN ('upgrade', 'downgrade', 'initial', 'manual')),
    change_reason TEXT,

    -- Qualification
    qualifying_points INTEGER,
    qualifying_spend NUMERIC(15, 2),
    qualifying_purchases INTEGER,

    -- Dates
    effective_date DATE NOT NULL DEFAULT CURRENT_DATE,
    valid_until DATE, -- For time-limited tier assignments

    -- Additional Information
    notes TEXT,
    metadata JSONB DEFAULT '{}',

    -- Audit
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_customer_tier_history_org_id ON customer_tier_history(organization_id);
CREATE INDEX idx_customer_tier_history_customer_id ON customer_tier_history(customer_id);
CREATE INDEX idx_customer_tier_history_tier_id ON customer_tier_history(tier_id);
CREATE INDEX idx_customer_tier_history_effective_date ON customer_tier_history(effective_date);

COMMENT ON TABLE customer_tier_history IS 'Track customer loyalty tier changes and upgrades/downgrades';

-- ============================================================================
-- ADD TIER TRACKING TO CUSTOMERS TABLE
-- ============================================================================

-- Add current tier to customers table
ALTER TABLE customers ADD COLUMN IF NOT EXISTS current_tier_id UUID REFERENCES loyalty_tiers(id) ON DELETE SET NULL;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tier_since DATE;
ALTER TABLE customers ADD COLUMN IF NOT EXISTS tier_expiry_date DATE;

-- Create index
CREATE INDEX IF NOT EXISTS idx_customers_current_tier ON customers(current_tier_id) WHERE deleted_at IS NULL;

-- ============================================================================
-- HELPER FUNCTION: Update Reward Redemption Count
-- ============================================================================

CREATE OR REPLACE FUNCTION update_reward_redemption_count()
RETURNS TRIGGER AS $$
BEGIN
    -- Increment redeemed count when a new redemption is created
    IF TG_OP = 'INSERT' AND NEW.status IN ('approved', 'fulfilled') THEN
        UPDATE loyalty_rewards
        SET total_redeemed = total_redeemed + 1
        WHERE id = NEW.reward_id;
    END IF;

    -- Decrement if redemption is cancelled
    IF TG_OP = 'UPDATE' AND OLD.status IN ('approved', 'fulfilled') AND NEW.status = 'cancelled' THEN
        UPDATE loyalty_rewards
        SET total_redeemed = GREATEST(0, total_redeemed - 1)
        WHERE id = NEW.reward_id;
    END IF;

    -- Increment if status changes to approved/fulfilled
    IF TG_OP = 'UPDATE' AND OLD.status NOT IN ('approved', 'fulfilled') AND NEW.status IN ('approved', 'fulfilled') THEN
        UPDATE loyalty_rewards
        SET total_redeemed = total_redeemed + 1
        WHERE id = NEW.reward_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for reward redemption count
CREATE TRIGGER trigger_update_reward_redemption_count
    AFTER INSERT OR UPDATE ON loyalty_redemptions
    FOR EACH ROW
    EXECUTE FUNCTION update_reward_redemption_count();

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

-- Enable RLS on all new tables
ALTER TABLE loyalty_tiers ENABLE ROW LEVEL SECURITY;
ALTER TABLE loyalty_tier_benefits ENABLE ROW LEVEL SECURITY;
ALTER TABLE loyalty_points_rules ENABLE ROW LEVEL SECURITY;
ALTER TABLE loyalty_rewards ENABLE ROW LEVEL SECURITY;
ALTER TABLE loyalty_redemptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE loyalty_points_transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE customer_tier_history ENABLE ROW LEVEL SECURITY;

-- ============================================================================
-- RLS POLICIES: LOYALTY_TIERS
-- ============================================================================

CREATE POLICY loyalty_tiers_super_admin_all
    ON loyalty_tiers FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_tiers_select_own_org
    ON loyalty_tiers FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tiers_insert_own_org
    ON loyalty_tiers FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tiers_update_own_org
    ON loyalty_tiers FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tiers_delete_own_org
    ON loyalty_tiers FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOYALTY_TIER_BENEFITS
-- ============================================================================

CREATE POLICY loyalty_tier_benefits_super_admin_all
    ON loyalty_tier_benefits FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_tier_benefits_select_own_org
    ON loyalty_tier_benefits FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tier_benefits_insert_own_org
    ON loyalty_tier_benefits FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tier_benefits_update_own_org
    ON loyalty_tier_benefits FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_tier_benefits_delete_own_org
    ON loyalty_tier_benefits FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOYALTY_POINTS_RULES
-- ============================================================================

CREATE POLICY loyalty_points_rules_super_admin_all
    ON loyalty_points_rules FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_points_rules_select_own_org
    ON loyalty_points_rules FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_points_rules_insert_own_org
    ON loyalty_points_rules FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY loyalty_points_rules_update_own_org
    ON loyalty_points_rules FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_points_rules_delete_own_org
    ON loyalty_points_rules FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOYALTY_REWARDS
-- ============================================================================

CREATE POLICY loyalty_rewards_super_admin_all
    ON loyalty_rewards FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_rewards_select_own_org
    ON loyalty_rewards FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_rewards_insert_own_org
    ON loyalty_rewards FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY loyalty_rewards_update_own_org
    ON loyalty_rewards FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_rewards_delete_own_org
    ON loyalty_rewards FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOYALTY_REDEMPTIONS
-- ============================================================================

CREATE POLICY loyalty_redemptions_super_admin_all
    ON loyalty_redemptions FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_redemptions_select_own_org
    ON loyalty_redemptions FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_redemptions_insert_own_org
    ON loyalty_redemptions FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY loyalty_redemptions_update_own_org
    ON loyalty_redemptions FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_redemptions_delete_own_org
    ON loyalty_redemptions FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: LOYALTY_POINTS_TRANSACTIONS
-- ============================================================================

CREATE POLICY loyalty_points_transactions_super_admin_all
    ON loyalty_points_transactions FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY loyalty_points_transactions_select_own_org
    ON loyalty_points_transactions FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY loyalty_points_transactions_insert_own_org
    ON loyalty_points_transactions FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

-- ============================================================================
-- RLS POLICIES: CUSTOMER_TIER_HISTORY
-- ============================================================================

CREATE POLICY customer_tier_history_super_admin_all
    ON customer_tier_history FOR ALL TO PUBLIC
    USING (is_super_admin());

CREATE POLICY customer_tier_history_select_own_org
    ON customer_tier_history FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY customer_tier_history_insert_own_org
    ON customer_tier_history FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

-- Migration completed successfully
