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
	g.Use(middlewares.JWTMiddleware) // Protegemos las rutas con autenticación JWT
	//routes
	g.GET(paths.Client.Index.Route, h.getAll)
	g.GET(paths.Client.Search.Route, h.SearchClientHandler)
	g.POST(paths.Client.New.Route, h.Create)
	g.PUT(paths.Client.Edit.Route, h.Edit)
	g.DELETE(paths.Client.Delete.Route, h.Delete)
}

func (h *ClientHandler) getAll(c echo.Context) error {
	// h.repo = repository.NewClientRepository(h.DB)
	clients, err := h.repo.GetAllClients(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al obtener los clientes")
	}
	return c.JSON(http.StatusCreated, clients)
}

// ESTE es el método que Echo va a ejecutar
func (h *ClientHandler) SearchClientHandler(c echo.Context) error {
	// 1. Extraer los datos de la petición HTTP
	param := c.QueryParam("q")

	// 2. Extraer el contexto de la request
	ctx := c.Request().Context()

	// 3. Llamar al repositorio para que haga el trabajo pesado
	clients, err := h.repo.SearchClient(ctx, param)

	if err != nil {
		// Si la base de datos falla, el handler decide qué error HTTP mostrar
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 4. Devolver la respuesta HTTP exitosa
	return c.JSON(http.StatusOK, clients)
}
func (h *ClientHandler) Create(c echo.Context) error {
	// 1. Extraemos la sesión y validamos que esté autenticado
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	// 2. Usamos tu DTO para leer el JSON
	var input InputClient
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	// 3. Validamos que cumpla las reglas (ej. validate:"required")
	// if err := c.Validate(&input); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, &input)
	// }

	// 4. Mapeamos el DTO al modelo de la base de datos
	newClient := models.Client{
		Name: input.Name,
	}
	fmt.Print(userSess)
	// 5. Llamamos al repositorio pasándole el modelo real y el UserID
	createdClient, err := h.repo.CreateClient(c, newClient, userSess.UserID)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Error al crear el cliente")
	// }
	if err != nil {
		fmt.Printf("🛑 ERROR FATAL EN BASE DE DATOS: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Error DB: "+err.Error())
	}
	// 6. Retornamos éxito
	return c.JSON(http.StatusCreated, createdClient)
}

func (h *ClientHandler) Edit(c echo.Context) error {
	// 1. Extraemos la sesión y validamos que esté autenticado
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	// 2. Obtenemos el ID del cliente
	idParam := c.Param("id")
	clientID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	// 3. Leemos el DTO de actualización
	var input InputClient
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Formato JSON inválido")
	}

	// 4. Validamos que cumpla las reglas
	if err := c.Validate(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Faltan campos obligatorios")
	}

	// 5. Mapeamos el DTO al modelo
	updatedClient := models.Client{
		ID:   clientID,
		Name: input.Name,
	}

	// 6. Llamamos al repositorio para actualizar
	client, err := h.repo.UpdateClient(c, updatedClient, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al actualizar el cliente")
	}

	// 7. Retornamos éxito
	return c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Delete(c echo.Context) error {
	// 1. Extraemos la sesión y validamos que esté autenticado
	userSess := middlewares.GetRequestUserSession(c)
	if !userSess.IsAuthenticated {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Usuario no autenticado",
		})
	}

	// 2. Obtenemos el ID del cliente
	idParam := c.Param("id")
	clientID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID inválido")
	}

	// 3. Llamamos al repositorio para eliminar
	err = h.repo.DeleteClient(c, clientID, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error al eliminar el cliente")
	}

	// 4. Retornamos éxito
	return c.JSON(http.StatusOK, map[string]string{"message": "Cliente eliminado correctamente"})
}
