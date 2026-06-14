package routes

import (
	"server/internal/http/handlers"
	"server/internal/http/handlers/auth"
	"server/internal/http/handlers/clients"
	"server/internal/http/handlers/files"
	"server/internal/http/handlers/projects"
	"server/internal/http/handlers/users"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *pgxpool.Pool) {

	auth.RegisterRoutes(e, db)

	clients.RegisterRoutes(e, db)
	projects.RegisterRoutes(e, db)
	files.RegisterRoutes(e, db)
	users.RegisterRoutes(e, db)
	handlers.RegisterRoutes(e, db)
}
