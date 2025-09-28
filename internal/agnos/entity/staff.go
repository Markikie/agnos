package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Staff struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Username  string    `gorm:"column:username;unique;not null"`
	Password  string    `gorm:"column:password;not null"`
	Hospital  string    `gorm:"column:hospital;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (s *Staff) TableName() string {
	return "tbl_staff"
}

func (s *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
