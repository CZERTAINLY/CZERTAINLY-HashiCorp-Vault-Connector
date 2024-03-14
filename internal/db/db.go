package db

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/config"
	"CZERTAINLY-HashiCorp-Vault-Connector/internal/logger"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var log = logger.Get()

func ConnectDB(config config.Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s %s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name, config.Database.Props)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.Database.Schema + ".",
			SingularTable: false,
		}})
	return
}

func MigrateDB(config config.Config) {
	log.Info("Migrating database")
	// search_path=public&x-migrations-table=hvault_migrations migration table name and schema, migration table must be in public schema if we want to create schema automatically
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?search_path=%s&x-migrations-table=hvault_migrations", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name, config.Database.Schema)
	m, err := migrate.New(
		"file://migrations",
		connectionString,
	)
	if err != nil {
		log.Error(err.Error())
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			{
				log.Error(err.Error())
			}
		}
	}
}
