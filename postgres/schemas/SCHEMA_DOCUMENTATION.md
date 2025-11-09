# POS SAAS Database Schema Documentation

## Overview

This document describes the complete database schema for the POS SAAS system, including all tables, relationships, and design decisions.

## Architecture

### Multi-Tenancy Strategy

The system uses a **hybrid multi-tenancy approach**:

1. **Shared Tables** (public schema):
   - Organizations
   - Users
   - Roles & Permissions
   - Audit Logs

2. **Tenant-Specific Tables**:
   - Products, Categories
   - Customers
   - Sales, Payments
   - Inventory

This approach provides:
- Strong data isolation between tenants
- Flexibility for per-tenant customization
- Efficient resource utilization
- Simplified backup and restore per tenant

## Database Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        ORGANIZATIONS                             │
│  (Tenant/Company Master)                                        │
└──────────────┬──────────────────────────────────────────────────┘
               │
               ├─────────────┐
               │             │
       ┌───────▼──────┐ ┌───▼──────────┐
       │    USERS     │ │  CATEGORIES  │
       │              │ │              │
       └───┬──────────┘ └──────┬───────┘
           │                   │
    ┌──────▼────────┐         │
    │  USER_ROLES   │         │
    │               │         │
    └──────┬────────┘         │
           │            ┌─────▼──────┐
    ┌──────▼────────┐  │  PRODUCTS  │
    │    ROLES      │  │            │
    │               │  └─────┬──────┘
    └──────┬────────┘        │
           │                 │
    ┌──────▼────────────┐    │
    │ ROLE_PERMISSIONS  │    │
    │                   │    │
    └──────┬────────────┘    │
           │                 │
    ┌──────▼────────┐        │
    │  PERMISSIONS  │        │
    └───────────────┘        │
                             │
    ┌────────────┐           │
    │ CUSTOMERS  │           │
    └─────┬──────┘           │
          │                  │
    ┌─────▼──────────────────▼─────┐
    │         SALES                 │
    │  (Transaction Header)         │
    └─────┬─────────────────────────┘
          │
    ┌─────▼──────────┐  ┌────────────────────┐
    │   SALE_ITEMS   │  │     PAYMENTS       │
    │  (Line Items)  │  │                    │
    └────────────────┘  └────────────────────┘
          │
    ┌─────▼──────────────────────┐
    │  INVENTORY_TRANSACTIONS    │
    │  (Stock Movement Audit)    │
    └────────────────────────────┘
```

## Table Descriptions

### Core System Tables

#### 1. organizations
**Purpose**: Multi-tenant organization/company master
**Key Fields**:
- `id` - Unique identifier (UUID)
- `slug` - URL-friendly identifier for schema naming
- `status` - trial, active, suspended, cancelled
- `plan` - Subscription plan
- `max_users`, `max_products`, `max_locations` - Quotas

**Indexes**:
- Primary key on `id`
- Unique on `slug`
- Index on `status`

**Design Notes**:
- Each organization represents a separate tenant
- Slug is used for creating tenant-specific schemas
- Settings stored as JSONB for flexibility

---

#### 2. users
**Purpose**: System users with organization mapping
**Key Fields**:
- `id` - Unique identifier (UUID)
- `organization_id` - Foreign key to organizations
- `email` - Email address (unique per organization)
- `status` - active, inactive, suspended, pending
- `full_name` - Generated column from first_name + last_name

**Relationships**:
- Belongs to: `organizations`
- Has many: `user_roles`, `sales` (as cashier)

**Security Features**:
- Password hash storage
- 2FA support
- Failed login tracking
- Account lockout mechanism

**Indexes**:
- Primary key on `id`
- Unique on `(organization_id, email)`
- Index on `email`, `status`, `full_name`

---

#### 3. roles
**Purpose**: Role definitions for RBAC
**Key Fields**:
- `id` - Unique identifier (UUID)
- `organization_id` - NULL for system roles, or organization-specific
- `slug` - URL-friendly identifier
- `is_system_role` - Cannot be deleted if true
- `is_default` - Assigned to new users

**System Roles**:
- Super Admin (full access)
- Admin (organization admin)
- Manager (store manager)
- Cashier (default for sales)
- Viewer (read-only)

**Relationships**:
- Belongs to: `organizations` (optional)
- Has many: `role_permissions`, `user_roles`

---

#### 4. permissions
**Purpose**: Granular permission definitions
**Key Fields**:
- `slug` - Unique identifier (e.g., 'products.create')
- `resource` - The resource (e.g., 'products')
- `action` - The action (e.g., 'create')
- `category` - Grouping (e.g., 'inventory')

**Permission Categories**:
- `inventory` - Product and stock management
- `sales` - Sales transactions
- `customers` - Customer management
- `admin` - User and system management
- `reports` - Reporting and analytics

---

#### 5. role_permissions
**Purpose**: Maps permissions to roles
**Type**: Junction table
**Relationships**:
- Many-to-many between `roles` and `permissions`

---

#### 6. user_roles
**Purpose**: Maps users to roles
**Type**: Junction table
**Relationships**:
- Many-to-many between `users` and `roles`

---

#### 7. audit_logs
**Purpose**: System-wide audit trail
**Key Fields**:
- `action` - What happened
- `resource_type` - What was affected
- `resource_id` - Which record
- `old_values`, `new_values`, `changes` - What changed

**Use Cases**:
- Compliance and auditing
- Debugging
- User activity tracking
- Change history

**Performance Notes**:
- Consider partitioning by `created_at` for large datasets
- Indexes on `organization_id`, `user_id`, `resource_type`

---

### POS Tables

#### 8. categories
**Purpose**: Hierarchical product categorization
**Key Fields**:
- `parent_id` - Self-referencing for hierarchy
- `level` - Depth in hierarchy
- `path` - Full path (e.g., 'electronics/phones/smartphones')
- `sort_order` - Display order

**Features**:
- Unlimited nesting levels
- Path-based hierarchy for easy queries
- Soft deletes

**Relationships**:
- Belongs to: `organizations`
- Has many: `products`, `categories` (children)
- Belongs to: `categories` (parent)

---

#### 9. products
**Purpose**: Product/item catalog
**Key Fields**:
- `sku` - Stock Keeping Unit (unique per org)
- `barcode` - EAN/UPC barcode
- `cost_price` - Purchase cost
- `selling_price` - Retail price
- `current_stock` - Current inventory level
- `track_inventory` - Enable/disable inventory tracking

**Features**:
- Multiple price points (cost, selling, compare_at)
- Tax rate support
- Inventory tracking (optional)
- Service vs physical products
- Composite/bundle products
- Product variants support
- Custom fields via JSONB

**Relationships**:
- Belongs to: `organizations`, `categories`
- Has many: `sale_items`, `inventory_transactions`

**Constraints**:
- Unique `(organization_id, sku)`
- Unique `(organization_id, barcode)`
- Positive price check

---

#### 10. customers
**Purpose**: Customer master data
**Key Fields**:
- `customer_code` - Internal identifier
- `full_name` - Generated from first + last name
- `loyalty_points` - Loyalty program points
- `credit_limit` - Credit allowed
- `outstanding_balance` - Current balance
- `total_purchases`, `total_orders` - Cached statistics

**Features**:
- Complete contact information
- Business customer support (tax_number)
- Loyalty program integration
- Credit management
- Purchase history tracking
- Custom fields via JSONB

**Relationships**:
- Belongs to: `organizations`
- Has many: `sales`

---

#### 11. sales
**Purpose**: Sales transaction header
**Key Fields**:
- `sale_number` - Receipt/Invoice number
- `transaction_type` - sale, return, exchange, void
- `subtotal`, `tax_amount`, `discount_amount`, `total_amount`
- `payment_status` - pending, completed, failed, refunded
- `transaction_date` - When transaction occurred

**Transaction Flow**:
1. Create sale (pending)
2. Add line items (sale_items)
3. Process payment(s)
4. Complete sale
5. Update inventory

**Relationships**:
- Belongs to: `organizations`, `customers`, `users` (cashier)
- Has many: `sale_items`, `payments`, `inventory_transactions`

**Features**:
- Multiple payment methods
- Partial payments
- Discounts (percentage or fixed)
- Transaction types for returns/exchanges

---

#### 12. sale_items
**Purpose**: Line items for sales
**Key Fields**:
- `product_id` - Reference to product (can be NULL if deleted)
- `product_name`, `product_sku` - Snapshot at time of sale
- `quantity` - Quantity sold
- `unit_price` - Price per unit (snapshot)
- `cost_price` - Cost at time of sale

**Design Notes**:
- Stores product snapshot to preserve history
- Allows product deletion without losing sale history
- Individual item discounts supported

**Relationships**:
- Belongs to: `sales`, `products`, `organizations`

**Constraints**:
- Positive quantity and prices

---

#### 13. payments
**Purpose**: Payment records
**Key Fields**:
- `payment_method` - cash, card, mobile_money, bank_transfer, other
- `payment_status` - pending, completed, failed, refunded
- `amount` - Payment amount
- `transaction_id` - External gateway reference

**Features**:
- Multiple payment methods per sale
- Split payments
- Payment gateway integration support
- Card details (last 4 digits)
- Mobile money/bank transfer details

**Relationships**:
- Belongs to: `sales`, `organizations`

---

#### 14. inventory_transactions
**Purpose**: Complete inventory audit trail
**Key Fields**:
- `transaction_type` - purchase, sale, adjustment, return, transfer, waste
- `quantity` - Change amount (positive or negative)
- `balance_after` - Stock level after transaction
- `sale_id` - Link to sale if applicable

**Use Cases**:
- Track all stock movements
- Audit trail for compliance
- Calculate stock valuation
- Identify discrepancies
- Generate inventory reports

**Relationships**:
- Belongs to: `organizations`, `products`, `sales` (optional)

**Design Notes**:
- Immutable log (no updates/deletes)
- Tracks balance after each transaction
- Links to sales for automatic tracking

---

## Data Types and Enums

### Custom Types

```sql
-- User status
user_status: 'active', 'inactive', 'suspended', 'pending'

-- Organization status
organization_status: 'trial', 'active', 'suspended', 'cancelled'

-- Payment status
payment_status: 'pending', 'completed', 'failed', 'refunded', 'cancelled'

-- Payment method
payment_method: 'cash', 'card', 'mobile_money', 'bank_transfer', 'other'

-- Transaction type
transaction_type: 'sale', 'return', 'exchange', 'void'

-- Inventory transaction type
inventory_transaction_type: 'purchase', 'sale', 'adjustment', 'return', 'transfer', 'waste'
```

## Common Patterns

### Audit Fields
All tables include:
```sql
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
deleted_at TIMESTAMP WITH TIME ZONE  -- For soft deletes
created_by UUID REFERENCES users(id)
updated_by UUID REFERENCES users(id)
```

### Soft Deletes
- Most tables use `deleted_at` for soft deletes
- Indexes exclude deleted records: `WHERE deleted_at IS NULL`
- Preserves referential integrity
- Allows data recovery

### JSONB Fields
Used for flexibility:
- `settings` - Configurable settings
- `metadata` - Additional data
- `custom_fields` - User-defined fields

### Auto-Update Triggers
All tables with `updated_at` have trigger:
```sql
CREATE TRIGGER update_{table}_updated_at
    BEFORE UPDATE ON {table}
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

## Indexes Strategy

### Primary Indexes
- All tables use UUID primary keys
- Foreign keys are indexed

### Performance Indexes
- Frequently queried fields (email, sku, barcode)
- Status fields for filtering
- Date fields for range queries
- Full-text search on names

### Partial Indexes
- Exclude soft-deleted records
- Active records only

## Security Considerations

### Row-Level Security (Future)
- Implement RLS for tenant isolation
- User-based access control
- Organization-based filtering

### Password Storage
- BCrypt hashed passwords
- Minimum 10 rounds
- Never store plain text

### Audit Trail
- All sensitive operations logged
- IP address tracking
- User agent capture

## Performance Optimization

### Query Optimization
- Use indexes effectively
- Avoid N+1 queries
- Use EXPLAIN ANALYZE
- Monitor slow queries

### Caching Strategy
- Cache product catalog
- Cache user permissions
- Cache organization settings

### Partitioning (Future)
- Partition `audit_logs` by date
- Partition `inventory_transactions` by date
- Partition `sales` by date

## Backup Strategy

### Regular Backups
- Daily full backups
- Hourly incremental backups
- Point-in-time recovery enabled

### Per-Tenant Backups
- Schema-based backups
- Easier restore for individual tenants
- Compliance with data retention policies

## Migration Strategy

### Version Control
- All changes via migrations
- Sequential versioning (V001, V002, etc.)
- Timestamp in filename
- Never modify existing migrations

### Rollback Strategy
- Test migrations in staging
- Keep rollback scripts
- Document breaking changes

---

### Additional POS Tables (Implemented in V004)

#### 15. suppliers
**Purpose**: Supplier/vendor master data
**Key Fields**:
- `supplier_code` - Unique supplier identifier
- `name` - Supplier business name
- `contact_person` - Primary contact name
- `payment_terms` - Payment terms (e.g., Net 30)
- `credit_limit` - Maximum credit allowed
- `outstanding_balance` - Current balance owed

**Features**:
- Complete contact and address information
- Tax number for business suppliers
- Payment terms management
- Credit limit tracking
- Purchase statistics (total purchases, total orders)
- Last order date tracking

**Relationships**:
- Belongs to: `organizations`
- Has many: `purchase_orders`

---

#### 16. purchase_orders
**Purpose**: Purchase orders for inventory procurement
**Key Fields**:
- `po_number` - Purchase order number
- `supplier_id` - Reference to supplier
- `po_date` - Order date
- `expected_delivery_date` - Expected delivery
- `status` - draft, pending, approved, ordered, partial, received, cancelled
- `payment_status` - pending, partial, paid, overdue

**Features**:
- Multi-status workflow (draft → approved → ordered → received)
- Financial tracking (subtotal, tax, discount, shipping)
- Payment tracking
- Approval workflow
- Delivery date management

**Relationships**:
- Belongs to: `organizations`, `suppliers`
- Has many: `purchase_order_items`

---

#### 17. purchase_order_items
**Purpose**: Line items for purchase orders
**Key Fields**:
- `product_id` - Reference to product
- `product_name`, `product_sku` - Product snapshot
- `quantity_ordered` - Quantity ordered
- `quantity_received` - Quantity received (partial receives)
- `unit_cost` - Purchase cost per unit

**Features**:
- Product snapshot preservation
- Partial receive tracking
- Discount and tax per line item
- Unit of measure support

**Relationships**:
- Belongs to: `purchase_orders`, `products`, `organizations`

---

#### 18. locations
**Purpose**: Store locations and branches
**Key Fields**:
- `location_code` - Unique location identifier
- `name` - Location name
- `location_type` - store, warehouse, headquarters, kiosk, online, other
- `manager_user_id` - Location manager
- `is_primary` - Primary location flag
- `allow_sales`, `allow_purchases` - Capability flags

**Features**:
- Multi-location support
- Complete address information
- Business hours tracking (JSONB)
- Timezone support
- Location-specific tax rates
- Manager assignment

**Relationships**:
- Belongs to: `organizations`
- Has many: `products`, `sales`, `expenses`, `shifts`
- Belongs to: `users` (manager)

---

#### 19. product_variants
**Purpose**: Product variations (size, color, etc.)
**Key Fields**:
- `product_id` - Parent product reference
- `variant_name` - Variant display name
- `sku`, `barcode` - Variant-specific identifiers
- `attributes` - Variation attributes (JSONB, e.g., {"size": "L", "color": "Blue"})
- `is_default` - Default variant flag

**Features**:
- Flexible attributes via JSONB
- Variant-specific pricing (can override parent)
- Independent inventory tracking per variant
- Physical properties (weight, dimensions)
- Variant-specific images
- Sort order for display

**Relationships**:
- Belongs to: `products`, `organizations`
- Has many: `sale_items`

---

#### 20. promotions
**Purpose**: Promotions and discount campaigns
**Key Fields**:
- `promotion_code` - Unique promotion code
- `promotion_type` - percentage, fixed_amount, buy_x_get_y, bundle, quantity_discount
- `discount_value` - Discount amount or percentage
- `applies_to` - all, specific_products, specific_categories, cart_total
- `start_date`, `end_date` - Promotion validity period

**Features**:
- Multiple promotion types
- Product/category targeting
- Minimum purchase requirements
- Usage limits (total and per customer)
- Date range validation
- Combinable promotions
- Priority system for applying multiple promotions
- Buy X Get Y rules

**Relationships**:
- Belongs to: `organizations`
- Has many: `promotion_usage`

---

#### 21. promotion_usage
**Purpose**: Track promotion usage per transaction
**Key Fields**:
- `promotion_id` - Reference to promotion
- `sale_id` - Reference to sale
- `customer_id` - Customer who used promotion
- `discount_amount` - Actual discount applied
- `used_at` - When promotion was used

**Use Cases**:
- Track promotion effectiveness
- Enforce usage limits
- Generate promotion reports
- Customer promotion history

**Relationships**:
- Belongs to: `promotions`, `sales`, `customers`, `organizations`

---

#### 22. expenses
**Purpose**: Business expense tracking
**Key Fields**:
- `expense_number` - Unique expense identifier
- `expense_date` - Date of expense
- `category` - Expense category (rent, utilities, salaries, supplies, marketing, other)
- `payee_name` - Who was paid
- `amount`, `tax_amount`, `total_amount` - Financial amounts
- `status` - pending, approved, paid, rejected, cancelled

**Features**:
- Category and subcategory classification
- Multiple payment methods
- Approval workflow
- Receipt and attachment tracking
- Purchase order reference
- Location-specific expenses

**Relationships**:
- Belongs to: `organizations`, `locations` (optional), `purchase_orders` (optional)
- Has references: `users` (created_by, approved_by)

---

#### 23. shifts
**Purpose**: Cashier shifts and cash register management
**Key Fields**:
- `shift_number` - Unique shift identifier
- `user_id` - Cashier for this shift
- `location_id` - Location where shift occurred
- `start_time`, `end_time` - Shift time range
- `status` - open, closed, suspended
- `opening_cash`, `expected_cash`, `actual_cash` - Cash reconciliation
- `cash_difference` - Difference between expected and actual

**Features**:
- Cash reconciliation tracking
- Sales summary per shift
- Payment method breakdown (JSONB)
- Transaction statistics
- Opening and closing notes
- Shift management workflow

**Relationships**:
- Belongs to: `organizations`, `locations`, `users` (cashier and closed_by)
- Has many: `sales` (via shift_id)

**Use Cases**:
- Cash register management
- Cashier accountability
- End-of-day reconciliation
- Shift performance tracking

---

## Enhanced Relationships (V004 Updates)

### New Foreign Key References
The following foreign keys were added to existing tables:

- `products.location_id` → `locations.id` - Track product location
- `sales.location_id` → `locations.id` - Track sale location
- `sales.shift_id` → `shifts.id` - Link sales to shifts
- `sale_items.product_variant_id` → `product_variants.id` - Support variant sales

---

## Future Enhancements

### Planned Features
1. **Inter-Location Transfers**
   - Transfer inventory between locations
   - Transfer orders and tracking
   - In-transit inventory management

2. **Advanced Reporting Tables**
   - Materialized views for analytics
   - Pre-aggregated data
   - Business intelligence dashboards

3. **Integrations**
   - E-commerce platforms
   - Accounting software (QuickBooks, Xero)
   - Payment gateways
   - Email/SMS notifications

4. **Advanced Features**
   - Subscription products
   - Rental/lease management
   - Booking/reservations
   - Service appointments

5. **Customer Loyalty Program**
   - Points redemption rules
   - Loyalty tiers
   - Rewards catalog
   - Member benefits

6. **Advanced Inventory**
   - Serial number tracking
   - Batch/lot tracking
   - Expiration date management
   - Cycle counting

## Appendix

### Naming Conventions
- Tables: plural, lowercase (e.g., `products`)
- Columns: snake_case (e.g., `created_at`)
- Indexes: `idx_{table}_{column(s)}` (e.g., `idx_products_sku`)
- Foreign keys: `{table}_id` (e.g., `organization_id`)
- Junction tables: `{table1}_{table2}` (e.g., `user_roles`)

### SQL Standards
- Use `TIMESTAMP WITH TIME ZONE` for all timestamps
- Use `NUMERIC` for currency (avoid FLOAT/REAL)
- Use `UUID` for primary keys
- Use `TEXT` for unlimited strings
- Use `VARCHAR(n)` for limited strings

### Documentation Updates
- Update this document with schema changes
- Document design decisions
- Include migration rationale
- Keep examples current
