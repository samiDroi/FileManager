package models

type File struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	ProyectID int    `json:"proyect_id"`
	TypeID    int    `json:"type_id"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
}
