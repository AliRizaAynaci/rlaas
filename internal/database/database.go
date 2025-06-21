package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

type Service interface {
	Health() map[string]string
	Close() error
	Pool() *pgxpool.Pool
}

type service struct {
	pool *pgxpool.Pool
}

// New creates or returns a singleton database service
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	dbURL := buildDSN()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		return nil
	}

	dbInstance = &service{
		pool: pool,
	}

	fmt.Println("Database connection established")
	fmt.Printf("Connected to %s\n", dbURL)
	return dbInstance
}

func GetInstance() (Service, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}

	dbService := New()

	dbInstance = &service{pool: dbService.Pool()}
	return dbInstance, nil
}

func (s *service) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *service) Close() error {
	if s.pool != nil {
		s.pool.Close()
		log.Println("Disconnected from database")
	}
	return nil
}

func (s *service) Health() map[string]string {
	stats := s.pool.Stat()
	return map[string]string{
		"status":              "up",
		"total_connections":   strconv.FormatInt(int64(stats.TotalConns()), 10),
		"idle_connections":    strconv.FormatInt(int64(stats.IdleConns()), 10),
		"used_connections":    strconv.FormatInt(int64(stats.AcquiredConns()), 10),
		"max_connections":     strconv.FormatInt(int64(stats.MaxConns()), 10),
		"acquire_count":       strconv.FormatInt(stats.AcquireCount(), 10),
		"acquire_duration":    stats.AcquireDuration().String(),
		"canceled_acquire":    strconv.FormatInt(stats.CanceledAcquireCount(), 10),
		"empty_acquire_count": strconv.FormatInt(stats.EmptyAcquireCount(), 10),
	}
}

func buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		username,
		password,
		host,
		port,
		database,
		schema,
	)
}
