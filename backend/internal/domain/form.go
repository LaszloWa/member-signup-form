package domain

type MemberType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type FormDetails struct {
	ClubID            string       `json:"clubId"`
	MemberTypes       []MemberType `json:"memberTypes"`
	FormID            string       `json:"formId"`
	Title             string       `json:"title"`
	RegistrationOpens string       `json:"registrationOpens"`
}
