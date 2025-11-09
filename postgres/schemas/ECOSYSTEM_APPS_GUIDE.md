# POS Ecosystem Apps - Database Integration Guide

Complete guide showing how each of the 16 applications in the POS ecosystem interacts with the PostgreSQL database.

---

## Table of Contents

1. [Backend & Admin Web Panel](#1-backend--admin-web-panel)
2. [Main POS/Register](#2-main-posregister)
3. [Waiter App (Tablet)](#3-waiter-app-tablet)
4. [Kitchen Display System (KDS)](#4-kitchen-display-system-kds)
5. [Customer Display](#5-customer-display)
6. [Self-Service Kiosk](#6-self-service-kiosk)
7. [Table/QR Code Ordering](#7-tableqr-code-ordering)
8. [Retail POS](#8-retail-pos)
9. [Stock & Inventory Handheld](#9-stock--inventory-handheld)
10. [Price Checker](#10-price-checker)
11. [Delivery Driver App](#11-delivery-driver-app)
12. [Owner Analytics App](#12-owner-analytics-app)
13. [Franchise/HQ Control Panel](#13-franchisehq-control-panel)
14. [Staff Time & Attendance](#14-staff-time--attendance)
15. [Maintenance/Setup App](#15-maintenancesetup-app)
16. [Customer Mobile App](#16-customer-mobile-app)

---

## 1. Backend & Admin Web Panel

**Purpose**: Centralized administration, configuration, and reporting for the entire organization.

### Key Tables Used:

#### Core Administration
- **organizations** - Manage tenant settings, subscriptions, quotas
- **users** - Create and manage system users
- **roles** & **permissions** - Configure RBAC
- **user_roles** - Assign roles to users
- **audit_logs** - View complete system audit trail
- **locations** - Manage store locations and branches

#### Product & Inventory Management
- **categories** - Organize product hierarchy
- **products** - Master product catalog
- **product_variants** - Manage product variations
- **suppliers** - Vendor management
- **purchase_orders** & **purchase_order_items** - Procurement workflow

#### Restaurant Configuration
- **floor_plans** - Design restaurant layouts
- **table_sections** - Organize dining areas
- **restaurant_tables** - Configure table settings
- **kitchen_stations** - Setup kitchen workflows
- **modifier_groups** & **modifiers** - Menu customization options
- **courses** - Define course timing

#### Delivery Setup
- **delivery_zones** - Configure delivery areas and fees
- **delivery_drivers** - Manage driver profiles

#### Staff Management
- **employee_schedules** - Create work schedules
- **tip_pools** - Configure tip distribution
- **devices** - Register and manage hardware
- **printer_configurations** - Setup printer routing

#### Reporting & Analytics
- **mv_daily_sales_summary** - Sales performance
- **mv_product_performance** - Best sellers
- **mv_customer_analytics** - Customer insights
- **mv_inventory_valuation** - Stock value
- **mv_location_performance** - Store comparisons
- **mv_promotion_effectiveness** - Campaign ROI

### Key Operations:
- CRUD operations on all master data
- Bulk imports/exports
- Report generation and scheduling
- User permission management
- System configuration

---

## 2. Main POS/Register

**Purpose**: Traditional point-of-sale terminal for retail and quick service transactions.

### Key Tables Used:

#### Transaction Processing
- **sales** - Complete sale transactions
- **sale_items** - Line items with products
- **payments** - Payment processing (cash, card, etc.)
- **customers** - Customer lookup and creation
- **shifts** - Cashier shift management

#### Product Lookup
- **products** - Product search and pricing
- **product_variants** - Variant selection
- **categories** - Product browsing
- **promotions** & **promotion_usage** - Apply discounts

#### Inventory
- **inventory_transactions** - Auto-update stock on sale
- **product_serial_numbers** - Serial number capture
- **product_batches** - Batch/lot selection

#### Restaurant Mode
- **orders** - Create dine-in/takeout orders
- **order_items** - Order line items
- **restaurant_tables** - Table selection for dine-in
- **order_item_modifiers** - Apply customizations

#### Employee & Device
- **users** - Cashier authentication
- **time_clock_entries** - Clock in/out
- **devices** - Terminal registration
- **printer_configurations** - Print receipts

### Typical Workflow:
1. Cashier clocks in via `time_clock_entries`
2. Open shift in `shifts` table
3. Add items to cart (query `products`, `product_variants`)
4. Apply promotions from `promotions`
5. Process payment via `payments`
6. Create sale record in `sales` + `sale_items`
7. Update inventory via `inventory_transactions`
8. Print receipt via configured printer
9. Close shift with cash reconciliation

---

## 3. Waiter App (Tablet)

**Purpose**: Mobile ordering app for restaurant waitstaff.

### Key Tables Used:

#### Order Management
- **orders** - Create and manage orders
- **order_items** - Add menu items
- **order_item_modifiers** - Customer customizations
- **restaurant_tables** - Table management
- **reservations** - Check reservations

#### Menu & Modifiers
- **products** - Browse menu
- **categories** - Menu sections
- **modifier_groups** & **modifiers** - Customization options
- **product_modifier_groups** - Available modifiers per product
- **courses** - Course selection for kitchen timing

#### Customer Info
- **customers** - Link orders to customers
- **loyalty_points_transactions** - Award loyalty points

#### Kitchen Communication
- **kitchen_tickets** - Send orders to kitchen
- **kitchen_stations** - Route items to stations

#### Staff
- **users** - Waiter authentication (waiter_id on orders)
- **tip_distributions** - Track tips

### Typical Workflow:
1. Waiter logs in via `users`
2. Select table from `restaurant_tables`
3. Create new order in `orders` (status: draft)
4. Add items to `order_items` with modifiers from `order_item_modifiers`
5. Submit order (status: submitted → sent_to_kitchen)
6. System creates `kitchen_tickets` for each station
7. Monitor order status updates from kitchen
8. Mark items as served in `order_items`
9. Process payment and close order

---

## 4. Kitchen Display System (KDS)

**Purpose**: Real-time order display for kitchen staff, organized by station.

### Key Tables Used:

#### Ticket Display
- **kitchen_tickets** - Active tickets to display
- **order_items** - Items on each ticket
- **order_item_modifiers** - Special instructions
- **orders** - Order context (table, covers, timing)

#### Station Routing
- **kitchen_stations** - Filter tickets by station
- **courses** - Course-based ordering

#### Status Management
- **order_items** - Update item status (pending → fired → preparing → ready)
- **orders** - Update order status based on items

#### Configuration
- **devices** - KDS device registration
- **printer_configurations** - Auto-print tickets

### Typical Workflow:
1. KDS device queries `v_active_kitchen_tickets` view filtered by `kitchen_station_id`
2. Display tickets ordered by priority and fire time
3. Kitchen staff acknowledges ticket (update `kitchen_tickets.status`)
4. Update `order_items.status` as items are prepared
5. Mark ticket ready when all items complete
6. Bump ticket (remove from screen) when served

### View Used:
```sql
SELECT * FROM v_active_kitchen_tickets
WHERE kitchen_station_id = :station_id
ORDER BY priority DESC, fired_at ASC;
```

---

## 5. Customer Display

**Purpose**: Facing display showing customer's cart and total during checkout.

### Key Tables Used:

#### Current Transaction
- **sale_items** - Items being rung up
- **sales** - Running total
- **promotions** - Applied discounts
- **products** - Product names and images

#### Loyalty
- **customers** - Customer name if identified
- **loyalty_points_transactions** - Points to be earned

### Typical Workflow:
1. Display updates in real-time as cashier scans items
2. Query current `sale_items` for active sale
3. Show product images from `products.image_url`
4. Display total from `sales.total_amount`
5. Show loyalty points to be earned
6. Thank you message after payment

---

## 6. Self-Service Kiosk

**Purpose**: Customer-facing kiosk for self-ordering (restaurant or retail).

### Key Tables Used:

#### Product Browsing
- **categories** - Browse menu by category
- **products** - View items with images and descriptions
- **product_variants** - Select size/color
- **modifier_groups** & **modifiers** - Customizations
- **promotions** - Display available deals

#### Order Creation
- **orders** - Create customer order
- **order_items** - Add items to order
- **order_item_modifiers** - Apply customizations

#### Payment
- **payments** - Process payment (card reader integration)
- **customers** - Optional customer identification

#### Order Type Selection
- **order_type** field: 'dine_in', 'takeout', 'delivery'
- **restaurant_tables** - Optional table selection for dine-in
- **delivery_zones** & **customer_addresses** - For delivery orders

### Typical Workflow:
1. Customer selects order type (dine-in/takeout/delivery)
2. Browse `products` by `categories`
3. Select modifiers from available `modifier_groups`
4. Add items to cart (create `order` with status: draft)
5. Review order
6. Process payment via `payments`
7. Submit order (status: submitted)
8. Display order number and estimated time
9. Send to kitchen via `kitchen_tickets`

---

## 7. Table/QR Code Ordering

**Purpose**: Customers scan QR code at table to order from their phone.

### Key Tables Used:

#### Table Context
- **restaurant_tables** - Identify table from QR code
- **floor_plans** & **table_sections** - Table location info

#### Menu & Ordering
- **categories** & **products** - Browse menu
- **modifier_groups** & **modifiers** - Customizations
- **orders** - Create order linked to table
- **order_items** - Add items
- **courses** - Course selection

#### Customer
- **customers** - Optional customer account
- **loyalty_points_transactions** - Earn points

#### Payment
- **payments** - Online payment processing
- **order** status: 'submitted' after payment

### Typical Workflow:
1. Customer scans QR code (contains table_id)
2. Query `restaurant_tables` to get table info
3. Create new `order` with table_id
4. Customer browses menu and adds items
5. Customer pays online via `payments`
6. Order status updated to 'submitted'
7. Kitchen receives order via `kitchen_tickets`
8. Waiter delivers food to table

---

## 8. Retail POS

**Purpose**: Standard retail checkout with barcode scanning (non-restaurant).

### Key Tables Used:

#### Product Lookup
- **products** - Barcode/SKU lookup
- **product_variants** - Variant-specific barcodes
- **categories** - Product categorization

#### Transaction
- **sales** - Sale header
- **sale_items** - Scanned items
- **payments** - Payment processing
- **customers** - Customer lookup

#### Inventory
- **inventory_transactions** - Stock deduction
- **product_serial_numbers** - Track serialized items
- **product_batches** - Batch selection for expiry tracking
- **locations** - Multi-location inventory

#### Promotions
- **promotions** - Auto-apply eligible promotions
- **promotion_usage** - Track usage limits

### Typical Workflow:
1. Scan barcode → lookup in `products` or `product_variants`
2. Add to `sale_items`
3. Apply automatic promotions from `promotions`
4. Tender payment → create `payments`
5. Complete sale → insert `sales`
6. Deduct inventory → `inventory_transactions`
7. Print receipt

---

## 9. Stock & Inventory Handheld

**Purpose**: Mobile device for stock counting, receiving, and transfers.

### Key Tables Used:

#### Inventory Management
- **products** - Product lookup
- **product_variants** - Variant identification
- **inventory_transactions** - Record all movements
- **locations** - Multi-location support

#### Receiving
- **purchase_orders** & **purchase_order_items** - PO receiving
- **product_batches** - Create batches on receive
- **product_serial_numbers** - Register serial numbers

#### Transfers
- **inventory_transfers** & **inventory_transfer_items** - Inter-location transfers
- **batch_transactions** - Track batch movements

#### Counting
- **cycle_counts** & **cycle_count_items** - Physical inventory counts
- **stock_adjustment_reasons** - Reason codes for adjustments

### Typical Workflow:

**Receiving:**
1. Scan PO barcode → query `purchase_orders`
2. Scan product barcodes
3. Enter quantity received
4. Create `product_batches` if applicable
5. Update `purchase_order_items.quantity_received`
6. Create `inventory_transactions` (type: purchase)

**Cycle Count:**
1. Start new `cycle_count`
2. Scan products and enter counted qty in `cycle_count_items`
3. System calculates variance vs system qty
4. Review variances
5. Approve count → create `inventory_transactions` for adjustments

**Transfer:**
1. Create `inventory_transfer` (from_location, to_location)
2. Scan items to transfer
3. Submit transfer (status: approved)
4. At destination, receive transfer
5. Create `inventory_transactions` at both locations

---

## 10. Price Checker

**Purpose**: Customer-facing device to check product prices and information.

### Key Tables Used:

#### Product Information
- **products** - Price, description, image
- **product_variants** - Variant-specific pricing
- **categories** - Product classification
- **promotions** - Active promotions

#### Inventory
- **inventory_transactions** - Stock availability (aggregated)
- **locations** - Check stock at current location

### Typical Workflow:
1. Customer scans barcode
2. Query `products` or `product_variants` by barcode
3. Display:
   - Product name and image
   - Current price
   - Active promotions
   - Stock availability (in stock / out of stock)
4. Show related products from same `category`

---

## 11. Delivery Driver App

**Purpose**: Mobile app for delivery drivers to manage deliveries.

### Key Tables Used:

#### Driver Profile
- **delivery_drivers** - Driver profile and stats
- **driver_shifts** - Clock in/out for shift

#### Deliveries
- **delivery_assignments** - Assigned deliveries
- **orders** - Order details
- **customer_addresses** - Delivery address and notes
- **delivery_zones** - Zone information

#### Status Updates
- **delivery_assignments** - Update status (accepted → picked_up → in_transit → arrived → delivered)
- **order_tracking_events** - Create tracking events for customer

#### Navigation
- **customer_addresses** (latitude/longitude) - GPS navigation
- **delivery_assignments.route_info** - Optimized route

### Typical Workflow:
1. Driver logs in → query `delivery_drivers`
2. Start shift → insert `driver_shifts`
3. View assigned deliveries from `delivery_assignments`
4. Accept delivery (status: accepted)
5. Navigate to restaurant (pickup location)
6. Mark as picked up (status: picked_up)
7. Navigate to customer (GPS from `customer_addresses`)
8. Update location periodically
9. Arrive at customer (status: arrived)
10. Complete delivery:
    - Capture signature/photo
    - Update `delivery_assignments.status` = 'delivered'
    - Collect cash if COD
11. System automatically creates `order_tracking_events`

---

## 12. Owner Analytics App

**Purpose**: Mobile/web dashboard for business owners to monitor performance.

### Key Tables Used:

#### Sales Analytics
- **mv_daily_sales_summary** - Daily sales trends
- **mv_location_performance** - Compare store performance
- **sales** - Detailed transaction history
- **sale_items** - Item-level analysis

#### Product Performance
- **mv_product_performance** - Best/worst sellers
- **products** - Product profitability
- **categories** - Category performance

#### Customer Insights
- **mv_customer_analytics** - Customer segmentation (RFM)
- **customers** - Customer database
- **loyalty_points_transactions** - Loyalty engagement

#### Inventory
- **mv_inventory_valuation** - Stock value
- **inventory_transactions** - Movement history

#### Promotions
- **mv_promotion_effectiveness** - Campaign ROI
- **promotions** & **promotion_usage** - Promotion analysis

#### Staff & Operations
- **employee_schedules** - Labor scheduling
- **time_clock_entries** - Actual hours worked
- **shifts** - Cash reconciliation
- **tip_distributions** & **staff_commissions** - Staff compensation

### Key Metrics:
- Total sales (daily, weekly, monthly)
- Sales by location, product, category
- Average transaction value
- Customer count and retention
- Inventory turnover
- Labor cost percentage
- Promotion effectiveness
- Top performing staff

---

## 13. Franchise/HQ Control Panel

**Purpose**: Multi-tenant dashboard for franchise owners managing multiple locations.

### Key Tables Used:

#### Organization Management
- **organizations** - Manage franchises/brands
- **locations** - All store locations
- **users** - Staff across all locations

#### Cross-Location Analytics
- **mv_location_performance** - Compare locations
- **mv_daily_sales_summary** - Aggregated sales
- **mv_product_performance** - System-wide product analysis
- **mv_inventory_valuation** - Total inventory value

#### Procurement
- **suppliers** - Centralized vendor management
- **purchase_orders** - PO tracking across locations

#### Operations
- **inventory_transfers** - Inter-location transfers
- **delivery_zones** - Coverage maps
- **delivery_drivers** - Fleet management

#### Standardization
- **products** & **categories** - Master catalog
- **modifier_groups** & **modifiers** - Standardized customizations
- **promotions** - Chain-wide campaigns

### Key Features:
- Drill-down from organization → location → individual transactions
- Compare location performance
- Centralized menu and pricing management
- Bulk operations across locations
- Consolidated reporting
- Franchise fee calculation

---

## 14. Staff Time & Attendance

**Purpose**: Dedicated app for employee time tracking and scheduling.

### Key Tables Used:

#### Scheduling
- **employee_schedules** - Work schedules
- **users** - Employee profiles
- **locations** - Work locations

#### Time Tracking
- **time_clock_entries** - Clock in/out records
- **devices** - Time clock device/kiosk

#### Payroll Data
- **v_employee_work_hours** - Calculated hours worked
- **tip_distributions** - Tips earned
- **staff_commissions** - Commissions earned

#### Shifts
- **shifts** - Register shifts (for cashiers)
- **driver_shifts** - Delivery driver shifts

### Typical Workflow:

**Employee:**
1. Employee approaches time clock (tablet/kiosk)
2. Authenticates (PIN, badge, biometric)
3. Clock in → insert `time_clock_entries` (type: clock_in)
4. System checks against `employee_schedules` (late/early)
5. At end of shift, clock out → insert (type: clock_out)
6. View schedule from `employee_schedules`

**Manager:**
1. Create schedules in `employee_schedules`
2. View attendance reports from `time_clock_entries`
3. Approve manual time corrections
4. Export payroll data from `v_employee_work_hours`
5. Review variances (scheduled vs actual)

---

## 15. Maintenance/Setup App

**Purpose**: Technical app for IT staff to configure devices and troubleshoot.

### Key Tables Used:

#### Device Management
- **devices** - Register and configure devices
- **printer_configurations** - Printer routing setup

#### Network Configuration
- **devices.ip_address**, **connection_string** - Network settings
- **devices.last_heartbeat_at** - Device health monitoring

#### Kitchen Setup
- **kitchen_stations** - Configure stations
- **printer_configurations** - Link printers to stations

#### Restaurant Layout
- **floor_plans** - Upload floor plan images
- **table_sections** - Define sections
- **restaurant_tables** - Position tables on floor plan

#### Testing
- **audit_logs** - View system events
- All tables (read-only) - Troubleshooting data issues

### Typical Operations:
1. Register new device in `devices`
2. Configure printer routing in `printer_configurations`
3. Test device connectivity (heartbeat)
4. View audit logs for troubleshooting
5. Setup kitchen stations and routing
6. Configure table layouts
7. Device firmware updates (metadata stored in `devices.metadata`)

---

## 16. Customer Mobile App

**Purpose**: Consumer app for online ordering, loyalty, and order tracking.

### Key Tables Used:

#### Customer Profile
- **customers** - Profile and preferences
- **customer_addresses** - Saved delivery addresses
- **loyalty_tiers** - Current tier
- **loyalty_points_transactions** - Points balance and history

#### Ordering
- **products** & **categories** - Browse menu
- **orders** - Place orders (source: 'mobile_app')
- **order_items** & **order_item_modifiers** - Customize orders
- **delivery_zones** - Check delivery availability
- **promotions** - Available offers

#### Payment
- **payments** - Payment processing
- **customers.metadata** - Saved payment methods

#### Order Tracking
- **orders** - Order status
- **order_tracking_events** - Real-time updates
- **delivery_assignments** - Delivery progress
- **delivery_drivers** - Driver info (name, photo, vehicle)

#### Loyalty
- **loyalty_rewards** - Available rewards
- **loyalty_redemptions** - Redeem rewards
- **loyalty_tier_benefits** - Tier perks
- **customer_tier_history** - Tier progress

#### Reservations
- **restaurant_tables** - View available tables
- **reservations** - Make reservations

### Typical Workflow:

**Browse & Order:**
1. Customer logs in → query `customers`
2. Browse menu from `products` and `categories`
3. Add items to cart (create `order` with status: draft)
4. Select delivery or pickup
5. For delivery: choose from `customer_addresses` or add new
6. Apply promotion code from `promotions`
7. Process payment → `payments`
8. Submit order → status: submitted
9. Receive confirmation

**Track Order:**
1. Query `orders` by customer_id
2. Display status from `order_tracking_events`
3. For delivery: show driver location from `delivery_drivers.current_location`
4. Show estimated time from `delivery_assignments.estimated_delivery_time`

**Loyalty:**
1. View points balance (SUM of `loyalty_points_transactions`)
2. View current tier from `loyalty_tiers`
3. Browse `loyalty_rewards`
4. Redeem reward → create `loyalty_redemptions`
5. View tier benefits from `loyalty_tier_benefits`

---

## Database Views for Applications

### Views Commonly Used Across Apps:

#### v_active_kitchen_tickets
Used by: KDS, Waiter App, Admin Panel
```sql
-- Real-time active tickets for kitchen displays
SELECT * FROM v_active_kitchen_tickets
WHERE kitchen_station_id = :station_id;
```

#### v_open_orders_summary
Used by: Waiter App, Main POS, Admin Panel
```sql
-- Summary of in-progress orders
SELECT * FROM v_open_orders_summary
WHERE location_id = :location_id;
```

#### v_active_deliveries
Used by: Admin Panel, Delivery Management
```sql
-- Real-time delivery tracking
SELECT * FROM v_active_deliveries
WHERE driver_id = :driver_id OR location_id = :location_id;
```

#### v_driver_performance
Used by: Admin Panel, Driver App, Analytics
```sql
-- Driver statistics and earnings
SELECT * FROM v_driver_performance
WHERE driver_id = :driver_id;
```

#### v_employee_work_hours
Used by: Time & Attendance, Payroll, Admin Panel
```sql
-- Calculate hours worked
SELECT * FROM v_employee_work_hours
WHERE employee_id = :employee_id
AND work_date BETWEEN :start_date AND :end_date;
```

#### v_active_devices
Used by: Maintenance App, Admin Panel
```sql
-- Monitor device health
SELECT * FROM v_active_devices
WHERE location_id = :location_id;
```

---

## Common Integration Patterns

### Pattern 1: Real-Time Order Flow
```
Customer Mobile App → orders table → Kitchen Tickets
                                   ↓
                              KDS Updates → order_items.status
                                   ↓
                              Waiter App sees ready items
                                   ↓
                              Mark served → order.status = 'completed'
                                   ↓
                              Convert to sale record
```

### Pattern 2: Multi-Location Inventory
```
Retail POS (Location A) → inventory_transactions (SALE)
                              ↓
                        Deduct from Location A stock
                              ↓
                        Trigger reorder if below threshold
                              ↓
                        Create purchase_order
                              ↓
                        Handheld receives at warehouse
                              ↓
                        Create inventory_transfer to Location A
```

### Pattern 3: Loyalty Integration
```
Customer Mobile App → Place order
                         ↓
                    Link to customers record
                         ↓
                    Main POS processes payment
                         ↓
                    Create loyalty_points_transaction (EARNED)
                         ↓
                    Check tier upgrade threshold
                         ↓
                    Update customer.current_tier_id if qualified
                         ↓
                    Create customer_tier_history record
```

---

## Security & RLS (Row-Level Security)

All tables implement RLS policies ensuring:
- **Super Admins** can access all data across organizations
- **Regular users** only see data for their `organization_id`
- Drivers only see their assigned deliveries
- Customers only see their own orders and profile

Example RLS policy applied to every table:
```sql
CREATE POLICY org_isolation_select ON table_name
    FOR SELECT TO PUBLIC
    USING (
        organization_id IN (
            SELECT organization_id FROM users
            WHERE id = auth.uid() AND deleted_at IS NULL
        )
    );
```

---

## Performance Considerations

### Indexes for High-Traffic Queries

Each app's common queries are optimized with indexes:
- `orders(organization_id, status, order_date)` - Order lists
- `order_items(order_id, status)` - Order details
- `kitchen_tickets(kitchen_station_id, status)` - KDS queries
- `delivery_assignments(driver_id, status)` - Driver app
- `products(organization_id, barcode)` - Barcode lookups
- `inventory_transactions(product_id, location_id, transaction_date)` - Stock queries

### Materialized Views

Heavy analytics queries use pre-computed materialized views:
- Refresh scheduled (e.g., hourly, daily)
- Function: `refresh_all_analytics_views()`
- Used by: Owner Analytics, Admin Panel, Franchise HQ

---

## Conclusion

This comprehensive POS ecosystem database supports **16 different applications** with a unified data model:

✅ **50+ core tables** covering all aspects of restaurant and retail operations
✅ **Complete multi-tenancy** with organization-level isolation
✅ **Real-time operations** (orders, kitchen, delivery)
✅ **Advanced features** (loyalty, inventory tracking, staff management)
✅ **Robust security** with RLS policies on every table
✅ **Optimized performance** with materialized views and strategic indexing

Each application has clear, well-defined access patterns to the database, enabling independent development while maintaining data consistency across the entire ecosystem.
