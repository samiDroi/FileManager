package repository

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TypeFileRepository struct {
	db *pgxpool.Pool
}

func NewTypeFileRepository(db *pgxpool.Pool) *TypeFileRepository {
	return &TypeFileRepository{db: db}
}

func (r *TypeFileRepository) GetAllTypeFiles(ctx context.Context) ([]*dtos.FileTypeResponse, error) {
	q := `SELECT id, name FROM file_types`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var typeFiles []*dtos.FileTypeResponse
	for rows.Next() {
		var typeFile dtos.FileTypeResponse
		if err := rows.Scan(&typeFile.ID, &typeFile.Name); err != nil {
			return nil, err
		}
		typeFiles = append(typeFiles, &typeFile)
	}
	return typeFiles, nil
}

func (r *TypeFileRepository) CreateTypeFile(ctx context.Context, fileType models.FileType, id int) (*dtos.FileTypeResponse, error) {

	q := `INSERT INTO file_types (name,created_by, updated_by) VALUES ($1, $2, $3) RETURNING id, name`
	var typeFile dtos.FileTypeResponse
	err := r.db.QueryRow(ctx, q, fileType.Name, id, id).Scan(&typeFile.ID, &typeFile.Name)

	if err != nil {
		return nil, err
	}
	return &typeFile, nil
}

// BulkCreateTypeFiles inserts multiple file types into the database in a single transaction and returns the inserted records.
func (r *TypeFileRepository) BulkCreateTypeFiles(ctx context.Context, fileTypes []models.FileType, id int) ([]*dtos.FileTypeResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	//if an error occurs, rollback the transaction
	defer tx.Rollback(ctx)

	q := `INSERT INTO file_types (name, created_by, updated_by) VALUES ($1, $2, $3) RETURNING id, name`
	var results []*dtos.FileTypeResponse
	// Loop through the file types array and execute the insert query for each
	for _, ft := range fileTypes {
		var typeFile dtos.FileTypeResponse
		err := tx.QueryRow(ctx, q, ft.Name, id, id).Scan(&typeFile.ID, &typeFile.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, &typeFile)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return results, nil
}
