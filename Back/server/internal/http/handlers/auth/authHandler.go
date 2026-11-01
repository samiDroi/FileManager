package auth

import (
	"fmt"
	"net/http"
	"server/internal/domain/models"
	"server/internal/http/api"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// there is one just RegisterRoutes because they're in the same package, so there is the resetPassword's routes
func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	// A. Construimos la cadena de mando de este módulo
	repo := repository.NewUserRepository(db)
	svc := services.NewAuthService(repo)
	h := &AuthHandler{service: svc}

	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.POST("/logout", h.Logout, middlewares.IsAuthenticated)
	e.POST("/refresh", h.RefreshToken)

	rh := NewResetPasswordHandler(db)
	e.POST("/forgot-password", rh.ForgotPassword)
	e.POST("/recovery-password", rh.RecoveryPassword)
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)

	// Validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}

	// Adapt the request dto to the user model
	userDto := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Call service to register the user
	user, err := h.service.RegisterUser(c, userDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al guardar en base de datos: " + err.Error(),
		})
	}

	// Generar AMBOS tokens para auto-login
	sessionToken, refreshToken, err := createTokens(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al generar la sesión: " + err.Error(),
		})
	}

	// Guardar las dos cookies
	api.SetTokenCookies(c, sessionToken, refreshToken)

	// Exit response
	return c.JSON(http.StatusCreated, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)

	// Validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}

	fmt.Print(req.Username, req.Password)

	// Call service to login the user
	user, err := h.service.LoginUser(c, req.Username, req.Password)
	if err != nil {
		fmt.Printf("DEBUG: Error en Scan de base de datos: %v\n", err)
		return err
	}

	// Generar AMBOS tokens
	sessionToken, refreshToken, err := createTokens(user.ID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error al generar la sesión: " + err.Error(),
		})
	}

	// Guardar las dos cookies
	api.SetTokenCookies(c, sessionToken, refreshToken)

	// Exit response
	return c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	// Ahora limpia ambas cookies (session_token y refresh_token)
	api.ClearCookies(c)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sesión cerrada correctamente",
	})
}

// Helper to create tokens
func createTokens(userID int, username string) (string, string, error) {
	sessionToken, refreshToken, err := api.GenerateTokens(userID, username)
	if err != nil {
		return "", "", err
	}
	return sessionToken, refreshToken, nil
}

// RefreshTokenHandler lee el refresh token, lo valida y emite nuevos tokens
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// 1. Leer la cookie del Refresh Token
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		// Si no hay cookie, el usuario debe iniciar sesión desde cero
		return echo.NewHTTPError(http.StatusUnauthorized, "No se encontró el refresh token")
	}

	// 2. Validar el Refresh Token
	// Reutilizamos tu función ValidateToken que verifica la firma y la fecha de expiración
	claims, err := api.ValidateToken(cookie.Value)
	if err != nil {
		// Si el refresh token expiró (pasaron los 7 días) o es inválido
		// Limpiamos las cookies por seguridad y forzamos el re-login
		api.ClearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "Sesión expirada. Por favor, inicia sesión de nuevo.")
	}

	// 3. Generar nuevos tokens
	// Usamos los datos (UserID y Username) que venían dentro del Refresh Token válido
	newSessionToken, newRefreshToken, err := api.GenerateTokens(claims.UserID, claims.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error interno al generar la sesión")
	}

	// 4. Actualizar las cookies en el navegador
	api.SetTokenCookies(c, newSessionToken, newRefreshToken)

	// 5. Devolver una respuesta de éxito
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sesión renovada exitosamente",
	})
}
