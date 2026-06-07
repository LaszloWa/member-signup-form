package domain

import "time"

type SubmissionInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phoneNumber"`
	BirthDate    string `json:"birthDate"`
	MemberTypeID string `json:"memberTypeId"`
	ClubID       string `json:"clubId"`
	FormID       string `json:"formId"`
}

type Submission struct {
	SubmissionID string    `json:"submissionId"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phoneNumber"`
	BirthDate    string    `json:"birthDate"`
	MemberTypeID string    `json:"memberTypeId"`
	ClubID       string    `json:"clubId"`
	FormID       string    `json:"formId"`
	CreatedAt    time.Time `json:"createdAt"`
}
