package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a DB connection and returns *gorm.DB
func Connect(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
}

// Migrate runs AutoMigrate for the given models
func Migrate(db *gorm.DB, models ...interface{}) error {
	// run only if MIGRATE_ON_START=true (env-toggled)
	if os.Getenv("MIGRATE_ON_START") != "true" {
		return nil
	}
	return db.AutoMigrate(models...)
}
