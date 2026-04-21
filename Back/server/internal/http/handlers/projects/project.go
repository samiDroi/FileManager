package projects

import (
	"fmt"
	"net/http"
	"server/internal/domain/models"
	"server/internal/http/middlewares"
	"server/internal/paths"
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
	g := e.Group(paths.Project.Prefix, middlewares.IsAuthenticated)
	g.GET(paths.Project.Index.Route, h.GetAll)
	g.GET(paths.Project.Search.Route, h.Search)
	g.POST(paths.Project.New.Route, h.Create)
	g.PUT(paths.Project.Update.Route, h.Update)
}

func (p *projectHandler) GetAll(c echo.Context) error {
	projects, err := p.repo.GetAllProjects(c.Request().Context())
	fmt.Print(projects)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	response := make([]projectResponse, len(projects))
	for i, project := range projects {
		response[i] = mapProjectToDtoToResponse(project)
	}
	return c.JSON(http.StatusCreated, response)
}
func (p *projectHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	ctx := c.Request().Context()

	projects, err := p.repo.Search(ctx, param)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := make([]projectResponse, len(projects))
	for i, project := range projects {
		response[i] = mapProjectToDtoToResponse(project)
	}
	return c.JSON(http.StatusOK, response)
}
func (h *projectHandler) Create(c echo.Context) error {
	var req projectRequest
	userSess := middlewares.GetRequestUserSession(c)

	c.Logger().Info("=== CREATE PROJECT DEBUG ===")
	c.Logger().Info("userSess.UserID: ", userSess.UserID)
	c.Logger().Info("userSess.IsAuthenticated: ", userSess.IsAuthenticated)

	if err := c.Bind(&req); err != nil {
		c.Logger().Error("Bind error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "El formato de los datos es incorrecto",
		})
	}

	c.Logger().Info("req.Name: ", req.Name)
	c.Logger().Info("req.ClientID: ", req.ClientID)

	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	createdProject, err := h.repo.CreateProject(c.Request().Context(), &projectModel, userSess.UserID)
	if err != nil {
		c.Logger().Error("Error creando proyecto: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ocurrió un error al guardar el proyecto",
		})
	}

	initialFiles := 0

	clientName, _ := h.repo.GetClientName(c.Request().Context(), createdProject.ClientID)

	response := projectResponse{
		ID:         createdProject.ID,
		Name:       createdProject.Name,
		ClientId:   createdProject.ClientID,
		ClientName: clientName,
		CountFiles: &initialFiles,
		CreatedAt:  utils.CustomTime{Time: createdProject.CreatedAt},
	}

	return c.JSON(http.StatusCreated, response)
}
func (h *projectHandler) Update(c echo.Context) error {
	idParam := c.Param("id")
	projectID, _ := strconv.Atoi(idParam)

	var req projectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("datos inválidos: %v", err),
		})
	}

	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	err := h.repo.UpdateProject(c.Request().Context(), &projectModel, projectID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "error en DB"})
	}

	updatedData, err := h.repo.GetProjectById(c.Request().Context(), projectID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]string{"message": "Actualizado (pero no pude recargar datos)"})
	}

	response := mapProjectToDtoToResponse(updatedData)

	return c.JSON(http.StatusOK, response)
}

func mapProjectToDtoToResponse(project *models.ProjectToDto) projectResponse {
	return projectResponse{
		ID:         project.ID,
		Name:       project.Name,
		ClientId:   project.ClientID,
		ClientName: project.ClientName,
		CountFiles: project.CountFiles,
		CreatedAt:  utils.CustomTime{Time: project.CreatedAt},
	}
}