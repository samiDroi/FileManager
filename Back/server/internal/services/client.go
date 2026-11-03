package services

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/store/repository"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(repo *repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) GetAllClients(ctx context.Context) ([]dtos.ClientResponse, error) {
	return s.repo.GetAllClients(ctx)
}

func (s *ClientService) SearchClient(ctx context.Context, searchParam string) ([]dtos.ClientResponse, error) {
	return s.repo.SearchClient(ctx, searchParam)
}

func (s *ClientService) CreateClient(ctx context.Context, req dtos.ClientRequest, userID int) (*dtos.ClientResponse, error) {
	var contactInfo *models.ContactInfo
	if req.Email != nil || req.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			Email:        req.Email,
			ContactPhone: req.ContactPhone,
		}
	}

	newClient := models.Client{
		Name:        req.Name,
		ContactInfo: contactInfo,
	}

	return s.repo.CreateClient(ctx, newClient, userID)
}

func (s *ClientService) UpdateClient(ctx context.Context, clientID int, req dtos.ClientRequest, userID int) (*dtos.ClientResponse, error) {
	var contactInfo *models.ContactInfo
	if req.Email != nil || req.ContactPhone != nil {
		contactInfo = &models.ContactInfo{
			ClientID:     clientID,
			Email:        req.Email,
			ContactPhone: req.ContactPhone,
		}
	}

	updatedClient := models.Client{
		ID:          clientID,
		Name:        req.Name,
		ContactInfo: contactInfo,
	}

	return s.repo.UpdateClient(ctx, updatedClient, userID)
}

func (s *ClientService) DeleteClient(ctx context.Context, clientID int64, userID int) error {
	return s.repo.DeleteClient(ctx, clientID, userID)
}

func (s *ClientService) GetClientsWithProjects(ctx context.Context) ([]dtos.ClientWithProjectsResponse, error) {
	return s.repo.GetClientsWithProjects(ctx)
}
