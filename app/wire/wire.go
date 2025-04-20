//go:build wireinject
// +build wireinject

package wire

import (
	"cron-outbox/internal/server"
	"github.com/google/wire"
)

func InitializeArticleService() (*server.App, error) {
	wire.Build(ServiceSet)
	return nil, nil
}
