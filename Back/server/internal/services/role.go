package services

import (
	"context"
	"errors"
	"strings"

	"server/internal/domain/dtos"
	"server/internal/store/repository"
)

var ErrInvalidRoleData = errors.New("role name is required")
var ErrSystemRole = errors.New("the Administrator role cannot be deleted")
var ErrSystemRoleRename = errors.New("the Administrator role cannot be renamed")

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CanManage(ctx context.Context, userID int) (bool, error) {
	return s.repo.CanManage(ctx, userID)
}

func (s *RoleService) GetPermissions(ctx context.Context) ([]dtos.PermissionResponse, error) {
	return s.repo.GetPermissions(ctx)
}

func (s *RoleService) GetRoles(ctx context.Context) ([]dtos.RoleResponse, error) {
	return s.repo.GetRoles(ctx)
}

func (s *RoleService) CreateRole(ctx context.Context, req dtos.RoleRequest, userID int) (*dtos.RoleResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrInvalidRoleData
	}
	return s.repo.CreateRole(ctx, req, userID)
}

func (s *RoleService) UpdateRole(ctx context.Context, roleID int, req dtos.RoleRequest, userID int) (*dtos.RoleResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrInvalidRoleData
	}

	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if isSystemRole(role.Name) && !strings.EqualFold(role.Name, req.Name) {
		return nil, ErrSystemRoleRename
	}

	return s.repo.UpdateRole(ctx, roleID, req, userID)
}

func isSystemRole(name string) bool {
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	return normalizedName == "administrador" || normalizedName == "admin"
}

func (s *RoleService) DeleteRole(ctx context.Context, roleID int) error {
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return err
	}
	if isSystemRole(role.Name) {
		return ErrSystemRole
	}
	return s.repo.DeleteRole(ctx, roleID)
}

func (s *RoleService) GetUserPermissions(ctx context.Context, userID int) (*dtos.UserPermissionsResponse, error) {
	return s.repo.GetUserPermissions(ctx, userID)
}
