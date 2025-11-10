# POS SAAS Database Structure

This directory contains the PostgreSQL database schema, migrations, and seed data for the Point of Sale SAAS system.

## Directory Structure

```
postgres/
├── migrations/          # Database migration files (DDL)
├── schemas/            # Schema definitions and documentation
├── seed_data/          # DML files with dummy data for testing
└── scripts/            # Utility scripts for database management
```

## Migration System

### Naming Convention
Migration files follow this pattern:
```
V{version}_{timestamp}_{description}.sql
```

Example: `V001_20251109_create_core_tenant_tables.sql`

### Migration Order
Migrations are executed in version order. Each migration should:
- Be idempotent where possible (use `IF NOT EXISTS`)
- Have a clear rollback strategy
- Be atomic (single transaction if possible)
- Include comments explaining the changes
- **ALWAYS include Row-Level Security (RLS) policies for new tables**

### Running Migrations

```bash
# Run all pending migrations
psql -U postgres -d pos_saas -f scripts/run_migrations.sh

# Run a specific migration
psql -U postgres -d pos_saas -f migrations/V001_20251109_create_core_tenant_tables.sql
```

## Database Design Principles

### Multi-Tenancy (SAAS)
- **Schema-per-tenant approach**: Each organization gets its own schema
- **Shared tables**: Core system tables in `public` schema
- **Tenant isolation**: Row-level security and schema separation

### Flexibility & Maintainability
- **Modular design**: Each module/feature in separate migration
- **Soft deletes**: Use `deleted_at` instead of hard deletes
- **Audit trails**: Track who created/modified records
- **Extensible**: JSON columns for custom fields
- **Versioning**: Track schema versions

### Row-Level Security (RLS) - MANDATORY

**CRITICAL**: All tables MUST have Row-Level Security (RLS) enabled for data isolation.

#### RLS Requirements
- ✅ Every table with `organization_id` MUST have RLS policies
- ✅ All migrations creating new tables MUST include RLS policies
- ✅ Use helper functions: `current_user_organization_id()`, `is_super_admin()`
- ✅ Set user context before queries: `SELECT set_user_context(user_id, org_id, is_admin)`

#### Quick RLS Template
```sql
-- Enable RLS on your table
ALTER TABLE your_table ENABLE ROW LEVEL SECURITY;

-- Super admin bypass (ALWAYS include first)
CREATE POLICY your_table_super_admin_all
    ON your_table FOR ALL TO PUBLIC
    USING (is_super_admin());

-- Organization-scoped policies
CREATE POLICY your_table_select_own_org
    ON your_table FOR SELECT TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY your_table_insert_own_org
    ON your_table FOR INSERT TO PUBLIC
    WITH CHECK (organization_id = current_user_organization_id());

CREATE POLICY your_table_update_own_org
    ON your_table FOR UPDATE TO PUBLIC
    USING (organization_id = current_user_organization_id());

CREATE POLICY your_table_delete_own_org
    ON your_table FOR DELETE TO PUBLIC
    USING (organization_id = current_user_organization_id());
```

#### Setting User Context in Application
```sql
-- At the start of each database session:
SELECT set_user_context(
    'user-uuid'::UUID,           -- current user ID
    'organization-uuid'::UUID,   -- user's organization ID
    FALSE                        -- is super admin (TRUE only for system admins)
);

-- Now all queries are automatically filtered by organization
SELECT * FROM products;  -- Only returns user's organization products
```

**📖 For complete RLS documentation, see [schemas/RLS_POLICY_GUIDE.md](schemas/RLS_POLICY_GUIDE.md)**

### Key Tables Structure

#### Core System Tables (public schema)
- `organizations` - Tenant/company information
- `users` - System users with organization mapping
- `roles` - Role definitions
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mapping
- `user_roles` - User-role mapping

#### POS Tables (per-tenant schema)
- `products` - Product catalog
- `categories` - Product categories
- `customers` - Customer information
- `sales` - Sales transactions
- `sale_items` - Line items for sales
- `payments` - Payment records
- `inventory` - Stock management
- `inventory_transactions` - Stock movement tracking

#### E-Invoicing Tables (V014+)
- `e_invoicing_documents` - Central repository for all e-invoicing documents (ZATCA, ETA, etc.)
- `e_invoicing_document_events` - Complete audit trail of e-invoicing events

### E-Invoicing Integration (ZATCA & ETA)

The system includes comprehensive support for electronic invoicing compliance with multiple tax authorities:

#### Supported Authorities
- **ZATCA** (Saudi Arabia) - Standard & Simplified invoices with QR codes, hash chaining, and cryptographic stamps
- **ETA** (Egypt) - JSON/XML format with CAdES-BES digital signatures
- **Extensible** - Architecture supports additional authorities

#### Key Features

**Authority-Agnostic Design**
- Polymorphic source references (works with `sales`, `customer_invoices`, or any invoice table)
- Flexible metadata storage for authority-specific requirements
- Unified status workflow across all authorities

**ZATCA Features**
- Invoice counter value (ICV) tracking
- Hash chaining for blockchain-style verification
- Cryptographic stamps and QR code generation
- PIH (Previous Invoice Hash) compliance
- Support for both Standard (B2B) and Simplified (B2C) invoices
- Excise tax handling for tobacco, soft drinks, etc.

**ETA Features**
- Document type versioning
- CAdES-BES digital signature support
- Long ID assignment after acceptance
- Receiver type classification (Business, Person, Foreigner)
- Egyptian GS1/EGS item code tracking

**Extended Tables**

The following tables have been extended with e-invoicing fields:

**Customers Table**
```sql
-- Tax Registration
tax_registration_number
tax_registration_type
tax_registration_country
is_tax_registered

-- ZATCA-specific address fields
zatca_building_number
zatca_street_name
zatca_district
zatca_city_name
zatca_postal_zone

-- ETA-specific fields
eta_receiver_type
eta_receiver_id
eta_governorate
eta_region_city

-- Flexible address storage
structured_address JSONB
```

**Products Table**
```sql
-- Standard item codes
standard_item_code
standard_item_code_type  -- GS1, EGS, GTIN, etc.
harmonized_system_code

-- Standard UOM codes
standard_uom_code  -- UN/ECE Recommendation 20
standard_uom_name

-- Tax classification
tax_category_code  -- S, Z, E, O
default_vat_rate
tax_exemption_reason

-- Authority-specific fields
zatca_is_excise_taxable
eta_gs1_code
```

**Sales Table**
```sql
-- E-invoicing status tracking
is_e_invoice_required
e_invoice_status
e_invoice_document_id
e_invoicing_metadata JSONB
```

#### Helper Views

Five helper views are available for querying e-invoicing data:

- `view_e_invoices_with_source` - E-invoices with resolved source details
- `view_e_invoice_event_history` - Complete event history with user details
- `view_zatca_invoices` - ZATCA-specific invoices with ZATCA fields
- `view_eta_invoices` - ETA-specific invoices with ETA fields
- `view_failed_e_invoices` - Failed invoices requiring attention/retry

#### Usage Example

```sql
-- Create an e-invoicing document for a sale
INSERT INTO e_invoicing_documents (
    organization_id,
    source_table,
    source_id,
    authority,
    country_code,
    document_type,
    document_number,
    status
) VALUES (
    'org-uuid',
    'sales',
    'sale-uuid',
    'ZATCA',
    'SAU',
    'standard',
    'INV-2024-001',
    'pending'
);

-- Track submission event
INSERT INTO e_invoicing_document_events (
    organization_id,
    e_invoicing_document_id,
    event_type,
    new_status,
    event_description
) VALUES (
    'org-uuid',
    'doc-uuid',
    'submitted',
    'submitted',
    'Document submitted to ZATCA'
);

-- Query all pending invoices
SELECT * FROM view_e_invoices_with_source
WHERE status = 'pending' AND authority = 'ZATCA';

-- Find failed invoices needing retry
SELECT * FROM view_failed_e_invoices
WHERE should_retry = true;
```

#### Implementation Notes

- All changes are non-destructive and additive
- Uses `IF NOT EXISTS` for idempotent migrations
- RLS policies applied for multi-tenant isolation
- Extensive seed data provided for testing (see `007_seed_e_invoicing_data.sql`)
- Status workflow: draft → pending → submitted → accepted/rejected/cancelled

## Common Fields

All tables include these standard fields:
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
deleted_at TIMESTAMP WITH TIME ZONE -- For soft deletes
```

## Seed Data

Seed data files are located in `seed_data/` and include:
- Sample organizations
- Default users and roles
- Sample products and categories
- Test transaction data

## Scripts

Utility scripts in `scripts/`:
- `init_database.sql` - Initialize database and extensions
- `run_migrations.sh` - Execute all migrations
- `reset_database.sh` - Reset database (development only)
- `backup_database.sh` - Backup utilities

## Development Workflow

1. **Create Migration**: Add new migration file with incremented version
2. **Test Migration**: Run in development environment
3. **Add Seed Data**: Create corresponding seed data if needed
4. **Document Changes**: Update this README if needed
5. **Commit**: Commit migration files to version control

## Production Considerations

- Always backup before running migrations
- Test migrations in staging first
- Use transactions for data migrations
- Monitor long-running migrations
- Keep rollback scripts ready

## Connection String Example

```bash
postgresql://username:password@localhost:5432/pos_saas
```

## Environment Variables

```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pos_saas
DB_USER=postgres
DB_PASSWORD=your_password
```
