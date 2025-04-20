package main

import (
	"cron-outbox/wire"
	"log"
)

func main() {
	app, err := wire.InitializeArticleService()
	if err != nil {
		log.Fatal("Error on startup service", err)
	}

	app.Start()
}
