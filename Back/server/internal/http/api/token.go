package api

import (
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

// GenerateTokens genera AMBOS tokens: Access (15m) y Refresh (7d)
func GenerateTokens(userID int, username string) (string, string, error) {
	if config.Config.SecretKey == "" {
		return "", "", fmt.Errorf("SECRET_KEY no está configurada")
	}

	secretKey := []byte(config.Config.SecretKey)

	// 1. ACCESS TOKEN: 15 min token
	accessClaims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// 2. REFRESH TOKEN: long duration (7 days)
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
func SetTokenCookies(c echo.Context, accessToken, refreshToken string) {
	// Cookie Access Token
	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   config.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// Cookie Refresh Token
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   config.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/auth/refresh", // Seguridad extra: solo se envía a la ruta de renovación
	})
}

func GetTokenCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("session_token")
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "not session token found")
	}
	return cookie.Value, nil
}

// validate the token and return claims if the token is valid
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		//verify signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(config.Config.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token inválido o expirado")
	}

	return claims, nil
}

// ClearCookies limpia ambas cookies
func ClearCookies(c echo.Context) {
	names := []string{"session_token", "refresh_token"}
	for _, name := range names {
		cookie := &http.Cookie{
			Name:     name,
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			Secure:   config.Config.CookieSecure,
			SameSite: http.SameSiteStrictMode,
		}
		c.SetCookie(cookie)
	}
}
