package services

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/store/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]dtos.UserResponse, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) SearchUsers(ctx context.Context, searchParam string) ([]dtos.UserResponse, error) {
	return s.repo.SearchUsers(ctx, searchParam)
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*dtos.UserResponse, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *UserService) CreateUser(ctx context.Context, req dtos.CreateUserRequest, createdBy int) (*dtos.CreatedUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateUserWithRole(ctx, req.Username, req.Email, string(hashedPassword), req.RoleID, createdBy)
}

func (s *UserService) UpdateUser(ctx context.Context, req *dtos.UpdateUserRequest, userID, updatedBy int) (*dtos.UserResponse, error) {
	return s.repo.UpdateUser(ctx, req, userID, updatedBy)
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	return s.repo.DeleteUser(ctx, userID)
}
