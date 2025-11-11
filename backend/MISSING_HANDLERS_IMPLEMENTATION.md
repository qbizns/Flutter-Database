# Missing HTTP Handlers Implementation

## Summary

This document describes the HTTP handlers that have been implemented for Phase 3B accounting features but are not yet active in the API due to repository layer compatibility issues.

## Implemented Handler Files

The following handler files have been created in `backend/internal/http/rest/`:

### 1. posting_handlers.go (3 endpoints)
- `POST /posting/post` - Post a document to the general ledger
- `GET /posting/rules` - Get posting rules for document types
- `GET /posting/audit` - Get posting audit trail

### 2. journal_entry_handlers.go (4 endpoints)
- `GET /journal-entries` - List journal entries with filtering
- `POST /journal-entries` - Create a new journal entry
- `GET /journal-entries/{id}` - Get a journal entry by ID
- `POST /journal-entries/{id}/post` - Post a journal entry to GL

### 3. report_handlers.go (3 endpoints)
- `GET /reports/balance-sheet` - Generate balance sheet report
- `GET /reports/income-statement` - Generate income statement
- `GET /reports/trial-balance` - Generate trial balance

### 4. accounting_handlers.go (12 endpoints)

**Chart of Accounts:**
- `GET /chart-of-accounts` - List accounts
- `POST /chart-of-accounts` - Create account
- `GET /chart-of-accounts/{id}` - Get account
- `PATCH /chart-of-accounts/{id}` - Update account
- `DELETE /chart-of-accounts/{id}` - Delete account

**Fiscal Years:**
- `GET /fiscal-years` - List fiscal years
- `POST /fiscal-years` - Create fiscal year
- `GET /fiscal-years/{id}` - Get fiscal year
- `PATCH /fiscal-years/{id}` - Update fiscal year
- `POST /fiscal-years/{id}/close` - Close fiscal year

**Accounting Periods:**
- `GET /accounting-periods` - List accounting periods
- `POST /accounting-periods` - Create accounting period

### 5. receivables_handlers.go (8 endpoints)

**Customer Invoices:**
- `GET /customer-invoices` - List customer invoices
- `POST /customer-invoices` - Create customer invoice
- `GET /customer-invoices/{id}` - Get customer invoice
- `PATCH /customer-invoices/{id}` - Update customer invoice
- `DELETE /customer-invoices/{id}` - Delete customer invoice

**Customer Payments:**
- `GET /customer-payments` - List customer payments
- `POST /customer-payments` - Create customer payment
- `GET /customer-payments/{id}` - Get customer payment

### 6. payables_handlers.go (8 endpoints)

**Vendor Bills:**
- `GET /vendor-bills` - List vendor bills
- `POST /vendor-bills` - Create vendor bill
- `GET /vendor-bills/{id}` - Get vendor bill
- `PATCH /vendor-bills/{id}` - Update vendor bill
- `DELETE /vendor-bills/{id}` - Delete vendor bill

**Vendor Payments:**
- `GET /vendor-payments` - List vendor payments
- `POST /vendor-payments` - Create vendor payment
- `GET /vendor-payments/{id}` - Get vendor payment

## Current Status

**BLOCKED** - Handlers are implemented but not active due to repository layer incompatibility.

## Issue Description

The codebase has an inconsistency in database connection handling:

- **New pattern**: Repositories use `*postgres.DB` (wraps `*pgxpool.Pool`)
  - Example: `ProductRepository`, `PostingConceptRepository`

- **Old pattern**: Some repositories use `*sqlx.DB` or `*sql.DB`
  - `AccountingRepository` uses `*sqlx.DB`
  - `ReceivablesRepository` uses `*sql.DB`
  - `PayablesRepository` uses `*sql.DB`

The handlers are written to use the new `*postgres.DB` pattern, but the repositories they depend on haven't been migrated yet.

## Next Steps to Activate Handlers

### Option 1: Migrate Repositories (Recommended)

Migrate the following repositories to use `*postgres.DB`:

1. **AccountingRepository** (`internal/repository/postgres/accounting_repository.go`)
   - Change: `func NewAccountingRepository(db *sqlx.DB)`
   - To: `func NewAccountingRepository(db *DB)`
   - Update all internal queries to use pgx instead of sqlx

2. **ReceivablesRepository** (`internal/repository/postgres/receivables_repository.go`)
   - Change: `func NewReceivablesRepository(db *sql.DB)`
   - To: `func NewReceivablesRepository(db *DB)`
   - Update all internal queries to use pgx

3. **PayablesRepository** (`internal/repository/postgres/payables_repository.go`)
   - Change: `func NewPayablesRepository(db *sql.DB)`
   - To: `func NewPayablesRepository(db *DB)`
   - Update all internal queries to use pgx

4. **PostingRepository** - Implement missing methods
   - Add repository methods required by `posting.Repository` interface
   - Implement: `LoadDocument`, `GetPostingRules`, `CreateJournalEntry`, etc.

### Option 2: Quick Fix (Temporary)

Add wrapper functions in `internal/repository/postgres/db.go`:

```go
// GetSQLXDB returns a sqlx.DB wrapper (temporary)
func (db *DB) GetSQLXDB() *sqlx.DB {
    // Create sqlx wrapper around pgxpool
    // This is a temporary bridge solution
}

// GetSQLDB returns a sql.DB wrapper (temporary)
func (db *DB) GetSQLDB() *sql.DB {
    // Use stdlib() method if available
}
```

Then update handlers to use these wrappers until full migration is complete.

### Option 3: Add Repository Adapters

Create adapter pattern to bridge old and new repository styles without changing existing code.

## Compilation Errors Resolved

The following issues were found and need to be addressed:

1. ❌ Repository constructor signature mismatches
2. ❌ `FiscalYear.FiscalYear` is `string` not `int` (line 331 in accounting_handlers.go)
3. ❌ Missing `PostingRepository.LoadDocument()` and other posting methods
4. ⚠️  Duplicate middleware declarations (ratelimit.go, security_headers.go)

## Testing Plan

Once repositories are migrated:

1. Unit test each handler
2. Integration test with test database
3. Test RLS (Row Level Security) with different organization contexts
4. Test authorization for all endpoints
5. Verify posting engine workflow end-to-end
6. Test financial reports with sample data

## Estimated Effort

- **Repository Migration**: 4-6 hours
  - AccountingRepository: 2 hours
  - ReceivablesRepository: 1.5 hours
  - PayablesRepository: 1.5 hours
  - PostingRepository methods: 1 hour

- **Testing**: 3-4 hours
  - Unit tests: 1 hour
  - Integration tests: 2-3 hours

- **Total**: 7-10 hours

## Benefits Once Complete

- **38 new API endpoints** for Phase 3B accounting features
- Complete posting engine integration
- Full accounts receivable and payable management
- Financial reporting capabilities
- Chart of accounts management
- Fiscal year and period tracking

## References

- Handler implementations: `backend/internal/http/rest/*_handlers.go`
- Routes (commented): `backend/cmd/api/main.go` lines 177-233
- Repository interfaces: `backend/internal/domain/*/service.go`
- Database wrapper: `backend/internal/repository/postgres/db.go`
