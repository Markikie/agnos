package app

import (
	"github.com/Markikie/agnos/internal/agnos/router"
	"github.com/gin-gonic/gin"
)

func NewRouter(ginEngine *gin.Engine, handler *Handler) {
	router.NewStaffRouter(ginEngine, handler.StaffHandler)
}
