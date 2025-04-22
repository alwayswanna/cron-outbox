package wire

import (
	"cron-outbox/internal/configuration"
	"cron-outbox/internal/cron"
	"cron-outbox/internal/db"
	"cron-outbox/internal/db/repository"
	"cron-outbox/internal/server"
	"cron-outbox/internal/service"
	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(
	configuration.NewProperties,
	db.NewDatabaseConnection,
	repository.NewArticleRepository,
	repository.NewOutboxMessageRepository,
	service.NewKafkaProducerService,
	service.NewArticleService,
	cron.NewOutBoxMessageScheduler,
	server.NewApp,
)
