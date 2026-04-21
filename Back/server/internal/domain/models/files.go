package models

import "time"

type File struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	ProjectID int    `json:"project_id"`
	TypeID    int    `json:"type_id"`
	Year      int    `json:"year"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}

type FilesToDto struct {
	File         File
	ProjectName  string
	ProjectID  int
	CreatedAt  time.Time
	FileType   string
}
