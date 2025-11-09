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
