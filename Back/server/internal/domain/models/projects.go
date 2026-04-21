package models

import "time"

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	ClientID  int       `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
}
type ProjectToDto struct {
	ID         int       `json:"id"`
	Name      string    `json:"name"`
	ClientID  int       `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
	ClientName string   `json:"client_name"`
	CountFiles *int     `json:"count_files,omitempty"`
}
