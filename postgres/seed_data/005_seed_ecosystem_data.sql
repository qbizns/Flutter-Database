-- =====================================================
-- Seed Data: 005 - Complete Ecosystem Features
-- Description: Sample data for V009-V012 (Restaurant, Kitchen, Delivery, Staff & Devices)
-- Dependencies: Previous seed files (001-004)
-- =====================================================

\echo 'Loading ecosystem seed data (005)...'

-- =====================================================
-- V009: Restaurant & Table Management
-- =====================================================

-- Floor Plans
INSERT INTO floor_plans (id, organization_id, location_id, floor_name, floor_level, layout_config, is_active) VALUES
-- Demo Retail Store (Coffee shop setup)
('e1f1f1f1-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Demo Retail Store' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Main Warehouse' LIMIT 1),
 'Main Dining', 1,
 '{"width": 800, "height": 600, "background": "floor_plan_main.png"}', true),

-- Coffee Corner
('e1f1f1f1-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Ground Floor', 1,
 '{"width": 600, "height": 500}', true),

('e1f1f1f1-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Outdoor Patio', 1,
 '{"width": 400, "height": 400}', true);

-- Table Sections
INSERT INTO table_sections (id, organization_id, floor_plan_id, section_name, section_type, display_order) VALUES
-- Coffee Corner sections
('e2s2s2s2-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e1f1f1f1-2222-2222-2222-222222222222',
 'Window Seats', 'regular', 1),

('e2s2s2s2-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e1f1f1f1-2222-2222-2222-222222222222',
 'Cozy Corner', 'vip', 2),

('e2s2s2s2-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e1f1f1f1-3333-3333-3333-333333333333',
 'Outdoor', 'outdoor', 3);

-- Restaurant Tables
INSERT INTO restaurant_tables (id, organization_id, location_id, floor_plan_id, section_id, table_number, table_name, max_capacity, min_capacity, table_shape, position_x, position_y, status) VALUES
-- Coffee Corner tables
('e3t3t3t3-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'e1f1f1f1-2222-2222-2222-222222222222',
 'e2s2s2s2-1111-1111-1111-111111111111',
 '1', 'Table 1', 2, 1, 'round', 100, 100, 'available'),

('e3t3t3t3-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'e1f1f1f1-2222-2222-2222-222222222222',
 'e2s2s2s2-1111-1111-1111-111111111111',
 '2', 'Table 2', 4, 2, 'rectangle', 300, 100, 'occupied'),

('e3t3t3t3-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'e1f1f1f1-2222-2222-2222-222222222222',
 'e2s2s2s2-2222-2222-2222-222222222222',
 '10', 'VIP 1', 6, 4, 'rectangle', 100, 300, 'available'),

('e3t3t3t3-4444-4444-4444-444444444444',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'e1f1f1f1-3333-3333-3333-333333333333',
 'e2s2s2s2-3333-3333-3333-333333333333',
 'P1', 'Patio 1', 4, 2, 'square', 50, 50, 'available');

-- Reservations
INSERT INTO reservations (id, organization_id, location_id, customer_id, table_id, reservation_number, reservation_date, reservation_time, party_size, status, special_requests, contact_name, contact_phone) VALUES
('e4r4r4r4-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM customers WHERE email = 'john.smith@email.com' LIMIT 1),
 'e3t3t3t3-3333-3333-3333-333333333333',
 'RES-20251109-001', CURRENT_DATE + 1, '18:00:00', 6, 'confirmed',
 'Birthday celebration, need birthday cake', 'John Smith', '+201234567890'),

('e4r4r4r4-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM customers WHERE email = 'sarah.j@email.com' LIMIT 1),
 'e3t3t3t3-4444-4444-4444-444444444444',
 'RES-20251109-002', CURRENT_DATE, '12:30:00', 4, 'seated',
 'Prefer outdoor seating', 'Sarah Johnson', '+201234567891');

-- Modifier Groups
INSERT INTO modifier_groups (id, organization_id, group_name, selection_type, min_selections, max_selections, is_required, display_order) VALUES
('e5m5m5m5-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Coffee Size', 'single', 1, 1, true, 1),

('e5m5m5m5-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Milk Options', 'single', 0, 1, false, 2),

('e5m5m5m5-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Extra Toppings', 'multiple', 0, 5, false, 3),

('e5m5m5m5-4444-4444-4444-444444444444',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Sandwich Extras', 'multiple', 0, 10, false, 4);

-- Modifiers
INSERT INTO modifiers (id, organization_id, modifier_group_id, modifier_name, price_adjustment, price_type, is_default, display_order) VALUES
-- Coffee sizes
('e6m6m6m6-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-1111-1111-1111-111111111111',
 'Small', 0.00, 'add', true, 1),
('e6m6m6m6-1112-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-1111-1111-1111-111111111111',
 'Medium', 5.00, 'add', false, 2),
('e6m6m6m6-1113-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-1111-1111-1111-111111111111',
 'Large', 10.00, 'add', false, 3),

-- Milk options
('e6m6m6m6-2221-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-2222-2222-2222-222222222222',
 'Whole Milk', 0.00, 'add', true, 1),
('e6m6m6m6-2222-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-2222-2222-2222-222222222222',
 'Almond Milk', 5.00, 'add', false, 2),
('e6m6m6m6-2223-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-2222-2222-2222-222222222222',
 'Oat Milk', 5.00, 'add', false, 3),

-- Toppings
('e6m6m6m6-3331-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-3333-3333-3333-333333333333',
 'Extra Shot Espresso', 8.00, 'add', false, 1),
('e6m6m6m6-3332-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-3333-3333-3333-333333333333',
 'Whipped Cream', 3.00, 'add', false, 2),
('e6m6m6m6-3333-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'e5m5m5m5-3333-3333-3333-333333333333',
 'Caramel Drizzle', 3.00, 'add', false, 3);

-- Product Modifier Groups (link products to modifier groups)
INSERT INTO product_modifier_groups (id, organization_id, product_id, modifier_group_id, is_required) VALUES
-- Link coffee products to modifiers (assuming product exists from previous seed data)
('e7p7p7p7-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM products WHERE product_name LIKE '%Coffee%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'e5m5m5m5-1111-1111-1111-111111111111', true),

('e7p7p7p7-1112-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM products WHERE product_name LIKE '%Coffee%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'e5m5m5m5-2222-2222-2222-222222222222', false),

('e7p7p7p7-1113-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM products WHERE product_name LIKE '%Coffee%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'e5m5m5m5-3333-3333-3333-333333333333', false);

-- Courses
INSERT INTO courses (id, organization_id, course_name, course_type, display_order, typical_duration_minutes, fire_delay_minutes) VALUES
('e8c8c8c8-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Beverages', 'beverage', 1, 5, 0),

('e8c8c8c8-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Appetizers', 'appetizer', 2, 10, 0),

('e8c8c8c8-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Main Course', 'main', 3, 20, 10),

('e8c8c8c8-4444-4444-4444-444444444444',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'Desserts', 'dessert', 4, 5, 20);

-- =====================================================
-- V010: Kitchen Operations
-- =====================================================

-- Kitchen Stations
INSERT INTO kitchen_stations (id, organization_id, location_id, station_name, station_code, station_type, color_code, is_active, display_order) VALUES
('e9k9k9k9-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Barista Station', 'BARISTA', 'bar', '#8B4513', true, 1),

('e9k9k9k9-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Sandwich & Grill', 'GRILL', 'kitchen', '#FF4500', true, 2),

('e9k9k9k9-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Dessert Station', 'DESSERT', 'dessert', '#FFB6C1', true, 3);

-- Orders (in-progress restaurant orders)
INSERT INTO orders (id, organization_id, location_id, order_number, display_number, order_type, table_id, customer_id, waiter_id, status, order_date, submitted_at, covers, subtotal, tax_amount, total_amount, customer_notes) VALUES
('f1o1o1o1-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'ORD-20251109-001', 1, 'dine_in',
 'e3t3t3t3-2222-2222-2222-222222222222',
 (SELECT id FROM customers WHERE email = 'john.smith@email.com' LIMIT 1),
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 'preparing', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '14 minutes', 4,
 120.00, 12.00, 132.00, 'No sugar in coffee'),

('f1o1o1o1-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'ORD-20251109-002', 2, 'takeout',
 NULL,
 (SELECT id FROM customers WHERE email = 'sarah.j@email.com' LIMIT 1),
 (SELECT id FROM users WHERE email = 'cashier1@demoretail.com' LIMIT 1),
 'ready', NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '9 minutes', 1,
 85.00, 8.50, 93.50, NULL);

-- Order Items
INSERT INTO order_items (id, organization_id, order_id, product_id, item_name, quantity, unit_price, course_id, kitchen_station_id, status, line_total, fired_at, started_preparing_at) VALUES
-- Order 1 items
('f2i2i2i2-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'f1o1o1o1-1111-1111-1111-111111111111',
 (SELECT id FROM products WHERE product_name LIKE '%Coffee%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'Cappuccino', 2, 30.00, 'e8c8c8c8-1111-1111-1111-111111111111',
 'e9k9k9k9-1111-1111-1111-111111111111', 'ready', 60.00,
 NOW() - INTERVAL '14 minutes', NOW() - INTERVAL '13 minutes'),

('f2i2i2i2-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'f1o1o1o1-1111-1111-1111-111111111111',
 (SELECT id FROM products WHERE product_name LIKE '%Sandwich%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'Club Sandwich', 2, 30.00, 'e8c8c8c8-3333-3333-3333-333333333333',
 'e9k9k9k9-2222-2222-2222-222222222222', 'preparing', 60.00,
 NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '9 minutes'),

-- Order 2 items
('f2i2i2i2-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'f1o1o1o1-2222-2222-2222-222222222222',
 (SELECT id FROM products WHERE product_name LIKE '%Coffee%' AND organization_id = (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1) LIMIT 1),
 'Latte', 3, 25.00, 'e8c8c8c8-1111-1111-1111-111111111111',
 'e9k9k9k9-1111-1111-1111-111111111111', 'ready', 75.00,
 NOW() - INTERVAL '9 minutes', NOW() - INTERVAL '8 minutes');

-- Kitchen Tickets
INSERT INTO kitchen_tickets (id, organization_id, location_id, order_id, kitchen_station_id, ticket_number, display_sequence, status, table_number, order_type, covers, waiter_name, created_at, fired_at) VALUES
('f3k3k3k3-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'f1o1o1o1-1111-1111-1111-111111111111',
 'e9k9k9k9-1111-1111-1111-111111111111',
 'TKT-001', 1, 'preparing', '2', 'dine_in', 4, 'Manager Demo',
 NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '14 minutes'),

('f3k3k3k3-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'f1o1o1o1-1111-1111-1111-111111111111',
 'e9k9k9k9-2222-2222-2222-222222222222',
 'TKT-002', 2, 'preparing', '2', 'dine_in', 4, 'Manager Demo',
 NOW() - INTERVAL '10 minutes', NOW() - INTERVAL '10 minutes');

-- Update order items with kitchen ticket references
UPDATE order_items SET kitchen_ticket_id = 'f3k3k3k3-1111-1111-1111-111111111111'
WHERE id = 'f2i2i2i2-1111-1111-1111-111111111111';

UPDATE order_items SET kitchen_ticket_id = 'f3k3k3k3-2222-2222-2222-222222222222'
WHERE id = 'f2i2i2i2-2222-2222-2222-222222222222';

-- =====================================================
-- V011: Delivery & Online Ordering
-- =====================================================

-- Delivery Zones
INSERT INTO delivery_zones (id, organization_id, location_id, zone_name, zone_code, base_delivery_fee, minimum_order_amount, estimated_delivery_time_minutes, is_active) VALUES
('g1z1z1z1-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Downtown Zone 1', 'DT1', 15.00, 50.00, 20, true),

('g1z1z1z1-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Downtown Zone 2', 'DT2', 25.00, 50.00, 30, true),

('g1z1z1z1-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Suburbs', 'SUB1', 40.00, 100.00, 45, true);

-- Delivery Drivers
INSERT INTO delivery_drivers (id, organization_id, driver_code, full_name, phone, email, vehicle_type, vehicle_make, vehicle_model, vehicle_color, license_plate, status, total_deliveries, successful_deliveries, rating, rating_count, is_available, hire_date) VALUES
('g2d2d2d2-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'DRV-001', 'Ahmed Mohamed', '+201001234567', 'ahmed.driver@email.com',
 'motorcycle', 'Honda', 'CB150R', 'Red', 'ABC 1234', 'active',
 127, 125, 4.75, 89, true, '2024-01-15'),

('g2d2d2d2-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'DRV-002', 'Mohamed Ali', '+201007654321', 'mohamed.driver@email.com',
 'motorcycle', 'Yamaha', 'YZF-R15', 'Blue', 'XYZ 5678', 'active',
 93, 91, 4.85, 67, false, '2024-03-20'),

('g2d2d2d2-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'DRV-003', 'Hassan Ibrahim', '+201009876543', 'hassan.driver@email.com',
 'car', 'Hyundai', 'Accent', 'White', 'DEF 9012', 'active',
 45, 44, 4.60, 35, true, '2024-06-10');

-- Driver Shifts
INSERT INTO driver_shifts (id, organization_id, location_id, driver_id, shift_date, scheduled_start_time, scheduled_end_time, actual_start_time, status, total_deliveries, total_distance_km, total_earnings) VALUES
('g3s3s3s3-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'g2d2d2d2-1111-1111-1111-111111111111',
 CURRENT_DATE, '08:00:00', '16:00:00', CURRENT_TIMESTAMP - INTERVAL '6 hours',
 'started', 8, 45.5, 120.00),

('g3s3s3s3-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'g2d2d2d2-3333-3333-3333-333333333333',
 CURRENT_DATE, '10:00:00', '18:00:00', CURRENT_TIMESTAMP - INTERVAL '4 hours',
 'started', 5, 32.0, 85.00);

-- Customer Addresses
INSERT INTO customer_addresses (id, organization_id, customer_id, address_label, address_line1, address_line2, city, postal_code, delivery_zone_id, is_default, location_notes) VALUES
('g4a4a4a4-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM customers WHERE email = 'john.smith@email.com' LIMIT 1),
 'Home', '15 El Tahrir Street', 'Apartment 5A', 'Cairo', '11511',
 'g1z1z1z1-1111-1111-1111-111111111111', true,
 'Ring bell twice, building has elevator'),

('g4a4a4a4-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM customers WHERE email = 'sarah.j@email.com' LIMIT 1),
 'Office', '45 Zamalek Street', 'Office 302', 'Cairo', '11211',
 'g1z1z1z1-2222-2222-2222-222222222222', true,
 'Call on arrival, reception will direct'),

('g4a4a4a4-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM customers WHERE email = 'john.smith@email.com' LIMIT 1),
 'Office', '20 Mohandessin Street', 'Floor 8', 'Cairo', '12411',
 'g1z1z1z1-2222-2222-2222-222222222222', false,
 'Ask security for visitor pass');

-- Delivery Orders (add delivery-specific orders)
INSERT INTO orders (id, organization_id, location_id, order_number, display_number, order_type, customer_id, delivery_zone_id, customer_address_id, delivery_address, delivery_fee, status, order_date, submitted_at, covers, subtotal, tax_amount, total_amount, source) VALUES
('g5o5o5o5-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'DEL-20251109-001', 101, 'delivery',
 (SELECT id FROM customers WHERE email = 'john.smith@email.com' LIMIT 1),
 'g1z1z1z1-1111-1111-1111-111111111111',
 'g4a4a4a4-1111-1111-1111-111111111111',
 '15 El Tahrir Street, Apartment 5A, Cairo', 15.00,
 'preparing', NOW() - INTERVAL '20 minutes', NOW() - INTERVAL '19 minutes', 1,
 100.00, 10.00, 125.00, 'mobile_app');

-- Delivery Assignments
INSERT INTO delivery_assignments (id, organization_id, order_id, driver_id, driver_shift_id, delivery_zone_id, customer_address_id, delivery_address, status, assigned_at, accepted_at, distance_km, delivery_fee, driver_commission, payment_method) VALUES
('g6d6d6d6-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'g5o5o5o5-1111-1111-1111-111111111111',
 'g2d2d2d2-1111-1111-1111-111111111111',
 'g3s3s3s3-1111-1111-1111-111111111111',
 'g1z1z1z1-1111-1111-1111-111111111111',
 'g4a4a4a4-1111-1111-1111-111111111111',
 '15 El Tahrir Street, Apartment 5A, Cairo',
 'in_transit', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '14 minutes',
 5.2, 15.00, 10.00, 'cash');

-- Order Tracking Events
INSERT INTO order_tracking_events (organization_id, order_id, delivery_assignment_id, event_type, event_timestamp, event_message, actor_type) VALUES
((SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'g5o5o5o5-1111-1111-1111-111111111111', NULL,
 'order_placed', NOW() - INTERVAL '20 minutes',
 'Your order has been placed successfully', 'customer'),

((SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'g5o5o5o5-1111-1111-1111-111111111111', NULL,
 'order_confirmed', NOW() - INTERVAL '19 minutes',
 'Restaurant has confirmed your order', 'system'),

((SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'g5o5o5o5-1111-1111-1111-111111111111',
 'g6d6d6d6-1111-1111-1111-111111111111',
 'picked_up', NOW() - INTERVAL '10 minutes',
 'Driver has picked up your order', 'driver'),

((SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 'g5o5o5o5-1111-1111-1111-111111111111',
 'g6d6d6d6-1111-1111-1111-111111111111',
 'out_for_delivery', NOW() - INTERVAL '8 minutes',
 'Your order is on the way', 'system');

-- =====================================================
-- V012: Staff & Device Management
-- =====================================================

-- Employee Schedules
INSERT INTO employee_schedules (id, organization_id, location_id, employee_id, schedule_date, shift_type, position, scheduled_start_time, scheduled_end_time, break_duration_minutes, status) VALUES
-- Today's schedules
('h1s1s1s1-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 CURRENT_DATE, 'regular', 'Manager', '08:00:00', '17:00:00', 60, 'started'),

('h1s1s1s1-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'cashier1@demoretail.com' LIMIT 1),
 CURRENT_DATE, 'regular', 'Cashier/Waiter', '09:00:00', '18:00:00', 60, 'started'),

-- Tomorrow's schedules
('h1s1s1s1-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 CURRENT_DATE + 1, 'regular', 'Manager', '08:00:00', '17:00:00', 60, 'scheduled');

-- Time Clock Entries
INSERT INTO time_clock_entries (id, organization_id, location_id, employee_id, schedule_id, entry_type, entry_timestamp, scheduled_timestamp, is_late, variance_minutes) VALUES
-- Manager clock in
('h2t2t2t2-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 'h1s1s1s1-1111-1111-1111-111111111111',
 'clock_in', CURRENT_DATE + TIME '07:58:00',
 CURRENT_DATE + TIME '08:00:00', false, -2),

-- Cashier clock in (late)
('h2t2t2t2-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'cashier1@demoretail.com' LIMIT 1),
 'h1s1s1s1-2222-2222-2222-222222222222',
 'clock_in', CURRENT_DATE + TIME '09:12:00',
 CURRENT_DATE + TIME '09:00:00', true, 12);

-- Devices
INSERT INTO devices (id, organization_id, location_id, device_code, device_name, device_type, manufacturer, model, serial_number, status, connection_type, ip_address, last_online_at, last_heartbeat_at) VALUES
-- POS Terminals
('h3d3d3d3-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'POS-01', 'Front Counter Terminal', 'pos_terminal', 'Sunmi', 'T2',
 'SN-POS-001-T2', 'active', 'network', '192.168.1.101',
 NOW(), NOW() - INTERVAL '30 seconds'),

('h3d3d3d3-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'TAB-01', 'Waiter Tablet 1', 'tablet', 'Samsung', 'Galaxy Tab A',
 'SN-TAB-001-GTA', 'active', 'wifi', '192.168.1.201',
 NOW(), NOW() - INTERVAL '1 minute'),

-- Kitchen Display
('h3d3d3d3-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'KDS-01', 'Kitchen Display 1', 'kds_display', 'Generic', 'TouchScreen 22"',
 'SN-KDS-001-TS22', 'active', 'network', '192.168.1.151',
 NOW(), NOW() - INTERVAL '45 seconds'),

-- Printers
('h3d3d3d3-4444-4444-4444-444444444444',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'PRT-RCPT-01', 'Receipt Printer - Counter', 'receipt_printer', 'Epson', 'TM-T88VI',
 'SN-PRT-001-T88', 'active', 'network', '192.168.1.111',
 NOW(), NOW() - INTERVAL '20 seconds'),

('h3d3d3d3-5555-5555-5555-555555555555',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'PRT-KIT-01', 'Kitchen Printer', 'kitchen_printer', 'Epson', 'TM-U220',
 'SN-PRT-002-U220', 'active', 'network', '192.168.1.112',
 NOW(), NOW() - INTERVAL '15 seconds'),

-- Handheld Scanner
('h3d3d3d3-6666-6666-6666-666666666666',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'SCAN-01', 'Inventory Scanner', 'handheld_scanner', 'Zebra', 'MC3300',
 'SN-SCAN-001-MC33', 'active', 'wifi', '192.168.1.180',
 NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours');

-- Update time clock entries with device references
UPDATE time_clock_entries SET device_id = 'h3d3d3d3-1111-1111-1111-111111111111'
WHERE id IN ('h2t2t2t2-1111-1111-1111-111111111111', 'h2t2t2t2-2222-2222-2222-222222222222');

-- Printer Configurations
INSERT INTO printer_configurations (id, organization_id, location_id, printer_device_id, document_type, filter_kitchen_station_id, number_of_copies, auto_print, is_active, paper_size) VALUES
-- Receipt printer for customer receipts
('h4p4p4p4-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'h3d3d3d3-4444-4444-4444-444444444444',
 'receipt', NULL, 1, true, true, 'thermal_80mm'),

-- Kitchen printer for barista station
('h4p4p4p4-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'h3d3d3d3-5555-5555-5555-555555555555',
 'kitchen_ticket', 'e9k9k9k9-1111-1111-1111-111111111111',
 1, true, true, 'thermal_80mm'),

-- Kitchen printer for grill station
('h4p4p4p4-3333-3333-3333-333333333333',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'h3d3d3d3-5555-5555-5555-555555555555',
 'kitchen_ticket', 'e9k9k9k9-2222-2222-2222-222222222222',
 1, true, true, 'thermal_80mm');

-- Tip Pools
INSERT INTO tip_pools (id, organization_id, location_id, pool_name, pool_type, distribution_method, eligible_positions, is_active) VALUES
('h5t5t5t5-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Daily Service Staff Pool', 'daily', 'hours_worked',
 ARRAY['Waiter', 'Cashier', 'Barista'], true),

('h5t5t5t5-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'Kitchen Staff Pool', 'weekly', 'equal',
 ARRAY['Chef', 'Cook', 'Kitchen Helper'], true);

-- Tip Distributions
INSERT INTO tip_distributions (id, organization_id, location_id, tip_pool_id, distribution_date, employee_id, source_type, tip_amount, distribution_amount, distribution_percentage, payment_status) VALUES
-- Today's tips for cashier
('h6t6t6t6-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'h5t5t5t5-1111-1111-1111-111111111111',
 CURRENT_DATE,
 (SELECT id FROM users WHERE email = 'cashier1@demoretail.com' LIMIT 1),
 'pool', 150.00, 75.00, 50.00, 'pending'),

-- Today's tips for manager
('h6t6t6t6-2222-2222-2222-222222222222',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 'h5t5t5t5-1111-1111-1111-111111111111',
 CURRENT_DATE,
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 'pool', 150.00, 75.00, 50.00, 'pending');

-- Staff Commissions
INSERT INTO staff_commissions (id, organization_id, location_id, employee_id, commission_date, period_start, period_end, source_type, commission_type, commission_rate, sales_amount, commission_amount, status) VALUES
-- Weekly commission for manager
('h7c7c7c7-1111-1111-1111-111111111111',
 (SELECT id FROM organizations WHERE organization_name = 'Coffee Corner' LIMIT 1),
 (SELECT id FROM locations WHERE location_name = 'Downtown Location' LIMIT 1),
 (SELECT id FROM users WHERE email = 'manager@demoretail.com' LIMIT 1),
 CURRENT_DATE, CURRENT_DATE - 6, CURRENT_DATE,
 'target', 'percentage', 2.00, 25000.00, 500.00, 'approved');

\echo 'Ecosystem seed data (005) loaded successfully!'
\echo 'Added: Floor plans, tables, reservations, modifiers, courses, kitchen stations,'
\echo '       orders, kitchen tickets, delivery zones, drivers, customer addresses,'
\echo '       delivery assignments, devices, printers, schedules, time clock, tips, and commissions'
