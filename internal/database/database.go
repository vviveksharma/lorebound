package database

import (
	"fmt"

	"github.com/vviveksharma/lorebound/internal/logger"
	"github.com/vviveksharma/lorebound/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func dsn() string {
	sslMode := utils.Getenv("DB_SSLMODE", "prefer")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		utils.MustGetenv("DB_HOST"),
		utils.MustGetenv("DB_USER"),
		utils.MustGetenv("DB_PASSWORD"),
		utils.MustGetenv("DB_NAME"),
		utils.Getenv("DB_PORT", "5432"),
		sslMode,
	)
}

func Connect() *gorm.DB {
	var db *gorm.DB
	var err error

	db, err = gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		logger.Error("database connection failed", "err", err)
		return nil
	}
	logger.Info("database connected")
	return db
}
