package dtos

import "server/internal/utils"

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	//
	// Password      string           `json:"-"`
	CreatedBy     string           `json:"createdBy"`
	UpdatedBy     string           `json:"updatedBy"`
	UploadedFiles int              `json:"uploadedFiles"`
	CreatedAt     utils.CustomTime `json:"created_at"`
	UpdatedAt     utils.CustomTime `json:"updated_at"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLoggedResetPwd struct {
	Email           string `json:"email" validate:"required,email"`
	CurrentPassword string `json:"current_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}
type UserForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}
type UserForgotPassword struct {
	Token string `json:"token" validate:"required"`
	// Email string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
