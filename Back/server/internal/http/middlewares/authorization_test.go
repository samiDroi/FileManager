package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestAuthorizationMiddlewaresRequireAuthenticatedUser(t *testing.T) {
	tests := []struct {
		name       string
		middleware echo.MiddlewareFunc
	}{
		{
			name:       "admin",
			middleware: RequireAdmin(nil),
		},
		{
			name:       "permission",
			middleware: RequirePermission(nil, "archivos", "editar"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), recorder)
			handler := test.middleware(func(c echo.Context) error {
				t.Fatal("protected handler ran without an authenticated user")
				return nil
			})

			if err := handler(ctx); err == nil {
				t.Fatal("expected unauthorized response")
			} else if httpError, ok := err.(*echo.HTTPError); !ok || httpError.Code != http.StatusUnauthorized {
				t.Fatalf("got error %v, want HTTP 401", err)
			}
		})
	}
}

func TestRequirePermissionRejectsEmptyPermissionConfiguration(t *testing.T) {
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	handler := RequirePermission(nil, " ", "editar")(func(c echo.Context) error {
		t.Fatal("protected handler ran with invalid permission configuration")
		return nil
	})

	if err := handler(ctx); err == nil {
		t.Fatal("expected configuration error")
	} else if httpError, ok := err.(*echo.HTTPError); !ok || httpError.Code != http.StatusInternalServerError {
		t.Fatalf("got error %v, want HTTP 500", err)
	}
}
