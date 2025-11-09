-- ============================================================================
-- Seed Data: Enhanced Features (V005-V007)
-- Description: Seeds inventory transfers, advanced inventory, and loyalty program data
-- ============================================================================

BEGIN;

-- ============================================================================
-- BYPASS RLS FOR SEEDING
-- ============================================================================

SELECT set_user_context(NULL, NULL, TRUE);

-- ============================================================================
-- INVENTORY TRANSFERS (V005)
-- ============================================================================

INSERT INTO inventory_transfers (id, organization_id, transfer_number, from_location_id, to_location_id, status, transfer_date, requested_date, approved_date, shipped_date, expected_delivery_date, carrier, tracking_number, reason, requested_by, approved_by, shipped_by)
VALUES
    -- Demo Retail Store: Transfer from warehouse to Manhattan store
    ('it000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'TR-001', 'loc00001-0000-0000-0000-000000000003', 'loc00001-0000-0000-0000-000000000001', 'received', '2025-11-01', '2025-11-01 09:00:00', '2025-11-01 10:00:00', '2025-11-01 14:00:00', '2025-11-01', 'Internal Transport', 'TR-20251101-001', 'Restock for weekend sales', '10000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004'),

    -- Demo Retail Store: Transfer from Manhattan to Brooklyn
    ('it000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'TR-002', 'loc00001-0000-0000-0000-000000000001', 'loc00001-0000-0000-0000-000000000002', 'in_transit', '2025-11-08', '2025-11-08 10:00:00', '2025-11-08 11:00:00', '2025-11-08 15:00:00', '2025-11-09', 'Internal Transport', 'TR-20251108-001', 'Balance inventory between stores', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000003'),

    -- Coffee Corner: Transfer between locations
    ('it000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'TR-CC-001', 'loc00002-0000-0000-0000-000000000001', 'loc00002-0000-0000-0000-000000000002', 'received', '2025-11-05', '2025-11-05 08:00:00', '2025-11-05 08:30:00', '2025-11-05 09:00:00', '2025-11-05', 'Own Vehicle', 'CC-TR-001', 'Restock Mission location', '20000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000004')
ON CONFLICT DO NOTHING;

-- Transfer items
INSERT INTO inventory_transfer_items (id, organization_id, inventory_transfer_id, product_id, product_name, product_sku, quantity_requested, quantity_shipped, quantity_received, unit_cost, total_cost, item_status)
VALUES
    -- TR-001 items
    ('iti00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'it000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000003', 'AirPods Pro', 'AIRPODS-PRO', 10, 10, 10, 179.00, 1790.00, 'received'),
    ('iti00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'it000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000004', 'Universal Phone Case', 'PHONE-CASE', 50, 50, 50, 5.00, 250.00, 'received'),
    ('iti00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'it000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000005', 'USB-C Cable 2m', 'USB-C-CABLE', 30, 30, 30, 3.00, 90.00, 'received'),

    -- TR-002 items (in transit)
    ('iti00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'it000001-0000-0000-0000-000000000002', 'p0000001-0000-0000-0000-000000000006', 'Men''s Cotton T-Shirt', 'MENS-TSHIRT', 20, 20, 0, 8.00, 160.00, 'shipped'),
    ('iti00001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', 'it000001-0000-0000-0000-000000000002', 'p0000001-0000-0000-0000-000000000007', 'Women''s Summer Dress', 'WOMENS-DRESS', 15, 15, 0, 25.00, 375.00, 'shipped'),

    -- TR-CC-001 items
    ('iti00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'it000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000005', 'Butter Croissant', 'CROISSANT', 100, 100, 100, 1.20, 120.00, 'received'),
    ('iti00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'it000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000006', 'Blueberry Muffin', 'MUFFIN', 80, 80, 80, 1.00, 80.00, 'received')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCT SERIAL NUMBERS (V006)
-- ============================================================================

INSERT INTO product_serial_numbers (id, organization_id, product_id, serial_number, status, purchase_date, purchase_cost, warranty_start_date, warranty_end_date, warranty_provider, location_id)
VALUES
    -- iPhone 15 Pro serial numbers
    ('psn00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'IP15P-2025-001-AAA', 'in_stock', '2025-10-15', 899.00, '2025-10-15', '2026-10-15', 'Apple Inc.', 'loc00001-0000-0000-0000-000000000001'),
    ('psn00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'IP15P-2025-001-BBB', 'in_stock', '2025-10-15', 899.00, '2025-10-15', '2026-10-15', 'Apple Inc.', 'loc00001-0000-0000-0000-000000000001'),
    ('psn00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'IP15P-2025-001-CCC', 'sold', '2025-10-15', 899.00, '2025-10-15', '2026-10-15', 'Apple Inc.', 'loc00001-0000-0000-0000-000000000001'),

    -- Samsung S24 serial numbers
    ('psn00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'SGS24-2025-001-XXX', 'in_stock', '2025-10-20', 699.00, '2025-10-20', '2026-10-20', 'Samsung Electronics', 'loc00001-0000-0000-0000-000000000001'),
    ('psn00001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'SGS24-2025-001-YYY', 'in_stock', '2025-10-20', 699.00, '2025-10-20', '2026-10-20', 'Samsung Electronics', 'loc00001-0000-0000-0000-000000000001'),
    ('psn00001-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'SGS24-2025-001-ZZZ', 'sold', '2025-10-20', 699.00, '2025-10-20', '2026-10-20', 'Samsung Electronics', 'loc00001-0000-0000-0000-000000000001'),

    -- MacBook Pro serial numbers
    ('psn00003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'p0000003-0000-0000-0000-000000000001', 'MBP16-2025-001-AAA', 'in_stock', '2025-10-01', 2199.00, '2025-10-01', '2026-10-01', 'Apple Inc.', 'loc00003-0000-0000-0000-000000000001'),
    ('psn00003-0000-0000-0000-000000000002', '33333333-3333-3333-3333-333333333333', 'p0000003-0000-0000-0000-000000000001', 'MBP16-2025-001-BBB', 'in_stock', '2025-10-01', 2199.00, '2025-10-01', '2026-10-01', 'Apple Inc.', 'loc00003-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCT BATCHES (V006)
-- ============================================================================

INSERT INTO product_batches (id, organization_id, product_id, batch_number, status, initial_quantity, current_quantity, manufacturing_date, expiration_date, received_date, supplier_id, unit_cost, total_cost, quality_status, location_id)
VALUES
    -- Coffee Corner: Croissant batches
    ('pb000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000005', 'CROIS-20251108', 'active', 100, 85, '2025-11-08', '2025-11-10', '2025-11-08', 'su000002-0000-0000-0000-000000000002', 1.20, 120.00, 'passed', 'loc00002-0000-0000-0000-000000000001'),
    ('pb000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000005', 'CROIS-20251107', 'active', 100, 25, '2025-11-07', '2025-11-09', '2025-11-07', 'su000002-0000-0000-0000-000000000002', 1.20, 120.00, 'passed', 'loc00002-0000-0000-0000-000000000001'),

    -- Coffee Corner: Muffin batches
    ('pb000002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000006', 'MUFF-20251108', 'active', 80, 65, '2025-11-08', '2025-11-11', '2025-11-08', 'su000002-0000-0000-0000-000000000002', 1.00, 80.00, 'passed', 'loc00002-0000-0000-0000-000000000001'),
    ('pb000002-0000-0000-0000-000000000004', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000006', 'MUFF-20251106', 'expired', 80, 0, '2025-11-06', '2025-11-09', '2025-11-06', 'su000002-0000-0000-0000-000000000002', 1.00, 80.00, 'passed', 'loc00002-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- CYCLE COUNTS (V006)
-- ============================================================================

INSERT INTO cycle_counts (id, organization_id, location_id, count_number, count_date, count_type, status, scheduled_date, started_at, completed_at, total_items_planned, total_items_counted, items_with_variance, created_by, counted_by, approved_by)
VALUES
    -- Demo Retail Store: Completed cycle count
    ('cc000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', 'CC-2025-001', '2025-11-01', 'cycle', 'completed', '2025-11-01', '2025-11-01 09:00:00', '2025-11-01 14:30:00', 7, 7, 2, '10000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000002'),

    -- Coffee Corner: In progress cycle count
    ('cc000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', 'CC-2025-CC-001', '2025-11-09', 'cycle', 'in_progress', '2025-11-09', '2025-11-09 08:00:00', NULL, 6, 4, 0, '20000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000003', NULL)
ON CONFLICT DO NOTHING;

-- Cycle count items
INSERT INTO cycle_count_items (id, organization_id, cycle_count_id, product_id, product_name, product_sku, system_quantity, counted_quantity, unit_cost, status, counted_at, counted_by)
VALUES
    -- CC-2025-001 items
    ('cci00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'cc000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000003', 'AirPods Pro', 'AIRPODS-PRO', 50, 48, 179.00, 'adjusted', '2025-11-01 10:15:00', '10000000-0000-0000-0000-000000000004'),
    ('cci00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'cc000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000004', 'Universal Phone Case', 'PHONE-CASE', 200, 203, 5.00, 'adjusted', '2025-11-01 10:30:00', '10000000-0000-0000-0000-000000000004'),
    ('cci00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'cc000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000005', 'USB-C Cable 2m', 'USB-C-CABLE', 150, 150, 3.00, 'counted', '2025-11-01 10:45:00', '10000000-0000-0000-0000-000000000004'),

    -- CC-2025-CC-001 items (in progress)
    ('cci00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'cc000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000005', 'Butter Croissant', 'CROISSANT', 85, 85, 1.20, 'counted', '2025-11-09 08:30:00', '20000000-0000-0000-0000-000000000003'),
    ('cci00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'cc000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000006', 'Blueberry Muffin', 'MUFFIN', 65, 65, 1.00, 'counted', '2025-11-09 08:45:00', '20000000-0000-0000-0000-000000000003')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOYALTY TIERS (V007)
-- ============================================================================

INSERT INTO loyalty_tiers (id, organization_id, tier_code, tier_name, tier_level, description, points_threshold, annual_spend_threshold, points_multiplier, discount_percentage, tier_color, is_active, is_default, sort_order)
VALUES
    -- Demo Retail Store Tiers
    ('lt000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'BRONZE', 'Bronze Member', 1, 'Entry level membership', 0, 0, 1.00, 0, '#CD7F32', true, true, 1),
    ('lt000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'SILVER', 'Silver Member', 2, 'Earn 1.25x points on purchases', 500, 500.00, 1.25, 5.00, '#C0C0C0', true, false, 2),
    ('lt000001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'GOLD', 'Gold Member', 3, 'Earn 1.5x points and get 10% discount', 1500, 1500.00, 1.50, 10.00, '#FFD700', true, false, 3),
    ('lt000001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'PLATINUM', 'Platinum Member', 4, 'Earn 2x points and get 15% discount', 5000, 5000.00, 2.00, 15.00, '#E5E4E2', true, false, 4),

    -- Coffee Corner Tiers
    ('lt000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'REGULAR', 'Regular Customer', 1, 'Welcome to our loyalty program', 0, 0, 1.00, 0, '#8B4513', true, true, 1),
    ('lt000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'PREMIUM', 'Premium Coffee Lover', 2, 'Earn 1.5x points and free upgrades', 300, 200.00, 1.50, 10.00, '#D2691E', true, false, 2),
    ('lt000002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'VIP', 'VIP Coffee Club', 3, 'Earn 2x points and exclusive benefits', 1000, 500.00, 2.00, 15.00, '#4B3621', true, false, 3)
ON CONFLICT DO NOTHING;

-- Update customers with tiers
UPDATE customers SET current_tier_id = 'lt000001-0000-0000-0000-000000000002', tier_since = '2025-01-15' WHERE id = 'cu000001-0000-0000-0000-000000000001'; -- Alice - Silver
UPDATE customers SET current_tier_id = 'lt000001-0000-0000-0000-000000000003', tier_since = '2025-03-01' WHERE id = 'cu000001-0000-0000-0000-000000000002'; -- Bob - Gold
UPDATE customers SET current_tier_id = 'lt000001-0000-0000-0000-000000000001', tier_since = '2025-05-10' WHERE id = 'cu000001-0000-0000-0000-000000000003'; -- Carol - Bronze
UPDATE customers SET current_tier_id = 'lt000001-0000-0000-0000-000000000004', tier_since = '2024-11-01' WHERE id = 'cu000001-0000-0000-0000-000000000004'; -- David - Platinum

UPDATE customers SET current_tier_id = 'lt000002-0000-0000-0000-000000000002', tier_since = '2025-02-01' WHERE id = 'cu000002-0000-0000-0000-000000000001'; -- Frank - Premium
UPDATE customers SET current_tier_id = 'lt000002-0000-0000-0000-000000000003', tier_since = '2024-12-01' WHERE id = 'cu000002-0000-0000-0000-000000000002'; -- Grace - VIP

-- ============================================================================
-- LOYALTY TIER BENEFITS (V007)
-- ============================================================================

INSERT INTO loyalty_tier_benefits (id, organization_id, tier_id, benefit_code, benefit_name, benefit_description, benefit_type, discount_value, discount_type, is_active, sort_order)
VALUES
    -- Silver tier benefits
    ('ltb00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000002', 'BIRTHDAY_10', 'Birthday Bonus', 'Get 10% off during your birthday month', 'birthday_bonus', 10, 'percentage', true, 1),
    ('ltb00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000002', 'EARLY_ACCESS', 'Early Access to Sales', 'Get notified 24 hours before public sales', 'early_access', NULL, NULL, true, 2),

    -- Gold tier benefits
    ('ltb00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000003', 'BIRTHDAY_20', 'Birthday Bonus', 'Get 20% off during your birthday month', 'birthday_bonus', 20, 'percentage', true, 1),
    ('ltb00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000003', 'FREE_SHIPPING', 'Free Shipping', 'Free shipping on all online orders', 'free_shipping', NULL, NULL, true, 2),
    ('ltb00001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000003', 'PRIORITY_SUPPORT', 'Priority Customer Support', 'Get priority assistance from our team', 'priority_support', NULL, NULL, true, 3),

    -- Platinum tier benefits
    ('ltb00001-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000004', 'BIRTHDAY_30', 'Birthday Bonus', 'Get 30% off during your birthday month', 'birthday_bonus', 30, 'percentage', true, 1),
    ('ltb00001-0000-0000-0000-000000000007', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000004', 'FREE_SHIPPING_PRIORITY', 'Free Express Shipping', 'Free express shipping on all orders', 'free_shipping', NULL, NULL, true, 2),
    ('ltb00001-0000-0000-0000-000000000008', '11111111-1111-1111-1111-111111111111', 'lt000001-0000-0000-0000-000000000004', 'EXCLUSIVE_ACCESS', 'Exclusive Product Access', 'First access to limited edition products', 'exclusive_products', NULL, NULL, true, 3)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOYALTY POINTS RULES (V007)
-- ============================================================================

INSERT INTO loyalty_points_rules (id, organization_id, rule_code, rule_name, description, rule_type, points_per_amount, fixed_points, multiplier, applies_to, is_active, priority, start_date)
VALUES
    -- Demo Retail Store rules
    ('lpr00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'PURCHASE_POINTS', 'Purchase Points', 'Earn 1 point for every $1 spent', 'purchase', 1.00, NULL, 1.00, 'all', true, 1, '2025-01-01'),
    ('lpr00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'SIGNUP_BONUS', 'Sign Up Bonus', 'Get 100 points when you join', 'signup', NULL, 100, 1.00, 'all', true, 1, '2025-01-01'),
    ('lpr00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'BIRTHDAY_BONUS', 'Birthday Bonus', 'Get 200 bonus points on your birthday', 'birthday', NULL, 200, 1.00, 'all', true, 1, '2025-01-01'),
    ('lpr00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'REFERRAL_BONUS', 'Referral Bonus', 'Get 500 points for each friend you refer', 'referral', NULL, 500, 1.00, 'all', true, 1, '2025-01-01'),

    -- Coffee Corner rules
    ('lpr00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'PURCHASE_POINTS_CC', 'Purchase Points', 'Earn 2 points for every $1 spent', 'purchase', 2.00, NULL, 1.00, 'all', true, 1, '2025-01-01'),
    ('lpr00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'SIGNUP_BONUS_CC', 'Welcome Bonus', 'Get 50 points when you join', 'signup', NULL, 50, 1.00, 'all', true, 1, '2025-01-01')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOYALTY REWARDS (V007)
-- ============================================================================

INSERT INTO loyalty_rewards (id, organization_id, reward_code, reward_name, description, reward_type, points_cost, reward_value, discount_percentage, discount_amount, is_active, available_from, total_available, max_redemptions_per_customer, sort_order, is_featured)
VALUES
    -- Demo Retail Store rewards
    ('lr000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'DISC_10', '$10 Off Purchase', 'Get $10 off your next purchase of $50 or more', 'discount_fixed', 500, 10.00, NULL, 10.00, true, '2025-01-01', NULL, 2, 1, true),
    ('lr000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'DISC_25', '$25 Off Purchase', 'Get $25 off your next purchase of $100 or more', 'discount_fixed', 1000, 25.00, NULL, 25.00, true, '2025-01-01', NULL, 2, 2, true),
    ('lr000001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'DISC_50', '$50 Off Purchase', 'Get $50 off your next purchase of $200 or more', 'discount_fixed', 2000, 50.00, NULL, 50.00, true, '2025-01-01', NULL, 1, 3, false),
    ('lr000001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'FREE_CASE', 'Free Phone Case', 'Get a free universal phone case', 'free_product', 300, 19.99, NULL, NULL, true, '2025-01-01', 100, 1, 4, true),

    -- Coffee Corner rewards
    ('lr000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'FREE_COFFEE', 'Free Coffee', 'Get a free coffee of your choice', 'free_product', 100, 5.00, NULL, NULL, true, '2025-01-01', NULL, NULL, 1, true),
    ('lr000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'FREE_PASTRY', 'Free Pastry', 'Get a free pastry with any drink purchase', 'free_product', 150, 3.50, NULL, NULL, true, '2025-01-01', NULL, NULL, 2, true),
    ('lr000002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'DISC_5', '$5 Off', 'Get $5 off your purchase', 'discount_fixed', 200, 5.00, NULL, 5.00, true, '2025-01-01', NULL, 2, 3, false)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOYALTY REDEMPTIONS (V007)
-- ============================================================================

INSERT INTO loyalty_redemptions (id, organization_id, customer_id, reward_id, redemption_number, redemption_date, points_redeemed, status, fulfillment_status, fulfilled_at)
VALUES
    ('lrd00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000002', 'lr000001-0000-0000-0000-000000000001', 'RDM-001', '2025-10-15 14:30:00', 500, 'fulfilled', 'completed', '2025-10-15 14:35:00'),
    ('lrd00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000004', 'lr000001-0000-0000-0000-000000000002', 'RDM-002', '2025-10-20 11:15:00', 1000, 'fulfilled', 'completed', '2025-10-20 11:20:00'),
    ('lrd00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'cu000002-0000-0000-0000-000000000001', 'lr000002-0000-0000-0000-000000000001', 'RDM-CC-001', '2025-11-05 09:00:00', 100, 'fulfilled', 'completed', '2025-11-05 09:05:00'),
    ('lrd00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'cu000002-0000-0000-0000-000000000002', 'lr000002-0000-0000-0000-000000000002', 'RDM-CC-002', '2025-11-06 10:30:00', 150, 'fulfilled', 'completed', '2025-11-06 10:35:00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOYALTY POINTS TRANSACTIONS (V007)
-- ============================================================================

INSERT INTO loyalty_points_transactions (id, organization_id, customer_id, transaction_type, points, balance_after, description, transaction_date)
VALUES
    -- Demo Retail Store transactions
    ('lpt00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000001', 'earned', 100, 100, 'Signup bonus', '2025-01-15 10:00:00'),
    ('lpt00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000001', 'earned', 50, 150, 'Purchase points from transaction', '2025-02-01 14:30:00'),

    ('lpt00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000002', 'earned', 100, 100, 'Signup bonus', '2025-03-01 09:00:00'),
    ('lpt00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000002', 'earned', 200, 300, 'Purchase points from transaction', '2025-03-15 16:45:00'),
    ('lpt00001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000002', 'redeemed', -500, -200, 'Redeemed $10 off reward', '2025-10-15 14:30:00'),
    ('lpt00001-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111', 'cu000001-0000-0000-0000-000000000002', 'earned', 450, 250, 'Purchase points from recent transactions', '2025-11-01 12:00:00'),

    -- Coffee Corner transactions
    ('lpt00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'cu000002-0000-0000-0000-000000000001', 'earned', 50, 50, 'Welcome bonus', '2025-02-01 08:00:00'),
    ('lpt00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'cu000002-0000-0000-0000-000000000001', 'earned', 150, 200, 'Purchase points', '2025-02-15 09:30:00'),
    ('lpt00002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'cu000002-0000-0000-0000-000000000001', 'redeemed', -100, 100, 'Redeemed free coffee', '2025-11-05 09:00:00')
ON CONFLICT DO NOTHING;

COMMIT;

-- ============================================================================
-- SUMMARY
-- ============================================================================

DO $$
DECLARE
    transfer_count INTEGER;
    serial_count INTEGER;
    batch_count INTEGER;
    cycle_count INTEGER;
    tier_count INTEGER;
    reward_count INTEGER;
    redemption_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO transfer_count FROM inventory_transfers WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO serial_count FROM product_serial_numbers WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO batch_count FROM product_batches WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO cycle_count FROM cycle_counts WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO tier_count FROM loyalty_tiers WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO reward_count FROM loyalty_rewards WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO redemption_count FROM loyalty_redemptions;

    RAISE NOTICE '=================================================';
    RAISE NOTICE 'Enhanced Features Data Seeded Successfully';
    RAISE NOTICE '=================================================';
    RAISE NOTICE 'Inventory Transfers: %', transfer_count;
    RAISE NOTICE 'Serial Numbers: %', serial_count;
    RAISE NOTICE 'Product Batches: %', batch_count;
    RAISE NOTICE 'Cycle Counts: %', cycle_count;
    RAISE NOTICE 'Loyalty Tiers: %', tier_count;
    RAISE NOTICE 'Loyalty Rewards: %', reward_count;
    RAISE NOTICE 'Loyalty Redemptions: %', redemption_count;
    RAISE NOTICE '=================================================';
END $$;
