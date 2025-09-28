package router

import (
	"github.com/Markikie/agnos/internal/agnos/handler"
	"github.com/Markikie/agnos/internal/agnos/middleware"
	"github.com/gin-gonic/gin"
)

func NewPatientRouter(
	ginEngine *gin.Engine,
	handler handler.PatientHandler,
) {
	patientRouter := ginEngine.Group("/patient")
	
	// Apply authentication middleware
	patientRouter.Use(middleware.AuthMiddleware("your-secret-key"))
	
	patientRouter.POST("/search", handler.SearchPatients)
}
