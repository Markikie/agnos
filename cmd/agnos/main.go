package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Markikie/agnos/internal/agnos/entity"
	"github.com/Markikie/agnos/internal/agnos/handler"
	"github.com/Markikie/agnos/internal/agnos/repository"
	"github.com/Markikie/agnos/internal/agnos/router"
	"github.com/Markikie/agnos/internal/agnos/service"
)

func main() {
	// Database connection
	dsn := "host=localhost user=agnos password=password dbname=agnos port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	err = db.AutoMigrate(&entity.Staff{}, &entity.Patient{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize repositories
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// Initialize services
	staffService := service.NewStaffService(staffRepo)
	patientService := service.NewPatientService(patientRepo)

	// Initialize handlers
	staffHandler := handler.NewStaffHandler(staffService)
	patientHandler := handler.NewPatientHandler(patientService)

	// Initialize Gin
	app := gin.Default()

	// Health check endpoint
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Agnos Hospital Middleware API",
			"status":  "healthy",
		})
	})

	// Setup routes
	router.NewStaffRouter(app, staffHandler)
	router.NewPatientRouter(app, patientHandler)

	// Start server
	log.Println("Starting server on :8081...")
	err = app.Run(":8081")

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
