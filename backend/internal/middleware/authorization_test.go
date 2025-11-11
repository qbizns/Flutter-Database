package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/domain/auth"
	"github.com/your-org/pos-backend/internal/logging"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
)

// Mock repositories
type mockRoleRepo struct {
	roles map[uuid.UUID]*auth.Role
}

func (m *mockRoleRepo) Create(ctx context.Context, role *auth.Role) error {
	m.roles[role.ID] = role
	return nil
}

func (m *mockRoleRepo) Get(ctx context.Context, id uuid.UUID) (*auth.Role, error) {
	if role, ok := m.roles[id]; ok {
		return role, nil
	}
	return nil, nil
}

func (m *mockRoleRepo) GetBySlug(ctx context.Context, slug string) (*auth.Role, error) {
	return nil, nil
}

func (m *mockRoleRepo) List(ctx context.Context, filters auth.RoleFilters) ([]auth.Role, error) {
	return nil, nil
}

func (m *mockRoleRepo) Count(ctx context.Context, filters auth.RoleFilters) (int64, error) {
	return 0, nil
}

func (m *mockRoleRepo) GetDefaultRole(ctx context.Context, orgID uuid.UUID) (*auth.Role, error) {
	return nil, nil
}

func (m *mockRoleRepo) Update(ctx context.Context, role *auth.Role) error {
	return nil
}

func (m *mockRoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockPermissionRepo struct {
	permissions map[uuid.UUID]*auth.Permission
}

func (m *mockPermissionRepo) Create(ctx context.Context, perm *auth.Permission) error {
	m.permissions[perm.ID] = perm
	return nil
}

func (m *mockPermissionRepo) Get(ctx context.Context, id uuid.UUID) (*auth.Permission, error) {
	if perm, ok := m.permissions[id]; ok {
		return perm, nil
	}
	return nil, nil
}

func (m *mockPermissionRepo) GetBySlug(ctx context.Context, slug string) (*auth.Permission, error) {
	return nil, nil
}

func (m *mockPermissionRepo) List(ctx context.Context, filters auth.PermissionFilters) ([]auth.Permission, error) {
	return nil, nil
}

func (m *mockPermissionRepo) Count(ctx context.Context, filters auth.PermissionFilters) (int64, error) {
	return 0, nil
}

func (m *mockPermissionRepo) Update(ctx context.Context, perm *auth.Permission) error {
	return nil
}

func (m *mockPermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type mockUserRoleRepo struct {
	userRoles map[uuid.UUID][]auth.Role
}

func (m *mockUserRoleRepo) Assign(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	return nil
}

func (m *mockUserRoleRepo) Revoke(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return nil
}

func (m *mockUserRoleRepo) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]auth.Role, error) {
	if roles, ok := m.userRoles[userID]; ok {
		return roles, nil
	}
	return []auth.Role{}, nil
}

func (m *mockUserRoleRepo) GetRoleUsers(ctx context.Context, roleID uuid.UUID) ([]auth.User, error) {
	return nil, nil
}

func (m *mockUserRoleRepo) HasRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockUserRoleRepo) DeleteAllUserRoles(ctx context.Context, userID uuid.UUID) error {
	return nil
}

type mockRolePermissionRepo struct {
	rolePermissions map[uuid.UUID][]auth.Permission
}

func (m *mockRolePermissionRepo) Assign(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	return nil
}

func (m *mockRolePermissionRepo) Revoke(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	return nil
}

func (m *mockRolePermissionRepo) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	if perms, ok := m.rolePermissions[roleID]; ok {
		return perms, nil
	}
	return []auth.Permission{}, nil
}

func (m *mockRolePermissionRepo) GetPermissionRoles(ctx context.Context, permissionID uuid.UUID) ([]auth.Role, error) {
	return nil, nil
}

func (m *mockRolePermissionRepo) HasPermission(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *mockRolePermissionRepo) DeleteAllRolePermissions(ctx context.Context, roleID uuid.UUID) error {
	return nil
}

func setupAuthMiddleware() (*AuthorizationMiddleware, *mockUserRoleRepo, *mockRolePermissionRepo, *mockPermissionRepo) {
	logger, _ := logging.NewLogger("error", "json")

	roleRepo := &mockRoleRepo{roles: make(map[uuid.UUID]*auth.Role)}
	permRepo := &mockPermissionRepo{permissions: make(map[uuid.UUID]*auth.Permission)}
	userRoleRepo := &mockUserRoleRepo{userRoles: make(map[uuid.UUID][]auth.Role)}
	rolePermRepo := &mockRolePermissionRepo{rolePermissions: make(map[uuid.UUID][]auth.Permission)}

	config := DefaultAuthorizationConfig()
	config.CacheExpiration = 100 * time.Millisecond // Short expiration for tests

	am := NewAuthorizationMiddleware(roleRepo, permRepo, userRoleRepo, rolePermRepo, logger, config)

	return am, userRoleRepo, rolePermRepo, permRepo
}

func TestRequirePermission_WithPermission(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	// Setup permission
	perm := &auth.Permission{
		ID:       permID,
		Resource: "customer",
		Action:   "create",
	}
	permRepo.permissions[permID] = perm

	// Setup user role
	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Manager"},
	}

	// Setup role permission
	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	handler := am.RequirePermission("customer", "create")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("POST", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRequirePermission_WithoutPermission(t *testing.T) {
	am, userRoleRepo, _, _ := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()

	// User has role but no permissions
	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "User"},
	}

	handler := am.RequirePermission("customer", "create")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestRequirePermission_NoUserID(t *testing.T) {
	am, _, _, _ := setupAuthMiddleware()

	handler := am.RequirePermission("customer", "create")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/customers", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestHasPermission_WildcardAction(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	// Setup wildcard permission (all actions on customer resource)
	perm := &auth.Permission{
		ID:       permID,
		Resource: "customer",
		Action:   "*",
	}
	permRepo.permissions[permID] = perm

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Manager"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	// Should have permission for any action on customer
	hasCreate, err := am.HasPermission(context.Background(), userID, "customer", "create")
	if err != nil || !hasCreate {
		t.Error("should have customer:create permission with wildcard action")
	}

	hasDelete, err := am.HasPermission(context.Background(), userID, "customer", "delete")
	if err != nil || !hasDelete {
		t.Error("should have customer:delete permission with wildcard action")
	}
}

func TestHasPermission_WildcardAll(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	// Setup wildcard permission (*:* = superadmin)
	perm := &auth.Permission{
		ID:       permID,
		Resource: "*",
		Action:   "*",
	}
	permRepo.permissions[permID] = perm

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Superadmin"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	// Should have permission for anything
	hasPermission, err := am.HasPermission(context.Background(), userID, "anything", "do")
	if err != nil || !hasPermission {
		t.Error("should have permission for any resource and action with *:*")
	}
}

func TestRequireAnyPermission_HasOne(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	// Setup one permission
	perm := &auth.Permission{
		ID:       permID,
		Resource: "customer",
		Action:   "read",
	}
	permRepo.permissions[permID] = perm

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Viewer"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	handler := am.RequireAnyPermission(
		Permission{"customer", "create"},
		Permission{"customer", "read"}, // Has this one
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequireAnyPermission_HasNone(t *testing.T) {
	am, userRoleRepo, _, _ := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "User"},
	}

	handler := am.RequireAnyPermission(
		Permission{"customer", "create"},
		Permission{"customer", "delete"},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestRequireAllPermissions_HasAll(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID1 := uuid.New()
	permID2 := uuid.New()

	// Setup two permissions
	perm1 := &auth.Permission{
		ID:       permID1,
		Resource: "customer",
		Action:   "create",
	}
	perm2 := &auth.Permission{
		ID:       permID2,
		Resource: "customer",
		Action:   "read",
	}
	permRepo.permissions[permID1] = perm1
	permRepo.permissions[permID2] = perm2

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Manager"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm1, *perm2}

	handler := am.RequireAllPermissions(
		Permission{"customer", "create"},
		Permission{"customer", "read"},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRequireAllPermissions_MissingOne(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	// Setup only one permission
	perm := &auth.Permission{
		ID:       permID,
		Resource: "customer",
		Action:   "read",
	}
	permRepo.permissions[permID] = perm

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Viewer"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	handler := am.RequireAllPermissions(
		Permission{"customer", "create"}, // Missing this
		Permission{"customer", "read"},
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/customers", nil)
	ctx := appctx.WithUserID(req.Context(), userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestPermissionCaching(t *testing.T) {
	am, userRoleRepo, rolePermRepo, permRepo := setupAuthMiddleware()

	userID := uuid.New()
	roleID := uuid.New()
	permID := uuid.New()

	perm := &auth.Permission{
		ID:       permID,
		Resource: "customer",
		Action:   "read",
	}
	permRepo.permissions[permID] = perm

	userRoleRepo.userRoles[userID] = []auth.Role{
		{ID: roleID, Name: "Viewer"},
	}

	rolePermRepo.rolePermissions[roleID] = []auth.Permission{*perm}

	// First call - should cache
	hasPermission1, err := am.HasPermission(context.Background(), userID, "customer", "read")
	if err != nil || !hasPermission1 {
		t.Fatal("should have permission")
	}

	// Second call - should use cache
	hasPermission2, err := am.HasPermission(context.Background(), userID, "customer", "read")
	if err != nil || !hasPermission2 {
		t.Fatal("should have permission from cache")
	}

	// Invalidate cache
	am.InvalidateUserCache(userID)

	// Third call - should fetch again
	hasPermission3, err := am.HasPermission(context.Background(), userID, "customer", "read")
	if err != nil || !hasPermission3 {
		t.Fatal("should have permission after cache invalidation")
	}
}

func TestParsePermissionString(t *testing.T) {
	tests := []struct {
		input          string
		expectedRes    string
		expectedAction string
		expectError    bool
	}{
		{"customer:create", "customer", "create", false},
		{"user:read", "user", "read", false},
		{"*:*", "*", "*", false},
		{"invalid", "", "", true},
		{"too:many:parts", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res, action, err := ParsePermissionString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if res != tt.expectedRes {
				t.Errorf("expected resource %s, got %s", tt.expectedRes, res)
			}

			if action != tt.expectedAction {
				t.Errorf("expected action %s, got %s", tt.expectedAction, action)
			}
		})
	}
}
