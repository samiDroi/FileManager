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
var ErrInvalidAccessType = errors.New("access_type debe ser leer, editar o descargar")
var ErrFileAccessDenied = errors.New("no cuentas con acceso suficiente sobre el archivo")

type FileService struct {
	repo *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) GetAllFiles(ctx context.Context, userID int) ([]dtos.FileResponse, error) {
	return s.repo.GetAllFiles(ctx, userID)
}

func (s *FileService) GetAllFilesByProject(ctx context.Context, projectID int, userID int) ([]dtos.FileResponse, error) {
	return s.repo.GetAllFilesByProject(ctx, projectID, userID)
}

func (s *FileService) SearchFiles(ctx context.Context, searchParam string, userID int) ([]dtos.FileResponse, error) {
	return s.repo.SearchFiles(ctx, searchParam, userID)
}

func (s *FileService) GetFileByID(ctx context.Context, fileID int, userID int) (dtos.FileResponse, error) {
	return s.repo.GetFileByID(ctx, fileID, userID)
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

func (s *FileService) DeleteFile(ctx context.Context, fileID, userID int) error {
	allowed, err := s.repo.CanDeleteFile(ctx, fileID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrFileAccessDenied
	}
	return s.repo.DeleteFile(ctx, fileID)
}

func (s *FileService) EditFile(ctx context.Context, fileID int, req dtos.FileRequest, userID int) (dtos.FileResponse, error) {
	allowed, err := s.repo.CanEditFile(ctx, fileID, userID)
	if err != nil {
		return dtos.FileResponse{}, err
	}
	if !allowed {
		return dtos.FileResponse{}, ErrFileAccessDenied
	}
	return s.repo.EditFile(ctx, fileID, req, userID)
}

func (s *FileService) CanDownloadFile(ctx context.Context, fileID, userID int) (bool, error) {
	return s.repo.CanDownloadFile(ctx, fileID, userID)
}

func (s *FileService) GrantUserAccess(ctx context.Context, fileID, userID int, accessType string, grantedBy int) error {
	if !isValidAccessType(accessType) {
		return ErrInvalidAccessType
	}
	if err := s.requireAccessManagement(ctx, fileID, grantedBy); err != nil {
		return err
	}
	return s.repo.UpsertUserAccess(ctx, fileID, userID, accessType, grantedBy)
}

func (s *FileService) RevokeUserAccess(ctx context.Context, fileID, userID, grantedBy int) error {
	if err := s.requireAccessManagement(ctx, fileID, grantedBy); err != nil {
		return err
	}
	return s.repo.DeleteUserAccess(ctx, fileID, userID)
}

func (s *FileService) GrantGroupAccess(ctx context.Context, fileID, groupID int, accessType string, grantedBy int) error {
	if !isValidAccessType(accessType) {
		return ErrInvalidAccessType
	}
	if err := s.requireAccessManagement(ctx, fileID, grantedBy); err != nil {
		return err
	}
	return s.repo.UpsertGroupAccess(ctx, fileID, groupID, accessType, grantedBy)
}

func (s *FileService) RevokeGroupAccess(ctx context.Context, fileID, groupID, grantedBy int) error {
	if err := s.requireAccessManagement(ctx, fileID, grantedBy); err != nil {
		return err
	}
	return s.repo.DeleteGroupAccess(ctx, fileID, groupID)
}

func (s *FileService) GetFileAccess(ctx context.Context, fileID, userID int) (*dtos.FileAccessResponse, error) {
	if err := s.requireAccessManagement(ctx, fileID, userID); err != nil {
		return nil, err
	}
	return s.repo.GetFileAccess(ctx, fileID)
}

func (s *FileService) requireAccessManagement(ctx context.Context, fileID, userID int) error {
	allowed, err := s.repo.CanManageFileAccess(ctx, fileID, userID)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrFileAccessDenied
	}
	return nil
}

func isValidAccessType(accessType string) bool {
	switch accessType {
	case "leer", "editar", "descargar":
		return true
	default:
		return false
	}
}
