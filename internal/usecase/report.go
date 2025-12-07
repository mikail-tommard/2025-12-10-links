package usecase

import (
	"context"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type ReportGenerator interface {
	GenerateReport(batches []*domain.LinkBatch) ([]byte, error)
}

type ReportService struct {
	repo     BatchRepository
	reporter ReportGenerator
}

func NewReportService(repo BatchRepository, reported ReportGenerator) *ReportService {
	return &ReportService{
		repo:     repo,
		reporter: reported,
	}
}

func (s *ReportService) GenerateReportForBatches(ctx context.Context, ids []domain.BatchID) ([]byte, error) {
	batches, err := s.repo.GetBatches(ctx, ids)
	if err != nil {
		return nil, err
	}

	bytes, err := s.reporter.GenerateReport(batches)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
