package sqlite

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
)

func TestSave_DuplicateScopedKey(t *testing.T) {
	repo := newTestRepository(t)
	t.Cleanup(func() {
		_ = repo.Close()
	})

	ctx := context.Background()
	first := testSubmission("test@example.com", "+47 12345678", "1990-04-21", "form-1", "club-1")
	if err := repo.Save(ctx, first); err != nil {
		t.Fatalf("expected first save to succeed, got %v", err)
	}

	duplicate := testSubmission("TEST@example.com", "+47 12345678", "1990-04-21", "form-1", "club-1")
	err := repo.Save(ctx, duplicate)
	if !errors.Is(err, repository.ErrDuplicateSubmission) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestSave_DuplicateScopedKey_PhoneFormattingVariants(t *testing.T) {
	repo := newTestRepository(t)
	t.Cleanup(func() {
		_ = repo.Close()
	})

	ctx := context.Background()
	first := testSubmission("person@example.com", "+47 901 23 456", "1990-04-21", "form-1", "club-1")
	if err := repo.Save(ctx, first); err != nil {
		t.Fatalf("expected first save to succeed, got %v", err)
	}

	duplicate := testSubmission("person@example.com", "+47(901)-23-456", "1990-04-21", "form-1", "club-1")
	err := repo.Save(ctx, duplicate)
	if !errors.Is(err, repository.ErrDuplicateSubmission) {
		t.Fatalf("expected duplicate error for phone-format variant, got %v", err)
	}
}

func TestList_PersistsAcrossRepositoryReopen(t *testing.T) {
	dir := t.TempDir()
	databasePath := filepath.Join(dir, "backend.db")

	repo, err := NewSubmissionRepository(databasePath)
	if err != nil {
		t.Fatalf("create repository: %v", err)
	}

	submission := testSubmission("test@example.com", "+47 12345678", "1990-04-21", "form-1", "club-1")
	if err := repo.Save(context.Background(), submission); err != nil {
		t.Fatalf("save submission: %v", err)
	}

	if err := repo.Close(); err != nil {
		t.Fatalf("close repository: %v", err)
	}

	reopened, err := NewSubmissionRepository(databasePath)
	if err != nil {
		t.Fatalf("reopen repository: %v", err)
	}
	t.Cleanup(func() {
		_ = reopened.Close()
	})

	items, err := reopened.List(context.Background())
	if err != nil {
		t.Fatalf("list submissions: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 submission, got %d", len(items))
	}
	if items[0].Email != "test@example.com" {
		t.Fatalf("unexpected email in persisted submission: %q", items[0].Email)
	}
}

func newTestRepository(t *testing.T) *SubmissionRepository {
	t.Helper()
	databasePath := filepath.Join(t.TempDir(), "backend.db")
	repo, err := NewSubmissionRepository(databasePath)
	if err != nil {
		t.Fatalf("create test repository: %v", err)
	}
	return repo
}

func testSubmission(email, phoneNumber, birthDate, formID, clubID string) domain.Submission {
	return domain.Submission{
		SubmissionID: "sub-id-" + time.Now().UTC().Format("150405.000000"),
		Name:         "Tester",
		Email:        email,
		PhoneNumber:  phoneNumber,
		BirthDate:    birthDate,
		MemberTypeID: "member-type",
		ClubID:       clubID,
		FormID:       formID,
		CreatedAt:    time.Now().UTC(),
	}
}
