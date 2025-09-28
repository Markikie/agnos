package app

import (
	"github.com/Markikie/agnos/internal/agnos/service"
)

type Service struct {
	PatientService service.PatientService
	StaffService   service.StaffService
}

func NewService(repository *Repository) *Service {
	return &Service{
		PatientService: service.NewPatientService(repository.PatientRepository),
		StaffService:   service.NewStaffService(repository.StaffRepository),
	}
}
