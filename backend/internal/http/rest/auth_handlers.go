package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/auth"
	"github.com/your-org/pos-backend/internal/config"
	domainauth "github.com/your-org/pos-backend/internal/domain/auth"
	"github.com/your-org/pos-backend/internal/logging"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
	"go.uber.org/zap"
)

// LoginHandler handles user login with explicit organization ID
func LoginHandler(
	cfg *config.Config,
	userService *domainauth.UserService,
	tokenService *auth.TokenService,
	logger *logging.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse request with organization ID
		var req struct {
			Email          string    `json:"email" validate:"required,email"`
			Password       string    `json:"password" validate:"required,min=6"`
			OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode login request", zap.Error(err))
			respondError(w, logger, apperrors.BadRequest("invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			logger.Error("login request validation failed", zap.Error(err))
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Get user by email within organization
		user, err := userService.GetByEmail(ctx, req.OrganizationID, req.Email)
		if err != nil {
			if apperrors.IsNotFound(err) {
				logger.Warn("login attempt for non-existent user",
					zap.String("email", req.Email),
					zap.String("org_id", req.OrganizationID.String()))
				respondError(w, logger, apperrors.Unauthorized("invalid credentials"))
				return
			}
			logger.Error("failed to get user", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		// Check if user is locked
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			logger.Warn("login attempt for locked user",
				zap.String("user_id", user.ID.String()),
				zap.Time("locked_until", *user.LockedUntil))
			respondError(w, logger, apperrors.Forbidden("account is temporarily locked due to too many failed login attempts"))
			return
		}

		// Check if user is suspended or inactive
		if user.Status == domainauth.StatusInactive {
			respondError(w, logger, apperrors.Forbidden("account is inactive"))
			return
		}

		// Verify password
		if err := userService.VerifyPassword(req.Password, user.PasswordHash); err != nil {
			// Record failed login attempt
			ip := getClientIP(r)
			if recErr := userService.RecordLoginAttempt(ctx, user.ID, ip, false); recErr != nil {
				logger.Error("failed to record login attempt", zap.Error(recErr))
			}

			logger.Warn("invalid password",
				zap.String("user_id", user.ID.String()),
				zap.String("email", user.Email))
			respondError(w, logger, apperrors.Unauthorized("invalid credentials"))
			return
		}

		// Record successful login
		ip := getClientIP(r)
		if err := userService.RecordLoginAttempt(ctx, user.ID, ip, true); err != nil {
			logger.Error("failed to record successful login", zap.Error(err))
			// Don't fail the login for this
		}

		// Extract role names
		roleNames := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roleNames[i] = role.Name
		}

		// Generate JWT tokens
		tokens, err := tokenService.GenerateTokenPair(
			user.ID,
			user.OrganizationID,
			user.Email,
			roleNames,
		)
		if err != nil {
			logger.Error("failed to generate tokens", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		// Build response using existing LoginResponse type
		response := LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    tokens.ExpiresAt,
			User: UserInfo{
				ID:             user.ID,
				Email:          user.Email,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				OrganizationID: user.OrganizationID,
			},
		}

		logger.Info("user logged in successfully",
			zap.String("user_id", user.ID.String()),
			zap.String("email", user.Email))

		respondJSON(w, http.StatusOK, response)
	}
}

// RegisterHandler handles user registration
func RegisterHandler(
	cfg *config.Config,
	userService *domainauth.UserService,
	tokenService *auth.TokenService,
	logger *logging.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse request using existing RegisterRequest type
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode register request", zap.Error(err))
			respondError(w, logger, apperrors.BadRequest("invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			logger.Error("register request validation failed", zap.Error(err))
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Validate password strength
		if err := validatePasswordStrength(req.Password); err != nil {
			respondError(w, logger, err)
			return
		}

		// Normalize email
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Create user object
		user := &domainauth.User{
			OrganizationID: req.OrganizationID,
			Email:          req.Email,
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Status:         domainauth.StatusActive, // or StatusPending if email verification is required
		}

		// Create user (service handles password hashing and validation)
		if err := userService.Create(ctx, user, req.Password); err != nil {
			logger.Error("failed to create user", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		// Get user with roles (assigned by Create method)
		user, err := userService.Get(ctx, user.ID)
		if err != nil {
			logger.Error("failed to get created user", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		// Extract role names
		roleNames := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roleNames[i] = role.Name
		}

		// Generate JWT tokens
		tokens, err := tokenService.GenerateTokenPair(
			user.ID,
			user.OrganizationID,
			user.Email,
			roleNames,
		)
		if err != nil {
			logger.Error("failed to generate tokens", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		// Build response
		response := LoginResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    tokens.ExpiresAt,
			User: UserInfo{
				ID:             user.ID,
				Email:          user.Email,
				FirstName:      user.FirstName,
				LastName:       user.LastName,
				OrganizationID: user.OrganizationID,
			},
		}

		logger.Info("user registered successfully",
			zap.String("user_id", user.ID.String()),
			zap.String("email", user.Email))

		respondJSON(w, http.StatusCreated, response)
	}
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(
	tokenService *auth.TokenService,
	logger *logging.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		var req struct {
			RefreshToken string `json:"refresh_token" validate:"required"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode refresh token request", zap.Error(err))
			respondError(w, logger, apperrors.BadRequest("invalid request body"))
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			respondError(w, logger, apperrors.ValidationFailed(err.Error()))
			return
		}

		// Refresh tokens
		tokens, err := tokenService.RefreshAccessToken(req.RefreshToken)
		if err != nil {
			logger.Warn("failed to refresh token", zap.Error(err))
			respondError(w, logger, err)
			return
		}

		logger.Info("token refreshed successfully")
		respondJSON(w, http.StatusOK, tokens)
	}
}

// validatePasswordStrength validates password strength
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return apperrors.ValidationFailed("password must be at least 8 characters long")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return apperrors.ValidationFailed("password must contain uppercase, lowercase, and numeric characters")
	}

	return nil
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}
