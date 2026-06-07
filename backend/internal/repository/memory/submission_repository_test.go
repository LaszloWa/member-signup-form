package memory

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
)

func TestSave_DuplicateScopedKey(t *testing.T) {
	repo := NewSubmissionRepository()
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
	repo := NewSubmissionRepository()
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

func TestSave_AllowsSameIdentityInDifferentForm(t *testing.T) {
	repo := NewSubmissionRepository()
	ctx := context.Background()

	first := testSubmission("test@example.com", "+47 12345678", "1990-04-21", "form-1", "club-1")
	if err := repo.Save(ctx, first); err != nil {
		t.Fatalf("expected first save to succeed, got %v", err)
	}

	second := testSubmission("test@example.com", "+47 12345678", "1990-04-21", "form-2", "club-1")
	if err := repo.Save(ctx, second); err != nil {
		t.Fatalf("expected second save in different form to succeed, got %v", err)
	}
}

func TestSave_ConcurrentDuplicateAttempts(t *testing.T) {
	repo := NewSubmissionRepository()
	ctx := context.Background()
	const attempts = 12

	errorsOut := make(chan error, attempts)
	var wg sync.WaitGroup

	for i := 0; i < attempts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := repo.Save(ctx, testSubmission("dupe@example.com", "+47 11111111", "1995-09-01", "form-1", "club-1"))
			errorsOut <- err
		}(i)
	}

	wg.Wait()
	close(errorsOut)

	successCount := 0
	duplicateCount := 0
	for err := range errorsOut {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, repository.ErrDuplicateSubmission):
			duplicateCount++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly one successful save, got %d", successCount)
	}
	if duplicateCount != attempts-1 {
		t.Fatalf("expected %d duplicate errors, got %d", attempts-1, duplicateCount)
	}
}

func testSubmission(email, phoneNumber, birthDate, formID, clubID string) domain.Submission {
	return domain.Submission{
		SubmissionID: "sub-id",
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
