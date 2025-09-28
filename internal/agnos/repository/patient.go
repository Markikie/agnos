package repository

import (
	"github.com/Markikie/agnos/internal/agnos/entity"
	"gorm.io/gorm"
)

type PatientRepository interface {
	Create(patient *entity.Patient) error
	Search(filters map[string]interface{}) ([]*entity.Patient, error)
	GetByID(id string) (*entity.Patient, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{
		db: db,
	}
}

func (r *patientRepository) Create(patient *entity.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) Search(filters map[string]interface{}) ([]*entity.Patient, error) {
	var patients []*entity.Patient
	query := r.db.Model(&entity.Patient{})
	
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "national_id":
				query = query.Where("national_id = ?", value)
			case "passport_id":
				query = query.Where("passport_id = ?", value)
			case "first_name":
				query = query.Where("first_name_th ILIKE ? OR first_name_en ILIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
			case "middle_name":
				query = query.Where("middle_name_th ILIKE ? OR middle_name_en ILIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
			case "last_name":
				query = query.Where("last_name_th ILIKE ? OR last_name_en ILIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
			case "date_of_birth":
				query = query.Where("date_of_birth = ?", value)
			case "phone_number":
				query = query.Where("phone_number = ?", value)
			case "email":
				query = query.Where("email = ?", value)
			}
		}
	}
	
	err := query.Find(&patients).Error
	return patients, err
}

func (r *patientRepository) GetByID(id string) (*entity.Patient, error) {
	var patient entity.Patient
	err := r.db.Where("id = ?", id).First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}
