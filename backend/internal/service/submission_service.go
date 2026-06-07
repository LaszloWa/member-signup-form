package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
	"spondtest/backend/internal/validate"
)

var ErrDuplicateSubmission = errors.New("duplicate submission")

type SubmissionService struct {
	formService *FormService
	repo        repository.SubmissionRepository
}

func NewSubmissionService(formService *FormService, repo repository.SubmissionRepository) *SubmissionService {
	return &SubmissionService{formService: formService, repo: repo}
}

func (s *SubmissionService) Create(ctx context.Context, input domain.SubmissionInput) (domain.Submission, map[string]string, error) {
	normalized := validate.NormalizeSubmissionInput(input)
	validationErrors := validate.SubmissionInput(normalized, s.formService)
	if len(validationErrors) > 0 {
		return domain.Submission{}, validationErrors, nil
	}

	now := time.Now().UTC()
	submission := domain.Submission{
		SubmissionID: randomHex(16),
		Name:         normalized.Name,
		Email:        normalized.Email,
		PhoneNumber:  normalized.PhoneNumber,
		BirthDate:    normalized.BirthDate,
		MemberTypeID: normalized.MemberTypeID,
		ClubID:       normalized.ClubID,
		FormID:       normalized.FormID,
		CreatedAt:    now,
	}

	if err := s.repo.Save(ctx, submission); err != nil {
		if errors.Is(err, repository.ErrDuplicateSubmission) {
			return domain.Submission{}, nil, ErrDuplicateSubmission
		}

		return domain.Submission{}, nil, err
	}

	return submission, nil, nil
}

func randomHex(bytesLen int) string {
	buffer := make([]byte, bytesLen)
	if _, err := rand.Read(buffer); err != nil {
		fallback := time.Now().UTC().Format("20060102150405.000000000")
		return fallback
	}
	return hex.EncodeToString(buffer)
}
