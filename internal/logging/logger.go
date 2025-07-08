package logging

import (
	"log/slog"
	"os"
)

// determine environment with a fallback
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// baseHandler writes JSON-formatted logs to stdout, including source location.
var baseHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
})

// L is the global logger preconfigured with default fields.
var L = slog.New(baseHandler).With(
	slog.String("service", "rlaas"),          // service name
	slog.String("env", getEnv("ENV", "dev")), // default to "dev"
)
