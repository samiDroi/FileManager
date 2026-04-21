package projects

import (
	"server/internal/utils"
)

// what frontend send us
type projectRequest struct {
	Name     string `json:"name"`
	ClientID int    `json:"client_id"`
}

// what back send to the frontend
type projectResponse struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	ClientId   int              `json:"client_id"`
	ClientName string          `json:"client_name"`
	CountFiles *int             `json:"count_files"`
	CreatedAt  utils.CustomTime `json:"created_at"`
}
