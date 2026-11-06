package auth

import (
	"context"
	"errors"
	"net/http"
	"server/internal/domain/models"
	"server/internal/http/api"
	"server/internal/services"
	"server/internal/store/repository"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AuthHandler struct {
	service     *services.AuthService
	sessionRepo *repository.SessionRepository
	roleRepo    *repository.RoleRepository
}

func NewAuthHandler(s *services.AuthService, sessionRepo *repository.SessionRepository, roleRepo *repository.RoleRepository) *AuthHandler {
	return &AuthHandler{service: s, sessionRepo: sessionRepo, roleRepo: roleRepo}
}

// there is one just RegisterRoutes because they're in the same package, so there is the resetPassword's routes
func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	// A. Construimos la cadena de mando de este módulo
	repo := repository.NewUserRepository(db)
	svc := services.NewAuthService(repo)
	h := NewAuthHandler(svc, repository.NewSessionRepository(db), repository.NewRoleRepository(db))

	authLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      10.0 / 60.0,
			Burst:     10,
			ExpiresIn: 3 * time.Minute,
		},
	))

	e.POST("/register", h.Register, authLimiter)
	e.POST("/login", h.Login, authLimiter)
	e.POST("/logout", h.Logout)
	e.POST("/refresh", h.RefreshToken)

	rh := NewResetPasswordHandler(db)
	e.POST("/forgot-password", rh.ForgotPassword, authLimiter)
	e.POST("/recovery-password", rh.RecoveryPassword, authLimiter)
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

	response, err := h.createSession(c, user)
	if err != nil {
		c.Logger().Errorf("register: no se pudo crear la sesión para el usuario %d: %v", user.ID, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al generar la sesión")
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)

	// Validate the request json body
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es inválido",
		})
	}

	// Call service to login the user
	user, err := h.service.LoginUser(c, req.Username, req.Password)
	if err != nil {
		return err
	}

	response, err := h.createSession(c, user)
	if err != nil {
		c.Logger().Errorf("login: no se pudo crear la sesión para el usuario %d: %v", user.ID, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al generar la sesión")
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie(api.RefreshTokenCookieName)
	if err == nil && cookie.Value != "" {
		session, lookupErr := h.sessionRepo.GetByTokenHash(
			c.Request().Context(),
			api.HashRefreshToken(cookie.Value),
		)
		if lookupErr == nil {
			if err := h.sessionRepo.RevokeFamily(c.Request().Context(), session.FamilyID); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo cerrar la sesión")
			}
		} else if !errors.Is(lookupErr, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo cerrar la sesión")
		}
	}
	api.ClearCookies(c)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Sesión cerrada correctamente",
	})
}

func (h *AuthHandler) createSession(c echo.Context, user *models.User) (*UserResponse, error) {
	userPermissions, err := h.roleRepo.GetUserPermissions(c.Request().Context(), user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := api.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := api.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	familyID, err := api.GenerateSessionFamilyID()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(api.RefreshTokenLifetime)
	if err := h.sessionRepo.Create(
		c.Request().Context(),
		user.ID,
		api.HashRefreshToken(refreshToken),
		familyID,
		expiresAt,
		c.RealIP(),
		c.Request().UserAgent(),
	); err != nil {
		return nil, err
	}

	api.ClearLegacyAccessTokenCookie(c)
	api.SetRefreshTokenCookie(c, refreshToken)
	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		AccessToken: accessToken,
		IsAdmin:     userPermissions.IsAdmin,
		Roles:       userPermissions.Roles,
		Permissions: userPermissions.Permissions,
	}, nil
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie(api.RefreshTokenCookieName)
	if err != nil {
		api.ClearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "No se encontró el refresh token")
	}

	session, err := h.sessionRepo.GetByTokenHash(
		c.Request().Context(),
		api.HashRefreshToken(cookie.Value),
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			api.ClearCookies(c)
			return echo.NewHTTPError(http.StatusUnauthorized, "Sesión expirada. Por favor, inicia sesión de nuevo.")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo validar la sesión")
	}

	now := time.Now()
	if !session.ExpiresAt.After(now) {
		_ = h.sessionRepo.RevokeFamily(c.Request().Context(), session.FamilyID)
		api.ClearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "Sesión expirada. Por favor, inicia sesión de nuevo.")
	}

	if session.IsRevoked {
		if session.ReplacedByID != nil &&
			session.RevokedAt != nil &&
			now.Sub(*session.RevokedAt) <= api.RefreshTokenGracePeriod {
			return h.respondWithAccessToken(c, session.UserID, session.Username)
		}

		_ = h.sessionRepo.RevokeFamily(c.Request().Context(), session.FamilyID)
		api.ClearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "Sesión comprometida. Por favor, inicia sesión de nuevo.")
	}

	refreshToken, err := api.GenerateRefreshToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error interno al generar el token")
	}

	err = h.sessionRepo.Rotate(
		c.Request().Context(),
		session.ID,
		api.HashRefreshToken(refreshToken),
		session.ExpiresAt,
		c.RealIP(),
		c.Request().UserAgent(),
	)
	if err != nil {
		if errors.Is(err, repository.ErrSessionAlreadyRotated) {
			return h.respondWithAccessToken(c, session.UserID, session.Username)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo renovar la sesión")
	}

	api.ClearLegacyAccessTokenCookie(c)
	api.SetRefreshTokenCookie(c, refreshToken)
	return h.respondWithAccessToken(c, session.UserID, session.Username)
}

func (h *AuthHandler) respondWithAccessToken(c echo.Context, userID int, username string) error {
	return h.buildAccessTokenResponse(c.Request().Context(), userID, username, c)
}

func (h *AuthHandler) buildAccessTokenResponse(ctx context.Context, userID int, username string, c echo.Context) error {
	accessToken, err := api.GenerateAccessToken(userID, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error interno al generar el token")
	}

	userPermissions, err := h.roleRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "No se pudieron cargar los permisos del usuario")
	}

	return c.JSON(http.StatusOK, TokenResponse{
		AccessToken: accessToken,
		IsAdmin:     userPermissions.IsAdmin,
		Roles:       userPermissions.Roles,
		Permissions: userPermissions.Permissions,
	})
}
