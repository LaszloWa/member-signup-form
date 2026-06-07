package service

import (
	"context"
	"errors"
	"testing"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
)

type stubSubmissionRepository struct {
	saved         []domain.Submission
	nextSaveError error
}

func (s *stubSubmissionRepository) Save(_ context.Context, submission domain.Submission) error {
	if s.nextSaveError != nil {
		return s.nextSaveError
	}

	s.saved = append(s.saved, submission)
	return nil
}

func (s *stubSubmissionRepository) List(_ context.Context) ([]domain.Submission, error) {
	result := make([]domain.Submission, len(s.saved))
	copy(result, s.saved)
	return result, nil
}

func TestCreate_MapsDuplicateRepositoryError(t *testing.T) {
	formService, err := NewFormService(DefaultFormDetails())
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	repo := &stubSubmissionRepository{nextSaveError: repository.ErrDuplicateSubmission}
	svc := NewSubmissionService(formService, repo)

	_, validationErrors, err := svc.Create(context.Background(), validInput())
	if !errors.Is(err, ErrDuplicateSubmission) {
		t.Fatalf("expected duplicate service error, got %v", err)
	}
	if len(validationErrors) != 0 {
		t.Fatalf("expected no validation errors, got %+v", validationErrors)
	}
}

func TestCreate_NormalizesBeforeSave(t *testing.T) {
	formService, err := NewFormService(DefaultFormDetails())
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	repo := &stubSubmissionRepository{}
	svc := NewSubmissionService(formService, repo)

	input := validInput()
	input.Email = "  TEST@Example.Com  "
	input.PhoneNumber = "  +47 12345678  "
	input.BirthDate = " 1990-04-21 "

	_, validationErrors, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no service error, got %v", err)
	}
	if len(validationErrors) != 0 {
		t.Fatalf("expected no validation errors, got %+v", validationErrors)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected one saved submission, got %d", len(repo.saved))
	}

	saved := repo.saved[0]
	if saved.Email != "test@example.com" {
		t.Fatalf("expected normalized email, got %q", saved.Email)
	}
	if saved.PhoneNumber != "+47 12345678" {
		t.Fatalf("expected trimmed phone number, got %q", saved.PhoneNumber)
	}
	if saved.BirthDate != "1990-04-21" {
		t.Fatalf("expected trimmed birth date, got %q", saved.BirthDate)
	}
}

func validInput() domain.SubmissionInput {
	defaults := DefaultFormDetails()
	return domain.SubmissionInput{
		Name:         "Test User",
		Email:        "test@example.com",
		PhoneNumber:  "+47 12345678",
		BirthDate:    "1990-04-21",
		MemberTypeID: defaults.MemberTypes[0].ID,
		ClubID:       defaults.ClubID,
		FormID:       defaults.FormID,
	}
}
