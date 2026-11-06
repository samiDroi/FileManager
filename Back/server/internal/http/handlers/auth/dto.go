package auth

// angular request body for registration/login
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// Angular response after registration/login
type UserResponse struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	AccessToken string   `json:"access_token"`
	IsAdmin     bool     `json:"is_admin"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type TokenResponse struct {
	AccessToken string   `json:"access_token"`
	IsAdmin     bool     `json:"is_admin"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
