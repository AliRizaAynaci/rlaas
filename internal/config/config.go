package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DSN  string
	JWT  string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port: env("PORT", "8080"),
		DSN:  buildDSN(),
		JWT:  env("JWT_SECRET", "super-secret-change-me"),
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func buildDSN() string {
	host := env("DB_HOST", "localhost")
	port := env("DB_PORT", "5432")
	user := env("DB_USERNAME", "postgres")
	pass := env("DB_PASSWORD", "password")
	db := env("DB_DATABASE", "rlaas")
	sch := env("DB_SCHEMA", "public")

	return "postgres://" + user + ":" + pass + "@" + host + ":" + port +
		"/" + db + "?sslmode=disable&search_path=" + sch
}
