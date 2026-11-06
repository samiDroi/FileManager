package clients

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

type ClientHandler struct {
	service *services.ClientService
}

func NewClientHandler(service *services.ClientService) *ClientHandler {
	return &ClientHandler{
		service: service,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewClientRepository(db)
	svc := services.NewClientService(repo)
	h := NewClientHandler(svc)

	e.GET("/client", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/client/search", h.Search, middlewares.IsAuthenticated)
	e.POST("/client/new", h.Create, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "clientes", "crear"))
	e.PUT("/client/edit/:id", h.Edit, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "clientes", "editar"))
	e.DELETE("/client/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "clientes", "eliminar"))
	e.GET("/client/with-projects", h.GetClientsWithProjects, middlewares.IsAuthenticated)
}

func (h *ClientHandler) GetAll(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	clients, err := h.service.GetAllClients(c.Request().Context(), userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, clients)
}

func (h *ClientHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	userSess := middlewares.GetRequestUserSession(c)
	clients, err := h.service.SearchClient(c.Request().Context(), param, userSess.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Fallo en BD"})
	}
	return c.JSON(http.StatusOK, clients)
}

func (h *ClientHandler) Create(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	var req dtos.ClientRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	createdClient, err := h.service.CreateClient(c.Request().Context(), req, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error DB: "+err.Error())
	}
	return c.JSON(http.StatusCreated, createdClient)
}

func (h *ClientHandler) Edit(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	idParam := c.Param("id")
	clientID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	var req dtos.ClientRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	client, err := h.service.UpdateClient(c.Request().Context(), clientID, req, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al actualizar el cliente")
	}

	return c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Delete(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	idParam := c.Param("id")
	clientID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	err = h.service.DeleteClient(c.Request().Context(), clientID, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al eliminar el cliente")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cliente eliminado correctamente"})
}

func (h *ClientHandler) GetClientsWithProjects(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	clients, err := h.service.GetClientsWithProjects(c.Request().Context(), userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, clients)
}
