package database

import (
	"fmt"
	"log"
	"ponial/internal/models"
	"ponial/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	var err error
	var gormLogger logger.Interface

	if cfg.AppMode == "debug" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database object: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(100000)

	log.Println("Database connected successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	return DB.AutoMigrate(
		&models.User{},
		&models.YandexAI{},
		&models.Workspace{},
	)
}

func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
