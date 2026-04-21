package clients

type InputClient struct {
	Name string `json:"name" validate:"required"`
	//optional Data
	ContactPhone  *string `json:"contact_phone,omitempty"`
	Email         *string `json:"email,omitempty"`
	Countprojects *int    `json:"countprojects,omitempty"`
}
