package memory

import (
	"context"
	"sync"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
)

type SubmissionRepository struct {
	mu          sync.RWMutex
	submissions []domain.Submission
	index       map[string]struct{}
}

func NewSubmissionRepository() *SubmissionRepository {
	return &SubmissionRepository{
		submissions: make([]domain.Submission, 0),
		index:       make(map[string]struct{}),
	}
}

func (r *SubmissionRepository) Save(_ context.Context, submission domain.Submission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := repository.DuplicateScopeFromSubmission(submission).Key()
	if _, exists := r.index[key]; exists {
		return repository.ErrDuplicateSubmission
	}

	r.submissions = append(r.submissions, submission)
	r.index[key] = struct{}{}
	return nil
}

func (r *SubmissionRepository) List(_ context.Context) ([]domain.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Submission, len(r.submissions))
	copy(result, r.submissions)
	return result, nil
}
