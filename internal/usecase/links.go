package usecase

import (
	"context"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type BatchRepository interface {
	// Вернуть новый id для батча
	NextBatchD(ctx context.Context) (domain.BatchID, error)
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
	Repo BatchRepository
	Checker Checker
}

func NewLinkService(repo BatchRepository, checker Checker) *LinksService {
	return &LinksService{
		Repo: repo,
		Checker: checker,
	}
}