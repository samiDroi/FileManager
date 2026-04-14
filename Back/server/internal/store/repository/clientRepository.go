package repository

import (
	"context"
	// "net/http"
	"server/internal/domain/models"

	// "strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

func (c *ClientRepository) GetAllClients(ctx echo.Context) ([]*models.Client, error) {
	q := `SELECT id, name FROM "Client"`
	rows, err := c.db.Query(ctx.Request().Context(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []*models.Client{}
	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

func (c *ClientRepository) CreateClient(ctx echo.Context, client models.Client, userID int) (*models.Client, error) {
	q := `INSERT INTO "Client"(name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING id`

	// Usamos QueryRow que es perfecto para devolver un solo valor (el ID)
	err := c.db.QueryRow(ctx.Request().Context(), q, client.Name, userID, userID).Scan(&client.ID)

	if err != nil {
		// Devolvemos un error estándar de Go
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepository) UpdateClient(ctx echo.Context, client models.Client, userID int) (*models.Client, error) {
	q := `UPDATE "Client" SET name = $1, updated_by = $2 WHERE id = $3 RETURNING id`

	err := c.db.QueryRow(ctx.Request().Context(), q, client.Name, userID, client.ID).Scan(&client.ID)

	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepository) DeleteClient(ctx echo.Context, clientID int64, userID int) error {
	q := `DELETE FROM "Client" WHERE id = $1`

	_, err := c.db.Exec(ctx.Request().Context(), q, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRepository) SearchClient(ctx context.Context, searchParam string) ([]models.Client, error) {
	results := make([]models.Client, 0)

	if searchParam == "" {
		return results, nil // Retorna vacío sin consultar a la BD
	}

	searchPattern := "%" + searchParam + "%"
	query := `SELECT id, name FROM "Client" WHERE name ILIKE $1`

	rows, err := c.db.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name); err != nil {
			return nil, err
		}
		results = append(results, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
