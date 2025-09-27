package entity

import "github.com/google/uuid"

type Staff struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Username string    `gorm:"column:username;unique"`
	Password string    `gorm:"column:password"`
	Hospital string    `gorm:"column:hospital"`
}
