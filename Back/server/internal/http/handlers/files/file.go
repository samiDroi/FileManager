package files

import (
	"errors"
	"net/http"
	"path/filepath"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type FileHandler struct {
	service *services.FileService
}

func NewFileHandler(service *services.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewFileRepository(db)
	svc := services.NewFileService(repo)
	h := NewFileHandler(svc)

	e.GET("/files", h.GetAll, middlewares.IsAuthenticated)
	e.GET("/files/project/:id", h.GetAllByProject, middlewares.IsAuthenticated)
	e.POST("/files/upload", h.Upload, middlewares.IsAuthenticated)
	e.DELETE("/files/:id", h.Delete, middlewares.IsAuthenticated)
	e.GET("/files/search", h.Search, middlewares.IsAuthenticated)
	e.GET("/files/:id/download", h.Download, middlewares.IsAuthenticated)
	e.GET("/files/:id/preview", h.GetPreviewFile, middlewares.IsAuthenticated)
	e.PUT("/files/edit/:id", h.Edit, middlewares.IsAuthenticated)
}

func (h *FileHandler) GetAll(c echo.Context) error {
	files, err := h.service.GetAllFiles(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, files)
}

func (h *FileHandler) GetAllByProject(c echo.Context) error {
	id_param := c.Param("id")
	project_id, err := strconv.Atoi(id_param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	files, err := h.service.GetAllFilesByProject(c.Request().Context(), project_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, files)
}

func (h *FileHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	files, err := h.service.SearchFiles(c.Request().Context(), param)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, files)
}

func (h *FileHandler) Download(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "invalid file ID"})
	}

	file, err := h.service.GetFileByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "file not found"})
	}

	ext := filepath.Ext(file.Path)
	downloadName := file.Tittle + ext

	return c.Attachment(file.Path, downloadName)
}

func (h *FileHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "file is required"})
	}

	var req dtos.FileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.service.ValidateUpload(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userSess := middlewares.GetRequestUserSession(c)

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer src.Close()

	newFile, err := h.service.UploadFile(c.Request().Context(), req, file.Filename, src, userSess.UserID)
	if err != nil {
		if errors.Is(err, services.ErrInvalidFileData) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, newFile)
}

func (h *FileHandler) Delete(c echo.Context) error {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err = h.service.DeleteFile(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "File deleted"})
}

func (h *FileHandler) Edit(c echo.Context) error {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	var req dtos.FileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	updated, err := h.service.EditFile(c.Request().Context(), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *FileHandler) GetPreviewFile(c echo.Context) error {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	file, err := h.service.GetFileByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "file not found"})
	}
	return c.Inline(file.Path, file.Tittle+filepath.Ext(file.Path))
}
