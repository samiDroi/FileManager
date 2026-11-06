package services

import (
	"context"
	"errors"
	"strings"

	"server/internal/domain/dtos"
	"server/internal/store/repository"
)

var ErrInvalidGroupData = errors.New("group name is required")
var ErrInvalidGroupRole = errors.New("invalid group role")
var ErrGroupAccessDenied = errors.New("no cuentas con acceso suficiente sobre el grupo")

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CanManage(ctx context.Context, userID int) (bool, error) {
	return s.repo.CanManage(ctx, userID)
}

func (s *GroupService) GetAllGroups(ctx context.Context, userID int) ([]dtos.GroupResponse, error) {
	canManageAll, err := s.repo.CanManage(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAllGroups(ctx, userID, canManageAll)
}

func (s *GroupService) GetGroupByID(ctx context.Context, groupID, userID int) (*dtos.GroupResponse, error) {
	group, err := s.repo.GetGroupByID(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if group.RoleInGroup == "" && !group.CanManage {
		return nil, ErrGroupAccessDenied
	}
	return group, nil
}

func (s *GroupService) CreateGroup(ctx context.Context, req dtos.GroupRequest, userID int) (*dtos.GroupResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrInvalidGroupData
	}
	return s.repo.CreateGroup(ctx, req, userID)
}

func (s *GroupService) UpdateGroup(ctx context.Context, groupID int, req dtos.GroupRequest, userID int) (*dtos.GroupResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrInvalidGroupData
	}
	if err := s.requireGroupManager(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.UpdateGroup(ctx, groupID, req, userID)
}

func (s *GroupService) DeleteGroup(ctx context.Context, groupID, userID int) error {
	canManage, err := s.repo.CanManage(ctx, userID)
	if err != nil {
		return err
	}
	if !canManage {
		return ErrGroupAccessDenied
	}
	return s.repo.DeleteGroup(ctx, groupID)
}

func (s *GroupService) GetGroupUsers(ctx context.Context, groupID, userID int) ([]dtos.GroupUserResponse, error) {
	if err := s.requireGroupViewer(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetGroupUsers(ctx, groupID)
}

func (s *GroupService) GetAvailableUsers(ctx context.Context, groupID, userID int) ([]dtos.AvailableUserResponse, error) {
	if err := s.requireGroupManager(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetAvailableUsers(ctx, groupID)
}

func (s *GroupService) AddGroupUser(ctx context.Context, groupID int, req dtos.AddGroupUserRequest, userID int) (*dtos.GroupUserResponse, error) {
	if req.UserID <= 0 {
		return nil, errors.New("user is required")
	}
	if err := s.requireGroupManager(ctx, groupID, userID); err != nil {
		return nil, err
	}

	roleInGroup := strings.ToLower(strings.TrimSpace(req.RoleInGroup))
	if roleInGroup == "" {
		roleInGroup = "lector"
	}
	switch roleInGroup {
	case "admin", "editor", "lector":
	default:
		return nil, ErrInvalidGroupRole
	}

	return s.repo.AddGroupUser(ctx, groupID, req.UserID, roleInGroup, userID)
}

func (s *GroupService) RemoveGroupUser(ctx context.Context, groupID int, targetUserID, userID int) error {
	if err := s.requireGroupManager(ctx, groupID, userID); err != nil {
		return err
	}
	return s.repo.RemoveGroupUser(ctx, groupID, targetUserID)
}

func (s *GroupService) GetGroupFiles(ctx context.Context, groupID, userID int) ([]dtos.GroupFileResponse, error) {
	if err := s.requireGroupViewer(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetGroupFiles(ctx, groupID)
}

func (s *GroupService) GetGroupFilesForManagement(ctx context.Context, groupID, userID int) ([]dtos.GroupFileResponse, error) {
	if err := s.requireGroupFileManager(ctx, groupID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetGroupFiles(ctx, groupID)
}

func (s *GroupService) CanEditGroupFiles(ctx context.Context, groupID, userID int) (bool, error) {
	return s.repo.CanEditGroupFiles(ctx, groupID, userID)
}

func (s *GroupService) RevokeGroupFile(ctx context.Context, groupID, fileID, userID int) error {
	if err := s.requireGroupFileManager(ctx, groupID, userID); err != nil {
		return err
	}
	return s.repo.RevokeGroupFile(ctx, groupID, fileID)
}

func (s *GroupService) requireGroupViewer(ctx context.Context, groupID, userID int) error {
	role, err := s.repo.GetUserGroupRole(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if role != "" {
		return nil
	}
	canManage, err := s.repo.CanManageGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if canManage {
		return nil
	}
	return ErrGroupAccessDenied
}

func (s *GroupService) requireGroupManager(ctx context.Context, groupID, userID int) error {
	allowed, err := s.repo.CanManageGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrGroupAccessDenied
	}
	return nil
}

func (s *GroupService) requireGroupFileManager(ctx context.Context, groupID, userID int) error {
	allowed, err := s.repo.CanEditGroupFiles(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrGroupAccessDenied
	}
	return nil
}
