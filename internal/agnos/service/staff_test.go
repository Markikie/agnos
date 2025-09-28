package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/google/uuid"
)

// MockStaffRepository is a mock implementation of StaffRepository
type MockStaffRepository struct {
	mock.Mock
}

func (m *MockStaffRepository) Create(staff *entity.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *MockStaffRepository) GetByUsernameAndHospital(username, hospital string) (*entity.Staff, error) {
	args := m.Called(username, hospital)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Staff), args.Error(1)
}

func (m *MockStaffRepository) GetByID(id string) (*entity.Staff, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Staff), args.Error(1)
}

func TestStaffService_CreateStaff_Success(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	// Mock that staff doesn't exist
	mockRepo.On("GetByUsernameAndHospital", "testuser", "hospital-a").Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*entity.Staff")).Return(nil)

	staff, err := service.CreateStaff("testuser", "password123", "hospital-a")

	assert.NoError(t, err)
	assert.NotNil(t, staff)
	assert.Equal(t, "testuser", staff.Username)
	assert.Equal(t, "hospital-a", staff.Hospital)
	assert.NotEmpty(t, staff.Password)
	
	// Verify password is hashed
	err = bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte("password123"))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestStaffService_CreateStaff_UserExists(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	existingStaff := &entity.Staff{
		ID:       uuid.New(),
		Username: "testuser",
		Hospital: "hospital-a",
	}

	mockRepo.On("GetByUsernameAndHospital", "testuser", "hospital-a").Return(existingStaff, nil)

	staff, err := service.CreateStaff("testuser", "password123", "hospital-a")

	assert.Error(t, err)
	assert.Nil(t, staff)
	assert.Contains(t, err.Error(), "already exists")

	mockRepo.AssertExpectations(t)
}

func TestStaffService_Login_Success(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingStaff := &entity.Staff{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
		Hospital: "hospital-a",
	}

	mockRepo.On("GetByUsernameAndHospital", "testuser", "hospital-a").Return(existingStaff, nil)

	staff, err := service.Login("testuser", "password123", "hospital-a")

	assert.NoError(t, err)
	assert.NotNil(t, staff)
	assert.Equal(t, existingStaff.ID, staff.ID)
	assert.Equal(t, "testuser", staff.Username)

	mockRepo.AssertExpectations(t)
}

func TestStaffService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	mockRepo.On("GetByUsernameAndHospital", "testuser", "hospital-a").Return(nil, errors.New("not found"))

	staff, err := service.Login("testuser", "wrongpassword", "hospital-a")

	assert.Error(t, err)
	assert.Nil(t, staff)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestStaffService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingStaff := &entity.Staff{
		ID:       uuid.New(),
		Username: "testuser",
		Password: string(hashedPassword),
		Hospital: "hospital-a",
	}

	mockRepo.On("GetByUsernameAndHospital", "testuser", "hospital-a").Return(existingStaff, nil)

	staff, err := service.Login("testuser", "wrongpassword", "hospital-a")

	assert.Error(t, err)
	assert.Nil(t, staff)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestStaffService_GetStaffByID_Success(t *testing.T) {
	mockRepo := new(MockStaffRepository)
	service := NewStaffService(mockRepo)

	staffID := uuid.New().String()
	existingStaff := &entity.Staff{
		ID:       uuid.MustParse(staffID),
		Username: "testuser",
		Hospital: "hospital-a",
	}

	mockRepo.On("GetByID", staffID).Return(existingStaff, nil)

	staff, err := service.GetStaffByID(staffID)

	assert.NoError(t, err)
	assert.NotNil(t, staff)
	assert.Equal(t, staffID, staff.ID.String())

	mockRepo.AssertExpectations(t)
}
