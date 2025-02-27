package config

import "os"

var ServerPort = getEnv("SERVER_PORT", "2083")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

