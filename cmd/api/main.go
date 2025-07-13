package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/AliRizaAynaci/rlaas/internal/app"
	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/logging"
)

func gracefulShutdown(app *fiber.App, done chan<- bool) {
	// Listen for interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done() // block until signal
	logging.L.Info("Shutting down gracefully, press Ctrl+C again to force")

	// Allow second Ctrl+C to force exit
	stop()

	// Give active connections 5 s to finish
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(timeoutCtx); err != nil {
		logging.L.Error("Forced shutdown error", "err", err)
	}

	logging.L.Info("Server exiting")
	done <- true // notify main() that shutdown is complete
}

func main() {
	/* ------------ infra ------------ */
	limiter.InitSharding()
	logging.L.Info("Redis sharding initialized")

	/* ------------ build Fiber app ------------ */
	app := app.New() // all wiring (DB, routes, etc.) inside

	/* ------------ graceful shutdown ------------ */
	done := make(chan bool, 1)
	go gracefulShutdown(app, done)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local fallback
	}
	logging.L.Info("Starting server on port", "port", port)

	if err := app.Listen(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logging.L.Error("http server error", "err", err)
	}

	<-done
	logging.L.Info("Graceful shutdown complete")
}
