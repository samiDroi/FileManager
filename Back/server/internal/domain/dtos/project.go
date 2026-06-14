package dtos

import (
	"server/internal/utils"
)

type ProjectRequest struct {
	Name     string `json:"name"`
	ClientID int    `json:"client_id"`
}

type ProjectResponse struct {
	ID          int               `json:"id"`
	Name       string            `json:"name"`
	ClientID   int               `json:"client_id"`
	ClientName string            `json:"client_name"`
	CountFiles *int             `json:"count_files"`
	CreatedAt  utils.CustomTime `json:"created_at"`
}