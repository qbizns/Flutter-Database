-- ============================================================================
-- Migration: V021 - Add Organization Feature Flags
-- Description: Enable/disable modules per organization (accounting, e-invoicing, etc.)
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- ORGANIZATIONS TABLE - Add Features JSONB Column
-- ============================================================================

ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS features JSONB DEFAULT '{
        "accounting": false,
        "e_invoicing": false,
        "advanced_inventory": true,
        "multi_location": true,
        "multi_currency": false,
        "loyalty_program": true,
        "delivery_management": false,
        "restaurant_mode": false,
        "table_management": false,
        "kitchen_display": false,
        "online_ordering": false,
        "api_access": false
    }'::jsonb;

CREATE INDEX IF NOT EXISTS idx_organizations_features ON organizations USING gin(features);

COMMENT ON COLUMN organizations.features IS 'Feature flags for enabled/disabled modules (accounting, e-invoicing, inventory, etc.)';

-- ============================================================================
-- ORGANIZATION_FEATURES TABLE (Normalized Alternative)
-- Description: Alternative normalized approach for feature flags
-- ============================================================================

CREATE TABLE IF NOT EXISTS organization_features (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Organization (for multi-tenancy)
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Feature Configuration
    feature_key VARCHAR(100) NOT NULL,  -- 'accounting', 'e_invoicing', 'multi_currency', etc.
    is_enabled BOOLEAN DEFAULT FALSE,
    is_available BOOLEAN DEFAULT TRUE,  -- Feature available in subscription plan

    -- Configuration
    configuration JSONB DEFAULT '{}',   -- Feature-specific settings
    limits JSONB DEFAULT '{}',          -- e.g., {"max_users": 10, "max_locations": 5}

    -- Effective Dates
    enabled_at TIMESTAMP WITH TIME ZONE,
    disabled_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,  -- For time-limited trials

    -- Notes
    notes TEXT,

    -- Metadata
    metadata JSONB DEFAULT '{}',

    -- Audit Fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    updated_by UUID REFERENCES users(id),

    -- Constraints
    UNIQUE(organization_id, feature_key)
);

-- Indexes
CREATE INDEX idx_organization_features_organization_id ON organization_features(organization_id);
CREATE INDEX idx_organization_features_feature_key ON organization_features(feature_key);
CREATE INDEX idx_organization_features_is_enabled ON organization_features(is_enabled) WHERE is_enabled = TRUE;

-- Trigger
CREATE TRIGGER update_organization_features_updated_at
    BEFORE UPDATE ON organization_features
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE organization_features IS 'Normalized feature flags per organization (alternative to JSONB column)';
COMMENT ON COLUMN organization_features.feature_key IS 'Unique feature identifier (accounting, e_invoicing, etc.)';
COMMENT ON COLUMN organization_features.is_enabled IS 'Whether feature is currently enabled';
COMMENT ON COLUMN organization_features.is_available IS 'Whether feature is available in subscription plan';
COMMENT ON COLUMN organization_features.configuration IS 'Feature-specific configuration settings';
COMMENT ON COLUMN organization_features.limits IS 'Feature usage limits (max_users, max_locations, etc.)';

-- ============================================================================
-- HELPER FUNCTION: Check if Feature is Enabled
-- ============================================================================

CREATE OR REPLACE FUNCTION is_feature_enabled(
    p_organization_id UUID,
    p_feature_key VARCHAR
) RETURNS BOOLEAN AS $$
DECLARE
    v_is_enabled BOOLEAN;
    v_features JSONB;
BEGIN
    -- First check normalized table
    SELECT is_enabled INTO v_is_enabled
    FROM organization_features
    WHERE organization_id = p_organization_id
      AND feature_key = p_feature_key
      AND deleted_at IS NULL
    LIMIT 1;

    IF FOUND THEN
        RETURN v_is_enabled;
    END IF;

    -- Fall back to JSONB column
    SELECT features INTO v_features
    FROM organizations
    WHERE id = p_organization_id;

    IF v_features IS NOT NULL AND v_features ? p_feature_key THEN
        RETURN (v_features->>p_feature_key)::BOOLEAN;
    END IF;

    -- Default to false if not found
    RETURN FALSE;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION is_feature_enabled IS 'Check if a feature is enabled for an organization (checks both normalized table and JSONB column)';

-- ============================================================================
-- HELPER FUNCTION: Enable/Disable Feature
-- ============================================================================

CREATE OR REPLACE FUNCTION set_feature_enabled(
    p_organization_id UUID,
    p_feature_key VARCHAR,
    p_is_enabled BOOLEAN,
    p_user_id UUID DEFAULT NULL
) RETURNS JSONB AS $$
DECLARE
    v_feature_id UUID;
    v_result JSONB;
BEGIN
    -- Insert or update in normalized table
    INSERT INTO organization_features (
        organization_id,
        feature_key,
        is_enabled,
        enabled_at,
        disabled_at,
        created_by,
        updated_by
    ) VALUES (
        p_organization_id,
        p_feature_key,
        p_is_enabled,
        CASE WHEN p_is_enabled THEN CURRENT_TIMESTAMP ELSE NULL END,
        CASE WHEN NOT p_is_enabled THEN CURRENT_TIMESTAMP ELSE NULL END,
        p_user_id,
        p_user_id
    )
    ON CONFLICT (organization_id, feature_key)
    DO UPDATE SET
        is_enabled = p_is_enabled,
        enabled_at = CASE WHEN p_is_enabled THEN CURRENT_TIMESTAMP ELSE organization_features.enabled_at END,
        disabled_at = CASE WHEN NOT p_is_enabled THEN CURRENT_TIMESTAMP ELSE organization_features.disabled_at END,
        updated_at = CURRENT_TIMESTAMP,
        updated_by = p_user_id
    RETURNING id INTO v_feature_id;

    -- Also update JSONB column for consistency
    UPDATE organizations
    SET features = jsonb_set(
        COALESCE(features, '{}'::jsonb),
        ARRAY[p_feature_key],
        to_jsonb(p_is_enabled)
    ),
    updated_at = CURRENT_TIMESTAMP
    WHERE id = p_organization_id;

    v_result := jsonb_build_object(
        'success', TRUE,
        'organization_id', p_organization_id,
        'feature_key', p_feature_key,
        'is_enabled', p_is_enabled,
        'feature_id', v_feature_id,
        'updated_at', CURRENT_TIMESTAMP
    );

    RETURN v_result;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION set_feature_enabled IS 'Enable or disable a feature for an organization (updates both normalized table and JSONB column)';

-- ============================================================================
-- HELPER VIEW: Organization Features Summary
-- ============================================================================

CREATE OR REPLACE VIEW view_organization_features AS
SELECT
    o.id as organization_id,
    o.name as organization_name,
    o.features as features_json,
    -- Extract common features from JSONB
    (o.features->>'accounting')::BOOLEAN as accounting_enabled,
    (o.features->>'e_invoicing')::BOOLEAN as e_invoicing_enabled,
    (o.features->>'multi_currency')::BOOLEAN as multi_currency_enabled,
    (o.features->>'advanced_inventory')::BOOLEAN as advanced_inventory_enabled,
    (o.features->>'loyalty_program')::BOOLEAN as loyalty_program_enabled,
    (o.features->>'delivery_management')::BOOLEAN as delivery_management_enabled,
    (o.features->>'restaurant_mode')::BOOLEAN as restaurant_mode_enabled,
    o.base_currency_code,
    o.is_active as organization_active,
    o.created_at
FROM organizations o
WHERE o.deleted_at IS NULL;

COMMENT ON VIEW view_organization_features IS 'Summary of organization features with extracted boolean flags';

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================

ALTER TABLE organization_features ENABLE ROW LEVEL SECURITY;

CREATE POLICY organization_features_tenant_isolation ON organization_features
    USING (organization_id IN (
        SELECT organization_id
        FROM user_organizations
        WHERE user_id = auth.uid()
    ));

COMMENT ON POLICY organization_features_tenant_isolation ON organization_features
    IS 'Ensure users can only access features for their organizations';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V021 completed successfully!';
    RAISE NOTICE 'Organization Feature Flags created:';
    RAISE NOTICE ' - Added features JSONB column to organizations';
    RAISE NOTICE ' - Created organization_features table';
    RAISE NOTICE ' - is_feature_enabled() function';
    RAISE NOTICE ' - set_feature_enabled() function';
    RAISE NOTICE ' - view_organization_features view';
    RAISE NOTICE ' - RLS policies applied';
    RAISE NOTICE '============================================';
END $$;
