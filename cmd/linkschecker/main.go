package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikail-tommard/2025-12-10-links/internal/adapter/checker"
	"github.com/mikail-tommard/2025-12-10-links/internal/adapter/httpapi"
	"github.com/mikail-tommard/2025-12-10-links/internal/adapter/report"
	"github.com/mikail-tommard/2025-12-10-links/internal/adapter/storage"
	"github.com/mikail-tommard/2025-12-10-links/internal/usecase"
)

const (
	addr      = ":8080"
	statePath = "data/state.json"
)

func main() {
	repo, err := storage.NewFileBatchRepository(statePath)
	if err != nil {
		log.Fatal("failed to create repository: %w", err)
	}
	linkChecker := checker.NewHTTPChecker(repo, 5*time.Second, 5)

	reporter := report.NewReportGenerator("Links check report", time.RFC3339)

	linksService := usecase.NewLinksService(repo, linkChecker)
	reportService := usecase.NewReportService(repo, reporter)

	api := httpapi.NewServer(linksService, reportService)

	srv := &http.Server{
		Addr:    addr,
		Handler: api.Handler(),
	}

	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %w", err)
	}
}
