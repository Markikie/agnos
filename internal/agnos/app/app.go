package app

import (
	"log"

	"github.com/Markikie/agnos/internal/agnos"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := env.Parse(&agnos.Env); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	Config *Config
}

func NewApp() *App {
	config := NewConfig()
	ginEngine := gin.New()

	NewMiddleware(ginEngine)
	repository := NewRepository(config)
	service := NewService(repository)
	handler := NewHandler(service)
	NewRouter(ginEngine, handler)
	return &App{
		Config: config,
	}
}
