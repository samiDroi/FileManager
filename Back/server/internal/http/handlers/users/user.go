package users

import (
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewUserRepository(db)
	svc := services.NewUserService(repo)
	h := NewUserHandler(svc)
	e.GET("/users", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/users/search", h.Search, middlewares.IsAuthenticated)
	e.GET("/users/:id", h.GetByID, middlewares.IsAuthenticated)
	e.PUT("/users/edit/:id", h.Update, middlewares.IsAuthenticated)
	e.DELETE("/users/:id", h.Delete, middlewares.IsAuthenticated)
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
	user, err := h.service.UpdateUser(c.Request().Context(), &req, paramID)
	if err != nil {
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
