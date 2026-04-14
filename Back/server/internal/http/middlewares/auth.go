package middlewares

import (
	"context"
	"net/http"
	"time"

	// IMPORTANTE: Cambia estas rutas por las de tu proyecto
	"server/internal/domain/models"

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
		userSess, isValid := GetSessionUser(c)

		if !isValid {
			// ✅ CAMBIO PARA ANGULAR: Devolver 401 Unauthorized en JSON
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "No autenticado o sesión expirada",
			})
		}

		c.Set("user", userSess)
		ctx := context.WithValue(c.Request().Context(), "user", userSess)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

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
