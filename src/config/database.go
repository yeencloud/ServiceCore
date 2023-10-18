package config

import (
	types2 "github.com/yeencloud/ServiceCore/src/domain/types"
)

type Database struct {
	Host     types2.Host
	Port     types2.Port
	User     string
	Password string
	Database string
}

func (config *Config) getDatabase() {
	if config.Database != nil {
		return
	}

	DatabaseURL := config.GetEnvStringOrDefault("DATABASE_URL", "")
	if DatabaseURL == "" {
		return
	}

	database := &Database{
		Host:     types2.Host(DatabaseURL),
		Port:     types2.Port(config.RequireEnvInt("DATABASE_PORT")),
		User:     config.RequireEnvString("DATABASE_USER"),
		Password: config.RequireEnvString("DATABASE_PASSWORD"),
		Database: config.RequireEnvString("DATABASE_NAME"),
	}

	config.Database = database
}

func (config *Config) getMetricsDatabase() {

}
