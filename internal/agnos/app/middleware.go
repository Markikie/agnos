package app

import (
	middlewareConfig "github.com/Markikie/agnos/internal/agnos/middleware/config"
	"github.com/gin-gonic/gin"
)

func NewMiddleware(ginEngine *gin.Engine) {
	ginEngine.Use(middlewareConfig.Logger())
}
