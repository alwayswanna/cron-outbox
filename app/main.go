package main

import (
	"cron-outbox/wire"
	"github.com/rs/zerolog/log"
)

func main() {
	app, err := wire.InitializeArticleService()
	if err != nil {
		log.Err(err).Msg("failed to initialize article service")
	}

	app.Start()
	log.Info().Msg("article service successfully started")
}
