package clients

type InputClient struct {
	Name string `json:"name" validate:"required"`
}
