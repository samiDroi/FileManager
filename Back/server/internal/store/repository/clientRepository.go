package repository

import (
	"context"
	"server/internal/domain/models"
	"server/internal/utils"
	"time"

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
	q := `SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
		  (SELECT COUNT(*) FROM "Proyect" p WHERE p.client_id = c.id) as count_proyects
		  FROM "Client" c 
		  LEFT JOIN "Contact_Info" ci ON c.id = ci.client_id`
	rows, err := c.db.Query(ctx.Request().Context(), q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []*models.Client{}
	for rows.Next() {
		var client models.Client
		var email, phone *string
		var countProyects *int
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&client.ID, &client.Name, &createdAt, &updatedAt, &email, &phone, &countProyects); err != nil {
			return nil, err
		}
		client.CreatedAt = utils.ToCustomTime(createdAt)
		client.UpdatedAt = utils.ToCustomTime(updatedAt)
		if email != nil || phone != nil {
			client.ContactInfo = &models.ContactInfo{
				ClientID:     client.ID,
				Email:        email,
				ContactPhone: phone,
			}
		}
		client.CountProyects = countProyects
		clients = append(clients, &client)
	}
	return clients, nil
}

func (c *ClientRepository) CreateClient(ctx echo.Context, client models.Client, userID int) (*models.Client, error) {
	tx, err := c.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx.Request().Context())

	q := `INSERT INTO "Client"(name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx.Request().Context(), q, client.Name, userID, userID).Scan(&client.ID)
	if err != nil {
		return nil, err
	}

	if client.ContactInfo != nil {
		contactQ := `INSERT INTO "Contact_Info"(client_id, email, contact_phone, created_by, updated_by) 
					 VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx.Request().Context(), contactQ, client.ID, client.ContactInfo.Email, client.ContactInfo.ContactPhone, userID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx.Request().Context()); err != nil {
		return nil, err
	}

	return c.GetClientByID(ctx, client.ID)
}

func (c *ClientRepository) UpdateClient(ctx echo.Context, client models.Client, userID int) (*models.Client, error) {
	tx, err := c.db.Begin(ctx.Request().Context())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx.Request().Context())

	q := `UPDATE "Client" SET name = $1, updated_by = $2 WHERE id = $3 RETURNING id`
	err = tx.QueryRow(ctx.Request().Context(), q, client.Name, userID, client.ID).Scan(&client.ID)
	if err != nil {
		return nil, err
	}

	if client.ContactInfo != nil {
		contactQ := `INSERT INTO "Contact_Info"(client_id, email, contact_phone, created_by, updated_by)
					 VALUES ($1, $2, $3, $4, $5)
					 ON CONFLICT (client_id) DO UPDATE SET 
						email = EXCLUDED.email, 
						contact_phone = EXCLUDED.contact_phone,
						updated_by = EXCLUDED.updated_by`
		_, err = tx.Exec(ctx.Request().Context(), contactQ, client.ID, client.ContactInfo.Email, client.ContactInfo.ContactPhone, userID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx.Request().Context()); err != nil {
		return nil, err
	}

	return c.GetClientByID(ctx, client.ID)
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
		return results, nil
	}

	searchPattern := "%" + searchParam + "%"
	query := `SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
			  (SELECT COUNT(*) FROM "Proyect" p WHERE p.client_id = c.id) as count_proyects
			  FROM "Client" c 
			  LEFT JOIN "Contact_Info" ci ON c.id = ci.client_id
			  WHERE c.name ILIKE $1`

	rows, err := c.db.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var client models.Client
		var email, phone *string
		var countProyects *int
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&client.ID, &client.Name, &createdAt, &updatedAt, &email, &phone, &countProyects); err != nil {
			return nil, err
		}
		client.CreatedAt = utils.ToCustomTime(createdAt)
		client.UpdatedAt = utils.ToCustomTime(updatedAt)
		if email != nil || phone != nil {
			client.ContactInfo = &models.ContactInfo{
				ClientID:     client.ID,
				Email:        email,
				ContactPhone: phone,
			}
		}
		client.CountProyects = countProyects
		results = append(results, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (c *ClientRepository) GetClientByID(ctx echo.Context, clientID int) (*models.Client, error) {
	q := `SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
		  (SELECT COUNT(*) FROM "Proyect" p WHERE p.client_id = c.id) as count_proyects
		  FROM "Client" c 
		  LEFT JOIN "Contact_Info" ci ON c.id = ci.client_id
		  WHERE c.id = $1`

	var client models.Client
	var email, phone *string
	var countProyects *int
	var createdAt, updatedAt time.Time

	err := c.db.QueryRow(ctx.Request().Context(), q, clientID).Scan(
		&client.ID, &client.Name, &createdAt, &updatedAt, &email, &phone, &countProyects)
	if err != nil {
		return nil, err
	}

	client.CreatedAt = utils.ToCustomTime(createdAt)
	client.UpdatedAt = utils.ToCustomTime(updatedAt)

	if email != nil || phone != nil {
		client.ContactInfo = &models.ContactInfo{
			ClientID:     client.ID,
			Email:        email,
			ContactPhone: phone,
		}
	}
	client.CountProyects = countProyects

	return &client, nil
}

// func (c *ClientRepository) countProyects(ctx echo.Context, clientID int) (int, error) {
// 	var total int
// 	q := `SELECT COUNT(*) FROM Proyect WHERE client_id = $1`

// 	err := c.db.QueryRow(ctx.Request().Context(), q, clientID).Scan(&total)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return total, err

// }
