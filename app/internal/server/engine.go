package server

import (
	"cron-outbox/internal/configuration"
	"fmt"
	"github.com/gin-gonic/gin"
)

type App struct {
	server     *gin.Engine
	properties *configuration.Properties
}

func (app *App) Start() {
	err := app.server.Run(fmt.Sprintf(":%d", app.properties.AppProperties.Port))
	if err != nil {
		panic(err)
	}
}

func NewApp(properties *configuration.Properties) *App {
	router := gin.Default()

	router.POST("/api/v1/article", func(c *gin.Context) {
		// TODO: implement.
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	return &App{server: router, properties: properties}
}
