package db

import (
	"cron-outbox/internal/configuration"
	"fmt"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)
import "gorm.io/driver/postgres"

type DatabaseConnection struct {
	db *gorm.DB
}

func (db *DatabaseConnection) GetDB() *gorm.DB {
	return db.db
}

func NewDatabaseConnection(properties *configuration.Properties) *DatabaseConnection {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		properties.DatabaseProperties.Host,
		properties.DatabaseProperties.Username,
		properties.DatabaseProperties.Password,
		properties.DatabaseProperties.DatabaseName,
		properties.DatabaseProperties.Port,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Err(err).Msg("failed to connect to database")
		panic(err)
	}

	log.Info().Msg("successfully connected to database")
	return &DatabaseConnection{db: db}
}
