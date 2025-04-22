package server

import (
	"cron-outbox/internal/configuration"
	"cron-outbox/internal/cron"
	"cron-outbox/internal/model"
	"cron-outbox/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type App struct {
	server         *gin.Engine
	properties     *configuration.Properties
	appCron        *cron.AppCron
	articleService *service.ArticleServiceImpl
}

func (app *App) Start() {
	app.appCron.StartupScheduledProcesses()
	err := app.server.Run(fmt.Sprintf(":%d", app.properties.AppProperties.Port))
	if err != nil {
		log.Err(err).Msg("failed to start server")
		panic(err)
	}

	defer app.appCron.StopScheduledProcesses()
}

func NewApp(properties *configuration.Properties, articleService *service.ArticleServiceImpl, appCron *cron.AppCron) *App {
	gin.SetMode(properties.GetGinProfile())
	router := gin.Default()

	router.POST("/api/v1/article", func(c *gin.Context) {
		var request model.ArticleRequest
		err := c.ShouldBindBodyWithJSON(&request)
		if err != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprintf("Error binding request: %s", err.Error()),
			})
			return
		}
		err = articleService.ProceedSaveArticle(request)
		if err != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprintf("Error saving article: %s", err.Error()),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "Success",
		})

	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	log.Info().Msgf("Starting server on port %d", properties.AppProperties.Port)
	return &App{server: router, properties: properties, articleService: articleService, appCron: appCron}
}
