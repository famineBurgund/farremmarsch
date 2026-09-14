package config

import (
	"fmt"
	"log"
	"os"

	"go.yaml.in/yaml/v2"
)

type Config struct {
	App      App      `yaml: "app"`
	HTTP     HTTP     `yaml: "http"`
	Postgres Postgres `yaml: "postgres"`
	MinIO    MinIO    `yaml: "minio"`
	Keycloak Keycloak `yaml: "keycloak"`
	Log      Log      `yaml: "log"`
}

type App struct {
	Name string `yaml: "name"`
	Env  string `yaml: "env"`
}

type HTTP struct {
	Port         string `yaml: "port"`
	ReadTimeout  int    `yaml: "read_timeout"`
	WriteTimeout int    `yaml: "write_timeout"`
}

type Postgres struct {
	Host     string `yaml: "host"`
	Port     string `yaml: "port"`
	User     string `yaml: "user"`
	Password string `yaml: "password"`
	DBName   string `yaml: "dbname"`
	SSLMode  string `yaml: "sslmode"`
}

type MinIO struct {
	Endpoint        string `yaml: "endpoint"`
	AccessKeyID     string `yaml: "access_key_id"`
	SecretAccessKey string `yaml: "secret_access_key"`
	BucketName      string `yaml: "bucket_name"`
}

type Keycloak struct {
	Enabled       bool   `yaml: "enabled"`
	AuthServerURL string `yaml: "auth_server_url"`
	Realm         string `yaml: "realm"`
	ClientID      string `yaml: "client_id"`
	ClientSecret  string `yaml: "client_secret"`
}

type Log struct {
	Level    string `yaml: "level"`
	FilePath string `yaml: "file_path"`
}

func Load(path string) (*Config, error) {
	cfg, err := loadConfigFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg, nil
}

func loadConfigFromFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
	return &cfg, nil
}

func (p *Postgres) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
	return dsn
}

func PostgresDSN(path string) (string, error) {
	cfg, err := Load(path)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if cfg.Postgres.Password == "" {
		pg, ok := os.LookupEnv("PGPASSWORD")
		if !ok {
			log.Fatalf("Failed to load config: PGPASSWORD environment variable is not set")
		}
		cfg.Postgres.Password = pg
	}
	if cfg.Postgres.Host == "" || cfg.Postgres.Port == "" || cfg.Postgres.User == "" || cfg.Postgres.DBName == "" || cfg.Postgres.SSLMode == "" {
		log.Fatalf("Failed to load config: Postgres configuration is incomplete")
	}
	return cfg.Postgres.GetDSN(), nil
}
