package app

import (
	"fmt"
	"log"

	"github.com/Markikie/agnos/internal/agnos"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB *gorm.DB
}

func NewConfig() *Config {
	return &Config{
		DB: ConnectDB(),
	}
}

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		agnos.Env.Database.Host,
		agnos.Env.Database.Port,
		agnos.Env.Database.User,
		agnos.Env.Database.Password,
		agnos.Env.Database.DBName,
	)

	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	return dbClient
}
