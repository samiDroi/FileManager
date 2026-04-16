package clients

import (
	"fmt"
	"net/http"
	"server/internal/domain/models"
	"server/internal/paths"
	"strconv"

	"server/internal/http/middlewares"
	"server/internal/store/repository"

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

	g := e.Group(paths.Client.Prefix)
	g.Use(middlewares.JWTMiddleware)
	g.GET(paths.Client.Index.Route, h.getAll)
	g.GET(paths.Client.Search.Route, h.SearchClientHandler)
	g.POST(paths.Client.New.Route, h.Create)
	g.PUT(paths.Client.Edit.Route, h.Edit)
	g.DELETE(paths.Client.Delete.Route, h.Delete)
}

func (h *ClientHandler) getAll(c echo.Context) error {
	clients, err := h.repo.GetAllClients(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, clients)
}

func (h *ClientHandler) SearchClientHandler(c echo.Context) error {
	param := c.QueryParam("q")
	ctx := c.Request().Context()

	clients, err := h.repo.SearchClient(ctx, param)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

	var input InputClient
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	var contactInfo *models.ContactInfo
	if input.Email != nil || input.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			Email:        input.Email,
			ContactPhone: input.ContactPhone,
		}
	}

	newClient := models.Client{
		Name:        input.Name,
		ContactInfo: contactInfo,
	}
	fmt.Print(userSess)
	createdClient, err := h.repo.CreateClient(c, newClient, userSess.UserID)
	if err != nil {
		fmt.Printf("🛑 ERROR FATAL EN BASE DE DATOS: %v\n", err)
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

	var input InputClient
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	var contactInfo *models.ContactInfo
	if input.Email != nil || input.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			ClientID:     clientID,
			Email:        input.Email,
			ContactPhone: input.ContactPhone,
		}
	}

	updatedClient := models.Client{
		ID:          clientID,
		Name:        input.Name,
		ContactInfo: contactInfo,
	}

	client, err := h.repo.UpdateClient(c, updatedClient, userSess.UserID)
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

	err = h.repo.DeleteClient(c, clientID, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al eliminar el cliente")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cliente eliminado correctamente"})
}
