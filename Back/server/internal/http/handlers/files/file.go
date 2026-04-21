package files

import (
	"net/http"
	"server/internal/domain/models"
	"server/internal/store/repository"
	"server/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	db   *pgxpool.Pool
	repo *repository.FileRepository
}

func NewFileHandler(db *pgxpool.Pool) *FileHandler {
	return &FileHandler{
		db:   db,
		repo: repository.NewFileRepository(db),
	}
}

func RegisterRoutes() {

}
func (h *FileHandler) getAll(c echo.Context) error {
	files, err := h.repo.GetAllFiles(c.Request().Context())

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]error{"error": err})
	}
	filesToMap := make([]FileResponse, len(files))
	for i, file := range files {
		filesToMap[i] = MapToDto(file)
	}
	return c.JSON(http.StatusOK, filesToMap)
}

func MapToDto(file *models.FilesToDto) FileResponse {
	return FileResponse{
		ID:          file.File.ID,
		ProjectID:  file.File.ProjectID,
		Title:      file.File.Title,
		Year:       file.File.Year,
		FileType:   file.FileType,
		CreatedAt:  utils.CustomTime{Time: file.CreatedAt},
		ProjectName: file.ProjectName,
	}
}
