package services

import (
	"context"
	"server/internal/domain/dtos"
	"server/internal/store/repository/summary"
)

type StatsService struct {
	repo *summary.StatsRepository
}

func NewStatsService(repo *summary.StatsRepository) *StatsService {
	return &StatsService{repo: repo}
}

func (s *StatsService) GetFileTypesStats(ctx context.Context) ([]dtos.FileTypeCountStats, error) {
	return s.repo.GetFileTypesStats(ctx)
}
