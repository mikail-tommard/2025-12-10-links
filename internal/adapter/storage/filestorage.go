package storage

import (
	"sync"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type FileBatchRepository struct {
	mu sync.RWMutex
	batches map[domain.BatchID]*domain.LinkBatch
	nextID domain.BatchID
	filepath string
}

type batchState struct {
	NextID domain.BatchID `json:"next_id"`
	Batches []*domain.LinkBatch `json:"batches"`
}

func NewFileBatchRepository(filepath string) *FileBatchRepository {
	return &FileBatchRepository{
		mu: sync.RWMutex{},
		batches: map[domain.BatchID]*domain.LinkBatch{},
		nextID: 0,
		filepath: filepath,
	}
}
