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
	e.POST("/files/upload", h.Upload, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "archivos", "crear"))
	e.DELETE("/files/:id", h.Delete, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "archivos", "eliminar"))
	e.GET("/files/search", h.Search, middlewares.IsAuthenticated)
	e.GET("/files/:id/download", h.Download, middlewares.IsAuthenticated)
	e.GET("/files/:id/preview", h.GetPreviewFile, middlewares.IsAuthenticated)
	e.PUT("/files/edit/:id", h.Edit, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "archivos", "editar"))
	e.GET("/files/:id/access", h.GetAccess, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
	e.POST("/files/:id/access/users", h.GrantUserAccess, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
	e.DELETE("/files/:id/access/users/:userId", h.RevokeUserAccess, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
	e.POST("/files/:id/access/groups", h.GrantGroupAccess, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
	e.DELETE("/files/:id/access/groups/:groupId", h.RevokeGroupAccess, middlewares.IsAuthenticated, middlewares.RequirePermission(db, "accesos", "otorgar"))
}

func (h *FileHandler) GetAll(c echo.Context) error {
	userSess := middlewares.GetRequestUserSession(c)
	files, err := h.service.GetAllFiles(c.Request().Context(), userSess.UserID)
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

	userSess := middlewares.GetRequestUserSession(c)
	files, err := h.service.GetAllFilesByProject(c.Request().Context(), project_id, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, files)
}

func (h *FileHandler) Search(c echo.Context) error {
	param := c.QueryParam("q")
	userSess := middlewares.GetRequestUserSession(c)
	files, err := h.service.SearchFiles(c.Request().Context(), param, userSess.UserID)
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

	userSess := middlewares.GetRequestUserSession(c)
	allowed, err := h.service.CanDownloadFile(c.Request().Context(), id, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": "No cuentas con permiso de descarga sobre este archivo"})
	}

	file, err := h.service.GetFileByID(c.Request().Context(), id, userSess.UserID)
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
	userSess := middlewares.GetRequestUserSession(c)
	err = h.service.DeleteFile(c.Request().Context(), id, userSess.UserID)
	if err != nil {
		if errors.Is(err, services.ErrFileAccessDenied) {
			return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
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
	userSess := middlewares.GetRequestUserSession(c)
	updated, err := h.service.EditFile(c.Request().Context(), id, req, userSess.UserID)
	if err != nil {
		if errors.Is(err, services.ErrFileAccessDenied) {
			return echo.NewHTTPError(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *FileHandler) GetAccess(c echo.Context) error {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de archivo inválido")
	}
	userSess := middlewares.GetRequestUserSession(c)
	access, err := h.service.GetFileAccess(c.Request().Context(), fileID, userSess.UserID)
	if err != nil {
		return fileAccessError(err)
	}
	return c.JSON(http.StatusOK, access)
}

func (h *FileHandler) GrantUserAccess(c echo.Context) error {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de archivo inválido")
	}

	var req dtos.GrantFileUserAccessRequest
	if err := c.Bind(&req); err != nil || req.UserID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id es obligatorio")
	}

	userSess := middlewares.GetRequestUserSession(c)
	if err := h.service.GrantUserAccess(c.Request().Context(), fileID, req.UserID, req.AccessType, userSess.UserID); err != nil {
		return fileAccessError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Acceso otorgado correctamente"})
}

func (h *FileHandler) RevokeUserAccess(c echo.Context) error {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de archivo inválido")
	}
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de usuario inválido")
	}

	userSess := middlewares.GetRequestUserSession(c)
	if err := h.service.RevokeUserAccess(c.Request().Context(), fileID, userID, userSess.UserID); err != nil {
		return fileAccessError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Acceso revocado correctamente"})
}

func (h *FileHandler) GrantGroupAccess(c echo.Context) error {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de archivo inválido")
	}

	var req dtos.GrantFileGroupAccessRequest
	if err := c.Bind(&req); err != nil || req.GroupID <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "group_id es obligatorio")
	}

	userSess := middlewares.GetRequestUserSession(c)
	if err := h.service.GrantGroupAccess(c.Request().Context(), fileID, req.GroupID, req.AccessType, userSess.UserID); err != nil {
		return fileAccessError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Acceso al grupo otorgado correctamente"})
}

func (h *FileHandler) RevokeGroupAccess(c echo.Context) error {
	fileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de archivo inválido")
	}
	groupID, err := strconv.Atoi(c.Param("groupId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID de grupo inválido")
	}

	userSess := middlewares.GetRequestUserSession(c)
	if err := h.service.RevokeGroupAccess(c.Request().Context(), fileID, groupID, userSess.UserID); err != nil {
		return fileAccessError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Acceso al grupo revocado correctamente"})
}

func fileAccessError(err error) error {
	switch {
	case errors.Is(err, services.ErrInvalidAccessType):
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrFileAccessDenied):
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func (h *FileHandler) GetPreviewFile(c echo.Context) error {
	id_param := c.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	userSess := middlewares.GetRequestUserSession(c)
	file, err := h.service.GetFileByID(c.Request().Context(), id, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "file not found"})
	}
	return c.Inline(file.Path, file.Tittle+filepath.Ext(file.Path))
}
