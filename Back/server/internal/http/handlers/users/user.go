package users

import (
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db   *pgxpool.Pool
	repo *repository.UserRepository
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		db:   db,
		repo: repository.NewUserRepository(db),
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewUserHandler(db)
	e.GET("/users", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/users/search", h.Search, middlewares.IsAuthenticated)
	e.GET("/users/:id", h.GetByID, middlewares.IsAuthenticated)
	e.PUT("/users/edit/:id", h.Update, middlewares.IsAuthenticated)
	e.DELETE("/users/:id", h.Delete, middlewares.IsAuthenticated)
}

func (h *UserHandler) Search(c echo.Context) error {
	query := c.QueryParam("query")
	users, err := h.repo.SearchUsers(c.Request().Context(), query)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetAll(c echo.Context) error {
	users, err := h.repo.GetAllUsers(c.Request().Context())
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
	user, err := h.repo.GetUserByID(c.Request().Context(), paramID)
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
	user, err := h.repo.UpdateUser(c.Request().Context(), &req, paramID)
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
	err = h.repo.DeleteUser(c.Request().Context(), paramID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
