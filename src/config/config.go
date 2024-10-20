package config

import (
	"os"
	"strconv"
)

type PostgresConfig struct {
	DbName string
	DbUser string
	DbPass string
	DbPort int
	DbHost string
}

type Config struct {
	Postgres      PostgresConfig
	TgAPIkey      string
	WeatherAPIkey string
}

func New() *Config {
	return &Config{
		Postgres: PostgresConfig{
			DbName: getEnv("POSTGRES_DB", ""),
			DbUser: getEnv("POSTGRES_USER", ""),
			DbPass: getEnv("POSTGRES_PASSWORD", ""),
			DbPort: getEnvAsInt("POSTGRES_PORT", 8000),
			DbHost: getEnv("POSTGRES_HOST", ""),
		},
		TgAPIkey:      getEnv("TG_API_KEY", ""),
		WeatherAPIkey: getEnv("WEATHER_API_KEY", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
