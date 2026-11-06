package dtos

import "server/internal/utils"

type PermissionResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermissionIDs []int  `json:"permission_ids"`
}

type RoleResponse struct {
	ID          int                  `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   utils.CustomTime     `json:"created_at"`
	UpdatedAt   utils.CustomTime     `json:"updated_at"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupResponse struct {
	ID          int              `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	FileCount   int              `json:"file_count"`
	UserCount   int              `json:"user_count"`
	RoleInGroup string           `json:"role_in_group"`
	CanManage   bool             `json:"can_manage"`
	CreatedAt   utils.CustomTime `json:"created_at"`
}

type GroupUserResponse struct {
	ID          int              `json:"id"`
	Username    string           `json:"username"`
	Email       string           `json:"email"`
	RoleInGroup string           `json:"role_in_group"`
	AddedAt     utils.CustomTime `json:"added_at"`
}

type AddGroupUserRequest struct {
	UserID      int    `json:"user_id"`
	RoleInGroup string `json:"role_in_group"`
}

type AvailableUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GroupFileResponse struct {
	ID          int              `json:"id"`
	Tittle      string           `json:"tittle"`
	ProjectName string           `json:"project_name"`
	FileType    string           `json:"file_type"`
	AccessType  string           `json:"access_type"`
	GrantedAt   utils.CustomTime `json:"granted_at"`
}

type UserPermissionsResponse struct {
	UserID      int      `json:"user_id"`
	IsAdmin     bool     `json:"is_admin"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
