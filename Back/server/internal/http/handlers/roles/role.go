package roles

import (
	"errors"
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	service *services.RoleService
}

func NewRoleHandler(service *services.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewRoleRepository(db)
	svc := services.NewRoleService(repo)
	h := NewRoleHandler(svc)

	e.GET("/permissions", h.GetPermissions, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.GET("/permissions/me", h.GetMyPermissions, middlewares.IsAuthenticated)
	e.GET("/roles", h.GetAll, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.POST("/roles/new", h.Create, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.PUT("/roles/edit/:id", h.Update, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.DELETE("/roles/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
}

func (h *RoleHandler) authorize(c echo.Context) (int, error) {
	userSession := middlewares.GetRequestUserSession(c)
	if !userSession.IsAuthenticated {
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Usuario no autenticado")
	}

	canManage, err := h.service.CanManage(c.Request().Context(), userSession.UserID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if !canManage {
		return 0, echo.NewHTTPError(http.StatusForbidden, "No cuentas con el permiso roles.gestionar")
	}
	return userSession.UserID, nil
}

func (h *RoleHandler) GetMyPermissions(c echo.Context) error {
	userSession := middlewares.GetRequestUserSession(c)
	if !userSession.IsAuthenticated {
		return echo.NewHTTPError(http.StatusUnauthorized, "Usuario no autenticado")
	}

	permissions, err := h.service.GetUserPermissions(c.Request().Context(), userSession.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, permissions)
}

func (h *RoleHandler) GetPermissions(c echo.Context) error {
	if _, err := h.authorize(c); err != nil {
		return err
	}

	permissions, err := h.service.GetPermissions(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, permissions)
}

func (h *RoleHandler) GetAll(c echo.Context) error {
	if _, err := h.authorize(c); err != nil {
		return err
	}

	roles, err := h.service.GetRoles(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) Create(c echo.Context) error {
	userID, err := h.authorize(c)
	if err != nil {
		return err
	}

	var req dtos.RoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	role, err := h.service.CreateRole(c.Request().Context(), req, userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRoleData) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, role)
}

func (h *RoleHandler) Update(c echo.Context) error {
	userID, err := h.authorize(c)
	if err != nil {
		return err
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	var req dtos.RoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	role, err := h.service.UpdateRole(c.Request().Context(), roleID, req, userID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRoleData) || errors.Is(err, services.ErrSystemRoleRename) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) Delete(c echo.Context) error {
	if _, err := h.authorize(c); err != nil {
		return err
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	err = h.service.DeleteRole(c.Request().Context(), roleID)
	if err != nil {
		if errors.Is(err, services.ErrSystemRole) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Rol eliminado correctamente"})
}
