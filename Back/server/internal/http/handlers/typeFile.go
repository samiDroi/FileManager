package handlers

import (
	"net/http"
	"server/internal/domain/dtos"
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type TypeFileHandler struct {
	service *services.TypeFileService
}

func NewTypeFileHandler(service *services.TypeFileService) *TypeFileHandler {
	return &TypeFileHandler{
		service: service,
	}
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := repository.NewTypeFileRepository(db)
	svc := services.NewTypeFileService(repo)
	h := NewTypeFileHandler(svc)

	e.GET("/type-files", h.GetAllTypeFiles, middlewares.IsAuthenticated)
	e.POST("/type-files/new", h.CreateTypeFile, middlewares.IsAuthenticated, middlewares.RequireAdmin(db), middlewares.RequirePermission(db, "tipos_archivo", "crear"))
}

func (h *TypeFileHandler) GetAllTypeFiles(ctx echo.Context) error {
	rows, err := h.service.GetAllTypeFiles(ctx.Request().Context())
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

	typesFiles, err := h.service.BulkCreateTypeFiles(ctx.Request().Context(), req, userSess.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, typesFiles)
}
