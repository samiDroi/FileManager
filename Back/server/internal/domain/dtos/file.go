package dtos

import (
	"server/internal/utils"
)

type FileRequest struct {
	ProjectID int    `form:"project_id"`
	Tittle    string `form:"tittle"`
	Year      string `form:"year"`
	TypeID    int    `form:"type_id"`
}

type FileResponse struct {
	ID          int              `json:"id"`
	ProjectID   int              `json:"project_id"`
	TypeID      int              `json:"type_id"`
	Tittle      string           `json:"tittle"`
	Path        string           `json:"path"`
	Year        int              `json:"year"`
	FileType    string           `json:"file_type"`
	CreatedAt   utils.CustomTime `json:"created_at"`
	ProjectName string           `json:"project_name"`
}
