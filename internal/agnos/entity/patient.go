package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	ID           uuid.UUID `grom:"column:id;type:uuid;primaryKey"`
	FirstNameTH  string    `gorm:"column:first_name_th"`
	MiddleNameTH string    `gorm:"column:middle_name_th"`
	LastNameTH   string    `gorm:"column:last_name_th"`
	FirstNameEN  string    `gorm:"column:first_name_en"`
	MiddleNameEN string    `gorm:"column:middle_name_en"`
	LastNameEN   string    `gorm:"column:last_name_en"`
	DateOfBirth  time.Time `gorm:"column:date_of_birth"`
	PatientHN    string    `gorm:"column:patient_hn"`
	NationalID   string    `gorm:"column:national_id;unique"`
	PassportID   string    `gorm:"column:passport_id;unique"`
	PhoneNumber  string    `gorm:"column:phone_number"`
	Email        string    `gorm:"column:email"`
	Gender       string    `gorm:"column:gender"`
}

func (e *Patient) TableName() string {
	return "tbl_patients"
}

func (e *Patient) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return
}
