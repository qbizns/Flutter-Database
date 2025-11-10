# Database Migrations

This directory contains database migration management for the POS backend.

## Directory Structure

```
db/
├── migrations-postgres/     -> Symlink to ../postgres/migrations (23 files)
├── migrations-accounting/   -> Symlink to ../accounting/migrations (13 files)
└── README.md               -> This file
```

## Migration Files

### PostgreSQL Migrations (23 files)
Core POS system tables:
- V001-V004: Core tenant and POS tables with RLS
- V005-V006: Inventory management
- V007: Loyalty program
- V008: Reporting and analytics
- V009-V011: Restaurant operations
- V012-V023: Extended features (delivery, gift cards, taxes, etc.)

### Accounting Migrations (13 files)
Accounting and financial tables:
- V001-V002: Core accounting (chart of accounts, AP/AR, assets)
- V003: Odoo integration extensions
- V004-V008: POS-accounting integration
- V009-V013: Posting engine and validation

**Total: 36 migration files**

## Prerequisites

1. **golang-migrate** must be installed:
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

2. **PostgreSQL** must be running (via Docker or local):
   ```bash
   make docker-up  # Start PostgreSQL + Redis
   ```

3. **DATABASE_URL** environment variable (optional, has default):
   ```bash
   export DATABASE_URL="postgres://user:pass@host:port/dbname?sslmode=disable"
   ```

## Migration Commands

### Apply All Migrations
```bash
make migrate-up
```
Applies all pending migrations from both postgres and accounting folders.

### Check Migration Status
```bash
make migrate-status
```
Shows current migration version for both migration sets.

### Rollback Last Migration
```bash
make migrate-down
```
Rolls back the last migration step from both sets (CAREFUL!).

### Rollback All Migrations (DANGEROUS!)
```bash
make migrate-down-all
```
Drops all tables and resets to clean state. Requires confirmation.

### Force Migration Version
```bash
make migrate-force VERSION=5
```
Forces the migration version to a specific number (use when migration state is inconsistent).

### Go to Specific Version
```bash
make migrate-goto VERSION=10
```
Migrates up or down to a specific version.

## Migration Workflow

### 1. Clean Database Setup
```bash
# Start PostgreSQL
make docker-up

# Check status (should show no migrations)
make migrate-status

# Apply all migrations
make migrate-up

# Verify
make migrate-status
# Should show:
# PostgreSQL migrations status: 23
# Accounting migrations status: 13
```

### 2. Development Workflow
```bash
# After pulling new migrations
make migrate-up

# Check current status
make migrate-status

# If something goes wrong
make migrate-down    # Rollback last step
```

### 3. Testing on Clean DB
```bash
# Reset database
make migrate-down-all  # Type 'yes' to confirm

# Reapply all
make migrate-up

# Run tests
make test
```

## Common Issues

### Issue: "Dirty database version"
**Solution**: Use `make migrate-force VERSION=X` to force the version, then retry.

### Issue: "No change" or migration already applied
**Solution**: Check `make migrate-status` to see current version.

### Issue: Connection refused
**Solution**: Ensure PostgreSQL is running (`make docker-up`).

### Issue: Permission denied
**Solution**: Check DATABASE_URL credentials match your PostgreSQL setup.

## Database Connection Details

Default connection (when using `make docker-up`):
- **Host**: localhost
- **Port**: 5432
- **User**: postgres
- **Password**: postgres
- **Database**: pos_saas
- **SSL Mode**: disable (local development)

## Migration File Naming Convention

Existing files use:
```
VXXX_YYYYMMDD_description.sql
```

Example:
- `V001_20251109_create_core_tenant_tables.sql`
- `V023_20251110_create_backend_infrastructure.sql`

## Security Notes

1. **Never run migrations directly on production** without testing on staging first
2. **Always backup production database** before running migrations
3. **Use SSL in production**: Change `sslmode=disable` to `sslmode=require`
4. **Review migration SQL** before applying to understand changes
5. **Test rollback procedures** in staging environment

## Next Steps

1. Test migrations on clean database (Task 3A.1)
2. Verify all tables are created correctly
3. Update docker-compose.yml to auto-run migrations on startup (optional)
4. Add migration checks to CI/CD pipeline
