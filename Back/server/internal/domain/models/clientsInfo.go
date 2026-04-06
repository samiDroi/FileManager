package models

type ClientsInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	ClientID  int    `json:"client_id"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}
