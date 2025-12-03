package domain

import (
	"errors"
	"time"
)

type BatchID int
type LinkStatus string
type BatchStatus string

const (
	StatusUnknown LinkStatus = "unknown"
	StatusAvailable LinkStatus = "available"
	StatusUnavailable LinkStatus = "unavailable"
)

const (
	BatchStatusCreated BatchStatus = "created"
	BatchStatusInProgress BatchStatus = "in_progress"
	BatchStatusDone BatchStatus = "done"
	BatchStatusFailed BatchStatus = "failed"
)

type Link struct {
	URL string
}

type LinkResult struct {
	Link Link
	Status LinkStatus
	Error string
}

type LinkBatch struct {
	ID BatchID
	Links []Link
	Results []LinkResult
	Status BatchStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewLinkBatch(id BatchID, urls []string) (*LinkBatch, error) {
	if len(urls) == 0 {
		return nil, errors.New("urls must not be empty")
	}
	now := time.Now()

	links := make([]Link, 0, len(urls))

	for _, url := range urls {
		if url == "" {
			continue
		}
		links = append(links, Link{
			URL: url,
		})
	}

	if len(links) == 0 {
		return nil, errors.New("no valid urls provided")
	}

	results := make([]LinkResult, 0, len(links))
	for _, link := range links {
		results = append(results, LinkResult{
			Link: link,
			Status: StatusUnknown,
			Error: "",
		})
	}

	return &LinkBatch{
		ID: id,
		Links: links,
		Results: results,
		Status: BatchStatusCreated,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (b *LinkBatch) allLinksProcessed() bool {
	for _, result := range b.Results {
		if result.Status == StatusUnknown {
			return false
		}
	}
	return true
}

func (b *LinkBatch) StartProcessing() error {
	if b.Status != BatchStatusCreated {
		return errors.New("batch must be in 'created' status")
	}
	b.Status = BatchStatusInProgress
	b.UpdatedAt = time.Now()
	return nil
}

func (b *LinkBatch) SetResults(url string, status LinkStatus, errText string) {
	for i, link := range b.Links {
		if link.URL == url {
			b.Results[i].Status = status
			b.Results[i].Error = errText
			b.UpdatedAt = time.Now()
			break
		}
	}

	if b.Status == BatchStatusInProgress && b.allLinksProcessed() {
		b.Status = BatchStatusDone
		b.UpdatedAt = time.Now()
	}
}