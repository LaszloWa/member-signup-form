package main

import (
	"testing"

	"spondtest/backend/internal/repository/memory"
	sqliterepo "spondtest/backend/internal/repository/sqlite"
	"spondtest/backend/internal/service"
)

func TestDefaultFormConfigurationIsValid(t *testing.T) {
	_, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		t.Fatalf("expected defaults to be valid, got %v", err)
	}
}

func TestBuildSubmissionRepositoryFromEnv_Memory(t *testing.T) {
	t.Setenv("BACKEND_REPOSITORY", "memory")

	repo, err := buildSubmissionRepositoryFromEnv()
	if err != nil {
		t.Fatalf("expected memory repository, got error: %v", err)
	}

	if _, ok := repo.(*memory.SubmissionRepository); !ok {
		t.Fatalf("expected memory repository type, got %T", repo)
	}
}

func TestBuildSubmissionRepositoryFromEnv_SQLite(t *testing.T) {
	t.Setenv("BACKEND_REPOSITORY", "sqlite")
	t.Setenv("BACKEND_SQLITE_PATH", t.TempDir()+"/backend.db")

	repo, err := buildSubmissionRepositoryFromEnv()
	if err != nil {
		t.Fatalf("expected sqlite repository, got error: %v", err)
	}

	if _, ok := repo.(*sqliterepo.SubmissionRepository); !ok {
		t.Fatalf("expected sqlite repository type, got %T", repo)
	}
}
