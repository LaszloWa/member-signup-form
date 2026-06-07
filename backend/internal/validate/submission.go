package validate

import (
	"net/mail"
	"regexp"
	"strings"
	"time"

	"spondtest/backend/internal/domain"
)

var phonePattern = regexp.MustCompile(`^[0-9+()\-\s]{6,20}$`)

type FormValidator interface {
	GetFormDetails() domain.FormDetails
	IsValidMemberType(memberTypeID string) bool
}

func NormalizeSubmissionInput(in domain.SubmissionInput) domain.SubmissionInput {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.PhoneNumber = strings.TrimSpace(in.PhoneNumber)
	in.BirthDate = strings.TrimSpace(in.BirthDate)
	in.MemberTypeID = strings.TrimSpace(in.MemberTypeID)
	in.ClubID = strings.TrimSpace(in.ClubID)
	in.FormID = strings.TrimSpace(in.FormID)
	return in
}

func SubmissionInput(in domain.SubmissionInput, formService FormValidator) map[string]string {
	errs := map[string]string{}

	if in.Name == "" {
		errs["name"] = "name is required"
	} else if len(in.Name) > 120 {
		errs["name"] = "name must be 120 characters or fewer"
	}

	if in.Email == "" {
		errs["email"] = "email is required"
	} else {
		parsedAddress, err := mail.ParseAddress(in.Email)
		if err != nil || parsedAddress.Address != in.Email {
			errs["email"] = "email is invalid"
		}
	}

	if in.PhoneNumber == "" {
		errs["phoneNumber"] = "phoneNumber is required"
	} else if !phonePattern.MatchString(in.PhoneNumber) {
		errs["phoneNumber"] = "phoneNumber format is invalid"
	}

	if in.BirthDate == "" {
		errs["birthDate"] = "birthDate is required"
	} else {
		birthDate, err := time.Parse("2006-01-02", in.BirthDate)
		if err != nil {
			errs["birthDate"] = "birthDate must be in YYYY-MM-DD format"
		} else {
			now := time.Now().UTC()
			if birthDate.After(now) {
				errs["birthDate"] = "birthDate cannot be in the future"
			}
			if birthDate.Before(now.AddDate(-120, 0, 0)) {
				errs["birthDate"] = "birthDate is unrealistically old"
			}
		}
	}

	form := formService.GetFormDetails()
	if in.MemberTypeID == "" {
		errs["memberTypeId"] = "memberTypeId is required"
	} else if !formService.IsValidMemberType(in.MemberTypeID) {
		errs["memberTypeId"] = "memberTypeId is invalid"
	}

	if in.ClubID == "" {
		errs["clubId"] = "clubId is required"
	} else if in.ClubID != form.ClubID {
		errs["clubId"] = "clubId does not match form"
	}

	if in.FormID == "" {
		errs["formId"] = "formId is required"
	} else if in.FormID != form.FormID {
		errs["formId"] = "formId does not match form"
	}

	registrationOpensAt, err := time.Parse(time.RFC3339, form.RegistrationOpens)
	if err == nil && registrationOpensAt.After(time.Now().UTC()) {
		errs["registrationOpens"] = "registration is not open yet"
	}

	return errs
}
