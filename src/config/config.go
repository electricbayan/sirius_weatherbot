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

type RedisConfig struct {
	RedisHost string
	RedisPort int
	RedisPass string
	RedisDb   string
}

type Config struct {
	Postgres       PostgresConfig
	Redis          RedisConfig
	TgAPIkey       string
	WeatherAPIkey  string
	GeocoderAPIkey string
	Host           string
}

func New() *Config {
	return &Config{
		Postgres: PostgresConfig{
			DbName: getEnv("POSTGRES_DB", ""),
			DbUser: getEnv("POSTGRES_USER", ""),
			DbPass: getEnv("POSTGRES_PASSWORD", ""),
			DbPort: 2828,
			DbHost: "172.28.0.2",
		},
		Redis: RedisConfig{
			RedisHost: "172.28.0.2",
			RedisPort: getEnvAsInt("REDIS_PORT", 6379),
			RedisPass: getEnv("REDIS_PASSWORD", ""),
			RedisDb:   getEnv("REDIS_DB", ""),
		},
		TgAPIkey:       getEnv("TG_API_KEY", ""),
		WeatherAPIkey:  getEnv("WEATHER_API_KEY", ""),
		Host:           getEnv("HOST", ""),
		GeocoderAPIkey: getEnv("GEOCODER_API_KEY", ""),
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
