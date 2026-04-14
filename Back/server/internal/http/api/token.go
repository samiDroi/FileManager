package api

import (
	"fmt"
	"server/internal/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Claims estructura personalizada con información del usuario
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerarToken genera un JWT token con los datos del usuario
func GenerateToken(userID int, username, email string) (string, error) {
	// Verificar que tengamos SECRET_KEY configurada
	if config.Config.SecretKey == "" {
		return "", fmt.Errorf("SECRET_KEY no está configurada")
	}

	// Crear los claims con la información del usuario
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	// create jwt with signature HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar el token con la SECRET_KEY
	secretKey := []byte(config.Config.SecretKey)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return tokenString, nil
}

// ValidarToken valida un JWT token y retorna los claims si es válido
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(config.Config.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear el token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}

// ObtenerTokenDesdeBearerAuth extrae el token del header Authorization
func GetTokenByBearerAuth(c echo.Context) (string, error) {
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return "", echo.NewHTTPError(401, "token de autenticación no proporcionado")
	}

	// Validar que sea Bearer token
	const bearerScheme = "Bearer "
	if len(auth) < len(bearerScheme) || auth[:len(bearerScheme)] != bearerScheme {
		return "", echo.NewHTTPError(401, "formato de token inválido")
	}

	return auth[len(bearerScheme):], nil
}
