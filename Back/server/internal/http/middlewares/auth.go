package middlewares

import (
	"context"
	"net/http"
	"strings"
	"time"

	"server/internal/domain/models"
	"server/internal/http/api"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var DefaultExpiresSessionTime time.Duration = 24 * time.Hour * 90
var SessionExpirationTime int = 24 * int(time.Hour) * 7

func GetRequestUserSession(c echo.Context) models.UserSession {
	// Extraer directamente de los values del contexto que establece JWTMiddleware
	userIDVal := c.Get("user_id")
	usernameVal := c.Get("username")

	// Validar que tengamos los valores necesarios
	if userIDVal == nil {
		c.Logger().Error("user_id no encontrado en contexto - usuario no autenticado")
		return models.UserSession{IsAuthenticated: false}
	}

	// Convertir con seguridad
	userID, ok := userIDVal.(int)
	if !ok {
		c.Logger().Error("user_id tiene tipo incorrecto en contexto")
		return models.UserSession{IsAuthenticated: false}
	}

	username, ok := usernameVal.(string)
	if !ok {
		username = "" // Username es menos crítico
	}

	return models.UserSession{
		UserID:          userID,
		Username:        username,
		IsAuthenticated: true,
	}
}

func GetSessionUser(c echo.Context) (models.UserSession, bool) {
	sess, err := session.Get("session", c)
	if err != nil {
		c.Logger().Error("Error recuperando sesión: ", err)
		return models.UserSession{}, false
	}

	userSess, ok := sess.Values["user_session"].(models.UserSession)
	if !ok || !userSess.IsAuthenticated {
		return models.UserSession{}, false
	}

	return userSess, true
}

// Middleware para proteger rutas (Solo usuarios logueados)
func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
		scheme, tokenString, ok := strings.Cut(authorization, " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(tokenString) == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token de acceso no encontrado"})
		}

		claims, err := api.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
		}

		userSess := models.UserSession{
			UserID:          claims.UserID,
			Username:        claims.Username,
			IsAuthenticated: true,
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		ctx := context.WithValue(c.Request().Context(), "user", userSess)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func RequireAdmin(db *pgxpool.Pool) echo.MiddlewareFunc {
	roleRepo := repository.NewRoleRepository(db)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := authenticatedUserID(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Usuario no autenticado")
			}

			isAdmin, err := roleRepo.IsAdmin(c.Request().Context(), userID)
			if err != nil {
				c.Logger().Errorf("No se pudo validar el rol administrador del usuario %d: %v", userID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo validar el acceso")
			}
			if !isAdmin {
				return echo.NewHTTPError(http.StatusForbidden, "Se requiere el rol administrador")
			}
			return next(c)
		}
	}
}

func RequirePermission(db *pgxpool.Pool, resource, action string) echo.MiddlewareFunc {
	roleRepo := repository.NewRoleRepository(db)
	permission := strings.TrimSpace(resource) + "." + strings.TrimSpace(action)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.TrimSpace(resource) == "" || strings.TrimSpace(action) == "" {
				return echo.NewHTTPError(http.StatusInternalServerError, "Middleware de permisos mal configurado")
			}

			userID, ok := authenticatedUserID(c)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Usuario no autenticado")
			}

			allowed, err := roleRepo.HasPermission(c.Request().Context(), userID, permission)
			if err != nil {
				c.Logger().Errorf("No se pudo validar el permiso %s del usuario %d: %v", permission, userID, err)
				return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo validar el acceso")
			}
			if !allowed {
				return echo.NewHTTPError(http.StatusForbidden, "No cuentas con el permiso "+permission)
			}
			return next(c)
		}
	}
}

func authenticatedUserID(c echo.Context) (int, bool) {
	userID, ok := c.Get("user_id").(int)
	return userID, ok && userID > 0
}

// func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userSess, isValid := GetSessionUser(c)

// 		if !isValid {
// 			// ✅ CAMBIO PARA ANGULAR: Devolver 401 Unauthorized en JSON
// 			return c.JSON(http.StatusUnauthorized, map[string]string{
// 				"error": "No autenticado o sesión expirada",
// 			})
// 		}

// 		c.Set("user", userSess)
// 		ctx := context.WithValue(c.Request().Context(), "user", userSess)
// 		c.SetRequest(c.Request().WithContext(ctx))

// 		return next(c)
// 	}
// }

// // Middleware para proteger rutas (Solo Super Admin)
// func IsSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		userSess := GetRequestUserSession(c)

// 		// if !userSess.IsSuperAdmin() {
// 		// 	// ✅ CAMBIO PARA ANGULAR: Devolver 403 Forbidden en JSON
// 		// 	return c.JSON(http.StatusForbidden, map[string]string{
// 		// 		"error": "No cuentas con los permisos necesarios",
// 		// 	})
// 		// }

// 		return next(c)
// 	}
// }
