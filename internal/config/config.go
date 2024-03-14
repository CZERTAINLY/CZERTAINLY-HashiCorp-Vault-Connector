package config

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"os"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Name	 string
		Host     string
		Port     string
		Username string
		Password string
		Props    string
		Schema   string
	}
}

var config Config

func Get() Config {
	l := logger.Get()

	config.Server.Port = os.Getenv("SERVER_PORT")
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Username = os.Getenv("DB_USERNAME")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Props = os.Getenv("DB_PROPS")
	config.Database.Name = os.Getenv("DB_DB")
	config.Database.Schema = os.Getenv("DB_SCHEMA")

	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}

	if config.Database.Port == "" {
		config.Database.Port = "5432"
	}

	if config.Database.Schema == "" {
		l.Info("DB_SCHEMA is mandatory to set!")
	}

	if config.Database.Username == "" {
		l.Info("DB_USERNAME is mandatory to set!")
	}

	if config.Database.Password == "" {
		l.Info("DB_PASSWORD is mandatory to set!")
	}

	return config
}
