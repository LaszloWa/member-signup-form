package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	httpx "spondtest/backend/internal/http"
	"spondtest/backend/internal/http/handlers"
	"spondtest/backend/internal/repository"
	"spondtest/backend/internal/repository/memory"
	sqliterepo "spondtest/backend/internal/repository/sqlite"
	"spondtest/backend/internal/service"
)

func main() {
	formService, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		log.Fatalf("invalid form configuration: %v", err)
	}

	submissionRepo, err := buildSubmissionRepositoryFromEnv()
	if err != nil {
		log.Fatalf("failed to initialize submission repository: %v", err)
	}

	if closer, ok := submissionRepo.(interface{ Close() error }); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("repository close failed: %v", err)
			}
		}()
	}

	submissionService := service.NewSubmissionService(formService, submissionRepo)

	h := handlers.NewHandler(formService, submissionService)
	router := httpx.NewRouter(h, httpx.RouterOptions{AllowedOrigins: httpx.ParseAllowedOrigins(os.Getenv("BACKEND_ALLOWED_ORIGINS"))})

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		log.Printf("backend listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

func buildSubmissionRepositoryFromEnv() (repository.SubmissionRepository, error) {
	backendRepository := strings.ToLower(strings.TrimSpace(os.Getenv("BACKEND_REPOSITORY")))
	if backendRepository == "" {
		backendRepository = "sqlite"
	}

	switch backendRepository {
	case "memory":
		return memory.NewSubmissionRepository(), nil
	case "sqlite":
		sqlitePath := strings.TrimSpace(os.Getenv("BACKEND_SQLITE_PATH"))
		if sqlitePath == "" {
			sqlitePath = "./data/backend.db"
		}

		dir := filepath.Dir(sqlitePath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create sqlite directory: %w", err)
		}

		repo, err := sqliterepo.NewSubmissionRepository(sqlitePath)
		if err != nil {
			return nil, err
		}

		return repo, nil
	default:
		return nil, fmt.Errorf("unsupported BACKEND_REPOSITORY %q (supported: memory, sqlite)", backendRepository)
	}
}
