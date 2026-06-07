package repository

import (
	"context"
	"errors"

	"spondtest/backend/internal/domain"
)

var ErrDuplicateSubmission = errors.New("duplicate submission")

type SubmissionRepository interface {
	Save(ctx context.Context, submission domain.Submission) error
	List(ctx context.Context) ([]domain.Submission, error)
}
