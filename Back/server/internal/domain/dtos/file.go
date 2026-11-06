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
	CanEdit     bool             `json:"can_edit"`
	CanDelete   bool             `json:"can_delete"`
	CanDownload bool             `json:"can_download"`
}

type GrantFileUserAccessRequest struct {
	UserID     int    `json:"user_id"`
	AccessType string `json:"access_type"`
}

type GrantFileGroupAccessRequest struct {
	GroupID    int    `json:"group_id"`
	AccessType string `json:"access_type"`
}

type FileUserAccessResponse struct {
	UserID     int              `json:"user_id"`
	Username   string           `json:"username"`
	AccessType string           `json:"access_type"`
	GrantedBy  int              `json:"granted_by"`
	GrantedAt  utils.CustomTime `json:"granted_at"`
}

type FileGroupAccessResponse struct {
	GroupID    int              `json:"group_id"`
	GroupName  string           `json:"group_name"`
	AccessType string           `json:"access_type"`
	GrantedBy  int              `json:"granted_by"`
	GrantedAt  utils.CustomTime `json:"granted_at"`
}

type FileAccessResponse struct {
	Users  []FileUserAccessResponse  `json:"users"`
	Groups []FileGroupAccessResponse `json:"groups"`
}

type AddGroupFileRequest struct {
	FileID     int    `json:"file_id"`
	AccessType string `json:"access_type"`
}
