package validate_test

import (
	"testing"
	"time"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/service"
	"spondtest/backend/internal/validate"
)

func TestSubmissionInput_Valid(t *testing.T) {
	formService, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	in := domain.SubmissionInput{
		Name:         "Test User",
		Email:        "test@example.com",
		PhoneNumber:  "+47 12345678",
		BirthDate:    "1990-04-21",
		MemberTypeID: "8FE4113D4E4020E0DCF887803A886981",
		ClubID:       "britsport",
		FormID:       "B171388180BC457D9887AD92B6CCFC86",
	}

	errs := validate.SubmissionInput(in, formService)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %+v", errs)
	}
}

func TestSubmissionInput_Invalid(t *testing.T) {
	formService, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	in := domain.SubmissionInput{
		Name:         "",
		Email:        "not-an-email",
		PhoneNumber:  "bad phone #",
		BirthDate:    "3020-01-01",
		MemberTypeID: "unknown",
		ClubID:       "wrong",
		FormID:       "wrong",
	}

	errs := validate.SubmissionInput(in, formService)
	if len(errs) < 7 {
		t.Fatalf("expected many validation errors, got %+v", errs)
	}
}

func TestSubmissionInput_EmailValidationRules(t *testing.T) {
	formService, err := service.NewFormService(service.DefaultFormDetails())
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	base := domain.SubmissionInput{
		Name:         "Test User",
		PhoneNumber:  "+47 12345678",
		BirthDate:    "1990-04-21",
		MemberTypeID: "8FE4113D4E4020E0DCF887803A886981",
		ClubID:       "britsport",
		FormID:       "B171388180BC457D9887AD92B6CCFC86",
	}

	tests := []struct {
		name            string
		email           string
		expectEmailErr  bool
		expectNormalized string
	}{
		{
			name:            "accepts plain email",
			email:           "person@example.com",
			expectEmailErr:  false,
			expectNormalized: "person@example.com",
		},
		{
			name:            "normalizes uppercase email before validation",
			email:           "  PERSON@EXAMPLE.COM  ",
			expectEmailErr:  false,
			expectNormalized: "person@example.com",
		},
		{
			name:            "rejects display name email format",
			email:           "Person <person@example.com>",
			expectEmailErr:  true,
			expectNormalized: "person <person@example.com>",
		},
		{
			name:            "rejects multiple addresses",
			email:           "person@example.com, other@example.com",
			expectEmailErr:  true,
			expectNormalized: "person@example.com, other@example.com",
		},
		{
			name:            "rejects invalid email",
			email:           "not-an-email",
			expectEmailErr:  true,
			expectNormalized: "not-an-email",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := base
			input.Email = tc.email
			normalizedInput := validate.NormalizeSubmissionInput(input)

			if normalizedInput.Email != tc.expectNormalized {
				t.Fatalf("expected normalized email %q, got %q", tc.expectNormalized, normalizedInput.Email)
			}

			errs := validate.SubmissionInput(normalizedInput, formService)
			_, hasEmailError := errs["email"]
			if hasEmailError != tc.expectEmailErr {
				t.Fatalf("expected email error %v, got errors %+v", tc.expectEmailErr, errs)
			}
		})
	}
}

func TestSubmissionInput_RegistrationNotOpenYet(t *testing.T) {
	futureForm := service.DefaultFormDetails()
	futureForm.RegistrationOpens = time.Now().UTC().Add(24 * time.Hour).Format(time.RFC3339)

	formService, err := service.NewFormService(futureForm)
	if err != nil {
		t.Fatalf("failed to create form service: %v", err)
	}

	in := domain.SubmissionInput{
		Name:         "Test User",
		Email:        "test@example.com",
		PhoneNumber:  "+47 12345678",
		BirthDate:    "1990-04-21",
		MemberTypeID: futureForm.MemberTypes[0].ID,
		ClubID:       futureForm.ClubID,
		FormID:       futureForm.FormID,
	}

	errs := validate.SubmissionInput(in, formService)
	if errs["registrationOpens"] != "registration is not open yet" {
		t.Fatalf("expected registrationOpens error, got %+v", errs)
	}
}
