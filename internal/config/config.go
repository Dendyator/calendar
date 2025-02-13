package config

import (
	"log"
	"time"

	"github.com/spf13/viper" //nolint:depguard
)

type Config struct {
	Server    ServerConfig
	GRPC      GRPCConfig
	Logger    LoggerConfig
	Database  DatabaseConfig
	RabbitMQ  RabbitMQConfig
	Scheduler SchedulerConfig
	Sender    SenderConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type GRPCConfig struct {
	Host string
	Port string
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type RabbitMQConfig struct {
	DSN string
}

type SchedulerConfig struct {
	Interval time.Duration
}

type SenderConfig struct{}

func LoadConfig(configPath string) Config {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	if intervalStr := viper.GetString("scheduler.interval"); intervalStr != "" {
		config.Scheduler.Interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			log.Fatalf("Invalid scheduler interval format, %v", err)
		}
	}

	return config
}
