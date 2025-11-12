# RLS Security Pattern - CRITICAL UPDATE

## ⚠️ Security Issue Fixed

**Problem**: Using `SET app.current_organization_id` at session level caused context leakage between pooled connections.

**Solution**: Use `SET LOCAL` within transactions to ensure context is transaction-scoped.

## ✅ New Secure Pattern

### Database Layer (db.go)

```go
// ALWAYS use WithOrgContext for RLS-protected operations
err := db.WithOrgContext(ctx, orgID.String(), func(tx pgx.Tx) error {
    // All queries here are RLS-protected
    return repo.Create(ctx, tx, entity)
})
```

### Repository Layer

```go
// Repository methods should accept pgx.Tx
func (r *SaleRepository) Create(ctx context.Context, tx pgx.Tx, sale *domain.Sale) error {
    query := `
        INSERT INTO sales (id, organization_id, customer_id, total_amount, ...)
        VALUES ($1, $2, $3, $4, ...)
    `
    _, err := tx.Exec(ctx, query, sale.ID, sale.OrganizationID, ...)
    return err
}
```

### Handler Layer

```go
func CreateSaleHandler(db *postgres.DB, logger *logging.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        orgID := appctx.MustGetOrganizationID(ctx)

        var req CreateSaleRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            respondError(w, logger, apperrors.BadRequest("invalid request"))
            return
        }

        // Use WithOrgContext - RLS automatically applied
        var saleID uuid.UUID
        err := db.WithOrgContext(ctx, orgID.String(), func(tx pgx.Tx) error {
            repo := postgres.NewSaleRepository(db)
            sale := req.ToSale(orgID)

            if err := repo.Create(ctx, tx, sale); err != nil {
                return err
            }

            saleID = sale.ID
            return nil
        })

        if err != nil {
            respondError(w, logger, err)
            return
        }

        respondJSON(w, http.StatusCreated, map[string]interface{}{
            "id": saleID,
        })
    }
}
```

## 🔒 Security Benefits

1. **Transaction-Scoped**: Context cleared after transaction
2. **Validation**: Organization existence verified before context set
3. **No Leakage**: Connection pooling safe
4. **Type-Safe**: Compile-time checking via function signatures

## 📝 Migration Checklist

- [ ] Update all handlers to use `WithOrgContext`
- [ ] Update all repository methods to accept `pgx.Tx`
- [ ] Remove direct calls to `SetOrganizationContext`
- [ ] Add integration tests for RLS isolation
- [ ] Review all queries for RLS coverage

## ⚠️ DO NOT

```go
// ❌ WRONG - Session-level context
db.Pool.Exec(ctx, "SET app.current_organization_id = $1", orgID)

// ❌ WRONG - No transaction
repo.Create(ctx, nil, sale)

// ❌ WRONG - No organization context
db.WithTx(ctx, func(tx pgx.Tx) error {
    return repo.Create(ctx, tx, sale) // RLS not set!
})
```

## ✅ DO

```go
// ✅ CORRECT
db.WithOrgContext(ctx, orgID.String(), func(tx pgx.Tx) error {
    return repo.Create(ctx, tx, sale)
})
```
