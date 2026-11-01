package clients

import (
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/http/middlewares"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	DB   *pgxpool.Pool
	repo *repository.ClientRepository
}

func NewClientHandler(db *pgxpool.Pool) *ClientHandler {
	return &ClientHandler{
		DB:   db,
		repo: repository.NewClientRepository(db),
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewClientHandler(db)

	e.GET("/client", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/client/search", h.Search, middlewares.IsAuthenticated)
	e.POST("/client/new", h.Create, middlewares.IsAuthenticated)
	e.PUT("/client/edit/:id", h.Edit, middlewares.IsAuthenticated)
	e.DELETE("/client/:id", h.Delete, middlewares.IsAuthenticated)
	e.GET("/client/with-projects", h.GetClientsWithProjects, middlewares.IsAuthenticated)
}

func (h *ClientHandler) GetAll(c echo.Context) error {
	clients, err := h.repo.GetAllClients(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, clients)
}

func (h *ClientHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	clients, err := h.repo.SearchClient(c.Request().Context(), param)
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

	var contactInfo *models.ContactInfo
	if req.Email != nil || req.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			Email:        req.Email,
			ContactPhone: req.ContactPhone,
		}
	}

	newClient := models.Client{
		Name:        req.Name,
		ContactInfo: contactInfo,
	}

	createdClient, err := h.repo.CreateClient(c.Request().Context(), newClient, userSess.UserID)
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

	var contactInfo *models.ContactInfo
	if req.Email != nil || req.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			ClientID:     clientID,
			Email:        req.Email,
			ContactPhone: req.ContactPhone,
		}
	}

	updatedClient := models.Client{
		ID:          clientID,
		Name:        req.Name,
		ContactInfo: contactInfo,
	}

	client, err := h.repo.UpdateClient(c.Request().Context(), updatedClient, userSess.UserID)
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

	err = h.repo.DeleteClient(c.Request().Context(), clientID, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al eliminar el cliente")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cliente eliminado correctamente"})
}
func (h *ClientHandler) GetClientsWithProjects(c echo.Context) error {
	clients, err := h.repo.GetClientsWithProjects(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, clients)
}
