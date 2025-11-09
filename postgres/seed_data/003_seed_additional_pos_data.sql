-- ============================================================================
-- Seed Data: Additional POS Data
-- Description: Seeds suppliers, locations, product variants, promotions, expenses, and shifts
-- ============================================================================

BEGIN;

-- ============================================================================
-- BYPASS RLS FOR SEEDING
-- ============================================================================

-- Set super admin context to bypass RLS policies during seeding
SELECT set_user_context(NULL, NULL, TRUE);

-- ============================================================================
-- SUPPLIERS
-- ============================================================================

INSERT INTO suppliers (id, organization_id, supplier_code, name, contact_person, email, phone, address, city, state, country, postal_code, payment_terms, credit_limit, status)
VALUES
    -- Demo Retail Store Suppliers
    ('su000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'SUP-001', 'TechWholesale Inc', 'John Anderson', 'john@techwholesale.com', '+1-555-0101', '500 Tech Park Drive', 'San Jose', 'CA', 'USA', '95113', 'Net 30', 50000.00, 'active'),
    ('su000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'SUP-002', 'Fashion Distributors LLC', 'Sarah Chen', 'sarah@fashiondist.com', '+1-555-0102', '200 Fashion Avenue', 'New York', 'NY', 'USA', '10018', 'Net 45', 30000.00, 'active'),
    ('su000001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'SUP-003', 'Mobile Devices Direct', 'Mike Johnson', 'mike@mobiledev.com', '+1-555-0103', '800 Wireless Blvd', 'Dallas', 'TX', 'USA', '75201', 'Net 30', 75000.00, 'active'),

    -- Coffee Corner Suppliers
    ('su000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'SUP-CC-001', 'Premium Coffee Roasters', 'Emma Williams', 'emma@premiumcoffee.com', '+1-555-0201', '100 Bean Street', 'Seattle', 'WA', 'USA', '98101', 'Net 15', 15000.00, 'active'),
    ('su000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'SUP-CC-002', 'Fresh Bakery Supplies', 'David Lee', 'david@freshbakery.com', '+1-555-0202', '50 Pastry Lane', 'Portland', 'OR', 'USA', '97201', 'Cash on Delivery', 5000.00, 'active'),

    -- Tech Gadgets Pro Suppliers
    ('su000003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'SUP-TG-001', 'Computer Components Corp', 'Lisa Brown', 'lisa@compcomponents.com', '+1-555-0301', '300 Silicon Valley Drive', 'Cupertino', 'CA', 'USA', '95014', 'Net 30', 100000.00, 'active')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- LOCATIONS / BRANCHES
-- ============================================================================

INSERT INTO locations (id, organization_id, location_code, name, location_type, phone, email, manager_user_id, address_line1, city, state, country, postal_code, timezone, is_active, is_primary, allow_sales, allow_purchases)
VALUES
    -- Demo Retail Store Locations
    ('loc00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'LOC-001', 'Main Store - Manhattan', 'store', '+1-555-1001', 'manhattan@demoretail.com', '10000000-0000-0000-0000-000000000002', '123 Fifth Avenue', 'New York', 'NY', 'USA', '10010', 'America/New_York', TRUE, TRUE, TRUE, TRUE),
    ('loc00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'LOC-002', 'Brooklyn Branch', 'store', '+1-555-1002', 'brooklyn@demoretail.com', '10000000-0000-0000-0000-000000000004', '456 Brooklyn Heights', 'Brooklyn', 'NY', 'USA', '11201', 'America/New_York', TRUE, FALSE, TRUE, FALSE),
    ('loc00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'LOC-003', 'Central Warehouse', 'warehouse', '+1-555-1003', 'warehouse@demoretail.com', NULL, '789 Industrial Park', 'Jersey City', 'NJ', 'USA', '07302', 'America/New_York', TRUE, FALSE, FALSE, TRUE),

    -- Coffee Corner Locations
    ('loc00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'CC-001', 'Coffee Corner - Downtown', 'store', '+1-555-2001', 'downtown@coffeecorner.com', '20000000-0000-0000-0000-000000000002', '101 Market Street', 'San Francisco', 'CA', 'USA', '94103', 'America/Los_Angeles', TRUE, TRUE, TRUE, TRUE),
    ('loc00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'CC-002', 'Coffee Corner - Mission', 'store', '+1-555-2002', 'mission@coffeecorner.com', NULL, '303 Valencia Street', 'San Francisco', 'CA', 'USA', '94110', 'America/Los_Angeles', TRUE, FALSE, TRUE, FALSE),

    -- Tech Gadgets Pro Locations
    ('loc00003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'TG-001', 'Tech Gadgets Pro - Main Store', 'store', '+1-555-3001', 'main@techgadgets.com', '30000000-0000-0000-0000-000000000002', '101 Congress Avenue', 'Austin', 'TX', 'USA', '78701', 'America/Chicago', TRUE, TRUE, TRUE, TRUE)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PRODUCT VARIANTS
-- ============================================================================

INSERT INTO product_variants (id, organization_id, product_id, variant_name, sku, barcode, attributes, cost_price, selling_price, current_stock, is_active, is_default, sort_order)
VALUES
    -- iPhone 15 Pro Variants (Colors)
    ('pv000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'iPhone 15 Pro - Natural Titanium', 'IPHONE-15-PRO-NAT', '1234567890011', '{"color": "Natural Titanium", "storage": "128GB"}', 899.00, 999.00, 10, TRUE, TRUE, 1),
    ('pv000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'iPhone 15 Pro - Blue Titanium', 'IPHONE-15-PRO-BLU', '1234567890012', '{"color": "Blue Titanium", "storage": "128GB"}', 899.00, 999.00, 8, TRUE, FALSE, 2),
    ('pv000001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'iPhone 15 Pro - White Titanium', 'IPHONE-15-PRO-WHT', '1234567890013', '{"color": "White Titanium", "storage": "128GB"}', 899.00, 999.00, 7, TRUE, FALSE, 3),

    -- Samsung Galaxy S24 Variants
    ('pv000001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'Samsung S24 - Onyx Black', 'SAMSUNG-S24-BLK', '1234567890021', '{"color": "Onyx Black", "storage": "256GB"}', 699.00, 849.00, 12, TRUE, TRUE, 1),
    ('pv000001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'Samsung S24 - Marble Gray', 'SAMSUNG-S24-GRY', '1234567890022', '{"color": "Marble Gray", "storage": "256GB"}', 699.00, 849.00, 10, TRUE, FALSE, 2),
    ('pv000001-0000-0000-0000-000000000006', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'Samsung S24 - Cobalt Violet', 'SAMSUNG-S24-VIO', '1234567890023', '{"color": "Cobalt Violet", "storage": "256GB"}', 699.00, 849.00, 8, TRUE, FALSE, 3),

    -- Men's T-Shirt Variants (Sizes)
    ('pv000001-0000-0000-0000-000000000007', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000006', 'Men''s T-Shirt - Small', 'MENS-TSHIRT-S', '1234567890061', '{"size": "S", "color": "Navy Blue"}', 8.00, 24.99, 25, TRUE, FALSE, 1),
    ('pv000001-0000-0000-0000-000000000008', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000006', 'Men''s T-Shirt - Medium', 'MENS-TSHIRT-M', '1234567890062', '{"size": "M", "color": "Navy Blue"}', 8.00, 24.99, 35, TRUE, TRUE, 2),
    ('pv000001-0000-0000-0000-000000000009', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000006', 'Men''s T-Shirt - Large', 'MENS-TSHIRT-L', '1234567890063', '{"size": "L", "color": "Navy Blue"}', 8.00, 24.99, 40, TRUE, FALSE, 3),

    -- Women's Dress Variants (Sizes)
    ('pv000001-0000-0000-0000-000000000010', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000007', 'Women''s Dress - Small', 'WOMENS-DRESS-S', '1234567890071', '{"size": "S", "color": "Floral"}', 25.00, 59.99, 10, TRUE, FALSE, 1),
    ('pv000001-0000-0000-0000-000000000011', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000007', 'Women''s Dress - Medium', 'WOMENS-DRESS-M', '1234567890072', '{"size": "M", "color": "Floral"}', 25.00, 59.99, 20, TRUE, TRUE, 2),
    ('pv000001-0000-0000-0000-000000000012', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000007', 'Women''s Dress - Large', 'WOMENS-DRESS-L', '1234567890073', '{"size": "L", "color": "Floral"}', 25.00, 59.99, 10, TRUE, FALSE, 3)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PROMOTIONS
-- ============================================================================

INSERT INTO promotions (id, organization_id, promotion_code, name, description, promotion_type, discount_value, applies_to, minimum_purchase_amount, start_date, end_date, is_active, is_combinable, priority, max_uses_total)
VALUES
    -- Demo Retail Store Promotions
    ('prm00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'WELCOME10', 'Welcome Discount', 'Get 10% off your first purchase', 'percentage', 10.00, 'all', 50.00, '2025-01-01 00:00:00+00', '2025-12-31 23:59:59+00', TRUE, TRUE, 1, 1000),
    ('prm00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'SUMMER25', 'Summer Sale', '25% off on clothing items', 'percentage', 25.00, 'specific_categories', 0, '2025-06-01 00:00:00+00', '2025-08-31 23:59:59+00', TRUE, FALSE, 2, NULL),
    ('prm00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'TECH50', 'Tech Sale - $50 Off', '$50 off on tech products over $500', 'fixed_amount', 50.00, 'specific_categories', 500.00, '2025-11-01 00:00:00+00', '2025-11-30 23:59:59+00', TRUE, TRUE, 1, 500),
    ('prm00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'BUY2GET1', 'Buy 2 Get 1 Free - Accessories', 'Buy 2 accessories, get 1 free', 'buy_x_get_y', 0, 'specific_categories', 0, '2025-01-01 00:00:00+00', '2025-12-31 23:59:59+00', TRUE, FALSE, 3, NULL),

    -- Coffee Corner Promotions
    ('prm00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'LOYALTY15', 'Loyalty Member Special', '15% off for loyalty members', 'percentage', 15.00, 'all', 0, '2025-01-01 00:00:00+00', '2025-12-31 23:59:59+00', TRUE, TRUE, 1, NULL),
    ('prm00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'MORNING5', 'Morning Rush', '$5 off purchases over $20 before 10 AM', 'fixed_amount', 5.00, 'all', 20.00, '2025-01-01 00:00:00+00', '2025-12-31 23:59:59+00', TRUE, FALSE, 1, NULL),

    -- Tech Gadgets Pro Promotions
    ('prm00003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'GAMING20', 'Gaming Gear Sale', '20% off all gaming accessories', 'percentage', 20.00, 'specific_categories', 0, '2025-11-01 00:00:00+00', '2025-11-30 23:59:59+00', TRUE, TRUE, 1, NULL)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PURCHASE ORDERS
-- ============================================================================

INSERT INTO purchase_orders (id, organization_id, supplier_id, po_number, po_date, expected_delivery_date, status, subtotal, tax_amount, total_amount, payment_status)
VALUES
    -- Demo Retail Store Purchase Orders
    ('po000001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'su000001-0000-0000-0000-000000000003', 'PO-2025-001', '2025-11-01', '2025-11-10', 'received', 27470.00, 2747.00, 30217.00, 'paid'),
    ('po000001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'su000001-0000-0000-0000-000000000001', 'PO-2025-002', '2025-11-05', '2025-11-15', 'ordered', 8950.00, 895.00, 9845.00, 'pending'),

    -- Coffee Corner Purchase Orders
    ('po000002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'su000002-0000-0000-0000-000000000001', 'PO-CC-001', '2025-11-01', '2025-11-05', 'received', 2500.00, 125.00, 2625.00, 'paid'),
    ('po000002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'su000002-0000-0000-0000-000000000002', 'PO-CC-002', '2025-11-08', '2025-11-10', 'pending', 1200.00, 60.00, 1260.00, 'pending')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PURCHASE ORDER ITEMS
-- ============================================================================

INSERT INTO purchase_order_items (id, organization_id, purchase_order_id, product_id, product_name, product_sku, quantity_ordered, quantity_received, unit_cost, line_total)
VALUES
    -- PO-2025-001 Items (iPhone and Samsung)
    ('poi00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'po000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000001', 'iPhone 15 Pro', 'IPHONE-15-PRO', 15, 15, 899.00, 13485.00),
    ('poi00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'po000001-0000-0000-0000-000000000001', 'p0000001-0000-0000-0000-000000000002', 'Samsung Galaxy S24', 'SAMSUNG-S24', 20, 20, 699.00, 13980.00),

    -- PO-2025-002 Items (AirPods)
    ('poi00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'po000001-0000-0000-0000-000000000002', 'p0000001-0000-0000-0000-000000000003', 'AirPods Pro', 'AIRPODS-PRO', 50, 0, 179.00, 8950.00),

    -- Coffee Corner PO Items
    ('poi00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'po000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000005', 'Butter Croissant', 'CROISSANT', 1000, 1000, 1.20, 1200.00),
    ('poi00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'po000002-0000-0000-0000-000000000001', 'p0000002-0000-0000-0000-000000000006', 'Blueberry Muffin', 'MUFFIN', 1000, 1000, 1.00, 1000.00),
    ('poi00002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'po000002-0000-0000-0000-000000000002', 'p0000002-0000-0000-0000-000000000005', 'Butter Croissant', 'CROISSANT', 1000, 0, 1.20, 1200.00)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- EXPENSES
-- ============================================================================

INSERT INTO expenses (id, organization_id, location_id, expense_number, expense_date, category, subcategory, payee_name, payment_method, amount, tax_amount, total_amount, status)
VALUES
    -- Demo Retail Store Expenses
    ('exp00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', 'EXP-001', '2025-11-01', 'rent', 'Store Rent', 'City Real Estate Management', 'bank_transfer', 5000.00, 0, 5000.00, 'paid'),
    ('exp00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', 'EXP-002', '2025-11-02', 'utilities', 'Electricity', 'City Power Company', 'bank_transfer', 450.00, 45.00, 495.00, 'paid'),
    ('exp00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', 'EXP-003', '2025-11-03', 'marketing', 'Social Media Ads', 'Digital Marketing Agency', 'card', 800.00, 80.00, 880.00, 'paid'),
    ('exp00001-0000-0000-0000-000000000004', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000002', 'EXP-004', '2025-11-01', 'rent', 'Store Rent', 'Brooklyn Properties LLC', 'bank_transfer', 3500.00, 0, 3500.00, 'paid'),
    ('exp00001-0000-0000-0000-000000000005', '11111111-1111-1111-1111-111111111111', NULL, 'EXP-005', '2025-11-05', 'supplies', 'Office Supplies', 'Office Depot', 'card', 250.00, 25.00, 275.00, 'paid'),

    -- Coffee Corner Expenses
    ('exp00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', 'EXP-CC-001', '2025-11-01', 'rent', 'Cafe Rent', 'SF Commercial Properties', 'bank_transfer', 4000.00, 0, 4000.00, 'paid'),
    ('exp00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', 'EXP-CC-002', '2025-11-02', 'utilities', 'Water & Gas', 'SF Utilities', 'bank_transfer', 200.00, 20.00, 220.00, 'paid'),
    ('exp00002-0000-0000-0000-000000000003', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', 'EXP-CC-003', '2025-11-05', 'supplies', 'Coffee Machine Repair', 'Cafe Equipment Repair', 'cash', 350.00, 35.00, 385.00, 'paid'),

    -- Tech Gadgets Pro Expenses
    ('exp00003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'loc00003-0000-0000-0000-000000000001', 'EXP-TG-001', '2025-11-01', 'rent', 'Store Rent', 'Austin Commercial Real Estate', 'bank_transfer', 3000.00, 0, 3000.00, 'paid'),
    ('exp00003-0000-0000-0000-000000000002', '33333333-3333-3333-3333-333333333333', 'loc00003-0000-0000-0000-000000000001', 'EXP-TG-002', '2025-11-03', 'marketing', 'Google Ads Campaign', 'Google LLC', 'card', 1200.00, 120.00, 1320.00, 'paid')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SHIFTS
-- ============================================================================

INSERT INTO shifts (id, organization_id, location_id, user_id, shift_number, start_time, end_time, status, opening_cash, expected_cash, actual_cash, cash_difference, total_sales, total_transactions, closed_by, closed_at)
VALUES
    -- Demo Retail Store Shifts
    ('shf00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000003', 'SHIFT-001', '2025-11-08 09:00:00+00', '2025-11-08 17:00:00+00', 'closed', 500.00, 1980.00, 1975.00, -5.00, 3895.89, 3, '10000000-0000-0000-0000-000000000002', '2025-11-08 17:15:00+00'),
    ('shf00001-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000004', 'SHIFT-002', '2025-11-08 17:00:00+00', '2025-11-09 01:00:00+00', 'closed', 500.00, 780.00, 785.00, 5.00, 1849.96, 2, '10000000-0000-0000-0000-000000000002', '2025-11-09 01:10:00+00'),
    ('shf00001-0000-0000-0000-000000000003', '11111111-1111-1111-1111-111111111111', 'loc00001-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000003', 'SHIFT-003', '2025-11-09 09:00:00+00', NULL, 'open', 500.00, NULL, NULL, NULL, 0, 0, NULL, NULL),

    -- Coffee Corner Shifts
    ('shf00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000003', 'SHIFT-CC-001', '2025-11-08 06:00:00+00', '2025-11-08 14:00:00+00', 'closed', 200.00, 650.00, 648.50, -1.50, 775.00, 35, '20000000-0000-0000-0000-000000000002', '2025-11-08 14:15:00+00'),
    ('shf00002-0000-0000-0000-000000000002', '22222222-2222-2222-2222-222222222222', 'loc00002-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000004', 'SHIFT-CC-002', '2025-11-08 14:00:00+00', '2025-11-08 22:00:00+00', 'closed', 200.00, 420.00, 422.00, 2.00, 544.50, 28, '20000000-0000-0000-0000-000000000002', '2025-11-08 22:10:00+00'),

    -- Tech Gadgets Pro Shifts
    ('shf00003-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'loc00003-0000-0000-0000-000000000001', '30000000-0000-0000-0000-000000000003', 'SHIFT-TG-001', '2025-11-08 10:00:00+00', '2025-11-08 18:00:00+00', 'closed', 300.00, 450.00, 448.00, -2.00, 2578.99, 2, '30000000-0000-0000-0000-000000000002', '2025-11-08 18:05:00+00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- PROMOTION USAGE (Examples)
-- ============================================================================

-- Some sample promotion usage records
INSERT INTO promotion_usage (id, organization_id, promotion_id, customer_id, discount_amount, used_at)
VALUES
    ('pru00001-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'prm00001-0000-0000-0000-000000000001', 'cu000001-0000-0000-0000-000000000001', 12.69, '2025-11-08 10:30:00+00'),
    ('pru00002-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'prm00002-0000-0000-0000-000000000001', 'cu000002-0000-0000-0000-000000000001', 5.62, '2025-11-08 08:15:00+00')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Update inventory_transactions table to include purchase orders
-- ============================================================================

-- Add some inventory transactions for the received purchase orders
INSERT INTO inventory_transactions (id, organization_id, product_id, transaction_type, quantity, unit_cost, balance_after, reference_number, notes)
VALUES
    ('it000001-0000-0000-0000-000000000015', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000001', 'purchase', 15, 899.00, 40, 'PO-2025-001', 'Purchase order received'),
    ('it000001-0000-0000-0000-000000000016', '11111111-1111-1111-1111-111111111111', 'p0000001-0000-0000-0000-000000000002', 'purchase', 20, 699.00, 50, 'PO-2025-001', 'Purchase order received'),
    ('it000002-0000-0000-0000-000000000015', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000005', 'purchase', 1000, 1.20, 1050, 'PO-CC-001', 'Purchase order received'),
    ('it000002-0000-0000-0000-000000000016', '22222222-2222-2222-2222-222222222222', 'p0000002-0000-0000-0000-000000000006', 'purchase', 1000, 1.00, 1040, 'PO-CC-001', 'Purchase order received')
ON CONFLICT DO NOTHING;

COMMIT;

-- ============================================================================
-- SUMMARY
-- ============================================================================

-- Display summary of seeded data
DO $$
DECLARE
    supplier_count INTEGER;
    location_count INTEGER;
    variant_count INTEGER;
    promotion_count INTEGER;
    po_count INTEGER;
    expense_count INTEGER;
    shift_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO supplier_count FROM suppliers WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO location_count FROM locations WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO variant_count FROM product_variants WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO promotion_count FROM promotions WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO po_count FROM purchase_orders WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO expense_count FROM expenses WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO shift_count FROM shifts;

    RAISE NOTICE '=================================================';
    RAISE NOTICE 'Additional POS Data Seeded Successfully';
    RAISE NOTICE '=================================================';
    RAISE NOTICE 'Suppliers: %', supplier_count;
    RAISE NOTICE 'Locations: %', location_count;
    RAISE NOTICE 'Product Variants: %', variant_count;
    RAISE NOTICE 'Promotions: %', promotion_count;
    RAISE NOTICE 'Purchase Orders: %', po_count;
    RAISE NOTICE 'Expenses: %', expense_count;
    RAISE NOTICE 'Shifts: %', shift_count;
    RAISE NOTICE '=================================================';
END $$;
