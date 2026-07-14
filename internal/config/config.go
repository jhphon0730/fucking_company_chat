package config

import (
	"os"
	"path"
	"sync"

	"github.com/joho/godotenv"
)

type postgres struct {
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
	SSL_MODE    string
	TIMEZONE    string
}

type Config struct {
	PORT            string
	GIN_MODE        string
	BUILD_VERSION   string
	COMMIT_HASH     string
	DB_AUTO_MIGRATE string

	BCRYPT_COST string
	JWT_SECRET  string

	UPLOAD_DIR string
	BACKUP_DIR string

	POSTGRES postgres
}

var (
	config_instance *Config
	once            sync.Once
)

func LoadConfig() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	_ = godotenv.Load(path.Join(wd, ".env"))

	return &Config{
		PORT:            getEnv("PORT", "5001"),
		GIN_MODE:        getEnv("GIN_MODE", "debug"),
		BUILD_VERSION:   getEnv("BUILD_VERSION", "dev"),
		COMMIT_HASH:     getEnv("COMMIT_HASH", "none"),
		DB_AUTO_MIGRATE: getEnv("DB_AUTO_MIGRATE", "false"),

		BCRYPT_COST: getEnv("BCRYPT_COST", "10"),
		JWT_SECRET:  getEnv("JWT_SECRET", "mysecretkey"),

		UPLOAD_DIR: getEnv("UPLOAD_DIR", "./uploads"),
		BACKUP_DIR: getEnv("BACKUP_DIR", "./backups"),

		POSTGRES: postgres{
			DB_HOST:     getEnv("DB_HOST", "localhost"),
			DB_USER:     getEnv("DB_USER", "postgres"),
			DB_PASSWORD: getEnv("DB_PASSWORD", "postgres"),
			DB_NAME:     getEnv("DB_NAME", "postgres"),
			DB_PORT:     getEnv("DB_PORT", "5432"),
			SSL_MODE:    getEnv("SSL_MODE", "disable"),
			TIMEZONE:    getEnv("TIMEZONE", "Asia/Seoul"),
		},
	}, nil
}

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		config_instance, err = LoadConfig()
	})

	return config_instance, err
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
