package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	AppPort  string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	NBABaseURL string
	NBATimeout int
	CacheTTL   int

	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	redisDB, _    := strconv.Atoi(getEnv("REDIS_DB", "0"))
	nbaTimeout, _ := strconv.Atoi(getEnv("NBA_TIMEOUT", "15"))
	cacheTTL, _   := strconv.Atoi(getEnv("CACHE_TTL_SECONDS", "3600"))

	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "airball"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
		NBABaseURL:    getEnv("NBA_BASE_URL", "https://stats.nba.com/stats"),
		NBATimeout:    nbaTimeout,
		CacheTTL:      cacheTTL,
		JWTSecret:     getEnv("JWT_SECRET", "change-me"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
