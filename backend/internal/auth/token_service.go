package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/your-org/pos-backend/internal/config"
	apperrors "github.com/your-org/pos-backend/internal/pkg/errors"
)

// TokenService handles JWT token generation and validation
type TokenService struct {
	jwtSecret               string
	accessTokenDuration     time.Duration
	refreshTokenDuration    time.Duration
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// NewTokenService creates a new token service
func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{
		jwtSecret:               cfg.JWT.Secret,
		accessTokenDuration:     cfg.JWT.AccessTokenDuration,
		refreshTokenDuration:    cfg.JWT.RefreshTokenDuration,
	}
}

// GenerateTokenPair generates both access and refresh tokens for a user
func (s *TokenService) GenerateTokenPair(userID, orgID uuid.UUID, email string, roles []string) (*TokenPair, error) {
	// Generate access token
	accessToken, expiresAt, err := s.generateToken(userID, orgID, email, roles, s.accessTokenDuration, false)
	if err != nil {
		return nil, apperrors.InternalError(err)
	}

	// Generate refresh token
	refreshToken, _, err := s.generateToken(userID, orgID, email, roles, s.refreshTokenDuration, true)
	if err != nil {
		return nil, apperrors.InternalError(err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// generateToken generates a JWT token
func (s *TokenService) generateToken(
	userID, orgID uuid.UUID,
	email string,
	roles []string,
	duration time.Duration,
	isRefresh bool,
) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(duration)

	claims := &Claims{
		UserID:         userID.String(),
		OrganizationID: orgID.String(),
		Email:          email,
		Roles:          roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "pos-backend",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.Unauthorized("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, apperrors.Unauthorized("invalid token")
	}

	if !token.Valid {
		return nil, apperrors.Unauthorized("expired or invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, apperrors.Unauthorized("invalid token claims")
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (s *TokenService) RefreshAccessToken(refreshTokenString string) (*TokenPair, error) {
	// Validate the refresh token
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Parse UUIDs from claims
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid user_id in token")
	}

	orgID, err := uuid.Parse(claims.OrganizationID)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid organization_id in token")
	}

	// Generate new token pair
	return s.GenerateTokenPair(userID, orgID, claims.Email, claims.Roles)
}
