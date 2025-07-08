package server

import (
	"fmt"
	"github.com/AliRizaAynaci/rlaas/internal/logging"
	"github.com/AliRizaAynaci/rlaas/internal/middleware"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/AliRizaAynaci/rlaas/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	logging.L.Info("Starting RLaaS server", // log at startup
		slog.Int("port", port), // attach port as key–value
	)

	srv := &Server{
		port: port,

		db: database.New(),
	}

	loggedHandler := middleware.LoggingMiddleware(srv.RegisterRoutes())

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      loggedHandler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
