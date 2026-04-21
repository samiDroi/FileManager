package repository

import (
	"context"
	"server/internal/domain/models"

	// Ya no necesitamos importar "server/internal/utils" aquí

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

// ✅ Cambio: ctx context.Context en lugar de echo.Context
func (c *ClientRepository) GetAllClients(ctx context.Context) ([]*models.Client, error) {
	q := `SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
          (SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id) as count_projects
          FROM clients c 
          LEFT JOIN "contact_info" ci ON c.id = ci.client_id`

	rows, err := c.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]*models.Client, 0)
	for rows.Next() {
		var client models.Client
		var email, phone *string

		// ✅ Escaneo directo, sin variables temporales de fecha
		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.CreatedAt,
			&client.UpdatedAt,
			&email,
			&phone,
			&client.Countprojects,
		); err != nil {
			return nil, err
		}

		if email != nil || phone != nil {
			client.ContactInfo = &models.ContactInfo{
				ClientID:     client.ID,
				Email:        email,
				ContactPhone: phone,
			}
		}
		clients = append(clients, &client)
	}
	return clients, nil
}

// ✅ Cambio: ctx context.Context
func (c *ClientRepository) CreateClient(ctx context.Context, client models.Client, userID int) (*models.Client, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO clients(name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(ctx, q, client.Name, userID, userID).Scan(&client.ID)
	if err != nil {
		return nil, err
	}

	if client.ContactInfo != nil {
		contactQ := `INSERT INTO "contact_info"(client_id, email, contact_phone, created_by, updated_by) 
                     VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx, contactQ, client.ID, client.ContactInfo.Email, client.ContactInfo.ContactPhone, userID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return c.GetClientByID(ctx, client.ID)
}

// ✅ Cambio: ctx context.Context
func (c *ClientRepository) UpdateClient(ctx context.Context, client models.Client, userID int) (*models.Client, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `UPDATE clients SET name = $1, updated_by = $2 WHERE id = $3 RETURNING id`
	err = tx.QueryRow(ctx, q, client.Name, userID, client.ID).Scan(&client.ID)
	if err != nil {
		return nil, err
	}

	if client.ContactInfo != nil {
		contactQ := `INSERT INTO "contact_info"(client_id, email, contact_phone, created_by, updated_by)
                     VALUES ($1, $2, $3, $4, $5)
                     ON CONFLICT (client_id) DO UPDATE SET 
                        email = EXCLUDED.email, 
                        contact_phone = EXCLUDED.contact_phone,
                        updated_by = EXCLUDED.updated_by`
		_, err = tx.Exec(ctx, contactQ, client.ID, client.ContactInfo.Email, client.ContactInfo.ContactPhone, userID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return c.GetClientByID(ctx, client.ID)
}

// ✅ Cambio: ctx context.Context
func (c *ClientRepository) DeleteClient(ctx context.Context, clientID int64, userID int) error {
	q := `DELETE FROM clients WHERE id = $1`

	_, err := c.db.Exec(ctx, q, clientID)
	if err != nil {
		return err
	}
	return nil
}

// Este ya estaba bien
func (c *ClientRepository) SearchClient(ctx context.Context, searchParam string) ([]models.Client, error) {
	results := make([]models.Client, 0)

	if searchParam == "" {
		return results, nil
	}

	searchPattern := "%" + searchParam + "%"

	query := `
        SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
               (SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id) as count_projects
        FROM clients c 
        LEFT JOIN "contact_info" ci ON c.id = ci.client_id
        WHERE c.name ILIKE $1
    `

	rows, err := c.db.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var client models.Client
		var email, phone *string

		if err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.CreatedAt,
			&client.UpdatedAt,
			&email,
			&phone,
			&client.Countprojects,
		); err != nil {
			return nil, err
		}

		if email != nil || phone != nil {
			client.ContactInfo = &models.ContactInfo{
				ClientID:     client.ID,
				Email:        email,
				ContactPhone: phone,
			}
		}

		results = append(results, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// ✅ Cambio: ctx context.Context y limpieza del scan
func (c *ClientRepository) GetClientByID(ctx context.Context, clientID int) (*models.Client, error) {
	q := `SELECT c.id, c.name, c.created_at, c.updated_at, ci.email, ci.contact_phone, 
          (SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id) as count_projects
          FROM clients c 
          LEFT JOIN "contact_info" ci ON c.id = ci.client_id
          WHERE c.id = $1`

	var client models.Client
	var email, phone *string

	// ✅ Escaneo directo al struct
	err := c.db.QueryRow(ctx, q, clientID).Scan(
		&client.ID,
		&client.Name,
		&client.CreatedAt,
		&client.UpdatedAt,
		&email,
		&phone,
		&client.Countprojects,
	)
	if err != nil {
		return nil, err
	}

	if email != nil || phone != nil {
		client.ContactInfo = &models.ContactInfo{
			ClientID:     client.ID,
			Email:        email,
			ContactPhone: phone,
		}
	}

	return &client, nil
}
