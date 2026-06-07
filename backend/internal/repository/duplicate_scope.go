package repository

import (
	"strings"

	"spondtest/backend/internal/domain"
)

type DuplicateScope struct {
	Email     string
	Phone     string
	BirthDate string
	FormID    string
	ClubID    string
}

func DuplicateScopeFromSubmission(submission domain.Submission) DuplicateScope {
	return DuplicateScope{
		Email:     strings.TrimSpace(strings.ToLower(submission.Email)),
		Phone:     normalizePhoneForDuplicate(submission.PhoneNumber),
		BirthDate: strings.TrimSpace(submission.BirthDate),
		FormID:    strings.TrimSpace(submission.FormID),
		ClubID:    strings.TrimSpace(submission.ClubID),
	}
}

func normalizePhoneForDuplicate(value string) string {
	value = strings.TrimSpace(value)
	replacer := strings.NewReplacer(
		" ", "",
		"-", "",
		"(", "",
		")", "",
	)
	return replacer.Replace(value)
}

func (s DuplicateScope) Key() string {
	return strings.Join([]string{
		s.Email,
		s.Phone,
		s.BirthDate,
		s.FormID,
		s.ClubID,
	}, "|")
}
