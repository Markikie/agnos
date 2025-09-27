package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World (HTTPS)")
	})

	// ฟังที่ port 443 ใช้ cert/key ที่สร้างไว้
	err := app.RunTLS(":443",
		"./ssl/hospital-a.api.co.th.crt",
		"./ssl/hospital-a.api.co.th.key")

	if err != nil {
		panic(err)
	}
}
