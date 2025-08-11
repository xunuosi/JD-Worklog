package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPass      string
	DBName      string
	JWTSecret   string
	CORSOrigins []string
}

func Load() *Config {
	c := &Config{
		DBHost:      getEnv("DB_HOST", "127.0.0.1"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "worklog"),
		DBPass:      getEnv("DB_PASS", "worklog@123"),
		DBName:      getEnv("DB_NAME", "worklog"),
		JWTSecret:   getEnv("JWT_SECRET", "devsecretchangeit"),
		CORSOrigins: splitCsv(getEnv("CORS_ORIGINS", "http://localhost:5173")),
	}
	log.Printf("config loaded: %+v", *c)
	return c
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func splitCsv(s string) []string {
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return []string{s}
	}
	return parts
}
