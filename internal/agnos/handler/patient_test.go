package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Markikie/agnos/internal/agnos/api/request"
	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/google/uuid"
)

// MockPatientService is a mock implementation of PatientService
type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) SearchPatients(filters map[string]interface{}, staffHospital string) ([]*entity.Patient, error) {
	args := m.Called(filters, staffHospital)
	return args.Get(0).([]*entity.Patient), args.Error(1)
}

func (m *MockPatientService) GetPatientFromHospitalAPI(id, hospital string) (*entity.Patient, error) {
	args := m.Called(id, hospital)
	return args.Get(0).(*entity.Patient), args.Error(1)
}

func TestPatientHandler_SearchPatients_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPatientService)
	handler := PatientHandler{
		patientService: mockService,
	}

	patients := []*entity.Patient{
		{
			ID:          uuid.New(),
			FirstNameTH: "สมชาย",
			LastNameTH:  "ใจดี",
			FirstNameEN: "Somchai",
			LastNameEN:  "Jaidee",
			NationalID:  "1234567890123",
			Gender:      "M",
		},
	}

	expectedFilters := map[string]interface{}{
		"national_id": "1234567890123",
	}

	mockService.On("SearchPatients", expectedFilters, "hospital-a").Return(patients, nil)

	reqBody := request.PatientSearchRequest{
		NationalID: "1234567890123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/patient/search", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("hospital", "hospital-a") // Simulate auth middleware

	handler.SearchPatients(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "patients")
	assert.Contains(t, response, "count")
	assert.Equal(t, float64(1), response["count"])
	
	mockService.AssertExpectations(t)
}

func TestPatientHandler_SearchPatients_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPatientService)
	handler := PatientHandler{
		patientService: mockService,
	}

	reqBody := request.PatientSearchRequest{
		NationalID: "1234567890123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/patient/search", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// No hospital set in context - simulating missing auth

	handler.SearchPatients(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHandler_SearchPatients_MultipleFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPatientService)
	handler := PatientHandler{
		patientService: mockService,
	}

	patients := []*entity.Patient{}

	expectedFilters := map[string]interface{}{
		"first_name":     "John",
		"last_name":      "Doe",
		"phone_number":   "0812345678",
		"date_of_birth":  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockService.On("SearchPatients", expectedFilters, "hospital-a").Return(patients, nil)

	reqBody := request.PatientSearchRequest{
		FirstName:    "John",
		LastName:     "Doe",
		PhoneNumber:  "0812345678",
		DateOfBirth:  "1990-01-01",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/patient/search", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("hospital", "hospital-a")

	handler.SearchPatients(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(0), response["count"])
	
	mockService.AssertExpectations(t)
}

func TestPatientHandler_SearchPatients_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPatientService)
	handler := PatientHandler{
		patientService: mockService,
	}

	req, _ := http.NewRequest("POST", "/patient/search", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("hospital", "hospital-a")

	handler.SearchPatients(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
