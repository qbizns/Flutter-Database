# Enhanced POS Features (V005-V008)

## Overview

Migrations V005-V008 add advanced POS capabilities including inter-location transfers, advanced inventory management, enhanced loyalty programs, and comprehensive reporting infrastructure.

---

## V005: Inter-Location Transfers

### Tables Added

#### 1. inventory_transfers
**Purpose**: Track inventory movement between locations with complete workflow

**Status Workflow**:
```
draft → requested → approved → in_transit → partially_received → received
                                  ↓
                              cancelled/rejected
```

**Key Features**:
- Complete transfer workflow with approval process
- Carrier and tracking number support
- Expected vs actual delivery tracking
- Shipping cost tracking
- Multi-user approval chain

**Key Fields**:
- `transfer_number` - Unique transfer identifier
- `from_location_id`, `to_location_id` - Source and destination
- `status` - Current workflow status
- `tracking_number` - Carrier tracking
- `requested_by`, `approved_by`, `shipped_by`, `received_by` - Audit trail

#### 2. inventory_transfer_items
**Purpose**: Line items for transfers with quantity tracking

**Key Features**:
- Quantity requested vs shipped vs received
- Variance tracking and reporting
- Automatic status updates
- Cost tracking for valuation

**Automatic Functions**:
- `calculate_transfer_item_variance()` - Auto-calculates variances
- `update_transfer_status()` - Updates header status based on item completion

#### 3. v_inventory_transfer_summary (View)
**Purpose**: Summary view of transfers with aggregated data

---

## V006: Advanced Inventory Management

### Tables Added

#### 1. product_serial_numbers
**Purpose**: Track individual items by serial number

**Use Cases**:
- Electronics warranty tracking
- Equipment maintenance history
- Asset management
- Theft prevention

**Key Features**:
- Complete lifecycle tracking (in_stock → sold → returned)
- Warranty information (start, end, provider)
- Purchase and sale linkage
- Location tracking

#### 2. product_batches
**Purpose**: Batch/lot tracking with expiration management

**Use Cases**:
- Food products (expiration dates)
- Pharmaceuticals (lot tracking)
- Chemicals (batch safety)
- Quality control

**Key Features**:
- Manufacturing and expiration date tracking
- Initial vs current quantity management
- Quality control status (pending, passed, failed, quarantine)
- Supplier batch number linkage
- Automatic expiration alerts

#### 3. batch_transactions
**Purpose**: Complete audit trail of all batch movements

**Transaction Types**:
- sale, adjustment, return, waste, transfer, expiration

#### 4. cycle_counts
**Purpose**: Physical inventory cycle counting

**Count Types**:
- cycle - Regular rotating counts
- full - Complete inventory count
- spot - Random sample counts
- blind - Count without system quantities

**Key Features**:
- Planned vs actual tracking
- Variance analysis
- Multi-user workflow (planner, counter, approver)
- Statistics aggregation

#### 5. cycle_count_items
**Purpose**: Individual product counts within a cycle count

**Key Features**:
- System vs counted quantity comparison
- Automatic variance calculation
- Variance percentage and value
- Recount support
- Adjustment tracking

**Automatic Functions**:
- `calculate_cycle_count_variance()` - Auto-calculates all variance metrics

#### 6. stock_adjustment_reasons
**Purpose**: Pre-defined reasons for inventory adjustments

**System Reasons** (Pre-loaded):
- DAMAGE - Damaged goods
- THEFT - Theft/shrinkage
- EXPIRED - Expired products
- FOUND - Found inventory
- RETURN_SUPPLIER - Return to supplier
- SAMPLE - Product samples
- PROMOTION - Promotional giveaway
- COUNT_ERROR - Counting error
- QUALITY_FAIL - Failed quality check
- RESTOCK - Restocking adjustment

---

## V007: Enhanced Loyalty Program

### Tables Added

#### 1. loyalty_tiers
**Purpose**: Customer loyalty tiers/levels (Bronze, Silver, Gold, Platinum)

**Key Features**:
- Points threshold for tier promotion
- Annual spend threshold
- Points multiplier (e.g., 1.5x for Gold)
- Automatic discount percentage
- Visual customization (color, icon, badge)

#### 2. loyalty_tier_benefits
**Purpose**: Specific benefits for each tier

**Benefit Types**:
- discount - Percentage or fixed discount
- free_shipping - Free or discounted shipping
- birthday_bonus - Birthday month benefits
- early_access - Early sale access
- priority_support - VIP support
- exclusive_products - Limited edition access
- free_product - Free item benefits

#### 3. loyalty_points_rules
**Purpose**: Rules for earning points across different activities

**Rule Types**:
- purchase - Points per dollar spent
- signup - Welcome bonus
- birthday - Birthday bonus
- referral - Refer a friend
- review - Product review points
- social_share - Social media sharing
- manual - Manual point adjustments

**Key Features**:
- Flexible applicability (all, specific products, categories, tiers)
- Minimum purchase requirements
- Maximum points limits (per transaction, day, month)
- Date range validation
- Priority-based rule application

#### 4. loyalty_rewards
**Purpose**: Catalog of rewards customers can redeem

**Reward Types**:
- discount_percentage - Percentage off purchase
- discount_fixed - Fixed dollar amount off
- free_product - Free product reward
- free_shipping - Free shipping benefit
- gift_card - Gift card reward
- experience - Special experience (VIP event, etc.)

**Key Features**:
- Points cost tracking
- Availability limits (total and per customer)
- Tier restrictions
- Date range availability
- Featured rewards
- Automatic redemption counting

#### 5. loyalty_redemptions
**Purpose**: Track reward redemptions and fulfillment

**Workflow**:
```
pending → approved → fulfilled
    ↓
cancelled/expired
```

**Key Features**:
- Redemption number tracking
- Sale linkage (if used in transaction)
- Expiry date management
- Fulfillment tracking
- Multi-status workflow

#### 6. loyalty_points_transactions
**Purpose**: Complete audit trail of all points movements

**Transaction Types**:
- earned - Points earned
- redeemed - Points spent
- expired - Points expired
- adjusted - Manual adjustment
- bonus - Bonus points
- refunded - Refund points

#### 7. customer_tier_history
**Purpose**: Track customer tier changes over time

**Change Types**:
- upgrade - Tier promotion
- downgrade - Tier demotion
- initial - First tier assignment
- manual - Manual tier change

**Key Features**:
- Previous tier tracking
- Qualification metrics (points, spend, purchases)
- Effective date management
- Time-limited tier assignments

### Enhanced Tables

#### customers (updated)
**New Fields**:
- `current_tier_id` - Current loyalty tier
- `tier_since` - When they reached current tier
- `tier_expiry_date` - Tier expiration (optional)

---

## V008: Reporting & Analytics Infrastructure

### Materialized Views

#### 1. mv_daily_sales_summary
**Purpose**: Daily sales aggregation by organization and location

**Metrics**:
- Total transactions and unique customers
- Gross sales, discounts, tax, net sales
- Average transaction value
- Average customer spend
- Total items sold
- Payment method breakdown (JSONB)

**Refresh**: Recommended hourly or after end-of-day

#### 2. mv_product_performance
**Purpose**: Product sales performance and inventory analysis

**Time Periods**:
- Last 30 days
- Last 90 days
- All time

**Metrics**:
- Transactions and units sold
- Revenue and profit
- Profit margin percentage
- Days of stock remaining
- Last sale date

**Use Cases**:
- Identify top sellers
- Find slow-moving inventory
- Calculate reorder quantities
- Analyze profitability

#### 3. mv_customer_analytics
**Purpose**: Customer behavior analysis with RFM segmentation

**Metrics**:
- Lifetime value (LTV)
- Average order value
- Purchase frequency (30d, 90d, 365d)
- Recency metrics (days since last purchase)
- Total items purchased
- Favorite categories (top 3)

**Customer Segments** (Auto-calculated):
- Active - Purchased within 30 days
- At Risk - Purchased 30-90 days ago
- Dormant - Purchased 90-180 days ago
- Lost - No purchase in 180+ days

#### 4. mv_inventory_valuation
**Purpose**: Current inventory valuation and stock status

**Calculations**:
- Inventory value at cost
- Inventory value at retail
- Potential profit
- Stock status (in_stock, low_stock, out_of_stock, not_tracked)
- Variant and batch counts

#### 5. mv_location_performance
**Purpose**: Location-based performance metrics

**Metrics**:
- Sales (30 days and all time)
- Product count and inventory value
- Expenses (30 days)
- Shift count (30 days)
- Last sale date

#### 6. mv_promotion_effectiveness
**Purpose**: Promotion campaign effectiveness and ROI

**Metrics**:
- Total uses and unique customers
- Total discount given
- Total sales value
- Average sale value
- ROI ratio (sales / discount given)

### Helper Views (Regular Views)

#### v_top_selling_products_30d
Top products ranked by revenue in last 30 days

#### v_low_stock_alerts
Products requiring reorder attention

#### v_customer_segments_summary
Customer segmentation summary with metrics per segment

### Functions

#### refresh_all_analytics_views(concurrent_refresh BOOLEAN)
**Purpose**: Refresh all materialized views

**Parameters**:
- `concurrent_refresh` - TRUE for non-blocking refresh (requires unique indexes)

**Returns**: Table with view name, status, and execution time

**Usage**:
```sql
-- Refresh all views concurrently
SELECT * FROM refresh_all_analytics_views(true);

-- Refresh all views (blocking but faster)
SELECT * FROM refresh_all_analytics_views(false);
```

### Scheduled Refresh (Optional - requires pg_cron)

Automated refresh schedules:
- **Hourly**: Concurrent refresh for real-time dashboards
- **Daily (midnight)**: Full refresh for comprehensive reports

---

## Database Growth Summary

### Total Tables by Migration

- **V001**: 7 tables (Core tenant system)
- **V002**: 7 tables (POS core)
- **V003**: 0 tables (RLS policies only)
- **V004**: 9 tables (Additional POS features)
- **V005**: 2 tables (Inter-location transfers)
- **V006**: 6 tables (Advanced inventory)
- **V007**: 7 tables (Enhanced loyalty)
- **V008**: 6 materialized views, 3 helper views

**Grand Total**: 38 tables + 9 views + 10+ helper functions

### Total Seed Data Records

- Organizations: 3
- Users: 7
- Products: 16
- Product Variants: 12
- Customers: 10
- Sales: 8
- Suppliers: 6
- Locations: 6
- Promotions: 7
- Expenses: 10
- Shifts: 6
- **NEW - Inventory Transfers**: 3
- **NEW - Serial Numbers**: 8
- **NEW - Product Batches**: 4
- **NEW - Cycle Counts**: 2
- **NEW - Loyalty Tiers**: 7
- **NEW - Loyalty Rewards**: 7
- **NEW - Loyalty Redemptions**: 4

---

## Performance Considerations

### Materialized Views

**Pros**:
- Fast query performance (pre-computed)
- Complex aggregations done once
- Reduce load on transactional tables

**Cons**:
- Require periodic refresh
- Use additional storage
- Refresh can be resource-intensive

**Best Practices**:
1. Refresh during low-traffic periods
2. Use concurrent refresh for zero downtime
3. Monitor refresh execution time
4. Consider partitioning for large datasets

### Indexes

All new tables have comprehensive indexing:
- Primary keys (UUID)
- Foreign keys
- Status columns
- Date columns
- Frequently queried fields
- Unique constraints

### RLS Policies

All tables have complete Row-Level Security:
- Super admin bypass
- Organization-scoped access
- Separate policies for SELECT, INSERT, UPDATE, DELETE

---

## Migration Dependencies

```
V001 (Core Tenants)
  ↓
V002 (POS Core) → V003 (RLS)
  ↓                    ↓
V004 (Additional POS) ←┘
  ↓
V005 (Transfers) → V006 (Adv Inventory) → V007 (Loyalty) → V008 (Analytics)
```

---

## Next Steps

### Option 1: API Development
Build REST API layer for Flutter app integration

### Option 2: Flutter Models
Generate Dart models and data access layer

### Option 3: Additional Features
- Payment gateway integration
- Email/SMS notifications
- E-commerce platform integration
- Advanced reporting dashboards
- Mobile app sync infrastructure
- Offline SQLite schema generation

---

## Testing Queries

### Test Inter-Location Transfers
```sql
-- View transfer summary
SELECT * FROM v_inventory_transfer_summary
WHERE organization_id = '11111111-1111-1111-1111-111111111111';

-- Check transfer items with variance
SELECT * FROM inventory_transfer_items
WHERE variance_quantity != 0;
```

### Test Serial Number Tracking
```sql
-- Find all in-stock serialized products
SELECT * FROM product_serial_numbers
WHERE status = 'in_stock' AND deleted_at IS NULL;

-- Track product lifecycle
SELECT serial_number, status, purchase_date, sale_date
FROM product_serial_numbers
WHERE product_id = 'p0000001-0000-0000-0000-000000000001';
```

### Test Batch Tracking
```sql
-- Find expiring batches
SELECT * FROM product_batches
WHERE expiration_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '7 days'
AND status = 'active';

-- Batch transaction history
SELECT * FROM batch_transactions
WHERE batch_id = 'pb000002-0000-0000-0000-000000000001'
ORDER BY transaction_date DESC;
```

### Test Loyalty Program
```sql
-- Customer tier distribution
SELECT lt.tier_name, COUNT(*) as customer_count
FROM customers c
JOIN loyalty_tiers lt ON c.current_tier_id = lt.id
WHERE c.deleted_at IS NULL
GROUP BY lt.tier_name;

-- Points transactions by customer
SELECT customer_id, transaction_type, SUM(points) as total_points
FROM loyalty_points_transactions
GROUP BY customer_id, transaction_type;

-- Reward redemption report
SELECT r.reward_name, COUNT(*) as redemptions, SUM(lr.points_redeemed) as total_points
FROM loyalty_redemptions lr
JOIN loyalty_rewards r ON lr.reward_id = r.id
WHERE lr.status = 'fulfilled'
GROUP BY r.reward_name;
```

### Test Analytics Views
```sql
-- Refresh all analytics
SELECT * FROM refresh_all_analytics_views(false);

-- Daily sales trend
SELECT sale_date, gross_sales, total_transactions
FROM mv_daily_sales_summary
WHERE organization_id = '11111111-1111-1111-1111-111111111111'
ORDER BY sale_date DESC
LIMIT 30;

-- Top selling products
SELECT * FROM v_top_selling_products_30d
WHERE organization_id = '11111111-1111-1111-1111-111111111111'
LIMIT 10;

-- Customer segmentation
SELECT * FROM v_customer_segments_summary
WHERE organization_id = '11111111-1111-1111-1111-111111111111';

-- Low stock alerts
SELECT * FROM v_low_stock_alerts
WHERE organization_id = '11111111-1111-1111-1111-111111111111';
```

---

**Status**: ✅ Production Ready
**Total Implementation**: 4 Migrations, 23 New Tables, 9 Views, Multiple Helper Functions
