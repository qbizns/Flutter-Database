package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	// CSRFTokenLength is the length of CSRF tokens in bytes
	CSRFTokenLength = 32

	// CSRFTokenHeader is the header name for CSRF token
	CSRFTokenHeader = "X-CSRF-Token"

	// CSRFCookieName is the name of the CSRF cookie
	CSRFCookieName = "csrf_token"

	// CSRFTokenLifetime is how long a CSRF token is valid
	CSRFTokenLifetime = 24 * time.Hour
)

var (
	// ErrNoCSRFToken is returned when no CSRF token is provided
	ErrNoCSRFToken = errors.New("CSRF token not provided")

	// ErrInvalidCSRFToken is returned when CSRF token validation fails
	ErrInvalidCSRFToken = errors.New("invalid CSRF token")
)

// CSRFConfig holds configuration for CSRF protection
type CSRFConfig struct {
	// TokenLength is the length of CSRF tokens in bytes (default: 32)
	TokenLength int

	// TokenLifetime is how long tokens are valid (default: 24h)
	TokenLifetime time.Duration

	// CookieName is the name of the CSRF cookie (default: "csrf_token")
	CookieName string

	// CookiePath is the path for the CSRF cookie (default: "/")
	CookiePath string

	// CookieDomain is the domain for the CSRF cookie
	CookieDomain string

	// CookieSecure sets the Secure flag on the cookie (default: true in production)
	CookieSecure bool

	// CookieSameSite sets the SameSite attribute (default: Strict)
	CookieSameSite http.SameSite

	// SkipPaths is a list of paths to skip CSRF validation
	SkipPaths []string

	// SafeMethods is a list of HTTP methods that don't require CSRF validation
	// Default: GET, HEAD, OPTIONS, TRACE
	SafeMethods []string
}

// DefaultCSRFConfig returns a default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:    CSRFTokenLength,
		TokenLifetime:  CSRFTokenLifetime,
		CookieName:     CSRFCookieName,
		CookiePath:     "/",
		CookieSecure:   true,
		CookieSameSite: http.SameSiteLaxMode, // Lax allows cross-origin GET requests
		SafeMethods:    []string{"GET", "HEAD", "OPTIONS", "TRACE"},
		SkipPaths:      []string{"/health", "/metrics", "/auth/login", "/auth/register"},
	}
}

// csrfToken represents a CSRF token with metadata
type csrfToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFProtection provides CSRF protection middleware
type CSRFProtection struct {
	config CSRFConfig
	tokens sync.Map // map[string]*csrfToken for token storage
	mu     sync.RWMutex
}

// NewCSRFProtection creates a new CSRF protection middleware
func NewCSRFProtection(config CSRFConfig) *CSRFProtection {
	// Set defaults
	if config.TokenLength == 0 {
		config.TokenLength = CSRFTokenLength
	}
	if config.TokenLifetime == 0 {
		config.TokenLifetime = CSRFTokenLifetime
	}
	if config.CookieName == "" {
		config.CookieName = CSRFCookieName
	}
	if config.CookiePath == "" {
		config.CookiePath = "/"
	}
	if len(config.SafeMethods) == 0 {
		config.SafeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
	}

	csrf := &CSRFProtection{
		config: config,
	}

	// Start cleanup goroutine
	go csrf.cleanupExpiredTokens()

	return csrf
}

// Middleware returns the CSRF protection middleware
func (c *CSRFProtection) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF validation for configured paths
		if c.shouldSkipPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Skip CSRF validation for safe methods
		if c.isSafeMethod(r.Method) {
			// For safe methods, generate and send CSRF token
			token, err := c.getOrCreateToken(r)
			if err != nil {
				http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
				return
			}
			c.setCSRFCookie(w, token)
			w.Header().Set(CSRFTokenHeader, token)
			next.ServeHTTP(w, r)
			return
		}

		// For unsafe methods, validate CSRF token
		if err := c.validateToken(r); err != nil {
			http.Error(w, "CSRF validation failed: "+err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// generateToken generates a new random CSRF token
func (c *CSRFProtection) generateToken() (string, error) {
	bytes := make([]byte, c.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// getOrCreateToken gets an existing valid token or creates a new one
func (c *CSRFProtection) getOrCreateToken(r *http.Request) (string, error) {
	// Try to get token from cookie
	cookie, err := r.Cookie(c.config.CookieName)
	if err == nil && cookie.Value != "" {
		// Validate existing token
		if tokenData, ok := c.tokens.Load(cookie.Value); ok {
			token := tokenData.(*csrfToken)
			if time.Now().Before(token.ExpiresAt) {
				return token.Token, nil
			}
			// Token expired, remove it
			c.tokens.Delete(cookie.Value)
		}
	}

	// Generate new token
	tokenStr, err := c.generateToken()
	if err != nil {
		return "", err
	}

	token := &csrfToken{
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(c.config.TokenLifetime),
	}
	c.tokens.Store(tokenStr, token)

	return tokenStr, nil
}

// validateToken validates the CSRF token from the request
func (c *CSRFProtection) validateToken(r *http.Request) error {
	// Get token from header
	headerToken := r.Header.Get(CSRFTokenHeader)
	if headerToken == "" {
		// Try to get from form data
		if err := r.ParseForm(); err == nil {
			headerToken = r.FormValue("csrf_token")
		}
	}

	if headerToken == "" {
		return ErrNoCSRFToken
	}

	// Get token from cookie
	cookie, err := r.Cookie(c.config.CookieName)
	if err != nil || cookie.Value == "" {
		return ErrNoCSRFToken
	}

	// Double-submit cookie pattern: compare header token with cookie token
	if !secureCompare(headerToken, cookie.Value) {
		return ErrInvalidCSRFToken
	}

	// Verify token exists in our store and is not expired
	tokenData, ok := c.tokens.Load(headerToken)
	if !ok {
		return ErrInvalidCSRFToken
	}

	token := tokenData.(*csrfToken)
	if time.Now().After(token.ExpiresAt) {
		c.tokens.Delete(headerToken)
		return ErrInvalidCSRFToken
	}

	return nil
}

// setCSRFCookie sets the CSRF token cookie
func (c *CSRFProtection) setCSRFCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     c.config.CookieName,
		Value:    token,
		Path:     c.config.CookiePath,
		Domain:   c.config.CookieDomain,
		MaxAge:   int(c.config.TokenLifetime.Seconds()),
		Secure:   c.config.CookieSecure,
		HttpOnly: true, // Prevent JavaScript access
		SameSite: c.config.CookieSameSite,
	}
	http.SetCookie(w, cookie)
}

// shouldSkipPath checks if a path should skip CSRF validation
func (c *CSRFProtection) shouldSkipPath(path string) bool {
	for _, skipPath := range c.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// isSafeMethod checks if an HTTP method is considered safe
func (c *CSRFProtection) isSafeMethod(method string) bool {
	for _, safeMethod := range c.config.SafeMethods {
		if method == safeMethod {
			return true
		}
	}
	return false
}

// cleanupExpiredTokens periodically removes expired tokens
func (c *CSRFProtection) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		c.tokens.Range(func(key, value interface{}) bool {
			token := value.(*csrfToken)
			if now.After(token.ExpiresAt) {
				c.tokens.Delete(key)
			}
			return true
		})
	}
}

// secureCompare performs constant-time comparison of two strings
func secureCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// GetCSRFToken extracts the CSRF token from a request (for testing purposes)
func GetCSRFToken(r *http.Request) string {
	return r.Header.Get(CSRFTokenHeader)
}

// SetCSRFToken sets the CSRF token in the response header
func SetCSRFToken(w http.ResponseWriter, token string) {
	w.Header().Set(CSRFTokenHeader, token)
}
