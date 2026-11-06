package projects

import (
	"errors"
	"fmt"
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type projectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *projectHandler {
	return &projectHandler{
		service: service,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewProjectRepository(db)
	svc := services.NewProjectService(repo)
	h := NewProjectHandler(svc)
	e.GET("/project", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/project/search", h.Search, middlewares.IsAuthenticated)
	e.POST("/project/new", h.Create, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "proyectos", "crear"))
	e.PUT("/project/edit/:id", h.Update, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "proyectos", "editar"))
	e.DELETE("/project/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "proyectos", "eliminar"))
}

func (h *projectHandler) GetAll(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	projects, err := h.service.GetAllProjects(c.Request().Context(), userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *projectHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	userSess := middlewares.GetRequestUserSession(c)
	projects, err := h.service.SearchProjects(c.Request().Context(), param, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *projectHandler) Create(c echo.Context) error {
	var req dtos.ProjectRequest
	userSess := middlewares.GetRequestUserSession(c)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es incorrecto",
		})
	}

	createdProject, err := h.service.CreateProject(c.Request().Context(), req, userSess.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ocurrió un error al guardar el proyecto",
		})
	}

	return c.JSON(http.StatusCreated, createdProject)
}

func (h *projectHandler) Delete(c echo.Context) error {
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de proyecto inválido")
	}

	if err := h.service.DeleteProject(c.Request().Context(), projectID); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Proyecto no encontrado")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Proyecto eliminado correctamente"})
}

func (h *projectHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	projectID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "invalid project ID"})
	}

	var req dtos.ProjectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("datos inválidos: %v", err),
		})
	}

	updatedData, err := h.service.UpdateProject(c.Request().Context(), projectID, req)
	if err != nil {
		if errors.Is(err, services.ErrProjectReload) {
			return c.JSON(http.StatusOK, map[string]string{"message": "Actualizado (pero no pude recargar datos)"})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "error en DB"})
	}

	return c.JSON(http.StatusOK, updatedData)
}
