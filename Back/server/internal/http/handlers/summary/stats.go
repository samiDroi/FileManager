// estadistics to show in graphics
package summary

import (
	"server/internal/http/middlewares"
	"server/internal/services"
	"server/internal/store/repository/summary"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	service *services.StatsService
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	repo := summary.NewStatsRepository(db)
	svc := services.NewStatsService(repo)
	h := NewStatsHandler(svc)
	e.GET("/stats/file-types", h.GetFileTypesStats, middlewares.IsAuthenticated, middlewares.RequireAdmin(db))
}

func NewStatsHandler(service *services.StatsService) *StatsHandler {
	return &StatsHandler{
		service: service,
	}
}

func (h *StatsHandler) GetFileTypesStats(c echo.Context) error {
	haveStats, err := h.service.GetFileTypesStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return c.JSON(200, haveStats)
}
