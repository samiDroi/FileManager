package repository

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository struct {
	db *pgxpool.Pool
}

func NewClientRepository(db *pgxpool.Pool) *ClientRepository {
	return &ClientRepository{db: db}
}

func (c *ClientRepository) GetAllClients(ctx context.Context, userID int) ([]dtos.ClientResponse, error) {
	q := `SELECT c.id, c.name, ci.email, ci.contact_phone, 
          COALESCE((SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id), 0) as count_projects,
          c.created_at
          FROM clients c 
          LEFT JOIN "contact_info" ci ON c.id = ci.client_id
          WHERE ` + clientVisibilityPredicate

	rows, err := c.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []dtos.ClientResponse
	for rows.Next() {
		var client dtos.ClientResponse
		var createdAt time.Time
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.ContactPhone, &client.CountProjects, &createdAt); err != nil {
			return nil, err
		}
		client.CreatedAt = utils.ToCustomTime(createdAt)
		clients = append(clients, client)
	}
	return clients, nil
}

func (c *ClientRepository) SearchClient(ctx context.Context, searchParam string, userID int) ([]dtos.ClientResponse, error) {
	if searchParam == "" {
		return []dtos.ClientResponse{}, nil
	}

	searchPattern := "%" + searchParam + "%"

	query := `
        SELECT c.id, c.name, ci.email, ci.contact_phone, 
               COALESCE((SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id), 0) as count_projects,
               c.created_at
        FROM clients c 
        LEFT JOIN "contact_info" ci ON c.id = ci.client_id
        WHERE c.name ILIKE $2
          AND ` + clientVisibilityPredicate

	rows, err := c.db.Query(ctx, query, userID, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []dtos.ClientResponse
	for rows.Next() {
		var client dtos.ClientResponse
		var createdAt time.Time
		if err := rows.Scan(&client.ID, &client.Name, &client.Email, &client.ContactPhone, &client.CountProjects, &createdAt); err != nil {
			return nil, err
		}
		client.CreatedAt = utils.ToCustomTime(createdAt)
		clients = append(clients, client)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *ClientRepository) CreateClient(ctx context.Context, client models.Client, userID int) (*dtos.ClientResponse, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO clients(name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING id`
	var clientID int
	err = tx.QueryRow(ctx, q, client.Name, userID, userID).Scan(&clientID)
	if err != nil {
		return nil, err
	}

	if client.ContactInfo != nil {
		contactQ := `INSERT INTO "contact_info"(client_id, email, contact_phone, created_by, updated_by) 
                     VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx, contactQ, clientID, client.ContactInfo.Email, client.ContactInfo.ContactPhone, userID, userID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return c.GetClientByID(ctx, clientID)
}

func (c *ClientRepository) UpdateClient(ctx context.Context, client models.Client, userID int) (*dtos.ClientResponse, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := `UPDATE clients SET name = $1, updated_by = $2 WHERE id = $3`
	_, err = tx.Exec(ctx, q, client.Name, userID, client.ID)
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

func (c *ClientRepository) DeleteClient(ctx context.Context, clientID int64, userID int) error {
	q := `DELETE FROM clients WHERE id = $1`
	_, err := c.db.Exec(ctx, q, clientID)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientRepository) GetClientByID(ctx context.Context, clientID int) (*dtos.ClientResponse, error) {
	q := `SELECT c.id, c.name, ci.email, ci.contact_phone, 
          COALESCE((SELECT COUNT(*) FROM "projects" p WHERE p.client_id = c.id), 0) as count_projects,
          c.created_at
          FROM clients c 
          LEFT JOIN "contact_info" ci ON c.id = ci.client_id
          WHERE c.id = $1`

	var client dtos.ClientResponse
	var createdAt time.Time
	err := c.db.QueryRow(ctx, q, clientID).Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.ContactPhone,
		&client.CountProjects,
		&createdAt,
	)
	client.CreatedAt = utils.ToCustomTime(createdAt)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) GetClientsWithProjects(ctx context.Context, userID int) ([]dtos.ClientWithProjectsResponse, error) {
	q := `SELECT c.id, c.name, p.id as project_id, p.name as project_name
		 FROM clients c
		 JOIN projects p ON c.id = p.client_id
		 WHERE ` + clientVisibilityPredicate + `
		   AND ` + projectVisibilityPredicate

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []dtos.ClientWithProjectsResponse
	for rows.Next() {
		//tempory variable to scan the row data
		var client dtos.ClientWithProjectsResponse
		err := rows.Scan(
			&client.ID,
			&client.Name,
			&client.ProjectID,
			&client.ProjectName,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}
