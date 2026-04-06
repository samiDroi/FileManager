package routes

import (
	"server/internal/http/handlers/auth"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *pgxpool.Pool) {
	auth.RegisterRoutes(e, db)
}
