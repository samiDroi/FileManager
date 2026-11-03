package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"server/internal/domain/dtos"
	"server/internal/store/repository"
	"time"
)

var ErrInvalidFileData = errors.New("project_id, tittle and year are required")

type FileService struct {
	repo *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) GetAllFiles(ctx context.Context) ([]dtos.FileResponse, error) {
	return s.repo.GetAllFiles(ctx)
}

func (s *FileService) GetAllFilesByProject(ctx context.Context, projectID int) ([]dtos.FileResponse, error) {
	return s.repo.GetAllFilesByProject(ctx, projectID)
}

func (s *FileService) SearchFiles(ctx context.Context, searchParam string) ([]dtos.FileResponse, error) {
	return s.repo.SearchFiles(ctx, searchParam)
}

func (s *FileService) GetFileByID(ctx context.Context, fileID int) (dtos.FileResponse, error) {
	return s.repo.GetFileByID(ctx, fileID)
}

func (s *FileService) ValidateUpload(req dtos.FileRequest) error {
	if req.ProjectID == 0 || req.Tittle == "" || req.Year == "" {
		return ErrInvalidFileData
	}
	return nil
}

func (s *FileService) UploadFile(ctx context.Context, req dtos.FileRequest, originalFilename string, src io.Reader, userID int) (dtos.FileResponse, error) {
	if err := s.ValidateUpload(req); err != nil {
		return dtos.FileResponse{}, err
	}

	ext := filepath.Ext(originalFilename)
	safeFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	return s.repo.SaveFile(ctx, req, safeFileName, src, userID)
}

func (s *FileService) DeleteFile(ctx context.Context, fileID int) error {
	return s.repo.DeleteFile(ctx, fileID)
}

func (s *FileService) EditFile(ctx context.Context, fileID int, req dtos.FileRequest) (dtos.FileResponse, error) {
	return s.repo.EditFile(ctx, fileID, req)
}
