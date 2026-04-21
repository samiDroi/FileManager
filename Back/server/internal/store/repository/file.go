package repository

import (
	"context"
	"server/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

// constructor
func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) GetAllFileTypes(ctx context.Context) ([]*models.FileType, error) {
	q := "SELECT id,name FROM File_Type"
	rows, err := r.db.Query(ctx, q)

	if err != nil {
		return nil, err
	}
	//closing the rows after we finish processing them
	defer rows.Close()
	//adapting the data to the struct model
	fileTypes := []*models.FileType{}
	for rows.Next() {
		fileType := &models.FileType{}
		if err := rows.Scan(&fileType.ID, &fileType.Name); err != nil {
			return nil, err
		}
		fileTypes = append(fileTypes, fileType)
	}
	return fileTypes, nil

}
func (r *FileRepository) GetAllFiles(ctx context.Context) ([]*models.FilesToDto, error) {
	q := `SELECT f.id, f.title, f.path, f.project_id, f.type_id, f.year, f.created_by, f.updated_by, f.created_at, p.id, p.name, ft.name
			FROM files f
			LEFT JOIN projects p ON p.id = f.project_id
			LEFT JOIN file_types ft ON ft.id = f.type_id
		`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := []*models.FilesToDto{}
	for rows.Next() {
		var f models.FilesToDto
		err := rows.Scan(
			&f.File.ID,
			&f.File.Title,
			&f.File.Path,
			&f.File.ProjectID,
			&f.File.TypeID,
			&f.File.Year,
			&f.File.CreatedBy,
			&f.File.UpdatedBy,
			&f.CreatedAt,
			&f.ProjectID,
			&f.ProjectName,
			&f.FileType,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	return files, nil
}
