package middlewares

import (
	"net/http"
	"server/internal/http/api"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware valida que las peticiones tengan un token JWT válido
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el token del header Authorization
		token, err := api.GetTokenByBearerAuth(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token de autenticación no válido",
			})
		}

		// Validar el token
		claims, err := api.ValidateToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Token inválido o expirado: " + err.Error(),
			})
		}

		// Agregar los claims al contexto para que los handlers puedan acceder a ellos
		c.Set("user", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		return next(c)
	}
}
