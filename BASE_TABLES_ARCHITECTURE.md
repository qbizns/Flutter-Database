# Base Tables Architecture

## 🏗️ Plugin Architecture Overview

This system follows a **plugin architecture** where:
- **POS System** (`postgres/` folder) contains all **BASE TABLES**
- **Accounting Module** (`accounting/` folder) is a **PLUGIN** that references base tables via foreign keys
- **No data duplication** - accounting never creates its own copies of customers, suppliers, products, etc.

```
┌─────────────────────────────────────────────────────────────┐
│                    POS BASE TABLES (postgres/)              │
│                                                             │
│  - organizations                                            │
│  - users                                                    │
│  - customers                                                │
│  - suppliers                                                │
│  - products                                                 │
│  - locations                                                │
│  - sales                                                    │
│  - purchase_orders                                          │
│  - inventory_transactions                                   │
│  - payments                                                 │
│  - ... (all POS operations)                                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ Foreign Key References Only
                         │ (No Data Duplication)
                         │
┌────────────────────────▼────────────────────────────────────┐
│              ACCOUNTING PLUGIN (accounting/)                │
│                                                             │
│  - chart_of_accounts (references organizations)            │
│  - journal_entries (references organizations, users)       │
│  - vendor_bills (references suppliers, organizations)      │
│  - customer_invoices (references customers, organizations) │
│  - fixed_assets (references products, locations)           │
│  - ... (accounting-specific tables only)                   │
└─────────────────────────────────────────────────────────────┘
```

## 📊 Base Tables in POS System (postgres/)

### Core Tenant Tables (V001)
| Table | Location | Purpose |
|-------|----------|---------|
| `organizations` | V001_create_core_tenant_tables.sql | Multi-tenant isolation |
| `users` | V001_create_core_tenant_tables.sql | User authentication & authorization |
| `roles` | V001_create_core_tenant_tables.sql | Role-based access control |
| `permissions` | V001_create_core_tenant_tables.sql | Granular permissions |
| `user_organizations` | V001_create_core_tenant_tables.sql | User-org relationships |

### POS Core Tables (V002)
| Table | Location | Purpose |
|-------|----------|---------|
| `products` | V002_create_pos_core_tables.sql | Product master data |
| `product_categories` | V002_create_pos_core_tables.sql | Product categorization |
| `units_of_measure` | V002_create_pos_core_tables.sql | UOM definitions |
| `customers` | V002_create_pos_core_tables.sql | **Customer master data** |
| `sales` | V002_create_pos_core_tables.sql | **Sales transactions** |
| `sale_items` | V002_create_pos_core_tables.sql | Sales line items |
| `payments` | V002_create_pos_core_tables.sql | Payment records |
| `inventory_transactions` | V002_create_pos_core_tables.sql | Inventory movements |

### Additional POS Tables (V004)
| Table | Location | Purpose |
|-------|----------|---------|
| `suppliers` | V004_create_additional_pos_tables.sql | **Supplier/Vendor master data** |
| `purchase_orders` | V004_create_additional_pos_tables.sql | **Purchase orders** |
| `purchase_order_items` | V004_create_additional_pos_tables.sql | PO line items |
| `locations` | V004_create_additional_pos_tables.sql | **Warehouse/store locations** |
| `product_variants` | V004_create_additional_pos_tables.sql | Product variations |
| `promotions` | V004_create_additional_pos_tables.sql | Promotional campaigns |
| `expenses` | V004_create_additional_pos_tables.sql | Operating expenses |
| `pos_sessions` | V004_create_additional_pos_tables.sql | Cash register sessions |

## 🔗 Accounting Module References

### How Accounting References Base Tables

The accounting module **never duplicates** base tables. Instead, it references them:

#### Vendor Bills Reference Suppliers
```sql
-- accounting/migrations/V002_create_ap_ar_assets.sql
CREATE TABLE vendor_bills (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    supplier_id UUID NOT NULL REFERENCES suppliers(id),  -- ← References postgres/suppliers
    bill_date DATE NOT NULL,
    total_amount NUMERIC(15, 2),
    ...
);
```

#### Customer Invoices Reference Customers
```sql
-- accounting/migrations/V002_create_ap_ar_assets.sql
CREATE TABLE customer_invoices (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    customer_id UUID NOT NULL REFERENCES customers(id),  -- ← References postgres/customers
    invoice_date DATE NOT NULL,
    total_amount NUMERIC(15, 2),
    ...
);
```

#### Fixed Assets Reference Products & Locations
```sql
-- accounting/migrations/V002_create_ap_ar_assets.sql
CREATE TABLE fixed_assets (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    product_id UUID REFERENCES products(id),      -- ← References postgres/products
    location_id UUID REFERENCES locations(id),     -- ← References postgres/locations
    acquisition_date DATE NOT NULL,
    acquisition_cost NUMERIC(15, 2),
    ...
);
```

#### POS Account Mappings Reference Products & Suppliers
```sql
-- accounting/migrations/V004_create_pos_account_mappings.sql
CREATE TABLE pos_account_mappings (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    source_type VARCHAR(50),  -- 'product', 'supplier', 'customer', etc.
    source_id UUID,           -- ← References products(id), suppliers(id), customers(id)
    account_id UUID NOT NULL REFERENCES chart_of_accounts(id),
    ...
);
```

## ✅ Architecture Validation

### Base Tables Checklist

All required base tables exist in `postgres/`:

- [x] **organizations** (V001) - Multi-tenant isolation
- [x] **users** (V001) - User management
- [x] **customers** (V002) - Customer master data
- [x] **suppliers** (V004) - Supplier/Vendor master data
- [x] **products** (V002) - Product master data
- [x] **locations** (V004) - Warehouse/store locations
- [x] **sales** (V002) - Sales transactions
- [x] **purchase_orders** (V004, V022) - Purchase orders

### Foreign Key References

All accounting tables correctly reference base tables:

- [x] `vendor_bills.supplier_id` → `suppliers(id)`
- [x] `vendor_payments.supplier_id` → `suppliers(id)`
- [x] `customer_invoices.customer_id` → `customers(id)`
- [x] `customer_payments.customer_id` → `customers(id)`
- [x] `fixed_assets.product_id` → `products(id)`
- [x] `fixed_assets.location_id` → `locations(id)`
- [x] `pos_account_mappings.source_id` → `products(id)`, `suppliers(id)`, `customers(id)`, etc.
- [x] `inventory_valuation_settings.organization_id` → `organizations(id)`

### No Data Duplication

✅ **Confirmed**: The accounting module does NOT create its own copies of:
- customers
- suppliers
- products
- locations
- organizations
- users

✅ **Confirmed**: All references use foreign keys to base tables in `postgres/`

## 📁 Directory Structure

```
Flutter-Database/
├── postgres/                          # BASE TABLES (Foundation)
│   ├── migrations/
│   │   ├── V001_*_core_tenant_tables.sql      # organizations, users
│   │   ├── V002_*_pos_core_tables.sql         # customers, products, sales
│   │   ├── V004_*_additional_pos_tables.sql   # suppliers, locations, purchase_orders
│   │   └── V013-V022_*_pos_features.sql       # Advanced POS features
│   ├── seed_data/
│   └── schemas/
│
├── accounting/                        # PLUGIN MODULE (References base tables)
│   ├── migrations/
│   │   ├── V001_*_accounting_core.sql         # COA, GL, JE (references organizations)
│   │   ├── V002_*_ap_ar_assets.sql            # Vendor/Customer bills (references suppliers/customers)
│   │   ├── V003_*_odoo_extensions.sql         # Advanced features
│   │   ├── V004-V010_*_pos_integration.sql    # POS bridges (references POS tables)
│   │   └── V011-V013_*_posting_engine.sql     # Configuration-driven posting
│   ├── seed_data/
│   └── schemas/
│
└── BASE_TABLES_ARCHITECTURE.md        # This document
```

## 🎯 Design Principles

### 1. Single Source of Truth
- **Base tables** (customers, suppliers, products) exist ONLY in `postgres/`
- **Accounting tables** exist ONLY in `accounting/`
- No duplication = no data synchronization issues

### 2. Foreign Key Integrity
- All references use PostgreSQL foreign keys
- Database enforces referential integrity
- Cascading deletes configured appropriately

### 3. Plugin Independence
- POS system works without accounting module
- Accounting module cannot exist without POS base tables
- Clean separation of concerns

### 4. Multi-Tenancy
- All tables include `organization_id`
- Row-Level Security (RLS) policies enforce isolation
- Each organization sees only its own data

## 🔍 Verification Queries

### Check All Base Table References
```sql
-- Find all foreign key references from accounting to postgres tables
SELECT
    tc.table_schema,
    tc.table_name,
    kcu.column_name,
    ccu.table_schema AS foreign_table_schema,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_schema = 'accounting'
    AND ccu.table_schema = 'public'
ORDER BY tc.table_name, kcu.column_name;
```

### Verify No Duplicate Tables
```sql
-- Check for tables that exist in both accounting and public schemas
SELECT
    a.table_name AS duplicate_table
FROM information_schema.tables a
INNER JOIN information_schema.tables b
    ON a.table_name = b.table_name
WHERE a.table_schema = 'accounting'
    AND b.table_schema = 'public'
ORDER BY a.table_name;
-- Expected result: Empty (no duplicates)
```

## 🚀 Usage Pattern

### Creating a Vendor Bill (Accounting Module)

```typescript
// ✅ CORRECT: Reference existing supplier from base table
async function createVendorBill(supplierId: string, amount: number) {
    // 1. Verify supplier exists in base table (postgres/suppliers)
    const supplier = await db.query('SELECT * FROM suppliers WHERE id = $1', [supplierId]);

    // 2. Create vendor bill in accounting (references supplier)
    const vendorBill = await db.query(`
        INSERT INTO accounting.vendor_bills (
            organization_id,
            supplier_id,        -- ← References public.suppliers(id)
            bill_date,
            total_amount
        ) VALUES ($1, $2, $3, $4)
        RETURNING *
    `, [orgId, supplierId, billDate, amount]);

    return vendorBill;
}

// ❌ INCORRECT: Creating supplier in accounting module
async function createVendorBillWrong() {
    // Don't do this! Suppliers belong in postgres/suppliers
    await db.query('INSERT INTO accounting.suppliers ...');  // ❌ Wrong!
}
```

### Posting a Sale to Accounting

```typescript
// ✅ CORRECT: Post existing sale to accounting
async function postSaleToAccounting(saleId: string) {
    // 1. Get sale from base table (postgres/sales)
    const sale = await db.query('SELECT * FROM sales WHERE id = $1', [saleId]);

    // 2. Get customer from base table (postgres/customers)
    const customer = await db.query('SELECT * FROM customers WHERE id = $1', [sale.customer_id]);

    // 3. Create journal entry in accounting (references sale & customer)
    const je = await db.query(`
        INSERT INTO accounting.journal_entries (
            organization_id,
            reference_table,     -- 'sales'
            reference_id,        -- sale.id
            description
        ) VALUES ($1, 'sales', $2, $3)
        RETURNING *
    `, [orgId, saleId, `Sale to ${customer.full_name}`]);

    // 4. Update sale with accounting reference
    await db.query(`
        UPDATE sales
        SET accounting_journal_entry_id = $1,
            posted_to_accounting_at = CURRENT_TIMESTAMP
        WHERE id = $2
    `, [je.id, saleId]);
}
```

## 📝 Summary

✅ **Architecture is correctly implemented**
✅ **All base tables exist in `postgres/`**
✅ **Accounting module only references base tables (no duplication)**
✅ **Foreign key integrity enforced by database**
✅ **Single source of truth for customers, suppliers, products**
✅ **Clean plugin separation**

The system follows best practices for modular design with proper separation of concerns and referential integrity.
