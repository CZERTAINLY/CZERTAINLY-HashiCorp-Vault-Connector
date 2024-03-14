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
		Name     string
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
	config.Database.Host = os.Getenv("DATABASE_HOST")
	config.Database.Port = os.Getenv("DATABASE_PORT")
	config.Database.Username = os.Getenv("DATABASE_USER")
	config.Database.Password = os.Getenv("DATABASE_PASSWORD")
	config.Database.Props = os.Getenv("DATABASE_PROPS")
	config.Database.Name = os.Getenv("DATABASE_NAME")
	config.Database.Schema = os.Getenv("DATABASE_SCHEMA")

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
		l.Info("DATABASE_SCHEMA is mandatory to set!")
	}

	if config.Database.Username == "" {
		l.Info("DATABASE_USER is mandatory to set!")
	}

	if config.Database.Password == "" {
		l.Info("DATABASE_PASSWORD is mandatory to set!")
	}

	if config.Database.Name == "" {
		l.Info("DATABASE_NAME is mandatory to set!")
	}

	return config
}
