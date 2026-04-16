package repository

import (
	// "net/http"
	"server/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ProyectRepository struct {
	db *pgxpool.Pool
}

// constructor
func NewProyectRepository(db *pgxpool.Pool) *ProyectRepository {
	return &ProyectRepository{db: db}
}
func (p *ProyectRepository) getAllProyects(c echo.Context) ([]*models.Proyect, error) {
	q := `SELECT p.id,p.name,p.created_at,p.updated_at,c.name,
			 (SELECT COUNT(*) FROM "File" f WHERE f.proyect_id = p.id) as count_files
			FROM Proyect p 
			LEFT JOIN Client c ON c.id = p.client_id`
	rows, err := p.db.Query(c.Request().Context(), q)
	if err != nil {
		return nil, err
	}
	//close the conection to the pool once we did the consult
	defer rows.Close()

	proyects := []*models.Proyect{}
	for rows.Next() {
		//model struct to fill model's array
		var proyect models.Proyect
		var clientName *string

		if err := rows.Scan(&proyect.ID, &proyect.Name, &proyect.CreatedAt, &proyect.UpdatedAt, &clientName, &proyect.CountFiles); err != nil {
			return nil, err
		}
		//append the row to the general slice
		proyects = append(proyects, &proyect)
	}
	return proyects, nil
}
