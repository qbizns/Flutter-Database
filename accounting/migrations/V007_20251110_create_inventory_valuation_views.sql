-- ============================================================================
-- Migration: V007 - Create Inventory Valuation Views
-- Description: Reporting views for inventory valuation and COGS analysis
-- Author: System
-- Date: 2025-11-10
-- ============================================================================

BEGIN;

-- ============================================================================
-- VIEW: Current Inventory Valuation by Product
-- ============================================================================

CREATE OR REPLACE VIEW view_inventory_valuation_by_product AS
SELECT
    icl.organization_id,
    icl.product_id,
    icl.location_id,
    SUM(icl.remaining_quantity) as total_quantity,
    SUM(icl.unit_cost * icl.remaining_quantity) as total_value,
    CASE
        WHEN SUM(icl.remaining_quantity) > 0
        THEN SUM(icl.unit_cost * icl.remaining_quantity) / SUM(icl.remaining_quantity)
        ELSE 0
    END as weighted_average_cost,
    MIN(icl.unit_cost) as min_unit_cost,
    MAX(icl.unit_cost) as max_unit_cost,
    COUNT(DISTINCT icl.id) as cost_layer_count,
    MIN(icl.layer_date) as oldest_layer_date,
    MAX(icl.layer_date) as newest_layer_date
FROM inventory_cost_layers icl
WHERE icl.remaining_quantity > 0
  AND icl.deleted_at IS NULL
GROUP BY icl.organization_id, icl.product_id, icl.location_id;

COMMENT ON VIEW view_inventory_valuation_by_product IS 'Current inventory valuation with weighted average cost per product/location';

-- ============================================================================
-- VIEW: Inventory Valuation Summary by Location
-- ============================================================================

CREATE OR REPLACE VIEW view_inventory_valuation_by_location AS
SELECT
    icl.organization_id,
    icl.location_id,
    COUNT(DISTINCT icl.product_id) as unique_products,
    SUM(icl.remaining_quantity) as total_quantity,
    SUM(icl.unit_cost * icl.remaining_quantity) as total_inventory_value,
    SUM(icl.unit_cost * icl.remaining_quantity) FILTER (WHERE icl.layer_date < CURRENT_DATE - INTERVAL '90 days') as slow_moving_value,
    SUM(icl.unit_cost * icl.remaining_quantity) FILTER (WHERE icl.layer_date < CURRENT_DATE - INTERVAL '180 days') as aged_inventory_value
FROM inventory_cost_layers icl
WHERE icl.remaining_quantity > 0
  AND icl.deleted_at IS NULL
GROUP BY icl.organization_id, icl.location_id;

COMMENT ON VIEW view_inventory_valuation_by_location IS 'Inventory valuation summary by location with aging analysis';

-- ============================================================================
-- VIEW: COGS by Period (from posted journal entries)
-- ============================================================================

CREATE OR REPLACE VIEW view_cogs_by_period AS
SELECT
    gl.organization_id,
    DATE_TRUNC('month', gl.transaction_date) as period_month,
    DATE_TRUNC('year', gl.transaction_date) as period_year,
    coa.account_code,
    coa.account_name,
    SUM(gl.debit_amount) as total_cogs_debit,
    SUM(gl.credit_amount) as total_cogs_credit,
    SUM(gl.debit_amount - gl.credit_amount) as net_cogs
FROM general_ledger gl
INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
INNER JOIN account_types at ON coa.account_type_id = at.id
WHERE at.type_name IN ('Cost of Sales', 'COGS', 'Cost of Goods Sold')
  OR coa.account_code LIKE '5%'  -- Typical COGS account range
GROUP BY
    gl.organization_id,
    DATE_TRUNC('month', gl.transaction_date),
    DATE_TRUNC('year', gl.transaction_date),
    coa.account_code,
    coa.account_name
ORDER BY period_year DESC, period_month DESC;

COMMENT ON VIEW view_cogs_by_period IS 'Cost of Goods Sold summarized by period from general ledger';

-- ============================================================================
-- VIEW: Inventory Turnover Analysis
-- ============================================================================

CREATE OR REPLACE VIEW view_inventory_turnover AS
WITH inventory_avg AS (
    SELECT
        organization_id,
        product_id,
        AVG(total_value) as avg_inventory_value
    FROM (
        SELECT
            organization_id,
            product_id,
            DATE_TRUNC('month', layer_date) as month,
            SUM(unit_cost * remaining_quantity) as total_value
        FROM inventory_cost_layers
        WHERE deleted_at IS NULL
          AND layer_date >= CURRENT_DATE - INTERVAL '12 months'
        GROUP BY organization_id, product_id, DATE_TRUNC('month', layer_date)
    ) monthly_inv
    GROUP BY organization_id, product_id
),
cogs_annual AS (
    SELECT
        gl.organization_id,
        -- Would need product mapping from journal entry line metadata
        SUM(gl.debit_amount - gl.credit_amount) as annual_cogs
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    WHERE at.type_name IN ('Cost of Sales', 'COGS')
      AND gl.transaction_date >= CURRENT_DATE - INTERVAL '12 months'
    GROUP BY gl.organization_id
)
SELECT
    ia.organization_id,
    ia.product_id,
    ia.avg_inventory_value,
    ca.annual_cogs,
    CASE
        WHEN ia.avg_inventory_value > 0
        THEN ca.annual_cogs / ia.avg_inventory_value
        ELSE 0
    END as turnover_ratio,
    CASE
        WHEN ca.annual_cogs > 0
        THEN (ia.avg_inventory_value / ca.annual_cogs) * 365
        ELSE 0
    END as days_inventory_outstanding
FROM inventory_avg ia
CROSS JOIN cogs_annual ca ON ia.organization_id = ca.organization_id;

COMMENT ON VIEW view_inventory_turnover IS 'Inventory turnover ratio and days inventory outstanding';

-- ============================================================================
-- VIEW: Cost Layer Aging Report
-- ============================================================================

CREATE OR REPLACE VIEW view_cost_layer_aging AS
SELECT
    icl.organization_id,
    icl.product_id,
    icl.location_id,
    icl.layer_date,
    CURRENT_DATE - icl.layer_date as age_in_days,
    CASE
        WHEN CURRENT_DATE - icl.layer_date <= 30 THEN '0-30 days'
        WHEN CURRENT_DATE - icl.layer_date <= 60 THEN '31-60 days'
        WHEN CURRENT_DATE - icl.layer_date <= 90 THEN '61-90 days'
        WHEN CURRENT_DATE - icl.layer_date <= 180 THEN '91-180 days'
        WHEN CURRENT_DATE - icl.layer_date <= 365 THEN '181-365 days'
        ELSE 'Over 1 year'
    END as age_bucket,
    icl.remaining_quantity,
    icl.unit_cost,
    icl.unit_cost * icl.remaining_quantity as layer_value,
    icl.source_transaction_type,
    icl.source_reference
FROM inventory_cost_layers icl
WHERE icl.remaining_quantity > 0
  AND icl.deleted_at IS NULL
ORDER BY icl.organization_id, icl.product_id, icl.layer_date;

COMMENT ON VIEW view_cost_layer_aging IS 'Inventory cost layer aging analysis for identifying slow-moving stock';

-- ============================================================================
-- VIEW: Inventory Valuation vs GL Balance Reconciliation
-- ============================================================================

CREATE OR REPLACE VIEW view_inventory_reconciliation AS
WITH inventory_valuation AS (
    SELECT
        organization_id,
        SUM(unit_cost * remaining_quantity) as total_inventory_value
    FROM inventory_cost_layers
    WHERE remaining_quantity > 0
      AND deleted_at IS NULL
    GROUP BY organization_id
),
gl_inventory_balance AS (
    SELECT
        gl.organization_id,
        SUM(gl.debit_amount - gl.credit_amount) as gl_inventory_balance
    FROM general_ledger gl
    INNER JOIN chart_of_accounts coa ON gl.account_id = coa.id
    INNER JOIN account_types at ON coa.account_type_id = at.id
    WHERE at.type_name IN ('Inventory', 'Current Assets')
      AND (coa.account_code LIKE '1%' OR coa.account_name ILIKE '%inventory%')
    GROUP BY gl.organization_id
)
SELECT
    COALESCE(iv.organization_id, gl.organization_id) as organization_id,
    COALESCE(iv.total_inventory_value, 0) as inventory_valuation,
    COALESCE(gl.gl_inventory_balance, 0) as gl_balance,
    COALESCE(gl.gl_inventory_balance, 0) - COALESCE(iv.total_inventory_value, 0) as variance,
    CASE
        WHEN ABS(COALESCE(gl.gl_inventory_balance, 0) - COALESCE(iv.total_inventory_value, 0)) < 0.01 THEN 'Balanced'
        WHEN COALESCE(gl.gl_inventory_balance, 0) > COALESCE(iv.total_inventory_value, 0) THEN 'GL Higher'
        ELSE 'Valuation Higher'
    END as variance_status
FROM inventory_valuation iv
FULL OUTER JOIN gl_inventory_balance gl ON iv.organization_id = gl.organization_id;

COMMENT ON VIEW view_inventory_reconciliation IS 'Reconciliation between inventory cost layers and GL inventory accounts';

-- ============================================================================
-- VIEW: FIFO Cost Layer Simulation
-- ============================================================================

CREATE OR REPLACE VIEW view_fifo_cost_layers AS
SELECT
    icl.organization_id,
    icl.product_id,
    icl.location_id,
    icl.layer_date,
    icl.unit_cost,
    icl.remaining_quantity,
    icl.unit_cost * icl.remaining_quantity as layer_value,
    ROW_NUMBER() OVER (
        PARTITION BY icl.organization_id, icl.product_id, icl.location_id
        ORDER BY icl.layer_date ASC, icl.created_at ASC
    ) as fifo_sequence
FROM inventory_cost_layers icl
WHERE icl.remaining_quantity > 0
  AND icl.deleted_at IS NULL
ORDER BY icl.organization_id, icl.product_id, icl.location_id, icl.layer_date;

COMMENT ON VIEW view_fifo_cost_layers IS 'Cost layers ordered for FIFO consumption (oldest first)';

-- ============================================================================
-- COMPLETION
-- ============================================================================

COMMIT;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '============================================';
    RAISE NOTICE 'Migration V007 completed successfully!';
    RAISE NOTICE 'Inventory Valuation Views created:';
    RAISE NOTICE ' - view_inventory_valuation_by_product';
    RAISE NOTICE ' - view_inventory_valuation_by_location';
    RAISE NOTICE ' - view_cogs_by_period';
    RAISE NOTICE ' - view_inventory_turnover';
    RAISE NOTICE ' - view_cost_layer_aging';
    RAISE NOTICE ' - view_inventory_reconciliation';
    RAISE NOTICE ' - view_fifo_cost_layers';
    RAISE NOTICE '============================================';
END $$;
