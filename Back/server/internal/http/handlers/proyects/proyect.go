package proyects

import (
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type proyectHandler struct {
	db   *pgxpool.Pool
	repo *repository.ProyectRepository
}
