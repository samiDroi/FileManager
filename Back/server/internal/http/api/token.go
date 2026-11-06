package api

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Claims struct personalized with user info
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	// Email    string `json:"email"`
	jwt.RegisteredClaims
}

const RefreshTokenCookieName = "refresh_token"

const AccessTokenLifetime = 15 * time.Minute
const RefreshTokenLifetime = 7 * 24 * time.Hour
const RefreshTokenGracePeriod = 30 * time.Second

func GenerateAccessToken(userID int, username string) (string, error) {
	if config.Config.SecretKey == "" {
		return "", fmt.Errorf("SECRET_KEY no está configurada")
	}

	secretKey := []byte(config.Config.SecretKey)
	now := time.Now()
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secretKey)
}

func GenerateRefreshToken() (string, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("generar refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(secret), nil
}

func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateSessionFamilyID() (string, error) {
	familyID := make([]byte, 16)
	if _, err := rand.Read(familyID); err != nil {
		return "", fmt.Errorf("generar family id: %w", err)
	}
	return hex.EncodeToString(familyID), nil
}

func GenerateResetPasswordToken(userID int, currentPasswordHash string) (string, error) {
	if config.Config.SecretKey == "" {
		return "", fmt.Errorf("SECRET_KEY no está configurada")
	}

	secretKey := []byte(config.Config.SecretKey + currentPasswordHash)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
}

func SetRefreshTokenCookie(c echo.Context, refreshToken string) {
	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(RefreshTokenLifetime),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func ClearLegacyAccessTokenCookie(c echo.Context) {
	clearCookie(c, "session_token")
}

// validate the token and return claims if the token is valid
func ValidateToken(tokenString string) (*Claims, error) {
	if config.Config.SecretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY no está configurada")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(config.Config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token inválido o expirado")
	}

	return claims, nil
}

// ClearCookies removes the refresh cookie and any legacy access-token cookie.
func ClearCookies(c echo.Context) {
	ClearLegacyAccessTokenCookie(c)
	clearCookie(c, RefreshTokenCookieName)
}

func clearCookie(c echo.Context, name string) {
	c.SetCookie(&http.Cookie{
		Name:     name,
		Value:    "",
		Expires:  time.Unix(1, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
