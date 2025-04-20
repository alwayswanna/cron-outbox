package configuration

import (
	"github.com/spf13/viper"
	"log"
)

type Properties struct {
	KafkaProperties struct {
		Host string `mapstructure:"host"`
		Port int32  `mapstructure:"port"`
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
}

func NewProperties() *Properties {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("resources")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var configurationProperties Properties
	if err := viper.Unmarshal(&configurationProperties); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &configurationProperties
}
