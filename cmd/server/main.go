package main

import (
	"context"
	"database/sql"
	"famiBurgund/internal/config"
	"famiBurgund/pkg/logger"
	"flag"
	"log"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.New(logger.Config{
		Level:      cfg.Log.Level,
		FilePath:   cfg.Log.FilePath,
		MaxSizeMB:  cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAgeDays: cfg.Log.MaxAgeDays,
	})
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Starting server with the following configuration:", zap.String("app_name", cfg.App.Name), zap.String("env", cfg.App.Env))

	var (
		db    *gorm.DB
		sqlDB *sql.DB
	)
	db, err = gorm.Open(postgres.Open(cfg.Postgres.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}
	defer sqlDB.Close()

	if err != sqlDB.PingContext(context.Background()) {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	appLogger.Info("Successfully connected to the database")
}
