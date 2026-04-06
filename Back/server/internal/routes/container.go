package routes

import (
	"server/internal/http/handlers/auth"
	"server/internal/services"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Container es un contenedor centralizado de todas las dependencias
// Así evitas repetir la inicialización en setup.go
type Container struct {
	// ============ REPOSITORIOS ============
	UserRepo *repository.UserRepository

	// ============ SERVICIOS ============
	AuthService *services.AuthService

	// ============ HANDLERS ============
	AuthHandler *auth.AuthHandler
}

// NewContainer inicializa todas las dependencias de una sola vez
func NewContainer(db *pgxpool.Pool) *Container {
	// 1. Crear repositorios
	userRepo := repository.NewUserRepository(db)

	// 2. Crear servicios (pasando los repositorios necesarios)
	authService := services.NewAuthService(userRepo)

	// 3. Crear handlers (pasando los servicios necesarios)
	authHandler := auth.NewAuthHandler(authService)

	// 4. Devolver todo empaquetado
	return &Container{
		UserRepo:    userRepo,
		AuthService: authService,
		AuthHandler: authHandler,
	}
}
