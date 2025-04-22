package cron

import (
	"cron-outbox/internal/configuration"
	"cron-outbox/internal/service"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type AppCron struct {
	cronManager *cron.Cron
}

func (a *AppCron) StartupScheduledProcesses() {
	a.cronManager.Start()
}

func (a *AppCron) StopScheduledProcesses() {
	a.cronManager.Stop()
}

func NewOutBoxMessageScheduler(properties *configuration.Properties, articleService *service.ArticleServiceImpl) *AppCron {
	cronManager := cron.New()
	_, err := cronManager.AddFunc(properties.Cron.OutboxMessageScheduler, func() {
		go articleService.ProceedSendArticle()
	})

	if err != nil {
		log.Err(err).Msg("failed to add outbox message scheduler job")
		panic(err)
	}

	log.Info().Msg("successfully added outbox message scheduler job")
	return &AppCron{cronManager: cronManager}
}
