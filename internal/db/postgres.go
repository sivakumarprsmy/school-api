package db

import (
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/periasamy/school-api/internal/models"
)

func NewDB(dsn string, logLevel string, log zerolog.Logger) (*gorm.DB, error) {
	gormLogLevel := gormlogger.Silent
	if logLevel == "debug" {
		gormLogLevel = gormlogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&models.Student{}); err != nil {
		return nil, err
	}

	log.Info().Msg("connected to PostgreSQL and ran migrations")
	return db, nil
}
