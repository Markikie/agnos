package request

type PatientSearchRequest struct {
	NationalID   string `json:"national_id,omitempty"`
	PassportID   string `json:"passport_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	MiddleName   string `json:"middle_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	DateOfBirth  string `json:"date_of_birth,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Email        string `json:"email,omitempty"`
}
