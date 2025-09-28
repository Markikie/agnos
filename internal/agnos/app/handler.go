package app

import "github.com/Markikie/agnos/internal/agnos/handler"

type Handler struct {
	StaffHandler   handler.StaffHandler
	PatientHandler handler.PatientHandler
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		StaffHandler:   handler.NewStaffHandler(service.StaffService),
		PatientHandler: handler.NewPatientHandler(service.PatientService),
	}
}
