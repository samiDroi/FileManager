package summary

import (
	"context"
	"server/internal/domain/dtos"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetFileTypesStats(ctx context.Context) ([]dtos.FileTypeCountStats, error) {
	q := `SELECT file_type, COUNT(*) FROM files GROUP BY file_type`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fileStats []dtos.FileTypeCountStats
	for rows.Next() {
		var stat dtos.FileTypeCountStats
		if err := rows.Scan(&stat.FileType, &stat.Count); err != nil {
			return nil, err
		}
		fileStats = append(fileStats, stat)
	}

	return fileStats, nil
}
