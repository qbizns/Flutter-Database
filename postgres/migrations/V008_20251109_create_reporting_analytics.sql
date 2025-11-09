-- Migration V008: Reporting & Analytics Infrastructure
-- Created: 2025-11-09
-- Description: Adds materialized views, reporting tables, and analytics infrastructure
-- Dependencies: V001, V002, V003, V004, V005, V006, V007

-- ============================================================================
-- DAILY SALES SUMMARY (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_daily_sales_summary AS
SELECT
    s.organization_id,
    DATE(s.transaction_date) as sale_date,
    s.location_id,
    l.name as location_name,
    -- Sales Metrics
    COUNT(DISTINCT s.id) as total_transactions,
    COUNT(DISTINCT s.customer_id) as unique_customers,
    SUM(s.total_amount) as gross_sales,
    SUM(s.discount_amount) as total_discounts,
    SUM(s.tax_amount) as total_tax,
    SUM(s.total_amount - s.discount_amount) as net_sales,
    -- Average Metrics
    AVG(s.total_amount) as avg_transaction_value,
    SUM(s.total_amount) / NULLIF(COUNT(DISTINCT s.customer_id), 0) as avg_customer_spend,
    -- Item Metrics
    SUM((SELECT SUM(quantity) FROM sale_items si WHERE si.sale_id = s.id)) as total_items_sold,
    SUM((SELECT COUNT(*) FROM sale_items si WHERE si.sale_id = s.id)) as total_line_items,
    -- Payment Methods (aggregated)
    jsonb_object_agg(
        COALESCE(p.payment_method, 'unknown'),
        COALESCE(SUM(p.amount), 0)
    ) FILTER (WHERE p.id IS NOT NULL) as payment_breakdown,
    -- Timestamp
    MAX(s.updated_at) as last_updated
FROM sales s
LEFT JOIN locations l ON s.location_id = l.id
LEFT JOIN payments p ON s.id = p.sale_id
WHERE s.deleted_at IS NULL
GROUP BY s.organization_id, DATE(s.transaction_date), s.location_id, l.name;

-- Indexes for daily sales summary
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_daily_sales_org_date_location
    ON mv_daily_sales_summary(organization_id, sale_date, COALESCE(location_id, '00000000-0000-0000-0000-000000000000'::UUID));
CREATE INDEX IF NOT EXISTS idx_mv_daily_sales_date ON mv_daily_sales_summary(sale_date DESC);
CREATE INDEX IF NOT EXISTS idx_mv_daily_sales_location ON mv_daily_sales_summary(location_id);

COMMENT ON MATERIALIZED VIEW mv_daily_sales_summary IS 'Daily aggregated sales metrics by organization and location';

-- ============================================================================
-- PRODUCT PERFORMANCE (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_product_performance AS
SELECT
    p.id as product_id,
    p.organization_id,
    p.name as product_name,
    p.sku,
    c.name as category_name,
    p.selling_price,
    p.cost_price,
    p.current_stock,
    -- Sales Metrics (Last 30 days)
    COUNT(DISTINCT si.sale_id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as transactions_last_30d,
    SUM(si.quantity) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as units_sold_last_30d,
    SUM(si.line_total) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as revenue_last_30d,
    -- Sales Metrics (Last 90 days)
    COUNT(DISTINCT si.sale_id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '90 days') as transactions_last_90d,
    SUM(si.quantity) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '90 days') as units_sold_last_90d,
    SUM(si.line_total) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '90 days') as revenue_last_90d,
    -- All Time Metrics
    COUNT(DISTINCT si.sale_id) as total_transactions,
    SUM(si.quantity) as total_units_sold,
    SUM(si.line_total) as total_revenue,
    SUM(si.line_total - (si.cost_price * si.quantity)) as total_profit,
    -- Profitability
    CASE
        WHEN SUM(si.line_total) > 0
        THEN ((SUM(si.line_total - (si.cost_price * si.quantity)) / SUM(si.line_total)) * 100)
        ELSE 0
    END as profit_margin_percentage,
    -- Inventory Metrics
    CASE
        WHEN p.track_inventory AND SUM(si.quantity) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') > 0
        THEN p.current_stock / (SUM(si.quantity) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') / 30.0)
        ELSE NULL
    END as days_of_stock_remaining,
    -- Timestamps
    MAX(s.transaction_date) as last_sale_date,
    CURRENT_TIMESTAMP as last_updated
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN sale_items si ON p.id = si.product_id
LEFT JOIN sales s ON si.sale_id = s.id AND s.deleted_at IS NULL
WHERE p.deleted_at IS NULL
GROUP BY p.id, p.organization_id, p.name, p.sku, c.name, p.selling_price, p.cost_price, p.current_stock, p.track_inventory;

-- Indexes for product performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_product_performance_product_id ON mv_product_performance(product_id);
CREATE INDEX IF NOT EXISTS idx_mv_product_performance_org_id ON mv_product_performance(organization_id);
CREATE INDEX IF NOT EXISTS idx_mv_product_performance_revenue_30d ON mv_product_performance(revenue_last_30d DESC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_mv_product_performance_units_30d ON mv_product_performance(units_sold_last_30d DESC NULLS LAST);

COMMENT ON MATERIALIZED VIEW mv_product_performance IS 'Product sales performance metrics and inventory analysis';

-- ============================================================================
-- CUSTOMER ANALYTICS (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_customer_analytics AS
SELECT
    c.id as customer_id,
    c.organization_id,
    c.customer_code,
    c.full_name,
    c.email,
    c.loyalty_points,
    lt.tier_name as loyalty_tier,
    -- Purchase Metrics
    COUNT(DISTINCT s.id) as total_orders,
    SUM(s.total_amount) as lifetime_value,
    AVG(s.total_amount) as avg_order_value,
    SUM(s.total_amount) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as spend_last_30d,
    SUM(s.total_amount) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '90 days') as spend_last_90d,
    SUM(s.total_amount) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '365 days') as spend_last_year,
    -- Frequency
    COUNT(DISTINCT s.id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as orders_last_30d,
    COUNT(DISTINCT s.id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '90 days') as orders_last_90d,
    COUNT(DISTINCT s.id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '365 days') as orders_last_year,
    -- Recency
    MAX(s.transaction_date) as last_purchase_date,
    CURRENT_DATE - MAX(s.transaction_date)::DATE as days_since_last_purchase,
    -- First Purchase
    MIN(s.transaction_date) as first_purchase_date,
    -- Items
    SUM((SELECT SUM(quantity) FROM sale_items si WHERE si.sale_id = s.id)) as total_items_purchased,
    -- Favorite Categories
    (
        SELECT jsonb_agg(jsonb_build_object('category', cat_name, 'count', purchase_count))
        FROM (
            SELECT c2.name as cat_name, COUNT(*) as purchase_count
            FROM sale_items si2
            JOIN sales s2 ON si2.sale_id = s2.id AND s2.customer_id = c.id
            JOIN products p2 ON si2.product_id = p2.id
            LEFT JOIN categories c2 ON p2.category_id = c2.id
            WHERE s2.deleted_at IS NULL
            GROUP BY c2.name
            ORDER BY COUNT(*) DESC
            LIMIT 3
        ) favorite_cats
    ) as favorite_categories,
    -- RFM Segments (simplified)
    CASE
        WHEN CURRENT_DATE - MAX(s.transaction_date)::DATE <= 30 THEN 'Active'
        WHEN CURRENT_DATE - MAX(s.transaction_date)::DATE <= 90 THEN 'At Risk'
        WHEN CURRENT_DATE - MAX(s.transaction_date)::DATE <= 180 THEN 'Dormant'
        ELSE 'Lost'
    END as customer_segment,
    -- Timestamp
    CURRENT_TIMESTAMP as last_updated
FROM customers c
LEFT JOIN sales s ON c.id = s.customer_id AND s.deleted_at IS NULL
LEFT JOIN loyalty_tiers lt ON c.current_tier_id = lt.id
WHERE c.deleted_at IS NULL
GROUP BY c.id, c.organization_id, c.customer_code, c.full_name, c.email, c.loyalty_points, lt.tier_name;

-- Indexes for customer analytics
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_customer_analytics_customer_id ON mv_customer_analytics(customer_id);
CREATE INDEX IF NOT EXISTS idx_mv_customer_analytics_org_id ON mv_customer_analytics(organization_id);
CREATE INDEX IF NOT EXISTS idx_mv_customer_analytics_ltv ON mv_customer_analytics(lifetime_value DESC NULLS LAST);
CREATE INDEX IF NOT EXISTS idx_mv_customer_analytics_segment ON mv_customer_analytics(customer_segment);

COMMENT ON MATERIALIZED VIEW mv_customer_analytics IS 'Customer behavior analytics with RFM segmentation and lifetime value';

-- ============================================================================
-- INVENTORY VALUATION (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_inventory_valuation AS
SELECT
    p.id as product_id,
    p.organization_id,
    p.location_id,
    l.name as location_name,
    p.name as product_name,
    p.sku,
    c.name as category_name,
    p.current_stock,
    p.cost_price,
    p.selling_price,
    -- Valuation
    (p.current_stock * p.cost_price) as inventory_value_at_cost,
    (p.current_stock * p.selling_price) as inventory_value_at_retail,
    (p.current_stock * (p.selling_price - p.cost_price)) as potential_profit,
    -- Stock Status
    CASE
        WHEN NOT p.track_inventory THEN 'not_tracked'
        WHEN p.current_stock <= 0 THEN 'out_of_stock'
        WHEN p.current_stock <= p.low_stock_threshold THEN 'low_stock'
        ELSE 'in_stock'
    END as stock_status,
    -- Product Variants (if any)
    (SELECT COUNT(*) FROM product_variants pv WHERE pv.product_id = p.id AND pv.deleted_at IS NULL) as variant_count,
    (SELECT SUM(current_stock) FROM product_variants pv WHERE pv.product_id = p.id AND pv.deleted_at IS NULL) as variant_total_stock,
    -- Active Batches (if any)
    (SELECT COUNT(*) FROM product_batches pb WHERE pb.product_id = p.id AND pb.status = 'active' AND pb.deleted_at IS NULL) as active_batch_count,
    -- Timestamp
    CURRENT_TIMESTAMP as last_updated
FROM products p
LEFT JOIN locations l ON p.location_id = l.id
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.deleted_at IS NULL AND p.track_inventory = true;

-- Indexes for inventory valuation
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_inventory_valuation_product_id ON mv_inventory_valuation(product_id);
CREATE INDEX IF NOT EXISTS idx_mv_inventory_valuation_org_id ON mv_inventory_valuation(organization_id);
CREATE INDEX IF NOT EXISTS idx_mv_inventory_valuation_location ON mv_inventory_valuation(location_id);
CREATE INDEX IF NOT EXISTS idx_mv_inventory_valuation_stock_status ON mv_inventory_valuation(stock_status);
CREATE INDEX IF NOT EXISTS idx_mv_inventory_valuation_value ON mv_inventory_valuation(inventory_value_at_cost DESC);

COMMENT ON MATERIALIZED VIEW mv_inventory_valuation IS 'Current inventory valuation and stock status analysis';

-- ============================================================================
-- LOCATION PERFORMANCE (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_location_performance AS
SELECT
    l.id as location_id,
    l.organization_id,
    l.location_code,
    l.name as location_name,
    l.location_type,
    l.is_active,
    -- Sales Metrics (Last 30 days)
    COUNT(DISTINCT s.id) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as transactions_last_30d,
    SUM(s.total_amount) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as revenue_last_30d,
    AVG(s.total_amount) FILTER (WHERE s.transaction_date >= CURRENT_DATE - INTERVAL '30 days') as avg_transaction_last_30d,
    -- Sales Metrics (All Time)
    COUNT(DISTINCT s.id) as total_transactions,
    SUM(s.total_amount) as total_revenue,
    AVG(s.total_amount) as avg_transaction_value,
    -- Inventory Metrics
    (SELECT COUNT(*) FROM products p WHERE p.location_id = l.id AND p.deleted_at IS NULL) as product_count,
    (SELECT SUM(current_stock * cost_price) FROM products p WHERE p.location_id = l.id AND p.deleted_at IS NULL) as inventory_value,
    -- Expense Metrics (Last 30 days)
    (SELECT SUM(total_amount) FROM expenses e WHERE e.location_id = l.id AND e.expense_date >= CURRENT_DATE - INTERVAL '30 days' AND e.deleted_at IS NULL) as expenses_last_30d,
    -- Staff Metrics
    (SELECT COUNT(*) FROM shifts sh WHERE sh.location_id = l.id AND sh.start_time >= CURRENT_DATE - INTERVAL '30 days') as shifts_last_30d,
    -- Timestamp
    MAX(s.transaction_date) as last_sale_date,
    CURRENT_TIMESTAMP as last_updated
FROM locations l
LEFT JOIN sales s ON l.id = s.location_id AND s.deleted_at IS NULL
WHERE l.deleted_at IS NULL
GROUP BY l.id, l.organization_id, l.location_code, l.name, l.location_type, l.is_active;

-- Indexes for location performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_location_performance_location_id ON mv_location_performance(location_id);
CREATE INDEX IF NOT EXISTS idx_mv_location_performance_org_id ON mv_location_performance(organization_id);
CREATE INDEX IF NOT EXISTS idx_mv_location_performance_revenue_30d ON mv_location_performance(revenue_last_30d DESC NULLS LAST);

COMMENT ON MATERIALIZED VIEW mv_location_performance IS 'Location-based performance metrics and operational analytics';

-- ============================================================================
-- PROMOTION EFFECTIVENESS (Materialized View)
-- ============================================================================

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_promotion_effectiveness AS
SELECT
    pr.id as promotion_id,
    pr.organization_id,
    pr.promotion_code,
    pr.name as promotion_name,
    pr.promotion_type,
    pr.start_date,
    pr.end_date,
    pr.is_active,
    -- Usage Metrics
    COUNT(DISTINCT pu.id) as total_uses,
    COUNT(DISTINCT pu.customer_id) as unique_customers,
    SUM(pu.discount_amount) as total_discount_given,
    AVG(pu.discount_amount) as avg_discount_per_use,
    -- Sales Metrics
    COUNT(DISTINCT pu.sale_id) as total_sales,
    SUM(s.total_amount) as total_sales_value,
    AVG(s.total_amount) as avg_sale_value,
    -- ROI Metrics (simplified)
    CASE
        WHEN SUM(pu.discount_amount) > 0
        THEN (SUM(s.total_amount) / SUM(pu.discount_amount))
        ELSE NULL
    END as roi_ratio,
    -- Timestamp
    MAX(pu.used_at) as last_used_at,
    CURRENT_TIMESTAMP as last_updated
FROM promotions pr
LEFT JOIN promotion_usage pu ON pr.id = pu.promotion_id
LEFT JOIN sales s ON pu.sale_id = s.id AND s.deleted_at IS NULL
WHERE pr.deleted_at IS NULL
GROUP BY pr.id, pr.organization_id, pr.promotion_code, pr.name, pr.promotion_type, pr.start_date, pr.end_date, pr.is_active;

-- Indexes for promotion effectiveness
CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_promotion_effectiveness_promo_id ON mv_promotion_effectiveness(promotion_id);
CREATE INDEX IF NOT EXISTS idx_mv_promotion_effectiveness_org_id ON mv_promotion_effectiveness(organization_id);
CREATE INDEX IF NOT EXISTS idx_mv_promotion_effectiveness_total_uses ON mv_promotion_effectiveness(total_uses DESC);

COMMENT ON MATERIALIZED VIEW mv_promotion_effectiveness IS 'Promotion campaign effectiveness and ROI analysis';

-- ============================================================================
-- HELPER FUNCTION: Refresh All Materialized Views
-- ============================================================================

CREATE OR REPLACE FUNCTION refresh_all_analytics_views(concurrent_refresh BOOLEAN DEFAULT true)
RETURNS TABLE(view_name TEXT, refresh_status TEXT, execution_time INTERVAL) AS $$
DECLARE
    start_time TIMESTAMP;
    end_time TIMESTAMP;
    view_record RECORD;
BEGIN
    FOR view_record IN
        SELECT matviewname FROM pg_matviews WHERE schemaname = 'public'
    LOOP
        start_time := clock_timestamp();

        BEGIN
            IF concurrent_refresh THEN
                EXECUTE format('REFRESH MATERIALIZED VIEW CONCURRENTLY %I', view_record.matviewname);
            ELSE
                EXECUTE format('REFRESH MATERIALIZED VIEW %I', view_record.matviewname);
            END IF;

            end_time := clock_timestamp();

            view_name := view_record.matviewname;
            refresh_status := 'SUCCESS';
            execution_time := end_time - start_time;
            RETURN NEXT;

        EXCEPTION WHEN OTHERS THEN
            end_time := clock_timestamp();

            view_name := view_record.matviewname;
            refresh_status := 'FAILED: ' || SQLERRM;
            execution_time := end_time - start_time;
            RETURN NEXT;
        END;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION refresh_all_analytics_views IS 'Refresh all materialized views and return execution status';

-- ============================================================================
-- SCHEDULED REFRESH (pg_cron integration - optional)
-- ============================================================================

-- Note: This requires pg_cron extension to be enabled
-- Uncomment the following lines if pg_cron is available:
--
-- -- Refresh analytics views every hour
-- SELECT cron.schedule(
--     'refresh-analytics-hourly',
--     '0 * * * *', -- Every hour
--     $$SELECT refresh_all_analytics_views(true)$$
-- );
--
-- -- Refresh analytics views at midnight (more thorough)
-- SELECT cron.schedule(
--     'refresh-analytics-daily',
--     '0 0 * * *', -- Every day at midnight
--     $$SELECT refresh_all_analytics_views(false)$$
-- );

-- ============================================================================
-- REPORTING HELPER VIEWS (Regular Views for Dynamic Queries)
-- ============================================================================

-- Top selling products (last 30 days)
CREATE OR REPLACE VIEW v_top_selling_products_30d AS
SELECT
    organization_id,
    product_name,
    sku,
    units_sold_last_30d,
    revenue_last_30d,
    transactions_last_30d,
    ROW_NUMBER() OVER (PARTITION BY organization_id ORDER BY revenue_last_30d DESC) as rank
FROM mv_product_performance
WHERE revenue_last_30d > 0
ORDER BY organization_id, revenue_last_30d DESC;

COMMENT ON VIEW v_top_selling_products_30d IS 'Top selling products ranked by revenue in last 30 days';

-- Low stock alerts
CREATE OR REPLACE VIEW v_low_stock_alerts AS
SELECT
    organization_id,
    location_name,
    product_name,
    sku,
    current_stock,
    stock_status,
    inventory_value_at_cost,
    days_of_stock_remaining
FROM mv_inventory_valuation
INNER JOIN mv_product_performance USING (product_id)
WHERE stock_status IN ('low_stock', 'out_of_stock')
ORDER BY organization_id, CASE stock_status WHEN 'out_of_stock' THEN 1 WHEN 'low_stock' THEN 2 END, days_of_stock_remaining;

COMMENT ON VIEW v_low_stock_alerts IS 'Products with low or out of stock status requiring attention';

-- Customer segments summary
CREATE OR REPLACE VIEW v_customer_segments_summary AS
SELECT
    organization_id,
    customer_segment,
    COUNT(*) as customer_count,
    SUM(lifetime_value) as segment_lifetime_value,
    AVG(lifetime_value) as avg_lifetime_value,
    SUM(orders_last_30d) as segment_orders_last_30d,
    AVG(orders_last_30d) as avg_orders_per_customer_30d
FROM mv_customer_analytics
GROUP BY organization_id, customer_segment
ORDER BY organization_id, customer_segment;

COMMENT ON VIEW v_customer_segments_summary IS 'Customer segmentation summary with key metrics per segment';

-- ============================================================================
-- INITIAL VIEW REFRESH
-- ============================================================================

-- Refresh all views on migration (non-concurrent for initial population)
SELECT refresh_all_analytics_views(false);

-- Migration completed successfully
