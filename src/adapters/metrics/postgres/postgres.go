package postgres

import (
	"fmt"
	"github.com/yeencloud/ServiceCore/src/config"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Database struct {
	engine      *gorm.DB
	serviceName string
}

func StartGormDatabase(config *config.Database, serviceName string) *Database {
	if config == nil {
		log.Fatalln("Database config is invalid")
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := gorm.Open(pg.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err)
	}
	db.Logger.LogMode(logger.Info)

	return &Database{
		engine:      db,
		serviceName: serviceName,
	}
}
