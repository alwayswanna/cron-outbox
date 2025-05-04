package configuration

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Properties struct {
	KafkaProperties struct {
		DeliveryTimeout int `mapstructure:"delivery-timeout"`
		Servers         []struct {
			Host string `mapstructure:"host"`
			Port int32  `mapstructure:"port"`
		} `mapstructure:"servers"`
		Producers []Producers `mapstructure:"producers"`
	} `mapstructure:"kafka"`

	AppProperties struct {
		Port int32 `mapstructure:"port"`
	} `mapstructure:"app"`

	DatabaseProperties struct {
		Host         string `mapstructure:"host"`
		Port         int32  `mapstructure:"port"`
		DatabaseName string `mapstructure:"database-name"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
	} `mapstructure:"database"`

	Cron struct {
		OutboxMessageScheduler string `mapstructure:"outbox-message-scheduler"`
	} `mapstructure:"cron"`

	Profile string `mapstructure:"profile"`
}

type Producers struct {
	TopicName         string `mapstructure:"topic-name"`
	TransactionID     string `mapstructure:"transaction-id"`
	EnableIdempotence bool   `mapstructure:"enable-idempotence"`
	ClientId          string `mapstructure:"client-id"`
	Retries           int    `mapstructure:"retries"`
	Ack               string `mapstructure:"ack"`
}

func (p *Properties) GetKafkaBootstrapServers() string {
	bootstrapServers := ""
	for _, server := range p.KafkaProperties.Servers {
		bootstrapServers += fmt.Sprintf("%s:%d,", server.Host, server.Port)
	}

	return bootstrapServers
}

func (p *Properties) GetGinProfile() string {
	switch p.Profile {
	case "prod":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	case "local":
		return gin.DebugMode
	default:
		return gin.DebugMode
	}
}

func NewProperties() *Properties {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("resources")

	if err := viper.ReadInConfig(); err != nil {
		log.Err(err).Msg("failed to read configuration")
		panic(err)
	}

	var configurationProperties Properties
	if err := viper.Unmarshal(&configurationProperties); err != nil {
		log.Err(err).Msg("failed to unmarshal configuration properties")
		panic(err)
	}

	log.Info().Msg("successfully loaded configuration properties")
	return &configurationProperties
}
