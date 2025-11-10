-- ============================================================================
-- Migration: V017 - Add Currency Support to POS Tables
-- Description: Add multi-currency fields to organizations and key POS transaction tables
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- ORGANIZATIONS TABLE - Add Base Currency
-- ============================================================================

ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS base_currency_code VARCHAR(3) DEFAULT 'USD',
    ADD COLUMN IF NOT EXISTS currency_display_format VARCHAR(50) DEFAULT 'symbol',  -- 'symbol', 'code', 'both'
    ADD COLUMN IF NOT EXISTS decimal_separator VARCHAR(1) DEFAULT '.',
    ADD COLUMN IF NOT EXISTS thousands_separator VARCHAR(1) DEFAULT ',';

CREATE INDEX IF NOT EXISTS idx_organizations_base_currency_code ON organizations(base_currency_code);

COMMENT ON COLUMN organizations.base_currency_code IS 'ISO 4217 currency code (USD, EUR, GBP, SAR, EGP, etc.)';
COMMENT ON COLUMN organizations.currency_display_format IS 'How to display currency: symbol ($), code (USD), or both ($USD)';

-- ============================================================================
-- SALES TABLE - Add Currency Fields
-- ============================================================================

ALTER TABLE sales
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),
    ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8) DEFAULT 1.0,
    ADD COLUMN IF NOT EXISTS base_currency_total NUMERIC(12, 2);  -- Total in organization base currency

CREATE INDEX IF NOT EXISTS idx_sales_currency_code ON sales(currency_code) WHERE currency_code IS NOT NULL;

COMMENT ON COLUMN sales.currency_code IS 'Transaction currency (defaults to organization base currency)';
COMMENT ON COLUMN sales.exchange_rate IS 'Exchange rate to base currency at transaction time';
COMMENT ON COLUMN sales.base_currency_total IS 'Total amount in base currency for reporting';

-- ============================================================================
-- PURCHASE_ORDERS TABLE - Add Currency Fields
-- ============================================================================

ALTER TABLE purchase_orders
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),
    ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8) DEFAULT 1.0,
    ADD COLUMN IF NOT EXISTS base_currency_total NUMERIC(12, 2);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_currency_code ON purchase_orders(currency_code) WHERE currency_code IS NOT NULL;

COMMENT ON COLUMN purchase_orders.currency_code IS 'Purchase order currency';
COMMENT ON COLUMN purchase_orders.exchange_rate IS 'Exchange rate to base currency';

-- ============================================================================
-- PAYMENTS TABLE - Add Currency Fields
-- ============================================================================

ALTER TABLE payments
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),
    ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8) DEFAULT 1.0,
    ADD COLUMN IF NOT EXISTS base_currency_amount NUMERIC(12, 2);

CREATE INDEX IF NOT EXISTS idx_payments_currency_code ON payments(currency_code) WHERE currency_code IS NOT NULL;

COMMENT ON COLUMN payments.currency_code IS 'Payment currency';
COMMENT ON COLUMN payments.base_currency_amount IS 'Payment amount in base currency';

-- ============================================================================
-- GOODS_RECEIPTS TABLE - Add Currency Fields
-- ============================================================================

ALTER TABLE goods_receipts
    ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),
    ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8) DEFAULT 1.0,
    ADD COLUMN IF NOT EXISTS base_currency_total NUMERIC(12, 2);

CREATE INDEX IF NOT EXISTS idx_goods_receipts_currency_code ON goods_receipts(currency_code) WHERE currency_code IS NOT NULL;

COMMENT ON COLUMN goods_receipts.currency_code IS 'Receipt currency (usually matches purchase order)';

-- ============================================================================
-- EXPENSES TABLE - Add Currency Fields (if exists)
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'expenses') THEN
        ALTER TABLE expenses
            ADD COLUMN IF NOT EXISTS currency_code VARCHAR(3),
            ADD COLUMN IF NOT EXISTS exchange_rate NUMERIC(20, 8) DEFAULT 1.0,
            ADD COLUMN IF NOT EXISTS base_currency_amount NUMERIC(12, 2);

        CREATE INDEX IF NOT EXISTS idx_expenses_currency_code ON expenses(currency_code) WHERE currency_code IS NOT NULL;
    END IF;
END $$;

-- ============================================================================
-- HELPER FUNCTION: Get Current Exchange Rate
-- ============================================================================

CREATE OR REPLACE FUNCTION get_current_exchange_rate(
    p_from_currency VARCHAR,
    p_to_currency VARCHAR,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS NUMERIC AS $$
DECLARE
    v_rate NUMERIC;
BEGIN
    -- If same currency, return 1
    IF p_from_currency = p_to_currency THEN
        RETURN 1.0;
    END IF;

    -- Try to get rate from accounting.currency_rates if accounting module is available
    -- For now, return 1.0 as default (application layer should provide rates)
    -- In production, this would query exchange rate services or accounting.currency_rates

    RETURN 1.0;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_current_exchange_rate IS 'Get exchange rate between two currencies (placeholder - implement with actual rate source)';

-- ============================================================================
-- HELPER FUNCTION: Convert to Base Currency
-- ============================================================================

CREATE OR REPLACE FUNCTION convert_to_base_currency(
    p_amount NUMERIC,
    p_from_currency VARCHAR,
    p_to_currency VARCHAR,
    p_exchange_rate NUMERIC DEFAULT NULL,
    p_date DATE DEFAULT CURRENT_DATE
) RETURNS NUMERIC AS $$
DECLARE
    v_rate NUMERIC;
    v_result NUMERIC;
BEGIN
    -- If same currency, return original amount
    IF p_from_currency = p_to_currency THEN
        RETURN p_amount;
    END IF;

    -- Use provided exchange rate or fetch current rate
    IF p_exchange_rate IS NOT NULL THEN
        v_rate := p_exchange_rate;
    ELSE
        v_rate := get_current_exchange_rate(p_from_currency, p_to_currency, p_date);
    END IF;

    -- Convert
    v_result := p_amount * v_rate;

    RETURN ROUND(v_result, 2);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION convert_to_base_currency IS 'Convert amount from one currency to another using exchange rate';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V017 completed successfully!';
    RAISE NOTICE 'Added currency support to:';
    RAISE NOTICE ' - organizations (base_currency_code)';
    RAISE NOTICE ' - sales (currency_code, exchange_rate)';
    RAISE NOTICE ' - purchase_orders (currency_code, exchange_rate)';
    RAISE NOTICE ' - payments (currency_code, exchange_rate)';
    RAISE NOTICE ' - goods_receipts (currency_code, exchange_rate)';
    RAISE NOTICE ' - expenses (currency_code, exchange_rate)';
    RAISE NOTICE ' - Helper functions for currency conversion';
    RAISE NOTICE '============================================';
END $$;
