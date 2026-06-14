package repository

import (
	"context"
	"log"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/utils"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAllProjects(ctx context.Context) ([]dtos.ProjectResponse, error) {
	q := `SELECT p.id, p.name, p.client_id, p.created_at, c.name,
		  (SELECT COUNT(*) FROM "files" f WHERE f.project_id = p.id) as count_files
		  FROM projects p 
		  LEFT JOIN clients c ON c.id = p.client_id`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []dtos.ProjectResponse
	for rows.Next() {
		var p dtos.ProjectResponse
		var createdAt time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.ClientID, &createdAt, &p.ClientName, &p.CountFiles); err != nil {
			return nil, err
		}
		p.CreatedAt = utils.ToCustomTime(createdAt)
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *ProjectRepository) Search(ctx context.Context, searchParam string) ([]dtos.ProjectResponse, error) {
	if searchParam == "" {
		return []dtos.ProjectResponse{}, nil
	}

	searchPattern := "%" + searchParam + "%"

	query := `
        SELECT p.id, p.name, p.client_id, p.created_at, c.name,
 		  (SELECT COUNT(*) FROM "files" f WHERE f.project_id = p.id) as count_files
 		  FROM projects p 
 		  LEFT JOIN clients c ON c.id = p.client_id
        WHERE p.name ILIKE $1 OR c.name ILIKE $1
    `

	rows, err := r.db.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []dtos.ProjectResponse
	for rows.Next() {
		var p dtos.ProjectResponse
		var createdAt time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.ClientID, &createdAt, &p.ClientName, &p.CountFiles); err != nil {
			return nil, err
		}
		p.CreatedAt = utils.ToCustomTime(createdAt)
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project *models.Project, userID int) (*models.Project, error) {
	log.Printf("=== CREATE PROJECT REPO DEBUG ===")
	log.Printf("ClientID: %d, Name: %s, UserID: %d", project.ClientID, project.Name, userID)

	q := `INSERT INTO public."projects" (client_id, name, created_by, updated_by) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`
	var projectID int
	var createdAt time.Time
	err := r.db.QueryRow(ctx, q, project.ClientID, project.Name, userID, userID).Scan(&projectID, &createdAt)
	if err != nil {
		log.Printf("ERROR inserting project: %v", err)
		return nil, err
	}
	log.Printf("Project created with ID: %d", projectID)
	project.ID = projectID
	project.CreatedAt = createdAt
	return project, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project *models.Project, projectID int) error {
	q := `UPDATE projects 
	      SET name = $1, client_id = $2 
	      WHERE id = $3`
	_, err := r.db.Exec(ctx, q, project.Name, project.ClientID, projectID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectRepository) GetProjectById(ctx context.Context, projectID int) (*dtos.ProjectResponse, error) {
	q := `SELECT p.id, p.name, p.client_id, p.created_at, c.name,
		  (SELECT COUNT(*) FROM "files" f WHERE f.project_id = p.id) as count_files
		  FROM projects p 
		  LEFT JOIN clients c ON c.id = p.client_id
		  WHERE p.id = $1`

	var p dtos.ProjectResponse
	var createdAt time.Time
	err := r.db.QueryRow(ctx, q, projectID).Scan(
		&p.ID,
		&p.Name,
		&p.ClientID,
		&createdAt,
		&p.ClientName,
		&p.CountFiles,
	)
	p.CreatedAt = utils.ToCustomTime(createdAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepository) GetClientName(ctx context.Context, clientID int) (string, error) {
	var clientName string
	err := r.db.QueryRow(ctx, `SELECT name FROM clients WHERE id = $1`, clientID).Scan(&clientName)
	if err != nil {
		return "", err
	}
	return clientName, nil
}