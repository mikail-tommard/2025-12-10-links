package usecase

import (
	"context"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type BatchRepository interface {
	// Вернуть новый id для батча
	NextBatchID(ctx context.Context) (domain.BatchID, error)
	// Сохранить или обновить батч
	SaveBatch(ctx context.Context, batch *domain.LinkBatch) error
	// Получить батч по id
	GetBatch(ctx context.Context, id domain.BatchID) (*domain.LinkBatch, error)
	// Получить несколько батчей по id
	GetBatches(ctx context.Context, id []domain.BatchID) ([]*domain.LinkBatch, error)
}

type Checker interface {
	// Проверка ссылок в батче и обвновляет их результаты
	CheckBatch(ctx context.Context, batch *domain.LinkBatch) error
}

type LinksService struct {
	repo    BatchRepository
	checker Checker
}

func NewLinksService(repo BatchRepository, checker Checker) *LinksService {
	return &LinksService{
		repo:    repo,
		checker: checker,
	}
}

func (s *LinksService) CreateAndCheckBatch(ctx context.Context, urls []string) (*domain.LinkBatch, error) {
	batchID, err := s.repo.NextBatchID(ctx)
	if err != nil {
		return nil, err
	}

	batch, err := domain.NewLinkBatch(batchID, urls)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveBatch(ctx, batch); err != nil {
		return nil, err
	}

	if err := batch.StartProcessing(); err != nil {
		return nil, err
	}

	if err := s.repo.SaveBatch(ctx, batch); err != nil {
		return nil, err
	}

	if err := s.checker.CheckBatch(ctx, batch); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetBatch(ctx, batch.ID)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
