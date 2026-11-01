// estadistics to show in graphics
package summary

import (
	"server/internal/store/repository/summary"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type StatsHandler struct {
	db   *pgxpool.Pool
	repo *summary.StatsRepository
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool) {
	h := NewStatsHandler(db)
	e.GET("/stats/file-types", h.GetFileTypesStats)
}
func NewStatsHandler(db *pgxpool.Pool) *StatsHandler {
	return &StatsHandler{
		db:   db,
		repo: summary.NewStatsRepository(db),
	}
}

func (h *StatsHandler) GetFileTypesStats(c echo.Context) error {
	haveStats, err := h.repo.GetFileTypesStats(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	return c.JSON(200, haveStats)
}
