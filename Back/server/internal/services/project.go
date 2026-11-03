package services

import (
	"context"
	"errors"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/store/repository"
	"server/internal/utils"
)

var ErrProjectReload = errors.New("project updated but could not be reloaded")

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]dtos.ProjectResponse, error) {
	return s.repo.GetAllProjects(ctx)
}

func (s *ProjectService) SearchProjects(ctx context.Context, searchParam string) ([]dtos.ProjectResponse, error) {
	return s.repo.Search(ctx, searchParam)
}

func (s *ProjectService) CreateProject(ctx context.Context, req dtos.ProjectRequest, userID int) (*dtos.ProjectResponse, error) {
	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	createdProject, err := s.repo.CreateProject(ctx, &projectModel, userID)
	if err != nil {
		return nil, err
	}

	clientName, _ := s.repo.GetClientName(ctx, createdProject.ClientID)

	initialFiles := 0
	response := dtos.ProjectResponse{
		ID:         createdProject.ID,
		Name:       createdProject.Name,
		ClientID:   createdProject.ClientID,
		ClientName: clientName,
		CountFiles: &initialFiles,
		CreatedAt:  utils.ToCustomTime(createdProject.CreatedAt),
	}

	return &response, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectID int, req dtos.ProjectRequest) (*dtos.ProjectResponse, error) {
	projectModel := models.Project{
		Name:     req.Name,
		ClientID: req.ClientID,
	}

	if err := s.repo.UpdateProject(ctx, &projectModel, projectID); err != nil {
		return nil, err
	}

	updatedData, err := s.repo.GetProjectById(ctx, projectID)
	if err != nil {
		return nil, ErrProjectReload
	}

	return updatedData, nil
}
