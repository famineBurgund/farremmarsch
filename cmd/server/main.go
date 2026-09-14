package main

import (
	"context"
	"database/sql"
	"famiBurgund/internal/config"
	"flag"
	"log"
	"os"

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

	log := log.New(os.Stdout, "", log.LstdFlags)
	log.Info("Starting server with the following configuration:")
	log.Infof("Postgres: %+v", cfg.Postgres)
	log.Infof("Keycloak: %+v", cfg.Keycloak)
	log.Infof("Log: %+v", cfg.Log)

	// Start the server with the loaded configuration
	// ...

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
	log.Info("Successfully connected to the database")
}
