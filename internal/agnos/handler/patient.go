package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/Markikie/agnos/internal/agnos/api/request"
	"github.com/Markikie/agnos/internal/agnos/service"
)

type PatientHandler struct {
	patientService service.PatientService
}

func NewPatientHandler(
	patientService service.PatientService,
) PatientHandler {
	return PatientHandler{
		patientService: patientService,
	}
}

func (h *PatientHandler) SearchPatients(c *gin.Context) {
	var req request.PatientSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get staff info from context (set by auth middleware)
	staffHospital, exists := c.Get("hospital")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Staff hospital not found"})
		return
	}

	// Build search filters
	filters := make(map[string]interface{})
	if req.NationalID != "" {
		filters["national_id"] = req.NationalID
	}
	if req.PassportID != "" {
		filters["passport_id"] = req.PassportID
	}
	if req.FirstName != "" {
		filters["first_name"] = req.FirstName
	}
	if req.MiddleName != "" {
		filters["middle_name"] = req.MiddleName
	}
	if req.LastName != "" {
		filters["last_name"] = req.LastName
	}
	if req.DateOfBirth != "" {
		// Parse date
		if dob, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			filters["date_of_birth"] = dob
		}
	}
	if req.PhoneNumber != "" {
		filters["phone_number"] = req.PhoneNumber
	}
	if req.Email != "" {
		filters["email"] = req.Email
	}

	patients, err := h.patientService.SearchPatients(filters, staffHospital.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"patients": patients,
		"count":    len(patients),
	})
}
