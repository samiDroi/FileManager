package users

import (
	"errors"
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service  *services.UserService
	roleRepo *repository.RoleRepository
}

func NewUserHandler(service *services.UserService, roleRepo *repository.RoleRepository) *UserHandler {
	return &UserHandler{
		service:  service,
		roleRepo: roleRepo,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewUserRepository(db)
	svc := services.NewUserService(repo)
	h := NewUserHandler(svc, repository.NewRoleRepository(db))
	e.GET("/users", h.GetAll, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.GET("/users/search", h.Search, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.GET("/users/assignable", h.GetAll, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
	e.GET("/users/:id", h.GetByID, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
	e.POST("/users", h.Create, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "usuarios", "crear"))
	e.PUT("/users/edit/:id", h.Update, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "usuarios", "editar"))
	e.DELETE("/users/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "usuarios", "eliminar"))
}

func (h *UserHandler) Create(c echo.Context) error {
	var req dtos.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}
	if req.Username == "" || req.Email == "" || len(req.Password) < 6 || req.RoleID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Usuario, correo, contraseña y rol válidos son obligatorios")
	}

	userSession := middlewares.GetRequestUserSession(c)
	user, err := h.service.CreateUser(c.Request().Context(), req, userSession.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "El rol seleccionado no existe")
		}
		c.Logger().Errorf("No se pudo crear el usuario: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo crear el usuario")
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	users, err := h.service.SearchUsers(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	userID := c.Param("id")
	paramID, err := strconv.Atoi(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	userSession := middlewares.GetRequestUserSession(c)
	if userSession.UserID != paramID {
		allowed, err := h.roleRepo.HasPermission(c.Request().Context(), userSession.UserID, "usuarios.leer")
		if err != nil {
			c.Logger().Errorf("No se pudo validar el permiso usuarios.leer: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo validar el acceso")
		}
		if !allowed {
			return echo.NewHTTPError(http.StatusForbidden, "No cuentas con el permiso usuarios.leer")
		}
	}

	user, err := h.service.GetUserByID(c.Request().Context(), paramID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c echo.Context) error {
	userID := c.Param("id")
	paramID, err := strconv.Atoi(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	var req dtos.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if req.Username == "" || req.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Username and email are required")
	}
	if req.RoleID != nil {
		userSession := middlewares.GetRequestUserSession(c)
		canManageRoles, err := h.roleRepo.CanManage(c.Request().Context(), userSession.UserID)
		if err != nil {
			c.Logger().Errorf("No se pudo validar la asignación de rol: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "No se pudo validar el acceso")
		}
		if !canManageRoles {
			return echo.NewHTTPError(http.StatusForbidden, "Se requiere el permiso roles.gestionar para cambiar el rol")
		}
		if *req.RoleID <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Role ID is invalid")
		}
	}
	userSession := middlewares.GetRequestUserSession(c)
	user, err := h.service.UpdateUser(c.Request().Context(), &req, paramID, userSession.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "User or role not found")
		}
		c.Logger().Errorf("No se pudo actualizar el usuario %d: %v", paramID, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c echo.Context) error {
	userID := c.Param("id")
	paramID, err := strconv.Atoi(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}
	err = h.service.DeleteUser(c.Request().Context(), paramID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
