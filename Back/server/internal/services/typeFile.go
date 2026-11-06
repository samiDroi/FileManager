package services

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/store/repository"
)

type TypeFileService struct {
	repo *repository.TypeFileRepository
}

func NewTypeFileService(repo *repository.TypeFileRepository) *TypeFileService {
	return &TypeFileService{repo: repo}
}

func (s *TypeFileService) GetAllTypeFiles(ctx context.Context) ([]*dtos.FileTypeResponse, error) {
	return s.repo.GetAllTypeFiles(ctx)
}

func (s *TypeFileService) BulkCreateTypeFiles(ctx context.Context, req dtos.BulkCreateFileTypeRequest, userID int) ([]*dtos.FileTypeResponse, error) {
	fileTypes := make([]models.FileType, 0, len(req.Types))
	for _, t := range req.Types {
		fileTypes = append(fileTypes, models.FileType{Name: t.Name})
	}

	return s.repo.BulkCreateTypeFiles(ctx, fileTypes, userID)
}
