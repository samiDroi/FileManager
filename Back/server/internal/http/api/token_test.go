package api

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/config"

	"github.com/labstack/echo/v4"
)

func TestGenerateAccessTokenExpiresIn15Minutes(t *testing.T) {
	originalSecret := config.Config.SecretKey
	config.Config.SecretKey = "test-secret"
	defer func() { config.Config.SecretKey = originalSecret }()

	encoded, err := GenerateAccessToken(17, "test-user")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ValidateToken(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != 17 || claims.Username != "test-user" {
		t.Fatalf("unexpected access-token claims: %#v", claims)
	}
	if got := claims.ExpiresAt.Sub(claims.IssuedAt.Time); got != AccessTokenLifetime {
		t.Fatalf("access-token lifetime = %s, want %s", got, AccessTokenLifetime)
	}
}

func TestGenerateRefreshTokenIsOpaqueAndHashed(t *testing.T) {
	first, err := GenerateRefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateRefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("refresh tokens must be random")
	}

	decoded, err := base64.RawURLEncoding.DecodeString(first)
	if err != nil {
		t.Fatalf("refresh token is not URL-safe base64: %v", err)
	}
	if len(decoded) != 32 {
		t.Fatalf("refresh token entropy = %d bytes, want 32", len(decoded))
	}

	hash := HashRefreshToken(first)
	if len(hash) != 64 || hash == first {
		t.Fatalf("unexpected SHA-256 hash representation: %q", hash)
	}
	if _, err := hex.DecodeString(hash); err != nil {
		t.Fatalf("refresh token hash is not hexadecimal: %v", err)
	}
	if HashRefreshToken(first) != hash {
		t.Fatal("hashing the same refresh token must be deterministic")
	}
}

func TestSetRefreshTokenCookieSecurityAttributes(t *testing.T) {
	e := echo.New()
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodPost, "/login", nil), recorder)

	SetRefreshTokenCookie(ctx, "opaque-token")
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies, want 1", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != RefreshTokenCookieName || cookie.Value != "opaque-token" {
		t.Fatalf("unexpected cookie: %#v", cookie)
	}
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode {
		t.Fatalf("cookie security flags are not strict: %#v", cookie)
	}
	if remaining := time.Until(cookie.Expires); remaining < RefreshTokenLifetime-time.Minute || remaining > RefreshTokenLifetime {
		t.Fatalf("cookie expiration is %s, want approximately %s", remaining, RefreshTokenLifetime)
	}
}

func TestClearCookiesExpiresRefreshCookie(t *testing.T) {
	e := echo.New()
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodPost, "/logout", nil), recorder)

	ClearCookies(ctx)
	var refreshCookie *http.Cookie
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == RefreshTokenCookieName {
			refreshCookie = cookie
			break
		}
	}
	if refreshCookie == nil || refreshCookie.Value != "" || refreshCookie.MaxAge >= 0 {
		t.Fatalf("refresh cookie was not cleared: %#v", refreshCookie)
	}
}
