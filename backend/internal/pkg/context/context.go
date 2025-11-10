package context

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	orgIDKey    contextKey = "organization_id"
	requestIDKey contextKey = "request_id"
	rolesKey    contextKey = "roles"
	scopesKey   contextKey = "scopes"
)

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

// MustGetUserID retrieves user ID from context or panics
func MustGetUserID(ctx context.Context) uuid.UUID {
	userID, ok := GetUserID(ctx)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}

// WithOrganizationID adds organization ID to context
func WithOrganizationID(ctx context.Context, orgID uuid.UUID) context.Context {
	return context.WithValue(ctx, orgIDKey, orgID)
}

// GetOrganizationID retrieves organization ID from context
func GetOrganizationID(ctx context.Context) (uuid.UUID, bool) {
	orgID, ok := ctx.Value(orgIDKey).(uuid.UUID)
	return orgID, ok
}

// MustGetOrganizationID retrieves organization ID from context or panics
func MustGetOrganizationID(ctx context.Context) uuid.UUID {
	orgID, ok := GetOrganizationID(ctx)
	if !ok {
		panic("organization_id not found in context")
	}
	return orgID
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}

// WithRoles adds user roles to context
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey, roles)
}

// GetRoles retrieves user roles from context
func GetRoles(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(rolesKey).([]string)
	return roles, ok
}

// WithScopes adds API scopes to context
func WithScopes(ctx context.Context, scopes []string) context.Context {
	return context.WithValue(ctx, scopesKey, scopes)
}

// GetScopes retrieves API scopes from context
func GetScopes(ctx context.Context) ([]string, bool) {
	scopes, ok := ctx.Value(scopesKey).([]string)
	return scopes, ok
}

// HasScope checks if context has a specific scope
func HasScope(ctx context.Context, scope string) bool {
	scopes, ok := GetScopes(ctx)
	if !ok {
		return false
	}
	for _, s := range scopes {
		if s == scope || s == "admin" {
			return true
		}
	}
	return false
}

// HasRole checks if context has a specific role
func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetRoles(ctx)
	if !ok {
		return false
	}
	for _, r := range roles {
		if r == role || r == "admin" {
			return true
		}
	}
	return false
}
