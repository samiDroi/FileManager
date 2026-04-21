package files

import (
	"server/internal/utils"
)

type FileRequest struct {
	ProjectID int    `json:"project_id"`
	Title     string `json:"title"`
	Year      string `json:"year"`
}

type FileResponse struct {
	ID          int              `json:"id" `
	ProjectID   int              `json:"project_id"`
	Title       string           `json:"title"`
	Year        int              `json:"year"`
	FileType   string           `json:"file_type"`
	CreatedAt  utils.CustomTime `json:"created_at"`
	ProjectName string         `json:"project_name"`
}
