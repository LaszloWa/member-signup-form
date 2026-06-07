package service

import (
	"errors"
	"strings"
	"time"

	"spondtest/backend/internal/domain"
)

type FormService struct {
	form domain.FormDetails
}

func DefaultFormDetails() domain.FormDetails {
	return domain.FormDetails{
		ClubID: "britsport",
		MemberTypes: []domain.MemberType{
			{ID: "8FE4113D4E4020E0DCF887803A886981", Name: "Active Member"},
			{ID: "4237C55C5CC3B4B082CBF2540612778E", Name: "Social Member"},
		},
		FormID:            "B171388180BC457D9887AD92B6CCFC86",
		Title:             "Coding camp summer 2025",
		RegistrationOpens: "2024-12-16T00:00:00Z",
	}
}

func NewFormService(form domain.FormDetails) (*FormService, error) {
	if err := validateFormDetails(form); err != nil {
		return nil, err
	}

	return &FormService{form: form}, nil
}

func validateFormDetails(form domain.FormDetails) error {
	if strings.TrimSpace(form.ClubID) == "" {
		return errors.New("clubId must not be empty")
	}

	if strings.TrimSpace(form.FormID) == "" {
		return errors.New("formId must not be empty")
	}

	if strings.TrimSpace(form.Title) == "" {
		return errors.New("title must not be empty")
	}

	if _, err := time.Parse(time.RFC3339, form.RegistrationOpens); err != nil {
		return errors.New("registrationOpens must be a valid RFC3339 timestamp")
	}

	if len(form.MemberTypes) == 0 {
		return errors.New("memberTypes must contain at least one option")
	}

	for _, memberType := range form.MemberTypes {
		if strings.TrimSpace(memberType.ID) == "" || strings.TrimSpace(memberType.Name) == "" {
			return errors.New("memberTypes entries must include non-empty id and name")
		}
	}

	return nil
}

func (s *FormService) GetFormDetails() domain.FormDetails {
	return s.form
}

func (s *FormService) IsValidMemberType(memberTypeID string) bool {
	for _, mt := range s.form.MemberTypes {
		if mt.ID == memberTypeID {
			return true
		}
	}
	return false
}
