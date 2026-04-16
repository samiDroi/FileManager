package proyects

//what frontend send us
type proyectRequest struct {
	Name string `json:"name" validate:"required"`
}

//what back send to the frontend
type proyectResponse struct {
}
