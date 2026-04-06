package models

type Proyect struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ClientID  int    `json:"client_id"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}
