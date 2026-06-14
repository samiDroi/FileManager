package handlers

import (
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/domain/models"
	"server/internal/http/middlewares"
	"server/internal/paths"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type TypeFileHandler struct {
	db   *pgxpool.Pool
	repo *repository.TypeFileRepository
}

func NewTypeFileHandler(db *pgxpool.Pool) *TypeFileHandler {
	return &TypeFileHandler{
		db:   db,
		repo: repository.NewTypeFileRepository(db),
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewTypeFileHandler(db)
	g := e.Group(paths.TypeFile.Prefix)
	g.Use(middlewares.IsAuthenticated)
	g.GET(paths.TypeFile.Index.Route, h.GetAllTypeFiles)
	g.POST(paths.TypeFile.New.Route, h.CreateTypeFile)
}

func (h *TypeFileHandler) GetAllTypeFiles(ctx echo.Context) error {
	rows, err := h.repo.GetAllTypeFiles(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if rows == nil {
		rows = make([]*dtos.FileTypeResponse, 0)
	}
	return ctx.JSON(http.StatusOK, rows)
}

func (h *TypeFileHandler) CreateTypeFile(ctx echo.Context) error {
	userSess := middlewares.GetRequestUserSession(ctx)

	var req dtos.BulkCreateFileTypeRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var fileTypes []models.FileType
	for _, t := range req.Types {
		fileTypes = append(fileTypes, models.FileType{Name: t.Name})
	}

	typesFiles, err := h.repo.BulkCreateTypeFiles(ctx.Request().Context(), fileTypes, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, typesFiles)
}
