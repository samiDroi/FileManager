package models

type Client struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}
