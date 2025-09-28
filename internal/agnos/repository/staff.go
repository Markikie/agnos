package repository

import (
	"github.com/Markikie/agnos/internal/agnos/entity"
	"gorm.io/gorm"
)

type StaffRepository interface {
	Create(staff *entity.Staff) error
	GetByUsernameAndHospital(username, hospital string) (*entity.Staff, error)
	GetByID(id string) (*entity.Staff, error)
}

type staffRepository struct {
	db *gorm.DB
}
func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{
		db: db,
	}
}

func (r *staffRepository) Create(staff *entity.Staff) error {
	return r.db.Create(staff).Error
}

func (r *staffRepository) GetByUsernameAndHospital(username, hospital string) (*entity.Staff, error) {
	var staff entity.Staff
	err := r.db.Where("username = ? AND hospital = ?", username, hospital).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) GetByID(id string) (*entity.Staff, error) {
	var staff entity.Staff
	err := r.db.Where("id = ?", id).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}