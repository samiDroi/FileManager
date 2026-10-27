package files

import (
	"fmt"
	"net/http"
	"path/filepath"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/paths"
	"server/internal/store/repository"
	"strconv"
	"time"

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

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewFileHandler(db)
	g := e.Group("/files")
	g.Use(middlewares.IsAuthenticated)
	g.GET(paths.File.Index.Route, h.GetAll)
	g.GET(paths.FileSpecials.GetByIds.Route, h.GetAllByProject)
	g.POST(paths.File.Upload.Route, h.Upload)
	g.DELETE(paths.File.Delete.Route, h.Delete)
	g.GET(paths.File.Search.Route, h.Search)
	g.GET(paths.File.Download.Route, h.Download)
	g.GET(paths.File.Preview.Route, h.GetPreviewFile)
	g.PUT(paths.File.Edit.Route, h.Edit)
}

func (h *FileHandler) GetAll(c echo.Context) error {
	files, err := h.repo.GetAllFiles(c.Request().Context())
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

	files, err := h.repo.GetAllFilesByProject(c.Request().Context(), project_id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, files)
}
func (h *FileHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	files, err := h.repo.SearchFiles(c.Request().Context(), param)
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

	file, err := h.repo.GetFileByID(c.Request().Context(), id)
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
	//validate with the dto
	var req dtos.FileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	//valid dto fields
	if req.ProjectID == 0 || req.Tittle == "" || req.Year == "" {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": "project_id, tittle and year are required"})
	}
	userSess := middlewares.GetRequestUserSession(c)

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer src.Close()
	//secure name to avoid conflicts and security issues, but keep the original extension
	ext := filepath.Ext(file.Filename)
	safeFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// newFile, err := h.repo.SaveFile(c.Request().Context(), req.ProjectID, req.Tittle, req.Year, safeFileName, src)
	newFile, err := h.repo.SaveFile(c.Request().Context(), req, safeFileName, src, userSess.UserID)
	if err != nil {
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
	err = h.repo.DeleteFile(c.Request().Context(), id)
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
	updated, err := h.repo.EditFile(c.Request().Context(), id, req)
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
	file, err := h.repo.GetFileByID(c.Request().Context(), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "file not found"})
	}
	return c.Inline(file.Path, file.Tittle+filepath.Ext(file.Path))
}
