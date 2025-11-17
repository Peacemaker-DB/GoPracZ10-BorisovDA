package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Port      string
	JWTSecret []byte
	JWTTTL    time.Duration

	RefreshSecret []byte
	RefreshTTL    time.Duration
}

func Load() Config {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	ttl := os.Getenv("JWT_TTL")
	if ttl == "" {
		ttl = "24h"
	}
	dur, err := time.ParseDuration(ttl)
	if err != nil {
		log.Fatal("bad JWT_TTL")
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "refresh-dev" // можно по-учебному
	}

	refreshTTL := os.Getenv("REFRESH_TTL")
	if refreshTTL == "" {
		refreshTTL = "168h" // 7 days
	}

	refreshDur, err := time.ParseDuration(refreshTTL)
	if err != nil {
		log.Fatal("bad REFRESH_TTL")
	}

	return Config{
		Port:          ":" + port,
		JWTSecret:     []byte(secret),
		JWTTTL:        dur,
		RefreshSecret: []byte(refreshSecret),
		RefreshTTL:    refreshDur,
	}
}
