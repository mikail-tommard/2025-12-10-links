package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
)

type FileBatchRepository struct {
	mu       sync.RWMutex
	batches  map[domain.BatchID]*domain.LinkBatch
	nextID   domain.BatchID
	filepath string
}

type batchState struct {
	NextID  domain.BatchID      `json:"next_id"`
	Batches []*domain.LinkBatch `json:"batches"`
}

func NewFileBatchRepository(filepath string) (*FileBatchRepository, error) {
	r := &FileBatchRepository{
		batches:  map[domain.BatchID]*domain.LinkBatch{},
		nextID:   0,
		filepath: filepath,
	}

	if err := r.loadState(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *FileBatchRepository) loadState() error {
	file, err := os.Open(r.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	state := batchState{}
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.batches = map[domain.BatchID]*domain.LinkBatch{}
	for _, v := range state.Batches {
		r.batches[v.ID] = v
	}

	r.nextID = state.NextID
	if r.nextID == 0 {
		r.nextID = 1
	}
	return nil
}

func (r *FileBatchRepository) saveState() error {
	batches := make([]*domain.LinkBatch, 0, len(r.batches))
	for _, v := range r.batches {
		batches = append(batches, v)
	}

	state := batchState{
		NextID:  r.nextID,
		Batches: batches,
	}

	tmpFile, err := os.CreateTemp(r.filepath, "batch-state-*.tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if err := json.NewEncoder(tmpFile).Encode(&state); err != nil {
		return err
	}

	if err := os.Rename(tmpFile.Name(), r.filepath); err != nil {
		return err
	}
	return nil
}

func (r *FileBatchRepository) NextBatchID(ctx context.Context) (domain.BatchID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++
	return r.nextID, nil
}

func (r *FileBatchRepository) SaveBatch(ctx context.Context, batch *domain.LinkBatch) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.batches[batch.ID] = batch

	if batch.ID >= r.nextID {
		r.nextID = batch.ID + 1
	}

	return r.saveState()
}

func (r *FileBatchRepository) GetBatch(ctx context.Context, id domain.BatchID) (*domain.LinkBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	batch, ok := r.batches[id]
	if !ok {
		return nil, fmt.Errorf("batch %d not found", id)
	}

	return batch, nil
}

func (r *FileBatchRepository) GetBatches(ctx context.Context, ids []domain.BatchID) ([]*domain.LinkBatch, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.LinkBatch, 0, len(ids))

	for _, id := range ids {
		batch, ok := r.batches[id]
		if !ok {
			return nil, fmt.Errorf("batch %d not found", id)
		}
		result = append(result, batch)
	}
	return result, nil
}
