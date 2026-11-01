package projects

import (
	"fmt"
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/http/middlewares"
	"server/internal/store/repository"
	"server/internal/utils"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type projectHandler struct {
	db   *pgxpool.Pool
	repo *repository.ProjectRepository
}

func NewProjectHandler(db *pgxpool.Pool) *projectHandler {
	return &projectHandler{
		db:   db,
		repo: repository.NewProjectRepository(db),
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewProjectHandler(db)
	e.GET("/project", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/project/search", h.Search, middlewares.IsAuthenticated)
	e.POST("/project/new", h.Create, middlewares.IsAuthenticated)
	e.PUT("/project/edit/:id", h.Update, middlewares.IsAuthenticated)
}

func (h *projectHandler) GetAll(c echo.Context) error {
	projects, err := h.repo.GetAllProjects(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, projects)
}

func (h *projectHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	projects, err := h.repo.Search(c.Request().Context(), param)
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

	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	createdProject, err := h.repo.CreateProject(c.Request().Context(), &projectModel, userSess.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ocurrió un error al guardar el proyecto",
		})
	}

	clientName, _ := h.repo.GetClientName(c.Request().Context(), createdProject.ClientID)

	initialFiles := 0
	response := dtos.ProjectResponse{
		ID:         createdProject.ID,
		Name:       createdProject.Name,
		ClientID:   createdProject.ClientID,
		ClientName: clientName,
		CountFiles: &initialFiles,
		CreatedAt:  utils.ToCustomTime(createdProject.CreatedAt),
	}

	return c.JSON(http.StatusCreated, response)
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

	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	err = h.repo.UpdateProject(c.Request().Context(), &projectModel, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "error en DB"})
	}

	updatedData, err := h.repo.GetProjectById(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Actualizado (pero no pude recargar datos)"})
	}

	return c.JSON(http.StatusOK, updatedData)
}
