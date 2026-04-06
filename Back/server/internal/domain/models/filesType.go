package models

type FileType struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}
