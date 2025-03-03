// config/config.go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string

	// Log Settings
	LogLocal int // 0: don't log locally , 1: log locally
	LogDB    int // 0: don't log to db, 1: log to db
	LogFilePath string

	// DB settings
	DBEnabled  int    // 0: deactivate DB, 1: Activate DB
	DBType     string // "postgres", "mysql", "sqlite", "mongodb", "mssql", "oracle"
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "2083"),

		// Log settings
		LogLocal:   getEnvAsInt("LOG_LOCAL", 1),
		LogDB:      getEnvAsInt("LOG_DB", 1),
		LogFilePath: getEnv("LOG_FILE_PATH", "./app.log"),

		// DB settings
		DBEnabled:  getEnvAsInt("DB_ENABLED", 1),
		DBType:     getEnv("DB_TYPE", "sqlite"), // sqllite as default
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "appdb"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

