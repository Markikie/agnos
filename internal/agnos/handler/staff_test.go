package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Markikie/agnos/internal/agnos/api/request"
	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/google/uuid"
)

// MockStaffService is a mock implementation of StaffService
type MockStaffService struct {
	mock.Mock
}

func (m *MockStaffService) CreateStaff(username, password, hospital string) (*entity.Staff, error) {
	args := m.Called(username, password, hospital)
	return args.Get(0).(*entity.Staff), args.Error(1)
}

func (m *MockStaffService) Login(username, password, hospital string) (*entity.Staff, error) {
	args := m.Called(username, password, hospital)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Staff), args.Error(1)
}

func (m *MockStaffService) GetStaffByID(id string) (*entity.Staff, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Staff), args.Error(1)
}

func TestStaffHandler_CreateStaff_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := StaffHandler{
		staffService: mockService,
		jwtSecret:    "test-secret",
	}

	staff := &entity.Staff{
		ID:       uuid.New(),
		Username: "testuser",
		Hospital: "hospital-a",
	}

	mockService.On("CreateStaff", "testuser", "password123", "hospital-a").Return(staff, nil)

	reqBody := request.StaffRequest{
		Username: "testuser",
		Password: "password123",
		Hospital: "hospital-a",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateStaff(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestStaffHandler_CreateStaff_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := StaffHandler{
		staffService: mockService,
		jwtSecret:    "test-secret",
	}

	reqBody := request.StaffRequest{
		Username: "", // Missing username
		Password: "password123",
		Hospital: "hospital-a",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateStaff(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStaffHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := StaffHandler{
		staffService: mockService,
		jwtSecret:    "test-secret",
	}

	staff := &entity.Staff{
		ID:       uuid.New(),
		Username: "testuser",
		Hospital: "hospital-a",
	}

	mockService.On("Login", "testuser", "password123", "hospital-a").Return(staff, nil)

	reqBody := request.LoginStaffRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/staff/login?hospital=hospital-a", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "access_token")
	
	mockService.AssertExpectations(t)
}

func TestStaffHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := StaffHandler{
		staffService: mockService,
		jwtSecret:    "test-secret",
	}

	mockService.On("Login", "testuser", "wrongpassword", "hospital-a").Return(nil, assert.AnError)

	reqBody := request.LoginStaffRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/staff/login?hospital=hospital-a", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}
