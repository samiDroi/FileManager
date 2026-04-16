package models

import "time"

type Proyect struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ClientID   int       `json:"client_id"`
	CountFiles *int      `json:"count_files,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  int       `json:"created_by"`
	UpdatedBy  int       `json:"updated_by"`
}
