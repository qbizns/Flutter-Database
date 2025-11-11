package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/auth"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	"go.uber.org/zap"
)

// PermissionChecker defines the interface for checking user permissions
type PermissionChecker interface {
	HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error)
}

// AuthorizationMiddleware provides permission-based authorization
type AuthorizationMiddleware struct {
	roleRepo       auth.RoleRepository
	permRepo       auth.PermissionRepository
	userRoleRepo   auth.UserRoleRepository
	rolePermRepo   auth.RolePermissionRepository
	logger         *logging.Logger

	// Cache for user permissions (expires after duration)
	permissionCache      sync.Map // map[string]*cachedPermissions
	cacheExpiration      time.Duration
	enableCache          bool
}

// cachedPermissions holds cached permission data
type cachedPermissions struct {
	permissions []auth.Permission
	expiresAt   time.Time
	mu          sync.RWMutex
}

// AuthorizationConfig holds configuration for authorization middleware
type AuthorizationConfig struct {
	EnableCache     bool
	CacheExpiration time.Duration
}

// DefaultAuthorizationConfig returns default configuration
func DefaultAuthorizationConfig() AuthorizationConfig {
	return AuthorizationConfig{
		EnableCache:     true,
		CacheExpiration: 5 * time.Minute, // Cache permissions for 5 minutes
	}
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware(
	roleRepo auth.RoleRepository,
	permRepo auth.PermissionRepository,
	userRoleRepo auth.UserRoleRepository,
	rolePermRepo auth.RolePermissionRepository,
	logger *logging.Logger,
	config AuthorizationConfig,
) *AuthorizationMiddleware {
	am := &AuthorizationMiddleware{
		roleRepo:        roleRepo,
		permRepo:        permRepo,
		userRoleRepo:    userRoleRepo,
		rolePermRepo:    rolePermRepo,
		logger:          logger,
		enableCache:     config.EnableCache,
		cacheExpiration: config.CacheExpiration,
	}

	// Start cache cleanup goroutine
	if am.enableCache {
		go am.cleanupExpiredCache()
	}

	return am
}

// RequirePermission creates middleware that checks for a specific permission
func (am *AuthorizationMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get user ID from context
			userID, err := appctx.GetUserIDOrError(ctx)
			if err != nil {
				am.forbidden(w, "authentication required")
				return
			}

			// Check permission
			hasPermission, err := am.HasPermission(ctx, userID, resource, action)
			if err != nil {
				am.logger.Error("permission check failed",
					zap.Error(err),
					zap.String("user_id", userID.String()),
					zap.String("resource", resource),
					zap.String("action", action),
				)
				am.internalError(w, "permission check failed")
				return
			}

			if !hasPermission {
				am.logger.Warn("permission denied",
					zap.String("user_id", userID.String()),
					zap.String("resource", resource),
					zap.String("action", action),
				)
				am.forbidden(w, fmt.Sprintf("requires permission: %s:%s", resource, action))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission creates middleware that checks for any of the specified permissions
func (am *AuthorizationMiddleware) RequireAnyPermission(permissions ...struct{ Resource, Action string }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, err := appctx.GetUserIDOrError(ctx)
			if err != nil {
				am.forbidden(w, "authentication required")
				return
			}

			// Check if user has any of the required permissions
			for _, perm := range permissions {
				hasPermission, err := am.HasPermission(ctx, userID, perm.Resource, perm.Action)
				if err != nil {
					am.logger.Error("permission check failed",
						zap.Error(err),
						zap.String("user_id", userID.String()),
						zap.String("resource", perm.Resource),
						zap.String("action", perm.Action),
					)
					continue
				}

				if hasPermission {
					next.ServeHTTP(w, r)
					return
				}
			}

			am.forbidden(w, "insufficient permissions")
		})
	}
}

// RequireAllPermissions creates middleware that checks for all specified permissions
func (am *AuthorizationMiddleware) RequireAllPermissions(permissions ...struct{ Resource, Action string }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, err := appctx.GetUserIDOrError(ctx)
			if err != nil {
				am.forbidden(w, "authentication required")
				return
			}

			// Check if user has all required permissions
			for _, perm := range permissions {
				hasPermission, err := am.HasPermission(ctx, userID, perm.Resource, perm.Action)
				if err != nil {
					am.logger.Error("permission check failed",
						zap.Error(err),
						zap.String("user_id", userID.String()),
						zap.String("resource", perm.Resource),
						zap.String("action", perm.Action),
					)
					am.internalError(w, "permission check failed")
					return
				}

				if !hasPermission {
					am.forbidden(w, fmt.Sprintf("requires permission: %s:%s", perm.Resource, perm.Action))
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// ResourceOwnerOrPermission checks if user owns the resource OR has the permission
// Useful for endpoints like "can edit own profile OR has admin permission"
func (am *AuthorizationMiddleware) ResourceOwnerOrPermission(
	getResourceOwnerID func(*http.Request) (uuid.UUID, error),
	resource, action string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, err := appctx.GetUserIDOrError(ctx)
			if err != nil {
				am.forbidden(w, "authentication required")
				return
			}

			// Check if user owns the resource
			resourceOwnerID, err := getResourceOwnerID(r)
			if err == nil && resourceOwnerID == userID {
				next.ServeHTTP(w, r)
				return
			}

			// If not owner, check permission
			hasPermission, err := am.HasPermission(ctx, userID, resource, action)
			if err != nil {
				am.logger.Error("permission check failed", zap.Error(err))
				am.internalError(w, "permission check failed")
				return
			}

			if !hasPermission {
				am.forbidden(w, "access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HasPermission checks if a user has a specific permission
func (am *AuthorizationMiddleware) HasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	// Get user permissions (from cache or database)
	permissions, err := am.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check if user has the required permission
	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
		// Check for wildcard permissions
		if perm.Resource == resource && perm.Action == "*" {
			return true, nil
		}
		if perm.Resource == "*" && perm.Action == action {
			return true, nil
		}
		if perm.Resource == "*" && perm.Action == "*" {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissions retrieves all permissions for a user
func (am *AuthorizationMiddleware) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	// Check cache first if enabled
	if am.enableCache {
		cacheKey := userID.String()
		if cached, ok := am.permissionCache.Load(cacheKey); ok {
			cp := cached.(*cachedPermissions)
			cp.mu.RLock()
			defer cp.mu.RUnlock()

			if time.Now().Before(cp.expiresAt) {
				return cp.permissions, nil
			}
			// Cache expired, will fetch fresh data
		}
	}

	// Get user roles
	roles, err := am.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	if len(roles) == 0 {
		return []auth.Permission{}, nil
	}

	// Collect all permissions from all roles
	permissionMap := make(map[uuid.UUID]auth.Permission)
	for _, role := range roles {
		perms, err := am.rolePermRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}

		for _, perm := range perms {
			permissionMap[perm.ID] = perm
		}
	}

	// Convert map to slice
	permissions := make([]auth.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	// Cache the result if caching is enabled
	if am.enableCache {
		cacheKey := userID.String()
		am.permissionCache.Store(cacheKey, &cachedPermissions{
			permissions: permissions,
			expiresAt:   time.Now().Add(am.cacheExpiration),
		})
	}

	return permissions, nil
}

// InvalidateUserCache invalidates the permission cache for a specific user
func (am *AuthorizationMiddleware) InvalidateUserCache(userID uuid.UUID) {
	if am.enableCache {
		am.permissionCache.Delete(userID.String())
	}
}

// InvalidateAllCache clears the entire permission cache
func (am *AuthorizationMiddleware) InvalidateAllCache() {
	if am.enableCache {
		am.permissionCache.Range(func(key, value interface{}) bool {
			am.permissionCache.Delete(key)
			return true
		})
	}
}

// cleanupExpiredCache periodically removes expired cache entries
func (am *AuthorizationMiddleware) cleanupExpiredCache() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		am.permissionCache.Range(func(key, value interface{}) bool {
			cp := value.(*cachedPermissions)
			cp.mu.RLock()
			expired := now.After(cp.expiresAt)
			cp.mu.RUnlock()

			if expired {
				am.permissionCache.Delete(key)
			}
			return true
		})
	}
}

func (am *AuthorizationMiddleware) forbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "FORBIDDEN",
		"message": message,
	})
}

func (am *AuthorizationMiddleware) internalError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "INTERNAL_ERROR",
		"message": message,
	})
}

// Common permission constants
const (
	// Resources
	ResourceUser         = "user"
	ResourceRole         = "role"
	ResourcePermission   = "permission"
	ResourceOrganization = "organization"
	ResourceCustomer     = "customer"
	ResourceSupplier     = "supplier"
	ResourceProduct      = "product"
	ResourceCategory     = "category"
	ResourceLocation     = "location"
	ResourceSale         = "sale"
	ResourceAccount      = "account"
	ResourceJournalEntry = "journal_entry"
	ResourceFiscalYear   = "fiscal_year"
	ResourceInvoice      = "invoice"
	ResourcePayment      = "payment"
	ResourceReport       = "report"

	// Actions
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionList   = "list"
	ActionPost   = "post"
	ActionReverse = "reverse"
	ActionApprove = "approve"
	ActionReject  = "reject"
	ActionExport  = "export"
	ActionImport  = "import"
	ActionManage  = "manage" // Full CRUD access
)

// Helper type for permission checks
type Permission struct {
	Resource string
	Action   string
}

// Common permission helpers
var (
	// User permissions
	PermUserCreate = Permission{ResourceUser, ActionCreate}
	PermUserRead   = Permission{ResourceUser, ActionRead}
	PermUserUpdate = Permission{ResourceUser, ActionUpdate}
	PermUserDelete = Permission{ResourceUser, ActionDelete}
	PermUserManage = Permission{ResourceUser, ActionManage}

	// Role permissions
	PermRoleCreate = Permission{ResourceRole, ActionCreate}
	PermRoleRead   = Permission{ResourceRole, ActionRead}
	PermRoleUpdate = Permission{ResourceRole, ActionUpdate}
	PermRoleDelete = Permission{ResourceRole, ActionDelete}
	PermRoleManage = Permission{ResourceRole, ActionManage}

	// Customer permissions
	PermCustomerCreate = Permission{ResourceCustomer, ActionCreate}
	PermCustomerRead   = Permission{ResourceCustomer, ActionRead}
	PermCustomerUpdate = Permission{ResourceCustomer, ActionUpdate}
	PermCustomerDelete = Permission{ResourceCustomer, ActionDelete}

	// Sale permissions
	PermSaleCreate = Permission{ResourceSale, ActionCreate}
	PermSaleRead   = Permission{ResourceSale, ActionRead}
	PermSaleUpdate = Permission{ResourceSale, ActionUpdate}
	PermSaleDelete = Permission{ResourceSale, ActionDelete}

	// Journal Entry permissions
	PermJournalEntryCreate  = Permission{ResourceJournalEntry, ActionCreate}
	PermJournalEntryRead    = Permission{ResourceJournalEntry, ActionRead}
	PermJournalEntryUpdate  = Permission{ResourceJournalEntry, ActionUpdate}
	PermJournalEntryDelete  = Permission{ResourceJournalEntry, ActionDelete}
	PermJournalEntryPost    = Permission{ResourceJournalEntry, ActionPost}
	PermJournalEntryReverse = Permission{ResourceJournalEntry, ActionReverse}

	// Report permissions
	PermReportView   = Permission{ResourceReport, ActionRead}
	PermReportExport = Permission{ResourceReport, ActionExport}
)

// ParsePermissionString parses a permission string like "resource:action" into components
func ParsePermissionString(permStr string) (resource, action string, err error) {
	parts := strings.Split(permStr, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid permission format, expected 'resource:action'")
	}
	return parts[0], parts[1], nil
}
