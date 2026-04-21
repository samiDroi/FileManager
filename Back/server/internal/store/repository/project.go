package repository

import (
	"context"
	"log"
	"server/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

// constructor
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}
func (r *ProjectRepository) GetAllProjects(c context.Context) ([]*models.ProjectToDto, error) {
	q := `SELECT p.id, p.name, p.client_id, p.created_at, c.name,
			 (SELECT COUNT(*) FROM "files" f WHERE f.project_id = p.id) as count_files
			FROM projects p 
			LEFT JOIN clients c ON c.id = p.client_id`
	rows, err := r.db.Query(c, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	Projects := []*models.ProjectToDto{}
	for rows.Next() {
		var Project models.ProjectToDto

		if err := rows.Scan(&Project.ID, &Project.Name, &Project.ClientID, &Project.CreatedAt, &Project.ClientName, &Project.CountFiles); err != nil {
			return nil, err
		}
		Projects = append(Projects, &Project)
	}
	return Projects, nil
}
func (r *ProjectRepository) Search(ctx context.Context, searchParam string) ([]*models.ProjectToDto, error) {
	results := []*models.ProjectToDto{}

	if searchParam == "" {
		return results, nil
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

	for rows.Next() {
		var p models.ProjectToDto

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.ClientID,
			&p.CreatedAt,
			&p.ClientName,
			&p.CountFiles,
		); err != nil {
			return nil, err
		}
		results = append(results, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *ProjectRepository) CreateProject(c context.Context, project *models.Project, userID int) (*models.Project, error) {
	log.Printf("=== CREATE PROJECT REPO DEBUG ===")
	log.Printf("ClientID: %d, Name: %s, UserID: %d", project.ClientID, project.Name, userID)

	q := `INSERT INTO public."projects" (client_id, name, created_by, updated_by) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at`
	var projectID int
	var createdAt time.Time
	err := r.db.QueryRow(c, q, project.ClientID, project.Name, userID, userID).Scan(&projectID, &createdAt)
	if err != nil {
		log.Printf("ERROR inserting project: %v", err)
		return nil, err
	}
	log.Printf("Project created with ID: %d", projectID)
	project.ID = projectID
	project.CreatedAt = createdAt
	return project, nil
}

func (r *ProjectRepository) UpdateProject(c context.Context, project *models.Project, projectID int) error {
	q := `UPDATE projects 
	      SET name = $1, client_id = $2 
	      WHERE id = $3`
	_, err := r.db.Exec(c, q, project.Name, project.ClientID, projectID)

	if err != nil {
		return err
	}

	return nil
}
func (r *ProjectRepository) GetProjectById(c context.Context, projectID int) (*models.ProjectToDto, error) {
	q := `SELECT p.id, p.name, p.client_id, p.created_at, c.name,
		(SELECT COUNT(*) FROM "files" f WHERE f.project_id = p.id) as count_files
		FROM projects p 
		LEFT JOIN clients c ON c.id = p.client_id
		WHERE p.id = $1`
	var project models.ProjectToDto
	err := r.db.QueryRow(c, q, projectID).Scan(
		&project.ID,
		&project.Name,
		&project.ClientID,
		&project.CreatedAt,
		&project.ClientName,
		&project.CountFiles,
	)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetClientName(c context.Context, clientID int) (string, error) {
	var clientName string
	err := r.db.QueryRow(c, `SELECT name FROM clients WHERE id = $1`, clientID).Scan(&clientName)
	if err != nil {
		return "", err
	}
	return clientName, nil
}
