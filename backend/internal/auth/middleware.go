package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
)

// Claims represents JWT claims
type Claims struct {
	UserID         string   `json:"user_id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	Roles          []string `json:"roles"`
	jwt.RegisteredClaims
}

// Middleware provides authentication middleware
type Middleware struct {
	jwtSecret string
}

// NewMiddleware creates a new auth middleware
func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{
		jwtSecret: jwtSecret,
	}
}

// Authenticate validates JWT token and sets user context
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.unauthorized(w, "missing authorization header")
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.unauthorized(w, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.Unauthorized("invalid signing method")
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			m.unauthorized(w, "invalid or expired token")
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			m.unauthorized(w, "invalid token claims")
			return
		}

		// Parse UUIDs
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			m.unauthorized(w, "invalid user_id in token")
			return
		}

		orgID, err := uuid.Parse(claims.OrganizationID)
		if err != nil {
			m.unauthorized(w, "invalid organization_id in token")
			return
		}

		// Add to context
		ctx := r.Context()
		ctx = appctx.WithUserID(ctx, userID)
		ctx = appctx.WithOrganizationID(ctx, orgID)
		ctx = appctx.WithRoles(ctx, claims.Roles)

		// Continue with enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthenticate validates JWT token if present
func (m *Middleware) OptionalAuthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Try to authenticate, but don't fail if invalid
		m.Authenticate(next).ServeHTTP(w, r)
	})
}

// RequireRole checks if user has required role
func (m *Middleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !appctx.HasRole(r.Context(), role) {
				m.forbidden(w, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireScope checks if API key has required scope
func (m *Middleware) RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !appctx.HasScope(r.Context(), scope) {
				m.forbidden(w, "insufficient scope")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Middleware) unauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "` + message + `"}`))
}

func (m *Middleware) forbidden(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error": "` + message + `"}`))
}
