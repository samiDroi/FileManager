package clients

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
	DB *pgxpool.Pool
}

func (h *ClientHandler) CreateHandler(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"msg": "Cliente creado"})
}
