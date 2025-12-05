package checker

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mikail-tommard/2025-12-10-links/internal/domain"
	"github.com/mikail-tommard/2025-12-10-links/internal/usecase"
)

type HTTPChecker struct {
	repo usecase.BatchRepository
	client *http.Client
	maxWorkers int
}

func NewHTTPChecker(repo usecase.BatchRepository, timeout time.Duration, maxWorkers int) *HTTPChecker {
	client := http.Client{
		Timeout: timeout,
	}

	if maxWorkers <= 0 {
		maxWorkers = 5
	}

	return &HTTPChecker{
		repo: repo,
		client: &client,
		maxWorkers: maxWorkers,
	}
}

func (c *HTTPChecker) CheckBatch(ctx context.Context, batch *domain.LinkBatch) error {
	sem := make(chan struct{}, c.maxWorkers)
	wg := sync.WaitGroup{}

	for _, link := range batch.Links {
		sem <- struct{}{}
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			defer func(){<-sem}()

			status := domain.StatusUnavailable
			errText := ""

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				errText := err.Error()
				batch.SetResults(url, status, errText)
				return
			}

			resp, err := c.client.Do(req)
			if err != nil {
				errText = err.Error()
				batch.SetResults(url, status, errText)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				status = domain.StatusAvailable
			} else {
				errText = fmt.Sprintf("unepexted status code: %d", resp.StatusCode)
			}
			
			batch.SetResults(url, status, errText)
		}(link.URL)
	}
	wg.Wait()

	if err := c.repo.SaveBatch(ctx, batch); err != nil {
		return err
	}
	return nil
}