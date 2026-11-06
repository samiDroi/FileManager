package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"server/internal/config"
	"server/internal/domain/dtos"
	"server/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const fileCapabilitySelect = `,
	` + fileEditPredicate + ` AS can_edit,
	` + fileDeletePredicate + ` AS can_delete,
	` + fileDownloadPredicate + ` AS can_download`

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) GetAllFileTypes(ctx context.Context) ([]*dtos.FileResponse, error) {
	q := "SELECT id, name FROM File_Type"
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dtos.FileResponse
	for rows.Next() {
		var f dtos.FileResponse
		if err := rows.Scan(&f.ID, &f.FileType); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return toPtrSlice(files), nil
}

func toPtrSlice[T any](slice []T) []*T {
	result := make([]*T, len(slice))
	for i := range slice {
		result[i] = &slice[i]
	}
	return result
}

func (r *FileRepository) GetAllFiles(ctx context.Context, userID int) ([]dtos.FileResponse, error) {
	q := `SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
			  COALESCE(p.name, '') as project_name,
			  COALESCE(ft.name, '') as file_type` + fileCapabilitySelect + `
		  FROM files f
		  LEFT JOIN projects p ON p.id = f.project_id
		  LEFT JOIN file_types ft ON ft.id = f.type_id
		  WHERE ` + fileVisibilityPredicate

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dtos.FileResponse
	for rows.Next() {
		var f dtos.FileResponse
		var createdAt time.Time
		err := rows.Scan(
			&f.ID,
			&f.Tittle,
			&f.Path,
			&f.ProjectID,
			&f.TypeID,
			&f.Year,
			&createdAt,
			&f.ProjectName,
			&f.FileType,
			&f.CanEdit,
			&f.CanDelete,
			&f.CanDownload,
		)
		if err != nil {
			return nil, err
		}
		f.CreatedAt = utils.ToCustomTime(createdAt)
		files = append(files, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

func (r *FileRepository) GetAllFilesByProject(ctx context.Context, id int, userID int) ([]dtos.FileResponse, error) {
	q := `SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
			  COALESCE(p.name, '') as project_name,
			  COALESCE(ft.name, '') as file_type` + fileCapabilitySelect + `
		  FROM files f
		  LEFT JOIN projects p ON p.id = f.project_id
		  LEFT JOIN file_types ft ON ft.id = f.type_id
		  WHERE ` + fileVisibilityPredicate + `
		    AND f.project_id = $2`

	rows, err := r.db.Query(ctx, q, userID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dtos.FileResponse
	for rows.Next() {
		var f dtos.FileResponse
		var createdAt time.Time
		err := rows.Scan(
			&f.ID,
			&f.Tittle,
			&f.Path,
			&f.ProjectID,
			&f.TypeID,
			&f.Year,
			&createdAt,
			&f.ProjectName,
			&f.FileType,
			&f.CanEdit,
			&f.CanDelete,
			&f.CanDownload,
		)
		if err != nil {
			return nil, err
		}
		f.CreatedAt = utils.ToCustomTime(createdAt)
		files = append(files, f)
	}
	return files, nil

}

func (r *FileRepository) GetFileByID(ctx context.Context, id int, userID int) (dtos.FileResponse, error) {
	q := `SELECT f.tittle, f.path` + fileCapabilitySelect + `
		  FROM files f
		  WHERE f.id = $2
		    AND ` + fileVisibilityPredicate

	var f dtos.FileResponse
	err := r.db.QueryRow(ctx, q, userID, id).Scan(
		&f.Tittle,
		&f.Path,
		&f.CanEdit,
		&f.CanDelete,
		&f.CanDownload,
	)
	if err != nil {
		return dtos.FileResponse{}, err
	}
	return f, nil
}

func (r *FileRepository) CanEditFile(ctx context.Context, fileID, userID int) (bool, error) {
	return r.fileCapability(ctx, fileEditPredicate, fileID, userID)
}

func (r *FileRepository) CanDeleteFile(ctx context.Context, fileID, userID int) (bool, error) {
	return r.fileCapability(ctx, fileDeletePredicate, fileID, userID)
}

func (r *FileRepository) CanDownloadFile(ctx context.Context, fileID, userID int) (bool, error) {
	return r.fileCapability(ctx, fileDownloadPredicate, fileID, userID)
}

func (r *FileRepository) CanManageFileAccess(ctx context.Context, fileID, userID int) (bool, error) {
	return r.fileCapability(ctx, fileDeletePredicate, fileID, userID)
}

func (r *FileRepository) fileCapability(ctx context.Context, predicate string, fileID, userID int) (bool, error) {
	var allowed bool
	err := r.db.QueryRow(ctx,
		`SELECT `+predicate+`
		 FROM files f
		 WHERE f.id = $2`,
		userID, fileID,
	).Scan(&allowed)
	return allowed, err
}

func (r *FileRepository) SaveFile(ctx context.Context, req dtos.FileRequest, filename string, src io.Reader, userID int) (dtos.FileResponse, error) {
	uploadPath := config.Config.UploadPath
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to resolve upload path: %w", err)
	}
	if err := os.MkdirAll(absUploadPath, 0755); err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to create directory: %w", err)
	}
	filePath := filepath.Join(absUploadPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to save file to disk: %w", err)
	}

	// Convertimos el Year que dejaste como string en el DTO a int (usando strconv)
	yearInt := 0
	if req.Year != "" {
		yearInt, _ = strconv.Atoi(req.Year)
	}

	// Si no envían TypeID, por defecto usamos 1
	typeID := req.TypeID
	if typeID == 0 {
		typeID = 1
	}

	var newFileID int
	var createdAt time.Time // Declaramos esto para guardar la fecha que genere la BD

	// INSERTAMOS. Usamos req.TypeID ($4).
	// Agregamos RETURNING id, created_at para llenar el DTO completo
	err = r.db.QueryRow(ctx,
		`INSERT INTO files (tittle, path, project_id, type_id, year, created_by, updated_by) 
         VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		req.Tittle, filePath, req.ProjectID, typeID, yearInt, userID, userID,
	).Scan(&newFileID, &createdAt)

	// Si la BD falla, borramos el archivo físico recién creado
	if err != nil {
		os.Remove(filePath)
		return dtos.FileResponse{}, fmt.Errorf("database error, file discarded: %w", err)
	}

	// Consultar los nombres resueltos con JOIN
	var resp dtos.FileResponse
	err = r.db.QueryRow(ctx,
		`SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
		        COALESCE(p.name, '') as project_name,
		        COALESCE(ft.name, '') as file_type
		 FROM files f
		 LEFT JOIN projects p ON p.id = f.project_id
		 LEFT JOIN file_types ft ON ft.id = f.type_id
		 WHERE f.id = $1`, newFileID,
	).Scan(
		&resp.ID, &resp.Tittle, &resp.Path, &resp.ProjectID,
		&resp.TypeID, &resp.Year, &createdAt,
		&resp.ProjectName, &resp.FileType,
	)
	if err != nil {
		return dtos.FileResponse{
			ID: newFileID, ProjectID: req.ProjectID, Tittle: req.Tittle,
			Path: filePath, Year: yearInt, TypeID: typeID,
			CreatedAt: utils.ToCustomTime(createdAt),
		}, nil
	}
	resp.CreatedAt = utils.ToCustomTime(createdAt)
	resp.CanEdit = true
	resp.CanDelete = true
	resp.CanDownload = true
	return resp, nil
}
func (r *FileRepository) SearchFiles(ctx context.Context, query string, userID int) ([]dtos.FileResponse, error) {
	q := `SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
			  COALESCE(p.name, '') as project_name,
			  COALESCE(ft.name, '') as file_type` + fileCapabilitySelect + `
		  FROM files f
		  LEFT JOIN projects p ON p.id = f.project_id
		  LEFT JOIN file_types ft ON ft.id = f.type_id
		  WHERE ` + fileVisibilityPredicate + `
		    AND (f.tittle ILIKE '%' || $2 || '%' OR
			 p.name ILIKE '%' || $2 || '%' OR
			 ft.name ILIKE '%' || $2 || '%' OR
			 CAST(f.year AS TEXT) ILIKE '%' || $2 || '%')`

	rows, err := r.db.Query(ctx, q, userID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []dtos.FileResponse
	for rows.Next() {
		var f dtos.FileResponse
		var createdAt time.Time
		err := rows.Scan(
			&f.ID,
			&f.Tittle,
			&f.Path,
			&f.ProjectID,
			&f.TypeID,
			&f.Year,
			&createdAt,
			&f.ProjectName,
			&f.FileType,
			&f.CanEdit,
			&f.CanDelete,
			&f.CanDownload,
		)
		if err != nil {
			return nil, err
		}
		f.CreatedAt = utils.ToCustomTime(createdAt)
		files = append(files, f)
	}
	return files, nil

}
func (r *FileRepository) DeleteFile(ctx context.Context, id int) error {
	var filePath string
	err := r.db.QueryRow(ctx, "SELECT path FROM files WHERE id = $1", id).Scan(&filePath)
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}

	_, err = r.db.Exec(ctx, "DELETE FROM files WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete file from database: %w", err)
	}

	if filePath != "" {
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("[WARN] failed to delete file from disk: %v\n", err)
		}
	}

	return nil
}
func (r *FileRepository) EditFile(ctx context.Context, id int, req dtos.FileRequest, userID int) (dtos.FileResponse, error) {
	var currentPath string
	err := r.db.QueryRow(ctx, "SELECT path FROM files WHERE id = $1", id).Scan(&currentPath)
	if err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to get current file path: %w", err)
	}

	uploadBase := config.Config.UploadPath
	if uploadBase == "" {
		uploadBase = "./uploads"
	}
	absUploadBase, _ := filepath.Abs(uploadBase)

	// Renombrar archivo en disco solo si el título cambió realmente
	if req.Tittle != "" {
		ext := filepath.Ext(currentPath)
		baseName := req.Tittle
		if strings.HasSuffix(strings.ToLower(baseName), strings.ToLower(ext)) {
			baseName = strings.TrimSuffix(baseName, ext)
		}
		newName := baseName + ext
		newPath := filepath.Join(filepath.Dir(currentPath), newName)

		if !strings.EqualFold(filepath.Base(currentPath), newName) {
			srcPath := currentPath
			if _, statErr := os.Stat(srcPath); statErr != nil {
				altPath := filepath.Join(absUploadBase, filepath.Base(currentPath))
				if _, altErr := os.Stat(altPath); altErr == nil {
					srcPath = altPath
					newPath = filepath.Join(absUploadBase, newName)
				} else {
					fmt.Printf("[WARN] file not found on disk, skipping rename: %s\n", currentPath)
					currentPath = newPath
				}
			}

			if _, statErr := os.Stat(srcPath); statErr == nil {
				if err := os.Rename(srcPath, newPath); err != nil {
					fmt.Printf("[WARN] failed to rename file on disk, skipping: %v\n", err)
				} else {
					currentPath = newPath
				}
			}
		}
	}

	// Construimos query dinámica para actualizar solo los campos enviados
	var setClauses []string
	var args []interface{}
	argIdx := 1

	if req.Tittle != "" {
		setClauses = append(setClauses, fmt.Sprintf("tittle = $%d", argIdx))
		args = append(args, req.Tittle)
		argIdx++
	}
	if req.ProjectID != 0 {
		setClauses = append(setClauses, fmt.Sprintf("project_id = $%d", argIdx))
		args = append(args, req.ProjectID)
		argIdx++
	}
	if req.TypeID != 0 {
		setClauses = append(setClauses, fmt.Sprintf("type_id = $%d", argIdx))
		args = append(args, req.TypeID)
		argIdx++
	}
	if req.Year != "" {
		yearInt, _ := strconv.Atoi(req.Year)
		setClauses = append(setClauses, fmt.Sprintf("year = $%d", argIdx))
		args = append(args, yearInt)
		argIdx++
	}
	if currentPath != "" {
		setClauses = append(setClauses, fmt.Sprintf("path = $%d", argIdx))
		args = append(args, currentPath)
		argIdx++
	}

	if len(setClauses) > 0 {
		args = append(args, id)
		q := fmt.Sprintf("UPDATE files SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argIdx)
		updated, err := r.db.Exec(ctx, q, args...)
		if err != nil {
			return dtos.FileResponse{}, fmt.Errorf("failed to update file in database: %w", err)
		}
		if updated.RowsAffected() == 0 {
			return dtos.FileResponse{}, fmt.Errorf("no file found with id: %d", id)
		}
	}

	// Reconstruimos la respuesta completa desde la BD
	var resp dtos.FileResponse
	var createdAt time.Time
	q := `SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
			  COALESCE(p.name, '') as project_name,
			  COALESCE(ft.name, '') as file_type` + fileCapabilitySelect + `
		  FROM files f
		  LEFT JOIN projects p ON p.id = f.project_id
		  LEFT JOIN file_types ft ON ft.id = f.type_id
		  WHERE f.id = $2`
	err = r.db.QueryRow(ctx, q, userID, id).Scan(
		&resp.ID, &resp.Tittle, &resp.Path, &resp.ProjectID,
		&resp.TypeID, &resp.Year, &createdAt,
		&resp.ProjectName, &resp.FileType,
		&resp.CanEdit, &resp.CanDelete, &resp.CanDownload,
	)
	if err != nil {
		return dtos.FileResponse{}, fmt.Errorf("failed to read updated file: %w", err)
	}
	resp.CreatedAt = utils.ToCustomTime(createdAt)
	return resp, nil
}

func (r *FileRepository) UpsertUserAccess(ctx context.Context, fileID, userID int, accessType string, grantedBy int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO files_users (file_id, user_id, access_type, granted_by, created_by, updated_by)
		 VALUES ($1, $2, $3::file_access_enum, $4, $4, $4)
		 ON CONFLICT (file_id, user_id) DO UPDATE
		 SET access_type = EXCLUDED.access_type,
		     granted_by = EXCLUDED.granted_by,
		     granted_at = CURRENT_TIMESTAMP,
		     updated_at = CURRENT_TIMESTAMP,
		     updated_by = EXCLUDED.updated_by`,
		fileID, userID, accessType, grantedBy,
	)
	return err
}

func (r *FileRepository) DeleteUserAccess(ctx context.Context, fileID, userID int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM files_users WHERE file_id = $1 AND user_id = $2`,
		fileID, userID,
	)
	return err
}

func (r *FileRepository) UpsertGroupAccess(ctx context.Context, fileID, groupID int, accessType string, grantedBy int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO files_groups (file_id, group_id, access_type, granted_by, created_by, updated_by)
		 VALUES ($1, $2, $3::file_access_enum, $4, $4, $4)
		 ON CONFLICT (file_id, group_id) DO UPDATE
		 SET access_type = EXCLUDED.access_type,
		     granted_by = EXCLUDED.granted_by,
		     granted_at = CURRENT_TIMESTAMP,
		     updated_at = CURRENT_TIMESTAMP,
		     updated_by = EXCLUDED.updated_by`,
		fileID, groupID, accessType, grantedBy,
	)
	return err
}

func (r *FileRepository) DeleteGroupAccess(ctx context.Context, fileID, groupID int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM files_groups WHERE file_id = $1 AND group_id = $2`,
		fileID, groupID,
	)
	return err
}

func (r *FileRepository) GetFileAccess(ctx context.Context, fileID int) (*dtos.FileAccessResponse, error) {
	response := &dtos.FileAccessResponse{
		Users:  make([]dtos.FileUserAccessResponse, 0),
		Groups: make([]dtos.FileGroupAccessResponse, 0),
	}

	userRows, err := r.db.Query(ctx,
		`SELECT fu.user_id, u.username, fu.access_type::text, fu.granted_by, fu.granted_at
		 FROM files_users fu
		 JOIN "user" u ON u.id = fu.user_id
		 WHERE fu.file_id = $1
		 ORDER BY u.username`,
		fileID,
	)
	if err != nil {
		return nil, err
	}
	for userRows.Next() {
		var access dtos.FileUserAccessResponse
		var grantedAt time.Time
		if err := userRows.Scan(&access.UserID, &access.Username, &access.AccessType, &access.GrantedBy, &grantedAt); err != nil {
			userRows.Close()
			return nil, err
		}
		access.GrantedAt = utils.ToCustomTime(grantedAt)
		response.Users = append(response.Users, access)
	}
	if err := userRows.Err(); err != nil {
		userRows.Close()
		return nil, err
	}
	userRows.Close()

	groupRows, err := r.db.Query(ctx,
		`SELECT fg.group_id, g.name, fg.access_type::text, fg.granted_by, fg.granted_at
		 FROM files_groups fg
		 JOIN "groups" g ON g.id = fg.group_id
		 WHERE fg.file_id = $1
		 ORDER BY g.name`,
		fileID,
	)
	if err != nil {
		return nil, err
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var access dtos.FileGroupAccessResponse
		var grantedAt time.Time
		if err := groupRows.Scan(&access.GroupID, &access.GroupName, &access.AccessType, &access.GrantedBy, &grantedAt); err != nil {
			return nil, err
		}
		access.GrantedAt = utils.ToCustomTime(grantedAt)
		response.Groups = append(response.Groups, access)
	}
	if err := groupRows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

// func (r *FileRepository) GetFilesByIds(ctx context.Context, ids []int) ([]dtos.FileResponse, error) {
// 	if len(ids) == 0 {
// 		return []dtos.FileResponse{}, nil
// 	}

// 	q := `SELECT f.id, f.tittle, f.path, f.project_id, f.type_id, f.year, f.created_at,
// 			  COALESCE(p.name, '') as project_name,
// 			  COALESCE(ft.name, '') as file_type
// 		  FROM files f
// 		  LEFT JOIN projects p ON p.id = f.project_id
// 		  LEFT JOIN file_types ft ON ft.id = f.type_id
// 		  WHERE f.id = ANY($1)`

// 	rows, err := r.db.Query(ctx, q, ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var files []dtos.FileResponse
// 	for rows.Next() {
// 		var f dtos.FileResponse
// 		var createdAt time.Time
// 		err := rows.Scan(
// 			&f.ID,
// 			&f.tittle,
// 			&f.Path,
// 			&f.ProjectID,
// 			&f.Year,
// 			&createdAt,
// 			&f.ProjectName,
// 			&f.FileType,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		f.CreatedAt = utils.ToCustomTime(createdAt)
// 		files = append(files, f)
// 	}
// 	return files, nil
// }
