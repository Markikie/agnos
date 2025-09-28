package app

import "github.com/Markikie/agnos/internal/agnos/repository"

type Repository struct {
	PatientRepository repository.PatientRepository
	StaffRepository   repository.StaffRepository
}

func NewRepository(config *Config) *Repository {
	return &Repository{
		PatientRepository: repository.NewPatientRepository(config.DB),
		StaffRepository:   repository.NewStaffRepository(config.DB),
	}
}
