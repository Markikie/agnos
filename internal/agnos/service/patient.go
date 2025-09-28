package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/Markikie/agnos/internal/agnos/repository"
)

type PatientService interface {
	SearchPatients(filters map[string]interface{}, staffHospital string) ([]*entity.Patient, error)
	GetPatientFromHospitalAPI(id, hospital string) (*entity.Patient, error)
}

type patientService struct {
	patientRepository repository.PatientRepository
}

type HospitalAPatientResponse struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

func NewPatientService(
	patientRepository repository.PatientRepository,
) PatientService {
	return &patientService{
		patientRepository: patientRepository,
	}
}

func (s *patientService) SearchPatients(filters map[string]interface{}, staffHospital string) ([]*entity.Patient, error) {
	// First search in local database
	patients, err := s.patientRepository.Search(filters)
	if err != nil {
		return nil, err
	}

	// If searching by national_id or passport_id and no local results, try Hospital API
	if len(patients) == 0 {
		if nationalID, ok := filters["national_id"].(string); ok && nationalID != "" {
			if apiPatient, err := s.GetPatientFromHospitalAPI(nationalID, staffHospital); err == nil {
				// Save to local database for future searches
				s.patientRepository.Create(apiPatient)
				patients = append(patients, apiPatient)
			}
		} else if passportID, ok := filters["passport_id"].(string); ok && passportID != "" {
			if apiPatient, err := s.GetPatientFromHospitalAPI(passportID, staffHospital); err == nil {
				// Save to local database for future searches
				s.patientRepository.Create(apiPatient)
				patients = append(patients, apiPatient)
			}
		}
	}

	return patients, nil
}

func (s *patientService) GetPatientFromHospitalAPI(id, hospital string) (*entity.Patient, error) {
	var apiURL string
	
	switch hospital {
	case "hospital-a":
		apiURL = fmt.Sprintf("https://hospital-a.api.co.th/patient/search/%s", id)
	default:
		return nil, fmt.Errorf("unsupported hospital: %s", hospital)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hospital API returned status: %d", resp.StatusCode)
	}

	var hospitalResp HospitalAPatientResponse
	if err := json.NewDecoder(resp.Body).Decode(&hospitalResp); err != nil {
		return nil, err
	}

	// Parse date of birth
	dob, err := time.Parse("2006-01-02", hospitalResp.DateOfBirth)
	if err != nil {
		return nil, err
	}

	patient := &entity.Patient{
		FirstNameTH:  hospitalResp.FirstNameTH,
		MiddleNameTH: hospitalResp.MiddleNameTH,
		LastNameTH:   hospitalResp.LastNameTH,
		FirstNameEN:  hospitalResp.FirstNameEN,
		MiddleNameEN: hospitalResp.MiddleNameEN,
		LastNameEN:   hospitalResp.LastNameEN,
		DateOfBirth:  dob,
		PatientHN:    hospitalResp.PatientHN,
		NationalID:   hospitalResp.NationalID,
		PassportID:   hospitalResp.PassportID,
		PhoneNumber:  hospitalResp.PhoneNumber,
		Email:        hospitalResp.Email,
		Gender:       hospitalResp.Gender,
	}

	return patient, nil
}
