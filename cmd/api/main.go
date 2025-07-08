package main

import (
	"context"
	"fmt"
	"github.com/AliRizaAynaci/rlaas/internal/limiter"
	"github.com/AliRizaAynaci/rlaas/internal/logging"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/AliRizaAynaci/rlaas/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logging.L.Info("Shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logging.L.Error("Forced shutdown error",
			slog.String("error", err.Error()),
		)
	}

	logging.L.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	limiter.InitSharding()
	logging.L.Info("Redis sharding initialized")
	server := server.NewServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	logging.L.Info("Graceful shutdown complete")
}
