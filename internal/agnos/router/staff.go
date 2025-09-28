package router

import (
	"github.com/Markikie/agnos/internal/agnos/handler"
	"github.com/gin-gonic/gin"
)

func NewStaffRouter(
	ginEngine *gin.Engine,
	handler handler.StaffHandler,
) {
	staffRouter := ginEngine.Group("/staff")

	staffRouter.POST("/create", handler.CreateStaff)
	staffRouter.POST("/login", handler.Login)
}
