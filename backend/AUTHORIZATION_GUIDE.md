# Authorization & RBAC Guide

## Overview

This backend implements a comprehensive Role-Based Access Control (RBAC) system with fine-grained permissions. The system supports:

- Resource-based permissions (Resource + Action pattern)
- Role-permission mappings
- User-role assignments
- Wildcard permissions for administrative access
- Permission caching for performance
- Multi-tenancy support (organization-scoped)

## Architecture

### Core Components

1. **Permission**: Defines what action can be performed on what resource
2. **Role**: Collection of permissions
3. **User**: System actor who can be assigned roles
4. **Organization**: Multi-tenancy boundary

### Database Schema

```
users
  └── user_roles ──> roles
                        └──> role_permissions ──> permissions
```

## Permission Model

### Permission Structure

```go
type Permission struct {
    ID          uuid.UUID
    Name        string          // "Create Customer"
    Slug        string          // "customer_create"
    Description string
    Resource    string          // "customer"
    Action      string          // "create"
    Category    string          // "sales"
}
```

### Resource + Action Pattern

Permissions follow the format: `resource:action`

**Resources:**
- `user` - User management
- `role` - Role management
- `permission` - Permission management
- `organization` - Organization settings
- `customer` - Customer records
- `supplier` - Supplier records
- `product` - Product catalog
- `category` - Product categories
- `location` - Store locations
- `sale` - Sales transactions
- `account` - Chart of accounts
- `journal_entry` - Journal entries
- `fiscal_year` - Fiscal year management
- `invoice` - Customer invoices
- `payment` - Payments (AR/AP)
- `report` - Financial reports

**Actions:**
- `create` - Create new records
- `read` - View/retrieve records
- `update` - Modify existing records
- `delete` - Remove records
- `list` - List/search records
- `post` - Post to general ledger
- `reverse` - Reverse journal entries
- `approve` - Approve transactions
- `reject` - Reject transactions
- `export` - Export data
- `import` - Import data
- `manage` - Full CRUD access

**Wildcards:**
- `customer:*` - All actions on customer resource
- `*:read` - Read action on all resources
- `*:*` - Superadmin (all actions on all resources)

## Role System

### System Roles

Pre-defined roles that ship with the system:

#### 1. **Superadmin** (`superadmin`)
- Permission: `*:*`
- Full system access
- Can manage organizations, users, roles

#### 2. **Administrator** (`admin`)
- Permissions: `*:manage` (all resources except system config)
- Can manage all business data
- Cannot modify system roles

#### 3. **Manager** (`manager`)
- Permissions:
  - `customer:*`, `supplier:*`, `product:*`
  - `sale:*`, `invoice:*`, `payment:*`
  - `journal_entry:create,read,update`
  - `account:read`, `report:*`
- Full access to daily operations
- Cannot post to GL or close periods

#### 4. **Accountant** (`accountant`)
- Permissions:
  - `account:*`, `journal_entry:*`
  - `fiscal_year:*`, `invoice:read`, `payment:read`
  - `report:*`
- Full accounting access
- Can post and reverse entries
- Can close periods

#### 5. **Sales** (`sales`)
- Permissions:
  - `customer:*`, `product:read`, `category:read`
  - `sale:*`, `invoice:create,read,update`
  - `payment:create,read`
- Sales and customer management
- Cannot delete transactions
- Cannot access accounting

#### 6. **Viewer** (`viewer`)
- Permissions: `*:read`, `*:list`
- Read-only access to all modules
- Cannot create, update, or delete
- Useful for external auditors

### Custom Roles

Organizations can create custom roles with specific permission combinations.

```go
// Example: Inventory Manager
customRole := &Role{
    Name:         "Inventory Manager",
    OrganizationID: &orgID,
    IsSystemRole: false,
    Permissions: []Permission{
        {Resource: "product", Action: "manage"},
        {Resource: "category", Action: "manage"},
        {Resource: "location", Action: "manage"},
        {Resource: "supplier", Action: "read"},
        {Resource: "report", Action: "export"},
    },
}
```

## Using Authorization Middleware

### 1. Require Single Permission

```go
// Require customer:create permission
r.With(authz.RequirePermission("customer", "create")).
    Post("/api/customers", handlers.CreateCustomer)

// Using constants
r.With(authz.RequirePermission(ResourceCustomer, ActionCreate)).
    Post("/api/customers", handlers.CreateCustomer)
```

### 2. Require Any Permission (OR logic)

```go
// Require either customer:create OR customer:manage
r.With(authz.RequireAnyPermission(
    Permission{"customer", "create"},
    Permission{"customer", "manage"},
)).Post("/api/customers", handlers.CreateCustomer)
```

### 3. Require All Permissions (AND logic)

```go
// Require both journal_entry:create AND journal_entry:post
r.With(authz.RequireAllPermissions(
    Permission{"journal_entry", "create"},
    Permission{"journal_entry", "post"},
)).Post("/api/journal-entries/{id}/post", handlers.PostJournalEntry)
```

### 4. Resource Owner OR Permission

Useful for endpoints where users can edit their own data OR have admin permissions:

```go
// Users can edit their own profile OR have user:update permission
r.With(authz.ResourceOwnerOrPermission(
    func(r *http.Request) (uuid.UUID, error) {
        return uuid.Parse(chi.URLParam(r, "id"))
    },
    "user",
    "update",
)).Put("/api/users/{id}", handlers.UpdateUser)
```

## Integration Example

### main.go Setup

```go
// Initialize repositories
roleRepo := postgres.NewRoleRepository(db)
permRepo := postgres.NewPermissionRepository(db)
userRoleRepo := postgres.NewUserRoleRepository(db)
rolePermRepo := postgres.NewRolePermissionRepository(db)

// Create authorization middleware
authzConfig := middleware.DefaultAuthorizationConfig()
authz := middleware.NewAuthorizationMiddleware(
    roleRepo,
    permRepo,
    userRoleRepo,
    rolePermRepo,
    logger,
    authzConfig,
)

// Protected routes
r.Route("/api", func(r chi.Router) {
    // Require authentication for all /api routes
    r.Use(authMiddleware.Authenticate)

    // Customer endpoints
    r.With(authz.RequirePermission("customer", "list")).
        Get("/customers", customerHandlers.List)

    r.With(authz.RequirePermission("customer", "create")).
        Post("/customers", customerHandlers.Create)

    r.With(authz.RequirePermission("customer", "read")).
        Get("/customers/{id}", customerHandlers.Get)

    r.With(authz.RequirePermission("customer", "update")).
        Put("/customers/{id}", customerHandlers.Update)

    r.With(authz.RequirePermission("customer", "delete")).
        Delete("/customers/{id}", customerHandlers.Delete)

    // Accounting endpoints
    r.With(authz.RequirePermission("journal_entry", "create")).
        Post("/journal-entries", accountingHandlers.CreateJournalEntry)

    r.With(authz.RequireAllPermissions(
        Permission{"journal_entry", "read"},
        Permission{"journal_entry", "post"},
    )).Post("/journal-entries/{id}/post", accountingHandlers.PostJournalEntry)

    // Report endpoints - require any report permission
    r.With(authz.RequireAnyPermission(
        Permission{"report", "read"},
        Permission{"report", "export"},
    )).Get("/reports/{type}", reportHandlers.GetReport)
})
```

## Permission Caching

The authorization middleware includes built-in caching to minimize database queries:

```go
config := middleware.AuthorizationConfig{
    EnableCache:     true,
    CacheExpiration: 5 * time.Minute, // Cache permissions for 5 minutes
}

authz := middleware.NewAuthorizationMiddleware(
    roleRepo, permRepo, userRoleRepo, rolePermRepo, logger, config,
)
```

### Cache Invalidation

```go
// Invalidate specific user's cache (e.g., after role change)
authz.InvalidateUserCache(userID)

// Invalidate all cached permissions (e.g., after system-wide role update)
authz.InvalidateAllCache()
```

### Cache Considerations

- **TTL**: Default 5 minutes, configurable
- **Automatic Cleanup**: Expired entries removed every 10 minutes
- **Thread-Safe**: Uses `sync.Map` for concurrent access
- **Invalidation**: Auto-invalidate on role/permission changes

## Common Permission Sets

### E-Commerce / Retail POS

```go
// Store Manager
manager_perms := []string{
    "customer:*", "product:*", "category:*",
    "sale:*", "location:read",
    "report:read", "report:export",
}

// Cashier
cashier_perms := []string{
    "customer:create", "customer:read",
    "product:read", "sale:create", "sale:read",
    "payment:create",
}

// Stock Clerk
stock_perms := []string{
    "product:create", "product:read", "product:update",
    "category:read", "supplier:read",
    "location:read",
}
```

### Accounting System

```go
// Chief Accountant
chief_accountant_perms := []string{
    "account:*", "journal_entry:*",
    "fiscal_year:*", "report:*",
    "invoice:*", "payment:*",
}

// Junior Accountant
junior_accountant_perms := []string{
    "account:read", "journal_entry:create", "journal_entry:read",
    "invoice:read", "payment:read",
    "report:read",
}

// Accounts Receivable
ar_perms := []string{
    "customer:read", "invoice:*",
    "payment:create", "payment:read", "payment:update",
    "report:read",
}

// Accounts Payable
ap_perms := []string{
    "supplier:read", "invoice:create", "invoice:read",
    "payment:*", "report:read",
}
```

## Testing Permissions

```go
func TestRequirePermission(t *testing.T) {
    // Setup test user with specific permissions
    userID := uuid.New()
    roleID := uuid.New()

    // Assign role with customer:create permission
    userRoleRepo.Assign(ctx, userID, roleID, nil)
    rolePermRepo.Assign(ctx, roleID, permissionID)

    // Test endpoint with permission check
    req := httptest.NewRequest("POST", "/api/customers", body)
    ctx := appctx.WithUserID(req.Context(), userID)
    req = req.WithContext(ctx)

    handler := authz.RequirePermission("customer", "create")(
        http.HandlerFunc(createCustomerHandler),
    )

    w := httptest.NewRecorder()
    handler.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Security Best Practices

### 1. Principle of Least Privilege
- Assign minimum necessary permissions
- Use specific permissions over wildcards
- Regular permission audits

### 2. Separation of Duties
- Different roles for sensitive operations
- Example: Creator cannot be approver
- Posting to GL requires separate role

### 3. Multi-Factor for Sensitive Operations
```go
// Require additional verification for critical operations
r.With(authz.RequirePermission("fiscal_year", "close")).
    With(mfaMiddleware.Require()).
    Post("/api/fiscal-years/{id}/close", handlers.CloseFiscalYear)
```

### 4. Audit Logging
```go
// Log all permission checks
logger.Info("permission check",
    zap.String("user_id", userID.String()),
    zap.String("resource", resource),
    zap.String("action", action),
    zap.Bool("granted", hasPermission),
)
```

### 5. Organization Isolation
- All resources scoped to organization
- Permissions checked with organization context
- No cross-organization access without explicit permission

## Permission Migration

### Adding New Permissions

```sql
-- Add new permission
INSERT INTO permissions (id, name, slug, resource, action, category)
VALUES (
    gen_random_uuid(),
    'Export Financial Reports',
    'report_export',
    'report',
    'export',
    'reporting'
);

-- Assign to existing roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.slug = 'accountant'
AND p.slug = 'report_export';
```

### Updating Role Permissions

```go
// Programmatic role update
func UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
    // Remove all existing permissions
    if err := rolePermRepo.DeleteAllRolePermissions(ctx, roleID); err != nil {
        return err
    }

    // Assign new permissions
    for _, permID := range permissionIDs {
        if err := rolePermRepo.Assign(ctx, roleID, permID); err != nil {
            return err
        }
    }

    // Invalidate cache for affected users
    authz.InvalidateAllCache()

    return nil
}
```

## Troubleshooting

### Permission Denied Errors

1. **Check user roles:**
   ```sql
   SELECT r.* FROM roles r
   JOIN user_roles ur ON ur.role_id = r.id
   WHERE ur.user_id = $1;
   ```

2. **Check role permissions:**
   ```sql
   SELECT p.* FROM permissions p
   JOIN role_permissions rp ON rp.permission_id = p.id
   WHERE rp.role_id = $1;
   ```

3. **Check cache:**
   ```go
   // Clear cache and retry
   authz.InvalidateUserCache(userID)
   ```

### Performance Issues

1. **Enable caching** (if not already enabled)
2. **Increase cache TTL** for stable permissions
3. **Add database indexes:**
   ```sql
   CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
   CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
   ```

## API Endpoints

### Permission Management

```
GET    /api/permissions          - List all permissions
POST   /api/permissions          - Create new permission (superadmin only)
GET    /api/permissions/{id}     - Get permission details
PUT    /api/permissions/{id}     - Update permission
DELETE /api/permissions/{id}     - Delete permission
```

### Role Management

```
GET    /api/roles                - List all roles
POST   /api/roles                - Create new role
GET    /api/roles/{id}           - Get role details
PUT    /api/roles/{id}           - Update role
DELETE /api/roles/{id}           - Delete role

GET    /api/roles/{id}/permissions        - List role permissions
POST   /api/roles/{id}/permissions        - Assign permission to role
DELETE /api/roles/{id}/permissions/{pid}  - Revoke permission from role
```

### User Role Assignment

```
GET    /api/users/{id}/roles     - List user roles
POST   /api/users/{id}/roles     - Assign role to user
DELETE /api/users/{id}/roles/{rid} - Revoke role from user
```

## Conclusion

This RBAC system provides enterprise-grade authorization with:
- ✅ Fine-grained permissions
- ✅ Flexible role system
- ✅ Performance-optimized caching
- ✅ Multi-tenancy support
- ✅ Wildcard permissions
- ✅ Comprehensive testing

For questions or issues, refer to the inline documentation in `internal/middleware/authorization.go`.
