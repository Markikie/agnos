package service

import (
	"errors"
	"time"

	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/Markikie/agnos/internal/agnos/repository"
	"golang.org/x/crypto/bcrypt"
)

type StaffService interface {
	CreateStaff(username, password, hospital string) (*entity.Staff, error)
	Login(username, password, hospital string) (*entity.Staff, error)
	GetStaffByID(id string) (*entity.Staff, error)
}

type staffService struct {
	staffRepository repository.StaffRepository
}

func NewStaffService(
	staffRepository repository.StaffRepository,
) StaffService {
	return &staffService{
		staffRepository: staffRepository,
	}
}

func (s *staffService) CreateStaff(username, password, hospital string) (*entity.Staff, error) {
	// Check if staff already exists
	existingStaff, _ := s.staffRepository.GetByUsernameAndHospital(username, hospital)
	if existingStaff != nil {
		return nil, errors.New("staff with this username already exists in this hospital")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	staff := &entity.Staff{
		Username:  username,
		Password:  string(hashedPassword),
		Hospital:  hospital,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.staffRepository.Create(staff)
	if err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Login(username, password, hospital string) (*entity.Staff, error) {
	staff, err := s.staffRepository.GetByUsernameAndHospital(username, hospital)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return staff, nil
}

func (s *staffService) GetStaffByID(id string) (*entity.Staff, error) {
	return s.staffRepository.GetByID(id)
}
