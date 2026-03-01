package webserver

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT claims structure
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled       bool   `json:"enabled"`
	SecretKey     string `json:"secret_key"`
	TokenDuration int    `json:"token_duration"` // in minutes
}

// NewAuthConfig creates default auth configuration
func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Enabled:       false,
		SecretKey:     generateRandomSecret(32),
		TokenDuration: 60, // 1 hour by default
	}
}

// generateRandomSecret generates a cryptographically secure random secret
func generateRandomSecret(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based secret if crypto/rand fails
		return fmt.Sprintf("triageprof-secret-%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

// GenerateToken creates a new JWT token
func (ac *AuthConfig) GenerateToken(username, role string) (string, error) {
	if !ac.Enabled {
		return "", errors.New("authentication is disabled")
	}

	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ac.TokenDuration) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "triageprof",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ac.SecretKey))
}

// ValidateToken validates a JWT token
func (ac *AuthConfig) ValidateToken(tokenString string) (*Claims, error) {
	if !ac.Enabled {
		return &Claims{Username: "anonymous", Role: "viewer"}, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ac.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(*Claims); ok {
			return claims, nil
		}
	}

	return nil, errors.New("invalid token")
}

// ExtractTokenFromRequest extracts JWT token from request headers or query parameters
func ExtractTokenFromRequest(r *http.Request) string {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// Check query parameter
	return r.URL.Query().Get("token")
}

// AuthMiddleware is middleware for JWT authentication
func AuthMiddleware(ac *AuthConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip authentication for health and root endpoints
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next(w, r)
			return
		}

		// Extract token
		tokenString := ExtractTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Validate token
		_, err := ac.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		r = r.WithContext(r.Context())
		// TODO: Add claims to context for future use

		next(w, r)
	}
}

// GenerateTokenHandler handles token generation requests
type TokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (ac *AuthConfig) GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if !ac.Enabled {
		http.Error(w, "Authentication is disabled", http.StatusServiceUnavailable)
		return
	}

	// For demo purposes, allow any username/password
	// In production, this should validate against a user store
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	// Default role if not specified
	if req.Role == "" {
		req.Role = "viewer"
	}

	// Generate token
	token, err := ac.GenerateToken(req.Username, req.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":    token,
		"expires_in": ac.TokenDuration * 60, // seconds
		"username": req.Username,
		"role":     req.Role,
	})
}