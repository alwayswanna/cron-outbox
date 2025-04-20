package db

import (
	"cron-outbox/internal/configuration"
	"fmt"
	"gorm.io/gorm"
	"log"
)
import "gorm.io/driver/postgres"

type DatabaseConnection struct {
	db *gorm.DB
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
		log.Fatal(err)
	}

	return &DatabaseConnection{db: db}
}
