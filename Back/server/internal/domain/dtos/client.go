package dtos

import (
	"server/internal/utils"
)

type ClientRequest struct {
	Name         string  `json:"name" validate:"required"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Email        *string `json:"email,omitempty"`
}

type ClientResponse struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	ContactPhone  *string          `json:"contact_phone,omitempty"`
	Email         *string          `json:"email,omitempty"`
	CountProjects *int             `json:"count_projects,omitempty"`
	CreatedAt     utils.CustomTime `json:"created_at"`
}
type ClientWithProjectsResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ProjectID   int    `json:"project_id"`
	ProjectName string `json:"project_name"`
}
