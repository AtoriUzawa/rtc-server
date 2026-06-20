// Package jwt provides JWT (JSON Web Token) utility functions.
//
// Design goals:
//   - Provide unified token generation and parsing
//   - Support access, refresh, and custom token scenarios
//   - Keep the interface simple, avoid misuse of low-level Claims
//
// Usage:
//  1. Call Init to initialize configuration
//  2. Use GenerateAccessToken / GenerateRefreshToken to generate tokens
//  3. Use ParseToken to parse and validate tokens
//
// Note:
//   - This package uses HS256 (symmetric encryption)
//   - secret must be stored securely (environment variables recommended)
//   - This package does not handle token storage (e.g., Redis blacklist)
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	// secret is the HMAC signing key (must be initialized)
	secret []byte

	// accessExp access token validity duration
	accessExp time.Duration

	// refreshExp refresh token validity duration
	refreshExp time.Duration
)

var (
	// ErrTokenExpired token has expired
	ErrTokenExpired = errors.New("token 已过期")

	// ErrTokenInvalid token is invalid (signature error/format error etc.)
	ErrTokenInvalid = errors.New("token 无效")

	// ErrNotInit Init has not been called
	ErrNotInit = errors.New("jwt 未初始化")
)

// TokenType represents the type of a JWT token (access, refresh, or custom).
type TokenType string

const (
	// TokenTypeAccess is the access token type used for API authentication.
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh is the refresh token type used for obtaining new access tokens.
	TokenTypeRefresh TokenType = "refresh"
	// TokenTypeCustom is the custom token type for special-purpose tokens (e.g. WebSocket handshake).
	TokenTypeCustom TokenType = "custom"
)

// Claims custom JWT claims structure.
//
// Field descriptions:
//   - UserID: unique user identifier
//   - JTI: unique token ID (used for replay prevention / blacklist control)
//   - Params: extension fields (can carry business info like room_id)
//   - RegisteredClaims: standard JWT fields (exp/iat etc.)
//
// Note:
//   - Params is optional, will not appear in token if unset
//   - JTI is unique per generation, recommended to use with Redis for token control
type Claims struct {
	UserID string         `json:"uid"`              // user ID
	JTI    string         `json:"jti"`              // unique token identifier
	Type   TokenType      `json:"type"`             // token type
	Scene  string         `json:"scene"`            // business scenario
	Params map[string]any `json:"params,omitempty"` // extension parameters

	jwt.RegisteredClaims
}

// Init initializes the JWT configuration.
//
// Parameters:
//   - sec: HMAC signing key (recommended >= 32 bytes)
//   - accExp: access token validity duration
//   - refExp: refresh token validity duration
//
// Note:
//   - Must be called before use, otherwise ErrNotInit is returned
//   - Recommended to initialize once at application startup
func Init(sec string, accExp, refExp time.Duration) {
	secret = []byte(sec)
	accessExp = accExp
	refreshExp = refExp
}

// GenerateAccessToken generates an access token.
//
// Parameters:
//   - userID: user ID
//   - scene: business type
//   - params: optional extension parameters (e.g., permissions, business fields)
//
// Returns:
//   - token string
//   - jti (unique token identifier)
//   - error
//
// Use cases:
//   - Issued after user login
//   - Used as API authentication credential
func GenerateAccessToken(userID string, scene string, params ...map[string]any) (string, string, error) {
	return generateToken(userID, accessExp, TokenTypeAccess, scene, params...)
}

// GenerateRefreshToken generates a refresh token.
//
// Parameters:
//   - userID: user ID
//   - scene: business type
//
// Returns:
//   - token string
//   - jti
//   - error
//
// Use cases:
//   - Used to refresh access tokens
//   - Recommended for use only in the /refresh endpoint
func GenerateRefreshToken(userID string, scene string) (string, string, error) {
	return generateToken(userID, refreshExp, TokenTypeRefresh, scene)
}

// GenerateCustomToken generates a custom-purpose token (e.g., WebSocket temporary token).
//
// Parameters:
//   - userID: user ID
//   - exp: expiration duration (recommended short, e.g., 30s ~ 1min)
//   - tp: token type
//   - scene: business type
//   - params: extension parameters (e.g., room_id, scene, etc.)
//
// Returns:
//   - token string
//   - jti
//   - error
//
// Note:
//   - Strongly recommended for one-time or short-term scenarios (e.g., WS handshake)
//   - Recommended to use with JTI for one-time validation (replay attack prevention)
func GenerateCustomToken(userID string, exp time.Duration, scene string, params ...map[string]any) (string, string, error) {
	return generateToken(userID, exp, TokenTypeCustom, scene, params...)
}

// generateToken unified token generation logic (internal use).
//
// Functionality:
//   - Build Claims
//   - Set exp / iat
//   - Generate unique JTI
//   - Sign using HS256
//
// Note:
//   - Not exposed externally, to avoid bypassing security policies
func generateToken(userID string, exp time.Duration, tp TokenType, scene string, params ...map[string]any) (string, string, error) {
	if len(secret) == 0 {
		return "", "", ErrNotInit
	}

	now := time.Now()

	var p map[string]any
	if len(params) > 0 {
		p = params[0]
	}

	claims := &Claims{
		UserID: userID,
		JTI:    uuid.NewString(),
		Type:   tp,
		Scene:  scene,
		Params: p,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	str, err := token.SignedString(secret)
	if err != nil {
		return "", "", err
	}

	return str, claims.JTI, nil
}

// ParseToken parses and validates a JWT token.
//
// Parameters:
//   - tokenStr: JWT string
//
// Returns:
//   - Claims (parsed payload)
//   - error
//
// Error descriptions:
//   - ErrTokenExpired: token has expired
//   - ErrTokenInvalid: token is invalid (signature error/format error etc.)
//   - ErrNotInit: not initialized
//
// Note:
//   - Only validates signature and expiration
//   - Does not include business validation (e.g., whether JTI is revoked)
//   - Recommended to call in middleware
func ParseToken(tokenStr string) (*Claims, error) {
	if len(secret) == 0 {
		return nil, ErrNotInit
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		// Verify the signing algorithm to prevent algorithm substitution attacks
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名方式无效")
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) ||
			errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// String retrieves a string value from Claims.Params by key.
//
// Parameters:
//   - k: parameter key (corresponds to field name in Params)
//
// Returns:
//   - string: the corresponding string value (if present and type matches)
//   - bool: whether the retrieval was successful
//
// Use cases:
//   - Retrieve business data from JWT extension field Params
//   - Such as room_id / role / scene and other string-type parameters
//
// Note:
//   - Only supports string type assertion
//   - Returns false if key does not exist or type does not match
//   - Params may be nil, no need to check before calling
func (c *Claims) String(k string) (string, bool) {
	v, ok := c.Params[k]
	if !ok || v == nil {
		return "", false
	}

	result, ok := v.(string)
	if !ok {
		return "", false
	}

	return result, true
}
