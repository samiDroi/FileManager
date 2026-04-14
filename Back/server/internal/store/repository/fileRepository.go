package repository

import (
	"server/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type FileRepository struct {
	db *pgxpool.Pool
}

// constructor
func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (f *FileRepository) getAllFileTypes(ctx echo.Context) ([]*models.FileType, error) {
	q := "SELECT id,name FROM File_Type"
	rows, err := f.db.Query(ctx.Request().Context(), q)

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
